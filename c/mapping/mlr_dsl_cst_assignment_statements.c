#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
typedef struct _srec_assignment_state_t {
	char*             srec_lhs_field_name;
	rval_evaluator_t* prhs_evaluator;
} srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_srec_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(srec_assignment_state_t));

	pstate->prhs_evaluator = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_FIELD_NAME);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren != NULL);

	pstate->srec_lhs_field_name = pleft->text;
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_srec_assignment,
		free_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_srec_assignment(mlr_dsl_cst_statement_t* pstatement) {
	srec_assignment_state_t* pstate = pstatement->pvstate;

	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	srec_assignment_state_t* pstate = pstatement->pvstate;

	char* srec_lhs_field_name = pstate->srec_lhs_field_name;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	// Write typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	//
	// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the input had
	// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
	// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
	// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
	// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
	// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
	// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
	// here. And putting something statically allocated minimizes copying/freeing.
	if (mv_is_present(&val)) {
		lhmsmv_put(pvars->ptyped_overlay, srec_lhs_field_name, &val, FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, srec_lhs_field_name, "bug", NO_FREE);
	} else {
		mv_free(&val);
	}
}

// ================================================================
typedef struct _indirect_srec_assignment_state_t {
	rval_evaluator_t* plhs_evaluator;
	rval_evaluator_t* prhs_evaluator;
} indirect_srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_indirect_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_indirect_srec_assignment;

// ----------------------------------------------------------------
// $ mlr --from ../data/small put -v '$[@x] = 1'
// AST ROOT:
// text="block", type=STATEMENT_BLOCK:
//     text="=", type=INDIRECT_SREC_ASSIGNMENT:
//         text="oosvar_keylist", type=OOSVAR_KEYLIST:
//             text="x", type=STRING_LITERAL.
//         text="1", type=NUMERIC_LITERAL.

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_indirect_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	indirect_srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(indirect_srec_assignment_state_t));

	pstate->prhs_evaluator = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	pstate->plhs_evaluator = rval_evaluator_alloc_from_ast(pleft,  pcst->pfmgr, type_inferencing, context_flags);
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_indirect_srec_assignment,
		free_indirect_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_indirect_srec_assignment(mlr_dsl_cst_statement_t* pstatement) {
	indirect_srec_assignment_state_t* pstate = pstatement->pvstate;

	pstate->plhs_evaluator->pfree_func(pstate->plhs_evaluator);
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_indirect_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	indirect_srec_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* plhs_evaluator = pstate->plhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;

	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	char free_flags;
	char* srec_lhs_field_name = mv_format_val(&lval, &free_flags);

	// Write typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	//
	// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the input had
	// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
	// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
	// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
	// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
	// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
	// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
	// here. And putting something statically allocated minimizes copying/freeing.
	if (mv_is_present(&rval)) {
		lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(srec_lhs_field_name), &rval,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, mlr_strdup_or_die(srec_lhs_field_name), "bug", FREE_ENTRY_KEY | FREE_ENTRY_KEY);
	} else {
		mv_free(&rval);
	}

	if (free_flags) {
		free(srec_lhs_field_name);
	}
}

// ================================================================
typedef struct _local_variable_definition_state_t {
	char*              lhs_variable_name;
	int                lhs_frame_relative_index;
	int                lhs_type_mask;
	rval_evaluator_t*  prhs_evaluator;
	rxval_evaluator_t* prhs_xevaluator;
} local_variable_definition_state_t;

static mlr_dsl_cst_statement_handler_t handle_local_variable_definition_from_val;
static mlr_dsl_cst_statement_handler_t handle_local_variable_definition_from_xval;
static mlr_dsl_cst_statement_freer_t free_local_variable_definition;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_local_variable_definition( // xxx XXX mapvars next
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 type_mask)
{
	local_variable_definition_state_t* pstate = mlr_malloc_or_die(
		sizeof(local_variable_definition_state_t));

	pstate->lhs_variable_name        = NULL;
	pstate->lhs_frame_relative_index = MD_UNUSED_INDEX;
	pstate->lhs_type_mask            = 0;
	pstate->prhs_evaluator           = NULL;
	pstate->prhs_xevaluator          = NULL;

	mlr_dsl_ast_node_t* pname_node = pnode->pchildren->phead->pvvalue;
	pstate->lhs_variable_name = pname_node->text;
	MLR_INTERNAL_CODING_ERROR_IF(pname_node->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->lhs_frame_relative_index = pname_node->vardef_frame_relative_index;
	pstate->lhs_type_mask = type_mask;

	mlr_dsl_cst_statement_handler_t* pstatement_handler = NULL;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	switch (prhs_node->type) {
	case MD_AST_NODE_TYPE_MAP_LITERAL:
	case MD_AST_NODE_TYPE_FULL_SREC:
	case MD_AST_NODE_TYPE_FULL_OOSVAR:
	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(
			prhs_node, pcst->pfmgr, type_inferencing, context_flags);
		pstatement_handler = handle_local_variable_definition_from_xval;
		break;

	default:
		pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);
		pstatement_handler = handle_local_variable_definition_from_val;
		break;
	}

	return mlr_dsl_cst_statement_valloc(
		pnode,
		pstatement_handler,
		free_local_variable_definition,
		pstate);
}

