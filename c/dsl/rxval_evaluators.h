#ifndef RXVAL_EVALUATORS_H
#define RXVAL_EVALUATORS_H

#include <stdio.h>
#include "containers/mvfuncs.h"
#include "containers/xvfuncs.h"
#include "dsl/mlr_dsl_ast.h"
#include "dsl/rval_evaluator.h"
#include "dsl/rxval_evaluator.h"
#include "dsl/function_manager.h"

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
//   type(s). See also comments in containers/mlrval.h.
//
// xxx update comment w/r/t rxvals
// ================================================================

// ================================================================
// rxval_expr_evaluators.c

// ----------------------------------------------------------------
// Topmost functions:

rxval_evaluator_t* rxval_evaluator_alloc_from_ast(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

// Next level:
rxval_evaluator_t* rxval_evaluator_alloc_from_map_literal(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_function_callsite(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_nonindexed_local_variable(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_indexed_local_variable(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_oosvar_keylist(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_full_oosvar(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_full_srec(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_wrapping_rval(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

// ================================================================
// rxval_func_evaluators.c

rxval_evaluator_t* rxval_evaluator_alloc_from_variadic_func(
	xv_variadic_func_t* pfunc,
	rxval_evaluator_t** pargs,
	int                 nargs);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_x_func(
	xv_unary_func_t*   pfunc,
	rxval_evaluator_t* parg1);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_m_func(
	xv_unary_func_t*   pfunc,
	rxval_evaluator_t* parg1);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mx_func(
	xv_binary_func_t*  pfunc,
	rxval_evaluator_t* parg1,
	rxval_evaluator_t* parg2);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ms_func(
	xv_binary_func_t*  pfunc,
	rxval_evaluator_t* parg1,
	rxval_evaluator_t* parg2);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ss_func(
	xv_binary_func_t*  pfunc,
	rxval_evaluator_t* parg1,
	rxval_evaluator_t* parg2);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mss_func(
	xv_ternary_func_t* pfunc,
	rxval_evaluator_t* parg1,
	rxval_evaluator_t* parg2,
	rxval_evaluator_t* parg3);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_sss_func(
	xv_ternary_func_t* pfunc,
	rxval_evaluator_t* parg1,
	rxval_evaluator_t* parg2,
	rxval_evaluator_t* parg3);

rxval_evaluator_t* rxval_evaluator_alloc_from_A_x_func(
	xv_unary_func_t*   pfunc,
	rxval_evaluator_t* parg1, char* desc);

#endif // RXVAL_EVALUATORS_H
