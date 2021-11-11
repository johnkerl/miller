package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameNothing = "nothing"

var NothingSetup = TransformerSetup{
	Verb:         verbNameNothing,
	ParseCLIFunc: transformerNothingParseCLI,
	UsageFunc:    transformerNothingUsage,
	IgnoresInput: false,
}

func transformerNothingUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameNothing)
	fmt.Fprintf(o, "Drops all input records. Useful for testing, or after tee/print/etc. have\n")
	fmt.Fprintf(o, "produced other output.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerNothingParseCLI(
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
			transformerNothingUsage(os.Stdout, true, 0)

		} else {
			transformerNothingUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerNothing()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerNothing struct {
	// stateless
}

func NewTransformerNothing() (*TransformerNothing, error) {
	return &TransformerNothing{}, nil
}

func (tr *TransformerNothing) Transform(
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
