// ================================================================
// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.
// ================================================================

package auxents

import (
	"fmt"
	"os"
	"path"

	"miller/src/auxents/regtest"
	"miller/src/auxents/repl"
)

// ----------------------------------------------------------------
type tAuxMain func(args []string) int
type tAuxUsage func(verbName string, o *os.File, exitCode int)

type tAuxLookupEntry struct {
	name  string
	main  tAuxMain
	usage tAuxUsage
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var _AUX_LOOKUP_TABLE = []tAuxLookupEntry{}

func init() {
	_AUX_LOOKUP_TABLE = []tAuxLookupEntry{
		{"aux-list", auxListMain, auxListUsage},
		{"hex", hexMain, hexUsage},
		{"lecat", lecatMain, lecatUsage},
		{"regtest", regtest.RegTestMain, regtest.RegTestUsage},
		{"repl", repl.ReplMain, repl.ReplUsage},
		{"termcvt", termcvtMain, termcvtUsage},
		{"unhex", unhexMain, unhexUsage},
	}
}

// ----------------------------------------------------------------
func mlrExeName() string {
	return path.Base(os.Args[0])
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
