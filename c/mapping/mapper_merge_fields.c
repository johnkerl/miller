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
#include "containers/mlrval.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

char* merge_fields_fake_acc_name_for_setups = "__setup_done__";

// ================================================================
struct _merge_fields_t; // forward reference for method definitons
typedef void merge_fields_dingest_func_t(void* pvstate, double val);
typedef void merge_fields_ningest_func_t(void* pvstate, mv_t* pval);
typedef void merge_fields_emit_func_t(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec);
typedef void merge_fields_free_func_t(struct _merge_fields_t* pmerge_fields);

typedef struct _merge_fields_t {
	void* pvstate;
	merge_fields_dingest_func_t* pdingest_func;
	merge_fields_ningest_func_t* pningest_func;
	merge_fields_emit_func_t*    pemit_func;
	merge_fields_free_func_t*    pfree_func; // virtual destructor
} merge_fields_t;

typedef merge_fields_t* merge_fields_alloc_func_t(char* value_field_name, char* merge_fields_name, int allow_int_float);

typedef struct _mapper_merge_fields_state_t {
	ap_state_t*     pargp;
	slls_t*         paccumulator_names;
	string_array_t* pvalue_field_names;
	int             allow_int_float;
} mapper_merge_fields_state_t;

// ----------------------------------------------------------------
static void      mapper_merge_fields_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_merge_fields_alloc(ap_state_t* pargp, slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	int allow_int_float);
static void      mapper_merge_fields_free(mapper_t* pmapper);
static sllv_t*   mapper_merge_fields_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_merge_fields_emit_all(lrec_t* pinrec, mapper_merge_fields_state_t* pstate);
//static lrec_t*   mapper_merge_fields_emit(mapper_merge_fields_state_t* pstate, lrec_t* poutrec,
//	char* value_field_name, char* merge_fields_name, lhmsv_t* acc_field_to_acc_state);

static merge_fields_t* merge_fields_count_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float);
static merge_fields_t* merge_fields_sum_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float);
static merge_fields_t* merge_fields_mean_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float);
static merge_fields_t* merge_fields_min_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float);
static merge_fields_t* merge_fields_max_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float);

static merge_fields_t* make_acc(char* value_field_name, char* merge_fields_name, int allow_int_float);
static void make_accs(char* value_field_name, slls_t* paccumulator_names, int allow_int_float, lhmsv_t* acc_field_to_acc_state);

// ----------------------------------------------------------------
typedef struct _acc_lookup_t {
	char* name;
	merge_fields_alloc_func_t* palloc_func;
	char* desc;
} merge_fields_lookup_t;
static merge_fields_lookup_t merge_fields_lookup_table[] = {
	{"count",    merge_fields_count_alloc,    "Count instances of fields"},
	{"sum",      merge_fields_sum_alloc,      "Compute sums of specified fields"},
	{"mean",     merge_fields_mean_alloc,     "Compute averages (sample means) of specified fields"},
	{"min",      merge_fields_min_alloc,      "Compute minimum values of specified fields"},
	{"max",      merge_fields_max_alloc,      "Compute maximum values of specified fields"},
};
static int merge_fields_lookup_table_length = sizeof(merge_fields_lookup_table) / sizeof(merge_fields_lookup_table[0]);

// ----------------------------------------------------------------
mapper_setup_t mapper_merge_fields_setup = {
	.verb        = "merge-fields",
	.pusage_func = mapper_merge_fields_usage,
	.pparse_func = mapper_merge_fields_parse_cli
};

// ----------------------------------------------------------------
static void mapper_merge_fields_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-- UNDER CONSTRUCTION --\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-a {sum,count,...}  Names of accumulators. One or more of:\n");
	for (int i = 0; i < merge_fields_lookup_table_length; i++) {
		fprintf(o, "  %-9s %s\n", merge_fields_lookup_table[i].name, merge_fields_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics\n");
	fprintf(o, "-k          xxx put description here.\n");
	fprintf(o, "-F          Computes integerable things (e.g. count) in floating point.\n");
	fprintf(o, "Example: %s %s -a min,max -f 'bytes_.*'\n", argv0, verb);
}

static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv) {
	slls_t*         paccumulator_names    = NULL;
	string_array_t* pvalue_field_names    = NULL;
	int             allow_int_float       = TRUE;
	// xxx -k flag

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate,  "-a", &paccumulator_names);
	ap_define_string_array_flag(pstate, "-f", &pvalue_field_names);
	ap_define_false_flag(pstate,        "-F", &allow_int_float);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_merge_fields_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (paccumulator_names == NULL || pvalue_field_names == NULL) {
		mapper_merge_fields_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_merge_fields_alloc(pstate, paccumulator_names, pvalue_field_names, allow_int_float);
}

