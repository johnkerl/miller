package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/output"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameTee = "tee"

var TeeSetup = TransformerSetup{
	Verb:         verbNameTee,
	UsageFunc:    transformerTeeUsage,
	ParseCLIFunc: transformerTeeParseCLI,
	IgnoresInput: false,
}

func transformerTeeUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {filename}\n", lib.MlrExeName(), verbNameTee)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o,
		`-a    Append to existing file, if any, rather than overwriting.
-p    Treat filename as a pipe-to command.
Any of the output-format command-line flags (see mlr -h). Example: using
  mlr --icsv --opprint put '...' then tee --ojson ./mytap.dat then stats1 ...
the input is CSV, the output is pretty-print tabular, but the tee-file output
is written in JSON format.

-h|--help Show this message.
`)
	if doExit {
		os.Exit(exitCode)
	}
}

func transformerTeeParseCLI(
	pargi *int,
	argc int,
	args []string,
	mainOptions *cliutil.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	filenameOrCommand := ""
	appending := false
	piping := false
	// TODO: make sure this is a full nested-struct copy.
	var recordWriterOptions *cliutil.TWriterOptions = nil
	if mainOptions != nil {
		copyThereof := mainOptions.WriterOptions // struct copy
		recordWriterOptions = &copyThereof
	}

	// Parse local flags.

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerTeeUsage(os.Stdout, true, 0)

		} else if opt == "-a" {
			appending = true
			piping = false

		} else if opt == "-p" {
			appending = false
			piping = true

		} else {
			// This is inelegant. For error-proofing we advance argi already in our
			// loop (so individual if-statements don't need to). However,
			// ParseWriterOptions expects it unadvanced.
			wargi := argi - 1
			if cliutil.ParseWriterOptions(args, argc, &wargi, recordWriterOptions) {
				// This lets mlr main and mlr tee have different output formats.
				// Nothing else to handle here.
				argi = wargi
			} else {
				transformerTeeUsage(os.Stderr, true, 1)
			}
		}
	}

	// Get the filename/command from the command line, after the flags
	if argi >= argc {
		transformerTeeUsage(os.Stderr, true, 1)
	}
	filenameOrCommand = args[argi]
	argi++

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

// ----------------------------------------------------------------
type TransformerTee struct {
	filenameOrCommandForDisplay string
	fileOutputHandler           *output.FileOutputHandler
}

func NewTransformerTee(
	appending bool,
	piping bool,
	filenameOrCommand string,
	recordWriterOptions *cliutil.TWriterOptions,
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

func (tr *TransformerTee) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		err := tr.fileOutputHandler.WriteRecordAndContext(inrecAndContext)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: error writing to tee \"%s\":\n",
				lib.MlrExeName(), tr.filenameOrCommandForDisplay,
			)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		outputChannel <- inrecAndContext
	} else {
		err := tr.fileOutputHandler.Close()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: error closing tee \"%s\":\n",
				lib.MlrExeName(), tr.filenameOrCommandForDisplay,
			)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		outputChannel <- inrecAndContext
	}
}
