#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

#define SORT_BY_KEY     0x10
#define SORT_BY_VALUE   0x20

#define SORT_NUMERIC    0x40
#define SORT_DESCENDING 0x80

typedef struct _mapper_sort_within_records_state_t {
	// Input parameters
	int sort_by;
	int sort_how;
	int reverse;
} mapper_sort_within_records_state_t;

//// Each sort key is string or number; use union to save space.
//typedef struct _typed_sort_within_records_key_t {
//	union {
//		char*  s;
//		double d;
//	} u;
//} typed_sort_within_records_key_t;

// ----------------------------------------------------------------
static void      mapper_sort_within_records_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_sort_within_records_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_sort_within_records_alloc();
static void      mapper_sort_within_records_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_sort_within_records_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

//static typed_sort_within_records_key_t* parse_sort_within_records_keys(slls_t* pkey_field_values, int* sort_params, context_t* pctx);

//// qsort is non-reentrant but qsort_r isn't portable. But since Miller is
//// single-threaded, even if we've got one sort chained to another, only one is
//// active at a time. We adopt the convention that we set the sort params
//// right before the sort.
//static int* pcmp_sort_within_records_params  = NULL;
//static int  cmp_params_length = 0;
//static int pbucket_comparator(const void* pva, const void* pvb);

// ----------------------------------------------------------------
mapper_setup_t mapper_sort_within_records_setup = {
	.verb = "sort-within-records",
	.pusage_func = mapper_sort_within_records_usage,
	.pparse_func = mapper_sort_within_records_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_sort_within_records_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s {flags}\n", argv0, verb);
	fprintf(o, "Flags:\n");
	fprintf(o, "  -k    Sort by keys\n");
	fprintf(o, "  -v    Sort by values\n");
	fprintf(o, "  -n    Sort numerically\n");
	fprintf(o, "  -f    Sort lexically (default)\n");
	fprintf(o, "  -r    Reverse sort\n");
	fprintf(o, "  -nk   Shorthand for -n -k\n");
	fprintf(o, "  -nv   Shorthand for -n -v\n");
	fprintf(o, "  -rk   Shorthand for -r -k\n");
	fprintf(o, "  -rv   Shorthand for -r -v\n");
	fprintf(o, "  -nrk  Shorthand for -n -r -k\n");
	fprintf(o, "  -nrv  Shorthand for -n -r -v\n");
}

static mapper_t* mapper_sort_within_records_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 3) {
		mapper_sort_within_records_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	while ((argc - *pargi) >= 1 && argv[*pargi][0] == '-') {
		*pargi += 1;

// xxx ap_state_t

//		if (streq(flag, "-f")) {
//		} else if (streq(flag, "-n")) {
//		} else if (streq(flag, "-nf")) {
//		} else if (streq(flag, "-r")) {
//		} else if (streq(flag, "-nr")) {
//		} else {
//			mapper_sort_within_records_usage(stderr, argv[0], verb);
//		}
//		slls_t* pnames_for_flag = slls_from_line(value, ',', FALSE);
//		// E.g. with "-nr a,b,c", replicate the "-nr" flag three times.
//		for (sllse_t* pe = pnames_for_flag->phead; pe != NULL; pe = pe->pnext) {
//			slls_append_no_free(pnames, pe->value);
//			slls_append_no_free(pflags, flag);
//		}
//		slls_free(pnames_for_flag);
	}

//	if (pnames->length < 1)
//		mapper_sort_within_records_usage(stderr, argv[0], verb);

	// xxx libify
	// Convert the list such as ["-nf","-nf","-r","-r","-r"] into an array of
	// bit-flags, one per sort-key field.
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

	return mapper_sort_within_records_alloc();
}

