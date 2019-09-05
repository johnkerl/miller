#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"

typedef struct _mapper_format_state_t {
	ap_state_t* pargp;
	char* string_format;
	char* int_format;
	char* float_format;
} mapper_format_state_t;

#define DEFAULT_STRING_FORMAT "%s"
#define DEFAULT_INT_FORMAT    "%lld"
#define DEFAULT_FLOAT_FORMAT  "%lf"

static void      mapper_format_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_format_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_format_alloc(ap_state_t* pargp, char* string_format, char* int_format, char* float_format);
static void      mapper_format_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_format_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_format_setup = {
	.verb = "format",
	.pusage_func = mapper_format_usage,
	.pparse_func = mapper_format_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_format_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* string_format = DEFAULT_STRING_FORMAT;
	char* int_format    = DEFAULT_INT_FORMAT;
	char* float_format  = DEFAULT_FLOAT_FORMAT;

	if ((argc - *pargi) < 1) {
		mapper_format_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "-s", &string_format);
	ap_define_string_flag(pstate, "-i", &int_format);
	ap_define_string_flag(pstate, "-f", &float_format);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_format_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_format_alloc(pstate, string_format, int_format, float_format);
	return pmapper;
}

static void mapper_format_usage(FILE* o, char* argv0, char* verb) {
	// xxx fmtnum is fine-grained. this is the sledgehammer which guesses for you.
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Passes input records directly to output. Most useful for format conversion.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-g {comma-separated field name(s)} When used with -n/-N, writes record-counters\n");
	fprintf(o, "          keyed by specified field name(s).\n");
	fprintf(o, "-v        Write a low-level record-structure dump to stderr.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_format_alloc(ap_state_t* pargp, char* string_format, char* int_format, char* float_format)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_format_state_t* pstate    = mlr_malloc_or_die(sizeof(mapper_format_state_t));
	pstate->pargp                 = pargp;
	pstate->string_format         = string_format;
	pstate->int_format            = int_format;
	pstate->float_format          = float_format;
	pmapper->pvstate              = pstate;

	pmapper->pprocess_func        = NULL;
	pmapper->pprocess_func        = mapper_format_process;

	pmapper->pfree_func           = mapper_format_free;
	return pmapper;
}
static void mapper_format_free(mapper_t* pmapper, context_t* _) {
	mapper_format_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_format_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_format_state_t* pstate = (mapper_format_state_t*)pvstate;
	if (pinrec != NULL) {
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			long long int_value;
			double float_value;
			if (pe->value == NULL) {
				printf("WTFC\n");
				continue;
			}
			char* string_value = pe->value;
			int is_int = mlr_try_int_from_string(string_value, &int_value);
			int is_float = mlr_try_float_from_string(string_value, &float_value);

			if (is_int) {
				lrece_update_value(pe, mlr_alloc_string_from_ll_and_format(int_value, pstate->int_format), TRUE);
			} else if (is_float) {
				lrece_update_value(pe, mlr_alloc_string_from_double(float_value, pstate->float_format), TRUE);
			} else {
				lrece_update_value(pe,
					mlr_alloc_string_from_string_and_format(string_value, pstate->string_format),
					TRUE
				);
			}

			// xxx int-to-float flag ...
		}
		return sllv_single(pinrec);
	} else { // end of record stream
		return sllv_single(NULL);
	}
}
