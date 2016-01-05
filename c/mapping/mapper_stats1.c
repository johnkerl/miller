#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/string_array.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "containers/mlrval.h"
#include "mapping/mappers.h"
#include "mapping/stats1_accumulators.h"
#include "cli/argparse.h"

static char* fake_acc_name_for_setups = "__setup_done__";

// ----------------------------------------------------------------
typedef struct _mapper_stats1_state_t {
	ap_state_t* pargp;
	slls_t*         paccumulator_names;
	string_array_t* pvalue_field_names;     // parameter
	string_array_t* pvalue_field_values;    // scratch space used per-record
	slls_t*         pgroup_by_field_names;  // parameter
	lhmslv_t*       groups;
	int             do_iterative_stats;
	int             allow_int_float;
} mapper_stats1_state_t;

static void      mapper_stats1_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_stats1_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_stats1_alloc(ap_state_t* pargp, slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int do_iterative_stats, int allow_int_float);
static void      mapper_stats1_free(mapper_t* pmapper);
static sllv_t*   mapper_stats1_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_stats1_ingest(lrec_t* pinrec, mapper_stats1_state_t* pstate);
static sllv_t*   mapper_stats1_emit_all(mapper_stats1_state_t* pstate);
static lrec_t*   mapper_stats1_emit(mapper_stats1_state_t* pstate, lrec_t* poutrec,
	char* value_field_name, char* stats1_name, lhmsv_t* acc_field_to_acc_state);

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
		fprintf(o, "  %-9s %s\n", stats1_lookup_table[i].name, stats1_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics\n");
	fprintf(o, "-g {d,e,f}  Optional group-by-field names\n");
	fprintf(o, "-s          Print iterative stats. Useful in tail -f contexts (in which\n");
	fprintf(o, "            case please avoid pprint-format output since end of input\n");
	fprintf(o, "            stream will never be seen).\n");
	fprintf(o, "-F          Computes integerable things (e.g. count) in floating point.\n");
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
	int             allow_int_float       = TRUE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate,  "-a", &paccumulator_names);
	ap_define_string_array_flag(pstate, "-f", &pvalue_field_names);
	ap_define_string_list_flag(pstate,  "-g", &pgroup_by_field_names);
	ap_define_true_flag(pstate,         "-s", &do_iterative_stats);
	ap_define_false_flag(pstate,        "-F", &allow_int_float);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_stats1_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (paccumulator_names == NULL || pvalue_field_names == NULL) {
		mapper_stats1_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_stats1_alloc(pstate, paccumulator_names, pvalue_field_names, pgroup_by_field_names,
		do_iterative_stats, allow_int_float);
}

// ----------------------------------------------------------------
static mapper_t* mapper_stats1_alloc(ap_state_t* pargp, slls_t* paccumulator_names, string_array_t* pvalue_field_names,
	slls_t* pgroup_by_field_names, int do_iterative_stats, int allow_int_float)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_stats1_state_t* pstate  = mlr_malloc_or_die(sizeof(mapper_stats1_state_t));

	pstate->pargp                 = pargp;
	pstate->paccumulator_names    = paccumulator_names;
	pstate->pvalue_field_names    = pvalue_field_names;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->pvalue_field_values   = string_array_alloc(pvalue_field_names->length);
	pstate->groups                = lhmslv_alloc();
	pstate->do_iterative_stats    = do_iterative_stats;
	pstate->allow_int_float       = allow_int_float;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_stats1_process;
	pmapper->pfree_func    = mapper_stats1_free;

	return pmapper;
}

static void mapper_stats1_free(mapper_t* pmapper) {
	mapper_stats1_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->paccumulator_names);
	string_array_free(pstate->pvalue_field_names);
	string_array_free(pstate->pvalue_field_values);
	slls_free(pstate->pgroup_by_field_names);

	// lhmslv_free and lhmsv_free will free the hashmap keys; we need to free
	// the void-star hashmap values.
	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
		lhmsv_t* pgroup_to_acc_field = pa->pvvalue;
		for (lhmsve_t* pb = pgroup_to_acc_field->phead; pb != NULL; pb = pb->pnext) {
			lhmsv_t* pacc_field_to_acc_state = pb->pvvalue;
			for (lhmsve_t* pc = pacc_field_to_acc_state->phead; pc != NULL; pc = pc->pnext) {
				if (streq(pc->key, fake_acc_name_for_setups))
					continue;
				stats1_t* pstats1 = pc->pvvalue;
				pstats1->pfree_func(pstats1);
			}
			lhmsv_free(pacc_field_to_acc_state);
		}
		lhmsv_free(pgroup_to_acc_field);
	}
	lhmslv_free(pstate->groups);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
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

	lhmsv_t* pgroup_to_acc_field = lhmslv_get(pstate->groups, pgroup_by_field_values);
	if (pgroup_to_acc_field == NULL) {
		pgroup_to_acc_field = lhmsv_alloc();
		lhmslv_put(pstate->groups, slls_copy(pgroup_by_field_values), pgroup_to_acc_field);
	}

	// for x=1 and y=2
	int n = pstate->pvalue_field_names->length;
	for (int i = 0; i < n; i++) {
		char* value_field_name = pstate->pvalue_field_names->strings[i];
		char* value_field_sval = pstate->pvalue_field_values->strings[i];

		lhmsv_t* acc_field_to_acc_state = lhmsv_get(pgroup_to_acc_field, value_field_name);
		if (acc_field_to_acc_state == NULL) {
			acc_field_to_acc_state = lhmsv_alloc();
			lhmsv_put(pgroup_to_acc_field, value_field_name, acc_field_to_acc_state);
		}

		// Look up presence of all accumulators at this level's hashmap.
		char* presence = lhmsv_get(acc_field_to_acc_state, fake_acc_name_for_setups);
		if (presence == NULL) {
			make_accs(value_field_name, pstate->paccumulator_names, pstate->allow_int_float, acc_field_to_acc_state);
			lhmsv_put(acc_field_to_acc_state, fake_acc_name_for_setups, fake_acc_name_for_setups);
		}

		if (value_field_sval == NULL)
			continue;

		int have_dval = FALSE;
		int have_nval = FALSE;
		double value_field_dval = -999.0;
		mv_t   value_field_nval = mv_from_null();

		// There isn't a one-to-one mapping between user-specified stats1_names
		// and internal stats1_t's. Here in the ingestor we feed each datum
		// into a stats1_t.  In the emitter, we loop over the stats1_names in
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
					value_field_nval = pstate->allow_int_float
						? mv_scan_number_or_die(value_field_sval)
						: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
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
static sllv_t* mapper_stats1_emit_all(mapper_stats1_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	for (lhmslve_t* pa = pstate->groups->phead; pa != NULL; pa = pa->pnext) {
		slls_t* pgroup_by_field_values = pa->key;
		lrec_t* poutrec = lrec_unbacked_alloc();

		// Add in a=s,b=t fields:
		sllse_t* pb = pstate->pgroup_by_field_names->phead;
		sllse_t* pc =         pgroup_by_field_values->phead;
		for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
			lrec_put(poutrec, pb->value, pc->value, NO_FREE);
		}

		// Add in fields such as x_sum=#, y_count=#, etc.:
		lhmsv_t* pgroup_to_acc_field = pa->pvvalue;
		// for "x", "y"
		for (lhmsve_t* pd = pgroup_to_acc_field->phead; pd != NULL; pd = pd->pnext) {
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
