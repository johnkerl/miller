#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "mapping/rval_evaluators.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

static rval_evaluator_t* rval_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	int type_inferencing, function_lookup_t* fcn_lookup_table);

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	return rval_evaluator_alloc_from_ast_aux(pnode, type_inferencing, FUNCTION_LOOKUP_TABLE);
}

static rval_evaluator_t* rval_evaluator_alloc_from_ast_aux(mlr_dsl_ast_node_t* pnode,
	int type_inferencing, function_lookup_t* fcn_lookup_table)
{
	if (pnode->pchildren == NULL) { // leaf node
		if (pnode->type == MD_AST_NODE_TYPE_FIELD_NAME) {
			return rval_evaluator_alloc_from_field_name(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
			return rval_evaluator_alloc_from_oosvar_name(pnode->text);
		} else if (pnode->type == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
			return rval_evaluator_alloc_from_strnum_literal(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_BOOLEAN_LITERAL) {
			return rval_evaluator_alloc_from_boolean_literal(pnode->text);
		} else if (pnode->type == MD_AST_NODE_TYPE_REGEXI) {
			return rval_evaluator_alloc_from_strnum_literal(pnode->text, type_inferencing);
		} else if (pnode->type == MD_AST_NODE_TYPE_CONTEXT_VARIABLE) {
			return rval_evaluator_alloc_from_context_variable(pnode->text);
		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}

	} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		return rval_evaluator_alloc_from_oosvar_level_keys(pnode);

	} else if (pnode->type == MD_AST_NODE_TYPE_ENV) {
		return rval_evaluator_alloc_from_environment(pnode, type_inferencing);

	} else { // operator/function
		if ((pnode->type != MD_AST_NODE_TYPE_FUNCTION_NAME)
		&& (pnode->type != MD_AST_NODE_TYPE_OPERATOR)) {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d (node type %s).\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__, mlr_dsl_ast_node_describe_type(pnode->type));
			exit(1);
		}
		char* func_name = pnode->text;

		int user_provided_arity = pnode->pchildren->length;

		check_arity_with_report(fcn_lookup_table, func_name, user_provided_arity);

		rval_evaluator_t* pevaluator = NULL;
		if (user_provided_arity == 0) {
			pevaluator = rval_evaluator_alloc_from_zary_func_name(func_name);
		} else if (user_provided_arity == 1) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing, fcn_lookup_table);
			pevaluator = rval_evaluator_alloc_from_unary_func_name(func_name, parg1);
		} else if (user_provided_arity == 2) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_regex_arg2_func_name(func_name,
					parg1, parg2_node->text, FALSE);
			} else if ((streq(func_name, "=~") || streq(func_name, "!=~")) && type2 == MD_AST_NODE_TYPE_REGEXI) {
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT);
			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_binary_func_name(func_name, parg1, parg2);
			}

		} else if (user_provided_arity == 3) {
			mlr_dsl_ast_node_t* parg1_node = pnode->pchildren->phead->pvvalue;
			mlr_dsl_ast_node_t* parg2_node = pnode->pchildren->phead->pnext->pvvalue;
			mlr_dsl_ast_node_t* parg3_node = pnode->pchildren->phead->pnext->pnext->pvvalue;
			int type2 = parg2_node->type;

			if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MD_AST_NODE_TYPE_STRNUM_LITERAL) {
				// sub/gsub-regex special case:
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					FALSE, parg3);

			} else if ((streq(func_name, "sub") || streq(func_name, "gsub")) && type2 == MD_AST_NODE_TYPE_REGEXI) {
				// sub/gsub-regex special case:
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_regex_arg2_func_name(func_name, parg1, parg2_node->text,
					TYPE_INFER_STRING_FLOAT_INT, parg3);

			} else {
				// regexes can still be applied here, e.g. if the 2nd argument is a non-terminal AST: however
				// the regexes will be compiled record-by-record rather than once at alloc time, which will
				// be slower.
				rval_evaluator_t* parg1 = rval_evaluator_alloc_from_ast_aux(parg1_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg2 = rval_evaluator_alloc_from_ast_aux(parg2_node, type_inferencing,
					fcn_lookup_table);
				rval_evaluator_t* parg3 = rval_evaluator_alloc_from_ast_aux(parg3_node, type_inferencing,
					fcn_lookup_table);
				pevaluator = rval_evaluator_alloc_from_ternary_func_name(func_name, parg1, parg2, parg3);
			}
		} else {
			fprintf(stderr, "Miller: internal coding error:  arity for function name \"%s\" misdetected.\n",
				func_name);
			exit(1);
		}
		if (pevaluator == NULL) {
			fprintf(stderr, "Miller: unrecognized function name \"%s\".\n", func_name);
			exit(1);
		}
		return pevaluator;
	}
}

