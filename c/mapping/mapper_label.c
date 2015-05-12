#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_label_state_t {
	slls_t* pnames;
} mapper_label_state_t;

// ----------------------------------------------------------------
// xxx comment why (--inidx)
// xxx comment what happens after end-of-list (or before)

// xxx make all such file-static
static sllv_t* mapper_label_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_label_state_t* pstate = (mapper_label_state_t*)pvstate;
		lrece_t* pe = pinrec->phead;
		sllse_t* pn = pstate->pnames->phead;
		for ( ; pe != NULL && pn != NULL; pe = pe->pnext, pn = pn->pnext) {
			char* old_name = pe->key;
			char* new_name = pn->value;
			lrec_rename(pinrec, old_name, new_name);
		}
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static void mapper_label_free(void* pvstate) {
	mapper_label_state_t* pstate = (mapper_label_state_t*)pvstate;
	if (pstate->pnames != NULL)
		slls_free(pstate->pnames);
}

static mapper_t* mapper_label_alloc(slls_t* pnames) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_label_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_label_state_t));
	pstate->pnames   = pnames;

	pmapper->pvstate              = (void*)pstate;
	pmapper->pmapper_process_func = mapper_label_func;
	pmapper->pmapper_free_func    = mapper_label_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_label_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s {new1,new2,new3,...}\n", argv0, verb);
	fprintf(stdout, "Given n comma-separated names, renames the first n fields of each record to\n");
	fprintf(stdout, "have the specified name. (Fields past the nth are left with their original\n");
	fprintf(stdout, "names.) Particularly useful with --inidx, to give useful names to otherwise\n");
	fprintf(stdout, "integer-indexed fields.\n");
}

static mapper_t* mapper_label_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 2) {
		mapper_label_usage(argv[0], argv[*pargi]);
		return NULL;
	}

	slls_t* pnames = slls_from_line(argv[*pargi+1], ',', FALSE);

	*pargi += 2;
	return mapper_label_alloc(pnames);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_label_setup = {
	.verb = "label",
	.pusage_func = mapper_label_usage,
	.pparse_func = mapper_label_parse_cli
};
