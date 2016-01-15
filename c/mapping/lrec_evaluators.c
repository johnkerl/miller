#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
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
//   lrec_evaluators.c to invoke functions here with mlrvals of the correct
//   type(s).
// ================================================================

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_b_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_b_b_state_t;

mv_t lrec_evaluator_b_b_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_b_b_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_b_b_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_b_b_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_b_b_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_b_b_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_b_b_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_b_b_func;
	pevaluator->pfree_func = lrec_evaluator_b_b_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_bb_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_bb_state_t;

mv_t lrec_evaluator_b_bb_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_b_bb_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_BOOL)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);

	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_BOOL)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_b_bb_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_b_bb_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_b_bb_func;
	pevaluator->pfree_func = lrec_evaluator_b_bb_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_z_state_t {
	mv_zary_func_t* pfunc;
} lrec_evaluator_x_z_state_t;

mv_t lrec_evaluator_x_z_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_z_state_t* pstate = pvstate;

	return pstate->pfunc();
}
static void lrec_evaluator_x_z_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_z_func(mv_zary_func_t* pfunc) {
	lrec_evaluator_x_z_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_z_state_t));
	pstate->pfunc = pfunc;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_z_func;
	pevaluator->pfree_func = lrec_evaluator_x_z_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_f_state_t {
	mv_unary_func_t* pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_f_state_t;

mv_t lrec_evaluator_f_f_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_f_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_FLOAT)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_f_f_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_f_f_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_f_f_func;
	pevaluator->pfree_func = lrec_evaluator_f_f_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_n_state_t {
	mv_unary_func_t* pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_x_n_state_t;

mv_t lrec_evaluator_x_n_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_n_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_x_n_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_n_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_n_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_x_n_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_n_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_n_func;
	pevaluator->pfree_func = lrec_evaluator_x_n_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_i_state_t {
	mv_unary_func_t* pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_i_i_state_t;

mv_t lrec_evaluator_i_i_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_i_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_i_i_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_i_i_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_i_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_i_i_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_i_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_i_i_func;
	pevaluator->pfree_func = lrec_evaluator_i_i_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_ff_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_f_ff_state_t;

mv_t lrec_evaluator_f_ff_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_f_ff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_float_nullable(&val2);
	NULL_OR_ERROR_OUT(val2);

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_f_ff_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_f_ff_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_f_ff_func;
	pevaluator->pfree_func = lrec_evaluator_f_ff_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_n_nn_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_n_nn_state_t;

mv_t lrec_evaluator_n_nn_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_n_nn_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);

	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	mv_set_number_nullable(&val2);
	NULL_OR_ERROR_OUT(val2);

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_n_nn_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_n_nn_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_n_nn_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_n_nn_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_n_nn_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_n_nn_func;
	pevaluator->pfree_func = lrec_evaluator_n_nn_free;

	return pevaluator;
}

// ----------------------------------------------------------------
// This is for min/max which can return non-null when one argument is null --
// in comparison to other functions which return null if *any* argument is
// null.

typedef struct _lrec_evaluator_n_nn_nullable_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_n_nn_nullable_state_t;

mv_t lrec_evaluator_n_nn_nullable_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_n_nn_nullable_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_number_nullable(&val1);
	ERROR_OUT(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_number_nullable(&val2);
	ERROR_OUT(val2);

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_n_nn_nullable_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_n_nn_nullable_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_n_nn_nullable_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_n_nn_nullable_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_n_nn_nullable_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_n_nn_nullable_func;
	pevaluator->pfree_func = lrec_evaluator_n_nn_nullable_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_fff_state_t {
	mv_ternary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_f_fff_state_t;

mv_t lrec_evaluator_f_fff_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_f_fff_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	mv_set_float_nullable(&val2);
	NULL_OR_ERROR_OUT(val2);

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg3->pvstate);
	mv_set_float_nullable(&val3);
	NULL_OR_ERROR_OUT(val3);

	return pstate->pfunc(&val1, &val2, &val3);
}
static void lrec_evaluator_f_fff_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_f_fff_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_fff_func(mv_ternary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	lrec_evaluator_f_fff_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_fff_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_f_fff_func;
	pevaluator->pfree_func = lrec_evaluator_f_fff_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_ii_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_i_ii_state_t;

mv_t lrec_evaluator_i_ii_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_i_ii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mv_set_int_nullable(&val1);
	NULL_OUT(val1);
	if (val1.type != MT_INT)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	mv_set_int_nullable(&val2);
	NULL_OUT(val2);
	if (val2.type != MT_INT)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_i_ii_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_i_ii_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_ii_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_i_ii_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_ii_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_i_ii_func;
	pevaluator->pfree_func = lrec_evaluator_i_ii_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_iii_state_t {
	mv_ternary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_i_iii_state_t;

mv_t lrec_evaluator_i_iii_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_i_iii_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mv_set_int_nullable(&val1);
	NULL_OUT(val1);
	if (val1.type != MT_INT)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	mv_set_int_nullable(&val2);
	NULL_OUT(val2);
	if (val2.type != MT_INT)
		return MV_ERROR;

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT(val3);
	mv_set_int_nullable(&val3);
	NULL_OUT(val3);
	if (val3.type != MT_INT)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2, &val3);
}
static void lrec_evaluator_i_iii_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_i_iii_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_iii_func(mv_ternary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	lrec_evaluator_i_iii_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_iii_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_i_iii_func;
	pevaluator->pfree_func = lrec_evaluator_i_iii_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_ternop_state_t {
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_ternop_state_t;

mv_t lrec_evaluator_ternop_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_ternop_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mv_set_boolean_strict(&val1);

	return val1.u.boolv
		? pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate)
		: pstate->parg3->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg3->pvstate);
}
static void lrec_evaluator_ternop_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_ternop_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_ternop(lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	lrec_evaluator_ternop_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_ternop_state_t));
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;
	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_ternop_func;
	pevaluator->pfree_func = lrec_evaluator_ternop_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_s_state_t;

