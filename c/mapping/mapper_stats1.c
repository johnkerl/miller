#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/mlrstat.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/string_array.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmsi.h"
#include "containers/mixutil.h"
#include "containers/percentile_keeper.h"
#include "mapping/mappers.h"
#include "mapping/mlr_val.h"
#include "cli/argparse.h"

#define DO_STDDEV 0xc1
#define DO_VAR    0xc2
#define DO_MEANEB 0xc3

// ================================================================
typedef void stats1_dingest_func_t(void* pvstate, double val);
typedef void stats1_ningest_func_t(void* pvstate, mv_t* pval);
typedef void stats1_singest_func_t(void* pvstate, char*  val);
typedef void stats1_emit_func_t(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec);

typedef struct _stats1_t {
	void* pvstate;
	stats1_dingest_func_t* pdingest_func;
	stats1_ningest_func_t* pningest_func;
	stats1_singest_func_t* psingest_func;
	stats1_emit_func_t*    pemit_func;
} stats1_t;

typedef stats1_t* stats1_alloc_func_t(char* value_field_name, char* stats1_name);

typedef struct _mapper_stats1_state_t {
	slls_t*         paccumulator_names;
	string_array_t* pvalue_field_names;     // parameter
	string_array_t* pvalue_field_values;    // scratch space used per-record
	slls_t*         pgroup_by_field_names;  // parameter
	lhmslv_t*       groups;
	int             do_iterative_stats;
} mapper_stats1_state_t;

// ----------------------------------------------------------------
static void      mapper_stats1_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_stats1_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_stats1_alloc(slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int do_iterative_stats);
static void      mapper_stats1_free(void* pvstate);
static sllv_t*   mapper_stats1_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_stats1_ingest(lrec_t* pinrec, mapper_stats1_state_t* pstate);
static sllv_t*   mapper_stats1_emit_all(mapper_stats1_state_t* pstate);
static lrec_t*   mapper_stats1_emit(mapper_stats1_state_t* pstate, lrec_t* poutrec,
	char* value_field_name, char* stats1_name, lhmsv_t* acc_field_to_acc_state);

static stats1_t* stats1_count_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_mode_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_sum_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_mean_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_stddev_var_meaneb_alloc(char* value_field_name, char* stats1_name, int do_which);
static stats1_t* stats1_stddev_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_var_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_meaneb_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_skewness_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_kurtosis_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_min_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_max_alloc(char* value_field_name, char* stats1_name);
static stats1_t* stats1_percentile_alloc(char* value_field_name, char* stats1_name);

static stats1_t* make_acc(char* value_field_name, char* stats1_name);
static void make_accs(char* value_field_name, slls_t* paccumulator_names, lhmsv_t* acc_field_to_acc_state);

