#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

static void handle_for_oosvar_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t           submap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_frame_type_masks,
	int                      prest_for_k_count);

static void handle_for_local_map_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t           submap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_frame_type_masks,
	int                      prest_for_k_count);

static void handle_for_map_literal_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t*          psubmap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_frame_type_masks,
	int                      prest_for_k_count);

// ================================================================
typedef struct _for_oosvar_state_t {
	char** k_variable_names;
	int*   k_frame_relative_indices;
	int*   k_type_masks;
	int    k_count;

	char*  v_variable_name;
	int    v_frame_relative_index;
	int    v_type_mask;

    sllv_t* ptarget_keylist_evaluators;
} for_oosvar_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_oosvar;
static mlr_dsl_cst_statement_freer_t free_for_oosvar;

// ----------------------------------------------------------------
// $ mlr -n put -v 'for((k1,k2,k3),v in @a["4"][$5]) { $6 = 7; $8 = 9}'
// AST ROOT:
// text="block", type=STATEMENT_BLOCK:
//     text="for", type=FOR_OOSVAR:
//         text="key_and_value_variables", type=FOR_VARIABLES:
//             text="key_variables", type=FOR_VARIABLES:
//                 text="k1", type=UNTYPED_LOCAL_DEFINITION.
//                 text="k2", type=UNTYPED_LOCAL_DEFINITION.
//                 text="k3", type=UNTYPED_LOCAL_DEFINITION.
//             text="v", type=UNTYPED_LOCAL_DEFINITION.
//         text="oosvar_keylist", type=OOSVAR_KEYLIST:
//             text="a", type=STRING_LITERAL.
//             text="4", type=STRING_LITERAL.
//             text="5", type=FIELD_NAME.
//         text="for_loop_oosvar_block", type=STATEMENT_BLOCK:
//             text="=", type=SREC_ASSIGNMENT:
//                 text="6", type=FIELD_NAME.
//                 text="7", type=NUMERIC_LITERAL.
//             text="=", type=SREC_ASSIGNMENT:
//                 text="8", type=FIELD_NAME.
//                 text="9", type=NUMERIC_LITERAL.

