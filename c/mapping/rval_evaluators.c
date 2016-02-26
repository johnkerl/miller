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
// NOTES:
//
// * Code here evaluates right-hand-side values (rvals) and return mlrvals (mv_t).
//
// * This is used by mlr filter and mlr put.
//
// * Unlike most files in Miller which are read top-down (with sufficient
//   static prototypes at the top of the file to keep the compiler happy),
//   please read this one from the bottom up.
//
// * Comparison to mlrval.c: the latter is functions from mlrval(s) to
//   mlrval; in this file we have the higher-level notion of evaluating lrec
//   objects, using mlrval.c to do so.
//
// * There are two kinds of lrec-evaluators here: those with _x_ in their names
//   which accept various types of mlrval, with disposition-matrices in
//   mlrval.c functions, and those with _i_/_f_/_b_/_s_ (int, float, boolean,
//   string) which either type-check or type-coerce their arguments, invoking
//   type-specific functions in mlrval.c.  Those with _n_ take int or float
//   and also use disposition matrices.  In all cases it's the job of
//   rval_evaluators.c to invoke functions here with mlrvals of the correct
//   type(s).
// ================================================================

// ----------------------------------------------------------------
typedef struct _rval_evaluator_b_b_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_b_b_state_t;

mv_t rval_evaluator_b_b_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_b_b_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOL)
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
mv_t rval_evaluator_b_bb_and_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOL)
		return mv_error();
	if (val1.u.boolv == FALSE)
		return val1;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_BOOL)
		return mv_error();

	return val2;
}

// This is different from most of the lrec-evaluator functions in that it does short-circuiting:
// since is logical OR, the LHS is not evaluated if the RHS is true.
mv_t rval_evaluator_b_bb_or_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOL)
		return mv_error();
	if (val1.u.boolv == TRUE)
		return val1;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_BOOL)
		return mv_error();

	return val2;
}

mv_t rval_evaluator_b_bb_xor_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_b_bb_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_BOOL)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_BOOL)
		return mv_error();

	return mv_from_bool(val1.u.boolv ^ val2.u.boolv);
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

mv_t rval_evaluator_x_z_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
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

mv_t rval_evaluator_f_f_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_f_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

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

mv_t rval_evaluator_x_n_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_n_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
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

mv_t rval_evaluator_i_i_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_i_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

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

mv_t rval_evaluator_f_ff_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
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

mv_t rval_evaluator_x_xx_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_xx_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);

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

mv_t rval_evaluator_x_xx_nullable_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_xx_nullable_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_number_nullable(&val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
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

mv_t rval_evaluator_f_fff_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_f_fff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_float_nullable(&val2);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val2);

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg3->pvstate);
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

mv_t rval_evaluator_i_ii_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_i_ii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_int_nullable(&val1);
	NULL_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_INT)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_int_nullable(&val2);
	NULL_OUT_FOR_NUMBERS(val2);
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

mv_t rval_evaluator_i_iii_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_i_iii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_int_nullable(&val1);
	NULL_OUT_FOR_NUMBERS(val1);
	if (val1.type != MT_INT)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_int_nullable(&val2);
	NULL_OUT_FOR_NUMBERS(val2);
	if (val2.type != MT_INT)
		return mv_error();

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg3->pvstate);
	mv_set_int_nullable(&val3);
	NULL_OUT_FOR_NUMBERS(val3);
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

mv_t rval_evaluator_ternop_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_ternop_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);
	mv_set_boolean_strict(&val1);

	return val1.u.boolv
		? pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate)
		: pstate->parg3->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg3->pvstate);
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

mv_t rval_evaluator_s_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_s_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
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
typedef struct _rval_evaluator_s_f_state_t {
	mv_unary_func_t*  pfunc;
	rval_evaluator_t* parg1;
} rval_evaluator_s_f_state_t;

mv_t rval_evaluator_s_f_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_s_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

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

mv_t rval_evaluator_s_i_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_s_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

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

mv_t rval_evaluator_f_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_f_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
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

mv_t rval_evaluator_i_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_i_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
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

mv_t rval_evaluator_x_x_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_x_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

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
typedef struct _rval_evaluator_x_ns_state_t {
	mv_binary_func_t* pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_x_ns_state_t;

mv_t rval_evaluator_x_ns_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_ns_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT_FOR_NUMBERS(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (val2.type != MT_STRING && val2.type != MT_VOID) {
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

mv_t rval_evaluator_x_ss_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_ss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (val2.type != MT_STRING && val2.type != MT_VOID)
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

mv_t rval_evaluator_x_ssc_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_ssc_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (val2.type != MT_STRING && val2.type != MT_VOID)
		return mv_error();
	return pstate->pfunc(&val1, &val2, ppregex_captures);
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

mv_t rval_evaluator_x_sr_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_sr_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
		return mv_error();

	return pstate->pfunc(&val1, &pstate->regex, pstate->psb, ppregex_captures);
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
typedef struct _rval_evaluator_s_xs_state_t {
	mv_binary_func_t*  pfunc;
	rval_evaluator_t* parg1;
	rval_evaluator_t* parg2;
} rval_evaluator_s_xs_state_t;

mv_t rval_evaluator_s_xs_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_s_xs_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (val2.type != MT_STRING && val2.type != MT_VOID)
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

mv_t rval_evaluator_s_sss_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_s_sss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID)
		return mv_error();

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val2);
	if (val2.type != MT_STRING && val2.type != MT_VOID) {
		mv_free(&val1);
		return mv_error();
	}

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val3);
	if (val3.type != MT_STRING && val3.type != MT_VOID) {
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

mv_t rval_evaluator_x_srs_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_x_srs_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val1);
	if (val1.type != MT_STRING && val1.type != MT_VOID) // xxx make a macro for this
		return mv_error();

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, poosvars, ppregex_captures, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT_FOR_STRINGS(val3);
	if (val3.type != MT_STRING && val3.type != MT_VOID) {
		mv_free(&val1);
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

// ================================================================
typedef struct _rval_evaluator_field_name_state_t {
	char* field_name;
} rval_evaluator_field_name_state_t;

mv_t rval_evaluator_field_name_func_string_only(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_void();
		} else {
			// string points into lrec memory and is valid as long as the lrec is.
			return mv_from_string_no_free(string);
		}
	}
}

mv_t rval_evaluator_field_name_func_string_float(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_void();
		} else {
			double fltv;
			if (mlr_try_float_from_string(string, &fltv)) {
				return mv_from_float(fltv);
			} else {
				// string points into lrec memory and is valid as long as the lrec is.
				return mv_from_string_no_free(string);
			}
		}
	}
}

