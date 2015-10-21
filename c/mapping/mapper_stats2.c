#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/mlrmath.h"
#include "lib/mlrstat.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhms2v.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "containers/dvector.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DO_CORR       0x11
#define DO_COV        0x22
#define DO_COVX       0x33
#define DO_LINREG_PCA 0x44

// ----------------------------------------------------------------
typedef void stats2_ingest_func_t(void* pvstate, double x, double y);
typedef void   stats2_emit_func_t(void* pvstate, char* name1, char* name2, lrec_t* poutrec);
typedef void    stats2_fit_func_t(void* pvstate, double x, double y, lrec_t* poutrec);

typedef struct _stats2_t {
	void* pvstate;
	stats2_ingest_func_t* pingest_func;
	stats2_emit_func_t*   pemit_func;
	stats2_fit_func_t*    pfit_func;
} stats2_t;

typedef struct _mapper_stats2_state_t {
	slls_t* paccumulator_names;
	slls_t* pvalue_field_name_pairs;
	slls_t* pgroup_by_field_names;

	lhmslv_t* acc_groups;
	lhmslv_t* record_groups;
	int     do_verbose;
	int     do_iterative_stats;
	int     do_hold_and_fit;
} mapper_stats2_state_t;

typedef stats2_t* stats2_alloc_func_t(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);

// ----------------------------------------------------------------
static void      mapper_stats2_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_stats2_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_stats2_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_name_pairs,
	slls_t* pgroup_by_field_names, int do_verbose, int do_iterative_stats, int do_hold_and_fit);
static void      mapper_stats2_free(void* pvstate);
static sllv_t*   mapper_stats2_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_stats2_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats2_state_t* pstate);
static sllv_t*   mapper_stats2_emit_all(mapper_stats2_state_t* pstate);
static void      mapper_stats2_emit(mapper_stats2_state_t* pstate, lrec_t* pinrec,
	char* value_field_name_1, char* value_field_name_2, lhmsv_t* acc_fields_to_acc_state);
static sllv_t*   mapper_stats2_fit_all(mapper_stats2_state_t* pstate);

