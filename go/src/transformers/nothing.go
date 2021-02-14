package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameNothing = "nothing"

var NothingSetup = transforming.TransformerSetup{
	Verb:         verbNameNothing,
	ParseCLIFunc: transformerNothingParseCLI,
	UsageFunc:    transformerNothingUsage,
	IgnoresInput: false,
}

func transformerNothingParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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

	transformer, _ := NewTransformerNothing()

	*pargi = argi
	return transformer
}

func transformerNothingUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameNothing)
	fmt.Fprintf(o, "Drops all input records. Useful for testing, or after tee/print/etc. have\n")
	fmt.Fprintf(o, "produced other output.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerNothing struct {
	// stateless
}

func NewTransformerNothing() (*TransformerNothing, error) {
	return &TransformerNothing{}, nil
}

func (this *TransformerNothing) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext
	}
}
