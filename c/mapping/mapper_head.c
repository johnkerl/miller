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

typedef struct _mapper_head_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	unsigned long long head_count;
	unsigned long long unkeyed_record_count;
	lhmslv_t* pcounts_by_group;
} mapper_head_state_t;

static void      mapper_head_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_head_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_head_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, unsigned long long head_count);
static void      mapper_head_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_head_process_unkeyed(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_head_process_keyed(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_head_setup = {
	.verb = "head",
	.pusage_func = mapper_head_usage,
	.pparse_func = mapper_head_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_head_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-n {count}    Head count to print; default 10\n");
	fprintf(o, "-g {a,b,c}    Optional group-by-field names for head counts\n");
	fprintf(o, "Passes through the first n records, optionally by category.\n");
	fprintf(o, "Without -g, ceases consuming more input (i.e. is fast) when n\n");
	fprintf(o, "records have been read.\n");
}

static mapper_t* mapper_head_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int     head_count            = 10;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &head_count);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_head_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_head_alloc(pstate, pgroup_by_field_names, head_count);
}

// ----------------------------------------------------------------
static mapper_t* mapper_head_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, unsigned long long head_count) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_head_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_head_state_t));

	pstate->pargp                  = pargp;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->head_count             = head_count;
	pstate->unkeyed_record_count   = 0LL;
	pstate->pcounts_by_group       = lhmslv_alloc();

	pmapper->pvstate        = pstate;
	pmapper->pprocess_func  = pgroup_by_field_names->length == 0
		? mapper_head_process_unkeyed
		: mapper_head_process_keyed;
	pmapper->pfree_func     = mapper_head_free;

	return pmapper;
}

static void mapper_head_free(mapper_t* pmapper, context_t* _) {
	mapper_head_state_t* pstate = pmapper->pvstate;
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
static sllv_t* mapper_head_process_unkeyed(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_head_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		pstate->unkeyed_record_count++;
		if (pstate->unkeyed_record_count <= pstate->head_count) {
			return sllv_single(pinrec);
		} else {
			pctx->force_eof = TRUE;
			lrec_free(pinrec);
			return NULL;
		}
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_head_process_keyed(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_head_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) {
			lrec_free(pinrec);
			return NULL;
		} else {
			unsigned long long* pcount_for_group = lhmslv_get(pstate->pcounts_by_group,
				pgroup_by_field_values);
			if (pcount_for_group == NULL) {
				pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount_for_group = 0LL;
				lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values),
					pcount_for_group, FREE_ENTRY_KEY);
			}
			slls_free(pgroup_by_field_values);
			(*pcount_for_group)++;
			if (*pcount_for_group <= pstate->head_count) {
				return sllv_single(pinrec);
			} else {
				lrec_free(pinrec);
				return NULL;
			}
		}
	} else {
		return sllv_single(NULL);
	}
}
