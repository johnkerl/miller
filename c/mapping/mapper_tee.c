#include "cli/argparse.h"
#include "containers/sllv.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

typedef struct _mapper_tee_state_t {
	ap_state_t* pargp;
	char* output_file_name;
	FILE* output_stream;
	int flush_every_record;
	lrec_writer_t* plrec_writer;
} mapper_tee_state_t;

#define DEFAULT_COUNTER_FIELD_NAME "n"

static void      mapper_tee_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_tee_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_tee_alloc(ap_state_t* pargp, int do_append, int flush_every_record,
	char* output_file_name);
static void      mapper_tee_free(mapper_t* pmapper);
static sllv_t*   mapper_tee_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_tee_setup = {
	.verb = "tee",
	.pusage_func = mapper_tee_usage,
	.pparse_func = mapper_tee_parse_cli
};

// ----------------------------------------------------------------
static mapper_t* mapper_tee_parse_cli(int* pargi, int argc, char** argv) {
	int   do_append = FALSE;
	int   flush_every_record = FALSE;

	if ((argc - *pargi) < 1) {
		mapper_tee_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-a",   &do_append);
	ap_define_true_flag(pstate, "-f",   &flush_every_record);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_tee_usage(stderr, argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_tee_usage(stderr, argv[0], verb);
		return NULL;
	}
	char* output_file_name = argv[*pargi];
	*pargi += 1;

	mapper_t* pmapper = mapper_tee_alloc(pstate, do_append, flush_every_record, output_file_name);
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
}

// ----------------------------------------------------------------
static mapper_t* mapper_tee_alloc(ap_state_t* pargp, int do_append, int flush_every_record,
	char* output_file_name)
{
	FILE* fp = fopen(output_file_name, do_append ? "a" : "w");
	if (fp == NULL) {
		perror("fopen");
		fprintf(stderr, "%s: fopen error on \"%s\".\n", MLR_GLOBALS.bargv0, output_file_name);
		exit(1);
	}

	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_tee_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tee_state_t));
	pstate->pargp              = pargp;
	pstate->output_file_name   = output_file_name;
	pstate->output_stream      = fp;
	pstate->flush_every_record = flush_every_record;
	pstate->plrec_writer       = NULL;
	pmapper->pvstate           = pstate;
	pmapper->pprocess_func     = mapper_tee_process;
	pmapper->pfree_func        = mapper_tee_free;
	return pmapper;
}
static void mapper_tee_free(mapper_t* pmapper) {
	mapper_tee_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_tee_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tee_state_t* pstate = (mapper_tee_state_t*)pvstate;

	// mapper_tee_alloc is called from the CLI-parser and cli_opts isn't finalized until that
	// returns. So we cannot do this in mapper_tee_alloc.
	if (pstate->plrec_writer == NULL) {
		cli_opts_t* popts = MLR_GLOBALS.popts;
		pstate->plrec_writer = lrec_writer_alloc(popts->ofile_fmt, popts->ors, popts->ofs, popts->ops,
			popts->headerless_csv_output, popts->oquoting, popts->left_align_pprint,
			popts->right_justify_xtab_value, popts->json_flatten_separator, popts->quote_json_values_always,
			popts->stack_json_output_vertically, popts->wrap_json_output_in_outer_list);
		if (pstate->plrec_writer == NULL) {
			fprintf(stderr, "%s: internal coding error detected in file \"%s\" at line %d.\n",
				MLR_GLOBALS.bargv0, __FILE__, __LINE__);
			exit(1);
		}
	}

	if (pinrec != NULL) {
		// Copy the record since the lrec-writer will free it, and we need the original
		// to return as stream output.
		lrec_t* pcopy = lrec_copy(pinrec);
		pstate->plrec_writer->pprocess_func(pstate->output_stream, pcopy, pstate->plrec_writer->pvstate);
		if (pstate->flush_every_record)
			fflush(pstate->output_stream);
		return sllv_single(pinrec);
	} else {
		pstate->plrec_writer->pprocess_func(pstate->output_stream, NULL, pstate->plrec_writer->pvstate);
		if (fclose(pstate->output_stream) != 0) {
			perror("fclose");
			fprintf(stderr, "%s: fclose error on \"%s\".\n", MLR_GLOBALS.bargv0, pstate->output_file_name);
			exit(1);
		}
		return sllv_single(NULL);
	}
}
