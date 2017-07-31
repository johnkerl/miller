#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "lib/mvfuncs.h"
#include "containers/slls.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

typedef struct _mapper_sec2gmtdate_state_t {
	slls_t*  pfield_names;
} mapper_sec2gmtdate_state_t;

static void      mapper_sec2gmtdate_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_sec2gmtdate_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_sec2gmtdate_alloc(slls_t* pfield_names);
static void      mapper_sec2gmtdate_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_sec2gmtdate_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_sec2gmtdate_setup = {
	.verb = "sec2gmtdate",
	.pusage_func = mapper_sec2gmtdate_usage,
	.pparse_func = mapper_sec2gmtdate_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_sec2gmtdate_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s {comma-separated list of field names}\n", argv0, verb);
	fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n");
	fprintf(o, "corresponding GMT year-month-day timestamp; leaves non-numbers as-is.\n");
	fprintf(o, "This is nothing more than a keystroke-saver for the sec2gmtdate function:\n");
	fprintf(o, "  %s %s time1,time2\n", argv0, verb);
	fprintf(o, "is the same as\n");
	fprintf(o, "  %s put '$time1=sec2gmtdate($time1);$time2=sec2gmtdate($time2)'\n", argv0);
}

// ----------------------------------------------------------------
static mapper_t* mapper_sec2gmtdate_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	if ((argc - *pargi) < 2) {
		mapper_sec2gmtdate_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	// verb:
	(*pargi)++;
	// field names:
	char* field_names_string = argv[(*pargi)++];
	slls_t* pfield_names = slls_from_line(field_names_string, ',', FALSE);

	return mapper_sec2gmtdate_alloc(pfield_names);
}

// ----------------------------------------------------------------
static mapper_t* mapper_sec2gmtdate_alloc(slls_t* pfield_names)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sec2gmtdate_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sec2gmtdate_state_t));
	pstate->pfield_names = pfield_names;
	pmapper->pprocess_func = mapper_sec2gmtdate_process;
	pmapper->pvstate       = (void*)pstate;
	pmapper->pfree_func    = mapper_sec2gmtdate_free;

	return pmapper;
}

static void mapper_sec2gmtdate_free(mapper_t* pmapper, context_t* _) {
	mapper_sec2gmtdate_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pfield_names);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_sec2gmtdate_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of stream
		return sllv_single(NULL);

	mapper_sec2gmtdate_state_t* pstate = (mapper_sec2gmtdate_state_t*)pvstate;

	for (sllse_t* pe = pstate->pfield_names->phead; pe != NULL; pe = pe->pnext) {
		char* name = pe->value;
		char* sval = lrec_get(pinrec, name);
		if (sval == NULL)
			continue;

		if (*sval == 0) {
			lrec_put(pinrec, name, "", NO_FREE);
		} else {
			mv_t mval = mv_scan_number_nullable(sval);
			if (!mv_is_error(&mval)) {
				mv_t stamp = time_string_from_seconds(&mval, ISO8601_DATE_FORMAT);
				lrec_put(pinrec, name, stamp.u.strv, FREE_ENTRY_VALUE);
			}
		}
	}
	return sllv_single(pinrec);
}
