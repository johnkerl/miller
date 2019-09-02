#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/hss.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

typedef struct _mapper_label_state_t {
	slls_t* pnames_as_list;
	hss_t* pnames_as_set; // needed for efficient deduplication
} mapper_label_state_t;

static void      mapper_label_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_label_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_label_alloc(slls_t* pnames_as_list, hss_t* pnames_as_set);
static void      mapper_label_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_label_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_label_setup = {
	.verb = "label",
	.pusage_func = mapper_label_usage,
	.pparse_func = mapper_label_parse_cli,
	.ignores_input = FALSE,
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

static mapper_t* mapper_label_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int argi = *pargi;
	if ((argc - argi) < 2) {
		mapper_label_usage(stderr, argv[0], argv[argi]);
		return NULL;
	}

	char* verb = argv[argi];
	char* names_as_string = argv[argi+1];

	slls_t* pnames_as_list = slls_from_line(names_as_string, ',', FALSE);
	hss_t* pnames_as_set = hss_from_slls(pnames_as_list);

	if (slls_size(pnames_as_list) != hss_size(pnames_as_set)) {
		fprintf(stderr, "%s %s: labels must be unique; got \"%s\"\n",
			MLR_GLOBALS.bargv0, verb, names_as_string);
		exit(1);
	}

	*pargi += 2;
	return mapper_label_alloc(pnames_as_list, pnames_as_set);
}

// ----------------------------------------------------------------
static mapper_t* mapper_label_alloc(slls_t* pnames_as_list, hss_t* pnames_as_set) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_label_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_label_state_t));
	pstate->pnames_as_list = pnames_as_list;
	pstate->pnames_as_set = pnames_as_set;

	pmapper->pvstate       = (void*)pstate;
	pmapper->pprocess_func = mapper_label_process;
	pmapper->pfree_func    = mapper_label_free;

	return pmapper;
}

static void mapper_label_free(mapper_t* pmapper, context_t* _) {
	mapper_label_state_t* pstate = pmapper->pvstate;
	if (pstate->pnames_as_list != NULL)
		slls_free(pstate->pnames_as_list);
	if (pstate->pnames_as_set != NULL)
		hss_free(pstate->pnames_as_set);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_label_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		mapper_label_state_t* pstate = (mapper_label_state_t*)pvstate;
		lrec_label(pinrec, pstate->pnames_as_list, pstate->pnames_as_set);
		return sllv_single(pinrec);
	}
	else {
		return sllv_single(NULL);
	}
}
