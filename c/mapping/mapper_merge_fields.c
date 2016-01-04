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
#include "containers/mlrval.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define MERGE_BY_NAME_LIST  0xef01
#define MERGE_BY_NAME_REGEX 0xef02
#define MERGE_BY_COLLAPSING 0xef03
#define MERGE_UNSPECIFIED   0xef04

// ================================================================
struct _merge_fields_t; // forward reference for method definitions
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
	slls_t*  paccumulator_names;
	slls_t*  pvalue_field_names;
	sllv_t*  pvalue_field_regexes;
	int      allow_int_float;
	int      keep_input_fields;
	string_builder_t* psb;
} mapper_merge_fields_state_t;

// ----------------------------------------------------------------
static void      mapper_merge_fields_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_merge_fields_alloc(slls_t* paccumulator_names, int do_which,
	slls_t* pvalue_field_names, char* output_field_basename, int allow_int_float, int keep_input_fields);
static void      mapper_merge_fields_free(mapper_t* pmapper);
static sllv_t*   mapper_merge_fields_process_by_name_list(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_merge_fields_process_by_name_regex(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_merge_fields_process_by_collapsing(lrec_t* pinrec, context_t* pctx, void* pvstate);

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
	//for (int i = 0; i < merge_fields_lookup_table_length; i++) {
		//fprintf(o, "  %-9s %s\n", merge_fields_lookup_table[i].name, merge_fields_lookup_table[i].desc);
	//}
	fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics\n");
	fprintf(o, "-r {a,b,c}  xxx describe me\n");
	fprintf(o, "-c {a,b,c}  xxx describe me\n");
	fprintf(o, "-o {x}      xxx describe me\n");
	fprintf(o, "-k          xxx put description here.\n");
	fprintf(o, "-F          Computes integerable things (e.g. count) in floating point.\n");
	fprintf(o, "Example: %s %s -a min,max -f 'bytes_.*'\n", argv0, verb);
}

// mlr merge-fields -k -a min,p50,max -f a,b,c -o foo
// mlr merge-fields -k -a min,p50,max -r 'bytes_.*,byte_.*' -o bytes
// mlr merge-fields -c in_,out_ -a sum

static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* paccumulator_names    = NULL;
	slls_t* pvalue_field_names    = NULL;
	char*   output_field_basename = NULL;
	int     allow_int_float       = TRUE;
	int     keep_input_fields     = FALSE;
	int     do_which              = MERGE_UNSPECIFIED;

	char* verb = argv[(*pargi)++];

	int argi = *pargi;

	while (argi < argc && argv[argi][0] == '-') {

		if (streq(argv[argi], "-a")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			if (pvalue_field_names != NULL)
				slls_free(pvalue_field_names);
			paccumulator_names = slls_from_line(argv[argi+1], ',', FALSE);
			argi += 2;

		} else if (streq(argv[argi], "-f")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			if (pvalue_field_names != NULL)
				slls_free(pvalue_field_names);
			pvalue_field_names = slls_from_line(argv[argi+1], ',', FALSE);
			do_which = MERGE_BY_NAME_LIST;
			argi += 2;
		} else if (streq(argv[argi], "-r")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			if (pvalue_field_names != NULL)
				slls_free(pvalue_field_names);
			pvalue_field_names = slls_from_line(argv[argi+1], ',', FALSE);
			do_which = MERGE_BY_NAME_REGEX;
			argi += 2;
		} else if (streq(argv[argi], "-c")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			if (pvalue_field_names != NULL) {
				slls_free(pvalue_field_names);
				pvalue_field_names = NULL;
			}
			pvalue_field_names = slls_from_line(argv[argi+1], ',', FALSE);
			do_which = MERGE_BY_COLLAPSING;
			argi += 2;

		} else if (streq(argv[argi], "-o")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			output_field_basename = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "-k")) {
			keep_input_fields = TRUE;
			argi += 1;
		} else if (streq(argv[argi], "-F")) {
			allow_int_float = FALSE;
			argi += 1;
		} else {
			mapper_merge_fields_usage(stderr, argv[0], verb);
			return NULL;
		}
	}

	if (paccumulator_names == NULL) {
		mapper_merge_fields_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pvalue_field_names == NULL) {
		mapper_merge_fields_usage(stderr, argv[0], verb);
		return NULL;
	}

	*pargi = argi;
	return mapper_merge_fields_alloc(paccumulator_names, do_which,
		pvalue_field_names, output_field_basename, allow_int_float, keep_input_fields);
}