mv_t lrec_evaluator_s_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_s_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_s_s_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_s_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_s_s_func;
	pevaluator->pfree_func = lrec_evaluator_s_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_f_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_f_state_t;

mv_t lrec_evaluator_s_f_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_s_f_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	mv_set_float_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_s_f_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_s_f_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_f_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_f_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_s_f_func;
	pevaluator->pfree_func = lrec_evaluator_s_f_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_i_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_s_i_state_t;

mv_t lrec_evaluator_s_i_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_s_i_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);

	mv_set_int_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_s_i_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_s_i_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_s_i_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_s_i_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_s_i_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_s_i_func;
	pevaluator->pfree_func = lrec_evaluator_s_i_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_f_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_f_s_state_t;

mv_t lrec_evaluator_f_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_f_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_f_s_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_f_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_f_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_f_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_f_s_func;
	pevaluator->pfree_func = lrec_evaluator_f_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_i_s_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_i_s_state_t;

mv_t lrec_evaluator_i_s_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_i_s_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_i_s_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_i_s_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_i_s_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_i_s_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_i_s_func;
	pevaluator->pfree_func = lrec_evaluator_i_s_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_x_state_t {
	mv_unary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
} lrec_evaluator_x_x_state_t;

mv_t lrec_evaluator_x_x_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_x_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);

	return pstate->pfunc(&val1);
}
static void lrec_evaluator_x_x_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_x_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_x_func(mv_unary_func_t* pfunc, lrec_evaluator_t* parg1) {
	lrec_evaluator_x_x_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_x_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_x_func;
	pevaluator->pfree_func = lrec_evaluator_x_x_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_b_xx_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_b_xx_state_t;

mv_t lrec_evaluator_b_xx_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_b_xx_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_b_xx_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_b_xx_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_b_xx_func;
	pevaluator->pfree_func = lrec_evaluator_b_xx_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_ns_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_x_ns_state_t;

mv_t lrec_evaluator_x_ns_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_ns_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	mv_set_number_nullable(&val1);
	NULL_OR_ERROR_OUT(val1);

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING) {
		mv_free(&val1);
		return MV_ERROR;
	}

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_x_ns_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_ns_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_ns_func(mv_binary_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_x_ns_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_ns_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_ns_func;
	pevaluator->pfree_func = lrec_evaluator_x_ns_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_ss_state_t {
	mv_binary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_x_ss_state_t;

mv_t lrec_evaluator_x_ss_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_ss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_x_ss_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_ss_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_x_ss_func;
	pevaluator->pfree_func = lrec_evaluator_x_ss_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_ssc_state_t {
	mv_binary_arg3_capture_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_x_ssc_state_t;

mv_t lrec_evaluator_x_ssc_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_ssc_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2, pregex_captures);
}
static void lrec_evaluator_x_ssc_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_ssc_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_ssc_func(mv_binary_arg3_capture_func_t* pfunc,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2)
{
	lrec_evaluator_x_ssc_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_ssc_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;
	pstate->parg2 = parg2;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_ssc_func;
	pevaluator->pfree_func = lrec_evaluator_x_ssc_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_sr_state_t {
	mv_binary_arg2_regex_func_t* pfunc;
	lrec_evaluator_t*             parg1;
	regex_t                       regex;
	string_builder_t*             psb;
} lrec_evaluator_x_sr_state_t;

mv_t lrec_evaluator_x_sr_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_sr_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &pstate->regex, pstate->psb, pregex_captures);
}
static void lrec_evaluator_x_sr_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_sr_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	regfree(&pstate->regex);
	sb_free(pstate->psb);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_sr_func(mv_binary_arg2_regex_func_t* pfunc,
	lrec_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	lrec_evaluator_x_sr_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_sr_state_t));
	pstate->pfunc = pfunc;
	pstate->parg1 = parg1;

	int cflags = ignore_case ? REG_ICASE : 0;
	regcomp_or_die(&pstate->regex, regex_string, cflags);
	pstate->psb = sb_alloc(MV_SB_ALLOC_LENGTH);

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_sr_func;
	pevaluator->pfree_func = lrec_evaluator_x_sr_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_xs_state_t {
	mv_binary_func_t*  pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
} lrec_evaluator_s_xs_state_t;

mv_t lrec_evaluator_s_xs_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_s_xs_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING)
		return MV_ERROR;

	return pstate->pfunc(&val1, &val2);
}
static void lrec_evaluator_s_xs_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_s_xs_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_s_xs_func;
	pevaluator->pfree_func = lrec_evaluator_s_xs_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_s_sss_state_t {
	mv_ternary_func_t* pfunc;
	lrec_evaluator_t* parg1;
	lrec_evaluator_t* parg2;
	lrec_evaluator_t* parg3;
} lrec_evaluator_s_sss_state_t;

