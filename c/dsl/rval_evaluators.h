// ================================================================
// NOTES:
//
// * Code here evaluates right-hand-side values (rvals) and return mlrvals (mv_t).
//
// * This is used by mlr filter and mlr put.
//
// * Comparison to mlrval.c: the latter is functions from mlrval(s) to mlrval;
//   in this file we have the higher-level notion of evaluating lrec objects,
//   using mlrval.c to do so.
//
// * Functions prototyped here evaluate right-hand-side values (rvals) and
//   return mlrvals (mv_t).  This is for scalar-valued contexts: almost all
//   expressions except for rxval contexts.
//
//   Values propagating up through the concrete syntax tree are always
//   dynamically allocated: e.g. in '$c = $a . $b' the $a and $b are copied out
//   as ephemerals; in the concat function their concatenation is computed and
//   the ephemeral input arguments are freed; then the result is stored in field
//   $c.
//
//   This is distinct from rxvals which are copy-on-write:
//   expression-intermediate values are not always ephemeral.
//
// * There are two kinds of lrec-evaluators here: those with _x_ in their names
//   which accept various types of mlrval, with disposition-matrices in
//   mlrval.c functions, and those with _i_/_f_/_b_/_s_ (int, float, boolean,
//   string) which either type-check or type-coerce their arguments, invoking
//   type-specific functions in mlrval.c.  Those with _n_ take int or float
//   and also use disposition matrices.  In all cases it's the job of
//   rval_evaluators.c to invoke functions here with mlrvals of the correct
//   type(s). See also comments in lib/mlrval.h.
//
// ================================================================

#ifndef RVAL_EVALUATORS_H
#define RVAL_EVALUATORS_H

#include <stdio.h>
#include "lib/mvfuncs.h"
#include "containers/xvfuncs.h"
#include "dsl/mlr_dsl_ast.h"
#include "dsl/rval_evaluator.h"
#include "dsl/function_manager.h"

// ================================================================
// rval_expr_evaluators.c
// ================================================================

// Topmost function:
rval_evaluator_t* rval_evaluator_alloc_from_ast(
	mlr_dsl_ast_node_t* past, fmgr_t* pfmgr, int type_inferencing, int context_flags);

// Next level:
rval_evaluator_t* rval_evaluator_alloc_from_field_name(char* field_name, int type_inferencing);
rval_evaluator_t* rval_evaluator_alloc_from_indirect_field_name(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags);
rval_evaluator_t* rval_evaluator_alloc_from_oosvar_keylist(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags);
rval_evaluator_t* rval_evaluator_alloc_from_local_map_keylist(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags);

// This is used for evaluating strings and numbers in literal expressions, e.g. '$x = "abc"'
// or '$x = "left_\1". The values are subject to replacement with regex captures. See comments
// in mapper_put for more information.
//
// Compare rval_evaluator_alloc_from_string which doesn't do regex replacement: it is intended for
// oosvar names on expression left-hand sides (outside of this file).
rval_evaluator_t* rval_evaluator_alloc_from_string_literal(char* string);
rval_evaluator_t* rval_evaluator_alloc_from_numeric_literal(char* string);
rval_evaluator_t* rval_evaluator_alloc_from_boolean_literal(char* string);

// This is intended only for oosvar names on expression left-hand sides.
// Compare rval_evaluator_alloc_from_numeric_literal.
rval_evaluator_t* rval_evaluator_alloc_from_string(char* string);

rval_evaluator_t* rval_evaluator_alloc_from_boolean(int boolval);
rval_evaluator_t* rval_evaluator_alloc_from_environment(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags);
rval_evaluator_t* rval_evaluator_alloc_from_NF();
rval_evaluator_t* rval_evaluator_alloc_from_NR();
rval_evaluator_t* rval_evaluator_alloc_from_FNR();
rval_evaluator_t* rval_evaluator_alloc_from_FILENAME();
rval_evaluator_t* rval_evaluator_alloc_from_FILENUM();
rval_evaluator_t* rval_evaluator_alloc_from_PI();
rval_evaluator_t* rval_evaluator_alloc_from_E();
rval_evaluator_t* rval_evaluator_alloc_from_IPS();
rval_evaluator_t* rval_evaluator_alloc_from_IFS();
rval_evaluator_t* rval_evaluator_alloc_from_IRS();
rval_evaluator_t* rval_evaluator_alloc_from_OPS();
rval_evaluator_t* rval_evaluator_alloc_from_OFS();
rval_evaluator_t* rval_evaluator_alloc_from_ORS();
rval_evaluator_t* rval_evaluator_alloc_from_context_variable(char* variable_name);
rval_evaluator_t* rval_evaluator_alloc_from_local_variable(int vardef_frame_relative_index);

