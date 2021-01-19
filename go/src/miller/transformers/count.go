package transformers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/transforming"
	"miller/types"
)

// ----------------------------------------------------------------
var CountSetup = transforming.TransformerSetup{
	Verb:         "count",
	ParseCLIFunc: transformerCountParseCLI,
	IgnoresInput: false,
}

func transformerCountParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) transforming.IRecordTransformer {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional group-by-field names for counts, e.g. a,b,c",
	)

	pShowCountsOnly := flagSet.Bool(
		"n",
		false,
		"Show only the number of distinct values. Not interesting without -g.",
	)

	pOutputFieldName := flagSet.String(
		"o",
		"count",
		`Field name for output count`,
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		transformerCountUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	transformer, _ := NewTransformerCount(
		*pGroupByFieldNames,
		*pShowCountsOnly,
		*pOutputFieldName,
	)

	*pargi = argi
	return transformer
}

func transformerCountUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Prints number of records, optionally grouped by distinct values for specified field names.
`)
	fmt.Fprintf(o, "Options:\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		if f.Name == "g" {
			fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage) // f.Name, f.Value
		} else {
			fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
		}
	})
}

// ----------------------------------------------------------------
type TransformerCount struct {
	// input
	groupByFieldNameList []string
	showCountsOnly       bool
	outputFieldName      string

	// state
	recordTransformerFunc transforming.RecordTransformerFunc
	ungroupedCount        int64
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
	groupByFieldNames string,
	showCountsOnly bool,
	outputFieldName string,
) (*TransformerCount, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &TransformerCount{
		groupByFieldNameList: groupByFieldNameList,
		showCountsOnly:       showCountsOnly,
		outputFieldName:      outputFieldName,

		ungroupedCount: 0,
		groupedCounts:  lib.NewOrderedMap(),
		groupingValues: lib.NewOrderedMap(),
	}

	if len(groupByFieldNameList) == 0 {
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
		mcount := types.MlrvalFromInt64(this.ungroupedCount)
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
			this.groupByFieldNameList,
		)
		if !ok { // Current record does not have specified fields; ignore
			return
		}

		if !this.groupedCounts.Has(groupingKey) {
			var count int64 = 1
			this.groupedCounts.Put(groupingKey, count)
			this.groupingValues.Put(groupingKey, selectedValues)
		} else {
			this.groupedCounts.Put(
				groupingKey,
				this.groupedCounts.Get(groupingKey).(int64)+1,
			)
		}

	} else {
		if this.showCountsOnly {
			newrec := types.NewMlrmapAsRecord()
			mcount := types.MlrvalFromInt64(this.groupedCounts.FieldCount)
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
					newrec.PutCopy(this.groupByFieldNameList[i], groupingValueForKey)
					i++
				}

				countForGroup := outer.Value.(int64)
				mcount := types.MlrvalFromInt64(countForGroup)
				newrec.PutCopy(this.outputFieldName, &mcount)

				outrecAndContext := types.NewRecordAndContext(newrec, &inrecAndContext.Context)
				outputChannel <- outrecAndContext
			}
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
