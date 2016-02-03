#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/string_builder.h"
#include "containers/lhmss.h"
#include "containers/sllv.h"
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

// mlr reshape --wide-to-long -i X,Y,Z   -o item,price
// mlr reshape --wide-to-long -r '[XYZ]' -o item,price
// * slls_t* input_field_names
// * slls_t* input_field_regexes`
// * char* output_key_field_name
// * char* output_value_field_name
//
// for each record:
//   make a map of removed field names & values, while removing them
//   for each k-v pair:
//     lrec-copy pinrec -> poutrec
//     lrec-put(k,v)
//     append to poutrecs
//   lrec_free(pinrec)

// mlr reshape --long-to-wide -i item,price
// * char* input_key_field_name
// * char* input_value_field_name

// ================================================================

#define RENAME_SB_ALLOC_LENGTH 16

typedef struct _mapper_reshape_state_t {
	ap_state_t* pargp;

	// for wide-to-long:
	slls_t* input_field_names;
	sllv_t* input_field_regexes;
	char* output_key_field_name;
	char* output_value_field_name;
} mapper_reshape_state_t;

static void      mapper_reshape_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_reshape_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_reshape_alloc(ap_state_t* pargp/*, lhmss_t* pold_to_new, int do_regexes*/);
static void      mapper_reshape_free(mapper_t* pmapper);
static sllv_t*   mapper_reshape_wide_to_long_no_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reshape_wide_to_long_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_reshape_long_to_wide_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_reshape_setup = {
	.verb = "reshape",
	.pusage_func = mapper_reshape_usage,
	.pparse_func = mapper_reshape_parse_cli
};

// ----------------------------------------------------------------
static void mapper_reshape_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {old1,new1,old2,new2,...}\n", argv0, verb);
	fprintf(o, "xxx under construction\n");
//	fprintf(o, "Renames specified fields.\n");
//	fprintf(o, "Options:\n");
//	fprintf(o, "-r         Treat old field  names as regular expressions. \"ab\", \"a.*b\"\n");
//	fprintf(o, "           will match any field name containing the substring \"ab\" or\n");
//	fprintf(o, "           matching \"a.*b\", respectively; anchors of the form \"^ab$\",\n");
//	fprintf(o, "           \"^a.*b$\" may be used. New field names may be plain strings,\n");
//	fprintf(o, "           or may contain capture groups of the form \"\\1\" through\n");
//	fprintf(o, "           \"\\9\". Wrapping the regex in double quotes is optional, but\n");
//	fprintf(o, "           is required if you wish to follow it with 'i' to indicate\n");
//	fprintf(o, "           case-insensitivity.\n");
//	fprintf(o, "Examples:\n");
//	fprintf(o, "%s %s -f old_name,new_name'\n", argv0, verb);
//	fprintf(o, "%s %s -f old_name_1,new_name_1,old_name_2,new_name_2'\n", argv0, verb);
//	fprintf(o, "%s %s -r 'Date_[0-9]+,Date,'  Rename all such fields to be \"Date\"\n", argv0, verb);
//	fprintf(o, "%s %s -r '\"Date_[0-9]+\",Date' Same\n", argv0, verb);
//	fprintf(o, "%s %s -r 'Date_([0-9]+).*,\\1' Rename all such fields to be of the form 20151015\n", argv0, verb);
//	fprintf(o, "%s %s -r '\"name\"i,Name'       Rename \"name\", \"Name\", \"NAME\", etc. to \"Name\"\n", argv0, verb);
}

static mapper_t* mapper_reshape_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* input_field_names         = NULL;
	slls_t* input_field_regex_strings = NULL;
	slls_t* output_field_names        = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-i", &input_field_names);
	ap_define_string_list_flag(pstate, "-r", &input_field_regex_strings);
	ap_define_string_list_flag(pstate, "-o", &output_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_reshape_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_reshape_alloc(pstate/*, pold_to_new, do_regexes*/);
}

// ----------------------------------------------------------------
static mapper_t* mapper_reshape_alloc(ap_state_t* pargp/*, lhmss_t* pold_to_new, int do_regexes*/) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_reshape_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_reshape_state_t));

	pstate->pargp = pargp;
