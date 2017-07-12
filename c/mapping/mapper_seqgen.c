#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "containers/sllv.h"
#include "containers/mvfuncs.h"

typedef struct _mapper_seqgen_state_t {
	ap_state_t* pargp;
	char* field_name;
	mv_t start;
	mv_t stop;
	mv_t step;
	int continue_cmp;
} mapper_seqgen_state_t;

#define DEFAULT_FIELD_NAME   "i"
#define DEFAULT_START_STRING "1"
#define DEFAULT_STOP_STRING  "100"
#define DEFAULT_STEP_STRING  "1"

static void      mapper_seqgen_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_seqgen_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_seqgen_alloc(ap_state_t* pargp, char* field_name, mv_t start, mv_t stop, mv_t step);
static void      mapper_seqgen_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_seqgen_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_seqgen_setup = {
	.verb = "seqgen",
	.pusage_func = mapper_seqgen_usage,
	.pparse_func = mapper_seqgen_parse_cli,
	.ignores_input = TRUE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_seqgen_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* field_name = DEFAULT_FIELD_NAME;
	char* start_string = DEFAULT_START_STRING;
	char* stop_string  = DEFAULT_STOP_STRING;
	char* step_string  = DEFAULT_STEP_STRING;

	if ((argc - *pargi) < 1) {
		mapper_seqgen_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "--start", &start_string);
	ap_define_string_flag(pstate, "--stop", &stop_string);
	ap_define_string_flag(pstate, "--step", &step_string);
	ap_define_string_flag(pstate, "-f", &field_name);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_seqgen_usage(stderr, argv[0], verb);
		return NULL;
	}

	mv_t start = mv_scan_number_or_die(start_string);
	mv_t stop  = mv_scan_number_or_die(stop_string);
	mv_t step  = mv_scan_number_or_die(step_string);
	mv_t zero  = mv_from_int(0);

	if (mveq(&step, &zero)) {
		if (!mveq(&start, &stop)) {
			fprintf(stderr, "%s %s: step must not be zero unless start == stop.\n", MLR_GLOBALS.bargv0, verb);
			fprintf(stderr, "Got start=%s, stop=%s, end=%s.\n", start_string, stop_string, step_string);
		}
		return NULL;
	}

	mapper_t* pmapper = mapper_seqgen_alloc(pstate, field_name, start, stop, step);
	return pmapper;
}

static void mapper_seqgen_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Produces a sequence of counters.  Discards the input record stream. Produces\n");
	fprintf(o, "output as specified by the following options:\n");
	fprintf(o, "-f {name} Field name for counters; default \"%s\".\n", DEFAULT_FIELD_NAME);
	fprintf(o, "--start {number} Inclusive start value; default \"%s\".\n", DEFAULT_START_STRING);
	fprintf(o, "--stop  {number} Inclusive stop value; default \"%s\".\n", DEFAULT_STOP_STRING);
	fprintf(o, "--step  {number} Step value; default \"%s\".\n", DEFAULT_STEP_STRING);
	fprintf(o, "Start, stop, and/or step may be floating-point. Output is integer if start,\n");
	fprintf(o, "stop, and step are all integers. Step may be negative. It may not be zero\n");
	fprintf(o, "unless start == stop.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_seqgen_alloc(ap_state_t* pargp, char* field_name, mv_t start, mv_t stop, mv_t step) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_seqgen_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_seqgen_state_t));

	pstate->pargp          = pargp;
	pstate->field_name     = field_name;
	pstate->start          = start;
	pstate->stop           = stop;
	pstate->step           = step;
	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_seqgen_process;
	pmapper->pfree_func    = mapper_seqgen_free;

	return pmapper;
}
static void mapper_seqgen_free(mapper_t* pmapper, context_t* _) {
	mapper_seqgen_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_seqgen_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	// Only produce data at end of input stream. Discard the input stream.
	if (pinrec != NULL) {
		pctx->force_eof = TRUE;
		lrec_free(pinrec);
		return NULL;
	}

	mapper_seqgen_state_t* pstate = pvstate;

	int continue_cmp = 1;
	mv_t zero  = mv_from_int(0);
	if (mv_i_nn_lt(&pstate->step, &zero)) {
		continue_cmp = -1;
	}

	sllv_t* poutrecs = sllv_alloc();
	for (
		mv_t counter = pstate->start;
		mv_nn_comparator(&counter, &pstate->stop) != continue_cmp;
		counter = x_xx_plus_func(&counter, &pstate->step)
	) {
		lrec_t* poutrec = lrec_unbacked_alloc();
		lrec_put(poutrec, pstate->field_name, mv_alloc_format_val(&counter), FREE_ENTRY_VALUE);
		sllv_append(poutrecs, poutrec);
	}
	sllv_append(poutrecs, NULL);

	return poutrecs;
}
