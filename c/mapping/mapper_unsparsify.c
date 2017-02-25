#include <stdio.h>
#include "lib/mlrutil.h"
#include "cli/argparse.h"
#include "containers/lhmsi.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_unsparsify_state_t {
	lhmsi_t* key_names;
	sllv_t* records;
	char*   filler;
	ap_state_t* pargp;
} mapper_unsparsify_state_t;

static void      mapper_unsparsify_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_unsparsify_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_unsparsify_alloc(ap_state_t* pargp, char* filler);
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
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Prints records with the union of field names over all input records.\n");
	fprintf(o, "For field names absent in a given record but present in others, fills in\n");
	fprintf(o, "a value. This verb retains all input before producing any output.\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "--fill-with {filler string}  What to fill absent fields with. Defaults to\n");
	fprintf(o, "                             the empty string.\n");
	fprintf(o, "\n");
	fprintf(o, "Example: if the input is two records, one being 'a=1,b=2' and the other\n");
	fprintf(o, "being 'b=3,c=4', then the output is the two records 'a=1,b=2,c=' and\n");
	fprintf(o, "'a=,b=3,c=4'.\n");
}

static mapper_t* mapper_unsparsify_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* filler = "";

	if ((argc - *pargi) < 1) {
		mapper_unsparsify_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pargp = ap_alloc();
	ap_define_string_flag(pargp, "--fill-with", &filler);
	if (!ap_parse(pargp, verb, pargi, argc, argv)) {
		mapper_unsparsify_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_unsparsify_alloc(pargp, filler);
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_unsparsify_alloc(ap_state_t* pargp, char* filler) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_unsparsify_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_unsparsify_state_t));
	pstate->records = sllv_alloc();
	pstate->key_names = lhmsi_alloc();
	pstate->filler = filler;
	pstate->pargp = pargp;

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
	ap_free(pstate->pargp);
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
					lrec_put(poutrec, mlr_strdup_or_die(key), pstate->filler, FREE_ENTRY_KEY);
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
