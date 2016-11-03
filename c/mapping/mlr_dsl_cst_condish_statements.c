#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
typedef struct _conditional_block_state_t {
	rval_evaluator_t* pexpression_evaluator;
} conditional_block_state_t;

static mlr_dsl_cst_statement_handler_t handle_conditional_block;
static mlr_dsl_cst_statement_freer_t free_conditional_block;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_conditional_block(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	conditional_block_state_t* pstate = mlr_malloc_or_die(sizeof(conditional_block_state_t));

	pstate->pexpression_evaluator = NULL;

	// Right node is a list of statements to be executed if the left evaluates to true.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	pstate->pexpression_evaluator = rval_evaluator_alloc_from_ast(
		pleft, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pright->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pright->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pchild_statement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock->pstatements, pchild_statement);
	}

	mlr_dsl_cst_block_handler_t* pblock_handler = (context_flags & IN_BREAKABLE)
		? handle_statement_block_with_break_continue
		: mlr_dsl_cst_handle_statement_block;

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_conditional_block,
		pblock,
		pblock_handler,
		free_conditional_block,
		pstate);
}

// ----------------------------------------------------------------
// xxx move all frees between allocs & handles. and header-file order too.

static void free_conditional_block(mlr_dsl_cst_statement_t* pstatement) { // conditional_block
	conditional_block_state_t* pstate = pstatement->pvstate;

	pstate->pexpression_evaluator->pfree_func(pstate->pexpression_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_conditional_block(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	conditional_block_state_t* pstate = pstatement->pvstate;

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);

	rval_evaluator_t* pexpression_evaluator = pstate->pexpression_evaluator;

	mv_t val = pexpression_evaluator->pprocess_func(pexpression_evaluator->pvstate, pvars);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		if (val.u.boolv) {
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		}
	}

	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}

// ================================================================
typedef struct _if_head_state_t {
	sllv_t* pif_chain_statements;
} if_head_state_t;

typedef struct _if_item_state_t {
	rval_evaluator_t* pexpression_evaluator;
} if_item_state_t;

static mlr_dsl_cst_statement_handler_t handle_if_head;
static mlr_dsl_cst_statement_freer_t free_if_head;
static mlr_dsl_cst_statement_freer_t free_if_item;

static mlr_dsl_cst_statement_t* alloc_if_item(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pitemnode,
	mlr_dsl_ast_node_t* pexprnode,
	mlr_dsl_ast_node_t* plistnode,
	int                 type_inferencing,
	int                 context_flags);

