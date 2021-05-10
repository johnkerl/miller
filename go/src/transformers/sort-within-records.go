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
const verbNameSortWithinRecords = "sort-within-records"

var SortWithinRecordsSetup = transforming.TransformerSetup{
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
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameSortWithinRecords)
	fmt.Fprintln(o, "Outputs records sorted lexically ascending by keys.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerSortWithinRecordsParseCLI(
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
			transformerSortWithinRecordsUsage(os.Stdout, true, 0)

		} else {
			transformerSortWithinRecordsUsage(os.Stderr, true, 1)
		}
	}

	// TODO: allow sort by key or value?
	// TODO: allow sort ascendending/descending?

	transformer, err := NewTransformerSortWithinRecords()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerSortWithinRecords struct {
}

func NewTransformerSortWithinRecords() (*TransformerSortWithinRecords, error) {

	this := &TransformerSortWithinRecords{}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerSortWithinRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKey()
	}
	outputChannel <- inrecAndContext // including end-of-stream marker
}
