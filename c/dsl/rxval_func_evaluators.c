#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "dsl/context_flags.h"
#include "dsl/rval_evaluators.h"

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_variadic_state_t {
	xv_variadic_func_t*  pfunc;
	rxval_evaluator_t**  pargs;
	boxed_xval_t*        pbxvals;
	int                  nargs;
} rxval_evaluator_variadic_state_t;

static boxed_xval_t rxval_evaluator_variadic_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_variadic_state_t* pstate = pvstate;
	int nargs = pstate->nargs;

	for (int i = 0; i < nargs; i++) {
		rxval_evaluator_t* parg = pstate->pargs[i];
		boxed_xval_t* pbxval = &pstate->pbxvals[i];
		*pbxval = parg->pprocess_func(parg->pvstate, pvars);
		// xxx map-check ...
	}

	boxed_xval_t bxrv = pstate->pfunc(pstate->pbxvals, nargs);

	for (int i = 0; i < nargs; i++) {
		boxed_xval_t* pbxval = &pstate->pbxvals[i];
		if (pbxval->is_ephemeral) {
			mlhmmv_xvalue_free(&pbxval->xval);
		}
	}
	return bxrv;
}

static void rxval_evaluator_variadic_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_variadic_state_t* pstate = pxevaluator->pvstate;

	for (int i = 0; i < pstate->nargs; i++)
		pstate->pargs[i]->pfree_func(pstate->pargs[i]);
	free(pstate->pargs);

	// xxx check ephemeral ...
	for (int i = 0; i < pstate->nargs; i++)
		mlhmmv_xvalue_free(&pstate->pbxvals[i].xval);
	free(pstate->pbxvals);

	free(pstate);

	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_variadic_func(
	xv_variadic_func_t* pfunc,
	rxval_evaluator_t** pargs,
	int nargs)
{
	rxval_evaluator_variadic_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_variadic_state_t));
	pstate->pfunc = pfunc;
	pstate->pargs = pargs;
	pstate->nargs = nargs;
	pstate->pbxvals  = mlr_malloc_or_die(nargs * sizeof(boxed_xval_t));

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_variadic_func;
	pxevaluator->pfree_func    = rxval_evaluator_variadic_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_x_state_t {
	xv_unary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
} rxval_evaluator_x_x_state_t;

static boxed_xval_t rxval_evaluator_x_x_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_x_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	boxed_xval_t bxrv = pstate->pfunc(&bxval1);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}

	return bxrv;
}

