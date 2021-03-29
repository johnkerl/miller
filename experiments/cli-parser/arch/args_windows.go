// +build windows

package arch

import (
	"fmt"
	"os"

	shellquote "github.com/kballard/go-shellquote"
	"golang.org/x/sys/windows"
)

func GetMainArgs() []string {
	commandLine := windows.UTF16PtrToString(windows.GetCommandLine())

	split, err := shellquote.Split(commandLine)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s: could not parse obtain Windows command line: %v\n",
			os.Args[0], err,
		)
		os.Exit(1)
	}
	return split
}
