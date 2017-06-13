#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/mlrstat.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmss.h"
#include "containers/lhmsll.h"
#include "containers/percentile_keeper.h"
#include "containers/mvfuncs.h"
#include "mapping/stats1_accumulators.h"

// ----------------------------------------------------------------
void make_stats1_accs(
	char*    value_field_name,            // input
	slls_t*  paccumulator_names,          // input
	int      allow_int_float,             // input
	int      do_interpolated_percentiles, // input
	lhmsv_t* acc_field_to_acc_state_in,   // output
	lhmsv_t* acc_field_to_acc_state_out)  // output
{
	stats1_acc_t* ppercentile_acc = NULL;
	for (sllse_t* pc = paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
		// for "sum", "count"
		char* stats1_acc_name = pc->value;

		// For percentiles there is one unique accumulator given (for example) five distinct
		// names p0,p25,p50,p75,p100.  The input accumulators are unique: only one
		// percentile-keeper. There are multiple output accumulators: each references the same
		// underlying percentile-keeper but with distinct parameters.  Hence the "_in" and "_out" maps.
		if (is_percentile_acc_name(stats1_acc_name)) {
			if (ppercentile_acc == NULL) {
				ppercentile_acc = stats1_percentile_alloc(value_field_name, stats1_acc_name, allow_int_float,
					do_interpolated_percentiles);
				if (ppercentile_acc == NULL) {
					fprintf(stderr, "%s stats1: accumulator \"%s\" not found.\n",
						MLR_GLOBALS.bargv0, stats1_acc_name);
					exit(1);
				}
				lhmsv_put(acc_field_to_acc_state_in, stats1_acc_name, ppercentile_acc, NO_FREE);
			} else {
				stats1_percentile_reuse(ppercentile_acc);
			}
			lhmsv_put(acc_field_to_acc_state_out, stats1_acc_name, ppercentile_acc, NO_FREE);
		} else {
			stats1_acc_t* pstats1_acc = make_stats1_acc(value_field_name, stats1_acc_name, allow_int_float,
				do_interpolated_percentiles);
			if (pstats1_acc == NULL) {
				fprintf(stderr, "%s stats1: accumulator \"%s\" not found.\n",
					MLR_GLOBALS.bargv0, stats1_acc_name);
				exit(1);
			}
			lhmsv_put(acc_field_to_acc_state_in, stats1_acc_name, pstats1_acc, NO_FREE);
			lhmsv_put(acc_field_to_acc_state_out, stats1_acc_name, pstats1_acc, NO_FREE);
		}
	}
}

stats1_acc_t* make_stats1_acc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	for (int i = 0; i < stats1_acc_lookup_table_length; i++)
		if (streq(stats1_acc_name, stats1_acc_lookup_table[i].name))
			return stats1_acc_lookup_table[i].palloc_func(value_field_name, stats1_acc_name, allow_int_float,
				do_interpolated_percentiles);
	return NULL;
}

int is_percentile_acc_name(char* stats1_acc_name) {
	if (streq(stats1_acc_name, "median"))
		return TRUE;
	double percentile;
	// sscanf(stats1_acc_name, "p%lf", &percentile) allows "p74x" et al. which isn't ok.
	if (stats1_acc_name[0] != 'p')
		return FALSE;
	if (!mlr_try_float_from_string(&stats1_acc_name[1], &percentile))
		return FALSE;
	if (percentile < 0.0 || percentile > 100.0) {
		fprintf(stderr, "%s stats1: percentile \"%s\" outside range [0,100].\n",
			MLR_GLOBALS.bargv0, stats1_acc_name);
		exit(1);
	}
	return TRUE;
}

// ----------------------------------------------------------------
typedef struct _stats1_count_state_t {
	mv_t counter;
	mv_t one;
	char* output_field_name;
} stats1_count_state_t;

