#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmsll.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DEFAULT_OUTPUT_FIELD_NAME "count"

typedef struct _mapper_count_similar_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	lhmslv_t* pcounts_by_group;
	lhmslv_t* precord_lists_by_group;
	char* output_field_name;
} mapper_count_similar_state_t;

static void mapper_count_similar_usage(
	FILE* o,
	char* argv0,
	char* verb);

static mapper_t* mapper_count_similar_parse_cli(
	int* pargi,
	int argc,
	char** argv,
	cli_reader_opts_t* _,
	cli_writer_opts_t* __);

static mapper_t* mapper_count_similar_alloc(
	ap_state_t* pargp,
	slls_t* pgroup_by_field_names,
	char* output_field_name);

static void mapper_count_similar_free(
	mapper_t* pmapper,
	context_t* _);

static sllv_t* mapper_count_similar_process(
	lrec_t* pinrec,
	context_t* pctx,
	void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_count_similar_setup = {
	.verb = "count-similar",
	.pusage_func = mapper_count_similar_usage,
	.pparse_func = mapper_count_similar_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_count_similar_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Ingests all records, then emits each record augmented by a count of \n");
	fprintf(o, "the number of other records having the same group-by field values.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-g {d,e,f} Group-by-field names for counts.\n");
	fprintf(o, "-o {name}  Field name for output count. Default \"%s\".\n",
		DEFAULT_OUTPUT_FIELD_NAME);
}

static mapper_t* mapper_count_similar_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* pgroup_by_field_names = NULL;
	char*   output_field_name = DEFAULT_OUTPUT_FIELD_NAME;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_string_flag(pstate,      "-o", &output_field_name);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_count_similar_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pgroup_by_field_names == NULL) {
		mapper_count_similar_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_count_similar_alloc(pstate, pgroup_by_field_names,
		output_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_count_similar_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	char* output_field_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_count_similar_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_count_similar_state_t));

	pstate->pargp                  = pargp;
	pstate->pgroup_by_field_names  = pgroup_by_field_names;
	pstate->pcounts_by_group       = lhmslv_alloc();
	pstate->precord_lists_by_group = lhmslv_alloc();
	pstate->output_field_name      = output_field_name;

	pmapper->pvstate = pstate;
	pmapper->pprocess_func = mapper_count_similar_process;
	pmapper->pfree_func = mapper_count_similar_free;

	return pmapper;
}

static void mapper_count_similar_free(mapper_t* pmapper, context_t* _) {
	mapper_count_similar_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pgroup_by_field_names);

	// lhmslv_free will free the keys: we only need to free the void-star values.
	for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
		unsigned long long* pcount = pa->pvvalue;
		free(pcount);
	}
	lhmslv_free(pstate->pcounts_by_group);

	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
		sllv_t* precord_list_for_group = pa->pvvalue;
		// outrecs were freed by caller of mapper_tail_process. Here, just free
		// the sllv container itself.
		sllv_free(precord_list_for_group);
	}
	lhmslv_free(pstate->precord_lists_by_group);

	pstate->pgroup_by_field_names = NULL;
	pstate->pcounts_by_group = NULL;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_count_similar_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_count_similar_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) {
			lrec_free(pinrec);
			return NULL;
		}

		unsigned long long* pcount = lhmslv_get(pstate->pcounts_by_group, pgroup_by_field_values);
		if (pcount == NULL) {
			pcount = mlr_malloc_or_die(sizeof(unsigned long long));
			*pcount = 1LL;
			lhmslv_put(pstate->pcounts_by_group, slls_copy(pgroup_by_field_values), pcount, FREE_ENTRY_KEY);
		} else {
			(*pcount)++;
		}

		sllv_t* precord_list_for_group = lhmslv_get(pstate->precord_lists_by_group, pgroup_by_field_values);
		if (precord_list_for_group == NULL) {
			precord_list_for_group = sllv_alloc();
			lhmslv_put(pstate->precord_lists_by_group, slls_copy(pgroup_by_field_values), precord_list_for_group,
				FREE_ENTRY_KEY);

		}
		sllv_append(precord_list_for_group, pinrec);

		slls_free(pgroup_by_field_values);
		return NULL;

	} else {
		sllv_t* poutrecs = sllv_alloc();
		for (lhmslve_t* pa = pstate->precord_lists_by_group->phead; pa != NULL; pa = pa->pnext) {
			unsigned long long* pcount = lhmslv_get(pstate->pcounts_by_group, pa->key);
			if (pcount == NULL) {
				fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
					MLR_GLOBALS.bargv0, __FILE__, __LINE__);
				exit(1);
			}

			sllv_t* precord_list_for_group = pa->pvvalue;
			while (precord_list_for_group->phead != NULL) {
				lrec_t* poutrec = sllv_pop(precord_list_for_group);

				char* scount = mlr_alloc_string_from_ll(*pcount);
				lrec_put(poutrec, pstate->output_field_name, scount, FREE_ENTRY_VALUE);
				sllv_append(poutrecs, poutrec);
			}
		}
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}
