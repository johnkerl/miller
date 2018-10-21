#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"

typedef struct _mapper_skip_trivial_records_state_t {
	ap_state_t* pargp;
} mapper_skip_trivial_records_state_t;

static void      mapper_skip_trivial_records_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_skip_trivial_records_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_skip_trivial_records_alloc(ap_state_t* pargp);
static void      mapper_skip_trivial_records_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_skip_trivial_records_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_skip_trivial_records_setup = {
	.verb = "skip-trivial-records",
	.pusage_func = mapper_skip_trivial_records_usage,
	.pparse_func = mapper_skip_trivial_records_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_skip_trivial_records_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_skip_trivial_records_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_skip_trivial_records_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_skip_trivial_records_alloc(pstate);
	return pmapper;
}

static void mapper_skip_trivial_records_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Passes through all records except:\n");
	fprintf(o, "* Those with zero fields\n");
	fprintf(o, "* Those for which all fields have empty value\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_skip_trivial_records_alloc(ap_state_t* pargp)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_skip_trivial_records_state_t* pstate    = mlr_malloc_or_die(sizeof(mapper_skip_trivial_records_state_t));
	pstate->pargp                 = pargp;
	pmapper->pvstate              = pstate;

	pmapper->pprocess_func = NULL;
	pmapper->pprocess_func = mapper_skip_trivial_records_process;
	pmapper->pfree_func    = mapper_skip_trivial_records_free;
	return pmapper;
}
static void mapper_skip_trivial_records_free(mapper_t* pmapper, context_t* _) {
	mapper_skip_trivial_records_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_skip_trivial_records_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) { // not end of record-stream
		int has_any = FALSE;
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			if (*pe->value != 0) {
				has_any = TRUE;
				break;
			}
		}
		if (has_any) {
			return sllv_single(pinrec);
		} else {
			lrec_free(pinrec);
			return NULL;
		}
	} else { // end of record-stream
		return sllv_single(NULL);
	}
}