static stats2_t* make_stats2            (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_linreg_pca_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_linreg_ols_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_r2_alloc        (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_logireg_alloc   (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_corr_cov_alloc  (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_which, int do_verbose);
static stats2_t* stats2_corr_alloc      (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_cov_alloc       (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);
static stats2_t* stats2_covx_alloc      (char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose);

// ----------------------------------------------------------------
typedef struct _stats2_lookup_t {
	char* name;
	stats2_alloc_func_t* pnew_func;
	char* desc;
} stats2_lookup_t;
static stats2_lookup_t stats2_lookup_table[] = {
	{"linreg-pca", stats2_linreg_pca_alloc, "Linear regression using principal component analysis"},
	{"linreg-ols", stats2_linreg_ols_alloc, "Linear regression using ordinary least squares"},
	{"r2",         stats2_r2_alloc,         "Quality metric for linreg-ols (linreg-pca emits its own)"},
	{"logireg",    stats2_logireg_alloc,    "Logistic regression"},
	{"corr",       stats2_corr_alloc,       "Sample correlation"},
	{"cov",        stats2_cov_alloc,        "Sample covariance"},
	{"covx",       stats2_covx_alloc,       "Sample-covariance matrix"},
};
static int stats2_lookup_table_length = sizeof(stats2_lookup_table) / sizeof(stats2_lookup_table[0]);

// ----------------------------------------------------------------
mapper_setup_t mapper_stats2_setup = {
	.verb = "stats2",
	.pusage_func = mapper_stats2_usage,
	.pparse_func = mapper_stats2_parse_cli
};

// ----------------------------------------------------------------
static void mapper_stats2_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-a {linreg-ols,corr,...}  Names of accumulators: one or more of:\n");
	for (int i = 0; i < stats2_lookup_table_length; i++) {
		fprintf(o, "  %-12s %s\n", stats2_lookup_table[i].name, stats2_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c,d}   Value-field name-pairs on which to compute statistics.\n");
	fprintf(o, "               There must be an even number of names.\n");
	fprintf(o, "-g {e,f,g}     Optional group-by-field names.\n");
	fprintf(o, "-v             Print additional output for linreg-pca.\n");
	fprintf(o, "-s             Print iterative stats. Useful in tail -f contexts (in which\n");
	fprintf(o, "               case please avoid pprint-format output since end of input\n");
	fprintf(o, "               stream will never be seen).\n");
	fprintf(o, "--fit          Rather than printing regression parameters, applies them to\n");
	fprintf(o, "               the input data to compute new fit fields. All input records are\n");
	fprintf(o, "               held in memory until end of input stream. Has effect only for\n");
	fprintf(o, "               linreg-ols, linreg-pca, and logireg.\n");
	fprintf(o, "Only one of -s or --fit may be used.\n");
	fprintf(o, "Example: %s %s -a linreg-pca -f x,y\n", argv0, verb);
	fprintf(o, "Example: %s %s -a linreg-ols,r2 -f x,y -g size,shape\n", argv0, verb);
	fprintf(o, "Example: %s %s -a corr -f x,y\n", argv0, verb);
}

static mapper_t* mapper_stats2_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* paccumulator_names    = NULL;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();
	int     do_verbose            = FALSE;
	int     do_iterative_stats    = FALSE;
	int     do_hold_and_fit       = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-a", &paccumulator_names);
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate,        "-v", &do_verbose);
	ap_define_true_flag(pstate,        "-s", &do_iterative_stats);
	ap_define_true_flag(pstate,        "--fit", &do_hold_and_fit);
	// xxx abend here or elsewhere if hold & fit requested for non-supporting accumulator

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_stats2_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (do_iterative_stats && do_hold_and_fit) {
		mapper_stats2_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (paccumulator_names == NULL || pvalue_field_names == NULL) {
		mapper_stats2_usage(stderr, argv[0], verb);
		return NULL;
	}
	if ((pvalue_field_names->length % 2) != 0) {
		mapper_stats2_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_stats2_alloc(paccumulator_names, pvalue_field_names, pgroup_by_field_names,
		do_verbose, do_iterative_stats, do_hold_and_fit);
}

// ----------------------------------------------------------------
static mapper_t* mapper_stats2_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_name_pairs,
	slls_t* pgroup_by_field_names, int do_verbose, int do_iterative_stats, int do_hold_and_fit)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats2_state_t* pstate   = mlr_malloc_or_die(sizeof(mapper_stats2_state_t));
	pstate->paccumulator_names      = paccumulator_names;
	pstate->pvalue_field_name_pairs = pvalue_field_name_pairs; // caller validates length is even
	pstate->pgroup_by_field_names   = pgroup_by_field_names;
	pstate->acc_groups              = lhmslv_alloc();
	pstate->record_groups           = lhmslv_alloc();
	pstate->do_verbose              = do_verbose;
	pstate->do_iterative_stats      = do_iterative_stats;
	pstate->do_hold_and_fit         = do_hold_and_fit;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_stats2_process;
	pmapper->pfree_func    = mapper_stats2_free;

	return pmapper;
}

static void mapper_stats2_free(void* pvstate) {
	mapper_stats2_state_t* pstate = pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_name_pairs);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->acc_groups);
	lhmslv_free(pstate->record_groups);
}

// ================================================================
// Given: accumulate corr,cov on values x,y group by a,b.
// Example input:       Example output:
//   a b x y            a b x_corr x_cov y_corr y_cov
//   s t 1 2            s t 2       6    2      8
//   u v 3 4            u v 1       3    1      4
//   s t 5 6            u w 1       7    1      9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   ["s","t"] : {                    <--- group-by field names
//     ["x","y"] : {                  <--- value field names
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
// }
// ================================================================

