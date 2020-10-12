package mappers

import (
	"container/list"
	"flag"
	"fmt"
	"os"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
	"miller/types"
)

// ----------------------------------------------------------------
var CountSimilarSetup = mapping.MapperSetup{
	Verb:         "count-similar",
	ParseCLIFunc: mapperCountSimilarParseCLI,
	IgnoresInput: false,
}

func mapperCountSimilarParseCLI(
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

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Group-by-field names for counts, e.g. a,b,c",
	)

	pCounterFieldName := flagSet.String(
		"o",
		"count",
		"Field name for output count.",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperCountSimilarUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	if *pGroupByFieldNames == "" {
		mapperCountSimilarUsage(os.Stderr, args[0], verb, flagSet)
		os.Exit(1)
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperCountSimilar(
		*pGroupByFieldNames,
		*pCounterFieldName,
	)

	*pargi = argi
	return mapper
}

func mapperCountSimilarUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprintf(o, "Ingests all records, then emits each record augmented by a count of\n")
	fmt.Fprintf(o, "the number of other records having the same group-by field values.\n")
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		if f.Name == "g" {
			fmt.Fprintf(o, " -%v %v\n", f.Name, f.Usage)
		} else {
			fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage)
		}
	})
}

// ----------------------------------------------------------------
type MapperCountSimilar struct {
	// Input:
	groupByFieldNameList []string
	counterFieldName     string

	// State:
	recordListsByGroup *lib.OrderedMap // map from string to *list.List
}

// ----------------------------------------------------------------
func NewMapperCountSimilar(
	groupByFieldNames string,
	counterFieldName string,
) (*MapperCountSimilar, error) {
	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")
	this := &MapperCountSimilar{
		groupByFieldNameList: groupByFieldNameList,
		counterFieldName:     counterFieldName,
		recordListsByGroup:   lib.NewOrderedMap(),
	}
	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperCountSimilar) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
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
			mgroupSize := types.MlrvalFromInt64(int64(groupSize))
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				recordAndContext := inner.Value.(*types.RecordAndContext)
				recordAndContext.Record.PutCopy(&this.counterFieldName, &mgroupSize)

				outputChannel <- recordAndContext
			}
		}

		outputChannel <- inrecAndContext // Emit the stream-terminating null record
	}
}
