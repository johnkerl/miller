package mappers

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miller/clitypes"
	"miller/lib"
	"miller/mapping"
)

// ----------------------------------------------------------------
var HeadSetup = mapping.MapperSetup{
	Verb:         "head",
	ParseCLIFunc: mapperHeadParseCLI,
	IgnoresInput: false,
}

func mapperHeadParseCLI(
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

	//Usage: mlr head [options]
	//-n {count}    Head count to print; default 10
	//-g {a,b,c}    Optional group-by-field names for head counts

	pHeadCount := flagSet.Uint64(
		"n",
		10,
		`Head count to print`,
	)

	pGroupByFieldNames := flagSet.String(
		"g",
		"",
		"Optional group-by-field names for head counts, e.g. a,b,c",
	)

	flagSet.Usage = func() {
		ostream := os.Stderr
		if errorHandling == flag.ContinueOnError { // help intentionally requested
			ostream = os.Stdout
		}
		mapperHeadUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	//	if *pGroupByFieldNames == "" {
	//		flagSet.Usage()
	//		os.Exit(1)
	//	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperHead(
		*pHeadCount,
		*pGroupByFieldNames,
	)

	*pargi = argi
	return mapper
}

func mapperHeadUsage(
	o *os.File,
	argv0 string,
	verb string,
	flagSet *flag.FlagSet,
) {
	fmt.Fprintf(o, "Usage: %s %s [options]\n", argv0, verb)
	fmt.Fprint(o,
		`Passes through the first n records, optionally by category.  Without -g, ceases
consuming more input (i.e. is fast) when n records have been read.
`)
	// flagSet.PrintDefaults() doesn't let us control stdout vs stderr
	flagSet.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(o, " -%v (default %v) %v\n", f.Name, f.Value, f.Usage) // f.Name, f.Value
	})
}

// ----------------------------------------------------------------
type MapperHead struct {
	// input
	headCount            uint64
	groupByFieldNameList []string
	groupByFieldNameSet  map[string]bool

	// state
	recordMapperFunc   mapping.RecordMapperFunc
	unkeyedRecordCount uint64
	keyedRecordCounts map[string]uint64
}

func NewMapperHead(
	headCount uint64,
	groupByFieldNames string,
) (*MapperHead, error) {

	// xxx util function
	groupByFieldNameList := make([]string, 0)
	if groupByFieldNames != "" {
		groupByFieldNameList = strings.Split(groupByFieldNames, ",")
	}

	// xxx make/find-reuse util func
	groupByFieldNameSet := make(map[string]bool)
	for _, groupByFieldName := range groupByFieldNameList {
		groupByFieldNameSet[groupByFieldName] = true
	}

	this := &MapperHead{
		headCount:            headCount,
		groupByFieldNameList: groupByFieldNameList,
		groupByFieldNameSet:  groupByFieldNameSet,

		unkeyedRecordCount: 0,
		keyedRecordCounts: make(map[string]uint64),
	}

	if len(groupByFieldNameList) == 0 {
		this.recordMapperFunc = this.mapUnkeyed
	} else {
		this.recordMapperFunc = this.mapKeyed
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperHead) Map(
	inrecAndContext *lib.RecordAndContext,
	outputChannel chan<- *lib.RecordAndContext,
) {
	this.recordMapperFunc(inrecAndContext, outputChannel)
}

func (this *MapperHead) mapUnkeyed(
	inrecAndContext *lib.RecordAndContext,
	outputChannel chan<- *lib.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		this.unkeyedRecordCount++
		if this.unkeyedRecordCount <= this.headCount {
			outputChannel <- inrecAndContext
		}
	} else {
		outputChannel <- inrecAndContext
	}
}

func (this *MapperHead) mapKeyed(
	inrecAndContext *lib.RecordAndContext,
	outputChannel chan<- *lib.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

		//		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
		//			pstate->pgroup_by_field_names);
		//		if pgroup_by_field_values == nil {
		//			return
		//		}
		//
		//		uint64* pcount_for_group = lhmslv_get(pstate->pcounts_by_group,
		//			pgroup_by_field_values);
		//		if (pcount_for_group == NULL) {
		//			pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
		//			*pcount_for_group = 0LL;
		//			lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values),
		//				pcount_for_group, FREE_ENTRY_KEY);
		//		}
		//		slls_free(pgroup_by_field_values);
		//		(*pcount_for_group)++;
		//		if (*pcount_for_group <= pstate->head_count) {
		//			outputChannel <- inrecAndContext // xxx stub
		//		}

		outputChannel <- inrecAndContext // xxx stub

	} else {
		outputChannel <- inrecAndContext
	}
}
