#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/top_keeper.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
typedef struct _mapper_top_state_t {
	slls_t* pvalue_field_names;
	slls_t* pgroup_by_field_names;
	int show_full_records;
	int top_count;
	double sign; // for +1 for max; -1 for min
	lhmslv_t* pmaps_level_1;
} mapper_top_state_t;

// ----------------------------------------------------------------
sllv_t* mapper_top_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_top_state_t* pstate = pvstate;
	if (pinrec != NULL) {

		// ["s", "t"]
		// xxx make value_field_values into a hashmap. then accept partial population on that.
		// xxx but retain full-population requirement on group-by.
		// e.g. if accumulating stats of x,y on a,b then skip row with x,y,a but process row with x,a,b.
		slls_t* pvalue_field_values    = mlr_selected_values_from_record(pinrec, pstate->pvalue_field_names);
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

		// xxx cmt
		if (pvalue_field_values->length != pstate->pvalue_field_names->length) {
			lrec_free(pinrec);
			return NULL;
		}
		if (pgroup_by_field_values->length != pstate->pgroup_by_field_names->length) {
			lrec_free(pinrec);
			return NULL;
		}

		lhmsv_t* pmaps_level_2 = lhmslv_get(pstate->pmaps_level_1, pgroup_by_field_values);
		if (pmaps_level_2 == NULL) {
			pmaps_level_2 = lhmsv_alloc();
			lhmslv_put(pstate->pmaps_level_1, slls_copy(pgroup_by_field_values), pmaps_level_2);
		}

		sllse_t* pa = pstate->pvalue_field_names->phead;
		sllse_t* pb =         pvalue_field_values->phead;
		// for "x", "y" and "1", "2"
		for ( ; pa != NULL && pb != NULL; pa = pa->pnext, pb = pb->pnext) {
			char* value_field_name = pa->value;
			char* value_field_sval = pb->value;
			double value_field_dval = mlr_double_from_string_or_die(value_field_sval);

			// xxx rename: for group and value-field-name
			top_keeper_t* ptop_keeper_for_group = lhmsv_get(pmaps_level_2, value_field_name);
			if (ptop_keeper_for_group == NULL) {
				ptop_keeper_for_group = top_keeper_alloc(pstate->top_count);
				lhmsv_put(pmaps_level_2, value_field_name, ptop_keeper_for_group);
			}

			top_keeper_add(ptop_keeper_for_group, value_field_dval * pstate->sign, pinrec);
		}

		lrec_free(pinrec);
		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->pmaps_level_1->phead; pa != NULL; pa = pa->pnext) {

			if (pstate->show_full_records) {
				// xxx if i do it this way (entire record) then there can only be one value column.
				// xxx assert on that, or at the very least note how confusing it is.
				lhmsv_t* pmaps_level_2 = pa->value;
				for (lhmsve_t* pd = pmaps_level_2->phead; pd != NULL; pd = pd->pnext) {
					top_keeper_t* ptop_keeper_for_group = pd->value;
					for (int i = 0;  i < ptop_keeper_for_group->size; i++)
						sllv_add(poutrecs, ptop_keeper_for_group->top_precords[i]);
				}
			}

			else {
				slls_t* pgroup_by_field_values = pa->key;
				for (int i = 0; i < pstate->top_count; i++) {
					lrec_t* poutrec = lrec_unbacked_alloc();

					// Add in a=s,b=t fields:
					sllse_t* pb = pstate->pgroup_by_field_names->phead;
					sllse_t* pc =         pgroup_by_field_values->phead;
					for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
						lrec_put(poutrec, pb->value, pc->value, 0);
					}

					char* sidx = mlr_alloc_string_from_ull(i+1);
					lrec_put(poutrec, "top_idx", sidx, LREC_FREE_ENTRY_VALUE);
					free(sidx);

					// Add in fields such as x_top_1=#
					lhmsv_t* pmaps_level_2 = pa->value;
					// for "x", "y"
					for (lhmsve_t* pd = pmaps_level_2->phead; pd != NULL; pd = pd->pnext) {
						char* value_field_name = pd->key;
						top_keeper_t* ptop_keeper_for_group = pd->value;

						char* key = mlr_paste_2_strings(value_field_name, "_top");
						if (i < ptop_keeper_for_group->size) {
							// xxx temp fmt
							double dval = ptop_keeper_for_group->top_values[i] * pstate->sign;
							char* sval = mlr_alloc_string_from_double(dval, pctx->statx.ofmt);
							lrec_put(poutrec, key, sval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
							free(sval);
						} else {
							lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
						}
					}
					sllv_add(poutrecs, poutrec);
				}
			}
		}

		sllv_add(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static void mapper_top_free(void* pvstate) {
	mapper_top_state_t* pstate = pvstate;
	slls_free(pstate->pvalue_field_names);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->pmaps_level_1);
}

mapper_t* mapper_top_alloc(slls_t* pvalue_field_names, slls_t* pgroup_by_field_names,
	int top_count, int do_max, int show_full_records)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_top_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_top_state_t));

	pstate->pvalue_field_names    = slls_copy(pvalue_field_names);
	pstate->pgroup_by_field_names = slls_copy(pgroup_by_field_names);
	pstate->show_full_records     = show_full_records;
	pstate->top_count             = top_count;
	pstate->sign                  = do_max ? 1.0 : -1.0;
	pstate->pmaps_level_1         = lhmslv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_top_func;
	pmapper->pmapper_free_func    = mapper_top_free;

	return pmapper;
}

// ----------------------------------------------------------------
void mapper_top_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-f {a,b,c}    Value-field names for top counts\n");
	fprintf(stdout, "-g {d,e,f}    Group-by-field names for top counts\n");
	fprintf(stdout, "-n {count}    Top n records to print; default 1\n");
	fprintf(stdout, "-a            Print all fields for top-value records; default is\n");
	fprintf(stdout, "              to print only value and group-by fields.\n");
	fprintf(stdout, "--min         Print top smallest values; default is top largest values\n");
}

mapper_t* mapper_top_parse_cli(int* pargi, int argc, char** argv) {
	int     top_count             = 1;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();
	int     show_full_records     = FALSE;
	int     do_max                = TRUE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &top_count);
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate, "-a", &show_full_records);
	ap_define_true_flag(pstate, "--max", &do_max);
	ap_define_false_flag(pstate, "--min", &do_max);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_top_usage(argv[0], verb);
		return NULL;
	}

	if (pvalue_field_names == NULL) {
		mapper_top_usage(argv[0], verb);
		return NULL;
	}

	return mapper_top_alloc(pvalue_field_names, pgroup_by_field_names,
		top_count, do_max, show_full_records);
}

mapper_setup_t mapper_top_setup = {
	.verb = "top",
	.pusage_func = mapper_top_usage,
	.pparse_func = mapper_top_parse_cli
};
