#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
typedef struct _mapper_tail_state_t {
	slls_t* pgroup_by_field_names;
	unsigned long long tail_count;
	lhmslv_t* precord_lists_by_group;
} mapper_tail_state_t;

// ----------------------------------------------------------------
sllv_t* mapper_tail_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tail_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
		sllv_t* precord_list_for_group = lhmslv_get(pstate->precord_lists_by_group, pgroup_by_field_values);
		if (precord_list_for_group == NULL) {
			precord_list_for_group = sllv_alloc();
			lhmslv_put(pstate->precord_lists_by_group, slls_copy(pgroup_by_field_values), precord_list_for_group);
		}
		if (precord_list_for_group->length >= pstate->tail_count) {
			lrec_t* porec = sllv_pop(precord_list_for_group);
			if (porec != NULL)
				lrec_free(porec);
		}
		sllv_add(precord_list_for_group, pinrec);

		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
			sllv_t* precord_list_for_group = pa->value;
			for (sllve_t* pb = precord_list_for_group->phead; pb != NULL; pb = pb->pnext) {
				sllv_add(poutrecs, pb->pvdata);
			}
		}
		sllv_add(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static void mapper_tail_free(void* pvstate) {
	mapper_tail_state_t* pstate = pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	if (pstate->precord_lists_by_group != NULL)
		// xxx free the void-star payloads 1st
		lhmslv_free(pstate->precord_lists_by_group);
}

// xxx make these all static?
mapper_t* mapper_tail_alloc(slls_t* pgroup_by_field_names, unsigned long long tail_count) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_tail_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tail_state_t));

	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->tail_count             = tail_count;
	pstate->precord_lists_by_group = lhmslv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_tail_func;
	pmapper->pmapper_free_func    = mapper_tail_free;

	return pmapper;
}

// ----------------------------------------------------------------
void mapper_tail_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-n {count}    Tail count to print; default 10\n");
	fprintf(stdout, "-g {a,b,c}    Group-by-field names for tail counts\n");
}

mapper_t* mapper_tail_parse_cli(int* pargi, int argc, char** argv) {
	int     tail_count             = 10;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &tail_count);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_tail_usage(argv[0], verb);
		return NULL;
	}

	return mapper_tail_alloc(pgroup_by_field_names, tail_count);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_tail_setup = {
	.verb = "tail",
	.pusage_func = mapper_tail_usage,
	.pparse_func = mapper_tail_parse_cli
};
