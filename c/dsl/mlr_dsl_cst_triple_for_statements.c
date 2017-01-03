#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

static mlr_dsl_cst_statement_handler_t handle_triple_for;
static mlr_dsl_cst_statement_freer_t free_triple_for;

typedef struct _triple_for_state_t {
	sllv_t* ptriple_for_start_statements;
	sllv_t* ptriple_for_pre_continuation_statements;
	rval_evaluator_t* ptriple_for_continuation_evaluator;
	sllv_t* ptriple_for_update_statements;
} triple_for_state_t;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_triple_for(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_ast_node_t* pstart_statements_node        = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pcontinuation_statements_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pupdate_statements_node       = pnode->pchildren->phead->pnext->pnext->pvvalue;
	mlr_dsl_ast_node_t* pbody_statements_node         = pnode->pchildren->phead->pnext->pnext->pnext->pvvalue;

	triple_for_state_t* pstate = mlr_malloc_or_die(sizeof(triple_for_state_t));

	pstate->ptriple_for_start_statements = sllv_alloc();
	for (sllve_t* pe = pstart_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pstate->ptriple_for_start_statements,
			mlr_dsl_cst_alloc_statement(
				pcst, pbody_ast_node, type_inferencing, context_flags & ~IN_BREAKABLE));
	}

	// Continuation statements are split into the final boolean, and the statements before (if any).
	pstate->ptriple_for_pre_continuation_statements = sllv_alloc();
	// Empty continuation for triple-for is an implicit TRUE.
	if (pcontinuation_statements_node->pchildren->length == 0) {
		pstate->ptriple_for_continuation_evaluator = rval_evaluator_alloc_from_boolean(TRUE);
	} else {
		for (
			sllve_t* pe = pcontinuation_statements_node->pchildren->phead;
			pe != NULL && pe->pnext != NULL;
			pe = pe->pnext
		)
		{
			mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
			sllv_append(pstate->ptriple_for_pre_continuation_statements,
				mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
				type_inferencing, context_flags & ~IN_BREAKABLE));
		}
		mlr_dsl_ast_node_t* pfinal_continuation_statement_node =
			pcontinuation_statements_node->pchildren->ptail->pvvalue;
		if (mlr_dsl_ast_node_cannot_be_bare_boolean(pfinal_continuation_statement_node)) {
			fprintf(stderr,
				"%s: the final triple-for continutation statement must be a bare boolean.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		pstate->ptriple_for_continuation_evaluator = rval_evaluator_alloc_from_ast(
			pfinal_continuation_statement_node, pcst->pfmgr,
			type_inferencing, (context_flags & ~IN_BREAKABLE) | IN_TRIPLE_FOR_CONTINUE);
	}

	pstate->ptriple_for_update_statements = sllv_alloc();
	for (sllve_t* pe = pupdate_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pstate->ptriple_for_update_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags & ~IN_BREAKABLE));
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pbody_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_triple_for,
		pblock,
		mlr_dsl_cst_handle_statement_block_with_break_continue,
		free_triple_for,
		pstate);
}

static void free_triple_for(mlr_dsl_cst_statement_t* pstatement) {

	triple_for_state_t* pstate = pstatement->pvstate;

	if (pstate->ptriple_for_start_statements != NULL) {
		for (sllve_t* pe = pstate->ptriple_for_start_statements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* ps = pe->pvvalue;
			mlr_dsl_cst_statement_free(ps);
		}
		sllv_free(pstate->ptriple_for_start_statements);
	}

	if (pstate->ptriple_for_pre_continuation_statements != NULL) {
		for (sllve_t* pe = pstate->ptriple_for_pre_continuation_statements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* ps = pe->pvvalue;
			mlr_dsl_cst_statement_free(ps);
		}
		sllv_free(pstate->ptriple_for_pre_continuation_statements);
	}

	if (pstate->ptriple_for_continuation_evaluator != NULL) {
		pstate->ptriple_for_continuation_evaluator->pfree_func(pstate->ptriple_for_continuation_evaluator);
	}

	if (pstate->ptriple_for_update_statements != NULL) {
		for (sllve_t* pe = pstate->ptriple_for_update_statements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* ps = pe->pvvalue;
			mlr_dsl_cst_statement_free(ps);
		}
		sllv_free(pstate->ptriple_for_update_statements);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_triple_for(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	triple_for_state_t* pstate = pstatement->pvstate;

	// Start statements
	mlr_dsl_cst_handle_statement_list(pstate->ptriple_for_start_statements, pvars, pcst_outputs);

	while (TRUE) {

		// Continuation statements:
		// * all but the last one are simply executed ...
		mlr_dsl_cst_handle_statement_list(pstate->ptriple_for_pre_continuation_statements, pvars, pcst_outputs);
		// * ... and the last one is used to determine continuation:
		rval_evaluator_t* pev = pstate->ptriple_for_continuation_evaluator;
		mv_t val = pev->pprocess_func(pev->pvstate, pvars);
		if (mv_is_non_null(&val))
			mv_set_boolean_strict(&val);
		if (!val.u.boolv)
			break;

		// Body statements
		mlr_dsl_cst_handle_statement_block_with_break_continue(pstatement->pblock, pvars, pcst_outputs);

		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}

		// Update statements
		mlr_dsl_cst_handle_statement_list(pstate->ptriple_for_update_statements, pvars, pcst_outputs);
	}

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}