// ----------------------------------------------------------------
static void free_local_variable_definition(mlr_dsl_cst_statement_t* pstatement) {
	local_variable_definition_state_t* pstate = pstatement->pvstate;

	if (pstate->prhs_evaluator != NULL) {
		pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_local_variable_definition_from_val( // xxx mapvar
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	local_variable_definition_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator; // xxx mapvar
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_present(&val)) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_define(pframe,
			pstate->lhs_variable_name, pstate->lhs_frame_relative_index,
			pstate->lhs_type_mask, val);
	} else {
		mv_free(&val);
	}
}

// ----------------------------------------------------------------
static void handle_local_variable_definition_from_xval( // xxx mapvar
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	local_variable_definition_state_t* pstate = pstatement->pvstate;

	rxval_evaluator_t* prhs_xevaluator = pstate->prhs_xevaluator; // xxx mapvar
	mlhmmv_value_t xval = prhs_xevaluator->pprocess_func(prhs_xevaluator->pvstate, pvars);

	if (!xval.is_terminal || mv_is_present(&xval.u.mlrval)) { // xxx funcify
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_xdefine(pframe, // xxx rename
			pstate->lhs_variable_name, pstate->lhs_frame_relative_index, pstate->lhs_type_mask,
			xval);
	} else {
		mlhmmv_free_submap(xval); // xxx rename
	}
}

// ================================================================
typedef struct _nonindexed_local_variable_assignment_state_t {
	char*              lhs_variable_name; // For error messages only: stack-index is computed by stack-allocator:
	int                lhs_frame_relative_index;
	rxval_evaluator_t* prhs_evaluator;
} nonindexed_local_variable_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_nonindexed_local_variable_assignment;
static mlr_dsl_cst_statement_freer_t free_nonindexed_local_variable_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_nonindexed_local_variable_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	nonindexed_local_variable_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		nonindexed_local_variable_assignment_state_t));

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren != NULL);

	pstate->lhs_variable_name = pleft->text;
	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->lhs_frame_relative_index = pleft->vardef_frame_relative_index;
	pstate->prhs_evaluator = rxval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_nonindexed_local_variable_assignment,
		free_nonindexed_local_variable_assignment,
		pstate);
}

// xxx libify

//static rval_evaluator_t* fmgr_alloc_from_udf_callsite(fmgr_t* pfmgr, udf_defsite_state_t* pdefsite_state,
//	mlr_dsl_ast_node_t* pnode, char* function_name, int arity, int type_inferencing, int context_flags)
//{
//	// xxx XXX mapvar fmgr_alloc_from_udf_callsite. extend or clone.
//	rval_evaluator_t* pudf_callsite_evaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
//	rval_evaluator_udf_callsite_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_udf_callsite_state_t));
//
//	pstate->arity = pnode->pchildren->length;
//
//	pstate->pevals = mlr_malloc_or_die(pstate->arity * sizeof(rval_evaluator_t*));
//	int i = 0;
//	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
//		mlr_dsl_ast_node_t* parg_node = pe->pvvalue;
//		pstate->pevals[i] = rval_evaluator_alloc_from_ast(parg_node,
//			pfmgr, type_inferencing, context_flags);
//	}
//
//	pstate->args = mlr_malloc_or_die(pstate->arity * sizeof(mv_t));
//	for (i = 0; i < pstate->arity; i++) {
//		pstate->args[i] = mv_absent();
//	}
//
//	pstate->pdefsite_state = pdefsite_state;
//
//	pudf_callsite_evaluator->pvstate = pstate;
//	pudf_callsite_evaluator->pprocess_func = rval_evaluator_udf_callsite_process;
//	pudf_callsite_evaluator->pfree_func = rval_evaluator_udf_callsite_free;
//
//	return pudf_callsite_evaluator;
//}

