package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameCheck = "check"

var CheckSetup = transforming.TransformerSetup{
	Verb:         verbNameCheck,
	ParseCLIFunc: transformerCheckParseCLI,
	UsageFunc:    transformerCheckUsage,
	IgnoresInput: false,
}

func transformerCheckParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
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
			transformerCheckUsage(os.Stdout, true, 0)

		} else {
			transformerCheckUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerCheck()

	*pargi = argi
	return transformer
}

func transformerCheckUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCheck)
	fmt.Fprintf(o, "Consumes records without printing any output.\n")
	fmt.Fprintf(o, "Useful for doing a well-formatted check on input data.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerCheck struct {
	// stateless
}

func NewTransformerCheck() (*TransformerCheck, error) {
	return &TransformerCheck{}, nil
}

func (this *TransformerCheck) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if inrecAndContext.EndOfStream {
		outputChannel <- inrecAndContext
	}
}
