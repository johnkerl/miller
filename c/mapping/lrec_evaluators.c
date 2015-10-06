#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/lrec_evaluators.h"

// ================================================================
// NOTES:
//
// * This is used by mlr filter and mlr put.
//
// * Unlike most files in Miller which are read top-down (with sufficient
//   static prototypes at the top of the file to keep the compiler happy),
//   please read this one from the bottom up.
//
// * Comparison to mlr_val.c: the latter is functions from mlr_val(s) to
//   mlr_val; in this file we have the higher-level notion of evaluating lrec
//   objects, using mlr_val.c to do so.
//
// * There are two kinds of lrec-evaluators here: those with _x_ in their names
//   which accept various types of mlr_val, with disposition-matrices in
//   mlr_val.c functions, and those with _i_/_f_/_b_/_s_ (int, float, boolean,
//   string) which either type-check or type-coerce their arguments, invoking
//   type-specific functions in mlr_val.c.  In either case it's the job of
//   lrec_evaluators.c to invoke functions here with mlr_vals of the correct
//   type(s).
// ================================================================

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_b_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_b_b_state_t;

mv_t lrec_evaluator_b_b_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_b_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_b_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_b_b_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_b_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_b_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_bb_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_bb_state_t;

mv_t lrec_evaluator_b_bb_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_bb_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);

	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_bb_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_b_bb_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_bb_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_bb_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_z_state_t {
	mv_zary_func_t* pfunc;
} lrec_evaluator_f_z_state_t;

mv_t lrec_evaluator_f_z_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_z_state_t* pstate = pvstate;

	return pstate->pfunc();
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_z_func(mv_zary_func_t* pfunc) {
	lrec_evaluator_f_z_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_z_state_t));
	pstate->pfunc = pfunc;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_z_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_f_state_t {
	mv_unary_func_t* pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_f_state_t;

mv_t lrec_evaluator_f_f_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	mt_get_double_nullable(&val1);
	NULL_OUT(val1);
	if (val1.type != MT_FLOAT)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_f_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_ff_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_f_ff_state_t;

mv_t lrec_evaluator_f_ff_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mt_get_double_nullable(&val1);
	NULL_OUT(val1);
	if (val1.type != MT_FLOAT)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	mt_get_double_nullable(&val2);
	NULL_OUT(val2);
	if (val2.type != MT_FLOAT)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_ff_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_f_ff_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_ff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_ff_func;

	return pevaluator;
}

// This is for min/max which can return non-null when one argument is null --
// in comparison to other functions which return null if *any* argument is
// null.
mv_t lrec_evaluator_f_ff_nullable_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	ERROR_OUT(val1);
	mt_get_double_nullable(&val1);
	if (val1.type != MT_FLOAT && val1.type != MT_NULL)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	ERROR_OUT(val2);
	mt_get_double_nullable(&val2);
	if (val2.type != MT_FLOAT && val2.type != MT_NULL)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_ff_nullable_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_f_ff_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_ff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_ff_nullable_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_s_state_t;

mv_t lrec_evaluator_s_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_f_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_f_state_t;

mv_t lrec_evaluator_s_f_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type == MT_FLOAT) {
		;
	} else if (val1.type == MT_INT) {
		val1.type = MT_FLOAT;
		val1.u.fltv = (double)val1.u.intv;
	} else {
		return MV_ERROR;
	}

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_f_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_i_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_i_state_t;

mv_t lrec_evaluator_s_i_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type == MT_INT) {
		;
	} else if (val1.type == MT_FLOAT) {
		val1.type = MT_INT;
		val1.u.intv = (long long)val1.u.fltv;
	} else {
		return MV_ERROR;
	}

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_i_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_i_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_i_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_i_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_s_state_t;

mv_t lrec_evaluator_f_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_f_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_f_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_i_s_state_t;

mv_t lrec_evaluator_i_s_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_i_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_i_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_i_s_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_x_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_x_x_state_t;

mv_t lrec_evaluator_x_x_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_x_x_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_x_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_x_x_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_x_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_x_x_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_xx_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_xx_state_t;

mv_t lrec_evaluator_b_xx_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_b_xx_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_xx_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_b_xx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_xx_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_b_xx_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_ss_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_x_ss_state_t;

