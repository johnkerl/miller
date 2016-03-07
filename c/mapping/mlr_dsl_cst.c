#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing);
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	char*             output_field_name,
	sllv_t*           poosvar_lhs_keylist_evaluators,
	sllv_t*           poosvar_lhs_namelist_evaluators,
	int               all_flag,
	rval_evaluator_t* prhs_evaluator,
	sllv_t*           pcond_statements,
	sllv_t*           poosvar_rhs_keylist_evaluators);

static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem);

static void mlr_dsl_cst_node_evaluate_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_oosvar_to_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_full_srec_from_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_unset(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_unset_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_emitf(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_emit(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_emit_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_dump(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_filter(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_conditional_block(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

static void mlr_dsl_cst_node_evaluate_bare_boolean(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

// ----------------------------------------------------------------
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
		sllv_append(pstatements, pstatement);
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

	pstatement->pitems = sllv_alloc();
	pstatement->pevaluator = NULL;

	if (past->type == MD_AST_NODE_TYPE_SREC_ASSIGNMENT) {
		if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_FIELD_NAME) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else if (pleft->pchildren != NULL) {
			fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			pleft->text,
			NULL,
			NULL,
			FALSE,
			rval_evaluator_alloc_from_ast(pright, type_inferencing),
			NULL,
			NULL));
		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_srec_assignment;

	} else if (past->type == MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT) {
		sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pleft->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		if (pleft->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			sllv_append(poosvar_lhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(mlr_strdup_or_die(pleft->text)));
		} else {
			mlr_dsl_ast_node_t* pnode = pleft;
			while (TRUE) {
				// Example AST:
				// % mlr put -v '@x[1]["2"][$3][@4]=5' /dev/null
				// = (oosvar_assignment):
				//     [] (oosvar_level_key):
				//         [] (oosvar_level_key):
				//             [] (oosvar_level_key):
				//                 [] (oosvar_level_key):
				//                     x (oosvar_name).
				//                     1 (strnum_literal).
				//                 2 (strnum_literal).
				//             3 (field_name).
				//         4 (oosvar_name).
				//     5 (strnum_literal).
				//
				// Here past is the =; pright is the 5; pleft is the string of bracket references
				// ending at the oosvar name.
				//
				// Bracket operators come in from the right. So the highest AST node is the rightmost map, and the
				// lowest is the oosvar name. Hence sllv_prepend rather than sllv_append.
				if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
					mlr_dsl_ast_node_t* pkeynode = pnode->pchildren->phead->pnext->pvvalue;
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
				} else {
					// Oosvar expressions are of the form '@name[$index1][@index2+3][4]["five"].  The first one
					// (name) is special: syntactically, it's outside the brackets, although that issue is for the
					// parser to handle. Here, it's special since it's always a string, never an expression that
					// evaluates to string.
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnode->text)));
				}
				if (pnode->pchildren == NULL)
					break;
				pnode = pnode->pchildren->phead->pvvalue;
			}
		}

		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			poosvar_lhs_keylist_evaluators,
			NULL,
			FALSE,
			rval_evaluator_alloc_from_ast(pright, type_inferencing),
			NULL,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_oosvar_to_oosvar_assignment; // xxx temp

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_oosvar_assignment;

	} else if (past->type == MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT) {
		sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pleft->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
		if (pright->type != MD_AST_NODE_TYPE_FULL_SREC) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		if (pleft->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			sllv_append(poosvar_lhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(mlr_strdup_or_die(pleft->text)));
		} else {
			// xxx make a private helper method out of this
			mlr_dsl_ast_node_t* pnode = pleft;
			while (TRUE) {
				if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
					mlr_dsl_ast_node_t* pkeynode = pnode->pchildren->phead->pnext->pvvalue;
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
				} else {
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnode->text)));
				}
				if (pnode->pchildren == NULL)
					break;
				pnode = pnode->pchildren->phead->pvvalue;
			}
		}

		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			poosvar_lhs_keylist_evaluators,
			NULL,
			FALSE,
			NULL,
			NULL,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_oosvar_from_full_srec_assignment;

	} else if (past->type == MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT) {
		sllv_t* poosvar_rhs_keylist_evaluators = sllv_alloc();

		mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

		if (pleft->type != MD_AST_NODE_TYPE_FULL_SREC) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
		if (pright->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pright->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

		if (pright->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			sllv_append(poosvar_rhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(mlr_strdup_or_die(pright->text)));
		} else {
			// xxx make a private helper method out of this
			mlr_dsl_ast_node_t* pnode = pright;
			while (TRUE) {
				if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
					mlr_dsl_ast_node_t* pkeynode = pnode->pchildren->phead->pnext->pvvalue;
					sllv_prepend(poosvar_rhs_keylist_evaluators,
						rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
				} else {
					sllv_prepend(poosvar_rhs_keylist_evaluators,
						rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnode->text)));
				}
				if (pnode->pchildren == NULL)
					break;
				pnode = pnode->pchildren->phead->pvvalue;
			}
		}

		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			NULL,
			NULL,
			FALSE,
			NULL,
			NULL,
			poosvar_rhs_keylist_evaluators));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_full_srec_from_oosvar_assignment;

	} else if (past->type == MD_AST_NODE_TYPE_UNSET) {

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_unset;
		for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnode = pe->pvvalue;

			if (pnode->type == MD_AST_NODE_TYPE_FIELD_NAME) {
				sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
					pnode->text,
					NULL,
					NULL,
					FALSE,
					NULL,
					NULL,
					NULL));

			} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
				sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();
				sllv_append(poosvar_lhs_keylist_evaluators,
					rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnode->text)));
				sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
					pnode->text,
					poosvar_lhs_keylist_evaluators,
					NULL,
					FALSE,
					NULL,
					NULL,
					NULL));

			// xxx brief cmts here paralleling emit
			} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
				sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();
				mlr_dsl_ast_node_t* pwalker = pnode;
				while (TRUE) {
					if (pwalker->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
						mlr_dsl_ast_node_t* pkeynode = pwalker->pchildren->phead->pnext->pvvalue;
						sllv_prepend(poosvar_lhs_keylist_evaluators,
							rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
					} else {
						sllv_prepend(poosvar_lhs_keylist_evaluators,
							rval_evaluator_alloc_from_string(mlr_strdup_or_die(pwalker->text)));
					}
					if (pwalker->pchildren == NULL)
						break;
					pwalker = pwalker->pchildren->phead->pvvalue;
				}

				sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
					pwalker->text,
					poosvar_lhs_keylist_evaluators,
					NULL,
					FALSE,
					NULL,
					NULL,
					NULL));

			} else if (pnode->type == MD_AST_NODE_TYPE_ALL) {
				sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
					NULL,
					NULL,
					NULL,
					TRUE,
					NULL,
					NULL,
					NULL));
				// The grammar allows only 'unset all', not 'unset @x, all, $y'.
				// So if 'all' appears at all, it's the only name.
				pstatement->pevaluator = mlr_dsl_cst_node_evaluate_unset_all;

			} else {
				fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
					MLR_GLOBALS.argv0, __FILE__, __LINE__);
				exit(1);
			}
		}

	} else if (past->type == MD_AST_NODE_TYPE_EMITF) {
		// Loop over oosvar names to emit in e.g. 'emitf @a, @b, @c'.
		for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnode = pe->pvvalue;
			sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				pnode->text,
				NULL,
				NULL,
				FALSE,
				rval_evaluator_alloc_from_ast(pnode, type_inferencing),
				NULL,
				NULL));
		}

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_emitf;

	} else if (past->type == MD_AST_NODE_TYPE_EMIT) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;

		// The grammar allows only 'emit all', not 'emit @x, all, $y'.
		// So if 'all' appears at all, it's the only name.
		if (pnode->type == MD_AST_NODE_TYPE_ALL) {
			sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				NULL,
				NULL,
				NULL,
				TRUE,
				NULL,
				NULL,
				NULL));

			pstatement->pevaluator = mlr_dsl_cst_node_evaluate_emit_all;

		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			// First argument is oosvar name. Remainings evaluate to string,
			// e.g. 'emit @sums, "color", "shape"'.
			mlr_dsl_ast_node_t* pnamenode = past->pchildren->phead->pvvalue;

			sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();
			sllv_append(poosvar_lhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(mlr_strdup_or_die(pnamenode->text)));

			sllv_t* poosvar_lhs_namelist_evaluators = sllv_alloc();
			for (sllve_t* pe = past->pchildren->phead->pnext; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
				sllv_append(poosvar_lhs_namelist_evaluators,
					rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
			}

			sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				pnamenode->text,
				poosvar_lhs_keylist_evaluators,
				poosvar_lhs_namelist_evaluators,
				FALSE,
				NULL,
				NULL,
				NULL));

			pstatement->pevaluator = mlr_dsl_cst_node_evaluate_emit;

		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {

			// First argument is keyed-oosvar name. Remainings evaluate to string,
			// e.g. 'emit @sums, "color", "shape"'.
			mlr_dsl_ast_node_t* pnamenode = past->pchildren->phead->pvvalue;

			// $ mlr put -q -v 'end{emit @v, "a", "b","c"}' ...
			// AST END STATEMENTS (1):
			// emit (emit):
			//     v (oosvar_name).
			//     a (strnum_literal).
			//     b (strnum_literal).
			//     c (strnum_literal).

			// mlr put -q -v 'end{emit @v[1][2], "a", "b","c"}' ...
			// AST END STATEMENTS (1):
			// emit (emit):
			//     [] (oosvar_level_key):
			//         [] (oosvar_level_key):
			//             v (oosvar_name).
			//             1 (strnum_literal).
			//         2 (strnum_literal).
			//     a (strnum_literal).
			//     b (strnum_literal).
			//     c (strnum_literal).

			// xxx brief cmts here paralleling emit
			sllv_t* poosvar_lhs_keylist_evaluators = sllv_alloc();
			mlr_dsl_ast_node_t* pwalker = pnode;

			while (TRUE) {
				if (pwalker->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
					mlr_dsl_ast_node_t* pkeynode = pwalker->pchildren->phead->pnext->pvvalue;
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
				} else {
					sllv_prepend(poosvar_lhs_keylist_evaluators,
						rval_evaluator_alloc_from_string(mlr_strdup_or_die(pwalker->text)));
				}
				if (pwalker->pchildren == NULL)
					break;
				pwalker = pwalker->pchildren->phead->pvvalue;
			}

			sllv_t* poosvar_lhs_namelist_evaluators = sllv_alloc();

			for (sllve_t* pe = past->pchildren->phead->pnext; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
				sllv_append(poosvar_lhs_keylist_evaluators,
					rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
			}

			sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				pnamenode->text,
				poosvar_lhs_keylist_evaluators,
				poosvar_lhs_namelist_evaluators,
				FALSE,
				NULL,
				NULL,
				NULL));

			pstatement->pevaluator = mlr_dsl_cst_node_evaluate_emit;

		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

	} else if (past->type == MD_AST_NODE_TYPE_CONDITIONAL_BLOCK) {
		// First child node is the AST for the boolean expression. Remaining child nodes are statements
		// to be executed if it evaluates to true.
		mlr_dsl_ast_node_t* pfirst  = past->pchildren->phead->pvvalue;
		sllv_t* pcond_statements = sllv_alloc();

		for (sllve_t* pe = past->pchildren->phead->pnext; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
			mlr_dsl_cst_statement_t *pstatement = cst_statement_alloc(pbody_ast_node, type_inferencing);
			sllv_append(pcond_statements, pstatement);
		}

		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			NULL,
			NULL,
			FALSE,
			rval_evaluator_alloc_from_ast(pfirst, type_inferencing),
			pcond_statements,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_conditional_block;

	} else if (past->type == MD_AST_NODE_TYPE_FILTER) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			NULL,
			NULL,
			FALSE,
			rval_evaluator_alloc_from_ast(pnode, type_inferencing),
			NULL,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_filter;

	} else if (past->type == MD_AST_NODE_TYPE_DUMP) {
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			NULL,
			NULL,
			FALSE,
			NULL,
			NULL,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_dump;

	} else { // Bare-boolean statement
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			NULL,
			NULL,
			NULL,
			FALSE,
			rval_evaluator_alloc_from_ast(past, type_inferencing),
			NULL,
			NULL));

		pstatement->pevaluator = mlr_dsl_cst_node_evaluate_bare_boolean;
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
	char*             output_field_name,
	sllv_t*           poosvar_lhs_keylist_evaluators,
	sllv_t*           poosvar_lhs_namelist_evaluators,
	int               all_flag,
	rval_evaluator_t* prhs_evaluator,
	sllv_t*           pcond_statements,
	sllv_t*           poosvar_rhs_keylist_evaluators)
{
	mlr_dsl_cst_statement_item_t* pitem = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_item_t));
	pitem->output_field_name = output_field_name == NULL ? NULL : mlr_strdup_or_die(output_field_name);
	pitem->poosvar_lhs_keylist_evaluators  = poosvar_lhs_keylist_evaluators;
	pitem->poosvar_lhs_namelist_evaluators = poosvar_lhs_namelist_evaluators;
	pitem->all_flag                        = all_flag;
	pitem->prhs_evaluator                  = prhs_evaluator;
	pitem->pcond_statements                = pcond_statements;
	pitem->poosvar_rhs_keylist_evaluators  = poosvar_rhs_keylist_evaluators;
	return pitem;
}

