package cli

import (
	"fmt"
	"os"
)

// For flags with values, e.g. ["-n" "10"], while we're looking at the "-n"
// this let us see if the "10" slot exists.
func checkArgCount(args []string, argi int, argc int, n int) {
	if (argc - argi) < n {
		fmt.Fprintf(os.Stderr, "%s: option \"%s\" missing argument(s).\n", args[0], args[argi])
		mainUsageShort()
		os.Exit(1)
	}
}

func SeparatorFromArg(arg string) string {
	// xxx stub
	return arg
	//	char* chars = lhmss_get(get_desc_to_chars_map(), arg);
	//	if (chars != nil) // E.g. crlf
	//		return mlr_strdup_or_die(chars);
	//	else // E.g. '\r\n'
	//		return mlr_alloc_unbackslash(arg);
}
