// ================================================================
// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.
// ================================================================

package terminals

import (
	"fmt"
	"os"
	"runtime"

	"github.com/johnkerl/miller/internal/pkg/terminals/help"
	"github.com/johnkerl/miller/internal/pkg/terminals/regtest"
	"github.com/johnkerl/miller/internal/pkg/terminals/repl"
	"github.com/johnkerl/miller/internal/pkg/version"
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
		{"terminal-list", terminalListMain},
		{"help", help.HelpMain},
		{"regtest", regtest.RegTestMain},
		{"repl", repl.ReplMain},
		{"version", showVersion},
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

func Dispatch(args []string) {
	if len(args) < 1 {
		return
	}
	terminal := args[0]

	for _, entry := range _TERMINAL_LOOKUP_TABLE {
		if terminal == entry.name {
			os.Exit(entry.main(args))
		}
	}
	fmt.Fprintf(os.Stderr, "mlr: terminal \"%s\" not found.\n", terminal)
	os.Exit(1)
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