// ----------------------------------------------------------------
// Example parser-input:
//
//   if (NR == 9) {
//       $x = 10;
//       $x = 11
//   } elif (NR == 12) {
//       $x = 13;
//       $x = 14
//   } else {
//       $x = 15;
//       $x = 16
//   };
//
// Corresponding parser-output AST:
//   if_head (if_head):
//       if (if_item):
//           == (operator):
//               NR (context_variable).
//               9 (numeric_literal).
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   10 (numeric_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   11 (numeric_literal).
//       elif (if_item):
//           == (operator):
//               NR (context_variable).
//               12 (numeric_literal).
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   13 (numeric_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   14 (numeric_literal).
//       else (if_item):
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   15 (numeric_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   16 (numeric_literal).

mlr_dsl_cst_statement_t* alloc_if_head(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	if_head_state_t* pstate = mlr_malloc_or_die(sizeof(if_head_state_t));

	pstate->pif_chain_statements = sllv_alloc();

	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		// For if and elif:
		// * Left subnode is the AST for the boolean expression.
		// * Right subnode is a list of statements to be executed if the left evaluates to true.
		// For else:
		// * Sole subnode is a list of statements to be executed.
		mlr_dsl_ast_node_t* pitemnode = pe->pvvalue;
		mlr_dsl_ast_node_t* pexprnode = NULL;
		mlr_dsl_ast_node_t* plistnode = NULL;
		if (pitemnode->pchildren->length == 2) {
			pexprnode = pitemnode->pchildren->phead->pvvalue;
			plistnode = pitemnode->pchildren->phead->pnext->pvvalue;
		} else {
			pexprnode = NULL;
			plistnode = pitemnode->pchildren->phead->pvvalue;
		}

		sllv_append(pstate->pif_chain_statements,
			alloc_if_item(pcst, pitemnode, pexprnode, plistnode, type_inferencing, context_flags)
		);
	}

	mlr_dsl_cst_block_handler_t* pblock_handler = (context_flags & IN_BREAKABLE)
		?  handle_statement_block_with_break_continue
		: mlr_dsl_cst_handle_statement_block;

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_if_head,
		NULL,
		pblock_handler,
		free_if_head,
		pstate);
}

static mlr_dsl_cst_statement_t* alloc_if_item(mlr_dsl_cst_t* pcst,
	mlr_dsl_ast_node_t* pitemnode, mlr_dsl_ast_node_t* pexprnode, mlr_dsl_ast_node_t* plistnode,
	int type_inferencing, int context_flags)
{
	if_item_state_t* pstate = mlr_malloc_or_die(sizeof(if_item_state_t));

	MLR_INTERNAL_CODING_ERROR_IF(plistnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(plistnode->subframe_var_count);

	for (sllve_t* pe = plistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pchild_statement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock->pstatements, pchild_statement);
	}

	pstate->pexpression_evaluator = pexprnode != NULL
		? rval_evaluator_alloc_from_ast(pexprnode, pcst->pfmgr,
			type_inferencing, context_flags) // if-statement or elif-statement
		: rval_evaluator_alloc_from_boolean(TRUE); // else-statement

	return mlr_dsl_cst_statement_valloc_with_block(
		pitemnode,
		NULL, // handled by the containing if-head evaluator
		pblock,
		NULL, // handled by the containing if-head evaluator
		free_if_item,
		pstate);
}

// ----------------------------------------------------------------
// xxx move all frees between allocs & handles. and header-file order too.

static void free_if_head(mlr_dsl_cst_statement_t* pstatement) {
	if_head_state_t* pstate = pstatement->pvstate;

	if (pstate->pif_chain_statements != NULL) {
		for (sllve_t* pe = pstate->pif_chain_statements->phead; pe != NULL; pe = pe->pnext)
			mlr_dsl_cst_statement_free(pe->pvvalue);
		sllv_free(pstate->pif_chain_statements);
	}

	free(pstate);
}

static void free_if_item(mlr_dsl_cst_statement_t* pstatement) {
	if_item_state_t* pstate = pstatement->pvstate;

	pstate->pexpression_evaluator->pfree_func(pstate->pexpression_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_if_head(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	if_head_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->pif_chain_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pitem_statement = pe->pvvalue;
		if_item_state_t* pitem_state = pitem_statement->pvstate;
		rval_evaluator_t* pexpression_evaluator = pitem_state->pexpression_evaluator;

		mv_t val = pexpression_evaluator->pprocess_func(pexpression_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
				local_stack_subframe_enter(pframe, pitem_statement->pblock->subframe_var_count);

				pstatement->pblock_handler(pitem_statement->pblock, pvars, pcst_outputs);

				local_stack_subframe_exit(pframe, pitem_statement->pblock->subframe_var_count);
				break;
			}
		}
	}
}

// ================================================================
typedef struct _while_state_t {
	rval_evaluator_t* pexpression_evaluator;
} while_state_t;

static mlr_dsl_cst_statement_handler_t handle_while;
static mlr_dsl_cst_statement_freer_t free_while;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_while(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	while_state_t* pstate = mlr_malloc_or_die(sizeof(while_state_t));

	pstate->pexpression_evaluator = NULL;

	// Left child node is the AST for the boolean expression.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pright->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pright->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pchild_statement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock->pstatements, pchild_statement);
	}

	pstate->pexpression_evaluator = rval_evaluator_alloc_from_ast(
		pleft, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_while,
		pblock,
		handle_statement_block_with_break_continue,
		free_while,
		pstate);
}

// ----------------------------------------------------------------
// xxx move all frees between allocs & handles. and header-file order too.

static void free_while(mlr_dsl_cst_statement_t* pstatement) { // xxx
	while_state_t* pstate = pstatement->pvstate;

	pstate->pexpression_evaluator->pfree_func(pstate->pexpression_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_while(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	while_state_t* pstate = pstatement->pvstate;

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	rval_evaluator_t* pexpression_evaluator = pstate->pexpression_evaluator;

	while (TRUE) {
		mv_t val = pexpression_evaluator->pprocess_func(pexpression_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
					break;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}
			} else {
				break;
			}
		} else {
			break;
		}
	}

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}


// ================================================================
typedef struct _do_while_state_t {
	rval_evaluator_t* pexpression_evaluator;
} do_while_state_t;

static mlr_dsl_cst_statement_handler_t handle_do_while;
static mlr_dsl_cst_statement_freer_t free_do_while;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_do_while(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	do_while_state_t* pstate = mlr_malloc_or_die(sizeof(do_while_state_t));

	pstate->pexpression_evaluator = NULL;

	// Left child node is the list of statements in the body.
	// Right child node is the AST for the boolean expression.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pright->subframe_var_count);

	for (sllve_t* pe = pleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pchild_statement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock->pstatements, pchild_statement);
	}

	pstate->pexpression_evaluator = rval_evaluator_alloc_from_ast(
		pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_do_while,
		pblock,
		handle_statement_block_with_break_continue,
		free_do_while,
		pstate);
}

// ----------------------------------------------------------------
// xxx move all frees between allocs & handles. and header-file order too.

static void free_do_while(mlr_dsl_cst_statement_t* pstatement) { // xxx
	do_while_state_t* pstate = pstatement->pvstate;

	pstate->pexpression_evaluator->pfree_func(pstate->pexpression_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_do_while(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	do_while_state_t* pstate = pstatement->pvstate;

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	rval_evaluator_t* pexpression_evaluator = pstate->pexpression_evaluator;

	while (TRUE) {
		pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			// don't skip the boolean test
		}

		mv_t val = pexpression_evaluator->pprocess_func(pexpression_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (!val.u.boolv) {
				break;
			}
		} else {
			break;
		}
	}

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}
