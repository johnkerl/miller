// ================================================================
// TOOO
// ================================================================

package regtest

import (
	"fmt"
	"os"
	"path"
	"strings"
)

// ================================================================
func RegTestUsage(verbName string, o *os.File, exitCode int) {
	exeName := path.Base(os.Args[0])
	fmt.Fprintf(o, "Usage: %s %s [options] [one or more directories/files]\n", exeName, verbName)
	fmt.Fprintf(o, "If no directories/files are specified, the directory %s is used by default.\n", DefaultPath)
	fmt.Fprintf(o, "Recursively walks the directory/ies looking for foo.cmd files having Miller command-lines,\n")
	fmt.Fprintf(o, "with foo.expout and foo.experr files having expected stdout and stderr, respectively.\n")
	fmt.Fprintf(o, "If foo.should-fail exists and is a file, the command is expected to exit non-zero back to\n")
	fmt.Fprintf(o, "the shell.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "[none] Print directory-level pass/fails, and overall pass/fail.\n")
	fmt.Fprintf(o, "-v     Also include pass/fail at command-file level.\n")
	fmt.Fprintf(o, "-vv    Also include pass/fail reasons for each command-file.\n")
	fmt.Fprintf(o, "-vvv   Also include full stdout/stderr/exit-code for each command-file.\n")
	os.Exit(exitCode)
}

// Here the args are the full Miller command line: "mlr regtest --foo bar".
func RegTestMain(args []string) int {

	exeName := path.Base(args[0])
	verbName := args[1]
	argc := len(args)
	argi := 2
	verbosityLevel := 0

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		arg := args[argi]

		if !strings.HasPrefix(arg, "-") {
			break // No more flag options to process
		}
		argi++

		if arg == "-h" || arg == "--help" {
			RegTestUsage(verbName, os.Stdout, 0)

		} else if arg == "-v" {
			verbosityLevel++
		} else if arg == "-vv" {
			verbosityLevel += 2
		} else if arg == "-vvv" {
			verbosityLevel += 3
		} else if arg == "-vvvv" {
			verbosityLevel += 4

		} else {
			RegTestUsage(verbName, os.Stderr, 1)
		}
	}
	paths := args[argi:]

	regtester := NewRegTester(
		exeName,
		verbName,
		verbosityLevel,
	)

	ok := regtester.Execute(paths)

	if !ok {
		return 1
	}

	return 0
}
