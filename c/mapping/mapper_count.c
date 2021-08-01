#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmsll.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/lhmsll.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DEFAULT_OUTPUT_FIELD_NAME "count"

typedef struct _mapper_count_state_t {
	ap_state_t* pargp;
	slls_t*   pgroup_by_field_names;
	long long ungrouped_count;
	lhmslv_t* pcounts_by_group;
	int show_counts_only;
	char* output_field_name;
} mapper_count_state_t;

// ----------------------------------------------------------------
static void mapper_count_usage(
	FILE* o,
	char* argv0,
	char* verb);

static mapper_t* mapper_count_parse_cli(
	int* pargi,
	int argc,
	char** argv,
	cli_reader_opts_t* _,
	cli_writer_opts_t* __);

static mapper_t* mapper_count_alloc(
	ap_state_t* pargp,
	slls_t* pgroup_by_field_names,
	int show_counts_only,
	char* output_field_name);

static void mapper_count_free(
	mapper_t* pmapper,
	context_t* _);

static sllv_t* mapper_count_process_ungrouped(
	lrec_t* pinrec,
	context_t* pctx,
	void* pvstate);

static sllv_t* mapper_count_process_grouped(
	lrec_t* pinrec,
	context_t* pctx,
	void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_count_setup = {
	.verb = "count",
	.pusage_func = mapper_count_usage,
	.pparse_func = mapper_count_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_count_usage(
	FILE* o,
	char* argv0,
	char* verb)
{
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Prints number of records, optionally grouped by distinct values for specified field names.\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-g {a,b,c}    Field names for distinct count.\n");
	fprintf(o, "-n            Show only the number of distinct values. Not interesting without -g.\n");
	fprintf(o, "-o {name}     Field name for output count. Default \"%s\".\n", DEFAULT_OUTPUT_FIELD_NAME);
}

// ----------------------------------------------------------------
static mapper_t* mapper_count_parse_cli(
	int* pargi,
	int argc,
	char** argv,
	cli_reader_opts_t* _,
	cli_writer_opts_t* __)
{
	slls_t* pgroup_by_field_names = NULL;
	char*   output_field_name = DEFAULT_OUTPUT_FIELD_NAME;
	int     show_counts_only = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);
	ap_define_string_flag(pstate,      "-o", &output_field_name);
	ap_define_true_flag(pstate,        "-n", &show_counts_only);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_count_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_count_alloc(pstate, pgroup_by_field_names, show_counts_only,
		output_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_count_alloc(
	ap_state_t* pargp,
	slls_t* pgroup_by_field_names,
	int show_counts_only,
	char* output_field_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_count_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_count_state_t));

	pstate->pargp                 = pargp;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->output_field_name     = output_field_name;
	pstate->ungrouped_count       = 0LL;
	pstate->pcounts_by_group      = lhmslv_alloc();
	pstate->show_counts_only      = show_counts_only;

	pmapper->pvstate = pstate;

	if (pgroup_by_field_names == NULL) {
		pmapper->pprocess_func = mapper_count_process_ungrouped;
	} else{
		pmapper->pprocess_func = mapper_count_process_grouped;
	}

	pmapper->pfree_func = mapper_count_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_count_free(
	mapper_t* pmapper,
	context_t* _)
{
	mapper_count_state_t* pstate = pmapper->pvstate;

	slls_free(pstate->pgroup_by_field_names);
	pstate->pgroup_by_field_names = NULL;

	// lhmslv_free will free the keys: we only need to free the void-star values.
	if (pstate->pcounts_by_group != NULL) {
		for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
			unsigned long long* pcount = pa->pvvalue;
			free(pcount);
		}
		lhmslv_free(pstate->pcounts_by_group);
		pstate->pcounts_by_group = NULL;
	}

	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_count_process_ungrouped(
	lrec_t* pinrec,
	context_t* pctx,
	void* pvstate)
{
	mapper_count_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		pstate->ungrouped_count++;
		lrec_free(pinrec);
		return NULL;
	} else { // end of record stream
		lrec_t* poutrec = lrec_unbacked_alloc();
		lrec_put(poutrec, pstate->output_field_name,
			mlr_alloc_string_from_ll(pstate->ungrouped_count), FREE_ENTRY_VALUE);
		return sllv_single(poutrec);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_count_process_grouped(
	lrec_t* pinrec,
	context_t* pctx,
	void* pvstate)
{
	mapper_count_state_t* pstate = pvstate;
	if (pinrec != NULL) { // not end of record stream
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
			slls_t* pcopy = slls_copy(pgroup_by_field_values);
			lhmslv_put(pstate->pcounts_by_group, pcopy, pcount, FREE_ENTRY_KEY);
			lrec_free(pinrec);
		} else {
			(*pcount)++;
			lrec_free(pinrec);
		}

		return NULL;

	} else { // end of record stream
		sllv_t* poutrecs = sllv_alloc();

		if (pstate->show_counts_only) {
			lrec_t* poutrec = lrec_unbacked_alloc();

			unsigned long long count = (unsigned long long)lhmslv_size(pstate->pcounts_by_group);
			lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(count),
				FREE_ENTRY_VALUE);

			sllv_append(poutrecs, poutrec);
		} else {
			for (lhmslve_t* pa = pstate->pcounts_by_group->phead; pa != NULL; pa = pa->pnext) {
				lrec_t* poutrec = lrec_unbacked_alloc();

				slls_t* pgroup_by_field_values = pa->key;

				sllse_t* pb = pstate->pgroup_by_field_names->phead;
				sllse_t* pc =         pgroup_by_field_values->phead;
				for ( ; pb != NULL && pc != NULL; pb = pb->pnext, pc = pc->pnext) {
					lrec_put(poutrec, pb->value, pc->value, NO_FREE);
				}

				unsigned long long* pcount = pa->pvvalue;
				lrec_put(poutrec, pstate->output_field_name, mlr_alloc_string_from_ull(*pcount),
					FREE_ENTRY_VALUE);

				sllv_append(poutrecs, poutrec);
			}
		}

		sllv_append(poutrecs, NULL);
		return poutrecs;

	}
}
