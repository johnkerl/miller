#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "lib/mvfuncs.h"
#include "containers/slls.h"
#include "containers/mixutil.h"
#include "cli/argparse.h"
#include "mapping/mappers.h"

typedef struct _mapper_sec2gmt_state_t {
	slls_t*  pfield_names;
	char* format_string;
} mapper_sec2gmt_state_t;

static void      mapper_sec2gmt_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_sec2gmt_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_sec2gmt_alloc(slls_t* pfield_names, int num_decimal_places);
static void      mapper_sec2gmt_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_sec2gmt_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_sec2gmt_setup = {
	.verb = "sec2gmt",
	.pusage_func = mapper_sec2gmt_usage,
	.pparse_func = mapper_sec2gmt_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_sec2gmt_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options] {comma-separated list of field names}\n", argv0, verb);
	fprintf(o, "Replaces a numeric field representing seconds since the epoch with the\n");
	fprintf(o, "corresponding GMT timestamp; leaves non-numbers as-is. This is nothing\n");
	fprintf(o, "more than a keystroke-saver for the sec2gmt function:\n");
	fprintf(o, "  %s %s time1,time2\n", argv0, verb);
	fprintf(o, "is the same as\n");
	fprintf(o, "  %s put '$time1=sec2gmt($time1);$time2=sec2gmt($time2)'\n", argv0);
	fprintf(o, "Options:\n");
	fprintf(o, "-1 through -9: format the seconds using 1..9 decimal places, respectively.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_sec2gmt_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int num_decimal_places = 0;

	char* verb = argv[*pargi];
	if ((argc - *pargi) < 1) {
		mapper_sec2gmt_usage(stderr, argv[0], verb);
		return NULL;
	}
	*pargi += 1;

    ap_state_t* pstate = ap_alloc();
	ap_define_int_value_flag(pstate, "-1", 1, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-2", 2, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-3", 3, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-4", 4, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-5", 5, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-6", 6, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-7", 7, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-8", 8, &num_decimal_places);
	ap_define_int_value_flag(pstate, "-9", 9, &num_decimal_places);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_sec2gmt_usage(stderr, argv[0], verb);
		return NULL;
	}
    ap_free(pstate);

	if ((argc - *pargi) < 1) {
		mapper_sec2gmt_usage(stderr, argv[0], verb);
		return NULL;
	}

	char* field_names_string = argv[(*pargi)++];
	slls_t* pfield_names = slls_from_line(field_names_string, ',', FALSE);

	return mapper_sec2gmt_alloc(pfield_names, num_decimal_places);
}

// ----------------------------------------------------------------
static mapper_t* mapper_sec2gmt_alloc(slls_t* pfield_names, int num_decimal_places) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_sec2gmt_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_sec2gmt_state_t));
	pstate->pfield_names   = pfield_names;

	switch(num_decimal_places) {
	case 0: pstate->format_string  = ISO8601_TIME_FORMAT;   break;
	case 1: pstate->format_string  = ISO8601_TIME_FORMAT_1; break;
	case 2: pstate->format_string  = ISO8601_TIME_FORMAT_2; break;
	case 3: pstate->format_string  = ISO8601_TIME_FORMAT_3; break;
	case 4: pstate->format_string  = ISO8601_TIME_FORMAT_4; break;
	case 5: pstate->format_string  = ISO8601_TIME_FORMAT_5; break;
	case 6: pstate->format_string  = ISO8601_TIME_FORMAT_6; break;
	case 7: pstate->format_string  = ISO8601_TIME_FORMAT_7; break;
	case 8: pstate->format_string  = ISO8601_TIME_FORMAT_8; break;
	case 9: pstate->format_string  = ISO8601_TIME_FORMAT_9; break;
	default: MLR_INTERNAL_CODING_ERROR(); break;
	}

	pmapper->pprocess_func = mapper_sec2gmt_process;
	pmapper->pvstate       = (void*)pstate;
	pmapper->pfree_func    = mapper_sec2gmt_free;

	return pmapper;
}

static void mapper_sec2gmt_free(mapper_t* pmapper, context_t* _) {
	mapper_sec2gmt_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->pfield_names);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_sec2gmt_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of stream
		return sllv_single(NULL);

	mapper_sec2gmt_state_t* pstate = (mapper_sec2gmt_state_t*)pvstate;

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
				mv_t stamp = time_string_from_seconds(&mval, pstate->format_string, TIME_FROM_SECONDS_GMT);
				lrec_put(pinrec, name, stamp.u.strv, FREE_ENTRY_VALUE);
			}
		}
	}
	return sllv_single(pinrec);
}
