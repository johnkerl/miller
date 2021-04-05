// ================================================================
// Handling single quotes and double quotes is different on Windows unless
// particular care is taken, which is what this file does.
// ================================================================

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

	// Things on the command line include: args[0], which is fine as-is;
	// various flags like -x or --xyz which are fine as-is; DSL expressions
	// such as '$a = $b . "ccc"'.  We don't want to give back the result of
	// shellquote.Split since that will remove the backslashes from things like
	// ..\data\foo\dat or C:\foo\bar.baz.
	rawCommandLine := windows.UTF16PtrToString(windows.GetCommandLine())
	newArgs, err := shellquote.Split(rawCommandLine)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"mlr: internal error: could not parse Windows raw command line: %v\n",
			err,
		)
	}

	for i, oldArg := range os.Args {
		if strings.HasPrefix(oldArg, "'") && strings.HasSuffix(oldArg, "'") {
			os.Args[i] = newArgs[i]
		}
	}
	return os.Args
}
