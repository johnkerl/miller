package transformers

import (
	"container/list"
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameCountSimilar = "count-similar"

var CountSimilarSetup = transforming.TransformerSetup{
	Verb:         verbNameCountSimilar,
	ParseCLIFunc: transformerCountSimilarParseCLI,
	UsageFunc:    transformerCountSimilarUsage,
	IgnoresInput: false,
}

func transformerCountSimilarParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TReaderOptions,
	__ *cliutil.TWriterOptions,
) transforming.IRecordTransformer {

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
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCountSimilarUsage(os.Stdout, true, 0)

		} else if opt == "-g" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			counterFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCountSimilarUsage(os.Stderr, true, 1)
		}
	}

	if groupByFieldNames == nil {
		transformerCountSimilarUsage(os.Stderr, true, 1)
	}

	transformer, _ := NewTransformerCountSimilar(
		groupByFieldNames,
		counterFieldName,
	)

	*pargi = argi
	return transformer
}

func transformerCountSimilarUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCountSimilar)
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
	this := &TransformerCountSimilar{
		groupByFieldNames:  groupByFieldNames,
		counterFieldName:   counterFieldName,
		recordListsByGroup: lib.NewOrderedMap(),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerCountSimilar) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok { // This particular record doesn't have the specified fields; ignore
			return
		}

		irecordListForGroup := this.recordListsByGroup.Get(groupingKey)
		if irecordListForGroup == nil { // first time
			irecordListForGroup = list.New()
			this.recordListsByGroup.Put(groupingKey, irecordListForGroup)
		}
		recordListForGroup := irecordListForGroup.(*list.List)

		recordListForGroup.PushBack(inrecAndContext)
	} else {

		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			// TODO: make 64-bit friendly
			groupSize := recordListForGroup.Len()
			mgroupSize := types.MlrvalFromInt(int(groupSize))
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				recordAndContext := inner.Value.(*types.RecordAndContext)
				recordAndContext.Record.PutCopy(this.counterFieldName, &mgroupSize)

				outputChannel <- recordAndContext
			}
		}

		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
