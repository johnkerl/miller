#ifndef FUNCTION_MANAGER_H
#define FUNCTION_MANAGER_H

#include "containers/mlrval.h"
#include "containers/lhmsv.h"
#include "containers/mlr_dsl_ast.h"
#include "mapping/rval_evaluator.h"
#include "mapping/type_inference.h"

// ----------------------------------------------------------------
// Things a user-defined function (however it is implemented) needs in order to
// be called: pvstate it its own state (whatever that is), and it defines its
// own process and free functions implementing this interface.

struct _udf_defsite_state_t; // Forward reference

typedef mv_t udf_defsite_process_func_t(void* pvstate, int arity, mv_t* pargs, variables_t* pvars);
typedef void udf_defsite_free_func_t(struct _udf_defsite_state_t* pdefsite_state);

typedef struct _udf_defsite_state_t {
	void* pvstate;
	int   arity;
	udf_defsite_process_func_t* pprocess_func;
	udf_defsite_free_func_t* pfree_func;
} udf_defsite_state_t;

// ----------------------------------------------------------------
// Holds built-in functions as well as user-defined functions

struct _function_lookup_t; // Private to the .c file
typedef struct _fmgr_t {
	struct _function_lookup_t * function_lookup_table; // Built-ins
	lhmsv_t* pudf_names_to_defsite_states;             // UDFs
} fmgr_t;

// ----------------------------------------------------------------
fmgr_t* fmgr_alloc();

void fmgr_free(fmgr_t* pfmgr);

void fmgr_install_udf(fmgr_t* pfmgr, char* name, int arity, udf_defsite_state_t* pdefsitate_state);

// Callsites as defined by AST nodes
rval_evaluator_t* fmgr_alloc_from_operator_or_function_call(fmgr_t* pfmgr, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags);

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
void fmgr_list_functions(fmgr_t* pfmgr, FILE* output_stream, char* leader);

// Pass function_name == NULL to get usage for all functions:
void fmgr_function_usage(fmgr_t* pfmgr, FILE* output_stream, char* function_name);

void fmgr_list_all_functions_raw(fmgr_t* pfmgr, FILE* output_stream);

#endif // FUNCTION_MANAGER_H
