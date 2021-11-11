// ================================================================
// Handling single quotes and double quotes is different on Windows unless
// particular care is taken, which is what this file does.
// ================================================================

//go:build windows
// +build windows

package platform

import (
	"fmt"
	"os"
	"strings"

	shellquote "github.com/kballard/go-shellquote"
	"golang.org/x/sys/windows"
)

// GetArgs returns a copy of os.Args, as-is except for any arg wrapped in single quotes.
// This is for compatibility Linux/Unix/MacOS.
//
// For example:
//
//   mlr --icsv --ojson put '$filename = $basename . ".ext"' ..\data\foo.csv
//
// The Windows way to say that is
//
//   mlr --icsv --ojson put "$filename = $basename . """.ext"""" ..\data\foo.csv
//
// This function makes it possible to say the former, or the latter.

func GetArgs() []string {
	// If this code is running in MSYS2: MSYS2 does the right thing already and
	// we won't try to improve on that.  This is true regardless of whether we
	// were compiled inside MSYS2 or outside; what matters is whether we're
	// _running_ inside MSYS2 or outside. (I.e.  this is necessarily a run-time
	// check, not a compile-time check.)
	msystem := os.Getenv("MSYSTEM")
	if msystem != "" {
		return os.Args
	}

	//printArgs(os.Args, "ORIGINAL")

	regrouped, ok := regroupForSingleQuote(os.Args)
	if !ok {
		return os.Args
	}
	//printArgs(regrouped, "REGROUPED")

	// Things on the command line include: args[0], which is fine as-is;
	// various flags like -x or --xyz which are fine as-is; DSL expressions
	// such as '$a = $b . "ccc"'.  We don't want to give back the result of
	// shellquote.Split since that will remove the backslashes from things like
	// ..\data\foo\dat or C:\foo\bar.baz.
	rawCommandLine := windows.UTF16PtrToString(windows.GetCommandLine())
	splitArgs, err := shellquote.Split(rawCommandLine)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"mlr: internal error: could not parse Windows raw command line: %v\n",
			err,
		)
	}

	retargs := make([]string, 0)

	// TODO err/stetret if lens uneq

	for i, oldArg := range regrouped {
		if strings.HasPrefix(oldArg, "'") && strings.HasSuffix(oldArg, "'") {
			retargs = append(retargs, splitArgs[i])
		} else {
			retargs = append(retargs, regrouped[i])
		}
	}
	//printArgs(retargs, "NEW")
	return retargs
}

// ----------------------------------------------------------------
func printArgs(args []string, description string) {
	fmt.Printf("%s:\n", description)
	for i, arg := range args {
		fmt.Printf("%d %s\n", i, arg)
	}
	fmt.Println()
}

// ----------------------------------------------------------------
func regroupForSingleQuote(inargs []string) ([]string, bool) {
	outargs := make([]string, 0, len(inargs))
	inside := false
	var concat string

	// TODO: comment
	// TODO: UT this all
	for _, inarg := range inargs {
		if !inside {
			if !strings.HasPrefix(inarg, "'") {
				// Current arg is not single-quoted, and not inside a single-quoted region
				outargs = append(outargs, inarg)

			} else {

				// Start of single-quoted region
				if strings.HasSuffix(inarg, "'") {
					// Start and end of single-quoted region, like '$y=$x'
					outargs = append(outargs, inarg)
				} else {
					// Start but not end of single-quoted region, like '$y=
					inside = true
					concat = inarg
				}

			}

		} else {
			if !strings.HasSuffix(inarg, "'") {
				// Continuation of single-quoted region
				concat = concat + " " + inarg

			} else {
				// End of single-quoted region
				inside = false
				concat = concat + " " + inarg
				outargs = append(outargs, concat)
				concat = ""
			}
		}
	}

	// TODO: error not bool?
	if inside {
		return nil, false
	}

	return outargs, true
}
