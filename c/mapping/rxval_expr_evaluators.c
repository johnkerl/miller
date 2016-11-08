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
#include "mapping/context_flags.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

rxval_evaluator_t* rxval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{

//	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	if (pnode->pchildren == NULL) {
//		// leaf node
//		switch(pnode->type) {
//
//		case MD_AST_NODE_TYPE_FIELD_NAME:
//			if (context_flags & IN_BEGIN_OR_END) {
//				fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
//					MLR_GLOBALS.bargv0);
//				exit(1);
//			}
//			return rval_evaluator_alloc_from_field_name(pnode->text, type_inferencing);
//			break;
//
//		case MD_AST_NODE_TYPE_STRING_LITERAL:
//			// In input data such as echo x=3,y=4 | mlr put '$z=$x+$y', the 3 and 4 are strings
//			// which need parsing as integers. But in DSL expression literals such as 'put $z = "3" + 4'
//			// the "3" should not.
//			return rval_evaluator_alloc_from_numeric_literal(pnode->text, TYPE_INFER_STRING_ONLY);
//			break;
//
//		case MD_AST_NODE_TYPE_NUMERIC_LITERAL:
//			// In input data such as echo x=3,y=4 | mlr put '$z=$x+$y', the 3 and 4 are strings
//			// which need parsing as integers. But in DSL expression literals such as 'put $z = "3" + 4'
//			// the "3" should not.
//			return rval_evaluator_alloc_from_numeric_literal(pnode->text, type_inferencing);
//			break;
//
//		case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:
//			return rval_evaluator_alloc_from_boolean_literal(pnode->text);
//			break;
//
//		case MD_AST_NODE_TYPE_REGEXI:
//			return rval_evaluator_alloc_from_numeric_literal(pnode->text, type_inferencing);
//			break;
//
//		case MD_AST_NODE_TYPE_CONTEXT_VARIABLE:
//			return rval_evaluator_alloc_from_context_variable(pnode->text);
//			break;
//
//		case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
//			return rval_evaluator_alloc_from_local_variable(pnode->vardef_frame_relative_index);
//			break;
//
//		default:
//			MLR_INTERNAL_CODING_ERROR();
//			return NULL; // not reached
//			break;
//		}
//
//	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	} else if (pnode->type == MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME) {
//		if (context_flags & IN_BEGIN_OR_END) {
//			fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
//				MLR_GLOBALS.bargv0);
//			exit(1);
//		}
//		return rval_evaluator_alloc_from_indirect_field_name(pnode->pchildren->phead->pvvalue, pfmgr,
//			type_inferencing, context_flags);
//
//	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
//		return rval_evaluator_alloc_from_oosvar_keylist(pnode, pfmgr, type_inferencing, context_flags);
//
//	} else if (pnode->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE) {
//		return rval_evaluator_alloc_from_local_map_keylist(pnode, pfmgr, type_inferencing, context_flags);
//
//	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	} else if (pnode->type == MD_AST_NODE_TYPE_ENV) {
//		return rval_evaluator_alloc_from_environment(pnode, pfmgr, type_inferencing, context_flags);
//
//	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//	} else {
//		MLR_INTERNAL_CODING_ERROR_IF((pnode->type != MD_AST_NODE_TYPE_FUNC_CALLSITE)
//			&& (pnode->type != MD_AST_NODE_TYPE_OPERATOR));
//		return fmgr_alloc_from_operator_or_function_call(pfmgr, pnode, type_inferencing, context_flags);
//	}

	return NULL;
}
