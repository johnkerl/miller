#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/string_builder.h"
#include "containers/lhmss.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define RENAME_SB_ALLOC_LENGTH 16

typedef struct _regex_pair_t {
	regex_t regex;
	char*   replacement;
} regex_pair_t;

typedef struct _mapper_rename_state_t {
	ap_state_t* pargp;
	lhmss_t* pold_to_new;
	sllv_t*  pregex_pairs;
	string_builder_t* psb;
	int      do_gsub;
} mapper_rename_state_t;

static void      mapper_rename_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_rename_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_rename_alloc(ap_state_t* pargp, lhmss_t* pold_to_new, int do_regexes, int do_gsub);
static void      mapper_rename_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_rename_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_rename_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_rename_setup = {
	.verb = "rename",
	.pusage_func = mapper_rename_usage,
	.pparse_func = mapper_rename_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_rename_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {old1,new1,old2,new2,...}\n", argv0, verb);
	fprintf(o, "Renames specified fields.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-r         Treat old field  names as regular expressions. \"ab\", \"a.*b\"\n");
	fprintf(o, "           will match any field name containing the substring \"ab\" or\n");
	fprintf(o, "           matching \"a.*b\", respectively; anchors of the form \"^ab$\",\n");
	fprintf(o, "           \"^a.*b$\" may be used. New field names may be plain strings,\n");
	fprintf(o, "           or may contain capture groups of the form \"\\1\" through\n");
	fprintf(o, "           \"\\9\". Wrapping the regex in double quotes is optional, but\n");
	fprintf(o, "           is required if you wish to follow it with 'i' to indicate\n");
	fprintf(o, "           case-insensitivity.\n");
	fprintf(o, "-g         Do global replacement within each field name rather than\n");
	fprintf(o, "           first-match replacement.\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "%s %s old_name,new_name'\n", argv0, verb);
	fprintf(o, "%s %s old_name_1,new_name_1,old_name_2,new_name_2'\n", argv0, verb);
	fprintf(o, "%s %s -r 'Date_[0-9]+,Date,'  Rename all such fields to be \"Date\"\n", argv0, verb);
	fprintf(o, "%s %s -r '\"Date_[0-9]+\",Date' Same\n", argv0, verb);
	fprintf(o, "%s %s -r 'Date_([0-9]+).*,\\1' Rename all such fields to be of the form 20151015\n", argv0, verb);
	fprintf(o, "%s %s -r '\"name\"i,Name'       Rename \"name\", \"Name\", \"NAME\", etc. to \"Name\"\n", argv0, verb);
}

static mapper_t* mapper_rename_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int do_regexes = FALSE;
	int do_gsub = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-r", &do_regexes);
	ap_define_true_flag(pstate, "-g", &do_gsub);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_rename_usage(stderr, argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_rename_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (do_gsub)
		do_regexes = TRUE;

	slls_t* pnames = slls_from_line(argv[*pargi], ',', FALSE);
	if ((pnames->length % 2) != 0) {
		fprintf(stderr, "%s %s: name-list must have even length; got \"%s\".\n",
			MLR_GLOBALS.bargv0, verb, argv[*pargi]);
		mapper_rename_usage(stderr, argv[0], verb);
		return NULL;
	}
	lhmss_t* pold_to_new = lhmss_alloc();
	for (sllse_t* pe = pnames->phead; pe != NULL; pe = pe->pnext->pnext) {
		lhmss_put(pold_to_new, mlr_strdup_or_die(pe->value), mlr_strdup_or_die(pe->pnext->value),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}
	slls_free(pnames);

	*pargi += 1;
	return mapper_rename_alloc(pstate, pold_to_new, do_regexes, do_gsub);
}

// ----------------------------------------------------------------
static mapper_t* mapper_rename_alloc(ap_state_t* pargp, lhmss_t* pold_to_new, int do_regexes, int do_gsub) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_rename_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_rename_state_t));

	pstate->pargp = pargp;
	if (do_regexes) {
		pmapper->pprocess_func = mapper_rename_regex_process;
		pstate->pold_to_new    = pold_to_new;
		pstate->pregex_pairs   = sllv_alloc();

		for (lhmsse_t* pe = pold_to_new->phead; pe != NULL; pe = pe->pnext) {
			char* regex_string = pe->key;
			char* replacement  = pe->value;

			regex_pair_t* ppair = mlr_malloc_or_die(sizeof(regex_pair_t));
			regcomp_or_die_quoted(&ppair->regex, regex_string, 0);
			ppair->replacement = replacement;
			sllv_append(pstate->pregex_pairs, ppair);
		}

		pstate->psb     = sb_alloc(RENAME_SB_ALLOC_LENGTH);
		pstate->do_gsub = do_gsub;
	} else {
		pmapper->pprocess_func = mapper_rename_process;
		pstate->pold_to_new    = pold_to_new;
		pstate->pregex_pairs   = NULL;
		pstate->psb            = NULL;
		pstate->do_gsub        = FALSE;
	}
	pmapper->pfree_func = mapper_rename_free;

	pmapper->pvstate = (void*)pstate;
	return pmapper;
}

static void mapper_rename_free(mapper_t* pmapper, context_t* _) {
	mapper_rename_state_t* pstate = pmapper->pvstate;
	lhmss_free(pstate->pold_to_new);
	if (pstate->pregex_pairs != NULL) {
		for (sllve_t* pe = pstate->pregex_pairs->phead; pe != NULL; pe = pe->pnext) {
			regex_pair_t* ppair = pe->pvvalue;
			regfree(&ppair->regex);
			// replacement is in pthe old_to_new list, already freed
			free(ppair);
		}
		sllv_free(pstate->pregex_pairs);
	}
	sb_free(pstate->psb);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_rename_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_rename_state_t* pstate = (mapper_rename_state_t*)pvstate;
		for (lhmsse_t* pe = pstate->pold_to_new->phead; pe != NULL; pe = pe->pnext) {
			char* old_name = pe->key;
			char* new_name = pe->value;
			if (lrec_get(pinrec, old_name) != NULL) {
				lrec_rename(pinrec, old_name, new_name, FALSE);
			}
		}
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}

static sllv_t* mapper_rename_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_rename_state_t* pstate = (mapper_rename_state_t*)pvstate;

		for (sllve_t* pe = pstate->pregex_pairs->phead; pe != NULL; pe = pe->pnext) {
			regex_pair_t* ppair = pe->pvvalue;
			regex_t* pregex = &ppair->regex;
			char* replacement = ppair->replacement;
			for (lrece_t* pf = pinrec->phead; pf != NULL; pf = pf->pnext) {
				int matched = FALSE;
				int all_captured = FALSE;
				char* old_name = pf->key;
				if (pstate->do_gsub) {
					char free_flags = NO_FREE;
					char* new_name = regex_gsub(old_name, pregex, pstate->psb, replacement, &matched,
						&all_captured, &free_flags);
					int new_needs_freeing = FALSE;
					if (free_flags & FREE_ENTRY_VALUE)
						new_needs_freeing = TRUE;
					if (matched)
						lrec_rename(pinrec, old_name, new_name, new_needs_freeing);
				} else {
					char* new_name = regex_sub(old_name, pregex, pstate->psb, replacement, &matched,
						&all_captured);
					if (matched) {
						lrec_rename(pinrec, old_name, new_name, TRUE);
					} else {
						free(new_name);
					}
				}
			}
		}

		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}
