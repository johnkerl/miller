#include <regex.h>
#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/sllv.h"

typedef struct _mapper_grep_state_t {
	int exclude;
	regex_t regex;
} mapper_grep_state_t;

static sllv_t*   mapper_grep_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_grep_free(void* pvstate);
static mapper_t* mapper_grep_alloc(char* regex_string, int exclude, int ignore_case);
static void      mapper_grep_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_grep_parse_cli(int* pargi, int argc, char** argv);

// ----------------------------------------------------------------
mapper_setup_t mapper_grep_setup = {
	.verb = "grep",
	.pusage_func = mapper_grep_usage,
	.pparse_func = mapper_grep_parse_cli
};

// ----------------------------------------------------------------
static mapper_t* mapper_grep_parse_cli(int* pargi, int argc, char** argv) {
	char* regex_string = NULL;
	int   exclude = FALSE;
	int   ignore_case = FALSE;

	if ((argc - *pargi) < 1) {
		mapper_grep_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-v", &exclude);
	ap_define_true_flag(pstate, "-i", &ignore_case);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_grep_usage(stderr, argv[0], verb);
		return NULL;
	}

	if ((argc - *pargi) < 1) {
		mapper_grep_usage(stderr, argv[0], verb);
		return NULL;
	}

	regex_string = argv[(*pargi)++];

	mapper_t* pmapper = mapper_grep_alloc(regex_string, exclude, ignore_case);
	return pmapper;
}
static void mapper_grep_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "[under construction]\n");
	fprintf(o, "[note limitations; most useful for when you don't know the field name(s) for filter]\n");
	fprintf(o, "[note mlr --odkvp ... | grep ... | grep mlr --idkvp]\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_grep_alloc(char* regex_string, int exclude, int ignore_case) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_grep_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_grep_state_t));
	int cflags = REG_NOSUB;
	if (ignore_case)
		cflags |= REG_ICASE;
	regcomp_or_die_quoted(&pstate->regex, regex_string, cflags);
	pstate->exclude = exclude;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_grep_process;
	pmapper->pfree_func    = mapper_grep_free;
	return pmapper;
}
static void mapper_grep_free(void* pvstate) {
	mapper_grep_state_t* pstate = (mapper_grep_state_t*)pvstate;
	regfree(&pstate->regex);
}

// ----------------------------------------------------------------
static sllv_t* mapper_grep_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return sllv_single(NULL);

	mapper_grep_state_t* pstate = (mapper_grep_state_t*)pvstate;

	// xxx temp -- should be from CLI
	char* ors = "\n";
	char* ofs = ",";
	char* ops = "=";

	char* line = lrec_sprint(pinrec, ors, ofs, ops);

	int matches = regmatch_or_die(&pstate->regex, line, 0, NULL);
	sllv_t* poutrecs = NULL;
	if (matches ^ pstate->exclude) {
		poutrecs = sllv_single(pinrec);
}
	free(line);
	return poutrecs;
}
