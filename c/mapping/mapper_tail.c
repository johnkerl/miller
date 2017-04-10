#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_tail_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	unsigned long long tail_count;
	lhmslv_t* precord_lists_by_group;
} mapper_tail_state_t;

static void      mapper_tail_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_tail_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_tail_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, unsigned long long tail_count);
static void      mapper_tail_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_tail_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_tail_setup = {
	.verb = "tail",
	.pusage_func = mapper_tail_usage,
	.pparse_func = mapper_tail_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_tail_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-n {count}    Tail count to print; default 10\n");
	fprintf(o, "-g {a,b,c}    Optional group-by-field names for tail counts\n");
	fprintf(o, "Passes through the last n records, optionally by category.\n");
}

static mapper_t* mapper_tail_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int     tail_count            = 10;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &tail_count);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_tail_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_tail_alloc(pstate, pgroup_by_field_names, tail_count);
}

// ----------------------------------------------------------------
static mapper_t* mapper_tail_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names, unsigned long long tail_count) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_tail_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_tail_state_t));

	pstate->pargp                  = pargp;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->tail_count             = tail_count;
	pstate->precord_lists_by_group = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_tail_process;
	pmapper->pfree_func    = mapper_tail_free;

	return pmapper;
}

static void mapper_tail_free(mapper_t* pmapper, context_t* _) {
	mapper_tail_state_t* pstate = pmapper->pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
		sllv_t* precord_list_for_group = pa->pvvalue;
		// outrecs were freed by caller of mapper_tail_process. Here, just free
		// the sllv container itself.
		sllv_free(precord_list_for_group);
	}
	lhmslv_free(pstate->precord_lists_by_group);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_tail_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_tail_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values != NULL) {
			sllv_t* precord_list_for_group = lhmslv_get(pstate->precord_lists_by_group, pgroup_by_field_values);
			if (precord_list_for_group == NULL) {
				precord_list_for_group = sllv_alloc();
				lhmslv_put(pstate->precord_lists_by_group, slls_copy(pgroup_by_field_values),
					precord_list_for_group, FREE_ENTRY_KEY);
			}
			slls_free(pgroup_by_field_values);
			if (precord_list_for_group->length >= pstate->tail_count) {
				lrec_t* porec = sllv_pop(precord_list_for_group);
				lrec_free(porec);
			}
			sllv_append(precord_list_for_group, pinrec);
		} else {
			lrec_free(pinrec);
		}
		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
			sllv_t* precord_list_for_group = pa->pvvalue;
			sllv_transfer(poutrecs, precord_list_for_group);
		}
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}
