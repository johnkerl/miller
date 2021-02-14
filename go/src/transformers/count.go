package transformers

import (
	"fmt"
	"os"
	"strings"

	"miller/src/cliutil"
	"miller/src/lib"
	"miller/src/transforming"
	"miller/src/types"
)

// ----------------------------------------------------------------
const verbNameCount = "count"

var CountSetup = transforming.TransformerSetup{
	Verb:         verbNameCount,
	ParseCLIFunc: transformerCountParseCLI,
	UsageFunc:    transformerCountUsage,
	IgnoresInput: false,
}

func transformerCountParseCLI(
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
			groupByFieldNames = cliutil.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-n" {
			showCountsOnly = true

		} else if opt == "-o" {
			outputFieldName = cliutil.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else {
			transformerCountUsage(os.Stderr, true, 1)
		}
	}

	transformer, _ := NewTransformerCount(
		groupByFieldNames,
		showCountsOnly,
		outputFieldName,
	)

	*pargi = argi
	return transformer
}

func transformerCountUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", lib.MlrExeName(), verbNameCount)
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

// ----------------------------------------------------------------
type TransformerCount struct {
	// input
	groupByFieldNames []string
	showCountsOnly    bool
	outputFieldName   string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
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

	this := &TransformerCount{
		groupByFieldNames: groupByFieldNames,
		showCountsOnly:    showCountsOnly,
		outputFieldName:   outputFieldName,

		ungroupedCount: 0,
		groupedCounts:  lib.NewOrderedMap(),
		groupingValues: lib.NewOrderedMap(),
	}

	if groupByFieldNames == nil {
		this.recordTransformerFunc = this.countUngrouped
	} else {
		this.recordTransformerFunc = this.countGrouped
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *TransformerCount) Transform(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordTransformerFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *TransformerCount) countUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		this.ungroupedCount++
	} else {
		newrec := types.NewMlrmapAsRecord()
		mcount := types.MlrvalFromInt(this.ungroupedCount)
		newrec.PutCopy(this.outputFieldName, &mcount)

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (this *TransformerCount) countGrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, selectedValues, ok := inrec.GetSelectedValuesAndJoined(
			this.groupByFieldNames,
		)
		if !ok { // Current record does not have specified fields; ignore
			return
		}

		if !this.groupedCounts.Has(groupingKey) {
			var count int = 1
			this.groupedCounts.Put(groupingKey, count)
			this.groupingValues.Put(groupingKey, selectedValues)
		} else {
			this.groupedCounts.Put(
				groupingKey,
				this.groupedCounts.Get(groupingKey).(int)+1,
			)
		}

	} else {
		if this.showCountsOnly {
			newrec := types.NewMlrmapAsRecord()
			mcount := types.MlrvalFromInt(this.groupedCounts.FieldCount)
			newrec.PutCopy(this.outputFieldName, &mcount)

			outrecAndContext := types.NewRecordAndContext(newrec, &inrecAndContext.Context)
			outputChannel <- outrecAndContext

		} else {
			for outer := this.groupedCounts.Head; outer != nil; outer = outer.Next {
				groupingKey := outer.Key
				newrec := types.NewMlrmapAsRecord()

				// Example:
				// * Suppose group-by fields are a,b.
				// * Record has a=foo,b=bar
				// * Grouping key is "foo,bar"
				// * Grouping values for key is ["foo", "bar"]
				// Here we populate a record with "a=foo,b=bar".

				groupingValuesForKey := this.groupingValues.Get(groupingKey).([]*types.Mlrval)
				i := 0
				for _, groupingValueForKey := range groupingValuesForKey {
					newrec.PutCopy(this.groupByFieldNames[i], groupingValueForKey)
					i++
				}

				countForGroup := outer.Value.(int)
				mcount := types.MlrvalFromInt(countForGroup)
				newrec.PutCopy(this.outputFieldName, &mcount)

				outrecAndContext := types.NewRecordAndContext(newrec, &inrecAndContext.Context)
				outputChannel <- outrecAndContext
			}
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
