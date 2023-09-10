package climain

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/pkg/lib"
	shellquote "github.com/kballard/go-shellquote"
)

// maybeInterpolateDashS supports Miller scripts with shebang lines like
//
//	#!/usr/bin/env mlr -s
//	--csv tac then filter '
//	  NR % 2 == 1
//	'
//
// invoked as
//
//	scriptfile input1.csv input2.csv
//
// The "-s" flag must be the very first command-line argument after "mlr" for
// two reasons:
// * This is how shebang lines work
// * There are Miller verbs with -s flags and we don't want to disrupt their behavior.
func maybeInterpolateDashS(args []string) ([]string, error) {
	stripComments := true

	if len(args) < 2 {
		return args, nil
	}
	if args[1] == "-s" {
		stripComments = true
	} else if args[1] == "--s-no-comment-strip" {
		stripComments = false
	} else { // Normal case
		return args, nil
	}
	if len(args) < 3 {
		return nil, fmt.Errorf("mlr: -s flag requires a filename after it.")
	}

	// mlr -s scriptfile input1.csv input2.csv
	// 0   1  2          3          4
	arg0 := args[0]
	filename := args[2]
	remainingArgs := args[3:]

	// Read the bytes in the filename given after -s.
	byteContents, rerr := ioutil.ReadFile(filename)
	if rerr != nil {
		return nil, fmt.Errorf("mlr: cannot read %s: %v", filename, rerr)
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

	if stripComments {
		re := regexp.MustCompile(`#.*`)
		for i, _ := range lines {
			lines[i] = re.ReplaceAllString(lines[i], "")
		}
	}

	// Re-join lines to strings, and pass off to a shell-parser to split into
	// an args[]-style array.
	contents = strings.Join(lines, "\n")
	argsFromFile, err := shellquote.Split(contents)
	if err != nil {
		return nil, fmt.Errorf("mlr: cannot parse %s: %v", filename, err)
	}

	// Join "mlr", the args from the script-file contents, and all the remaining arguments
	// on the original command line after "mlr -s {scriptfile}"
	newArgs := []string{arg0}
	newArgs = append(newArgs, argsFromFile...)

	// So people can have verb-chains in their shebang file and flags like
	// --icsv --ojson after
	newArgs = append(newArgs, "--")
	newArgs = append(newArgs, remainingArgs...)

	return newArgs, nil
}
