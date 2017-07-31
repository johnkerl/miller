#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "lib/string_array.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "lib/mvfuncs.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DEFAULT_STRING_ALPHA "0.5"

// ----------------------------------------------------------------
struct _step_t; // forward reference for method declarations

typedef void step_dprocess_func_t(void* pvstate, double fltv, lrec_t* prec);
typedef void step_nprocess_func_t(void* pvstate, mv_t* pval,  lrec_t* prec);
typedef void step_sprocess_func_t(void* pvstate, char*  strv, lrec_t* prec);
typedef void step_zprocess_func_t(void* pvstate,              lrec_t* prec);
typedef void step_free_func_t(struct _step_t* pstep);

typedef struct _step_t {
	void* pvstate;
	step_dprocess_func_t* pdprocess_func;
	step_nprocess_func_t* pnprocess_func;
	step_sprocess_func_t* psprocess_func;
	step_zprocess_func_t* pzprocess_func;
	step_free_func_t*     pfree_func;
} step_t;

typedef step_t* step_alloc_func_t(char* input_field_name, int allow_int_float,
	slls_t* pstring_alphas, slls_t* pewma_suffixes);

typedef struct _mapper_step_state_t {
	ap_state_t*     pargp;
	slls_t*         pstepper_names;
	string_array_t* pvalue_field_names;    // parameter
	string_array_t* pvalue_field_values;   // scratch space used per-record
	slls_t*         pgroup_by_field_names; // parameter
	lhmslv_t*       groups;
	int             allow_int_float;
	slls_t*         pstring_alphas;
	slls_t*         pewma_suffixes;
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
static mapper_t* mapper_step_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_step_alloc(ap_state_t* pargp, slls_t* pstepper_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int allow_int_float, slls_t* pstring_alphas, slls_t* pewma_suffixes);
static void      mapper_step_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_step_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

static step_t* step_delta_alloc      (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_shift_alloc      (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_from_first_alloc (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_ratio_alloc      (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_rsum_alloc       (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_counter_alloc    (char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2);
static step_t* step_ewma_alloc       (char* input_field_name, int unused,
	slls_t* pstring_alphas, slls_t* pewma_suffixes);

static step_t* make_step(char* step_name, char* input_field_name, int allow_int_float,
	slls_t* pstring_alphas, slls_t* pewma_suffixes);

typedef struct _step_lookup_t {
	char* name;
	step_alloc_func_t* palloc_func;
	char* desc;
} step_lookup_t;
static step_lookup_t step_lookup_table[] = {
	{"delta",      step_delta_alloc,      "Compute differences in field(s) between successive records"},
	{"shift",      step_shift_alloc,      "Include value(s) in field(s) from previous record, if any"},
	{"from-first", step_from_first_alloc, "Compute differences in field(s) from first record"},
	{"ratio",      step_ratio_alloc,      "Compute ratios in field(s) between successive records"},
	{"rsum",       step_rsum_alloc,       "Compute running sums of field(s) between successive records"},
	{"counter",    step_counter_alloc,    "Count instances of field(s) between successive records"},
	{"ewma",       step_ewma_alloc,       "Exponentially weighted moving average over successive records"},
};
static int step_lookup_table_length = sizeof(step_lookup_table) / sizeof(step_lookup_table[0]);

// ----------------------------------------------------------------
mapper_setup_t mapper_step_setup = {
	.verb = "step",
	.pusage_func = mapper_step_usage,
	.pparse_func = mapper_step_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_step_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Computes values dependent on the previous record, optionally grouped\n");
	fprintf(o, "by category.\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-a {delta,rsum,...}   Names of steppers: comma-separated, one or more of:\n");
	for (int i = 0; i < step_lookup_table_length; i++) {
		fprintf(o, "  %-8s %s\n", step_lookup_table[i].name, step_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c} Value-field names on which to compute statistics\n");
	fprintf(o, "-g {d,e,f} Optional group-by-field names\n");
	fprintf(o, "-F         Computes integerable things (e.g. counter) in floating point.\n");
	fprintf(o, "-d {x,y,z} Weights for ewma. 1 means current sample gets all weight (no\n");
	fprintf(o, "           smoothing), near under under 1 is light smoothing, near over 0 is\n");
	fprintf(o, "           heavy smoothing. Multiple weights may be specified, e.g.\n");
	fprintf(o, "           \"%s %s -a ewma -f sys_load -d 0.01,0.1,0.9\". Default if omitted\n", argv0, verb);
	fprintf(o, "           is \"-d %s\".\n", DEFAULT_STRING_ALPHA);
	fprintf(o, "-o {a,b,c} Custom suffixes for EWMA output fields. If omitted, these default to\n");
	fprintf(o, "           the -d values. If supplied, the number of -o values must be the same\n");
	fprintf(o, "           as the number of -d values.\n");
	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s -a rsum -f request_size\n", argv0, verb);
	fprintf(o, "  %s %s -a delta -f request_size -g hostname\n", argv0, verb);
	fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -f x,y\n", argv0, verb);
	fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y\n", argv0, verb);
	fprintf(o, "  %s %s -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name\n", argv0, verb);
	fprintf(o, "\n");
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html#filter or\n");
	fprintf(o, "https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average\n");
	fprintf(o, "for more information on EWMA.\n");
}

static mapper_t* mapper_step_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t*         pstepper_names        = NULL;
	string_array_t* pvalue_field_names    = NULL;
	slls_t*         pgroup_by_field_names = slls_alloc();
	slls_t*         pstring_alphas        = slls_single_no_free(DEFAULT_STRING_ALPHA);
	slls_t*         pewma_suffixes        = NULL;
	int             allow_int_float       = TRUE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate,  "-a", &pstepper_names);
	ap_define_string_array_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate,  "-g", &pgroup_by_field_names);
	ap_define_string_list_flag(pstate,  "-d", &pstring_alphas);
	ap_define_string_list_flag(pstate,  "-o", &pewma_suffixes);
	ap_define_false_flag(pstate,        "-F", &allow_int_float);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_step_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pstepper_names == NULL || pvalue_field_names == NULL) {
		mapper_step_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (pstring_alphas != NULL && pewma_suffixes != NULL) {
		if (pewma_suffixes->length != pstring_alphas->length) {
			mapper_step_usage(stderr, argv[0], verb);
			return NULL;
		}
	}

	return mapper_step_alloc(pstate, pstepper_names, pvalue_field_names, pgroup_by_field_names,
		allow_int_float, pstring_alphas, pewma_suffixes);
}

// ----------------------------------------------------------------
static mapper_t* mapper_step_alloc(ap_state_t* pargp, slls_t* pstepper_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int allow_int_float, slls_t* pstring_alphas, slls_t* pewma_suffixes)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_step_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_step_state_t));

	pstate->pargp                 = pargp;
	pstate->pstepper_names        = pstepper_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pvalue_field_values   = string_array_alloc(pvalue_field_names->length);
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->groups                = lhmslv_alloc();
	pstate->allow_int_float       = allow_int_float;
	pstate->pstring_alphas        = pstring_alphas;
	pstate->pewma_suffixes        = pewma_suffixes;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_step_process;
	pmapper->pfree_func    = mapper_step_free;

	return pmapper;
}

static void mapper_step_free(mapper_t* pmapper, context_t* _) {
	mapper_step_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pstepper_names);
	string_array_free(pstate->pvalue_field_names);
	string_array_free(pstate->pvalue_field_values);
	slls_free(pstate->pgroup_by_field_names);
	slls_free(pstate->pstring_alphas);
	slls_free(pstate->pewma_suffixes);

	// lhmslv_free and lhmsv_free will free the hashmap keys; we need to free
	// the void-star hashmap values.
	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
		lhmsv_t* pgroup_to_acc_field = pa->pvvalue;
		for (lhmsve_t* pb = pgroup_to_acc_field->phead; pb != NULL; pb = pb->pnext) {
			lhmsv_t* pacc_field_to_acc_state = pb->pvvalue;
			for (lhmsve_t* pc = pacc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
				step_t* pstep = pc->pvvalue;
				pstep->pfree_func(pstep);
			}
			lhmsv_free(pacc_field_to_acc_state);
		}
		lhmsv_free(pgroup_to_acc_field);
	}
	lhmslv_free(pstate->groups);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_step_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_step_state_t* pstate = pvstate;
	if (pinrec == NULL)
		return sllv_single(NULL);

	// ["s", "t"]
	mlr_reference_values_from_record_into_string_array(pinrec, pstate->pvalue_field_names, pstate->pvalue_field_values);
	slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

	if (pgroup_by_field_values == NULL) {
		slls_free(pgroup_by_field_values);
		return sllv_single(pinrec);
	}

	lhmsv_t* pgroup_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (pgroup_to_acc_field == NULL) {
		pgroup_to_acc_field = lhmsv_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), pgroup_to_acc_field, FREE_ENTRY_KEY);
	}
	slls_free(pgroup_by_field_values);

	// for x=1 and y=2
	int n = pstate->pvalue_field_names->length;
	for (int i = 0; i < n; i++) {
		char* value_field_name = pstate->pvalue_field_names->strings[i];
		char* value_field_sval = pstate->pvalue_field_values->strings[i];
		if (value_field_sval == NULL) // Key not present
			continue;

		int have_dval = FALSE;
		int have_nval = FALSE;
		double value_field_dval = -999.0;
		mv_t   value_field_nval = mv_absent();

		lhmsv_t* pacc_field_to_acc_state = lhmsv_get(pgroup_to_acc_field, value_field_name);
		if (pacc_field_to_acc_state == NULL) {
			pacc_field_to_acc_state = lhmsv_alloc();
			lhmsv_put(pgroup_to_acc_field, value_field_name, pacc_field_to_acc_state, NO_FREE);
		}

		// for "delta", "rsum"
		sllse_t* pc = pstate->pstepper_names->phead;
		for ( ; pc != NULL; pc = pc->pnext) {
			char* step_name = pc->value;
			step_t* pstep = lhmsv_get(pacc_field_to_acc_state, step_name);
			if (pstep == NULL) {
				pstep = make_step(step_name, value_field_name, pstate->allow_int_float,
					pstate->pstring_alphas, pstate->pewma_suffixes);
				if (pstep == NULL) {
					fprintf(stderr, "mlr step: stepper \"%s\" not found.\n",
						step_name);
					exit(1);
				}
				lhmsv_put(pacc_field_to_acc_state, step_name, pstep, NO_FREE);
			}

			if (*value_field_sval == 0) { // Key present with null value
				if (pstep->pzprocess_func != NULL) {
					pstep->pzprocess_func(pstep->pvstate, pinrec);
				}
			} else {

				if (pstep->pdprocess_func != NULL) {
					if (!have_dval) {
						value_field_dval = mlr_double_from_string_or_die(value_field_sval);
						have_dval = TRUE;
					}
					pstep->pdprocess_func(pstep->pvstate, value_field_dval, pinrec);
				}

				if (pstep->pnprocess_func != NULL) {
					if (!have_nval) {
						value_field_nval = pstate->allow_int_float
							? mv_scan_number_or_die(value_field_sval)
							: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
						have_nval = TRUE;
					}
					pstep->pnprocess_func(pstep->pvstate, &value_field_nval, pinrec);
				}

				if (pstep->psprocess_func != NULL) {
					pstep->psprocess_func(pstep->pvstate, value_field_sval, pinrec);
				}
			}
		}
	}
	return sllv_single(pinrec);
}

