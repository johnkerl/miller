#include "lib/mlrutil.h"
#include "lib/string_builder.h"
#include "containers/lrec.h"
#include "containers/string_array.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

typedef struct _mapper_bar_state_t {
	string_array_t*  pfield_names;
	char    fill_char;
	char    oob_char;
	char    blank_char;
	double  lo;
	double  hi;
	int     width;
	char**  bars;
	sllv_t* precords; // only for auto mode
} mapper_bar_state_t;

static sllv_t*   mapper_bar_process_no_auto(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_bar_process_auto(lrec_t* pinrec, context_t* pctx, void* pvstate);
static void      mapper_bar_free(void* pvstate);
static mapper_t* mapper_bar_alloc(string_array_t* pfield_names,
	char fill_char, char oob_char, char blank_char, double lo, double hi,
	int width, int do_auto);
static void      mapper_bar_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_bar_parse_cli(int* pargi, int argc, char** argv);

#define DEFAULT_FILL_CHAR  '*'
#define DEFAULT_OOB_CHAR   '#'
#define DEFAULT_BLANK_CHAR '.'
#define DEFAULT_LO         0.0
#define DEFAULT_HI         100.0
#define DEFAULT_WIDTH      40

#define SB_ALLOC_LENGTH    128

// ----------------------------------------------------------------
mapper_setup_t mapper_bar_setup = {
	.verb = "bar",
	.pusage_func = mapper_bar_usage,
	.pparse_func = mapper_bar_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_bar_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Replaces a numeric field with a number of asterisks, allowing for cheesy\n");
	fprintf(o, "bar plots. These align best with --opprint or --oxtab output format.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-f   {a,b,c}      Field names to convert to bars.\n");
	fprintf(o, "-c   {character}  Fill character: default '%c'.\n", DEFAULT_FILL_CHAR);
	fprintf(o, "-x   {character}  Out-of-bounds character: default '%c'.\n", DEFAULT_OOB_CHAR);
	fprintf(o, "-b   {character}  Blank character: default '%c'.\n", DEFAULT_BLANK_CHAR);
	fprintf(o, "--lo {lo}         Lower-limit value for min-width bar: default '%lf'.\n", DEFAULT_LO);
	fprintf(o, "--hi {hi}         Upper-limit value for max-width bar: default '%lf'.\n", DEFAULT_HI);
	fprintf(o, "-w   {n}          Bar-field width: default '%d'.\n", DEFAULT_WIDTH);
	fprintf(o, "--auto            Automatically computes limits, ignoring --lo and --hi.\n");
	fprintf(o, "                  Holds all records in memory before producing any output.\n");
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
	int    do_auto      = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_array_flag(pstate, "-f",     &pfield_names);
	ap_define_string_flag(pstate,       "-c",     &fill_string);
	ap_define_string_flag(pstate,       "-x",     &oob_string);
	ap_define_string_flag(pstate,       "-b",     &blank_string);
	ap_define_float_flag(pstate,        "--lo",   &lo);
	ap_define_float_flag(pstate,        "--hi",   &hi);
	ap_define_int_flag(pstate,          "-w",     &width);
	ap_define_true_flag(pstate,         "--auto", &do_auto);

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

	return mapper_bar_alloc(pfield_names, fill_char, oob_char, blank_char,
		lo, hi, width, do_auto);
}

// ----------------------------------------------------------------
static mapper_t* mapper_bar_alloc(string_array_t* pfield_names,
	char fill_char, char oob_char, char blank_char, double lo, double hi,
	int width, int do_auto)
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
	pstate->precords = do_auto ? sllv_alloc() : NULL;

	pmapper->pprocess_func = do_auto
		? mapper_bar_process_auto
		: mapper_bar_process_no_auto;
	pmapper->pvstate    = (void*)pstate;
	pmapper->pfree_func = mapper_bar_free;

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
static sllv_t* mapper_bar_process_no_auto(lrec_t* pinrec, context_t* pctx, void* pvstate) {
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
		lrec_put(pinrec, name, pstate->bars[idx], NO_FREE);
	}
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_bar_process_auto(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_bar_state_t* pstate = (mapper_bar_state_t*)pvstate;

	if (pinrec != NULL) { // not end of stream
		sllv_add(pstate->precords, pinrec);
		return NULL;
	}

	// end of stream
	int n = pstate->pfield_names->length;
	string_builder_t* psb = sb_alloc(SB_ALLOC_LENGTH);

	// Loop over field names to be barred
	for (int i = 0; i < n; i++) {
		char* name = pstate->pfield_names->strings[i];
		double lo = 0.0;
		double hi = 0.0;

		// First pass computes lo and hi from the data
		int j = 0;
		for (sllve_t* pe = pstate->precords->phead; pe != NULL; pe = pe->pnext, j++) {
			lrec_t* prec = pe->pvdata;
			char* sval = lrec_get(prec, name);
			if (sval == NULL)
				continue;
			double dval = mlr_double_from_string_or_die(sval);
			if (j == 0 || dval < lo)
				lo = dval;
			if (j == 0 || dval > hi)
				hi = dval;
		}

		// Second pass applies the bars. There is some redundant computation
		// which could be hoisted out of the loop for performance ... but this
		// verb computes data solely for visual inspection and I take the
		// nominal use case to be tens or hundreds of records. So, optimization
		// isn't worth the effort here.
		char* slo = mlr_alloc_string_from_double(lo, "%g");
		char* shi = mlr_alloc_string_from_double(hi, "%g");

		for (sllve_t* pe = pstate->precords->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* prec = pe->pvdata;
			char* sval = lrec_get(prec, name);
			if (sval == NULL)
				continue;
			double dval = mlr_double_from_string_or_die(sval);

			int idx = (int)(pstate->width * (dval - lo) / (hi - lo));
			if (idx < 0)
				idx = 0;
			if (idx > pstate->width)
				idx = pstate->width;
			sb_append_string(psb, "[");
			sb_append_string(psb, slo);
			sb_append_string(psb, "]");
			sb_append_string(psb, pstate->bars[idx]);
			sb_append_string(psb, "[");
			sb_append_string(psb, shi);
			sb_append_string(psb, "]");
			lrec_put(prec, name, sb_finish(psb), FREE_ENTRY_VALUE);
		}

	}

	sb_free(psb);
	sllv_add(pstate->precords, NULL);
	return pstate->precords;
}