mv_t lrec_evaluator_s_sss_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_s_sss_state_t* pstate = pvstate;
	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val2 = pstate->parg2->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg2->pvstate);
	NULL_OR_ERROR_OUT(val2);
	if (val2.type != MT_STRING) {
		mv_free(&val1);
		return MV_ERROR;
	}

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT(val3);
	if (val3.type != MT_STRING) {
		mv_free(&val1);
		mv_free(&val2);
		return MV_ERROR;
	}

	return pstate->pfunc(&val1, &val2, &val3);
}
static void lrec_evaluator_s_sss_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_s_sss_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	pstate->parg2->pfree_func(pstate->parg2);
	pstate->parg3->pfree_func(pstate->parg3);
	free(pstate);
	free(pevaluator);
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
	pevaluator->pprocess_func = lrec_evaluator_s_sss_func;
	pevaluator->pfree_func = lrec_evaluator_s_sss_free;

	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _lrec_evaluator_x_srs_state_t {
	mv_ternary_arg2_regex_func_t* pfunc;
	lrec_evaluator_t*             parg1;
	regex_t                       regex;
	lrec_evaluator_t*             parg3;
	string_builder_t*             psb;
} lrec_evaluator_x_srs_state_t;

mv_t lrec_evaluator_x_srs_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_x_srs_state_t* pstate = pvstate;

	mv_t val1 = pstate->parg1->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg1->pvstate);
	NULL_OR_ERROR_OUT(val1);
	if (val1.type != MT_STRING)
		return MV_ERROR;

	mv_t val3 = pstate->parg3->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pstate->parg3->pvstate);
	NULL_OR_ERROR_OUT(val3);
	if (val3.type != MT_STRING) {
		mv_free(&val1);
		return MV_ERROR;
	}

	return pstate->pfunc(&val1, &pstate->regex, pstate->psb, &val3);
}
static void lrec_evaluator_x_srs_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_x_srs_state_t* pstate = pevaluator->pvstate;
	pstate->parg1->pfree_func(pstate->parg1);
	regfree(&pstate->regex);
	pstate->parg3->pfree_func(pstate->parg3);
	sb_free(pstate->psb);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_x_srs_func(mv_ternary_arg2_regex_func_t* pfunc,
	lrec_evaluator_t* parg1, char* regex_string, int ignore_case, lrec_evaluator_t* parg3)
{
	lrec_evaluator_x_srs_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_x_srs_state_t));
	pstate->pfunc = pfunc;

	pstate->parg1 = parg1;

	int cflags = ignore_case ? REG_ICASE : 0;
	regcomp_or_die(&pstate->regex, regex_string, cflags);
	pstate->psb = sb_alloc(MV_SB_ALLOC_LENGTH);

	pstate->parg3 = parg3;

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = lrec_evaluator_x_srs_func;
	pevaluator->pfree_func = lrec_evaluator_x_srs_free;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_field_name_state_t {
	char* field_name;
} lrec_evaluator_field_name_state_t;

mv_t lrec_evaluator_field_name_func_string_only(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_field_name_state_t* pstate = pvstate;
	// xxx comment here ...
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// xxx comment
		return mv_copy(poverlay); // xxx mem-mgmt for strings ...
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return MV_NULL;
		} else {
			// string points into AST memory and is valid as long as the AST is.
			return mv_from_string_no_free(string);
		}
	}
}

mv_t lrec_evaluator_field_name_func_string_float(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_field_name_state_t* pstate = pvstate;
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// xxx comment
		return mv_copy(poverlay); // xxx mem-mgmt for strings ...
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return MV_NULL;
		} else {
			double fltv;
			if (mlr_try_float_from_string(string, &fltv)) {
				return mv_from_float(fltv);
			} else {
				// string points into AST memory and is valid as long as the AST is.
				return mv_from_string_no_free(string);
			}
		}
	}
}

mv_t lrec_evaluator_field_name_func_string_float_int(lrec_t* prec, lhmsv_t* ptyped_overlay,
	string_array_t* pregex_captures, context_t* pctx, void* pvstate)
{
	lrec_evaluator_field_name_state_t* pstate = pvstate;
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// xxx comment
		return mv_copy(poverlay); // xxx mem-mgmt for strings ...
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return MV_NULL;
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
static void lrec_evaluator_field_name_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_field_name_state_t* pstate = pevaluator->pvstate;
	free(pstate->field_name);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_field_name(char* field_name, int type_inferencing) {
	lrec_evaluator_field_name_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_field_name_state_t));
	pstate->field_name = mlr_strdup_or_die(field_name);

	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	switch (type_inferencing) {
	case TYPE_INFER_STRING_ONLY:
		pevaluator->pprocess_func = lrec_evaluator_field_name_func_string_only;
		break;
	case TYPE_INFER_STRING_FLOAT:
		pevaluator->pprocess_func = lrec_evaluator_field_name_func_string_float;
		break;
	case TYPE_INFER_STRING_FLOAT_INT:
		pevaluator->pprocess_func = lrec_evaluator_field_name_func_string_float_int;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
		break;
	}
	pevaluator->pfree_func = lrec_evaluator_field_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _lrec_evaluator_literal_state_t {
	mv_t literal;
} lrec_evaluator_literal_state_t;

