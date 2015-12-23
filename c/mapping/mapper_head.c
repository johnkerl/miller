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
	lhmslv_t* precord_lists_by_group;
} mapper_head_state_t;

static sllv_t*   mapper_head_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_head_free(mapper_t* pmapper);
static mapper_t* mapper_head_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, unsigned long long head_count);
static void      mapper_head_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_head_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_head_setup = {
	.verb = "head",
	.pusage_func = mapper_head_usage,
	.pparse_func = mapper_head_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_head_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-n {count}    Head count to print; default 10\n");
	fprintf(o, "-g {a,b,c}    Optional group-by-field names for head counts\n");
	fprintf(o, "Passes through the first n records, optionally by category.\n");
}

static mapper_t* mapper_head_parse_cli(int* pargi, int argc, char** argv) {
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
	pstate->precord_lists_by_group = lhmslv_alloc();

	pmapper->pvstate        = pstate;
	pmapper->pprocess_func  = mapper_head_process;
	pmapper->pfree_func     = mapper_head_free;

	return pmapper;
}

static void mapper_head_free(mapper_t* pmapper) {
	mapper_head_state_t* pstate = pmapper->pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
		unsigned long long* pcount_for_group = pa->pvvalue;
		free(pcount_for_group);
	}
	lhmslv_free(pstate->precord_lists_by_group);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_head_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_head_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) {
			lrec_free(pinrec);
			return NULL;
		} else {
			unsigned long long* pcount_for_group = lhmslv_get(pstate->precord_lists_by_group, pgroup_by_field_values);
			if (pcount_for_group == NULL) {
				pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount_for_group = 0LL;
				lhmslv_put(pstate->precord_lists_by_group, slls_copy(pgroup_by_field_values), pcount_for_group);
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
