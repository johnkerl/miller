#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "cli/mlrcli.h"
#include "lib/context.h"
#include "containers/lhmsmv.h"
#include "containers/local_stack.h"
#include "containers/loop_stack.h"
#include "containers/type_decl.h"
#include "dsl/mlr_dsl_ast.h"
#include "dsl/mlr_dsl_blocked_ast.h"
#include "dsl/rval_evaluators.h"
#include "dsl/rxval_evaluators.h"
#include "dsl/function_manager.h"
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

// ----------------------------------------------------------------
// dsl/mlr_dsl_stack_allocate.c
// Two-pass stack allocator which operates on the block-structured AST
// before the CST is build (mlr_dsl_stack_allocate.c).
void blocked_ast_allocate_locals(blocked_ast_t* paast, int trace);

// ----------------------------------------------------------------
// Forward references for virtual-function prototypes
struct _mlr_dsl_cst_t;
struct _mlr_dsl_cst_statement_t;
struct _subr_defsite_t;

// Parameter bag to reduce parameter-marshaling
typedef struct _cst_outputs_t {
	int*    pshould_emit_rec;
	sllv_t* poutrecs;
	char*   oosvar_flatten_separator;
	cli_writer_opts_t* pwriter_opts;
} cst_outputs_t;

// ----------------------------------------------------------------
typedef struct _cst_statement_block_t {
	int subframe_var_count;
	sllv_t* pstatements;
} cst_statement_block_t;

cst_statement_block_t* cst_statement_block_alloc(int subframe_var_count);
void cst_statement_block_free(cst_statement_block_t* pblock, context_t* pctx);

// ----------------------------------------------------------------
typedef struct _cst_top_level_statement_block_t {
	local_stack_frame_t* pframe;
	int max_var_depth;
	cst_statement_block_t* pblock;
} cst_top_level_statement_block_t;

cst_top_level_statement_block_t* cst_top_level_statement_block_alloc(int max_var_depth, int subframe_var_count);
void cst_top_level_statement_block_free(cst_top_level_statement_block_t* pblock, context_t* pctx);

// ----------------------------------------------------------------
// Generic handler for a statement.

// Handler for statement lists: begin/main/end; cond/if/for/while/do-while.
typedef void mlr_dsl_cst_block_handler_t(
	cst_statement_block_t* pblock,
	variables_t*           pvars,
	cst_outputs_t*         pcst_outputs);

// ----------------------------------------------------------------
// mlr_dsl_cst_statement_t is a base class extended by all manner of subclasses.
// The following are for their method pointers.
typedef struct _mlr_dsl_cst_statement_t* mlr_dsl_cst_statement_allocator_t(
	struct _mlr_dsl_cst_t* pcst,
	mlr_dsl_ast_node_t*    pnode,
	int                    type_inferencing,
	int                    context_flags);

typedef void mlr_dsl_cst_statement_handler_t(
	struct _mlr_dsl_cst_statement_t* pstatement,
	variables_t*                     pvars,
	cst_outputs_t*                   pcst_outputs);

typedef void mlr_dsl_cst_statement_freer_t(
	struct _mlr_dsl_cst_statement_t* pstatement,
	context_t* pctx);

// ----------------------------------------------------------------
// MLR_DSL_CST_STATEMENT OBJECT

typedef struct _mlr_dsl_cst_statement_t {

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Common to most or all statement types:

	// For trace-mode.
	mlr_dsl_ast_node_t* past_node;

	// Function-pointer for the handler of the given statement type, e.g. srec-assignment, while-loop, etc.
	mlr_dsl_cst_statement_handler_t* pstatement_handler;

	// Subclass destructor. It should free whatever's in the pvstate but it should not
	// free the pstatement itself.
	mlr_dsl_cst_statement_freer_t* pstatement_freer;

	// The reason for this being a function pointer is that there are two variants of
	// statement-list handlers: one for inside loop bodies which has to check
	// break/continue flags after each statement, and another for outside loop bodies
	// which doesn't need to check those. (This is a micro-optimization.) For bodyless
	// statements (e.g. assignment) this is null.
	cst_statement_block_t* pblock;
	mlr_dsl_cst_block_handler_t* pblock_handler;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Specific to each statement type:

	void* pvstate;

} mlr_dsl_cst_statement_t;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// For use by the statement-subclass constructors

