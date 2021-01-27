package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/output"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameTee = "tee"

var TeeSetup = transforming.TransformerSetup{
	Verb:         verbNameTee,
	ParseCLIFunc: transformerTeeParseCLI,
	IgnoresInput: false,
}

func transformerTeeParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	mainRecordWriterOptions *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	filenameOrCommand := ""
	appending := false
	piping := false
	// TODO: make sure this is a full nested-struct copy.
	var recordWriterOptions *clitypes.TWriterOptions = nil
	if mainRecordWriterOptions != nil {
		copyThereof := *mainRecordWriterOptions
		recordWriterOptions = &copyThereof
	}

	// Parse local flags.

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerTeeUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else if args[argi] == "-a" {
			appending = true
			piping = false
			argi += 1

		} else if args[argi] == "-p" {
			appending = false
			piping = true
			argi += 1

		} else if clitypes.ParseWriterOptions(args, argc, &argi, recordWriterOptions) {
			// This lets mlr main and mlr tee have different output formats.
			// Nothing else to handle here.

		} else {
			transformerTeeUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	// Get the filename/command from the command line, after the flags
	if argi >= argc {
		transformerTeeUsage(os.Stderr, true, 1)
		os.Exit(1)
	}
	filenameOrCommand = args[argi]
	argi += 1

	transformer, err := NewTransformerTee(
		appending,
		piping,
		filenameOrCommand,
		recordWriterOptions,
	)
	if err != nil {
		// Error message already printed out
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

func transformerTeeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {filename}\n", os.Args[0], verbNameTee)
	fmt.Fprintf(o,
		`-a    Append to existing file, if any, rather than overwriting.
-p    Treat filename as a pipe-to command.
Any of the output-format command-line flags (see mlr -h). Example: using
  mlr --icsv --opprint put '...' then tee --ojson ./mytap.dat then stats1 ...
the input is CSV, the output is pretty-print tabular, but the tee-file output
is written in JSON format.

-h|--help Show this message.
`)
}

// ----------------------------------------------------------------
type TransformerTee struct {
	filenameOrCommandForDisplay string
	fileOutputHandler           *output.FileOutputHandler
}

func NewTransformerTee(
	appending bool,
	piping bool,
	filenameOrCommand string,
	recordWriterOptions *clitypes.TWriterOptions,
) (*TransformerTee, error) {
	var fileOutputHandler *output.FileOutputHandler = nil
	var err error = nil
	filenameOrCommandForDisplay := filenameOrCommand
	if piping {
		fileOutputHandler, err = output.NewPipeWriteOutputHandler(filenameOrCommand, recordWriterOptions)
		filenameOrCommandForDisplay = "| " + filenameOrCommand
	} else if appending {
		fileOutputHandler, err = output.NewFileAppendOutputHandler(filenameOrCommand, recordWriterOptions)
		filenameOrCommandForDisplay = ">> " + filenameOrCommand
	} else {
		fileOutputHandler, err = output.NewFileWriteOutputHandler(filenameOrCommand, recordWriterOptions)
		filenameOrCommandForDisplay = "> " + filenameOrCommand
	}
	if err != nil {
		return nil, err
	}

	return &TransformerTee{
		filenameOrCommandForDisplay: filenameOrCommandForDisplay,
		fileOutputHandler:           fileOutputHandler,
	}, nil
}

func (this *TransformerTee) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		err := this.fileOutputHandler.WriteRecordAndContext(inrecAndContext)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: error writing to tee \"%s\":\n",
				os.Args[0], this.filenameOrCommandForDisplay,
			)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputChannel <- inrecAndContext
	} else {
		err := this.fileOutputHandler.Close()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: error closing tee \"%s\":\n",
				os.Args[0], this.filenameOrCommandForDisplay,
			)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		outputChannel <- inrecAndContext
	}
}