mlr_dsl_cst_statement_t* alloc_for_oosvar(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_oosvar_state_t* pstate = mlr_malloc_or_die(sizeof(for_oosvar_state_t));

	pstate->k_variable_names           = NULL;
	pstate->k_frame_relative_indices   = NULL;
	pstate->k_type_masks               = NULL;
	pstate->k_count                    = 0;

	pstate->v_variable_name            = NULL;
	pstate->v_frame_relative_index     = MD_UNUSED_INDEX;
	pstate->v_type_mask                = 0;

	pstate->ptarget_keylist_evaluators = NULL;

	// Left child node is list of bound variables.
	// - Left subnode is namelist for key boundvars.
	// - Right subnode is name for value boundvar.
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubleft  = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubright = pleft->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	pstate->k_variable_names = mlr_malloc_or_die(sizeof(char*) * psubleft->pchildren->length);
	pstate->k_frame_relative_indices = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_type_masks = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_count = 0;
	for (sllve_t* pe = psubleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnamenode = pe->pvvalue;
		MLR_INTERNAL_CODING_ERROR_IF(pnamenode->vardef_frame_relative_index == MD_UNUSED_INDEX);
		pstate->k_variable_names[pstate->k_count] = pnamenode->text;
		pstate->k_frame_relative_indices[pstate->k_count] = pnamenode->vardef_frame_relative_index;
		pstate->k_type_masks[pstate->k_count] = mlr_dsl_ast_node_type_to_type_mask(pnamenode->type);
		pstate->k_count++;
	}
	pstate->v_variable_name = psubright->text;
	MLR_INTERNAL_CODING_ERROR_IF(psubright->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->v_frame_relative_index = psubright->vardef_frame_relative_index;
	pstate->v_type_mask = mlr_dsl_ast_node_type_to_type_mask(psubright->type);

	pstate->ptarget_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_oosvar,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_oosvar,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_oosvar(mlr_dsl_cst_statement_t* pstatement) {
	for_oosvar_state_t* pstate = pstatement->pvstate;

	free(pstate->k_variable_names);
	free(pstate->k_frame_relative_indices);
	free(pstate->k_type_masks);
	for (sllve_t* pe = pstate->ptarget_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	sllv_free(pstate->ptarget_keylist_evaluators);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_oosvar(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_oosvar_state_t* pstate = pstatement->pvstate;

	// Evaluate the keylist: e.g. in 'for ((k1, k2), v in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k1/k2/v are bound within the for-loop.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* ptarget_keylist = evaluate_list(pstate->ptarget_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
		loop_stack_push(pvars->ploop_stack);

		// Locate and copy the submap indexed by the keylist. E.g. in 'for ((k1, k2), v in @a[3][$4]) { ... }', the
		// submap is indexed by ["a", 3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.
		mlhmmv_xvalue_t submap = mlhmmv_root_copy_xvalue(pvars->poosvars, ptarget_keylist);

		if (!submap.is_terminal && submap.pnext_level != NULL) {
			// Recurse over the for-k-names, e.g. ["k1", "k2"], on each call descending one level
			// deeper into the submap.  Note there must be at least one k-name so we are assuming
			// the for-loop within handle_for_oosvar_aux was gone through once & thus
			// handle_statement_block_with_break_continue was called through there.

			handle_for_oosvar_aux(pstatement, pvars, pcst_outputs, submap,
				pstate->k_variable_names, pstate->k_frame_relative_indices,
				pstate->k_type_masks, pstate->k_count);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}
		}

		mlhmmv_xvalue_free(submap);

		loop_stack_pop(pvars->ploop_stack);
		local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
	}
	sllmv_free(ptarget_keylist);
}

static void handle_for_oosvar_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t           submap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_type_masks,
	int                      prest_for_k_count)
{
	for_oosvar_state_t* pstate = pstatement->pvstate;

	if (prest_for_k_count > 0) { // Keep recursing over remaining k-names

		if (submap.is_terminal) {
			// The submap was too shallow for the user-specified k-names; there are no terminals here.
		} else {
			// Loop over keys at this submap level:
			for (mlhmmv_level_entry_t* pe = submap.pnext_level->phead; pe != NULL; pe = pe->pnext) {
				// Bind the k-name to the entry-key mlrval:
				local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
				local_stack_frame_define_scalar(pframe, prest_for_k_variable_names[0], prest_for_k_frame_relative_indices[0],
					prest_for_k_type_masks[0], mv_copy(&pe->level_key));
				// Recurse into the next-level submap:
				handle_for_oosvar_aux(pstatement, pvars, pcst_outputs, pe->level_value,
					&prest_for_k_variable_names[1], &prest_for_k_frame_relative_indices[1], &prest_for_k_type_masks[1],
					prest_for_k_count - 1);

				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					// Bit cleared in recursive caller
					return;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}

			}
		}

	} else { // End of recursion: k-names have all been used up

		if (!submap.is_terminal) {
			// The submap was too deep for the user-specified k-names; there are no terminals here.
		} else {
			// Bind the v-name to the terminal mlrval:
			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			local_stack_frame_define_scalar(pframe, pstate->v_variable_name, pstate->v_frame_relative_index,
				pstate->v_type_mask, mv_copy(&submap.terminal_mlrval));
			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		}

	}
}

