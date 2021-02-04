// ================================================================
// Just playing around -- nothing serious here.
// ================================================================

package repl

import (
	"fmt"
	"os"
	"path"
	"strings"

	"miller/cliutil"
)

// ================================================================
func ReplUsage(verbName string, o *os.File, exitCode int) {
	fmt.Fprintf(o, "Usage: %s %s with no arguments\n", mlrExeName(), verbName)
	os.Exit(exitCode)
}

func mlrExeName() string {
	return path.Base(os.Args[0])
}

// args are the full Miller command line: "mlr repl foo bar".
func ReplMain(args []string) int {
	exeName := args[0]
	replName := args[1]
	argc := len(args)
	argi := 2

	astPrintMode := ASTPrintNone
	options := cliutil.DefaultOptions()

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process
		}

		if args[argi] == "-h" || args[argi] == "--help" {
			fmt.Println("help stub")
			os.Exit(0)
			// transformerPutUsage(os.Stdout, true, 0)

		} else if args[argi] == "-v" {
			astPrintMode = ASTPrintIndent
			argi++
		} else if args[argi] == "-d" {
			astPrintMode = ASTPrintParex
			argi++
		} else if args[argi] == "-D" {
			astPrintMode = ASTPrintParexOneLine
			argi++

		} else if cliutil.ParseReaderWriterOptions(args, argc, &argi, &options.ReaderOptions, &options.WriterOptions) {

		} else if cliutil.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {

		} else if cliutil.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {

		} else {
			fmt.Println("help stub")
			os.Exit(1)
			// transformerPutUsage(os.Stderr, true, 1)
		}
	}

	repl, err := NewRepl(
		exeName,
		replName,
		astPrintMode,
		&options,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	filenames := args[argi:]
	if len(filenames) > 0 {
		repl.openFiles(filenames)
	}

	repl.handleSession(os.Stdin)
	return 0
}
