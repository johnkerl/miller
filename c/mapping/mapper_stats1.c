#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlrstat.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmsi.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
typedef void acc_dingest_func_t(void* pvstate, double val);
typedef void acc_singest_func_t(void* pvstate, char*  val);
typedef void acc_emit_func_t(void* pvstate, char* value_field_name, lrec_t* poutrec);

typedef struct _acc_t {
	void* pvstate;
	acc_singest_func_t* psingest_func;
	acc_dingest_func_t* pdingest_func;
	acc_emit_func_t*    pemit_func;
} acc_t;

typedef acc_t* acc_alloc_func_t(static_context_t* pstatx);

// ----------------------------------------------------------------
typedef struct _acc_count_state_t {
	unsigned long long count;
	static_context_t* pstatx;
} acc_count_state_t;
void acc_count_singest(void* pvstate, char* val) {
	acc_count_state_t* pstate = pvstate;
	pstate->count++;
}
void acc_count_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_count_state_t* pstate = pvstate;
	char* key = mlr_paste_2_strings(value_field_name, "_count");
	char* val = mlr_alloc_string_from_ull(pstate->count);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
acc_t* acc_count_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_count_state_t* pstate = mlr_malloc_or_die(sizeof(acc_count_state_t));
	pstate->count    = 0LL;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = &acc_count_singest;
	pacc->pdingest_func = NULL;
	pacc->pemit_func  = &acc_count_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_sum_state_t {
	double sum;
	static_context_t* pstatx;
} acc_sum_state_t;
void acc_sum_dingest(void* pvstate, double val) {
	acc_sum_state_t* pstate = pvstate;
	pstate->sum += val;
}
void acc_sum_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_sum_state_t* pstate = pvstate;
	char* key = mlr_paste_2_strings(value_field_name, "_sum");
	char* val = mlr_alloc_string_from_double(pstate->sum, pstate->pstatx->ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
acc_t* acc_sum_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_sum_state_t* pstate = mlr_malloc_or_die(sizeof(acc_sum_state_t));
	pstate->sum      = 0LL;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_sum_dingest;
	pacc->pemit_func  = &acc_sum_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_avg_state_t {
	double sum;
	unsigned long long count;
	static_context_t* pstatx;
} acc_avg_state_t;
void acc_avg_dingest(void* pvstate, double val) {
	acc_avg_state_t* pstate = pvstate;
	pstate->sum   += val;
	pstate->count++;
}
void acc_avg_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_avg_state_t* pstate = pvstate;
	double quot = pstate->sum / pstate->count;
	char* key = mlr_paste_2_strings(value_field_name, "_avg");
	char* val = mlr_alloc_string_from_double(quot, pstate->pstatx->ofmt);
	lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
}
acc_t* acc_avg_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_avg_state_t* pstate = mlr_malloc_or_die(sizeof(acc_avg_state_t));
	pstate->sum    = 0.0;
	pstate->count  = 0LL;
	pstate->pstatx = pstatx;
	pacc->pvstate  = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_avg_dingest;
	pacc->pemit_func  = &acc_avg_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_stddev_avgeb_state_t {
	unsigned long long count;
	double sumx;
	double sumx2;
	int    do_avgeb;
	static_context_t* pstatx;
} acc_stddev_avgeb_state_t;
void acc_stddev_avgeb_dingest(void* pvstate, double val) {
	acc_stddev_avgeb_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
}
void acc_stddev_avgeb_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_stddev_avgeb_state_t* pstate = pvstate;
	char* key = mlr_paste_2_strings(value_field_name, pstate->do_avgeb ? "_avgeb" : "_stddev");
	if (pstate->count < 2LL) {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	} else {
		double output = mlr_get_stddev(pstate->count, pstate->sumx, pstate->sumx2);
		if (pstate->do_avgeb)
			output = output / sqrt(pstate->count);
		char* val =  mlr_alloc_string_from_double(output, pstate->pstatx->ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	}
}

