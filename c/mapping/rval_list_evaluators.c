#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/rval_evaluators.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

sllmv_t* evaluate_list(sllv_t* pevaluators, variables_t* pvars, int* pall_non_null_or_error) {
	sllmv_t* pmvs = sllmv_alloc();
	int all_non_null_or_error = TRUE;
	for (sllve_t* pe = pevaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pevaluator = pe->pvvalue;
		mv_t mv = pevaluator->pprocess_func(pevaluator->pvstate, pvars);
		if (mv_is_null_or_error(&mv)) {
			all_non_null_or_error = FALSE;
			break;
		}
		// Don't free the mlrval since its memory will be managed by the sllmv.
		sllmv_add_with_free(pmvs, &mv);
	}

	*pall_non_null_or_error = all_non_null_or_error;
	return pmvs;
}
