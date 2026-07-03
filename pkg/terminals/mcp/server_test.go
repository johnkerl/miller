package mcp

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// connectTestClient stands up the server and a client over in-memory
// transports and returns the connected client session.
func connectTestClient(t *testing.T, config *serverConfig) *mcpsdk.ClientSession {
	t.Helper()
	if config == nil {
		config = &serverConfig{
			timeoutSeconds: defaultTimeoutSeconds,
			maxOutputBytes: defaultMaxOutputBytes,
		}
	}
	server := newServer(config)

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	serverSession, err := server.Connect(context.Background(), serverTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = serverSession.Wait() })

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test-client", Version: "0"}, nil)
	clientSession, err := client.Connect(context.Background(), clientTransport, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = clientSession.Close() })

	return clientSession
}

// callToolInto invokes a tool and unmarshals its structured content into out.
func callToolInto(t *testing.T, session *mcpsdk.ClientSession, name string, arguments map[string]any, out any) *mcpsdk.CallToolResult {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcpsdk.CallToolParams{
		Name:      name,
		Arguments: arguments,
	})
	require.NoError(t, err)
	require.False(t, result.IsError, "tool %s returned IsError; content: %+v", name, result.Content)
	b, err := json.Marshal(result.StructuredContent)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(b, out))
	return result
}

func TestListToolsAndAnnotations(t *testing.T) {
	session := connectTestClient(t, nil)

	result, err := session.ListTools(context.Background(), &mcpsdk.ListToolsParams{})
	require.NoError(t, err)

	byName := map[string]*mcpsdk.Tool{}
	for _, tool := range result.Tools {
		byName[tool.Name] = tool
	}
	for _, name := range []string{"list_capabilities", "which", "validate_dsl", "describe_data", "run"} {
		require.Contains(t, byName, name)
	}

	assert.True(t, byName["list_capabilities"].Annotations.ReadOnlyHint)
	assert.True(t, byName["which"].Annotations.ReadOnlyHint)
	require.NotNil(t, byName["run"].Annotations.DestructiveHint)
	assert.True(t, *byName["run"].Annotations.DestructiveHint)
}

func TestListCapabilitiesIndex(t *testing.T) {
	session := connectTestClient(t, nil)

	var output listCapabilitiesOutput
	callToolInto(t, session, "list_capabilities", map[string]any{"index": true}, &output)

	assert.NotEmpty(t, output.MlrVersion)
	assert.Greater(t, output.CatalogSchemaVersion, 0)
	assert.NotEmpty(t, output.Index)
	assert.Empty(t, output.Verbs)

	kinds := map[string]bool{}
	for _, entry := range output.Index {
		kinds[entry.Kind] = true
	}
	for _, kind := range []string{"verb", "function", "flag", "keyword"} {
		assert.True(t, kinds[kind], "index missing kind %s", kind)
	}
}

func TestListCapabilitiesKindAndNames(t *testing.T) {
	session := connectTestClient(t, nil)

	var output listCapabilitiesOutput
	callToolInto(t, session, "list_capabilities",
		map[string]any{"kind": "verb", "names": []string{"cat"}}, &output)
	require.Len(t, output.Verbs, 1)
	assert.Equal(t, "cat", output.Verbs[0].Name)
	assert.NotEmpty(t, output.Verbs[0].UsageText)

	var full listCapabilitiesOutput
	callToolInto(t, session, "list_capabilities", map[string]any{}, &full)
	assert.NotEmpty(t, full.Verbs)
	assert.NotEmpty(t, full.Functions)
	assert.NotEmpty(t, full.Flags)
	assert.NotEmpty(t, full.Keywords)
}

func TestListCapabilitiesBadKind(t *testing.T) {
	session := connectTestClient(t, nil)

	result, err := session.CallTool(context.Background(), &mcpsdk.CallToolParams{
		Name:      "list_capabilities",
		Arguments: map[string]any{"kind": "nonesuch"},
	})
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestWhich(t *testing.T) {
	session := connectTestClient(t, nil)

	var output whichOutput
	callToolInto(t, session, "which", map[string]any{"query": "join two files on a key"}, &output)
	require.NotEmpty(t, output.Results)
	assert.True(t, output.Confident)
	assert.Equal(t, "join", output.Results[0].Name)
}

func TestPlaybookPromptAndResource(t *testing.T) {
	session := connectTestClient(t, nil)

	promptResult, err := session.GetPrompt(context.Background(), &mcpsdk.GetPromptParams{
		Name: playbookPromptName,
	})
	require.NoError(t, err)
	require.Len(t, promptResult.Messages, 1)
	text := promptResult.Messages[0].Content.(*mcpsdk.TextContent).Text
	assert.Contains(t, text, "Miller agent playbook")

	resourceResult, err := session.ReadResource(context.Background(), &mcpsdk.ReadResourceParams{
		URI: playbookResourceURI,
	})
	require.NoError(t, err)
	require.Len(t, resourceResult.Contents, 1)
	assert.Equal(t, text, resourceResult.Contents[0].Text)
}

// ----------------------------------------------------------------
// Subprocess-backed tools. These need a built mlr binary; they are skipped
// when the repo-root binary is absent (e.g. `go test` before `make build`).

// repoRootMlrPath locates the mlr binary built at the repository root,
// relative to this source file.
func repoRootMlrPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	exeName := "mlr"
	if runtime.GOOS == "windows" {
		exeName = "mlr.exe"
	}
	path := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", exeName)
	if _, err := os.Stat(path); err != nil {
		t.Skipf("no built mlr binary at %s; run make build first", path)
	}
	return path
}