// ----------------------------------------------------------------
// Lookups for all but percentiles, which are a special case.
typedef struct _acc_lookup_t {
	char* name;
	stats1_alloc_func_t* pnew_func;
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

// ----------------------------------------------------------------
mapper_setup_t mapper_stats1_setup = {
	.verb        = "stats1",
	.pusage_func = mapper_stats1_usage,
	.pparse_func = mapper_stats1_parse_cli
};

// ----------------------------------------------------------------
static void mapper_stats1_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Options:\n");
	fprintf(o, "-a {sum,count,...}  Names of accumulators: p10 p25.2 p50 p98 p100 etc. and/or\n");
	fprintf(o, "                    one or more of:\n");
	for (int i = 0; i < stats1_lookup_table_length; i++) {
		fprintf(o, "  %-7s %s\n", stats1_lookup_table[i].name, stats1_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics\n");
	fprintf(o, "-g {d,e,f}  Optional group-by-field names\n");
	fprintf(o, "-s          Print iterative stats. Useful in tail -f contexts (in which\n");
	fprintf(o, "            case please avoid pprint-format output since end of input\n");
	fprintf(o, "            stream will never be seen).\n");
	fprintf(o, "Example: %s %s -a min,p10,p50,p90,max -f value -g size,shape\n", argv0, verb);
	fprintf(o, "Example: %s %s -a count,mode -f size\n", argv0, verb);
	fprintf(o, "Example: %s %s -a count,mode -f size -g shape\n", argv0, verb);
	fprintf(o, "Notes:\n");
	fprintf(o, "* p50 is a synonym for median.\n");
	fprintf(o, "* min and max output the same results as p0 and p100, respectively, but use\n");
	fprintf(o, "  less memory.\n");
	fprintf(o, "* count and mode allow text input; the rest require numeric input.\n");
	fprintf(o, "  In particular, 1 and 1.0 are distinct text for count and mode.\n");
	fprintf(o, "* When there are mode ties, the first-encountered datum wins.\n");
}

static mapper_t* mapper_stats1_parse_cli(int* pargi, int argc, char** argv) {
	slls_t*         paccumulator_names    = NULL;
	string_array_t* pvalue_field_names    = NULL;
	slls_t*         pgroup_by_field_names = slls_alloc();
	int             do_iterative_stats    = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate,  "-a", &paccumulator_names);
	ap_define_string_array_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate,  "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate,         "-s", &do_iterative_stats);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_stats1_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (paccumulator_names == NULL || pvalue_field_names == NULL) {
		mapper_stats1_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_stats1_alloc(paccumulator_names, pvalue_field_names, pgroup_by_field_names,
		do_iterative_stats);
}

// ----------------------------------------------------------------
static mapper_t* mapper_stats1_alloc(slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int do_iterative_stats)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats1_state_t* pstate  = mlr_malloc_or_die(sizeof(mapper_stats1_state_t));

	pstate->paccumulator_names     = paccumulator_names;
	pstate->pvalue_field_names     = pvalue_field_names;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->pvalue_field_values    = string_array_alloc(pvalue_field_names->length);
	pstate->groups                 = lhmslv_alloc();
	pstate->do_iterative_stats     = do_iterative_stats;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_stats1_process;
	pmapper->pfree_func    = mapper_stats1_free;

	return pmapper;
}

static void mapper_stats1_free(void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	slls_free(pstate->paccumulator_names);
	string_array_free(pstate->pvalue_field_names);
	string_array_free(pstate->pvalue_field_values);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->groups);
}

// ================================================================
// Given: accumulate count,sum on values x,y group by a,b.
// Example input:       Example output:
//   a b x y            a b x_count x_sum y_count y_sum
//   s t 1 2            s t 2       6     2       8
//   u v 3 4            u v 1       3     1       4
//   s t 5 6            u w 1       7     1       9
//   u w 7 9
//
// Multilevel hashmap structure:
// {
//   ["s","t"] : {                <--- group-by field names
//     ["x"] : {                  <--- value field names
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//   },
//   ["u","v"] : {
//     ["x"] : {
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//   },
//   ["u","w"] : {
//     ["x"] : {
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats1_count_t object,
//       "sum"   : stats1_sum_t  object
//     },
//   },
// }
// ================================================================

char* fake_acc_name_for_setups = "__setup_done__";

// In the iterative case, add to the current record its current group's stats fields.
// In the non-iterative case, produce output only at the end of the input stream.
static sllv_t* mapper_stats1_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_stats1_ingest(pinrec, pstate);
		if (pstate->do_iterative_stats) {
			// The input record is modified in this case, with new fields appended
			return sllv_single(pinrec);
		} else {
			lrec_free(pinrec);
			return NULL;
		}
	} else if (!pstate->do_iterative_stats) {
		return mapper_stats1_emit_all(pstate);
	} else {
		return NULL;
	}
}

static stats1_t* make_acc(char* value_field_name, char* stats1_name) {
	for (int i = 0; i < stats1_lookup_table_length; i++)
		if (streq(stats1_name, stats1_lookup_table[i].name))
			return stats1_lookup_table[i].pnew_func(value_field_name, stats1_name);
	return NULL;
}