acc_t* acc_stddev_avgeb_alloc(int do_avgeb, static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_stddev_avgeb_state_t* pstate = mlr_malloc_or_die(sizeof(acc_stddev_avgeb_state_t));
	pstate->count       = 0LL;
	pstate->sumx        = 0.0;
	pstate->sumx2       = 0.0;
	pstate->do_avgeb    = do_avgeb;
	pstate->pstatx      = pstatx;
	pacc->pvstate       = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_stddev_avgeb_dingest;
	pacc->pemit_func    = &acc_stddev_avgeb_emit;
	return pacc;
}
acc_t* acc_stddev_alloc(static_context_t* pstatx) {
	return acc_stddev_avgeb_alloc(FALSE, pstatx);
}
acc_t* acc_avgeb_alloc(static_context_t* pstatx) {
	return acc_stddev_avgeb_alloc(TRUE, pstatx);
}

// ----------------------------------------------------------------
typedef struct _acc_min_state_t {
	int have_min;
	double min;
	static_context_t* pstatx;
} acc_min_state_t;
void acc_min_dingest(void* pvstate, double val) {
	acc_min_state_t* pstate = pvstate;
	if (pstate->have_min) {
		if (val < pstate->min)
			pstate->min = val;
	} else {
		pstate->have_min = TRUE;
		pstate->min = val;
	}
}
void acc_min_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_min_state_t* pstate = pvstate;
	char* key = mlr_paste_2_strings(value_field_name, "_min");
	if (pstate->have_min) {
		char* val = mlr_alloc_string_from_double(pstate->min, pstate->pstatx->ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	} else {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	}
}
acc_t* acc_min_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_min_state_t* pstate = mlr_malloc_or_die(sizeof(acc_min_state_t));
	pstate->have_min = FALSE;
	pstate->min      = -999.0;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_min_dingest;
	pacc->pemit_func  = &acc_min_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_max_state_t {
	int have_max;
	double max;
	static_context_t* pstatx;
} acc_max_state_t;
void acc_max_dingest(void* pvstate, double val) {
	acc_max_state_t* pstate = pvstate;
	if (pstate->have_max) {
		if (val > pstate->max)
			pstate->max = val;
	} else {
		pstate->have_max = TRUE;
		pstate->max = val;
	}
}
void acc_max_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_max_state_t* pstate = pvstate;
	char* key = mlr_paste_2_strings(value_field_name, "_max");
	if (pstate->have_max) {
		char* val = mlr_alloc_string_from_double(pstate->max, pstate->pstatx->ofmt);
		lrec_put(poutrec, key, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	} else {
		lrec_put(poutrec, key, "", LREC_FREE_ENTRY_KEY);
	}
}
acc_t* acc_max_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_max_state_t* pstate = mlr_malloc_or_die(sizeof(acc_max_state_t));
	pstate->have_max = FALSE;
	pstate->max      = -999.0;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_max_dingest;
	pacc->pemit_func  = &acc_max_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_mode_state_t {
	lhmsi_t* pcounts_for_value;
	static_context_t* pstatx;
} acc_mode_state_t;
// mode on strings? what about "1.0" and "1" and "1.0000" ??
void acc_mode_singest(void* pvstate, char* val) {
	acc_mode_state_t* pstate = pvstate;
	lhmsie_t* pe = lhmsi_get_entry(pstate->pcounts_for_value, val);
	if (pe == NULL) {
		// xxx at the moment, lhmsi does a strdup so we needn't.
		lhmsi_put(pstate->pcounts_for_value, val, 1);
	} else {
		pe->value++;
	}
}
void acc_mode_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
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
	char* key = mlr_paste_2_strings(value_field_name, "_mode");
	lrec_put(poutrec, key, max_key, LREC_FREE_ENTRY_KEY);
}
// xxx somewhere note a crucial assumption: while pctx is passed in
// to the mapper on a per-row basis, we stash it here on first use and
// use it on subsequent rows. assumptions:
// * the address doesn't change
// * the content we use (namely, ofmt) isn't row-dependent
// Option 1:
// * modify make_acc to special-case p{n}. needs multi-level hashmap keys
// * do it outside make_acc; requires separate hash maps for percentiles/deciles/quartiles/etc.
acc_t* acc_mode_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_mode_state_t* pstate = mlr_malloc_or_die(sizeof(acc_mode_state_t));
	pstate->pcounts_for_value = lhmsi_alloc();
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = &acc_mode_singest;
	pacc->pdingest_func = NULL;
	pacc->pemit_func  = &acc_mode_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_foo_state_t {
	int have_foo;
	double foo;
	static_context_t* pstatx;
} acc_foo_state_t;
void acc_foo_dingest(void* pvstate, double val) {
	acc_foo_state_t* pstate = pvstate;
	if (pstate->have_foo) {
		if (val > pstate->foo)
			pstate->foo = val;
	} else {
		pstate->have_foo = TRUE;
		pstate->foo = val;
	}
}
void acc_foo_emit(void* pvstate, char* value_field_name, lrec_t* poutrec) {
	acc_foo_state_t* pstate = pvstate;
	char* key1 = mlr_paste_2_strings(value_field_name, "_foo");
	char* key2 = mlr_paste_2_strings(value_field_name, "_bar");
	if (pstate->have_foo) {
		char* val = mlr_alloc_string_from_double(pstate->foo, pstate->pstatx->ofmt);
		lrec_put(poutrec, key1, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
		val = mlr_alloc_string_from_double(-pstate->foo, pstate->pstatx->ofmt);
		lrec_put(poutrec, key2, val, LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
	} else {
		lrec_put(poutrec, key1, "", LREC_FREE_ENTRY_KEY);
		lrec_put(poutrec, key2, "", LREC_FREE_ENTRY_KEY);
	}
}
acc_t* acc_foo_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_foo_state_t* pstate = mlr_malloc_or_die(sizeof(acc_foo_state_t));
	pstate->have_foo = FALSE;
	pstate->foo      = -999.0;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psingest_func = NULL;
	pacc->pdingest_func = &acc_foo_dingest;
	pacc->pemit_func  = &acc_foo_emit;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_lookup_t {
	char* name;
	acc_alloc_func_t* pnew_func;
} acc_lookup_t;
static acc_lookup_t acc_lookup_table[] = {
	{"count",  acc_count_alloc},
	{"sum",    acc_sum_alloc},
	{"avg",    acc_avg_alloc},
	{"stddev", acc_stddev_alloc},
	{"avgeb",  acc_avgeb_alloc},
	{"min",    acc_min_alloc},
	{"max",    acc_max_alloc},
	{"mode",   acc_mode_alloc},
	{"foo",    acc_foo_alloc},
};
static int acc_lookup_table_length = sizeof(acc_lookup_table) / sizeof(acc_lookup_table[0]);

// xxx make this a hashmap?
// xxx what if acc_name is p50? need:
// * here and here alone is cross-dependence between accumulators
// * if there are min,p10,p50,avg,p90,max then the values array should be
//   shared between p10,p50,p90
static acc_t* make_acc(char* acc_name, static_context_t* pstatx) {
	for (int i = 0; i < acc_lookup_table_length; i++)
		if (streq(acc_name, acc_lookup_table[i].name))
			return acc_lookup_table[i].pnew_func(pstatx);
	return NULL;
}

// ================================================================
typedef struct _mapper_stats1_state_t {
	slls_t* paccumulator_names;
	slls_t* pvalue_field_names;
	slls_t* pgroup_by_field_names;

	lhmslv_t* groups;
} mapper_stats1_state_t;

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
//   ["s","t"] : {
//     ["x"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//   },
//   ["u","v"] : {
//     ["x"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//   },
//   ["u","w"] : {
//     ["x"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//     ["y"] : {
//       "count" : C stats2_count_t object,
//       "sum"   : C stats2_sum_t  object
//     },
//   },
// }

// ----------------------------------------------------------------
static void mapper_stats1_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats1_state_t* pstate);
static sllv_t* mapper_stats1_emit(mapper_stats1_state_t* pstate);
static void make_accs(
	slls_t*    paccumulator_names,      // Input
	context_t* pctx,                    // Input
	lhmsv_t*   acc_field_to_acc_state); // Output
char* fake_acc_name_for_setups = "__setup_done__";

sllv_t* mapper_stats1_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_stats1_ingest(pinrec, pctx, pstate);
		lrec_free(pinrec);
		return NULL;
	} else {
		return mapper_stats1_emit(pstate);
	}
}

