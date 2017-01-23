#ifndef KEYLIST_EVALUATORS_H
#define KEYLIST_EVALUATORS_H

#include "containers/sllv.h"
#include "dsl/mlr_dsl_ast.h"
#include "dsl/function_manager.h"

sllv_t* allocate_keylist_evaluators_from_ast_node(
	mlr_dsl_ast_node_t* pnode,
	fmgr_t*             pfmgr,
	int                 type_inferencing,
	int                 context_flags);

#endif // KEYLIST_EVALUATORS_H
