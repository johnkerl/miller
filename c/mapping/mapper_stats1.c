#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/mlrstat.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmsi.h"
#include "containers/mixutil.h"
#include "containers/percentile_keeper.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
typedef void acc_dingest_func_t(void* pvstate, double val);
typedef void acc_singest_func_t(void* pvstate, char*  val);
typedef void acc_emit_func_t(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec);

typedef struct _acc_t {
	void* pvstate;
	acc_singest_func_t* psingest_func;
	acc_dingest_func_t* pdingest_func;
	acc_emit_func_t*    pemit_func;
} acc_t;

typedef acc_t* acc_alloc_func_t();

// ----------------------------------------------------------------
typedef struct _acc_count_state_t {
	unsigned long long count;
} acc_count_state_t;
static void acc_count_singest(void* pvstate, char* val) {
	acc_count_state_t* pstate = pvstate;
	pstate->count++;
}
static void acc_count_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_count_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	char* val = mlr_alloc_string_from_ull(pstate->count);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static acc_t* acc_count_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_count_state_t* pstate = mlr_malloc_or_die(sizeof(acc_count_state_t));
	pstate->count       = 0LL;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = &acc_count_singest;
	pacc->pdingest_func = NULL;
	pacc->pemit_func    = &acc_count_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_mode_state_t {
	lhmsi_t* pcounts_for_value;
} acc_mode_state_t;
// mode on strings? what about "1.0" and "1" and "1.0000" ??
static void acc_mode_singest(void* pvstate, char* val) {
	acc_mode_state_t* pstate = pvstate;
	lhmsie_t* pe = lhmsi_get_entry(pstate->pcounts_for_value, val);
	if (pe == NULL) {
		// xxx at the moment, lhmsi does a strdup so we needn't.
		lhmsi_put(pstate->pcounts_for_value, val, 1);
	} else {
		pe->value++;
	}
}
static void acc_mode_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_mode_state_t* pstate = pvstate;
	int max_count = 0;
	char* max_key = "";
	for (lhmsie_t* pe = pstate->pcounts_for_value->phead; pe != NULL; pe = pe->pnext) {
		int count = pe->value;
		if (count > max_count) {
			max_key = pe->key;
			max_count = count;
		}
	}
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	lrec_put(poutrec, key, max_key, LREC_FREE_ENTRY_KEY);
}
static acc_t* acc_mode_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_mode_state_t* pstate = mlr_malloc_or_die(sizeof(acc_mode_state_t));
	pstate->pcounts_for_value = lhmsi_alloc();
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = &acc_mode_singest;
	pacc->pdingest_func = NULL;
	pacc->pemit_func    = &acc_mode_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_sum_state_t {
	double sum;
} acc_sum_state_t;
static void acc_sum_dingest(void* pvstate, double val) {
	acc_sum_state_t* pstate = pvstate;
	pstate->sum += val;
}
static void acc_sum_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_sum_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	char* val = mlr_alloc_string_from_double(pstate->sum, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static acc_t* acc_sum_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_sum_state_t* pstate = mlr_malloc_or_die(sizeof(acc_sum_state_t));
	pstate->sum         = 0LL;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_sum_dingest;
	pacc->pemit_func    = &acc_sum_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_mean_state_t {
	double sum;
	unsigned long long count;
} acc_mean_state_t;
static void acc_mean_dingest(void* pvstate, double val) {
	acc_mean_state_t* pstate = pvstate;
	pstate->sum   += val;
	pstate->count++;
}
static void acc_mean_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_mean_state_t* pstate = pvstate;
	double quot = pstate->sum / pstate->count;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	char* val = mlr_alloc_string_from_double(quot, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static acc_t* acc_mean_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_mean_state_t* pstate = mlr_malloc_or_die(sizeof(acc_mean_state_t));
	pstate->sum         = 0.0;
	pstate->count       = 0LL;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_mean_dingest;
	pacc->pemit_func    = &acc_mean_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_stddev_meaneb_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	int    do_meaneb;
} acc_stddev_meaneb_state_t;
static void acc_stddev_meaneb_dingest(void* pvstate, double val) {
	acc_stddev_meaneb_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
}
// xxx recast all of these in terms of providable outputs
static void acc_stddev_meaneb_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_stddev_meaneb_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	if (pstate->count < 2LL) {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	} else {
		double output = mlr_get_stddev(pstate->count, pstate->sumx, pstate->sumx2);
		if (pstate->do_meaneb)
			output = output / sqrt(pstate->count);
		char* val =  mlr_alloc_string_from_double(output, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

static acc_t* acc_stddev_meaneb_alloc(int do_meaneb) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_stddev_meaneb_state_t* pstate = mlr_malloc_or_die(sizeof(acc_stddev_meaneb_state_t));
	pstate->count       = 0LL;
	pstate->sumx        = 0.0;
	pstate->sumx2       = 0.0;
	pstate->do_meaneb   = do_meaneb;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_stddev_meaneb_dingest;
	pacc->pemit_func    = &acc_stddev_meaneb_emit;
	return pacc;
}
static acc_t* acc_stddev_alloc() {
	return acc_stddev_meaneb_alloc(FALSE);
}
static acc_t* acc_meaneb_alloc() {
	return acc_stddev_meaneb_alloc(TRUE);
}

// ----------------------------------------------------------------
typedef struct _acc_min_state_t {
	int have_min;
	double min;
} acc_min_state_t;
static void acc_min_dingest(void* pvstate, double val) {
	acc_min_state_t* pstate = pvstate;
	if (pstate->have_min) {
		if (val < pstate->min)
			pstate->min = val;
	} else {
		pstate->have_min = TRUE;
		pstate->min = val;
	}
}
static void acc_min_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_min_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	if (pstate->have_min) {
		char* val = mlr_alloc_string_from_double(pstate->min, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	} else {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	}
}
static acc_t* acc_min_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_min_state_t* pstate = mlr_malloc_or_die(sizeof(acc_min_state_t));
	pstate->have_min    = FALSE;
	pstate->min         = -999.0;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_min_dingest;
	pacc->pemit_func    = &acc_min_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_max_state_t {
	int have_max;
	double max;
} acc_max_state_t;
static void acc_max_dingest(void* pvstate, double val) {
	acc_max_state_t* pstate = pvstate;
	if (pstate->have_max) {
		if (val > pstate->max)
			pstate->max = val;
	} else {
		pstate->have_max = TRUE;
		pstate->max = val;
	}
}
static void acc_max_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_max_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
	if (pstate->have_max) {
		char* val = mlr_alloc_string_from_double(pstate->max, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	} else {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	}
}
static acc_t* acc_max_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_max_state_t* pstate = mlr_malloc_or_die(sizeof(acc_max_state_t));
	pstate->have_max    = FALSE;
	pstate->max         = -999.0;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_max_dingest;
	pacc->pemit_func    = &acc_max_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_percentile_state_t {
	percentile_keeper_t* ppercentile_keeper;
} acc_percentile_state_t;
static void acc_percentile_dingest(void* pvstate, double val) {
	acc_percentile_state_t* pstate = pvstate;
	percentile_keeper_ingest(pstate->ppercentile_keeper, val);
}
static void acc_percentile_emit(void* pvstate, char* value_field_name, char* acc_name, lrec_t* poutrec) {
	acc_percentile_state_t* pstate = pvstate;
	char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);

	double p;
	(void)sscanf(acc_name, "p%lf", &p); // Assuming this was range-checked earlier on to be in [0,100].
	double v = percentile_keeper_emit(pstate->ppercentile_keeper, p);
	char* s = mlr_alloc_string_from_double(v, MLR_GLOBALS.ofmt);
	lrec_put(poutrec, key, s, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
static acc_t* acc_percentile_alloc() {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_percentile_state_t* pstate = mlr_malloc_or_die(sizeof(acc_percentile_state_t));
	pstate->ppercentile_keeper = percentile_keeper_alloc();
	pacc->pvstate        = (void*)pstate;
	pacc->psingest_func  = NULL;
	pacc->pdingest_func  = &acc_percentile_dingest;
	pacc->pemit_func     = &acc_percentile_emit;
	return pacc;
}

// ----------------------------------------------------------------
// Lookups for all but percentiles, which are a special case.
typedef struct _acc_lookup_t {
	char* name;
	acc_alloc_func_t* pnew_func;
} acc_lookup_t;
static acc_lookup_t acc_lookup_table[] = {
	{"count",  acc_count_alloc},
	{"mode",   acc_mode_alloc},
	{"sum",    acc_sum_alloc},
	{"mean",   acc_mean_alloc},
	{"stddev", acc_stddev_alloc},
	{"meaneb", acc_meaneb_alloc},
	{"min",    acc_min_alloc},
	{"max",    acc_max_alloc},
};
static int acc_lookup_table_length = sizeof(acc_lookup_table) / sizeof(acc_lookup_table[0]);

static acc_t* make_acc(char* acc_name) {
	for (int i = 0; i < acc_lookup_table_length; i++)
		if (streq(acc_name, acc_lookup_table[i].name))
			return acc_lookup_table[i].pnew_func();
	return NULL;
}

// ================================================================
typedef struct _mapper_stats1_state_t {
	slls_t* paccumulator_names;
	slls_t* pvalue_field_names;
	slls_t* pgroup_by_field_names;

	lhmslv_t* groups;
} mapper_stats1_state_t;

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
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//   },
//   ["u","v"] : {
//     ["x"] : {
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//   },
//   ["u","w"] : {
//     ["x"] : {
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : stats2_count_t object,
//       "sum"   : stats2_sum_t  object
//     },
//   },
// }
// ================================================================

// ----------------------------------------------------------------
static void mapper_stats1_ingest(lrec_t* pinrec, mapper_stats1_state_t* pstate);
static sllv_t* mapper_stats1_emit(mapper_stats1_state_t* pstate);
static void make_accs(
	slls_t*    paccumulator_names,      // Input
	lhmsv_t*   acc_field_to_acc_state); // Output
char* fake_acc_name_for_setups = "__setup_done__";

static sllv_t* mapper_stats1_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_stats1_ingest(pinrec, pstate);
		lrec_free(pinrec);
		return NULL;
	} else {
		return mapper_stats1_emit(pstate);
	}
}

// ----------------------------------------------------------------
static void mapper_stats1_ingest(lrec_t* pinrec, mapper_stats1_state_t* pstate) {
	// E.g. ["s", "t"]
	// To do: make value_field_values into a hashmap. Then accept partial
	// population on that, but retain full-population requirement on group-by.
	// E.g. if accumulating stats of x,y on a,b then skip record with x,y,a but
	// process record with x,a,b.
	slls_t* pvalue_field_values    = mlr_selected_values_from_record(pinrec, pstate->pvalue_field_names);
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);
	if (pvalue_field_values->length != pstate->pvalue_field_names->length)
		return;
	if (pgroup_by_field_values->length != pstate->pgroup_by_field_names->length)
		return;

	lhmsv_t* group_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (group_to_acc_field == NULL) {
		group_to_acc_field = lhmsv_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), group_to_acc_field);
	}

	sllse_t* pa = pstate->pvalue_field_names->phead;
	sllse_t* pb =         pvalue_field_values->phead;
	// for x=1 and y=2
	for ( ; pa != NULL && pb != NULL; pa = pa->pnext, pb = pb->pnext) {
		char* value_field_name = pa->value;
		char* value_field_sval = pb->value;
		int   have_dval = FALSE;
		double value_field_dval = -999.0;

		lhmsv_t* acc_field_to_acc_state = lhmsv_get(group_to_acc_field, value_field_name);
		if (acc_field_to_acc_state == NULL) {
			acc_field_to_acc_state = lhmsv_alloc();
			lhmsv_put(group_to_acc_field, value_field_name, acc_field_to_acc_state);
		}

		// Look up presence of all accumulators at this level's hashmap.
		char* presence = lhmsv_get(acc_field_to_acc_state, fake_acc_name_for_setups);
		if (presence == NULL) {
			make_accs(pstate->paccumulator_names, acc_field_to_acc_state);
			lhmsv_put(acc_field_to_acc_state, fake_acc_name_for_setups, fake_acc_name_for_setups);
		}

		// There isn't a one-to-one mapping between user-specified acc_names
		// and internal acc_t's. Here in the ingestor we feed each datum into
		// an acc_t.  In the emitter, we loop over the acc_names in
		// user-specified order. Example: they ask for p10,mean,p90. Then there
		// is only one percentiles accumulator to be told about each point. In
		// the emitter it will be asked to produce output twice: once for the
		// 10th percentile & once for the 90th.
		for (lhmsve_t* pc = acc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
			char* acc_name = pc->key;
			if (streq(acc_name, fake_acc_name_for_setups))
				continue;
			acc_t* pacc = pc->pvvalue;

			if (pacc->psingest_func != NULL) {
				pacc->psingest_func(pacc->pvstate, value_field_sval);
			}
			if (pacc->pdingest_func != NULL) {
				if (!have_dval) {
					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
					have_dval = TRUE;
				}
				pacc->pdingest_func(pacc->pvstate, value_field_dval);
			}
		}
	}
}