// ----------------------------------------------------------------
static mapper_t* mapper_merge_fields_alloc(slls_t* paccumulator_names, int do_which,
	slls_t* pvalue_field_names, char* output_field_basename, int allow_int_float, int keep_input_fields)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_merge_fields_state_t* pstate  = mlr_malloc_or_die(sizeof(mapper_merge_fields_state_t));

	pstate->paccumulator_names   = paccumulator_names;
	pstate->pvalue_field_names   = pvalue_field_names;
	pstate->pvalue_field_regexes = sllv_alloc(); // xxx temp
	for (sllse_t* pa = pvalue_field_names->phead; pa != NULL; pa = pa->pnext) {
		char* value_field_name = pa->value;
		regex_t* pvalue_field_regex = mlr_malloc_or_die(sizeof(regex_t));
		regcomp_or_die(pvalue_field_regex, value_field_name, 0);
		sllv_add(pstate->pvalue_field_regexes, pvalue_field_regex);
	}
	pstate->allow_int_float      = allow_int_float;
	pstate->keep_input_fields    = keep_input_fields;
	pstate->psb                  = sb_alloc(32); // xxx need #define for length

	pmapper->pvstate = pstate;
	pmapper->pprocess_func = (do_which == MERGE_BY_NAME_LIST) ? mapper_merge_fields_process_by_name_list :
		(do_which == MERGE_BY_NAME_REGEX) ? mapper_merge_fields_process_by_name_regex :
		mapper_merge_fields_process_by_collapsing;
	// xxx split out x 3?
	pmapper->pfree_func = mapper_merge_fields_free;

	return pmapper;
}