// ================================================================
typedef struct _for_oosvar_key_only_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	char* v_variable_name;
	int   v_frame_relative_index;
	int   v_type_mask;

    sllv_t* ptarget_keylist_evaluators;
} for_oosvar_key_only_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_oosvar_key_only;
static mlr_dsl_cst_statement_freer_t free_for_oosvar_key_only;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_oosvar_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_oosvar_key_only_state_t* pstate = mlr_malloc_or_die(sizeof(for_oosvar_key_only_state_t));

	pstate->k_variable_name            = NULL;
	pstate->k_frame_relative_index     = MD_UNUSED_INDEX;
	pstate->k_type_mask                = 0;

	pstate->v_variable_name            = NULL;
	pstate->v_frame_relative_index     = MD_UNUSED_INDEX;
	pstate->v_type_mask                = 0;

	pstate->ptarget_keylist_evaluators = NULL;

	// Left child node is single bound variable
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->k_variable_name = pleft->text;
	pstate->k_frame_relative_index = pleft->vardef_frame_relative_index;
	pstate->k_type_mask = mlr_dsl_ast_node_type_to_type_mask(pleft->type);

	pstate->ptarget_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_oosvar_key_only,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_oosvar_key_only,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_oosvar_key_only(mlr_dsl_cst_statement_t* pstatement) {
	for_oosvar_key_only_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->ptarget_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	sllv_free(pstate->ptarget_keylist_evaluators);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_oosvar_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_oosvar_key_only_state_t* pstate = pstatement->pvstate;

	// Evaluate the keylist: e.g. in 'for (k in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k is bound within the for-loop.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* ptarget_keylist = evaluate_list(pstate->ptarget_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		// Locate the submap indexed by the keylist and copy its keys. E.g. in 'for (k1 in @a[3][$4]) { ... }', the
		// submap is indexed by ["a", 3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
		loop_stack_push(pvars->ploop_stack);

		sllv_t* pkeys = mlhmmv_root_copy_keys_from_submap(pvars->poosvars, ptarget_keylist);

		for (sllve_t* pe = pkeys->phead; pe != NULL; pe = pe->pnext) {
			// Bind the v-name to the terminal mlrval:
			local_stack_frame_define_scalar(pframe,
				pstate->k_variable_name, pstate->k_frame_relative_index,
				pstate->k_type_mask, mv_copy(pe->pvvalue));

			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}

			mv_free(pe->pvvalue);
			free(pe->pvvalue);
		}

		loop_stack_pop(pvars->ploop_stack);
		local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);

		sllv_free(pkeys);
	}
	sllmv_free(ptarget_keylist);
}

