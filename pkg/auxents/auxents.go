// ================================================================
// Support for Miller regression testing. Originally bash scripts; ported to Go
// for ease of Windows-native testing.
// ================================================================

package auxents

import (
	"fmt"
	"os"
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

// auxListMain is the handler for 'mlr aux-list'.
func auxListMain(args []string) int {
	ShowAuxEntries(os.Stdout)
	return 0
}

// ShowAuxEntries is a symbol is exported for 'mlr --help'.
func ShowAuxEntries(o *os.File) {
	fmt.Fprintf(o, "Available entries:\n")
	for _, entry := range _AUX_LOOKUP_TABLE {
		fmt.Fprintf(o, "  mlr %s\n", entry.name)
	}

	fmt.Fprintf(o, "For more information, please invoke mlr {subcommand} --help.\n")
}
