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
// -a min,max -v x,y -g a,b

// ["s","t"] |--> "x" |--> "sum" |--> step_sum_t* (as void*)
// level 1      level 2   level 3
// lhmslv_t     lhmsv_t   lhmsv_t
// step_sum_t implements interface:
//   void  dacc(double dval);
//   void  sacc(char*  sval);
//   char* get();

// ================================================================
typedef void  step_dput_func_t(void* pvstate, double dval, lrec_t* prec);
typedef void  step_sput_func_t(void* pvstate, char*  sval, lrec_t* prec);

typedef struct _step_t {
	void* pvstate;
	step_sput_func_t* psput_func;
	step_dput_func_t* pdput_func;
	static_context_t* pstatx;
} step_t;

typedef step_t* step_alloc_func_t(char* input_field_name, static_context_t* pstatx);

// ----------------------------------------------------------------
typedef struct _step_delta_state_t {
	double prev;
	int    have_prev;
	char*  output_field_name;
	static_context_t* pstatx;
} step_delta_state_t;
void step_delta_dput(void* pvstate, double dval, lrec_t* prec) {
	step_delta_state_t* pstate = pvstate;
	double delta = dval;
	if (pstate->have_prev) {
		delta = dval - pstate->prev;
	} else {
		pstate->have_prev = TRUE;
	}
	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_double(delta, pstate->pstatx->ofmt),
		LREC_FREE_ENTRY_VALUE);
	pstate->prev = dval;
}
step_t* step_delta_alloc(char* input_field_name, static_context_t* pstatx) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_delta_state_t* pstate = mlr_malloc_or_die(sizeof(step_delta_state_t));
	pstate->prev      = -999.0;
	pstate->have_prev = FALSE;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_delta");
	pstate->pstatx    = pstatx;
	pstep->pvstate    = (void*)pstate;
	pstep->psput_func = NULL;
	pstep->pdput_func = &step_delta_dput;
	return pstep;
}
// xxx step_delta_free et al.

// ----------------------------------------------------------------
typedef struct _step_rsum_state_t {
	double rsum;
	char*  output_field_name;
	static_context_t* pstatx;
} step_rsum_state_t;
void step_rsum_dput(void* pvstate, double dval, lrec_t* prec) {
	step_rsum_state_t* pstate = pvstate;
	pstate->rsum += dval;
	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_double(pstate->rsum, pstate->pstatx->ofmt),
		LREC_FREE_ENTRY_VALUE);
}
step_t* step_rsum_alloc(char* input_field_name, static_context_t* pstatx) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_rsum_state_t* pstate = mlr_malloc_or_die(sizeof(step_rsum_state_t));
	pstate->rsum      = 0.0;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_rsum");
	pstate->pstatx    = pstatx;
	pstep->pvstate    = (void*)pstate;
	pstep->psput_func = NULL;
	pstep->pdput_func = &step_rsum_dput;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_counter_state_t {
	unsigned long long counter;
	char*  output_field_name;
	static_context_t* pstatx;
} step_counter_state_t;
void step_counter_sput(void* pvstate, char* sval, lrec_t* prec) {
	step_counter_state_t* pstate = pvstate;
	pstate->counter++;
	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_ull(pstate->counter),
		LREC_FREE_ENTRY_VALUE);
}
step_t* step_counter_alloc(char* input_field_name, static_context_t* pstatx) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_counter_state_t* pstate = mlr_malloc_or_die(sizeof(step_counter_state_t));
	pstate->counter   = 0LL;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_counter");
	pstate->pstatx    = pstatx;
	pstep->pvstate    = (void*)pstate;
	pstep->psput_func = &step_counter_sput;
	pstep->pdput_func = NULL;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_lookup_t {
	char* name;
	step_alloc_func_t* pnew_func;
} step_lookup_t;
static step_lookup_t step_lookup_table[] = {
	{"delta",   step_delta_alloc},
	{"rsum",    step_rsum_alloc},
	{"counter", step_counter_alloc},
};
static int step_lookup_table_length = sizeof(step_lookup_table) / sizeof(step_lookup_table[0]);

static step_t* make_step(char* step_name, char* input_field_name, static_context_t* pstatx) {
	for (int i = 0; i < step_lookup_table_length; i++)
		if (streq(step_name, step_lookup_table[i].name))
			return step_lookup_table[i].pnew_func(input_field_name, pstatx);
	return NULL;
}

