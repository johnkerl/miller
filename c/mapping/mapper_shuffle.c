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
typedef struct _mapper_shuffle_state_t {
	ap_state_t* pargp;
	sllv_t*     precs;
} mapper_shuffle_state_t;

static void      mapper_shuffle_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_shuffle_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_shuffle_alloc(ap_state_t* pargp);
static void      mapper_shuffle_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_shuffle_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_shuffle_setup = {
	.verb = "shuffle",
	.pusage_func = mapper_shuffle_usage,
	.pparse_func = mapper_shuffle_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_shuffle_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s {no options}\n", argv0, verb);
	fprintf(o, "Outputs records randomly permuted. No output records are produced until\n");
	fprintf(o, "all input records are read.\n");
	fprintf(o, "See also %s bootstrap and %s sample.\n", argv0, argv0);
}

static mapper_t* mapper_shuffle_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_shuffle_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_shuffle_alloc(pstate);
}

// ----------------------------------------------------------------
static mapper_t* mapper_shuffle_alloc(ap_state_t* pargp) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_shuffle_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_shuffle_state_t));

	pstate->pargp = pargp;
	pstate->precs = sllv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_shuffle_process;
	pmapper->pfree_func    = mapper_shuffle_free;

	return pmapper;
}

static void mapper_shuffle_free(mapper_t* pmapper, context_t* _) {
	mapper_shuffle_state_t* pstate = pmapper->pvstate;
	// Records will have been freed by the emitter; here, free the list structure.
	sllv_free(pstate->precs);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_shuffle_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_shuffle_state_t* pstate = pvstate;

	// Not end of input stream: retain the record, and emit nothing until end of stream.
	if (pinrec != NULL) {
		sllv_append(pstate->precs, pinrec);
		return NULL;
	}

	sllv_t* poutrecs = sllv_alloc();

	// Knuth shuffle:
	// * Initial permutation is identity.
	// * Make a pseudorandom permutation using pseudorandom swaps in the image map.
	int n = pstate->precs->length;
	int* images = mlr_malloc_or_die(n * sizeof(int));
	for (int i = 0; i < n; i++)
		images[i] = i;

	int unused_start = 0;
	int num_unused   = n;
	for (int i = 0; i < n; i++) {
		// Select a pseudorandom element from the pool of unused images.
		int u = unused_start + num_unused * get_mtrand_double();
		int temp  = images[u];
		images[u] = images[i];
		images[i] = temp;

		// Decrease the size of the pool by 1.  (Yes, unused_start and k always have the same value.
		// Using two variables wastes neglible memory and makes the code easier to understand.)
		unused_start++;
		num_unused--;
	}

	// Make an array of pointers into the input list.
	lrec_t** record_array = mlr_malloc_or_die(n * sizeof(lrec_t**));
	sllve_t* pe = pstate->precs->phead;
	for (int i = 0; i < n; i++, pe = pe->pnext) {
		record_array[i] = pe->pvvalue;
	}

	// Transfer from input array to output list. Because permutations are one-to-one maps,
	// all input records have ownership transferred exactly once. So, there are no
	// records to copy here, or free here.
	for (int i = 0; i < n; i++) {
		sllv_append(poutrecs, record_array[images[i]]);
	}

	free(record_array);
	free(images);

	// Null-terminate the output list to signify end of stream.
	sllv_append(poutrecs, NULL);
	return poutrecs;
}

