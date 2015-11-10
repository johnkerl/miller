#include <regex.h>
#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlr_globals.h"
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
	fprintf(o, "Usage: %s %s [options] {pattern}\n", argv0, verb);
	fprintf(o, "Passes through records which match {pattern}.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-i    Use case-insensitive search.\n");
	fprintf(o, "-v    Invert: pass through records which do not match the pattern.\n");
	fprintf(o, "Note that \"%s filter\" is more powerful, but requires you to know field names.\n", argv0);
	fprintf(o, "By contrast, \"%s %s\" allows you to pattern-match the entire record. It does\n", argv0, verb);
	fprintf(o, "this by formatting each record in memory as DKVP, using command-line-specified\n");
	fprintf(o, "ORS/OFS/OPS, and matching the resulting line against the pattern specified\n");
	fprintf(o, "here. Not all the options to system grep are supported, and this command\n");
	fprintf(o, "is intended to be merely a keystroke-saver. To get all the features\n");
	fprintf(o, "of system grep, you can do \"%s --odkvp ... | grep ... | %s --idkvp ...\"\n", argv0, argv0);
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

	char* line = lrec_sprint(pinrec,
		MLR_GLOBALS.popts->ors,
		MLR_GLOBALS.popts->ofs,
		MLR_GLOBALS.popts->ops);

	int matches = regmatch_or_die(&pstate->regex, line, 0, NULL);
	sllv_t* poutrecs = NULL;
	if (matches ^ pstate->exclude) {
		poutrecs = sllv_single(pinrec);
	}
	free(line);
	return poutrecs;
}
