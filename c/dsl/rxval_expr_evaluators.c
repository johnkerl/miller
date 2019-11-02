#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "dsl/context_flags.h"
#include "dsl/keylist_evaluators.h"
#include "dsl/rval_evaluators.h"
#include "dsl/rxval_evaluators.h"
#include "dsl/function_manager.h"
#include "mapping/mapper.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

// ----------------------------------------------------------------
rxval_evaluator_t* rxval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_MAP_LITERAL:
		return rxval_evaluator_alloc_from_map_literal(
			pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		return rxval_evaluator_alloc_from_nonindexed_local_variable(
			pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
		return rxval_evaluator_alloc_from_indexed_local_variable(
			pnode, pfmgr, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_INDEXED_FUNCTION_CALLSITE:
		return rxval_evaluator_alloc_from_indexed_function_call(
			pnode, pfmgr, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
		return rxval_evaluator_alloc_from_oosvar_keylist(
			pnode, pfmgr, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_OOSVAR:
		return rxval_evaluator_alloc_from_full_oosvar(
			pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FULL_SREC:
		return rxval_evaluator_alloc_from_full_srec(
			pnode, pfmgr, type_inferencing, context_flags);
		return NULL;
		break;

	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:
		return fmgr_xalloc_provisional_from_operator_or_function_call(pfmgr, pnode, type_inferencing, context_flags);
		break;

	default:
		return rxval_evaluator_alloc_wrapping_rval(pnode, pfmgr, type_inferencing, context_flags);
		break;
	}
}

// ----------------------------------------------------------------
// Srec assignments have output that is rval not rxval (scalar, not map).  But for
// indexed function calls we have a function which produces rxval, indexed down by a
// keylist, to produce an rval as a final result. This needs special handling.
// Here we invoke the rxval evaluator, then unbox it and return the resulting rval.

typedef struct _rval_evaluator_indexed_function_call_state_t {
	rxval_evaluator_t* prxval_evaluator;
} rval_evaluator_indexed_function_call_state_t;

static mv_t rval_evaluator_indexed_function_call_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_indexed_function_call_state_t* pstate = pvstate;

	rxval_evaluator_t* prxval_evaluator = pstate->prxval_evaluator;
	boxed_xval_t boxed_xval = prxval_evaluator->pprocess_func(prxval_evaluator->pvstate, pvars);
	if (boxed_xval.xval.is_terminal) {
		return boxed_xval.xval.terminal_mlrval;
	} else {
		return mv_absent();
	}

}
static void rval_evaluator_indexed_function_call_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_indexed_function_call_state_t* pstate = pevaluator->pvstate;
	pstate->prxval_evaluator->pfree_func(pstate->prxval_evaluator);
	free(pstate);
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_indexed_function_call(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rval_evaluator_indexed_function_call_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_indexed_function_call_state_t)
	);
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	pstate->prxval_evaluator = rxval_evaluator_alloc_from_indexed_function_call(
		pnode, pfmgr, type_inferencing, context_flags);

	pevaluator->pprocess_func = rval_evaluator_indexed_function_call_func;
	pevaluator->pfree_func    = rval_evaluator_indexed_function_call_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// Map-literal input:
//
// {
//   "a" : 1,
//   "b" : {
//     "x" : 7,
//     "y" : 8,
//   },
//   "c" : 3,
// }
//
// Map-literal AST:
//
// $ mlr --from s put -v -q 'm={"a":NR,"b":{"x":999},"c":3};dump m'
// text="block", type=STATEMENT_BLOCK:
//     text="=", type=NONINDEXED_LOCAL_ASSIGNMENT:
//         text="m", type=NONINDEXED_LOCAL_VARIABLE.
//         text="map_literal", type=MAP_LITERAL:
//             text="mappair", type=MAP_LITERAL_PAIR:
//                 text="mapkey", type=MAP_LITERAL:
//                     text="a", type=STRING_LITERAL.
//                 text="mapval", type=MAP_LITERAL:
//                     text="NR", type=CONTEXT_VARIABLE.
//             text="mappair", type=MAP_LITERAL_PAIR:
//                 text="mapkey", type=MAP_LITERAL:
//                     text="b", type=STRING_LITERAL.
//                 text="mapval", type=MAP_LITERAL:
//                     text="map_literal", type=MAP_LITERAL:
//                         text="mappair", type=MAP_LITERAL_PAIR:
//                             text="mapkey", type=MAP_LITERAL:
//                                 text="x", type=STRING_LITERAL.
//                             text="mapval", type=MAP_LITERAL:
//                                 text="999", type=NUMERIC_LITERAL.
//             text="mappair", type=MAP_LITERAL_PAIR:
//                 text="mapkey", type=MAP_LITERAL:
//                     text="c", type=STRING_LITERAL.
//                 text="mapval", type=MAP_LITERAL:
//                     text="3", type=NUMERIC_LITERAL.
//     text="dump", type=DUMP:
//         text=">", type=FILE_WRITE:
//             text="stdout", type=STDOUT:
//         text="m", type=NONINDEXED_LOCAL_VARIABLE.

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
typedef struct _map_literal_list_evaluator_t {
	sllv_t* pkvpair_evaluators;
} map_literal_list_evaluator_t;

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
typedef struct _map_literal_kvpair_evaluator_t {
	rval_evaluator_t*             pkey_evaluator;
	int                           is_terminal;
	rxval_evaluator_t*            pxval_evaluator;
	map_literal_list_evaluator_t* plist_evaluator;
} map_literal_kvpair_evaluator_t;

// ----------------------------------------------------------------
static map_literal_list_evaluator_t* allocate_map_literal_evaluator_from_ast(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	map_literal_list_evaluator_t* plist_evaluator = mlr_malloc_or_die(
		sizeof(map_literal_list_evaluator_t));
	plist_evaluator->pkvpair_evaluators = sllv_alloc();
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_MAP_LITERAL);
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {

		map_literal_kvpair_evaluator_t* pkvpair = mlr_malloc_or_die(
			sizeof(map_literal_kvpair_evaluator_t));
		*pkvpair = (map_literal_kvpair_evaluator_t) {
			.pkey_evaluator  = NULL,
			.is_terminal     = TRUE,
			.pxval_evaluator = NULL,
			.plist_evaluator = NULL,
		};

		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		MLR_INTERNAL_CODING_ERROR_IF(pchild->type != MD_AST_NODE_TYPE_MAP_LITERAL_PAIR);

		mlr_dsl_ast_node_t* pleft = pchild->pchildren->phead->pvvalue;
		MLR_INTERNAL_CODING_ERROR_IF(pleft->type != MD_AST_NODE_TYPE_MAP_LITERAL_KEY);
		mlr_dsl_ast_node_t* pkeynode = pleft->pchildren->phead->pvvalue;
		pkvpair->pkey_evaluator = rval_evaluator_alloc_from_ast(pkeynode, pfmgr,
			type_inferencing, context_flags);

		mlr_dsl_ast_node_t* pright = pchild->pchildren->phead->pnext->pvvalue;
		mlr_dsl_ast_node_t* pvalnode = pright->pchildren->phead->pvvalue;
		if (pright->type == MD_AST_NODE_TYPE_MAP_LITERAL_VALUE) {
			pkvpair->pxval_evaluator = rxval_evaluator_alloc_from_ast(pvalnode, pfmgr,
				type_inferencing, context_flags);
		} else if (pright->type == MD_AST_NODE_TYPE_MAP_LITERAL) {
			pkvpair->is_terminal = FALSE;
			pkvpair->plist_evaluator = allocate_map_literal_evaluator_from_ast(
				pvalnode, pfmgr, type_inferencing, context_flags);
		} else {
			MLR_INTERNAL_CODING_ERROR();
		}

		sllv_append(plist_evaluator->pkvpair_evaluators, pkvpair);
	}
	return plist_evaluator;
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
typedef struct _rxval_evaluator_from_map_literal_state_t {
	map_literal_list_evaluator_t* proot_list_evaluator;
} rxval_evaluator_from_map_literal_state_t;

static void rxval_evaluator_from_map_literal_aux(
	rxval_evaluator_from_map_literal_state_t* pstate,
	map_literal_list_evaluator_t*             plist_evaluator,
	mlhmmv_level_t*                           plevel,
	variables_t*                              pvars)
{
	for (sllve_t* pe = plist_evaluator->pkvpair_evaluators->phead; pe != NULL; pe = pe->pnext) {
		map_literal_kvpair_evaluator_t* pkvpair = pe->pvvalue;

		// mlhmmv_level_put_terminal will copy keys and values
		mv_t mvkey = pkvpair->pkey_evaluator->pprocess_func(pkvpair->pkey_evaluator->pvstate, pvars);
		if (pkvpair->is_terminal) {
			boxed_xval_t boxed_xval = pkvpair->pxval_evaluator->pprocess_func(
				pkvpair->pxval_evaluator->pvstate, pvars);
			if (boxed_xval.is_ephemeral) {
				mlhmmv_level_put_xvalue_singly_keyed(plevel, &mvkey, &boxed_xval.xval);
			} else {
				mlhmmv_xvalue_t copy_xval = mlhmmv_xvalue_copy(&boxed_xval.xval);
				mlhmmv_level_put_xvalue_singly_keyed(plevel, &mvkey, &copy_xval);
			}
		} else {
			mlhmmv_level_t* pnext_level = mlhmmv_level_put_empty_map(plevel, &mvkey);
			rxval_evaluator_from_map_literal_aux(pstate, pkvpair->plist_evaluator, pnext_level, pvars);
		}
        mv_free(&mvkey);
	}
}

static boxed_xval_t rxval_evaluator_from_map_literal_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_map_literal_state_t* pstate = pvstate;

	mlhmmv_xvalue_t xval = mlhmmv_xvalue_alloc_empty_map();

	rxval_evaluator_from_map_literal_aux(pstate, pstate->proot_list_evaluator, xval.pnext_level, pvars);

	return (boxed_xval_t) {
		.xval = xval,
		.is_ephemeral = TRUE,
	};
}

static void rxval_evaluator_from_map_literal_free_aux(map_literal_list_evaluator_t* plist_evaluator) {
	for (sllve_t* pe = plist_evaluator->pkvpair_evaluators->phead; pe != NULL; pe = pe->pnext) {
		map_literal_kvpair_evaluator_t* pkvpair_evaluator = pe->pvvalue;
		if (pkvpair_evaluator->pkey_evaluator != NULL) {
			pkvpair_evaluator->pkey_evaluator->pfree_func(pkvpair_evaluator->pkey_evaluator);
		}
		if (pkvpair_evaluator->pxval_evaluator != NULL) {
			pkvpair_evaluator->pxval_evaluator->pfree_func(pkvpair_evaluator->pxval_evaluator);
		}
		if (pkvpair_evaluator->plist_evaluator != NULL) {
			rxval_evaluator_from_map_literal_free_aux(pkvpair_evaluator->plist_evaluator);
		}
		free(pkvpair_evaluator);
	}
	sllv_free(plist_evaluator->pkvpair_evaluators);
	free(plist_evaluator);
}

static void rxval_evaluator_from_map_literal_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_map_literal_state_t* pstate = prxval_evaluator->pvstate;
	rxval_evaluator_from_map_literal_free_aux(pstate->proot_list_evaluator);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_map_literal(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rxval_evaluator_from_map_literal_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_map_literal_state_t));
	pstate->proot_list_evaluator = allocate_map_literal_evaluator_from_ast(
		pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_map_literal_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_map_literal_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_from_nonindexed_local_variable_state_t {
	int vardef_frame_relative_index;
} rxval_evaluator_from_nonindexed_local_variable_state_t;

static boxed_xval_t rxval_evaluator_from_nonindexed_local_variable_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = pvstate;
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	mlhmmv_xvalue_t* pxval = local_stack_frame_ref_extended_from_nonindexed(
		pframe, pstate->vardef_frame_relative_index);
	if (pxval == NULL) {
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = FALSE,
		};
	} else {
		return (boxed_xval_t) {
			.xval = *pxval,
			.is_ephemeral = FALSE,
		};
	}
}