mv_t lrec_evaluator_non_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

mv_t lrec_evaluator_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	lrec_evaluator_literal_state_t* pstate = pvstate;

	char* input = pstate->literal.u.strv;
	int was_allocated = FALSE;
	char *output = interpolate_regex_captures(input, pregex_captures, &was_allocated);
	if (was_allocated) {
		return mv_from_string_with_free(output);
	} else {
		return mv_from_string_no_free(output);
	}
}
static void lrec_evaluator_literal_free(lrec_evaluator_t* pevaluator) {
	lrec_evaluator_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_literal(char* string, int type_inferencing) {
	lrec_evaluator_literal_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_evaluator_literal_state_t));
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));

	if (string == NULL) {
		pstate->literal = MV_NULL;
		pevaluator->pprocess_func = lrec_evaluator_non_string_literal_func;
	} else {
		long long intv;
		double fltv;

		pevaluator->pprocess_func = NULL;
		switch (type_inferencing) {
		case TYPE_INFER_STRING_ONLY:
			pstate->literal = mv_from_string_no_free(string);
			pevaluator->pprocess_func = lrec_evaluator_string_literal_func;
			break;

		case TYPE_INFER_STRING_FLOAT:
			if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = lrec_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = lrec_evaluator_string_literal_func;
			}
			break;

		case TYPE_INFER_STRING_FLOAT_INT:
			if (mlr_try_int_from_string(string, &intv)) {
				pstate->literal = mv_from_int(intv);
				pevaluator->pprocess_func = lrec_evaluator_non_string_literal_func;
			} else if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = lrec_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = lrec_evaluator_string_literal_func;
			}
			break;

		default:
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
			break;
		}
	}
	pevaluator->pfree_func = lrec_evaluator_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
mv_t lrec_evaluator_NF_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_int(prec->field_count);
}
static void lrec_evaluator_NF_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NF() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_NF_func;
	pevaluator->pfree_func = lrec_evaluator_NF_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_NR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->nr);
}
static void lrec_evaluator_NR_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_NR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_NR_func;
	pevaluator->pfree_func = lrec_evaluator_NR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FNR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->fnr);
}
static void lrec_evaluator_FNR_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FNR() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_FNR_func;
	pevaluator->pfree_func = lrec_evaluator_FNR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENAME_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_string_no_free(pctx->filename);
}
static void lrec_evaluator_FILENAME_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}

