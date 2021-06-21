package cli

import (
	"fmt"
	"os"

	"miller/src/dsl/cst"
	"miller/src/lib"
)

// ----------------------------------------------------------------
func mainUsageShort(o *os.File, exitCode int) {
	fmt.Fprintf(o, "Please run \"%s --help\" for detailed usage information.\n", lib.MlrExeName())
	os.Exit(exitCode)
}

// ----------------------------------------------------------------
func mainUsageHelpOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o, "  -h or --help                 Show this message.\n")
	fmt.Fprintf(o, "  --version                    Show the software version.\n")
	fmt.Fprintf(o, "  {verb name} --help           Show verb-specific help.\n")
	fmt.Fprintf(o, "  --help-all-verbs             Show help on all verbs.\n")
	fmt.Fprintf(o, "  -l or --list-all-verbs       List only verb names.\n")
	fmt.Fprintf(o, "  -L                           List only verb names, one per line.\n")
	fmt.Fprintf(o, "  -f or --help-all-functions   Show help on all built-in functions.\n")
	fmt.Fprintf(o, "  -F                           Show a bare listing of built-in functions by name.\n")
	fmt.Fprintf(o, "  -k or --help-all-keywords    Show help on all keywords.\n")
	fmt.Fprintf(o, "  -K                           Show a bare listing of keywords by name.\n")
}

func mainUsageFunctions(o *os.File) {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
	fmt.Fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", lib.MlrExeName())
}

func usageUnrecognizedVerb(argv0 string, arg string) {
	fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", argv0, arg)
	fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", argv0)
	os.Exit(1)
}