static void rxval_evaluator_from_nonindexed_local_variable_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = prxval_evaluator->pvstate;
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_nonindexed_local_variable(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_nonindexed_local_variable_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_nonindexed_local_variable_state_t));
	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_nonindexed_local_variable_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_nonindexed_local_variable_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_from_indexed_local_variable_state_t {
	int vardef_frame_relative_index;
	sllv_t* pkeylist_evaluators;
} rxval_evaluator_from_indexed_local_variable_state_t;

static boxed_xval_t rxval_evaluator_from_indexed_local_variable_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pkeylist_evaluators, pvars, &all_non_null_or_error);

	if (all_non_null_or_error) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		mlhmmv_xvalue_t* pxval = local_stack_frame_ref_extended_from_indexed(
			pframe, pstate->vardef_frame_relative_index, pmvkeys);
		sllmv_free(pmvkeys);
		if (pxval == NULL) {
			return (boxed_xval_t) {
				.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
				.is_ephemeral = FALSE,
			};
		} else {
			return (boxed_xval_t) {
				.xval = *pxval,
				.is_ephemeral = FALSE,
			};
		}
	} else {
		sllmv_free(pmvkeys);
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = TRUE,
		};
	}
}

static void rxval_evaluator_from_indexed_local_variable_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = prxval_evaluator->pvstate;
	for (sllve_t* pe = pstate->pkeylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* prval_evaluator = pe->pvvalue;
		prval_evaluator->pfree_func(prval_evaluator);
	}
	sllv_free(pstate->pkeylist_evaluators);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_indexed_local_variable(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_indexed_local_variable_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_indexed_local_variable_state_t));
	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;
	pstate->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_indexed_local_variable_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_indexed_local_variable_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_from_indexed_function_call_state_t {
	rxval_evaluator_t* pfunction_call_evaluator;
	sllv_t* pkeylist_evaluators;
} rxval_evaluator_from_indexed_function_call_state_t;