mlr_dsl_cst_statement_t* mlr_dsl_cst_statement_valloc(
	mlr_dsl_ast_node_t*                    past_node,
	mlr_dsl_cst_statement_handler_t*       pstatement_handler,
	mlr_dsl_cst_statement_freer_t*         pstatement_freer,
	void*                                  pvstate);

mlr_dsl_cst_statement_t* mlr_dsl_cst_statement_valloc_with_block(
	mlr_dsl_ast_node_t*                    past_node,
	mlr_dsl_cst_statement_handler_t*       pstatement_handler,
	cst_statement_block_t*                 pblock,
	mlr_dsl_cst_block_handler_t*           pblock_handler,
	mlr_dsl_cst_statement_freer_t*         pstatement_freer,
	void*                                  pvstate);

// ----------------------------------------------------------------
// MLR_DSL_CST OBJECT

typedef struct _mlr_dsl_cst_t {
	sllv_t* pbegin_blocks;
	cst_top_level_statement_block_t* pmain_block;
	sllv_t* pend_blocks;

	// Function manager for built-in functions as well as user-defined functions (which are CST-specific).
	fmgr_t* pfmgr;

	// Subroutine bodies
	lhmsv_t* psubr_defsites;

	// Subroutine callsites, used to bootstrap (e.g. subroutine f calls subroutine g before the latter
	// has been defined).
	sllv_t* psubr_callsite_statements_to_resolve;

	// fflush on emit/tee/print/dump
	int flush_every_record;

	// The CST object retains the AST pointer (in order to reuse its strings etc. with minimal copying)
	// and will free the AST in the CST destructor.
	blocked_ast_t* paast;
} mlr_dsl_cst_t;

// ----------------------------------------------------------------
// CONSTRUCTORS/DESTRUCTORS/METHODS

// Notes:
// * do_final_filter is FALSE for mlr put, TRUE for mlr filter.
// * negate_final_filter is TRUE for mlr filter -x.
// * The CST object strips nodes off the raw AST, constructed by the Lemon parser, in order
//   to do analysis on it. Nonetheless the caller should free what's left.
mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int print_ast, int trace_stack_allocation,
	int type_inferencing, int flush_every_record, int do_final_filter, int negate_final_filter);

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags);

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_final_filter_statement(mlr_dsl_cst_t* pcst,
	mlr_dsl_ast_node_t* pnode, int negate_final_filter, int type_inferencing, int context_flags);

void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst, context_t* pctx);
void mlr_dsl_cst_statement_free(mlr_dsl_cst_statement_t* pstatement, context_t* pctx);

