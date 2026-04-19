// Utilities for Miller verbs to share for command-line parsing.
// These return error instead of os.Exit, so callers (e.g. transformer ParseCLIFunc)
// can propagate errors to the CLI entrypoint layer.

package cli

import (
	"fmt"
	"strconv"

	"github.com/johnkerl/miller/v6/pkg/lib"
)

// VerbCheckArgCount returns an error if there aren't enough args remaining.
// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this lets us see if the "10" slot exists. The verb is nominally something
// from a ways earlier in args[]; the opt is nominally what's at args[argi-1].
// This function should be called with args[argi] pointing to the "10" slot.
func VerbCheckArgCount(verb string, opt string, args []string, argi int, argc int, n int) error {
	if (argc - argi) < n {
		return fmt.Errorf("%s %s: option \"%s\" missing argument(s)", "mlr", verb, opt)
	}
	return nil
}

// VerbGetStringArg ensures there is something in the value position and returns it.
// E.g. with ["-f", "a,b,c"], returns "a,b,c".
func VerbGetStringArg(verb string, opt string, args []string, pargi *int, argc int) (string, error) {
	if err := VerbCheckArgCount(verb, opt, args, *pargi, argc, 1); err != nil {
		return "", err
	}
	retval := args[*pargi]
	*pargi += 1
	return retval, nil
}

// VerbGetStringArrayArg ensures there is something in the value position,
// splits it on commas, and returns it. E.g. with ["-f", "a,b,c"], returns ["a","b","c"].
func VerbGetStringArrayArg(verb string, opt string, args []string, pargi *int, argc int) ([]string, error) {
	stringArg, err := VerbGetStringArg(verb, opt, args, pargi, argc)
	if err != nil {
		return nil, err
	}
	return lib.SplitString(stringArg, ","), nil
}

// VerbErrorf returns an error prefixed with "mlr {verb}: " so the entrypoint
// can print it without double-prefixing.
func VerbErrorf(verb, format string, args ...interface{}) error {
	return fmt.Errorf("mlr "+verb+": "+format, args...)
}

// VerbGetIntArg ensures there is something in the value position and parses it as int64.
// E.g. with ["-n", "10"], returns 10.
func VerbGetIntArg(verb string, opt string, args []string, pargi *int, argc int) (int64, error) {
	flag := args[*pargi]
	stringArg, err := VerbGetStringArg(verb, opt, args, pargi, argc)
	if err != nil {
		return 0, err
	}
	retval, err := strconv.ParseInt(stringArg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s %s: could not scan flag \"%s\" argument \"%s\" as int",
			"mlr", verb, flag, stringArg)
	}
	return retval, nil
}

// VerbGetFloatArg ensures there is something in the value position and parses it as float64.
// E.g. with ["-n", "10.3"], returns 10.3.
func VerbGetFloatArg(verb string, opt string, args []string, pargi *int, argc int) (float64, error) {
	flag := args[*pargi]
	stringArg, err := VerbGetStringArg(verb, opt, args, pargi, argc)
	if err != nil {
		return 0, err
	}
	retval, err := strconv.ParseFloat(stringArg, 64)
	if err != nil {
		return 0, fmt.Errorf("%s %s: could not scan flag \"%s\" argument \"%s\" as float",
			"mlr", verb, flag, stringArg)
	}
	return retval, nil
}
