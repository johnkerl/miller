#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
static void handle_return_void(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	pvars->return_state.returned = TRUE;
}

static void free_return_void(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
}

mlr_dsl_cst_statement_t* alloc_return_void(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_return_void,
		free_return_void,
		NULL);
}

// ================================================================
typedef struct _return_value_state_t {
	rxval_evaluator_t* preturn_value_xevaluator;
} return_value_state_t;

// ----------------------------------------------------------------
static void return_value_func(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	return_value_state_t* pstate = pstatement->pvstate;
	rxval_evaluator_t* pxev = pstate->preturn_value_xevaluator;
	boxed_xval_t retval = pxev->pprocess_func(pxev->pvstate, pvars);

	if (retval.is_ephemeral) {
		pvars->return_state.retval = retval;
	} else {
		pvars->return_state.retval = box_ephemeral_xval(mlhmmv_xvalue_copy(&retval.xval));
	}

	pvars->return_state.returned = TRUE;
}

// ----------------------------------------------------------------
static void return_value_free(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	return_value_state_t* pstate = pstatement->pvstate;
	rxval_evaluator_t* pxev = pstate->preturn_value_xevaluator;
	pxev->pfree_func(pxev);
	free(pstate);
}

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_return_value(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags)
{
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;

	return_value_state_t* pstate = mlr_malloc_or_die(sizeof(return_value_state_t));

	pstate->preturn_value_xevaluator = rxval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr,
		type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		return_value_func,
		return_value_free,
		pstate);
}
