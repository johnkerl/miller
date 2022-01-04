package climain

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/lib"
	shellquote "github.com/kballard/go-shellquote"
)

// maybeInterpolateDashS supports Miller scripts with shebang lines like
//   #!/usr/bin/env mlr -s
//   --csv tac then filter '
//     NR % 2 == 1
//   '
// invoked as
//   scriptfile input1.csv input2.csv
// The "-s" flag must be the very first command-line argument after "mlr" for
// two reasons:
// * This is how shebang lines work
// * There are Miller verbs with -s flags and we don't want to disrupt their behavior.
func maybeInterpolateDashS(args []string) []string {
	if len(args) < 2 {
		return args
	}
	if args[1] != "-s" { // Normal case
		return args
	}
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "mlr: -s flag requires a filename after it.\n")
		os.Exit(1)
	}

	// mlr -s scriptfile input1.csv input2.csv
	// 0   1  2          3          4
	arg0 := args[0]
	filename := args[2]
	remainingArgs := args[3:]

	// Read the bytes in the filename given after -s.
	byteContents, rerr := ioutil.ReadFile(filename)
	if rerr != nil {
		fmt.Fprintf(os.Stderr, "mlr: cannot read %s: %v\n", filename, rerr)
		os.Exit(1)
	}
	contents := string(byteContents)

	// Split into lines
	contents = strings.ReplaceAll(contents, "\r\n", "\n")
	lines := lib.SplitString(contents, "\n")

	// Remove the shebang line itself.
	if len(lines) >= 1 {
		if strings.HasPrefix(lines[0], "#!") {
			lines = lines[1:]
		}
	}

	// TODO: maybe support comment lines deeper within the script-file.
	// Make sure they're /^[\s]+#/ since we don't want to disrupt a "#" within
	// strings which are not actually comment characters.

	// Re-join lines to strings, and pass off to a shell-parser to split into
	// an args[]-style array.
	contents = strings.Join(lines, "\n")
	argsFromFile, err := shellquote.Split(contents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: cannot parse %s: %v\n", filename, err)
		os.Exit(1)
	}

	// Join "mlr", the args from the script-file contents, and all the remaining arguments
	// on the original command line after "mlr -s {scriptfile}"
	newArgs := make([]string, 1)
	newArgs[0] = arg0

	for _, argFromFile := range argsFromFile {
		newArgs = append(newArgs, argFromFile)
	}

	// So people can have verb-chains in their shebang file and flags like
	// --icsv --ojson after
	newArgs = append(newArgs, "--")

	for _, remainingArg := range remainingArgs {
		newArgs = append(newArgs, remainingArg)
	}

	return newArgs
}
