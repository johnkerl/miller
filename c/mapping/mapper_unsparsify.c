#include <stdio.h>
#include "lib/mlrutil.h"
#include "containers/lhmsi.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_unsparsify_state_t {
	lhmsi_t* key_names;
	sllv_t* records;
} mapper_unsparsify_state_t;

static void      mapper_unsparsify_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_unsparsify_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_unsparsify_alloc();
static void      mapper_unsparsify_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_unsparsify_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_unsparsify_setup = {
	.verb = "unsparsify",
	.pusage_func = mapper_unsparsify_usage,
	.pparse_func = mapper_unsparsify_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_unsparsify_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Prints records with the union of field names over all input records.\n");
	fprintf(o, "For field names absent in a given record but present in others, fills value with\n");
	fprintf(o, "empty string. This verb retains all input before producing any output.\n");
}

static mapper_t* mapper_unsparsify_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_unsparsify_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_unsparsify_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_unsparsify_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_unsparsify_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_unsparsify_state_t));
	pstate->records = sllv_alloc();
	pstate->key_names = lhmsi_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_unsparsify_process;
	pmapper->pfree_func    = mapper_unsparsify_free;

	return pmapper;
}

static void mapper_unsparsify_free(mapper_t* pmapper, context_t* _) {
	mapper_unsparsify_state_t* pstate = pmapper->pvstate;
	// Free the container
	sllv_free(pstate->records);
	lhmsi_free(pstate->key_names);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_unsparsify_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_unsparsify_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		// Not end of stream.
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			if (!lhmsi_has_key(pstate->key_names, pe->key)) {
				lhmsi_put(pstate->key_names, mlr_strdup_or_die(pe->key), 1, FREE_ENTRY_KEY);
			}
		}
		// The caller will free the outrecs
		sllv_append(pstate->records, pinrec);
		return NULL;
	}
	else {
		// End of stream.
		sllv_t* poutrecs = sllv_alloc();
		for (sllve_t* pe = pstate->records->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* pinrec = pe->pvvalue;
			lrec_t* poutrec = lrec_unbacked_alloc();
			for (lhmsie_t* pf = pstate->key_names->phead; pf != NULL; pf = pf->pnext) {
				char* key = pf->key;
				char* value = lrec_get(pinrec, key);
				if (value == NULL) {
					lrec_put(poutrec, mlr_strdup_or_die(key), "", FREE_ENTRY_KEY);
				} else {
					lrec_put(poutrec, mlr_strdup_or_die(key), mlr_strdup_or_die(value),
						FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
				}
			}
			sllv_append(poutrecs, poutrec);
			// Free the void-star payload
			lrec_free(pinrec);
		}

		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}
