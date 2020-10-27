#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/lhmsv.h"
#include "mapping/mappers.h"

#define SORT_BY_KEY     0x10
#define SORT_BY_VALUE   0x20

#define SORT_NUMERIC    0x40
#define SORT_DESCENDING 0x80

typedef struct _mapper_sort_within_records_state_t {
	// Input parameters
	// TODO:
	//int sort_by;
	//int sort_how;
	//int reverse;
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
	fprintf(o, "Usage: %s %s [no options]\n", argv0, verb);
	fprintf(o, "Outputs records sorted lexically ascending by keys.\n");
	// TODO
	//fprintf(o, "Usage: %s %s {flags}\n", argv0, verb);
	//fprintf(o, "Flags:\n");
	//fprintf(o, "  -k    Sort by keys\n");
	//fprintf(o, "  -v    Sort by values\n");
	//fprintf(o, "  -n    Sort numerically\n");
	//fprintf(o, "  -f    Sort lexically (default)\n");
	//fprintf(o, "  -r    Reverse sort\n");
	//fprintf(o, "  -nk   Shorthand for -n -k\n");
	//fprintf(o, "  -nv   Shorthand for -n -v\n");
	//fprintf(o, "  -rk   Shorthand for -r -k\n");
	//fprintf(o, "  -rv   Shorthand for -r -v\n");
	//fprintf(o, "  -nrk  Shorthand for -n -r -k\n");
	//fprintf(o, "  -nrv  Shorthand for -n -r -v\n");
}

static mapper_t* mapper_sort_within_records_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_sort_within_records_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	*pargi += 1;

	return mapper_sort_within_records_alloc();
}

// ----------------------------------------------------------------
static mapper_t* mapper_sort_within_records_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sort_within_records_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sort_within_records_state_t));

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
// For qsort
static int key_comparator(const void* pva, const void* pvb) {
    const char *pa = *(const char**)pva;
    const char *pb = *(const char**)pvb;
	return strcmp(pa, pb);
}

// ----------------------------------------------------------------
static sllv_t* mapper_sort_within_records_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {

	if (pinrec == NULL) { // end of record stream
		return sllv_single(NULL);
	}

	// Get keys as array
	lhmsv_t* keys_to_entries = lhmsv_alloc();
	int n = pinrec->field_count;
	char** keys_array = mlr_malloc_or_die(n * sizeof(char*));
	int i = 0;
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		keys_array[i] = pe->key;
		lhmsv_put(keys_to_entries, pe->key, pe, NO_FREE);
		i++;
	}

	// Sort the keys
	qsort(keys_array, n, sizeof(char*), key_comparator);

	// Make a new record
	lrec_t* poutrec = lrec_unbacked_alloc();

	for (i = 0; i < n; i++) {
		lrece_t* pe = lhmsv_get(keys_to_entries, keys_array[i]);
		lrec_put(
			poutrec,
			mlr_strdup_or_die(pe->key),
			mlr_strdup_or_die(pe->value),
			FREE_ENTRY_KEY | FREE_ENTRY_VALUE
		);
	}
	lrec_free(pinrec);
	lhmsv_free(keys_to_entries);
	free(keys_array);

	return sllv_single(poutrec);
}