static void stats1_count_singest(void* pvstate, char* val) {
	stats1_count_state_t* pstate = pvstate;
	pstate->counter = x_xx_plus_func(&pstate->counter, &pstate->one);

}
static void stats1_count_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_count_state_t* pstate = pvstate;
	if (copy_data)
		lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mv_alloc_format_val(&pstate->counter),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	else
		lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->counter),
			FREE_ENTRY_VALUE);
}
static void stats1_count_free(stats1_acc_t* pstats1_acc) {
	stats1_count_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_count_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc    = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_count_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_count_state_t));
	pstate->counter              = allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->one                  = allow_int_float ? mv_from_int(1LL) : mv_from_float(1.0);
	pstate->output_field_name    = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate         = (void*)pstate;
	pstats1_acc->pdingest_func   = NULL;
	pstats1_acc->pningest_func   = NULL;
	pstats1_acc->psingest_func   = stats1_count_singest;
	pstats1_acc->pemit_func      = stats1_count_emit;
	pstats1_acc->pfree_func      = stats1_count_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_mode_state_t {
	lhmsll_t* pcounts_for_value;
	char* output_field_name;
} stats1_mode_state_t;
// mode on strings: "1" and "1.0" and "1.0000" are distinct text.
static void stats1_mode_singest(void* pvstate, char* val) {
	stats1_mode_state_t* pstate = pvstate;
	lhmslle_t* pe = lhmsll_get_entry(pstate->pcounts_for_value, val);
	if (pe == NULL) {
		// lhmsll does a strdup so we needn't.
		lhmsll_put(pstate->pcounts_for_value, mlr_strdup_or_die(val), 1, FREE_ENTRY_KEY);
	} else {
		pe->value++;
	}
}
static void stats1_mode_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_mode_state_t* pstate = pvstate;
	int max_count = 0;
	char* max_key = "";
	for (lhmslle_t* pe = pstate->pcounts_for_value->phead; pe != NULL; pe = pe->pnext) {
		int count = pe->value;
		if (count > max_count) {
			max_key = pe->key;
			max_count = count;
		}
	}
	if (copy_data)
		lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mlr_strdup_or_die(max_key),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	else
		lrec_put(poutrec, pstate->output_field_name, max_key, NO_FREE);
}
static void stats1_mode_free(stats1_acc_t* pstats1_acc) {
	stats1_mode_state_t* pstate = pstats1_acc->pvstate;
	lhmsll_free(pstate->pcounts_for_value);
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_mode_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc   = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_mode_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_mode_state_t));
	pstate->pcounts_for_value   = lhmsll_alloc();
	pstate->output_field_name   = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate        = (void*)pstate;
	pstats1_acc->pdingest_func  = NULL;
	pstats1_acc->pningest_func  = NULL;
	pstats1_acc->psingest_func  = stats1_mode_singest;
	pstats1_acc->pemit_func     = stats1_mode_emit;
	pstats1_acc->pfree_func     = stats1_mode_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_antimode_state_t {
	lhmsll_t* pcounts_for_value;
	char* output_field_name;
} stats1_antimode_state_t;
// antimode on strings: "1" and "1.0" and "1.0000" are distinct text.
static void stats1_antimode_singest(void* pvstate, char* val) {
	stats1_antimode_state_t* pstate = pvstate;
	lhmslle_t* pe = lhmsll_get_entry(pstate->pcounts_for_value, val);
	if (pe == NULL) {
		// lhmsll does a strdup so we needn't.
		lhmsll_put(pstate->pcounts_for_value, mlr_strdup_or_die(val), 1, FREE_ENTRY_KEY);
	} else {
		pe->value++;
	}
}
static void stats1_antimode_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_antimode_state_t* pstate = pvstate;
	int min_count = 0;
	int have_min_count = FALSE;
	char* min_key = "";
	for (lhmslle_t* pe = pstate->pcounts_for_value->phead; pe != NULL; pe = pe->pnext) {
		int count = pe->value;
		if (!have_min_count || count < min_count) {
			min_key = pe->key;
			min_count = count;
			have_min_count = TRUE;
		}
	}
	if (copy_data)
		lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mlr_strdup_or_die(min_key),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	else
		lrec_put(poutrec, pstate->output_field_name, min_key, NO_FREE);
}
static void stats1_antimode_free(stats1_acc_t* pstats1_acc) {
	stats1_antimode_state_t* pstate = pstats1_acc->pvstate;
	lhmsll_free(pstate->pcounts_for_value);
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_antimode_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc   = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_antimode_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_antimode_state_t));
	pstate->pcounts_for_value   = lhmsll_alloc();
	pstate->output_field_name   = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate        = (void*)pstate;
	pstats1_acc->pdingest_func  = NULL;
	pstats1_acc->pningest_func  = NULL;
	pstats1_acc->psingest_func  = stats1_antimode_singest;
	pstats1_acc->pemit_func     = stats1_antimode_emit;
	pstats1_acc->pfree_func     = stats1_antimode_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_sum_state_t {
	mv_t sum;
	char* output_field_name;
	int allow_int_float;
} stats1_sum_state_t;
static void stats1_sum_ningest(void* pvstate, mv_t* pval) {
	stats1_sum_state_t* pstate = pvstate;
	pstate->sum = x_xx_plus_func(&pstate->sum, pval);
}
static void stats1_sum_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_sum_state_t* pstate = pvstate;
	if (copy_data)
		lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mv_alloc_format_val(&pstate->sum),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	else
		lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->sum),
			FREE_ENTRY_VALUE);
}
static void stats1_sum_free(stats1_acc_t* pstats1_acc) {
	stats1_sum_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_sum_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc  = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_sum_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_sum_state_t));
	pstate->allow_int_float    = allow_int_float;
	pstate->sum                = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = NULL;
	pstats1_acc->pningest_func = stats1_sum_ningest;
	pstats1_acc->psingest_func = NULL;
	pstats1_acc->pemit_func    = stats1_sum_emit;
	pstats1_acc->pfree_func    = stats1_sum_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_mean_state_t {
	double sum;
	unsigned long long count;
	char* output_field_name;
} stats1_mean_state_t;
static void stats1_mean_dingest(void* pvstate, double val) {
	stats1_mean_state_t* pstate = pvstate;
	pstate->sum   += val;
	pstate->count++;
}
static void stats1_mean_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_mean_state_t* pstate = pvstate;
	if (pstate->count == 0LL) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		double quot = pstate->sum / pstate->count;
		char* val = mlr_alloc_string_from_double(quot, MLR_GLOBALS.ofmt);
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
	}
}
static void stats1_mean_free(stats1_acc_t* pstats1_acc) {
	stats1_mean_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_mean_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc   = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_mean_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_mean_state_t));
	pstate->sum                 = 0.0;
	pstate->count               = 0LL;
	pstate->output_field_name   = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate        = (void*)pstate;
	pstats1_acc->pdingest_func  = stats1_mean_dingest;
	pstats1_acc->pningest_func  = NULL;
	pstats1_acc->psingest_func  = NULL;
	pstats1_acc->pemit_func     = stats1_mean_emit;
	pstats1_acc->pfree_func     = stats1_mean_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_stddev_var_meaneb_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	cumulant2o_t  do_which;
	char* output_field_name;
} stats1_stddev_var_meaneb_state_t;
static void stats1_stddev_var_meaneb_dingest(void* pvstate, double val) {
	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
}