// ----------------------------------------------------------------
static void mapper_stats1_ingest(lrec_t* pinrec, mapper_stats1_state_t* pstate) {
	// E.g. ["s", "t"]
	// To do: make value_field_values into a hashmap. Then accept partial
	// population on that, but retain full-population requirement on group-by.
	// E.g. if accumulating stats of x,y on a,b then skip record with x,y,a but
	// process record with x,a,b.
	mlr_reference_values_from_record(pinrec, pstate->pvalue_field_names, pstate->pvalue_field_values);
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

	if (pgroup_by_field_values == NULL) {
		slls_free(pgroup_by_field_values);
		return;
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

		lhmsv_t* acc_field_to_acc_state = lhmsv_get(group_to_acc_field, value_field_name);
		if (acc_field_to_acc_state == NULL) {
			acc_field_to_acc_state = lhmsv_alloc();
			lhmsv_put(group_to_acc_field, value_field_name, acc_field_to_acc_state);
		}

		// Look up presence of all accumulators at this level's hashmap.
		char* presence = lhmsv_get(acc_field_to_acc_state, fake_acc_name_for_setups);
		if (presence == NULL) {
			make_accs(value_field_name, pstate->paccumulator_names, acc_field_to_acc_state);
			lhmsv_put(acc_field_to_acc_state, fake_acc_name_for_setups, fake_acc_name_for_setups);
		}

		if (value_field_sval == NULL)
			continue;

		int have_dval = FALSE;
		int have_nval = FALSE;
		double value_field_dval = -999.0;
		mv_t   value_field_nval = mv_from_float(-888.0);

		// There isn't a one-to-one mapping between user-specified stats1_names
		// and internal stats1_t's. Here in the ingestor we feed each datum into
		// an stats1_t.  In the emitter, we loop over the stats1_names in
		// user-specified order. Example: they ask for p10,mean,p90. Then there
		// is only one percentiles accumulator to be told about each point. In
		// the emitter it will be asked to produce output twice: once for the
		// 10th percentile & once for the 90th.
		for (lhmsve_t* pc = acc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
			char* stats1_name = pc->key;
			if (streq(stats1_name, fake_acc_name_for_setups))
				continue;
			stats1_t* pstats1 = pc->pvvalue;

			if (pstats1->pdingest_func != NULL) {
				if (!have_dval) {
					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
					have_dval = TRUE;
				}
				pstats1->pdingest_func(pstats1->pvstate, value_field_dval);
			}
			if (pstats1->pningest_func != NULL) {
				if (!have_nval) {
					value_field_nval = mt_scan_number_or_die(value_field_sval);
					have_nval = TRUE;
				}
				pstats1->pningest_func(pstats1->pvstate, &value_field_nval);
			}
			if (pstats1->psingest_func != NULL) {
				pstats1->psingest_func(pstats1->pvstate, value_field_sval);
			}

			if (pstate->do_iterative_stats) {
				mapper_stats1_emit(pstate, pinrec, value_field_name, stats1_name, acc_field_to_acc_state);
			}
		}
	}
	slls_free(pgroup_by_field_values);
}

// ----------------------------------------------------------------
static int is_percentile_acc_name(char* stats1_name) {
	double percentile;
	// sscanf(stats1_name, "p%lf", &percentile) allows "p74x" et al. which isn't ok.
	if (stats1_name[0] != 'p')
		return FALSE;
	if (!mlr_try_float_from_string(&stats1_name[1], &percentile))
		return FALSE;
	if (percentile < 0.0 || percentile > 100.0) {
		fprintf(stderr, "%s stats1: percentile \"%s\" outside range [0,100].\n",
			MLR_GLOBALS.argv0, stats1_name);
		exit(1);
	}
	return TRUE;
}

