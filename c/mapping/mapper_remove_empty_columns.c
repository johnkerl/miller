#include <stdio.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmsi.h"
#include "mapping/mappers.h"

typedef struct _mapper_remove_empty_columns_state_t {
	sllv_t* precords;
	lhmsi_t*  pnames_with_nonempty_values;
} mapper_remove_empty_columns_state_t;

static void      mapper_remove_empty_columns_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_remove_empty_columns_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_remove_empty_columns_alloc();
static void      mapper_remove_empty_columns_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_remove_empty_columns_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_remove_empty_columns_setup = {
	.verb = "remove-empty-columns",
	.pusage_func = mapper_remove_empty_columns_usage,
	.pparse_func = mapper_remove_empty_columns_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_remove_empty_columns_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Omits fields which are empty on every input row. Non-streaming.\n");
}

static mapper_t* mapper_remove_empty_columns_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 1) {
		mapper_remove_empty_columns_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_remove_empty_columns_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_remove_empty_columns_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_remove_empty_columns_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_remove_empty_columns_state_t));
	pstate->precords = sllv_alloc();
	pstate->pnames_with_nonempty_values = lhmsi_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_remove_empty_columns_process;
	pmapper->pfree_func    = mapper_remove_empty_columns_free;

	return pmapper;
}

static void mapper_remove_empty_columns_free(mapper_t* pmapper, context_t* _) {
	mapper_remove_empty_columns_state_t* pstate = pmapper->pvstate;
	// Free the container
	sllv_free(pstate->precords);
	lhmsi_free(pstate->pnames_with_nonempty_values);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_remove_empty_columns_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_remove_empty_columns_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		// The caller will free the outrecs
		sllv_append(pstate->precords, pinrec);
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			if (pe->value[0] != 0) {
				if (!lhmsi_has_key(pstate->pnames_with_nonempty_values, pe->key)) {
					lhmsi_put(pstate->pnames_with_nonempty_values, mlr_strdup_or_die(pe->key), 1, FREE_ENTRY_KEY);
				}
			}
		}
		return NULL;
	}
	else {
		for (sllve_t* pe = pstate->precords->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* prec = pe->pvvalue;
			lrece_t* pf = prec->phead;
			while (pf != NULL) {
				if (lhmsi_has_key(pstate->pnames_with_nonempty_values, pf->key)) {
					pf = pf->pnext;
				} else {
					lrece_t* pnext = pf->pnext;
					lrec_unlink_and_free(prec, pf);
					pf = pnext;
				}
			}
		}

		sllv_append(pstate->precords, NULL);
		sllv_t* retval = pstate->precords;
		pstate->precords = sllv_alloc();
		return retval;
	}
}