static void stats1_stddev_var_meaneb_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		double output = mlr_get_var(pstate->count, pstate->sumx, pstate->sumx2);
		if (pstate->do_which == DO_STDDEV)
			output = sqrt(output);
		else if (pstate->do_which == DO_MEANEB)
			output = sqrt(output / pstate->count);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
	}
}
static void stats1_stddev_var_meaneb_free(stats1_acc_t* pstats1_acc) {
	stats1_stddev_var_meaneb_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}

stats1_acc_t* stats1_stddev_var_meaneb_alloc(char* value_field_name, char* stats1_acc_name, cumulant2o_t do_which) {
	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_stddev_var_meaneb_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_stddev_var_meaneb_state_t));
	pstate->count              = 0LL;
	pstate->sumx               = 0.0;
	pstate->sumx2              = 0.0;
	pstate->do_which           = do_which;
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = stats1_stddev_var_meaneb_dingest;
	pstats1_acc->pningest_func = NULL;
	pstats1_acc->psingest_func = NULL;
	pstats1_acc->pemit_func    = stats1_stddev_var_meaneb_emit;
	pstats1_acc->pfree_func    = stats1_stddev_var_meaneb_free;
	return pstats1_acc;
}
stats1_acc_t* stats1_stddev_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_STDDEV);
}
stats1_acc_t* stats1_var_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_VAR);
}
stats1_acc_t* stats1_meaneb_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_acc_name, DO_MEANEB);
}

