package mappers

// ================================================================
// OVERVIEW
//
// * Suppose we are sorting records lexically ascending on field "a" and then
//   numerically descending on field "x".
//
// * CLI syntax is "mlr sort -f a -nr x".
//
// * We first consume all input records and for each extract the string values
//   of fields a and x. For each uniq combination of a-value (e.g. "red",
//   "green", "blue") and x-value (e.g. "1", "1.0", "2.4") -- e.g.
//   pairs ["red","1"], and so on -- we keep a linked list of all the records
//   having those sort-key values, in the order encountered.
//
// * For each of those unique sort-key-value combinations, we also parse the
//   numerical fields at this point into an array of union-string-double.
//   E.g. the list ["red", "1.0"] maps to the array ["red", 1.0].
//
// * The pairing of parsed-value array the linked list of same-key-value records
//   is called a *bucket*. E.g the records
//     {"a":"red","b":"circle","x":"1.0","y":"3.9"}
//     {"a":"red","b":"square","x":"1.0","z":"5.7", "q":"even"}
//   would both land in the ["red","1.0"] bucket.
//
// * Buckets are retained in a hash map: the key is the string-list of the form
//   ["red","1.0"] and the value is the pairing of parsed-value array ["red",1.0]
//   and linked list of records.
//
// * Once all the input records are ingested into this hash map, we copy the
//   bucket-pointers into an array and sort it: this being the pairing of
//   parsed-value array and linked list of records. The comparator callback for
//   the sort walks through the parsed-value arrays one slot at a time,
//   looking at the first difference, e.g. if one has "a"="red" and the other
//   has "a"="blue". If the first field matches then the sort moves to the
//   second field, and so on.
//
// * Note in particular that string keys ["a":"red","x":"1"] and
//   ["a":"red","x":"1.0"] map to different buckets, but will sort equally.
//
// ================================================================

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
var SortSetup = mapping.MapperSetup{
	Verb:         "sort",
	ParseCLIFunc: mapperSortParseCLI,
	IgnoresInput: false,
}

func mapperSortParseCLI(
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
		mapperSortUsage(ostream, args[0], verb, flagSet)
	}
	flagSet.Parse(args[argi:])
	if errorHandling == flag.ContinueOnError { // help intentioally requested
		return nil
	}

	// Find out how many flags were consumed by this verb and advance for the
	// next verb
	argi = len(args) - len(flagSet.Args())

	// Get the sort field names from the command line
	if argi >= argc {
		flagSet.Usage()
		os.Exit(1)
	}
	groupByFieldNames := args[argi]
	argi += 1

	mapper, _ := NewMapperSort(
		groupByFieldNames,
	)

	*pargi = argi
	return mapper
}