// ----------------------------------------------------------------
static void make_accs(
	char*      value_field_name,       // input
	slls_t*    paccumulator_names,     // input
	lhmsv_t*   acc_field_to_acc_state) // output
{
	stats1_t* ppercentile_acc = NULL;
	for (sllse_t* pc = paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
		// for "sum", "count"
		char* stats1_name = pc->value;

		if (is_percentile_acc_name(stats1_name)) {
			if (ppercentile_acc == NULL) {
				ppercentile_acc = stats1_percentile_alloc(value_field_name, stats1_name);
			}
			lhmsv_put(acc_field_to_acc_state, stats1_name, ppercentile_acc);
		} else {
			stats1_t* pstats1 = make_acc(value_field_name, stats1_name);
			if (pstats1 == NULL) {
				fprintf(stderr, "%s stats1: accumulator \"%s\" not found.\n",
					MLR_GLOBALS.argv0, stats1_name);
				exit(1);
			}
			lhmsv_put(acc_field_to_acc_state, stats1_name, pstats1);
		}
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_stats1_emit_all(mapper_stats1_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
		slls_t* pgroup_by_field_values = pa->key;
		lrec_t* poutrec = lrec_unbacked_alloc();

		// Add in a=s,b=t fields:
		sllse_t* pb = pstate->pgroup_by_field_names->phead;
		sllse_t* pc =         pgroup_by_field_values->phead;
		for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
			lrec_put(poutrec, pb->value, pc->value, 0);
		}

		// Add in fields such as x_sum=#, y_count=#, etc.:
		lhmsv_t* group_to_acc_field = pa->pvvalue;
		// for "x", "y"
		for (lhmsve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
			char* value_field_name = pd->key;
			lhmsv_t* acc_field_to_acc_state = pd->pvvalue;

			for (sllse_t* pe = pstate->paccumulator_names->phead; pe != NULL; pe = pe->pnext) {
				char* stats1_name = pe->value;
				mapper_stats1_emit(pstate, poutrec, value_field_name, stats1_name, acc_field_to_acc_state);
			}
		}
		sllv_add(poutrecs, poutrec);
	}
	sllv_add(poutrecs, NULL);
	return poutrecs;
}

// ----------------------------------------------------------------
static lrec_t* mapper_stats1_emit(mapper_stats1_state_t* pstate, lrec_t* poutrec,
	char* value_field_name, char* stats1_name, lhmsv_t* acc_field_to_acc_state)
{
	// Add in fields such as x_sum=#, y_count=#, etc.:
	for (sllse_t* pe = pstate->paccumulator_names->phead; pe != NULL; pe = pe->pnext) {
		char* stats1_name = pe->value;
		if (streq(stats1_name, fake_acc_name_for_setups))
			continue;
		stats1_t* pstats1 = lhmsv_get(acc_field_to_acc_state, stats1_name);
		if (pstats1 == NULL) {
			fprintf(stderr, "%s stats1: internal coding error: stats1_name \"%s\" has gone missing.\n",
				MLR_GLOBALS.argv0, stats1_name);
			exit(1);
		}
		pstats1->pemit_func(pstats1->pvstate, value_field_name, stats1_name, poutrec);
	}
	return poutrec;
}

// ----------------------------------------------------------------
typedef struct _stats1_count_state_t {
	unsigned long long count;
	char* output_field_name;
} stats1_count_state_t;

