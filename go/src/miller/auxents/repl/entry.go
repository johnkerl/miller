// ================================================================
// This is the shell command-line entry point to the Miller REPL command line.
// E.g. at the shell prompt, you type 'mlr repl --json' -- this file will parse
// that.  It will then hand off control to a REPL session which will handle all
// subsequent REPL command-line statements you type in at a Miller REPL prompt.
//
// Example:
//
// bash$ mlr repl --json         <------------- this file handles this bit
// [mlr] :open myfile.json
// [mlr] :read
// [mlr] :context
// FILENAME="myfile.json",FILENUM=1,NR=1,FNR=1
// [mlr] $*
// {
//   "hostname": "localhost",
//   "pid": 12345
// }
// [mlr] :quit
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
	exeName := path.Base(os.Args[0])
	fmt.Fprintf(o, "Usage: %s %s [options] {zero or more data-file names}\n", exeName, verbName)

	// TODO: cliutil/UsageForReaderOptions
	// TODO: cliutil/UsageForWriterOptions
	// TODO: cliutil/UsageForReaderWriterOptions

	// TODO: maybe -f/-e as in put?
	// TODO: maybe -s as in put?
	// TODO: maybe -x as in put?
	// TODO: maybe -q as in put?

	fmt.Fprint(o,
		`-v Prints the expressions's AST (abstract syntax tree), which gives
full transparency on the precedence and associativity rules of
Miller's grammar, to stdout.

-d Like -v but uses a parenthesized-expression format for the AST. Then, exits without
   stream processing.

-D Like -d but with output all on one line.

or, ny --icsv, --ojson, etc. reader/writer options as for the main Miller command line.

-h|--help Show this message.

Any data-file names are opened just as if you had waited and typed :open {filenames}
at the Miller REPL prompt.
`)

	os.Exit(exitCode)
}

// Here the args are the full Miller command line: "mlr repl foo bar".
func ReplMain(args []string) int {
	exeName := path.Base(args[0])
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
			ReplUsage(replName, os.Stdout, 0)

		} else if args[argi] == "-v" {
			astPrintMode = ASTPrintIndent
			argi++
		} else if args[argi] == "-d" {
			astPrintMode = ASTPrintParex
			argi++
		} else if args[argi] == "-D" {
			astPrintMode = ASTPrintParexOneLine
			argi++

		} else if cliutil.ParseReaderWriterOptions(
			args, argc, &argi, &options.ReaderOptions, &options.WriterOptions,
		) {

		} else if cliutil.ParseReaderOptions(args, argc, &argi, &options.ReaderOptions) {

		} else if cliutil.ParseWriterOptions(args, argc, &argi, &options.WriterOptions) {

		} else {
			ReplUsage(replName, os.Stderr, 1)
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
