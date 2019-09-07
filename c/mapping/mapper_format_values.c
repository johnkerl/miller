#include "cli/argparse.h"
#include "mapping/mappers.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"

typedef struct _mapper_format_values_state_t {
	ap_state_t* pargp;
	char* string_format;
	char* int_format;
	char* float_format;
	int coerce_int_to_float;
} mapper_format_values_state_t;

#define DEFAULT_STRING_FORMAT "%s"
#define DEFAULT_INT_FORMAT    "%lld"
#define DEFAULT_FLOAT_FORMAT  "%lf"

static void      mapper_format_values_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_format_values_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_format_values_alloc(ap_state_t* pargp,
	char* string_format, char* int_format, char* float_format,
	int coerce_int_to_float);
static void      mapper_format_values_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_format_values_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_format_values_setup = {
	.verb = "format-values",
	.pusage_func = mapper_format_values_usage,
	.pparse_func = mapper_format_values_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static mapper_t* mapper_format_values_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	char* string_format = DEFAULT_STRING_FORMAT;
	char* int_format    = DEFAULT_INT_FORMAT;
	char* float_format  = DEFAULT_FLOAT_FORMAT;
	int   coerce_int_to_float = FALSE;

	if ((argc - *pargi) < 1) {
		mapper_format_values_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	char* verb = argv[*pargi];
	*pargi += 1;

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "-s", &string_format);
	ap_define_string_flag(pstate, "-i", &int_format);
	ap_define_string_flag(pstate, "-f", &float_format);
	ap_define_true_flag(pstate, "-n", &coerce_int_to_float);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_format_values_usage(stderr, argv[0], verb);
		return NULL;
	}

	mapper_t* pmapper = mapper_format_values_alloc(pstate, string_format, int_format, float_format,
		coerce_int_to_float);
	return pmapper;
}

static void mapper_format_values_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Applies format strings to all field values, depending on autodetected type.\n");
	fprintf(o, "* If a field value is detected to be integer, applies integer format.\n");
	fprintf(o, "* Else, if a field value is detected to be float, applies float format.\n");
	fprintf(o, "* Else, applies string format.\n");
	fprintf(o, "\n");
	fprintf(o, "Note: this is a low-keystroke way to apply formatting to many fields. To get\n");
	fprintf(o, "finer control, please see the fmtnum function within the mlr put DSL.\n");
	fprintf(o, "\n");
	fprintf(o, "Note: this verb lets you apply arbitrary format strings, which can produce\n");
	fprintf(o, "undefined behavior and/or program crashes.  See your system's \"man printf\".\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-i {integer format} Defaults to \"%s\".\n", DEFAULT_INT_FORMAT);
	fprintf(o, "                    Examples: \"%%06lld\", \"%%08llx\".\n");
	fprintf(o, "                    Note that Miller integers are long long so you must use\n");
	fprintf(o, "                    formats which apply to long long, e.g. with ll in them.\n");
	fprintf(o, "                    Undefined behavior results otherwise.\n");
	fprintf(o, "-f {float format}   Defaults to \"%s\".\n", DEFAULT_FLOAT_FORMAT);
	fprintf(o, "                    Examples: \"%%8.3lf\", \"%%.6le\".\n");
	fprintf(o, "                    Note that Miller floats are double-precision so you must\n");
	fprintf(o, "                    use formats which apply to double, e.g. with l[efg] in them.\n");
	fprintf(o, "                    Undefined behavior results otherwise.\n");
	fprintf(o, "-s {string format}  Defaults to \"%s\".\n", DEFAULT_STRING_FORMAT);
	fprintf(o, "                    Examples: \"_%%s\", \"%%08s\".\n");
	fprintf(o, "                    Note that you must use formats which apply to string, e.g.\n");
	fprintf(o, "                    with s in them. Undefined behavior results otherwise.\n");
	fprintf(o, "-n                  Coerce field values autodetected as int to float, and then\n");
	fprintf(o, "                    apply the float format.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_format_values_alloc(ap_state_t* pargp,
	char* string_format, char* int_format, char* float_format,
	int coerce_int_to_float)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));
	mapper_format_values_state_t* pstate    = mlr_malloc_or_die(sizeof(mapper_format_values_state_t));
	pstate->pargp                 = pargp;
	pstate->string_format         = string_format;
	pstate->int_format            = int_format;
	pstate->float_format          = float_format;
	pstate->coerce_int_to_float   = coerce_int_to_float;
	pmapper->pvstate              = pstate;

	pmapper->pprocess_func        = NULL;
	pmapper->pprocess_func        = mapper_format_values_process;

	pmapper->pfree_func           = mapper_format_values_free;
	return pmapper;
}
static void mapper_format_values_free(mapper_t* pmapper, context_t* _) {
	mapper_format_values_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_format_values_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_format_values_state_t* pstate = (mapper_format_values_state_t*)pvstate;
	if (pinrec == NULL) { // end of record stream
		return sllv_single(NULL);
	}

	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		long long int_value;
		double float_value;
		char* string_value = pe->value;
		int is_int = mlr_try_int_from_string(string_value, &int_value);
		int is_float = mlr_try_float_from_string(string_value, &float_value);

		if (is_int) {
			if (pstate->coerce_int_to_float) {
				lrece_update_value(pe, mlr_alloc_string_from_double((double)int_value, pstate->float_format), TRUE);
			} else {
				lrece_update_value(pe, mlr_alloc_string_from_ll_and_format(int_value, pstate->int_format), TRUE);
			}
		} else if (is_float) {
			lrece_update_value(pe, mlr_alloc_string_from_double(float_value, pstate->float_format), TRUE);
		} else {
			lrece_update_value(pe,
				mlr_alloc_string_from_string_and_format(string_value, pstate->string_format),
				TRUE
			);
		}
	}
	return sllv_single(pinrec);
}
