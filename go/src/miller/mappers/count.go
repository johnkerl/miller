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
var CountSetup = mapping.MapperSetup{
	Verb:         "count",
	ParseCLIFunc: mapperCountParseCLI,
	IgnoresInput: false,
}

func mapperCountParseCLI(
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
		mapperCountUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentionally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	mapper, _ := NewMapperCount(
		*pGroupByFieldNames,
		*pShowCountsOnly,
		*pOutputFieldName,
	)

	*pargi = argi
	return mapper
}

func mapperCountUsage(
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
type MapperCount struct {
	// input
	groupByFieldNameList []string
	showCountsOnly       bool
	outputFieldName      string

	// state
	recordMapperFunc mapping.RecordMapperFunc
	recordListsByGroup *lib.OrderedMap
	ungroupedCount int64
}

func NewMapperCount(
	groupByFieldNames string,
	showCountsOnly bool,
	outputFieldName string,
) (*MapperCount, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &MapperCount{
		groupByFieldNameList: groupByFieldNameList,
		showCountsOnly:       showCountsOnly,
		outputFieldName:      outputFieldName,

		recordListsByGroup: lib.NewOrderedMap(),
		ungroupedCount: 0,
	}

	if len(groupByFieldNameList) == 0 {
		this.recordMapperFunc = this.countUngrouped
	} else {
		this.recordMapperFunc = this.countGrouped
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperCount) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	this.recordMapperFunc(inrecAndContext, outputChannel)
}

// ----------------------------------------------------------------
func (this *MapperCount) countUngrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream
		this.ungroupedCount++;
	} else {
		newrec := types.NewMlrmapAsRecord()
		mcount := types.MlrvalFromInt64(this.ungroupedCount)
		newrec.PutCopy(&this.outputFieldName, &mcount)

		outputChannel <- types.NewRecordAndContext(newrec, &inrecAndContext.Context)

		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

//		return NULL;
//	} else { // end of record stream
//		lrec_t* poutrec = lrec_unbacked_alloc();
//		lrec_put(poutrec, pstate->output_field_name,
//			mlr_alloc_string_from_ll(pstate->ungrouped_count), FREE_ENTRY_VALUE);
//		return sllv_single(poutrec);
//	}
//}

// ----------------------------------------------------------------
func (this *MapperCount) countGrouped(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil { // not end of record stream

	} else {
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}



//// ----------------------------------------------------------------
//static sllv_t* mapper_count_process_grouped(
//	lrec_t* pinrec,
//	context_t* pctx,
//	void* pvstate)
//{
//	mapper_count_state_t* pstate = pvstate;
//	if (pinrec != NULL) { // not end of record stream
//		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
//			pstate->pgroup_by_field_names);
//		if (pgroup_by_field_values == NULL) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//
//		unsigned long long* pcount = lhmslv_get(pstate->pcounts_by_group, pgroup_by_field_values);
//		if (pcount == NULL) {
//			pcount = mlr_malloc_or_die(sizeof(unsigned long long));
//			*pcount = 1LL;
//			slls_t* pcopy = slls_copy(pgroup_by_field_values);
//			lhmslv_put(pstate->pcounts_by_group, pcopy, pcount, FREE_ENTRY_KEY);
//			lrec_free(pinrec);
//		} else {
//			(*pcount)++;
//			lrec_free(pinrec);
//		}
//
//		return NULL;
//
//	} else { // end of record stream
//		sllv_t* poutrecs = sllv_alloc();
//
//		if (pstate->show_counts_only) {
//			lrec_t* poutrec = lrec_unbacked_alloc();
//
//			unsigned long long count = (unsigned long long)lhmslv_size(pstate->pcounts_by_group);
//			lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(count),
//				FREE_ENTRY_VALUE);
//
//			sllv_append(poutrecs, poutrec);
//		} else {
//			for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
//				lrec_t* poutrec = lrec_unbacked_alloc();
//
//				slls_t* pgroup_by_field_values = pa->key;
//
//				sllse_t* pb = pstate->pgroup_by_field_names->phead;
//				sllse_t* pc =         pgroup_by_field_values->phead;
//				for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
//					lrec_put(poutrec, pb->value, pc->value, NO_FREE);
//				}
//
//				unsigned long long* pcount = pa->pvvalue;
//				lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(*pcount),
//					FREE_ENTRY_VALUE);
//
//				sllv_append(poutrecs, poutrec);
//			}
//		}
//
//		sllv_append(poutrecs, NULL);
//		return poutrecs;
//
//	}
//}

//	if x {
//		//		groupingKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
//		//		if !ok {
//		//			return
//		//		}
//		//
//		//		irecordListForGroup := this.recordListsByGroup.Get(groupingKey)
//		//		if irecordListForGroup == nil { // first time
//		//			irecordListForGroup = list.New()
//		//			this.recordListsByGroup.Put(groupingKey, irecordListForGroup)
//		//		}
//		//		recordListForGroup := irecordListForGroup.(*list.List)
//		//
//		//		recordListForGroup.PushBack(inrecAndContext)
//		//		for uint64(recordListForGroup.Len()) > this.countCount {
//		//			recordListForGroup.Remove(recordListForGroup.Front())
//		//		}
//	} else {
//		//		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
//		//			recordListForGroup := outer.Value.(*list.List)
//		//			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
//		//				outputChannel <- inner.Value.(*types.RecordAndContext)
//		//			}
//		//		}
//		outputChannel <- inrecAndContext // end-of-stream marker
//	}