// ----------------------------------------------------------------
typedef struct _stats1_skewness_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	double sumx3;
	char* output_field_name;
} stats1_skewness_state_t;
static void stats1_skewness_dingest(void* pvstate, double val) {
	stats1_skewness_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
	pstate->sumx3 += val*val*val;
}

static void stats1_skewness_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_skewness_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		double output = mlr_get_skewness(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
	}
}
static void stats1_skewness_free(stats1_acc_t* pstats1_acc) {
	stats1_skewness_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}

stats1_acc_t* stats1_skewness_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_skewness_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_skewness_state_t));
	pstate->count              = 0LL;
	pstate->sumx               = 0.0;
	pstate->sumx2              = 0.0;
	pstate->sumx3              = 0.0;
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = stats1_skewness_dingest;
	pstats1_acc->pningest_func = NULL;
	pstats1_acc->psingest_func = NULL;
	pstats1_acc->pemit_func    = stats1_skewness_emit;
	pstats1_acc->pfree_func    = stats1_skewness_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_kurtosis_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	double sumx3;
	double sumx4;
	char* output_field_name;
} stats1_kurtosis_state_t;
static void stats1_kurtosis_dingest(void* pvstate, double val) {
	stats1_kurtosis_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
	pstate->sumx3 += val*val*val;
	pstate->sumx4 += val*val*val*val;
}

static void stats1_kurtosis_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_kurtosis_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		double output = mlr_get_kurtosis(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3, pstate->sumx4);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), val, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
	}
}
static void stats1_kurtosis_free(stats1_acc_t* pstats1_acc) {
	stats1_kurtosis_state_t* pstate = pstats1_acc->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_kurtosis_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_kurtosis_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_kurtosis_state_t));
	pstate->count              = 0LL;
	pstate->sumx               = 0.0;
	pstate->sumx2              = 0.0;
	pstate->sumx3              = 0.0;
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = stats1_kurtosis_dingest;
	pstats1_acc->pningest_func = NULL;
	pstats1_acc->psingest_func = NULL;
	pstats1_acc->pemit_func    = stats1_kurtosis_emit;
	pstats1_acc->pfree_func    = stats1_kurtosis_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_min_state_t {
	mv_t min;
	char* output_field_name;
} stats1_min_state_t;
static void stats1_min_singest(void* pvstate, char* sval) {
	stats1_min_state_t* pstate = pvstate;
	mv_t val = mv_copy_type_infer_string_or_float_or_int(sval);
	pstate->min = x_xx_min_func(&pstate->min, &val);
}
static void stats1_min_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_min_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->min)) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mv_alloc_format_val(&pstate->min),
				FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->min),
				FREE_ENTRY_VALUE);
	}
}
static void stats1_min_free(stats1_acc_t* pstats1_acc) {
	stats1_min_state_t* pstate = pstats1_acc->pvstate;
	mv_free(&pstate->min);
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_min_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc  = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_min_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_min_state_t));
	pstate->min                = mv_absent();
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = NULL;
	pstats1_acc->pningest_func = NULL;
	pstats1_acc->psingest_func = stats1_min_singest;
	pstats1_acc->pemit_func    = stats1_min_emit;
	pstats1_acc->pfree_func    = stats1_min_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_max_state_t {
	mv_t max;
	char* output_field_name;
} stats1_max_state_t;
static void stats1_max_singest(void* pvstate, char* sval) {
	stats1_max_state_t* pstate = pvstate;
	mv_t val = mv_copy_type_infer_string_or_float_or_int(sval);
	pstate->max = x_xx_max_func(&pstate->max, &val);
}
static void stats1_max_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_max_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->max)) {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), "", FREE_ENTRY_KEY);
		else
			lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		if (copy_data)
			lrec_put(poutrec, mlr_strdup_or_die(pstate->output_field_name), mv_alloc_format_val(&pstate->max),
				FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		else
			lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->max),
				FREE_ENTRY_VALUE);
	}
}
static void stats1_max_free(stats1_acc_t* pstats1_acc) {
	stats1_max_state_t* pstate = pstats1_acc->pvstate;
	mv_free(&pstate->max);
	free(pstate->output_field_name);
	free(pstate);
	free(pstats1_acc);
}
stats1_acc_t* stats1_max_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc  = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_max_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_max_state_t));
	pstate->max                = mv_absent();
	pstate->output_field_name  = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);

	pstats1_acc->pvstate       = (void*)pstate;
	pstats1_acc->pdingest_func = NULL;
	pstats1_acc->pningest_func = NULL;
	pstats1_acc->psingest_func = stats1_max_singest;
	pstats1_acc->pemit_func    = stats1_max_emit;
	pstats1_acc->pfree_func    = stats1_max_free;
	return pstats1_acc;
}