// ----------------------------------------------------------------
static void mapper_stats1_ingest(lrec_t* pinrec, context_t* pctx, mapper_stats1_state_t* pstate) {
	// ["s", "t"]
	// xxx make value_field_values into a hashmap. then accept partial population on that.
	// xxx but retain full-population requirement on group-by.
	// e.g. if accumulating stats of x,y on a,b then skip row with x,y,a but process row with x,a,b.
	slls_t* pvalue_field_values    = mlr_selected_values_from_record(pinrec, pstate->pvalue_field_names);
	slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

	// xxx cmt
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
	// for "x", "y" and "1", "2"
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

		// xxx cmt
		char* presence = lhmsv_get(acc_field_to_acc_state, fake_acc_name_for_setups);
		if (presence == NULL) {
			make_accs(pstate->paccumulator_names, pctx, acc_field_to_acc_state);
			lhmsv_put(acc_field_to_acc_state, fake_acc_name_for_setups, fake_acc_name_for_setups);
		}

		for (lhmsve_t* pc = acc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
			char* acc_name = pc->key;
			if (streq(acc_name, fake_acc_name_for_setups))
				continue;
			acc_t* pacc = pc->value;

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
	return (1 == sscanf(acc_name, "p%lf", &percentile));
}

// ----------------------------------------------------------------
static void make_accs(
	slls_t*    paccumulator_names,     // Input
	context_t* pctx,                   // Input
	lhmsv_t*   acc_field_to_acc_state) // Output
{
	for (sllse_t* pc = paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
		// for "sum", "count"
		char* acc_name = pc->value;
		slls_t* ppercentile_names = slls_alloc();

		if (is_percentile_acc_name(acc_name)) {
			// crap. this isn't cool. it moves the order of the pnn's within the user-provided arg. :/
			slls_add_no_free(ppercentile_names, acc_name);
		} else {
			acc_t* pacc = make_acc(acc_name, &pctx->statx);
			if (pacc == NULL) {
				// xxx needs argv[0] from mlr_globals.
				fprintf(stderr, "mlr stats1: accumulator \"%s\" not found.\n",
					acc_name);
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
		lrec_t* poutrec = lrec_unbacked_alloc();

		slls_t* pgroup_by_field_values = pa->key;

		// Add in a=s,b=t fields:
		sllse_t* pb = pstate->pgroup_by_field_names->phead;
		sllse_t* pc =         pgroup_by_field_values->phead;
		for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
			lrec_put(poutrec, pb->value, pc->value, 0);
		}

		// Add in fields such as x_sum=#, y_count=#, etc.:
		lhmsv_t* group_to_acc_field = pa->value;
		// for "x", "y"
		for (lhmsve_t* pd = group_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
			char* value_field_name = pd->key;
			lhmsv_t* acc_field_to_acc_state = pd->value;
			// for "count", "sum"
			for (lhmsve_t* pe = acc_field_to_acc_state->phead; pe != NULL; pe = pe->pnext) {
				if (streq(pe->key, fake_acc_name_for_setups))
					continue;
				acc_t* pacc = pe->value;
				pacc->pemit_func(pacc->pvstate, value_field_name, poutrec);
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

mapper_t* mapper_stats1_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats1_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_stats1_state_t));

	pstate->paccumulator_names    = paccumulator_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->groups                = lhmslv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_stats1_func;
	pmapper->pmapper_free_func    = mapper_stats1_free;

	return pmapper;
}

// ----------------------------------------------------------------
// xxx argify the stdout/stderr in ALL usages
void mapper_stats1_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-a {sum,count,...}    Names of accumulators: one or more of\n");
	fprintf(stdout, "                     ");
	for (int i = 0; i < acc_lookup_table_length; i++) {
		fprintf(stdout, " %s", acc_lookup_table[i].name);
	}
	fprintf(stdout, "\n");
	fprintf(stdout, "-f {a,b,c}            Value-field names on which to compute statistics\n");
	fprintf(stdout, "-g {d,e,f}            Group-by-field names\n");
}

mapper_t* mapper_stats1_parse_cli(int* pargi, int argc, char** argv) {
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
