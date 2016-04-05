#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef enum _repeat_type_t {
	BY_COUNT,
	BY_FIELD_NAME,
} repeat_type_t;

typedef struct _mapper_repeat_state_t {
	ap_state_t* pargp;
	char* repeat_count_field_name;
	long long repeat_count;
} mapper_repeat_state_t;

static void      mapper_repeat_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_repeat_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_repeat_alloc(ap_state_t* pargp, long long repeat_count, char* repeat_count_field_name);
static void      mapper_repeat_free(mapper_t* pmapper);
static sllv_t*   mapper_repeat_process_nop(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_repeat_process_by_positive_count(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_repeat_process_by_field_name(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_repeat_setup = {
	.verb = "repeat",
	.pusage_func = mapper_repeat_usage,
	.pparse_func = mapper_repeat_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_repeat_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Copies input records to output records multiple times.\n");
	fprintf(o, "Options must be exactly one of the following:\n");
	fprintf(o, "  -n {repeat count}  Repeat each input record this many times.\n");
	fprintf(o, "  -f {field name}    Same, but take the repeat count from the specified\n");
	fprintf(o, "                     field name of each input record.\n");
	fprintf(o, "Example:\n");
	fprintf(o, "  echo x=0 | %s %s -n 4 then put '$x=urand()'\n", argv0, verb);
	fprintf(o, "produces:\n");
	fprintf(o, " x=0.488189\n");
	fprintf(o, " x=0.484973\n");
	fprintf(o, " x=0.704983\n");
	fprintf(o, " x=0.147311\n");
	fprintf(o, "Example:\n");
	fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f b\n", argv0, verb);
	fprintf(o, "produces:\n");
	fprintf(o, "  a=1,b=2,c=3\n");
	fprintf(o, "  a=1,b=2,c=3\n");
	fprintf(o, "Example:\n");
	fprintf(o, "  echo a=1,b=2,c=3 | %s %s -f c\n", argv0, verb);
	fprintf(o, "produces:\n");
	fprintf(o, "  a=1,b=2,c=3\n");
	fprintf(o, "  a=1,b=2,c=3\n");
	fprintf(o, "  a=1,b=2,c=3\n");
}

static const long long UNINIT_REPEAT_COUNT = -123457689LL;
static mapper_t* mapper_repeat_parse_cli(int* pargi, int argc, char** argv) {
	long long repeat_count = UNINIT_REPEAT_COUNT;
	char* repeat_count_field_name = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "-f", &repeat_count_field_name);
	ap_define_long_long_flag(pstate, "-n", &repeat_count);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_repeat_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (repeat_count == UNINIT_REPEAT_COUNT && repeat_count_field_name == NULL) {
		mapper_repeat_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (repeat_count != UNINIT_REPEAT_COUNT && repeat_count_field_name != NULL) {
		mapper_repeat_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_repeat_alloc(pstate, repeat_count, repeat_count_field_name);
}

// ----------------------------------------------------------------
static mapper_t* mapper_repeat_alloc(ap_state_t* pargp, long long repeat_count, char* repeat_count_field_name) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_repeat_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_repeat_state_t));

	pstate->pargp           = pargp;
	pstate->repeat_count    = repeat_count;
	pstate->repeat_count_field_name = repeat_count_field_name;

	pmapper->pvstate        = pstate;

	if (repeat_count_field_name != NULL)
		pmapper->pprocess_func  = mapper_repeat_process_by_field_name;
	else if (repeat_count >= 1LL)
		pmapper->pprocess_func  = mapper_repeat_process_by_positive_count;
	else
		pmapper->pprocess_func  = mapper_repeat_process_nop;

	pmapper->pfree_func     = mapper_repeat_free;

	return pmapper;
}

static void mapper_repeat_free(mapper_t* pmapper) {
	mapper_repeat_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_repeat_process_nop(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of record stream
		return sllv_single(NULL);
	else
		return NULL;
}

// ----------------------------------------------------------------
static sllv_t* mapper_repeat_process_by_positive_count(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_repeat_state_t* pstate = pvstate;

	if (pinrec == NULL) // End of record stream
		return sllv_single(NULL);

	sllv_t* poutrecs = sllv_alloc();
	// The input record can be baton-passed once to the output without multiple frees.
	// After that, we must copy it.
	sllv_append(poutrecs, pinrec);
	for (long long i = 1; i < pstate->repeat_count; i++) {
		sllv_append(poutrecs, lrec_copy(pinrec));
	}
	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_repeat_process_by_field_name(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_repeat_state_t* pstate = pvstate;

	if (pinrec == NULL) // End of record stream
		return sllv_single(NULL);

	char* repeat_count_field_svalue = lrec_get(pinrec, pstate->repeat_count_field_name);
	if (repeat_count_field_svalue == NULL) {
		lrec_free(pinrec);
		return NULL;
	}

	long long repeat_count;
	int ok = mlr_try_int_from_string(repeat_count_field_svalue, &repeat_count);
	if (!ok || repeat_count <= 0) {
		lrec_free(pinrec);
		return NULL;
	}

	sllv_t* poutrecs = sllv_alloc();
	// The input record can be baton-passed once to the output without multiple frees.
	// After that, we must copy it.
	sllv_append(poutrecs, pinrec);
	for (long long i = 1; i < repeat_count; i++) {
		sllv_append(poutrecs, lrec_copy(pinrec));
	}
	return poutrecs;
}
