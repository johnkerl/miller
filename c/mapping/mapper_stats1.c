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
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
// -a min,max -v x,y -g a,b

// ["s","t"] |--> "x" |--> "sum" |--> acc_sum_t* (as void*)
// level 1      level 2   level 3
// lhmslv_t     lhmsv_t   lhmsv_t
// acc_sum_t implements interface:
//   void  dacc(double dval);
//   void  sacc(char*  sval);
//   char* get();

// ================================================================
typedef void  acc_dput_func_t(void* pvstate, double val);
typedef void  acc_sput_func_t(void* pvstate, char*  val);
typedef char* acc_get_func_t (void* pvstate, char* pfree_flags);

typedef struct _acc_t {
	void* pvstate;
	acc_sput_func_t* psput_func;
	acc_dput_func_t* pdput_func;
	acc_get_func_t*  pget_func;
} acc_t;

typedef acc_t* acc_alloc_func_t(static_context_t* pstatx);

// ----------------------------------------------------------------
typedef struct _acc_count_state_t {
	unsigned long long count;
	static_context_t* pstatx;
} acc_count_state_t;
void acc_count_sput(void* pvstate, char* val) {
	acc_count_state_t* pstate = pvstate;
	pstate->count++;
}
char* acc_count_get(void* pvstate, char* pfree_flags) {
	acc_count_state_t* pstate = pvstate;
	*pfree_flags |= LREC_FREE_ENTRY_VALUE;
	return mlr_alloc_string_from_ull(pstate->count);
}
// xxx somewhere note a crucial assumption: while pctx is passed in
// to the mapper on a per-row basis, we stash it here on first use and
// use it on subsequent rows. assumptions:
// * the address doesn't change
// * the content we use (namely, ofmt) isn't row-dependent
acc_t* acc_count_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_count_state_t* pstate = mlr_malloc_or_die(sizeof(acc_count_state_t));
	pstate->count    = 0LL;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psput_func = &acc_count_sput;
	pacc->pdput_func = NULL;
	pacc->pget_func  = &acc_count_get;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_sum_state_t {
	double sum;
	static_context_t* pstatx;
} acc_sum_state_t;
void acc_sum_dput(void* pvstate, double val) {
	acc_sum_state_t* pstate = pvstate;
	pstate->sum += val;
}
char* acc_sum_get(void* pvstate, char* pfree_flags) {
	acc_sum_state_t* pstate = pvstate;
	*pfree_flags |= LREC_FREE_ENTRY_VALUE;
	return mlr_alloc_string_from_double(pstate->sum, pstate->pstatx->ofmt);
}
acc_t* acc_sum_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_sum_state_t* pstate = mlr_malloc_or_die(sizeof(acc_sum_state_t));
	pstate->sum      = 0LL;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psput_func = NULL;
	pacc->pdput_func = &acc_sum_dput;
	pacc->pget_func  = &acc_sum_get;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_avg_state_t {
	double sum;
	unsigned long long count;
	static_context_t* pstatx;
} acc_avg_state_t;
void acc_avg_dput(void* pvstate, double val) {
	acc_avg_state_t* pstate = pvstate;
	pstate->sum   += val;
	pstate->count++;
}
char* acc_avg_get(void* pvstate, char* pfree_flags) {
	acc_avg_state_t* pstate = pvstate;
	double quot = pstate->sum / pstate->count;
	// xxx decide: format "NaN" or "" when count is zeo. and, when *can* count be zero?
	*pfree_flags |= LREC_FREE_ENTRY_VALUE;
	return mlr_alloc_string_from_double(quot, pstate->pstatx->ofmt);
}
acc_t* acc_avg_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_avg_state_t* pstate = mlr_malloc_or_die(sizeof(acc_avg_state_t));
	pstate->sum    = 0.0;
	pstate->count  = 0LL;
	pstate->pstatx = pstatx;
	pacc->pvstate  = (void*)pstate;
	pacc->psput_func = NULL;
	pacc->pdput_func = &acc_avg_dput;
	pacc->pget_func  = &acc_avg_get;
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
void acc_stddev_avgeb_dput(void* pvstate, double val) {
	acc_stddev_avgeb_state_t* pstate = pvstate;
	pstate->count++;
	pstate->sumx  += val;
	pstate->sumx2 += val*val;
}
char* acc_stddev_avgeb_get(void* pvstate, char* pfree_flags) {
	acc_stddev_avgeb_state_t* pstate = pvstate;
	if (pstate->count < 2LL) {
		*pfree_flags &= ~LREC_FREE_ENTRY_VALUE;
		return "";
	} else {
		double output = mlr_get_stddev(pstate->count, pstate->sumx, pstate->sumx2);
		*pfree_flags |= LREC_FREE_ENTRY_VALUE;
		if (pstate->do_avgeb)
			output = output / sqrt(pstate->count);
		return mlr_alloc_string_from_double(output, pstate->pstatx->ofmt);
	}
}
acc_t* acc_stddev_avgeb_alloc(int do_avgeb, static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_stddev_avgeb_state_t* pstate = mlr_malloc_or_die(sizeof(acc_stddev_avgeb_state_t));
	pstate->count     = 0LL;
	pstate->sumx      = 0.0;
	pstate->sumx2     = 0.0;
	pstate->do_avgeb  = do_avgeb;
	pstate->pstatx    = pstatx;
	pacc->pvstate     = (void*)pstate;
	pacc->psput_func  = NULL;
	pacc->pdput_func  = &acc_stddev_avgeb_dput;
	pacc->pget_func   = &acc_stddev_avgeb_get;
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
void acc_min_dput(void* pvstate, double val) {
	acc_min_state_t* pstate = pvstate;
	if (pstate->have_min) {
		if (val < pstate->min)
			pstate->min = val;
	} else {
		pstate->have_min = TRUE;
		pstate->min = val;
	}
}
char* acc_min_get(void* pvstate, char* pfree_flags) {
	acc_min_state_t* pstate = pvstate;
	if (pstate->have_min) {
		*pfree_flags |= LREC_FREE_ENTRY_VALUE;
		return mlr_alloc_string_from_double(pstate->min, pstate->pstatx->ofmt);
	} else {
		*pfree_flags &= ~LREC_FREE_ENTRY_VALUE;
		return "";
	}
}
acc_t* acc_min_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_min_state_t* pstate = mlr_malloc_or_die(sizeof(acc_min_state_t));
	pstate->have_min = FALSE;
	pstate->min      = -999.0;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psput_func = NULL;
	pacc->pdput_func = &acc_min_dput;
	pacc->pget_func  = &acc_min_get;
	return pacc;
}

