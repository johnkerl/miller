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
const verbNameSortWithinRecords = "sort-within-records"

var SortWithinRecordsSetup = transforming.TransformerSetup{
	Verb:         verbNameSortWithinRecords,
	ParseCLIFunc: transformerSortWithinRecordsParseCLI,
	UsageFunc:    transformerSortWithinRecordsUsage,
	IgnoresInput: false,
}

func transformerSortWithinRecordsParseCLI(
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
			transformerSortWithinRecordsUsage(os.Stdout, true, 0)
			return nil // help intentionally requested

		} else {
			transformerSortWithinRecordsUsage(os.Stderr, true, 1)
			os.Exit(1)
		}
	}

	// TODO: allow sort by key or value?
	// TODO: allow sort ascendending/descending?

	transformer, _ := NewTransformerSortWithinRecords()

	*pargi = argi
	return transformer
}

func transformerSortWithinRecordsUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s, with no options.\n", os.Args[0], verbNameSortWithinRecords)
	fmt.Fprintln(o, "Outputs records sorted lexically ascending by keys.")

	if doExit {
		os.Exit(exitCode)
	}
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