static boxed_xval_t rxval_evaluator_from_indexed_function_call_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_indexed_function_call_state_t* pstate = pvstate;

	// This is for things like $b = splitnvx($a, ":")[1].
	// * First we evaluate the function call -- the splitnvx($a, ":") part.
	// * Second we evaluate the index/indices -- [1] part.
	// * Third we index function-call return value by the indices.

	// Evaluate the function call:
	rxval_evaluator_t* pfunction_call_evaluator = pstate->pfunction_call_evaluator;
	boxed_xval_t function_call_output = pfunction_call_evaluator->pprocess_func(
		pfunction_call_evaluator->pvstate, pvars);

	// Non-indexable if not a map.
	if (function_call_output.xval.is_terminal) {
		mv_free(&function_call_output.xval.terminal_mlrval);
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = TRUE,
		};
	}

	// Evaluate the indices:
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pkeylist_evaluators, pvars, &all_non_null_or_error);

	if (!all_non_null_or_error) {
		mlhmmv_xvalue_free(&function_call_output.xval);
		sllmv_free(pmvkeys);
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = TRUE,
		};
	}

	// Index the function-call return value by the indices:
	int lookup_error_reason_unused = FALSE;
	mlhmmv_xvalue_t* pxval = mlhmmv_level_look_up_and_ref_xvalue(
		function_call_output.xval.pnext_level, pmvkeys, &lookup_error_reason_unused);

	if (pxval == NULL) {
		mlhmmv_xvalue_free(&function_call_output.xval);
		sllmv_free(pmvkeys);
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = FALSE,
		};
	} else {
		// xxx copy out little bit, free full retval
		// mlhmmv_xvalue_free(&function_call_output.xval);
		sllmv_free(pmvkeys);
		return (boxed_xval_t) {
			.xval = *pxval,
			.is_ephemeral = FALSE,
		};
	}
}

