package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/cliutil"
	"miller/lib"
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
			transformerTacUsage(os.Stdout, true, 0)

		} else {
			transformerTacUsage(os.Stderr, true, 1)
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameTac)
	fmt.Fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

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
