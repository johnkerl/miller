#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/string_array.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "mapping/mlr_val.h"
#include "cli/argparse.h"

// ----------------------------------------------------------------
typedef void step_dprocess_func_t(void* pvstate, double fltv, lrec_t* prec);
typedef void step_nprocess_func_t(void* pvstate, mv_t* pval,  lrec_t* prec);
typedef void step_sprocess_func_t(void* pvstate, char*  strv, lrec_t* prec);

typedef struct _step_t {
	void* pvstate;
	step_dprocess_func_t* pdprocess_func;
	step_nprocess_func_t* pnprocess_func;
	step_sprocess_func_t* psprocess_func;
} step_t;

typedef step_t* step_alloc_func_t(char* input_field_name);

typedef struct _mapper_step_state_t {
	slls_t* pstepper_names;
	string_array_t* pvalue_field_names;  // parameter
	string_array_t* pvalue_field_values; // scratch space used per-record
	slls_t* pgroup_by_field_names;       // parameter
	lhmslv_t* groups;
} mapper_step_state_t;

// Multilevel hashmap structure:
// {
//   ["s","t"] : {              <--- group-by field names
//     ["x","y"] : {            <--- value field names
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : C stats2_corr_t object,
//       "cov"  : C stats2_cov_t  object
//     }
//   },
// }

// ----------------------------------------------------------------
static void      mapper_step_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_step_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_step_alloc(slls_t* pstepper_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names);
static void      mapper_step_free(void* pvstate);
static sllv_t*   mapper_step_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

static step_t* step_delta_alloc(char* input_field_name);
static step_t* step_ratio_alloc(char* input_field_name);
static step_t* step_rsum_alloc(char* input_field_name);
static step_t* step_counter_alloc(char* input_field_name);

static step_t* make_step(char* step_name, char* input_field_name);

typedef struct _step_lookup_t {
	char* name;
	step_alloc_func_t* pnew_func;
	char* desc;
} step_lookup_t;
static step_lookup_t step_lookup_table[] = {
	{"delta",   step_delta_alloc,   "Compute differences in field(s) between successive records"},
	{"ratio",   step_ratio_alloc,   "Compute ratios in field(s) between successive records"},
	{"rsum",    step_rsum_alloc,    "Compute running sums of field(s) between successive records"},
	{"counter", step_counter_alloc, "Count instances of field(s) between successive records"},
};
static int step_lookup_table_length = sizeof(step_lookup_table) / sizeof(step_lookup_table[0]);

// ----------------------------------------------------------------
mapper_setup_t mapper_step_setup = {
	.verb = "step",
	.pusage_func = mapper_step_usage,
	.pparse_func = mapper_step_parse_cli
};

// ----------------------------------------------------------------
static void mapper_step_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-a {delta,rsum,...}   Names of steppers: comma-separated, one or more of:\n");
	for (int i = 0; i < step_lookup_table_length; i++) {
		fprintf(o, "  %-8s %s\n", step_lookup_table[i].name, step_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c}            Value-field names on which to compute statistics\n");
	fprintf(o, "-g {d,e,f}            Optional group-by-field names\n");
	fprintf(o, "Computes values dependent on the previous record, optionally grouped\n");
	fprintf(o, "by category.\n");
}

static mapper_t* mapper_step_parse_cli(int* pargi, int argc, char** argv) {
	slls_t*         pstepper_names        = NULL;
	string_array_t* pvalue_field_names    = NULL;
	slls_t*         pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate,  "-a", &pstepper_names);
	ap_define_string_array_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate,  "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_step_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pstepper_names == NULL || pvalue_field_names == NULL) {
		mapper_step_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_step_alloc(pstepper_names, pvalue_field_names, pgroup_by_field_names);
}

// ----------------------------------------------------------------
static mapper_t* mapper_step_alloc(slls_t* pstepper_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_step_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_step_state_t));

	pstate->pstepper_names        = pstepper_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pvalue_field_values   = string_array_alloc(pvalue_field_names->length);
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->groups                = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_step_process;
	pmapper->pfree_func    = mapper_step_free;

	return pmapper;
}

static void mapper_step_free(void* pvstate) {
	mapper_step_state_t* pstate = pvstate;
	slls_free(pstate->pstepper_names);
	string_array_free(pstate->pvalue_field_names);
	string_array_free(pstate->pvalue_field_values);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->groups);
}

