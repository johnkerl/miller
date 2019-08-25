#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "keylist_evaluators.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
// Most statements have one item, except emit and unset.
struct _unset_item_t;
typedef void unset_item_handler_t(
	struct _unset_item_t* punset_item,
	variables_t*          pvars,
	cst_outputs_t*        pcst_outputs);

typedef struct _unset_item_t {
	unset_item_handler_t* punset_item_handler;
	int                   local_variable_frame_relative_index;
	char*                 srec_field_name;
	sllv_t*               pkeylist_evaluators;
	rval_evaluator_t*     psrec_field_name_evaluator;
} unset_item_t;

static unset_item_t* alloc_blank_unset_item();
static void free_unset_item(unset_item_t* punset_item);

static void handle_unset_nonindexed_local_variable(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_indexed_local_variable(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_oosvar(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_full_srec(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_indirect_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static void handle_unset_positional_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

// ================================================================
typedef struct _unset_state_t {
	sllv_t* punset_items;
} unset_state_t;

static mlr_dsl_cst_statement_handler_t handle_unset;
static mlr_dsl_cst_statement_freer_t free_unset;

static mlr_dsl_cst_statement_handler_t handle_unset;
static mlr_dsl_cst_statement_handler_t handle_unset_all;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_unset(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	unset_state_t* pstate = mlr_malloc_or_die(sizeof(unset_state_t));

	pstate->punset_items = sllv_alloc();

	mlr_dsl_cst_statement_handler_t* pstatement_handler = handle_unset;
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;

		if (pchild->type == MD_AST_NODE_TYPE_ALL || pchild->type == MD_AST_NODE_TYPE_FULL_OOSVAR) {
			// The grammar allows only 'unset all', not 'unset @x, all, $y'.
			// So if 'all' appears at all, it's the only name. Likewise with '@*'.
			pstatement_handler = handle_unset_all;

		} else if (pchild->type == MD_AST_NODE_TYPE_FULL_SREC) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables is not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_full_srec;
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_FIELD_NAME) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables is not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_srec_field_name;
			punset_item->srec_field_name = pchild->text;
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables are not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_indirect_srec_field_name;
			punset_item->psrec_field_name_evaluator = rval_evaluator_alloc_from_ast(
				pchild->pchildren->phead->pvvalue, pcst->pfmgr, type_inferencing, context_flags);
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables are not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_positional_srec_field_name;
			punset_item->psrec_field_name_evaluator = rval_evaluator_alloc_from_ast(
				pchild->pchildren->phead->pvvalue, pcst->pfmgr, type_inferencing, context_flags);
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_oosvar;
			punset_item->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
				pchild, pcst->pfmgr, type_inferencing, context_flags);
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE) {
			MLR_INTERNAL_CODING_ERROR_IF(pchild->vardef_frame_relative_index == MD_UNUSED_INDEX);
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_nonindexed_local_variable;
			punset_item->local_variable_frame_relative_index = pchild->vardef_frame_relative_index;
			sllv_append(pstate->punset_items, punset_item);

		} else if (pchild->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE) {
			MLR_INTERNAL_CODING_ERROR_IF(pchild->vardef_frame_relative_index == MD_UNUSED_INDEX);
			unset_item_t* punset_item = alloc_blank_unset_item();
			punset_item->punset_item_handler = handle_unset_indexed_local_variable;
			punset_item->local_variable_frame_relative_index = pchild->vardef_frame_relative_index;
			punset_item->pkeylist_evaluators = allocate_keylist_evaluators_from_ast_node(
				pchild, pcst->pfmgr, type_inferencing, context_flags);
			sllv_append(pstate->punset_items, punset_item);

		} else {
			MLR_INTERNAL_CODING_ERROR();
		}
	}
	return mlr_dsl_cst_statement_valloc(
		pnode,
		pstatement_handler,
		free_unset,
		pstate);
}