// In the iterative case, add to the current record its current group's stats fields.
// In the non-iteratiive case, produce output only at end of input stream.
static sllv_t* mapper_stats2_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats2_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_stats2_ingest(pinrec, pctx, pstate);
		if (pstate->do_iterative_stats) {
			// The input record is modified in this case, with new fields appended
			return sllv_single(pinrec);
		} else if (pstate->do_hold_and_fit) {
			// The input record is held by the ingestor
			return NULL;
		} else {
			lrec_free(pinrec);
			return NULL;
		}
	} else if (!pstate->do_iterative_stats) {
		if (!pstate->do_hold_and_fit) {
			return mapper_stats2_emit_all(pstate);
		} else {
			return mapper_stats2_fit_all(pstate);
		}
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static void mapper_stats2_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats2_state_t* pstate) {
	// ["s", "t"]
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
	if (pgroup_by_field_values->length != pstate->pgroup_by_field_names->length) {
		slls_free(pgroup_by_field_values);
		return;
	}

	lhms2v_t* group_to_acc_field = lhmslv_get(pstate->acc_groups, pgroup_by_field_values);
	if (group_to_acc_field == NULL) {
		group_to_acc_field = lhms2v_alloc();
		lhmslv_put(pstate->acc_groups, slls_copy(pgroup_by_field_values), group_to_acc_field);
	}

	if (pstate->do_hold_and_fit) { // Retain the input record in memory, for fitting and delivery at end of stream
		sllv_t* group_to_records = lhmslv_get(pstate->record_groups, pgroup_by_field_values);
		if (group_to_records == NULL) {
			group_to_records = sllv_alloc();
			lhmslv_put(pstate->record_groups, slls_copy(pgroup_by_field_values), group_to_records);
		}
		sllv_add(group_to_records, pinrec);
	}

	// for [["x","y"]]
	for (sllse_t* pa = pstate->pvalue_field_name_pairs->phead; pa != NULL; pa = pa->pnext->pnext) {
		char* value_field_name_1 = pa->value;
		char* value_field_name_2 = pa->pnext->value;

		lhmsv_t* acc_fields_to_acc_state = lhms2v_get(group_to_acc_field, value_field_name_1, value_field_name_2);
		if (acc_fields_to_acc_state == NULL) {
			acc_fields_to_acc_state = lhmsv_alloc();
			lhms2v_put(group_to_acc_field, value_field_name_1, value_field_name_2, acc_fields_to_acc_state);
		}

		char* sval1 = lrec_get(pinrec, value_field_name_1);
		if (sval1 == NULL)
			continue;

		char* sval2 = lrec_get(pinrec, value_field_name_2);
		if (sval2 == NULL)
			continue;

		// for ["corr", "cov"]
		sllse_t* pc = pstate->paccumulator_names->phead;
		for ( ; pc != NULL; pc = pc->pnext) {
			char* stats2_name = pc->value;
			stats2_t* pstats2 = lhmsv_get(acc_fields_to_acc_state, stats2_name);
			if (pstats2 == NULL) {
				pstats2 = make_stats2(value_field_name_1, value_field_name_2, stats2_name, pstate->do_verbose);
				if (pstats2 == NULL) {
					fprintf(stderr, "mlr stats2: accumulator \"%s\" not found.\n",
						stats2_name);
					exit(1);
				}
				lhmsv_put(acc_fields_to_acc_state, stats2_name, pstats2);
			}

			double dval1 = mlr_double_from_string_or_die(sval1);
			double dval2 = mlr_double_from_string_or_die(sval2);
			pstats2->pingest_func(pstats2->pvstate, dval1, dval2);
		}
		if (pstate->do_iterative_stats) {
			mapper_stats2_emit(pstate, pinrec, value_field_name_1, value_field_name_2,
				acc_fields_to_acc_state);
		}
	}

	slls_free(pgroup_by_field_values);
}

// ----------------------------------------------------------------
static sllv_t* mapper_stats2_emit_all(mapper_stats2_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->acc_groups->phead; pa != NULL; pa = pa->pnext) {
		lrec_t* poutrec = lrec_unbacked_alloc();

		// Add in a=s,b=t fields:
		slls_t* pgroup_by_field_values = pa->key;
		sllse_t* pb = pstate->pgroup_by_field_names->phead;
		sllse_t* pc =         pgroup_by_field_values->phead;
		for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
			lrec_put(poutrec, pb->value, pc->value, 0);
		}

		// Add in fields such as x_y_corr, etc.
		lhms2v_t* group_to_acc_field = pa->pvvalue;

		// For "x","y"
		for (lhms2ve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
			char*    value_field_name_1 = pd->key1;
			char*    value_field_name_2 = pd->key2;
			lhmsv_t* acc_fields_to_acc_state = pd->pvvalue;

			mapper_stats2_emit(pstate, poutrec, value_field_name_1, value_field_name_2,
				acc_fields_to_acc_state);

			// For "corr", "linreg"
			for (lhmsve_t* pe = acc_fields_to_acc_state->phead; pe != NULL; pe = pe->pnext) {
				stats2_t* pstats2 = pe->pvvalue;
				pstats2->pemit_func(pstats2->pvstate, value_field_name_1, value_field_name_2, poutrec);
			}
		}

		sllv_add(poutrecs, poutrec);
	}
	sllv_add(poutrecs, NULL);
	return poutrecs;
}

