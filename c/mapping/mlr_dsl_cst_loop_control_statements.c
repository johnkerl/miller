#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
static mlr_dsl_cst_statement_freer_t free_break;
static mlr_dsl_cst_statement_handler_t handle_break;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_break(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_break,
		free_break,
		NULL);
}

static void free_break(mlr_dsl_cst_statement_t* pstatement) {
}

// ----------------------------------------------------------------
static void handle_break(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	loop_stack_set(pvars->ploop_stack, LOOP_BROKEN);
}

// ================================================================
static mlr_dsl_cst_statement_handler_t handle_continue;
static mlr_dsl_cst_statement_freer_t free_continue;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_continue(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_continue,
		free_continue,
		NULL);
}

// ----------------------------------------------------------------
static void handle_continue(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	loop_stack_set(pvars->ploop_stack, LOOP_CONTINUED);
}

static void free_continue(mlr_dsl_cst_statement_t* pstatement) {
}