static step_t* make_step(char* step_name, char* input_field_name, int allow_int_float,
	slls_t* pstring_alphas, slls_t* pewma_suffixes)
{
	for (int i = 0; i < step_lookup_table_length; i++)
		if (streq(step_name, step_lookup_table[i].name))
			return step_lookup_table[i].palloc_func(input_field_name, allow_int_float,
				pstring_alphas, pewma_suffixes);
	return NULL;
}

// ----------------------------------------------------------------
typedef struct _step_delta_state_t {
	mv_t  prev;
	char* output_field_name;
	int   allow_int_float;
} step_delta_state_t;
static void step_delta_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
	step_delta_state_t* pstate = pvstate;
	mv_t delta;
	if (mv_is_null(&pstate->prev)) {
		delta = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	} else {
		delta = x_xx_minus_func(pnumv, &pstate->prev);
	}
	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&delta), FREE_ENTRY_VALUE);
	pstate->prev = *pnumv;
}
static void step_delta_zprocess(void* pvstate, lrec_t* prec) {
	step_delta_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_delta_free(step_t* pstep) {
	step_delta_state_t* pstate = pstep->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_delta_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_delta_state_t* pstate = mlr_malloc_or_die(sizeof(step_delta_state_t));
	pstate->prev = mv_absent();
	pstate->allow_int_float = allow_int_float;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_delta");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = step_delta_nprocess;
	pstep->psprocess_func = NULL;
	pstep->pzprocess_func = step_delta_zprocess;
	pstep->pfree_func     = step_delta_free;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_shift_state_t {
	char* prev;
	char* output_field_name;
	int   allow_int_float;
} step_shift_state_t;
static void step_shift_sprocess(void* pvstate, char* strv, lrec_t* prec) {
	step_shift_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, pstate->prev, FREE_ENTRY_VALUE);
	pstate->prev = mlr_strdup_or_die(strv);
}
static void step_shift_zprocess(void* pvstate, lrec_t* prec) {
	step_shift_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_shift_free(step_t* pstep) {
	step_shift_state_t* pstate = pstep->pvstate;
	free(pstate->prev);
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_shift_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_shift_state_t* pstate = mlr_malloc_or_die(sizeof(step_shift_state_t));
	pstate->prev = mlr_strdup_or_die("");
	pstate->allow_int_float = allow_int_float;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_shift");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = NULL;
	pstep->psprocess_func = step_shift_sprocess;
	pstep->pzprocess_func = step_shift_zprocess;
	pstep->pfree_func     = step_shift_free;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_from_first_state_t {
	mv_t  first;
	char* output_field_name;
	int   allow_int_float;
} step_from_first_state_t;
static void step_from_first_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
	step_from_first_state_t* pstate = pvstate;
	mv_t from_first;
	if (mv_is_null(&pstate->first)) {
		from_first = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
		pstate->first = *pnumv;
	} else {
		from_first = x_xx_minus_func(pnumv, &pstate->first);
	}
	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&from_first), FREE_ENTRY_VALUE);
}
static void step_from_first_zprocess(void* pvstate, lrec_t* prec) {
	step_from_first_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_from_first_free(step_t* pstep) {
	step_from_first_state_t* pstate = pstep->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_from_first_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_from_first_state_t* pstate = mlr_malloc_or_die(sizeof(step_from_first_state_t));
	pstate->first = mv_absent();
	pstate->allow_int_float = allow_int_float;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_from_first");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = step_from_first_nprocess;
	pstep->psprocess_func = NULL;
	pstep->pzprocess_func = step_from_first_zprocess;
	pstep->pfree_func     = step_from_first_free;
	return pstep;
}

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
		FREE_ENTRY_VALUE);
	pstate->prev = fltv;
}
static void step_ratio_zprocess(void* pvstate, lrec_t* prec) {
	step_ratio_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_ratio_free(step_t* pstep) {
	step_ratio_state_t* pstate = pstep->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_ratio_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_ratio_state_t* pstate = mlr_malloc_or_die(sizeof(step_ratio_state_t));
	pstate->prev          = -999.0;
	pstate->have_prev     = FALSE;
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_ratio");

	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = step_ratio_dprocess;
	pstep->pnprocess_func = NULL;
	pstep->psprocess_func = NULL;
	pstep->pzprocess_func = step_ratio_zprocess;
	pstep->pfree_func     = step_ratio_free;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_rsum_state_t {
	mv_t   rsum;
	char*  output_field_name;
	int    allow_int_float;
} step_rsum_state_t;
static void step_rsum_nprocess(void* pvstate, mv_t* pnumv, lrec_t* prec) {
	step_rsum_state_t* pstate = pvstate;
	pstate->rsum = x_xx_plus_func(&pstate->rsum, pnumv);
	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&pstate->rsum),
		FREE_ENTRY_VALUE);
}
static void step_rsum_zprocess(void* pvstate, lrec_t* prec) {
	step_rsum_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_rsum_free(step_t* pstep) {
	step_rsum_state_t* pstate = pstep->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_rsum_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_rsum_state_t* pstate = mlr_malloc_or_die(sizeof(step_rsum_state_t));
	pstate->allow_int_float = allow_int_float;
	pstate->rsum = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_rsum");
	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = step_rsum_nprocess;
	pstep->psprocess_func = NULL;
	pstep->pzprocess_func = step_rsum_zprocess;
	pstep->pfree_func     = step_rsum_free;
	return pstep;
}

// ----------------------------------------------------------------
typedef struct _step_counter_state_t {
	mv_t counter;
	mv_t one;
	char*  output_field_name;
} step_counter_state_t;
static void step_counter_sprocess(void* pvstate, char* strv, lrec_t* prec) {
	step_counter_state_t* pstate = pvstate;
	pstate->counter = x_xx_plus_func(&pstate->counter, &pstate->one);
	lrec_put(prec, pstate->output_field_name, mv_alloc_format_val(&pstate->counter),
		FREE_ENTRY_VALUE);
}
static void step_counter_zprocess(void* pvstate, lrec_t* prec) {
	step_counter_state_t* pstate = pvstate;
	lrec_put(prec, pstate->output_field_name, "", NO_FREE);
}
static void step_counter_free(step_t* pstep) {
	step_counter_state_t* pstate = pstep->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstep);
}
static step_t* step_counter_alloc(char* input_field_name, int allow_int_float, slls_t* unused1, slls_t* unused2) {
	step_t* pstep = mlr_malloc_or_die(sizeof(step_t));
	step_counter_state_t* pstate = mlr_malloc_or_die(sizeof(step_counter_state_t));
	pstate->counter = allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->one     = allow_int_float ? mv_from_int(1LL) : mv_from_float(1.0);
	pstate->output_field_name = mlr_paste_2_strings(input_field_name, "_counter");

	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = NULL;
	pstep->pnprocess_func = NULL;
	pstep->psprocess_func = step_counter_sprocess;
	pstep->pzprocess_func = step_counter_zprocess;
	pstep->pfree_func     = step_counter_free;
	return pstep;
}