// ----------------------------------------------------------------
static mapper_t* mapper_merge_fields_alloc(ap_state_t* pargp, slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	int allow_int_float)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_merge_fields_state_t* pstate  = mlr_malloc_or_die(sizeof(mapper_merge_fields_state_t));

	pstate->pargp                 = pargp;
	pstate->paccumulator_names    = paccumulator_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->allow_int_float       = allow_int_float;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_merge_fields_process;
	pmapper->pfree_func    = mapper_merge_fields_free;

	return pmapper;
}

static void mapper_merge_fields_free(mapper_t* pmapper) {
	mapper_merge_fields_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->paccumulator_names);
	string_array_free(pstate->pvalue_field_names);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

static sllv_t* mapper_merge_fields_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_merge_fields_state_t* pstate = pvstate;
	if (pinrec == NULL) { // end of input stream
		return NULL;
	}

	mapper_merge_fields_emit_all(pinrec, pstate);
	return sllv_single(pinrec);
}

static merge_fields_t* make_acc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	for (int i = 0; i < merge_fields_lookup_table_length; i++)
		if (streq(merge_fields_name, merge_fields_lookup_table[i].name))
			return merge_fields_lookup_table[i].palloc_func(value_field_name, merge_fields_name, allow_int_float);
	return NULL;
}

// ----------------------------------------------------------------
static void mapper_merge_fields_emit_all(lrec_t* pinrec, mapper_merge_fields_state_t* pstate) {
	make_accs(NULL, NULL, FALSE, NULL); // xxx temp stub

//	for (sllve_t* pe = pstate->pregex_pairs->phead; pe != NULL; pe = pe->pnext) {
//		regex_pair_t* ppair = pe->pvdata;
//		regex_t* pregex = &ppair->regex;
//		char* replacement = ppair->replacement;
//		for (lrece_t* pf = pinrec->phead; pf != NULL; pf = pf->pnext) {
//			int matched = FALSE;
//			int all_captured = FALSE;
//			char* old_name = pf->key;
//			// xxx clean this up. maybe free-flags into regex_sub. or maybe just a needs-freeing
//			// arg to both.
//			char* new_name = regex_sub(old_name, pregex, pstate->psb, replacement, &matched, &all_captured);
//			if (matched)
//				lrec_rename(pinrec, old_name, new_name, TRUE);
//			}
//		}
//	}

}

//	mlr_reference_values_from_record(pinrec, pstate->pvalue_field_names, pstate->pvalue_field_values);
//
//	// for x=1 and y=2
//	int n = pstate->pvalue_field_names->length;
//	for (int i = 0; i < n; i++) {
//		char* value_field_name = pstate->pvalue_field_names->strings[i];
//		char* value_field_sval = pstate->pvalue_field_values->strings[i];
//
//		if (value_field_sval == NULL)
//			continue;
//
//		int have_dval = FALSE;
//		int have_nval = FALSE;
//		double value_field_dval = -999.0;
//		mv_t   value_field_nval = mv_from_null();
//
//		for (lhmsve_t* pc = acc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
//			char* merge_fields_name = pc->key;
//			if (streq(merge_fields_name, merge_fields_fake_acc_name_for_setups))
//				continue;
//			merge_fields_t* pmerge_fields = pc->pvvalue;
//
//			if (pmerge_fields->pdingest_func != NULL) {
//				if (!have_dval) {
//					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
//					have_dval = TRUE;
//				}
//				pmerge_fields->pdingest_func(pmerge_fields->pvstate, value_field_dval);
//			}
//			if (pmerge_fields->pningest_func != NULL) {
//				if (!have_nval) {
//					value_field_nval = pstate->allow_int_float
//						? mv_scan_number_or_die(value_field_sval)
//						: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
//					have_nval = TRUE;
//				}
//				pmerge_fields->pningest_func(pmerge_fields->pvstate, &value_field_nval);
//			}
//
//
//			// Add in fields such as x_sum=#, y_count=#, etc.:
//			for (sllse_t* pe = pstate->paccumulator_names->phead; pe != NULL; pe = pe->pnext) {
//				char* merge_fields_name = pe->value;
//				if (streq(merge_fields_name, merge_fields_fake_acc_name_for_setups))
//					continue;
//				merge_fields_t* pmerge_fields = lhmsv_get(acc_field_to_acc_state, merge_fields_name);
//				if (pmerge_fields == NULL) {
//					fprintf(stderr, "%s merge_fields: internal coding error: merge_fields_name \"%s\" has gone missing.\n",
//						MLR_GLOBALS.argv0, merge_fields_name);
//					exit(1);
//				}
//				pmerge_fields->pemit_func(pmerge_fields->pvstate, value_field_name, merge_fields_name, poutrec);
//			}
//		}
//	}
//	slls_free(pgroup_by_field_values);