static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem) {
	if (pitem == NULL)
		return;
	free(pitem->output_field_name);
	if (pitem->prhs_evaluator != NULL)
		pitem->prhs_evaluator->pfree_func(pitem->prhs_evaluator);

	if (pitem->poosvar_lhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pitem->poosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* pevaluator = pe->pvvalue;
			pevaluator->pfree_func(pevaluator);
		}
		sllv_free(pitem->poosvar_lhs_keylist_evaluators);
	}

	if (pitem->poosvar_lhs_namelist_evaluators != NULL) {
		for (sllve_t* pe = pitem->poosvar_lhs_namelist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* pevaluator = pe->pvvalue;
			pevaluator->pfree_func(pevaluator);
		}
		sllv_free(pitem->poosvar_lhs_namelist_evaluators);
	}

	free(pitem);
}

// ----------------------------------------------------------------
void mlr_dsl_cst_evaluate(
	sllv_t*          pcst_statements,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		pstatement->pevaluator(pstatement, poosvars, pinrec, ptyped_overlay, ppregex_captures,
			pctx, pshould_emit_rec, poutrecs);
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
	char* output_field_name = pitem->output_field_name;
	rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	mv_t* pval = mlr_malloc_or_die(sizeof(mv_t));
	*pval = val;

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
	lhmsv_put(ptyped_overlay, output_field_name, pval, NO_FREE);
	lrec_put(pinrec, output_field_name, "bug", NO_FREE);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;

	rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay,
		poosvars, ppregex_captures, pctx, prhs_evaluator->pvstate);

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pitem->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error)
		mlhmmv_put(poosvars, pmvkeys, &rhs_value);
	sllmv_free(pmvkeys);
}

