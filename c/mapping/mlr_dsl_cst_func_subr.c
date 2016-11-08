#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/hss.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

static mlhmmv_value_t cst_udf_process_callback(void* pvstate, int arity, mv_t* args, variables_t* pvars);
static void cst_udf_free_callback(void* pvstate);

// ----------------------------------------------------------------
// $ cat def
//mlr --from s put -v '
//  def f(x,y,z) {
//    var a = 1;
//    $x = 2;
//    return a + y * 2;
//  }
//'

// $ def
// AST ROOT:
// text="list", type=statement_list:
//     text="f", type=def:
//         text="f", type=non_sigil_name:
//             text="x", type=non_sigil_name.
//             text="y", type=non_sigil_name.
//             text="z", type=non_sigil_name.
//         text="list", type=statement_list:
//             text="var", type=return:
//                 text="a", type=non_sigil_name.
//                 text="1", type=numeric_literal.
//             text="=", type=srec_assignment:
//                 text="x", type=field_name.
//                 text="2", type=numeric_literal.
//             text="return", type=return:
//                 text="+", type=operator:
//                     text="a", type=local_variable.
//                     text="*", type=operator:
//                         text="y", type=local_variable.
//                         text="2", type=numeric_literal.

udf_defsite_state_t* mlr_dsl_cst_alloc_udf(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_ast_node_t* pparameters_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pbody_node = pnode->pchildren->phead->pnext->pvvalue;

	cst_udf_state_t* pcst_udf_state = mlr_malloc_or_die(sizeof(cst_udf_state_t));

	pcst_udf_state->name = mlr_strdup_or_die(pnode->text);
	pcst_udf_state->arity = pparameters_node->pchildren->length;
	pcst_udf_state->parameter_names = mlr_malloc_or_die(pcst_udf_state->arity * sizeof(char*));
	pcst_udf_state->parameter_type_masks = mlr_malloc_or_die(pcst_udf_state->arity * sizeof(int));
	int ok = TRUE;
	hss_t* pnameset = hss_alloc();
	int i = 0;
	for (sllve_t* pe = pparameters_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;

		if (hss_has(pnameset, pparameter_node->text)) {
			fprintf(stderr, "%s: duplicate parameter name \"%s\" in function \"%s\".\n",
				MLR_GLOBALS.bargv0, pparameter_node->text, pnode->text);
			ok = FALSE;
		}
		hss_add(pnameset, pparameter_node->text);

		pcst_udf_state->parameter_names[i] = mlr_strdup_or_die(pparameter_node->text);
		pcst_udf_state->parameter_type_masks[i] = mlr_dsl_ast_node_type_to_type_mask(pparameter_node->type);
	}
	hss_free(pnameset);

	if (!ok) {
		fprintf(stderr, "Parameter names: ");
		for (sllve_t* pe = pparameters_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
			mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
			fprintf(stderr, "\"%s\"", pparameter_node->text);
			if (pe->pnext != NULL)
				fprintf(stderr, ", ");
		}
		fprintf(stderr, ".\n");
		exit(1);
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->max_var_depth == MD_UNUSED_INDEX);
	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	pcst_udf_state->ptop_level_block = cst_top_level_statement_block_alloc(pnode->max_var_depth,
		pnode->subframe_var_count);

	for (sllve_t* pe = pbody_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		if (pbody_ast_node->type == MD_AST_NODE_TYPE_RETURN_VOID) {
			fprintf(stderr,
				"%s: return statements within user-defined functions must return a value.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		sllv_append(pcst_udf_state->ptop_level_block->pblock->pstatements,
			mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node, type_inferencing, context_flags | IN_FUNC_DEF));
	}

	// Callback struct for the function manager to invoke the new function:
	udf_defsite_state_t* pdefsite_state = mlr_malloc_or_die(sizeof(udf_defsite_state_t));
	pdefsite_state->pvstate       = pcst_udf_state;
	pdefsite_state->name          = mlr_strdup_or_die(pnode->text);
	pdefsite_state->arity         = pcst_udf_state->arity;
	pdefsite_state->pprocess_func = cst_udf_process_callback;
	pdefsite_state->pfree_func    = cst_udf_free_callback;

	return pdefsite_state;
}

void mlr_dsl_cst_free_udf(cst_udf_state_t* pstate) {
	if (pstate == NULL)
		return;

	free(pstate->name);
	for (int i = 0; i < pstate->arity; i++)
		free(pstate->parameter_names[i]);
	free(pstate->parameter_names);
	free(pstate->parameter_type_masks);

	cst_top_level_statement_block_free(pstate->ptop_level_block);

	free(pstate);
}

