package transformers

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/lib"
	"github.com/johnkerl/miller/v6/pkg/mlrval"
	"github.com/johnkerl/miller/v6/pkg/types"
)

const verbNameMostFrequent = "most-frequent"
const verbNameLeastFrequent = "least-frequent"

const mostLeastFrequentDefaultMaxOutputLength = int64(10)
const mostLeastFrequentDefaultOutputFieldName = "count"

var MostFrequentSetup = TransformerSetup{
	Verb:         verbNameMostFrequent,
	UsageFunc:    transformerMostFrequentUsage,
	ParseCLIFunc: transformerMostFrequentParseCLI,
	IgnoresInput: false,
}

var LeastFrequentSetup = TransformerSetup{
	Verb:         verbNameLeastFrequent,
	UsageFunc:    transformerLeastFrequentUsage,
	ParseCLIFunc: transformerLeastFrequentParseCLI,
	IgnoresInput: false,
}

func transformerMostFrequentUsage(
	o *os.File,
) {
	argv0 := "mlr"
	verb := verbNameMostFrequent
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Shows the most frequently occurring distinct values for specified field names.\n")
	fmt.Fprintf(o, "The first entry is the statistical mode; the remaining are runners-up.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {one or more comma-separated field names}. Required flag.\n")
	fmt.Fprintf(o, "-n {count}. Optional flag defaulting to %d.\n", mostLeastFrequentDefaultMaxOutputLength)
	fmt.Fprintf(o, "-b          Suppress counts; show only field values.\n")
	fmt.Fprintf(o, "-o {name}   Field name for output count. Default \"%s\".\n", mostLeastFrequentDefaultOutputFieldName)
	fmt.Fprintf(o, "See also \"%s %s\".\n", argv0, "least-frequent")
}

func transformerLeastFrequentUsage(
	o *os.File,
) {
	argv0 := "mlr"
	verb := verbNameLeastFrequent
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Shows the least frequently occurring distinct values for specified field names.\n")
	fmt.Fprintf(o, "The first entry is the statistical anti-mode; the remaining are runners-up.\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {one or more comma-separated field names}. Required flag.\n")
	fmt.Fprintf(o, "-n {count}. Optional flag defaulting to %d.\n", mostLeastFrequentDefaultMaxOutputLength)
	fmt.Fprintf(o, "-b          Suppress counts; show only field values.\n")
	fmt.Fprintf(o, "-o {name}   Field name for output count. Default \"%s\".\n", mostLeastFrequentDefaultOutputFieldName)
	fmt.Fprintf(o, "See also \"%s %s\".\n", argv0, "most-frequent")
}

func transformerMostFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {
	return transformerMostOrLeastFrequentParseCLI(pargi, argc, args, true, transformerMostFrequentUsage, doConstruct)
}

func transformerLeastFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {
	return transformerMostOrLeastFrequentParseCLI(pargi, argc, args, false, transformerLeastFrequentUsage, doConstruct)
}

func transformerMostOrLeastFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	descending bool,
	usageFunc TransformerUsageFunc,
	doConstruct bool, // false for first pass of CLI-parse, true for second pass
) (RecordTransformer, error) {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var groupByFieldNames []string = nil
	maxOutputLength := mostLeastFrequentDefaultMaxOutputLength
	showCounts := true
	outputFieldName := mostLeastFrequentDefaultOutputFieldName

	var err error
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
			usageFunc(os.Stdout)
			return nil, cli.ErrHelpRequested

		} else if opt == "-f" {
			groupByFieldNames, err = cli.VerbGetStringArrayArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else if opt == "-n" {
			maxOutputLength, err = cli.VerbGetIntArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else if opt == "-b" {
			showCounts = false

		} else if opt == "-o" {
			outputFieldName, err = cli.VerbGetStringArg(verb, opt, args, &argi, argc)
			if err != nil {
				return nil, err
			}

		} else {
			return nil, cli.VerbErrorf(verb, "option \"%s\" not recognized", opt)
		}
	}

	if groupByFieldNames == nil {
		return nil, cli.VerbErrorf(verb, "-f field names required")
	}

	*pargi = argi
	if !doConstruct { // All transformers must do this for main command-line parsing
		return nil, nil
	}

	transformer, err := NewTransformerMostOrLeastFrequent(
		groupByFieldNames,
		maxOutputLength,
		showCounts,
		outputFieldName,
		descending,
	)
	if err != nil {
		return nil, err
	}

	return transformer, nil
}