// ----------------------------------------------------------------
static int is_percentile_acc_name(char* acc_name) {
	double percentile;
	if (sscanf(acc_name, "p%lf", &percentile) != 1)
		return FALSE;
	if (percentile < 0.0 || percentile > 100.0) {
		fprintf(stderr, "%s stats1: percentile \"%s\" outside range [0,100].\n",
			MLR_GLOBALS.argv0, acc_name);
		exit(1);
	}
	return TRUE;
}

// ----------------------------------------------------------------
static void make_accs(
	slls_t*    paccumulator_names,     // Input
	lhmsv_t*   acc_field_to_acc_state) // Output
{
	acc_t* ppercentile_acc = NULL;
	for (sllse_t* pc = paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
		// for "sum", "count"
		char* acc_name = pc->value;

		if (is_percentile_acc_name(acc_name)) {
			if (ppercentile_acc == NULL) {
				ppercentile_acc = acc_percentile_alloc();
			}
			lhmsv_put(acc_field_to_acc_state, acc_name, ppercentile_acc);
		} else {
			acc_t* pacc = make_acc(acc_name);
			if (pacc == NULL) {
				fprintf(stderr, "%s stats1: accumulator \"%s\" not found.\n",
					MLR_GLOBALS.argv0, acc_name);
				exit(1);
			}
			lhmsv_put(acc_field_to_acc_state, acc_name, pacc);
		}
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_stats1_emit(mapper_stats1_state_t* pstate) {
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
				char* acc_name = pe->value;
				if (streq(acc_name, fake_acc_name_for_setups))
					continue;
				acc_t* pacc = lhmsv_get(acc_field_to_acc_state, acc_name);
				if (pacc == NULL) {
					fprintf(stderr, "%s stats1: internal coding error: acc_name \"%s\" has gone missing.\n",
						MLR_GLOBALS.argv0, acc_name);
					exit(1);
				}
				pacc->pemit_func(pacc->pvstate, value_field_name, acc_name, poutrec);
			}
		}
		sllv_add(poutrecs, poutrec);
	}
	sllv_add(poutrecs, NULL);
	return poutrecs;
}

