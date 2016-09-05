#ifndef FUNCTION_MANAGER_H
#define FUNCTION_MANAGER_H

#include "containers/lhmsv.h"
#include "containers/mlr_dsl_ast.h" // xxx factor out this dependency?
#include "mapping/rval_evaluator.h"
#include "mapping/type_inference.h"

// ----------------------------------------------------------------
// xxx make static in fmgr

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

// xxx move to fcn manager, along with move functions -> methods there
typedef struct _function_lookup_t {
	func_class_t function_class;
	char*        function_name;
	int          arity;
	char*        usage_string;
} function_lookup_t;

extern function_lookup_t FUNCTION_LOOKUP_TABLE[]; // xxx rm

// ----------------------------------------------------------------

typedef struct _fmgr_t {
	function_lookup_t * function_lookup_table;
	lhmsv_t* pUDF_names_to_evaluators;
} fmgr_t;

fmgr_t* fmgr_alloc();
void fmgr_free(fmgr_t* pfmgr);
// xxx disallow redefine ?
void fmgr_install_UDF(fmgr_t* pfmgr, char* name, rval_evaluator_t* pevaluator);
rval_evaluator_t* fmgr_alloc_evaluator(fmgr_t* pfmgr, char* name);

void fmgr_list_functions(fmgr_t* pfmgr, FILE* output_stream, char* leader);
// Pass function_name == NULL to get usage for all functions:
void fmgr_function_usage(fmgr_t* pfmgr, FILE* output_stream, char* function_name);
void fmgr_list_all_functions_raw(fmgr_t* pfmgr, FILE* output_stream);

rval_evaluator_t* fmgr_alloc_from_operator_or_function(fmgr_t* pfmgr, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags);

#endif // FUNCTION_MANAGER_H