mv_t lrec_evaluator_x_ss_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_x_ss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_ss_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_x_ss_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_ss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_x_ss_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_xs_state_t {
	mv_binary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_s_xs_state_t;

mv_t lrec_evaluator_s_xs_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_xs_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_xs_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_s_xs_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_xs_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_xs_func;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_sss_state_t {
	mv_ternary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_s_sss_state_t;

mv_t lrec_evaluator_s_sss_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_s_sss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pevaluator_func(prec, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pevaluator_func(prec, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	mv_t val3 = pstate->parg3->pevaluator_func(prec, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT(val3);
	if (val3.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2, &val3);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_sss_func(mv_ternary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	lrec_evaluator_s_sss_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_sss_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_s_sss_func;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_field_name_state_t {
	char* field_name;
} lrec_evaluator_field_name_state_t;

mv_t lrec_evaluator_field_name_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_field_name_state_t* pstate = pvstate;
	char* string = lrec_get(prec, pstate->field_name);
	if (string == NULL) {
		return (mv_t) {.type = MT_NULL, .u.intv = 0};
	} else {
		double fltv;
		if (mlr_try_double_from_string(string, &fltv)) {
			return (mv_t) {.type = MT_FLOAT, .u.fltv = fltv};
		} else {
			return (mv_t) {.type = MT_STRING, .u.strv = mlr_strdup_or_die(string)};
		}
	}
}

lrec_evaluator_t* lrec_evaluator_alloc_from_field_name(char* field_name) {
	lrec_evaluator_field_name_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_field_name_state_t));
	pstate->field_name = mlr_strdup_or_die(field_name);

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pevaluator_func = lrec_evaluator_field_name_func;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_literal_state_t {
	mv_t literal;
} lrec_evaluator_literal_state_t;

mv_t lrec_evaluator_double_literal_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_literal_state_t* pstate = pvstate;
	return pstate->literal;
}
mv_t lrec_evaluator_string_literal_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	lrec_evaluator_literal_state_t* pstate = pvstate;
	// This is due to strdup-only semantics in mlrvals. If we implement a
	// free-flag as in slls and lrec, we could reduce some of the needless
	// strdups (at the cost of some code complexity).
	return (mv_t) {.type = MT_STRING, .u.strv = mlr_strdup_or_die(pstate->literal.u.strv)};
}

lrec_evaluator_t* lrec_evaluator_alloc_from_literal(char* string) {
	lrec_evaluator_literal_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_literal_state_t));
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));

	double fltv;
	if (mlr_try_double_from_string(string, &fltv)) {
		pstate->literal = (mv_t) {.type = MT_FLOAT, .u.fltv = fltv};
		pevaluator->pevaluator_func = lrec_evaluator_double_literal_func;
	} else {
		pstate->literal = (mv_t) {.type = MT_STRING, .u.strv = mlr_strdup_or_die(string)};
		pevaluator->pevaluator_func = lrec_evaluator_string_literal_func;
	}
	pevaluator->pvstate = pstate;

	return pevaluator;
}

// ================================================================
mv_t lrec_evaluator_NF_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = prec->field_count};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NF() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_NF_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_NR_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->nr};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_NR_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FNR_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->fnr};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FNR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FNR_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENAME_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_STRING, .u.strv = mlr_strdup_or_die(pctx->filename)};
}

