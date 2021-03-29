package entrypoint

import (
	"fmt"
	"os"

	"miller/src/auxents"
	"miller/src/cli"
	"miller/src/stream"
)

// ----------------------------------------------------------------
func Main() {

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
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
