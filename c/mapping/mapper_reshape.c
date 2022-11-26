#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/string_builder.h"
#include "containers/lhmss.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
// WIDE:
//          time           X          Y           Z
// 1  2009-01-01  0.65473572  2.4520609 -1.46570942
// 2  2009-01-02 -0.89248112  0.2154713 -2.05357735
// 3  2009-01-03  0.98012375  1.3179287  4.64248357
// 4  2009-01-04  0.35397376  3.3765645 -0.25237774
// 5  2009-01-05  2.19357813  1.3477511  0.09719105

// LONG:
//          time  item       price
// 1  2009-01-01     X  0.65473572
// 2  2009-01-02     X -0.89248112
// 3  2009-01-03     X  0.98012375
// 4  2009-01-04     X  0.35397376
// 5  2009-01-05     X  2.19357813
// 6  2009-01-01     Y  2.45206093
// 7  2009-01-02     Y  0.21547134
// 8  2009-01-03     Y  1.31792866
// 9  2009-01-04     Y  3.37656453
// 10 2009-01-05     Y  1.34775108
// 11 2009-01-01     Z -1.46570942
// 12 2009-01-02     Z -2.05357735
// 13 2009-01-03     Z  4.64248357
// 14 2009-01-04     Z -0.25237774
// 15 2009-01-05     Z  0.09719105

// ================================================================
typedef struct _mapper_reshape_state_t {
	ap_state_t* pargp;

	// for wide-to-long:
	slls_t* input_field_names;
	sllv_t* input_field_regexes;
	char* output_key_field_name;
	char* output_value_field_name;

	// for long-to-wide:
	char* split_out_key_field_name;
	char* split_out_value_field_name;
	lhmslv_t* other_keys_to_other_values_to_buckets;
} mapper_reshape_state_t;

typedef struct _reshape_bucket_t {
	lrec_t* prepresentative;
	lhmss_t* pairs;
} reshape_bucket_t;

static void      mapper_reshape_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_reshape_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_reshape_alloc(
	ap_state_t* pargp,
	slls_t* input_field_names,
	slls_t* input_field_regex_strings,
	char*   output_key_field_name,
	char*   output_value_field_name,
	char*   split_out_key_field_name,
	char*   split_out_value_field_name);
static void      mapper_reshape_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_reshape_wide_to_long_no_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reshape_wide_to_long_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reshape_long_to_wide_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

static reshape_bucket_t* reshape_bucket_alloc(lrec_t* prepresentative);
static void reshape_bucket_free(reshape_bucket_t* pbucket);

