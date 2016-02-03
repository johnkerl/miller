#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ----------------------------------------------------------------
typedef struct _sample_bucket_t {
	int nalloc;
	int nused;
	lrec_t** plrecs;
} sample_bucket_t;

sample_bucket_t* sample_bucket_alloc(int nalloc);
void             sample_bucket_free(sample_bucket_t* pbucket);
void             sample_bucket_handle(sample_bucket_t* pbucket, lrec_t* prec, int record_number);

// ----------------------------------------------------------------
typedef struct _mapper_sample_state_t {
	ap_state_t* pargp;
	slls_t* pgroup_by_field_names;
	unsigned long long sample_count;
	lhmslv_t* pbuckets_by_group;
} mapper_sample_state_t;

static void      mapper_sample_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_sample_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_sample_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	unsigned long long sample_count);
static void      mapper_sample_free(mapper_t* pmapper);
static sllv_t*   mapper_sample_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_sample_setup = {
	.verb = "sample",
	.pusage_func = mapper_sample_usage,
	.pparse_func = mapper_sample_parse_cli
};

// ----------------------------------------------------------------
static void mapper_sample_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Reservoir sampling (subsampling without replacement), optionally by category.\n");
	fprintf(o, "-k {count}    Required: number of records to output, total, or by group if using -g.\n");
	fprintf(o, "-g {a,b,c}    Optional: group-by-field names for samples.\n");
}

static mapper_t* mapper_sample_parse_cli(int* pargi, int argc, char** argv) {
	int     sample_count          = -1;
	slls_t* pgroup_by_field_names = slls_alloc();

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-k", &sample_count);
	ap_define_string_list_flag(pstate, "-g", &pgroup_by_field_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_sample_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (sample_count == -1) {
		mapper_sample_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_sample_alloc(pstate, pgroup_by_field_names, sample_count);
}

// ----------------------------------------------------------------
static mapper_t* mapper_sample_alloc(ap_state_t* pargp, slls_t* pgroup_by_field_names,
	unsigned long long sample_count)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sample_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sample_state_t));

	pstate->pargp                 = pargp;
	pstate->pgroup_by_field_names = pgroup_by_field_names;
	pstate->sample_count          = sample_count;
	pstate->pbuckets_by_group     = lhmslv_alloc();

	pmapper->pvstate              = pstate;
	pmapper->pprocess_func        = mapper_sample_process;
	pmapper->pfree_func           = mapper_sample_free;

	return pmapper;
}

static void mapper_sample_free(mapper_t* pmapper) {
	mapper_sample_state_t* pstate = pmapper->pvstate;
	if (pstate->pgroup_by_field_names != NULL)
		slls_free(pstate->pgroup_by_field_names);
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->pbuckets_by_group->phead; pa != NULL; pa = pa->pnext) {
		sample_bucket_t* pbucket = pa->pvvalue;
		sample_bucket_free(pbucket);
	}
	lhmslv_free(pstate->pbuckets_by_group);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_sample_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_sample_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pgroup_by_field_values = mlr_reference_selected_values_from_record(pinrec,
			pstate->pgroup_by_field_names);
		if (pgroup_by_field_values != NULL) {
			sample_bucket_t* pbucket = lhmslv_get(pstate->pbuckets_by_group, pgroup_by_field_values);
			if (pbucket == NULL) {
				pbucket = sample_bucket_alloc(pstate->sample_count);
				lhmslv_put(pstate->pbuckets_by_group, slls_copy(pgroup_by_field_values), pbucket,
					FREE_ENTRY_KEY);
			}
			sample_bucket_handle(pbucket, pinrec, pctx->nr);
			slls_free(pgroup_by_field_values);
		} else {
			lrec_free(pinrec);
		}
		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pa = pstate->pbuckets_by_group->phead; pa != NULL; pa = pa->pnext) {
			sample_bucket_t* pbucket = pa->pvvalue;
			for (int i = 0; i < pbucket->nused; i++) {
				sllv_append(poutrecs, pbucket->plrecs[i]);
				pbucket->plrecs[i] = NULL;
			}
			pbucket->nused = 0;
		}
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
sample_bucket_t* sample_bucket_alloc(int nalloc) {
	sample_bucket_t* pbucket = mlr_malloc_or_die(sizeof(sample_bucket_t));
	pbucket->nalloc = nalloc;
	pbucket->nused  = 0;
	pbucket->plrecs = mlr_malloc_or_die(nalloc * sizeof(lrec_t*));
	return pbucket;
}

void sample_bucket_free(sample_bucket_t* pbucket) {
	for (int i = 0; i < pbucket->nused; i++)
		lrec_free(pbucket->plrecs[i]);
	free(pbucket->plrecs);
	free(pbucket);
}

// This is the reservoir-sampling algorithm.
// Here we retain a pointer to an input record (if retained in the sample) or
// free it (if not retained in the sample).
void sample_bucket_handle(sample_bucket_t* pbucket, lrec_t* prec, int record_number) {
	if (pbucket->nused < pbucket->nalloc) {
		// Always accept new entries until the bucket is full
		pbucket->plrecs[pbucket->nused++] = prec;
	} else {
		int r = get_mtrand_int31() % record_number;
		if (r < pbucket->nalloc) {
			lrec_free(pbucket->plrecs[r]);
			pbucket->plrecs[r] = prec;
		} else {
			lrec_free(prec);
		}
	}
}
