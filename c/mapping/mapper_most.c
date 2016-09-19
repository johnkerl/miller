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

#define DEFAULT_MOST_COUNT 10LL

// xxx pick a better name ...
typedef struct _mapper_most_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	lhmslv_t* pcounts_by_group;
	long long most_count;
} mapper_most_state_t;

static void      mapper_most_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_most_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_most_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, long long most_count);
static void      mapper_most_free(mapper_t* pmapper);

static sllv_t* mapper_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

typedef struct _sort_pair_t {
	slls_t* pgroup_by_field_values;
	unsigned long long count;
} sort_pair_t;

// ----------------------------------------------------------------
mapper_setup_t mapper_most_setup = {
	.verb = "most",
	.pusage_func = mapper_most_usage,
	.pparse_func = mapper_most_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_most_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "xxx type me up.\n");
}

static mapper_t* mapper_most_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pgroup_by_field_names = NULL;
	long long most_count = DEFAULT_MOST_COUNT;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pgroup_by_field_names);
	ap_define_long_long_flag(pstate,   "-n", &most_count);
//	ap_define_true_flag(pstate,        "-n", &show_num_distinct_only);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_most_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pgroup_by_field_names == NULL) {
		mapper_most_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_most_alloc(pstate, pgroup_by_field_names, most_count);
}

// ----------------------------------------------------------------
static mapper_t* mapper_most_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, long long most_count) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_most_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_most_state_t));

	pstate->pargp                  = pargp;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->pcounts_by_group       = lhmslv_alloc();
	pstate->most_count             = most_count;

	pmapper->pvstate = pstate;
	pmapper->pprocess_func = mapper_most_process;
	pmapper->pfree_func = mapper_most_free;

	return pmapper;
}

static void mapper_most_free(mapper_t* pmapper) {
	mapper_most_state_t* pstate = pmapper->pvstate;
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
static sllv_t* mapper_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_most_state_t* pstate = pvstate;

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

		int input_length = pstate->pcounts_by_group->num_occupied;
		sort_pair_t* sort_pairs = mlr_malloc_or_die(input_length * sizeof(sort_pair_t));

		int i = 0;
		for (lhmslve_t* pe = pstate->pcounts_by_group->phead; pe != NULL; pe = pe->pnext) {
			sort_pairs[i].pgroup_by_field_values = pe->key;
			sort_pairs[i].count = *(long long *)pe->pvvalue;
			i++;
		}

		//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
		// xxx sort
		//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

		sllv_t* poutrecs = sllv_alloc();

		int output_length = (input_length < pstate->most_count) ? input_length : pstate->most_count;

		for (i = 0; i < output_length; i++) {
			lrec_t* poutrec = lrec_unbacked_alloc();
			slls_t* pgroup_by_field_values = sort_pairs[i].pgroup_by_field_values;
			sllse_t* pb = pstate->pgroup_by_field_names->phead;
			sllse_t* pc =         pgroup_by_field_values->phead;
			for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
				lrec_put(poutrec, pb->value, pc->value, NO_FREE);
			}

			lrec_put(poutrec, "count", mlr_alloc_string_from_ull(sort_pairs[i].count), FREE_ENTRY_VALUE);
			sllv_append(poutrecs, poutrec);
		}

		free(sort_pairs);

		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}
