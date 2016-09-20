#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DEFAULT_MAX_OUTPUT_LENGTH 10LL

typedef struct _mapper_most_or_least_state_t {
	ap_state_t* pargp;
	slls_t*     pgroup_by_field_names;
	lhmslv_t*   pcounts_by_group;
	long long   max_output_length;
	int         descending;
	int         show_counts;
} mapper_most_or_least_state_t;

static void mapper_most_usage(FILE*  o, char* argv0, char* verb);
static void mapper_least_usage(FILE* o, char* argv0, char* verb);

static mapper_t* mapper_most_parse_cli(int*  pargi, int argc, char** argv, cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_least_parse_cli(int* pargi, int argc, char** argv, cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_most_or_least_parse_cli(int* pargi, int argc, char** argv, int descending);

static mapper_t* mapper_most_or_least_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	long long max_output_length, int descending, int show_counts);
static void      mapper_most_or_least_free(mapper_t* pmapper);

static sllv_t*   mapper_most_or_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// qsort callbacks
static int descending_vcmp(const void* pva, const void* pvb);
static int ascending_vcmp(const void* pva, const void* pvb);

typedef struct _sort_pair_t {
	slls_t* pgroup_by_field_values;
	long long count; // signed, not unsigned, for sort-cmp callback
} sort_pair_t;

// ----------------------------------------------------------------
mapper_setup_t mapper_most_setup = {
	.verb          = "most",
	.pusage_func   = mapper_most_usage,
	.pparse_func   = mapper_most_parse_cli,
	.ignores_input = FALSE,
};

mapper_setup_t mapper_least_setup = {
	.verb          = "least",
	.pusage_func   = mapper_least_usage,
	.pparse_func   = mapper_least_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_most_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Shows the most frequently occurring distinct values for specified field names.\n");
	fprintf(o, "The first entry is the statistical mode; the remaining are runners-up.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-f {one or more comma-separated field names}. Required flag.\n");
	fprintf(o, "-n {count}. Optional flag defaulting to %lld.\n", DEFAULT_MAX_OUTPUT_LENGTH);
	fprintf(o, "-b          Suppress counts; show only field values.\n");
	fprintf(o, "See also \"%s %s\".\n", argv0, "least");
}

static void mapper_least_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Shows the least frequently occurring distinct values for specified field names.\n");
	fprintf(o, "The first entry is the statistical anti-mode; the remaining are runners-up.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-f {one or more comma-separated field names}. Required flag.\n");
	fprintf(o, "-n {count}. Optional flag defaulting to %lld.\n", DEFAULT_MAX_OUTPUT_LENGTH);
	fprintf(o, "-b          Suppress counts; show only field values.\n");
	fprintf(o, "See also \"%s %s\".\n", argv0, "most");
}

static mapper_t* mapper_most_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	return mapper_most_or_least_parse_cli(pargi, argc, argv, TRUE);
}

static mapper_t* mapper_least_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	return mapper_most_or_least_parse_cli(pargi, argc, argv, FALSE);
}

static mapper_t* mapper_most_or_least_parse_cli(int* pargi, int argc, char** argv, int descending) {
	slls_t*   pgroup_by_field_names = NULL;
	long long max_output_length     = DEFAULT_MAX_OUTPUT_LENGTH;
	int       show_counts           = TRUE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pgroup_by_field_names);
	ap_define_long_long_flag(pstate,   "-n", &max_output_length);
	ap_define_false_flag(pstate,       "-b", &show_counts);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_most_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pgroup_by_field_names == NULL) {
		mapper_most_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_most_or_least_alloc(pstate, pgroup_by_field_names, max_output_length, descending, show_counts);
}

// ----------------------------------------------------------------
static mapper_t* mapper_most_or_least_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	long long max_output_length, int descending, int show_counts)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_most_or_least_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_most_or_least_state_t));

	pstate->pargp                 = pargp;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->pcounts_by_group      = lhmslv_alloc();
	pstate->max_output_length     = max_output_length;
	pstate->descending            = descending;
	pstate->show_counts           = show_counts;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_most_or_least_process;
	pmapper->pfree_func    = mapper_most_or_least_free;

	return pmapper;
}

static void mapper_most_or_least_free(mapper_t* pmapper) {
	mapper_most_or_least_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pgroup_by_field_names);
	// lhmslv_free will free the keys: we only need to free the void-star values.
	for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
		unsigned long long* pcount = pa->pvvalue;
		free(pcount);
	}
	lhmslv_free(pstate->pcounts_by_group);
	pstate->pgroup_by_field_names = NULL;
	pstate->pcounts_by_group = NULL;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_most_or_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_most_or_least_state_t* pstate = pvstate;

	if (pinrec != NULL) { // Not end of input record stream
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values != NULL) {
			unsigned long long* pcount = lhmslv_get(pstate->pcounts_by_group, pgroup_by_field_values);
			if (pcount == NULL) {
				pcount = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount = 1LL;
				lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values), pcount, FREE_ENTRY_KEY);
			} else {
				(*pcount)++;
			}
			slls_free(pgroup_by_field_values);
		}
		lrec_free(pinrec);
		return NULL;

	} else { // End of input record stream

		// Copy keys and counters from hashmap to array for sorting
		int input_length = pstate->pcounts_by_group->num_occupied;
		sort_pair_t* sort_pairs = mlr_malloc_or_die(input_length * sizeof(sort_pair_t));
		int i = 0;
		for (lhmslve_t* pe = pstate->pcounts_by_group->phead; pe != NULL; pe = pe->pnext) {
			sort_pairs[i].pgroup_by_field_values = pe->key;
			sort_pairs[i].count = *(long long *)pe->pvvalue;
			i++;
		}

		// Sort by count
		qsort(sort_pairs, input_length, sizeof(sort_pair_t),
			pstate->descending ? descending_vcmp : ascending_vcmp);

		// Emit top n
		sllv_t* poutrecs = sllv_alloc();
		int output_length = (input_length < pstate->max_output_length) ? input_length : pstate->max_output_length;
		for (i = 0; i < output_length; i++) {
			lrec_t* poutrec = lrec_unbacked_alloc();
			slls_t* pgroup_by_field_values = sort_pairs[i].pgroup_by_field_values;
			sllse_t* pb = pstate->pgroup_by_field_names->phead;
			sllse_t* pc =         pgroup_by_field_values->phead;
			for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
				lrec_put(poutrec, pb->value, pc->value, NO_FREE);
			}
			if (pstate->show_counts) {
				lrec_put(poutrec, "count", mlr_alloc_string_from_ull(sort_pairs[i].count), FREE_ENTRY_VALUE);
			}
			sllv_append(poutrecs, poutrec);
		}
		sllv_append(poutrecs, NULL);

		free(sort_pairs);
		return poutrecs;
	}
}

static int descending_vcmp(const void* pva, const void* pvb) {
	const sort_pair_t* pa = pva;
	const sort_pair_t* pb = pvb;
	return pb->count - pa->count;
}
static int ascending_vcmp(const void* pva, const void* pvb) {
	const sort_pair_t* pa = pva;
	const sort_pair_t* pb = pvb;
	return pa->count - pb->count;
}
