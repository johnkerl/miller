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
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DO_CORR       0x11
#define DO_COV        0x22
#define DO_COVX       0x33
#define DO_LINREG_PCA 0x44

// ----------------------------------------------------------------
typedef void stats2_ingest_func_t(void* pvstate, double x, double y);
typedef void stats2_emit_func_t(void* pvstate, char* name1, char* name2, lrec_t* poutrec);

typedef struct _stats2_t {
	void* pvstate;
	stats2_ingest_func_t* pingest_func;
	stats2_emit_func_t*   pemit_func;
} stats2_t;

typedef stats2_t* stats2_alloc_func_t(int do_verbose);

// ----------------------------------------------------------------
typedef struct _stats2_linreg_ols_state_t {
	unsigned long long count;
	double sumx;
	double sumy;
	double sumx2;
	double sumxy;
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
	double m, b;

	mlr_get_linear_regression_ols(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumxy, pstate->sumy, &m, &b);

	char* key = mlr_paste_4_strings(name1, "_", name2, "_ols_m");
	char* val = mlr_alloc_string_from_double(m, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);

	key = mlr_paste_4_strings(name1, "_", name2, "_ols_b");
	val = mlr_alloc_string_from_double(b, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);

	key = mlr_paste_4_strings(name1, "_", name2, "_ols_n");
	val = mlr_alloc_string_from_ll(pstate->count);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static stats2_t* stats2_linreg_ols_alloc(int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_linreg_ols_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_linreg_ols_state_t));
	pstate->count = 0LL;
	pstate->sumx  = 0.0;
	pstate->sumy  = 0.0;
	pstate->sumx2 = 0.0;
	pstate->sumxy = 0.0;

	pstats2->pvstate = (void*)pstate;
	pstats2->pingest_func = &stats2_linreg_ols_ingest;
	pstats2->pemit_func = &stats2_linreg_ols_emit;
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
	char* suffix = "r2";
	char* key = mlr_paste_5_strings(name1, "_", name2, "_", suffix);
	if (pstate->count < 2LL) {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
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
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}
static stats2_t* stats2_r2_alloc(int do_verbose) {
	stats2_t* pstats2 = mlr_malloc_or_die(sizeof(stats2_t));
	stats2_r2_state_t* pstate = mlr_malloc_or_die(sizeof(stats2_r2_state_t));
	pstate->count     = 0LL;
	pstate->sumx      = 0.0;
	pstate->sumy      = 0.0;
	pstate->sumx2     = 0.0;
	pstate->sumxy     = 0.0;
	pstate->sumy2     = 0.0;

	pstats2->pvstate      = (void*)pstate;
	pstats2->pingest_func = &stats2_r2_ingest;
	pstats2->pemit_func   = &stats2_r2_emit;

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
		char* key00 = mlr_paste_4_strings(name1, "_", name1, "_covx");
		char* key01 = mlr_paste_4_strings(name1, "_", name2, "_covx");
		char* key10 = mlr_paste_4_strings(name2, "_", name1, "_covx");
		char* key11 = mlr_paste_4_strings(name2, "_", name2, "_covx");
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
		char* keym   = mlr_paste_4_strings(name1, "_", name2, "_pca_m");
		char* keyb   = mlr_paste_4_strings(name1, "_", name2, "_pca_b");
		char* keyn   = mlr_paste_4_strings(name1, "_", name2, "_pca_n");
		char* keyq   = mlr_paste_4_strings(name1, "_", name2, "_pca_quality");
		char* keyl1  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival1");
		char* keyl2  = mlr_paste_4_strings(name1, "_", name2, "_pca_eival2");
		char* keyv11 = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec11");
		char* keyv12 = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec12");
		char* keyv21 = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec21");
		char* keyv22 = mlr_paste_4_strings(name1, "_", name2, "_pca_eivec22");
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
		char* suffix = (pstate->do_which == DO_CORR) ? "corr" : "cov";
		char* key = mlr_paste_5_strings(name1, "_", name2, "_", suffix);
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
static stats2_t* stats2_corr_cov_alloc(int do_which, int do_verbose) {
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

	pstats2->pvstate      = (void*)pstate;
	pstats2->pingest_func = &stats2_corr_cov_ingest;
	pstats2->pemit_func   = &stats2_corr_cov_emit;

	return pstats2;
}
static stats2_t* stats2_corr_alloc(int do_verbose) {
	return stats2_corr_cov_alloc(DO_CORR, do_verbose);
}
static stats2_t* stats2_cov_alloc(int do_verbose) {
	return stats2_corr_cov_alloc(DO_COV, do_verbose);
}
static stats2_t* stats2_covx_alloc(int do_verbose) {
	return stats2_corr_cov_alloc(DO_COVX, do_verbose);
}
static stats2_t* stats2_linreg_pca_alloc(int do_verbose) {
	return stats2_corr_cov_alloc(DO_LINREG_PCA, do_verbose);
}

// ----------------------------------------------------------------
typedef struct _stats2_lookup_t {
	char* name;
	stats2_alloc_func_t* pnew_func;
} stats2_lookup_t;
static stats2_lookup_t stats2_lookup_table[] = {
	{"linreg-pca", stats2_linreg_pca_alloc},
	{"linreg-ols", stats2_linreg_ols_alloc},
	{"r2",         stats2_r2_alloc},
	{"corr",       stats2_corr_alloc},
	{"cov",        stats2_cov_alloc},
	{"covx",       stats2_covx_alloc},
};
static int stats2_lookup_table_length = sizeof(stats2_lookup_table) / sizeof(stats2_lookup_table[0]);

static stats2_t* make_stats2(char* stats2_name, int do_verbose) {
	for (int i = 0; i < stats2_lookup_table_length; i++)
		if (streq(stats2_name, stats2_lookup_table[i].name))
			return stats2_lookup_table[i].pnew_func(do_verbose);
	return NULL;
}

// ================================================================
typedef struct _mapper_stats2_state_t {
	slls_t* paccumulator_names;
	slls_t* pvalue_field_name_pairs;
	slls_t* pgroup_by_field_names;

	lhmslv_t* groups;
	int     do_verbose;
} mapper_stats2_state_t;

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
// ================================================================

// ----------------------------------------------------------------
static void mapper_stats2_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats2_state_t* pstate);
static sllv_t* mapper_stats2_emit(mapper_stats2_state_t* pstate);

