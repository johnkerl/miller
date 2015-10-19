#include "lib/mlrutil.h"
#include "containers/lhmslv.h"
#include "containers/sllv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

typedef struct _mapper_regularize_state_t {
	lhmslv_t* psorted_to_original;
} mapper_regularize_state_t;

static sllv_t*   mapper_regularize_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_regularize_free(void* pvstate);
static mapper_t* mapper_regularize_alloc();
static void      mapper_regularize_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_regularize_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_regularize_setup = {
	.verb = "regularize",
	.pusage_func = mapper_regularize_usage,
	.pparse_func = mapper_regularize_parse_cli
};

// ----------------------------------------------------------------
static void mapper_regularize_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "For records seen earlier in the data stream with same field names in\n");
	fprintf(o, "a different order, outputs them with field names in the previously\n");
	fprintf(o, "encountered order.\n");
	fprintf(o, "Example: input records a=1,c=2,b=3, then e=4,d=5, then c=7,a=6,b=8\n");
	fprintf(o, "output as              a=1,c=2,b=3, then e=4,d=5, then a=6,c=7,b=8\n");
}

static mapper_t* mapper_regularize_parse_cli(int* pargi, int argc, char** argv) {
    *pargi += 1;
	return mapper_regularize_alloc();
}

// ----------------------------------------------------------------
static mapper_t* mapper_regularize_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_regularize_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_regularize_state_t));
	pstate->psorted_to_original = lhmslv_alloc();

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_regularize_process;
	pmapper->pfree_func    = mapper_regularize_free;

	return pmapper;
}

static void mapper_regularize_free(void* pvstate) {
	mapper_regularize_state_t* pstate = (mapper_regularize_state_t*)pvstate;
	lhmslv_free(pstate->psorted_to_original);
}

// ----------------------------------------------------------------
static sllv_t* mapper_regularize_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_regularize_state_t* pstate = (mapper_regularize_state_t*)pvstate;
		slls_t* current_sorted_field_names = mlr_reference_keys_from_record(pinrec);
		slls_sort(current_sorted_field_names);
		slls_t* previous_sorted_field_names = lhmslv_get(pstate->psorted_to_original, current_sorted_field_names);
		if (previous_sorted_field_names == NULL) {
			previous_sorted_field_names = slls_copy(current_sorted_field_names);
			lhmslv_put(pstate->psorted_to_original, previous_sorted_field_names, mlr_copy_keys_from_record(pinrec));
			return sllv_single(pinrec);
		} else {
			lrec_t* poutrec = lrec_unbacked_alloc();
			for (sllse_t* pe = previous_sorted_field_names->phead; pe != NULL; pe = pe->pnext) {
				lrec_put(poutrec, pe->value, mlr_strdup_or_die(lrec_get(pinrec, pe->value)), LREC_FREE_ENTRY_VALUE);
			}
			lrec_free(pinrec);
			return sllv_single(poutrec);
		}
	}
	else {
		return sllv_single(NULL);
	}
}
