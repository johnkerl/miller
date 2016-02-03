#include <stdio.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_tac_state_t {
	sllv_t* records;
} mapper_tac_state_t;

static void      mapper_tac_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_tac_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_tac_alloc();
static void      mapper_tac_free(mapper_t* pmapper);
static sllv_t*   mapper_tac_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_tac_setup = {
	.verb = "tac",
	.pusage_func = mapper_tac_usage,
	.pparse_func = mapper_tac_parse_cli
};

// ----------------------------------------------------------------
static void mapper_tac_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Prints records in reverse order from the order in which they were encountered.\n");
}

static mapper_t* mapper_tac_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 1) {
		mapper_tac_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_tac_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_tac_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_tac_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tac_state_t));
	pstate->records = sllv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_tac_process;
	pmapper->pfree_func    = mapper_tac_free;

	return pmapper;
}

static void mapper_tac_free(mapper_t* pmapper) {
	mapper_tac_state_t* pstate = pmapper->pvstate;
	// Free the container
	sllv_free(pstate->records);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_tac_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tac_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		// The caller will free the outrecs
		sllv_append(pstate->records, pinrec);
		return NULL;
	}
	else {
		sllv_reverse(pstate->records);
		sllv_append(pstate->records, NULL);
		sllv_t* retval = pstate->records;
		pstate->records = sllv_alloc();
		return retval;
	}
}
