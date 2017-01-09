#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "containers/dvector.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define DVECTOR_INITIAL_SIZE 1024

typedef struct _mapper_histogram_state_t {
	ap_state_t* pargp;
	slls_t* value_field_names;
	double lo;
	int    nbins;
	double hi;
	double mul;
	lhmsv_t* pcounts_by_field;
	lhmsv_t* pvectors_by_field; // For auto-mode
	char*  output_prefix;
} mapper_histogram_state_t;

static void      mapper_histogram_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_histogram_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_histogram_alloc(ap_state_t* pargp, slls_t* value_field_names, double lo, int nbins, double hi,
	int do_auto, char* output_prefix);
static void      mapper_histogram_free(mapper_t* pmapper);

static void      mapper_histogram_ingest(lrec_t* pinrec, mapper_histogram_state_t* pstate);
static sllv_t*   mapper_histogram_emit(mapper_histogram_state_t* pstate);
static sllv_t*   mapper_histogram_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

static void      mapper_histogram_ingest_auto(lrec_t* pinrec, mapper_histogram_state_t* pstate);
static sllv_t*   mapper_histogram_emit_auto(mapper_histogram_state_t* pstate);
static sllv_t*   mapper_histogram_process_auto(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_histogram_setup = {
	.verb = "histogram",
	.pusage_func = mapper_histogram_usage,
	.pparse_func = mapper_histogram_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_histogram_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "-f {a,b,c}    Value-field names for histogram counts\n");
	fprintf(o, "--lo {lo}     Histogram low value\n");
	fprintf(o, "--hi {hi}     Histogram high value\n");
	fprintf(o, "--nbins {n}   Number of histogram bins\n");
	fprintf(o, "--auto        Automatically computes limits, ignoring --lo and --hi.\n");
	fprintf(o, "              Holds all values in memory before producing any output.\n");
	fprintf(o, "-o {prefix}   Prefix for output field name. Default: no prefix.\n");
	fprintf(o, "Just a histogram. Input values < lo or > hi are not counted.\n");
}

static mapper_t* mapper_histogram_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	slls_t* value_field_names = NULL;
	double lo = 0.0;
	double hi = 0.0;
	int nbins = 0;
	int do_auto = FALSE;
	char* output_prefix = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &value_field_names);
	ap_define_float_flag(pstate, "--lo",     &lo);
	ap_define_float_flag(pstate, "--hi",     &hi);
	ap_define_int_flag(pstate,   "--nbins",  &nbins);
	ap_define_true_flag(pstate,  "--auto",   &do_auto);
	ap_define_string_flag(pstate,  "-o",     &output_prefix);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_histogram_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (value_field_names == NULL) {
		mapper_histogram_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (nbins == 0) {
		mapper_histogram_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (lo == hi && !do_auto) {
		mapper_histogram_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_histogram_alloc(pstate, value_field_names, lo, nbins, hi, do_auto, output_prefix);
}

// ----------------------------------------------------------------
static mapper_t* mapper_histogram_alloc(ap_state_t* pargp, slls_t* value_field_names,
	double lo, int nbins, double hi, int do_auto, char* output_prefix)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_histogram_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_histogram_state_t));

	pstate->pargp = pargp;
	pstate->value_field_names = value_field_names;
	pstate->nbins = nbins;
	pstate->pcounts_by_field = lhmsv_alloc();
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		unsigned long long* pcounts = mlr_malloc_or_die(nbins * sizeof(unsigned long long));
		for (int i = 0; i < nbins; i++)
			pcounts[i] = 0LL;
		lhmsv_put(pstate->pcounts_by_field, value_field_name, pcounts, NO_FREE);
	}
	if (do_auto) {
		pstate->pvectors_by_field = lhmsv_alloc();
		for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
			char* value_field_name = pe->value;
			dvector_t* pvector = dvector_alloc(DVECTOR_INITIAL_SIZE);
			lhmsv_put(pstate->pvectors_by_field, value_field_name, pvector, NO_FREE);
		}
	} else {
		pstate->pvectors_by_field = NULL;
		pstate->lo  = lo;
		pstate->hi  = hi;
		pstate->mul = nbins / (hi - lo);
	}
	pstate->output_prefix = output_prefix;

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = do_auto ? mapper_histogram_process_auto : mapper_histogram_process;
	pmapper->pfree_func    = mapper_histogram_free;

	return pmapper;
}

static void mapper_histogram_free(mapper_t* pmapper) {
	mapper_histogram_state_t* pstate = pmapper->pvstate;
	slls_free(pstate->value_field_names);
	if (pstate->pcounts_by_field != NULL) {
		for (lhmsve_t* pe = pstate->pcounts_by_field->phead; pe != NULL; pe = pe->pnext) {
			unsigned long long* pcounts = pe->pvvalue;
			free(pcounts);
		}
		lhmsv_free(pstate->pcounts_by_field);
	}
	if (pstate->pvectors_by_field != NULL) {
		for (lhmsve_t* pe = pstate->pvectors_by_field->phead; pe != NULL; pe = pe->pnext) {
			dvector_t* pvector = pe->pvvalue;
			dvector_free(pvector);
		}
		lhmsv_free(pstate->pvectors_by_field);
	}
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_histogram_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_histogram_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_histogram_ingest(pinrec, pstate);
		lrec_free(pinrec);
		return NULL;
	}
	else {
		return mapper_histogram_emit(pstate);
	}
}

static void mapper_histogram_ingest(lrec_t* pinrec, mapper_histogram_state_t* pstate) {
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		char* strv = lrec_get(pinrec, value_field_name);
		unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);
		if (strv != NULL) {
			double val = mlr_double_from_string_or_die(strv);
			if ((val >= pstate->lo) && (val < pstate->hi)) {
				int idx = (int)((val-pstate->lo) * pstate->mul);
				pcounts[idx]++;
			} else if (val == pstate->hi) {
				int idx = pstate->nbins - 1;
				pcounts[idx]++;
			}
		}
	}
}

