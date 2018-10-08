#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmss.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_fill_down_state_t {
	ap_state_t* pargp;
	slls_t* pfill_down_field_names;
	lhmss_t* plast_non_null_values;
} mapper_fill_down_state_t;

static void      mapper_fill_down_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_fill_down_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_fill_down_alloc(ap_state_t* pargp, slls_t* pfill_down_field_names);
static void      mapper_fill_down_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_fill_down_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_fill_down_setup = {
	.verb = "fill-down",
	.pusage_func = mapper_fill_down_usage,
	.pparse_func = mapper_fill_down_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_fill_down_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}    Field names for fill-down\n");
	fprintf(o, "xxx description goes here.\n");
	// xxx null vs absent behavior
}

static mapper_t* mapper_fill_down_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pfill_down_field_names = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfill_down_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_fill_down_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (pfill_down_field_names == NULL) {
		mapper_fill_down_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_fill_down_alloc(pstate, pfill_down_field_names);
}

// ----------------------------------------------------------------
static mapper_t* mapper_fill_down_alloc(ap_state_t* pargp, slls_t* pfill_down_field_names) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_fill_down_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_fill_down_state_t));

	pstate->pargp                  = pargp;
	pstate->pfill_down_field_names = pfill_down_field_names;
	pstate->plast_non_null_values  = lhmss_alloc();

	pmapper->pvstate        = pstate;
	pmapper->pprocess_func  = mapper_fill_down_process;
	pmapper->pfree_func     = mapper_fill_down_free;

	return pmapper;
}

static void mapper_fill_down_free(mapper_t* pmapper, context_t* _) {
	mapper_fill_down_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pfill_down_field_names);
	lhmss_free(pstate->plast_non_null_values);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_fill_down_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_fill_down_state_t* pstate = pvstate;
	if (pinrec == NULL) { // end of record stream
		return sllv_single(NULL);
	}

	for (sllse_t* pe = pstate->pfill_down_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* pkey = pe->value;
		char* pvalue = lrec_get(pinrec, pkey);
		if (pvalue != NULL) {
			// xxx
			lhmss_put(pstate->plast_non_null_values,
				mlr_strdup_or_die(pkey),
				mlr_strdup_or_die(pvalue),
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		} else {
			// xxx
			lhmss_put(pstate->plast_non_null_values,
				mlr_strdup_or_die(pkey),
				mlr_strdup_or_die(pvalue),
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}
	}

	return sllv_single(pinrec);
}