static void mapper_stats2_emit(mapper_stats2_state_t* pstate, lrec_t* poutrec,
	char* value_field_name_1, char* value_field_name_2, lhmsv_t* acc_fields_to_acc_state)
{
	// For "corr", "linreg"
	for (lhmsve_t* pe = acc_fields_to_acc_state->phead; pe != NULL; pe = pe->pnext) {
		stats2_t* pstats2 = pe->pvvalue;
		pstats2->pemit_func(pstats2->pvstate, value_field_name_1, value_field_name_2, poutrec);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_stats2_fit_all(mapper_stats2_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->acc_groups->phead; pa != NULL; pa = pa->pnext) {
		slls_t* pgroup_by_field_values = pa->key;
		sllv_t* precords = lhmslv_get(pstate->record_groups, pgroup_by_field_values);

		while (precords->phead) {
			lrec_t* prec = sllv_pop(precords);

			lhms2v_t* group_to_acc_field = pa->pvvalue;

			// For "x","y"
			for (lhms2ve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
				char*    value_field_name_1 = pd->key1;
				char*    value_field_name_2 = pd->key2;
				lhmsv_t* acc_fields_to_acc_state = pd->pvvalue;

				// For "linreg-ols", "logireg"
				for (lhmsve_t* pe = acc_fields_to_acc_state->phead; pe != NULL; pe = pe->pnext) {
					stats2_t* pstats2 = pe->pvvalue;
					if (pstats2->pfit_func != NULL) {
						char* sx = lrec_get(prec, value_field_name_1);
						char* sy = lrec_get(prec, value_field_name_2);
						if (sx != NULL && sy != NULL) {
							double x = mlr_double_from_string_or_die(sx);
							double y = mlr_double_from_string_or_die(sy);
							pstats2->pfit_func(pstats2->pvstate, x, y, prec);
						}
					}
				}
			}

			sllv_add(poutrecs, prec);
		}
	}
	sllv_add(poutrecs, NULL);
	return poutrecs;
}

// ================================================================
// Given: accumulate corr,cov on values x,y group by a,b.
// Example input:       Example output:
//   a b x y            a b x_corr x_cov y_corr y_cov
//   s t 1 2            s t 2       6    2      8
//   u v 3 4            u v 1       3    1      4
//   s t 5 6            u w 1       7    1      9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   ["s","t"] : {                    <--- group-by field names
//     ["x","y"] : {                  <--- value field names
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","v"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
//   ["u","w"] : {
//     ["x","y"] : {
//       "corr" : stats2_corr_t object,
//       "cov"  : stats2_cov_t  object
//     }
//   },
// }
// ================================================================

// ----------------------------------------------------------------
static stats2_t* make_stats2(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	for (int i = 0; i < stats2_lookup_table_length; i++)
		if (streq(stats2_name, stats2_lookup_table[i].name))
			return stats2_lookup_table[i].pnew_func(value_field_name_1, value_field_name_2, stats2_name, do_verbose);
	return NULL;
}

// ----------------------------------------------------------------
typedef struct _stats2_linreg_ols_state_t {
	unsigned long long count;
	double sumx;
	double sumy;
	double sumx2;
	double sumxy;
	char*  m_output_field_name;
	char*  b_output_field_name;
	char*  n_output_field_name;

	char*  fit_output_field_name;
	int    fit_ready;
	double m;
	double b;

} stats2_linreg_ols_state_t;
static void stats2_linreg_ols_ingest(void* pvstate, double x, double y) {
	stats2_linreg_ols_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += x;
	pstate->sumy  += y;
	pstate->sumx2 += x*x;
	pstate->sumxy += x*y;
}

static void stats2_linreg_ols_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
	stats2_linreg_ols_state_t* pstate = pvstate;

	if (pstate->count < 2) {
		lrec_put(poutrec, pstate->m_output_field_name, "", LREC_FREE_ENTRY_KEY);
		lrec_put(poutrec, pstate->b_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double m, b;
		mlr_get_linear_regression_ols(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumxy, pstate->sumy, &m, &b);
		char* mval = mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt);
		char* bval = mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt);

		lrec_put(poutrec, pstate->m_output_field_name, mval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
		lrec_put(poutrec, pstate->b_output_field_name, bval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}

	char* nval = mlr_alloc_string_from_ll(pstate->count);
	lrec_put(poutrec, pstate->n_output_field_name, nval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}

static void stats2_linreg_ols_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
	stats2_linreg_ols_state_t* pstate = pvstate;

	if (!pstate->fit_ready) {
		mlr_get_linear_regression_ols(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumxy, pstate->sumy,
			&pstate->m, &pstate->b);
		pstate->fit_ready = TRUE;
	}

	if (pstate->count < 2) {
		lrec_put(poutrec, pstate->fit_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double yfit = pstate->m * x + pstate->b;
		char* sfit = mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->fit_output_field_name, sfit, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static stats2_t* stats2_linreg_ols_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_linreg_ols_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_linreg_ols_state_t));
	pstate->count = 0LL;
	pstate->sumx  = 0.0;
	pstate->sumy  = 0.0;
	pstate->sumx2 = 0.0;
	pstate->sumxy = 0.0;
	pstate->m_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_m");
	pstate->b_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_b");
	pstate->n_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_n");
	pstate->fit_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_ols_fit");
	pstate->fit_ready = FALSE;
	pstate->m         = -999.0;
	pstate->b         = -999.0;

	pstats2->pvstate = (void*)pstate;
	pstats2->pingest_func = stats2_linreg_ols_ingest;
	pstats2->pemit_func   = stats2_linreg_ols_emit;
	pstats2->pfit_func    = stats2_linreg_ols_fit;
	return pstats2;
}

// ----------------------------------------------------------------
#define LOGIREG_DVECTOR_INITIAL_SIZE 1024
typedef struct _stats2_logireg_state_t {
	dvector_t* pxs;
	dvector_t* pys;
	char*  m_output_field_name;
	char*  b_output_field_name;
	char*  n_output_field_name;
	char*  fit_output_field_name;
	int    fit_ready;
	double m;
	double b;
} stats2_logireg_state_t;
static void stats2_logireg_ingest(void* pvstate, double x, double y) {
	stats2_logireg_state_t* pstate = pvstate;
	dvector_append(pstate->pxs, x);
	dvector_append(pstate->pys, y);
}

static void stats2_logireg_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
	stats2_logireg_state_t* pstate = pvstate;

	if (pstate->pxs->size < 2) {
		lrec_put(poutrec, pstate->m_output_field_name, "", LREC_FREE_ENTRY_KEY);
		lrec_put(poutrec, pstate->b_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double m, b;
		mlr_logistic_regression(pstate->pxs->data, pstate->pys->data, pstate->pxs->size, &m, &b);
		char* mval = mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt);
		char* bval = mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt);

		lrec_put(poutrec, pstate->m_output_field_name, mval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
		lrec_put(poutrec, pstate->b_output_field_name, bval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}

	char* nval = mlr_alloc_string_from_ll(pstate->pxs->size);
	lrec_put(poutrec, pstate->n_output_field_name, nval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}

static void stats2_logireg_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
	stats2_logireg_state_t* pstate = pvstate;

	if (!pstate->fit_ready) {
		mlr_logistic_regression(pstate->pxs->data, pstate->pys->data, pstate->pxs->size, &pstate->m, &pstate->b);
		pstate->fit_ready = TRUE;
	}

	if (pstate->pxs->size < 2) {
		lrec_put(poutrec, pstate->fit_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double yfit = 1.0 / (1.0 + exp(-pstate->m*x - pstate->b));
		char* fitval = mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->fit_output_field_name, fitval, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static stats2_t* stats2_logireg_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_logireg_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_logireg_state_t));
	pstate->pxs = dvector_alloc(LOGIREG_DVECTOR_INITIAL_SIZE);
	pstate->pys = dvector_alloc(LOGIREG_DVECTOR_INITIAL_SIZE);
	pstate->m_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_m");
	pstate->b_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_b");
	pstate->n_output_field_name   = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_n");
	pstate->fit_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_logistic_fit");
	pstate->fit_ready = FALSE;
	pstate->m         = -999.0;
	pstate->b         = -999.0;

	pstats2->pvstate = (void*)pstate;
	pstats2->pingest_func = stats2_logireg_ingest;
	pstats2->pemit_func   = stats2_logireg_emit;
	pstats2->pfit_func    = stats2_logireg_fit;
	return pstats2;
}

// ----------------------------------------------------------------
// http://en.wikipedia.org/wiki/Pearson_product-moment_correlation_coefficient
// Alternatively, just use sqrt(corr) as defined above.

typedef struct _stats2_r2_state_t {
	unsigned long long count;
	double sumx;
	double sumy;
	double sumx2;
	double sumxy;
	double sumy2;
	char*  r2_output_field_name;
} stats2_r2_state_t;
static void stats2_r2_ingest(void* pvstate, double x, double y) {
	stats2_r2_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += x;
	pstate->sumy  += y;
	pstate->sumx2 += x*x;
	pstate->sumxy += x*y;
	pstate->sumy2 += y*y;
}
static void stats2_r2_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
	stats2_r2_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		lrec_put(poutrec, pstate->r2_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		unsigned long long n = pstate->count;
		double sumx  = pstate->sumx;
		double sumy  = pstate->sumy;
		double sumx2 = pstate->sumx2;
		double sumy2 = pstate->sumy2;
		double sumxy = pstate->sumxy;
		double numerator = n*sumxy - sumx*sumy;
		numerator = numerator * numerator;
		double denominator = (n*sumx2 - sumx*sumx) * (n*sumy2 - sumy*sumy);
		double output = numerator/denominator;
		char* val = mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->r2_output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}
static stats2_t* stats2_r2_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_r2_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_r2_state_t));
	pstate->count     = 0LL;
	pstate->sumx      = 0.0;
	pstate->sumy      = 0.0;
	pstate->sumx2     = 0.0;
	pstate->sumxy     = 0.0;
	pstate->sumy2     = 0.0;
	pstate->r2_output_field_name = mlr_paste_4_strings(value_field_name_1, "_", value_field_name_2, "_r2");

	pstats2->pvstate      = (void*)pstate;
	pstats2->pingest_func = stats2_r2_ingest;
	pstats2->pemit_func   = stats2_r2_emit;
	pstats2->pfit_func    = NULL;

	return pstats2;
}

