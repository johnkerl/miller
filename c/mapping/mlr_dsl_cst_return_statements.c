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