static void rxval_evaluator_from_indexed_function_call_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_indexed_function_call_state_t* pstate = prxval_evaluator->pvstate;
	for (sllve_t* pe = pstate->pkeylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* prval_evaluator = pe->pvvalue;
		prval_evaluator->pfree_func(prval_evaluator);
	}
	sllv_free(pstate->pkeylist_evaluators);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_indexed_function_call(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_indexed_function_call_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_indexed_function_call_state_t));

	// Example: 'foo(1,2,3)[4][5]' parses to AST
	//   text="foo", type=INDEXED_FUNCTION_CALLSITE:
	//       text="foo", type=FUNCTION_CALLSITE:
	//           text="1", type=NUMERIC_LITERAL.
	//           text="2", type=NUMERIC_LITERAL.
	//           text="3", type=NUMERIC_LITERAL.
	//       text="indexing", type=MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST:
	//           text="4", type=NUMERIC_LITERAL.
	//           text="5", type=NUMERIC_LITERAL.
	pstate->pfunction_call_evaluator = fmgr_xalloc_provisional_from_operator_or_function_call(
		pfmgr, pnode->pchildren->phead->pvvalue, type_inferencing, context_flags);
	pstate->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pnode->pchildren->phead->pnext->pvvalue, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_indexed_function_call_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_indexed_function_call_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_from_oosvar_keylist_state_t {
	sllv_t* pkeylist_evaluators;
} rxval_evaluator_from_oosvar_keylist_state_t;

