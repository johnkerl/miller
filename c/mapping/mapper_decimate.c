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

typedef struct _mapper_decimate_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	unsigned long long decimate_count;
	unsigned long long remainder_for_keep;
	lhmslv_t* pcounts_by_group;
} mapper_decimate_state_t;

static void      mapper_decimate_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_decimate_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_decimate_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	unsigned long long decimate_count, int keep_last);
static void      mapper_decimate_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_decimate_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_decimate_setup = {
	.verb = "decimate",
	.pusage_func = mapper_decimate_usage,
	.pparse_func = mapper_decimate_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_decimate_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-n {count}    Decimation factor; default 10\n");
	fprintf(o, "-b            Decimate by printing first of every n.\n");
	fprintf(o, "-e            Decimate by printing last of every n (default).\n");
	fprintf(o, "-g {a,b,c}    Optional group-by-field names for decimate counts\n");
	fprintf(o, "Passes through one of every n records, optionally by category.\n");
}

static mapper_t* mapper_decimate_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int     decimate_count = 10;
	int     keep_last      = TRUE;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &decimate_count);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_false_flag(pstate, "-b", &keep_last);
	ap_define_true_flag(pstate, "-e", &keep_last);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_decimate_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_decimate_alloc(pstate, pgroup_by_field_names, decimate_count, keep_last);
}

// ----------------------------------------------------------------
static mapper_t* mapper_decimate_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	unsigned long long decimate_count, int keep_last)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_decimate_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_decimate_state_t));

	pstate->pargp                  = pargp;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->decimate_count         = decimate_count;
	pstate->remainder_for_keep     = keep_last ? decimate_count - 1 : 0;
	pstate->pcounts_by_group = lhmslv_alloc();

	pmapper->pvstate        = pstate;
	pmapper->pprocess_func  = mapper_decimate_process;
	pmapper->pfree_func     = mapper_decimate_free;

	return pmapper;
}

static void mapper_decimate_free(mapper_t* pmapper, context_t* _) {
	mapper_decimate_state_t* pstate = pmapper->pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
		unsigned long long* pcount_for_group = pa->pvvalue;
		free(pcount_for_group);
	}
	lhmslv_free(pstate->pcounts_by_group);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_decimate_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_decimate_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) {
			return NULL;
		} else {
			unsigned long long* pcount_for_group = lhmslv_get(pstate->pcounts_by_group, pgroup_by_field_values);
			if (pcount_for_group == NULL) {
				pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount_for_group = 0LL;
				lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values), pcount_for_group,
					FREE_ENTRY_KEY);
			}

			unsigned long long remainder = *pcount_for_group % pstate->decimate_count;
			if (remainder == pstate->remainder_for_keep) {
				(*pcount_for_group)++;
				slls_free(pgroup_by_field_values);
				return sllv_single(pinrec);
			} else {
				(*pcount_for_group)++;
				lrec_free(pinrec);
				slls_free(pgroup_by_field_values);
				return NULL;
			}
		}
	} else {
		return sllv_single(NULL);
	}
}
