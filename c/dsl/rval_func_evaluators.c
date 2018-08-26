#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "dsl/context_flags.h"
#include "dsl/rval_evaluators.h"

// ----------------------------------------------------------------
typedef struct _rval_evaluator_variadic_state_t {
	mv_variadic_func_t* pfunc;
	rval_evaluator_t**  pargs;
	mv_t*               pmvs;
	int                 nargs;
} rval_evaluator_variadic_state_t;

static mv_t rval_evaluator_variadic_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_variadic_state_t* pstate = pvstate;
	int nargs = pstate->nargs;

	for (int i = 0; i < nargs; i++) {
		rval_evaluator_t* parg = pstate->pargs[i];
		mv_t* pmv = &pstate->pmvs[i];
		*pmv = parg->pprocess_func(parg->pvstate, pvars);
	}

	return pstate->pfunc(pstate->pmvs, nargs);
}

static void rval_evaluator_variadic_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_variadic_state_t* pstate = pevaluator->pvstate;

	for (int i = 0; i < pstate->nargs; i++)
		pstate->pargs[i]->pfree_func(pstate->pargs[i]);
	free(pstate->pargs);
	// contents already mv_freed by evaluator chains at process time
	free(pstate->pmvs);

	free(pstate);

	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_variadic_func(mv_variadic_func_t* pfunc, rval_evaluator_t** pargs, int nargs) {
	rval_evaluator_variadic_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_variadic_state_t));
	pstate->pfunc = pfunc;
	pstate->pargs = pargs;
	pstate->nargs = nargs;
	pstate->pmvs  = mlr_malloc_or_die(nargs * sizeof(mv_t));

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate       = pstate;
	pevaluator->pprocess_func = rval_evaluator_variadic_func;
	pevaluator->pfree_func    = rval_evaluator_variadic_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_b_b_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_b_b_state_t;

static mv_t rval_evaluator_b_b_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_b_b_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOLEAN)
		return mv_error();

	return pstate->pfunc(&val1);
}
static void rval_evaluator_b_b_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_b_b_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_b_b_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_b_b_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_b_b_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_b_b_func;
	pevaluator->pfree_func = rval_evaluator_b_b_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_b_bb_state_t {
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_b_bb_state_t;

// This is different from most of the lrec-evaluator functions in that it does short-circuiting:
// since is logical AND, the LHS is not evaluated if the RHS is false.
static mv_t rval_evaluator_b_bb_and_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type == MT_BOOLEAN) {
		if (val1.u.boolv == FALSE)
			return val1;
	} else if (val1.type != MT_ABSENT) {
		return mv_error();
	}

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type == MT_BOOLEAN) {
		return val2;
	} else if (val2.type == MT_ABSENT) {
		return val1;
	} else {
		return mv_error();
	}
}

// This is different from most of the lrec-evaluator functions in that it does short-circuiting:
// since is logical OR, the LHS is not evaluated if the RHS is true.
static mv_t rval_evaluator_b_bb_or_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type == MT_BOOLEAN) {
		if (val1.u.boolv == TRUE)
			return val1;
	} else if (val1.type != MT_ABSENT) {
		return mv_error();
	}

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type == MT_BOOLEAN) {
		return val2;
	} else if (val2.type == MT_ABSENT) {
		return val1;
	} else {
		return mv_error();
	}
}

static mv_t rval_evaluator_b_bb_xor_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOLEAN && val1.type != MT_ABSENT) {
		return mv_error();
	}

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type == MT_BOOLEAN) {
		if (val1.type == MT_BOOLEAN)
			return mv_from_bool(val1.u.boolv ^ val2.u.boolv);
		else
			return val2;
	} else if (val2.type == MT_ABSENT) {
		return val1;
	} else {
		return mv_error();
	}
}

static void rval_evaluator_b_bb_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_b_bb_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_b_bb_and_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2) {
	rval_evaluator_b_bb_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_b_bb_state_t));
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_b_bb_and_func;
	pevaluator->pfree_func = rval_evaluator_b_bb_free;

	return pevaluator;
}

rval_evaluator_t* rval_evaluator_alloc_from_b_bb_or_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2) {
	rval_evaluator_b_bb_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_b_bb_state_t));
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_b_bb_or_func;
	pevaluator->pfree_func = rval_evaluator_b_bb_free;

	return pevaluator;
}

