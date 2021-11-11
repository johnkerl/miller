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
const verbNameCountDistinct = "count-distinct"
const verbNameUniq = "uniq"
const uniqDefaultOutputFieldName = "count"

var CountDistinctSetup = TransformerSetup{
	Verb:         verbNameCountDistinct,
	UsageFunc:    transformerCountDistinctUsage,
	ParseCLIFunc: transformerCountDistinctParseCLI,
	IgnoresInput: false,
}

var UniqSetup = TransformerSetup{
	Verb:         verbNameUniq,
	UsageFunc:    transformerUniqUsage,
	ParseCLIFunc: transformerUniqParseCLI,
	IgnoresInput: false,
}

// ----------------------------------------------------------------
func transformerCountDistinctUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameCountDistinct
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Prints number of records having distinct values for specified field names.\n")
	fmt.Fprintf(o, "Same as uniq -c.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-f {a,b,c}    Field names for distinct count.\n")
	fmt.Fprintf(o, "-n            Show only the number of distinct values. Not compatible with -u.\n")
	fmt.Fprintf(o, "-o {name}     Field name for output count. Default \"%s\".\n", uniqDefaultOutputFieldName)
	fmt.Fprintf(o, "              Ignored with -u.\n")
	fmt.Fprintf(o, "-u            Do unlashed counts for multiple field names. With -f a,b and\n")
	fmt.Fprintf(o, "              without -u, computes counts for distinct combinations of a\n")
	fmt.Fprintf(o, "              and b field values. With -f a,b and with -u, computes counts\n")
	fmt.Fprintf(o, "              for distinct a field values and counts for distinct b field\n")
	fmt.Fprintf(o, "              values separately.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerCountDistinctParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fieldNames []string = nil
	showNumDistinctOnly := false
	outputFieldName := uniqDefaultOutputFieldName
	doLashed := true

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerCountDistinctUsage(os.Stdout, true, 0)

		} else if opt == "-g" || opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-n" {
			showNumDistinctOnly = true

		} else if opt == "-o" {
			outputFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-u" {
			doLashed = false

		} else {
			transformerCountDistinctUsage(os.Stderr, true, 1)
		}
	}

	if fieldNames == nil {
		transformerCountDistinctUsage(os.Stderr, true, 1)
	}
	if !doLashed && showNumDistinctOnly {
		transformerCountDistinctUsage(os.Stderr, true, 1)
	}

	showCounts := true
	uniqifyEntireRecords := false

	transformer, err := NewTransformerUniq(
		fieldNames,
		showCounts,
		showNumDistinctOnly,
		outputFieldName,
		doLashed,
		uniqifyEntireRecords,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
func transformerUniqUsage(
	o *os.File,
	doExit bool,
	exitCode int,
) {
	argv0 := "mlr"
	verb := verbNameUniq
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Prints distinct values for specified field names. With -c, same as\n")
	fmt.Fprintf(o, "count-distinct. For uniq, -f is a synonym for -g.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Options:\n")
	fmt.Fprintf(o, "-g {d,e,f}    Group-by-field names for uniq counts.\n")
	fmt.Fprintf(o, "-c            Show repeat counts in addition to unique values.\n")
	fmt.Fprintf(o, "-n            Show only the number of distinct values.\n")
	fmt.Fprintf(o, "-o {name}     Field name for output count. Default \"%s\".\n", uniqDefaultOutputFieldName)
	fmt.Fprintf(o, "-a            Output each unique record only once. Incompatible with -g.\n")
	fmt.Fprintf(o, "              With -c, produces unique records, with repeat counts for each.\n")
	fmt.Fprintf(o, "              With -n, produces only one record which is the unique-record count.\n")
	fmt.Fprintf(o, "              With neither -c nor -n, produces unique records.\n")

	if doExit {
		os.Exit(exitCode)
	}
}

func transformerUniqParseCLI(
	pargi *int,
	argc int,
	args []string,
	_ *cli.TOptions,
) IRecordTransformer {

	// Skip the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	var fieldNames []string = nil
	showCounts := false
	showNumDistinctOnly := false
	outputFieldName := uniqDefaultOutputFieldName
	uniqifyEntireRecords := false

	for argi < argc /* variable increment: 1 or 2 depending on flag */ {
		opt := args[argi]
		if !strings.HasPrefix(opt, "-") {
			break // No more flag options to process
		}
		argi++

		if opt == "-h" || opt == "--help" {
			transformerUniqUsage(os.Stdout, true, 0)

		} else if opt == "-g" || opt == "-f" {
			fieldNames = cli.VerbGetStringArrayArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-c" {
			showCounts = true

		} else if opt == "-n" {
			showNumDistinctOnly = true

		} else if opt == "-o" {
			outputFieldName = cli.VerbGetStringArgOrDie(verb, opt, args, &argi, argc)

		} else if opt == "-a" {
			uniqifyEntireRecords = true

		} else {
			transformerUniqUsage(os.Stderr, true, 1)
		}
	}

	if uniqifyEntireRecords {
		if fieldNames != nil {
			transformerUniqUsage(os.Stderr, true, 1)
		}
		if showCounts && showNumDistinctOnly {
			transformerUniqUsage(os.Stderr, true, 1)
		}
	} else {
		if fieldNames == nil {
			transformerUniqUsage(os.Stderr, true, 1)
		}
	}

	doLashed := true

	transformer, _ := NewTransformerUniq(
		fieldNames,
		showCounts,
		showNumDistinctOnly,
		outputFieldName,
		doLashed,
		uniqifyEntireRecords,
	)

	*pargi = argi
	return transformer
}

// ----------------------------------------------------------------
type TransformerUniq struct {
	fieldNames      []string
	showCounts      bool
	outputFieldName string

	// Example:
	// Input is:
	//   a=1,b=2,c=3
	//   a=4,b=5,c=6
	// Uniquing on fields ["a","b"]
	// uniqifiedRecordCounts:
	//   '{"a":1,"b":2,"c":3}' => 1
	//   '{"a":4,"b":5,"c":6}' => 1
	// uniqifiedRecords:
	//   '{"a":1,"b":2,"c":3}' => {"a":1,"b":2,"c":3}
	//   '{"a":4,"b":5,"c":6}' => {"a":4,"b":5,"c":6}
	// countsByGroup:
	//  "1,2" -> 1
	//  "4,5" -> 1
	// valuesByGroup:
	//  "1,2" -> [1,2]
	//  "4,5" -> [4,5]
	// unlashedCounts:
	//   "a" => "1" => 1
	//   "a" => "4" => 1
	//   ...
	// unlashedCountValues:
	//   "a" => "1" => 1
	//   "a" => "4" => 4
	uniqifiedRecordCounts *lib.OrderedMap // record-as-string -> counts
	uniqifiedRecords      *lib.OrderedMap // record-as-string -> records
	countsByGroup         *lib.OrderedMap // grouping key -> count
	valuesByGroup         *lib.OrderedMap // grouping key -> array of values
	unlashedCounts        *lib.OrderedMap // field name -> string field value -> count
	unlashedCountValues   *lib.OrderedMap // field name -> string field value -> typed field value

	recordTransformerFunc RecordTransformerFunc
}

// ----------------------------------------------------------------
func NewTransformerUniq(
	fieldNames []string,
	showCounts bool,
	showNumDistinctOnly bool,
	outputFieldName string,
	doLashed bool,
	uniqifyEntireRecords bool,
) (*TransformerUniq, error) {

	tr := &TransformerUniq{
		fieldNames:      fieldNames,
		showCounts:      showCounts,
		outputFieldName: outputFieldName,

		uniqifiedRecordCounts: lib.NewOrderedMap(),
		uniqifiedRecords:      lib.NewOrderedMap(),
		countsByGroup:         lib.NewOrderedMap(),
		valuesByGroup:         lib.NewOrderedMap(),
		unlashedCounts:        lib.NewOrderedMap(),
		unlashedCountValues:   lib.NewOrderedMap(),
	}

	if uniqifyEntireRecords {
		if showCounts {
			tr.recordTransformerFunc = tr.transformUniqifyEntireRecordsShowCounts
		} else if showNumDistinctOnly {
			tr.recordTransformerFunc = tr.transformUniqifyEntireRecordsShowNumDistinctOnly
		} else {
			tr.recordTransformerFunc = tr.transformUniqifyEntireRecords
		}
	} else if !doLashed {
		tr.recordTransformerFunc = tr.transformUnlashed
	} else if showNumDistinctOnly {
		tr.recordTransformerFunc = tr.transformNumDistinctOnly
	} else if showCounts {
		tr.recordTransformerFunc = tr.transformWithCounts
	} else {
		tr.recordTransformerFunc = tr.transformWithoutCounts
	}

	return tr, nil
}

// ----------------------------------------------------------------

func (tr *TransformerUniq) Transform(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	HandleDefaultDownstreamDone(inputDownstreamDoneChannel, outputDownstreamDoneChannel)
	tr.recordTransformerFunc(inrecAndContext, inputDownstreamDoneChannel, outputDownstreamDoneChannel, outputChannel)
}

// ----------------------------------------------------------------
// Print each unique record only once, with uniqueness counts.  This means
// non-streaming, with output at end of stream.
func (tr *TransformerUniq) transformUniqifyEntireRecordsShowCounts(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		recordAsString := inrec.String()
		icount, present := tr.uniqifiedRecordCounts.GetWithCheck(recordAsString)
		if !present { // first time seen
			tr.uniqifiedRecordCounts.Put(recordAsString, 1)
			tr.uniqifiedRecords.Put(recordAsString, inrecAndContext.Copy())
		} else { // have seen before
			tr.uniqifiedRecordCounts.Put(recordAsString, icount.(int)+1)
		}

	} else { // end of record stream

		for pe := tr.uniqifiedRecords.Head; pe != nil; pe = pe.Next {
			outrecAndContext := pe.Value.(*types.RecordAndContext)
			icount := tr.uniqifiedRecordCounts.Get(pe.Key)
			mcount := types.MlrvalFromInt(icount.(int))
			outrecAndContext.Record.PrependReference(tr.outputFieldName, mcount)
			outputChannel <- outrecAndContext
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}

}

// ----------------------------------------------------------------
// Print count of unique records.  This means non-streaming, with output at end
// of stream.
func (tr *TransformerUniq) transformUniqifyEntireRecordsShowNumDistinctOnly(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record
		recordAsString := inrec.String()
		if !tr.uniqifiedRecordCounts.Has(recordAsString) {
			tr.uniqifiedRecordCounts.Put(recordAsString, 1)
		}

	} else { // end of record stream
		outrec := types.NewMlrmapAsRecord()
		outrec.PutReference(
			tr.outputFieldName,
			types.MlrvalFromInt(tr.uniqifiedRecordCounts.FieldCount),
		)
		outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
// Print each unique record only once (on first occurrence).
func (tr *TransformerUniq) transformUniqifyEntireRecords(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		recordAsString := inrec.String()
		if !tr.uniqifiedRecordCounts.Has(recordAsString) {
			tr.uniqifiedRecordCounts.Put(recordAsString, 1)
			outputChannel <- inrecAndContext
		}

	} else { // end of record stream

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerUniq) transformUnlashed(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		for _, fieldName := range tr.fieldNames {
			var countsForFieldName *lib.OrderedMap = nil
			iCountsForFieldName, present := tr.unlashedCounts.GetWithCheck(fieldName)
			if !present {
				countsForFieldName = lib.NewOrderedMap()
				tr.unlashedCounts.Put(fieldName, countsForFieldName)
				tr.unlashedCountValues.Put(fieldName, lib.NewOrderedMap())
			} else {
				countsForFieldName = iCountsForFieldName.(*lib.OrderedMap)
			}

			fieldValue := inrec.Get(fieldName)
			if fieldValue != nil {
				fieldValueString := fieldValue.String()
				if !countsForFieldName.Has(fieldValueString) {
					countsForFieldName.Put(fieldValueString, 1)
					tr.unlashedCountValues.Get(fieldName).(*lib.OrderedMap).Put(fieldValueString, fieldValue.Copy())
				} else {
					countsForFieldName.Put(fieldValueString, countsForFieldName.Get(fieldValueString).(int)+1)
				}
			}
		}

	} else { // end of record stream

		for pe := tr.unlashedCounts.Head; pe != nil; pe = pe.Next {
			fieldName := pe.Key
			countsForFieldName := pe.Value.(*lib.OrderedMap)
			for pf := countsForFieldName.Head; pf != nil; pf = pf.Next {
				fieldValueString := pf.Key
				outrec := types.NewMlrmapAsRecord()
				outrec.PutReference("field", types.MlrvalFromString(fieldName))
				outrec.PutCopy(
					"value",
					tr.unlashedCountValues.Get(fieldName).(*lib.OrderedMap).Get(fieldValueString).(*types.Mlrval),
				)
				outrec.PutReference("count", types.MlrvalFromInt(pf.Value.(int)))
				outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
			}
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerUniq) transformNumDistinctOnly(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, ok := inrec.GetSelectedValuesJoined(tr.fieldNames)
		if ok {
			iCount, present := tr.countsByGroup.GetWithCheck(groupingKey)
			if !present {
				tr.countsByGroup.Put(groupingKey, 1)
			} else {
				tr.countsByGroup.Put(groupingKey, iCount.(int)+1)
			}
		}

	} else {
		outrec := types.NewMlrmapAsRecord()
		outrec.PutReference(
			"count",
			types.MlrvalFromInt(tr.countsByGroup.FieldCount),
		)
		outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerUniq) transformWithCounts(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, selectedValues, ok := inrec.GetSelectedValuesAndJoined(tr.fieldNames)
		if ok {
			iCount, present := tr.countsByGroup.GetWithCheck(groupingKey)
			if !present {
				tr.countsByGroup.Put(groupingKey, 1)
				tr.valuesByGroup.Put(groupingKey, selectedValues)
			} else {
				tr.countsByGroup.Put(groupingKey, iCount.(int)+1)
			}
		}

	} else { // end of record stream

		for pa := tr.countsByGroup.Head; pa != nil; pa = pa.Next {
			outrec := types.NewMlrmapAsRecord()
			valuesForGroup := tr.valuesByGroup.Get(pa.Key).([]*types.Mlrval)
			for i, fieldName := range tr.fieldNames {
				outrec.PutCopy(
					fieldName,
					valuesForGroup[i],
				)
			}
			if tr.showCounts {
				outrec.PutReference(
					tr.outputFieldName,
					types.MlrvalFromInt(pa.Value.(int)),
				)
			}
			outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)
		}

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

// ----------------------------------------------------------------
func (tr *TransformerUniq) transformWithoutCounts(
	inrecAndContext *types.RecordAndContext,
	inputDownstreamDoneChannel <-chan bool,
	outputDownstreamDoneChannel chan<- bool,
	outputChannel chan<- *types.RecordAndContext,
) {
	if !inrecAndContext.EndOfStream {
		inrec := inrecAndContext.Record

		groupingKey, selectedValues, ok := inrec.GetSelectedValuesAndJoined(tr.fieldNames)
		if !ok {
			return
		}

		iCount, present := tr.countsByGroup.GetWithCheck(groupingKey)
		if !present {
			tr.countsByGroup.Put(groupingKey, 1)
			tr.valuesByGroup.Put(groupingKey, selectedValues)
			outrec := types.NewMlrmapAsRecord()

			for i, fieldName := range tr.fieldNames {
				outrec.PutCopy(
					fieldName,
					selectedValues[i],
				)
			}

			outputChannel <- types.NewRecordAndContext(outrec, &inrecAndContext.Context)

		} else {
			tr.countsByGroup.Put(groupingKey, iCount.(int)+1)
		}

	} else { // end of record stream
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