mv_t rval_evaluator_field_name_func_string_float_int(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate) {
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_void();
		} else {
			long long intv;
			double fltv;
			if (mlr_try_int_from_string(string, &intv)) {
				return mv_from_int(intv);
			} else if (mlr_try_float_from_string(string, &fltv)) {
				return mv_from_float(fltv);
			} else {
				// string points into AST memory and is valid as long as the AST is.
				return mv_from_string_no_free(string);
			}
		}
	}
}
static void rval_evaluator_field_name_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_field_name_state_t* pstate = pevaluator->pvstate;
	free(pstate->field_name);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_field_name(char* field_name, int type_inferencing) {
	rval_evaluator_field_name_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_field_name_state_t));
	pstate->field_name = mlr_strdup_or_die(field_name);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	switch (type_inferencing) {
	case TYPE_INFER_STRING_ONLY:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_only;
		break;
	case TYPE_INFER_STRING_FLOAT:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_float;
		break;
	case TYPE_INFER_STRING_FLOAT_INT:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_float_int;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
		break;
	}
	pevaluator->pfree_func = rval_evaluator_field_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_oosvar_name_state_t {
	sllmv_t* pmvkeys;
} rval_evaluator_oosvar_name_state_t;

mv_t rval_evaluator_oosvar_name_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_oosvar_name_state_t* pstate = pvstate;
	int error = 0;
	mv_t* pval = mlhmmv_get(poosvars, pstate->pmvkeys, &error);
	if (pval != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy a value here to feed
		// into that. Otherwise the typed-overlay map in mapper_put would have its contents freed out from underneath it
		// by the evaluator functions.
		if (pval->type == MT_STRING && *pval->u.strv == 0)
			return mv_void();
		else
			return mv_copy(pval);
	} else {
		return mv_uninit();
	}
}

static void rval_evaluator_oosvar_name_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_oosvar_name_state_t* pstate = pevaluator->pvstate;
	free(pstate);
	free(pevaluator);
}

// This is used for evaluating @-variables that don't have brackets: e.g. @x vs. @x[$1].
// See comments above rval_evaluator_alloc_from_oosvar_level_keys for more information.
rval_evaluator_t* rval_evaluator_alloc_from_oosvar_name(char* oosvar_name) {
	rval_evaluator_oosvar_name_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_oosvar_name_state_t));
	mv_t mv_name = mv_from_string(oosvar_name, NO_FREE);
	pstate->pmvkeys = sllmv_single(&mv_name);
	mv_free(&mv_name);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_oosvar_name_func;
	pevaluator->pfree_func = rval_evaluator_oosvar_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_oosvar_level_keys_state_t {
	sllv_t* poosvar_rhs_keylist_evaluators;
} rval_evaluator_oosvar_level_keys_state_t;

mv_t rval_evaluator_oosvar_level_keys_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_oosvar_level_keys_state_t* pstate = pvstate;

	sllmv_t* pmvkeys = sllmv_alloc();
	int keys_ok = TRUE;
	for (sllve_t* pe = pstate->poosvar_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pmvkey_evaluator = pe->pvvalue;
		mv_t mvkey = pmvkey_evaluator->pprocess_func(prec, ptyped_overlay,
			poosvars, ppregex_captures, pctx, pmvkey_evaluator->pvstate);
		if (mv_is_null(&mvkey)) {
			keys_ok = FALSE;
			break;
		}
		// Don't free the mlrval since its memory will be managed by the sllmv.
		sllmv_add(pmvkeys, &mvkey);
	}

	mv_t rv = mv_uninit();
	if (keys_ok) {
		int error = 0;
		mv_t* pval = mlhmmv_get(poosvars, pmvkeys, &error);
		if (pval != NULL) {
			if (pval->type == MT_STRING && *pval->u.strv == 0)
				rv = mv_void();
			else
				rv = *pval;
		}
	}

	sllmv_free(pmvkeys);
	return rv;
}

static void rval_evaluator_oosvar_level_keys_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_oosvar_level_keys_state_t* pstate = pevaluator->pvstate;
	for (sllve_t* pe = pstate->poosvar_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pevaluator = pe->pvvalue;
		pevaluator->pfree_func(pevaluator);
	}
	free(pstate);
	free(pevaluator);
}

// ----------------------------------------------------------------
// Example AST:
//
// $ mlr put -v '$y = @x[1]["two"][$3+4][@5]' /dev/null
// = (srec_assignment):
//     y (field_name).
//     [] (oosvar_level_key):
//         [] (oosvar_level_key):
//             [] (oosvar_level_key):
//                 [] (oosvar_level_key):
//                     x (oosvar_name).
//                     1 (strnum_literal).
//                 two (strnum_literal).
//             + (operator):
//                 3 (field_name).
//                 4 (strnum_literal).
//         5 (oosvar_name).
//
// Here past is the =; pright is the 5; pleft is the string of bracket references
// ending at the oosvar name.
//
// The job of this allocator is to set up a linked list of evaluators, with the first position for the oosvar name,
// and the rest for each of the bracketed expressions.  This is used for when there *are* brackets; see
// rval_evaluator_alloc_from_oosvar_name for when there are no brackets.

rval_evaluator_t* rval_evaluator_alloc_from_oosvar_level_keys(mlr_dsl_ast_node_t* past) {
	rval_evaluator_oosvar_level_keys_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_oosvar_level_keys_state_t));

	sllv_t* poosvar_rhs_keylist_evaluators = sllv_alloc();
	mlr_dsl_ast_node_t* pnode = past;
	while (TRUE) {
		// Bracket operators come in from the right. So the highest AST node is the rightmost
		// map, and the lowest is the oosvar name. Hence sllv_prepend rather than sllv_append.
		if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			mlr_dsl_ast_node_t* pkeynode = pnode->pchildren->phead->pnext->pvvalue;
			sllv_prepend(poosvar_rhs_keylist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, TYPE_INFER_STRING_FLOAT_INT));
		} else {
			// Oosvar expressions are of the form '@name[$index1][@index2+3][4]["five"].  The first one (name) is
			// special: syntactically, it's outside the brackets, although that issue is for the parser to handle.
			// Here it's special since it's always a string, never an expression that evaluates to string.
			// Yet for the mlhmmv the first key isn't special.
			sllv_prepend(poosvar_rhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnode->text)));
		}
		if (pnode->pchildren == NULL)
				break;
		pnode = pnode->pchildren->phead->pvvalue;
	}
	pstate->poosvar_rhs_keylist_evaluators = poosvar_rhs_keylist_evaluators;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_oosvar_level_keys_func;
	pevaluator->pfree_func = rval_evaluator_oosvar_level_keys_free;

	return pevaluator;
}

// ================================================================
// This is used for evaluating strings and numbers in literal expressions, e.g. '$x = "abc"'
// or '$x = "left_\1". The values are subject to replacement with regex captures. See comments
// in mapper_put for more information.
//
// Compare rval_evaluator_alloc_from_string which doesn't do regex replacement: it is intended for
// oosvar names on expression left-hand sides (outside of this file).

typedef struct _rval_evaluator_strnum_literal_state_t {
	mv_t literal;
} rval_evaluator_strnum_literal_state_t;

