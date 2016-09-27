#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "cli/mlrcli.h"
#include "containers/mlr_dsl_ast.h"
#include "containers/lhmsmv.h"
#include "containers/bind_stack.h"
#include "containers/loop_stack.h"
#include "mapping/rval_evaluators.h"
#include "mapping/function_manager.h"
#include "output/multi_out.h"
#include "output/multi_lrec_writer.h"

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
struct _subr_defsite_t;

// Parameter bag to reduce parameter-marshaling
typedef struct _cst_outputs_t {
	int*    pshould_emit_rec;
	sllv_t* poutrecs;
	char*   oosvar_flatten_separator;
	int     flush_every_record; // fflush on emit/tee/print/dump
	cli_writer_opts_t* pwriter_opts;
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

// ----------------------------------------------------------------
// MLR_DSL_CST_STATEMENT OBJECT

// These hold all the member data needed to evaluate any CST statement. No one kind of statement
// uses all of them. They aren't expressed as a union since their count is small: there's one CST
// per mlr-put invocation, independent of the number of stream records processed.
//
// Difference between keylist and namelist: in emit @a[$b]["c"], "d", @e, the keylist is ["a", $b, "c"]
// and the namelist is ["d", @e].

typedef struct _mlr_dsl_cst_statement_t {

	// Function-pointer for the handler of the given statement type, e.g. srec-assignment, while-loop, etc.
	mlr_dsl_cst_node_handler_func_t* pnode_handler;

	// For subroutine callsites
	rval_evaluator_t** subr_callsite_argument_evaluators;
	mv_t* subr_callsite_arguments;
	struct _subr_callsite_t *psubr_callsite;
	struct _subr_defsite_t *psubr_defsite;

	// Definition of local variable within user-defined function. Uses prhs_evaluator for value.
	char* local_variable_name;

	// Return statement within user-defined function
	rval_evaluator_t* preturn_evaluator;

	// There are two variants of statement-list handlers: one for inside loop bodies which has to check break/continue
	// flags after each statement, and another for outside loop bodies which doesn't need to check those. (This is a
	// micro-optimization.) For bodyless statements (e.g. assignment) this is null.
	mlr_dsl_cst_statement_list_handler_t* pblock_handler;

	// Assignment to oosvar
	sllv_t* poosvar_lhs_keylist_evaluators;

	// unlashed emit and emitp; indices ["a", 1, $2] in 'for (k,v in @a[1][$2]) {...}'.
	sllv_t* pemit_keylist_evaluators;

	// lashed emit and emitp; indices ["a", 1, $2] in 'for (k,v in @a[1][$2]) {...}'.
	int num_emit_keylist_evaluators;
	sllv_t** ppemit_keylist_evaluators;

	// Assignment to local
	char* local_lhs_variable_name;

	// Assignment to srec
	char* srec_lhs_field_name;

	// Indirect assignment to srec
	rval_evaluator_t* psrec_lhs_evaluator;

	// Assignments to srec or oosvar, as well as the boolean expression in filter, cond, and bare-boolean
	rval_evaluator_t* prhs_evaluator;

	// For print/printn/eprint/eprintn
	FILE* stdfp;
	char* print_terminator;

	// For print-to-file and dump-to-file, and emit-to-file
	rval_evaluator_t* poutput_filename_evaluator;
	file_output_mode_t file_output_mode;
	multi_out_t* pmulti_out; // print-to-file and dump-to-file
	lrec_writer_t* psingle_lrec_writer; // emit/tee to stdout/stderr
	multi_lrec_writer_t* pmulti_lrec_writer; // emit-to-file

	// Assigning full srec from oosvar:
	sllv_t* poosvar_rhs_keylist_evaluators;

	// emit/emitp:
	sllv_t* pemit_oosvar_namelist_evaluators;

	// Vararg stuff for emit and unset
	sllv_t* pvarargs;

	// emit vs. emitp
	int do_full_prefixing;

	// Pattern-action blocks, while, for, etc.
	sllv_t* pblock_statements;

	// if-elif-elif-else:
	sllv_t* pif_chain_statements;

	// for-srec / for-oosvar:
	char* for_srec_k_name;
	slls_t* pfor_oosvar_k_names;
	char* for_v_name;
	type_infererenced_srec_field_getter_t* ptype_infererenced_srec_field_getter;

	// triple-for:
	sllv_t* ptriple_for_start_statements;
	rval_evaluator_t* ptriple_for_continuation_evaluator;
	sllv_t* ptriple_for_update_statements;

	// for any kind of statement-block
	bind_stack_frame_t* pframe;

} mlr_dsl_cst_statement_t;

// ----------------------------------------------------------------
// MLR_DSL_CST OBJECT

typedef struct _mlr_dsl_cst_t {
	sllv_t* pbegin_statements;
	sllv_t* pmain_statements;
	sllv_t* pend_statements;

	// Function manager for built-in functions as well as user-defined functions (which are CST-specific).
	fmgr_t* pfmgr;

	// Subroutine bodies
	lhmsv_t* psubr_defsites;

	// Subroutine callsites, used to bootstrap (e.g. subroutine f calls subroutine g before the latter
	// has been defined).
	sllv_t* psubr_callsite_statements_to_resolve;

	// For mlr filter which takes restricted syntax
	rval_evaluator_t* pfilter_evaluator;
} mlr_dsl_cst_t;

// ----------------------------------------------------------------
// CONSTRUCTORS/DESTRUCTORS/METHODS

// For mlr filter, which takes a subset of the syntax of mlr put. Namely, a single top-level
// bare-boolean statement.

mlr_dsl_cst_t* mlr_dsl_cst_alloc_filterable(mlr_dsl_ast_t* ptop, int type_inferencing);

mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing,
	int do_filter); // xxx temp

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags);

