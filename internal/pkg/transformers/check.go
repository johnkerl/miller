package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCheck = "check"

var CheckSetup = TransformerSetup{
	Verb:         verbNameCheck,
	UsageFunc:    transformerCheckUsage,
	ParseCLIFunc: transformerCheckParseCLI,
	IgnoresInput: false,
}

func transformerCheckUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCheck)
	fmt.Fprintf(o, "Consumes records without printing any output.\n")
	fmt.Fprintf(o, "Useful for doing a well-formatted check on input data.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCheckParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCheckUsage(os.Stdout, true, 0)

		} else {
			transformerCheckUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerCheck()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerCheck struct {
	// stateless
}

func NewTransformerCheck() (*TransformerCheck, error) {
	return &TransformerCheck{}, nil
}

func (tr *TransformerCheck) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext
	}
}
