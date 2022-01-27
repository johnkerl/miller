package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/internal/pkg/cli"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/mlrval"
	"github.com/johnkerl/miller/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCountSimilar = "count-similar"

var CountSimilarSetup = TransformerSetup{
	Verb:         verbNameCountSimilar,
	UsageFunc:    transformerCountSimilarUsage,
	ParseCLIFunc: transformerCountSimilarParseCLI,
	IgnoresInput: false,
}

func transformerCountSimilarUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCountSimilar)
	fmt.Fprintf(o, "Ingests all records, then emits each record augmented by a count of\n")
	fmt.Fprintf(o, "the number of other records having the same group-by field values.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Group-by-field names for counts, e.g. a,b,c\n")
	fmt.Fprintf(o, "-o {name} Field name for output-counts. Defaults to \"count\".\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCountSimilarParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var groupByFieldNames []string = nil
	counterFieldName := "count"

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
			transformerCountSimilarUsage(os.Stdout, true, 0)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			counterFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCountSimilarUsage(os.Stderr, true, 1)
		}
	}

	if groupByFieldNames == nil {
		transformerCountSimilarUsage(os.Stderr, true, 1)
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil
	}

	transformer, err := NewTransformerCountSimilar(
		groupByFieldNames,
		counterFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return transformer
}

// ----------------------------------------------------------------
type TransformerCountSimilar struct {
	// Input:
	groupByFieldNames []string
	counterFieldName  string

	// State:
	recordListsByGroup *lib.OrderedMap // map from string to *list.List
}

// ----------------------------------------------------------------
func NewTransformerCountSimilar(
	groupByFieldNames []string,
	counterFieldName string,
) (*TransformerCountSimilar, error) {
	tr := &TransformerCountSimilar{
		groupByFieldNames:  groupByFieldNames,
		counterFieldName:   counterFieldName,
		recordListsByGroup: lib.NewOrderedMap(),
	}
	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCountSimilar) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *list.List, // list of *types.RecordAndContext
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.groupByFieldNames)
		if !ok { // This particular record doesn't have the specified fields; ignore
			return
		}

		irecordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if irecordListForGroup == nil { // first time
			irecordListForGroup = list.New()
			tr.recordListsByGroup.Put(groupingKey, irecordListForGroup)
		}
		recordListForGroup := irecordListForGroup.(*list.List)

		recordListForGroup.PushBack(inrecAndContext)
	} else {

		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			// TODO: make 64-bit friendly
			groupSize := recordListForGroup.Len()
			mgroupSize := mlrval.FromInt(int64(groupSize))
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				recordAndContext := inner.Value.(*types.RecordAndContext)
				recordAndContext.Record.PutCopy(tr.counterFieldName, mgroupSize)

				outputRecordsAndContexts.PushBack(recordAndContext)
			}
		}

		outputRecordsAndContexts.PushBack(inrecAndContext) // Emit the stream-terminating null record
	}
}