// ================================================================
typedef struct _for_local_map_state_t {
	char** k_variable_names;
	int*   k_frame_relative_indices;
	int*   k_type_masks;
	int    k_count;

	char*  v_variable_name;
	int    v_frame_relative_index;
	int    v_type_mask;

	int    target_frame_relative_index;
    sllv_t* ptarget_keylist_evaluators;
} for_local_map_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_local_map;
static mlr_dsl_cst_statement_freer_t free_for_local_map;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_local_map(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_local_map_state_t* pstate = mlr_malloc_or_die(sizeof(for_local_map_state_t));

	pstate->k_variable_names            = NULL;
	pstate->k_frame_relative_indices    = NULL;
	pstate->k_type_masks                = NULL;
	pstate->k_count                     = 0;

	pstate->v_variable_name             = NULL;
	pstate->v_frame_relative_index      = MD_UNUSED_INDEX;
	pstate->v_type_mask                 = 0;

	pstate->target_frame_relative_index = MD_UNUSED_INDEX;
	pstate->ptarget_keylist_evaluators  = NULL;

	// Left child node is list of bound variables.
	//   Left subnode is namelist for key boundvars.
	//   Right subnode is name for value boundvar.
	// Middle child node is keylist for basepoint in the local mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubleft  = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubright = pleft->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	pstate->k_variable_names = mlr_malloc_or_die(sizeof(char*) * psubleft->pchildren->length);
	pstate->k_frame_relative_indices = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_type_masks = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_count = 0;
	for (sllve_t* pe = psubleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnamenode = pe->pvvalue;
		MLR_INTERNAL_CODING_ERROR_IF(pnamenode->vardef_frame_relative_index == MD_UNUSED_INDEX);
		pstate->k_variable_names[pstate->k_count] = pnamenode->text;
		pstate->k_frame_relative_indices[pstate->k_count] = pnamenode->vardef_frame_relative_index;
		pstate->k_type_masks[pstate->k_count] = mlr_dsl_ast_node_type_to_type_mask(pnamenode->type);
		pstate->k_count++;
	}
	pstate->v_variable_name = psubright->text;
	MLR_INTERNAL_CODING_ERROR_IF(psubright->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->v_frame_relative_index = psubright->vardef_frame_relative_index;
	pstate->v_type_mask = mlr_dsl_ast_node_type_to_type_mask(psubright->type);

	// xxx comment liberally
	MLR_INTERNAL_CODING_ERROR_IF(pmiddle->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->target_frame_relative_index = pmiddle->vardef_frame_relative_index;
	pstate->ptarget_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_local_map,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_local_map,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_local_map(mlr_dsl_cst_statement_t* pstatement) {
	for_local_map_state_t* pstate = pstatement->pvstate;

	free(pstate->k_variable_names);
	free(pstate->k_frame_relative_indices);
	free(pstate->k_type_masks);
	for (sllve_t* pe = pstate->ptarget_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	sllv_free(pstate->ptarget_keylist_evaluators);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_local_map(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_local_map_state_t* pstate = pstatement->pvstate;

	// Evaluate the keylist: e.g. in 'for ((k1, k2), v in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k1/k2/v are bound within the for-loop.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* ptarget_keylist = evaluate_list(pstate->ptarget_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {

		// In '(for a, b in c) { ... }' the 'c' is evaluated in the outer scope and
		// the a, b are bound within the inner scope.
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);

		// Locate and copy the submap indexed by the keylist. E.g. in 'for ((k1, k2), v in a[3][$4])
		// { ... }', the submap is first indexed by the stack-frame slot for "a", then further
		// indexed by [3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.

		mlhmmv_xvalue_t *psubmap = local_stack_frame_get_extended_from_indexed(pframe,
			pstate->target_frame_relative_index, ptarget_keylist);

		if (psubmap != NULL) {
			mlhmmv_xvalue_t submap = mlhmmv_xvalue_copy(psubmap);

			local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
			loop_stack_push(pvars->ploop_stack);

			if (!submap.is_terminal && submap.pnext_level != NULL) {
				// Recurse over the for-k-names, e.g. ["k1", "k2"], on each call descending one level
				// deeper into the submap.  Note there must be at least one k-name so we are assuming
				// the for-loop within handle_for_local_map_aux was gone through once & thus
				// handle_statement_block_with_break_continue was called through there.

				handle_for_local_map_aux(pstatement, pvars, pcst_outputs, submap,
					pstate->k_variable_names, pstate->k_frame_relative_indices,
					pstate->k_type_masks, pstate->k_count);

				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
				}
				if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}
			}

			mlhmmv_xvalue_free(submap);

			loop_stack_pop(pvars->ploop_stack);
			local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
		}
	}
	sllmv_free(ptarget_keylist);
}

static void handle_for_local_map_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t           submap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_type_masks,
	int                      prest_for_k_count)
{
	for_local_map_state_t* pstate = pstatement->pvstate;

	if (prest_for_k_count > 0) { // Keep recursing over remaining k-names

		if (submap.is_terminal) {
			// The submap was too shallow for the user-specified k-names; there are no terminals here.
		} else {
			// Loop over keys at this submap level:
			for (mlhmmv_level_entry_t* pe = submap.pnext_level->phead; pe != NULL; pe = pe->pnext) {
				// Bind the k-name to the entry-key mlrval:
				local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
				local_stack_frame_define_scalar(pframe, prest_for_k_variable_names[0], prest_for_k_frame_relative_indices[0],
					prest_for_k_type_masks[0], mv_copy(&pe->level_key));
				// Recurse into the next-level submap:
				handle_for_local_map_aux(pstatement, pvars, pcst_outputs, pe->level_value,
					&prest_for_k_variable_names[1], &prest_for_k_frame_relative_indices[1], &prest_for_k_type_masks[1],
					prest_for_k_count - 1);

				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					// Bit cleared in recursive caller
					return;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}

			}
		}

	} else { // End of recursion: k-names have all been used up

		if (!submap.is_terminal) {
			// The submap was too deep for the user-specified k-names; there are no terminals here.
		} else {
			// Bind the v-name to the terminal mlrval:
			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			local_stack_frame_define_scalar(pframe, pstate->v_variable_name, pstate->v_frame_relative_index,
				pstate->v_type_mask, mv_copy(&submap.terminal_mlrval));
			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		}

	}
}

