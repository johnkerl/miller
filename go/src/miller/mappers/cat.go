package mappers

import (
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var CatSetup = mapping.MapperSetup{
	Verb:         "cat",
	ParseCLIFunc: mapperCatParseCLI,
	IgnoresInput: false,
}

func mapperCatParseCLI(
	pargi *int,
	argc int,
	args []string,
	errorHandling flag.ErrorHandling, // ContinueOnError or ExitOnError
	_ *clitypes.TReaderOptions,
	__ *clitypes.TWriterOptions,
) mapping.IRecordMapper {

	// Get the verb name from the current spot in the mlr command line
	argi := *pargi
	verb := args[argi]
	argi++

	// Parse local flags
	flagSet := flag.NewFlagSet(verb, errorHandling)

	pDoCounters := flagSet.Bool(
		"n",
		false,
		"Prepend field \"n\" to each record with record-counter starting at 1",
	)

	pCounterFieldName := flagSet.String(
		"N",
		"",
		"Prepend field {name} to each record with record-counter starting at 1",
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional group-by-field names for counters, e.g. a,b,c",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperCatUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// TODO:
	//	fmt.Fprintf(o, "-v        Write a low-level record-structure dump to stderr.\n");

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperCat(
		*pDoCounters,
		pCounterFieldName,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return mapper
}

func mapperCatUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperCat struct {
	doCounters           bool
	groupByFieldNameList []string

	counter          int64
	countsByGroup    map[string]int64
	counterFieldName string

	recordMapperFunc mapping.RecordMapperFunc
}

// ----------------------------------------------------------------
func NewMapperCat(
	doCounters bool,
	pCounterFieldName *string,
	groupByFieldNames string,
) (*MapperCat, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	counterFieldName := "n"
	if *pCounterFieldName != "" {
		counterFieldName = *pCounterFieldName
		doCounters = true
	}

	this := &MapperCat{
		doCounters:           doCounters,
		groupByFieldNameList: groupByFieldNameList,
		counter:              0,
		countsByGroup:        make(map[string]int64),
		counterFieldName:     counterFieldName,
	}

	if !doCounters {
		this.recordMapperFunc = this.simpleCat
	} else {
		if groupByFieldNames == "" {
			this.recordMapperFunc = this.countersUngrouped
		} else {
			this.recordMapperFunc = this.countersGrouped
		}
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperCat) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordMapperFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *MapperCat) simpleCat(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (this *MapperCat) countersUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	record := inrecAndContext.Record
	if record != nil { // not end of record stream
		this.counter++
		key := this.counterFieldName
		value := types.MlrvalFromInt64(this.counter)
		record.PrependCopy(&key, &value)
	}
	outputChannel <- inrecAndContext
}

// ----------------------------------------------------------------
func (this *MapperCat) countersGrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
		var counter int64 = 0
		if !ok {
			// Treat as unkeyed
			this.counter++
			counter = this.counter
		} else {
			counter, ok = this.countsByGroup[groupingKey]
			if ok {
				counter++
			} else {
				counter = 1
			}
			this.countsByGroup[groupingKey] = counter
		}

		key := this.counterFieldName
		value := types.MlrvalFromInt64(counter)
		inrec.PrependCopy(&key, &value)
	}
	outputChannel <- inrecAndContext
}