rval_evaluator_t* rval_evaluator_alloc_from_b_bb_xor_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2) {
	rval_evaluator_b_bb_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_b_bb_state_t));
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_b_bb_xor_func;
	pevaluator->pfree_func = rval_evaluator_b_bb_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_z_state_t {
	mv_zary_func_t* pfunc;
} rval_evaluator_x_z_state_t;

static mv_t rval_evaluator_x_z_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_z_state_t* pstate = pvstate;
	return pstate->pfunc();
}
static void rval_evaluator_x_z_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_z_func(mv_zary_func_t* pfunc) {
	rval_evaluator_x_z_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_z_state_t));
	pstate->pfunc = pfunc;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_z_func;
	pevaluator->pfree_func = rval_evaluator_x_z_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_f_f_state_t {
	mv_unary_func_t* pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_f_f_state_t;

static mv_t rval_evaluator_f_f_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_f_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_FLOAT)
		return mv_error();

	return pstate->pfunc(&val1);
}
static void rval_evaluator_f_f_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_f_f_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_f_f_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_f_f_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_f_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_f_f_func;
	pevaluator->pfree_func = rval_evaluator_f_f_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_n_state_t {
	mv_unary_func_t* pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_x_n_state_t;

static mv_t rval_evaluator_x_n_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_n_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	return pstate->pfunc(&val1);
}
static void rval_evaluator_x_n_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_n_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_n_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_x_n_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_n_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_n_func;
	pevaluator->pfree_func = rval_evaluator_x_n_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_i_i_state_t {
	mv_unary_func_t* pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_i_i_state_t;

static mv_t rval_evaluator_i_i_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_i_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	return pstate->pfunc(&val1);
}
static void rval_evaluator_i_i_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_i_i_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_i_i_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_i_i_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_i_i_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_i_i_func;
	pevaluator->pfree_func = rval_evaluator_i_i_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_f_ff_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_f_ff_state_t;

static mv_t rval_evaluator_f_ff_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_float_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_f_ff_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_f_ff_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_f_ff_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_f_ff_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_f_ff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_f_ff_func;
	pevaluator->pfree_func = rval_evaluator_f_ff_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_xx_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_xx_state_t;

static mv_t rval_evaluator_x_xx_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_xx_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);

	// nullities handled by full disposition matrices
	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_x_xx_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_xx_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_xx_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_xx_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_xx_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_xx_func;
	pevaluator->pfree_func = rval_evaluator_x_xx_free;

	return pevaluator;
}

// ----------------------------------------------------------------
// This is for min/max which can return non-null when one argument is null --
// in comparison to other functions which return null if *any* argument is
// null.

typedef struct _rval_evaluator_x_xx_nullable_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_xx_nullable_state_t;

static mv_t rval_evaluator_x_xx_nullable_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_xx_nullable_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_number_nullable(&val1);

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_number_nullable(&val2);

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_x_xx_nullable_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_xx_nullable_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_xx_nullable_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_xx_nullable_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_xx_nullable_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_xx_nullable_func;
	pevaluator->pfree_func = rval_evaluator_x_xx_nullable_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_f_fff_state_t {
	mv_ternary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
	rval_evaluator_t* parg3;
} rval_evaluator_f_fff_state_t;

static mv_t rval_evaluator_f_fff_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_f_fff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_float_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);

	mv_t val3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
	mv_set_float_nullable(&val3);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val3);

	return pstate->pfunc(&val1, &val2, &val3);
}
static void rval_evaluator_f_fff_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_f_fff_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_f_fff_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	rval_evaluator_f_fff_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_f_fff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_f_fff_func;
	pevaluator->pfree_func = rval_evaluator_f_fff_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_i_ii_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_i_ii_state_t;

static mv_t rval_evaluator_i_ii_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_i_ii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_INT)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_int_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_INT)
		return mv_error();

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_i_ii_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_i_ii_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_i_ii_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_i_ii_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_i_ii_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_i_ii_func;
	pevaluator->pfree_func = rval_evaluator_i_ii_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_i_iii_state_t {
	mv_ternary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
	rval_evaluator_t* parg3;
} rval_evaluator_i_iii_state_t;

static mv_t rval_evaluator_i_iii_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_i_iii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_INT)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_int_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_INT)
		return mv_error();

	mv_t val3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
	mv_set_int_nullable(&val3);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val3);
	if (val3.type != MT_INT)
		return mv_error();

	return pstate->pfunc(&val1, &val2, &val3);
}
static void rval_evaluator_i_iii_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_i_iii_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_i_iii_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	rval_evaluator_i_iii_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_i_iii_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_i_iii_func;
	pevaluator->pfree_func = rval_evaluator_i_iii_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_ternop_state_t {
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
	rval_evaluator_t* parg3;
} rval_evaluator_ternop_state_t;

