#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"

// ----------------------------------------------------------------
static sllv_t* mapper_check_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	return sllv_single(NULL);
}

// ----------------------------------------------------------------
static void mapper_check_free(void* pvstate) {
}

static mapper_t* mapper_check_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate              = NULL;
	pmapper->pprocess_func = mapper_check_process;
	pmapper->pfree_func    = mapper_check_free;
	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_check_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s\n", argv0, verb);
	fprintf(stdout, "Consumes records without printing any output.\n");
	fprintf(stdout, "Useful for doing a well-formatted check on input data.\n");
}

static mapper_t* mapper_check_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 1) {
		mapper_check_usage(argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_check_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
mapper_setup_t mapper_check_setup = {
	.verb = "check",
	.pusage_func = mapper_check_usage,
	.pparse_func = mapper_check_parse_cli
};
