package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/output"
	"mlr/internal/pkg/types"
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
	fmt.Fprintf(o, "Usage: %s %s [options] {filename}\n", "mlr", verbNameTee)
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
	mainOptions *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	filenameOrCommand := ""
	appending := false
	piping := false
	// TODO: make sure this is a full nested-struct copy.
	var localOptions *cli.TOptions = nil
	if mainOptions != nil {
		copyThereof := *mainOptions // struct copy
		localOptions = &copyThereof
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
			largi := argi - 1
			if cli.FLAG_TABLE.Parse(args, argc, &largi, localOptions) {
				// This lets mlr main and mlr tee have different output formats.
				// Nothing else to handle here.
				argi = largi
			} else {
				transformerTeeUsage(os.Stderr, true, 1)
			}
		}
	}

	cli.FinalizeWriterOptions(&localOptions.WriterOptions)

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
		&localOptions.WriterOptions,
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
	recordWriterOptions *cli.TWriterOptions,
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
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {

	// If we receive a downstream-done flag from a transformer downstream from
	// us, read it to unblock their goroutine but -- unlike most other verbs --
	// do not forward the flag farther upstream.
	//
	// For example, 'mlr cut -f foo then head -n 10' on million-line input:
	// head can signal it's got 10 records, then write downStreamDone <- true,
	// then cut and record-reader can stop sending any more data. This makes
	// the UX response for head on huge files.
	//
	// But 'mlr cut -f foo then tee bar.txt then head -n 10' -- one does expect
	// bar.txt to have all the output from cut.
	select {
	case _ = <-inputDownstreamDoneChannel:
		// Do not write this to the coutputDownstreamDoneChannel, as other transformers do
		break
	default:
		break
	}

	if !inrecAndContext.EndOfStream {
		err := tr.fileOutputHandler.WriteRecordAndContext(inrecAndContext)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"%s: error writing to tee \"%s\":\n",
				"mlr", tr.filenameOrCommandForDisplay,
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
				"mlr", tr.filenameOrCommandForDisplay,
			)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		outputChannel <- inrecAndContext
	}
}