static boxed_xval_t rxval_evaluator_from_oosvar_keylist_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_oosvar_keylist_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pkeylist_evaluators, pvars, &all_non_null_or_error);

	if (all_non_null_or_error) {
		int lookup_error = FALSE;
		mlhmmv_xvalue_t* pxval = pmvkeys->phead == NULL
			? &pvars->poosvars->root_xvalue
			: mlhmmv_level_look_up_and_ref_xvalue(pvars->poosvars->root_xvalue.pnext_level, pmvkeys, &lookup_error);
		sllmv_free(pmvkeys);
		if (pxval != NULL) {
			return (boxed_xval_t) {
				.xval = *pxval,
				.is_ephemeral = FALSE,
			};
		} else {
			return (boxed_xval_t) {
				.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
				.is_ephemeral = TRUE,
			};
		}
	} else {
		sllmv_free(pmvkeys);
		return (boxed_xval_t) {
			.xval = mlhmmv_xvalue_wrap_terminal(mv_absent()),
			.is_ephemeral = TRUE,
		};
	}
}

static void rxval_evaluator_from_oosvar_keylist_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_oosvar_keylist_state_t* pstate = prxval_evaluator->pvstate;
	for (sllve_t* pe = pstate->pkeylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* prval_evaluator = pe->pvvalue;
		prval_evaluator->pfree_func(prval_evaluator);
	}
	sllv_free(pstate->pkeylist_evaluators);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_oosvar_keylist(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_oosvar_keylist_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_oosvar_keylist_state_t));
	pstate->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_oosvar_keylist_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_oosvar_keylist_free;

	return prxval_evaluator;
}

// ================================================================
static boxed_xval_t rxval_evaluator_from_full_oosvar_func(void* pvstate, variables_t* pvars) {
	return (boxed_xval_t) {
		.is_ephemeral = FALSE,
		.xval = pvars->poosvars->root_xvalue,
	};
}

