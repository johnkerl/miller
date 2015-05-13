#include <stdio.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_tac_state_t {
	sllv_t* records;
} mapper_tac_state_t;

// ----------------------------------------------------------------
static sllv_t* mapper_tac_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tac_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		sllv_add(pstate->records, pinrec);
		return NULL;
	}
	else {
		sllv_reverse(pstate->records);
		sllv_add(pstate->records, NULL);
		sllv_t* retval = pstate->records;
		pstate->records = sllv_alloc();
		return retval;
	}
}

// ----------------------------------------------------------------
static void mapper_tac_free(void* pvstate) {
	mapper_tac_state_t* pstate = pvstate;
	if (pstate->records != NULL)
		// xxx free the void-star payload
		sllv_free(pstate->records);
}

static mapper_t* mapper_tac_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_tac_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tac_state_t));
	pstate->records = sllv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_tac_process;
	pmapper->pmapper_free_func    = mapper_tac_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_tac_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s\n", argv0, verb);
}

static mapper_t* mapper_tac_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 1) {
		mapper_tac_usage(argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_tac_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
mapper_setup_t mapper_tac_setup = {
	.verb = "tac",
	.pusage_func = mapper_tac_usage,
	.pparse_func = mapper_tac_parse_cli
};