static sllv_t* mapper_stats2_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats2_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_stats2_ingest(pinrec, pctx, pstate);
		lrec_free(pinrec);
		return NULL;
	}
	else {
		return mapper_stats2_emit(pstate);
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

	lhms2v_t* group_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (group_to_acc_field == NULL) {
		group_to_acc_field = lhms2v_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), group_to_acc_field);
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
				pstats2 = make_stats2(stats2_name, pstate->do_verbose);
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
	}

	slls_free(pgroup_by_field_values);
}

// ----------------------------------------------------------------
static sllv_t* mapper_stats2_emit(mapper_stats2_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
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

// ----------------------------------------------------------------
static void mapper_stats2_free(void* pvstate) {
	mapper_stats2_state_t* pstate = pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_name_pairs);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->groups);
}

static mapper_t* mapper_stats2_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_name_pairs,
	slls_t* pgroup_by_field_names, int do_verbose)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats2_state_t* pstate   = mlr_malloc_or_die(sizeof(mapper_stats2_state_t));
	pstate->paccumulator_names      = paccumulator_names;
	pstate->pvalue_field_name_pairs = pvalue_field_name_pairs; // xxx validate length is even
	pstate->pgroup_by_field_names   = pgroup_by_field_names;
	pstate->groups                  = lhmslv_alloc();
	pstate->do_verbose              = do_verbose;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_stats2_process;
	pmapper->pfree_func    = mapper_stats2_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_stats2_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-a {linreg-ols,corr,...}  Names of accumulators: one or more of\n");
	fprintf(o, "             ");
	for (int i = 0; i < stats2_lookup_table_length; i++) {
		fprintf(o, " %s", stats2_lookup_table[i].name);
	}
	fprintf(o, "\n");
	fprintf(o, "              r2 is a quality metric for linreg-ols; linrec-pca outputs its\n");
	fprintf(o, "              own quality metric.\n");
	fprintf(o, "-f {a,b,c,d}  Value-field name-pairs on which to compute statistics.\n");
	fprintf(o, "              There must be an even number of names.\n");
	fprintf(o, "-g {e,f,g}    Optional group-by-field names.\n");
	fprintf(o, "-v            Print additional output for linreg-pca.\n");
	fprintf(o, "Example: %s %s -a linreg-pca -f x,y\n", argv0, verb);
	fprintf(o, "Example: %s %s -a linreg-ols,r2 -f x,y -g size,shape\n", argv0, verb);
	fprintf(o, "Example: %s %s -a corr -f x,y\n", argv0, verb);
}

static mapper_t* mapper_stats2_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* paccumulator_names    = NULL;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();
	int     do_verbose = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-a", &paccumulator_names);
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate,        "-v", &do_verbose);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
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

	return mapper_stats2_alloc(paccumulator_names, pvalue_field_names, pgroup_by_field_names, do_verbose);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_stats2_setup = {
	.verb = "stats2",
	.pusage_func = mapper_stats2_usage,
	.pparse_func = mapper_stats2_parse_cli
};
