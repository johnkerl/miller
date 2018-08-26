#include <stdio.h>
#include "cli/argparse.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_altkv_state_t {
	ap_state_t* pargp;
} mapper_altkv_state_t;

static void      mapper_altkv_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_altkv_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_altkv_alloc(ap_state_t* pargp);
static void      mapper_altkv_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_altkv_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_altkv_setup = {
	.verb = "altkv",
	.pusage_func = mapper_altkv_usage,
	.pparse_func = mapper_altkv_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_altkv_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [no options]\n", argv0, verb);
	fprintf(o, "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.\n");
}

static mapper_t* mapper_altkv_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_altkv_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}

	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_altkv_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_altkv_alloc(pstate);
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_altkv_alloc(ap_state_t* pargp) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_altkv_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_altkv_state_t));
	pstate->pargp   = pargp;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_altkv_process;
	pmapper->pfree_func    = mapper_altkv_free;

	return pmapper;
}

static void mapper_altkv_free(mapper_t* pmapper, context_t* _) {
	mapper_altkv_state_t* pstate = pmapper->pvstate;
	// Free the container
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_altkv_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) { // End of input stream: emit null.
		return sllv_single(NULL);
	}

	lrec_t* poutrec = lrec_unbacked_alloc();
	int output_field_number = 1;
	for (lrece_t* pe = pinrec->phead; pe != NULL; /* increment in loop */) {

		if (pe->pnext == NULL) { // Odd-numbered field count
			char* key = mlr_alloc_string_from_int(output_field_number);
			char* value = mlr_strdup_or_die(pe->value);
			lrec_put(poutrec, key, value, FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		} else {
			char* key = mlr_strdup_or_die(pe->value);
			char* value = mlr_strdup_or_die(pe->pnext->value);
			lrec_put(poutrec, key, value, FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}

		output_field_number++;
		pe = pe->pnext;
		if (pe == NULL) {// Odd-numbered field count
			break;
		}
		pe = pe->pnext;
	}

	return sllv_single(poutrec);
}
