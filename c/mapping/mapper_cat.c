#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"

typedef struct _mapper_cat_state_t {
	ap_state_t* pargp;
	char* counter_field_name;
	unsigned long long counter;
} mapper_cat_state_t;

#define DEFAULT_COUNTER_FIELD_NAME "n"

static void      mapper_cat_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_cat_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_cat_alloc(ap_state_t* pargp, int do_counters, char* counter_field_name);
static void      mapper_cat_free(mapper_t* pmapper);
static sllv_t*   mapper_cat_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_catn_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_cat_setup = {
	.verb = "cat",
	.pusage_func = mapper_cat_usage,
	.pparse_func = mapper_cat_parse_cli
};

// ----------------------------------------------------------------
static mapper_t* mapper_cat_parse_cli(int* pargi, int argc, char** argv) {
	char* default_counter_field_name = DEFAULT_COUNTER_FIELD_NAME;
	char* counter_field_name = NULL;
	int   do_counters = FALSE;

	if ((argc - *pargi) < 1) {
		mapper_cat_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-n",   &do_counters);
	ap_define_string_flag(pstate, "-N", &counter_field_name);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_cat_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (counter_field_name != NULL) {
		do_counters = TRUE;
	} else if (do_counters) {
		counter_field_name = default_counter_field_name;
	}

	mapper_t* pmapper = mapper_cat_alloc(pstate, do_counters, counter_field_name);
	return pmapper;
}

static void mapper_cat_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-n        Prepend field \"%s\" to each record with record-counter starting at 1\n",
		DEFAULT_COUNTER_FIELD_NAME);
	fprintf(o, "-N {name} Prepend field {name} to each record with record-counter starting at 1\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_cat_alloc(ap_state_t* pargp, int do_counters, char* counter_field_name) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_cat_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_cat_state_t));
	pstate->pargp              = pargp;
	pstate->counter_field_name = counter_field_name;
	pstate->counter            = 0LL;
	pmapper->pvstate           = pstate;
	pmapper->pprocess_func     = do_counters ? mapper_catn_process : mapper_cat_process;
	pmapper->pfree_func        = mapper_cat_free;
	return pmapper;
}
static void mapper_cat_free(mapper_t* pmapper) {
	mapper_cat_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_cat_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL)
		return sllv_single(pinrec);
	else
		return sllv_single(NULL);
}

// ----------------------------------------------------------------
static sllv_t* mapper_catn_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_cat_state_t* pstate = (mapper_cat_state_t*)pvstate;
	if (pinrec != NULL) {
		char* counter_field_value = mlr_alloc_string_from_ull(pstate->counter++);
		lrec_prepend(pinrec, pstate->counter_field_name, counter_field_value, FREE_ENTRY_VALUE);
		return sllv_single(pinrec);
	} else {
		return sllv_single(NULL);
	}
}