// Top-level entry point, e.g. from mapper_put.
void mlr_dsl_cst_handle_top_level_statement_blocks(
	sllv_t*      ptop_level_blocks, // block bodies for begins, main, ends
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

void mlr_dsl_cst_handle_top_level_statement_block(
	cst_top_level_statement_block_t* ptop_level_block,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs);

// Recursive entry point: block bodies for begin, main, end; cond, if, for, while.
void mlr_dsl_cst_handle_statement_block(
	cst_statement_block_t* pblock,
	variables_t*           pvars,
	cst_outputs_t*         pcst_outputs);

void mlr_dsl_cst_handle_statement_block_with_break_continue(
	cst_statement_block_t* pblock,
	variables_t*           pvars,
	cst_outputs_t*         pcst_outputs);

// Statement lists which are not curly-braced bodies: start/continuation/update statements for triple-for.
void mlr_dsl_cst_handle_statement_list(
	sllv_t*        pstatements,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

// ================================================================
// dsl/mlr_dsl_cst_func_subr.c

// ----------------------------------------------------------------
// cst_udf_state_t is data needed to execute the body of a user-defined function which is implemented by CST statements.
// udf_defsite_state_t is data needed for any user-defined function (no matter how implemented).
typedef struct _cst_udf_state_t {
	char*     name;
	int       arity;
	char**    parameter_names;
	int*      parameter_type_masks;
	cst_top_level_statement_block_t* ptop_level_block;
	char*     return_value_type_name;
	int       return_value_type_mask;
} cst_udf_state_t;

udf_defsite_state_t* mlr_dsl_cst_alloc_udf(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

void mlr_dsl_cst_free_udf(cst_udf_state_t* pstate, context_t* pctx);

// ----------------------------------------------------------------

typedef struct _subr_callsite_t {
	char* name;
	int   arity;
	int   type_inferencing;
	int   context_flags;
} subr_callsite_t;

typedef struct _subr_defsite_t {
	char*     name;
	int       arity;
	char**    parameter_names;
	int*      parameter_type_masks;
	cst_top_level_statement_block_t* ptop_level_block;
} subr_defsite_t;

subr_defsite_t* mlr_dsl_cst_alloc_subroutine(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

void mlr_dsl_cst_free_subroutine(subr_defsite_t* psubr_defsite, context_t* pctx);

// Invoked directly from the CST statement handler for a subroutine callsite.
// (Functions, by contrast, are invoked by callback from the right-hand-site-evaluator logic
// -- hence no execute-function method here.)
void mlr_dsl_cst_execute_subroutine(subr_defsite_t* pstate, variables_t* pvars,
	cst_outputs_t* pcst_outputs, int callsite_arity, boxed_xval_t* args);

// ================================================================
// For on-line help / manpage
// dsl/mlr_dsl_cst_keywords.c

void mlr_dsl_list_all_keywords_raw(FILE* output_stream);

// Pass function_name == NULL to get usage for all keywords:
void mlr_dsl_keyword_usage(FILE* output_stream, char* keyword);

// ================================================================
// Specific CST-statement subclasses

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_condish_statements.c
mlr_dsl_cst_statement_allocator_t alloc_conditional_block;
mlr_dsl_cst_statement_allocator_t alloc_if_head;
mlr_dsl_cst_statement_allocator_t alloc_while;
mlr_dsl_cst_statement_allocator_t alloc_do_while;
mlr_dsl_cst_statement_allocator_t alloc_bare_boolean;

mlr_dsl_cst_statement_t* alloc_filter(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

mlr_dsl_cst_statement_t* alloc_final_filter(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 negate_final_filter,
	int                 type_inferencing,
	int                 context_flags);

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_terminal_assignment_statements.c
mlr_dsl_cst_statement_allocator_t alloc_srec_assignment;
mlr_dsl_cst_statement_allocator_t alloc_indirect_srec_assignment;
mlr_dsl_cst_statement_allocator_t alloc_positional_srec_name_assignment;
mlr_dsl_cst_statement_allocator_t alloc_env_assignment;

// dsl/mlr_dsl_cst_map_assignment_statements.c
mlr_dsl_cst_statement_allocator_t alloc_full_srec_assignment;
mlr_dsl_cst_statement_t* alloc_local_variable_definition(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 type_mask);
mlr_dsl_cst_statement_allocator_t alloc_nonindexed_local_variable_assignment;
mlr_dsl_cst_statement_allocator_t alloc_indexed_local_variable_assignment;
mlr_dsl_cst_statement_allocator_t alloc_oosvar_assignment;
mlr_dsl_cst_statement_allocator_t alloc_full_oosvar_assignment;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_unset_statements.c
mlr_dsl_cst_statement_allocator_t alloc_unset;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_for_srec_statements.c
mlr_dsl_cst_statement_allocator_t alloc_for_srec;
mlr_dsl_cst_statement_allocator_t alloc_for_srec_key_only;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_for_map_statements.c
mlr_dsl_cst_statement_allocator_t alloc_for_map;
mlr_dsl_cst_statement_allocator_t alloc_for_map_key_only;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_triple_for_statements.c
mlr_dsl_cst_statement_allocator_t alloc_triple_for;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_loop_control_statements.c
mlr_dsl_cst_statement_allocator_t alloc_break;
mlr_dsl_cst_statement_allocator_t alloc_continue;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_return_statements.c
mlr_dsl_cst_statement_allocator_t alloc_return_void;  // For subroutines
mlr_dsl_cst_statement_allocator_t alloc_return_value; // For UDFs

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_output_statements.c

mlr_dsl_cst_statement_t* alloc_print(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	char*               print_terminator);

mlr_dsl_cst_statement_allocator_t alloc_tee;

mlr_dsl_cst_statement_allocator_t alloc_emitf;

mlr_dsl_cst_statement_t* alloc_emit(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 do_full_prefixing);

mlr_dsl_cst_statement_t* alloc_emit_lashed(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 do_full_prefixing);

mlr_dsl_cst_statement_allocator_t alloc_dump;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
// dsl/mlr_dsl_cst_func_subr.c

// When we allocate a callsite we can do so before the callee has been defined.
// Hence the two-step process, with the second step being an object-binding step.
mlr_dsl_cst_statement_allocator_t alloc_subr_callsite_statement;
void mlr_dsl_cst_resolve_subr_callsite(mlr_dsl_cst_t* pcst, mlr_dsl_cst_statement_t* pstatement);

#endif // MLR_DSL_CST_H
