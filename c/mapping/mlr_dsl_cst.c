#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"

static sllv_t* mlr_dsl_cst_alloc_from_statement_list(sllv_t* pasts, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing);
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static mlr_dsl_cst_statement_item_t* mlr_dsl_cst_statement_item_alloc(
	mlr_dsl_cst_lhs_type_t lhs_type,
	char* output_field_name,
	sllv_t* poosvar_lhs_keylist_evaluators,
	rval_evaluator_t* prhs_evaluator);
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

	pstatement->ast_node_type = past->type;
	pstatement->pitems = sllv_alloc();

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
			MLR_DSL_CST_LHS_TYPE_SREC,
			pleft->text,
			NULL,
			rval_evaluator_alloc_from_ast(pright, type_inferencing)));

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
			MLR_DSL_CST_LHS_TYPE_OOSVAR,
			NULL,
			poosvar_lhs_keylist_evaluators,
			rval_evaluator_alloc_from_ast(pright, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_FILTER) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_NONE,
			NULL,
			NULL,
			rval_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_GATE) {
		mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_NONE,
			NULL,
			NULL,
			rval_evaluator_alloc_from_ast(pnode, type_inferencing)));

	} else if (past->type == MD_AST_NODE_TYPE_EMIT) {
		// Loop over oosvar names to emit in e.g. 'emit @a, @b, @c'.
		for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnode = pe->pvvalue;
			sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
				MLR_DSL_CST_LHS_TYPE_OOSVAR,
				pnode->text,
				NULL,
				rval_evaluator_alloc_from_ast(pnode, type_inferencing)));
		}

	} else if (past->type == MD_AST_NODE_TYPE_DUMP) {
		// No arguments: the node-type alone suffices for the caller to be able to execute this.

	} else { // Bare-boolean statement
		sllv_append(pstatement->pitems, mlr_dsl_cst_statement_item_alloc(
			MLR_DSL_CST_LHS_TYPE_NONE,
			NULL,
			NULL,
			rval_evaluator_alloc_from_ast(past, type_inferencing)));
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
	mlr_dsl_cst_lhs_type_t lhs_type,
	char* output_field_name,
	sllv_t* poosvar_lhs_keylist_evaluators,
	rval_evaluator_t* prhs_evaluator)
{
	mlr_dsl_cst_statement_item_t* pitem = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_item_t));
	pitem->lhs_type = lhs_type;
	pitem->output_field_name = output_field_name == NULL ? NULL : mlr_strdup_or_die(output_field_name);
	pitem->poosvar_lhs_keylist_evaluators = poosvar_lhs_keylist_evaluators;
	pitem->prhs_evaluator = prhs_evaluator;
	return pitem;
}

static void cst_statement_item_free(mlr_dsl_cst_statement_item_t* pitem) {
	if (pitem == NULL)
		return;
	free(pitem->output_field_name);
	pitem->prhs_evaluator->pfree_func(pitem->prhs_evaluator);
	if (pitem->poosvar_lhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pitem->poosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* pevaluator = pe->pvvalue;
			pevaluator->pfree_func(pevaluator);
		}
		sllv_free(pitem->poosvar_lhs_keylist_evaluators);
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
	int*             pemit_rec,
	sllv_t*          poutrecs)
{
	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;

		// xxx temp
		int continuable = mlr_dsl_cst_node_evaluate(pstatement,
			poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pemit_rec, poutrecs);
		if (!continuable)
			break;
	}
}

// ----------------------------------------------------------------
int mlr_dsl_cst_node_evaluate(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pemit_rec,
	sllv_t*          poutrecs)
{

	// Do the evaluations, writing typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).

	mlr_dsl_ast_node_type_t node_type = pnode->ast_node_type;

	if (node_type == MD_AST_NODE_TYPE_SREC_ASSIGNMENT) {
		mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
		char* output_field_name = pitem->output_field_name;
		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		mv_t* pval = mlr_malloc_or_die(sizeof(mv_t));
		*pval = val;

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

	} else if (node_type == MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT) {
		mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;

		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;
		mv_t rhs_value = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay,
			poosvars, ppregex_captures, pctx, prhs_evaluator->pvstate);

		sllmv_t* pmvkeys = sllmv_alloc();
		int keys_ok = TRUE;
		for (sllve_t* pe = pitem->poosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* pmvkey_evaluator = pe->pvvalue;
			mv_t mvkey = pmvkey_evaluator->pprocess_func(pinrec, ptyped_overlay,
				poosvars, ppregex_captures, pctx, pmvkey_evaluator->pvstate);
			if (mv_is_null(&mvkey)) {
				keys_ok = FALSE;
				break;
			}
			// Don't free the mlrval since its memory will be managed by the sllmv.
			sllmv_add(pmvkeys, &mvkey);
		}

		if (keys_ok)
			mlhmmv_put(poosvars, pmvkeys, &rhs_value);

		sllmv_free(pmvkeys);

	} else if (node_type == MD_AST_NODE_TYPE_EMIT) {
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

	} else if (node_type == MD_AST_NODE_TYPE_DUMP) {
		mlhmmv_print_json_stacked(poosvars, FALSE);

	} else if (node_type == MD_AST_NODE_TYPE_FILTER) {
		mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (val.type != MT_NULL) {
			mv_set_boolean_strict(&val);
			*pemit_rec = val.u.boolv;
			if (!val.u.boolv) {
				return FALSE;
			}
		}

	} else if (node_type == MD_AST_NODE_TYPE_GATE) {
		mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (val.type == MT_NULL)
			return FALSE;
		mv_set_boolean_strict(&val);
		if (!val.u.boolv) {
			return FALSE;
		}

	} else if (node_type == MD_AST_NODE_TYPE_CONDITIONAL_BLOCK) {
		// xxx temp

		// xxx evaluate rhs.
		// if true: recurse
		//mlr_dsl_cst_evaluate(pnode->pcond_statements,
			//poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pemit_rec, poutrecs);

	} else { // Bare-boolean statement, or error.
		mlr_dsl_cst_statement_item_t* pitem = pnode->pitems->phead->pvvalue;
		rval_evaluator_t* prhs_evaluator = pitem->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (val.type != MT_NULL)
			mv_set_boolean_strict(&val);
	}

	return TRUE;
}
