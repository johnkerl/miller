package cli

import (
	"fmt"
	"os"
)

// CheckArgCount is for flags with values, e.g. ["-n" "10"], while we're
// looking at the "-n": this let us see if the "10" slot exists.
func CheckArgCount(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", "mlr", args[argi])
		fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", "mlr")
		os.Exit(1)
	}
}

// SeparatorFromArg is for letting people do things like `--ifs pipe`
// rather than `--ifs '|'`.
func SeparatorFromArg(name string) string {
	sep, ok := SEPARATOR_NAMES_TO_VALUES[name]
	if ok {
		return sep
	} else {
		return name
	}
}
