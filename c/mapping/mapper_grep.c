#include <regex.h>
#include "cli/argparse.h"
#include "cli/mlrcli.h"
#include "mapping/mappers.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "containers/sllv.h"

typedef struct _mapper_grep_state_t {
	ap_state_t* pargp;
	int exclude;
	regex_t regex;
	cli_writer_opts_t* pwriter_opts;
} mapper_grep_state_t;

static void      mapper_grep_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_grep_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_grep_alloc(ap_state_t* pargp, char* regex_string, int exclude, int ignore_case,
	cli_writer_opts_t* pwriter_opts);
static void      mapper_grep_free(mapper_t* pmapper);
static sllv_t*   mapper_grep_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_grep_setup = {
	.verb = "grep",
	.pusage_func = mapper_grep_usage,
	.pparse_func = mapper_grep_parse_cli
};

// ----------------------------------------------------------------
static mapper_t* mapper_grep_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* pwriter_opts)
{
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

	mapper_t* pmapper = mapper_grep_alloc(pstate, regex_string, exclude, ignore_case, pwriter_opts);
	return pmapper;
}
static void mapper_grep_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {regular expression}\n", argv0, verb);
	fprintf(o, "Passes through records which match {regex}.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-i    Use case-insensitive search.\n");
	fprintf(o, "-v    Invert: pass through records which do not match the regex.\n");
	fprintf(o, "Note that \"%s filter\" is more powerful, but requires you to know field names.\n", argv0);
	fprintf(o, "By contrast, \"%s %s\" allows you to regex-match the entire record. It does\n", argv0, verb);
	fprintf(o, "this by formatting each record in memory as DKVP, using command-line-specified\n");
	fprintf(o, "ORS/OFS/OPS, and matching the resulting line against the regex specified\n");
	fprintf(o, "here. In particular, the regex is not applied to the input stream: if you\n");
	fprintf(o, "have CSV with header line \"x,y,z\" and data line \"1,2,3\" then the regex will\n");
	fprintf(o, "be matched, not against either of these lines, but against the DKVP line\n");
	fprintf(o, "\"x=1,y=2,z=3\".  Furthermore, not all the options to system grep are supported,\n");
	fprintf(o, "and this command is intended to be merely a keystroke-saver. To get all the\n");
	fprintf(o, "features of system grep, you can do\n");
	fprintf(o, "  \"%s --odkvp ... | grep ... | %s --idkvp ...\"\n", argv0, argv0);
}

// ----------------------------------------------------------------
static mapper_t* mapper_grep_alloc(ap_state_t* pargp, char* regex_string, int exclude, int ignore_case,
	cli_writer_opts_t* pwriter_opts)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_grep_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_grep_state_t));
	pstate->pargp = pargp;
	int cflags = REG_NOSUB;
	if (ignore_case)
		cflags |= REG_ICASE;
	regcomp_or_die_quoted(&pstate->regex, regex_string, cflags);
	pstate->exclude = exclude;
	pstate->pwriter_opts = pwriter_opts;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_grep_process;
	pmapper->pfree_func    = mapper_grep_free;
	return pmapper;
}
static void mapper_grep_free(mapper_t* pmapper) {
	mapper_grep_state_t* pstate = pmapper->pvstate;
	regfree(&pstate->regex);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_grep_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of input stream
		return sllv_single(NULL);

	mapper_grep_state_t* pstate = (mapper_grep_state_t*)pvstate;

	char* line = lrec_sprint(pinrec,
		pstate->pwriter_opts->ors,
		pstate->pwriter_opts->ofs,
		pstate->pwriter_opts->ops);

	int matches = regmatch_or_die(&pstate->regex, line, 0, NULL);
	sllv_t* poutrecs = NULL;
	if (matches ^ pstate->exclude) {
		poutrecs = sllv_single(pinrec);
	} else {
		lrec_free(pinrec);
	}
	free(line);
	return poutrecs;
}
