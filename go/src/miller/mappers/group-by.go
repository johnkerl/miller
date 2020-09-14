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
var GroupBySetup = mapping.MapperSetup{
	Verb:         "group-by",
	ParseCLIFunc: mapperGroupByParseCLI,
	IgnoresInput: false,
}

func mapperGroupByParseCLI(
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

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperGroupByUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the group-by field names from the command line
	if argi >= argc {
		flagSet.Usage()
		os.Exit(1)
	}
	groupByFieldNames := args[argi]
	argi += 1

	mapper, _ := NewMapperGroupBy(
		groupByFieldNames,
	)

	*pargi = argi
	return mapper
}

func mapperGroupByUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Outputs records in batches having identical values at specified field names.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperGroupBy struct {
	// input
	groupByFieldNameList []string

	// state
	recordListsByGroup map[string]*list.List
}

func NewMapperGroupBy(
	groupByFieldNames string,
) (*MapperGroupBy, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &MapperGroupBy{
		groupByFieldNameList: groupByFieldNameList,

		recordListsByGroup: make(map[string]*list.List),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperGroupBy) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		groupByKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
		if !ok {
			return
		}

		recordListForGroup, present := this.recordListsByGroup[groupByKey]
		if !present { // first time
			recordListForGroup = list.New()
			this.recordListsByGroup[groupByKey] = recordListForGroup
		}

		recordListForGroup.PushBack(inrecAndContext)

	} else {
		for _, recordListForGroup := range this.recordListsByGroup {
			for entry := recordListForGroup.Front(); entry != nil; entry = entry.Next() {
				outputChannel <- entry.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}
