#ifndef RVAL_EVALUATORS_H
#define RVAL_EVALUATORS_H
#include <stdio.h>
#include "containers/mlr_dsl_ast.h"
#include "mapping/rval_evaluator.h"

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
// ================================================================

#define TYPE_INFER_STRING_FLOAT_INT 0xce08
#define TYPE_INFER_STRING_FLOAT     0xce09
#define TYPE_INFER_STRING_ONLY      0xce0a

typedef enum _func_class_t {
	FUNC_CLASS_ARITHMETIC,
	FUNC_CLASS_MATH,
	FUNC_CLASS_BOOLEAN,
	FUNC_CLASS_STRING,
	FUNC_CLASS_CONVERSION,
	FUNC_CLASS_TIME
} func_class_t;

typedef enum _arity_check_t {
	ARITY_CHECK_PASS,
	ARITY_CHECK_FAIL,
	ARITY_CHECK_NO_SUCH
} arity_check_t;

typedef struct _function_lookup_t {
	func_class_t function_class;
	char*        function_name;
	int          arity;
	char*        usage_string;
} function_lookup_t;

extern function_lookup_t FUNCTION_LOOKUP_TABLE[];

// ----------------------------------------------------------------
// rval_expr_evaluators.c

// Topmost function:
rval_evaluator_t* rval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* past, int type_inferencing);

// Next level fromrval_evaluator_alloc_from_ast:
rval_evaluator_t* rval_evaluator_alloc_from_field_name(char* field_name, int type_inferencing);

// This is used for evaluating @-variables that don't have brackets: e.g. @x vs. @x[$1].
// See comments above rval_evaluator_alloc_from_oosvar_level_keys for more information.
rval_evaluator_t* rval_evaluator_alloc_from_oosvar_name(char* oosvar_name);
// This is used for evaluating @-variables that don't have brackets: e.g. @x vs. @x[$1].
rval_evaluator_t* rval_evaluator_alloc_from_oosvar_level_keys(mlr_dsl_ast_node_t* past);

// This is used for evaluating strings and numbers in literal expressions, e.g. '$x = "abc"'
// or '$x = "left_\1". The values are subject to replacement with regex captures. See comments
// in mapper_put for more information.
//
// Compare rval_evaluator_alloc_from_string which doesn't do regex replacement: it is intended for
// oosvar names on expression left-hand sides (outside of this file).
rval_evaluator_t* rval_evaluator_alloc_from_strnum_literal(char* string, int type_inferencing);

// This is intended only for oosvar names on expression left-hand sides.
// Compare rval_evaluator_alloc_from_strnum_literal.
rval_evaluator_t* rval_evaluator_alloc_from_string(char* string);

rval_evaluator_t* rval_evaluator_alloc_from_boolean_literal(char* string);
rval_evaluator_t* rval_evaluator_alloc_from_boolean(int boolval);
rval_evaluator_t* rval_evaluator_alloc_from_environment(mlr_dsl_ast_node_t* pnode, int type_inferencing);
rval_evaluator_t* rval_evaluator_alloc_from_NF();
rval_evaluator_t* rval_evaluator_alloc_from_NR();
rval_evaluator_t* rval_evaluator_alloc_from_FNR();
rval_evaluator_t* rval_evaluator_alloc_from_FILENAME();
rval_evaluator_t* rval_evaluator_alloc_from_FILENUM();
rval_evaluator_t* rval_evaluator_alloc_from_PI();
rval_evaluator_t* rval_evaluator_alloc_from_E();
rval_evaluator_t* rval_evaluator_alloc_from_context_variable(char* variable_name);

// For unit test:
rval_evaluator_t* rval_evaluator_alloc_from_mlrval(mv_t* pval);

// ----------------------------------------------------------------
// rval_func_evaluators.c

void rval_evaluator_list_functions(FILE* output_stream, char* leader);
// Pass function_name == NULL to get usage for all functions:
void rval_evaluator_function_usage(FILE* output_stream, char* function_name);
void rval_evaluator_list_all_functions_raw(FILE* output_stream);
void check_arity_with_report(function_lookup_t fcn_lookup_table[], char* function_name,
	int user_provided_arity);

rval_evaluator_t* rval_evaluator_alloc_from_zary_func_name(char* function_name);
rval_evaluator_t* rval_evaluator_alloc_from_unary_func_name(char* fnnm, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_binary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_binary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case);
rval_evaluator_t* rval_evaluator_alloc_from_ternary_func_name(char* fnnm,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_ternary_regex_arg2_func_name(char* fnnm,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3);

// These have some shared code that would otherwise be duplicated per-function in containers/mlrval.c.

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
rval_evaluator_t* rval_evaluator_alloc_from_ternop(rval_evaluator_t* parg1, rval_evaluator_t* parg2,
	rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_s_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_s_f_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_s_i_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_f_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_i_s_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_x_x_func(mv_unary_func_t* pfunc, rval_evaluator_t* parg1);
rval_evaluator_t* rval_evaluator_alloc_from_x_ns_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_ss_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_ssc_func(mv_binary_arg3_capture_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_x_sr_func(mv_binary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case);
rval_evaluator_t* rval_evaluator_alloc_from_s_xs_func(mv_binary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2);
rval_evaluator_t* rval_evaluator_alloc_from_s_sss_func(mv_ternary_func_t* pfunc,
	rval_evaluator_t* parg1, rval_evaluator_t* parg2, rval_evaluator_t* parg3);
rval_evaluator_t* rval_evaluator_alloc_from_x_srs_func(mv_ternary_arg2_regex_func_t* pfunc,
	rval_evaluator_t* parg1, char* regex_string, int ignore_case, rval_evaluator_t* parg3);

// ----------------------------------------------------------------
// rval_list_evaluators.c

// Nominally for oosvar multikeys.
sllmv_t* evaluate_list(
	sllv_t*          pevaluators,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	mlhmmv_t*        poosvars,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pall_non_null_or_error);

#endif // LREC_FEVALUATORS_H
