#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

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
//   ane linked list of records.
//
// * Once all the input records are ingested into this hash map, we copy the
//   bucket-pointers into an array and sort it: this being the pairing of
//   parsed-value array and linked list of records. The comparator callback for
//   the sort walks through the parsed-value arrays one slot at a time,
//   looking at the first difference, e.g. if one has "a"="red" and the other
//   has "a"="blue". If the first field matches then the sort moves to the
//   second field, and so on.
//
// * Recall in particular that string keys ["a":"red","x":"1"] and
//   ["a":"red","x":"1.0"] map to different buckets, but will sort equally.
//
// ================================================================

// ----------------------------------------------------------------
#define SORT_NUMERIC    0x80
#define SORT_DESCENDING 0x40

// Each sort key is string or number; use union to save space.
typedef struct _typed_sort_key_t {
	union {
		char*  s;
		double d;
	} u;
} typed_sort_key_t;

typedef struct _bucket_t {
	typed_sort_key_t* typed_sort_keys;
	sllv_t*           precords;
} bucket_t;

// qsort is non-reentrant but qsort_r isn't portable. But since Miller is
// single-threaded, even if we've got one sort chained to another, only one is
// active at a time. We adopt the convention that we set the sort params
// right before the sort.
static int* pcmp_sort_params  = NULL;
static int  cmp_params_length = 0;
static int pbucket_comparator(const void* pva, const void* pvb) {
	// We are sorting an array of bucket_t*.
	const bucket_t** pba = (const bucket_t**)pva;
	const bucket_t** pbb = (const bucket_t**)pvb;
	typed_sort_key_t* akeys = (*pba)->typed_sort_keys;
	typed_sort_key_t* bkeys = (*pbb)->typed_sort_keys;
	for (int i = 0; i < cmp_params_length; i++) {
		int sort_param = pcmp_sort_params[i];
		if (sort_param & SORT_NUMERIC) {
			double d = akeys[i].u.d - bkeys[i].u.d;
			int s = (d < 0) ? -1 : (d > 0) ? 1 : 0;
			if (s != 0)
				return (sort_param & SORT_DESCENDING) ? -s : s;
		} else {
			int s = strcmp(akeys[i].u.s, bkeys[i].u.s);
			if (s != 0)
				return (sort_param & SORT_DESCENDING) ? -s : s;
		}
	}
	return 0;
}

// E.g. parse the list ["red","1.0"] into the array ["red",1.0].
static typed_sort_key_t* parse_sort_keys(slls_t* pkey_field_values, int* sort_params) {
	typed_sort_key_t* typed_sort_keys = mlr_malloc_or_die(pkey_field_values->length * sizeof(typed_sort_key_t));
	int i = 0;
	for (sllse_t* pe = pkey_field_values->phead; pe != NULL; pe = pe->pnext, i++) {
		if (sort_params[i] & SORT_NUMERIC) {
			if (sscanf(pe->value, "%lf", &typed_sort_keys[i].u.d) != 1) {
				// xxx to do: print some more context here, e.g. file name & line number
				fprintf(stderr, "Couldn't parse \"%s\" as number.\n", pe->value);
				exit(1);
			}
		} else {
			typed_sort_keys[i].u.s = pe->value;
		}
	}
	return typed_sort_keys;
}

