#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_cut_state_t {
	ap_state_t* pargp;
	slls_t*  pfield_name_list;
	hss_t*   pfield_name_set;
	regex_t* regexes;
	int      nregex;
	int      do_arg_order;
	int      do_complement;
} mapper_cut_state_t;

static void      mapper_cut_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_cut_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_cut_alloc(ap_state_t* pargp, slls_t* pfield_name_list,
	int do_arg_order, int do_complement, int do_regexes);
static void      mapper_cut_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_cut_process_no_regexes(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_cut_process_with_regexes(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_cut_setup = {
	.verb = "cut",
	.pusage_func = mapper_cut_usage,
	.pparse_func = mapper_cut_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_cut_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Passes through input records with specified fields included/excluded.\n");
	fprintf(o, "-f {a,b,c}       Field names to include for cut.\n");
	fprintf(o, "-o               Retain fields in the order specified here in the argument list.\n");
	fprintf(o, "                 Default is to retain them in the order found in the input data.\n");
	fprintf(o, "-x|--complement  Exclude, rather than include, field names specified by -f.\n");
	fprintf(o, "-r               Treat field names as regular expressions. \"ab\", \"a.*b\" will\n");
	fprintf(o, "                 match any field name containing the substring \"ab\" or matching\n");
	fprintf(o, "                 \"a.*b\", respectively; anchors of the form \"^ab$\", \"^a.*b$\" may\n");
	fprintf(o, "                 be used. The -o flag is ignored when -r is present.\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  %s %s -f hostname,status\n", argv0, verb);
	fprintf(o, "  %s %s -x -f hostname,status\n", argv0, verb);
	fprintf(o, "  %s %s -r -f '^status$,sda[0-9]'\n", argv0, verb);
	fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"'\n", argv0, verb);
	fprintf(o, "  %s %s -r -f '^status$,\"sda[0-9]\"i' (this is case-insensitive)\n", argv0, verb);
}

// ----------------------------------------------------------------
static mapper_t* mapper_cut_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pfield_name_list  = NULL;
	int     do_arg_order  = FALSE;
	int     do_complement = FALSE;
	int     do_regexes    = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f",    &pfield_name_list);
	ap_define_true_flag(pstate, "-o",           &do_arg_order);
	ap_define_true_flag(pstate, "-x",           &do_complement);
	ap_define_true_flag(pstate, "--complement", &do_complement);
	ap_define_true_flag(pstate, "-r",           &do_regexes);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_cut_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_name_list == NULL) {
		mapper_cut_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_cut_alloc(pstate, pfield_name_list, do_arg_order, do_complement, do_regexes);
}

// ----------------------------------------------------------------
static mapper_t* mapper_cut_alloc(ap_state_t* pargp, slls_t* pfield_name_list,
	int do_arg_order, int do_complement, int do_regexes)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_cut_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_cut_state_t));
	pstate->pargp = pargp;
	if (!do_regexes) {
		pstate->pfield_name_list   = pfield_name_list;
		slls_reverse(pstate->pfield_name_list);
		pstate->pfield_name_set    = hss_from_slls(pfield_name_list);
		pstate->nregex             = 0;
		pstate->regexes            = NULL;
		pmapper->pprocess_func     = mapper_cut_process_no_regexes;
	} else {
		pstate->pfield_name_list   = NULL;
		pstate->pfield_name_set    = NULL;
		pstate->nregex = pfield_name_list->length;
		pstate->regexes = mlr_malloc_or_die(pstate->nregex * sizeof(regex_t));
		int i = 0;
		for (sllse_t* pe = pfield_name_list->phead; pe != NULL; pe = pe->pnext, i++) {
			// Let them type in a.*b if they want, or "a.*b", or "a.*b"i.
			// Strip off the leading " and trailing " or "i.
			regcomp_or_die_quoted(&pstate->regexes[i], pe->value, REG_NOSUB);
		}
		slls_free(pfield_name_list);
		pmapper->pprocess_func = mapper_cut_process_with_regexes;
	}
	pstate->do_arg_order  = do_arg_order;
	pstate->do_complement = do_complement;

	pmapper->pvstate      = (void*)pstate;
	pmapper->pfree_func   = mapper_cut_free;

	return pmapper;
}

static void mapper_cut_free(mapper_t* pmapper, context_t* _) {
	mapper_cut_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pfield_name_list);
	hss_free(pstate->pfield_name_set);
	for (int i = 0; i < pstate->nregex; i++)
		regfree(&pstate->regexes[i]);
	free(pstate->regexes);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_cut_process_no_regexes(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_cut_state_t* pstate = (mapper_cut_state_t*)pvstate;
		if (!pstate->do_complement) {
			// Loop over the record and free the fields not in the
			// to-be-retained set, being careful about the fact that we're
			// modifying what we're looping over.
			for (lrece_t* pe = pinrec->phead; pe != NULL; /* next in loop */) {
				if (!hss_has(pstate->pfield_name_set, pe->key)) {
					lrece_t* pf = pe->pnext;
					lrec_remove(pinrec, pe->key);
					pe = pf;
				} else {
					pe = pe->pnext;
				}
			}
			if (pstate->do_arg_order) {
				// OK since the field-name list was reversed at construction time.
				for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
					char* field_name = pe->value;
					lrec_move_to_head(pinrec, field_name);
				}
			}
			return sllv_single(pinrec);
		} else {
			for (sllse_t* pe = pstate->pfield_name_list->phead; pe != NULL; pe = pe->pnext) {
				char* field_name = pe->value;
				lrec_remove(pinrec, field_name);
			}
			return sllv_single(pinrec);
		}
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_cut_process_with_regexes(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_cut_state_t* pstate = (mapper_cut_state_t*)pvstate;
		// Loop over the record and free the fields to be discarded, being
		// careful about the fact that we're modifying what we're looping over.
		for (lrece_t* pe = pinrec->phead; pe != NULL; /* next in loop */) {
			int matches_any = FALSE;
			for (int i = 0; i < pstate->nregex; i++) {
				if (regmatch_or_die(&pstate->regexes[i], pe->key, 0, NULL)) {
					matches_any = TRUE;
					break;
				}
			}
			if (matches_any ^ pstate->do_complement) {
				pe = pe->pnext;
			} else {
				lrece_t* pf = pe->pnext;
				lrec_remove(pinrec, pe->key);
				pe = pf;
			}
		}
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}