//	if (do_regexes) {
//		pmapper->pprocess_func = mapper_reshape_regex_process;
//		pstate->pold_to_new    = pold_to_new;
//		pstate->pregex_pairs   = sllv_alloc();
//
//		for (lhmsse_t* pe = pold_to_new->phead; pe != NULL; pe = pe->pnext) {
//			char* regex_string = pe->key;
//			char* replacement  = pe->value;
//
//			regex_pair_t* ppair = mlr_malloc_or_die(sizeof(regex_pair_t));
//			regcomp_or_die_quoted(&ppair->regex, regex_string, 0);
//			ppair->replacement = replacement;
//			sllv_append(pstate->pregex_pairs, ppair);
//		}
//
//		pstate->psb     = sb_alloc(RENAME_SB_ALLOC_LENGTH);
//	} else {
		pmapper->pprocess_func = mapper_reshape_wide_to_long_no_regex_process;
		pmapper->pprocess_func = mapper_reshape_wide_to_long_regex_process;
		pmapper->pprocess_func = mapper_reshape_long_to_wide_regex_process;
//		pstate->pold_to_new    = pold_to_new;
//		pstate->pregex_pairs   = NULL;
//		pstate->psb            = NULL;
//	}
	pmapper->pfree_func = mapper_reshape_free;

	pmapper->pvstate = (void*)pstate;
	return pmapper;
}

static void mapper_reshape_free(mapper_t* pmapper) {
	mapper_reshape_state_t* pstate = pmapper->pvstate;
//	lhmss_free(pstate->pold_to_new);
//	if (pstate->pregex_pairs != NULL) {
//		for (sllve_t* pe = pstate->pregex_pairs->phead; pe != NULL; pe = pe->pnext) {
//			regex_pair_t* ppair = pe->pvvalue;
//			regfree(&ppair->regex);
//			// replacement is in pthe old_to_new list, already freed
//			free(ppair);
//		}
//		sllv_free(pstate->pregex_pairs);
//	}
//	sb_free(pstate->psb);
//	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_reshape_wide_to_long_no_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL)
		return sllv_single(NULL);

	mapper_reshape_state_t* pstate = (mapper_reshape_state_t*)pvstate;

	sllv_t* poutrecs = sllv_alloc();
	lhmss_t* pairs = lhmss_alloc();
	for (sllse_t* pe = pstate->input_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* key = pe->value;
		char* value = lrec_get(pinrec, key);
		if (value != NULL)
			lhmss_put(pairs, key, value, NO_FREE);
	}

	// Unset the lrec keys after iterating over them, rather than during
	for (lhmsse_t* pf = pairs->phead; pf != NULL; pf = pf->pnext)
		lrec_remove(pinrec, pf->key);

	for (lhmsse_t* pf = pairs->phead; pf != NULL; pf = pf->pnext) {
		lrec_t* poutrec = lrec_copy(pinrec);
		lrec_put(poutrec, pstate->output_key_field_name, mlr_strdup_or_die(pf->key), FREE_ENTRY_VALUE);
		lrec_put(poutrec, pstate->output_value_field_name, mlr_strdup_or_die(pf->value), FREE_ENTRY_VALUE);
		sllv_append(poutrecs, poutrec);
	}

	lhmss_free(pairs);
	return poutrecs;
}

static sllv_t* mapper_reshape_wide_to_long_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL)
		return sllv_single(NULL);

	mapper_reshape_state_t* pstate = (mapper_reshape_state_t*)pvstate;

	sllv_t* poutrecs = sllv_alloc();
	lhmss_t* pairs = lhmss_alloc();

	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		for (sllve_t* pf = pstate->input_field_regexes->phead; pf != NULL; pf = pf->pnext) {
			regex_t* pregex = pf->pvvalue;
			if (regmatch_or_die(pregex, pe->key, 0, NULL)) {
				lhmss_put(pairs, pe->key, pe->value, NO_FREE);
				break;
			}
		}
	}

	// Unset the lrec keys after iterating over them, rather than during
	for (lhmsse_t* pg = pairs->phead; pg != NULL; pg = pg->pnext)
		lrec_remove(pinrec, pg->key);

	for (lhmsse_t* pg = pairs->phead; pg != NULL; pg = pg->pnext) {
		lrec_t* poutrec = lrec_copy(pinrec);
		lrec_put(poutrec, pstate->output_key_field_name, mlr_strdup_or_die(pg->key), FREE_ENTRY_VALUE);
		lrec_put(poutrec, pstate->output_value_field_name, mlr_strdup_or_die(pg->value), FREE_ENTRY_VALUE);
		sllv_append(poutrecs, poutrec);
	}

	lhmss_free(pairs);

	return poutrecs;
}

static sllv_t* mapper_reshape_long_to_wide_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	return sllv_single(NULL); // xxx stub
}
