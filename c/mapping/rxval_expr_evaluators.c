#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/rval_evaluators.h"
#include "mapping/function_manager.h"
#include "mapping/mlr_dsl_cst.h" // xxx only for allocate_keylist_evaluators_from_ast_node -- xxx move
#include "mapping/context_flags.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

// ----------------------------------------------------------------
rxval_evaluator_t* rxval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		//return rxval_evaluator_alloc_from_nonindexed_local_variable(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		//return rxval_evaluator_alloc_from_indexed_local_variable(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		//return rxval_evaluator_alloc_from_oosvar_keylist(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR:
		//return rxval_evaluator_alloc_from_full_oosvar(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FULL_SREC:
		//return rxval_evaluator_alloc_from_full_srec(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		// xxx XXX to do
		//return rxval_evaluator_alloc_from_function_callsite(
			//pnode, pfmgr, type_inferencing, context_flags);
		//return rxval_evaluator_alloc_wrapping_rval(pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_MAP_LITERAL:
		//return rxval_evaluator_alloc_from_map_literal(
			//pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	default:
		//return rxval_evaluator_alloc_wrapping_rval(pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;
	}
}
