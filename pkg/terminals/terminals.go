// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.

package terminals

import (
	"fmt"
	"os"
	"runtime"

	"github.com/johnkerl/miller/v6/pkg/terminals/completion"
	"github.com/johnkerl/miller/v6/pkg/terminals/help"
	"github.com/johnkerl/miller/v6/pkg/terminals/mcp"
	"github.com/johnkerl/miller/v6/pkg/terminals/registry"
	"github.com/johnkerl/miller/v6/pkg/terminals/regtest"
	"github.com/johnkerl/miller/v6/pkg/terminals/repl"
	"github.com/johnkerl/miller/v6/pkg/terminals/script"
	"github.com/johnkerl/miller/v6/pkg/terminals/skill"
	"github.com/johnkerl/miller/v6/pkg/version"
)

// tTerminalMain is a function-pointer type for the entrypoint handler for a given terminal,
// such as 'help' or 'regtest'.
type tTerminalMain func(args []string) int

type tTerminalLookupEntry struct {
	name string
	main tTerminalMain
}

// _TERMINAL_LOOKUP_TABLE is the lookup table for terminals. We get a Golang
// "initialization loop" if this is defined statically. So, we use a "package
// init" function.
var _TERMINAL_LOOKUP_TABLE = []tTerminalLookupEntry{}

func init() {
	_TERMINAL_LOOKUP_TABLE = []tTerminalLookupEntry{
		{registry.TerminalList, terminalListMain},
		{registry.Completion, completion.CompletionMain},
		{registry.Help, help.HelpMain},
		{registry.Mcp, mcp.McpMain},
		{registry.Regtest, regtest.RegTestMain},
		{registry.Repl, repl.ReplMain},
		{registry.Script, script.ScriptMain},
		{registry.Skill, skill.SkillMain},
		{registry.Version, showVersion},
		{registry.Which, help.WhichMain},
	}
}

func Dispatchable(arg string) bool {
	for _, entry := range _TERMINAL_LOOKUP_TABLE {
		if arg == entry.name {
			return true
		}
	}
	return false
}

// Dispatch runs the terminal named by args[0] and returns the process exit
// code for the caller to propagate (via lib.ExitRequest) up to the entrypoint
// layer.
func Dispatch(args []string) int {
	if len(args) < 1 {
		// Can't happen: the climain caller passes a non-empty terminal sequence.
		fmt.Fprintf(os.Stderr, "mlr: internal coding error: empty terminal sequence.\n")
		return 1
	}
	terminal := args[0]

	for _, entry := range _TERMINAL_LOOKUP_TABLE {
		if terminal == entry.name {
			return entry.main(args)
		}
	}
	fmt.Fprintf(os.Stderr, "mlr: terminal \"%s\" not found.\n", terminal)
	return 1
}

// terminalListMain is the handler for 'mlr terminal-list'.
func terminalListMain(args []string) int {
	ShowTerminalEntries(os.Stdout)
	return 0
}

// ShowTerminalEntries is a symbol is exported for 'mlr --help'.
func ShowTerminalEntries(o *os.File) {
	fmt.Fprintf(o, "Available subcommands:\n")
	for _, entry := range _TERMINAL_LOOKUP_TABLE {
		fmt.Fprintf(o, "  %s\n", entry.name)
	}

	fmt.Fprintf(o, "For more information, please invoke mlr {subcommand} --help.\n")
}

func showVersion(args []string) int {
	fmt.Printf("mlr version %s for %s/%s/%s\n",
		version.STRING, runtime.GOOS, runtime.GOARCH, runtime.Version())
	return 0
}
