#ifndef STATS1_ACCUMULATORS_H
#define STATS1_ACCUMULATORS_H

#include "containers/lrec.h"
#include "containers/slls.h"
#include "containers/lhmsv.h"

// ----------------------------------------------------------------
// These are used by mlr stats1 as well as mlr merge-fields.

// ----------------------------------------------------------------
// Types of second-order cumulant
typedef enum _cumulant2o_t {
	DO_STDDEV,
	DO_VAR,
	DO_MEANEB,
} cumulant2o_t;

// ----------------------------------------------------------------
struct _stats1_acc_t; // forward reference for method definitions
typedef void stats1_dingest_func_t(void* pvstate, double val);
typedef void stats1_ningest_func_t(void* pvstate, mv_t* pval);
typedef void stats1_singest_func_t(void* pvstate, char*  val);
// For mlr stats1, there is a single alloc at start including computation of
// output field name; arbitrary number of records output with that field name;
// then free. There it makes sense to point the poutrec key at the output field
// name in the acc state rather than copying it. For mlr merge-fields, by
// contrast, the outer/inner loops are reversed: arbitrary number of records;
// for each, allocate, set output field key/value, free; then, the record is
// output. There it is necessary to copy the key since it will be referenced
// after the accumulator is freed.
typedef void stats1_emit_func_t(void* pvstate, char* value_field_name, char* stats1_acc_name, int copy_data, lrec_t* poutrec);
typedef void stats1_free_func_t(struct _stats1_acc_t* pstats1_acc);

typedef struct _stats1_acc_t {
	void* pvstate;
	stats1_dingest_func_t* pdingest_func;
	stats1_ningest_func_t* pningest_func;
	stats1_singest_func_t* psingest_func;
	stats1_emit_func_t*    pemit_func;
	stats1_free_func_t*    pfree_func; // virtual destructor
} stats1_acc_t;

typedef stats1_acc_t* stats1_alloc_func_t(char* value_field_name, char* stats1_acc_name, int allow_int_float,
	int do_interpolated_percentiles);

// aif = allow_int_float
// dip = do_interpolated_percentiles
stats1_acc_t* stats1_count_alloc             (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_mode_alloc              (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_antimode_alloc          (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_sum_alloc               (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_mean_alloc              (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_stddev_var_meaneb_alloc (char* value_field_name, char* stats1_acc_name, cumulant2o_t do_which);
stats1_acc_t* stats1_stddev_alloc            (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_var_alloc               (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_meaneb_alloc            (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_skewness_alloc          (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_kurtosis_alloc          (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_min_alloc               (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_max_alloc               (char* value_field_name, char* stats1_acc_name, int aif, int dip);
stats1_acc_t* stats1_percentile_alloc        (char* value_field_name, char* stats1_acc_name, int aif, int dip);
void          stats1_percentile_reuse        (stats1_acc_t* pstats1_acc);


// For percentiles there is one unique accumulator given (for example) five distinct
// names p0,p25,p50,p75,p100.  The input accumulators are unique: only one
// percentile-keeper. There are multiple output accumulators: each references the same
// underlying percentile-keeper but with distinct parameters.  Hence the "_in" and "_out" maps.
void make_stats1_accs(
	char*    value_field_name,
	slls_t*  paccumulator_names,
	int      allow_int_float,
	int      do_interpolated_percentiles,
	lhmsv_t* acc_field_to_acc_state_in,
	lhmsv_t* acc_field_to_acc_state_out);

stats1_acc_t* make_stats1_acc(
	char* value_field_name,
	char* stats1_acc_name,
	int   allow_int_float,
	int   do_interpolated_percentiles);

int is_percentile_acc_name(char* stats1_acc_name);

// ----------------------------------------------------------------
// Lookups for all but percentiles, which are a special case.
typedef struct _stats1_acc_lookup_t {
	char* name;
	stats1_alloc_func_t* palloc_func;
	char* desc;
} stats1_acc_lookup_t;
static stats1_acc_lookup_t stats1_acc_lookup_table[] = {
	{"count",    stats1_count_alloc,    "Count instances of fields"},
	{"mode",     stats1_mode_alloc,     "Find most-frequently-occurring values for fields; first-found wins tie"},
	{"antimode", stats1_antimode_alloc, "Find least-frequently-occurring values for fields; first-found wins tie"},
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
static int stats1_acc_lookup_table_length = sizeof(stats1_acc_lookup_table) / sizeof(stats1_acc_lookup_table[0]);

#endif // STATS1_ACCUMULATORS_H