// ----------------------------------------------------------------
// Corr(X,Y) = Cov(X,Y) / sigma_X sigma_Y.
typedef struct _stats2_corr_cov_state_t {
	unsigned long long count;
	double sumx;
	double sumy;
	double sumx2;
	double sumxy;
	double sumy2;
	int    do_which;
	int    do_verbose;

	char*  covx_00_output_field_name;
	char*  covx_01_output_field_name;
	char*  covx_10_output_field_name;
	char*  covx_11_output_field_name;

	char*  pca_m_output_field_name;
	char*  pca_b_output_field_name;
	char*  pca_n_output_field_name;
	char*  pca_q_output_field_name;
	char* pca_l1_output_field_name;
	char* pca_l2_output_field_name;
	char* pca_v11_output_field_name;
	char* pca_v12_output_field_name;
	char* pca_v21_output_field_name;
	char* pca_v22_output_field_name;
	char* pca_fit_output_field_name;
	int   fit_ready;
	double m;
	double b;
	double q;

	char*  corr_output_field_name;
	char*   cov_output_field_name;

} stats2_corr_cov_state_t;
static void stats2_corr_cov_ingest(void* pvstate, double x, double y) {
	stats2_corr_cov_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += x;
	pstate->sumy  += y;
	pstate->sumx2 += x*x;
	pstate->sumxy += x*y;
	pstate->sumy2 += y*y;
}

