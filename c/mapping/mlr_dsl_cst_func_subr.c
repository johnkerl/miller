#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/hss.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

static mv_t cst_udf_process_callback(void* pvstate, int arity, mv_t* args, variables_t* pvars);
static void cst_udf_free_callback(void* pvstate);

// ----------------------------------------------------------------
// $ cat def
//mlr --from s put -v '
//  def f(x,y,z) {
//    local a = 1;
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
//             text="local", type=return:
//                 text="a", type=non_sigil_name.
//                 text="1", type=strnum_literal.
//             text="=", type=srec_assignment:
//                 text="x", type=field_name.
//                 text="2", type=strnum_literal.
//             text="return", type=return:
//                 text="+", type=operator:
//                     text="a", type=bound_variable.
//                     text="*", type=operator:
//                         text="y", type=bound_variable.
//                         text="2", type=strnum_literal.

udf_defsite_state_t* mlr_dsl_cst_alloc_udf(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_ast_node_t* pparameters_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pbody_node = pnode->pchildren->phead->pnext->pvvalue;

	cst_udf_state_t* pcst_udf_state = mlr_malloc_or_die(sizeof(cst_udf_state_t));

	pcst_udf_state->arity = pparameters_node->pchildren->length;
	pcst_udf_state->parameter_names = mlr_malloc_or_die(pcst_udf_state->arity * sizeof(char*));
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

	pcst_udf_state->pbound_variables = lhmsmv_alloc();

	pcst_udf_state->pblock_statements = sllv_alloc();

	for (sllve_t* pe = pbody_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		if (pbody_ast_node->type == MD_AST_NODE_TYPE_RETURN_VOID) {
			fprintf(stderr,
				"%s: return statements within user-defined functions must return a value.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		sllv_append(pcst_udf_state->pblock_statements,
			mlr_dsl_cst_alloc_statement(pbody_ast_node, pcst->pfmgr, pcst->psubroutine_states,
				type_inferencing, context_flags | IN_BINDABLE | IN_FUNC_DEF));
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

	for (int i = 0; i < pstate->arity; i++)
		free(pstate->parameter_names[i]);
	free(pstate->parameter_names);

	lhmsmv_free(pstate->pbound_variables);

	for (sllve_t* pe = pstate->pblock_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		mlr_dsl_cst_statement_free(pstatement);
	}
	sllv_free(pstate->pblock_statements);

	free(pstate);
}

// ----------------------------------------------------------------
// Callback function for the function manager to invoke into here

static mv_t cst_udf_process_callback(void* pvstate, int arity, mv_t* args, variables_t* pvars) {
	cst_udf_state_t* pstate = pvstate;
	mv_t retval = mv_absent();

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Bind parameters to arguments
	bind_stack_push_fenced(pvars->pbind_stack, pstate->pbound_variables);
	for (int i = 0; i < arity; i++) {
		lhmsmv_put(pstate->pbound_variables, pstate->parameter_names[i], &args[i], FREE_ENTRY_VALUE);
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Compute the function value
	cst_outputs_t* pcst_outputs = NULL; // Functions only produce output via their return values

	for (sllve_t* pe = pstate->pblock_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		if (pstatement->preturn_evaluator != NULL) {
			// return statement
			retval = pstatement->preturn_evaluator->pprocess_func(pstatement->preturn_evaluator->pvstate, pvars);
			break;
		} else {
			// anything else
			pstatement->pnode_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	bind_stack_pop(pvars->pbind_stack);

	return retval;
}

// ----------------------------------------------------------------
// Callback function for the function manager to invoke into here

static void cst_udf_free_callback(void* pvstate) {
	cst_udf_state_t* pstate = pvstate;
	mlr_dsl_cst_free_udf(pstate);
}

// ----------------------------------------------------------------
cst_subroutine_state_t* mlr_dsl_cst_alloc_subroutine(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_ast_node_t* pparameters_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pbody_node = pnode->pchildren->phead->pnext->pvvalue;

	int arity = pparameters_node->pchildren->length;
	cst_subroutine_state_t* pstate = mlr_malloc_or_die(sizeof(cst_subroutine_state_t));

	pstate->name = mlr_strdup_or_die(pparameters_node->text);

	pstate->arity = arity;

	pstate->parameter_names = mlr_malloc_or_die(arity * sizeof(char*));
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

	pstate->pbound_variables = lhmsmv_alloc();

	pstate->pblock_statements = sllv_alloc();

	for (sllve_t* pe = pbody_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		if (pbody_ast_node->type == MD_AST_NODE_TYPE_RETURN_VALUE) {
			fprintf(stderr,
				"%s: return statements within user-defined subroutines must not return a value.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		mlr_dsl_cst_statement_t* pstatement = mlr_dsl_cst_alloc_statement(pbody_ast_node,
			pcst->pfmgr, pcst->psubroutine_states, type_inferencing, context_flags | IN_BINDABLE | IN_SUBR_DEF);
		sllv_append(pstate->pblock_statements, pstatement);
	}

	return pstate;
}

void mlr_dsl_cst_free_subroutine(cst_subroutine_state_t* pstate) {
	if (pstate == NULL)
		return;

	free(pstate->name);

	for (int i = 0; i < pstate->arity; i++)
		free(pstate->parameter_names[i]);
	free(pstate->parameter_names);

	lhmsmv_free(pstate->pbound_variables);

	for (sllve_t* pe = pstate->pblock_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		mlr_dsl_cst_statement_free(pstatement);
	}
	sllv_free(pstate->pblock_statements);

	free(pstate);
}

// ----------------------------------------------------------------
void mlr_dsl_cst_execute_subroutine(cst_subroutine_state_t* pstate, variables_t* pvars,
	cst_outputs_t* pcst_outputs, int callsite_arity, mv_t* args)
{
	// Bind parameters to arguments
	bind_stack_push_fenced(pvars->pbind_stack, pstate->pbound_variables);

	for (int i = 0; i < pstate->arity; i++) {
		lhmsmv_put(pstate->pbound_variables, pstate->parameter_names[i], &args[i], FREE_ENTRY_VALUE);
	}

	for (sllve_t* pe = pstate->pblock_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		if (pstatement->is_return_void) {
			// return statement
			break;
		} else {
			// anything else
			pstatement->pnode_handler(pstatement, pvars, pcst_outputs);
			if (loop_stack_get(pvars->ploop_stack) != 0) {
				break;
			}
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	bind_stack_pop(pvars->pbind_stack);
}