// ----------------------------------------------------------------
static void free_nonindexed_local_variable_assignment(mlr_dsl_cst_statement_t* pstatement) {
	nonindexed_local_variable_assignment_state_t* pstate = pstatement->pvstate;

	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_nonindexed_local_variable_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	nonindexed_local_variable_assignment_state_t* pstate = pstatement->pvstate;

	rxval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mlhmmv_value_t xval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (!xval.is_terminal || mv_is_present(&xval.u.mlrval)) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_xassign_non_map(pframe, pstate->lhs_frame_relative_index, xval);
	} else {
		mlhmmv_free_submap(xval);
	}
}

// ================================================================
typedef struct _indexed_local_variable_assignment_state_t {
	char*             lhs_variable_name; // For error messages only: stack-index is computed by stack-allocator:
	int               lhs_frame_relative_index;
	sllv_t*           plhs_keylist_evaluators;
	rval_evaluator_t* prhs_evaluator;
} indexed_local_variable_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_indexed_local_variable_assignment;
static mlr_dsl_cst_statement_freer_t free_indexed_local_variable_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_indexed_local_variable_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	indexed_local_variable_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		indexed_local_variable_assignment_state_t));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren == NULL);

	pstate->lhs_variable_name = mlr_strdup_or_die(pleft->text);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->lhs_frame_relative_index = pleft->vardef_frame_relative_index;

	pstate->plhs_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pleft, pcst->pfmgr, type_inferencing, context_flags);

	switch(pright->type) { // xxx XXX mapvar cases

	case MD_AST_NODE_TYPE_MAP_LITERAL:
		// xxx map-build
		// xxx level-copy
		break;

	case MD_AST_NODE_TYPE_FULL_SREC:
		break;

	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		// xxx level-copy
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR:
		// xxx level-copy
		break;

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		// xxx level-copy
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		break;

	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		// xxx clone fmgr_alloc_from_operator_or_function_call
		// xxx parameter-marshal etc.
		// xxx retval etc.
		break;

	default:
		pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);
		break;

	}

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_indexed_local_variable_assignment,
		free_indexed_local_variable_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_indexed_local_variable_assignment(mlr_dsl_cst_statement_t* pstatement) {
	indexed_local_variable_assignment_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->plhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_indexed_local_variable_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	indexed_local_variable_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_present(&rhs_value)) {

		int all_non_null_or_error = TRUE;
		sllmv_t* pmvkeys = evaluate_list(pstate->plhs_keylist_evaluators, pvars,
			&all_non_null_or_error);
		if (all_non_null_or_error) {
			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			local_stack_frame_assign_map(pframe, pstate->lhs_frame_relative_index, pmvkeys, rhs_value);
		}
		sllmv_free(pmvkeys);

	} else {
		mv_free(&rhs_value);
	}
}

// ================================================================
// All assignments produce a mlrval on the RHS and store it on the left -- except if both LHS and RHS
// are oosvars in which case there are recursive copies, or in case of $* on the LHS or RHS.