// ----------------------------------------------------------------
static mapper_t* mapper_sort_within_records_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sort_within_records_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sort_within_records_state_t));

	//pstate->sort_params = sort_params;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_sort_within_records_process;
	pmapper->pfree_func    = mapper_sort_within_records_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_sort_within_records_free(mapper_t* pmapper, context_t* _) {
	mapper_sort_within_records_state_t* pstate = pmapper->pvstate;
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_sort_within_records_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	//mapper_sort_within_records_state_t* pstate = pvstate;

	// xxx stub
	if (pinrec != NULL)
		return sllv_single(pinrec);
	else
		return sllv_single(NULL);

//	if (pinrec != NULL) {
//		// Consume another input record.
//		slls_t* pkey_field_values = mlr_reference_selected_values_from_record(pinrec, pstate->pkey_field_names);
//		if (pkey_field_values == NULL) {
//			sllv_append(pstate->precords_missing_sort_within_records_keys, pinrec);
//		} else {
//			sort_bucket_t* pbucket = lhmslv_get(pstate->pbuckets_by_key_field_values, pkey_field_values);
//			if (pbucket == NULL) { // New key-field-value: new bucket and hash-map entry
//				slls_t* pkey_field_values_copy = slls_copy(pkey_field_values);
//				sort_bucket_t* pbucket = mlr_malloc_or_die(sizeof(sort_bucket_t));
//				pbucket->typed_sort_within_records_keys = parse_sort_within_records_keys(pkey_field_values_copy, pstate->sort_params, pctx);
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
//		pcmp_sort_within_records_params  = pstate->sort_params;
//		cmp_params_length = pstate->pkey_field_names->length;
//
//		qsort(pbucket_array, num_buckets, sizeof(sort_bucket_t*), pbucket_comparator);
//
//		pcmp_sort_within_records_params  = NULL;
//		cmp_params_length = 0;
//
//		// Emit each bucket's record
//		sllv_t* poutput = sllv_alloc();
//		for (i = 0; i < num_buckets; i++) {
//			sllv_t* plist = pbucket_array[i]->precords;
//			sllv_transfer(poutput, plist);
//			sllv_free(plist);
//		}
//		sllv_transfer(poutput, pstate->precords_missing_sort_within_records_keys);
//		free(pbucket_array);
//		sllv_append(poutput, NULL); // Signal end of output-record stream.
//		return poutput;
//	}

}

//static int pbucket_comparator(const void* pva, const void* pvb) {
//	// We are sorting an array of sort_bucket_t*.
//	const sort_bucket_t** pba = (const sort_bucket_t**)pva;
//	const sort_bucket_t** pbb = (const sort_bucket_t**)pvb;
//	typed_sort_within_records_key_t* akeys = (*pba)->typed_sort_within_records_keys;
//	typed_sort_within_records_key_t* bkeys = (*pbb)->typed_sort_within_records_keys;
//	for (int i = 0; i < cmp_params_length; i++) {
//		int sort_param = pcmp_sort_within_records_params[i];
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

//// E.g. parse the list ["red","1.0"] into the array ["red",1.0].
//static typed_sort_within_records_key_t* parse_sort_within_records_keys(slls_t* pkey_field_values, int* sort_params, context_t* pctx) {
//	typed_sort_within_records_key_t* typed_sort_within_records_keys = mlr_malloc_or_die(pkey_field_values->length * sizeof(typed_sort_within_records_key_t));
//	int i = 0;
//	for (sllse_t* pe = pkey_field_values->phead; pe != NULL; pe = pe->pnext, i++) {
//		if (sort_params[i] & SORT_NUMERIC) {
//			if (*pe->value == 0) { // null input value
//				typed_sort_within_records_keys[i].u.d = nan("");
//			} else if (!mlr_try_float_from_string(pe->value, &typed_sort_within_records_keys[i].u.d)) {
//				fprintf(stderr, "%s: couldn't parse \"%s\" as number in file \"%s\" record %lld.\n",
//					MLR_GLOBALS.bargv0, pe->value, pctx->filename, pctx->fnr);
//				exit(1);
//			}
//		} else {
//			typed_sort_within_records_keys[i].u.s = pe->value;
//		}
//	}
//	return typed_sort_within_records_keys;
//}
