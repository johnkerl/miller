// Server construction for `mlr mcp`: tool, prompt, and resource registration.

package mcp

import (
	"context"
	_ "embed"

	"github.com/johnkerl/miller/v6/pkg/version"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed SKILL.md
var playbookText string

const playbookPromptName = "miller-playbook"
const playbookResourceURI = "miller://playbook"

const serverInstructions = `Miller (mlr) processes CSV/TSV/JSON/etc. from the command line.
Work the loop: describe_data to learn the input's fields and values; which or
list_capabilities (index first, then one full entry) to pick a verb;
validate_dsl before running any put/filter expression; then run. On a run
error, branch on the structured error's kind/hint/did_you_mean rather than
re-guessing. Get the full playbook from the miller-playbook prompt or the
miller://playbook resource. Never invent flag or function names: everything
valid is in list_capabilities.`

func boolPtr(b bool) *bool { return &b }

// newServer assembles the MCP server: five tools, the playbook prompt, and
// the playbook resource.
func newServer(config *serverConfig) *mcpsdk.Server {
	server := mcpsdk.NewServer(
		&mcpsdk.Implementation{
			Name:    "miller",
			Title:   "Miller (mlr)",
			Version: version.STRING,
		},
		&mcpsdk.ServerOptions{
			Instructions: serverInstructions,
		},
	)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name: "list_capabilities",
		Description: "List Miller's verbs, DSL functions, main flags, and DSL keywords " +
			"as structured JSON (the `mlr help --as-json` catalog). " +
			"With index=true, returns only {kind, name, summary} triples across the whole surface -- " +
			"the cheap first call. With kind (and optionally names), returns full entries for " +
			"just those items. With no arguments, returns the entire catalog (large). " +
			"Results are fully cacheable against (mlr_version, catalog_schema_version).",
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: boolPtr(false)},
	}, listCapabilitiesHandler)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name: "which",
		Description: "Route a natural-language intent (e.g. \"join two files on a key\") to ranked " +
			"Miller capabilities: verbs, functions, flags, keywords. Returns {confident, results}; " +
			"when confident is true the top result's name matched a query word. " +
			"Use before drilling into list_capabilities.",
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: boolPtr(false)},
	}, whichHandler)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name: "validate_dsl",
		Description: "Parse and type-check a Miller put/filter DSL expression without reading any " +
			"input data (mlr put --explain). Returns {valid} on success, or {valid: false, error} " +
			"with a structured error document: kind, hint, did_you_mean. " +
			"Always validate before using an expression in run.",
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true, OpenWorldHint: boolPtr(false)},
	}, validateDSLHandler(config))

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name: "describe_data",
		Description: "Learn the shape of input data before constructing a command " +
			"(mlr describe): per-field name, types seen with counts, occurrence count, null count, " +
			"cardinality, min/max, and -- for low-cardinality fields -- every distinct value. " +
			"Copy field names and values from here instead of guessing them. " +
			"Pass input file paths via the files field, and optionally the input format " +
			"(csv, tsv, json, etc.) via input_format.",
		Annotations: &mcpsdk.ToolAnnotations{ReadOnlyHint: true, IdempotentHint: true},
	}, describeDataHandler(config))

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name: "run",
		Description: "Run an mlr command line. Pass argv after the leading \"mlr\" via the args field, " +
			"one element per shell word, e.g. [\"--icsv\", \"--ojson\", \"cat\", \"data.csv\"]; optional " +
			"inline input via the stdin_text field. Returns {exit_code, stdout, stderr} plus a parsed structured error " +
			"when the command fails. Output is truncated at the server's --max-output-bytes; " +
			"execution is time-limited. External-command execution (the DSL system/exec functions, " +
			"piped redirects, --prepipe) is disabled unless the server was started with --allow-shell. " +
			"Commands can still write files via tee/split/--from-less redirection outputs.",
		Annotations: &mcpsdk.ToolAnnotations{DestructiveHint: boolPtr(true)},
	}, runHandler(config))

	server.AddPrompt(&mcpsdk.Prompt{
		Name:        playbookPromptName,
		Title:       "Miller agent playbook",
		Description: "How to drive Miller effectively as an agent: the discover -> constrain -> validate -> run loop.",
	}, playbookPromptHandler)

	server.AddResource(&mcpsdk.Resource{
		URI:         playbookResourceURI,
		Name:        playbookPromptName,
		Title:       "Miller agent playbook",
		Description: "How to drive Miller effectively as an agent: the discover -> constrain -> validate -> run loop.",
		MIMEType:    "text/markdown",
	}, playbookResourceHandler)

	return server
}

func playbookPromptHandler(_ context.Context, _ *mcpsdk.GetPromptRequest) (*mcpsdk.GetPromptResult, error) {
	return &mcpsdk.GetPromptResult{
		Description: "Miller agent playbook",
		Messages: []*mcpsdk.PromptMessage{
			{
				Role:    "user",
				Content: &mcpsdk.TextContent{Text: playbookText},
			},
		},
	}, nil
}

func playbookResourceHandler(_ context.Context, _ *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      playbookResourceURI,
				MIMEType: "text/markdown",
				Text:     playbookText,
			},
		},
	}, nil
}
