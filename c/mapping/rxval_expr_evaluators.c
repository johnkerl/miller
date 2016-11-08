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

// ----------------------------------------------------------------
rxval_evaluator_t* rxval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{

	switch(pnode->type) {

	case MD_AST_NODE_TYPE_MAP_LITERAL:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FULL_SREC:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR:
		return NULL;
		break;

	default:
		return NULL;
		break;
	}

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
//		MLR_INTERNAL_CODING_ERROR_IF((pnode->type != MD_AST_NODE_TYPE_FUNCTION_CALLSITE)
//			&& (pnode->type != MD_AST_NODE_TYPE_OPERATOR));
//		return fmgr_alloc_from_operator_or_function_call(pfmgr, pnode, type_inferencing, context_flags);
//	}
}

//// ================================================================
//typedef struct _rval_evaluator_from_local_variable_state_t {
//	int vardef_frame_relative_index;
//} rval_evaluator_from_local_variable_state_t;

//mv_t rval_evaluator_from_local_variable_func(void* pvstate, variables_t* pvars) {
//	rval_evaluator_from_local_variable_state_t* pstate = pvstate;
//	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
//	mv_t val = local_stack_frame_get_non_map(pframe, pstate->vardef_frame_relative_index);
//	return mv_copy(&val);
//}

//static void rval_evaluator_from_local_variable_free(rval_evaluator_t* pevaluator) {
//	rval_evaluator_from_local_variable_state_t* pstate = pevaluator->pvstate;
//	free(pstate);
//	free(pevaluator);
//}

//rval_evaluator_t* rval_evaluator_alloc_from_local_variable(int vardef_frame_relative_index) {
//	rval_evaluator_from_local_variable_state_t* pstate = mlr_malloc_or_die(
//		sizeof(rval_evaluator_from_local_variable_state_t));
//	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
//
//	pstate->vardef_frame_relative_index = vardef_frame_relative_index;
//	pevaluator->pprocess_func    = rval_evaluator_from_local_variable_func;
//	pevaluator->pfree_func       = rval_evaluator_from_local_variable_free;
//
//	pevaluator->pvstate = pstate;
//	return pevaluator;
//}

//// ================================================================
//typedef struct _rval_evaluator_local_map_keylist_state_t {
//	int vardef_frame_relative_index;
//	sllv_t* plocal_map_rhs_keylist_evaluators;
//} rval_evaluator_local_map_keylist_state_t;
//
//mv_t rval_evaluator_local_map_keylist_func(void* pvstate, variables_t* pvars) {
//	rval_evaluator_local_map_keylist_state_t* pstate = pvstate;
//
//	int all_non_null_or_error = TRUE;
//	sllmv_t* pmvkeys = evaluate_list(pstate->plocal_map_rhs_keylist_evaluators, pvars, &all_non_null_or_error);
//
//	mv_t rv = mv_absent();
//	if (all_non_null_or_error) {
//		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
//		mv_t val = local_stack_frame_get_map(pframe, pstate->vardef_frame_relative_index, pmvkeys);
//		if (val.type == MT_STRING && *val.u.strv == 0)
//			rv = mv_empty();
//		else
//			rv = mv_copy(&val);
//	}
//
//	sllmv_free(pmvkeys);
//	return rv;
//}

//static void rval_evaluator_local_map_keylist_free(rval_evaluator_t* pevaluator) {
//	rval_evaluator_local_map_keylist_state_t* pstate = pevaluator->pvstate;
//	for (sllve_t* pe = pstate->plocal_map_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
//		rval_evaluator_t* pevaluator = pe->pvvalue;
//		pevaluator->pfree_func(pevaluator);
//	}
//	sllv_free(pstate->plocal_map_rhs_keylist_evaluators);
//	free(pstate);
//	free(pevaluator);
//}

//rval_evaluator_t* rval_evaluator_alloc_from_local_map_keylist(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
//	int type_inferencing, int context_flags)
//{
//	rval_evaluator_local_map_keylist_state_t* pstate = mlr_malloc_or_die(
//		sizeof(rval_evaluator_local_map_keylist_state_t));
//
//	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
//
//	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;
//
//	sllv_t* pkeylist_evaluators = sllv_alloc();
//	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
//		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
//		if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
//			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
//		} else {
//			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pfmgr,
//				type_inferencing, context_flags));
//		}
//	}
//	pstate->plocal_map_rhs_keylist_evaluators = pkeylist_evaluators;
//
//	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
//	pevaluator->pvstate = pstate;
//	pevaluator->pprocess_func = NULL;
//	pevaluator->pprocess_func = rval_evaluator_local_map_keylist_func;
//	pevaluator->pfree_func = rval_evaluator_local_map_keylist_free;
//
//	return pevaluator;
//}
