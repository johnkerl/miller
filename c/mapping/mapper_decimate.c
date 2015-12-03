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
	slls_t* pgroup_by_field_names;
	unsigned long long decimate_count;
	unsigned long long remainder_for_keep;
	lhmslv_t* precord_lists_by_group;
} mapper_decimate_state_t;

static sllv_t*   mapper_decimate_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_decimate_free(void* pvstate);
static mapper_t* mapper_decimate_alloc(slls_t* pgroup_by_field_names,
	unsigned long long decimate_count, int keep_last);
static void      mapper_decimate_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_decimate_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_decimate_setup = {
	.verb = "decimate",
	.pusage_func = mapper_decimate_usage,
	.pparse_func = mapper_decimate_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_decimate_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-n {count}    Decimation factor; default 10\n");
	fprintf(o, "-b            Decimate by printing first of every n.\n");
	fprintf(o, "-e            Decimate by printing last of every n (default).\n");
	fprintf(o, "-g {a,b,c}    Optional group-by-field names for decimate counts\n");
	fprintf(o, "Passes through the first n records, optionally by category.\n");
}

static mapper_t* mapper_decimate_parse_cli(int* pargi, int argc, char** argv) {
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

	return mapper_decimate_alloc(pgroup_by_field_names, decimate_count, keep_last);
}

// ----------------------------------------------------------------
static mapper_t* mapper_decimate_alloc(slls_t* pgroup_by_field_names,
	unsigned long long decimate_count, int keep_last)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_decimate_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_decimate_state_t));

	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->decimate_count         = decimate_count;
	pstate->remainder_for_keep     = keep_last ? decimate_count - 1 : 0;
	pstate->precord_lists_by_group = lhmslv_alloc();

	pmapper->pvstate        = pstate;
	pmapper->pprocess_func  = mapper_decimate_process;
	pmapper->pfree_func     = mapper_decimate_free;

	return pmapper;
}

static void mapper_decimate_free(void* pvstate) {
	mapper_decimate_state_t* pstate = (mapper_decimate_state_t*)pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	// xxx recursively free void-stars ... here & elsewhere.
	lhmslv_free(pstate->precord_lists_by_group);
}

// ----------------------------------------------------------------
static sllv_t* mapper_decimate_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_decimate_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) {
			return NULL;
		} else {
			unsigned long long* pcount_for_group = lhmslv_get(pstate->precord_lists_by_group, pgroup_by_field_values);
			if (pcount_for_group == NULL) {
				pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount_for_group = 0LL;
				lhmslv_put(pstate->precord_lists_by_group, slls_copy(pgroup_by_field_values), pcount_for_group);
			}

			unsigned long long remainder = *pcount_for_group % pstate->decimate_count;
			if (remainder == pstate->remainder_for_keep) {
				(*pcount_for_group)++;
				return sllv_single(pinrec);
			} else {
				(*pcount_for_group)++;
				lrec_free(pinrec);
				return NULL;
			}
		}
	} else {
		return sllv_single(NULL);
	}
}