// ================================================================
typedef struct _for_local_map_key_only_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	char* v_variable_name;
	int   v_frame_relative_index;
	int   v_type_mask;

	int   target_frame_relative_index;
    sllv_t* ptarget_keylist_evaluators;
} for_local_map_key_only_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_local_map_key_only;
static mlr_dsl_cst_statement_freer_t free_for_local_map_key_only;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_local_map_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_local_map_key_only_state_t* pstate = mlr_malloc_or_die(sizeof(for_local_map_key_only_state_t));

	pstate->k_variable_name             = NULL;
	pstate->k_frame_relative_index      = MD_UNUSED_INDEX;
	pstate->k_type_mask                 = 0;

	pstate->v_variable_name             = NULL;
	pstate->v_frame_relative_index      = MD_UNUSED_INDEX;
	pstate->v_type_mask                 = 0;

	pstate->target_frame_relative_index = MD_UNUSED_INDEX;
	pstate->ptarget_keylist_evaluators  = NULL;

	// Left child node is single bound variable
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->k_variable_name = pleft->text;
	pstate->k_frame_relative_index = pleft->vardef_frame_relative_index;
	pstate->k_type_mask = mlr_dsl_ast_node_type_to_type_mask(pleft->type);

	// xxx comment liberally
	MLR_INTERNAL_CODING_ERROR_IF(pmiddle->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->target_frame_relative_index = pmiddle->vardef_frame_relative_index;
	pstate->ptarget_keylist_evaluators = allocate_keylist_evaluators_from_ast_node( // xxx rename x 2
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_local_map_key_only,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_local_map_key_only,
		pstate);
}


// ----------------------------------------------------------------
static void free_for_local_map_key_only(mlr_dsl_cst_statement_t* pstatement) {
	for_local_map_key_only_state_t* pstate = pstatement->pvstate;

	for (sllve_t* pe = pstate->ptarget_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pev = pe->pvvalue;
		pev->pfree_func(pev);
	}
	sllv_free(pstate->ptarget_keylist_evaluators);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_local_map_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_local_map_key_only_state_t* pstate = pstatement->pvstate;

	// Evaluate the keylist: e.g. in 'for (k in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k is bound within the for-loop.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* ptarget_keylist = evaluate_list(pstate->ptarget_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		// Locate the submap indexed by the keylist and copy its keys. E.g. in 'for (k1 in @a[3][$4]) { ... }', the
		// submap is indexed by ["a", 3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);

		mlhmmv_xvalue_t *psubmap = local_stack_frame_get_extended_from_indexed(pframe,
			pstate->target_frame_relative_index, ptarget_keylist);
		sllv_t* pkeys = mlhmmv_xvalue_copy_keys_indexed(psubmap, NULL); // xxx refactor w/o null

		local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
		loop_stack_push(pvars->ploop_stack);

		for (sllve_t* pe = pkeys->phead; pe != NULL; pe = pe->pnext) {
			// Bind the v-name to the terminal mlrval:
			local_stack_frame_define_scalar(pframe,
				pstate->k_variable_name, pstate->k_frame_relative_index,
				pstate->k_type_mask, mv_copy(pe->pvvalue));

			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}

			mv_free(pe->pvvalue);
			free(pe->pvvalue);
		}

		loop_stack_pop(pvars->ploop_stack);
		local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);

		sllv_free(pkeys);
	}
	sllmv_free(ptarget_keylist);
}