lrec_evaluator_t* lrec_evaluator_alloc_from_FILENAME() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FILENAME_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENUM_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_INT, .u.intv = pctx->filenum};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FILENUM() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_FILENUM_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_PI_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_FLOAT, .u.fltv = M_PI};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_PI() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_PI_func;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_E_func(lrec_t* prec, context_t* pctx, void* pvstate) {
	return (mv_t) {.type = MT_FLOAT, .u.fltv = M_E};
}
lrec_evaluator_t* lrec_evaluator_alloc_from_E() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pevaluator_func = lrec_evaluator_E_func;
	return pevaluator;
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_context_variable(char* variable_name) {
	if        (streq(variable_name, "NF"))       { return lrec_evaluator_alloc_from_NF();
    } else if (streq(variable_name, "NR"))       { return lrec_evaluator_alloc_from_NR();
    } else if (streq(variable_name, "FNR"))      { return lrec_evaluator_alloc_from_FNR();
    } else if (streq(variable_name, "FILENAME")) { return lrec_evaluator_alloc_from_FILENAME();
    } else if (streq(variable_name, "FILENUM"))  { return lrec_evaluator_alloc_from_FILENUM();
    } else if (streq(variable_name, "PI"))       { return lrec_evaluator_alloc_from_PI();
    } else if (streq(variable_name, "E"))        { return lrec_evaluator_alloc_from_E();
	} else  { return NULL; 
	}
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_zary_func_name(char* function_name) {
	if        (streq(function_name, "urand")) {
		return lrec_evaluator_alloc_from_f_z_func(f_z_urand_func);
	} else if (streq(function_name, "systime")) {
		return lrec_evaluator_alloc_from_f_z_func(f_z_systime_func);
	} else  {
		return NULL;
	}
}

// ================================================================
typedef struct _function_lookup_t {
	int   function_class;
	char* function_name;
	int   arity;
	char* usage_string;
} function_lookup_t;

#define FUNC_CLASS_MATH       0xa0
#define FUNC_CLASS_BOOLEAN    0xa1
#define FUNC_CLASS_STRING     0xa3
#define FUNC_CLASS_CONVERSION 0xa4
#define FUNC_CLASS_TIME       0xa2

static function_lookup_t FUNCTION_LOOKUP_TABLE[] = {
	{ FUNC_CLASS_MATH, "abs",      1 , "Absolute value."},
	{ FUNC_CLASS_MATH, "acos",     1 , "Inverse trigonometric cosine."},
	{ FUNC_CLASS_MATH, "acosh",    1 , "Inverse hyperbolic cosine."},
	{ FUNC_CLASS_MATH, "asin",     1 , "Inverse trigonometric sine."},
	{ FUNC_CLASS_MATH, "asinh",    1 , "Inverse hyperbolic sine."},
	{ FUNC_CLASS_MATH, "atan",     1 , "One-argument arctangent."},
	{ FUNC_CLASS_MATH, "atan2",    2 , "Two-argument arctangent."},
	{ FUNC_CLASS_MATH, "atanh",    1 , "Inverse hyperbolic tangent."},
	{ FUNC_CLASS_MATH, "cbrt",     1 , "Cube root."},
	{ FUNC_CLASS_MATH, "ceil",     1 , "Ceiling: nearest integer at or above."},
	{ FUNC_CLASS_MATH, "cos",      1 , "Trigonometric cosine."},
	{ FUNC_CLASS_MATH, "cosh",     1 , "Hyperbolic cosine."},
	{ FUNC_CLASS_MATH, "erf",      1 , "Error function."},
	{ FUNC_CLASS_MATH, "erfc",     1 , "Complementary error function."},
	{ FUNC_CLASS_MATH, "exp",      1 , "Exponential function e**x."},
	{ FUNC_CLASS_MATH, "expm1",    1 , "e**x - 1."},
	{ FUNC_CLASS_MATH, "floor",    1 , "Floor: nearest integer at or below."},
	// See also http://johnkerl.org/doc/randuv.pdf for more about urand() -> other distributions
	{ FUNC_CLASS_MATH, "invqnorm", 1 , "Inverse of normal cumulative distribution function. Note that invqorm(urand()) is normally distributed."},
	{ FUNC_CLASS_MATH, "log",      1 , "Natural (base-e) logarithm."},
	{ FUNC_CLASS_MATH, "log10",    1 , "Base-10 logarithm."},
	{ FUNC_CLASS_MATH, "log1p",    1 , "log(1-x)."},
	{ FUNC_CLASS_MATH, "max",      2 , "max of two numbers; null loses"},
	{ FUNC_CLASS_MATH, "min",      2 , "min of two numbers; null loses"},
	{ FUNC_CLASS_MATH, "pow",      2 , "Exponentiation; same as **."},
	{ FUNC_CLASS_MATH, "qnorm",    1 , "Normal cumulative distribution function."},
	{ FUNC_CLASS_MATH, "round",    1 , "Round to nearest integer."},
	{ FUNC_CLASS_MATH, "roundm",   2 , "Round to nearest multiple of m: roundm($x,$m) is the same as round($x/$m)*$m"},
	{ FUNC_CLASS_MATH, "sin",      1 , "Trigonometric sine."},
	{ FUNC_CLASS_MATH, "sinh",     1 , "Hyperbolic sine."},
	{ FUNC_CLASS_MATH, "sqrt",     1 , "Square root."},
	{ FUNC_CLASS_MATH, "tan",      1 , "Trigonometric tangent."},
	{ FUNC_CLASS_MATH, "tanh",     1 , "Hyperbolic tangent."},
	{ FUNC_CLASS_MATH, "urand",    0 , "Floating-point numbers on the unit interval. Int-valued example: '$n=floor(20+urand()*11)'." },

	{ FUNC_CLASS_MATH, "+",       2 , "Addition."},
	{ FUNC_CLASS_MATH, "-",       1 , "Unary minus."},
	{ FUNC_CLASS_MATH, "-",       2 , "Subtraction."},
	{ FUNC_CLASS_MATH, "*",       2 , "Multiplication."},
	{ FUNC_CLASS_MATH, "/",       2 , "Division."},
	{ FUNC_CLASS_MATH, "%",       2 , "Remainder; never negative-valued."},
	{ FUNC_CLASS_MATH, "**",      2 , "Exponentiation; same as pow."},

	{ FUNC_CLASS_BOOLEAN, "=~",      2 , "String (left-hand side) matches regex (right-hand side) [under construction]."},
	{ FUNC_CLASS_BOOLEAN, "!=~",     2 , "String (left-hand side) does not match regex (right-hand side) [under construction]."},
	{ FUNC_CLASS_BOOLEAN, "==",      2 , "String/numeric equality. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "!=",      2 , "String/numeric inequality. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, ">",       2 , "String/numeric greater-than. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, ">=",      2 , "String/numeric greater-than-or-equals. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "<",       2 , "String/numeric less-than. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "<=",      2 , "String/numeric less-than-or-equals. Mixing number and string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "&&",      2 , "Logical AND."},
	{ FUNC_CLASS_BOOLEAN, "||",      2 , "Logical OR."},
	{ FUNC_CLASS_BOOLEAN, "!",       1 , "Logical negation."},

	{ FUNC_CLASS_STRING, "strlen",   1 , "String length."},
	{ FUNC_CLASS_STRING, "sub",      3 , "Example: '$name=sub($name, \"old\", \"new\")'. Regexes not supported."},
	{ FUNC_CLASS_STRING, "tolower",  1 , "Convert string to lowercase."},
	{ FUNC_CLASS_STRING, "toupper",  1 , "Convert string to uppercase."},
	{ FUNC_CLASS_STRING, ".",       2 , "String concatenation."},

	{ FUNC_CLASS_CONVERSION, "boolean",  1 , "Convert int/float/bool/string to boolean."},
	{ FUNC_CLASS_CONVERSION, "float",    1 , "Convert int/float/bool/string to float."},
	{ FUNC_CLASS_CONVERSION, "int",      1 , "Convert int/float/bool/string to int."},
	{ FUNC_CLASS_CONVERSION, "string",   1 , "Convert int/float/bool/string to string."},
	{ FUNC_CLASS_CONVERSION, "hexfmt",   1 , "Convert int to string, e.g. 255 to \"0xff\"."},
	{ FUNC_CLASS_CONVERSION, "fmtnum",   2 , "Convert int/float/bool to string using printf-style format string, e.g. \"%06lld\"."},

	{ FUNC_CLASS_TIME, "systime",   0 , "Floating-point seconds since the epoch, e.g. 1440768801.748936." },
	{ FUNC_CLASS_TIME, "sec2gmt",   1 , "Formats seconds since epoch (integer part only) as GMT timestamp, e.g. sec2gmt(1440768801.7) = \"2015-08-28T13:33:21Z\"."},
	{ FUNC_CLASS_TIME, "gmt2sec",   1 , "Parses GMT timestamp as integer seconds since epoch."},
	{ FUNC_CLASS_TIME, "sec2hms",   1 , "Formats integer seconds as in sec2hms(5000) = \"01:23:20\""},
	{ FUNC_CLASS_TIME, "sec2dhms",  1 , "Formats integer seconds as in sec2dhms(500000) = \"5d18h53m20s\""},
	{ FUNC_CLASS_TIME, "hms2sec",   1 , "Recovers integer seconds as in hms2sec(\"01:23:20\") = 5000"},
	{ FUNC_CLASS_TIME, "dhms2sec",  1 , "Recovers integer seconds as in dhms2sec(\"5d18h53m20s\") = 500000"},
	{ FUNC_CLASS_TIME, "fsec2hms",  1 , "Formats floating-point seconds as in fsec2hms(5000.25) = \"01:23:20.250000\""},
	{ FUNC_CLASS_TIME, "fsec2dhms", 1 , "Formats floating-point seconds as in fsec2dhms(500000.25) = \"5d18h53m20.250000s\""},
	{ FUNC_CLASS_TIME, "hms2fsec",  1 , "Recovers floating-point seconds as in hms2fsec(\"01:23:20.250000\") = 5000.250000"},
	{ FUNC_CLASS_TIME, "dhms2fsec", 1 , "Recovers floating-point seconds as in dhms2fsec(\"5d18h53m20.250000s\") = 500000.250000"},

	{  0, NULL,      -1 , NULL}, // table terminator
};

#define ARITY_CHECK_PASS    0xbb
#define ARITY_CHECK_FAIL    0xbc
#define ARITY_CHECK_NO_SUCH 0xbd

static int check_arity(function_lookup_t lookup_table[], char* function_name, int user_provided_arity, int *parity) {
	*parity = -1;
	int found_function_name = FALSE;
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &lookup_table[i];
		if (plookup->function_name == NULL)
			break;
		if (streq(function_name, plookup->function_name)) {
			found_function_name = TRUE;
			*parity = plookup->arity;
			if (user_provided_arity == plookup->arity) {
				return ARITY_CHECK_PASS;
			}
		}
	}
	if (found_function_name) {
		return ARITY_CHECK_FAIL;
	} else {
		return ARITY_CHECK_NO_SUCH;
	}
}

