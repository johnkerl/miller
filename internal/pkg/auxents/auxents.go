// ================================================================
// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.
// ================================================================

package auxents

import (
	"fmt"
	"os"
	"runtime"

	"github.com/johnkerl/miller/internal/pkg/auxents/help"
	"github.com/johnkerl/miller/internal/pkg/auxents/regtest"
	"github.com/johnkerl/miller/internal/pkg/auxents/repl"
	"github.com/johnkerl/miller/internal/pkg/version"
)

// tAuxMain is a function-pointer type for the entrypoint handler for a given auxent,
// such as 'help' or 'regtest'.
type tAuxMain func(args []string) int

type tAuxLookupEntry struct {
	name string
	main tAuxMain
}

// _AUX_LOOKUP_TABLE is the lookup table for auxents. We get a Golang
// "initialization loop" if this is defined statically. So, we use a "package
// init" function.
var _AUX_LOOKUP_TABLE = []tAuxLookupEntry{}

func init() {
	_AUX_LOOKUP_TABLE = []tAuxLookupEntry{
		{"aux-list", auxListMain},
		{"hex", hexMain},
		{"lecat", lecatMain},
		{"termcvt", termcvtMain},
		{"unhex", unhexMain},
		{"help", help.HelpMain},
		{"regtest", regtest.RegTestMain},
		{"repl", repl.ReplMain},
		{"version", showVersion},
	}
}

// Dispatch is called from Miller main. Here we indicate if argv[1] is handled
// by us, or not. If so, we handle it and exit, not returning control to Miller
// main.
func Dispatch(args []string) {
	if len(args) < 2 {
		return
	}
	verb := args[1]

	for _, entry := range _AUX_LOOKUP_TABLE {
		if verb == entry.name {
			os.Exit(entry.main(args))
		}
	}

	// Else, return control to main for the rest of Miller.
}

// auxListUsage shows the available auxents.
func auxListUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: mlr %s [options]\n", verbName)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

// auxListMain is the handler for 'mlr aux-list'.
func auxListMain(args []string) int {
	ShowAuxEntries(os.Stdout)
	return 0
}

// ShowAuxEntries is a symbol is exported for 'mlr --help'.
func ShowAuxEntries(o *os.File) {
	fmt.Fprintf(o, "Available subcommands:\n")
	for _, entry := range _AUX_LOOKUP_TABLE {
		fmt.Fprintf(o, "  %s\n", entry.name)
	}

	fmt.Fprintf(o, "For more information, please invoke mlr {subcommand} --help.\n")
}

func showVersion(args []string) int {
	fmt.Printf("mlr version %s for %s/%s/%s\n", version.STRING, runtime.GOOS, runtime.GOARCH, runtime.Version())
	return 0
}
