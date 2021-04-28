#include <stdio.h>
#include "cli/argparse.h"
#include "lib/mlrutil.h"
#include "containers/lhmsi.h"
#include "mapping/mappers.h"

#define DEFAULT_FILL_STRING "N/A"

typedef struct _mapper_fill_empty_state_t {
	char* fill_string;
} mapper_fill_empty_state_t;

static void      mapper_fill_empty_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_fill_empty_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_fill_empty_alloc();
static void      mapper_fill_empty_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_fill_empty_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_fill_empty_setup = {
	.verb = "fill-empty",
	.pusage_func = mapper_fill_empty_usage,
	.pparse_func = mapper_fill_empty_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_fill_empty_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Fills empty-string fields with specified fill-value.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-v {string} Fill-value: defaults to \"%s\"\n", DEFAULT_FILL_STRING);
}

static mapper_t* mapper_fill_empty_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* fill_string = DEFAULT_FILL_STRING;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "-v", &fill_string);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_fill_empty_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_fill_empty_alloc(fill_string);
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_fill_empty_alloc(char* fill_string) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_fill_empty_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_fill_empty_state_t));
	pstate->fill_string = fill_string;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_fill_empty_process;
	pmapper->pfree_func    = mapper_fill_empty_free;

	return pmapper;
}

static void mapper_fill_empty_free(mapper_t* pmapper, context_t* _) {
	mapper_fill_empty_state_t* pstate = pmapper->pvstate;
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_fill_empty_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_fill_empty_state_t* pstate = pvstate;

    if (pinrec == NULL) { // End of input stream: emit null.
        return sllv_single(NULL);
    }

	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		if (pe->value[0] == 0) {
			lrece_update_value(pe, pstate->fill_string, NO_FREE);
		}
	}

	return sllv_single(pinrec);
}
