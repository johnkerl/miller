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
typedef struct _local_non_map_variable_assignment_state_t {

	// For error messages only: stack-index is computed by stack-allocator:
	char* local_lhs_variable_name;

	int local_lhs_frame_relative_index;

	rval_evaluator_t* prhs_evaluator;

} local_non_map_variable_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_local_non_map_variable_assignment;
static mlr_dsl_cst_statement_freer_t free_local_non_map_variable_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_local_non_map_variable_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	local_non_map_variable_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		local_non_map_variable_assignment_state_t));

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_LOCAL_NON_MAP_VARIABLE);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren != NULL);

	pstate->local_lhs_variable_name = pleft->text;
	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->local_lhs_frame_relative_index = pleft->vardef_frame_relative_index;
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_local_non_map_variable_assignment,
		free_local_non_map_variable_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_local_non_map_variable_assignment(mlr_dsl_cst_statement_t* pstatement) {
	local_non_map_variable_assignment_state_t* pstate = pstatement->pvstate;

	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_local_non_map_variable_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	local_non_map_variable_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_present(&val)) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_assign_non_map(pframe, pstate->local_lhs_frame_relative_index, val);
	} else {
		mv_free(&val);
	}
}

// ================================================================
typedef struct _local_map_variable_assignment_state_t {

	// For error messages only: stack-index is computed by stack-allocator:
	char* local_lhs_variable_name; // xxx renames

	int local_lhs_frame_relative_index;

	sllv_t* plhs_keylist_evaluators;

	rval_evaluator_t* prhs_evaluator;

} local_map_variable_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_local_map_variable_assignment;
static mlr_dsl_cst_statement_freer_t free_local_map_variable_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_local_map_variable_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	local_map_variable_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(
		local_map_variable_assignment_state_t));

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_LOCAL_MAP_VARIABLE);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->pchildren == NULL);

	pstate->local_lhs_variable_name = mlr_strdup_or_die(pleft->text);
	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->local_lhs_frame_relative_index = pleft->vardef_frame_relative_index;

	pstate->plhs_keylist_evaluators = sllv_alloc();
	for (sllve_t* pe = pleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
			sllv_append(pstate->plhs_keylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
		} else {
			sllv_append(pstate->plhs_keylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr,
				type_inferencing, context_flags));
		}
	}

	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_local_map_variable_assignment,
		free_local_map_variable_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_local_map_variable_assignment(mlr_dsl_cst_statement_t* pstatement) {
	local_map_variable_assignment_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->plhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_local_map_variable_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	local_map_variable_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_present(&rhs_value)) {

		int all_non_null_or_error = TRUE;
		sllmv_t* pmvkeys = evaluate_list(pstate->plhs_keylist_evaluators, pvars,
			&all_non_null_or_error);
		if (all_non_null_or_error) {
			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			local_stack_frame_assign_map(pframe, pstate->local_lhs_frame_relative_index, pmvkeys, rhs_value);
		}
		sllmv_free(pmvkeys);

	} else {
		mv_free(&rhs_value);
	}
}