static void mapper_merge_fields_free(mapper_t* pmapper) {
	mapper_merge_fields_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_names);
	for (sllve_t* pa = pstate->pvalue_field_regexes->phead; pa != NULL; pa = pa->pnext) {
		regex_t* pvalue_field_regex = pa->pvdata;
		regfree(pvalue_field_regex);
	}
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_merge_fields_process_by_name_list(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return NULL;

	mapper_merge_fields_state_t* pstate = pvstate;
	sllv_t* paccs = sllv_alloc();
	for (sllse_t* pa = pstate->paccumulator_names->phead; pa != NULL; pa = pa->pnext) {
//		char* accumulator_name = pa->value;
//		acc_t* pacc = alloc(pstate->output_field_basename, accumulator_name, pstate->allow_int_float);
		void* pacc = NULL; // xxx temp
		sllv_add(paccs, pacc);
	}

	for (sllse_t* pb = pstate->pvalue_field_names->phead; pb != NULL; pb = pb->pnext) {
		char* field_name = pb->value;
//   matches = FALSE
//   for regex in regexes {
//     if field.key matches regex {
//       matches = TRUE
//       break
//     }
//   }
//   if matches {
//     accumulator1.ningest(field.value)
//     accumulator2.ningest(field.value)
//     accumulator3.ningest(field.value)
//   }
		if (!pstate->keep_input_fields)
			lrec_remove(pinrec, field_name);
	}

	for (sllve_t* pz = paccs->phead; pz != NULL; pz = pz->pnext) {
		// acc_t* pacc = pz->pvdata;
		// pacc->emit(pstate->output_field_basename, xxx acc_name, pinrec);
		// pacc->freefunc();
	}
	sllv_free(paccs);
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_merge_fields_process_by_name_regex(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return NULL;

	//mapper_merge_fields_state_t* pstate = pvstate;
	return sllv_single(pinrec);
}

// mlr merge-fields -k -a min,p50,max -r 'bytes_.*,byte_.*' -o bytes
// accumulator1 = alloc("bytes", "min", TRUE)
// accumulator2 = alloc("bytes", "p50", TRUE)
// accumulator3 = alloc("bytes", "max", TRUE)
// for field in inrec {
//   matches = FALSE
//   for regex in regexes {
//     if field.key matches regex {
//       matches = TRUE
//       break
//     }
//   }
//   if matches {
//     accumulator1.ningest(field.value)
//     accumulator2.ningest(field.value)
//     accumulator3.ningest(field.value)
//   }
//   if !keep
//     lrec_remove(pinrec, field.key)
// }
// accumulator1.emit("bytes", "min", inrec)
// accumulator2.emit("bytes", "p50", inrec)
// accumulator3.emit("bytes", "max", inrec)
// accumulator1.free
// accumulator2.free
// accumulator3.free

// ----------------------------------------------------------------
// mlr merge -c in_,out_ -a sum
// a_in_x  1     a_sum_x 3
// a_out_x 2     b_sum_y 4
// b_in_y  4     b_sum_x 8
// b_out_x 8

static sllv_t* mapper_merge_fields_process_by_collapsing(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return NULL;

	mapper_merge_fields_state_t* pstate = pvstate;
	lhmsv_t* short_names_to_acc_maps = lhmsv_alloc();

	for (lrece_t* pa = pinrec->phead; pa != NULL; /* increment inside loop */ ) {
		int matched = FALSE;
		for (sllve_t* pb = pstate->pvalue_field_regexes->phead; pb != NULL && !matched; pb = pb->pnext) {
			regex_t* pvalue_field_regex = pb->pvdata;

			char* short_name = regex_sub(pa->key, pvalue_field_regex, pstate->psb, "", &matched, NULL);
			if (matched) {
				lhmsv_t* acc_map_for_short_name = lhmsv_get(short_names_to_acc_maps, short_name);
				if (acc_map_for_short_name == NULL) { // First such
					acc_map_for_short_name = lhmsv_alloc();
					for (sllse_t* pc = pstate->paccumulator_names->phead; pc != NULL; pc = pc->pnext) {
						char* acc_name = pc->value;
//						acc_t* pacc = alloc(short_name, "min", pstate->allow_int_float);
						void* pacc = NULL; // xxx temp
//						// xxx implement free-flags here (& for all lhm's) for copy-reduction
						lhmsv_put(acc_map_for_short_name, acc_name, pacc);
					}
//					// xxx implement free-flags here (& for all lhm's) for copy-reduction
					lhmsv_put(short_names_to_acc_maps, short_name, acc_map_for_short_name);
				}
				for (lhmsve_t* pd = acc_map_for_short_name->phead; pd != NULL; pd = pd->pnext) {
//					acc_t* pacc = pd->pvvalue;
//					pacc->ningest(field.value)
				}
				if (!pstate->keep_input_fields) {
					// We are modifying the lrec while iterating over it.
					lrece_t* pnext = pa->pnext;
					lrec_unlink(pinrec, pa);
					pa = pnext;
				} else {
					pa = pa->pnext;
				}
				break;
			}
		}
		if (!matched)
			pa = pa->pnext;
	}

	for (lhmsve_t* pe = short_names_to_acc_maps->phead; pe != NULL; pe = pe->pnext) {
		//char* short_name = pe->key;
		lhmsv_t* acc_map_for_short_name = pe->pvvalue;
		for (lhmsve_t* pf = acc_map_for_short_name->phead; pf != NULL; pf = pf->pnext) {
			// acc.emit(short_name, acc_name, inrec)
			// acc.freefunc()
		}
		lhmsv_free(acc_map_for_short_name);
	}
	lhmsv_free(short_names_to_acc_maps);
	return sllv_single(pinrec);
}


// ================================================================
//static void mapper_merge_fields_emit_all(lrec_t* pinrec, mapper_merge_fields_state_t* pstate) {
//}

//static merge_fields_t* make_acc(char* value_field_name, char* merge_fields_name, int allow_int_float) {
//	for (int i = 0; i < merge_fields_lookup_table_length; i++)
//		if (streq(merge_fields_name, merge_fields_lookup_table[i].name))
//			return merge_fields_lookup_table[i].palloc_func(value_field_name, merge_fields_name, allow_int_float);
//	return NULL;
//}

// ----------------------------------------------------------------
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