// ----------------------------------------------------------------
static void free_unset(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	unset_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->punset_items->phead; pe != NULL; pe = pe->pnext) {
		free_unset_item(pe->pvvalue);
	}
	sllv_free(pstate->punset_items);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_unset(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	unset_state_t* pstate = pstatement->pvstate;
	for (sllve_t* pf = pstate->punset_items->phead; pf != NULL; pf = pf->pnext) {
		unset_item_t* punset_item = pf->pvvalue;
		punset_item->punset_item_handler(punset_item, pvars, pcst_outputs);
	}
}

static void handle_unset_all(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	sllmv_t* pempty = sllmv_alloc();
	mlhmmv_root_remove(pvars->poosvars, pempty);
	sllmv_free(pempty);
}

// ----------------------------------------------------------------
static unset_item_t* alloc_blank_unset_item() {
	unset_item_t* punset_item = mlr_malloc_or_die(sizeof(unset_item_t));

	punset_item->punset_item_handler                 = NULL;
	punset_item->local_variable_frame_relative_index = MD_UNUSED_INDEX;
	punset_item->srec_field_name                     = NULL;
	punset_item->pkeylist_evaluators                 = NULL;
	punset_item->psrec_field_name_evaluator          = NULL;

	return punset_item;
}

static void free_unset_item(unset_item_t* punset_item) {
	if (punset_item->pkeylist_evaluators != NULL) {
		for (sllve_t* pe = punset_item->pkeylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(punset_item->pkeylist_evaluators);
	}
	if (punset_item->psrec_field_name_evaluator != NULL) {
		punset_item->psrec_field_name_evaluator->pfree_func(punset_item->psrec_field_name_evaluator);
	}
	free(punset_item);
}

static void handle_unset_nonindexed_local_variable(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_frame_assign_terminal_nonindexed(pframe, punset_item->local_variable_frame_relative_index, mv_absent());
}

// As with oosvars, unset removes the key. E.g. if 'v = { 1:2, 3:4 }' then
// 'unset v[1]' results in 'v = { 3:4 }'.
static void handle_unset_indexed_local_variable(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(punset_item->pkeylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);

		// 'unset nonesuch[someindex]' requires the existence check first: else we'd be poking data
		// into the absent-value stack-frame-index-0 slot.
		mlhmmv_xvalue_t* pxval = local_stack_frame_ref_extended_from_indexed(pframe,
			punset_item->local_variable_frame_relative_index, NULL);
		if (pxval != NULL) {
			mlhmmv_level_remove(pxval->pnext_level, pmvkeys->phead);
		}
	}
	sllmv_free(pmvkeys);
}

static void handle_unset_oosvar(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(punset_item->pkeylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error)
		mlhmmv_root_remove(pvars->poosvars, pmvkeys);
	sllmv_free(pmvkeys);
}

static void handle_unset_full_srec(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	lrec_clear(pvars->pinrec);
}

static void handle_unset_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	lrec_remove(pvars->pinrec, punset_item->srec_field_name);
}

static void handle_unset_indirect_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	rval_evaluator_t* pevaluator = punset_item->psrec_field_name_evaluator;
	mv_t nameval = pevaluator->pprocess_func(pevaluator->pvstate, pvars);
	char free_flags = NO_FREE;
	char* field_name = mv_maybe_alloc_format_val(&nameval, &free_flags);
	lrec_remove(pvars->pinrec, field_name);
	if (free_flags & FREE_ENTRY_VALUE)
		free(field_name);
	mv_free(&nameval);
}

static void handle_unset_positional_srec_field_name(
	unset_item_t*  punset_item,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	rval_evaluator_t* pevaluator = punset_item->psrec_field_name_evaluator;
	mv_t nameval = pevaluator->pprocess_func(pevaluator->pvstate, pvars);
	if (!mv_is_int(&nameval)) {
		char free_flags = NO_FREE;
		char* text = mv_maybe_alloc_format_val(&nameval, &free_flags);
		fprintf(stderr, "%s: positional names must be integers; got \"%s\".\n", MLR_GLOBALS.bargv0, text);
		if (free_flags)
			free(text);
		exit(1);
	}
	// xxx typed overlay too!!
	int field_position = nameval.u.intv;
	lrec_remove_by_position(pvars->pinrec, field_position);
	mv_free(&nameval);
}