// ----------------------------------------------------------------
mapper_setup_t mapper_reshape_setup = {
	.verb = "reshape",
	.pusage_func = mapper_reshape_usage,
	.pparse_func = mapper_reshape_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_reshape_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Wide-to-long options:\n");
	fprintf(o, "  -i {input field names}   -o {key-field name,value-field name}\n");
	fprintf(o, "  -r {input field regexes} -o {key-field name,value-field name}\n");
	fprintf(o, "  These pivot/reshape the input data such that the input fields are removed\n");
	fprintf(o, "  and separate records are emitted for each key/value pair.\n");
	fprintf(o, "  Note: this works with tail -f and produces output records for each input\n");
	fprintf(o, "  record seen.\n");
	fprintf(o, "Long-to-wide options:\n");
	fprintf(o, "  -s {key-field name,value-field name}\n");
	fprintf(o, "  These pivot/reshape the input data to undo the wide-to-long operation.\n");
	fprintf(o, "  Note: this does not work with tail -f; it produces output records only after\n");
	fprintf(o, "  all input records have been read.\n");
	fprintf(o, "\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "\n");
	fprintf(o, "  Input file \"wide.txt\":\n");
	fprintf(o, "    time       X           Y\n");
	fprintf(o, "    2009-01-01 0.65473572  2.4520609\n");
	fprintf(o, "    2009-01-02 -0.89248112 0.2154713\n");
	fprintf(o, "    2009-01-03 0.98012375  1.3179287\n");
	fprintf(o, "\n");
	fprintf(o, "  %s --pprint %s -i X,Y -o item,value wide.txt\n", argv0, verb);
	fprintf(o, "    time       item value\n");
	fprintf(o, "    2009-01-01 X    0.65473572\n");
	fprintf(o, "    2009-01-01 Y    2.4520609\n");
	fprintf(o, "    2009-01-02 X    -0.89248112\n");
	fprintf(o, "    2009-01-02 Y    0.2154713\n");
	fprintf(o, "    2009-01-03 X    0.98012375\n");
	fprintf(o, "    2009-01-03 Y    1.3179287\n");
	fprintf(o, "\n");
	fprintf(o, "  %s --pprint %s -r '[A-Z]' -o item,value wide.txt\n", argv0, verb);
	fprintf(o, "    time       item value\n");
	fprintf(o, "    2009-01-01 X    0.65473572\n");
	fprintf(o, "    2009-01-01 Y    2.4520609\n");
	fprintf(o, "    2009-01-02 X    -0.89248112\n");
	fprintf(o, "    2009-01-02 Y    0.2154713\n");
	fprintf(o, "    2009-01-03 X    0.98012375\n");
	fprintf(o, "    2009-01-03 Y    1.3179287\n");
	fprintf(o, "\n");
	fprintf(o, "  Input file \"long.txt\":\n");
	fprintf(o, "    time       item value\n");
	fprintf(o, "    2009-01-01 X    0.65473572\n");
	fprintf(o, "    2009-01-01 Y    2.4520609\n");
	fprintf(o, "    2009-01-02 X    -0.89248112\n");
	fprintf(o, "    2009-01-02 Y    0.2154713\n");
	fprintf(o, "    2009-01-03 X    0.98012375\n");
	fprintf(o, "    2009-01-03 Y    1.3179287\n");
	fprintf(o, "\n");
	fprintf(o, "  %s --pprint %s -s item,value long.txt\n", argv0, verb);
	fprintf(o, "    time       X           Y\n");
	fprintf(o, "    2009-01-01 0.65473572  2.4520609\n");
	fprintf(o, "    2009-01-02 -0.89248112 0.2154713\n");
	fprintf(o, "    2009-01-03 0.98012375  1.3179287\n");
	fprintf(o, "See also %s nest.\n", argv0);
}

static mapper_t* mapper_reshape_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* input_field_names         = NULL;
	slls_t* input_field_regex_strings = NULL;
	slls_t* output_field_names        = NULL;
	slls_t* split_out_field_names     = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-i", &input_field_names);
	ap_define_string_list_flag(pstate, "-r", &input_field_regex_strings);
	ap_define_string_list_flag(pstate, "-o", &output_field_names);
	ap_define_string_list_flag(pstate, "-s", &split_out_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_reshape_usage(stderr, argv[0], verb);
		return NULL;
	}

	char* output_key_field_name      = NULL;
	char* output_value_field_name    = NULL;
	char* split_out_key_field_name   = NULL;
	char* split_out_value_field_name = NULL;

	if (split_out_field_names == NULL) {
		// wide to long
		if (input_field_names == NULL && input_field_regex_strings == NULL) {
			mapper_reshape_usage(stderr, argv[0], verb);
			return NULL;
		}

		if (output_field_names == NULL) {
			mapper_reshape_usage(stderr, argv[0], verb);
			return NULL;
		}
		if (output_field_names->length != 2) {
			mapper_reshape_usage(stderr, argv[0], verb);
			return NULL;
		}
		output_key_field_name   = mlr_strdup_or_die(output_field_names->phead->value);
		output_value_field_name = mlr_strdup_or_die(output_field_names->phead->pnext->value);
		slls_free(output_field_names);

	} else {
		// long to wide
		if (split_out_field_names->length != 2) {
			mapper_reshape_usage(stderr, argv[0], verb);
			return NULL;
		}
		split_out_key_field_name   = mlr_strdup_or_die(split_out_field_names->phead->value);
		split_out_value_field_name = mlr_strdup_or_die(split_out_field_names->phead->pnext->value);
		slls_free(split_out_field_names);
	}

	return mapper_reshape_alloc(pstate, input_field_names, input_field_regex_strings,
		output_key_field_name, output_value_field_name,
		split_out_key_field_name, split_out_value_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reshape_alloc(
	ap_state_t* pargp,
	slls_t* input_field_names,
	slls_t* input_field_regex_strings,
	char*   output_key_field_name,
	char*   output_value_field_name,
	char*   split_out_key_field_name,
	char*   split_out_value_field_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_reshape_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_reshape_state_t));

	pstate->pargp                      = pargp;
	pstate->input_field_names          = input_field_names;
	pstate->output_key_field_name      = output_key_field_name;
	pstate->output_value_field_name    = output_value_field_name;
	pstate->split_out_key_field_name   = split_out_key_field_name;
	pstate->split_out_value_field_name = split_out_value_field_name;

	if (input_field_regex_strings == NULL) {
		pstate->input_field_regexes = NULL;
	} else {
		pstate->input_field_regexes = sllv_alloc();
		for (sllse_t* pe = input_field_regex_strings->phead; pe != NULL; pe = pe->pnext) {
			regex_t* pregex = mlr_malloc_or_die(sizeof(regex_t));
			regcomp_or_die(pregex, pe->value, 0);
			sllv_append(pstate->input_field_regexes, pregex);
		}
		slls_free(input_field_regex_strings);
	}

	if (split_out_key_field_name == NULL) {
		if (pstate->input_field_regexes == NULL)
			pmapper->pprocess_func = mapper_reshape_wide_to_long_no_regex_process;
		else
			pmapper->pprocess_func = mapper_reshape_wide_to_long_regex_process;
		pstate->other_keys_to_other_values_to_buckets = NULL;
	} else {
		pmapper->pprocess_func = mapper_reshape_long_to_wide_process;
		pstate->other_keys_to_other_values_to_buckets = lhmslv_alloc();
	}

	pmapper->pfree_func = mapper_reshape_free;

	pmapper->pvstate = (void*)pstate;
	return pmapper;
}

