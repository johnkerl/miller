#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/string_array.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_bar_state_t {
	string_array_t*  pfield_names;
	char   fill_char;
	char   oob_char;
	char   blank_char;
	double lo;
	double hi;
	int    width;
	char** bars;
} mapper_bar_state_t;

static sllv_t*   mapper_bar_process(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_bar_free(void* pvstate);
static mapper_t* mapper_bar_alloc(string_array_t* pfield_names,
	char fill_char, char oob_char, char blank_char, double lo, double hi, int width);
static void      mapper_bar_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_bar_parse_cli(int* pargi, int argc, char** argv);

#define DEFAULT_FILL_CHAR  '*'
#define DEFAULT_OOB_CHAR   '#'
#define DEFAULT_BLANK_CHAR '.'
#define DEFAULT_LO         0.0
#define DEFAULT_HI         100.0
#define DEFAULT_WIDTH      40

// ----------------------------------------------------------------
mapper_setup_t mapper_bar_setup = {
	.verb = "bar",
	.pusage_func = mapper_bar_usage,
	.pparse_func = mapper_bar_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_bar_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Replaces a numeric field with a number of asterisks, allowing for cheesy bar plots.\n");
	fprintf(o, "These align best with --opprint or --oxtab output format.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-f   {a,b,c}      Field names to convert to bars.\n");
	fprintf(o, "-c   {character}  Fill character: default '%c'.\n", DEFAULT_FILL_CHAR);
	fprintf(o, "-x   {character}  Out-of-bounds character: default '%c'.\n", DEFAULT_OOB_CHAR);
	fprintf(o, "-b   {character}  Blank character: default '%c'.\n", DEFAULT_BLANK_CHAR);
	fprintf(o, "--lo {lo}         Lower-limit value for min-width bar: default '%lf'.\n", DEFAULT_LO);
	fprintf(o, "--hi {hi}         Upper-limit value for max-width bar: default '%lf'.\n", DEFAULT_HI);
	fprintf(o, "-w   {n}          Bar-field width: default '%d'.\n", DEFAULT_WIDTH);
}

// ----------------------------------------------------------------
static mapper_t* mapper_bar_parse_cli(int* pargi, int argc, char** argv) {
	string_array_t*  pfield_names = NULL;
	char*  fill_string  = NULL;
	char*  oob_string   = NULL;
	char*  blank_string = NULL;
	double lo           = DEFAULT_LO;
	double hi           = DEFAULT_HI;
	int    width        = DEFAULT_WIDTH;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_array_flag(pstate, "-f",   &pfield_names);
	ap_define_string_flag(pstate,       "-c",   &fill_string);
	ap_define_string_flag(pstate,       "-x",   &oob_string);
	ap_define_string_flag(pstate,       "-b",   &blank_string);
	ap_define_double_flag(pstate,       "--lo", &lo);
	ap_define_double_flag(pstate,       "--hi", &hi);
	ap_define_int_flag(pstate,          "-w",   &width);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_bar_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (pfield_names == NULL) {
		mapper_bar_usage(stderr, argv[0], verb);
		return NULL;
	}

	char fill_char  = DEFAULT_FILL_CHAR;
	char oob_char   = DEFAULT_OOB_CHAR;
	char blank_char = DEFAULT_BLANK_CHAR;

	if (fill_string != NULL) {
		if (strlen(fill_string) != 1) {
			mapper_bar_usage(stderr, argv[0], verb);
			return NULL;
		}
		fill_char = fill_string[0];
	}

	if (oob_string != NULL) {
		if (strlen(oob_string) != 1) {
			mapper_bar_usage(stderr, argv[0], verb);
			return NULL;
		}
		oob_char = oob_string[0];
	}

	if (blank_string != NULL) {
		if (strlen(blank_string) != 1) {
			mapper_bar_usage(stderr, argv[0], verb);
			return NULL;
		}
		blank_char = blank_string[0];
	}

	return mapper_bar_alloc(pfield_names, fill_char, oob_char, blank_char, lo, hi, width);
}

// ----------------------------------------------------------------
static mapper_t* mapper_bar_alloc(string_array_t* pfield_names,
	char fill_char, char oob_char, char blank_char, double lo, double hi, int width)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_bar_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_bar_state_t));
	pstate->pfield_names = pfield_names;
	pstate->fill_char    = fill_char;
	pstate->oob_char     = oob_char;
	pstate->blank_char   = blank_char;
	pstate->lo           = lo;
	pstate->hi           = hi;
	pstate->width        = width;
	pstate->bars         = mlr_malloc_or_die((pstate->width + 1) * sizeof(char*));
	for (int i = 0; i <= pstate->width; i++) {
		pstate->bars[i] = mlr_malloc_or_die(pstate->width + 1);
		char* bar = pstate->bars[i];
		memset(bar, pstate->blank_char, pstate->width);
		bar[pstate->width] = 0;
		if (i == 0) {
			bar[0] = pstate->oob_char;
		} else if (i < pstate->width) {
			memset(bar, pstate->fill_char, i);
		} else {
			memset(bar, pstate->fill_char, pstate->width);
			bar[pstate->width-1] = pstate->oob_char;
		}
		pstate->bars[i] = bar;
	}

	pmapper->pprocess_func = mapper_bar_process;
	pmapper->pvstate       = (void*)pstate;
	pmapper->pfree_func    = mapper_bar_free;

	return pmapper;
}

static void mapper_bar_free(void* pvstate) {
	mapper_bar_state_t* pstate = (mapper_bar_state_t*)pvstate;
	string_array_free(pstate->pfield_names);
	for (int i = 0; i <= pstate->width; i++)
		free(pstate->bars[i]);
	free(pstate->bars);
}

// ----------------------------------------------------------------
static sllv_t* mapper_bar_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // end of stream
		return sllv_single(NULL);

	mapper_bar_state_t* pstate = (mapper_bar_state_t*)pvstate;
	int n = pstate->pfield_names->length;
	for (int i = 0; i < n; i++) {
		char* name = pstate->pfield_names->strings[i];
		char* sval = lrec_get(pinrec, name);
		if (sval == NULL)
			continue;
		double dval = mlr_double_from_string_or_die(sval);
		int idx = (int)(pstate->width * (dval - pstate->lo) / (pstate->hi - pstate->lo));
		if (idx < 0)
			idx = 0;
		if (idx > pstate->width)
			idx = pstate->width;
		lrec_put_no_free(pinrec, name, pstate->bars[idx]);
	}
	return sllv_single(pinrec);
}
