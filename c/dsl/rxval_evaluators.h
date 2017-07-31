// ================================================================
// These evaluate right-hand-side extended values (rxvals) and return the same.
//
// For scalar (non-extended) right-hand side values, everything is ephemeral as
// it propagates up through the concrete syntax tree: e.g. in '$c = $a . $b' the
// $a and $b are copied out as ephemerals; in the concat function their
// concatenation is computed and the ephemeral input arguments are freed; then
// the result is stored in field $c.
//
// But for extended values (here) everything is copy-on-write:
// expression-intermediate values are not always ephemeral.  This is due to the
// size of the data involved. We can do dump or emit of a nested hashmap stored
// in an oosvar or local without copying it; we can do mapdiff of two map-valued
// variables while not modifying or copying either argument.
//
// The boxed_xval_t decorates mlhmmv_value_t (extended value) with an
// is_ephemeral flag.  The mlhmmv_value_t in turn has a map or a scalar.
// ================================================================

#ifndef RXVAL_EVALUATORS_H
#define RXVAL_EVALUATORS_H

#include <stdio.h>
#include "lib/mvfuncs.h"
#include "containers/xvfuncs.h"
#include "dsl/mlr_dsl_ast.h"
#include "dsl/rval_evaluator.h"
#include "dsl/rxval_evaluator.h"
#include "dsl/function_manager.h"

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
	xv_variadic_func_t*  pfunc,
	sllv_t*              parg_nodes,
	fmgr_t*              pfmgr,
	int                  type_inferencing,
	int                  context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_x_func(
	xv_unary_func_t*    pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_m_func(
	xv_unary_func_t*    pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mx_func(
	xv_binary_func_t*   pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	mlr_dsl_ast_node_t* parg2_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ms_func(
	xv_binary_func_t*   pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	mlr_dsl_ast_node_t* parg2_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_ss_func(
	xv_binary_func_t*   pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	mlr_dsl_ast_node_t* parg2_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_mss_func(
	xv_ternary_func_t*  pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	mlr_dsl_ast_node_t* parg2_node,
	mlr_dsl_ast_node_t* parg3_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_x_sss_func(
	xv_ternary_func_t*  pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	mlr_dsl_ast_node_t* parg2_node,
	mlr_dsl_ast_node_t* parg3_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

rxval_evaluator_t* rxval_evaluator_alloc_from_A_x_func(
	xv_unary_func_t*    pfunc,
	mlr_dsl_ast_node_t* parg1_node,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags,
	char*               desc);

#endif // RXVAL_EVALUATORS_H
