package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameGroupBy = "group-by"

var GroupBySetup = TransformerSetup{
	Verb:         verbNameGroupBy,
	UsageFunc:    transformerGroupByUsage,
	ParseCLIFunc: transformerGroupByParseCLI,
	IgnoresInput: false,
}

func transformerGroupByUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options] {comma-separated field names}\n", "mlr", verbNameGroupBy)
	fmt.Fprint(o, "Outputs records in batches having identical values at specified field names.")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
}

func transformerGroupByParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	argi++

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
			transformerGroupByUsage(os.Stdout)
			return nil, cli.ErrHelpRequested

		}
		return nil, cli.VerbErrorf(verbNameGroupBy, "option \"%s\" not recognized", opt)
	}

	// Get the group-by field names from the command line
	if argi >= argc {
		return nil, cli.VerbErrorf(verbNameGroupBy, "group-by field names required")
	}
	groupByFieldNames := lib.SplitString(args[argi], ",")
	argi++

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerGroupBy(
		groupByFieldNames,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerGroupBy struct {
	// input
	groupByFieldNames []string

	// state
	// map from string to record slices
	recordListsByGroup *lib.OrderedMap[*[]*types.RecordAndContext]
}

func NewTransformerGroupBy(
	groupByFieldNames []string,
) (*TransformerGroupBy, error) {

	tr := &TransformerGroupBy{
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap[*[]*types.RecordAndContext](),
	}

	return tr, nil
}

func (tr *TransformerGroupBy) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok {
			return
		}

		recordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil {
			records := []*types.RecordAndContext{}
			recordListForGroup = &records
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		*recordListForGroup = append(*recordListForGroup, inrecAndContext)

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, *outer.Value...)
		}
		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // end-of-stream marker
	}
}
