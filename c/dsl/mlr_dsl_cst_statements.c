#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
// The Lemon parser in parsing/mlr_dsl_parse.y builds up an abstract syntax tree
// specifically for the CST builder here.
//
// For clearer visuals on what the ASTs look like:
// * See parsing/mlr_dsl_parse.y
// * See reg_test/run's filter -v and put -v outputs, e.g. in reg_test/expected/out
// * Do "mlr -n put -v 'your expression goes here'"
// ================================================================

// ================================================================
cst_statement_block_t* cst_statement_block_alloc(int subframe_var_count) {
	cst_statement_block_t* pblock = mlr_malloc_or_die(sizeof(cst_statement_block_t));

	pblock->subframe_var_count = subframe_var_count;
	pblock->pstatements     = sllv_alloc();

	return pblock;
}

// ----------------------------------------------------------------
void cst_statement_block_free(cst_statement_block_t* pblock, context_t* pctx) {
	if (pblock == NULL)
		return;
	for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_free(pe->pvvalue, pctx);
	}
	sllv_free(pblock->pstatements);
	free(pblock);
}

// ================================================================
// ALLOCATORS
// ================================================================
cst_top_level_statement_block_t* cst_top_level_statement_block_alloc(int max_var_depth, int subframe_var_count) {
	cst_top_level_statement_block_t* ptop_level_block = mlr_malloc_or_die(sizeof(cst_top_level_statement_block_t));

	ptop_level_block->max_var_depth = max_var_depth;
	ptop_level_block->pframe        = local_stack_frame_alloc(max_var_depth);
	ptop_level_block->pblock        = cst_statement_block_alloc(subframe_var_count);

	return ptop_level_block;
}

// ----------------------------------------------------------------
void cst_top_level_statement_block_free(cst_top_level_statement_block_t* ptop_level_block, context_t* pctx) {
	if (ptop_level_block == NULL)
		return;
	local_stack_frame_free(ptop_level_block->pframe);
	cst_statement_block_free(ptop_level_block->pblock, pctx);
	free(ptop_level_block);
}

