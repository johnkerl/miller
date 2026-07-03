// Subprocess-backed MCP tools: validate_dsl, describe_data, and run. These
// shell out to the same mlr binary that is serving MCP (os.Executable()), so
// an agent sees byte-identical behavior to a terminal, and the CLI paths'
// os.Exit/global-state/panic behaviors stay isolated from the server process.
//
// Every subprocess runs with MLR_ERRORS_JSON=1 (parse errors come back as a
// structured JSON document on stderr) and, unless the server was started with
// --allow-shell, MLR_NO_SHELL=1 (the DSL system/exec functions, piped
// redirects, and --prepipe fail cleanly). The MLR_NO_SHELL gate is one-way in
// the child: an agent-supplied argv cannot re-enable it.

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// structuredError mirrors the JSON document `mlr --errors-json` emits on
// stderr (climain.StructuredError). It is redeclared here rather than
// imported: pkg/climain dispatches into pkg/terminals, which imports this
// package, so importing climain from here would be an import cycle -- and
// this is a wire shape parsed from a subprocess, not a shared internal type.
type structuredError struct {
	Error      string   `json:"error"`
	Kind       string   `json:"kind"`
	Token      string   `json:"token,omitempty"`
	Verb       string   `json:"verb,omitempty"`
	Hint       string   `json:"hint,omitempty"`
	DidYouMean []string `json:"did_you_mean,omitempty"`
}

// capWriter accumulates writes up to a byte limit, discarding (but counting)
// the excess. It never reports an error, so the subprocess is not killed by
// its own verbosity -- output beyond the cap is simply dropped and the result
// marked truncated.
type capWriter struct {
	buf       strings.Builder
	limit     int64
	truncated bool
}

func (w *capWriter) Write(p []byte) (int, error) {
	room := w.limit - int64(w.buf.Len())
	if room <= 0 {
		w.truncated = true
		return len(p), nil
	}
	if int64(len(p)) > room {
		w.truncated = true
		w.buf.Write(p[:room])
	} else {
		w.buf.Write(p)
	}
	return len(p), nil
}

// mlrResult is what runMlrSubprocess reports back to the tool handlers.
type mlrResult struct {
	exitCode        int
	stdout          string
	stdoutTruncated bool
	stderr          string
	stderrTruncated bool
	timedOut        bool
	structured      *structuredError
}

// runMlrSubprocess executes this same mlr binary with the given argv (not
// including the executable name), an optional stdin payload, and a wall-clock
// limit. Stdin is always redirected -- never inherited -- since the server's
// own stdin carries the MCP transport.
func runMlrSubprocess(
	ctx context.Context,
	config *serverConfig,
	argv []string,
	stdinText string,
	timeoutSeconds int64,
) (*mlrResult, error) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = config.timeoutSeconds
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, config.mlrExePath, argv...)
	cmd.Stdin = strings.NewReader(stdinText)

	stdout := &capWriter{limit: config.maxOutputBytes}
	stderr := &capWriter{limit: config.maxOutputBytes}
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	env := append(os.Environ(), "MLR_ERRORS_JSON=1")
	if !config.allowShell {
		env = append(env, "MLR_NO_SHELL=1")
	}
	cmd.Env = env

	err := cmd.Run()
	result := &mlrResult{
		stdout:          stdout.buf.String(),
		stdoutTruncated: stdout.truncated,
		stderr:          stderr.buf.String(),
		stderrTruncated: stderr.truncated,
		timedOut:        ctx.Err() == context.DeadlineExceeded,
	}

	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok && !result.timedOut {
			// The command could not be started at all.
			return nil, fmt.Errorf("could not run %s: %w", config.mlrExePath, err)
		}
		if ok {
			result.exitCode = exitErr.ExitCode()
		} else {
			result.exitCode = -1
		}
		result.structured = parseStructuredError(result.stderr)
	}

	return result, nil
}

// parseStructuredError extracts the MLR_ERRORS_JSON document from a failed
// command's stderr, if one is present. The document is the final thing
// written to stderr, so scan forward for a line starting the JSON object and
// try each candidate.
func parseStructuredError(stderrText string) *structuredError {
	lines := strings.Split(stderrText, "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			candidate := strings.Join(lines[i:], "\n")
			var se structuredError
			if err := json.Unmarshal([]byte(candidate), &se); err == nil && se.Error != "" {
				return &se
			}
		}
	}
	return nil
}

// ----------------------------------------------------------------
// validate_dsl

type validateDSLInput struct {
	Expression string `json:"expression" jsonschema:"The DSL expression to check e.g. '$z = $x + $y'."`
	Kind       string `json:"kind,omitempty" jsonschema:"put (default) or filter."`
}

type validateDSLOutput struct {
	Valid  bool             `json:"valid"`
	Error  *structuredError `json:"error,omitempty"`
	Stderr string           `json:"stderr,omitempty"`
}

