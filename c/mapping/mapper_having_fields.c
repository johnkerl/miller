#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmsi.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_having_fields_state_t {
	slls_t*  pfield_names;
	lhmsi_t* pfield_name_set; // xxx make this a true set now that i wrote one :/ ;)
	int      criterion;
} mapper_having_fields_state_t;

// ----------------------------------------------------------------
// record = a,b,c,d,e
// at least b,c
static sllv_t* mapper_having_fields_at_least_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL)
		return sllv_single(NULL);
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	int num_found = 0;
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (lhmsi_has_key(pstate->pfield_name_set, pe->key)) {
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
		if (!lhmsi_has_key(pstate->pfield_name_set, pe->key)) {
			lrec_free(pinrec);
			return NULL;
		}
	}
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_having_fields_at_most_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL)
		return sllv_single(NULL); // xxx cmt all of these, in all mappers
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (!lhmsi_has_key(pstate->pfield_name_set, pe->key)) {
			lrec_free(pinrec);
			return NULL; // xxx cmt all of these, in all mappers
		}
	}
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static void mapper_having_fields_free(void* pvstate) {
	mapper_having_fields_state_t* pstate = (mapper_having_fields_state_t*)pvstate;
	if (pstate->pfield_names != NULL)
		slls_free(pstate->pfield_names);
	if (pstate->pfield_name_set != NULL)
		lhmsi_free(pstate->pfield_name_set);
}

static mapper_t* mapper_having_fields_alloc(slls_t* pfield_names, int criterion) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_having_fields_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_having_fields_state_t));
	pstate->pfield_names    = pfield_names;
	pstate->pfield_name_set = lhmsi_alloc();
	for (sllse_t* pe = pfield_names->phead; pe != NULL; pe = pe->pnext)
		lhmsi_put(pstate->pfield_name_set, pe->value, 0);
	pstate->criterion = criterion;

	pmapper->pvstate = (void*)pstate;
	if (criterion == HAVING_FIELDS_AT_LEAST)
		pmapper->pprocess_func = mapper_having_fields_at_least_process;
	else if (criterion == HAVING_FIELDS_WHICH_ARE)
		pmapper->pprocess_func = mapper_having_fields_which_are_process;
	else if (criterion == HAVING_FIELDS_AT_MOST)
		pmapper->pprocess_func = mapper_having_fields_at_most_process;
	pmapper->pfree_func = mapper_having_fields_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_having_fields_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "--at-least  {a,b,c}\n");
	fprintf(stdout, "--which-are {a,b,c}\n");
	fprintf(stdout, "--at-most   {a,b,c}\n");
	fprintf(stdout, "Conditionally passes through records depending on each record's field names.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_having_fields_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_names  = NULL;
	int     criterion     = FALSE;

	char* verb = argv[(*pargi)++];

	int argi = *pargi;
	while (argi < argc && argv[argi][0] == '-') {
		if (streq(argv[argi], "--at-least")) {
			criterion = HAVING_FIELDS_AT_LEAST;
		} else if (streq(argv[argi], "--which-are")) {
			criterion = HAVING_FIELDS_WHICH_ARE;
		} else if (streq(argv[argi], "--at-most")) {
			criterion = HAVING_FIELDS_AT_MOST;
		} else {
			mapper_having_fields_usage(argv[0], verb);
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
		mapper_having_fields_usage(argv[0], verb);
		return NULL;
	}
	if (criterion == FALSE) {
		mapper_having_fields_usage(argv[0], verb);
		return NULL;
	}

	*pargi = argi;
	return mapper_having_fields_alloc(pfield_names, criterion);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_having_fields_setup = {
	.verb = "having-fields",
	.pusage_func = mapper_having_fields_usage,
	.pparse_func = mapper_having_fields_parse_cli,
};