// ================================================================
// The parser accepts many things that are invalid, e.g.
// * begin{end{}} -- begin/end not at top level
// * begin{$x=1} -- references to stream records at begin/end
// * break/continue outside of for/while/do-while
// * $x=x -- boundvars outside of for-loop variable bindings
//
// All of the above are enforced here by the CST builder, which takes the parser's output AST as
// input.  This is done (a) to keep the parser from being overly complex, and (b) so we can get much
// more informative error messages in C than in Lemon ('syntax error').
//
// In this file we set up left-hand sides for assignments, as well as right-hand sides for emit and
// unset.  Most right-hand sides are set up in rval_expr_evaluators.c so the context_flags are
// passed through to there as well.

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_FUNC_DEF:
		fprintf(stderr, "%s: func statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_SUBR_DEF:
		fprintf(stderr, "%s: subr statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_BEGIN:
		fprintf(stderr, "%s: begin statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_END:
		fprintf(stderr, "%s: end statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_RETURN_VALUE:
		if (!(context_flags & IN_FUNC_DEF)) {
			fprintf(stderr, "%s: return-value statements are only valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_return_value(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_RETURN_VOID:
		if (!(context_flags & IN_SUBR_DEF)) {
			fprintf(stderr, "%s: return-void statements are only valid within subr blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_return_void(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_SUBR_CALLSITE:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: subroutine calls are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_subr_callsite_statement(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:
		return alloc_conditional_block(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_IF_HEAD:
		return alloc_if_head(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_WHILE:
		return alloc_while(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_DO_WHILE:
		return alloc_do_while(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_SREC:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_for_srec(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_SREC_KEY_ONLY:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_for_srec_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_OOSVAR:
		return alloc_for_map(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY:
		return alloc_for_map_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_LOCAL_MAP:
		return alloc_for_map(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_LOCAL_MAP_KEY_ONLY:
		return alloc_for_map_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_MAP_LITERAL:
		return alloc_for_map(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_MAP_LITERAL_KEY_ONLY:
		return alloc_for_map_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_FOR_FUNC_RETVAL:
		return alloc_for_map(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_FUNC_RETVAL_KEY_ONLY:
		return alloc_for_map_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_TRIPLE_FOR:
		return alloc_triple_for(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_BREAK:
		if (!(context_flags & IN_BREAKABLE)) {
			fprintf(stderr, "%s: break statements are only valid within for, while, or do-while.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_break(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_CONTINUE:
		if (!(context_flags & IN_BREAKABLE)) {
			fprintf(stderr, "%s: break statements are only valid within for, while, or do-while.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_continue(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_ANY);

	case MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_NUMERIC);
		break;

	case MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_INT);
		break;

	case MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_FLOAT);
		break;

	case MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_BOOLEAN);
		break;

	case MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_STRING);
		break;

	case MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags, TYPE_MASK_MAP);
		break;

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT:
		return alloc_nonindexed_local_variable_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT:
		return alloc_indexed_local_variable_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_indirect_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_positional_srec_name_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:
		return alloc_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments from $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT:
		return alloc_full_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments from $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_full_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_full_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_ENV_ASSIGNMENT:
		return alloc_env_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_UNSET:
		return alloc_unset(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_TEE:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: tee statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_tee(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_EMITF:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitf statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emitf(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_EMITP:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitp statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit(pcst, pnode, type_inferencing, context_flags, TRUE);
		break;
	case MD_AST_NODE_TYPE_EMIT:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emit statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit(pcst, pnode, type_inferencing, context_flags, FALSE);
		break;

	case MD_AST_NODE_TYPE_EMITP_LASHED:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitp statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit_lashed(pcst, pnode, type_inferencing, context_flags, TRUE);
		break;
	case MD_AST_NODE_TYPE_EMIT_LASHED:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emit statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit_lashed(pcst, pnode, type_inferencing, context_flags, FALSE);
		break;

	case MD_AST_NODE_TYPE_FILTER:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: filter statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_MLR_FILTER) {
			fprintf(stderr, "%s filter: expressions must not also contain the \"filter\" keyword.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_filter(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_DUMP:
		return alloc_dump(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_PRINT:
		return alloc_print(pcst, pnode, type_inferencing, context_flags, "\n");
		break;

	case MD_AST_NODE_TYPE_PRINTN:
		return alloc_print(pcst, pnode, type_inferencing, context_flags, "");
		break;

	default:
		return alloc_bare_boolean(pcst, pnode, type_inferencing, context_flags);
		break;
	}
}

// ----------------------------------------------------------------
// mlr put and mlr filter are almost entirely the same code. The key difference is that the final
// statement for the latter must be a bare boolean expression.

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_final_filter_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int negate_final_filter, int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_FILTER:
		fprintf(stderr, "%s filter: expressions must not also contain the \"filter\" keyword.\n",
			MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_FUNC_DEF:
	case MD_AST_NODE_TYPE_SUBR_DEF:
	case MD_AST_NODE_TYPE_BEGIN:
	case MD_AST_NODE_TYPE_END:
	case MD_AST_NODE_TYPE_RETURN_VALUE:
	case MD_AST_NODE_TYPE_RETURN_VOID:
	case MD_AST_NODE_TYPE_SUBR_CALLSITE:
	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:
	case MD_AST_NODE_TYPE_IF_HEAD:
	case MD_AST_NODE_TYPE_WHILE:
	case MD_AST_NODE_TYPE_DO_WHILE:
	case MD_AST_NODE_TYPE_FOR_SREC:
	case MD_AST_NODE_TYPE_FOR_SREC_KEY_ONLY:
	case MD_AST_NODE_TYPE_FOR_OOSVAR:
	case MD_AST_NODE_TYPE_TRIPLE_FOR:
	case MD_AST_NODE_TYPE_BREAK:
	case MD_AST_NODE_TYPE_CONTINUE:
	case MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT:
	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME_ASSIGNMENT:
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT:
	case MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_UNSET:
	case MD_AST_NODE_TYPE_TEE:
	case MD_AST_NODE_TYPE_EMITF:
	case MD_AST_NODE_TYPE_EMITP:
	case MD_AST_NODE_TYPE_EMIT:
	case MD_AST_NODE_TYPE_EMITP_LASHED:
	case MD_AST_NODE_TYPE_EMIT_LASHED:
	case MD_AST_NODE_TYPE_DUMP:
	case MD_AST_NODE_TYPE_PRINT:
	case MD_AST_NODE_TYPE_PRINTN:
		fprintf(stderr, "%s: filter expressions must end in a final boolean statement.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	default:
		return alloc_final_filter(pcst, pnode, negate_final_filter, type_inferencing, context_flags);
		break;
	}
}

// ----------------------------------------------------------------
// For used by constructors of subclasses of mlr_dsl_cst_statement_t.

mlr_dsl_cst_statement_t* mlr_dsl_cst_statement_valloc(
	mlr_dsl_ast_node_t*              past_node,
	mlr_dsl_cst_statement_handler_t* pstatement_handler,
	mlr_dsl_cst_statement_freer_t*   pstatement_freer,
	void*                            pvstate)
{
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));
	pstatement->past_node           = past_node;
	pstatement->pstatement_handler  = pstatement_handler;
	pstatement->pblock              = NULL;
	pstatement->pblock_handler      = NULL;
	pstatement->pstatement_freer    = pstatement_freer;
	pstatement->pvstate             = pvstate;
	return pstatement;
}

mlr_dsl_cst_statement_t* mlr_dsl_cst_statement_valloc_with_block(
	mlr_dsl_ast_node_t*              past_node,
	mlr_dsl_cst_statement_handler_t* pstatement_handler,
	cst_statement_block_t*           pblock,
	mlr_dsl_cst_block_handler_t*     pblock_handler,
	mlr_dsl_cst_statement_freer_t*   pstatement_freer,
	void*                            pvstate)
{
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));
	pstatement->past_node           = past_node;
	pstatement->pstatement_handler  = pstatement_handler;
	pstatement->pblock              = pblock;
	pstatement->pblock_handler      = pblock_handler;
	pstatement->pstatement_freer    = pstatement_freer;
	pstatement->pvstate             = pvstate;
	return pstatement;
}

// ----------------------------------------------------------------
void mlr_dsl_cst_statement_free(mlr_dsl_cst_statement_t* pstatement, context_t* pctx) {

	if (pstatement->pstatement_freer != NULL) {
		pstatement->pstatement_freer(pstatement, pctx);
	}

	cst_statement_block_free(pstatement->pblock, pctx);

	free(pstatement);
}

// ================================================================
// Top-level entry point for statement-handling, e.g. from mapper_put.

void mlr_dsl_cst_handle_top_level_statement_blocks(
	sllv_t*      ptop_level_blocks, // block bodies for begins, main, ends
	variables_t* pvars,
	cst_outputs_t* pcst_outputs)
{
	for (sllve_t* pe = ptop_level_blocks->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_handle_top_level_statement_block(pe->pvvalue, pvars, pcst_outputs);
	}
}

void mlr_dsl_cst_handle_top_level_statement_block(
	cst_top_level_statement_block_t* ptop_level_block,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs)
{
	local_stack_push(pvars->plocal_stack, local_stack_frame_enter(ptop_level_block->pframe));
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, ptop_level_block->pblock->subframe_var_count);

	mlr_dsl_cst_handle_statement_block(ptop_level_block->pblock, pvars, pcst_outputs);

	local_stack_subframe_exit(pframe, ptop_level_block->pblock->subframe_var_count);
	local_stack_frame_exit(local_stack_pop(pvars->plocal_stack));
}

// ================================================================
// HANDLERS
// ================================================================
// This is for statement lists not recursively contained within a loop body -- including the
// main/begin/end statements.  Since there is no containing loop body, there is no need to check
// for break or continue flags after each statement.
void mlr_dsl_cst_handle_statement_block(
	cst_statement_block_t* pblock,
	variables_t*           pvars,
	cst_outputs_t*         pcst_outputs)
{
	if (pvars->trace_execution) {
		for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			fprintf(stdout, "TRACE ");
			mlr_dsl_ast_node_pretty_fprint(pstatement->past_node, stdout);
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
			if (pvars->return_state.returned) {
				break;
			}
		}
	} else {
		for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
			if (pvars->return_state.returned) {
				break;
			}
		}
	}
}

// This is for statement lists recursively contained within a loop body.
// It checks for break or continue flags after each statement.
void mlr_dsl_cst_handle_statement_block_with_break_continue(
	cst_statement_block_t* pblock,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	if (pvars->trace_execution) {
		for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			fprintf(stdout, "TRACE ");
			mlr_dsl_ast_node_pretty_fprint(pstatement->past_node, stdout);
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
			if (pvars->return_state.returned) {
				break;
			}
		}
	} else {
		for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
			if (pvars->return_state.returned) {
				break;
			}
		}
	}
}

// Triple-for start/continuation/update statement lists
void mlr_dsl_cst_handle_statement_list(
	sllv_t*        pstatements,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	if (pvars->trace_execution) {
		for (sllve_t* pe = pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			fprintf(stdout, "TRACE ");
			mlr_dsl_ast_node_pretty_fprint(pstatement->past_node, stdout);
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
		}
	} else {
		for (sllve_t* pe = pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
		}
	}
}
