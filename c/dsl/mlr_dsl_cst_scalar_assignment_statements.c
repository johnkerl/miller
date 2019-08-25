#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
typedef struct _srec_assignment_state_t {
	char*             srec_lhs_field_name;
	rval_evaluator_t* prhs_evaluator;
} srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_srec_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(srec_assignment_state_t));

	pstate->prhs_evaluator = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->type != MD_AST_NODE_TYPE_FIELD_NAME);
	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->pchildren != NULL);

	pstate->srec_lhs_field_name = plhs_node->text;
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_srec_assignment,
		free_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_srec_assignment(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	srec_assignment_state_t* pstate = pstatement->pvstate;

	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	srec_assignment_state_t* pstate = pstatement->pvstate;

	char* srec_lhs_field_name = pstate->srec_lhs_field_name;

	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	// Write typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	//
	// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the input had
	// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
	// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
	// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
	// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
	// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
	// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
	// here. And putting something statically allocated minimizes copying/freeing.
	if (mv_is_present(&val)) {
		lhmsmv_put(pvars->ptyped_overlay, srec_lhs_field_name, &val, FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, srec_lhs_field_name, "bug", NO_FREE);
	} else {
		mv_free(&val);
	}
}

// ================================================================
typedef struct _indirect_srec_assignment_state_t {
	rval_evaluator_t* plhs_evaluator;
	rval_evaluator_t* prhs_evaluator;
} indirect_srec_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_indirect_srec_assignment;
static mlr_dsl_cst_statement_freer_t free_indirect_srec_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_indirect_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	indirect_srec_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(indirect_srec_assignment_state_t));

	pstate->prhs_evaluator = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	pstate->plhs_evaluator = rval_evaluator_alloc_from_ast(plhs_node,  pcst->pfmgr, type_inferencing, context_flags);
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_indirect_srec_assignment,
		free_indirect_srec_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_indirect_srec_assignment(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	indirect_srec_assignment_state_t* pstate = pstatement->pvstate;

	pstate->plhs_evaluator->pfree_func(pstate->plhs_evaluator);
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_indirect_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	indirect_srec_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* plhs_evaluator = pstate->plhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;

	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	char free_flags;
	char* srec_lhs_field_name = mv_format_val(&lval, &free_flags);

	// Write typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	//
	// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the input had
	// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
	// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
	// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
	// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
	// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
	// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
	// here. And putting something statically allocated minimizes copying/freeing.
	if (mv_is_present(&rval)) {
		lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(srec_lhs_field_name), &rval,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, mlr_strdup_or_die(srec_lhs_field_name), "bug", FREE_ENTRY_KEY | FREE_ENTRY_KEY);
	} else {
		mv_free(&rval);
	}

	if (free_flags) {
		free(srec_lhs_field_name);
	}
}

// ================================================================
typedef struct _positional_srec_name_assignment_state_t {
	rval_evaluator_t* plhs_evaluator;
	rval_evaluator_t* prhs_evaluator;
} positional_srec_name_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_positional_srec_name_assignment;
static mlr_dsl_cst_statement_freer_t free_positional_srec_name_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_positional_srec_name_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	positional_srec_name_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(positional_srec_name_assignment_state_t));

	pstate->prhs_evaluator = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	pstate->plhs_evaluator = rval_evaluator_alloc_from_ast(plhs_node,  pcst->pfmgr, type_inferencing, context_flags);
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_positional_srec_name_assignment,
		free_positional_srec_name_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_positional_srec_name_assignment(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	positional_srec_name_assignment_state_t* pstate = pstatement->pvstate;

	pstate->plhs_evaluator->pfree_func(pstate->plhs_evaluator);
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_positional_srec_name_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	positional_srec_name_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* plhs_evaluator = pstate->plhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;

	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (!mv_is_int(&lval)) {
		char free_flags = NO_FREE;
		char* text = mv_maybe_alloc_format_val(&lval, &free_flags);
		fprintf(stderr, "%s: positional names must be integers; got \"%s\".\n", MLR_GLOBALS.bargv0, text);
		if (free_flags)
			free(text);
		exit(1);
	}

	if (mv_is_absent(&rval)) {
		return;
	} else if (!mv_is_string(&rval)) {
		char free_flags = NO_FREE;
		char* text = mv_maybe_alloc_format_val(&rval, &free_flags);
		fprintf(stderr, "%s: new positional names must be strings; got [%s].\n", MLR_GLOBALS.bargv0, text);
		if (free_flags)
			free(text);
		exit(1);
	}

	int srec_lhs_field_position = lval.u.intv;
	// xxx elsewhere (lrec logic: prohibit empty-string name -- ?
	char* new_name = rval.u.strv;

	// Before: srec is 'a=1,b=2,c=3'
	// Assignment: '$[[3]] = "X"'
	// After:  srec is 'a=1,b=2,X=3'
	if (mv_is_present(&rval)) {
		// xxx fix the lhmsmv_unset of old name on the same separate commit for fixing the unset bug
		char* old_name = lrec_get_key_by_position(pvars->pinrec, srec_lhs_field_position);
		if (old_name != NULL) {
			mv_t* poverlay = lhmsmv_get(pvars->ptyped_overlay, old_name);
			if (poverlay != NULL) {
				mv_t copy = mv_copy(poverlay);
				lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(new_name), &copy,
					FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			}
		}
		lrec_rename_at_position(pvars->pinrec, srec_lhs_field_position, mlr_strdup_or_die(new_name), TRUE);
		mv_free(&rval);
	}
}


// ================================================================
typedef struct _env_assignment_state_t {
	rval_evaluator_t* plhs_evaluator;
	rval_evaluator_t* prhs_evaluator;
} env_assignment_state_t;

static mlr_dsl_cst_statement_handler_t handle_env_assignment;
static mlr_dsl_cst_statement_freer_t free_env_assignment;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_env_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	env_assignment_state_t* pstate = mlr_malloc_or_die(sizeof(env_assignment_state_t));

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));

	mlr_dsl_ast_node_t* plhs_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pnext->pvvalue;

	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->type != MD_AST_NODE_TYPE_ENV);
	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->pchildren == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(plhs_node->pchildren->length != 2);
	mlr_dsl_ast_node_t* pnamenode  = plhs_node->pchildren->phead->pnext->pvvalue;

	pstate->plhs_evaluator = rval_evaluator_alloc_from_ast(pnamenode, pcst->pfmgr, type_inferencing, context_flags);
	pstate->prhs_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr, type_inferencing, context_flags);

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_env_assignment,
		free_env_assignment,
		pstate);
}

// ----------------------------------------------------------------
static void free_env_assignment(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	env_assignment_state_t* pstate = pstatement->pvstate;

	pstate->plhs_evaluator->pfree_func(pstate->plhs_evaluator);
	pstate->prhs_evaluator->pfree_func(pstate->prhs_evaluator);

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_env_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	env_assignment_state_t* pstate = pstatement->pvstate;

	rval_evaluator_t* plhs_evaluator = pstate->plhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstate->prhs_evaluator;
	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&lval) && mv_is_present(&rval)) {
		mlr_arch_setenv(
			mlr_strdup_or_die(mv_alloc_format_val(&lval)),
			mlr_strdup_or_die(mv_alloc_format_val(&rval)));
	}
	mv_free(&lval);
	mv_free(&rval);
}
