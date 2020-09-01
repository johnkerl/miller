#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"

static void      mapper_nothing_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_nothing_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_nothing_alloc();
static void      mapper_nothing_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_nothing_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_nothing_setup = {
	.verb = "nothing",
	.pusage_func = mapper_nothing_usage,
	.pparse_func = mapper_nothing_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_nothing_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_nothing_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	*pargi += 1;
	mapper_t* pmapper = mapper_nothing_alloc();
	return pmapper;
}

static void mapper_nothing_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Drops all input records. Useful for testing, or after tee/print/etc. have\n");
	fprintf(o, "produced other output.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_nothing_alloc(ap_state_t* pargp, int do_counters, char* counter_field_name) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	pmapper->pvstate       = NULL;
	pmapper->pprocess_func = mapper_nothing_process;
	pmapper->pfree_func    = mapper_nothing_free;
	return pmapper;
}
static void mapper_nothing_free(mapper_t* pmapper, context_t* _) {
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_nothing_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		lrec_free(pinrec);
		return NULL;
	} else {
		return sllv_single(NULL);
	}
}