// ================================================================
typedef struct _rval_evaluator_field_name_state_t {
	char* field_name;
} rval_evaluator_field_name_state_t;

static mv_t rval_evaluator_field_name_func_string_only(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_empty();
		} else {
			// string points into lrec memory and is valid as long as the lrec is.
			return mv_from_string_no_free(string);
		}
	}
}

static mv_t rval_evaluator_field_name_func_string_float(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_empty();
		} else {
			double fltv;
			if (mlr_try_float_from_string(string, &fltv)) {
				return mv_from_float(fltv);
			} else {
				// string points into lrec memory and is valid as long as the lrec is.
				return mv_from_string_no_free(string);
			}
		}
	}
}

static mv_t rval_evaluator_field_name_func_string_float_int(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate) {
	rval_evaluator_field_name_state_t* pstate = pvstate;
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsv_get(ptyped_overlay, pstate->field_name);
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		return mv_copy(poverlay);
	} else {
		char* string = lrec_get(prec, pstate->field_name);
		if (string == NULL) {
			return mv_absent();
		} else if (*string == 0) {
			return mv_empty();
		} else {
			long long intv;
			double fltv;
			if (mlr_try_int_from_string(string, &intv)) {
				return mv_from_int(intv);
			} else if (mlr_try_float_from_string(string, &fltv)) {
				return mv_from_float(fltv);
			} else {
				// string points into AST memory and is valid as long as the AST is.
				return mv_from_string_no_free(string);
			}
		}
	}
}
static void rval_evaluator_field_name_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_field_name_state_t* pstate = pevaluator->pvstate;
	free(pstate->field_name);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_field_name(char* field_name, int type_inferencing) {
	rval_evaluator_field_name_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_field_name_state_t));
	pstate->field_name = mlr_strdup_or_die(field_name);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	switch (type_inferencing) {
	case TYPE_INFER_STRING_ONLY:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_only;
		break;
	case TYPE_INFER_STRING_FLOAT:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_float;
		break;
	case TYPE_INFER_STRING_FLOAT_INT:
		pevaluator->pprocess_func = rval_evaluator_field_name_func_string_float_int;
		break;
	default:
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
		break;
	}
	pevaluator->pfree_func = rval_evaluator_field_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_oosvar_name_state_t {
	sllmv_t* pmvkeys;
} rval_evaluator_oosvar_name_state_t;

mv_t rval_evaluator_oosvar_name_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_oosvar_name_state_t* pstate = pvstate;
	int error = 0;
	mv_t* pval = mlhmmv_get_terminal(poosvars, pstate->pmvkeys, &error);
	if (pval != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy a value here to feed
		// into that. Otherwise the typed-overlay map in mapper_put would have its contents freed out from underneath it
		// by the evaluator functions.
		if (pval->type == MT_STRING && *pval->u.strv == 0)
			return mv_empty();
		else
			return mv_copy(pval);
	} else {
		return mv_absent();
	}
}

static void rval_evaluator_oosvar_name_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_oosvar_name_state_t* pstate = pevaluator->pvstate;
	sllmv_free(pstate->pmvkeys);
	free(pstate);
	free(pevaluator);
}