// ----------------------------------------------------------------
// https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
typedef struct _step_ewma_state_t {
	int     num_alphas;
	double* alphas;
	double* alphacompls;
	double* prevs;
	int     have_prevs;
	char**  output_field_names;
} step_ewma_state_t;
static void step_ewma_dprocess(void* pvstate, double fltv, lrec_t* prec) {
	step_ewma_state_t* pstate = pvstate;
	if (!pstate->have_prevs) {
		for (int i = 0; i < pstate->num_alphas; i++) {
			lrec_put(prec, pstate->output_field_names[i], mlr_alloc_string_from_double(fltv, MLR_GLOBALS.ofmt),
				FREE_ENTRY_VALUE);
			pstate->prevs[i] = fltv;
		}
		pstate->have_prevs = TRUE;
	} else {
		for (int i = 0; i < pstate->num_alphas; i++) {
			double curr = fltv;
			curr = pstate->alphas[i] * curr + pstate->alphacompls[i] * pstate->prevs[i];
			lrec_put(prec, pstate->output_field_names[i], mlr_alloc_string_from_double(curr, MLR_GLOBALS.ofmt),
				FREE_ENTRY_VALUE);
			pstate->prevs[i] = curr;
		}
	}
}
static void step_ewma_zprocess(void* pvstate, lrec_t* prec) {
	step_ewma_state_t* pstate = pvstate;
	for (int i = 0; i < pstate->num_alphas; i++)
		lrec_put(prec, pstate->output_field_names[i], "", NO_FREE);
}
static void step_ewma_free(step_t* pstep) {
	step_ewma_state_t* pstate = pstep->pvstate;
	for (int i = 0; i < pstate->num_alphas; i++) {
		free(pstate->output_field_names[i]);
	}
	free(pstate->alphas);
	free(pstate->alphacompls);
	free(pstate->prevs);
	free(pstate->output_field_names);
	free(pstate);
	free(pstep);
}

