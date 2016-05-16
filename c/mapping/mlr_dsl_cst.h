#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "containers/mlr_dsl_ast.h"
#include "rval_evaluators.h"
#include "containers/lhmsmv.h"
#include "containers/bind_stack.h"

// ================================================================
// Concrete syntax tree (CST) derived from an abstract syntax tree (AST).
//
// Statements are of the form:
//
// * Assignment of mlrval (i.e. result of expression evaluation, e.g. $name or f($x,$y)) to oosvar (out-of-stream
// variables, prefixed with @ sigil)
//
// * Assignment to srec (in-stream records, with field names prefixed with $ sigil)
//
// * Copying full srec ($* syntax) to/from an oosvar
//
// * Oosvar-to-oosvar assignments (recursively if RHS is non-terminal)
//
// * pattern-action statements: boolean expression with curly-braced statements which are executed only
//   when the boolean evaluates to true.
//
// * bare-boolean statements: no-ops unless they have side effects: namely, the matches/does-not-match
//   operators =~ and !=~ setting regex captures \1, \2, etc.
//
// * emit statements: these place oosvar key-value pairs into the output stream.  These can be of the following forms:
//
//   o 'emit @a; emit @b' which produce separate records such as a=3 and b=4
//
//   o 'emitf @a, @b' which produce records such as a=3,b=4
//
//   o For nested maps, 'emit @c, "x", "y"' in which case the first two map levels are pulled out and named "x" and "y"
//   in separate fields. See containers/mlhmmv.h for more information.
//
// Further, these statements are organized into three groups:
//
// * begin: executed once, before the first input record is read.
// * main:  executed for each input record.
// * end:   executed once, after the last input record is read.
//
// The exceptions being, of course, assignment to/from srec is disallowed for begin/end statements since those occur
// before/after stream processing, respectively.
// ================================================================

struct _mlr_dsl_cst_statement_t;

typedef void mlr_dsl_cst_node_handler_func_t(
	struct _mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

// Most statements have one item, except emit and unset.
typedef struct _mlr_dsl_cst_statement_vararg_t {
	char*   emitf_or_unset_srec_field_name;
	rval_evaluator_t* pemitf_arg_evaluator;
	sllv_t* punset_oosvar_keylist_evaluators;
} mlr_dsl_cst_statement_vararg_t;

// E.g. emit @a[$b]["c"], "d", @e: keylist is [$b, "c"] and namelist is ["d", @e].
typedef struct _mlr_dsl_cst_statement_t {
	// Function-pointer for the handler of the given statement type, e.g. srec-assignment, while-loop, etc.
	mlr_dsl_cst_node_handler_func_t* phandler;

	// Assignment to oosvar, emit, and emitp
	sllv_t* poosvar_lhs_keylist_evaluators;

	// Assignment to srec
	char* srec_lhs_field_name;

	// Assignments to srec or oosvar, as well as the boolean expression in filter, cond, and bare-boolean
	rval_evaluator_t* prhs_evaluator;

	// Assigning full srec from oosvar:
	sllv_t* poosvar_rhs_keylist_evaluators;

	// emit/emitp:
	sllv_t* pemit_oosvar_namelist_evaluators;

	// Vararg stuff for emit and unset
	sllv_t* pvarargs;

	// Pattern-action blocks, while, for, etc.
	sllv_t* pblock_statements;

	// if-elif-elif-else:
	sllv_t* pif_chain_statements;

	// for-srec:
	char* for_srec_k_name;
	char* for_srec_v_name;

	// xxx for-oosvar key-list of names

	// for-srec and for-oosvar:
	lhmsmv_t* pbound_variables;

} mlr_dsl_cst_statement_t;

typedef struct _mlr_dsl_cst_t {
	sllv_t* pbegin_statements;
	sllv_t* pmain_statements;
	sllv_t* pend_statements;
} mlr_dsl_cst_t;

// ----------------------------------------------------------------
// For mlr filter, which takes a subset of the syntax of mlr put. Namely, a single top-level
// bare-boolean statement.
mlr_dsl_ast_node_t* extract_filterable_statement(mlr_dsl_ast_t* past, int type_inferencing);

// ----------------------------------------------------------------
mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing);
void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst);

void mlr_dsl_cst_handle(
	sllv_t*          pcst_statements, // begin/main/end
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

#endif // MLR_DSL_CST_H
