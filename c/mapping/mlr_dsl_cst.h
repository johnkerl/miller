#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "containers/mlr_dsl_ast.h"
#include "rval_evaluators.h"
#include "containers/lhmsmv.h"
#include "containers/bind_stack.h"
#include "containers/loop_stack.h"

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

// Forward references for virtual-function prototypes
struct _mlr_dsl_cst_statement_t;
struct _mlr_dsl_cst_statement_vararg_t;

// Parameter bag to reduce parameter-marshaling
typedef struct _cst_outputs_t {
	int*    pshould_emit_rec;
	sllv_t* poutrecs;
	char*   oosvar_flatten_separator;
} cst_outputs_t;

// Generic handler for a statement.
typedef void mlr_dsl_cst_node_handler_func_t(
	struct _mlr_dsl_cst_statement_t* pnode,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

// Subhandler for emitf/unset vararg items: e.g. in 'unset @o, $s' there is one for the @o and one for the $s.
typedef void unset_vararg_handler_t(
	struct _mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs);

// Most statements have one item, except emit and unset.
typedef struct _mlr_dsl_cst_statement_vararg_t {
	unset_vararg_handler_t* punset_handler;
	char* emitf_or_unset_srec_field_name;
	rval_evaluator_t* punset_srec_field_name_evaluator;
	rval_evaluator_t* pemitf_arg_evaluator;
	sllv_t* punset_oosvar_keylist_evaluators;
} mlr_dsl_cst_statement_vararg_t;

// Handler for statement lists: begin/main/end; cond/if/for/while/do-while.
typedef void mlr_dsl_cst_statement_list_handler_t(
	sllv_t*      pcst_statements,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

// Diffrence between keylist and namelist: in emit @a[$b]["c"], "d", @e,
// the keylist is [$b, "c"] and the namelist is ["d", @e].
typedef struct _mlr_dsl_cst_statement_t {

	// Function-pointer for the handler of the given statement type, e.g. srec-assignment, while-loop, etc.
	mlr_dsl_cst_node_handler_func_t* pnode_handler;

	// There are two variants of statement-list handlers: one for inside loop bodies which has to check break/continue
	// flags after each statement, and another for outside loop bodies which doesn't need to check those. (This is a
	// micro-optimization.) For bodyless statements (e.g. assignment) this is null.
	mlr_dsl_cst_statement_list_handler_t* pblock_handler;

	// Assignment to oosvar, emit, and emitp; indices ["a", 1, $2] in 'for (k,v in @a[1][$2]) {...}'.
	sllv_t* poosvar_lhs_keylist_evaluators;

	// Assignment to srec
	char* srec_lhs_field_name;

	// Indirect assignment to srec
	rval_evaluator_t* psrec_lhs_evaluator;

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
	slls_t* pfor_oosvar_k_names;
	char* for_v_name;

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

void mlr_dsl_cst_handle_statement_list(
	sllv_t*      pcst_statements, // block bodies for begin, main, end; cond, if, for, while
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

#endif // MLR_DSL_CST_H
