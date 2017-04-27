#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mlr_globals.h"
#include "lib/mlrstat.h"
#include "cli/argparse.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "containers/mlrval.h"
#include "mapping/mappers.h"
#include "mapping/stats1_accumulators.h"

typedef enum _merge_by_t {
	MERGE_BY_NAME_LIST,
	MERGE_BY_NAME_REGEX,
	MERGE_BY_COLLAPSING,
	MERGE_UNSPECIFIED
} merge_by_t;

#define SB_ALLOC_LENGTH 32

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
	char*    output_field_basename;
	int      allow_int_float;
	int      do_interpolated_percentiles;
	int      keep_input_fields;
	string_builder_t* psb;
} mapper_merge_fields_state_t;

// ----------------------------------------------------------------
static void      mapper_merge_fields_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_merge_fields_alloc(slls_t* paccumulator_names, merge_by_t do_which,
	slls_t* pvalue_field_names, char* output_field_basename, int allow_int_float, int do_interpolated_percentiles,
	int keep_input_fields);
static void      mapper_merge_fields_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_merge_fields_process_by_name_list(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_merge_fields_process_by_name_regex(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_merge_fields_process_by_collapsing(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_merge_fields_setup = {
	.verb        = "merge-fields",
	.pusage_func = mapper_merge_fields_usage,
	.pparse_func = mapper_merge_fields_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_merge_fields_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Computes univariate statistics for each input record, accumulated across\n");
	fprintf(o, "specified fields.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-a {sum,count,...}  Names of accumulators. One or more of:\n");
	for (int i = 0; i < stats1_acc_lookup_table_length; i++) {
		fprintf(o, "  %-9s %s\n", stats1_acc_lookup_table[i].name, stats1_acc_lookup_table[i].desc);
	}
	fprintf(o, "-f {a,b,c}  Value-field names on which to compute statistics. Requires -o.\n");
	fprintf(o, "-r {a,b,c}  Regular expressions for value-field names on which to compute\n");
	fprintf(o, "            statistics. Requires -o.\n");
	fprintf(o, "-c {a,b,c}  Substrings for collapse mode. All fields which have the same names\n");
	fprintf(o, "            after removing substrings will be accumulated together. Please see\n");
	fprintf(o, "            examples below.\n");
	fprintf(o, "-i          Use interpolated percentiles, like R's type=7; default like type=1.\n");
	fprintf(o, "            Not sensical for string-valued fields.\n");
	fprintf(o, "-o {name}   Output field basename for -f/-r.\n");
	fprintf(o, "-k          Keep the input fields which contributed to the output statistics;\n");
	fprintf(o, "            the default is to omit them.\n");
	fprintf(o, "-F          Computes integerable things (e.g. count) in floating point.\n");
	fprintf(o, "\n");
	fprintf(o, "String-valued data make sense unless arithmetic on them is required,\n");
	fprintf(o, "e.g. for sum, mean, interpolated percentiles, etc. In case of mixed data,\n");
	fprintf(o, "numbers are less than strings.\n");
	fprintf(o, "\n");
	fprintf(o, "Example input data: \"a_in_x=1,a_out_x=2,b_in_y=4,b_out_x=8\".\n");
	fprintf(o, "Example: %s %s -a sum,count -f a_in_x,a_out_x -o foo\n", argv0, verb);
	fprintf(o, "  produces \"b_in_y=4,b_out_x=8,foo_sum=3,foo_count=2\" since \"a_in_x,a_out_x\" are\n");
	fprintf(o, "  summed over.\n");
	fprintf(o, "Example: %s %s -a sum,count -r in_,out_ -o bar\n", argv0, verb);
	fprintf(o, "  produces \"bar_sum=15,bar_count=4\" since all four fields are summed over.\n");
	fprintf(o, "Example: %s %s -a sum,count -c in_,out_\n", argv0, verb);
	fprintf(o, "  produces \"a_x_sum=3,a_x_count=2,b_y_sum=4,b_y_count=1,b_x_sum=8,b_x_count=1\"\n");
	fprintf(o, "  since \"a_in_x\" and \"a_out_x\" both collapse to \"a_x\", \"b_in_y\" collapses to\n");
	fprintf(o, "  \"b_y\", and \"b_out_x\" collapses to \"b_x\".\n");
}

static mapper_t* mapper_merge_fields_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t*    paccumulator_names          = NULL;
	slls_t*    pvalue_field_names          = NULL;
	char*      output_field_basename       = NULL;
	int        allow_int_float             = TRUE;
	int        do_interpolated_percentiles = FALSE;
	int        keep_input_fields           = FALSE;
	merge_by_t do_which                    = MERGE_UNSPECIFIED;

	char* verb = argv[(*pargi)++];

	int argi = *pargi;

	while (argi < argc && argv[argi][0] == '-') {

		if (streq(argv[argi], "-a")) {
			if (argc - argi < 2) {
				mapper_merge_fields_usage(stderr, argv[0], verb);
				return NULL;
			}
			if (paccumulator_names != NULL)
				slls_free(paccumulator_names);
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
		} else if (streq(argv[argi], "-i")) {
			do_interpolated_percentiles = TRUE;
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
	if (output_field_basename == NULL) {
		if (do_which == MERGE_BY_NAME_LIST || do_which == MERGE_BY_NAME_REGEX) {
			mapper_merge_fields_usage(stderr, argv[0], verb);
			return NULL;
		}
	}

	*pargi = argi;
	return mapper_merge_fields_alloc(paccumulator_names, do_which,
		pvalue_field_names, output_field_basename, allow_int_float, do_interpolated_percentiles,
		keep_input_fields);
}

// ----------------------------------------------------------------
static mapper_t* mapper_merge_fields_alloc(slls_t* paccumulator_names, merge_by_t do_which,
	slls_t* pvalue_field_names, char* output_field_basename, int allow_int_float, int do_interpolated_percentiles,
	int keep_input_fields)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_merge_fields_state_t* pstate  = mlr_malloc_or_die(sizeof(mapper_merge_fields_state_t));

	pstate->paccumulator_names   = paccumulator_names;
	pstate->pvalue_field_names   = pvalue_field_names;
	pstate->pvalue_field_regexes = sllv_alloc();
	for (sllse_t* pa = pvalue_field_names->phead; pa != NULL; pa = pa->pnext) {
		char* value_field_name = pa->value;
		regex_t* pvalue_field_regex = mlr_malloc_or_die(sizeof(regex_t));
		regcomp_or_die(pvalue_field_regex, value_field_name, 0);
		sllv_append(pstate->pvalue_field_regexes, pvalue_field_regex);
	}
	pstate->output_field_basename       = output_field_basename;
	pstate->allow_int_float             = allow_int_float;
	pstate->do_interpolated_percentiles = do_interpolated_percentiles;
	pstate->keep_input_fields           = keep_input_fields;
	pstate->psb                         = sb_alloc(SB_ALLOC_LENGTH);

	pmapper->pvstate = pstate;
	pmapper->pprocess_func = (do_which == MERGE_BY_NAME_LIST) ? mapper_merge_fields_process_by_name_list :
		(do_which == MERGE_BY_NAME_REGEX) ? mapper_merge_fields_process_by_name_regex :
		mapper_merge_fields_process_by_collapsing;
	pmapper->pfree_func = mapper_merge_fields_free;

	return pmapper;
}

static void mapper_merge_fields_free(mapper_t* pmapper, context_t* _) {
	mapper_merge_fields_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->paccumulator_names);
	slls_free(pstate->pvalue_field_names);
	for (sllve_t* pa = pstate->pvalue_field_regexes->phead; pa != NULL; pa = pa->pnext) {
		regex_t* pvalue_field_regex = pa->pvvalue;
		regfree(pvalue_field_regex);
		free(pvalue_field_regex);
	}
	sllv_free(pstate->pvalue_field_regexes);
	sb_free(pstate->psb);
	free(pstate);
	free(pmapper);
}

// ================================================================
static sllv_t* mapper_merge_fields_process_by_name_list(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return NULL;

	mapper_merge_fields_state_t* pstate = pvstate;
	// For percentiles there is one unique accumulator given (for example) five distinct
	// names p0,p25,p50,p75,p100.  The input accumulators are unique: only one
	// percentile-keeper. There are multiple output accumulators: each references the same
	// underlying percentile-keeper but with distinct parameters.
	lhmsv_t* pinaccs = lhmsv_alloc();
	lhmsv_t* poutaccs = lhmsv_alloc();

	make_stats1_accs(pstate->output_field_basename, pstate->paccumulator_names,
	    pstate->allow_int_float, pstate->do_interpolated_percentiles, pinaccs, poutaccs);

	for (sllse_t* pb = pstate->pvalue_field_names->phead; pb != NULL; pb = pb->pnext) {
		char* field_name = pb->value;
		char* value_field_sval = lrec_get(pinrec, field_name);
		if (value_field_sval == NULL) // Key not present
			continue;

		if (*value_field_sval == 0) { // Key present with null value
			if (!pstate->keep_input_fields)
				lrec_remove(pinrec, field_name);
			continue;
		}

		int have_dval = FALSE;
		int have_nval = FALSE;
		double value_field_dval = -999.0;
		mv_t   value_field_nval = mv_absent();

		for (lhmsve_t* pc = pinaccs->phead; pc != NULL; pc = pc->pnext) {
			stats1_acc_t* pacc = pc->pvvalue;

			if (pacc->pdingest_func != NULL) {
				if (!have_dval) {
					value_field_dval = mlr_double_from_string_or_die(value_field_sval);
					have_dval = TRUE;
				}
				pacc->pdingest_func(pacc->pvstate, value_field_dval);
			}
			if (pacc->pningest_func != NULL) {
				if (!have_nval) {
					value_field_nval = pstate->allow_int_float
						? mv_scan_number_or_die(value_field_sval)
						: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
					have_nval = TRUE;
				}
				pacc->pningest_func(pacc->pvstate, &value_field_nval);
			}
			if (pacc->psingest_func != NULL) {
				pacc->psingest_func(pacc->pvstate, value_field_sval);
			}
		}

		if (!pstate->keep_input_fields)
			lrec_remove(pinrec, field_name);
	}

	for (lhmsve_t* pz = poutaccs->phead; pz != NULL; pz = pz->pnext) {
		char* acc_name = pz->key;
		stats1_acc_t* pacc = pz->pvvalue;
		pacc->pemit_func(pacc->pvstate, pstate->output_field_basename, acc_name, TRUE, pinrec);
		pacc->pfree_func(pacc);
	}
	lhmsv_free(pinaccs);
	lhmsv_free(poutaccs);

	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_merge_fields_process_by_name_regex(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return NULL;

	mapper_merge_fields_state_t* pstate = pvstate;

	// For percentiles there is one unique accumulator given (for example) five distinct
	// names p0,p25,p50,p75,p100.  The input accumulators are unique: only one
	// percentile-keeper. There are multiple output accumulators: each references the same
	// underlying percentile-keeper but with distinct parameters.

	lhmsv_t* pinaccs = lhmsv_alloc();
	lhmsv_t* poutaccs = lhmsv_alloc();

	make_stats1_accs(pstate->output_field_basename, pstate->paccumulator_names,
	    pstate->allow_int_float, pstate->do_interpolated_percentiles, pinaccs, poutaccs);

	for (lrece_t* pb = pinrec->phead; pb != NULL; /* increment inside loop */ ) {
		char* field_name = pb->key;
		int matched = FALSE;
		for (sllve_t* pc = pstate->pvalue_field_regexes->phead; pc != NULL && !matched; pc = pc->pnext) {
			regex_t* pvalue_field_regex = pc->pvvalue;
			matched = regmatch_or_die(pvalue_field_regex, field_name, 0, NULL);
			if (matched) {
				char* value_field_sval = lrec_get(pinrec, field_name);
				if (value_field_sval != NULL) { // Key not present
					int have_dval = FALSE;
					int have_nval = FALSE;
					double value_field_dval = -999.0;
					mv_t   value_field_nval = mv_absent();

					if (*value_field_sval != 0) { // Key present with null value
						for (lhmsve_t* pd = pinaccs->phead; pd != NULL; pd = pd->pnext) {
							stats1_acc_t* pacc = pd->pvvalue;

							if (pacc->pdingest_func != NULL) {
								if (!have_dval) {
									value_field_dval = mlr_double_from_string_or_die(value_field_sval);
									have_dval = TRUE;
								}
								pacc->pdingest_func(pacc->pvstate, value_field_dval);
							}
							if (pacc->pningest_func != NULL) {
								if (!have_nval) {
									value_field_nval = pstate->allow_int_float
										? mv_scan_number_or_die(value_field_sval)
										: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
									have_nval = TRUE;
								}
								pacc->pningest_func(pacc->pvstate, &value_field_nval);
							}
							if (pacc->psingest_func != NULL) {
								pacc->psingest_func(pacc->pvstate, value_field_sval);
							}
						}
					}

					if (!pstate->keep_input_fields) {
						// We are modifying the lrec while iterating over it.
						lrece_t* pnext = pb->pnext;
						lrec_unlink_and_free(pinrec, pb);
						pb = pnext;
					} else {
						pb = pb->pnext;
					}
				} else {
					pb = pb->pnext;
				}
			}
		}
		if (!matched)
			pb = pb->pnext;
	}

	for (lhmsve_t* pz = poutaccs->phead; pz != NULL; pz = pz->pnext) {
		char* acc_name = pz->key;
		stats1_acc_t* pacc = pz->pvvalue;
		pacc->pemit_func(pacc->pvstate, pstate->output_field_basename, acc_name, TRUE, pinrec);
		pacc->pfree_func(pacc);
	}
	lhmsv_free(pinaccs);
	lhmsv_free(poutaccs);

	return sllv_single(pinrec);
}

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

	// For percentiles there is one unique accumulator given (for example) five distinct
	// names p0,p25,p50,p75,p100.  The input accumulators are unique: only one
	// percentile-keeper. There are multiple output accumulators: each references the same
	// underlying percentile-keeper but with distinct parameters.

	lhmsv_t* short_names_to_in_acc_maps = lhmsv_alloc();
	lhmsv_t* short_names_to_out_acc_maps = lhmsv_alloc();

	for (lrece_t* pa = pinrec->phead; pa != NULL; /* increment inside loop */ ) {
		char* field_name = pa->key;
		int matched = FALSE;
		for (sllve_t* pb = pstate->pvalue_field_regexes->phead; pb != NULL && !matched; pb = pb->pnext) {
			regex_t* pvalue_field_regex = pb->pvvalue;
			char* short_name = regex_sub(field_name, pvalue_field_regex, pstate->psb, "", &matched, NULL);
			if (matched) {
				lhmsv_t* in_acc_map_for_short_name = lhmsv_get(short_names_to_in_acc_maps, short_name);
				lhmsv_t* out_acc_map_for_short_name = lhmsv_get(short_names_to_out_acc_maps, short_name);
				if (out_acc_map_for_short_name == NULL) { // First such

					in_acc_map_for_short_name = lhmsv_alloc();
					out_acc_map_for_short_name = lhmsv_alloc();

					make_stats1_accs(short_name, pstate->paccumulator_names,
						pstate->allow_int_float, pstate->do_interpolated_percentiles,
						in_acc_map_for_short_name, out_acc_map_for_short_name);

					lhmsv_put(short_names_to_in_acc_maps, mlr_strdup_or_die(short_name), in_acc_map_for_short_name,
						FREE_ENTRY_KEY);
					lhmsv_put(short_names_to_out_acc_maps, mlr_strdup_or_die(short_name), out_acc_map_for_short_name,
						FREE_ENTRY_KEY);

				}

				char* value_field_sval = lrec_get(pinrec, field_name);
				if (value_field_sval != NULL) { // Key present

					if (*value_field_sval != 0) { // Key present with non-null value
						for (lhmsve_t* pd = in_acc_map_for_short_name->phead; pd != NULL; pd = pd->pnext) {
							stats1_acc_t* pacc = pd->pvvalue;

							int have_dval = FALSE;
							int have_nval = FALSE;
							double value_field_dval = -999.0;
							mv_t   value_field_nval = mv_absent();

							if (pacc->pdingest_func != NULL) {
								if (!have_dval) {
									value_field_dval = mlr_double_from_string_or_die(value_field_sval);
									have_dval = TRUE;
								}
								pacc->pdingest_func(pacc->pvstate, value_field_dval);
							}
							if (pacc->pningest_func != NULL) {
								if (!have_nval) {
									value_field_nval = pstate->allow_int_float
										? mv_scan_number_or_die(value_field_sval)
										: mv_from_float(mlr_double_from_string_or_die(value_field_sval));
									have_nval = TRUE;
								}
								pacc->pningest_func(pacc->pvstate, &value_field_nval);
							}
							if (pacc->psingest_func != NULL) {
								pacc->psingest_func(pacc->pvstate, value_field_sval);
							}
						}
					}

					if (!pstate->keep_input_fields) {
						// We are modifying the lrec while iterating over it.
						lrece_t* pnext = pa->pnext;
						lrec_unlink_and_free(pinrec, pa);
						pa = pnext;
					} else {
						pa = pa->pnext;
					}
				} else {
					pa = pa->pnext;
				}
				free(short_name);
			} else {
				free(short_name);
			}
		}
		if (!matched)
			pa = pa->pnext;
	}

	for (lhmsve_t* pe = short_names_to_out_acc_maps->phead; pe != NULL; pe = pe->pnext) {
		char* short_name = pe->key;
		lhmsv_t* out_acc_map_for_short_name = pe->pvvalue;
		for (lhmsve_t* pf = out_acc_map_for_short_name->phead; pf != NULL; pf = pf->pnext) {
			char* acc_name = pf->key;
			stats1_acc_t* pacc = pf->pvvalue;
			pacc->pemit_func(pacc->pvstate, short_name, acc_name, TRUE, pinrec);
		}
	}

	for (lhmsve_t* pe = short_names_to_out_acc_maps->phead; pe != NULL; pe = pe->pnext) {
		lhmsv_t* out_acc_map_for_short_name = pe->pvvalue;
		for (lhmsve_t* pf = out_acc_map_for_short_name->phead; pf != NULL; pf = pf->pnext) {
			stats1_acc_t* pacc = pf->pvvalue;
			pacc->pfree_func(pacc);
		}
		lhmsv_free(out_acc_map_for_short_name);
	}

	for (lhmsve_t* pe = short_names_to_in_acc_maps->phead; pe != NULL; pe = pe->pnext) {
		lhmsv_t* in_acc_map_for_short_name = pe->pvvalue;
		lhmsv_free(in_acc_map_for_short_name);
	}

	lhmsv_free(short_names_to_in_acc_maps);
	lhmsv_free(short_names_to_out_acc_maps);
	return sllv_single(pinrec);
}
