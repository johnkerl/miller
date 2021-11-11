package transformers

import (
	"fmt"
	"os"
	"strings"

	"mlr/internal/pkg/cli"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
const verbNameCount = "count"

var CountSetup = TransformerSetup{
	Verb:         verbNameCount,
	UsageFunc:    transformerCountUsage,
	ParseCLIFunc: transformerCountParseCLI,
	IgnoresInput: false,
}

func transformerCountUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", "mlr", verbNameCount)
	fmt.Fprint(o,
		`Prints number of records, optionally grouped by distinct values for specified field names.
`)
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {a,b,c} Optional group-by-field names for counts, e.g. a,b,c\n")
	fmt.Fprintf(o, "-n {n} Show only the number of distinct values. Not interesting without -g.\n")
	fmt.Fprintf(o, "-o {name} Field name for output-count. Default \"count\".\n")
	fmt.Fprintf(o, "-h|--help Show this message.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCountParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	var groupByFieldNames []string = nil
	showCountsOnly := false
	outputFieldName := "count"

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCountUsage(os.Stdout, true, 0)

		} else if opt == "-g" {
			groupByFieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-n" {
			showCountsOnly = true

		} else if opt == "-o" {
			outputFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCountUsage(os.Stderr, true, 1)
		}
	}

	transformer, err := NewTransformerCount(
		groupByFieldNames,
		showCountsOnly,
		outputFieldName,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerCount struct {
	// input
	groupByFieldNames []string
	showCountsOnly    bool
	outputFieldName   string

	// state
	recordTransformerFunc RecordTransformerFunc
	ungroupedCount        int
	// Example:
	// * Suppose group-by fields are a,b.
	// * One record has a=foo,b=bar
	// * Another record has a=baz,b=quux
	// * Map keys are strings "foo,bar" and "baz,quux".
	// * groupedCounts maps "foo,bar" to 1 and "baz,quux" to 1.
	// * groupByValues maps "foo,bar" to ["foo", "bar"] and "baz,quux" to ["baz", "quux"].
	groupedCounts  *lib.OrderedMap
	groupingValues *lib.OrderedMap
}

func NewTransformerCount(
	groupByFieldNames []string,
	showCountsOnly bool,
	outputFieldName string,
) (*TransformerCount, error) {

	tr := &TransformerCount{
		groupByFieldNames: groupByFieldNames,
		showCountsOnly:    showCountsOnly,
		outputFieldName:   outputFieldName,

		ungroupedCount: 0,
		groupedCounts:  lib.NewOrderedMap(),
		groupingValues: lib.NewOrderedMap(),
	}

	if groupByFieldNames == nil {
		tr.recordTransformerFunc = tr.countUngrouped
	} else {
		tr.recordTransformerFunc = tr.countGrouped
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerCount) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
func (tr *TransformerCount) countUngrouped(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		tr.ungroupedCount++
	} else {
		newrec := types.NewMlrmapAsRecord()
		newrec.PutCopy(tr.outputFieldName, types.MlrvalFromInt(tr.ungroupedCount))

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerCount) countGrouped(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, selectedValues, ok := inrec.GetSelectedValuesAndJoined(
			tr.groupByFieldNames,
		)
		if !ok { // Current record does not have specified fields; ignore
			return
		}

		if !tr.groupedCounts.Has(groupingKey) {
			var count int = 1
			tr.groupedCounts.Put(groupingKey, count)
			tr.groupingValues.Put(groupingKey, selectedValues)
		} else {
			tr.groupedCounts.Put(
				groupingKey,
				tr.groupedCounts.Get(groupingKey).(int)+1,
			)
		}

	} else {
		if tr.showCountsOnly {
			newrec := types.NewMlrmapAsRecord()
			newrec.PutCopy(tr.outputFieldName, types.MlrvalFromInt(tr.groupedCounts.FieldCount))

			outrecAndContext := types.NewRecordAndContext(newrec, &inrecAndContext.Context)
			outputChannel <- outrecAndContext

		} else {
			for outer := tr.groupedCounts.Head; outer != nil; outer = outer.Next {
				groupingKey := outer.Key
				newrec := types.NewMlrmapAsRecord()

				// Example:
				// * Suppose group-by fields are a,b.
				// * Record has a=foo,b=bar
				// * Grouping key is "foo,bar"
				// * Grouping values for key is ["foo", "bar"]
				// Here we populate a record with "a=foo,b=bar".

				groupingValuesForKey := tr.groupingValues.Get(groupingKey).([]*types.Mlrval)
				i := 0
				for _, groupingValueForKey := range groupingValuesForKey {
					newrec.PutCopy(tr.groupByFieldNames[i], groupingValueForKey)
					i++
				}

				countForGroup := outer.Value.(int)
				newrec.PutCopy(tr.outputFieldName, types.MlrvalFromInt(countForGroup))

				outrecAndContext := types.NewRecordAndContext(newrec, &inrecAndContext.Context)
				outputChannel <- outrecAndContext
			}
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
