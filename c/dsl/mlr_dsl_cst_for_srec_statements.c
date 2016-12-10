#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_freer_t free_for_srec;
static mlr_dsl_cst_statement_handler_t handle_for_srec;

// The variable names are used only for type-decl exceptions. Otherwise the
// names are replaced with frame-relative indices by the stack allocator.
typedef struct _for_srec_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	char* v_variable_name;
	int   v_frame_relative_index;
	int   v_type_mask;

	type_inferenced_srec_field_copy_getter_t* ptype_inferenced_srec_field_copy_getter;

} for_srec_state_t;

// ----------------------------------------------------------------
// $ mlr -n put -v 'for (k,v in $*) { $x=1; $y=2 }'
// AST ROOT:
// text="block", type=STATEMENT_BLOCK:
//     text="for", type=FOR_SREC:
//         text="variables", type=FOR_VARIABLES:
//             text="k", type=UNTYPED_LOCAL_DEFINITION.
//             text="v", type=UNTYPED_LOCAL_DEFINITION.
//         text="for_full_srec_block", type=STATEMENT_BLOCK:
//             text="=", type=SREC_ASSIGNMENT:
//                 text="x", type=FIELD_NAME.
//                 text="1", type=NUMERIC_LITERAL.
//             text="=", type=SREC_ASSIGNMENT:
//                 text="y", type=FIELD_NAME.
//                 text="2", type=NUMERIC_LITERAL.

mlr_dsl_cst_statement_t* alloc_for_srec(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_srec_state_t* pstate = mlr_malloc_or_die(sizeof(for_srec_state_t));

	pstate->k_variable_name        = NULL;
	pstate->k_frame_relative_index = 0;
	pstate->k_type_mask            = TYPE_MASK_ANY;
	pstate->v_variable_name        = NULL;
	pstate->v_frame_relative_index = 0;
	pstate->v_type_mask            = TYPE_MASK_ANY;
	pstate-> ptype_inferenced_srec_field_copy_getter = NULL;

	// Left child node is list of bound variables.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pknode = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvnode = pleft->pchildren->phead->pnext->pvvalue;

	if (streq(pknode->text, pvnode->text)) {
		fprintf(stderr, "%s: duplicate for-loop boundvars \"%s\" and \"%s\".\n",
			MLR_GLOBALS.bargv0, pknode->text, pvnode->text);
		exit(1);
	}

	pstate->k_variable_name = mlr_strdup_or_die(pknode->text);
	pstate->v_variable_name = mlr_strdup_or_die(pvnode->text);
	MLR_INTERNAL_CODING_ERROR_IF(pknode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	MLR_INTERNAL_CODING_ERROR_IF(pvnode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->k_frame_relative_index = pknode->vardef_frame_relative_index;
	pstate->v_frame_relative_index = pvnode->vardef_frame_relative_index;
	pstate->k_type_mask = mlr_dsl_ast_node_type_to_type_mask(pknode->type);
	pstate->v_type_mask = mlr_dsl_ast_node_type_to_type_mask(pvnode->type);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	pstate->ptype_inferenced_srec_field_copy_getter =
		(type_inferencing == TYPE_INFER_STRING_ONLY)      ? get_copy_srec_value_string_only_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT)     ? get_copy_srec_value_string_float_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT_INT) ? get_copy_srec_value_string_float_int_aux :
		NULL;
	MLR_INTERNAL_CODING_ERROR_IF(pstate->ptype_inferenced_srec_field_copy_getter == NULL);

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_srec,
		pblock,
		mlr_dsl_cst_handle_statement_block_with_break_continue,
		free_for_srec,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_srec(mlr_dsl_cst_statement_t* pstatement) {
	for_srec_state_t* pstate = pstatement->pvstate;
	free(pstate->k_variable_name);
	free(pstate->v_variable_name);
	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_srec(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_srec_state_t* pstate = pstatement->pvstate;

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	// Copy the lrec for the very likely case that it is being updated inside the for-loop.
	lrec_t* pcopyrec = lrec_copy(pvars->pinrec);
	lhmsmv_t* pcopyoverlay = lhmsmv_copy(pvars->ptyped_overlay);

	for (lrece_t* pe = pcopyrec->phead; pe != NULL; pe = pe->pnext) {

		mv_t mvkey = mv_from_string_no_free(pe->key);
		mv_t mvval = pstate->ptype_inferenced_srec_field_copy_getter(pe, pcopyoverlay);

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_define_terminal(pframe,
			pstate->k_variable_name, pstate->k_frame_relative_index,
			pstate->k_type_mask, mvkey);
		local_stack_frame_define_terminal(pframe,
			pstate->v_variable_name, pstate->v_frame_relative_index,
			pstate->v_type_mask, mvval);

		pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}
	}
	lhmsmv_free(pcopyoverlay);
	lrec_free(pcopyrec);

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_freer_t free_for_srec_key_only;
static mlr_dsl_cst_statement_handler_t handle_for_srec_key_only;

// The variable names are used only for type-decl exceptions. Otherwise the
// names are replaced with frame-relative indices by the stack allocator.
typedef struct _for_srec_key_only_state_t {
	char* k_variable_name;
	int   k_frame_relative_index;
	int   k_type_mask;

	type_inferenced_srec_field_copy_getter_t* ptype_inferenced_srec_field_copy_getter;

} for_srec_key_only_state_t;

// ----------------------------------------------------------------
// $ mlr -n put -v 'for (k in $*) { $x=1; $y=2 }'
//
// AST ROOT:
// text="block", type=STATEMENT_BLOCK:
//     text="for", type=FOR_SREC_KEY_ONLY:
//         text="variables", type=FOR_VARIABLES:
//             text="k", type=UNTYPED_LOCAL_DEFINITION.
//         text="for_full_srec_block", type=STATEMENT_BLOCK:
//             text="=", type=SREC_ASSIGNMENT:
//                 text="x", type=FIELD_NAME.
//                 text="1", type=NUMERIC_LITERAL.
//             text="=", type=SREC_ASSIGNMENT:
//                 text="y", type=FIELD_NAME.
//                 text="2", type=NUMERIC_LITERAL.

mlr_dsl_cst_statement_t* alloc_for_srec_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	for_srec_key_only_state_t* pstate = mlr_malloc_or_die(sizeof(for_srec_key_only_state_t));

	pstate->k_variable_name        = NULL;
	pstate->k_frame_relative_index = 0;
	pstate->k_type_mask            = TYPE_MASK_ANY;
	pstate-> ptype_inferenced_srec_field_copy_getter = NULL;

	// Left child node is list of bound variables.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pknode = pleft->pchildren->phead->pvvalue;

	pstate->k_variable_name = mlr_strdup_or_die(pknode->text);
	MLR_INTERNAL_CODING_ERROR_IF(pknode->vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->k_frame_relative_index = pknode->vardef_frame_relative_index;
	pstate->k_type_mask = mlr_dsl_ast_node_type_to_type_mask(pknode->type);

	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	cst_statement_block_t* pblock = cst_statement_block_alloc(pnode->subframe_var_count);

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	pstate->ptype_inferenced_srec_field_copy_getter =
		(type_inferencing == TYPE_INFER_STRING_ONLY)      ? get_copy_srec_value_string_only_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT)     ? get_copy_srec_value_string_float_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT_INT) ? get_copy_srec_value_string_float_int_aux :
		NULL;
	MLR_INTERNAL_CODING_ERROR_IF(pstate->ptype_inferenced_srec_field_copy_getter == NULL);

	return mlr_dsl_cst_statement_valloc_with_block(
		pnode,
		handle_for_srec_key_only,
		pblock,
		mlr_dsl_cst_handle_statement_block_with_break_continue,
		free_for_srec_key_only,
		pstate);
}

