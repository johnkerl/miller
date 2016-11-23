#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
static mlr_dsl_cst_statement_handler_t handle_return_void;
static mlr_dsl_cst_statement_freer_t free_return_void;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_void(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_void,
		free_return_void,
		NULL);
}

// ----------------------------------------------------------------
static void handle_return_void(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	pvars->return_state.returned = TRUE;
}

static void free_return_void(mlr_dsl_cst_statement_t* pstatement) {
}

// ================================================================
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_map_literal;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_local_non_map_variable;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_indexed_local_variable;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_oosvar;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_full_oosvar;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_full_srec; // xxx needs grammar support
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_function_callsite;
static mlr_dsl_cst_statement_allocator_t alloc_return_value_from_non_map_valued;

// ----------------------------------------------------------------
// xxx mapvar: special-case retval is map-literal: need vardef_frame_relative_index & keylist evaluators
// xxx mapvar: special-case retval is local map/non-map: need vardef_frame_relative_index & keylist evaluators
// xxx mapvar: special-case retval is oosvar/@*: need keylist evaluators
// xxx mapvar: special-case retval is $* ?
// xxx mapvar: what if 'return g(a,b)' and g is map-valued?

mlr_dsl_cst_statement_t* alloc_return_value(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;

	// $ mlr --from s put -v 'map v = {}'
	// AST ROOT:
	// text="block", type=STATEMENT_BLOCK:
	//     text="map", type=MAP_LOCAL_DEFINITION:
	//         text="v", type=NONINDEXED_LOCAL_VARIABLE.
	//         text="map_literal", type=MAP_LITERAL:

	switch (prhs_node->type) {
	case MD_AST_NODE_TYPE_MAP_LITERAL:
		return alloc_return_value_from_map_literal(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		return alloc_return_value_from_local_non_map_variable(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		return alloc_return_value_from_indexed_local_variable(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		return alloc_return_value_from_oosvar(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_FULL_OOSVAR:
		return alloc_return_value_from_full_oosvar(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_FULL_SREC:
		return alloc_return_value_from_full_srec(pcst, pnode, type_inferencing, context_flags);
		break;

	case  MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		return alloc_return_value_from_function_callsite(pcst, pnode, type_inferencing, context_flags);
		break;

	default:
		return alloc_return_value_from_non_map_valued(pcst, pnode, type_inferencing, context_flags);
		break;

	}
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_handler_t handle_return_value_from_local_non_map_variable;
static mlr_dsl_cst_statement_handler_t handle_return_value_from_indexed_local_variable;
static mlr_dsl_cst_statement_handler_t handle_return_value_from_oosvar;
static mlr_dsl_cst_statement_handler_t handle_return_value_from_full_oosvar;
static mlr_dsl_cst_statement_handler_t handle_return_value_from_full_srec; // xxx needs grammar support
static mlr_dsl_cst_statement_handler_t handle_return_value_from_function_callsite;
static mlr_dsl_cst_statement_handler_t handle_return_value_from_non_map_valued;

// ================================================================
typedef struct _return_value_from_local_non_map_variable_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_local_non_map_variable_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_local_non_map_variable;
static mlr_dsl_cst_statement_freer_t free_return_value_from_local_non_map_variable;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_map_literal(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_local_non_map_variable_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_local_non_map_variable_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_local_non_map_variable,
		free_return_value_from_local_non_map_variable,
		pstate);
}

static void free_return_value_from_local_non_map_variable(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_local_non_map_variable_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_local_non_map_variable(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_local_non_map_variable_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_local_non_map_variable_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_local_non_map_variable,
		free_return_value_from_local_non_map_variable,
		pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_local_non_map_variable(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_local_non_map_variable_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_indexed_local_variable_state_t {

	// For error messages only: stack-index is computed by stack-allocator:
	char* rhs_variable_name;
	int rhs_frame_relative_index;
	sllv_t* prhs_keylist_evaluators;

} return_value_from_indexed_local_variable_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_indexed_local_variable;
static mlr_dsl_cst_statement_freer_t free_return_value_from_indexed_local_variable;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_indexed_local_variable(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_indexed_local_variable_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_indexed_local_variable_state_t));

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;

	pstate->rhs_variable_name = prhs_node->text;

	MLR_INTERNAL_CODING_ERROR_IF(prhs_node->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->rhs_frame_relative_index = prhs_node->vardef_frame_relative_index;

	pstate->prhs_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_indexed_local_variable,
		free_return_value_from_indexed_local_variable,
		pstate);
}

static void free_return_value_from_indexed_local_variable(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_indexed_local_variable_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->prhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_indexed_local_variable(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_indexed_local_variable_state_t* pstate = pstatement->pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->prhs_keylist_evaluators, pvars,
		&all_non_null_or_error);
	if (all_non_null_or_error) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);

		mlhmmv_xvalue_t* pmvalue = local_stack_frame_get_extended_from_indexed(pframe,
			pstate->rhs_frame_relative_index, pmvkeys);

		if (pmvalue == NULL) {
			pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal(mv_absent());
		} else {
			pvars->return_state.retval = mlhmmv_xvalue_copy(pmvalue);
		}

	} else {
		pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal(mv_absent());
	}

	sllmv_free(pmvkeys);

	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_oosvar_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_oosvar_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_oosvar;
static mlr_dsl_cst_statement_freer_t free_return_value_from_oosvar;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_oosvar(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_oosvar_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_oosvar_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_oosvar,
		free_return_value_from_oosvar,
		pstate);
}

static void free_return_value_from_oosvar(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_oosvar_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_oosvar(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_oosvar_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_full_oosvar_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_full_oosvar_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_full_oosvar;
static mlr_dsl_cst_statement_freer_t free_return_value_from_full_oosvar;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_full_oosvar(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_full_oosvar_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_full_oosvar_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_full_oosvar,
		free_return_value_from_full_oosvar,
		pstate);
}

static void free_return_value_from_full_oosvar(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_full_oosvar_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_full_oosvar(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_full_oosvar_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_full_srec_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_full_srec_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_full_srec;
static mlr_dsl_cst_statement_freer_t free_return_value_from_full_srec;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_full_srec(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_full_srec_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_full_srec_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_full_srec,
		free_return_value_from_full_srec,
		pstate);
}

static void free_return_value_from_full_srec(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_full_srec_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_full_srec(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_full_srec_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_function_callsite_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_function_callsite_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_function_callsite;
static mlr_dsl_cst_statement_freer_t free_return_value_from_function_callsite;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_function_callsite(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_function_callsite_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_function_callsite_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, // xxx mapvars
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_function_callsite,
		free_return_value_from_function_callsite,
		pstate);
}

static void free_return_value_from_function_callsite(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_function_callsite_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_function_callsite(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_function_callsite_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}

// ================================================================
typedef struct _return_value_from_non_map_valued_state_t {
	rval_evaluator_t* preturn_value_evaluator;
} return_value_from_non_map_valued_state_t;

static mlr_dsl_cst_statement_handler_t handle_return_value_from_non_map_valued;
static mlr_dsl_cst_statement_freer_t free_return_value_from_non_map_valued;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value_from_non_map_valued(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	return_value_from_non_map_valued_state_t* pstate =
		mlr_malloc_or_die(sizeof(return_value_from_non_map_valued_state_t));

	pstate->preturn_value_evaluator = NULL;

	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstate->preturn_value_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr,
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_value_from_non_map_valued,
		free_return_value_from_non_map_valued,
		pstate);
}

static void free_return_value_from_non_map_valued(mlr_dsl_cst_statement_t* pstatement) {
	return_value_from_non_map_valued_state_t* pstate = pstatement->pvstate;

	pstate->preturn_value_evaluator->pfree_func(pstate->preturn_value_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_return_value_from_non_map_valued(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_from_non_map_valued_state_t* pstate = pstatement->pvstate;

	pvars->return_state.retval = mlhmmv_xvalue_wrap_terminal( // xxx mapvars
		pstate->preturn_value_evaluator->pprocess_func(
			pstate->preturn_value_evaluator->pvstate, pvars));
	pvars->return_state.returned = TRUE;
}
