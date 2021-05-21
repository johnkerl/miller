package transformers

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameMostFrequent = "most-frequent"
const verbNameLeastFrequent = "least-frequent"

const mostLeastFrequentDefaultMaxOutputLength = 10
const mostLeastFrequentDefaultOutputFieldName = "count"

var MostFrequentSetup = transforming.TransformerSetup{
	Verb:         verbNameMostFrequent,
	UsageFunc:    transformerMostFrequentUsage,
	ParseCLIFunc: transformerMostFrequentParseCLI,
	IgnoresInput: false,
}

var LeastFrequentSetup = transforming.TransformerSetup{
	Verb:         verbNameLeastFrequent,
	UsageFunc:    transformerLeastFrequentUsage,
	ParseCLIFunc: transformerLeastFrequentParseCLI,
	IgnoresInput: false,
}

func transformerMostFrequentUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := lib.MlrExeName()
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
	if doExit {
		os.Exit(exitCode)
	}
}

func transformerLeastFrequentUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := lib.MlrExeName()
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
	if doExit {
		os.Exit(exitCode)
	}
}

func transformerMostFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TOptions,
) transforming.IRecordTransformer {
	return transformerMostOrLeastFrequentParseCLI(pargi, argc, args, true, transformerMostFrequentUsage)
}

func transformerLeastFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cliutil.TOptions,
) transforming.IRecordTransformer {
	return transformerMostOrLeastFrequentParseCLI(pargi, argc, args, false, transformerLeastFrequentUsage)
}

func transformerMostOrLeastFrequentParseCLI(
	pargi *int,
	argc int,
	args []string,
	descending bool,
	usageFunc transforming.TransformerUsageFunc,
) transforming.IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var groupByFieldNames []string = nil
	maxOutputLength := mostLeastFrequentDefaultMaxOutputLength
	showCounts := true
	outputFieldName := mostLeastFrequentDefaultOutputFieldName

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			usageFunc(os.Stdout, true, 0)

		} else if opt == "-f" {
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-n" {
			maxOutputLength = cliutil.VerbGetIntArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-b" {
			showCounts = false

		} else if opt == "-o" {
			outputFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			usageFunc(os.Stderr, true, 1)
		}
	}

	if groupByFieldNames == nil {
		usageFunc(os.Stderr, true, 1)
		return nil
	}

	transformer, err := NewTransformerMostOrLeastFrequent(
		groupByFieldNames,
		maxOutputLength,
		showCounts,
		outputFieldName,
		descending,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerMostOrLeastFrequent struct {
	groupByFieldNames []string
	maxOutputLength   int
	showCounts        bool
	outputFieldName   string
	descending        bool
	countsByGroup     map[string]int
	valuesForGroup    map[string][]*types.Mlrval
}

type tMostOrLeastFrequentSortPair struct {
	count       int
	groupingKey string
}

// ----------------------------------------------------------------
func NewTransformerMostOrLeastFrequent(
	groupByFieldNames []string,
	maxOutputLength int,
	showCounts bool,
	outputFieldName string,
	descending bool,
) (*TransformerMostOrLeastFrequent, error) {
	this := &TransformerMostOrLeastFrequent{
		groupByFieldNames: groupByFieldNames,
		maxOutputLength:   maxOutputLength,
		showCounts:        showCounts,
		outputFieldName:   outputFieldName,
		descending:        descending,
		countsByGroup:     make(map[string]int),
		valuesForGroup:    make(map[string][]*types.Mlrval),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerMostOrLeastFrequent) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNames)
		if !ok {
			return
		}

		this.countsByGroup[groupingKey]++
		if this.valuesForGroup[groupingKey] == nil {
			selectedValues, _ := inrec.GetSelectedValues(this.groupByFieldNames)
			this.valuesForGroup[groupingKey] = selectedValues
		}

	} else {
		// TODO: Use a heap so this would be m log(n) not n log(n), where m is
		// the output length and n is the input length. (Each delete-max would
		// be O(log n) and there would be m of them.)

		// Copy keys and counters from hashmap to array for sorting
		inputLength := len(this.countsByGroup)

		sortPairs := make([]tMostOrLeastFrequentSortPair, inputLength)
		i := 0
		for groupingKey, count := range this.countsByGroup {
			sortPairs[i].groupingKey = groupingKey
			sortPairs[i].count = count
			i++
		}

		// Sort by count
		// Go sort API: for ascending sort, return true if element i < element j.
		if this.descending {

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
		if inputLength > this.maxOutputLength {
			outputLength = this.maxOutputLength
		}
		for i := 0; i < outputLength; i++ {
			outrec := types.NewMlrmapAsRecord()
			groupByFieldValues := this.valuesForGroup[sortPairs[i].groupingKey]
			for j, _ := range this.groupByFieldNames {
				outrec.PutCopy(
					this.groupByFieldNames[j],
					groupByFieldValues[j],
				)
			}

			if this.showCounts {
				outrec.PutReference(this.outputFieldName, types.MlrvalPointerFromInt(sortPairs[i].count))
			}
			outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // End-of-stream marker
	}
}