static mv_t rval_evaluator_ternop_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_ternop_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	mv_set_boolean_strict(&val1);

	return val1.u.boolv
		? pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars)
		: pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
}
static void rval_evaluator_ternop_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_ternop_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_ternop(rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	rval_evaluator_ternop_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_ternop_state_t));
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_ternop_func;
	pevaluator->pfree_func = rval_evaluator_ternop_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_s_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_s_s_state_t;

static mv_t rval_evaluator_s_s_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	return pstate->pfunc(&val1);
}
static void rval_evaluator_s_s_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_s_s_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_s_func;
	pevaluator->pfree_func = rval_evaluator_s_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_sii_state_t {
	mv_ternary_func_t*  pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
	rval_evaluator_t* parg3;
} rval_evaluator_s_sii_state_t;

static mv_t rval_evaluator_s_sii_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_sii_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	mv_set_int_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);

	mv_t val3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
	mv_set_int_nullable(&val3);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val3);

	return pstate->pfunc(&val1, &val2, &val3);
}
static void rval_evaluator_s_sii_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_sii_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_sii_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	rval_evaluator_s_sii_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_sii_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_sii_func;
	pevaluator->pfree_func = rval_evaluator_s_sii_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_f_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_s_f_state_t;

static mv_t rval_evaluator_s_f_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	return pstate->pfunc(&val1);
}
static void rval_evaluator_s_f_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_f_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_s_f_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_f_func;
	pevaluator->pfree_func = rval_evaluator_s_f_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_i_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_s_i_state_t;

static mv_t rval_evaluator_s_i_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	return pstate->pfunc(&val1);
}
static void rval_evaluator_s_i_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_i_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_i_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_s_i_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_i_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_i_func;
	pevaluator->pfree_func = rval_evaluator_s_i_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_f_s_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_f_s_state_t;

static mv_t rval_evaluator_f_s_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_f_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	return pstate->pfunc(&val1);
}
static void rval_evaluator_f_s_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_f_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_f_s_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_f_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_f_s_func;
	pevaluator->pfree_func = rval_evaluator_f_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_i_s_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_i_s_state_t;

static mv_t rval_evaluator_i_s_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_i_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	return pstate->pfunc(&val1);
}
static void rval_evaluator_i_s_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_i_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_i_s_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_i_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_i_s_func;
	pevaluator->pfree_func = rval_evaluator_i_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_x_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_x_x_state_t;

static mv_t rval_evaluator_x_x_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_x_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	// nullity handled by full disposition matrices
	return pstate->pfunc(&val1);
}
static void rval_evaluator_x_x_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_x_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_x_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1) {
	rval_evaluator_x_x_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_x_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_x_func;
	pevaluator->pfree_func = rval_evaluator_x_x_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_xi_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_xi_state_t;

static mv_t rval_evaluator_x_xi_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_xi_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);

	if (val2.type != MT_INT) {
		return mv_error();
	}

	// nullity of 1st argument handled by full disposition matrices
	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_x_xi_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_xi_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_xi_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_xi_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_xi_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_xi_func;
	pevaluator->pfree_func = rval_evaluator_x_xi_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_ns_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_ns_state_t;

static mv_t rval_evaluator_x_ns_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_ns_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (!mv_is_string_or_empty(&val2)) {
		mv_free(&val1);
		return mv_error();
	}

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_x_ns_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_ns_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_ns_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_ns_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_ns_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_ns_func;
	pevaluator->pfree_func = rval_evaluator_x_ns_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_ss_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_ss_state_t;

