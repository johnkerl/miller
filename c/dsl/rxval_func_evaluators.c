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

	return pstate->pfunc(pstate->pbxvals, nargs);
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

	return pstate->pfunc(&bxval1);
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

	return pstate->pfunc(&bxval1);
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

	return pstate->pfunc(&bxval1, &bxval2);
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

	return pstate->pfunc(&bxval1, &bxval2);
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
		return box_ephemeral_val(mv_error());
	}
	// xxx to-string ...

	if (!bxval2.xval.is_terminal) {
		return box_ephemeral_val(mv_error());
	}
	// xxx to-string ...

	return pstate->pfunc(&bxval1, &bxval2);
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

	return pstate->pfunc(&bxval1, &bxval2, &bxval3);
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