// ----------------------------------------------------------------
// xxx cmt
static void mlr_dsl_cst_node_evaluate_oosvar_to_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;

	int lhs_all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pitem->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &lhs_all_non_null_or_error);

	if (lhs_all_non_null_or_error) {
		int rhs_all_non_null_or_error = TRUE;
		sllmv_t* prhskeys = evaluate_list(pitem->poosvar_rhs_keylist_evaluators,
			pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &rhs_all_non_null_or_error);
		if (rhs_all_non_null_or_error) {
			mlhmmv_copy(poosvars, plhskeys, prhskeys);
		}
		sllmv_free(prhskeys);
	}

	sllmv_free(plhskeys);
}

static void mlr_dsl_cst_node_evaluate_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;

	int all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pitem->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {

		mlhmmv_level_t* plevel = mlhmmv_get_or_create_level(poosvars, plhskeys);
		if (plevel != NULL) {

			for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
				mv_t k = mv_from_string(pe->key, NO_FREE); // mlhmmv_level_put will copy
				sllmve_t e = { .value = k, .free_flags = 0, .pnext = NULL };
				mv_t* pomv = lhmsv_get(ptyped_overlay, pe->key);
				if (pomv != NULL) {
					mlhmmv_level_put(plevel, &e, pomv);
				} else {
					mv_t v = mv_from_string(pe->value, NO_FREE); // mlhmmv_level_put will copy
					mlhmmv_level_put(plevel, &e, &v);
				}
			}

		}
	}
	sllmv_free(plhskeys);
}