// For unit test:
rval_evaluator_t* rval_evaluator_alloc_from_mlrval(mv_t* pval);

// ================================================================
// rval_func_evaluators.c
// ================================================================

// These have some shared code that would otherwise be duplicated per-function in lib/mvfuncs.c.
rval_evaluator_t* rval_evaluator_alloc_from_variadic_func(mv_variadic_func_t* pfunc, rval_evaluator_t** pargs, int nargs);
rval_evaluator_t* rval_evaluator_alloc_from_b_b_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_b_bb_and_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_b_bb_or_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_b_bb_xor_func(rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_z_func(mv_zary_func_t* pfunc);
rval_evaluator_t* rval_evaluator_alloc_from_f_f_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_x_n_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_i_i_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_f_ff_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_xx_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_xx_nullable_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_f_fff_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_i_ii_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_i_iii_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_s_sii_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_ternop(rval_evaluator_t* parg1, rval_evaluator_t* parg2,
	rval_evaluator_t* parg3);

rval_evaluator_t* rval_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_s_i_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_x_x_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);

rval_evaluator_t* rval_evaluator_alloc_from_x_xi_func(
	mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1,
	rval_evaluator_t* parg2);

rval_evaluator_t* rval_evaluator_alloc_from_x_ns_func(
	mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1,
	rval_evaluator_t* parg2);

rval_evaluator_t* rval_evaluator_alloc_from_x_ss_func(
	mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1,
	rval_evaluator_t* parg2);

rval_evaluator_t* rval_evaluator_alloc_from_x_ssc_func(
	mv_binary_arg3_capture_func_t* pfunc,
	rval_evaluator_t* parg1,
	rval_evaluator_t* parg2);

rval_evaluator_t* rval_evaluator_alloc_from_x_sr_func(
	mv_binary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1,
	char* regex_string,
	int ignore_case);

rval_evaluator_t* rval_evaluator_alloc_from_x_se_func(
	mv_binary_arg2_regextract_func_t* pfunc,
	rval_evaluator_t* parg1,
	char* regex_string,
	int ignore_case);

rval_evaluator_t* rval_evaluator_alloc_from_s_xs_func(
	mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1,
	rval_evaluator_t* parg2);

rval_evaluator_t* rval_evaluator_alloc_from_s_sss_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);

rval_evaluator_t* rval_evaluator_alloc_from_x_srs_func(mv_ternary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3);

rval_evaluator_t* rval_evaluator_alloc_from_x_ses_func(mv_ternary_arg2_regextract_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3);

// ================================================================
// rval_list_evaluators.c
// ================================================================

// Nominally for oosvar multikeys.
sllmv_t* evaluate_list(sllv_t* pevaluators, variables_t* pvars, int* pall_non_null_or_error);

// ----------------------------------------------------------------
// Type-inferenced srec-field getters for the expression-evaluators, as well as for boundvars in srec for-loops.

// For RHS evaluation:
mv_t get_srec_value_string_only(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay);
mv_t get_srec_value_string_float(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay);
mv_t get_srec_value_string_float_int(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay);

// For boundvars in for-srec:
typedef mv_t type_inferenced_srec_field_copy_getter_t(lrece_t* pentry, lhmsmv_t* ptyped_overlay);
mv_t get_copy_srec_value_string_only_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay);
mv_t get_copy_srec_value_string_float_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay);
mv_t get_copy_srec_value_string_float_int_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay);

#endif // RVAL_EVALUATORS_H
