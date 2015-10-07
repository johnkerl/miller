#include <regex.h>
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/hss.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_having_fields_state_t {
	slls_t*  pfield_names;
	hss_t*   pfield_name_set;
	regex_t* regexes;
	int      nregex;
} mapper_having_fields_state_t;

static void      mapper_having_fields_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_having_fields_parse_cli(int* pargi, int argc, char** argv);

static mapper_t* mapper_having_fields_alloc(slls_t* pfield_names, int criterion, int do_regexes);
static void      mapper_having_fields_free(void* pvstate);

static sllv_t*   mapper_having_fields_at_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_having_fields_which_are_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_having_fields_at_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

//static sllv_t*   mapper_having_fields_at_least_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
//static sllv_t*   mapper_having_fields_which_are_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
//static sllv_t*   mapper_having_fields_at_most_regex_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_having_fields_setup = {
	.verb = "having-fields",
	.pusage_func = mapper_having_fields_usage,
	.pparse_func = mapper_having_fields_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_having_fields_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Conditionally passes through records depending on each record's field names.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "--at-least  {a,b,c}\n");
	fprintf(o, "--which-are {a,b,c}\n");
	fprintf(o, "--at-most   {a,b,c}\n");
	fprintf(o, "With -r option, the argument to --at-least, --which-are, or --at-most is treated as\n");
	fprintf(o, "a comma-separated list of regular expressions.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_having_fields_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_names  = NULL;
	int     criterion     = FALSE;
	int     do_regexes    = FALSE;

	char* verb = argv[(*pargi)++];

	int argi = *pargi;
	while (argi < argc && argv[argi][0] == '-') {
		if (streq(argv[argi], "-r")) {
			do_regexes = TRUE;
			argi++;
			continue;
		}

		if (streq(argv[argi], "--at-least")) {
			criterion = HAVING_FIELDS_AT_LEAST;
		} else if (streq(argv[argi], "--which-are")) {
			criterion = HAVING_FIELDS_WHICH_ARE;
		} else if (streq(argv[argi], "--at-most")) {
			criterion = HAVING_FIELDS_AT_MOST;
		} else {
			mapper_having_fields_usage(stderr, argv[0], verb);
			return NULL;
		}

		if (argc - argi < 2) {
			return NULL;
		}
		if (pfield_names != NULL)
			slls_free(pfield_names);
		pfield_names = slls_from_line(argv[argi+1], ',', FALSE);
		argi += 2;
	}

	if (pfield_names == NULL) {
		mapper_having_fields_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (criterion == FALSE) {
		mapper_having_fields_usage(stderr, argv[0], verb);
		return NULL;
	}

	*pargi = argi;
	return mapper_having_fields_alloc(pfield_names, criterion, do_regexes);
}

// ----------------------------------------------------------------
static mapper_t* mapper_having_fields_alloc(slls_t* pfield_names, int criterion, int do_regexes) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_having_fields_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_having_fields_state_t));

	pmapper->pvstate = (void*)pstate;

	if (do_regexes) {
		pstate->pfield_names    = NULL;
		pstate->pfield_name_set = hss_alloc();

//		if (criterion == HAVING_FIELDS_AT_LEAST)
//			pmapper->pprocess_func = mapper_having_fields_at_least_process_regex;
//		else if (criterion == HAVING_FIELDS_WHICH_ARE)
//			pmapper->pprocess_func = mapper_having_fields_which_are_process_regex;
//		else if (criterion == HAVING_FIELDS_AT_MOST)
//			pmapper->pprocess_func = mapper_having_fields_at_most_process_regex;
//		pmapper->pfree_func = mapper_having_fields_free;

	} else {
		pstate->pfield_names    = pfield_names;
		pstate->pfield_name_set = hss_alloc();
		pstate->regexes         = NULL;
		pstate->nregex          = 0;
		for (sllse_t* pe = pfield_names->phead; pe != NULL; pe = pe->pnext)
			hss_add(pstate->pfield_name_set, pe->value);

		if (criterion == HAVING_FIELDS_AT_LEAST)
			pmapper->pprocess_func = mapper_having_fields_at_least_process;
		else if (criterion == HAVING_FIELDS_WHICH_ARE)
			pmapper->pprocess_func = mapper_having_fields_which_are_process;
		else if (criterion == HAVING_FIELDS_AT_MOST)
			pmapper->pprocess_func = mapper_having_fields_at_most_process;
		pmapper->pfree_func = mapper_having_fields_free;
	}

	return pmapper;
}

static void mapper_having_fields_free(void* pvstate) {
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	if (pstate->pfield_names != NULL)
		slls_free(pstate->pfield_names);
	if (pstate->pfield_name_set != NULL)
		hss_free(pstate->pfield_name_set);
	if (pstate->regexes != NULL)
		for (int i = 0; i < pstate->nregex; i++)
			regfree(&pstate->regexes[i]);
}

// ----------------------------------------------------------------
// record = a,b,c,d,e
// at least b,c
static sllv_t* mapper_having_fields_at_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL)
		return sllv_single(NULL);
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	int num_found = 0;
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (hss_has(pstate->pfield_name_set, pe->key)) {
			num_found++;
			if (num_found == pstate->pfield_name_set->num_occupied)
				return sllv_single(pinrec);
		}
	}
	lrec_free(pinrec);
	return NULL;
}

static sllv_t* mapper_having_fields_which_are_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL)
		return sllv_single(NULL);
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	if (pinrec->field_count != pstate->pfield_name_set->num_occupied) {
		lrec_free(pinrec);
		return NULL;
	}
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (!hss_has(pstate->pfield_name_set, pe->key)) {
			lrec_free(pinrec);
			return NULL;
		}
	}
	return sllv_single(pinrec);
}

static sllv_t* mapper_having_fields_at_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL)
		return sllv_single(NULL);
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (!hss_has(pstate->pfield_name_set, pe->key)) {
			lrec_free(pinrec);
			return NULL;
		}
	}
	return sllv_single(pinrec);
}

//// ----------------------------------------------------------------
//static sllv_t* mapper_having_fields_at_least_process_regex(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	int num_found = 0;
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		for (int i = 0; i < pstate->nregex; i++) {
//			regex_t* pregex = &pstate->regexes[i];
//			xxx make a lib wrapper for regcomp and another for regexec
//			if 
//		}
//		if (hss_has(pstate->pfield_name_set, pe->key)) {
//			num_found++;
//			if (num_found == pstate->pfield_name_set->num_occupied)
//				return sllv_single(pinrec);
//		}
//	}
//	lrec_free(pinrec);
//	return NULL;
//}
//
//static sllv_t* mapper_having_fields_which_are_process_regex(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	if (pinrec->field_count != pstate->pfield_name_set->num_occupied) {
//		lrec_free(pinrec);
//		return NULL;
//	}
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (!hss_has(pstate->pfield_name_set, pe->key)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}
//
//static sllv_t* mapper_having_fields_at_most_process_regex(lrec_t* pinrec, context_t* pctx, void* pvstate) {
//	if (pinrec == NULL)
//		return sllv_single(NULL);
//	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
//	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
//		if (!hss_has(pstate->pfield_name_set, pe->key)) {
//			lrec_free(pinrec);
//			return NULL;
//		}
//	}
//	return sllv_single(pinrec);
//}