static step_t* step_ewma_alloc(char* input_field_name, int unused, slls_t* pstring_alphas, slls_t* pewma_suffixes) {
	step_t* pstep              = mlr_malloc_or_die(sizeof(step_t));

	step_ewma_state_t* pstate  = mlr_malloc_or_die(sizeof(step_ewma_state_t));
	int n                      = pstring_alphas->length;
	pstate->num_alphas         = n;
	pstate->alphas             = mlr_malloc_or_die(n * sizeof(double));
	pstate->alphacompls        = mlr_malloc_or_die(n * sizeof(double));
	pstate->prevs              = mlr_malloc_or_die(n * sizeof(double));
	pstate->have_prevs         = FALSE;
	pstate->output_field_names = mlr_malloc_or_die(n * sizeof(char*));
	slls_t* psuffixes = (pewma_suffixes == NULL) ? pstring_alphas : pewma_suffixes;
	sllse_t* pe = pstring_alphas->phead;
	sllse_t* pf = psuffixes->phead;
	for (int i = 0; i < n; i++, pe = pe->pnext, pf = pf->pnext) {
		char* string_alpha     = pe->value;
		char* suffix           = pf->value;
		pstate->alphas[i]      = mlr_double_from_string_or_die(string_alpha);
		pstate->alphacompls[i] = 1.0 - pstate->alphas[i];
		pstate->prevs[i]       = 0.0;
		pstate->output_field_names[i] = mlr_paste_3_strings(input_field_name, "_ewma_", suffix);
	}
	pstate->have_prevs = FALSE;

	pstep->pvstate        = (void*)pstate;
	pstep->pdprocess_func = step_ewma_dprocess;
	pstep->pnprocess_func = NULL;
	pstep->psprocess_func = NULL;
	pstep->pzprocess_func = step_ewma_zprocess;
	pstep->pfree_func     = step_ewma_free;
	return pstep;
}