static mv_t rval_evaluator_x_ss_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_ss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (!mv_is_string_or_empty(&val2))
		return mv_error();

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_x_ss_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_ss_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_ss_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_ss_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_ss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_ss_func;
	pevaluator->pfree_func = rval_evaluator_x_ss_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_ssc_state_t {
	mv_binary_arg3_capture_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_ssc_state_t;

static mv_t rval_evaluator_x_ssc_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_ssc_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (!mv_is_string_or_empty(&val2))
		return mv_error();
	return pstate->pfunc(&val1, &val2, pvars->ppregex_captures);
}
static void rval_evaluator_x_ssc_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_ssc_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_ssc_func(mv_binary_arg3_capture_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_x_ssc_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_ssc_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_ssc_func;
	pevaluator->pfree_func = rval_evaluator_x_ssc_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_sr_state_t {
	mv_binary_arg2_regex_func_t* pfunc;
	rval_evaluator_t*             parg1;
	regex_t                       regex;
	string_builder_t*             psb;
} rval_evaluator_x_sr_state_t;

static mv_t rval_evaluator_x_sr_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_sr_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	return pstate->pfunc(&val1, &pstate->regex, pstate->psb, pvars->ppregex_captures);
}
static void rval_evaluator_x_sr_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_sr_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	regfree(&pstate->regex);
	sb_free(pstate->psb);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_sr_func(mv_binary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	rval_evaluator_x_sr_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_sr_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	int cflags = ignore_case ? REG_ICASE : 0;
	regcomp_or_die(&pstate->regex, regex_string, cflags);
	pstate->psb = sb_alloc(MV_SB_ALLOC_LENGTH);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_sr_func;
	pevaluator->pfree_func = rval_evaluator_x_sr_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_se_state_t {
	mv_binary_arg2_regex_extract_func_t* pfunc;
	rval_evaluator_t*             parg1;
	regex_t                       regex;
} rval_evaluator_x_se_state_t;

static mv_t rval_evaluator_x_se_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_se_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);

	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	return pstate->pfunc(&val1, &pstate->regex);
}

static void rval_evaluator_x_se_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_se_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	regfree(&pstate->regex);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_se_func(mv_binary_arg2_regex_extract_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	rval_evaluator_x_se_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_se_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	int cflags = ignore_case ? REG_ICASE : 0;
	regcomp_or_die(&pstate->regex, regex_string, cflags);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_se_func;
	pevaluator->pfree_func = rval_evaluator_x_se_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_xs_state_t {
	mv_binary_func_t*  pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_s_xs_state_t;

static mv_t rval_evaluator_s_xs_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_xs_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (!mv_is_string_or_empty(&val2))
		return mv_error();

	return pstate->pfunc(&val1, &val2);
}
static void rval_evaluator_s_xs_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_xs_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_xs_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	rval_evaluator_s_xs_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_xs_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_xs_func;
	pevaluator->pfree_func = rval_evaluator_s_xs_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_s_sss_state_t {
	mv_ternary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
	rval_evaluator_t* parg3;
} rval_evaluator_s_sss_state_t;

static mv_t rval_evaluator_s_sss_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_s_sss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(pstate->parg2->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (!mv_is_string_or_empty(&val2)) {
		mv_free(&val1);
		return mv_error();
	}

	mv_t val3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val3);
	if (!mv_is_string_or_empty(&val3)) {
		mv_free(&val1);
		mv_free(&val2);
		return mv_error();
	}

	return pstate->pfunc(&val1, &val2, &val3);
}
static void rval_evaluator_s_sss_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_s_sss_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_s_sss_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	rval_evaluator_s_sss_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_s_sss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_s_sss_func;
	pevaluator->pfree_func = rval_evaluator_s_sss_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_x_srs_state_t {
	mv_ternary_arg2_regex_func_t* pfunc;
	rval_evaluator_t*             parg1;
	regex_t                       regex;
	rval_evaluator_t*             parg3;
	string_builder_t*             psb;
} rval_evaluator_x_srs_state_t;

static mv_t rval_evaluator_x_srs_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_x_srs_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(pstate->parg1->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (!mv_is_string_or_empty(&val1))
		return mv_error();

	mv_t val3 = pstate->parg3->pprocess_func(pstate->parg3->pvstate, pvars);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val3);
	if (!mv_is_string_or_empty(&val3)) {
		mv_free(&val3);
		return mv_error();
	}

	return pstate->pfunc(&val1, &pstate->regex, pstate->psb, &val3);
}
static void rval_evaluator_x_srs_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_x_srs_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	regfree(&pstate->regex);
	pstate->parg3->pfree_func(pstate->parg3);
	sb_free(pstate->psb);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_x_srs_func(mv_ternary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3)
{
	rval_evaluator_x_srs_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_x_srs_state_t));
	pstate->pfunc = pfunc;

	pstate->parg1 = parg1;

	int cflags = ignore_case ? REG_ICASE : 0;
	regcomp_or_die(&pstate->regex, regex_string, cflags);
	pstate->psb = sb_alloc(MV_SB_ALLOC_LENGTH);

	pstate->parg3 = parg3;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = rval_evaluator_x_srs_func;
	pevaluator->pfree_func = rval_evaluator_x_srs_free;

	return pevaluator;
}
