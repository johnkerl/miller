#include <stdio.h>
#include "cli/argparse.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

#define NOUT_EQUALS_NIN -1
typedef struct _mapper_bootstrap_state_t {
	ap_state_t* pargp;
	int nout;
	sllv_t* records;
} mapper_bootstrap_state_t;

static void      mapper_bootstrap_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_bootstrap_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_bootstrap_alloc(int nout, ap_state_t* pargp);
static void      mapper_bootstrap_free(mapper_t* pmapper);
static sllv_t*   mapper_bootstrap_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_bootstrap_setup = {
	.verb = "bootstrap",
	.pusage_func = mapper_bootstrap_usage,
	.pparse_func = mapper_bootstrap_parse_cli
};

// ----------------------------------------------------------------
static void mapper_bootstrap_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Emits an n-sample, with replacement, of the input records.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-n {number} Number of samples to output. Defaults to number of input records.\n");
	fprintf(o, "            Must be non-negative.\n");
}

static mapper_t* mapper_bootstrap_parse_cli(int* pargi, int argc, char** argv) {
	int nout = NOUT_EQUALS_NIN;
	if ((argc - *pargi) < 1) {
		mapper_bootstrap_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}

	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_int_flag(pstate, "-n", &nout);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_bootstrap_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (nout != NOUT_EQUALS_NIN && nout < 0) {
		mapper_bootstrap_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_bootstrap_alloc(nout, pstate);
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_bootstrap_alloc(int nout, ap_state_t* pargp) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_bootstrap_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_bootstrap_state_t));
	pstate->nout    = nout;
	pstate->pargp   = pargp;
	pstate->records = sllv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_bootstrap_process;
	pmapper->pfree_func    = mapper_bootstrap_free;

	return pmapper;
}

static void mapper_bootstrap_free(mapper_t* pmapper) {
	mapper_bootstrap_state_t* pstate = pmapper->pvstate;
	// Free the container
	sllv_free(pstate->records);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_bootstrap_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_bootstrap_state_t* pstate = pvstate;
	if (pinrec != NULL) { // Not end of input stream: consume an input record.
		// The caller will free the outrecs
		sllv_append(pstate->records, pinrec);
		return NULL;
	}

	// This is entirely straightforward except for memory management. The
	// Miller contract to respect is that each routine frees records or passes
	// them on.  This means lrec-writers will free lrecs they receive, after
	// writing them.  This also means we must free lrecs we don't pass on.
	//
	// Given nin input records, we produce nout output records, but sampling with replacement.
	// The memory-management criteria above mean:
	// * If an input lrec is not output at all, we must free it.
	// * If an input lrec is output once, we pass it on through and let the write free it.
	// * If an input lrec is output more than once, it would get double-freed (which is not OK)
	//   so for all repetitions past the first we must make a copy.
	// A used_flags[] array allows us to handle all of this.

	sllv_t* poutrecs = sllv_alloc();
	int nin = pstate->records->length;
	int nout = (pstate->nout == NOUT_EQUALS_NIN) ? nin : pstate->nout;
	if (nin == 0) {
		sllv_append(poutrecs, NULL);
		return poutrecs;
	}

	// Make an array of pointers into the input list. Mark each lrec as not yet output.
	lrec_t** record_array = mlr_malloc_or_die(nin * sizeof(lrec_t**));
	char* used_flags = mlr_malloc_or_die(nin * sizeof(char));
	sllve_t* pe = pstate->records->phead;
	for (int i = 0; i < nin; i++, pe = pe->pnext) {
		record_array[i] = pe->pvvalue;
		used_flags[i] = FALSE;
	}

	// Do the sample-with-replacment, reading from random indices in the input
	// array and appending to the output list.
	for (int i = 0; i < nout; i++) {
		int index = nin * get_mtrand_double();
		if (index >= nin)
			index = nin - 1;
		lrec_t* prec = record_array[index];
		if (used_flags[index]) { // Copy for repeated output.
			prec = lrec_copy(prec);
		} else { // First output of this record; remember it.
			used_flags[index] = TRUE;
		}
		sllv_append(poutrecs, prec);
	}

	// Free non-output records
	pe = pstate->records->phead;
	for (int i = 0; i < nin; i++, pe = pe->pnext)
		if (!used_flags[i])
			lrec_free(record_array[i]);
	free(used_flags);
	// Free the temp array
	free(record_array);

	// Null-terminate the output list to signify end of stream.
	sllv_append(poutrecs, NULL);
	return poutrecs;
}
