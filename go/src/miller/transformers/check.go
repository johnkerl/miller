package transformers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
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
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		if !strings.HasPrefix(args[argi], "-") {
			break // No more flag options to process

		} else if args[argi] == "-h" || args[argi] == "--help" {
			transformerCheckUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerCheckUsage(os.Stderr, true, 1)
			os.Exit(1)
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
	fmt.Fprintf(o, "Usage: %s %s, with no options\n", os.Args[0], verbNameCheck)
	fmt.Fprintf(o, "Consumes records without printing any output.\n")
	fmt.Fprintf(o, "Useful for doing a well-formatted check on input data.\n")

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