// This is used for evaluating @-variables that don't have brackets: e.g. @x vs. @x[$1].
// See comments above rval_evaluator_alloc_from_oosvar_level_keys for more information.
rval_evaluator_t* rval_evaluator_alloc_from_oosvar_name(char* oosvar_name) {
	rval_evaluator_oosvar_name_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_oosvar_name_state_t));
	mv_t mv_name = mv_from_string(oosvar_name, NO_FREE);
	pstate->pmvkeys = sllmv_single_no_free(&mv_name);
	mv_free(&mv_name);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_oosvar_name_func;
	pevaluator->pfree_func = rval_evaluator_oosvar_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_oosvar_level_keys_state_t {
	sllv_t* poosvar_rhs_keylist_evaluators;
} rval_evaluator_oosvar_level_keys_state_t;

mv_t rval_evaluator_oosvar_level_keys_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_oosvar_level_keys_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->poosvar_rhs_keylist_evaluators,
		prec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);

	mv_t rv = mv_absent();
	if (all_non_null_or_error) {
		int error = 0;
		mv_t* pval = mlhmmv_get_terminal(poosvars, pmvkeys, &error);
		if (pval != NULL) {
			if (pval->type == MT_STRING && *pval->u.strv == 0)
				rv = mv_empty();
			else
				rv = *pval;
		}
	}

	sllmv_free(pmvkeys);
	return rv;
}

static void rval_evaluator_oosvar_level_keys_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_oosvar_level_keys_state_t* pstate = pevaluator->pvstate;
	for (sllve_t* pe = pstate->poosvar_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pevaluator = pe->pvvalue;
		pevaluator->pfree_func(pevaluator);
	}
	sllv_free(pstate->poosvar_rhs_keylist_evaluators);
	free(pstate);
	free(pevaluator);
}

// ================================================================
// Example AST:
//
// $ mlr put -v '$y = @x[1]["two"][$3+4][@5]' /dev/null
// = (srec_assignment):
//     y (field_name).
//     [] (oosvar_level_key):
//         [] (oosvar_level_key):
//             [] (oosvar_level_key):
//                 [] (oosvar_level_key):
//                     x (oosvar_name).
//                     1 (strnum_literal).
//                 two (strnum_literal).
//             + (operator):
//                 3 (field_name).
//                 4 (strnum_literal).
//         5 (oosvar_name).
//
// Here past is the =; pright is the 5; pleft is the string of bracket references
// ending at the oosvar name.
//
// The job of this allocator is to set up a linked list of evaluators, with the first position for the oosvar name,
// and the rest for each of the bracketed expressions.  This is used for when there *are* brackets; see
// rval_evaluator_alloc_from_oosvar_name for when there are no brackets.

rval_evaluator_t* rval_evaluator_alloc_from_oosvar_level_keys(mlr_dsl_ast_node_t* past) {
	rval_evaluator_oosvar_level_keys_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_oosvar_level_keys_state_t));

	sllv_t* poosvar_rhs_keylist_evaluators = sllv_alloc();
	mlr_dsl_ast_node_t* pnode = past;
	while (TRUE) {
		// Bracket operators come in from the right. So the highest AST node is the rightmost
		// map, and the lowest is the oosvar name. Hence sllv_prepend rather than sllv_append.
		if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			mlr_dsl_ast_node_t* pkeynode = pnode->pchildren->phead->pnext->pvvalue;
			sllv_prepend(poosvar_rhs_keylist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, TYPE_INFER_STRING_FLOAT_INT));
		} else {
			// Oosvar expressions are of the form '@name[$index1][@index2+3][4]["five"].  The first one (name) is
			// special: syntactically, it's outside the brackets, although that issue is for the parser to handle.
			// Here it's special since it's always a string, never an expression that evaluates to string.
			// Yet for the mlhmmv the first key isn't special.
			sllv_prepend(poosvar_rhs_keylist_evaluators,
				rval_evaluator_alloc_from_string(pnode->text));
		}
		if (pnode->pchildren == NULL)
				break;
		pnode = pnode->pchildren->phead->pvvalue;
	}
	pstate->poosvar_rhs_keylist_evaluators = poosvar_rhs_keylist_evaluators;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_oosvar_level_keys_func;
	pevaluator->pfree_func = rval_evaluator_oosvar_level_keys_free;

	return pevaluator;
}