static void rxval_evaluator_x_x_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_x_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_x_func(xv_unary_func_t* pfunc, rxval_evaluator_t* parg1) {
	rxval_evaluator_x_x_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_x_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_x_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_x_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_m_state_t {
	xv_unary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
} rxval_evaluator_x_m_state_t;

static boxed_xval_t rxval_evaluator_x_m_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_m_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	if (bxval1.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	boxed_xval_t bxrv = pstate->pfunc(&bxval1);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_m_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_m_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_m_func(xv_unary_func_t* pfunc, rxval_evaluator_t* parg1) {
	rxval_evaluator_x_m_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_m_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_m_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_m_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_mx_state_t {
	xv_binary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
	rxval_evaluator_t* parg2;
} rxval_evaluator_x_mx_state_t;

static boxed_xval_t rxval_evaluator_x_mx_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_mx_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	boxed_xval_t bxval2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);

	if (bxval1.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval2.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	boxed_xval_t bxrv = pstate->pfunc(&bxval1, &bxval2);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	if (bxval2.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval2.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_mx_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_mx_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mx_func(xv_binary_func_t* pfunc,
	rxval_evaluator_t* parg1, rxval_evaluator_t* parg2)
{
	rxval_evaluator_x_mx_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_mx_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_mx_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_mx_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_ms_state_t {
	xv_binary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
	rxval_evaluator_t* parg2;
} rxval_evaluator_x_ms_state_t;

static boxed_xval_t rxval_evaluator_x_ms_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_ms_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	boxed_xval_t bxval2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);

	if (bxval1.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval2.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	// xxx to-string ...

	boxed_xval_t bxrv = pstate->pfunc(&bxval1, &bxval2);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	if (bxval2.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval2.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_ms_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_ms_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ms_func(xv_binary_func_t* pfunc,
	rxval_evaluator_t* parg1, rxval_evaluator_t* parg2)
{
	rxval_evaluator_x_ms_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_ms_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_ms_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_ms_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_ss_state_t {
	xv_binary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
	rxval_evaluator_t* parg2;
} rxval_evaluator_x_ss_state_t;

static boxed_xval_t rxval_evaluator_x_ss_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_ss_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	boxed_xval_t bxval2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);

	if (!bxval1.xval.is_terminal) {
		mv_free(&bxval1.xval.terminal_mlrval);
		return box_ephemeral_val(mv_error());
	}
	// xxx to-string ...
	if (!mv_is_string_or_empty(&bxval1.xval.terminal_mlrval)) {
		mv_free(&bxval1.xval.terminal_mlrval);
		return box_ephemeral_val(mv_error());
	}

	if (!bxval2.xval.is_terminal) {
		mv_free(&bxval2.xval.terminal_mlrval);
		return box_ephemeral_val(mv_error());
	}
	// xxx to-string ...
	if (!mv_is_string_or_empty(&bxval2.xval.terminal_mlrval)) {
		mv_free(&bxval2.xval.terminal_mlrval);
		return box_ephemeral_val(mv_error());
	}

	boxed_xval_t bxrv = pstate->pfunc(&bxval1, &bxval2);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	if (bxval2.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval2.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_ss_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_ss_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ss_func(xv_binary_func_t* pfunc,
	rxval_evaluator_t* parg1, rxval_evaluator_t* parg2)
{
	rxval_evaluator_x_ss_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_ss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_ss_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_ss_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_mss_state_t {
	xv_ternary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
	rxval_evaluator_t* parg2;
	rxval_evaluator_t* parg3;
} rxval_evaluator_x_mss_state_t;

static boxed_xval_t rxval_evaluator_x_mss_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_mss_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	boxed_xval_t bxval2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	boxed_xval_t bxval3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);

	if (bxval1.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval2.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval3.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	// xxx to-string ...

	boxed_xval_t bxrv = pstate->pfunc(&bxval1, &bxval2, &bxval3);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	if (bxval2.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval2.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_mss_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_mss_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mss_func(xv_ternary_func_t* pfunc,
	rxval_evaluator_t* parg1, rxval_evaluator_t* parg2, rxval_evaluator_t* parg3)
{
	rxval_evaluator_x_mss_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_mss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_mss_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_mss_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
typedef struct _rxval_evaluator_x_sss_state_t {
	xv_ternary_func_t*  pfunc;
	rxval_evaluator_t* parg1;
	rxval_evaluator_t* parg2;
	rxval_evaluator_t* parg3;
} rxval_evaluator_x_sss_state_t;

static boxed_xval_t rxval_evaluator_x_sss_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_x_sss_state_t* pstate = pvstate;
	boxed_xval_t bxval1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	boxed_xval_t bxval2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	boxed_xval_t bxval3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);

	if (!bxval1.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval2.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	if (!bxval3.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}

	// xxx to-string ...

	boxed_xval_t bxrv = pstate->pfunc(&bxval1, &bxval2, &bxval3);

	if (bxval1.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval1.xval);
	}
	if (bxval2.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval2.xval);
	}
	if (bxval3.is_ephemeral) { // xxx funcify as bxval_free_if_ephemeral or some such
		mlhmmv_xvalue_free(&bxval3.xval);
	}
	return bxrv;
}

static void rxval_evaluator_x_sss_free(rxval_evaluator_t* pxevaluator) {
	rxval_evaluator_x_sss_state_t* pstate = pxevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pxevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_x_sss_func(xv_ternary_func_t* pfunc,
	rxval_evaluator_t* parg1, rxval_evaluator_t* parg2, rxval_evaluator_t* parg3)
{
	rxval_evaluator_x_sss_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_x_sss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rxval_evaluator_t* pxevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pxevaluator->pvstate       = pstate;
	pxevaluator->pprocess_func = rxval_evaluator_x_sss_func;
	pxevaluator->pfree_func    = rxval_evaluator_x_sss_free;

	return pxevaluator;
}

// ----------------------------------------------------------------
// Does type-check assertion on the argument and returns it unmodified if the
// test passes.  Else throws an error.
typedef struct _rxval_evaluator_A_x_state_t {
	xv_unary_func_t*   pfunc;
	rxval_evaluator_t* parg1;
	char*              desc;
} rxval_evaluator_A_x_state_t;

static boxed_xval_t rxval_evaluator_A_x_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_A_x_state_t* pstate = pvstate;
	boxed_xval_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	boxed_xval_t ok = pstate->pfunc(&val1);

	if (!ok.xval.terminal_mlrval.u.boolv) {
		fprintf(stderr, "%s: %s type-assertion failed at NR=%lld FNR=%lld FILENAME=%s\n",
			MLR_GLOBALS.bargv0, pstate->desc, pvars->pctx->nr, pvars->pctx->fnr, pvars->pctx->filename);
		exit(1);
	}

	return val1;
}
static void rxval_evaluator_A_x_free(rxval_evaluator_t* pevaluator) {
	rxval_evaluator_A_x_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate->desc);
	free(pstate);
	free(pevaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_A_x_func(xv_unary_func_t* pfunc, rxval_evaluator_t* parg1, char* desc) {
	rxval_evaluator_A_x_state_t* pstate = mlr_malloc_or_die(sizeof(rxval_evaluator_A_x_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->desc  = mlr_strdup_or_die(desc);

	rxval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rxval_evaluator_A_x_func;
	pevaluator->pfree_func = rxval_evaluator_A_x_free;

	return pevaluator;
}
