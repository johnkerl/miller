package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameSortWithinRecords = "sort-within-records"

var SortWithinRecordsSetup = TransformerSetup{
	Verb:         verbNameSortWithinRecords,
	UsageFunc:    transformerSortWithinRecordsUsage,
	ParseCLIFunc: transformerSortWithinRecordsParseCLI,
	IgnoresInput: false,
}

func transformerSortWithinRecordsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSortWithinRecords)
	fmt.Fprintln(o, "Outputs records sorted lexically ascending by keys.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-r        Recursively sort subobjects/submaps, e.g. for JSON input.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSortWithinRecordsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++
	doRecurse := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSortWithinRecordsUsage(os.Stdout, true, 0)

		} else if opt == "-r" {
			doRecurse = true

		} else {
			transformerSortWithinRecordsUsage(os.Stderr, true, 1)
		}
	}

	// TODO: allow sort by key or value?
	// TODO: allow sort ascendending/descending?

	transformer, err := NewTransformerSortWithinRecords(doRecurse)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerSortWithinRecords struct {
	recordTransformerFunc RecordTransformerFunc
}

func NewTransformerSortWithinRecords(
	doRecurse bool,
) (*TransformerSortWithinRecords, error) {

	tr := &TransformerSortWithinRecords{}
	if doRecurse {
		tr.recordTransformerFunc = tr.transformRecursively
	} else {
		tr.recordTransformerFunc = tr.transformNonrecursively
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerSortWithinRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformNonrecursively(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKey()
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}

// ----------------------------------------------------------------
func (tr *TransformerSortWithinRecords) transformRecursively(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKeyRecursively()
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}