// ----------------------------------------------------------------
typedef struct _mapper_sort_state_t {
	// Input parameters
	slls_t* pkey_field_names; // Fields to sort on
	int*    sort_params;      // Lexical/numeric; ascending/descending
	int do_sort;              // If false, just do group-by
	// Sort state: buckets of like records.
	lhmslv_t* pbuckets_by_key_field_names;
} mapper_sort_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_sort_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_sort_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		// Consume another input record.
		slls_t* pkey_field_values = mlr_selected_values_from_record(pinrec, pstate->pkey_field_names);
		bucket_t* pbucket = lhmslv_get(pstate->pbuckets_by_key_field_names, pkey_field_values);
		if (pbucket == NULL) { // New key-field-value: new bucket and hash-map entry
			slls_t* pkey_field_values_copy = slls_copy(pkey_field_values);
			bucket_t* pbucket = mlr_malloc_or_die(sizeof(bucket_t)); // xxx free in the free func
			pbucket->typed_sort_keys = parse_sort_keys(pkey_field_values_copy, pstate->sort_params);
			pbucket->precords = sllv_alloc();
			sllv_add(pbucket->precords, pinrec);
			lhmslv_put(pstate->pbuckets_by_key_field_names, pkey_field_values_copy, pbucket);
		} else { // Previously seen key-field-value: append record to bucket
			sllv_add(pbucket->precords, pinrec);
		}
		return NULL;
	}
	else if (!pstate->do_sort) {
		// Group-by
		sllv_t* poutput = sllv_alloc();
		for (lhmslve_t* pe = pstate->pbuckets_by_key_field_names->phead; pe != NULL; pe = pe->pnext) {
			bucket_t* pbucket = pe->pvvalue;
			for (sllve_t* pf = pbucket->precords->phead; pf != NULL; pf = pf->pnext) {
				sllv_add(poutput, pf->pvdata);
			}
		}
		sllv_add(poutput, NULL);
		return poutput;
	} else {
		// End of input stream: sort bucket labels
		int num_buckets = pstate->pbuckets_by_key_field_names->num_occupied;
		bucket_t** pbucket_array = mlr_malloc_or_die(num_buckets * sizeof(bucket_t*));

		// Copy bucket-pointers to an array for qsort
		int i = 0;
		for (lhmslve_t* pe = pstate->pbuckets_by_key_field_names->phead; pe != NULL; pe = pe->pnext, i++) {
			pbucket_array[i] = pe->pvvalue;
		}

		pcmp_sort_params  = pstate->sort_params;
		cmp_params_length = pstate->pkey_field_names->length;

		// xxx to do: allow the user to specify the sort algorithm from the CLI, assuming portability
		qsort(pbucket_array, num_buckets, sizeof(bucket_t*), pbucket_comparator);
		//heapsort(pbucket_array, num_buckets, sizeof(bucket_t*), pbucket_comparator);
		//mergesort(pbucket_array, num_buckets, sizeof(bucket_t*), pbucket_comparator);

		pcmp_sort_params  = NULL;
		cmp_params_length = 0;

		// Emit each bucket's record
		sllv_t* poutput = sllv_alloc();
		for (i = 0; i < num_buckets; i++) {
			sllv_t* plist = pbucket_array[i]->precords;
			for (sllve_t* pf = plist->phead; pf != NULL; pf = pf->pnext) {
				sllv_add(poutput, pf->pvdata);
			}
		}
		free(pbucket_array);
		sllv_add(poutput, NULL); // Signal end of output-record stream.
		return poutput;
	}
}

// ----------------------------------------------------------------
// xxx to do: at end of sort, free and nullify the things populated within
// mapper_sort_process.  here, free the things allocated as setup time.
static void mapper_sort_free(void* pvstate) {
	mapper_sort_state_t* pstate = pvstate;
	if (pstate->pkey_field_names != NULL)
		slls_free(pstate->pkey_field_names);
	if (pstate->pbuckets_by_key_field_names != NULL)
		// xxx free void-star payloads 1st
		lhmslv_free(pstate->pbuckets_by_key_field_names);
	if (pstate->sort_params != NULL)
		free(pstate->sort_params);
}

static mapper_t* mapper_sort_alloc(slls_t* pkey_field_names, int* sort_params, int do_sort) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sort_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sort_state_t));

	pstate->pkey_field_names            = pkey_field_names;
	pstate->sort_params                 = sort_params;
	pstate->pbuckets_by_key_field_names = lhmslv_alloc();
	pstate->do_sort                     = do_sort;

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_sort_process;
	pmapper->pmapper_free_func    = mapper_sort_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_group_by_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s {comma-separated field names}\n", argv0, verb);
	fprintf(stdout, "Outputs records in batches having identical values at specified field names.\n");
}

