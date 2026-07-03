// Entrypoint for `mlr mcp`: a Model Context Protocol server exposing Miller
// to AI agents.
//
// Transport is JSON-RPC 2.0 over stdin/stdout ("stdio" in MCP terms): the MCP
// client -- Claude Code, Claude Desktop, Cursor, etc. -- spawns `mlr mcp` as a
// subprocess and speaks the protocol over its standard streams. No network
// port is opened. Register with e.g.:
//
//	claude mcp add miller -- mlr mcp
//
// Tools exposed (see server.go): list_capabilities, which, validate_dsl,
// describe_data, run. The first two are served in-process from the same
// registries behind `mlr help --as-json`; the rest shell out to this same mlr
// binary so agents see byte-identical behavior to a terminal. Subprocesses
// run with MLR_NO_SHELL=1 (unless --allow-shell) so an agent-constructed
// command line cannot execute external commands, plus MLR_ERRORS_JSON=1 so
// failures come back structured.

package mcp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

const verbNameMCP = "mcp"

const (
	defaultTimeoutSeconds = 60
	defaultMaxOutputBytes = 1024 * 1024
)

// serverConfig carries the `mlr mcp` command-line options into the tool
// handlers.
type serverConfig struct {
	// mlrExePath is the mlr binary used for the subprocess-backed tools
	// (validate_dsl, describe_data, run): this same executable, so tool
	// behavior is version-locked to the server.
	mlrExePath string

	// allowShell, when true, omits the MLR_NO_SHELL=1 injection on
	// subprocesses, re-enabling the DSL system/exec functions, piped
	// redirects, and --prepipe for agent-run commands.
	allowShell bool

	// timeoutSeconds is the default wall-clock limit for the run tool
	// (overridable per-call via its timeout_seconds input).
	timeoutSeconds int64

	// maxOutputBytes caps captured stdout/stderr per subprocess; output
	// beyond the cap is discarded and reported as truncated.
	maxOutputBytes int64
}

func mcpUsage(o *os.File) {
	fmt.Fprintf(o, "Usage: mlr mcp [options]\n")
	fmt.Fprintf(o, "Runs a Model Context Protocol (MCP) server exposing Miller to AI agents.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "The server speaks JSON-RPC over stdin/stdout (MCP \"stdio\" transport); it is\n")
	fmt.Fprintf(o, "meant to be spawned by an MCP client rather than run interactively. Example\n")
	fmt.Fprintf(o, "registration, for Claude Code:\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "  claude mcp add miller -- mlr mcp\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Tools exposed:\n")
	fmt.Fprintf(o, "  list_capabilities  the mlr help --as-json catalog/index, filterable\n")
	fmt.Fprintf(o, "  which              intent -> ranked matching verbs/functions/flags/keywords\n")
	fmt.Fprintf(o, "  validate_dsl       parse/type-check a put/filter expression without running it\n")
	fmt.Fprintf(o, "  describe_data      field names/types/cardinality/values of input data\n")
	fmt.Fprintf(o, "  run                run an mlr command line\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Also exposed: an agent playbook, as MCP prompt \"miller-playbook\" and MCP\n")
	fmt.Fprintf(o, "resource \"miller://playbook\", encoding the discover -> constrain -> validate\n")
	fmt.Fprintf(o, "-> run loop.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Commands started by the run/validate_dsl/describe_data tools are run with\n")
	fmt.Fprintf(o, "MLR_NO_SHELL=1 -- the DSL system/exec functions, piped redirects, and\n")
	fmt.Fprintf(o, "--prepipe/--prepipex fail cleanly -- and with MLR_ERRORS_JSON=1 so errors\n")
	fmt.Fprintf(o, "come back structured.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, " --allow-shell           Do not set MLR_NO_SHELL=1 on subprocesses, re-enabling\n")
	fmt.Fprintf(o, "                         the DSL system/exec functions for agent-run commands.\n")
	fmt.Fprintf(o, " --timeout {seconds}     Default wall-clock limit for the run tool (default %d).\n", defaultTimeoutSeconds)
	fmt.Fprintf(o, " --max-output-bytes {n}  Cap captured stdout/stderr per command (default %d).\n", defaultMaxOutputBytes)
	fmt.Fprintf(o, " -h or --help            Show this message.\n")
}

// McpMain is the entrypoint called by the terminals dispatcher for `mlr mcp`.
func McpMain(args []string) int {
	config := serverConfig{
		timeoutSeconds: defaultTimeoutSeconds,
		maxOutputBytes: defaultMaxOutputBytes,
	}

	args = args[1:] // strip "mcp"
	argi := 0
	argc := len(args)
	for argi < argc {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			fmt.Fprintf(os.Stderr, "mlr mcp: extraneous argument \"%s\".\n", opt)
			return 1
		}
		argi++
		switch opt {
		case "-h", "--help":
			mcpUsage(os.Stdout)
			return 0
		case "--allow-shell":
			config.allowShell = true
		case "--timeout":
			n, err := cli.VerbGetIntArg(verbNameMCP, opt, args, &argi, argc)
			if err != nil || n <= 0 {
				fmt.Fprintf(os.Stderr, "mlr mcp: option \"%s\" requires a positive integer argument.\n", opt)
				return 1
			}
			config.timeoutSeconds = n
		case "--max-output-bytes":
			n, err := cli.VerbGetIntArg(verbNameMCP, opt, args, &argi, argc)
			if err != nil || n <= 0 {
				fmt.Fprintf(os.Stderr, "mlr mcp: option \"%s\" requires a positive integer argument.\n", opt)
				return 1
			}
			config.maxOutputBytes = n
		default:
			fmt.Fprintf(os.Stderr, "mlr mcp: option \"%s\" not recognized.\n", opt)
			return 1
		}
	}

	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr mcp: could not locate the mlr executable: %v\n", err)
		return 1
	}
	config.mlrExePath = exePath

	server := newServer(&config)
	err = server.Run(context.Background(), &mcpsdk.StdioTransport{})
	if err != nil && !isClientDisconnect(err) {
		fmt.Fprintf(os.Stderr, "mlr mcp: %v\n", err)
		return 1
	}
	return 0
}

// isClientDisconnect reports whether the server-session error just means the
// client closed the connection (EOF on stdin) -- the normal way an MCP stdio
// session ends -- as opposed to a real transport failure.
func isClientDisconnect(err error) bool {
	if errors.Is(err, io.EOF) || errors.Is(err, mcpsdk.ErrConnectionClosed) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "EOF") || strings.Contains(msg, "connection closed")
}