typedef struct _oosvar_assignment_state_t {
	sllv_t*            plhs_keylist_evaluators;
	rval_evaluator_t*  prhs_evaluator;
	rxval_evaluator_t* prhs_xevaluator;
	sllv_t*            prhs_keylist_evaluators;
} oosvar_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_oosvar_assignment_from_val;
static mlr_dsl_cst_statement_handler_t handle_oosvar_assignment_from_xval;
// xxx rm static mlr_dsl_cst_statement_handler_t handle_oosvar_to_oosvar_assignment;
static mlr_dsl_cst_statement_freer_t free_oosvar_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_oosvar_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	oosvar_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(oosvar_assignment_state_t));

	pstate->plhs_keylist_evaluators = NULL;
	pstate->prhs_evaluator          = NULL;
	pstate->prhs_xevaluator         = NULL;
	pstate->prhs_keylist_evaluators = NULL;

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->type != MD_AST_NODE_TYPE_OOSVAR_KEYLIST);

	pstate->plhs_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		plhs_node, pcst->pfmgr, type_inferencing, context_flags);

	mlr_dsl_cst_statement_handler_t* pstatement_handler = NULL;

	switch (prhs_node->type) {
	case MD_AST_NODE_TYPE_MAP_LITERAL:
	case MD_AST_NODE_TYPE_FULL_SREC:
	case MD_AST_NODE_TYPE_FULL_OOSVAR:
	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(
			prhs_node, pcst->pfmgr, type_inferencing, context_flags);
		pstatement_handler = handle_oosvar_assignment_from_xval;
		break;

	default:
		pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);
		pstatement_handler = handle_oosvar_assignment_from_val;
		break;
	}

	return mlr_dsl_cst_statement_valloc(
		pnode,
		pstatement_handler,
		free_oosvar_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_oosvar_assignment(mlr_dsl_cst_statement_t* pstatement) {
	oosvar_assignment_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->plhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	if (pstate->prhs_evaluator != NULL) {
		pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);
	}
	if (pstate->prhs_xevaluator != NULL) {
		pstate->prhs_xevaluator->pfree_func(pstate->prhs_xevaluator);
	}
	if (pstate->prhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstate->prhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* pev = pe->pvvalue;
			pev->pfree_func(pev);
		}
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_oosvar_assignment_from_val(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	oosvar_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&rhs_value)) {
		int all_non_null_or_error = TRUE;
		sllmv_t* pmvkeys = evaluate_list(pstate->plhs_keylist_evaluators, pvars,
			&all_non_null_or_error);
		if (all_non_null_or_error)
			mlhmmv_put_terminal(pvars->poosvars, pmvkeys, &rhs_value);
		sllmv_free(pmvkeys);
	}
	mv_free(&rhs_value);
}

// ----------------------------------------------------------------
static void handle_oosvar_assignment_from_xval(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	oosvar_assignment_state_t* pstate = pstatement->pvstate;

	int lhs_all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pstate->plhs_keylist_evaluators, pvars,
		&lhs_all_non_null_or_error);

	if (lhs_all_non_null_or_error) {
		rxval_evaluator_t* prhs_xevaluator = pstate->prhs_xevaluator;
		mlhmmv_value_t xval = prhs_xevaluator->pprocess_func(prhs_xevaluator->pvstate, pvars);
		if (!xval.is_terminal || mv_is_present(&xval.u.mlrval)) { // xxx funcify
			mlhmmv_put_value_at_level_aux(pvars->poosvars->proot_level, plhskeys->phead, &xval); // xxx rename
		} else {
			mlhmmv_free_submap(xval); // xxx rename
		}
	}

	sllmv_free(plhskeys);
}

// ================================================================
// All assignments produce a mlrval on the RHS and store it on the left -- except if both LHS and RHS
// are oosvars in which case there are recursive copies, or in case of $* on the LHS or RHS.

typedef struct _oosvar_from_full_srec_assignment_state_t {
	sllv_t* plhs_keylist_evaluators;
} oosvar_from_full_srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_oosvar_from_full_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_oosvar_from_full_srec_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_oosvar_from_full_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	oosvar_from_full_srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		oosvar_from_full_srec_assignment_state_t));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_OOSVAR_KEYLIST);
	MLR_INTERNAL_CODING_ERROR_IF(pright->type != MD_AST_NODE_TYPE_FULL_SREC);

	// xxx XXX mapvar srec := map-lit
	// xxx XXX mapvar srec := non-indexed local
	// xxx XXX mapvar srec := indexed local
	// xxx XXX mapvar srec := udf call

	pstate->plhs_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pleft, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_oosvar_from_full_srec_assignment,
		free_oosvar_from_full_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_oosvar_from_full_srec_assignment(mlr_dsl_cst_statement_t* pstatement) {
	oosvar_from_full_srec_assignment_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->plhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	oosvar_from_full_srec_assignment_state_t* pstate = pstatement->pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pstate->plhs_keylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {

		mlhmmv_level_t* plevel = mlhmmv_get_or_create_level(pvars->poosvars, plhskeys);
		if (plevel != NULL) {

			mlhmmv_clear_level(plevel);

			for (lrece_t* pe = pvars->pinrec->phead; pe != NULL; pe = pe->pnext) {
				mv_t k = mv_from_string(pe->key, NO_FREE); // mlhmmv_put_terminal_from_level will copy
				sllmve_t e = { .value = k, .free_flags = 0, .pnext = NULL };
				mv_t* pomv = lhmsmv_get(pvars->ptyped_overlay, pe->key);
				if (pomv != NULL) {
					mlhmmv_put_terminal_from_level(plevel, &e, pomv);
				} else {
					mv_t v = mv_from_string(pe->value, NO_FREE); // mlhmmv_put_terminal_from_level will copy
					mlhmmv_put_terminal_from_level(plevel, &e, &v);
				}
			}

		}
	}
	sllmv_free(plhskeys);
}

