// Entry point for mlr script: executes a Miller DSL script with next()-driven
// record iteration. No prompt, no REPL. Script from -f or -e, data from
// filenames or stdin.
//
// Example:
//
//	mlr script -f sum.mlr data.csv
//	mlr script -e '@sum=0; while(next()){@sum+=$x}; print @sum' --csv data.csv

package script

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
)

func scriptUsage(verbName string, o *os.File, exitCode int) {
	exeName := path.Base(os.Args[0])
	fmt.Fprintf(o, "Usage: %s %s [options] -f {script file} | -e {expression} [zero or more data-file names]\n", exeName, verbName)
	fmt.Fprint(o,
		`-f {file}       Load DSL script from file. If file is a directory, load all *.mlr files.
-e {expression}  DSL script from command line. May be combined with -f.

--load {file}    Preload DSL before script (e.g. function library).
--mload {files}  Like --load but multiple files. Use -- to terminate list.

-w              Show warnings about uninitialized variables.
-z              Strict mode.

-h|--help       Show this message.

Or any --icsv, --ojson, etc. reader/writer options.

Data-file names (or stdin if none) are the record input stream.
`)
	os.Exit(exitCode)
}

func ScriptMain(args []string) int {
	scriptName := args[0]
	argc := len(args)
	argi := 1

	doWarnings := false
	strictMode := false
	options := cli.DefaultOptions()
	var dslStrings []string
	haveDSLStrings := false

	for argi < argc {
		if !strings.HasPrefix(args[argi], "-") {
			break
		}

		if args[argi] == "-h" || args[argi] == "--help" {
			scriptUsage(scriptName, os.Stdout, 0)
		} else if args[argi] == "-w" {
			doWarnings = true
			argi++
		} else if args[argi] == "-z" {
			strictMode = true
			argi++
		} else if args[argi] == "-f" {
			if argc-argi < 2 {
				scriptUsage(scriptName, os.Stderr, 1)
			}
			argi++
			filename, err := cli.VerbGetStringArg("script", "-f", args, &argi, argc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
			theseStrings, err := lib.LoadStringsFromFileOrDir(filename, ".mlr")
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr script: cannot load script from \"%s\": %v\n", filename, err)
				os.Exit(1)
			}
			dslStrings = append(dslStrings, theseStrings...)
			haveDSLStrings = true
		} else if args[argi] == "-e" {
			if argc-argi < 2 {
				scriptUsage(scriptName, os.Stderr, 1)
			}
			argi++
			expr, err := cli.VerbGetStringArg("script", "-e", args, &argi, argc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
				os.Exit(1)
			}
			dslStrings = append(dslStrings, expr)
			haveDSLStrings = true
		} else if args[argi] == "--load" {
			if argc-argi < 2 {
				scriptUsage(scriptName, os.Stderr, 1)
			}
			options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi+1])
			argi += 2
		} else if args[argi] == "--mload" {
			if argc-argi < 2 {
				scriptUsage(scriptName, os.Stderr, 1)
			}
			argi++
			for argi < argc && args[argi] != "--" {
				options.DSLPreloadFileNames = append(options.DSLPreloadFileNames, args[argi])
				argi++
			}
			if argi < argc && args[argi] == "--" {
				argi++
			}
		} else if cli.FLAG_TABLE.Parse(args, argc, &argi, options) {
		} else {
			scriptUsage(scriptName, os.Stderr, 1)
		}
	}

	if !haveDSLStrings {
		fmt.Fprintf(os.Stderr, "mlr script: -f or -e is required\n")
		scriptUsage(scriptName, os.Stderr, 1)
	}

	cli.FinalizeReaderOptions(&options.ReaderOptions)
	cli.FinalizeWriterOptions(&options.WriterOptions)
	options.WriterOptions.AutoFlatten = cli.DecideFinalFlatten(&options.WriterOptions)
	options.WriterOptions.AutoUnflatten = cli.DecideFinalUnflatten(options, [][]string{})

	filenames := args[argi:]

	scr, err := NewScript(options, doWarnings, strictMode, dslStrings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr script: %v\n", err)
		os.Exit(1)
	}

	scr.openFiles(filenames)
	err = scr.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr script: %v\n", err)
		os.Exit(1)
	}

	scr.bufferedRecordOutputStream.Flush()
	err = scr.closeBufferedOutputStream()
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr script: %v\n", err)
		os.Exit(1)
	}
	return 0
}