func subprocessTestConfig(t *testing.T) *serverConfig {
	return &serverConfig{
		mlrExePath:     repoRootMlrPath(t),
		timeoutSeconds: defaultTimeoutSeconds,
		maxOutputBytes: defaultMaxOutputBytes,
	}
}

func TestValidateDSL(t *testing.T) {
	session := connectTestClient(t, subprocessTestConfig(t))

	var good validateDSLOutput
	callToolInto(t, session, "validate_dsl", map[string]any{"expression": "$z = $x + $y"}, &good)
	assert.True(t, good.Valid)
	assert.Nil(t, good.Error)

	var bad validateDSLOutput
	callToolInto(t, session, "validate_dsl", map[string]any{"expression": "$z = $x +"}, &bad)
	assert.False(t, bad.Valid)
	require.NotNil(t, bad.Error)
	assert.Equal(t, "dsl-parse-error", bad.Error.Kind)

	result, err := session.CallTool(context.Background(), &mcpsdk.CallToolParams{
		Name:      "validate_dsl",
		Arguments: map[string]any{"expression": "$a=1", "kind": "nonesuch"},
	})
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestDescribeData(t *testing.T) {
	session := connectTestClient(t, subprocessTestConfig(t))

	csvPath := filepath.Join(t.TempDir(), "data.csv")
	require.NoError(t, os.WriteFile(csvPath, []byte("a,b\n1,x\n2,y\n2,x\n"), 0o644))

	var output describeDataOutput
	callToolInto(t, session, "describe_data",
		map[string]any{"files": []string{csvPath}, "input_format": "csv"}, &output)
	require.Len(t, output.Fields, 2)
	assert.Equal(t, "a", output.Fields[0]["field_name"])
	assert.Equal(t, float64(2), output.Fields[0]["distinct_count"])
}

func TestRun(t *testing.T) {
	session := connectTestClient(t, subprocessTestConfig(t))

	var ok runOutput
	callToolInto(t, session, "run", map[string]any{
		"args":       []string{"--icsv", "--ojson", "cat"},
		"stdin_text": "a,b\n1,2\n",
	}, &ok)
	assert.Equal(t, 0, ok.ExitCode)
	assert.Contains(t, ok.Stdout, `"a": 1`)
	assert.False(t, ok.StdoutTruncated)

	// Unknown verb: structured error parsed from MLR_ERRORS_JSON stderr.
	var bad runOutput
	callToolInto(t, session, "run", map[string]any{"args": []string{"frobnicate"}}, &bad)
	assert.NotEqual(t, 0, bad.ExitCode)
	require.NotNil(t, bad.Error)
	assert.Equal(t, "unknown-verb", bad.Error.Kind)
	assert.Equal(t, "frobnicate", bad.Error.Token)

	// system() is blocked by the injected MLR_NO_SHELL.
	var blocked runOutput
	callToolInto(t, session, "run", map[string]any{
		"args": []string{"-n", "put", `end{print system("echo owned")}`},
	}, &blocked)
	assert.Equal(t, 0, blocked.ExitCode)
	assert.Contains(t, blocked.Stdout, "(error)")
	assert.NotContains(t, blocked.Stdout, "owned")
}

func TestRunAllowShell(t *testing.T) {
	config := subprocessTestConfig(t)
	config.allowShell = true
	session := connectTestClient(t, config)

	var output runOutput
	callToolInto(t, session, "run", map[string]any{
		"args": []string{"-n", "put", `end{print system("echo owned")}`},
	}, &output)
	assert.Equal(t, 0, output.ExitCode)
	assert.Contains(t, output.Stdout, "owned")
}

func TestRunOutputCap(t *testing.T) {
	config := subprocessTestConfig(t)
	config.maxOutputBytes = 64
	session := connectTestClient(t, config)

	var output runOutput
	callToolInto(t, session, "run", map[string]any{
		"args": []string{"--ojson", "seqgen", "--start", "1", "--stop", "100"},
	}, &output)
	assert.Equal(t, 0, output.ExitCode)
	assert.True(t, output.StdoutTruncated)
	assert.LessOrEqual(t, len(output.Stdout), 64)
}

func TestRunTimeout(t *testing.T) {
	config := subprocessTestConfig(t)
	session := connectTestClient(t, config)

	var output runOutput
	callToolInto(t, session, "run", map[string]any{
		"args":            []string{"--ojson", "seqgen", "--start", "1", "--stop", "9000000000"},
		"timeout_seconds": 1,
	}, &output)
	assert.True(t, output.TimedOut)
	assert.NotEqual(t, 0, output.ExitCode)
}

// ----------------------------------------------------------------
// Pure helpers

func TestCapWriter(t *testing.T) {
	w := &capWriter{limit: 5}
	n, err := w.Write([]byte("abc"))
	require.NoError(t, err)
	assert.Equal(t, 3, n)
	n, err = w.Write([]byte("defgh"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "abcde", w.buf.String())
	assert.True(t, w.truncated)
}

func TestParseStructuredError(t *testing.T) {
	se := parseStructuredError("some preamble\n{\n  \"error\": \"boom\",\n  \"kind\": \"generic\"\n}\n")
	require.NotNil(t, se)
	assert.Equal(t, "boom", se.Error)
	assert.Equal(t, "generic", se.Kind)

	assert.Nil(t, parseStructuredError("plain prose error\n"))
	assert.Nil(t, parseStructuredError(""))
	assert.Nil(t, parseStructuredError("{\"unrelated\": true}"))
}

func TestPlaybookHasFrontmatter(t *testing.T) {
	assert.True(t, strings.HasPrefix(playbookText, "---\n"))
	assert.Contains(t, playbookText, "name: miller")
}