// ----------------------------------------------------------------
typedef struct _acc_max_state_t {
	int have_max;
	double max;
	static_context_t* pstatx;
} acc_max_state_t;
void acc_max_dput(void* pvstate, double val) {
	acc_max_state_t* pstate = pvstate;
	if (pstate->have_max) {
		if (val > pstate->max)
			pstate->max = val;
	} else {
		pstate->have_max = TRUE;
		pstate->max = val;
	}
}
char* acc_max_get(void* pvstate, char* pfree_flags) {
	acc_max_state_t* pstate = pvstate;
	if (pstate->have_max) {
		*pfree_flags &= ~LREC_FREE_ENTRY_VALUE;
		return mlr_alloc_string_from_double(pstate->max, pstate->pstatx->ofmt);
	} else {
		*pfree_flags |= LREC_FREE_ENTRY_VALUE;
		return "";
	}
}
acc_t* acc_max_alloc(static_context_t* pstatx) {
	acc_t* pacc = mlr_malloc_or_die(sizeof(acc_t));
	acc_max_state_t* pstate = mlr_malloc_or_die(sizeof(acc_max_state_t));
	pstate->have_max = FALSE;
	pstate->max      = -999.0;
	pstate->pstatx   = pstatx;
	pacc->pvstate    = (void*)pstate;
	pacc->psput_func = NULL;
	pacc->pdput_func = &acc_max_dput;
	pacc->pget_func  = &acc_max_get;
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
};
static int acc_lookup_table_length = sizeof(acc_lookup_table) / sizeof(acc_lookup_table[0]);

// xxx make this a hashmap?
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

	lhmslv_t* pmaps_level_1;

} mapper_stats1_state_t;

// given: accumulate count,sum on values x,y group by a,b
// example input:       example output:
//   a b x y            a b x_count x_sum y_count y_sum
//   s t 1 2            s t 2       6     2       8
//   u v 3 4            u v 1       3     1       4
//   s t 5 6            u w 1       7     1       9
//   u w 7 9

// ["s","t"] |--> "x" |--> "sum" |--> acc_sum_t* (as void*)
// level_1      level_2   level_3
// lhmslv_t     lhmsv_t   lhmsv_t
// acc_sum_t implements interface:
//   void  init();
//   void  dacc(double dval);
//   void  sacc(char*  sval);
//   char* get();

