// ================================================================
// Utilities for Miller verbs to share for command-line parsing.
// ================================================================

package clitypes

import (
	"fmt"
	"os"
	"strconv"

	"miller/lib"
)

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func VerbCheckArgCount(verb string, args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s %s: option \"%s\" missing argument(s).\n",
			args[0], verb, args[argi],
		)
		os.Exit(1)
	}
}

// E.g. with ["-f", "a,b,c"], makes sure there is something in the "a,b,c" position,
// and returns it.
func VerbGetStringArgOrDie(verb string, args []string, pargi *int, argc int) string {
	VerbCheckArgCount(verb, args, *pargi, argc, 2)
	retval := args[*pargi+1]
	*pargi += 2
	return retval
}

// E.g. with ["-f", "a,b,c"], makes sure there is something in the "a,b,c" position,
// splits it on commas, and returns it.
func VerbGetStringArrayArgOrDie(verb string, args []string, pargi *int, argc int) []string {
	return lib.SplitString(VerbGetStringArgOrDie(verb, args, pargi, argc), ",")
}

// E.g. with ["-n", "10"], makes sure there is something in the "10" position,
// scans it as int, and returns it.
func VerbGetIntArgOrDie(verb string, args []string, pargi *int, argc int) int {
	flag := args[*pargi]
	stringArg := VerbGetStringArgOrDie(verb, args, pargi, argc)
	retval, err := strconv.Atoi(stringArg)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"%s %s: could not scan flag \"%s\" argument \"%s\" as int.\n",
			os.Args[0], verb, flag, stringArg,
		)
		os.Exit(1)
	}
	return retval
}
