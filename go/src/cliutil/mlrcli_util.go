package cliutil

import (
	"fmt"
	"os"
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
// xxx temp -- still needs '\002' etc

var SEPARATOR_NAMES_TO_VALUES = map[string]string{
	"cr":        "\r",
	"crcr":      "\r\r",
	"newline":   "\n",
	"lf":        "\n",
	"lflf":      "\n\n",
	"crlf":      "\r\n",
	"crlfcrlf":  "\r\n\r\n",
	"tab":       "\t",
	"space":     " ",
	"comma":     ",",
	"pipe":      "|",
	"slash":     "/",
	"colon":     ":",
	"semicolon": ";",
	"equals":    "=",
}

func SeparatorFromArg(name string) string {
	sep, ok := SEPARATOR_NAMES_TO_VALUES[name]
	if !ok {
		// xxx temp
		//fmt.Fprintf(os.Stderr, "Miller: could not handle separator \"%s\".\n", name)
		//os.Exit(1)
		// It's OK if they do '--ifs ,' -- just pass it back.
		return name
	}
	return sep
	//	char* chars = lhmss_get(get_desc_to_chars_map(), arg);
	//	if (chars != nil) // E.g. crlf
	//		return mlr_strdup_or_die(chars);
	//	else // E.g. '\r\n'
	//		return mlr_alloc_unbackslash(arg);
}