static void rxval_evaluator_from_full_oosvar_free(rxval_evaluator_t* prxval_evaluator) {
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_full_oosvar(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pprocess_func = rxval_evaluator_from_full_oosvar_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_full_oosvar_free;

	return prxval_evaluator;
}

// ================================================================
static boxed_xval_t rxval_evaluator_from_full_srec_func(void* pvstate, variables_t* pvars) {
	boxed_xval_t boxed_xval = box_ephemeral_xval(mlhmmv_xvalue_alloc_empty_map());

	for (lrece_t* pe = pvars->pinrec->phead; pe != NULL; pe = pe->pnext) {
		// mlhmmv_level_put_terminal will copy mv keys and values so we needn't (and shouldn't)
		// duplicate them here.
		mv_t k = mv_from_string(pe->key, NO_FREE);
		mv_t* pomv = lhmsmv_get(pvars->ptyped_overlay, pe->key);
		if (pomv != NULL) {
			mlhmmv_level_put_terminal_singly_keyed(boxed_xval.xval.pnext_level, &k, pomv);
		} else {
			mv_t v = mv_from_string(pe->value, NO_FREE); // mlhmmv_level_put_terminal will copy
			mlhmmv_level_put_terminal_singly_keyed(boxed_xval.xval.pnext_level, &k, &v);
		}
	}

	return boxed_xval;
}

static void rxval_evaluator_from_full_srec_free(rxval_evaluator_t* prxval_evaluator) {
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_full_srec(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pprocess_func = rxval_evaluator_from_full_srec_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_full_srec_free;

	return prxval_evaluator;
}

// ================================================================
// xxx code dup w/ function_manager.c

// ================================================================
typedef struct _rxval_evaluator_from_function_callsite_state_t {
	rval_evaluator_t* prval_evaluator;
} rxval_evaluator_from_function_callsite_state_t;

static boxed_xval_t rxval_evaluator_from_function_callsite_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_from_function_callsite_state_t* pstate = pvstate;

	rval_evaluator_t* prval_evaluator = pstate->prval_evaluator;
	mv_t val = prval_evaluator->pprocess_func(prval_evaluator->pvstate, pvars);
	return (boxed_xval_t) {
		.xval = mlhmmv_xvalue_wrap_terminal(val),
		.is_ephemeral = FALSE, // verify reference semantics for RHS evaluators!
	};
}

static void rxval_evaluator_from_function_callsite_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_from_function_callsite_state_t* pstate = prxval_evaluator->pvstate;
	pstate->prval_evaluator->pfree_func(pstate->prval_evaluator);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_from_function_callsite(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	rxval_evaluator_from_function_callsite_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_from_function_callsite_state_t));

	pstate->prval_evaluator = rval_evaluator_alloc_from_ast(pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));

	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_from_function_callsite_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_from_function_callsite_free;

	return prxval_evaluator;
}

// ================================================================
typedef struct _rxval_evaluator_wrapping_rval_state_t {
	rval_evaluator_t* prval_evaluator;
} rxval_evaluator_wrapping_rval_state_t;

static boxed_xval_t rxval_evaluator_wrapping_rval_func(void* pvstate, variables_t* pvars) {
	rxval_evaluator_wrapping_rval_state_t* pstate = pvstate;
	rval_evaluator_t* prval_evaluator = pstate->prval_evaluator;
	mv_t val = prval_evaluator->pprocess_func(prval_evaluator->pvstate, pvars);
	return (boxed_xval_t) {
		.xval = mlhmmv_xvalue_wrap_terminal(val),
		.is_ephemeral = TRUE, // verify reference semantics for RHS evaluators!
	};
}

static void rxval_evaluator_wrapping_rval_free(rxval_evaluator_t* prxval_evaluator) {
	rxval_evaluator_wrapping_rval_state_t* pstate = prxval_evaluator->pvstate;
	pstate->prval_evaluator->pfree_func(pstate->prval_evaluator);
	free(pstate);
	free(prxval_evaluator);
}

rxval_evaluator_t* rxval_evaluator_alloc_wrapping_rval(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rxval_evaluator_wrapping_rval_state_t* pstate = mlr_malloc_or_die(
		sizeof(rxval_evaluator_wrapping_rval_state_t));
	pstate->prval_evaluator = rval_evaluator_alloc_from_ast(pnode, pfmgr, type_inferencing, context_flags);

	rxval_evaluator_t* prxval_evaluator = mlr_malloc_or_die(sizeof(rxval_evaluator_t));
	prxval_evaluator->pvstate       = pstate;
	prxval_evaluator->pprocess_func = rxval_evaluator_wrapping_rval_func;
	prxval_evaluator->pfree_func    = rxval_evaluator_wrapping_rval_free;

	return prxval_evaluator;
}