func mapperSortUsage(
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
// Example:
// * mlr sort -f a -n i
// * group-by field-name list is "a,i"
// * input record 'a=pan,b=pan,i=1,x=0.3467,y=0.7268'
//   o values at a,i are "pan",1
//   o grouping key for the ordered map from string to record-bucket is the string "pan,1"
//   o we also need a map from "pan,1" to the array of mlrvals ["pan", 1].
// * next input record 'a=eks,b=pan,i=2,x=0.7586,y=0.5221'
//   o values at a,i are "eks",2
//   o grouping key for the ordered map from string to record-bucket is the string "eks,2"
//   o we also need a map from "eks,2" to the array of mlrvals ["eks", 2].
// * what gets sorted are the bucket-heading arrays of mlrvals:
//   o make an array [ ("pan,1", ["pan", 1]), ("eks,2", ["eks", 2])
//   o sort that
// * output is simply for each slot in the array, emit each record in the bucket

type MapperSort struct {
	// Input
	groupByFieldNameList []string

	// State
	// Map from string to *list.List:
	recordListsByGroup *lib.OrderedMap
	// Map from string to []lib.Mlrval:
	bucketHeads *lib.OrderedMap
}

func NewMapperSort(
	groupByFieldNames string,
) (*MapperSort, error) {

	groupByFieldNameList := lib.SplitString(groupByFieldNames, ",")

	this := &MapperSort{
		groupByFieldNameList: groupByFieldNameList,

		recordListsByGroup: lib.NewOrderedMap(),
		bucketHeads:        lib.NewOrderedMap(),
	}

	return this, nil
}

// ----------------------------------------------------------------
func (this *MapperSort) Map(
	inrecAndContext *types.RecordAndContext,
	outputChannel chan<- *types.RecordAndContext,
) {
	inrec := inrecAndContext.Record
	if inrec != nil {
		// Not end of record stream

		// xxx need to factor out the joined and non-joined -- this verb needs both
		groupByKey, ok := inrec.GetSelectedValuesJoined(this.groupByFieldNameList)
		if !ok {
			return
		}

		recordListForGroup := this.recordListsByGroup.Get(groupByKey)
		if recordListForGroup == nil {
			recordListForGroup = list.New()
			this.recordListsByGroup.Put(groupByKey, recordListForGroup)
			this.bucketHeads.Put(groupByKey, recordListForGroup)
		}

		recordListForGroup.(*list.List).PushBack(inrecAndContext)
		// xxx put the bucket-heads data
		// bucketHeads.(*list.List).PushBack(xxx)

	} else {
		// End of record stream

		// xxx stub from group-by (no sorting):
		for outer := this.recordListsByGroup.Head; outer != nil; outer = outer.Next {
			recordListForGroup := outer.Value.(*list.List)
			for inner := recordListForGroup.Front(); inner != nil; inner = inner.Next() {
				outputChannel <- inner.Value.(*types.RecordAndContext)
			}
		}
		outputChannel <- inrecAndContext // end-of-stream marker
	}
}

//#define SORT_NUMERIC    0x80
//#define SORT_DESCENDING 0x40

//typedef struct _mapper_sort_state_t {
//	// Input parameters
//	slls_t* pkey_field_names; // Fields to sort on
//	int*    sort_params;      // Lexical/numeric; ascending/descending
//	// Sort state: buckets of like records.
//	lhmslv_t* pbuckets_by_key_field_values;
//	sllv_t*   precords_missing_sort_keys;
//} mapper_sort_state_t;

//typedef struct _sort_bucket_t {
//	typed_sort_key_t* typed_sort_keys;
//	sllv_t*           precords;
//} sort_bucket_t;

//// ----------------------------------------------------------------
//static void mapper_sort_usage(FILE* o, char* argv0, char* verb) {
//	fprintf(o, "Usage: %s %s {flags}\n", argv0, verb);
//	fprintf(o, "Flags:\n");
//	fprintf(o, "  -f  {comma-separated field names}  Lexical ascending\n");
//	fprintf(o, "  -n  {comma-separated field names}  Numerical ascending; nulls sort last\n");
//	fprintf(o, "  -nf {comma-separated field names}  Same as -n\n");
//	fprintf(o, "  -r  {comma-separated field names}  Lexical descending\n");
//	fprintf(o, "  -nr {comma-separated field names}  Numerical descending; nulls sort first\n");
//	fprintf(o, "Sorts records primarily by the first specified field, secondarily by the second\n");
//	fprintf(o, "field, and so on.  (Any records not having all specified sort keys will appear\n");
//	fprintf(o, "at the end of the output, in the order they were encountered, regardless of the\n");
//	fprintf(o, "specified sort order.) The sort is stable: records that compare equal will sort\n");
//	fprintf(o, "in the order they were encountered in the input record stream.\n");
//	fprintf(o, "\n");
//	fprintf(o, "Example:\n");
//	fprintf(o, "  %s %s -f a,b -nr x,y,z\n", argv0, verb);
//	fprintf(o, "which is the same as:\n");
//	fprintf(o, "  %s %s -f a -f b -nr x -nr y -nr z\n", argv0, verb);
//}

//static mapper_t* mapper_sort_parse_cli(int* pargi, int argc, char** argv,
//	cli_reader_opts_t* _, cli_writer_opts_t* __)
//{
//	if ((argc - *pargi) < 3) {
//		mapper_sort_usage(stderr, argv[0], argv[*pargi]);
//		return NULL;
//	}
//	char* verb = argv[*pargi];
//	*pargi += 1;
//	slls_t* pnames = slls_alloc();
//	slls_t* pflags = slls_alloc();
//
//	while ((argc - *pargi) >= 1 && argv[*pargi][0] == '-') {
//		if ((argc - *pargi) < 2)
//			mapper_sort_usage(stderr, argv[0], verb);
//		char* flag  = argv[*pargi];
//		char* value = argv[*pargi+1];
//		*pargi += 2;
//
//		if (streq(flag, "-f")) {
//		} else if (streq(flag, "-n")) {
//		} else if (streq(flag, "-nf")) {
//		} else if (streq(flag, "-r")) {
//		} else if (streq(flag, "-nr")) {
//		} else {
//			mapper_sort_usage(stderr, argv[0], verb);
//		}
//		slls_t* pnames_for_flag = slls_from_line(value, ',', FALSE);
//		// E.g. with "-nr a,b,c", replicate the "-nr" flag three times.
//		for (sllse_t* pe = pnames_for_flag->phead; pe != NULL; pe = pe->pnext) {
//			slls_append_no_free(pnames, pe->value);
//			slls_append_no_free(pflags, flag);
//		}
//		slls_free(pnames_for_flag);
//	}
//
//	if (pnames->length < 1)
//		mapper_sort_usage(stderr, argv[0], verb);
//
//	// Convert the list such as ["-nf","-nf","-r","-r","-r"] into an array of
//	// bit-flags, one per sort-key field.
//	int* opt_array = mlr_malloc_or_die(pnames->length * sizeof(int));
//	sllse_t* pe;
//	int di;
//	for (pe = pflags->phead, di = 0; pe != NULL; pe = pe->pnext, di++) {
//		char* flag = pe->value;
//		int opt =
//			streq(flag, "-nf") ? SORT_NUMERIC :
//			streq(flag, "-n")  ? SORT_NUMERIC :
//			streq(flag, "-r")  ? SORT_DESCENDING :
//			streq(flag, "-nr") ? SORT_NUMERIC|SORT_DESCENDING :
//			0;
//		opt_array[di] =opt;
//	}
//	slls_free(pflags);
//
//	return mapper_sort_alloc(pnames, opt_array, TRUE);
//}

//// ----------------------------------------------------------------
//static mapper_t* mapper_sort_alloc(slls_t* pkey_field_names, int* sort_params, int do_sort) {
//	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
//
//	mapper_sort_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sort_state_t));
//
//	pstate->pkey_field_names             = pkey_field_names;
//	pstate->sort_params                  = sort_params;
//	pstate->pbuckets_by_key_field_values = lhmslv_alloc();
//	pstate->precords_missing_sort_keys   = sllv_alloc();
//	pstate->do_sort                      = do_sort;
//
//	pmapper->pvstate       = pstate;
//	pmapper->pprocess_func = mapper_sort_process;
//	pmapper->pfree_func    = mapper_sort_free;
//
//	return pmapper;
//}
//
//// ----------------------------------------------------------------
//static void mapper_sort_free(mapper_t* pmapper, context_t* _) {
//	mapper_sort_state_t* pstate = pmapper->pvstate;
//	if (pstate->pkey_field_names != NULL)
//		slls_free(pstate->pkey_field_names);
//	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
//	for (lhmslve_t* pa = pstate->pbuckets_by_key_field_values->phead; pa != NULL; pa = pa->pnext) {
//		sort_bucket_t* pbucket = pa->pvvalue;
//		free(pbucket->typed_sort_keys);
//		free(pbucket);
//		// precords freed in emitter
//	}
//	lhmslv_free(pstate->pbuckets_by_key_field_values);
//	sllv_free(pstate->precords_missing_sort_keys);
//	free(pstate->sort_params);
//	free(pstate);
//	free(pmapper);
//}
//
//// ----------------------------------------------------------------
//static sllv_t* mapper_sort_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	mapper_sort_state_t* pstate = pvstate;
//	if (pinrec != NULL) {
//		// Consume another input record.
//		slls_t* pkey_field_values = mlr_reference_selected_values_from_record(pinrec, pstate->pkey_field_names);
//		if (pkey_field_values == NULL) {
//			sllv_append(pstate->precords_missing_sort_keys, pinrec);
//		} else {
//			sort_bucket_t* pbucket = lhmslv_get(pstate->pbuckets_by_key_field_values, pkey_field_values);
//			if (pbucket == NULL) { // New key-field-value: new bucket and hash-map entry
//				slls_t* pkey_field_values_copy = slls_copy(pkey_field_values);
//				sort_bucket_t* pbucket = mlr_malloc_or_die(sizeof(sort_bucket_t));
//				pbucket->typed_sort_keys = parse_sort_keys(pkey_field_values_copy, pstate->sort_params, pctx);
//				pbucket->precords = sllv_alloc();
//				sllv_append(pbucket->precords, pinrec);
//				lhmslv_put(pstate->pbuckets_by_key_field_values, pkey_field_values_copy, pbucket,
//					FREE_ENTRY_KEY);
//			} else { // Previously seen key-field-value: append record to bucket
//				sllv_append(pbucket->precords, pinrec);
//			}
//			slls_free(pkey_field_values);
//		}
//		return NULL;
//	} else {
//		// End of input stream: sort bucket labels
//		int num_buckets = pstate->pbuckets_by_key_field_values->num_occupied;
//		sort_bucket_t** pbucket_array = mlr_malloc_or_die(num_buckets * sizeof(sort_bucket_t*));
//
//		// Copy bucket-pointers to an array for qsort
//		int i = 0;
//		for (lhmslve_t* pe = pstate->pbuckets_by_key_field_values->phead; pe != NULL; pe = pe->pnext, i++) {
//			pbucket_array[i] = pe->pvvalue;
//		}
//
//		pcmp_sort_params  = pstate->sort_params;
//		cmp_params_length = pstate->pkey_field_names->length;
//
//		qsort(pbucket_array, num_buckets, sizeof(sort_bucket_t*), pbucket_comparator);
//
//		pcmp_sort_params  = NULL;
//		cmp_params_length = 0;
//
//		// Emit each bucket's record
//		sllv_t* poutput = sllv_alloc();
//		for (i = 0; i < num_buckets; i++) {
//			sllv_t* plist = pbucket_array[i]->precords;
//			sllv_transfer(poutput, plist);
//			sllv_free(plist);
//		}
//		sllv_transfer(poutput, pstate->precords_missing_sort_keys);
//		free(pbucket_array);
//		sllv_append(poutput, NULL); // Signal end of output-record stream.
//		return poutput;
//	}
//}
//
//static int pbucket_comparator(const void* pva, const void* pvb) {
//	// We are sorting an array of sort_bucket_t*.
//	const sort_bucket_t** pba = (const sort_bucket_t**)pva;
//	const sort_bucket_t** pbb = (const sort_bucket_t**)pvb;
//	typed_sort_key_t* akeys = (*pba)->typed_sort_keys;
//	typed_sort_key_t* bkeys = (*pbb)->typed_sort_keys;
//	for (int i = 0; i < cmp_params_length; i++) {
//		int sort_param = pcmp_sort_params[i];
//		if (sort_param & SORT_NUMERIC) {
//			double a = akeys[i].u.d;
//			double b = bkeys[i].u.d;
//			if (isnan(a)) { // null input value
//				if (!isnan(b)) {
//					return (sort_param & SORT_DESCENDING) ? -1 : 1;
//				}
//			} else if (isnan(b)) {
//					return (sort_param & SORT_DESCENDING) ? 1 : -1;
//			} else {
//				double d = a - b;
//				int s = (d < 0) ? -1 : (d > 0) ? 1 : 0;
//				if (s != 0)
//					return (sort_param & SORT_DESCENDING) ? -s : s;
//			}
//		} else {
//			int s = strcmp(akeys[i].u.s, bkeys[i].u.s);
//			if (s != 0)
//				return (sort_param & SORT_DESCENDING) ? -s : s;
//		}
//	}
//	return 0;
//}
//
//// E.g. parse the list ["red","1.0"] into the array ["red",1.0].
//static typed_sort_key_t* parse_sort_keys(slls_t* pkey_field_values, int* sort_params, context_t* pctx) {
//	typed_sort_key_t* typed_sort_keys = mlr_malloc_or_die(pkey_field_values->length * sizeof(typed_sort_key_t));
//	int i = 0;
//	for (sllse_t* pe = pkey_field_values->phead; pe != NULL; pe = pe->pnext, i++) {
//		if (sort_params[i] & SORT_NUMERIC) {
//			if (*pe->value == 0) { // null input value
//				typed_sort_keys[i].u.d = nan("");
//			} else if (!mlr_try_float_from_string(pe->value, &typed_sort_keys[i].u.d)) {
//				fprintf(stderr, "%s: couldn't parse \"%s\" as number in file \"%s\" record %lld.\n",
//					MLR_GLOBALS.bargv0, pe->value, pctx->filename, pctx->fnr);
//				exit(1);
//			}
//		} else {
//			typed_sort_keys[i].u.s = pe->value;
//		}
//	}
//	return typed_sort_keys;
//}