mv_t rval_evaluator_non_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_strnum_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

mv_t rval_evaluator_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_strnum_literal_state_t* pstate = pvstate;
	char* input = pstate->literal.u.strv;

	if (ppregex_captures == NULL || *ppregex_captures == NULL) {
		return mv_from_string_no_free(input);
	} else {
		int was_allocated = FALSE;
		char* output = interpolate_regex_captures(input, *ppregex_captures, &was_allocated);
		if (was_allocated)
			return mv_from_string_with_free(output);
		else
			return mv_from_string_no_free(output);
	}
}
static void rval_evaluator_strnum_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_strnum_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_strnum_literal(char* string, int type_inferencing) {
	rval_evaluator_strnum_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_strnum_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (string == NULL) {
		pstate->literal = mv_absent();
		pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
	} else {
		long long intv;
		double fltv;

		pevaluator->pprocess_func = NULL;
		switch (type_inferencing) {
		case TYPE_INFER_STRING_ONLY:
			pstate->literal = mv_from_string_no_free(string);
			pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			break;

		case TYPE_INFER_STRING_FLOAT:
			if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			}
			break;

		case TYPE_INFER_STRING_FLOAT_INT:
			if (mlr_try_int_from_string(string, &intv)) {
				pstate->literal = mv_from_int(intv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			}
			break;
		default:
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
			break;
		}
	}
	pevaluator->pfree_func = rval_evaluator_strnum_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ----------------------------------------------------------------
// This is intended only for oosvar names on expression left-hand sides (outside of this file).
// Compare rval_evaluator_alloc_from_strnum_literal.

typedef struct _rval_evaluator_string_state_t {
	char* string;
} rval_evaluator_string_state_t;

mv_t rval_evaluator_string_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_string_state_t* pstate = pvstate;
	return mv_from_string_no_free(pstate->string);
}
static void rval_evaluator_string_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_string_state_t* pstate = pevaluator->pvstate;
	free(pstate->string);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_string(char* string) {
	rval_evaluator_string_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_string_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	pstate->string            = mlr_strdup_or_die(string);
	pevaluator->pprocess_func = rval_evaluator_string_func;
	pevaluator->pfree_func    = rval_evaluator_string_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_boolean_literal_state_t {
	mv_t literal;
} rval_evaluator_boolean_literal_state_t;

mv_t rval_evaluator_boolean_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_boolean_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

static void rval_evaluator_boolean_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_boolean_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_boolean_literal(char* string) {
	rval_evaluator_boolean_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_boolean_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (streq(string, "true")) {
		pstate->literal = mv_from_true();
	} else if (streq(string, "false")) {
		pstate->literal = mv_from_false();
	} else {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	pevaluator->pprocess_func = rval_evaluator_boolean_literal_func;
	pevaluator->pfree_func = rval_evaluator_boolean_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// Example:
// $ mlr put -v '$y=ENV["X"]' ...
// AST BEGIN STATEMENTS (0):
// AST MAIN STATEMENTS (1):
// = (srec_assignment):
//     y (field_name).
//     env (env):
//         ENV (env).
//         X (strnum_literal).
// AST END STATEMENTS (0):

// ================================================================
typedef struct _rval_evaluator_environment_state_t {
	rval_evaluator_t* pname_evaluator;
} rval_evaluator_environment_state_t;

mv_t rval_evaluator_environment_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_environment_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(prec, ptyped_overlay,
		poosvars, ppregex_captures, pctx, pstate->pname_evaluator->pvstate);
	if (mv_is_null(&mvname)) {
		return mv_absent();
	}
	char free_flags;
	char* strname = mv_format_val(&mvname, &free_flags);
	char* strvalue = getenv(strname);
	if (strvalue == NULL) {
		mv_free(&mvname);
		if (free_flags & FREE_ENTRY_VALUE)
			free(strname);
		return mv_void();
	}
	mv_t rv = mv_from_string(strvalue, NO_FREE);
	mv_free(&mvname);
	if (free_flags & FREE_ENTRY_VALUE)
		free(strname);
	return rv;
}

static void rval_evaluator_environment_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_environment_state_t* pstate = pevaluator->pvstate;
	pstate->pname_evaluator->pfree_func(pstate->pname_evaluator);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_environment(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	rval_evaluator_environment_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_environment_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pnext->pvvalue;

	pstate->pname_evaluator = rval_evaluator_alloc_from_ast(pnamenode, type_inferencing);
	pevaluator->pprocess_func = rval_evaluator_environment_func;
	pevaluator->pfree_func = rval_evaluator_environment_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
mv_t rval_evaluator_NF_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(prec->field_count);
}
static void rval_evaluator_NF_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_NF() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_NF_func;
	pevaluator->pfree_func = rval_evaluator_NF_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_NR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->nr);
}
static void rval_evaluator_NR_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_NR() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_NR_func;
	pevaluator->pfree_func = rval_evaluator_NR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FNR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->fnr);
}
static void rval_evaluator_FNR_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_FNR() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FNR_func;
	pevaluator->pfree_func = rval_evaluator_FNR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FILENAME_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_string_no_free(pctx->filename);
}
static void rval_evaluator_FILENAME_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_FILENAME() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FILENAME_func;
	pevaluator->pfree_func = rval_evaluator_FILENAME_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FILENUM_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->filenum);
}
static void rval_evaluator_FILENUM_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_FILENUM() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FILENUM_func;
	pevaluator->pfree_func = rval_evaluator_FILENUM_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_PI_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_float(M_PI);
}
static void rval_evaluator_PI_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_PI() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_PI_func;
	pevaluator->pfree_func = rval_evaluator_PI_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_E_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_float(M_E);
}
static void rval_evaluator_E_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_E() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_E_func;
	pevaluator->pfree_func = rval_evaluator_E_free;
	return pevaluator;
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_context_variable(char* variable_name) {
	if        (streq(variable_name, "NF"))       { return rval_evaluator_alloc_from_NF();
	} else if (streq(variable_name, "NR"))       { return rval_evaluator_alloc_from_NR();
	} else if (streq(variable_name, "FNR"))      { return rval_evaluator_alloc_from_FNR();
	} else if (streq(variable_name, "FILENAME")) { return rval_evaluator_alloc_from_FILENAME();
	} else if (streq(variable_name, "FILENUM"))  { return rval_evaluator_alloc_from_FILENUM();
	} else if (streq(variable_name, "PI"))       { return rval_evaluator_alloc_from_PI();
	} else if (streq(variable_name, "E"))        { return rval_evaluator_alloc_from_E();
	} else  { return NULL;
	}
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_zary_func_name(char* function_name) {
	if        (streq(function_name, "urand")) {
		return rval_evaluator_alloc_from_x_z_func(f_z_urand_func);
	} else if (streq(function_name, "urand32")) {
		return rval_evaluator_alloc_from_x_z_func(i_z_urand32_func);
	} else if (streq(function_name, "systime")) {
		return rval_evaluator_alloc_from_x_z_func(f_z_systime_func);
	} else  {
		return NULL;
	}
}

// ================================================================
typedef enum _func_class_t {
	FUNC_CLASS_ARITHMETIC,
	FUNC_CLASS_MATH,
	FUNC_CLASS_BOOLEAN,
	FUNC_CLASS_STRING,
	FUNC_CLASS_CONVERSION,
	FUNC_CLASS_TIME
} func_class_t;

typedef struct _function_lookup_t {
	func_class_t function_class;
	char*        function_name;
	int          arity;
	char*        usage_string;
} function_lookup_t;

static function_lookup_t FUNCTION_LOOKUP_TABLE[] = {

	{ FUNC_CLASS_ARITHMETIC, "+",  2 , "Addition."},
	{ FUNC_CLASS_ARITHMETIC, "+",  1 , "Unary plus."},
	{ FUNC_CLASS_ARITHMETIC, "-",  2 , "Subtraction."},
	{ FUNC_CLASS_ARITHMETIC, "-",  1 , "Unary minus."},
	{ FUNC_CLASS_ARITHMETIC, "*",  2 , "Multiplication."},
	{ FUNC_CLASS_ARITHMETIC, "/",  2 , "Division."},
	{ FUNC_CLASS_ARITHMETIC, "//", 2 , "Integer division: rounds to negative (pythonic)."},
	{ FUNC_CLASS_ARITHMETIC, "%",  2 , "Remainder; never negative-valued (pythonic)."},
	{ FUNC_CLASS_ARITHMETIC, "**", 2 , "Exponentiation; same as pow, but as an infix\noperator."},
	{ FUNC_CLASS_ARITHMETIC, "|",  2 , "Bitwise OR."},
	{ FUNC_CLASS_ARITHMETIC, "^",  2 , "Bitwise XOR."},
	{ FUNC_CLASS_ARITHMETIC, "&",  2 , "Bitwise AND."},
	{ FUNC_CLASS_ARITHMETIC, "~",  1 , "Bitwise NOT. Beware '$y=~$x' since =~ is the\nregex-match operator: try '$y = ~$x'."},
	{ FUNC_CLASS_ARITHMETIC, "<<", 2 , "Bitwise left-shift."},
	{ FUNC_CLASS_ARITHMETIC, ">>", 2 , "Bitwise right-shift."},

	{ FUNC_CLASS_BOOLEAN, "==",      2 , "String/numeric equality. Mixing number and string\nresults in string compare."},
	{ FUNC_CLASS_BOOLEAN, "!=",      2 , "String/numeric inequality. Mixing number and string\nresults in string compare."},
	{ FUNC_CLASS_BOOLEAN, "=~",      2 , "String (left-hand side) matches regex (right-hand\nside), e.g. '$name =~ \"^a.*b$\"'."},
	{ FUNC_CLASS_BOOLEAN, "!=~",     2 , "String (left-hand side) does not match regex\n(right-hand side), e.g. '$name !=~ \"^a.*b$\"'."},
	{ FUNC_CLASS_BOOLEAN, ">",       2 , "String/numeric greater-than. Mixing number and string\nresults in string compare."},
	{ FUNC_CLASS_BOOLEAN, ">=",      2 , "String/numeric greater-than-or-equals. Mixing number\nand string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "<",       2 , "String/numeric less-than. Mixing number and string\nresults in string compare."},
	{ FUNC_CLASS_BOOLEAN, "<=",      2 , "String/numeric less-than-or-equals. Mixing number\nand string results in string compare."},
	{ FUNC_CLASS_BOOLEAN, "&&",      2 , "Logical AND."},
	{ FUNC_CLASS_BOOLEAN, "||",      2 , "Logical OR."},
	{ FUNC_CLASS_BOOLEAN, "^^",      2 , "Logical XOR."},
	{ FUNC_CLASS_BOOLEAN, "!",       1 , "Logical negation."},
	{ FUNC_CLASS_BOOLEAN, "? :",     3 , "Ternary operator."},

	{ FUNC_CLASS_CONVERSION, "isnull",    1 , "True if argument is null, false otherwise"},
	{ FUNC_CLASS_CONVERSION, "isnotnull", 1 , "False if argument is null, true otherwise."},
	{ FUNC_CLASS_CONVERSION, "boolean",   1 , "Convert int/float/bool/string to boolean."},
	{ FUNC_CLASS_CONVERSION, "float",     1 , "Convert int/float/bool/string to float."},
	{ FUNC_CLASS_CONVERSION, "fmtnum",    2 , "Convert int/float/bool to string using\nprintf-style format string, e.g. \"%06lld\"."},
	{ FUNC_CLASS_CONVERSION, "hexfmt",    1 , "Convert int to string, e.g. 255 to \"0xff\"."},
	{ FUNC_CLASS_CONVERSION, "int",       1 , "Convert int/float/bool/string to int."},
	{ FUNC_CLASS_CONVERSION, "string",    1 , "Convert int/float/bool/string to string."},

	{ FUNC_CLASS_STRING, ".",        2 , "String concatenation."},
	{ FUNC_CLASS_STRING, "gsub",     3 , "Example: '$name=gsub($name, \"old\", \"new\")'\n(replace all)."},
	{ FUNC_CLASS_STRING, "strlen",   1 , "String length."},
	{ FUNC_CLASS_STRING, "sub",      3 , "Example: '$name=sub($name, \"old\", \"new\")'\n(replace once)."},
	{ FUNC_CLASS_STRING, "tolower",  1 , "Convert string to lowercase."},
	{ FUNC_CLASS_STRING, "toupper",  1 , "Convert string to uppercase."},

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
	{ FUNC_CLASS_MATH, "invqnorm", 1 , "Inverse of normal cumulative distribution\nfunction. Note that invqorm(urand()) is normally distributed."},
	{ FUNC_CLASS_MATH, "log",      1 , "Natural (base-e) logarithm."},
	{ FUNC_CLASS_MATH, "log10",    1 , "Base-10 logarithm."},
	{ FUNC_CLASS_MATH, "log1p",    1 , "log(1-x)."},
	{ FUNC_CLASS_MATH, "logifit",  3 , "Given m and b from logistic regression, compute\nfit: $yhat=logifit($x,$m,$b)."},
	{ FUNC_CLASS_MATH, "madd",     3 , "a + b mod m (integers)"},
	{ FUNC_CLASS_MATH, "max",      2 , "max of two numbers; null loses"},
	{ FUNC_CLASS_MATH, "mexp",     3 , "a ** b mod m (integers)"},
	{ FUNC_CLASS_MATH, "min",      2 , "min of two numbers; null loses"},
	{ FUNC_CLASS_MATH, "mmul",     3 , "a * b mod m (integers)"},
	{ FUNC_CLASS_MATH, "msub",     3 , "a - b mod m (integers)"},
	{ FUNC_CLASS_MATH, "pow",      2 , "Exponentiation; same as **."},
	{ FUNC_CLASS_MATH, "qnorm",    1 , "Normal cumulative distribution function."},
	{ FUNC_CLASS_MATH, "round",    1 , "Round to nearest integer."},
	{ FUNC_CLASS_MATH, "roundm",   2 , "Round to nearest multiple of m: roundm($x,$m) is\nthe same as round($x/$m)*$m"},
	{ FUNC_CLASS_MATH, "sgn",      1 , "+1 for positive input, 0 for zero input, -1 for\nnegative input."},
	{ FUNC_CLASS_MATH, "sin",      1 , "Trigonometric sine."},
	{ FUNC_CLASS_MATH, "sinh",     1 , "Hyperbolic sine."},
	{ FUNC_CLASS_MATH, "sqrt",     1 , "Square root."},
	{ FUNC_CLASS_MATH, "tan",      1 , "Trigonometric tangent."},
	{ FUNC_CLASS_MATH, "tanh",     1 , "Hyperbolic tangent."},
	{ FUNC_CLASS_MATH, "urand",    0 , "Floating-point numbers on the unit interval.\nInt-valued example: '$n=floor(20+urand()*11)'." },
	{ FUNC_CLASS_MATH, "urand32",  0 , "Integer uniformly distributed 0 and 2**32-1\ninclusive." },
	{ FUNC_CLASS_MATH, "urandint", 2 , "Integer uniformly distributed between inclusive\ninteger endpoints." },

	{ FUNC_CLASS_TIME, "dhms2fsec", 1 , "Recovers floating-point seconds as in\ndhms2fsec(\"5d18h53m20.250000s\") = 500000.250000"},
	{ FUNC_CLASS_TIME, "dhms2sec",  1 , "Recovers integer seconds as in\ndhms2sec(\"5d18h53m20s\") = 500000"},
	{ FUNC_CLASS_TIME, "fsec2dhms", 1 , "Formats floating-point seconds as in\nfsec2dhms(500000.25) = \"5d18h53m20.250000s\""},
	{ FUNC_CLASS_TIME, "fsec2hms",  1 , "Formats floating-point seconds as in\nfsec2hms(5000.25) = \"01:23:20.250000\""},
	{ FUNC_CLASS_TIME, "gmt2sec",   1 , "Parses GMT timestamp as integer seconds since\nthe epoch."},
	{ FUNC_CLASS_TIME, "hms2fsec",  1 , "Recovers floating-point seconds as in\nhms2fsec(\"01:23:20.250000\") = 5000.250000"},
	{ FUNC_CLASS_TIME, "hms2sec",   1 , "Recovers integer seconds as in\nhms2sec(\"01:23:20\") = 5000"},
	{ FUNC_CLASS_TIME, "sec2dhms",  1 , "Formats integer seconds as in sec2dhms(500000)\n= \"5d18h53m20s\""},
	{ FUNC_CLASS_TIME, "sec2gmt",   1 , "Formats seconds since epoch (integer part)\nas GMT timestamp, e.g. sec2gmt(1440768801.7) = \"2015-08-28T13:33:21Z\"."},
	{ FUNC_CLASS_TIME, "sec2hms",   1 , "Formats integer seconds as in\nsec2hms(5000) = \"01:23:20\""},
	{ FUNC_CLASS_TIME, "strftime",  2 , "Formats seconds since epoch (integer part)\nas timestamp, e.g.\nstrftime(1440768801.7,\"%Y-%m-%dT%H:%M:%SZ\") = \"2015-08-28T13:33:21Z\"."},
	{ FUNC_CLASS_TIME, "strptime",  2 , "Parses timestamp as integer seconds since epoch,\ne.g. strptime(\"2015-08-28T13:33:21Z\",\"%Y-%m-%dT%H:%M:%SZ\") = 1440768801."},
	{ FUNC_CLASS_TIME, "systime",   0 , "Floating-point seconds since the epoch,\ne.g. 1440768801.748936." },

	{  0, NULL,      -1 , NULL}, // table terminator
};

typedef enum _arity_check_t {
	ARITY_CHECK_PASS,
	ARITY_CHECK_FAIL,
	ARITY_CHECK_NO_SUCH
} arity_check_t;

static arity_check_t check_arity(function_lookup_t lookup_table[], char* function_name,
	int user_provided_arity, int *parity)
{
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

static void check_arity_with_report(function_lookup_t fcn_lookup_table[], char* function_name,
	int user_provided_arity)
{
	int arity = -1;
	arity_check_t result = check_arity(fcn_lookup_table, function_name, user_provided_arity, &arity);
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

static char* function_class_to_desc(func_class_t function_class) {
	switch(function_class) {
	case FUNC_CLASS_ARITHMETIC: return "arithmetic"; break;
	case FUNC_CLASS_MATH:       return "math";       break;
	case FUNC_CLASS_BOOLEAN:    return "boolean";    break;
	case FUNC_CLASS_STRING:     return "string";     break;
	case FUNC_CLASS_CONVERSION: return "conversion"; break;
	case FUNC_CLASS_TIME:       return "time";       break;
	default:                    return "???";        break;
	}
}

void rval_evaluator_list_functions(FILE* o, char* leader) {
	char* separator = " ";
	int leaderlen = strlen(leader);
	int separatorlen = strlen(separator);
	int linelen = leaderlen;
	int j = 0;

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
void rval_evaluator_function_usage(FILE* output_stream, char* function_name) {
	int found = FALSE;
	char* fmt = "%s (class=%s #args=%d): %s\n";

	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL) // end of table
			break;
		if (function_name == NULL || streq(function_name, plookup->function_name)) {
			fprintf(output_stream, fmt, plookup->function_name,
				function_class_to_desc(plookup->function_class),
				plookup->arity, plookup->usage_string);
			found = TRUE;
		}
		if (function_name == NULL)
			fprintf(output_stream, "\n");
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

void rval_evaluator_list_all_functions_raw(FILE* output_stream) {
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL) // end of table
			break;
		printf("%s\n", plookup->function_name);
	}
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_unary_func_name(char* fnnm, rval_evaluator_t* parg1)  {
	if        (streq(fnnm, "!"))         { return rval_evaluator_alloc_from_b_b_func(b_b_not_func,       parg1);
	} else if (streq(fnnm, "+"))         { return rval_evaluator_alloc_from_x_x_func(x_x_upos_func,      parg1);
	} else if (streq(fnnm, "-"))         { return rval_evaluator_alloc_from_x_x_func(x_x_uneg_func,      parg1);
	} else if (streq(fnnm, "~"))         { return rval_evaluator_alloc_from_i_i_func(i_i_bitwise_not_func, parg1);
	} else if (streq(fnnm, "abs"))       { return rval_evaluator_alloc_from_x_x_func(x_x_abs_func,       parg1);
	} else if (streq(fnnm, "acos"))      { return rval_evaluator_alloc_from_f_f_func(f_f_acos_func,      parg1);
	} else if (streq(fnnm, "acosh"))     { return rval_evaluator_alloc_from_f_f_func(f_f_acosh_func,     parg1);
	} else if (streq(fnnm, "asin"))      { return rval_evaluator_alloc_from_f_f_func(f_f_asin_func,      parg1);
	} else if (streq(fnnm, "asinh"))     { return rval_evaluator_alloc_from_f_f_func(f_f_asinh_func,     parg1);
	} else if (streq(fnnm, "atan"))      { return rval_evaluator_alloc_from_f_f_func(f_f_atan_func,      parg1);
	} else if (streq(fnnm, "atanh"))     { return rval_evaluator_alloc_from_f_f_func(f_f_atanh_func,     parg1);
	} else if (streq(fnnm, "boolean"))   { return rval_evaluator_alloc_from_x_x_func(b_x_boolean_func,   parg1);
	} else if (streq(fnnm, "boolean"))   { return rval_evaluator_alloc_from_x_x_func(b_x_boolean_func,   parg1);
	} else if (streq(fnnm, "cbrt"))      { return rval_evaluator_alloc_from_f_f_func(f_f_cbrt_func,      parg1);
	} else if (streq(fnnm, "ceil"))      { return rval_evaluator_alloc_from_x_x_func(x_x_ceil_func,      parg1);
	} else if (streq(fnnm, "cos"))       { return rval_evaluator_alloc_from_f_f_func(f_f_cos_func,       parg1);
	} else if (streq(fnnm, "cosh"))      { return rval_evaluator_alloc_from_f_f_func(f_f_cosh_func,      parg1);
	} else if (streq(fnnm, "dhms2fsec")) { return rval_evaluator_alloc_from_f_s_func(f_s_dhms2fsec_func, parg1);
	} else if (streq(fnnm, "dhms2sec"))  { return rval_evaluator_alloc_from_f_s_func(i_s_dhms2sec_func,  parg1);
	} else if (streq(fnnm, "erf"))       { return rval_evaluator_alloc_from_f_f_func(f_f_erf_func,       parg1);
	} else if (streq(fnnm, "erfc"))      { return rval_evaluator_alloc_from_f_f_func(f_f_erfc_func,      parg1);
	} else if (streq(fnnm, "exp"))       { return rval_evaluator_alloc_from_f_f_func(f_f_exp_func,       parg1);
	} else if (streq(fnnm, "expm1"))     { return rval_evaluator_alloc_from_f_f_func(f_f_expm1_func,     parg1);
	} else if (streq(fnnm, "float"))     { return rval_evaluator_alloc_from_x_x_func(f_x_float_func,     parg1);
	} else if (streq(fnnm, "floor"))     { return rval_evaluator_alloc_from_x_x_func(x_x_floor_func,     parg1);
	} else if (streq(fnnm, "fsec2dhms")) { return rval_evaluator_alloc_from_s_f_func(s_f_fsec2dhms_func, parg1);
	} else if (streq(fnnm, "fsec2hms"))  { return rval_evaluator_alloc_from_s_f_func(s_f_fsec2hms_func,  parg1);
	} else if (streq(fnnm, "gmt2sec"))   { return rval_evaluator_alloc_from_i_s_func(i_s_gmt2sec_func,   parg1);
	} else if (streq(fnnm, "hexfmt"))    { return rval_evaluator_alloc_from_x_x_func(s_x_hexfmt_func,    parg1);
	} else if (streq(fnnm, "hms2fsec"))  { return rval_evaluator_alloc_from_f_s_func(f_s_hms2fsec_func,  parg1);
	} else if (streq(fnnm, "hms2sec"))   { return rval_evaluator_alloc_from_f_s_func(i_s_hms2sec_func,   parg1);
	} else if (streq(fnnm, "int"))       { return rval_evaluator_alloc_from_x_x_func(i_x_int_func,       parg1);
	} else if (streq(fnnm, "invqnorm"))  { return rval_evaluator_alloc_from_f_f_func(f_f_invqnorm_func,  parg1);
	} else if (streq(fnnm, "isnotnull")) { return rval_evaluator_alloc_from_x_x_func(b_x_isnotnull_func, parg1);
	} else if (streq(fnnm, "isnull"))    { return rval_evaluator_alloc_from_x_x_func(b_x_isnull_func,    parg1);
	} else if (streq(fnnm, "log"))       { return rval_evaluator_alloc_from_f_f_func(f_f_log_func,       parg1);
	} else if (streq(fnnm, "log10"))     { return rval_evaluator_alloc_from_f_f_func(f_f_log10_func,     parg1);
	} else if (streq(fnnm, "log1p"))     { return rval_evaluator_alloc_from_f_f_func(f_f_log1p_func,     parg1);
	} else if (streq(fnnm, "qnorm"))     { return rval_evaluator_alloc_from_f_f_func(f_f_qnorm_func,     parg1);
	} else if (streq(fnnm, "round"))     { return rval_evaluator_alloc_from_x_x_func(x_x_round_func,     parg1);
	} else if (streq(fnnm, "sec2dhms"))  { return rval_evaluator_alloc_from_s_i_func(s_i_sec2dhms_func,  parg1);
	} else if (streq(fnnm, "sec2gmt"))   { return rval_evaluator_alloc_from_x_n_func(s_n_sec2gmt_func,   parg1);
	} else if (streq(fnnm, "sec2hms"))   { return rval_evaluator_alloc_from_s_i_func(s_i_sec2hms_func,   parg1);
	} else if (streq(fnnm, "sgn"))       { return rval_evaluator_alloc_from_x_x_func(x_x_sgn_func,       parg1);
	} else if (streq(fnnm, "sin"))       { return rval_evaluator_alloc_from_f_f_func(f_f_sin_func,       parg1);
	} else if (streq(fnnm, "sinh"))      { return rval_evaluator_alloc_from_f_f_func(f_f_sinh_func,      parg1);
	} else if (streq(fnnm, "sqrt"))      { return rval_evaluator_alloc_from_f_f_func(f_f_sqrt_func,      parg1);
	} else if (streq(fnnm, "string"))    { return rval_evaluator_alloc_from_x_x_func(s_x_string_func,    parg1);
	} else if (streq(fnnm, "strlen"))    { return rval_evaluator_alloc_from_i_s_func(i_s_strlen_func,    parg1);
	} else if (streq(fnnm, "tan"))       { return rval_evaluator_alloc_from_f_f_func(f_f_tan_func,       parg1);
	} else if (streq(fnnm, "tanh"))      { return rval_evaluator_alloc_from_f_f_func(f_f_tanh_func,      parg1);
	} else if (streq(fnnm, "tolower"))   { return rval_evaluator_alloc_from_s_s_func(s_s_tolower_func,   parg1);
	} else if (streq(fnnm, "toupper"))   { return rval_evaluator_alloc_from_s_s_func(s_s_toupper_func,   parg1);

	} else return NULL;
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_binary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2)
{
	if        (streq(fnnm, "&&"))     { return rval_evaluator_alloc_from_b_bb_and_func(parg1, parg2);
	} else if (streq(fnnm, "||"))     { return rval_evaluator_alloc_from_b_bb_or_func (parg1, parg2);
	} else if (streq(fnnm, "^^"))     { return rval_evaluator_alloc_from_b_bb_xor_func(parg1, parg2);
	} else if (streq(fnnm, "=~"))     { return rval_evaluator_alloc_from_x_ssc_func(matches_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "!=~"))    { return rval_evaluator_alloc_from_x_ssc_func(does_not_match_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "=="))     { return rval_evaluator_alloc_from_x_xx_func(eq_op_func,             parg1, parg2);
	} else if (streq(fnnm, "!="))     { return rval_evaluator_alloc_from_x_xx_func(ne_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">"))      { return rval_evaluator_alloc_from_x_xx_func(gt_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">="))     { return rval_evaluator_alloc_from_x_xx_func(ge_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<"))      { return rval_evaluator_alloc_from_x_xx_func(lt_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<="))     { return rval_evaluator_alloc_from_x_xx_func(le_op_func,             parg1, parg2);
	} else if (streq(fnnm, "."))      { return rval_evaluator_alloc_from_x_ss_func(s_ss_dot_func,          parg1, parg2);
	} else if (streq(fnnm, "+"))      { return rval_evaluator_alloc_from_x_xx_func(x_xx_plus_func,         parg1, parg2);
	} else if (streq(fnnm, "-"))      { return rval_evaluator_alloc_from_x_xx_func(x_xx_minus_func,        parg1, parg2);
	} else if (streq(fnnm, "*"))      { return rval_evaluator_alloc_from_x_xx_func(x_xx_times_func,        parg1, parg2);
	} else if (streq(fnnm, "/"))      { return rval_evaluator_alloc_from_x_xx_func(x_xx_divide_func,       parg1, parg2);
	} else if (streq(fnnm, "//"))     { return rval_evaluator_alloc_from_x_xx_func(x_xx_int_divide_func,   parg1, parg2);
	} else if (streq(fnnm, "%"))      { return rval_evaluator_alloc_from_x_xx_func(x_xx_mod_func,          parg1, parg2);
	} else if (streq(fnnm, "**"))     { return rval_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "pow"))    { return rval_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "atan2"))  { return rval_evaluator_alloc_from_f_ff_func(f_ff_atan2_func,        parg1, parg2);
	} else if (streq(fnnm, "max"))    { return rval_evaluator_alloc_from_x_xx_nullable_func(x_xx_max_func, parg1, parg2);
	} else if (streq(fnnm, "min"))    { return rval_evaluator_alloc_from_x_xx_nullable_func(x_xx_min_func, parg1, parg2);
	} else if (streq(fnnm, "roundm")) { return rval_evaluator_alloc_from_x_xx_func(x_xx_roundm_func,       parg1, parg2);
	} else if (streq(fnnm, "fmtnum")) { return rval_evaluator_alloc_from_s_xs_func(s_xs_fmtnum_func,       parg1, parg2);
	} else if (streq(fnnm, "urandint")) { return rval_evaluator_alloc_from_i_ii_func(i_ii_urandint_func,   parg1, parg2);
	} else if (streq(fnnm, "|"))      { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_or_func,   parg1, parg2);
	} else if (streq(fnnm, "^"))      { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_xor_func,  parg1, parg2);
	} else if (streq(fnnm, "&"))      { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_and_func,  parg1, parg2);
	} else if (streq(fnnm, "<<"))     { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_lsh_func,  parg1, parg2);
	} else if (streq(fnnm, ">>"))     { return rval_evaluator_alloc_from_i_ii_func(i_ii_bitwise_rsh_func,  parg1, parg2);
	} else if (streq(fnnm, "strftime")) { return rval_evaluator_alloc_from_x_ns_func(s_ns_strftime_func,   parg1, parg2);
	} else if (streq(fnnm, "strptime")) { return rval_evaluator_alloc_from_x_ss_func(i_ss_strptime_func,   parg1, parg2);
	} else  { return NULL; }
}