// ----------------------------------------------------------------
static sllv_t* mapper_step_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_step_state_t* pstate = pvstate;
	if (pinrec == NULL)
		return sllv_single(NULL);

	// ["s", "t"]
	mlr_reference_values_from_record(pinrec, pstate->pvalue_field_names, pstate->pvalue_field_values);
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

	if (pgroup_by_field_values == NULL) {
		slls_free(pgroup_by_field_values);
		return sllv_single(pinrec);
	}

	lhmsv_t* group_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (group_to_acc_field == NULL) {
		group_to_acc_field = lhmsv_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), group_to_acc_field);
	}

	// for x=1 and y=2
	int n = pstate->pvalue_field_names->length;
	for (int i = 0; i < n; i++) {
		char* value_field_name = pstate->pvalue_field_names->strings[i];
		char* value_field_sval = pstate->pvalue_field_values->strings[i];
		if (value_field_sval == NULL)
			continue;

		int have_dval = FALSE;
		int have_nval = FALSE;
		double value_field_dval = -999.0;
		mv_t   value_field_nval = mv_from_int(-888LL);

		lhmsv_t* acc_field_to_acc_state = lhmsv_get(group_to_acc_field, value_field_name);
		if (acc_field_to_acc_state == NULL) {
			acc_field_to_acc_state = lhmsv_alloc();
			lhmsv_put(group_to_acc_field, value_field_name, acc_field_to_acc_state);
		}

		// for "delta", "rsum"
		sllse_t* pc = pstate->pstepper_names->phead;
		for ( ; pc != NULL; pc = pc->pnext) {
			char* step_name = pc->value;
			step_t* pstep = lhmsv_get(acc_field_to_acc_state, step_name);
			if (pstep == NULL) {
				pstep = make_step(step_name, value_field_name);
				if (pstep == NULL) {
					fprintf(stderr, "mlr step: stepper \"%s\" not found.\n",
						step_name);
					exit(1);
				}
				lhmsv_put(acc_field_to_acc_state, step_name, pstep);
			}

			if (pstep->pdprocess_func != NULL) {
				if (!have_dval) {
					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
					have_dval = TRUE;
				}
				pstep->pdprocess_func(pstep->pvstate, value_field_dval, pinrec);
			}

			if (pstep->pnprocess_func != NULL) {
				if (!have_nval) {
					value_field_nval = mv_scan_number_or_die(value_field_sval);
					have_nval = TRUE;
				}
				pstep->pnprocess_func(pstep->pvstate, &value_field_nval, pinrec);
			}

			if (pstep->psprocess_func != NULL) {
				pstep->psprocess_func(pstep->pvstate, value_field_sval, pinrec);
			}

		}
	}
	return sllv_single(pinrec);
}

static step_t* make_step(char* step_name, char* input_field_name) {
	for (int i = 0; i < step_lookup_table_length; i++)
		if (streq(step_name, step_lookup_table[i].name))
			return step_lookup_table[i].pnew_func(input_field_name);
	return NULL;
}

// ----------------------------------------------------------------
typedef struct _step_delta_state_t {
	mv_t  prev;
	char* output_field_name;
} step_delta_state_t;
static void step_delta_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
	step_delta_state_t* pstate = pvstate;
	mv_t delta;
	if (mv_is_null(&pstate->prev)) {
		delta = mv_from_int(0);
	} else {
		delta = n_nn_minus_func(&pstate->prev, pnumv);
	}
	lrec_put(prec, pstate->output_field_name, mv_format_val(&delta), LREC_FREE_ENTRY_VALUE);
	pstate->prev = *pnumv;
}
static step_t* step_delta_alloc(char* input_field_name) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_delta_state_t* pstate = mlr_malloc_or_die(sizeof(step_delta_state_t));
	pstate->prev = mv_from_null();
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_delta");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = step_delta_nprocess;
	pstep->psprocess_func = NULL;
	return pstep;
}
// xxx step_delta_free et al.

// ----------------------------------------------------------------
typedef struct _step_ratio_state_t {
	double prev;
	int    have_prev;
	char*  output_field_name;
} step_ratio_state_t;
static void step_ratio_dprocess(void* pvstate, double fltv, lrec_t* prec) {
	step_ratio_state_t* pstate = pvstate;
	double ratio = 1.0;
	if (pstate->have_prev) {
		ratio = fltv / pstate->prev;
	} else {
		pstate->have_prev = TRUE;
	}
	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_double(ratio, MLR_GLOBALS.ofmt),
		LREC_FREE_ENTRY_VALUE);
	pstate->prev = fltv;
}
static step_t* step_ratio_alloc(char* input_field_name) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_ratio_state_t* pstate = mlr_malloc_or_die(sizeof(step_ratio_state_t));
	pstate->prev          = -999.0;
	pstate->have_prev     = FALSE;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_ratio");

	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = step_ratio_dprocess;
	pstep->pnprocess_func = NULL;
	pstep->psprocess_func = NULL;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_rsum_state_t {
	mv_t   rsum;
	char*  output_field_name;
} step_rsum_state_t;

static void step_rsum_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
	step_rsum_state_t* pstate = pvstate;
	pstate->rsum = n_nn_plus_func(&pstate->rsum, pnumv);
	lrec_put(prec, pstate->output_field_name, mv_format_val(&pstate->rsum),
		LREC_FREE_ENTRY_VALUE);
}

static step_t* step_rsum_alloc(char* input_field_name) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_rsum_state_t* pstate = mlr_malloc_or_die(sizeof(step_rsum_state_t));
	pstate->rsum = mv_from_int(0LL);
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_rsum");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = step_rsum_nprocess;;
	pstep->psprocess_func = NULL;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_counter_state_t {
	unsigned long long counter;
	char*  output_field_name;
} step_counter_state_t;
static void step_counter_sprocess(void* pvstate, char* strv, lrec_t* prec) {
	step_counter_state_t* pstate = pvstate;
	pstate->counter++;
	lrec_put(prec, pstate->output_field_name, mlr_alloc_string_from_ull(pstate->counter),
		LREC_FREE_ENTRY_VALUE);
}
static step_t* step_counter_alloc(char* input_field_name) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_counter_state_t* pstate = mlr_malloc_or_die(sizeof(step_counter_state_t));
	pstate->counter       = 0LL;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_counter");

	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->psprocess_func = step_counter_sprocess;
	pstep->pnprocess_func = NULL;
	return pstep;
}
