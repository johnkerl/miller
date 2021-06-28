// ================================================================
// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.
// ================================================================

package auxents

import (
	"fmt"
	"os"

	"miller/src/auxents/help"
	"miller/src/auxents/regtest"
	"miller/src/auxents/repl"
)

// ----------------------------------------------------------------
type tAuxMain func(args []string) int

type tAuxLookupEntry struct {
	name string
	main tAuxMain
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
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
	}
}

// ----------------------------------------------------------------
func mlrExeName() string {
	// TODO:
	// This is ideal, so if someone has a 'mlr.debug' or somesuch, the messages will reflect that:

	// return path.Base(os.Args[0])

	// ... however it makes automated regression-testing hard, cross-platform. For example,
	// 'mlr' vs 'C:\something\something\mlr.exe'.
	return "mlr"
}

// ================================================================
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

	// Else, return control to mlr.go for the rest of Miller.
}

// ================================================================
func auxListUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", mlrExeName(), verbName)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h or --help: print this message\n")
	os.Exit(exitCode)
}

func auxListMain(args []string) int {
	ShowAuxEntries(os.Stdout)
	return 0
}

// This symbol is exported for 'mlr --help'.
func ShowAuxEntries(o *os.File) {
	fmt.Fprintf(o, "Available subcommands:\n")
	for _, entry := range _AUX_LOOKUP_TABLE {
		fmt.Fprintf(o, "  %s\n", entry.name)
	}

	fmt.Fprintf(
		o,
		"For more information, please invoke %s {subcommand} --help.\n",
		mlrExeName(),
	)
}