// ================================================================
// This is used for evaluating strings and numbers in literal expressions, e.g. '$x = "abc"'
// or '$x = "left_\1". The values are subject to replacement with regex captures. See comments
// in mapper_put for more information.
//
// Compare rval_evaluator_alloc_from_string which doesn't do regex replacement: it is intended for
// oosvar names on expression left-hand sides (outside of this file).

typedef struct _rval_evaluator_strnum_literal_state_t {
	mv_t literal;
} rval_evaluator_strnum_literal_state_t;

mv_t rval_evaluator_non_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_strnum_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

mv_t rval_evaluator_string_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_strnum_literal_state_t* pstate = pvstate;
	char* input = pstate->literal.u.strv;

	if (ppregex_captures == NULL || *ppregex_captures == NULL) {
		return mv_from_string_no_free(input);
	} else {
		int was_allocated = FALSE;
		char* output = interpolate_regex_captures(input, *ppregex_captures, &was_allocated);
		if (was_allocated)
			return mv_from_string_with_free(output);
		else
			return mv_from_string_no_free(output);
	}
}
static void rval_evaluator_strnum_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_strnum_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_strnum_literal(char* string, int type_inferencing) {
	rval_evaluator_strnum_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_strnum_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (string == NULL) {
		pstate->literal = mv_absent();
		pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
	} else {
		long long intv;
		double fltv;

		pevaluator->pprocess_func = NULL;
		switch (type_inferencing) {
		case TYPE_INFER_STRING_ONLY:
			pstate->literal = mv_from_string_no_free(string);
			pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			break;

		case TYPE_INFER_STRING_FLOAT:
			if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			}
			break;

		case TYPE_INFER_STRING_FLOAT_INT:
			if (mlr_try_int_from_string(string, &intv)) {
				pstate->literal = mv_from_int(intv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else if (mlr_try_float_from_string(string, &fltv)) {
				pstate->literal = mv_from_float(fltv);
				pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
			} else {
				pstate->literal = mv_from_string_no_free(string);
				pevaluator->pprocess_func = rval_evaluator_string_literal_func;
			}
			break;
		default:
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
			break;
		}
	}
	pevaluator->pfree_func = rval_evaluator_strnum_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// This is intended only for oosvar names on expression left-hand sides (outside of this file).
// Compare rval_evaluator_alloc_from_strnum_literal.

typedef struct _rval_evaluator_string_state_t {
	char* string;
} rval_evaluator_string_state_t;

mv_t rval_evaluator_string_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_string_state_t* pstate = pvstate;
	return mv_from_string_no_free(pstate->string);
}
static void rval_evaluator_string_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_string_state_t* pstate = pevaluator->pvstate;
	free(pstate->string);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_string(char* string) {
	rval_evaluator_string_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_string_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	pstate->string            = mlr_strdup_or_die(string);
	pevaluator->pprocess_func = rval_evaluator_string_func;
	pevaluator->pfree_func    = rval_evaluator_string_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_boolean_literal_state_t {
	mv_t literal;
} rval_evaluator_boolean_literal_state_t;

mv_t rval_evaluator_boolean_literal_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_boolean_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

static void rval_evaluator_boolean_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_boolean_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_boolean_literal(char* string) {
	rval_evaluator_boolean_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_boolean_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (streq(string, "true")) {
		pstate->literal = mv_from_true();
	} else if (streq(string, "false")) {
		pstate->literal = mv_from_false();
	} else {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	pevaluator->pprocess_func = rval_evaluator_boolean_literal_func;
	pevaluator->pfree_func = rval_evaluator_boolean_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// Example:
// $ mlr put -v '$y=ENV["X"]' ...
// AST BEGIN STATEMENTS (0):
// AST MAIN STATEMENTS (1):
// = (srec_assignment):
//     y (field_name).
//     env (env):
//         ENV (env).
//         X (strnum_literal).
// AST END STATEMENTS (0):

// ----------------------------------------------------------------
typedef struct _rval_evaluator_environment_state_t {
	rval_evaluator_t* pname_evaluator;
} rval_evaluator_environment_state_t;

mv_t rval_evaluator_environment_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	rval_evaluator_environment_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(prec, ptyped_overlay,
		poosvars, ppregex_captures, pctx, pstate->pname_evaluator->pvstate);
	if (mv_is_null(&mvname)) {
		return mv_absent();
	}
	char free_flags;
	char* strname = mv_format_val(&mvname, &free_flags);
	char* strvalue = getenv(strname);
	if (strvalue == NULL) {
		mv_free(&mvname);
		if (free_flags & FREE_ENTRY_VALUE)
			free(strname);
		return mv_empty();
	}
	mv_t rv = mv_from_string(strvalue, NO_FREE);
	mv_free(&mvname);
	if (free_flags & FREE_ENTRY_VALUE)
		free(strname);
	return rv;
}

static void rval_evaluator_environment_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_environment_state_t* pstate = pevaluator->pvstate;
	pstate->pname_evaluator->pfree_func(pstate->pname_evaluator);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_environment(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	rval_evaluator_environment_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_environment_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pnext->pvvalue;

	pstate->pname_evaluator = rval_evaluator_alloc_from_ast(pnamenode, type_inferencing);
	pevaluator->pprocess_func = rval_evaluator_environment_func;
	pevaluator->pfree_func = rval_evaluator_environment_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
mv_t rval_evaluator_NF_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(prec->field_count);
}
static void rval_evaluator_NF_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_NF() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_NF_func;
	pevaluator->pfree_func = rval_evaluator_NF_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_NR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->nr);
}
static void rval_evaluator_NR_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_NR() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_NR_func;
	pevaluator->pfree_func = rval_evaluator_NR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FNR_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->fnr);
}
static void rval_evaluator_FNR_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_FNR() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FNR_func;
	pevaluator->pfree_func = rval_evaluator_FNR_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FILENAME_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_string_no_free(pctx->filename);
}
static void rval_evaluator_FILENAME_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_FILENAME() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FILENAME_func;
	pevaluator->pfree_func = rval_evaluator_FILENAME_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_FILENUM_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_int(pctx->filenum);
}
static void rval_evaluator_FILENUM_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_FILENUM() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_FILENUM_func;
	pevaluator->pfree_func = rval_evaluator_FILENUM_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_PI_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_float(M_PI);
}
static void rval_evaluator_PI_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_PI() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_PI_func;
	pevaluator->pfree_func = rval_evaluator_PI_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_E_func(lrec_t* prec, lhmsv_t* ptyped_overlay, mlhmmv_t* poosvars,
	string_array_t** ppregex_captures, context_t* pctx, void* pvstate)
{
	return mv_from_float(M_E);
}
static void rval_evaluator_E_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_E() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_E_func;
	pevaluator->pfree_func = rval_evaluator_E_free;
	return pevaluator;
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_context_variable(char* variable_name) {
	if        (streq(variable_name, "NF"))       { return rval_evaluator_alloc_from_NF();
	} else if (streq(variable_name, "NR"))       { return rval_evaluator_alloc_from_NR();
	} else if (streq(variable_name, "FNR"))      { return rval_evaluator_alloc_from_FNR();
	} else if (streq(variable_name, "FILENAME")) { return rval_evaluator_alloc_from_FILENAME();
	} else if (streq(variable_name, "FILENUM"))  { return rval_evaluator_alloc_from_FILENUM();
	} else if (streq(variable_name, "PI"))       { return rval_evaluator_alloc_from_PI();
	} else if (streq(variable_name, "E"))        { return rval_evaluator_alloc_from_E();
	} else  { return NULL;
	}
}