static void check_arity_with_report(function_lookup_t function_lookup_table[], char* function_name,
	int user_provided_arity)
{
	int arity = -1;
	int result = check_arity(function_lookup_table, function_name, user_provided_arity, &arity);
	if (result == ARITY_CHECK_NO_SUCH) {
		fprintf(stderr, "Function name \"%s\" not found.\n", function_name);
		exit(1);
	}
	if (result == ARITY_CHECK_FAIL) {
		// More flexibly, I'd have a list of arities supported by each
		// function. But this is overkill: there are unary and binary minus,
		// and everything else has a single arity.
		if (streq(function_name, "-")) {
			fprintf(stderr, "Function named \"%s\" takes one argument or two; got %d.\n",
				function_name, user_provided_arity);
		} else {
		}
			fprintf(stderr, "Function named \"%s\" takes %d argument%s; got %d.\n",
				function_name, arity, (arity == 1) ? "" : "s", user_provided_arity);
		exit(1);
	}
}

static char* function_class_to_desc(int function_class) {
	switch(function_class) {
	case FUNC_CLASS_MATH:       return "math";       break;
	case FUNC_CLASS_BOOLEAN:    return "boolean";    break;
	case FUNC_CLASS_STRING:     return "string";     break;
	case FUNC_CLASS_CONVERSION: return "conversion"; break;
	case FUNC_CLASS_TIME:       return "time";       break;
	default:                    return "???";        break;
	}
}