func validateDSLHandler(config *serverConfig) mcpsdk.ToolHandlerFor[validateDSLInput, validateDSLOutput] {
	return func(
		ctx context.Context,
		_ *mcpsdk.CallToolRequest,
		input validateDSLInput,
	) (*mcpsdk.CallToolResult, validateDSLOutput, error) {
		kind := input.Kind
		if kind == "" {
			kind = "put"
		}
		if kind != "put" && kind != "filter" {
			return nil, validateDSLOutput{}, fmt.Errorf("kind must be put or filter, not %q", kind)
		}
		if input.Expression == "" {
			return nil, validateDSLOutput{}, fmt.Errorf("expression must be non-empty")
		}

		result, err := runMlrSubprocess(ctx, config, []string{kind, "--explain", input.Expression}, "", 0)
		if err != nil {
			return nil, validateDSLOutput{}, err
		}

		output := validateDSLOutput{Valid: result.exitCode == 0}
		if !output.Valid {
			output.Error = result.structured
			if output.Error == nil {
				output.Stderr = result.stderr
			}
		}
		return nil, output, nil
	}
}

// ----------------------------------------------------------------
// describe_data

type describeDataInput struct {
	Files       []string `json:"files" jsonschema:"Input file paths."`
	InputFormat string   `json:"input_format,omitempty" jsonschema:"csv / tsv / json / jsonl / dkvp / nidx / pprint / xtab etc.; Miller's default is dkvp."`
	MaxValues   int64    `json:"max_values,omitempty" jsonschema:"List a field's distinct values only if it has at most this many; 0 uses the describe verb's default of 20."`
	ExtraArgs   []string `json:"extra_args,omitempty" jsonschema:"Additional main-level flags placed before the verb e.g. [\"--ifs\" \";\"] or [\"--implicit-csv-header\"]."`
}

type describeDataOutput struct {
	Fields []map[string]any `json:"fields"`
	Error  *structuredError `json:"error,omitempty"`
	Stderr string           `json:"stderr,omitempty"`
}

func describeDataHandler(config *serverConfig) mcpsdk.ToolHandlerFor[describeDataInput, describeDataOutput] {
	return func(
		ctx context.Context,
		_ *mcpsdk.CallToolRequest,
		input describeDataInput,
	) (*mcpsdk.CallToolResult, describeDataOutput, error) {
		if len(input.Files) == 0 {
			return nil, describeDataOutput{}, fmt.Errorf("files must be non-empty")
		}

		argv := []string{}
		if input.InputFormat != "" {
			argv = append(argv, "-i", input.InputFormat)
		}
		argv = append(argv, input.ExtraArgs...)
		argv = append(argv, "--ojson", "describe")
		if input.MaxValues > 0 {
			argv = append(argv, "-n", strconv.FormatInt(input.MaxValues, 10))
		}
		argv = append(argv, input.Files...)

		result, err := runMlrSubprocess(ctx, config, argv, "", 0)
		if err != nil {
			return nil, describeDataOutput{}, err
		}
		if result.exitCode != 0 {
			return nil, describeDataOutput{Error: result.structured, Stderr: result.stderr},
				fmt.Errorf("mlr describe exited %d", result.exitCode)
		}

		var fields []map[string]any
		if err := json.Unmarshal([]byte(result.stdout), &fields); err != nil {
			return nil, describeDataOutput{}, fmt.Errorf("could not parse mlr describe output: %w", err)
		}
		return nil, describeDataOutput{Fields: fields}, nil
	}
}

// ----------------------------------------------------------------
// run

type runInput struct {
	Args           []string `json:"args" jsonschema:"The mlr argv without the leading executable name; one element per shell word e.g. [\"--icsv\" \"--ojson\" \"cat\" \"data.csv\"]."`
	StdinText      string   `json:"stdin_text,omitempty" jsonschema:"Text supplied to the command's stdin; commands never read the server's own stdin."`
	TimeoutSeconds int64    `json:"timeout_seconds,omitempty" jsonschema:"Wall-clock limit for this call; 0 uses the server default."`
}

type runOutput struct {
	ExitCode        int              `json:"exit_code"`
	Stdout          string           `json:"stdout"`
	StdoutTruncated bool             `json:"stdout_truncated,omitempty"`
	Stderr          string           `json:"stderr,omitempty"`
	StderrTruncated bool             `json:"stderr_truncated,omitempty"`
	TimedOut        bool             `json:"timed_out,omitempty"`
	Error           *structuredError `json:"error,omitempty"`
}

func runHandler(config *serverConfig) mcpsdk.ToolHandlerFor[runInput, runOutput] {
	return func(
		ctx context.Context,
		_ *mcpsdk.CallToolRequest,
		input runInput,
	) (*mcpsdk.CallToolResult, runOutput, error) {
		if len(input.Args) == 0 {
			return nil, runOutput{}, fmt.Errorf("args must be non-empty")
		}

		result, err := runMlrSubprocess(ctx, config, input.Args, input.StdinText, input.TimeoutSeconds)
		if err != nil {
			return nil, runOutput{}, err
		}

		return nil, runOutput{
			ExitCode:        result.exitCode,
			Stdout:          result.stdout,
			StdoutTruncated: result.stdoutTruncated,
			Stderr:          result.stderr,
			StderrTruncated: result.stderrTruncated,
			TimedOut:        result.timedOut,
			Error:           result.structured,
		}, nil
	}
}
