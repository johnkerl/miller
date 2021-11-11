// ================================================================
// Utilities for Miller verbs to share for command-line parsing.
// ================================================================

package cli

import (
	"fmt"
	"os"
	"strconv"

	"mlr/internal/pkg/lib"
)

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n" this let us see if the "10" slot exists.
// The verb is nominally something from a ways earlier in args[]; the opt is nominally what's at args[argi-1].
// So this function should be called with args[argi] pointing to the "10" slot.
func verbCheckArgCount(verb string, opt string, args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s %s: option \"%s\" missing argument(s).\n",
			"mlr", verb, opt,
		)
		os.Exit(1)
	}
}

// E.g. with ["-f", "a,b,c"], makes sure there is something in the "a,b,c" position, and returns it.
func VerbGetStringArgOrDie(verb string, opt string, args []string, pargi *int, argc int) string {
	verbCheckArgCount(verb, opt, args, *pargi, argc, 1)
	retval := args[*pargi]
	*pargi += 1
	return retval
}

// E.g. with ["-f", "a,b,c"], makes sure there is something in the "a,b,c" position,
// splits it on commas, and returns it.
func VerbGetStringArrayArgOrDie(verb string, opt string, args []string, pargi *int, argc int) []string {
	stringArg := VerbGetStringArgOrDie(verb, opt, args, pargi, argc)
	return lib.SplitString(stringArg, ",")
}

// E.g. with ["-n", "10"], makes sure there is something in the "10" position,
// scans it as int, and returns it.
func VerbGetIntArgOrDie(verb string, opt string, args []string, pargi *int, argc int) int {
	flag := args[*pargi]
	stringArg := VerbGetStringArgOrDie(verb, opt, args, pargi, argc)
	retval, err := strconv.Atoi(stringArg)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s %s: could not scan flag \"%s\" argument \"%s\" as int.\n",
			"mlr", verb, flag, stringArg,
		)
		os.Exit(1)
	}
	return retval
}

// E.g. with ["-n", "10.3"], makes sure there is something in the "10.3"
// position, scans it as float, and returns it.
func VerbGetFloatArgOrDie(verb string, opt string, args []string, pargi *int, argc int) float64 {
	flag := args[*pargi]
	stringArg := VerbGetStringArgOrDie(verb, opt, args, pargi, argc)
	retval, err := strconv.ParseFloat(stringArg, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s %s: could not scan flag \"%s\" argument \"%s\" as float.\n",
			"mlr", verb, flag, stringArg,
		)
		os.Exit(1)
	}
	return retval
}
