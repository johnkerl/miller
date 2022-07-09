package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
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
) IRecordTransformer {

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
			os.Exit(0)

		} else {
			transformerGroupByUsage(os.Stderr)
			os.Exit(1)
		}
	}

	// Get the group-by field names from the command line
	if argi >= argc {
		transformerGroupByUsage(os.Stderr)
		os.Exit(1)
	}
	groupByFieldNames := lib.SplitString(args[argi], ",")
	argi++

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerGroupBy(
		groupByFieldNames,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerGroupBy struct {
	// input
	groupByFieldNames []string

	// state
	// map from string to *list.List
	recordListsByGroup *lib.OrderedMap
}

func NewTransformerGroupBy(
	groupByFieldNames []string,
) (*TransformerGroupBy, error) {

	tr := &TransformerGroupBy{
		groupByFieldNames: groupByFieldNames,

		recordListsByGroup: lib.NewOrderedMap(),
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerGroupBy) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
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
			recordListForGroup = list.New()
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)

	} else {
		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputRecordsAndContexts.PushBack(inner.Value.(*types.RecordAndContext))
			}
		}
		outputRecordsAndContexts.PushBack(inrecAndContext) // end-of-stream marker
	}
}
