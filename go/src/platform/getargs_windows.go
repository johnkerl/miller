// +build windows

package platform

import (
	"fmt"
	"os"

	shellquote "github.com/kballard/go-shellquote"
	"golang.org/x/sys/windows"
)

func GetArgs() []string {

	// MSYS2 does the right thing already and we can't improve on that.  This
	// is true regardless of whether we were compiled inside MSYS2 or outside;
	// what matters is whether we're _running_ inside MSYS2 or outside. (I.e.
	// this is a run-time check, not a compile-time check.)
	msystem := os.Getenv("MSYSTEM")
	if msystem != "" {
		return os.Args
	}

	rawCommandLine := windows.UTF16PtrToString(windows.GetCommandLine())
	args, err := shellquote.Split(os.Args)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: internal error: could not parse Windows raw command line: %v\n",
			os.Args[0], err,
		)
	}
	return args
}