// ----------------------------------------------------------------
sllv_t* mapper_stats1_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		// ["s", "t"]
		// xxx make value_field_values into a hashmap. then accept partial population on that.
		// xxx but retain full-population requirement on group-by.
		// e.g. if accumulating stats of x,y on a,b then skip row with x,y,a but process row with x,a,b.
		slls_t* pvalue_field_values    = mlr_selected_values_from_record(pinrec, pstate->pvalue_field_names);
		slls_t* pgroup_by_field_values = mlr_selected_values_from_record(pinrec, pstate->pgroup_by_field_names);

		// xxx cmt
		if (pvalue_field_values->length != pstate->pvalue_field_names->length) {
			lrec_free(pinrec);
			return NULL;
		}
		if (pgroup_by_field_values->length != pstate->pgroup_by_field_names->length) {
			lrec_free(pinrec);
			return NULL;
		}

		lhmsv_t* pmaps_level_2 = lhmslv_get(pstate->pmaps_level_1, pgroup_by_field_values);
		if (pmaps_level_2 == NULL) {
			pmaps_level_2 = lhmsv_alloc();
			lhmslv_put(pstate->pmaps_level_1, slls_copy(pgroup_by_field_values), pmaps_level_2);
		}

		sllse_t* pa = pstate->pvalue_field_names->phead;
		sllse_t* pb =         pvalue_field_values->phead;
		// for "x", "y" and "1", "2"
		for ( ; pa != NULL && pb != NULL; pa = pa->pnext, pb = pb->pnext) {
			char* value_field_name = pa->value;
			char* value_field_sval = pb->value;
			int   have_dval = FALSE;
			double value_field_dval = -999.0;

			lhmsv_t* pmaps_level_3 = lhmsv_get(pmaps_level_2, value_field_name);
			if (pmaps_level_3 == NULL) {
				pmaps_level_3 = lhmsv_alloc();
				lhmsv_put(pmaps_level_2, value_field_name, pmaps_level_3);
			}

			// for "sum", "count"
			sllse_t* pc = pstate->paccumulator_names->phead;
			for ( ; pc != NULL; pc = pc->pnext) {
				char* acc_name = pc->value;
				acc_t* pacc = lhmsv_get(pmaps_level_3, acc_name);
				if (pacc == NULL) {
					pacc = make_acc(acc_name, &pctx->statx);
					if (pacc == NULL) {
						fprintf(stderr, "mlr stats1: accumulator \"%s\" not found.\n",
							acc_name);
						exit(1);
					}
					lhmsv_put(pmaps_level_3, acc_name, pacc);
				}
				if (pacc == NULL) {
					// xxx needs argv[0] somewhere ...
					fprintf(stderr, "stats1: could not find accumulator named \"%s\".\n", acc_name);
					exit(1);
				}

				if (pacc->psput_func != NULL) {
					pacc->psput_func(pacc->pvstate, value_field_sval);
				}
				if (pacc->pdput_func != NULL) {
					if (!have_dval) {
						value_field_dval = mlr_double_from_string_or_die(value_field_sval);
						have_dval = TRUE;
					}
					pacc->pdput_func(pacc->pvstate, value_field_dval);
				}
			}
		}

		lrec_free(pinrec);
		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->pmaps_level_1->phead; pa != NULL; pa = pa->pnext) {
			lrec_t* poutrec = lrec_unbacked_alloc();

			slls_t* pgroup_by_field_values = pa->key;

			// Add in a=s,b=t fields:
			sllse_t* pb = pstate->pgroup_by_field_names->phead;
			sllse_t* pc =         pgroup_by_field_values->phead;
			for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
				lrec_put(poutrec, pb->value, pc->value, 0);
			}

			// Add in fields such as x_sum=#, y_count=#, etc.:
			lhmsv_t* pmaps_level_2 = pa->value;
			// for "x", "y"
			for (lhmsve_t* pd = pmaps_level_2->phead; pd != NULL; pd = pd->pnext) {
				char* value_field_name = pd->key;
				lhmsv_t* pmaps_level_3 = pd->value;
				// for "count", "sum"
				for (lhmsve_t* pe = pmaps_level_3->phead; pe != NULL; pe = pe->pnext) {
					char* acc_name = pe->key;
					acc_t* pacc = pe->value;

					char free_flags = LREC_FREE_ENTRY_KEY;
					char* key = mlr_paste_3_strings(value_field_name, "_", acc_name);
					char* val = pacc->pget_func(pacc->pvstate, &free_flags);
					lrec_put(poutrec, key, val, free_flags);
				}
			}

			sllv_add(poutrecs, poutrec);
		}
		sllv_add(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static void mapper_stats1_free(void* pvstate) {
	mapper_stats1_state_t* pstate = pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_names);
	slls_free(pstate->pgroup_by_field_names);
	// xxx free the level-2's 1st
	lhmslv_free(pstate->pmaps_level_1);
}

mapper_t* mapper_stats1_alloc(slls_t* paccumulator_names, slls_t* pvalue_field_names, slls_t* pgroup_by_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats1_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_stats1_state_t));

	pstate->paccumulator_names    = paccumulator_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->pmaps_level_1         = lhmslv_alloc();

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