void lrec_evaluator_list_functions(FILE* o) {
	char* leader = "  ";
	char* separator = " ";
	int leaderlen = strlen(leader);
	int separatorlen = strlen(separator);
	int linelen = leaderlen;
	int j = 0;
	fprintf(o, "Functions for filter and put:\n");

	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		char* fname = plookup->function_name;
		if (fname == NULL)
			break;
		int fnamelen = strlen(fname);
		linelen += separatorlen + fnamelen;
		if (linelen >= 80) {
			fprintf(o, "\n");
			linelen = 0;
			linelen = leaderlen + separatorlen + fnamelen;
			j = 0;
		}
		if (j == 0)
			fprintf(o, "%s", leader);
		fprintf(o, "%s%s", separator, fname);
		j++;
	}
	fprintf(o, "\n");
}

// Pass function_name == NULL to get usage for all functions.
void lrec_evaluator_function_usage(FILE* output_stream, char* function_name) {
	int found = FALSE;
	char class_and_colon[128];
	char* fmt = (function_name == NULL)
		? "%-10s (%-11s #args=%d): %s\n"
		: "%s (%s #args=%d): %s\n";

	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL)
			break;
		if (function_name == NULL || streq(function_name, plookup->function_name)) {
			sprintf(class_and_colon, "%s:", function_class_to_desc(plookup->function_class));
			fprintf(output_stream, fmt, plookup->function_name, class_and_colon,
				plookup->arity, plookup->usage_string);
			found = TRUE;
		}
	}
	if (!found)
		fprintf(output_stream, "%s: no such function.\n", function_name);
	if (function_name == NULL) {
		fprintf(output_stream, "To set the seed for urand, you may specify decimal or hexadecimal 32-bit\n");
		fprintf(output_stream, "numbers of the form \"%s --seed 123456789\" or \"%s --seed 0xcafefeed\".\n",
			MLR_GLOBALS.argv0, MLR_GLOBALS.argv0);
		fprintf(output_stream, "Miller's built-in variables are NF, NR, FNR, FILENUM, and FILENAME (awk-like)\n");
		fprintf(output_stream, "along with the mathematical constants PI and E.\n");
	}
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_unary_func_name(char* fnnm, lrec_evaluator_t* parg1)  {
	if        (streq(fnnm, "!"))         { return lrec_evaluator_alloc_from_b_b_func(b_b_not_func,       parg1);
	} else if (streq(fnnm, "-"))         { return lrec_evaluator_alloc_from_f_f_func(f_f_uneg_func,      parg1);
	} else if (streq(fnnm, "abs"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_abs_func,       parg1);
	} else if (streq(fnnm, "acos"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_acos_func,      parg1);
	} else if (streq(fnnm, "acosh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_acosh_func,     parg1);
	} else if (streq(fnnm, "asin"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_asin_func,      parg1);
	} else if (streq(fnnm, "asinh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_asinh_func,     parg1);
	} else if (streq(fnnm, "atan"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_atan_func,      parg1);
	} else if (streq(fnnm, "atanh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_atanh_func,     parg1);
	} else if (streq(fnnm, "boolean"))   { return lrec_evaluator_alloc_from_x_x_func(b_x_boolean_func,   parg1);
	} else if (streq(fnnm, "cbrt"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_cbrt_func,      parg1);
	} else if (streq(fnnm, "ceil"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_ceil_func,      parg1);
	} else if (streq(fnnm, "cos"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_cos_func,       parg1);
	} else if (streq(fnnm, "cosh"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_cosh_func,      parg1);
	} else if (streq(fnnm, "erf"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_erf_func,       parg1);
	} else if (streq(fnnm, "erfc"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_erfc_func,      parg1);
	} else if (streq(fnnm, "exp"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_exp_func,       parg1);
	} else if (streq(fnnm, "expm1"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_expm1_func,     parg1);
	} else if (streq(fnnm, "float"))     { return lrec_evaluator_alloc_from_x_x_func(f_x_float_func,     parg1);
	} else if (streq(fnnm, "floor"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_floor_func,     parg1);
	} else if (streq(fnnm, "gmt2sec"))   { return lrec_evaluator_alloc_from_i_s_func(i_s_gmt2sec_func,   parg1);
	} else if (streq(fnnm, "hms2sec"))   { return lrec_evaluator_alloc_from_f_s_func(i_s_hms2sec_func,   parg1);
	} else if (streq(fnnm, "hms2fsec"))  { return lrec_evaluator_alloc_from_f_s_func(f_s_hms2fsec_func,  parg1);
	} else if (streq(fnnm, "dhms2sec"))  { return lrec_evaluator_alloc_from_f_s_func(i_s_dhms2sec_func,  parg1);
	} else if (streq(fnnm, "dhms2fsec")) { return lrec_evaluator_alloc_from_f_s_func(f_s_dhms2fsec_func, parg1);
	} else if (streq(fnnm, "hexfmt"))    { return lrec_evaluator_alloc_from_x_x_func(s_x_hexfmt_func,    parg1);
	} else if (streq(fnnm, "int"))       { return lrec_evaluator_alloc_from_x_x_func(i_x_int_func,       parg1);
	} else if (streq(fnnm, "log"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_log_func,       parg1);
	} else if (streq(fnnm, "log10"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_log10_func,     parg1);
	} else if (streq(fnnm, "log1p"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_log1p_func,     parg1);
	} else if (streq(fnnm, "qnorm"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_qnorm_func,     parg1);
	} else if (streq(fnnm, "invqnorm"))  { return lrec_evaluator_alloc_from_f_f_func(f_f_invqnorm_func,  parg1);
	} else if (streq(fnnm, "round"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_round_func,     parg1);
	} else if (streq(fnnm, "sec2gmt"))   { return lrec_evaluator_alloc_from_s_f_func(s_f_sec2gmt_func,   parg1);
	} else if (streq(fnnm, "sec2hms"))   { return lrec_evaluator_alloc_from_s_i_func(s_i_sec2hms_func,   parg1);
	} else if (streq(fnnm, "fsec2hms"))  { return lrec_evaluator_alloc_from_s_f_func(s_f_fsec2hms_func,  parg1);
	} else if (streq(fnnm, "sec2dhms"))  { return lrec_evaluator_alloc_from_s_i_func(s_i_sec2dhms_func,  parg1);
	} else if (streq(fnnm, "fsec2dhms")) { return lrec_evaluator_alloc_from_s_f_func(s_f_fsec2dhms_func, parg1);
	} else if (streq(fnnm, "sin"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_sin_func,       parg1);
	} else if (streq(fnnm, "sinh"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_sinh_func,      parg1);
	} else if (streq(fnnm, "sqrt"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_sqrt_func,      parg1);
	} else if (streq(fnnm, "string"))    { return lrec_evaluator_alloc_from_x_x_func(s_x_string_func,    parg1);
	} else if (streq(fnnm, "strlen"))    { return lrec_evaluator_alloc_from_i_s_func(i_s_strlen_func,    parg1);
	} else if (streq(fnnm, "tan"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_tan_func,       parg1);
	} else if (streq(fnnm, "tanh"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_tanh_func,      parg1);
	} else if (streq(fnnm, "tolower"))   { return lrec_evaluator_alloc_from_s_s_func(s_s_tolower_func,   parg1);
	} else if (streq(fnnm, "toupper"))   { return lrec_evaluator_alloc_from_s_s_func(s_s_toupper_func,   parg1);

	} else return NULL;
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_binary_func_name(char* fnnm,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	if        (streq(fnnm, "&&"))     { return lrec_evaluator_alloc_from_b_bb_func(b_bb_and_func,             parg1, parg2);
	} else if (streq(fnnm, "||"))     { return lrec_evaluator_alloc_from_b_bb_func(b_bb_or_func,              parg1, parg2);
	} else if (streq(fnnm, "=~"))     { return lrec_evaluator_alloc_from_x_ss_func(matches_op_func,           parg1, parg2);
	} else if (streq(fnnm, "!=~"))    { return lrec_evaluator_alloc_from_x_ss_func(does_not_match_op_func,    parg1, parg2);
	} else if (streq(fnnm, "=="))     { return lrec_evaluator_alloc_from_b_xx_func(eq_op_func,                parg1, parg2);
	} else if (streq(fnnm, "!="))     { return lrec_evaluator_alloc_from_b_xx_func(ne_op_func,                parg1, parg2);
	} else if (streq(fnnm, ">"))      { return lrec_evaluator_alloc_from_b_xx_func(gt_op_func,                parg1, parg2);
	} else if (streq(fnnm, ">="))     { return lrec_evaluator_alloc_from_b_xx_func(ge_op_func,                parg1, parg2);
	} else if (streq(fnnm, "<"))      { return lrec_evaluator_alloc_from_b_xx_func(lt_op_func,                parg1, parg2);
	} else if (streq(fnnm, "<="))     { return lrec_evaluator_alloc_from_b_xx_func(le_op_func,                parg1, parg2);
	} else if (streq(fnnm, "."))      { return lrec_evaluator_alloc_from_x_ss_func(s_ss_dot_func,             parg1, parg2);
	} else if (streq(fnnm, "+"))      { return lrec_evaluator_alloc_from_f_ff_func(f_ff_plus_func,            parg1, parg2);
	} else if (streq(fnnm, "-"))      { return lrec_evaluator_alloc_from_f_ff_func(f_ff_minus_func,           parg1, parg2);
	} else if (streq(fnnm, "*"))      { return lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func,           parg1, parg2);
	} else if (streq(fnnm, "/"))      { return lrec_evaluator_alloc_from_f_ff_func(f_ff_divide_func,          parg1, parg2);
	} else if (streq(fnnm, "**"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,             parg1, parg2);
	} else if (streq(fnnm, "pow"))    { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,             parg1, parg2);
	} else if (streq(fnnm, "%"))      { return lrec_evaluator_alloc_from_f_ff_func(f_ff_mod_func,             parg1, parg2);
	} else if (streq(fnnm, "atan2"))  { return lrec_evaluator_alloc_from_f_ff_func(f_ff_atan2_func,           parg1, parg2);
	} else if (streq(fnnm, "max"))    { return lrec_evaluator_alloc_from_f_ff_nullable_func(f_ff_max_func,    parg1, parg2);
	} else if (streq(fnnm, "min"))    { return lrec_evaluator_alloc_from_f_ff_nullable_func(f_ff_min_func,    parg1, parg2);
	} else if (streq(fnnm, "roundm")) { return lrec_evaluator_alloc_from_f_ff_nullable_func(f_ff_roundm_func, parg1, parg2);
	} else if (streq(fnnm, "fmtnum")) { return lrec_evaluator_alloc_from_s_xs_func(s_xs_fmtnum_func,          parg1, parg2);
	} else  { return NULL; }
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_ternary_func_name(char* fnnm,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	if (streq(fnnm, "sub")) { return lrec_evaluator_alloc_from_s_sss_func(s_sss_sub_func,   parg1, parg2, parg3);
	} else  { return NULL; }
}

// ================================================================
static lrec_evaluator_t* lrec_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	function_lookup_t* function_lookup_table)
{
	if (pnode->pchildren == NULL) { // leaf node
		if (pnode->type == MLR_DSL_AST_NODE_TYPE_FIELD_NAME) {
			return lrec_evaluator_alloc_from_field_name(pnode->text);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_LITERAL) {
			return lrec_evaluator_alloc_from_literal(pnode->text);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE) {
			return lrec_evaluator_alloc_from_context_variable(pnode->text);
		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
	} else { // operator/function
		if ((pnode->type != MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME)
		&& (pnode->type != MLR_DSL_AST_NODE_TYPE_OPERATOR)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
		char* func_name = pnode->text;

		int user_provided_arity = pnode->pchildren->length;

		check_arity_with_report(function_lookup_table, func_name, user_provided_arity);

		lrec_evaluator_t* pevaluator = NULL;
		if (user_provided_arity == 0) {
			pevaluator = lrec_evaluator_alloc_from_zary_func_name(func_name);
		} else if (user_provided_arity == 1) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_unary_func_name(func_name, parg1);
		} else if (user_provided_arity == 2) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_binary_func_name(func_name, parg1, parg2);
		} else if (user_provided_arity == 3) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvdata;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvdata;
			mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvdata;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, function_lookup_table);
			lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, function_lookup_table);
			lrec_evaluator_t* parg3 = lrec_evaluator_alloc_from_ast_aux(parg3_node, function_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_ternary_func_name(func_name, parg1, parg2, parg3);
		} else {
			fprintf(stderr, "Miller: internal coding error:  arity for function name \"%s\" misdetected.\n",
				func_name);
			exit(1);
		}
		if (pevaluator == NULL) {
			fprintf(stderr, "Miller: unrecognized function name \"%s\".\n", func_name);
			exit(1);
		}
		return pevaluator;
	}
}

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode) {
	lrec_evaluator_t* pevaluator = lrec_evaluator_alloc_from_ast_aux(pnode, FUNCTION_LOOKUP_TABLE);
	return pevaluator;
}