lrec_evaluator_t* lrec_evaluator_alloc_from_FILENAME() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_FILENAME_func;
	pevaluator->pfree_func = lrec_evaluator_FILENAME_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_FILENUM_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->filenum);
}
static void lrec_evaluator_FILENUM_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_FILENUM() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_FILENUM_func;
	pevaluator->pfree_func = lrec_evaluator_FILENUM_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_PI_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_float(M_PI);
}
static void lrec_evaluator_PI_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_PI() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_PI_func;
	pevaluator->pfree_func = lrec_evaluator_PI_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t lrec_evaluator_E_func(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate)
{
	return mv_from_float(M_E);
}
static void lrec_evaluator_E_free(lrec_evaluator_t* pevaluator) {
	free(pevaluator);
}
lrec_evaluator_t* lrec_evaluator_alloc_from_E() {
	lrec_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(lrec_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = lrec_evaluator_E_func;
	pevaluator->pfree_func = lrec_evaluator_E_free;
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
		return lrec_evaluator_alloc_from_x_z_func(f_z_urand_func);
	} else if (streq(function_name, "urand32")) {
		return lrec_evaluator_alloc_from_x_z_func(i_z_urand32_func);
	} else if (streq(function_name, "systime")) {
		return lrec_evaluator_alloc_from_x_z_func(f_z_systime_func);
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

#define FUNC_CLASS_ARITHMETIC 0xa0
#define FUNC_CLASS_MATH       0xa1
#define FUNC_CLASS_BOOLEAN    0xa2
#define FUNC_CLASS_STRING     0xa3
#define FUNC_CLASS_CONVERSION 0xa4
#define FUNC_CLASS_TIME       0xa5

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

static void check_arity_with_report(function_lookup_t fcn_lookup_table[], char* function_name,
	int user_provided_arity)
{
	int arity = -1;
	int result = check_arity(fcn_lookup_table, function_name, user_provided_arity, &arity);
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
	case FUNC_CLASS_ARITHMETIC: return "arithmetic"; break;
	case FUNC_CLASS_MATH:       return "math";       break;
	case FUNC_CLASS_BOOLEAN:    return "boolean";    break;
	case FUNC_CLASS_STRING:     return "string";     break;
	case FUNC_CLASS_CONVERSION: return "conversion"; break;
	case FUNC_CLASS_TIME:       return "time";       break;
	default:                    return "???";        break;
	}
}

void lrec_evaluator_list_functions(FILE* o, char* leader) {
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
void lrec_evaluator_function_usage(FILE* output_stream, char* function_name) {
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

void lrec_evaluator_list_all_functions_raw(FILE* output_stream) {
	for (int i = 0; ; i++) {
		function_lookup_t* plookup = &FUNCTION_LOOKUP_TABLE[i];
		if (plookup->function_name == NULL) // end of table
			break;
		printf("%s\n", plookup->function_name);
	}
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_unary_func_name(char* fnnm, lrec_evaluator_t* parg1)  {
	if        (streq(fnnm, "!"))         { return lrec_evaluator_alloc_from_b_b_func(b_b_not_func,       parg1);
	} else if (streq(fnnm, "+"))         { return lrec_evaluator_alloc_from_x_n_func(n_n_upos_func,      parg1);
	} else if (streq(fnnm, "-"))         { return lrec_evaluator_alloc_from_x_n_func(n_n_uneg_func,      parg1);
	} else if (streq(fnnm, "~"))         { return lrec_evaluator_alloc_from_i_i_func(i_i_bitwise_not_func, parg1);
	} else if (streq(fnnm, "abs"))       { return lrec_evaluator_alloc_from_x_n_func(n_n_abs_func,       parg1);
	} else if (streq(fnnm, "acos"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_acos_func,      parg1);
	} else if (streq(fnnm, "acosh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_acosh_func,     parg1);
	} else if (streq(fnnm, "asin"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_asin_func,      parg1);
	} else if (streq(fnnm, "asinh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_asinh_func,     parg1);
	} else if (streq(fnnm, "atan"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_atan_func,      parg1);
	} else if (streq(fnnm, "atanh"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_atanh_func,     parg1);
	} else if (streq(fnnm, "boolean"))   { return lrec_evaluator_alloc_from_x_x_func(b_x_boolean_func,   parg1);
	} else if (streq(fnnm, "boolean"))   { return lrec_evaluator_alloc_from_x_x_func(b_x_boolean_func,   parg1);
	} else if (streq(fnnm, "cbrt"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_cbrt_func,      parg1);
	} else if (streq(fnnm, "ceil"))      { return lrec_evaluator_alloc_from_x_n_func(n_n_ceil_func,      parg1);
	} else if (streq(fnnm, "cos"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_cos_func,       parg1);
	} else if (streq(fnnm, "cosh"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_cosh_func,      parg1);
	} else if (streq(fnnm, "dhms2fsec")) { return lrec_evaluator_alloc_from_f_s_func(f_s_dhms2fsec_func, parg1);
	} else if (streq(fnnm, "dhms2sec"))  { return lrec_evaluator_alloc_from_f_s_func(i_s_dhms2sec_func,  parg1);
	} else if (streq(fnnm, "erf"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_erf_func,       parg1);
	} else if (streq(fnnm, "erfc"))      { return lrec_evaluator_alloc_from_f_f_func(f_f_erfc_func,      parg1);
	} else if (streq(fnnm, "exp"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_exp_func,       parg1);
	} else if (streq(fnnm, "expm1"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_expm1_func,     parg1);
	} else if (streq(fnnm, "float"))     { return lrec_evaluator_alloc_from_x_x_func(f_x_float_func,     parg1);
	} else if (streq(fnnm, "floor"))     { return lrec_evaluator_alloc_from_x_n_func(n_n_floor_func,     parg1);
	} else if (streq(fnnm, "fsec2dhms")) { return lrec_evaluator_alloc_from_s_f_func(s_f_fsec2dhms_func, parg1);
	} else if (streq(fnnm, "fsec2hms"))  { return lrec_evaluator_alloc_from_s_f_func(s_f_fsec2hms_func,  parg1);
	} else if (streq(fnnm, "gmt2sec"))   { return lrec_evaluator_alloc_from_i_s_func(i_s_gmt2sec_func,   parg1);
	} else if (streq(fnnm, "hexfmt"))    { return lrec_evaluator_alloc_from_x_x_func(s_x_hexfmt_func,    parg1);
	} else if (streq(fnnm, "hms2fsec"))  { return lrec_evaluator_alloc_from_f_s_func(f_s_hms2fsec_func,  parg1);
	} else if (streq(fnnm, "hms2sec"))   { return lrec_evaluator_alloc_from_f_s_func(i_s_hms2sec_func,   parg1);
	} else if (streq(fnnm, "int"))       { return lrec_evaluator_alloc_from_x_x_func(i_x_int_func,       parg1);
	} else if (streq(fnnm, "invqnorm"))  { return lrec_evaluator_alloc_from_f_f_func(f_f_invqnorm_func,  parg1);
	} else if (streq(fnnm, "isnotnull")) { return lrec_evaluator_alloc_from_x_x_func(b_x_isnotnull_func, parg1);
	} else if (streq(fnnm, "isnull"))    { return lrec_evaluator_alloc_from_x_x_func(b_x_isnull_func,    parg1);
	} else if (streq(fnnm, "log"))       { return lrec_evaluator_alloc_from_f_f_func(f_f_log_func,       parg1);
	} else if (streq(fnnm, "log10"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_log10_func,     parg1);
	} else if (streq(fnnm, "log1p"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_log1p_func,     parg1);
	} else if (streq(fnnm, "qnorm"))     { return lrec_evaluator_alloc_from_f_f_func(f_f_qnorm_func,     parg1);
	} else if (streq(fnnm, "round"))     { return lrec_evaluator_alloc_from_x_n_func(n_n_round_func,     parg1);
	} else if (streq(fnnm, "sec2dhms"))  { return lrec_evaluator_alloc_from_s_i_func(s_i_sec2dhms_func,  parg1);
	} else if (streq(fnnm, "sec2gmt"))   { return lrec_evaluator_alloc_from_x_n_func(s_n_sec2gmt_func,   parg1);
	} else if (streq(fnnm, "sec2hms"))   { return lrec_evaluator_alloc_from_s_i_func(s_i_sec2hms_func,   parg1);
	} else if (streq(fnnm, "sgn"))       { return lrec_evaluator_alloc_from_x_n_func(n_n_sgn_func,       parg1);
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
	if        (streq(fnnm, "&&"))     { return lrec_evaluator_alloc_from_b_bb_func(b_bb_and_func,          parg1, parg2);
	} else if (streq(fnnm, "||"))     { return lrec_evaluator_alloc_from_b_bb_func(b_bb_or_func,           parg1, parg2);
	} else if (streq(fnnm, "^^"))     { return lrec_evaluator_alloc_from_b_bb_func(b_bb_xor_func,          parg1, parg2);
	} else if (streq(fnnm, "=~"))     { return lrec_evaluator_alloc_from_x_ssc_func(matches_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "!=~"))    { return lrec_evaluator_alloc_from_x_ssc_func(does_not_match_no_precomp_func, parg1, parg2);
	} else if (streq(fnnm, "=="))     { return lrec_evaluator_alloc_from_b_xx_func(eq_op_func,             parg1, parg2);
	} else if (streq(fnnm, "!="))     { return lrec_evaluator_alloc_from_b_xx_func(ne_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">"))      { return lrec_evaluator_alloc_from_b_xx_func(gt_op_func,             parg1, parg2);
	} else if (streq(fnnm, ">="))     { return lrec_evaluator_alloc_from_b_xx_func(ge_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<"))      { return lrec_evaluator_alloc_from_b_xx_func(lt_op_func,             parg1, parg2);
	} else if (streq(fnnm, "<="))     { return lrec_evaluator_alloc_from_b_xx_func(le_op_func,             parg1, parg2);
	} else if (streq(fnnm, "."))      { return lrec_evaluator_alloc_from_x_ss_func(s_ss_dot_func,          parg1, parg2);
	} else if (streq(fnnm, "+"))      { return lrec_evaluator_alloc_from_n_nn_func(n_nn_plus_func,         parg1, parg2);
	} else if (streq(fnnm, "-"))      { return lrec_evaluator_alloc_from_n_nn_func(n_nn_minus_func,        parg1, parg2);
	} else if (streq(fnnm, "*"))      { return lrec_evaluator_alloc_from_n_nn_func(n_nn_times_func,        parg1, parg2);
	} else if (streq(fnnm, "/"))      { return lrec_evaluator_alloc_from_n_nn_func(n_nn_divide_func,       parg1, parg2);
	} else if (streq(fnnm, "//"))     { return lrec_evaluator_alloc_from_n_nn_func(n_nn_int_divide_func,   parg1, parg2);
	} else if (streq(fnnm, "%"))      { return lrec_evaluator_alloc_from_n_nn_func(n_nn_mod_func,          parg1, parg2);
	} else if (streq(fnnm, "**"))     { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "pow"))    { return lrec_evaluator_alloc_from_f_ff_func(f_ff_pow_func,          parg1, parg2);
	} else if (streq(fnnm, "atan2"))  { return lrec_evaluator_alloc_from_f_ff_func(f_ff_atan2_func,        parg1, parg2);
	} else if (streq(fnnm, "max"))    { return lrec_evaluator_alloc_from_n_nn_nullable_func(n_nn_max_func, parg1, parg2);
	} else if (streq(fnnm, "min"))    { return lrec_evaluator_alloc_from_n_nn_nullable_func(n_nn_min_func, parg1, parg2);
	} else if (streq(fnnm, "roundm")) { return lrec_evaluator_alloc_from_n_nn_func(n_nn_roundm_func,       parg1, parg2);
	} else if (streq(fnnm, "fmtnum")) { return lrec_evaluator_alloc_from_s_xs_func(s_xs_fmtnum_func,       parg1, parg2);
	} else if (streq(fnnm, "urandint")) { return lrec_evaluator_alloc_from_i_ii_func(i_ii_urandint_func,   parg1, parg2);
	} else if (streq(fnnm, "|"))      { return lrec_evaluator_alloc_from_i_ii_func(i_ii_bitwise_or_func,   parg1, parg2);
	} else if (streq(fnnm, "^"))      { return lrec_evaluator_alloc_from_i_ii_func(i_ii_bitwise_xor_func,  parg1, parg2);
	} else if (streq(fnnm, "&"))      { return lrec_evaluator_alloc_from_i_ii_func(i_ii_bitwise_and_func,  parg1, parg2);
	} else if (streq(fnnm, "<<"))     { return lrec_evaluator_alloc_from_i_ii_func(i_ii_bitwise_lsh_func,  parg1, parg2);
	} else if (streq(fnnm, ">>"))     { return lrec_evaluator_alloc_from_i_ii_func(i_ii_bitwise_rsh_func,  parg1, parg2);
	} else if (streq(fnnm, "strftime")) { return lrec_evaluator_alloc_from_x_ns_func(s_ns_strftime_func,   parg1, parg2);
	} else if (streq(fnnm, "strptime")) { return lrec_evaluator_alloc_from_x_ss_func(i_ss_strptime_func,   parg1, parg2);
	} else  { return NULL; }
}

lrec_evaluator_t* lrec_evaluator_alloc_from_binary_regex_arg2_func_name(char* fnnm,
	lrec_evaluator_t* parg1, char* regex_string, int ignore_case)
{
	if        (streq(fnnm, "=~"))  {
		return lrec_evaluator_alloc_from_x_sr_func(matches_precomp_func,        parg1, regex_string, ignore_case);
	} else if (streq(fnnm, "!=~")) {
		return lrec_evaluator_alloc_from_x_sr_func(does_not_match_precomp_func, parg1, regex_string, ignore_case);
	} else  { return NULL; }
}

// ================================================================
lrec_evaluator_t* lrec_evaluator_alloc_from_ternary_func_name(char* fnnm,
	lrec_evaluator_t* parg1, lrec_evaluator_t* parg2, lrec_evaluator_t* parg3)
{
	if (streq(fnnm, "sub")) {
		return lrec_evaluator_alloc_from_s_sss_func(sub_no_precomp_func,  parg1, parg2, parg3);
	} else if (streq(fnnm, "gsub")) {
		return lrec_evaluator_alloc_from_s_sss_func(gsub_no_precomp_func, parg1, parg2, parg3);
	} else if (streq(fnnm, "logifit")) {
		return lrec_evaluator_alloc_from_f_fff_func(f_fff_logifit_func,   parg1, parg2, parg3);
	} else if (streq(fnnm, "madd")) {
		return lrec_evaluator_alloc_from_i_iii_func(i_iii_modadd_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "msub")) {
		return lrec_evaluator_alloc_from_i_iii_func(i_iii_modsub_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mmul")) {
		return lrec_evaluator_alloc_from_i_iii_func(i_iii_modmul_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "mexp")) {
		return lrec_evaluator_alloc_from_i_iii_func(i_iii_modexp_func,    parg1, parg2, parg3);
	} else if (streq(fnnm, "? :")) {
		return lrec_evaluator_alloc_from_ternop(parg1, parg2, parg3);
	} else  { return NULL; }
}

lrec_evaluator_t* lrec_evaluator_alloc_from_ternary_regex_arg2_func_name(char* fnnm,
	lrec_evaluator_t* parg1, char* regex_string, int ignore_case, lrec_evaluator_t* parg3)
{
	if (streq(fnnm, "sub"))  {
		return lrec_evaluator_alloc_from_x_srs_func(sub_precomp_func,  parg1, regex_string, ignore_case, parg3);
	} else if (streq(fnnm, "gsub"))  {
		return lrec_evaluator_alloc_from_x_srs_func(gsub_precomp_func, parg1, regex_string, ignore_case, parg3);
	} else  { return NULL; }
}

// ================================================================
static lrec_evaluator_t* lrec_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	int type_inferencing, function_lookup_t* fcn_lookup_table)
{
	if (pnode->pchildren == NULL) { // leaf node
		if (pnode->type == MLR_DSL_AST_NODE_TYPE_FIELD_NAME) {
			return lrec_evaluator_alloc_from_field_name(pnode->text, type_inferencing);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_LITERAL) {
			return lrec_evaluator_alloc_from_literal(pnode->text, type_inferencing);
		} else if (pnode->type == MLR_DSL_AST_NODE_TYPE_REGEXI) {
			return lrec_evaluator_alloc_from_literal(pnode->text, type_inferencing);
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

		check_arity_with_report(fcn_lookup_table, func_name, user_provided_arity);

		lrec_evaluator_t* pevaluator = NULL;
		if (user_provided_arity == 0) {
			pevaluator = lrec_evaluator_alloc_from_zary_func_name(func_name);
		} else if (user_provided_arity == 1) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
			pevaluator = lrec_evaluator_alloc_from_unary_func_name(func_name, parg1);
		} else if (user_provided_arity == 2) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MLR_DSL_AST_NODE_TYPE_LITERAL) {
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_binary_regex_arg2_func_name(func_name, parg1, parg2_node->text, FALSE);
			} else if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MLR_DSL_AST_NODE_TYPE_REGEXI) {
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_binary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT);
			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_binary_func_name(func_name, parg1, parg2);
			}

		} else if (user_provided_arity == 3) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MLR_DSL_AST_NODE_TYPE_LITERAL) {
				// sub/gsub-regex special case:
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				lrec_evaluator_t* parg3 = lrec_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text, FALSE, parg3);

			} else if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MLR_DSL_AST_NODE_TYPE_REGEXI) {
				// sub/gsub-regex special case:
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				lrec_evaluator_t* parg3 = lrec_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT, parg3);

			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				lrec_evaluator_t* parg1 = lrec_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
				lrec_evaluator_t* parg2 = lrec_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing, fcn_lookup_table);
				lrec_evaluator_t* parg3 = lrec_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing, fcn_lookup_table);
				pevaluator = lrec_evaluator_alloc_from_ternary_func_name(func_name, parg1, parg2, parg3);
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

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	lrec_evaluator_t* pevaluator = lrec_evaluator_alloc_from_ast_aux(pnode, type_inferencing, FUNCTION_LOOKUP_TABLE);
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
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	string_array_t* pregex_captures = NULL;

	mv_t val = pnr->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 888);

	val = pfnr->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pfnr->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_INT);
	mu_assert_lf(val.u.intv == 999);

	val = pfilename->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pfilename->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "filename-goes-here"));

	val = pfilenum->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pfilenum->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
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

	lrec_evaluator_t* ps       = lrec_evaluator_alloc_from_field_name("s", TYPE_INFER_STRING_FLOAT_INT);
	lrec_evaluator_t* pdef     = lrec_evaluator_alloc_from_literal("def", TYPE_INFER_STRING_FLOAT_INT);
	lrec_evaluator_t* pdot     = lrec_evaluator_alloc_from_x_ss_func(s_ss_dot_func, ps, pdef);
	lrec_evaluator_t* ptolower = lrec_evaluator_alloc_from_s_s_func(s_s_tolower_func, pdot);
	lrec_evaluator_t* ptoupper = lrec_evaluator_alloc_from_s_s_func(s_s_toupper_func, pdot);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "s", "abc", NO_FREE);
	printf("lrec s = %s\n", lrec_get(prec, "s"));

	mv_t val = ps->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, ps->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abc"));

	val = pdef->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pdef->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "def"));

	val = pdot->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pdot->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptolower->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, ptolower->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
	mu_assert_lf(val.type == MT_STRING);
	mu_assert_lf(val.u.strv != NULL);
	mu_assert_lf(streq(val.u.strv, "abcdef"));

	val = ptoupper->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, ptoupper->pvstate);
	printf("[%s] %s\n", mt_describe_type(val.type), mv_alloc_format_val(&val));
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

	lrec_evaluator_t* p2     = lrec_evaluator_alloc_from_literal("2.0", TYPE_INFER_STRING_FLOAT_INT);
	lrec_evaluator_t* px     = lrec_evaluator_alloc_from_field_name("x", TYPE_INFER_STRING_FLOAT_INT);
	lrec_evaluator_t* plogx  = lrec_evaluator_alloc_from_f_f_func(f_f_log10_func, px);
	lrec_evaluator_t* p2logx = lrec_evaluator_alloc_from_n_nn_func(n_nn_times_func, p2, plogx);
	lrec_evaluator_t* px2    = lrec_evaluator_alloc_from_n_nn_func(n_nn_times_func, px, px);
	lrec_evaluator_t* p4     = lrec_evaluator_alloc_from_n_nn_func(n_nn_times_func, p2, p2);

	mlr_dsl_ast_node_t* pxnode     = mlr_dsl_ast_node_alloc("x",  MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	mlr_dsl_ast_node_t* plognode   = mlr_dsl_ast_node_alloc_zary("log", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	mlr_dsl_ast_node_t* plogxnode  = mlr_dsl_ast_node_append_arg(plognode, pxnode);
	mlr_dsl_ast_node_t* p2node     = mlr_dsl_ast_node_alloc("2",   MLR_DSL_AST_NODE_TYPE_LITERAL);
	mlr_dsl_ast_node_t* p2logxnode = mlr_dsl_ast_node_alloc_binary("*", MLR_DSL_AST_NODE_TYPE_OPERATOR,
		p2node, plogxnode);

	lrec_evaluator_t*  pastr = lrec_evaluator_alloc_from_ast(p2logxnode, TYPE_INFER_STRING_FLOAT_INT);

	lrec_t* prec = lrec_unbacked_alloc();
	lhmsv_t* ptyped_overlay = lhmsv_alloc();
	string_array_t* pregex_captures = NULL;
	lrec_put(prec, "x", "4.5", NO_FREE);

	mv_t valp2     = p2->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, p2->pvstate);
	mv_t valp4     = p4->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, p4->pvstate);
	mv_t valpx     = px->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, px->pvstate);
	mv_t valpx2    = px2->pprocess_func(prec,    ptyped_overlay, pregex_captures, pctx, px2->pvstate);
	mv_t valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, pregex_captures, pctx, plogx->pvstate);
	mv_t valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, p2logx->pvstate);

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
		mv_describe_val(pastr->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, pastr->pvstate)));
	printf("\n");

	lrec_rename(prec, "x", "y", FALSE);

	valp2     = p2->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, p2->pvstate);
	valp4     = p4->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, p4->pvstate);
	valpx     = px->pprocess_func(prec,     ptyped_overlay, pregex_captures, pctx, px->pvstate);
	valpx2    = px2->pprocess_func(prec,    ptyped_overlay, pregex_captures, pctx, px2->pvstate);
	valplogx  = plogx->pprocess_func(prec,  ptyped_overlay, pregex_captures, pctx, plogx->pvstate);
	valp2logx = p2logx->pprocess_func(prec, ptyped_overlay, pregex_captures, pctx, p2logx->pvstate);

	printf("lrec   x        = %s\n", lrec_get(prec, "x"));
	printf("newval 2        = %s\n", mv_describe_val(valp2));
	printf("newval 4        = %s\n", mv_describe_val(valp4));
	printf("newval x        = %s\n", mv_describe_val(valpx));
	printf("newval x^2      = %s\n", mv_describe_val(valpx2));
	printf("newval log(x)   = %s\n", mv_describe_val(valplogx));
	printf("newval 2*log(x) = %s\n", mv_describe_val(valp2logx));

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