// ----------------------------------------------------------------
static void free_for_srec_key_only(mlr_dsl_cst_statement_t* pstatement) {
	for_srec_key_only_state_t* pstate = pstatement->pvstate;
	free(pstate->k_variable_name);
	free(pstate);
}

// ----------------------------------------------------------------
static void handle_for_srec_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for_srec_key_only_state_t* pstate = pstatement->pvstate;

	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	local_stack_subframe_enter(pframe, pstatement->pblock->subframe_var_count);
	loop_stack_push(pvars->ploop_stack);

	// Copy the lrec for the very likely case that it is being updated inside the for-loop.
	lrec_t* pcopyrec = lrec_copy(pvars->pinrec);
	lhmsmv_t* pcopyoverlay = lhmsmv_copy(pvars->ptyped_overlay);

	for (lrece_t* pe = pcopyrec->phead; pe != NULL; pe = pe->pnext) {

		mv_t mvkey = mv_from_string_no_free(pe->key);

		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		local_stack_frame_define_terminal(pframe,
			pstate->k_variable_name, pstate->k_frame_relative_index,
			pstate->k_type_mask, mvkey);

		pstatement->pblock_handler(pstatement->pblock, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}
	}
	lhmsmv_free(pcopyoverlay);
	lrec_free(pcopyrec);

	loop_stack_pop(pvars->ploop_stack);
	local_stack_subframe_exit(pframe, pstatement->pblock->subframe_var_count);
}