rval_evaluator_t* rval_evaluator_alloc_from_binary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	if        (streq(fnnm, "=~"))  {
		return rval_evaluator_alloc_from_x_sr_func(matches_precomp_func,        parg1, regex_string, ignore_case);
	} else if (streq(fnnm, "!=~")) {
		return rval_evaluator_alloc_from_x_sr_func(does_not_match_precomp_func, parg1, regex_string, ignore_case);
	} else  { return NULL; }
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_ternary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3)
{
	if (streq(fnnm, "sub")) {
		return rval_evaluator_alloc_from_s_sss_func(sub_no_precomp_func,  parg1, parg2, parg3);
	} else if (streq(fnnm, "gsub")) {
		return rval_evaluator_alloc_from_s_sss_func(gsub_no_precomp_func, parg1, parg2, parg3);
	} else if (streq(fnnm, "logifit")) {
		return rval_evaluator_alloc_from_f_fff_func(f_fff_logifit_func,   parg1, parg2, parg3);
	} else if (streq(fnnm, "madd")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modadd_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "msub")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modsub_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mmul")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modmul_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mexp")) {
		return rval_evaluator_alloc_from_i_iii_func(i_iii_modexp_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "? :")) {
		return rval_evaluator_alloc_from_ternop(parg1, parg2, parg3);
	} else  { return NULL; }
}