static void stats1_count_singest(void* pvstate, char* val) {
	stats1_count_state_t* pstate = pvstate;
	pstate->count++;
}
static void stats1_count_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_count_state_t* pstate = pvstate;
	char* val = mlr_alloc_string_from_ull(pstate->count);
	lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static stats1_t* stats1_count_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_count_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_count_state_t));
	pstate->count       = 0LL;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = NULL;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = stats1_count_singest;
	pstats1->pemit_func    = stats1_count_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_mode_state_t {
	lhmsi_t* pcounts_for_value;
	char* output_field_name;
} stats1_mode_state_t;
// mode on strings: "1" and "1.0" and "1.0000" are distinct text.
static void stats1_mode_singest(void* pvstate, char* val) {
	stats1_mode_state_t* pstate = pvstate;
	lhmsie_t* pe = lhmsi_get_entry(pstate->pcounts_for_value, val);
	if (pe == NULL) {
		// lhmsi does a strdup so we needn't.
		lhmsi_put(pstate->pcounts_for_value, val, 1);
	} else {
		pe->value++;
	}
}
static void stats1_mode_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_mode_state_t* pstate = pvstate;
	int max_count = 0;
	char* max_key = "";
	for (lhmsie_t* pe = pstate->pcounts_for_value->phead; pe != NULL; pe = pe->pnext) {
		int count = pe->value;
		if (count > max_count) {
			max_key = pe->key;
			max_count = count;
		}
	}
	lrec_put(poutrec, pstate->output_field_name, max_key, LREC_FREE_ENTRY_KEY);
}
static stats1_t* stats1_mode_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_mode_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_mode_state_t));
	pstate->pcounts_for_value = lhmsi_alloc();
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = NULL;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = stats1_mode_singest;
	pstats1->pemit_func    = stats1_mode_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_sum_state_t {
	double sum;
	char* output_field_name;
} stats1_sum_state_t;
static void stats1_sum_dingest(void* pvstate, double val) {
	stats1_sum_state_t* pstate = pvstate;
	pstate->sum += val;
}
static void stats1_sum_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_sum_state_t* pstate = pvstate;
	char* val = mlr_alloc_string_from_double(pstate->sum, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static stats1_t* stats1_sum_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_sum_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_sum_state_t));
	pstate->sum         = 0.0;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = stats1_sum_dingest;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_sum_emit;
	return pstats1;
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
static void stats1_mean_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_mean_state_t* pstate = pvstate;
	if (pstate->count == 0LL) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double quot = pstate->sum / pstate->count;
		char* val = mlr_alloc_string_from_double(quot, MLR_GLOBALS.ofmt);
		// xxx to do: the output field names should be freed by our free
		// method, & the flags shouldn't include LREC_FREE_ENTRY_KEY.
		lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}
static stats1_t* stats1_mean_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_mean_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_mean_state_t));
	pstate->sum         = 0.0;
	pstate->count       = 0LL;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = stats1_mean_dingest;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_mean_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_stddev_var_meaneb_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	int    do_which;
	char* output_field_name;
} stats1_stddev_var_meaneb_state_t;
static void stats1_stddev_var_meaneb_dingest(void* pvstate, double val) {
	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
}

static void stats1_stddev_var_meaneb_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_stddev_var_meaneb_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double output = mlr_get_var(pstate->count, pstate->sumx, pstate->sumx2);
		if (pstate->do_which == DO_STDDEV)
			output = sqrt(output);
		else if (pstate->do_which == DO_MEANEB)
			output = sqrt(output / pstate->count);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static stats1_t* stats1_stddev_var_meaneb_alloc(char* value_field_name, char* stats1_name, int do_which) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_stddev_var_meaneb_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_stddev_var_meaneb_state_t));
	pstate->count       = 0LL;
	pstate->sumx        = 0.0;
	pstate->sumx2       = 0.0;
	pstate->do_which    = do_which;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = stats1_stddev_var_meaneb_dingest;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_stddev_var_meaneb_emit;
	return pstats1;
}
static stats1_t* stats1_stddev_alloc(char* value_field_name, char* stats1_name) {
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_name, DO_STDDEV);
}
static stats1_t* stats1_var_alloc(char* value_field_name, char* stats1_name) {
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_name, DO_VAR);
}
static stats1_t* stats1_meaneb_alloc(char* value_field_name, char* stats1_name) {
	return stats1_stddev_var_meaneb_alloc(value_field_name, stats1_name, DO_MEANEB);
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

static void stats1_skewness_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_skewness_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double output = mlr_get_skewness(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static stats1_t* stats1_skewness_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_skewness_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_skewness_state_t));
	pstate->count       = 0LL;
	pstate->sumx        = 0.0;
	pstate->sumx2       = 0.0;
	pstate->sumx3       = 0.0;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = stats1_skewness_dingest;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_skewness_emit;
	return pstats1;
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