void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst);
void mlr_dsl_cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

// Top-level entry point, e.g. from mapper_put.
void mlr_dsl_cst_handle_base_statement_list(
	sllv_t*      pcst_statements, // block bodies for begin, main, end; cond, if, for, while
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

// Recursive entry point: block bodies for begin, main, end; cond, if, for, while.
void mlr_dsl_cst_handle_statement_list(
	sllv_t*      pcst_statements,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

// ================================================================
// mapping/mlr_dsl_cst_func_subr.c

// ----------------------------------------------------------------
// cst_udf_state_t is data needed to execute the body of a user-defined function which is implemented by CST statements.
// udf_defsite_state_t is data needed for any user-defined function (no matter how implemented).
typedef struct _cst_udf_state_t {
	int       arity;
	char**    parameter_names;
    bind_stack_frame_t* pframe;
	sllv_t*   pblock_statements;
} cst_udf_state_t;

udf_defsite_state_t* mlr_dsl_cst_alloc_udf(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

void mlr_dsl_cst_free_udf(cst_udf_state_t* pstate);

// ----------------------------------------------------------------

typedef struct _subr_callsite_t {
	char* name;
	int   arity;
	int   type_inferencing;
	int   context_flags;
} subr_callsite_t;

subr_callsite_t* subr_callsite_alloc(char* name, int arity, int type_inferencing, int context_flags);
void subr_callsite_free(subr_callsite_t* psubr_callsite);

typedef struct _subr_defsite_t {
	char*     name;
	int       arity;
	char**    parameter_names;
    bind_stack_frame_t* pframe;
	sllv_t*   pblock_statements;
} subr_defsite_t;

subr_defsite_t* mlr_dsl_cst_alloc_subroutine(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

void mlr_dsl_cst_free_subroutine(subr_defsite_t* psubr_defsite);

// Invoked directly from the CST statement handler for a subroutine callsite.
// (Functions, by contrast, are invoked by callback from the right-hand-site-evaluator logic
// -- hence no execute-function method here.)
void mlr_dsl_cst_execute_subroutine(subr_defsite_t* pstate, variables_t* pvars,
	cst_outputs_t* pcst_outputs, int callsite_arity, mv_t* args);

// ================================================================
// For on-line help / manpage
// mapping/mlr_dsl_cst_keywords.c

void mlr_dsl_list_all_keywords_raw(FILE* output_stream);

// Pass function_name == NULL to get usage for all keywords:
void mlr_dsl_keyword_usage(FILE* output_stream, char* keyword);

#endif // MLR_DSL_CST_H