rval_evaluator_t* rval_evaluator_alloc_from_ternary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3)
{
	if (streq(fnnm, "sub"))  {
		return rval_evaluator_alloc_from_x_srs_func(sub_precomp_func,  parg1, regex_string, ignore_case, parg3);
	} else if (streq(fnnm, "gsub"))  {
		return rval_evaluator_alloc_from_x_srs_func(gsub_precomp_func, parg1, regex_string, ignore_case, parg3);
	} else  { return NULL; }
}

// ================================================================
static rval_evaluator_t* rval_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	int type_inferencing, function_lookup_t* fcn_lookup_table)
{
	if (pnode->pchildren == NULL) { // leaf node
		if (pnode->type == MD_AST_NODE_TYPE_FIELD_NAME) {
			return rval_evaluator_alloc_from_field_name(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			return rval_evaluator_alloc_from_oosvar_name(pnode->text);
		} else if (pnode->type == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
			return rval_evaluator_alloc_from_strnum_literal(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_BOOLEAN_LITERAL) {
			return rval_evaluator_alloc_from_boolean_literal(pnode->text);
		} else if (pnode->type == MD_AST_NODE_TYPE_REGEXI) {
			return rval_evaluator_alloc_from_strnum_literal(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_CONTEXT_VARIABLE) {
			return rval_evaluator_alloc_from_context_variable(pnode->text);
		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

	} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		return rval_evaluator_alloc_from_oosvar_level_keys(pnode);

	} else if (pnode->type == MD_AST_NODE_TYPE_ENV) {
		return rval_evaluator_alloc_from_environment(pnode, type_inferencing);

	} else { // operator/function
		if ((pnode->type != MD_AST_NODE_TYPE_FUNCTION_NAME)
		&& (pnode->type != MD_AST_NODE_TYPE_OPERATOR)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d (node type %s).\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__, mlr_dsl_ast_node_describe_type(pnode->type));
			exit(1);
		}
		char* func_name = pnode->text;

		int user_provided_arity = pnode->pchildren->length;

		check_arity_with_report(fcn_lookup_table, func_name, user_provided_arity);

		rval_evaluator_t* pevaluator = NULL;
		if (user_provided_arity == 0) {
			pevaluator = rval_evaluator_alloc_from_zary_func_name(func_name);
		} else if (user_provided_arity == 1) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
			pevaluator = rval_evaluator_alloc_from_unary_func_name(func_name, parg1);
		} else if (user_provided_arity == 2) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_regex_arg2_func_name(func_name,
					parg1, parg2_node->text, FALSE);
			} else if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MD_AST_NODE_TYPE_REGEXI) {
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT);
			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_func_name(func_name, parg1, parg2);
			}

		} else if (user_provided_arity == 3) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
				// sub/gsub-regex special case:
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					FALSE, parg3);

			} else if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MD_AST_NODE_TYPE_REGEXI) {
				// sub/gsub-regex special case:
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT, parg3);

			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_func_name(func_name, parg1, parg2, parg3);
			}
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