static void stats2_corr_cov_emit(void* pvstate, char* name1, char* name2, lrec_t* poutrec) {
	stats2_corr_cov_state_t* pstate = pvstate;
	if (pstate->do_which == DO_COVX) {
		char* key00 = pstate->covx_00_output_field_name;
		char* key01 = pstate->covx_01_output_field_name;
		char* key10 = pstate->covx_10_output_field_name;
		char* key11 = pstate->covx_11_output_field_name;
		if (pstate->count < 2LL) {
			lrec_put(poutrec, key00, "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, key01, "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, key10, "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, key11, "", LREC_FREE_ENTRY_KEY);
		} else {
			double Q[2][2];
			mlr_get_cov_matrix(pstate->count,
				pstate->sumx, pstate->sumx2, pstate->sumy, pstate->sumy2, pstate->sumxy, Q);
			char* val00 = mlr_alloc_string_from_double(Q[0][0], MLR_GLOBALS.ofmt);
			char* val01 = mlr_alloc_string_from_double(Q[0][1], MLR_GLOBALS.ofmt);
			char* val10 = mlr_alloc_string_from_double(Q[1][0], MLR_GLOBALS.ofmt);
			char* val11 = mlr_alloc_string_from_double(Q[1][1], MLR_GLOBALS.ofmt);
			lrec_put(poutrec, key00, val00, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			lrec_put(poutrec, key01, val01, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			lrec_put(poutrec, key10, val10, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			lrec_put(poutrec, key11, val11, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
		}

	} else if (pstate->do_which == DO_LINREG_PCA) {
		char* keym   = pstate->pca_m_output_field_name;
		char* keyb   = pstate->pca_b_output_field_name;
		char* keyn   = pstate->pca_n_output_field_name;
		char* keyq   = pstate->pca_q_output_field_name;
		char* keyl1  = pstate->pca_l1_output_field_name;
		char* keyl2  = pstate->pca_l2_output_field_name;
		char* keyv11 = pstate->pca_v11_output_field_name;
		char* keyv12 = pstate->pca_v12_output_field_name;
		char* keyv21 = pstate->pca_v21_output_field_name;
		char* keyv22 = pstate->pca_v22_output_field_name;
		if (pstate->count < 2LL) {
			lrec_put(poutrec, keym,   "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, keyb,   "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, keyn,   "", LREC_FREE_ENTRY_KEY);
			lrec_put(poutrec, keyq,   "", LREC_FREE_ENTRY_KEY);
			if (pstate->do_verbose) {
				lrec_put(poutrec, keyl1,  "", LREC_FREE_ENTRY_KEY);
				lrec_put(poutrec, keyl2,  "", LREC_FREE_ENTRY_KEY);
				lrec_put(poutrec, keyv11, "", LREC_FREE_ENTRY_KEY);
				lrec_put(poutrec, keyv12, "", LREC_FREE_ENTRY_KEY);
				lrec_put(poutrec, keyv21, "", LREC_FREE_ENTRY_KEY);
				lrec_put(poutrec, keyv22, "", LREC_FREE_ENTRY_KEY);
			}
		} else {
			double Q[2][2];
			mlr_get_cov_matrix(pstate->count,
				pstate->sumx, pstate->sumx2, pstate->sumy, pstate->sumy2, pstate->sumxy, Q);

			double l1, l2;       // Eigenvalues
			double v1[2], v2[2]; // Eigenvectors
			mlr_get_real_symmetric_eigensystem(Q, &l1, &l2, v1, v2);

			double x_mean = pstate->sumx / pstate->count;
			double y_mean = pstate->sumy / pstate->count;
			double m, b, q;
			mlr_get_linear_regression_pca(l1, l2, v1, v2, x_mean, y_mean, &m, &b, &q);

			char free_flags = LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE;
			lrec_put(poutrec, keym, mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt), free_flags);
			lrec_put(poutrec, keyb, mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt), free_flags);
			lrec_put(poutrec, keyn, mlr_alloc_string_from_ll(pstate->count),           free_flags);
			lrec_put(poutrec, keyq, mlr_alloc_string_from_double(q, MLR_GLOBALS.ofmt), free_flags);
			if (pstate->do_verbose) {
				lrec_put(poutrec, keyl1,  mlr_alloc_string_from_double(l1,    MLR_GLOBALS.ofmt), free_flags);
				lrec_put(poutrec, keyl2,  mlr_alloc_string_from_double(l2,    MLR_GLOBALS.ofmt), free_flags);
				lrec_put(poutrec, keyv11, mlr_alloc_string_from_double(v1[0], MLR_GLOBALS.ofmt), free_flags);
				lrec_put(poutrec, keyv12, mlr_alloc_string_from_double(v1[1], MLR_GLOBALS.ofmt), free_flags);
				lrec_put(poutrec, keyv21, mlr_alloc_string_from_double(v2[0], MLR_GLOBALS.ofmt), free_flags);
				lrec_put(poutrec, keyv22, mlr_alloc_string_from_double(v2[1], MLR_GLOBALS.ofmt), free_flags);
			}
		}
	} else {
		char* key = (pstate->do_which == DO_CORR) ? pstate->corr_output_field_name : pstate->cov_output_field_name;
		if (pstate->count < 2LL) {
			lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
		} else {
			double output = mlr_get_cov(pstate->count, pstate->sumx, pstate->sumy, pstate->sumxy);
			if (pstate->do_which == DO_CORR) {
				double sigmax = sqrt(mlr_get_var(pstate->count, pstate->sumx, pstate->sumx2));
				double sigmay = sqrt(mlr_get_var(pstate->count, pstate->sumy, pstate->sumy2));
				output = output / sigmax / sigmay;
			}
			char* val = mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
			lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
		}
	}
}

static void linreg_pca_fit(void* pvstate, double x, double y, lrec_t* poutrec) {
	stats2_corr_cov_state_t* pstate = pvstate;

	if (!pstate->fit_ready) {
		double Q[2][2];
		mlr_get_cov_matrix(pstate->count,
			pstate->sumx, pstate->sumx2, pstate->sumy, pstate->sumy2, pstate->sumxy, Q);

		double l1, l2;       // Eigenvalues
		double v1[2], v2[2]; // Eigenvectors
		mlr_get_real_symmetric_eigensystem(Q, &l1, &l2, v1, v2);

		double x_mean = pstate->sumx / pstate->count;
		double y_mean = pstate->sumy / pstate->count;
		mlr_get_linear_regression_pca(l1, l2, v1, v2, x_mean, y_mean, &pstate->m, &pstate->b, &pstate->q);

		pstate->fit_ready = TRUE;
	}
	if (pstate->count < 2LL) {
		lrec_put(poutrec, pstate->pca_fit_output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double yfit = pstate->m * x + pstate->b;
		char free_flags = LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE;
		lrec_put(poutrec, pstate->pca_fit_output_field_name, mlr_alloc_string_from_double(yfit, MLR_GLOBALS.ofmt),
			free_flags);
	}
}

static stats2_t* stats2_corr_cov_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_which, int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_corr_cov_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_corr_cov_state_t));
	pstate->count      = 0LL;
	pstate->sumx       = 0.0;
	pstate->sumy       = 0.0;
	pstate->sumx2      = 0.0;
	pstate->sumxy      = 0.0;
	pstate->sumy2      = 0.0;
	pstate->do_which   = do_which;
	pstate->do_verbose = do_verbose;

	char* name1 = value_field_name_1;
	char* name2 = value_field_name_2;

	pstate->covx_00_output_field_name = mlr_paste_4_strings(name1, "_", name1, "_covx");
	pstate->covx_01_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_covx");
	pstate->covx_10_output_field_name = mlr_paste_4_strings(name2, "_", name1, "_covx");
	pstate->covx_11_output_field_name = mlr_paste_4_strings(name2, "_", name2, "_covx");

	pstate->pca_m_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_m");
	pstate->pca_b_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_b");
	pstate->pca_n_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_n");
	pstate->pca_q_output_field_name   = mlr_paste_4_strings(name1, "_", name2, "_pca_quality");
	pstate->pca_l1_output_field_name  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival1");
	pstate->pca_l2_output_field_name  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival2");
	pstate->pca_v11_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec11");
	pstate->pca_v12_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec12");
	pstate->pca_v21_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec21");
	pstate->pca_v22_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec22");
	pstate->pca_fit_output_field_name = mlr_paste_4_strings(name1, "_", name2, "_pca_fit");
	pstate->fit_ready = FALSE;
	pstate->m         = -999.0;
	pstate->b         = -999.0;

	pstate->corr_output_field_name    = mlr_paste_4_strings(name1, "_", name2, "_corr");
	pstate->cov_output_field_name     = mlr_paste_4_strings(name1, "_", name2, "_cov");

	pstats2->pvstate      = (void*)pstate;
	pstats2->pingest_func = stats2_corr_cov_ingest;
	pstats2->pemit_func   = stats2_corr_cov_emit;
	if (do_which == DO_LINREG_PCA)
		pstats2->pfit_func = linreg_pca_fit;
	else
		pstats2->pfit_func = NULL;

	return pstats2;
}
static stats2_t* stats2_corr_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_name, DO_CORR, do_verbose);
}
static stats2_t* stats2_cov_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_name, DO_COV, do_verbose);
}
static stats2_t* stats2_covx_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_name, DO_COVX, do_verbose);
}
static stats2_t* stats2_linreg_pca_alloc(char* value_field_name_1, char* value_field_name_2, char* stats2_name, int do_verbose) {
	return stats2_corr_cov_alloc(value_field_name_1, value_field_name_2, stats2_name, DO_LINREG_PCA, do_verbose);
}
