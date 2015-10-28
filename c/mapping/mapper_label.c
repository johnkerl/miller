#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"

typedef struct _mapper_label_state_t {
	slls_t* pnames;
} mapper_label_state_t;

static sllv_t*   mapper_label_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_label_free(void* pvstate);
static mapper_t* mapper_label_alloc(slls_t* pnames);
static void      mapper_label_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_label_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_label_setup = {
	.verb = "label",
	.pusage_func = mapper_label_usage,
	.pparse_func = mapper_label_parse_cli
};

// ----------------------------------------------------------------
static void mapper_label_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s {new1,new2,new3,...}\n", argv0, verb);
	fprintf(o, "Given n comma-separated names, renames the first n fields of each record to\n");
	fprintf(o, "have the respective name. (Fields past the nth are left with their original\n");
	fprintf(o, "names.) Particularly useful with --inidx or --implicit-csv-header, to give\n");
	fprintf(o, "useful names to otherwise integer-indexed fields.\n");
	fprintf(o, "Examples:\n");
	fprintf(o, "  \"echo 'a b c d' | %s --inidx --odkvp cat\"       gives \"1=a,2=b,3=c,4=d\"\n", argv0);
	fprintf(o, "  \"echo 'a b c d' | %s --inidx --odkvp %s s,t\" gives \"s=a,t=b,3=c,4=d\"\n", argv0, verb);

}

static mapper_t* mapper_label_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 2) {
		mapper_label_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}

	slls_t* pnames = slls_from_line(argv[*pargi+1], ',', FALSE);

	*pargi += 2;
	return mapper_label_alloc(pnames);
}

// ----------------------------------------------------------------
static mapper_t* mapper_label_alloc(slls_t* pnames) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_label_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_label_state_t));
	pstate->pnames = pnames;

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_label_process;
	pmapper->pfree_func    = mapper_label_free;

	return pmapper;
}

static void mapper_label_free(void* pvstate) {
	mapper_label_state_t* pstate = (mapper_label_state_t*)pvstate;
	if (pstate->pnames != NULL)
		slls_free(pstate->pnames);
}

// ----------------------------------------------------------------
static sllv_t* mapper_label_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_label_state_t* pstate = (mapper_label_state_t*)pvstate;
		lrece_t* pe = pinrec->phead;
		sllse_t* pn = pstate->pnames->phead;
		for ( ; pe != NULL && pn != NULL; pe = pe->pnext, pn = pn->pnext) {
			char* old_name = pe->key;
			char* new_name = pn->value;
			lrec_rename(pinrec, old_name, new_name, FALSE);
		}
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}
