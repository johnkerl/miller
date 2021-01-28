package transformers

import (
	"container/list"
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
const verbNameTac = "tac"

var TacSetup = transforming.TransformerSetup{
	Verb:         verbNameTac,
	ParseCLIFunc: transformerTacParseCLI,
	UsageFunc:    transformerTacUsage,
	IgnoresInput: false,
}

func transformerTacParseCLI(
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
			transformerTacUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerTacUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	transformer, _ := NewTransformerTac()

	*pargi = argi
	return transformer
}

func transformerTacUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameTac)
	fmt.Fprintf(
		o,
		"Prints records in reverse order from the order in which they were encountered.\n",
	)
	if doExit {
		os.Exit(exitCode)
	}
}

// ----------------------------------------------------------------
type TransformerTac struct {
	recordsAndContexts *list.List
}

func NewTransformerTac() (*TransformerTac, error) {
	return &TransformerTac{
		recordsAndContexts: list.New(),
	}, nil
}

func (this *TransformerTac) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.recordsAndContexts.PushFront(inrecAndContext)
	} else {
		// end of stream
		for e := this.recordsAndContexts.Front(); e != nil; e = e.Next() {
			outputChannel <- e.Value.(*types.RecordAndContext)
		}
		outputChannel <- types.NewEndOfStreamMarker(&inrecAndContext.Context)
	}
}