static void mlr_dsl_cst_node_evaluate_full_srec_from_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;

	lrec_clear(pinrec);
	// xxx clear the typed overlay too

	int all_non_null_or_error = TRUE;
	sllmv_t* prhskeys = evaluate_list(pitem->poosvar_rhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {
		int error = 0;
		mlhmmv_level_t* plevel = mlhmmv_get_level(poosvars, prhskeys, &error);
		if (plevel != NULL) {
			for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
				if (pentry->level_value.is_terminal) {
					// xxx else flatten!
					char* skey = mv_alloc_format_val(&pentry->level_key);
					mv_t* pval = mv_alloc_copy(&pentry->level_value.u.mlrval);

					// xxx xref to srec_assignment comments re typed_overlay.
					lhmsv_put(ptyped_overlay, mlr_strdup_or_die(skey), pval, FREE_ENTRY_KEY);
					lrec_put(pinrec, skey, "bug", FREE_ENTRY_KEY);
				}
			}
		}
	}
	sllmv_free(prhskeys);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_unset(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	for (sllve_t* pf = pnode->pitems->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_item_t* pitem = pf->pvvalue;
		if (pitem->poosvar_lhs_keylist_evaluators != NULL) {

			int all_non_null_or_error = TRUE;
			sllmv_t* pmvkeys = evaluate_list(pitem->poosvar_lhs_keylist_evaluators,
				pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);

			if (all_non_null_or_error)
				mlhmmv_remove(poosvars, pmvkeys);
			sllmv_free(pmvkeys);
		} else {
			lrec_remove(pinrec, pitem->output_field_name);
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_unset_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	sllmv_t* pempty = sllmv_alloc();
	mlhmmv_remove(poosvars, pempty);
	sllmv_free(pempty);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_emitf(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	lrec_t* prec_to_emit = lrec_unbacked_alloc();
	for (sllve_t* pf = pnode->pitems->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_item_t* pitem = pf->pvvalue;
		char* output_field_name = pitem->output_field_name;
		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

		// This is overkill ... the grammar allows only for oosvar names as args to emit.  So we could bypass
		// that and just hashmap-get keyed by output_field_name here.
		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);

		if (val.type == MT_STRING) {
			// Ownership transfer from (newly created) mlrval to (newly created) lrec.
			lrec_put(prec_to_emit, output_field_name, val.u.strv, val.free_flags);
		} else {
			char free_flags = NO_FREE;
			char* string = mv_format_val(&val, &free_flags);
			lrec_put(prec_to_emit, output_field_name, string, free_flags);
		}

	}
	sllv_append(poutrecs, prec_to_emit);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_emit(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pitem->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	sllmv_t* pmvnames = evaluate_list(pitem->poosvar_lhs_namelist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_to_lrecs(poosvars, pmvkeys, pmvnames, poutrecs);
	}
	sllmv_free(pmvkeys);
	sllmv_free(pmvnames);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_emit_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlhmmv_all_to_lrecs(poosvars, poutrecs);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_dump(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlhmmv_print_json_stacked(poosvars, FALSE);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_filter(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
	rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		*pshould_emit_rec = val.u.boolv;
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_conditional_block(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
	rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		if (val.u.boolv) {
			mlr_dsl_cst_evaluate(pitem->pcond_statements,
				poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs);
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_evaluate_bare_boolean(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs)
{
	mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
	rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val))
		mv_set_boolean_strict(&val);
}
