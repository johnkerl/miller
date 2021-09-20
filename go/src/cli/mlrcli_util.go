package cli

import (
	"fmt"
	"os"

	"mlr/src/lib"
)

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func CheckArgCount(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", "mlr", args[argi])
		fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", "mlr")
		os.Exit(1)
	}
}

// ----------------------------------------------------------------
// TODO: give symbolic name to all the RHSes

var SEPARATOR_NAMES_TO_VALUES = map[string]string{
	"colon":     ":",
	"comma":     ",",
	"cr":        "\\r",
	"crcr":      "\\r\\r",
	"crlf":      "\\r\\n",
	"crlfcrlf":  "\\r\\n\\r\\n",
	"equals":    "=",
	"lf":        "\\n",
	"lflf":      "\\n\\n",
	"newline":   "\\n",
	"pipe":      "|",
	"semicolon": ";",
	"slash":     "/",
	"space":     " ",
	"tab":       "\\t",

	"ascii_null": "\\x01",
	"ascii_soh":  "\\x02",
	"ascii_stx":  "\\x03",
	"ascii_etx":  "\\x04",

	"ascii_esc":  "\\x1b",
	"ascii_fs":   "\\x1c",
	"ascii_gs":   "\\x1d",
	"ascii_rs":   "\\x1e",
	"ascii_us":   "\\x1f",
}

func SeparatorFromArg(name string) string {
	sep, ok := SEPARATOR_NAMES_TO_VALUES[name]
	if !ok {
		// "\001" -> control-A, etc.
		return lib.UnbackslashStringLiteral(name)
	}
	return sep
}
