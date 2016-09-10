#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"

static void      mapper_check_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_check_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_check_alloc();
static void      mapper_check_free(mapper_t* pmapper);
static sllv_t*   mapper_check_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_check_setup = {
	.verb = "check",
	.pusage_func = mapper_check_usage,
	.pparse_func = mapper_check_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_check_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Consumes records without printing any output.\n");
	fprintf(o, "Useful for doing a well-formatted check on input data.\n");
}

static mapper_t* mapper_check_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_check_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_check_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_check_alloc() {
	mapper_t* pmapper      = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate       = NULL;
	pmapper->pprocess_func = mapper_check_process;
	pmapper->pfree_func    = mapper_check_free;
	return pmapper;
}
static void mapper_check_free(mapper_t* pmapper) {
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_check_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	lrec_free(pinrec);
	return sllv_single(NULL);
}
