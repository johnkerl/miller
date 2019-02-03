#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"

typedef struct _mapper_cat_state_t {
	ap_state_t* pargp;
	int verbose;
	char* counter_field_name;
	unsigned long long counter;
	slls_t* pgroup_by_field_names;
	lhmslv_t* pcounters_by_group;
} mapper_cat_state_t;

#define DEFAULT_COUNTER_FIELD_NAME "n"

static void      mapper_cat_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_cat_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_cat_alloc(ap_state_t* pargp, int do_counters, int verbose, char* counter_field_name,
	slls_t* pgroup_by_field_names);
static void      mapper_cat_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_cat_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_catn_process_ungrouped(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_catn_process_grouped(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_cat_setup = {
	.verb = "cat",
	.pusage_func = mapper_cat_usage,
	.pparse_func = mapper_cat_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_cat_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* default_counter_field_name = DEFAULT_COUNTER_FIELD_NAME;
	char* counter_field_name = NULL;
	int   do_counters = FALSE;
	int   verbose = FALSE;
	slls_t* pgroup_by_field_names = slls_alloc();

	if ((argc - *pargi) < 1) {
		mapper_cat_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-n",   &do_counters);
	ap_define_true_flag(pstate, "-v",   &verbose);
	ap_define_string_flag(pstate, "-N", &counter_field_name);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_cat_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (counter_field_name != NULL) {
		do_counters = TRUE;
	} else if (do_counters) {
		counter_field_name = default_counter_field_name;
	}

	mapper_t* pmapper = mapper_cat_alloc(pstate, do_counters, verbose, counter_field_name, pgroup_by_field_names);
	return pmapper;
}

static void mapper_cat_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-n        Prepend field \"%s\" to each record with record-counter starting at 1\n",
		DEFAULT_COUNTER_FIELD_NAME);
	fprintf(o, "-g {comma-separated field name(s)} When used with -n/-N, writes record-counters\n");
	fprintf(o, "          keyed by specified field name(s).\n");
	fprintf(o, "-v        Write a low-level record-structure dump to stderr.\n");
	fprintf(o, "-N {name} Prepend field {name} to each record with record-counter starting at 1\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_cat_alloc(ap_state_t* pargp, int do_counters, int verbose, char* counter_field_name,
	slls_t* pgroup_by_field_names)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_cat_state_t* pstate    = mlr_malloc_or_die(sizeof(mapper_cat_state_t));
	pstate->pargp                 = pargp;
	pstate->verbose               = verbose;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->counter_field_name    = counter_field_name;
	pstate->counter               = 0LL;
	pstate->pcounters_by_group    = lhmslv_alloc();
	pmapper->pvstate              = pstate;

	pmapper->pprocess_func = NULL;
	if (do_counters) {
		if (pgroup_by_field_names->length == 0) {
			pmapper->pprocess_func = mapper_catn_process_ungrouped;
		} else {
			pmapper->pprocess_func = mapper_catn_process_grouped;
		}
	} else {
		pmapper->pprocess_func = mapper_cat_process;
	}

	pmapper->pfree_func           = mapper_cat_free;
	return pmapper;
}
static void mapper_cat_free(mapper_t* pmapper, context_t* _) {
	mapper_cat_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pgroup_by_field_names);
	for (lhmslve_t* pe = pstate->pcounters_by_group->phead; pe != NULL; pe = pe->pnext) {
		free(pe->pvvalue);
	}
	lhmslv_free(pstate->pcounters_by_group);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_cat_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_cat_state_t* pstate = (mapper_cat_state_t*)pvstate;
	if (pstate->verbose) {
		lrec_dump_fp(pinrec, stderr);
	}
	if (pinrec != NULL)
		return sllv_single(pinrec);
	else
		return sllv_single(NULL);
}

// ----------------------------------------------------------------
static sllv_t* mapper_catn_process_ungrouped(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_cat_state_t* pstate = (mapper_cat_state_t*)pvstate;
	if (pstate->verbose) {
		lrec_dump_fp(pinrec, stderr);
	}
	if (pinrec != NULL) {
		char* counter_field_value = mlr_alloc_string_from_ull(++pstate->counter);
		lrec_prepend(pinrec, pstate->counter_field_name, counter_field_value, FREE_ENTRY_VALUE);
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_catn_process_grouped(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_cat_state_t* pstate = (mapper_cat_state_t*)pvstate;
	if (pstate->verbose) {
		lrec_dump_fp(pinrec, stderr);
	}
	if (pinrec != NULL) {

		unsigned long long counter = 0LL;

		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values == NULL) { // Treat as unkeyed
			counter = ++pstate->counter;
		} else {
			unsigned long long* pcount_for_group = lhmslv_get(pstate->pcounters_by_group,
				pgroup_by_field_values);
			if (pcount_for_group == NULL) {
				pcount_for_group = mlr_malloc_or_die(sizeof(unsigned long long));
				*pcount_for_group = 0LL;
				lhmslv_put(pstate->pcounters_by_group, slls_copy(pgroup_by_field_values),
					pcount_for_group, FREE_ENTRY_KEY);
			}
			slls_free(pgroup_by_field_values);
			(*pcount_for_group)++;
			counter = *pcount_for_group;
		}
		char* counter_field_value = mlr_alloc_string_from_ull(counter);
		lrec_prepend(pinrec, pstate->counter_field_name, counter_field_value, FREE_ENTRY_VALUE);
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}
