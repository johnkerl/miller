#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
typedef struct _mapper_histogram_state_t {
	slls_t* value_field_names;
	double lo;
	int    nbins;
	double hi;
	double mul;
	lhmsv_t* pcounts_by_field;
} mapper_histogram_state_t;

// ----------------------------------------------------------------
sllv_t* mapper_histogram_func(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_histogram_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
			char* value_field_name = pe->value;
			char* sval = lrec_get(pinrec, value_field_name);
			unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);
			if (sval != NULL) {
				double val = mlr_double_from_string_or_die(sval);
				if ((val >= pstate->lo) && (val < pstate->hi)) {
					int idx = (int)((val-pstate->lo) * pstate->mul);
					pcounts[idx]++;
				} else if (val == pstate->hi) {
					int idx = pstate->nbins - 1;
					pcounts[idx]++;
				}
			}
		}
		lrec_free(pinrec);

		return NULL;
	}
	else {
		sllv_t* poutrecs = sllv_alloc();

		lhmss_t* pcount_field_names = lhmss_alloc();
		for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
			char* value_field_name = pe->value;
			char* count_field_name = mlr_paste_3_strings(value_field_name, "_", "count");
			lhmss_put(pcount_field_names, value_field_name, count_field_name);
		}

		for (int i = 0; i < pstate->nbins; i++) {
			lrec_t* poutrec = lrec_unbacked_alloc();

			char* value = mlr_alloc_string_from_double(pstate->lo + i / pstate->mul, pctx->statx.ofmt);
			lrec_put(poutrec, "bin_lo", value, LREC_FREE_ENTRY_VALUE);
			free(value);

			value = mlr_alloc_string_from_double(pstate->lo + (i+1) / pstate->mul, pctx->statx.ofmt);
			lrec_put(poutrec, "bin_hi", value, LREC_FREE_ENTRY_VALUE);
			free(value);

			for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
				char* value_field_name = pe->value;
				unsigned long long* pcounts = lhmsv_get(pstate->pcounts_by_field, value_field_name);

				char* count_field_name = lhmss_get(pcount_field_names, value_field_name);

				value = mlr_alloc_string_from_ull(pcounts[i]);
				lrec_put(poutrec, count_field_name, value, LREC_FREE_ENTRY_VALUE);
				free(value);
			}

			sllv_add(poutrecs, poutrec);
		}

		lhmss_free(pcount_field_names);

		sllv_add(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static void mapper_histogram_free(void* pvstate) {
	mapper_histogram_state_t* pstate = (mapper_histogram_state_t*)pvstate;
	if (pstate->value_field_names != NULL)
		slls_free(pstate->value_field_names);
	if (pstate->pcounts_by_field != NULL)
		lhmsv_free(pstate->pcounts_by_field);
}

mapper_t* mapper_histogram_alloc(slls_t* value_field_names, double lo, int nbins, double hi) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_histogram_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_histogram_state_t));

	pstate->value_field_names = slls_copy(value_field_names);
	pstate->lo    = lo;
	pstate->nbins = nbins;
	pstate->hi    = hi;
	pstate->mul   = nbins / (hi - lo);
	pstate->pcounts_by_field = lhmsv_alloc();
	for (sllse_t* pe = pstate->value_field_names->phead; pe != NULL; pe = pe->pnext) {
		char* value_field_name = pe->value;
		unsigned long long* pcounts = mlr_malloc_or_die(nbins * sizeof(unsigned long long));
		for (int i = 0; i < nbins; i++)
			pcounts[i] = 0LL;
		lhmsv_put(pstate->pcounts_by_field, value_field_name, pcounts);
	}

	pmapper->pvstate              = pstate;
	pmapper->pmapper_process_func = mapper_histogram_func;
	pmapper->pmapper_free_func    = mapper_histogram_free;

	return pmapper;
}

// ----------------------------------------------------------------
void mapper_histogram_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-f {a,b,c}    Value-field names for histogram counts\n");
	fprintf(stdout, "--lo {lo}     Histogram low value\n");
	fprintf(stdout, "--hi {hi}     Histogram high value\n");
	fprintf(stdout, "--nbins {n}   Number of histogram bins\n");
}

mapper_t* mapper_histogram_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pvalue_field_names = NULL;
	double lo = 0.0;
	double hi = 0.0;
	int nbins = 0;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pvalue_field_names);
	ap_define_double_flag(pstate, "--lo", &lo);
	ap_define_double_flag(pstate, "--hi", &hi);
	ap_define_int_flag(pstate, "--nbins", &nbins);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_histogram_usage(argv[0], verb);
		return NULL;
	}

	if (pvalue_field_names == NULL) {
		mapper_histogram_usage(argv[0], verb);
		return NULL;
	}

	if ((lo == hi) || (nbins == 0)) {
		mapper_histogram_usage(argv[0], verb);
		return NULL;
	}

	return mapper_histogram_alloc(pvalue_field_names, lo, nbins, hi);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_histogram_setup = {
	.verb = "histogram",
	.pusage_func = mapper_histogram_usage,
	.pparse_func = mapper_histogram_parse_cli,
};
