#ifndef STATS1_ACCUMULATORS_H
#define STATS1_ACCUMULATORS_H

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//#include <math.h>
//#include "lib/mlrutil.h"
//#include "lib/mlr_globals.h"
//#include "lib/mlrstat.h"
//#include "containers/sllv.h"
//#include "containers/slls.h"
//#include "containers/string_array.h"
//#include "containers/lhmslv.h"
//#include "containers/lhmsv.h"
//#include "containers/lhmsi.h"
//#include "containers/mixutil.h"
//#include "containers/percentile_keeper.h"
//#include "containers/mlrval.h"
//#include "mapping/mappers.h"
//#include "cli/argparse.h"

// ----------------------------------------------------------------
struct _stats1_t; // forward reference for method definitions
typedef void stats1_dingest_func_t(void* pvstate, double val);
typedef void stats1_ningest_func_t(void* pvstate, mv_t* pval);
typedef void stats1_singest_func_t(void* pvstate, char*  val);
typedef void stats1_emit_func_t(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec);
typedef void stats1_free_func_t(struct _stats1_t* pstats1);

typedef struct _stats1_t {
	void* pvstate;
	stats1_dingest_func_t* pdingest_func;
	stats1_ningest_func_t* pningest_func;
	stats1_singest_func_t* psingest_func;
	stats1_emit_func_t*    pemit_func;
	stats1_free_func_t*    pfree_func; // virtual destructor
} stats1_t;

typedef stats1_t* stats1_alloc_func_t(char* value_field_name, char* stats1_name, int allow_int_float);

stats1_t* stats1_count_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_mode_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_sum_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_mean_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_stddev_var_meaneb_alloc(char* value_field_name, char* stats1_name, int do_which);
stats1_t* stats1_stddev_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_var_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_meaneb_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_skewness_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_kurtosis_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_min_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_max_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
stats1_t* stats1_percentile_alloc(char* value_field_name, char* stats1_name, int allow_int_float);
void      stats1_percentile_reuse(stats1_t* pstats1);

void make_accs(char* value_field_name, slls_t* paccumulator_names, int allow_int_float, lhmsv_t* acc_field_to_acc_state);
stats1_t* make_acc(char* value_field_name, char* stats1_name, int allow_int_float);

// ----------------------------------------------------------------
// Lookups for all but percentiles, which are a special case.
typedef struct _acc_lookup_t {
	char* name;
	stats1_alloc_func_t* palloc_func;
	char* desc;
} stats1_lookup_t;
static stats1_lookup_t stats1_lookup_table[] = {
	{"count",    stats1_count_alloc,    "Count instances of fields"},
	{"mode",     stats1_mode_alloc,     "Find most-frequently-occurring values for fields; first-found wins tie"},
	{"sum",      stats1_sum_alloc,      "Compute sums of specified fields"},
	{"mean",     stats1_mean_alloc,     "Compute averages (sample means) of specified fields"},
	{"stddev",   stats1_stddev_alloc,   "Compute sample standard deviation of specified fields"},
	{"var",      stats1_var_alloc,      "Compute sample variance of specified fields"},
	{"meaneb",   stats1_meaneb_alloc,   "Estimate error bars for averages (assuming no sample autocorrelation)"},
	{"skewness", stats1_skewness_alloc, "Compute sample skewness of specified fields"},
	{"kurtosis", stats1_kurtosis_alloc, "Compute sample kurtosis of specified fields"},
	{"min",      stats1_min_alloc,      "Compute minimum values of specified fields"},
	{"max",      stats1_max_alloc,      "Compute maximum values of specified fields"},
};
static int stats1_lookup_table_length = sizeof(stats1_lookup_table) / sizeof(stats1_lookup_table[0]);

#endif // STATS1_ACCUMULATORS_H
