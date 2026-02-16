package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameSortWithinRecords = "sort-within-records"

var SortWithinRecordsSetup = TransformerSetup{
	Verb:         verbNameSortWithinRecords,
	UsageFunc:    transformerSortWithinRecordsUsage,
	ParseCLIFunc: transformerSortWithinRecordsParseCLI,
	IgnoresInput: false,
}

func transformerSortWithinRecordsUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameSortWithinRecords)
	fmt.Fprintln(o, "Outputs records sorted lexically ascending by keys.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-r        Recursively sort subobjects/submaps, e.g. for JSON input.\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerSortWithinRecordsParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++
	doRecurse := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		if args[argi] == "--" {
			break // All transformers must do this so main-flags can follow verb-flags
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerSortWithinRecordsUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		} else if opt == "-r" {
			doRecurse = true

		} else {
			return nil, cli.VerbErrorf(verbNameSortWithinRecords, "option \"%s\" not recognized", opt)
		}
	}

	// TODO: allow sort by key or value?
	// TODO: allow sort ascendending/descending?

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerSortWithinRecords(doRecurse)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

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

func (tr *TransformerSortWithinRecords) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, outputRecordsAndContexts, inputDownstreamDoneChannel, outputDownstreamDoneChannel)
}

func (tr *TransformerSortWithinRecords) transformNonrecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKey()
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}

func (tr *TransformerSortWithinRecords) transformRecursively(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		inrec.SortByKeyRecursively()
	}
	*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // including end-of-stream marker
}