static sllv_t* mapper_histogram_emit(mapper_histogram_state_t* pstate) {
	sllv_t* poutrecs = sllv_alloc();

	lhmss_t* pcount_field_names = lhmss_alloc();
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		char* count_field_name = (pstate->output_prefix == NULL)
			? mlr_paste_3_strings(value_field_name, "_", "count")
			: mlr_paste_4_strings(pstate->output_prefix, value_field_name, "_", "count");
		lhmss_put(pcount_field_names, mlr_strdup_or_die(value_field_name), count_field_name,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}

	for (int i = 0; i < pstate->nbins; i++) {
		lrec_t* poutrec = lrec_unbacked_alloc();

		char* value = mlr_alloc_string_from_double(pstate->lo + i / pstate->mul, MLR_GLOBALS.ofmt);
		if (pstate->output_prefix == NULL) {
			lrec_put(poutrec, "bin_lo", value, FREE_ENTRY_VALUE);
		} else {
			lrec_put(poutrec, mlr_paste_2_strings(pstate->output_prefix, "bin_lo"), value,
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}

		value = mlr_alloc_string_from_double(pstate->lo + (i+1) / pstate->mul, MLR_GLOBALS.ofmt);
		if (pstate->output_prefix == NULL) {
			lrec_put(poutrec, "bin_hi", value, FREE_ENTRY_VALUE);
		} else {
			lrec_put(poutrec, mlr_paste_2_strings(pstate->output_prefix, "bin_hi"), value,
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}

		for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
			char* value_field_name = pe->value;
			unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);

			char* count_field_name = lhmss_get(pcount_field_names, value_field_name);

			value = mlr_alloc_string_from_ull(pcounts[i]);
			lrec_put(poutrec, mlr_strdup_or_die(count_field_name), value,
				FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		}

		sllv_append(poutrecs, poutrec);
	}

	lhmss_free(pcount_field_names);

	sllv_append(poutrecs, NULL);
	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_histogram_process_auto(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_histogram_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		mapper_histogram_ingest_auto(pinrec, pstate);
		lrec_free(pinrec);
		return NULL;
	}
	else {
		return mapper_histogram_emit_auto(pstate);
	}
}

static void mapper_histogram_ingest_auto(lrec_t* pinrec, mapper_histogram_state_t* pstate) {
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		char* strv = lrec_get(pinrec, value_field_name);
		dvector_t* pvector = lhmsv_get(pstate->pvectors_by_field, value_field_name);
		if (strv != NULL) {
			dvector_append(pvector, mlr_double_from_string_or_die(strv));
		}
	}
}

static sllv_t* mapper_histogram_emit_auto(mapper_histogram_state_t* pstate) {
	int have_lo_hi = FALSE;
	double lo = 0.0, hi = 1.0;
	int nbins = pstate->nbins;

	// Limits pass
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		dvector_t* pvector = lhmsv_get(pstate->pvectors_by_field, value_field_name);
		int n = pvector->size;
		for (int i = 0; i < n; i++) {
			double val = pvector->data[i];
			if (have_lo_hi) {
				if (lo > val)
					lo = val;
				if (hi < val)
					hi = val;
			} else {
				lo = val;
				hi = val;
				have_lo_hi = TRUE;
			}
		}
	}

	// Binning pass
	double mul = nbins / (hi - lo);
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		dvector_t* pvector = lhmsv_get(pstate->pvectors_by_field, value_field_name);
		unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);
		int n = pvector->size;
		for (int i = 0; i < n; i++) {
			double val = pvector->data[i];
			if ((val >= lo) && (val < hi)) {
				int idx = (int)((val-lo) * mul);
				pcounts[idx]++;
			} else if (val == hi) {
				int idx = nbins - 1;
				pcounts[idx]++;
			}
		}
	}

	// Emission pass
	sllv_t* poutrecs = sllv_alloc();
	lhmss_t* pcount_field_names = lhmss_alloc();
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		char* count_field_name = (pstate->output_prefix == NULL)
			? mlr_paste_3_strings(value_field_name, "_", "count")
			: mlr_paste_4_strings(pstate->output_prefix, value_field_name, "_", "count");
		lhmss_put(pcount_field_names, mlr_strdup_or_die(value_field_name), count_field_name,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}

	for (int i = 0; i < nbins; i++) {
		lrec_t* poutrec = lrec_unbacked_alloc();

		char* value = mlr_alloc_string_from_double(lo + i / mul, MLR_GLOBALS.ofmt);
		if (pstate->output_prefix == NULL) {
			lrec_put(poutrec, "bin_lo", value, FREE_ENTRY_VALUE);
		} else {
			lrec_put(poutrec, mlr_paste_2_strings(pstate->output_prefix, "bin_lo"), value,
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		}

		value = mlr_alloc_string_from_double(lo + (i+1) / mul, MLR_GLOBALS.ofmt);
		lrec_put(poutrec, "bin_hi", value, FREE_ENTRY_VALUE);

		for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
			char* value_field_name = pe->value;
			unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);
			char* count_field_name = lhmss_get(pcount_field_names, value_field_name);
			value = mlr_alloc_string_from_ull(pcounts[i]);
			lrec_put(poutrec, mlr_strdup_or_die(count_field_name), value, FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		}

		sllv_append(poutrecs, poutrec);
	}
	sllv_append(poutrecs, NULL);

	lhmss_free(pcount_field_names);
	return poutrecs;
}
