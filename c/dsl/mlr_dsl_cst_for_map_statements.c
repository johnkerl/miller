#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

static void handle_for_map_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t*         psubmap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_type_masks,
	int                      prest_for_k_count);

// ================================================================
typedef struct _for_map_state_t {
	char** k_variable_names;
	int*   k_frame_relative_indices;
	int*   k_type_masks;
	int    k_count;

	char*  v_variable_name;
	int    v_frame_relative_index;
	int    v_type_mask;

	rxval_evaluator_t* ptarget_xevaluator;
} for_map_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_map;
static mlr_dsl_cst_statement_freer_t free_for_map;

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

mlr_dsl_cst_statement_t* alloc_for_map(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_map_state_t* pstate = mlr_malloc_or_die(sizeof(for_map_state_t));

	pstate->k_variable_names           = NULL;
	pstate->k_frame_relative_indices   = NULL;
	pstate->k_type_masks               = NULL;
	pstate->k_count                    = 0;

	pstate->v_variable_name            = NULL;
	pstate->v_frame_relative_index     = MD_UNUSED_INDEX;
	pstate->v_type_mask                = 0;

	pstate->ptarget_xevaluator = NULL;

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
		handle_for_map,
		pblock,
		mlr_dsl_cst_handle_statement_block_with_break_continue,
		free_for_map,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_map(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	for_map_state_t* pstate = pstatement->pvstate;

	free(pstate->k_variable_names);
	free(pstate->k_frame_relative_indices);
	free(pstate->k_type_masks);

	pstate->ptarget_xevaluator->pfree_func(pstate->ptarget_xevaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_map(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_map_state_t* pstate = pstatement->pvstate;

	rxval_evaluator_t* ptarget_xevaluator = pstate->ptarget_xevaluator;

	// Evaluate the keylist: e.g. in 'for (k,v in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k is bound within the for-loop.
	boxed_xval_t boxed_xval = ptarget_xevaluator->pprocess_func(ptarget_xevaluator->pvstate, pvars);

	if (!boxed_xval.xval.is_terminal) { // is a map

		// Copy the map for the very likely case that it is being updated inside the for-loop.
		// But ephemerals (map-literals, function return values) aren't named and so can't
		// be modified and so don't need to be copied.
		mlhmmv_xvalue_t* pmap = &boxed_xval.xval;
		mlhmmv_xvalue_t  copy;
		if (!boxed_xval.is_ephemeral) {
			copy = mlhmmv_xvalue_copy(&boxed_xval.xval);
			pmap = &copy;
		}

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
		loop_stack_push(pvars->ploop_stack);

		if (!pmap->is_terminal && pmap->pnext_level != NULL) {
			// Recurse over the for-k-names, e.g. ["k1", "k2"], on each call descending one level
			// deeper into the map.  Note there must be at least one k-name so we are assuming
			// the for-loop within handle_for_map_aux was gone through once & thus
			// mlr_dsl_cst_handle_statement_block_with_break_continue was called through there.

			handle_for_map_aux(pstatement, pvars, pcst_outputs, pmap,
				pstate->k_variable_names, pstate->k_frame_relative_indices,
				pstate->k_type_masks, pstate->k_count);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}
		}

		if (!boxed_xval.is_ephemeral) {
			mlhmmv_xvalue_free(&copy);
		}

		loop_stack_pop(pvars->ploop_stack);
		local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
	}

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(&boxed_xval.xval);
	}
}

static void handle_for_map_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_xvalue_t*         psubmap,
	char**                   prest_for_k_variable_names,
	int*                     prest_for_k_frame_relative_indices,
	int*                     prest_for_k_type_masks,
	int                      prest_for_k_count)
{
	for_map_state_t* pstate = pstatement->pvstate;

	if (prest_for_k_count > 0) { // Keep recursing over remaining k-names

		if (psubmap->is_terminal) {
			// The submap was too shallow for the user-specified k-names; there are no terminals here.
		} else {
			// Loop over keys at this submap level:
			for (mlhmmv_level_entry_t* pe = psubmap->pnext_level->phead; pe != NULL; pe = pe->pnext) {
				// Bind the k-name to the entry-key mlrval:
				local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
				// xxx note copy/ref semantics
				local_stack_frame_define_terminal(pframe, prest_for_k_variable_names[0],
					prest_for_k_frame_relative_indices[0], prest_for_k_type_masks[0],
					mv_copy(&pe->level_key));
				// Recurse into the next-level submap:
				handle_for_map_aux(pstatement, pvars, pcst_outputs, &pe->level_xvalue,
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

		// Bind the v-name to the terminal mlrval:
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		// xxx note copy/ref semantics
		local_stack_frame_define_extended(pframe, pstate->v_variable_name, pstate->v_frame_relative_index,
			pstate->v_type_mask, mlhmmv_xvalue_copy(psubmap));
		// Execute the loop-body statements:
		pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);

	}
}

// ================================================================
typedef struct _for_map_key_only_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	rxval_evaluator_t* ptarget_xevaluator;
} for_map_key_only_state_t;

static mlr_dsl_cst_statement_handler_t handle_for_map_key_only;
static mlr_dsl_cst_statement_freer_t free_for_map_key_only;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_for_map_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_map_key_only_state_t* pstate = mlr_malloc_or_die(sizeof(for_map_key_only_state_t));

	pstate->k_variable_name            = NULL;
	pstate->k_frame_relative_index     = MD_UNUSED_INDEX;
	pstate->k_type_mask                = 0;

	pstate->ptarget_xevaluator         = NULL;

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
		handle_for_map_key_only,
		pblock,
		mlr_dsl_cst_handle_statement_block_with_break_continue,
		free_for_map_key_only,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_map_key_only(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	for_map_key_only_state_t* pstate = pstatement->pvstate;

	pstate->ptarget_xevaluator->pfree_func(pstate->ptarget_xevaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_map_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_map_key_only_state_t* pstate = pstatement->pvstate;
	rxval_evaluator_t* ptarget_xevaluator = pstate->ptarget_xevaluator;

	// Evaluate the keylist: e.g. in 'for (k in @a[$3][x]) { ... }', find the values of $3
	// and x for the current record and stack frame. The keylist bindings are outside the scope
	// of the for-loop, while the k is bound within the for-loop.
	boxed_xval_t boxed_xval = ptarget_xevaluator->pprocess_func(ptarget_xevaluator->pvstate, pvars);

	if (!boxed_xval.xval.is_terminal) { // is a map

		// Copy the map for the very likely case that it is being updated inside the for-loop.
		// But ephemerals (map-literals, function return values) aren't named and so can't
		// be modified and so don't need to be copied.
		mlhmmv_xvalue_t* pmap = &boxed_xval.xval;
		mlhmmv_xvalue_t  copy;
		if (!boxed_xval.is_ephemeral) {
			copy = mlhmmv_xvalue_copy(&boxed_xval.xval);
			pmap = &copy;
		}

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
		loop_stack_push(pvars->ploop_stack);

		sllv_t* pkeys = mlhmmv_xvalue_copy_keys_nonindexed(pmap);

		for (sllve_t* pe = pkeys->phead; pe != NULL; pe = pe->pnext) {
			// Bind the k-name to the current key:
			local_stack_frame_define_terminal(pframe,
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
		sllv_free(pkeys);
		if (!boxed_xval.is_ephemeral) {
			mlhmmv_xvalue_free(&copy);
		}

		loop_stack_pop(pvars->ploop_stack);
		local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
	}

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(&boxed_xval.xval);
	}
}