// ----------------------------------------------------------------
typedef struct _stats1_percentile_state_t {
	percentile_keeper_t* ppercentile_keeper;
	lhmss_t* poutput_field_names;
	int reference_count;
	percentile_keeper_emitter_t* ppercentile_keeper_emitter;
} stats1_percentile_state_t;
static void stats1_percentile_singest(void* pvstate, char* sval) {
	stats1_percentile_state_t* pstate = pvstate;
	mv_t val = mv_copy_type_infer_string_or_float_or_int(sval);
	percentile_keeper_ingest(pstate->ppercentile_keeper, val);
}

static void stats1_percentile_emit(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec) {
	stats1_percentile_state_t* pstate = pvstate;
	double p;

	if (stats1_acc_name[0] == 'm') { // Pre-validated to be either p{number} or median.
		p = 50.0;
	} else {
		// TODO: do the sscanf once at alloc time and store the double in the state struct for a minor perf gain.
		(void)sscanf(stats1_acc_name, "p%lf", &p); // Assuming this was range-checked earlier on to be in [0,100].
	}
	mv_t v = pstate->ppercentile_keeper_emitter(pstate->ppercentile_keeper, p);
	char* s = mv_alloc_format_val(&v);
	// For this type, one accumulator tracks many stats1_names, but a single value_field_name.
	char* output_field_name = lhmss_get(pstate->poutput_field_names, stats1_acc_name);
	if (output_field_name == NULL) {
		output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_acc_name);
		lhmss_put(pstate->poutput_field_names, mlr_strdup_or_die(stats1_acc_name),
			output_field_name, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}
	lrec_put(poutrec, mlr_strdup_or_die(output_field_name), s, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
}

static void stats1_percentile_free(stats1_acc_t* pstats1_acc) {
	stats1_percentile_state_t* pstate = pstats1_acc->pvstate;
	pstate->reference_count--;
	if (pstate->reference_count == 0) {
		percentile_keeper_free(pstate->ppercentile_keeper);
		lhmss_free(pstate->poutput_field_names);
		free(pstate);
		free(pstats1_acc);
	}
}
stats1_acc_t* stats1_percentile_alloc(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles)
{
	stats1_acc_t* pstats1_acc   = mlr_malloc_or_die(sizeof(stats1_acc_t));
	stats1_percentile_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_percentile_state_t));
	pstate->ppercentile_keeper  = percentile_keeper_alloc();
	pstate->poutput_field_names = lhmss_alloc();
	pstate->reference_count     = 1;
	pstate->ppercentile_keeper_emitter = (do_interpolated_percentiles)
		? percentile_keeper_emit_linearly_interpolated
		: percentile_keeper_emit_non_interpolated;

	pstats1_acc->pvstate        = (void*)pstate;
	pstats1_acc->pdingest_func  = NULL;
	pstats1_acc->pningest_func  = NULL;
	pstats1_acc->psingest_func  = stats1_percentile_singest;
	pstats1_acc->pemit_func     = stats1_percentile_emit;
	pstats1_acc->pfree_func     = stats1_percentile_free;
	return pstats1_acc;
}
void stats1_percentile_reuse(stats1_acc_t* pstats1_acc) {
	stats1_percentile_state_t* pstate = pstats1_acc->pvstate;
	pstate->reference_count++;
}
