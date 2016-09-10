#include "cli/mlrcli.h"
#include "containers/sllv.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

typedef struct _mapper_tee_state_t {
	char* output_file_name;
	FILE* output_stream;
	int flush_every_record;
	lrec_writer_t* plrec_writer;
	cli_writer_opts_t* pwriter_opts;
} mapper_tee_state_t;

#define DEFAULT_COUNTER_FIELD_NAME "n"

static void      mapper_tee_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_tee_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_tee_alloc(int do_append, int flush_every_record,
	char* output_file_name, cli_writer_opts_t* pwriter_opts, cli_writer_opts_t* pmain_writer_opts);
static void      mapper_tee_free(mapper_t* pmapper);
static sllv_t*   mapper_tee_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_tee_setup = {
	.verb = "tee",
	.pusage_func = mapper_tee_usage,
	.pparse_func = mapper_tee_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_tee_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* pmain_writer_opts)
{
	int   do_append = FALSE;
	int   flush_every_record = TRUE;
	cli_writer_opts_t* pwriter_opts = mlr_malloc_or_die(sizeof(cli_writer_opts_t));
	cli_writer_opts_init(pwriter_opts);

	int argi = *pargi;
	if ((argc - argi) < 1) {
		mapper_tee_usage(stderr, argv[0], argv[argi]);
		return NULL;
	}
	char* verb = argv[argi++];

	for (; argi < argc; /* variable increment: 1 or 2 depending on flag */) {

		if (argv[argi][0] != '-') {
			break; // No more flag options to process

		} else if (cli_handle_writer_options(argv, argc, &argi, pwriter_opts)) {
			// handled

		} else if (streq(argv[argi], "-a")) {
			do_append = TRUE;
			argi++;

		} else if (streq(argv[argi], "--no-fflush") || streq(argv[argi], "--no-flush")) {
			flush_every_record = FALSE;
			argi++;

		} else {
			mapper_tee_usage(stderr, argv[0], verb);
			return NULL;
		}

	}

	if ((argc - argi) < 1) {
		mapper_tee_usage(stderr, argv[0], verb);
		return NULL;
	}
	char* output_file_name = argv[argi++];

	*pargi = argi;

	mapper_t* pmapper = mapper_tee_alloc(do_append, flush_every_record, output_file_name,
		pwriter_opts, pmain_writer_opts);
	return pmapper;
}

static void mapper_tee_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {filename}\n", argv0, verb);
	fprintf(o, "Passes through input records (like %s cat) but also writes to specified output\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "file, using output-format flags from the command line (e.g. --ocsv). See also\n");
	fprintf(o, "the \"tee\" keyword within %s put, which allows data-dependent filenames.\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "Options:\n");
	fprintf(o, "-a:          append to existing file, if any, rather than overwriting.\n");
	fprintf(o, "--no-fflush: don't call fflush() after every record.\n");
	fprintf(o, "Any of the output-format command-line flags (see %s -h). Example: using\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "  %s --icsv --opprint put '...' then tee --ojson ./mytap.dat then stats1 ...\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "the input is CSV, the output is pretty-print tabular, but the tee-file output\n");
	fprintf(o, "is written in JSON format.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_tee_alloc(int do_append, int flush_every_record,
	char* output_file_name, cli_writer_opts_t* pwriter_opts, cli_writer_opts_t* pmain_writer_opts)
{
	FILE* fp = fopen(output_file_name, do_append ? "a" : "w");
	if (fp == NULL) {
		perror("fopen");
		fprintf(stderr, "%s: fopen error on \"%s\".\n", MLR_GLOBALS.bargv0, output_file_name);
		exit(1);
	}

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_tee_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tee_state_t));
	pstate->output_file_name   = output_file_name;
	pstate->output_stream      = fp;
	pstate->flush_every_record = flush_every_record;
	pstate->pwriter_opts       = pwriter_opts;

	cli_merge_writer_opts(pstate->pwriter_opts, pmain_writer_opts);
	pstate->plrec_writer = lrec_writer_alloc_or_die(pstate->pwriter_opts);

	pmapper->pvstate           = pstate;
	pmapper->pprocess_func     = mapper_tee_process;
	pmapper->pfree_func        = mapper_tee_free;
	return pmapper;
}
static void mapper_tee_free(mapper_t* pmapper) {
	mapper_tee_state_t* pstate = pmapper->pvstate;
	pstate->plrec_writer->pfree_func(pstate->plrec_writer);
	free(pstate->pwriter_opts);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_tee_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tee_state_t* pstate = (mapper_tee_state_t*)pvstate;

	if (pinrec != NULL) {
		// Copy the record since the lrec-writer will free it, and we need the original
		// to return as stream output.
		lrec_t* pcopy = lrec_copy(pinrec);
		pstate->plrec_writer->pprocess_func(pstate->plrec_writer->pvstate, pstate->output_stream, pcopy);
		if (pstate->flush_every_record)
			fflush(pstate->output_stream);
		return sllv_single(pinrec);
	} else {
		pstate->plrec_writer->pprocess_func(pstate->plrec_writer->pvstate, pstate->output_stream, NULL);
		if (fclose(pstate->output_stream) != 0) {
			perror("fclose");
			fprintf(stderr, "%s: fclose error on \"%s\".\n", MLR_GLOBALS.bargv0, pstate->output_file_name);
			exit(1);
		}
		return sllv_single(NULL);
	}
}
