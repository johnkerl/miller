package transformers

import (
	"fmt"
	"os"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameCountSimilar = "count-similar"

var CountSimilarSetup = TransformerSetup{
	Verb:         verbNameCountSimilar,
	UsageFunc:    transformerCountSimilarUsage,
	ParseCLIFunc: transformerCountSimilarParseCLI,
	IgnoresInput: false,
}

func transformerCountSimilarUsage(
	o *os.File,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCountSimilar)
	fmt.Fprintf(o, "Ingests all records, then emits each record augmented by a count of\n")
	fmt.Fprintf(o, "the number of other records having the same group-by field values.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Group-by-field names for counts, e.g. a,b,c\n")
	fmt.Fprintf(o, "-o {name} Field name for output-counts. Defaults to \"count\".\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")
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
			transformerCountSimilarUsage(os.Stdout)
			os.Exit(0)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-o" {
			counterFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCountSimilarUsage(os.Stderr)
			os.Exit(1)
		}
	}

	if groupByFieldNames == nil {
		transformerCountSimilarUsage(os.Stderr)
		os.Exit(1)
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

type TransformerCountSimilar struct {
	// Input:
	groupByFieldNames []string
	counterFieldName  string

	// State:
	recordListsByGroup *lib.OrderedMap[*[]*types.RecordAndContext] // map from string to records
}

func NewTransformerCountSimilar(
	groupByFieldNames []string,
	counterFieldName string,
) (*TransformerCountSimilar, error) {
	tr := &TransformerCountSimilar{
		groupByFieldNames:  groupByFieldNames,
		counterFieldName:   counterFieldName,
		recordListsByGroup: lib.NewOrderedMap[*[]*types.RecordAndContext](),
	}
	return tr, nil
}


func (tr *TransformerCountSimilar) Transform(
	inrecAndContext *types.RecordAndContext,
	outputRecordsAndContexts *[]*types.RecordAndContext, // list of *types.RecordAndContext
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

		recordListForGroup := tr.recordListsByGroup.Get(groupingKey)
		if recordListForGroup == nil { // first time
			records := make([]*types.RecordAndContext, 0)
			recordListForGroup = &records
			tr.recordListsByGroup.Put(groupingKey, recordListForGroup)
		}

		*recordListForGroup = append(*recordListForGroup, inrecAndContext)
	} else {

		for outer := tr.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value
			// TODO: make 64-bit friendly
			groupSize := len(*recordListForGroup)
			mgroupSize := mlrval.FromInt(int64(groupSize))
			for _, recordAndContext := range *recordListForGroup {
				recordAndContext.Record.PutCopy(tr.counterFieldName, mgroupSize)

				*outputRecordsAndContexts = append(*outputRecordsAndContexts, recordAndContext)
			}
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // Emit the stream-terminating null record
	}
}