rval_evaluator_t* rval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	return rval_evaluator_alloc_from_ast_aux(pnode, type_inferencing, FUNCTION_LOOKUP_TABLE);
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
	printf("-- TEST_RVAL_EVALUATORS test1 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* pnr       = rval_evaluator_alloc_from_NR();
	rval_evaluator_t* pfnr      = rval_evaluator_alloc_from_FNR();
	rval_evaluator_t* pfilename = rval_evaluator_alloc_from_FILENAME();
	rval_evaluator_t* pfilenum  = rval_evaluator_alloc_from_FILENUM();

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;

	mv_t val = pnr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 888);

	val = pfnr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 999);

	val = pfilename->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfilename->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "filename-goes-here"));

	val = pfilenum->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pfilenum->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 123);

	return 0;
}

// ----------------------------------------------------------------
static char * test2() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test2 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* ps       = rval_evaluator_alloc_from_field_name("s", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* pdef     = rval_evaluator_alloc_from_strnum_literal("def", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* pdot     = rval_evaluator_alloc_from_x_ss_func(s_ss_dot_func, ps, pdef);
	rval_evaluator_t* ptolower = rval_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	rval_evaluator_t* ptoupper = rval_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "s", "abc", NO_FREE);
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	mv_t val = ps->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ps->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abc"));

	val = pdef->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pdef->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "def"));

	val = pdot->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pdot->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptolower->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ptolower->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptoupper->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, ptoupper->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "ABCDEF"));

	return 0;
}