type TransformerMostOrLeastFrequent struct {
	groupByFieldNames []string
	maxOutputLength   int64
	showCounts        bool
	outputFieldName   string
	descending        bool
	countsByGroup     *lib.OrderedMap[int64] // map[string]int
	valuesForGroup    map[string][]*mlrval.Mlrval
}

type tMostOrLeastFrequentSortPair struct {
	count       int64
	groupingKey string
}

func NewTransformerMostOrLeastFrequent(
	groupByFieldNames []string,
	maxOutputLength int64,
	showCounts bool,
	outputFieldName string,
	descending bool,
) (*TransformerMostOrLeastFrequent, error) {
	tr := &TransformerMostOrLeastFrequent{
		groupByFieldNames: groupByFieldNames,
		maxOutputLength:   maxOutputLength,
		showCounts:        showCounts,
		outputFieldName:   outputFieldName,
		descending:        descending,
		countsByGroup:     lib.NewOrderedMap[int64](),
		valuesForGroup:    make(map[string][]*mlrval.Mlrval),
	}

	return tr, nil
}

func (tr *TransformerMostOrLeastFrequent) Transform(
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

		iCount, ok := tr.countsByGroup.GetWithCheck(groupingKey)
		if !ok {
			tr.countsByGroup.Put(groupingKey, int64(1))
		} else {
			tr.countsByGroup.Put(groupingKey, iCount+1)
		}
		if tr.valuesForGroup[groupingKey] == nil {
			selectedValues, _ := inrec.GetSelectedValues(tr.groupByFieldNames)
			tr.valuesForGroup[groupingKey] = selectedValues
		}

	} else {
		// TODO: Use a heap so this would be m log(n) not n log(n), where m is
		// the output length and n is the input length. (Each delete-max would
		// be O(log n) and there would be m of them.)

		// Copy keys and counters from hashmap to array for sorting
		inputLength := tr.countsByGroup.FieldCount

		sortPairs := make([]tMostOrLeastFrequentSortPair, inputLength)
		i := 0
		for pe := tr.countsByGroup.Head; pe != nil; pe = pe.Next {
			groupingKey := pe.Key
			count := pe.Value
			sortPairs[i].groupingKey = groupingKey
			sortPairs[i].count = count
			i++
		}

		// Sort by count
		// Go sort API: for ascending sort, return true if element i < element j.
		if tr.descending {

			sort.Slice(sortPairs, func(i, j int) bool {
				return sortPairs[i].count > sortPairs[j].count
			})

		} else {

			sort.Slice(sortPairs, func(i, j int) bool {
				return sortPairs[i].count < sortPairs[j].count
			})

		}

		// Emit top n
		outputLength := inputLength
		if inputLength > tr.maxOutputLength {
			outputLength = tr.maxOutputLength
		}
		for i := int64(0); i < outputLength; i++ {
			outrec := mlrval.NewMlrmapAsRecord()
			groupByFieldValues := tr.valuesForGroup[sortPairs[i].groupingKey]
			for j := range tr.groupByFieldNames {
				outrec.PutCopy(
					tr.groupByFieldNames[j],
					groupByFieldValues[j],
				)
			}

			if tr.showCounts {
				outrec.PutReference(tr.outputFieldName, mlrval.FromInt(sortPairs[i].count))
			}
			*outputRecordsAndContexts = append(*outputRecordsAndContexts, types.NewRecordAndContext(outrec, &inrecAndContext.Context))
		}

		*outputRecordsAndContexts = append(*outputRecordsAndContexts, inrecAndContext) // End-of-stream marker
	}
}
