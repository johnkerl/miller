#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing);
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	char* output_field_name, int is_oosvar, lrec_evaluator_t* pevaluator);
static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem);

// ----------------------------------------------------------------
// At present (initial oosvar experimens, January 2016) the begin/main/end are organized as follows:
//
// Input:
//
//   mlr put 'begin @a = 1; begin @b = 0; @b = $x + @a; end emit @a, @b'
//
// i.e. there are separate begin keywords one per statement, rather than
//
//   mlr put 'begin { @a = 1; @b = 0}; @b = $x + @a; end { emit @a, @b }'
//
// Corresponding list of per-statement ASTs:
//   begin (begin):
//       = (oosvar_assignment):
//           a (oosvar_name).
//           1 (strnum_literal).
//   begin (begin):
//       = (oosvar_assignment):
//           b (oosvar_name).
//           0 (strnum_literal).
//   = (oosvar_assignment):
//       b (oosvar_name).
//       + (operator):
//           x (field_name).
//           a (oosvar_name).
//   end (end):
//       emit (emit):
//           a (oosvar_name).
//           b (oosvar_name).
//
// (Note that the AST input is a list of per-statement ASTs, rather than a single root-node AST with per-statement child
// nodes.)
//
// So our job here at present is to loop through the per-statement ASTs, splitting them out by begin/main/end.

mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing) {
	mlr_dsl_cst_t* pcst = mlr_malloc_or_die(sizeof(mlr_dsl_cst_t));

	pcst->pbegin_statements = mlr_dsl_cst_alloc_from_statement_list(past->pbegin_statements, type_inferencing);
	pcst->pmain_statements  = mlr_dsl_cst_alloc_from_statement_list(past->pmain_statements,  type_inferencing);
	pcst->pend_statements   = mlr_dsl_cst_alloc_from_statement_list(past->pend_statements,   type_inferencing);

	return pcst;
}

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing) {
	sllv_t* pstatements = sllv_alloc();
	for (sllve_t* pe = pasts->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* past = pe->pvvalue;
		mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc(past, type_inferencing);
		sllv_add(pstatements, pstatement);
	}
	return pstatements;
}

void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst) {
	if (pcst == NULL)
		return;
	for (sllve_t* pe = pcst->pbegin_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pmain_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pend_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	sllv_free(pcst->pbegin_statements);
	sllv_free(pcst->pmain_statements);
	sllv_free(pcst->pend_statements);
	free(pcst);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));

	pstatement->ast_node_type = past->type;
	pstatement->pitems = sllv_alloc();

	if (past->type == MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT) {
		if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MLR_DSL_AST_NODE_TYPE_FIELD_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			pleft->text,
			FALSE,
			lrec_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MLR_DSL_AST_NODE_TYPE_OOSVAR_ASSIGNMENT) {
		if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MLR_DSL_AST_NODE_TYPE_OOSVAR_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			pleft->text,
			TRUE,
			lrec_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MLR_DSL_AST_NODE_TYPE_FILTER) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			TRUE,
			lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MLR_DSL_AST_NODE_TYPE_GATE) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			TRUE,
			lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MLR_DSL_AST_NODE_TYPE_EMIT) {
		// Loop over oosvar names to emit in e.g. 'emit @a, @b, @c'.
		for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnode = pe->pvvalue;
			sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				pnode->text,
				TRUE,
				lrec_evaluator_alloc_from_ast(pnode, type_inferencing)));
		}

	} else { // Bare-boolean statement
		sllv_add(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			TRUE,
			lrec_evaluator_alloc_from_ast(past, type_inferencing)));
	}

	return pstatement;
}

static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement) {
	for (sllve_t* pe = pstatement->pitems->phead; pe != NULL; pe = pe->pnext)
		cst_statement_item_free(pe->pvvalue);
	sllv_free(pstatement->pitems);
	free(pstatement);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	char* output_field_name, int is_oosvar, lrec_evaluator_t* pevaluator)
{
	mlr_dsl_cst_statement_item_t* pitem = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_item_t));
	pitem->output_field_name = output_field_name == NULL ? NULL : mlr_strdup_or_die(output_field_name);
	pitem->pevaluator = pevaluator;
	pitem->is_oosvar = is_oosvar;
	return pitem;
}

static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem) {
	if (pitem == NULL)
		return;
	free(pitem->output_field_name);
	pitem->pevaluator->pfree_func(pitem->pevaluator);
	free(pitem);
}
