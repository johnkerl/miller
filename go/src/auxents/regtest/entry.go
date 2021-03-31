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
	fmt.Fprintf(o, "Usage: %s %s [options] {TODO}\n", exeName, verbName)
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