// ----------------------------------------------------------------
static void make_accs(
	char*      value_field_name,       // input
	slls_t*    paccumulator_names,     // input
	int        allow_int_float,        // input
	lhmsv_t*   acc_field_to_acc_state) // output
{
	for (sllse_t* pc = paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
		// for "sum", "count"
		char* merge_fields_name = pc->value;

		merge_fields_t* pmerge_fields = make_acc(value_field_name, merge_fields_name, allow_int_float);
		if (pmerge_fields == NULL) {
			fprintf(stderr, "%s merge_fields: accumulator \"%s\" not found.\n",
				MLR_GLOBALS.argv0, merge_fields_name);
			exit(1);
		}
		lhmsv_put(acc_field_to_acc_state, merge_fields_name, pmerge_fields);
	}
}

// ----------------------------------------------------------------
typedef struct _merge_fields_count_state_t {
	mv_t counter;
	mv_t one;
	char* output_field_name;
} merge_fields_count_state_t;

static void merge_fields_count_emit(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec) {
	merge_fields_count_state_t* pstate = pvstate;
	lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->counter),
		FREE_ENTRY_VALUE);
}
static void merge_fields_count_free(merge_fields_t* pmerge_fields) {
	merge_fields_count_state_t* pstate = pmerge_fields->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pmerge_fields);
}
static merge_fields_t* merge_fields_count_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	merge_fields_t* pmerge_fields = mlr_malloc_or_die(sizeof(merge_fields_t));
	merge_fields_count_state_t* pstate = mlr_malloc_or_die(sizeof(merge_fields_count_state_t));
	pstate->counter = allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->one     = allow_int_float ? mv_from_int(1LL) : mv_from_float(1.0);
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", merge_fields_name);

	pmerge_fields->pvstate       = (void*)pstate;
	pmerge_fields->pdingest_func = NULL;
	pmerge_fields->pningest_func = NULL;
	pmerge_fields->pemit_func    = merge_fields_count_emit;
	pmerge_fields->pfree_func    = merge_fields_count_free;
	return pmerge_fields;
}

// ----------------------------------------------------------------
typedef struct _merge_fields_sum_state_t {
	mv_t sum;
	char* output_field_name;
	int allow_int_float;
} merge_fields_sum_state_t;
static void merge_fields_sum_ningest(void* pvstate, mv_t* pval) {
	merge_fields_sum_state_t* pstate = pvstate;
	pstate->sum = n_nn_plus_func(&pstate->sum, pval);
}
static void merge_fields_sum_emit(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec) {
	merge_fields_sum_state_t* pstate = pvstate;
	lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->sum),
		FREE_ENTRY_VALUE);
}
static void merge_fields_sum_free(merge_fields_t* pmerge_fields) {
	merge_fields_sum_state_t* pstate = pmerge_fields->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pmerge_fields);
}
static merge_fields_t* merge_fields_sum_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	merge_fields_t* pmerge_fields = mlr_malloc_or_die(sizeof(merge_fields_t));
	merge_fields_sum_state_t* pstate = mlr_malloc_or_die(sizeof(merge_fields_sum_state_t));
	pstate->allow_int_float = allow_int_float;
	pstate->sum = pstate->allow_int_float ? mv_from_int(0LL) : mv_from_float(0.0);
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", merge_fields_name);
	pmerge_fields->pvstate       = (void*)pstate;
	pmerge_fields->pdingest_func = NULL;
	pmerge_fields->pningest_func = merge_fields_sum_ningest;
	pmerge_fields->pemit_func    = merge_fields_sum_emit;
	pmerge_fields->pfree_func    = merge_fields_sum_free;
	return pmerge_fields;
}

