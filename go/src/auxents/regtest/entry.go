// ================================================================
// TOOO
// ================================================================

package regtest

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// ================================================================
func regTestUsage(verbName string, o *os.File, exitCode int) {
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
	fmt.Fprintf(o, "-m {...} Specify name of Miller executable to use.\n")
	fmt.Fprintf(o, "-c       Shorthand for -m ../c/mlr.\n")
	fmt.Fprintf(o, "-p       Create the .expout and .experr files, rather than checking them.\n")
	fmt.Fprintf(o, "-v       Also include pass/fail at command-file level.\n")
	fmt.Fprintf(o, "-vv      Also include pass/fail reasons for each command-file.\n")
	fmt.Fprintf(o, "-vvv     Also include full stdout/stderr/exit-code for each command-file.\n")
	fmt.Fprintf(o, "-j       Just show the Miller command-line, put/filter script if any, and output.\n")
	fmt.Fprintf(o, "-s {n}   After running tests, re-run first n failed .cmd files with verbosity level 3.\n")
	fmt.Fprintf(o, "-S       After running tests, re-run all failed .cmd files with verbosity level 3.\n")
	os.Exit(exitCode)
}

// Here the args are the full Miller command line: "mlr regtest --foo bar".
func RegTestMain(args []string) int {

	exeName := args[0]
	verbName := args[1]
	argc := len(args)
	argi := 2
	verbosityLevel := 0
	doPopulate := false
	plainMode := false
	firstNFailsToShow := 0

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		arg := args[argi]

		if !strings.HasPrefix(arg, "-") {
			break // No more flag options to process
		}
		argi++

		if arg == "-h" || arg == "--help" {
			regTestUsage(verbName, os.Stdout, 0)

		} else if arg == "-m" {
			if argi >= argc {
				regTestUsage(verbName, os.Stderr, 1)
			}
			exeName = args[argi]
			argi++

		} else if arg == "-c" {
			exeName = "../c/mlr"

		} else if arg == "-g" {
			exeName = "../go/mlr"

		} else if arg == "-s" {
			if argi >= argc {
				regTestUsage(verbName, os.Stderr, 1)
			}
			temp, err := strconv.Atoi(args[argi])
			if err != nil {
				regTestUsage(verbName, os.Stderr, 1)
			}
			firstNFailsToShow = temp
			argi++

		} else if arg == "-S" {
			firstNFailsToShow = 1000000000

		} else if arg == "-p" {
			doPopulate = true

		} else if arg == "-v" {
			verbosityLevel++

		} else if arg == "-j" {
			plainMode = true

		} else {
			regTestUsage(verbName, os.Stderr, 1)
		}
	}
	paths := args[argi:]

	regtester := NewRegTester(
		exeName,
		doPopulate,
		verbosityLevel,
		plainMode,
		firstNFailsToShow,
	)

	ok := regtester.Execute(paths)

	if !ok {
		return 1
	}

	return 0
}