// ================================================================
typedef struct _mapper_step_state_t {
	slls_t* pstepper_names;
	slls_t* pvalue_field_names;
	slls_t* pgroup_by_field_names;

	lhmslv_t* pmaps_level_1;

} mapper_step_state_t;

// given: step rsum,delta values x,y group by a,b
// example input:       example output:
//   a b x y            a b x_count x_sum y_count y_sum
//   s t 1 2            s t 2       6     2       8
//   u v 3 4            u v 1       3     1       4
//   s t 5 6            u w 1       7     1       9
//   u w 7 9

// ["s","t"] |--> "x" |--> "sum" |--> step_sum_t* (as void*)
// level_1      level_2   level_3
// lhmslv_t     lhmsv_t   lhmsv_t
// step_sum_t implements interface:
//   void  init();
//   void  dacc(double dval);
//   void  sacc(char*  sval);
//   char* get();

// ----------------------------------------------------------------
sllv_t* mapper_step_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_step_state_t* pstate = pvstate;
	if (pinrec == NULL)
		return sllv_single(NULL);

	// ["s", "t"]
	// xxx make value_field_values into a hashmap. then accept partial population on that.
	// xxx but retain full-population requirement on group-by.
	// e.g. if accumulating stats of x,y on a,b then skip row with x,y,a but process row with x,a,b.
	slls_t* pvalue_field_values    = mlr_selected_values_from_record(pinrec, pstate->pvalue_field_names); // xxx this can be slls or sllv
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

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
		int   have_dval = FALSE;
		double value_field_dval = -999.0;

		lhmsv_t* pmaps_level_3 = lhmsv_get(pmaps_level_2, value_field_name);
		if (pmaps_level_3 == NULL) {
			pmaps_level_3 = lhmsv_alloc();
			lhmsv_put(pmaps_level_2, value_field_name, pmaps_level_3);
		}

		// for "delta", "rsum"
		sllse_t* pc = pstate->pstepper_names->phead;
		for ( ; pc != NULL; pc = pc->pnext) {
			char* step_name = pc->value;
			step_t* pstep = lhmsv_get(pmaps_level_3, step_name);
			if (pstep == NULL) {
				pstep = make_step(step_name, value_field_name, &pctx->statx);
				if (pstep == NULL) {
					fprintf(stderr, "mlr step: stepper \"%s\" not found.\n",
						step_name);
					exit(1);
				}
				lhmsv_put(pmaps_level_3, step_name, pstep);
			}

			if (pstep->psput_func != NULL) {
				pstep->psput_func(pstep->pvstate, value_field_sval, pinrec);
			}
			if (pstep->pdput_func != NULL) {
				if (!have_dval) {
					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
					have_dval = TRUE;
				}
				pstep->pdput_func(pstep->pvstate, value_field_dval, pinrec);
			}
		}
	}
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static void mapper_step_free(void* pvstate) {
	mapper_step_state_t* pstate = pvstate;
	slls_free(pstate->pstepper_names);
	slls_free(pstate->pvalue_field_names);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->pmaps_level_1);
}

mapper_t* mapper_step_alloc(slls_t* pstepper_names, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_step_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_step_state_t));

	pstate->pstepper_names        = pstepper_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->pmaps_level_1         = lhmslv_alloc();

	pmapper->pvstate      = pstate;
	pmapper->pmapper_process_func = mapper_step_func;
	pmapper->pmapper_free_func  = mapper_step_free;

	return pmapper;
}

// ----------------------------------------------------------------
void mapper_step_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-a {delta,rsum,...}    Names of steppers: one or more of\n");
	fprintf(stdout, "                     ");
	for (int i = 0; i < step_lookup_table_length; i++) {
		fprintf(stdout, " %s", step_lookup_table[i].name);
	}
	fprintf(stdout, "\n");
	fprintf(stdout, "-f {a,b,c}            Value-field names on which to compute statistics\n");
	fprintf(stdout, "-g {d,e,f}            Group-by-field names\n");
}

mapper_t* mapper_step_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pstepper_names        = NULL;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-a", &pstepper_names);
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_step_usage(argv[0], verb);
		return NULL;
	}

	if (pstepper_names == NULL || pvalue_field_names == NULL) {
		mapper_step_usage(argv[0], verb);
		return NULL;
	}

	return mapper_step_alloc(pstepper_names, pvalue_field_names, pgroup_by_field_names);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_step_setup = {
	.verb = "step",
	.pusage_func = mapper_step_usage,
	.pparse_func = mapper_step_parse_cli
};