// ----------------------------------------------------------------
typedef struct _merge_fields_mean_state_t {
	double sum;
	unsigned long long count;
	char* output_field_name;
} merge_fields_mean_state_t;
static void merge_fields_mean_dingest(void* pvstate, double val) {
	merge_fields_mean_state_t* pstate = pvstate;
	pstate->sum   += val;
	pstate->count++;
}
static void merge_fields_mean_emit(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec) {
	merge_fields_mean_state_t* pstate = pvstate;
	if (pstate->count == 0LL) {
		lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		double quot = pstate->sum / pstate->count;
		char* val = mlr_alloc_string_from_double(quot, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, pstate->output_field_name, val, FREE_ENTRY_VALUE);
	}
}
static void merge_fields_mean_free(merge_fields_t* pmerge_fields) {
	merge_fields_mean_state_t* pstate = pmerge_fields->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pmerge_fields);
}
static merge_fields_t* merge_fields_mean_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	merge_fields_t* pmerge_fields = mlr_malloc_or_die(sizeof(merge_fields_t));
	merge_fields_mean_state_t* pstate = mlr_malloc_or_die(sizeof(merge_fields_mean_state_t));
	pstate->sum         = 0.0;
	pstate->count       = 0LL;
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", merge_fields_name);

	pmerge_fields->pvstate       = (void*)pstate;
	pmerge_fields->pdingest_func = merge_fields_mean_dingest;
	pmerge_fields->pningest_func = NULL;
	pmerge_fields->pemit_func    = merge_fields_mean_emit;
	pmerge_fields->pfree_func    = merge_fields_mean_free;
	return pmerge_fields;
}

// ----------------------------------------------------------------
typedef struct _merge_fields_min_state_t {
	mv_t min;
	char* output_field_name;
} merge_fields_min_state_t;
static void merge_fields_min_ningest(void* pvstate, mv_t* pval) {
	merge_fields_min_state_t* pstate = pvstate;
	pstate->min = n_nn_min_func(&pstate->min, pval);
}
static void merge_fields_min_emit(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec) {
	merge_fields_min_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->min)) {
		lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->min),
			FREE_ENTRY_VALUE);
	}
}
static void merge_fields_min_free(merge_fields_t* pmerge_fields) {
	merge_fields_min_state_t* pstate = pmerge_fields->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pmerge_fields);
}
static merge_fields_t* merge_fields_min_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	merge_fields_t* pmerge_fields = mlr_malloc_or_die(sizeof(merge_fields_t));
	merge_fields_min_state_t* pstate = mlr_malloc_or_die(sizeof(merge_fields_min_state_t));
	pstate->min = mv_from_null();
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", merge_fields_name);
	pmerge_fields->pvstate       = (void*)pstate;
	pmerge_fields->pdingest_func = NULL;
	pmerge_fields->pningest_func = merge_fields_min_ningest;
	pmerge_fields->pemit_func    = merge_fields_min_emit;
	pmerge_fields->pfree_func    = merge_fields_min_free;
	return pmerge_fields;
}

// ----------------------------------------------------------------
typedef struct _merge_fields_max_state_t {
	mv_t max;
	char* output_field_name;
} merge_fields_max_state_t;
static void merge_fields_max_ningest(void* pvstate, mv_t* pval) {
	merge_fields_max_state_t* pstate = pvstate;
	pstate->max = n_nn_max_func(&pstate->max, pval);
}
static void merge_fields_max_emit(void* pvstate, char* value_field_name, char* merge_fields_name, lrec_t* poutrec) {
	merge_fields_max_state_t* pstate = pvstate;
	if (mv_is_null(&pstate->max)) {
		lrec_put(poutrec, pstate->output_field_name, "", NO_FREE);
	} else {
		lrec_put(poutrec, pstate->output_field_name, mv_alloc_format_val(&pstate->max),
			FREE_ENTRY_VALUE);
	}
}
static void merge_fields_max_free(merge_fields_t* pmerge_fields) {
	merge_fields_max_state_t* pstate = pmerge_fields->pvstate;
	free(pstate->output_field_name);
	free(pstate);
	free(pmerge_fields);
}
static merge_fields_t* merge_fields_max_alloc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
	merge_fields_t* pmerge_fields = mlr_malloc_or_die(sizeof(merge_fields_t));
	merge_fields_max_state_t* pstate = mlr_malloc_or_die(sizeof(merge_fields_max_state_t));
	pstate->max = mv_from_null();
	pstate->output_field_name = mlr_paste_3_strings(value_field_name, "_", merge_fields_name);
	pmerge_fields->pvstate       = (void*)pstate;
	pmerge_fields->pdingest_func = NULL;
	pmerge_fields->pningest_func = merge_fields_max_ningest;
	pmerge_fields->pemit_func    = merge_fields_max_emit;
	pmerge_fields->pfree_func    = merge_fields_max_free;
	return pmerge_fields;
}