static mapper_t* mapper_group_by_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 2) {
		mapper_group_by_usage(argv[0], argv[*pargi]);
		return NULL;
	}

	slls_t* pnames = slls_from_line(argv[*pargi+1], ',', FALSE);
	int* opt_array = mlr_malloc_or_die(pnames->length * sizeof(int));
	for (int i = 0; i < pnames->length; i++)
		opt_array[i] = 0;

	*pargi += 2;
	return mapper_sort_alloc(pnames, opt_array, FALSE);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_group_by_setup = {
	.verb = "group-by",
	.pusage_func = mapper_group_by_usage,
	.pparse_func = mapper_group_by_parse_cli
};

// ----------------------------------------------------------------
static void mapper_sort_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s {flags}\n", argv0, verb);
	fprintf(stdout, "Flags:\n");
	fprintf(stdout, "  -f  {comma-separated field names}  Lexical ascending\n");
	fprintf(stdout, "  -n  {comma-separated field names}  Numerical ascending\n");
	fprintf(stdout, "  -nf {comma-separated field names}  Numerical ascending\n");
	fprintf(stdout, "  -r  {comma-separated field names}  Lexical descending\n");
	fprintf(stdout, "  -nr {comma-separated field names}  Numerical descending\n");
	fprintf(stdout, "Sorts records primarily by the first specified field, secondarily by the second field, and so on.\n");
	fprintf(stdout, "Example:\n");
	fprintf(stdout, "  %s %s -f a,b -nr x,y,z\n", argv0, verb);
	fprintf(stdout, "which is the same as:\n");
	fprintf(stdout, "  %s %s -f a -f b -nr x -nr y -nr z\n", argv0, verb);
}

static mapper_t* mapper_sort_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 3) {
		mapper_sort_usage(argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;
	slls_t* pnames = slls_alloc();
	slls_t* pflags = slls_alloc();

	while ((argc - *pargi) >= 1 && argv[*pargi][0] == '-') {
		if ((argc - *pargi) < 2)
			mapper_sort_usage(argv[0], verb);
		char* flag  = argv[*pargi];
		char* value = argv[*pargi+1];
		*pargi += 2;

		if (streq(flag, "-f")) {
		} else if (streq(flag, "-n")) {
		} else if (streq(flag, "-nf")) {
		} else if (streq(flag, "-r")) {
		} else if (streq(flag, "-nr")) {
		} else {
			mapper_sort_usage(argv[0], verb);
		}
		slls_t* pnames_for_flag = slls_from_line(value, ',', FALSE);
		// E.g. with "-nr a,b,c", replicate the "-nr" flag three times.
		for (sllse_t* pe = pnames_for_flag->phead; pe != NULL; pe = pe->pnext) {
			slls_add_no_free(pnames, pe->value);
			slls_add_no_free(pflags, flag);
		}
		slls_free(pnames_for_flag);
	}

	if (pnames->length < 1)
		mapper_sort_usage(argv[0], verb);

	// Convert the list such as ["-nf","-nf","-r","-r","-r"] into an array of
	// bit-flags, one per sort-key field.
	int* opt_array = mlr_malloc_or_die(pnames->length * sizeof(int));
	sllse_t* pe;
	int di;
	for (pe = pflags->phead, di = 0; pe != NULL; pe = pe->pnext, di++) {
		char* flag = pe->value;
		int opt =
			streq(flag, "-nf") ? SORT_NUMERIC :
			streq(flag, "-n")  ? SORT_NUMERIC :
			streq(flag, "-r")  ? SORT_DESCENDING :
			streq(flag, "-nr") ? SORT_NUMERIC|SORT_DESCENDING :
			0;
		opt_array[di] =opt;
	}
	slls_free(pflags);

	return mapper_sort_alloc(pnames, opt_array, TRUE);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_sort_setup = {
	.verb = "sort",
	.pusage_func = mapper_sort_usage,
	.pparse_func = mapper_sort_parse_cli
};