// ----------------------------------------------------------------
static char * test3() {
	printf("\n");
	printf("-- TEST_RVAL_EVALUATORS test3 ENTER\n");
	context_t ctx = {.nr = 888, .fnr = 999, .filenum = 123, .filename = "filename-goes-here"};
	context_t* pctx = &ctx;

	rval_evaluator_t* p2     = rval_evaluator_alloc_from_strnum_literal("2.0", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* px     = rval_evaluator_alloc_from_field_name("x", TYPE_INFER_STRING_FLOAT_INT);
	rval_evaluator_t* plogx  = rval_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	rval_evaluator_t* p2logx = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, plogx);
	rval_evaluator_t* px2    = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, px, px);
	rval_evaluator_t* p4     = rval_evaluator_alloc_from_x_xx_func(x_xx_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MD_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MD_AST_NODE_TYPE_FUNCTION_NAME);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MD_AST_NODE_TYPE_STRNUM_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	rval_evaluator_t*  pastr = rval_evaluator_alloc_from_ast(p2logxnode, TYPE_INFER_STRING_FLOAT_INT);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	mlhmmv_t* poosvars = mlhmmv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "x", "4.5", NO_FREE);

	mv_t valp2     = p2->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p2->pvstate);
	mv_t valp4     = p4->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p4->pvstate);
	mv_t valpx     = px->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, px->pvstate);
	mv_t valpx2    = px2->pprocess_func(prec,    ptyped_overlay, poosvars, &pregex_captures, pctx, px2->pvstate);
	mv_t valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, poosvars, &pregex_captures, pctx, plogx->pvstate);
	mv_t valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, p2logx->pvstate);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

	printf("XXX %s\n", mt_describe_type(valp2.type));
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
	printf("newval AST      = %s\n",
		mv_describe_val(pastr->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, pastr->pvstate)));
	printf("\n");

	lrec_rename(prec, "x", "y", FALSE);

	valp2     = p2->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p2->pvstate);
	valp4     = p4->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, p4->pvstate);
	valpx     = px->pprocess_func(prec,     ptyped_overlay, poosvars, &pregex_captures, pctx, px->pvstate);
	valpx2    = px2->pprocess_func(prec,    ptyped_overlay, poosvars, &pregex_captures, pctx, px2->pvstate);
	valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, poosvars, &pregex_captures, pctx, plogx->pvstate);
	valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, poosvars, &pregex_captures, pctx, p2logx->pvstate);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

	mu_assert_lf(valp2.type     == MT_FLOAT);
	mu_assert_lf(valp4.type     == MT_FLOAT);
	mu_assert_lf(valpx.type     == MT_ABSENT);
	mu_assert_lf(valpx2.type    == MT_ERROR);
	mu_assert_lf(valplogx.type  == MT_ABSENT);
	mu_assert_lf(valp2logx.type == MT_ERROR);

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

// test_rval_evaluators has the MinUnit inside rval_evaluators, as it tests
// many private methods. (The other option is to make them all public.)
int test_rval_evaluators_main(int argc, char **argv) {
	mlr_global_init(argv[0], "%lf", NULL);

	printf("TEST_RVAL_EVALUATORS ENTER\n");
	char *result = all_tests();
	printf("\n");
	if (result != 0) {
		printf("Not all unit tests passed\n");
	}
	else {
		printf("TEST_RVAL_EVALUATORS: ALL UNIT TESTS PASSED\n");
	}
	printf("Tests      passed: %d of %d\n", tests_run - tests_failed, tests_run);
	printf("Assertions passed: %d of %d\n", assertions_run - assertions_failed, assertions_run);

	return result != 0;
}