static void mapper_reshape_free(mapper_t* pmapper, context_t* _) {
	mapper_reshape_state_t* pstate = pmapper->pvstate;

	slls_free(pstate->input_field_names);

	free(pstate->output_key_field_name);
	free(pstate->output_value_field_name);

	free(pstate->split_out_key_field_name);
	free(pstate->split_out_value_field_name);

	if (pstate->input_field_regexes != NULL) {
		for (sllve_t* pe = pstate->input_field_regexes->phead; pe != NULL; pe = pe->pnext) {
			regex_t* pregex = pe->pvvalue;
			regfree(pregex);
			free(pregex);
		}
		sllv_free(pstate->input_field_regexes);
	}

	if (pstate->other_keys_to_other_values_to_buckets != NULL) {
		for (lhmslve_t* pe = pstate->other_keys_to_other_values_to_buckets->phead; pe != NULL; pe = pe->pnext) {
			lhmslv_t* other_values_to_buckets = pe->pvvalue;
			for (lhmslve_t* pf = other_values_to_buckets->phead; pf != NULL; pf = pf->pnext) {
				reshape_bucket_t* pbucket = pf->pvvalue;
				reshape_bucket_free(pbucket);
			}
			lhmslv_free(other_values_to_buckets);
		}
		lhmslv_free(pstate->other_keys_to_other_values_to_buckets);
	}

	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_reshape_wide_to_long_no_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);

	mapper_reshape_state_t* pstate = (mapper_reshape_state_t*)pvstate;

	sllv_t* poutrecs = sllv_alloc();
	lhmss_t* pairs = lhmss_alloc();
	for (sllse_t* pe = pstate->input_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* key = pe->value;
		char* value = lrec_get(pinrec, key);
		if (value != NULL) {
			lhmss_put(pairs, mlr_strdup_or_die(key), mlr_strdup_or_die(value), FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}
	}

	// Unset the lrec keys after iterating over them, rather than during
	for (lhmsse_t* pf = pairs->phead; pf != NULL; pf = pf->pnext)
		lrec_remove(pinrec, pf->key);

	if (pairs->num_occupied == 0) {
		sllv_append(poutrecs, pinrec);
	} else {
		for (lhmsse_t* pf = pairs->phead; pf != NULL; pf = pf->pnext) {
			lrec_t* poutrec = lrec_copy(pinrec);
			lrec_put(poutrec, pstate->output_key_field_name, mlr_strdup_or_die(pf->key), FREE_ENTRY_VALUE);
			lrec_put(poutrec, pstate->output_value_field_name, mlr_strdup_or_die(pf->value), FREE_ENTRY_VALUE);
			sllv_append(poutrecs, poutrec);
		}
		lrec_free(pinrec);
	}

	lhmss_free(pairs);
	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_reshape_wide_to_long_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);

	mapper_reshape_state_t* pstate = (mapper_reshape_state_t*)pvstate;

	sllv_t* poutrecs = sllv_alloc();
	lhmss_t* pairs = lhmss_alloc();

	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		for (sllve_t* pf = pstate->input_field_regexes->phead; pf != NULL; pf = pf->pnext) {
			regex_t* pregex = pf->pvvalue;
			if (regmatch_or_die(pregex, pe->key, 0, NULL)) {
				// Ownership-transfer of the about-to-be-freed key-value pairs from lrec to lhmss
				lhmss_put(pairs, pe->key, pe->value, pe->free_flags);
				pe->free_flags = NO_FREE;
				break;
			}
		}
	}

	// Unset the lrec keys after iterating over them, rather than during
	for (lhmsse_t* pg = pairs->phead; pg != NULL; pg = pg->pnext)
		lrec_remove(pinrec, pg->key);

	if (pairs->num_occupied == 0) {
		sllv_append(poutrecs, pinrec);
	} else {
		for (lhmsse_t* pf = pairs->phead; pf != NULL; pf = pf->pnext) {
			lrec_t* poutrec = lrec_copy(pinrec);
			lrec_put(poutrec, pstate->output_key_field_name, mlr_strdup_or_die(pf->key), FREE_ENTRY_VALUE);
			lrec_put(poutrec, pstate->output_value_field_name, mlr_strdup_or_die(pf->value), FREE_ENTRY_VALUE);
			sllv_append(poutrecs, poutrec);
		}
		lrec_free(pinrec);
	}

	lhmss_free(pairs);

	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_reshape_long_to_wide_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_reshape_state_t* pstate = (mapper_reshape_state_t*)pvstate;

	if (pinrec != NULL) { // Not end of input stream
		char* split_out_key_field_value   = lrec_get(pinrec, pstate->split_out_key_field_name);
		char* split_out_value_field_value = lrec_get(pinrec, pstate-> split_out_value_field_name);
		if (split_out_key_field_value == NULL || split_out_value_field_value == NULL)
			return sllv_single(pinrec);
		split_out_key_field_value   = mlr_strdup_or_die(split_out_key_field_value);
		split_out_value_field_value = mlr_strdup_or_die(split_out_value_field_value);
		lrec_remove(pinrec, pstate->split_out_key_field_name);
		lrec_remove(pinrec, pstate->split_out_value_field_name);

		slls_t* other_keys = mlr_reference_keys_from_record(pinrec);
		lhmslv_t* other_values_to_buckets = lhmslv_get(pstate->other_keys_to_other_values_to_buckets, other_keys);
		if (other_values_to_buckets == NULL) {
			other_values_to_buckets = lhmslv_alloc();
			lhmslv_put(pstate->other_keys_to_other_values_to_buckets,
				slls_copy(other_keys), other_values_to_buckets, FREE_ENTRY_KEY);
		}

		slls_t* other_values = mlr_reference_values_from_record(pinrec);
		reshape_bucket_t* pbucket = lhmslv_get(other_values_to_buckets, other_values);
		if (pbucket == NULL) {
			pbucket = reshape_bucket_alloc(pinrec);
			lhmslv_put(other_values_to_buckets, slls_copy(other_values), pbucket, FREE_ENTRY_KEY);
		} else {
			lrec_free(pinrec);
		}
		lhmss_put(pbucket->pairs, split_out_key_field_value, split_out_value_field_value,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);

		slls_free(other_values);
		slls_free(other_keys);

		return NULL;

	} else { // end of input stream
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pe = pstate->other_keys_to_other_values_to_buckets->phead; pe != NULL; pe = pe->pnext) {
			lhmslv_t* other_values_to_buckets = pe->pvvalue;
			for (lhmslve_t* pf = other_values_to_buckets->phead; pf != NULL; pf = pf->pnext) {
				reshape_bucket_t* pbucket = pf->pvvalue;
				lrec_t* poutrec = pbucket->prepresentative;
				pbucket->prepresentative = NULL; // ownership transfer
				for (lhmsse_t* pg = pbucket->pairs->phead; pg != NULL; pg = pg->pnext) {
					// Strings in these lrecs are backed by our multi-level hashmaps which aren't freed by our free
					// method until shutdown time (in particular, after all outrecs are emitted).
					lrec_put(poutrec, pg->key, pg->value, NO_FREE);
				}
				sllv_append(poutrecs, poutrec);
			}
		}

		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static reshape_bucket_t* reshape_bucket_alloc(lrec_t* prepresentative) {
	reshape_bucket_t* pbucket = mlr_malloc_or_die(sizeof(reshape_bucket_t));
	pbucket->prepresentative = prepresentative;
	pbucket->pairs = lhmss_alloc();
	return pbucket;
}
static void reshape_bucket_free(reshape_bucket_t* pbucket) {
	lrec_free(pbucket->prepresentative);
	lhmss_free(pbucket->pairs);
	free(pbucket);
}
