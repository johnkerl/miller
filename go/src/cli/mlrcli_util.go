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
var SEPARATOR_NAMES_TO_VALUES = map[string]string{
	"colon":     ":",
	"comma":     ",",
	"cr":        "\r",
	"crcr":      "\r\r",
	"crlf":      "\r\n",
	"crlfcrlf":  "\r\n\r\n",
	"equals":    "=",
	"lf":        "\n",
	"lflf":      "\n\n",
	"newline":   "\n",
	"pipe":      "|",
	"semicolon": ";",
	"slash":     "/",
	"space":     " ",
	"tab":       "\t",
}

var SEPARATOR_NAMES_TO_VALUES_FOR_ONLINE_HELP = map[string]string{
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
}

func SeparatorFromArg(name string) string {
	sep, ok := SEPARATOR_NAMES_TO_VALUES[name]
	if !ok {
		// "\001" -> control-A, etc.
		return lib.UnbackslashStringLiteral(name)
	}
	return sep
}