// ================================================================
#include "lib/minunit.h"

// ----------------------------------------------------------------
int tests_run         = 0;
int tests_failed      = 0;
int assertions_run    = 0;
int assertions_failed = 0;

// ----------------------------------------------------------------
static char * test1() {
	printf("\n");
	printf("-- TEST_LREC_EVALUATORS test1 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	lrec_evaluator_t* pnr       = lrec_evaluator_alloc_from_NR();
	lrec_evaluator_t* pfnr      = lrec_evaluator_alloc_from_FNR();
	lrec_evaluator_t* pfilename = lrec_evaluator_alloc_from_FILENAME();
	lrec_evaluator_t* pfilenum  = lrec_evaluator_alloc_from_FILENUM();

	lrec_t* prec = lrec_unbacked_alloc();

	mv_t val = pnr->pevaluator_func(prec, pctx, pnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 888);

	val = pfnr->pevaluator_func(prec, pctx, pfnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 999);

	val = pfilename->pevaluator_func(prec, pctx, pfilename->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "filename-goes-here"));

	val = pfilenum->pevaluator_func(prec, pctx, pfilenum->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 123);

	return 0;
}

// ----------------------------------------------------------------
static char * test2() {
	printf("\n");
	printf("-- TEST_LREC_EVALUATORS test2 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	lrec_evaluator_t* ps       = lrec_evaluator_alloc_from_field_name("s");
	lrec_evaluator_t* pdef     = lrec_evaluator_alloc_from_literal("def");
	lrec_evaluator_t* pdot     = lrec_evaluator_alloc_from_x_ss_func(s_ss_dot_func, ps, pdef);
	lrec_evaluator_t* ptolower = lrec_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	lrec_evaluator_t* ptoupper = lrec_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put_no_free(prec, "s", "abc");
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	mv_t val = ps->pevaluator_func(prec, pctx, ps->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abc"));

	val = pdef->pevaluator_func(prec, pctx, pdef->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "def"));

	val = pdot->pevaluator_func(prec, pctx, pdot->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptolower->pevaluator_func(prec, pctx, ptolower->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptoupper->pevaluator_func(prec, pctx, ptoupper->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mt_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "ABCDEF"));

	return 0;
}

// ----------------------------------------------------------------
static char * test3() {
	printf("\n");
	printf("-- TEST_LREC_EVALUATORS test3 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	lrec_evaluator_t* p2     = lrec_evaluator_alloc_from_literal("2.0");
	lrec_evaluator_t* px     = lrec_evaluator_alloc_from_field_name("x");
	lrec_evaluator_t* plogx  = lrec_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	lrec_evaluator_t* p2logx = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, p2, plogx);
	lrec_evaluator_t* px2    = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, px, px);
	lrec_evaluator_t* p4     = lrec_evaluator_alloc_from_f_ff_func(f_ff_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MLR_DSL_AST_NODE_TYPE_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MLR_DSL_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	lrec_evaluator_t*  pastr = lrec_evaluator_alloc_from_ast(p2logxnode);

	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put_no_free(prec, "x", "4.5");

	mv_t valp2     = p2->pevaluator_func(prec,     pctx, p2->pvstate);
	mv_t valp4     = p4->pevaluator_func(prec,     pctx, p4->pvstate);
	mv_t valpx     = px->pevaluator_func(prec,     pctx, px->pvstate);
	mv_t valpx2    = px2->pevaluator_func(prec,    pctx, px2->pvstate);
	mv_t valplogx  = plogx->pevaluator_func(prec,  pctx, plogx->pvstate);
	mv_t valp2logx = p2logx->pevaluator_func(prec, pctx, p2logx->pvstate);

    printf("lrec   x        = %s\n", lrec_get(prec, "x"));
    printf("newval 2        = %s\n", mt_describe_val(valp2));
    printf("newval 4        = %s\n", mt_describe_val(valp4));
    printf("newval x        = %s\n", mt_describe_val(valpx));
    printf("newval x^2      = %s\n", mt_describe_val(valpx2));
    printf("newval log(x)   = %s\n", mt_describe_val(valplogx));
    printf("newval 2*log(x) = %s\n", mt_describe_val(valp2logx));

	mu_assert_lf(valp2.type     == MT_FLOAT);
	mu_assert_lf(valp4.type     == MT_FLOAT);
	mu_assert_lf(valpx.type     == MT_FLOAT);
	mu_assert_lf(valpx2.type    == MT_FLOAT);
	mu_assert_lf(valplogx.type  == MT_FLOAT);
	mu_assert_lf(valp2logx.type == MT_FLOAT);

	mu_assert_lf(valp2.u.fltv     == 2.0);
	mu_assert_lf(valp4.u.fltv     == 4.0);
	mu_assert_lf(valpx.u.fltv     == 4.5);
	mu_assert_lf(valpx2.u.fltv    == 20.25);
	mu_assert_lf(fabs(valplogx.u.fltv  - 0.653213) < 1e-5);
	mu_assert_lf(fabs(valp2logx.u.fltv - 1.306425) < 1e-5);

	mlr_dsl_ast_node_print(p2logxnode);
	printf("newval AST      = %s\n",  mt_describe_val(pastr->pevaluator_func(prec, pctx, pastr->pvstate)));
	printf("\n");

	lrec_rename(prec, "x", "y");

	valp2     = p2->pevaluator_func(prec,     pctx, p2->pvstate);
	valp4     = p4->pevaluator_func(prec,     pctx, p4->pvstate);
	valpx     = px->pevaluator_func(prec,     pctx, px->pvstate);
	valpx2    = px2->pevaluator_func(prec,    pctx, px2->pvstate);
	valplogx  = plogx->pevaluator_func(prec,  pctx, plogx->pvstate);
	valp2logx = p2logx->pevaluator_func(prec, pctx, p2logx->pvstate);

    printf("lrec   x        = %s\n", lrec_get(prec, "x"));
    printf("newval 2        = %s\n", mt_describe_val(valp2));
    printf("newval 4        = %s\n", mt_describe_val(valp4));
    printf("newval x        = %s\n", mt_describe_val(valpx));
    printf("newval x^2      = %s\n", mt_describe_val(valpx2));
    printf("newval log(x)   = %s\n", mt_describe_val(valplogx));
    printf("newval 2*log(x) = %s\n", mt_describe_val(valp2logx));

	mu_assert_lf(valp2.type     == MT_FLOAT);
	mu_assert_lf(valp4.type     == MT_FLOAT);
	mu_assert_lf(valpx.type     == MT_NULL);
	mu_assert_lf(valpx2.type    == MT_NULL);
	mu_assert_lf(valplogx.type  == MT_NULL);
	mu_assert_lf(valp2logx.type == MT_NULL);

	mu_assert_lf(valp2.u.fltv     == 2.0);
	mu_assert_lf(valp4.u.fltv     == 4.0);

	return 0;
}

// ================================================================
static char * all_tests() {
	mu_run_test(test1);
	mu_run_test(test2);
	mu_run_test(test3);
	return 0;
}

// test_lrec_evaluators has the MinUnit inside lrec_evaluators, as it tests
// many private methods. (The other option is to make them all public.)
int test_lrec_evaluators_main(int argc, char **argv) {
	mlr_global_init(argv[0], "%lf", NULL);

	printf("TEST_LREC_EVALUATORS ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_LREC_EVALUATORS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