// ----------------------------------------------------------------
// Callback function for the function manager to invoke into here

static mlhmmv_value_t cst_udf_process_callback(void* pvstate, int arity, mv_t* args, variables_t* pvars) {
	cst_udf_state_t* pstate = pvstate;
	cst_top_level_statement_block_t* ptop_level_block = pstate->ptop_level_block;
	mlhmmv_value_t retval = mlhmmv_value_transfer_terminal(mv_absent());

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Push stack and bind parameters to arguments

	local_stack_frame_t* pframe = local_stack_frame_enter(ptop_level_block->pframe);
	local_stack_push(pvars->plocal_stack, pframe);
	local_stack_subframe_enter(pframe, ptop_level_block->pblock->subframe_var_count);

	for (int i = 0; i < arity; i++) {
		// Absent-null is by convention at slot 0 of the frame, and arguments are next.
		// Hence starting the loop at 1.
		local_stack_frame_define(pframe, pstate->parameter_names[i], i+1,
			pstate->parameter_type_masks[i], args[i]); // xxx mapvars
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Compute the function value
	cst_outputs_t* pcst_outputs = NULL; // Functions only produce output via their return values

	if (pvars->trace_execution) {
		fprintf(stderr, "TRACE ENTER FUNC %s\n", pstate->name);
		for (sllve_t* pe = ptop_level_block->pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			fprintf(stderr, "TRACE ");
			mlr_dsl_ast_node_pretty_fprint(pstatement->past_node, stderr);
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			if (pvars->return_state.returned) {
				retval = pvars->return_state.retval; // xxx mapvar
				pvars->return_state.retval = mlhmmv_value_transfer_terminal(mv_absent());
				pvars->return_state.returned = FALSE;
				break;
			}
		}
		fprintf(stderr, "TRACE EXIT FUNC %s\n", pstate->name);
	} else {
		for (sllve_t* pe = ptop_level_block->pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			if (pvars->return_state.returned) {
				retval = pvars->return_state.retval; // xxx mapvar
				pvars->return_state.retval = mlhmmv_value_transfer_terminal(mv_absent());
				pvars->return_state.returned = FALSE;
				break;
			}
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Pop stack
	local_stack_subframe_exit(pframe, ptop_level_block->pblock->subframe_var_count);
	local_stack_frame_exit(local_stack_pop(pvars->plocal_stack));

	return retval;
}

// ----------------------------------------------------------------
// Callback function for the function manager to invoke into here

static void cst_udf_free_callback(void* pvstate) {
	cst_udf_state_t* pstate = pvstate;
	mlr_dsl_cst_free_udf(pstate);
}

// ================================================================
typedef struct _subr_callsite_statement_state_t {
	rval_evaluator_t** subr_callsite_argument_evaluators; // xxx mapvar
	mv_t*              subr_callsite_arguments;           // xxx mapvar
	subr_callsite_t    *psubr_callsite;
	subr_defsite_t     *psubr_defsite;

} subr_callsite_statement_state_t;

static mlr_dsl_cst_statement_handler_t handle_subr_callsite_statement;
static mlr_dsl_cst_statement_freer_t free_subr_callsite_statement;

static subr_callsite_t* subr_callsite_alloc(char* name, int arity, int type_inferencing, int context_flags);
static void subr_callsite_free(subr_callsite_t* psubr_callsite);

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_subr_callsite_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	subr_callsite_statement_state_t* pstate = mlr_malloc_or_die(sizeof(subr_callsite_statement_state_t));

	pstate->subr_callsite_argument_evaluators = NULL;
	pstate->subr_callsite_arguments           = NULL;
	pstate->psubr_callsite                    = NULL;
	pstate->psubr_defsite                     = NULL;

	mlr_dsl_ast_node_t* pname_node = pnode->pchildren->phead->pvvalue;
	int callsite_arity = pname_node->pchildren->length;

	pstate->psubr_callsite = subr_callsite_alloc(pname_node->text, callsite_arity,
		type_inferencing, context_flags);

	pstate->subr_callsite_argument_evaluators = mlr_malloc_or_die(callsite_arity * sizeof(rval_evaluator_t*));
	pstate->subr_callsite_arguments = mlr_malloc_or_die(callsite_arity * sizeof(mv_t));

	int i = 0;
	for (sllve_t* pe = pname_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pargument_node = pe->pvvalue;
		pstate->subr_callsite_argument_evaluators[i] = rval_evaluator_alloc_from_ast(pargument_node,
			pcst->pfmgr, type_inferencing, context_flags);
	}

	mlr_dsl_cst_statement_t* pstatement = mlr_dsl_cst_statement_valloc(
		pnode,
		handle_subr_callsite_statement,
		free_subr_callsite_statement,
		pstate);

	// Remember this callsite to be resolved later, after all subroutine definitions have been done.
	sllv_append(pcst->psubr_callsite_statements_to_resolve, pstatement);

	return pstatement;
}

// ----------------------------------------------------------------
void mlr_dsl_cst_resolve_subr_callsite(mlr_dsl_cst_t* pcst, mlr_dsl_cst_statement_t* pstatement) {
	subr_callsite_statement_state_t* pstate = pstatement->pvstate;

	subr_callsite_t* psubr_callsite = pstate->psubr_callsite;
	subr_defsite_t* psubr_defsite = lhmsv_get(pcst->psubr_defsites, psubr_callsite->name);
	if (psubr_defsite == NULL) {
		fprintf(stderr, "%s: subroutine \"%s\" not found.\n", MLR_GLOBALS.bargv0, psubr_callsite->name);
		exit(1);
	}
	if (psubr_defsite->arity != psubr_callsite->arity) {
		fprintf(stderr, "%s: subroutine \"%s\" expects argument count %d but argument count %d was provided.\n",
			MLR_GLOBALS.bargv0, psubr_callsite->name, psubr_defsite->arity, psubr_callsite->arity);
		exit(1);
	}
	pstate->psubr_defsite = psubr_defsite;
}

// ----------------------------------------------------------------
static void handle_subr_callsite_statement( // XXX
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	subr_callsite_statement_state_t* pstate = pstatement->pvstate;

	for (int i = 0; i < pstate->psubr_callsite->arity; i++) {
		rval_evaluator_t* pev = pstate->subr_callsite_argument_evaluators[i];
		pstate->subr_callsite_arguments[i] = pev->pprocess_func(pev->pvstate, pvars); // xxx mapvars
	}

	mlr_dsl_cst_execute_subroutine(pstate->psubr_defsite, pvars, pcst_outputs,
		pstate->psubr_callsite->arity, pstate->subr_callsite_arguments);
}

// ----------------------------------------------------------------
static void free_subr_callsite_statement(mlr_dsl_cst_statement_t* pstatement) { // subr_callsite_statement
	subr_callsite_statement_state_t* pstate = pstatement->pvstate;

	// xxx pre-federation
	if (pstate->subr_callsite_argument_evaluators != NULL) {
		for (int i = 0; i < pstate->psubr_callsite->arity; i++) {
			rval_evaluator_t* phandler = pstate->subr_callsite_argument_evaluators[i];
			phandler->pfree_func(phandler);
		}
		free(pstate->subr_callsite_argument_evaluators);
	}

	if (pstate->subr_callsite_arguments != NULL) {
		// mv_frees are done by the local-stack container which owns the mlrvals it contains
		free(pstate->subr_callsite_arguments);
	}
	subr_callsite_free(pstate->psubr_callsite);

	free(pstate);
}

// ----------------------------------------------------------------
static subr_callsite_t* subr_callsite_alloc(char* name, int arity, int type_inferencing, int context_flags) {
	subr_callsite_t* psubr_callsite  = mlr_malloc_or_die(sizeof(subr_callsite_t));
	psubr_callsite->name             = mlr_strdup_or_die(name);
	psubr_callsite->arity            = arity;
	psubr_callsite->type_inferencing = type_inferencing;
	psubr_callsite->context_flags    = context_flags;
	return psubr_callsite;
}

static void subr_callsite_free(subr_callsite_t* psubr_callsite) {
	if (psubr_callsite == NULL)
		return;
	free(psubr_callsite->name);
	free(psubr_callsite);
}

// ================================================================
subr_defsite_t* mlr_dsl_cst_alloc_subroutine(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode, // xxx rename?
	int type_inferencing, int context_flags)
{
	mlr_dsl_ast_node_t* pparameters_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pbody_node = pnode->pchildren->phead->pnext->pvvalue;

	int arity = pparameters_node->pchildren->length;
	subr_defsite_t* pstate = mlr_malloc_or_die(sizeof(subr_defsite_t));

	pstate->name = mlr_strdup_or_die(pparameters_node->text);

	pstate->arity = arity;

	pstate->parameter_names = mlr_malloc_or_die(arity * sizeof(char*));
	pstate->parameter_type_masks = mlr_malloc_or_die(arity * sizeof(int));
	int ok = TRUE;
	hss_t* pnameset = hss_alloc();
	int i = 0;
	for (sllve_t* pe = pparameters_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;

		if (hss_has(pnameset, pparameter_node->text)) {
			fprintf(stderr, "%s: duplicate parameter name \"%s\" in subroutine \"%s\".\n",
				MLR_GLOBALS.bargv0, pparameter_node->text, pnode->text);
			ok = FALSE;
		}
		hss_add(pnameset, pparameter_node->text);

		pstate->parameter_names[i] = mlr_strdup_or_die(pparameter_node->text);
		pstate->parameter_type_masks[i] = mlr_dsl_ast_node_type_to_type_mask(pparameter_node->type);
	}
	hss_free(pnameset);

	if (!ok) {
		fprintf(stderr, "Parameter names: ");
		for (sllve_t* pe = pparameters_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
			mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
			fprintf(stderr, "\"%s\"", pparameter_node->text);
			if (pe->pnext != NULL)
				fprintf(stderr, ", ");
		}
		fprintf(stderr, ".\n");
		exit(1);
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->max_var_depth == MD_UNUSED_INDEX);
	MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
	pstate->ptop_level_block = cst_top_level_statement_block_alloc(pnode->max_var_depth, pnode->subframe_var_count);

	for (sllve_t* pe = pbody_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		if (pbody_ast_node->type == MD_AST_NODE_TYPE_RETURN_VALUE) {
			fprintf(stderr,
				"%s: return statements within user-defined subroutines must not return a value.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		mlr_dsl_cst_statement_t* pstatement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags | IN_SUBR_DEF);
		sllv_append(pstate->ptop_level_block->pblock->pstatements, pstatement);
	}

	return pstate;
}

// ----------------------------------------------------------------
void mlr_dsl_cst_free_subroutine(subr_defsite_t* pstate) {
	if (pstate == NULL)
		return;

	free(pstate->name);

	for (int i = 0; i < pstate->arity; i++)
		free(pstate->parameter_names[i]);
	free(pstate->parameter_names);
	free(pstate->parameter_type_masks);

	cst_top_level_statement_block_free(pstate->ptop_level_block);

	free(pstate);
}

// ----------------------------------------------------------------
void mlr_dsl_cst_execute_subroutine(subr_defsite_t* pstate, variables_t* pvars, // xxx mv_t -> mlhmmv_value_t
	cst_outputs_t* pcst_outputs, int callsite_arity, mv_t* args)
{
	cst_top_level_statement_block_t* ptop_level_block = pstate->ptop_level_block;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Push stack and bind parameters to arguments
	local_stack_frame_t* pframe = local_stack_frame_enter(ptop_level_block->pframe);
	local_stack_push(pvars->plocal_stack, pframe);
	local_stack_subframe_enter(pframe, ptop_level_block->pblock->subframe_var_count);

	for (int i = 0; i < pstate->arity; i++) {
		// Absent-null is by convention at slot 0 of the frame, and arguments are next.
		// Hence starting the loop at 1.
		local_stack_frame_define(pframe, pstate->parameter_names[i], i+1,
			pstate->parameter_type_masks[i], args[i]);
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Execute the subroutine body

	if (pvars->trace_execution) {
		fprintf(stderr, "TRACE ENTER SUBR %s\n", pstate->name);
		for (sllve_t* pe = pstate->ptop_level_block->pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			fprintf(stderr, "TRACE ");
			mlr_dsl_ast_node_pretty_fprint(pstatement->past_node, stderr);
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			if (pvars->return_state.returned) {
				pvars->return_state.returned = FALSE;
				break;
			}
		}
		fprintf(stderr, "TRACE EXIT SUBR %s\n", pstate->name);
	} else {
		for (sllve_t* pe = pstate->ptop_level_block->pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			pstatement->pstatement_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
			if (pvars->return_state.returned) {
				pvars->return_state.returned = FALSE;
				break;
			}
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Pop stack
	local_stack_subframe_exit(pframe, ptop_level_block->pblock->subframe_var_count);
	local_stack_frame_exit(local_stack_pop(pvars->plocal_stack));
}