static void stats1_kurtosis_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_kurtosis_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		double output = mlr_get_kurtosis(pstate->count, pstate->sumx, pstate->sumx2, pstate->sumx3, pstate->sumx4);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->output_field_name, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static stats1_t* stats1_kurtosis_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_kurtosis_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_kurtosis_state_t));
	pstate->count       = 0LL;
	pstate->sumx        = 0.0;
	pstate->sumx2       = 0.0;
	pstate->sumx3       = 0.0;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);

	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = stats1_kurtosis_dingest;
	pstats1->pningest_func = NULL;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_kurtosis_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_min_state_t {
	mv_t min;
	char* output_field_name;
} stats1_min_state_t;
static void stats1_min_ningest(void* pvstate, mv_t* pval) {
	stats1_min_state_t* pstate = pvstate;
	pstate->min = n_nn_min_func(&pstate->min, pval);
}
static void stats1_min_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_min_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->min)) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		lrec_put(poutrec, pstate->output_field_name, mt_format_val(&pstate->min),
			LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}
static stats1_t* stats1_min_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_min_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_min_state_t));
	pstate->min = mv_from_null();
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);
	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = NULL;
	pstats1->pningest_func = stats1_min_ningest;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_min_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_max_state_t {
	mv_t max;
	char* output_field_name;
} stats1_max_state_t;
static void stats1_max_ningest(void* pvstate, mv_t* pval) {
	stats1_max_state_t* pstate = pvstate;
	pstate->max = n_nn_min_func(&pstate->max, pval);
}
static void stats1_max_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_max_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->max)) {
		lrec_put(poutrec, pstate->output_field_name, "", LREC_FREE_ENTRY_KEY);
	} else {
		lrec_put(poutrec, pstate->output_field_name, mt_format_val(&pstate->max),
			LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}
static stats1_t* stats1_max_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_max_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_max_state_t));
	pstate->max = mv_from_null();
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);
	pstats1->pvstate       = (void*)pstate;
	pstats1->pdingest_func = NULL;
	pstats1->pningest_func = stats1_max_ningest;
	pstats1->psingest_func = NULL;
	pstats1->pemit_func    = stats1_max_emit;
	return pstats1;
}

// ----------------------------------------------------------------
typedef struct _stats1_percentile_state_t {
	percentile_keeper_t* ppercentile_keeper;
} stats1_percentile_state_t;
static void stats1_percentile_dingest(void* pvstate, double val) {
	stats1_percentile_state_t* pstate = pvstate;
	percentile_keeper_ingest(pstate->ppercentile_keeper, val);
}
static void stats1_percentile_emit(void* pvstate, char* value_field_name, char* stats1_name, lrec_t* poutrec) {
	stats1_percentile_state_t* pstate = pvstate;

	double p;
	(void)sscanf(stats1_name, "p%lf", &p); // Assuming this was range-checked earlier on to be in [0,100].
	double v = percentile_keeper_emit(pstate->ppercentile_keeper, p);
	char* s = mlr_alloc_string_from_double(v, MLR_GLOBALS.ofmt);
	// For this type, one accumulator track many stats1_names.
	char* output_field_name = mlr_paste_3_strings(value_field_name, "_", stats1_name);
	lrec_put(poutrec, output_field_name, s, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static stats1_t* stats1_percentile_alloc(char* value_field_name, char* stats1_name) {
	stats1_t* pstats1 = mlr_malloc_or_die(sizeof(stats1_t));
	stats1_percentile_state_t* pstate = mlr_malloc_or_die(sizeof(stats1_percentile_state_t));
	pstate->ppercentile_keeper = percentile_keeper_alloc();

	pstats1->pvstate        = (void*)pstate;
	pstats1->pdingest_func  = stats1_percentile_dingest;
	pstats1->pningest_func  = NULL;
	pstats1->psingest_func  = NULL;
	pstats1->pemit_func     = stats1_percentile_emit;
	return pstats1;
}