// ----------------------------------------------------------------
static void mapper_stats1_free(void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_names);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->groups);
}

static mapper_t* mapper_stats1_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats1_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_stats1_state_t));

	pstate->paccumulator_names    = paccumulator_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->groups                = lhmslv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_stats1_process;
	pmapper->pmapper_free_func    = mapper_stats1_free;

	return pmapper;
}

// ----------------------------------------------------------------
// xxx argify the stdout/stderr in ALL usages
static void mapper_stats1_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-a {sum,count,...}    Names of accumulators: p10 p25.2 p50 p98 p100 etc. and/or one or more of\n");
	fprintf(stdout, "                     ");
	for (int i = 0; i < acc_lookup_table_length; i++) {
		fprintf(stdout, " %s", acc_lookup_table[i].name);
	}
	fprintf(stdout, "\n");
	fprintf(stdout, "Options:\n");
	fprintf(stdout, "-f {a,b,c}            Value-field names on which to compute statistics\n");
	fprintf(stdout, "-g {d,e,f}            Group-by-field names\n");
	fprintf(stdout, "Example: %s %s -f value\n", argv0, verb);
	fprintf(stdout, "Example: %s %s -f value -g size,shape\n", argv0, verb);
	fprintf(stdout, "Notes:\n");
	fprintf(stdout, "* p50 is a synonym for median.\n");
	fprintf(stdout, "* min and max output the same results as p0 and p100, respectively, but use less memory.\n");
	fprintf(stdout, "* count and mode allow text input; the rest require numeric input. In particular, 1 and 1.0\n");
	fprintf(stdout, "  are distinct text for count and mode.\n");
	fprintf(stdout, "* When there are mode ties, the first-encountered datum wins.\n");
}

static mapper_t* mapper_stats1_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* paccumulator_names    = NULL;
	slls_t* pvalue_field_names    = NULL;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-a", &paccumulator_names);
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_stats1_usage(argv[0], verb);
		return NULL;
	}

	if (paccumulator_names == NULL || pvalue_field_names == NULL) {
		mapper_stats1_usage(argv[0], verb);
		return NULL;
	}

	return mapper_stats1_alloc(paccumulator_names, pvalue_field_names, pgroup_by_field_names);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_stats1_setup = {
	.verb        = "stats1",
	.pusage_func = mapper_stats1_usage,
	.pparse_func = mapper_stats1_parse_cli
};
