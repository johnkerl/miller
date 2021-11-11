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

	"mlr/internal/pkg/cli"
)

// ================================================================
func replUsage(verbName string, o *os.File, exitCode int) {
	exeName := path.Base(os.Args[0])
	fmt.Fprintf(o, "Usage: %s %s [options] {zero or more data-file names}\n", exeName, verbName)

	// TODO: cli/UsageForReaderOptions
	// TODO: cli/UsageForWriterOptions
	// TODO: cli/UsageForReaderWriterOptions

	// TODO: maybe -f/-e as in put?
	// TODO: maybe -s as in put?
	// TODO: maybe -x as in put?
	// TODO: maybe -q as in put?

	fmt.Fprint(o,
		`-v Prints the expressions's AST (abstract syntax tree), which gives
   full transparency on the precedence and associativity rules of
   Miller's grammar, to stdout.

-d Like -v but uses a parenthesized-expression format for the AST.

-D Like -d but with output all on one line.

-w Show warnings about uninitialized variables

-q Don't show startup banner

-s Don't show prompts

--load {DSL script file} Load script file before presenting the prompt.
   If the name following --load is a directory, load all "*.mlr" files
   in that directory.

--mload {DSL script files} -- Like --load but works with more than one filename,
   e.g. '--mload *.mlr --'.

-h|--help Show this message.

Or any --icsv, --ojson, etc. reader/writer options as for the main Miller command line.

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

	showStartupBanner := true
	showPrompts := true
	astPrintMode := ASTPrintNone
	doWarnings := false
	options := cli.DefaultOptions()

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process
		}

		if args[argi] == "-h" || args[argi] == "--help" {
			replUsage(replName, os.Stdout, 0)

		} else if args[argi] == "-q" {
			showStartupBanner = false
			argi++
		} else if args[argi] == "-s" {
			showPrompts = false
			argi++
		} else if args[argi] == "-v" {
			astPrintMode = ASTPrintIndent
			argi++
		} else if args[argi] == "-d" {
			astPrintMode = ASTPrintParex
			argi++
		} else if args[argi] == "-D" {
			astPrintMode = ASTPrintParexOneLine
			argi++
		} else if args[argi] == "-w" {
			doWarnings = true
			argi++

		} else if args[argi] == "--load" {
			if argc-argi < 2 {
				replUsage(replName, os.Stderr, 1)
			}
			options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi+1])
			argi += 2

		} else if args[argi] == "--mload" {
			if argc-argi < 2 {
				replUsage(replName, os.Stderr, 1)
			}
			argi += 1
			for argi < argc && args[argi] != "--" {
				options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi])
				argi += 1
			}
			if args[argi] == "--" {
				argi += 1
			}

		} else if cli.FLAG_TABLE.Parse(args, argc, &argi, &options) {

		} else {
			replUsage(replName, os.Stderr, 1)
		}
	}

	cli.FinalizeReaderOptions(&options.ReaderOptions)
	cli.FinalizeWriterOptions(&options.WriterOptions)

	// --auto-flatten is on by default. But if input and output formats are both JSON,
	// then we don't need to actually do anything. See also mlrcli_parse.go.
	options.WriterOptions.AutoFlatten = cli.DecideFinalFlatten(&options.WriterOptions)
	options.WriterOptions.AutoUnflatten = cli.DecideFinalUnflatten(&options)

	repl, err := NewRepl(
		exeName,
		replName,
		showStartupBanner,
		showPrompts,
		astPrintMode,
		doWarnings,
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