// ================================================================
typedef struct _for_map_literal_state_t {
	char** k_variable_names;
	int*   k_frame_relative_indices;
	int*   k_type_masks;
	int    k_count;

	char*  v_variable_name;
	int    v_frame_relative_index;
	int    v_type_mask;

	rxval_evaluator_t* ptarget_xevaluator;
} for_map_literal_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_map_literal;
static mlr_dsl_cst_statement_freer_t free_for_map_literal;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_map_literal(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_map_literal_state_t* pstate = mlr_malloc_or_die(sizeof(for_map_literal_state_t));

	pstate->k_variable_names            = NULL;
	pstate->k_frame_relative_indices    = NULL;
	pstate->k_type_masks                = NULL;
	pstate->k_count                     = 0;

	pstate->v_variable_name             = NULL;
	pstate->v_frame_relative_index      = MD_UNUSED_INDEX;
	pstate->v_type_mask                 = 0;

	pstate->ptarget_xevaluator          = NULL;

	// Left child node is list of bound variables.
	//   Left subnode is namelist for key boundvars.
	//   Right subnode is name for value boundvar.
	// Middle child node is keylist for basepoint in the local mlhmmv. // xxx update
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubleft  = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubright = pleft->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	pstate->k_variable_names = mlr_malloc_or_die(sizeof(char*) * psubleft->pchildren->length);
	pstate->k_frame_relative_indices = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_type_masks = mlr_malloc_or_die(sizeof(int) * psubleft->pchildren->length);
	pstate->k_count = 0;
	for (sllve_t* pe = psubleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnamenode = pe->pvvalue;
		MLR_INTERNAL_CODING_ERROR_IF(pnamenode->vardef_frame_relative_index == MD_UNUSED_INDEX);
		pstate->k_variable_names[pstate->k_count] = pnamenode->text;
		pstate->k_frame_relative_indices[pstate->k_count] = pnamenode->vardef_frame_relative_index;
		pstate->k_type_masks[pstate->k_count] = mlr_dsl_ast_node_type_to_type_mask(pnamenode->type);
		pstate->k_count++;
	}
	pstate->v_variable_name = psubright->text;
	MLR_INTERNAL_CODING_ERROR_IF(psubright->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->v_frame_relative_index = psubright->vardef_frame_relative_index;
	pstate->v_type_mask = mlr_dsl_ast_node_type_to_type_mask(psubright->type);

	// xxx comment liberally
	pstate->ptarget_xevaluator = rxval_evaluator_alloc_from_ast(
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_map_literal,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_map_literal,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_map_literal(mlr_dsl_cst_statement_t* pstatement) {
	for_map_literal_state_t* pstate = pstatement->pvstate;

	free(pstate->k_variable_names);
	free(pstate->k_frame_relative_indices);
	free(pstate->k_type_masks);
	pstate->ptarget_xevaluator->pfree_func(pstate->ptarget_xevaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_map_literal(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_map_literal_state_t* pstate = pstatement->pvstate;

	// Evaluate the keylist: e.g. in 'for ((k1, k2), v in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k1/k2/v are bound within the for-loop.

	// In '(for a, b in c) { ... }' the 'c' is evaluated in the outer scope and
	// the a, b are bound within the inner scope.
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);

	// Locate and copy the submap indexed by the keylist. E.g. in 'for ((k1, k2), v in a[3][$4])
	// { ... }', the submap is first indexed by the stack-frame slot for "a", then further
	// indexed by [3, $4].  Copy it for the very likely case that it is being updated inside the
	// for-loop.

	boxed_xval_t boxed_xval = pstate->ptarget_xevaluator->pprocess_func(
		pstate->ptarget_xevaluator->pvstate, pvars);

	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	if (!boxed_xval.xval.is_terminal && boxed_xval.xval.pnext_level != NULL) {
		// Recurse over the for-k-names, e.g. ["k1", "k2"], on each call descending one level
		// deeper into the submap.  Note there must be at least one k-name so we are assuming
		// the for-loop within handle_for_map_literal_aux was gone through once & thus
		// handle_statement_block_with_break_continue was called through there.

		handle_for_map_literal_aux(pstatement, pvars, pcst_outputs, &boxed_xval.xval,
			pstate->k_variable_names, pstate->k_frame_relative_indices,
			pstate->k_type_masks, pstate->k_count);

		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
		}
		if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}
	}

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(boxed_xval.xval);
	}

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}

static void handle_for_map_literal_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t*          psubmap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_type_masks,
	int                      prest_for_k_count)
{
	for_map_literal_state_t* pstate = pstatement->pvstate;

	if (prest_for_k_count > 0) { // Keep recursing over remaining k-names

		if (psubmap->is_terminal) {
			// The submap was too shallow for the user-specified k-names; there are no terminals here.
		} else {
			// Loop over keys at this submap level:
			for (mlhmmv_level_entry_t* pe = psubmap->pnext_level->phead; pe != NULL; pe = pe->pnext) {
				// Bind the k-name to the entry-key mlrval:
				local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
				local_stack_frame_define_scalar(pframe, prest_for_k_variable_names[0], prest_for_k_frame_relative_indices[0],
					prest_for_k_type_masks[0], mv_copy(&pe->level_key));
				// Recurse into the next-level submap:
				handle_for_map_literal_aux(pstatement, pvars, pcst_outputs, &pe->level_value,
					&prest_for_k_variable_names[1], &prest_for_k_frame_relative_indices[1], &prest_for_k_type_masks[1],
					prest_for_k_count - 1);

				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					// Bit cleared in recursive caller
					return;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}

			}
		}

	} else { // End of recursion: k-names have all been used up

		if (!psubmap->is_terminal) {
			// The submap was too deep for the user-specified k-names; there are no terminals here.
		} else {
			// Bind the v-name to the terminal mlrval:
			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			local_stack_frame_define_scalar(pframe, pstate->v_variable_name, pstate->v_frame_relative_index,
				pstate->v_type_mask, mv_copy(&psubmap->terminal_mlrval));
			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		}

	}
}

