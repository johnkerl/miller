package entrypoint

import (
	"fmt"
	"os"

	"miller/src/platform"

	"miller/src/auxents"
	"miller/src/cli"
	"miller/src/lib"
	"miller/src/stream"
)

// ----------------------------------------------------------------
func Main() {
	// Special handling for Windows so we can do things like:
	//
	//   mlr put '$a = $b . "cd \"efg\" hi"' foo.dat
	//
	// as on Linux/Unix/MacOS.
	os.Args = platform.GetArgs()

	// Expand "-xyz" into "-x -y -z" while leaving "--xyz" intact. This is a
	// keystroke-saver for the user.
	//
	// This is OK to do globally here since Miller is quite consistent (in
	// main, verbs, and auxents) that multi-character options start with two
	// dashes, e.g. "--csv". (The sole exception is the sort verb's -nf/-nr
	// which are handled specially there.)
	os.Args = lib.Getoptify(os.Args)

	// 'mlr repl' or 'mlr lecat' or any other non-miller-per-se toolery which
	// is delivered (for convenience) within the mlr executable. If argv[1] is
	// found then this function will not return.
	auxents.Dispatch(os.Args)

	options, recordTransformers, err := cli.ParseCommandLine(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}

	err = stream.Stream(options, recordTransformers)
	if err != nil {
		fmt.Fprintln(os.Stderr, os.Args[0], ": ", err)
		os.Exit(1)
	}
}