// ================================================================
// All assignments produce a mlrval on the RHS and store it on the left -- except if both LHS and RHS
// are oosvars in which case there are recursive copies, or in case of $* on the LHS or RHS.

typedef struct _full_srec_assignment_state_t {
	rxval_evaluator_t* prhs_xevaluator;
} full_srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_full_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_full_srec_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_full_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	full_srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		full_srec_assignment_state_t));

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->type != MD_AST_NODE_TYPE_FULL_SREC);

	pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(
		prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_full_srec_assignment,
		free_full_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_full_srec_assignment(mlr_dsl_cst_statement_t* pstatement) {
	full_srec_assignment_state_t* pstate = pstatement->pvstate;

	pstate->prhs_xevaluator->pfree_func(pstate->prhs_xevaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_full_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	full_srec_assignment_state_t* pstate = pstatement->pvstate;

	lrec_clear(pvars->pinrec);
	lhmsmv_clear(pvars->ptyped_overlay);

	rxval_evaluator_t* prhs_xevaluator = pstate->prhs_xevaluator;
	mlhmmv_value_t mapval = prhs_xevaluator->pprocess_func(prhs_xevaluator->pvstate, pvars);

	if (mapval.is_terminal) {
		mlhmmv_free_submap(mapval); // xxx rename
	} else {

		for (mlhmmv_level_entry_t* pe = mapval.u.pnext_level->phead; pe != NULL; pe = pe->pnext) {

			mv_t* pkey = &pe->level_key;
			mlhmmv_value_t* pval = &pe->level_value;

			if (pval->is_terminal) { // xxx else collapse-down using json separator?

				char* skey = mv_alloc_format_val(pkey);
				mv_t val = mv_copy(&pval->u.mlrval);
				// Write typed mlrval output to the typed overlay rather than into the lrec
				// (which holds only string values).
				//
				// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the
				// input had "x"=>"abc","y"=>"def" but a previous statement had set "y"=>7.4 and
				// "z"=>"ghi", then an expression right-hand side referring to $y would get the
				// floating-point value 7.4. So we don't need to lrec_put the value here, and
				// moreover should not for two reasons: (1) there is a performance hit of doing
				// throwaway number-to-string formatting -- it's better to do it once at the
				// end; (2) having the string values doubly owned by the typed overlay and the
				// lrec would result in double frees, or awkward bookkeeping. However, the NR
				// variable evaluator reads prec->field_count, so we need to put something here.
				// And putting something statically allocated minimizes copying/freeing.
				lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(skey), &val,
					FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
				lrec_put(pvars->pinrec, skey, "bug", FREE_ENTRY_KEY);
			}
		}
	}
}

// ================================================================
typedef struct _env_assignment_state_t {
	rval_evaluator_t* plhs_evaluator;
	rval_evaluator_t* prhs_evaluator;
} env_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_env_assignment;
static mlr_dsl_cst_statement_freer_t free_env_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_env_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	env_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(env_assignment_state_t));

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_ENV);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren->length != 2);
	mlr_dsl_ast_node_t* pnamenode  = pleft->pchildren->phead->pnext->pvvalue;

	pstate->plhs_evaluator = rval_evaluator_alloc_from_ast(pnamenode, pcst->pfmgr, type_inferencing, context_flags);
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_env_assignment,
		free_env_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_env_assignment(mlr_dsl_cst_statement_t* pstatement) {
	env_assignment_state_t* pstate = pstatement->pvstate;

	pstate->plhs_evaluator->pfree_func(pstate->plhs_evaluator);
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_env_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	env_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* plhs_evaluator = pstate->plhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&lval) && mv_is_present(&rval)) {
		setenv(mlr_strdup_or_die(mv_alloc_format_val(&lval)), mlr_strdup_or_die(mv_alloc_format_val(&rval)), 1);
	}
	mv_free(&lval);
	mv_free(&rval);
}