// ================================================================
typedef struct _for_map_literal_key_only_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	char* v_variable_name;
	int   v_frame_relative_index;
	int   v_type_mask;

	rxval_evaluator_t* ptarget_xevaluator;
} for_map_literal_key_only_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_map_literal_key_only;
static mlr_dsl_cst_statement_freer_t free_for_map_literal_key_only;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_map_literal_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_map_literal_key_only_state_t* pstate = mlr_malloc_or_die(sizeof(for_map_literal_key_only_state_t));

	pstate->k_variable_name        = NULL;
	pstate->k_frame_relative_index = MD_UNUSED_INDEX;
	pstate->k_type_mask            = 0;

	pstate->v_variable_name        = NULL;
	pstate->v_frame_relative_index = MD_UNUSED_INDEX;
	pstate->v_type_mask            = 0;

	pstate->ptarget_xevaluator     = NULL;

	// Left child node is single bound variable
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(pleft->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->k_variable_name = pleft->text;
	pstate->k_frame_relative_index = pleft->vardef_frame_relative_index;
	pstate->k_type_mask = mlr_dsl_ast_node_type_to_type_mask(pleft->type);

	pstate->ptarget_xevaluator = rxval_evaluator_alloc_from_ast(
		pmiddle, pcst->pfmgr, type_inferencing, context_flags);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_map_literal_key_only,
		pblock,
		handle_statement_block_with_break_continue,
		free_for_map_literal_key_only,
		pstate);
}


// ----------------------------------------------------------------
static void free_for_map_literal_key_only(mlr_dsl_cst_statement_t* pstatement) {
	for_map_literal_key_only_state_t* pstate = pstatement->pvstate;

	pstate->ptarget_xevaluator->pfree_func(pstate->ptarget_xevaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_map_literal_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_map_literal_key_only_state_t* pstate = pstatement->pvstate;

	boxed_xval_t boxed_xval = pstate->ptarget_xevaluator->pprocess_func(
		pstate->ptarget_xevaluator->pvstate, pvars);

	sllv_t* pkeys = mlhmmv_xvalue_copy_keys_indexed(&boxed_xval.xval, NULL); // xxx refactor w/o null

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	for (sllve_t* pe = pkeys->phead; pe != NULL; pe = pe->pnext) {
		// Bind the v-name to the terminal mlrval:
		local_stack_frame_define_scalar(pframe,
			pstate->k_variable_name, pstate->k_frame_relative_index,
			pstate->k_type_mask, mv_copy(pe->pvvalue));

		// Execute the loop-body statements:
		pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);

		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
		}
		if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}

		mv_free(pe->pvvalue);
		free(pe->pvvalue);
	}

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);

	sllv_free(pkeys);

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(boxed_xval.xval); // xxx rename
	}
}
