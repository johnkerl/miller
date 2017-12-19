#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <ctype.h> // for tolower(), toupper()
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mtrand.h"
#include "mapping/mapper.h"
#include "dsl/rval_evaluators.h"
#include "dsl/function_manager.h"
#include "dsl/context_flags.h"

// ================================================================
// See comments in rval_evaluators.h
// ================================================================

// ================================================================
// The grammar permits certain statements which are syntactically invalid, (a) because it's awkward to handle
// there, and (b) because we get far better control over error messages here (vs. 'syntax error').
// The context flags are used as the CST is built from the AST, for CST-build-time validation.
// This semantic analysis isn't a separate pass through the AST or CST since it's done while the
// CST is being constructed.

rval_evaluator_t* rval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	if (pnode->pchildren == NULL) {
		// leaf node
		switch(pnode->type) {

		case MD_AST_NODE_TYPE_FIELD_NAME:
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			return rval_evaluator_alloc_from_field_name(pnode->text, type_inferencing);
			break;

		case MD_AST_NODE_TYPE_STRING_LITERAL:
			// In input data such as echo x=3,y=4 | mlr put '$z=$x+$y', the 3 and 4 are strings
			// which need parsing as integers. But in DSL expression literals such as 'put $z = "3" + 4'
			// the "3" should not.
			return rval_evaluator_alloc_from_string_literal(pnode->text);
			break;

		case MD_AST_NODE_TYPE_NUMERIC_LITERAL:
			// In input data such as echo x=3,y=4 | mlr put '$z=$x+$y', the 3 and 4 are strings
			// which need parsing as integers. But in DSL expression literals such as 'put $z = "3" + 4'
			// the "3" should not.
			return rval_evaluator_alloc_from_numeric_literal(pnode->text);
			break;

		case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:
			return rval_evaluator_alloc_from_boolean_literal(pnode->text);
			break;

		case MD_AST_NODE_TYPE_REGEXI:
			return rval_evaluator_alloc_from_string_literal(pnode->text);
			break;

		case MD_AST_NODE_TYPE_CONTEXT_VARIABLE:
			return rval_evaluator_alloc_from_context_variable(pnode->text);
			break;

		case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
			return rval_evaluator_alloc_from_local_variable(pnode->vardef_frame_relative_index);
			break;

		case MD_AST_NODE_TYPE_FULL_SREC:
			fprintf(stderr, "%s: $* is not valid within scalar contexts.\n",
				MLR_GLOBALS.bargv0);
			exit(1);

		case MD_AST_NODE_TYPE_FULL_OOSVAR:
			fprintf(stderr, "%s: @* is not valid within scalar contexts.\n",
				MLR_GLOBALS.bargv0);
			exit(1);

		case MD_AST_NODE_TYPE_MAP_LITERAL:
			fprintf(stderr, "%s: map-literals are not valid within scalar contexts.\n",
				MLR_GLOBALS.bargv0);
			exit(1);

		default:
			MLR_INTERNAL_CODING_ERROR();
			return NULL; // not reached
			break;
		}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	} else if (pnode->type == MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME) {
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: statements involving $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return rval_evaluator_alloc_from_indirect_field_name(pnode->pchildren->phead->pvvalue, pfmgr,
			type_inferencing, context_flags);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
		return rval_evaluator_alloc_from_oosvar_keylist(pnode, pfmgr, type_inferencing, context_flags);

	} else if (pnode->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE) {
		return rval_evaluator_alloc_from_local_map_keylist(pnode, pfmgr, type_inferencing, context_flags);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	} else if (pnode->type == MD_AST_NODE_TYPE_ENV) {
		return rval_evaluator_alloc_from_environment(pnode, pfmgr, type_inferencing, context_flags);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	} else if ((pnode->type == MD_AST_NODE_TYPE_FUNCTION_CALLSITE) || (pnode->type == MD_AST_NODE_TYPE_OPERATOR)) {
		return fmgr_alloc_provisional_from_operator_or_function_call(pfmgr, pnode, type_inferencing, context_flags);

	} else if (pnode->type == MD_AST_NODE_TYPE_FULL_SREC) {
		fprintf(stderr, "%s: $* is not valid within scalar contexts.\n",
			MLR_GLOBALS.bargv0);
		exit(1);

	} else if (pnode->type == MD_AST_NODE_TYPE_FULL_OOSVAR) {
		fprintf(stderr, "%s: @* is not valid within scalar contexts.\n",
			MLR_GLOBALS.bargv0);
		exit(1);

	} else if (pnode->type == MD_AST_NODE_TYPE_MAP_LITERAL) {
		fprintf(stderr, "%s: map-literals are not valid within scalar contexts.\n",
			MLR_GLOBALS.bargv0);
		exit(1);

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// This is the fall-through which typically gets hit when you update the AST-producing grammar but
	// haven't yet implemented the CST handler for it.
	} else {
		MLR_INTERNAL_CODING_ERROR();
		return NULL; // not reached
	}
}

// ================================================================
typedef struct _rval_evaluator_field_name_state_t {
	char* field_name;
} rval_evaluator_field_name_state_t;

static mv_t rval_evaluator_field_name_func_string_only(void* pvstate, variables_t* pvars) {
	rval_evaluator_field_name_state_t* pstate = pvstate;
	return get_srec_value_string_only(pstate->field_name, pvars->pinrec, pvars->ptyped_overlay);
}

static mv_t rval_evaluator_field_name_func_string_float(void* pvstate, variables_t* pvars) {
	rval_evaluator_field_name_state_t* pstate = pvstate;
	return get_srec_value_string_float(pstate->field_name, pvars->pinrec, pvars->ptyped_overlay);
}

static mv_t rval_evaluator_field_name_func_string_float_int(void* pvstate, variables_t* pvars) {
	rval_evaluator_field_name_state_t* pstate = pvstate;
	return get_srec_value_string_float_int(pstate->field_name, pvars->pinrec, pvars->ptyped_overlay);
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
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
	pevaluator->pfree_func = rval_evaluator_field_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_indirect_field_name_state_t {
	rval_evaluator_t* pname_evaluator;
} rval_evaluator_indirect_field_name_state_t;

static mv_t rval_evaluator_indirect_field_name_func_string_only(void* pvstate, variables_t* pvars) {
	rval_evaluator_indirect_field_name_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(pstate->pname_evaluator->pvstate, pvars);
	if (mv_is_null(&mvname)) {
		mv_free(&mvname);
		return mv_absent();
	}
	char free_flags = NO_FREE;
	char* indirect_field_name = mv_maybe_alloc_format_val(&mvname, &free_flags);

	mv_t rv = get_srec_value_string_only(indirect_field_name, pvars->pinrec, pvars->ptyped_overlay);
	if (free_flags & FREE_ENTRY_VALUE)
		free(indirect_field_name);
	mv_free(&mvname);
	return rv;
}

static mv_t rval_evaluator_indirect_field_name_func_string_float(void* pvstate, variables_t* pvars) {
	rval_evaluator_indirect_field_name_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(pstate->pname_evaluator->pvstate, pvars);
	if (mv_is_null(&mvname)) {
		mv_free(&mvname);
		return mv_absent();
	}
	char free_flags = NO_FREE;
	char* indirect_field_name = mv_maybe_alloc_format_val(&mvname, &free_flags);

	mv_t rv = get_srec_value_string_float(indirect_field_name, pvars->pinrec, pvars->ptyped_overlay);

	if (free_flags & FREE_ENTRY_VALUE)
		free(indirect_field_name);
	mv_free(&mvname);
	return rv;
}

static mv_t rval_evaluator_indirect_field_name_func_string_float_int(void* pvstate, variables_t* pvars) {
	rval_evaluator_indirect_field_name_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(pstate->pname_evaluator->pvstate, pvars);
	if (mv_is_null(&mvname)) {
		mv_free(&mvname);
		return mv_absent();
	}
	char free_flags = NO_FREE;
	char* indirect_field_name = mv_maybe_alloc_format_val(&mvname, &free_flags);

	mv_t rv = get_srec_value_string_float_int(indirect_field_name, pvars->pinrec, pvars->ptyped_overlay);

	if (free_flags & FREE_ENTRY_VALUE)
		free(indirect_field_name);
	mv_free(&mvname);
	return rv;
}

static void rval_evaluator_indirect_field_name_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_indirect_field_name_state_t* pstate = pevaluator->pvstate;
	pstate->pname_evaluator->pfree_func(pstate->pname_evaluator);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_indirect_field_name(mlr_dsl_ast_node_t* pnamenode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rval_evaluator_indirect_field_name_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_indirect_field_name_state_t));

	pstate->pname_evaluator = rval_evaluator_alloc_from_ast(pnamenode, pfmgr, type_inferencing, context_flags);

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	switch (type_inferencing) {
	case TYPE_INFER_STRING_ONLY:
		pevaluator->pprocess_func = rval_evaluator_indirect_field_name_func_string_only;
		break;
	case TYPE_INFER_STRING_FLOAT:
		pevaluator->pprocess_func = rval_evaluator_indirect_field_name_func_string_float;
		break;
	case TYPE_INFER_STRING_FLOAT_INT:
		pevaluator->pprocess_func = rval_evaluator_indirect_field_name_func_string_float_int;
		break;
	default:
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
	pevaluator->pfree_func = rval_evaluator_indirect_field_name_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_oosvar_keylist_state_t {
	sllv_t* poosvar_rhs_keylist_evaluators;
} rval_evaluator_oosvar_keylist_state_t;

mv_t rval_evaluator_oosvar_keylist_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_oosvar_keylist_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->poosvar_rhs_keylist_evaluators, pvars, &all_non_null_or_error);

	mv_t rv = mv_absent();
	if (all_non_null_or_error) {
		int error = 0;
		mv_t* pval = mlhmmv_root_look_up_and_ref_terminal(pvars->poosvars, pmvkeys, &error);
		if (pval != NULL) {
			if (pval->type == MT_STRING && *pval->u.strv == 0)
				rv = mv_empty();
			else
				rv = mv_copy(pval);
		}
	}

	sllmv_free(pmvkeys);
	return rv;
}

static void rval_evaluator_oosvar_keylist_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_oosvar_keylist_state_t* pstate = pevaluator->pvstate;
	for (sllve_t* pe = pstate->poosvar_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pevaluator = pe->pvvalue;
		pevaluator->pfree_func(pevaluator);
	}
	sllv_free(pstate->poosvar_rhs_keylist_evaluators);
	free(pstate);
	free(pevaluator);
}

// Example AST:
//
// $ mlr -n put -v '$y = @x[1]["two"][$3+4][@5]'
// list (statement_list):
//     = (srec_assignment):
//         y (field_name).
//         oosvar_keylist (oosvar_keylist):
//             x (string_literal).
//             1 (numeric_literal).
//             two (numeric_literal).
//             + (operator):
//                 3 (field_name).
//                 4 (numeric_literal).
//             oosvar_keylist (oosvar_keylist):
//                 5 (string_literal).

rval_evaluator_t* rval_evaluator_alloc_from_oosvar_keylist(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rval_evaluator_oosvar_keylist_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_oosvar_keylist_state_t));

	sllv_t* pkeylist_evaluators = sllv_alloc();

	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
		} else {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pfmgr,
				type_inferencing, context_flags));
		}
	}
	pstate->poosvar_rhs_keylist_evaluators = pkeylist_evaluators;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_oosvar_keylist_func;
	pevaluator->pfree_func = rval_evaluator_oosvar_keylist_free;

	return pevaluator;
}

// ================================================================
typedef struct _rval_evaluator_local_map_keylist_state_t {
	int vardef_frame_relative_index;
	sllv_t* plocal_map_rhs_keylist_evaluators;
} rval_evaluator_local_map_keylist_state_t;

mv_t rval_evaluator_local_map_keylist_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_local_map_keylist_state_t* pstate = pvstate;

	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->plocal_map_rhs_keylist_evaluators, pvars, &all_non_null_or_error);

	mv_t rv = mv_absent();
	if (all_non_null_or_error) {
		local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
		mv_t val = local_stack_frame_ref_terminal_from_indexed(pframe, pstate->vardef_frame_relative_index, pmvkeys);
		if (val.type == MT_STRING && *val.u.strv == 0)
			rv = mv_empty();
		else
			rv = mv_copy(&val);
	}

	sllmv_free(pmvkeys);
	return rv;
}

static void rval_evaluator_local_map_keylist_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_local_map_keylist_state_t* pstate = pevaluator->pvstate;
	for (sllve_t* pe = pstate->plocal_map_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
		rval_evaluator_t* pevaluator = pe->pvvalue;
		pevaluator->pfree_func(pevaluator);
	}
	sllv_free(pstate->plocal_map_rhs_keylist_evaluators);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_local_map_keylist(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rval_evaluator_local_map_keylist_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_local_map_keylist_state_t));

	MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index == MD_UNUSED_INDEX);

	pstate->vardef_frame_relative_index = pnode->vardef_frame_relative_index;

	sllv_t* pkeylist_evaluators = sllv_alloc();
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
		} else {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pfmgr,
				type_inferencing, context_flags));
		}
	}
	pstate->plocal_map_rhs_keylist_evaluators = pkeylist_evaluators;

	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = pstate;
	pevaluator->pprocess_func = NULL;
	pevaluator->pprocess_func = rval_evaluator_local_map_keylist_func;
	pevaluator->pfree_func = rval_evaluator_local_map_keylist_free;

	return pevaluator;
}

// ================================================================
// This is used for evaluating numbers in literal expressions, e.g. '$x = 4'

typedef struct _rval_evaluator_numeric_literal_state_t {
	mv_t literal;
} rval_evaluator_numeric_literal_state_t;

mv_t rval_evaluator_non_string_literal_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_numeric_literal_state_t* pstate = pvstate;
	return pstate->literal;
}

mv_t rval_evaluator_string_literal_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_numeric_literal_state_t* pstate = pvstate;
	char* input = pstate->literal.u.strv;

	if (pvars->ppregex_captures == NULL || *pvars->ppregex_captures == NULL) {
		return mv_from_string_no_free(input);
	} else {
		int was_allocated = FALSE;
		char* output = interpolate_regex_captures(input, *pvars->ppregex_captures, &was_allocated);
		if (was_allocated)
			return mv_from_string_with_free(output);
		else
			return mv_from_string_no_free(output);
	}
}
static void rval_evaluator_numeric_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_numeric_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

// How to handle echo a=1,b=2.0 | mlr put {flag} '$s = $a; $t = $b; $u = 3; $v = 4.0', where {flag} is -S, -F, or
// neither:
// * (no flag) TYPE_INFER_STRING_FLOAT_INT: a and s = int 1,      b and t = float 2.0,    u = int 3, v = float 4.0
// * (-F flag) TYPE_INFER_STRING_FLOAT:     a and s = float 1.0,  b and t = float 2.0,    u = int 3, v = float 4.0
// * (-S flag) TYPE_INFER_STRING_ONLY:      a and s = string "1", b and t = string "2.0", u = int 3, v = float 4.0
// The -S/-F flags for put/filter are for type inferencing in record data, not in literal expressions.

rval_evaluator_t* rval_evaluator_alloc_from_numeric_literal(char* string) {
	rval_evaluator_numeric_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_numeric_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (string == NULL) {
		pstate->literal = mv_absent();
		pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
	} else {
		long long intv;
		double fltv;

		pevaluator->pprocess_func = NULL;

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
	}
	pevaluator->pfree_func = rval_evaluator_numeric_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// This is used for evaluating strings and numbers in literal expressions, e.g. '$x = "abc"'
// or '$x = "left_\1". The values are subject to replacement with regex captures. See comments
// in mapper_put for more information.
//
// Compare rval_evaluator_alloc_from_string which doesn't do regex replacement: it is intended for
// oosvar names on expression left-hand sides (outside of this file).

typedef struct _rval_evaluator_string_literal_state_t {
	mv_t literal;
} rval_evaluator_string_literal_state_t;

static void rval_evaluator_string_literal_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_string_literal_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_string_literal(char* string) {
	rval_evaluator_string_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_string_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	if (string == NULL) {
		pstate->literal = mv_absent();
		pevaluator->pprocess_func = rval_evaluator_non_string_literal_func;
	} else {
		pstate->literal = mv_from_string_no_free(string);
		pevaluator->pprocess_func = rval_evaluator_string_literal_func;
	}
	pevaluator->pfree_func = rval_evaluator_string_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// This is intended only for oosvar names on expression left-hand sides (outside of this file).
// Compare rval_evaluator_alloc_from_string_literal.

typedef struct _rval_evaluator_string_state_t {
	char* string;
} rval_evaluator_string_state_t;

mv_t rval_evaluator_string_func(void* pvstate, variables_t* pvars) {
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

mv_t rval_evaluator_boolean_literal_func(void* pvstate, variables_t* pvars) {
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
		MLR_INTERNAL_CODING_ERROR();
	}
	pevaluator->pprocess_func = rval_evaluator_boolean_literal_func;
	pevaluator->pfree_func = rval_evaluator_boolean_literal_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

rval_evaluator_t* rval_evaluator_alloc_from_boolean(int boolval) {
	rval_evaluator_boolean_literal_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_boolean_literal_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	pstate->literal = mv_from_bool(boolval);
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
//         X (numeric_literal).
// AST END STATEMENTS (0):

// ----------------------------------------------------------------
typedef struct _rval_evaluator_environment_state_t {
	rval_evaluator_t* pname_evaluator;
} rval_evaluator_environment_state_t;

mv_t rval_evaluator_environment_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_environment_state_t* pstate = pvstate;

	mv_t mvname = pstate->pname_evaluator->pprocess_func(pstate->pname_evaluator->pvstate, pvars);
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

rval_evaluator_t* rval_evaluator_alloc_from_environment(mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr,
	int type_inferencing, int context_flags)
{
	rval_evaluator_environment_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_environment_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pnext->pvvalue;

	pstate->pname_evaluator = rval_evaluator_alloc_from_ast(pnamenode, pfmgr, type_inferencing, context_flags);
	pevaluator->pprocess_func = rval_evaluator_environment_func;
	pevaluator->pfree_func = rval_evaluator_environment_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
mv_t rval_evaluator_NF_func(void* pvstate, variables_t* pvars) {
	return mv_from_int(pvars->pinrec->field_count);
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
mv_t rval_evaluator_NR_func(void* pvstate, variables_t* pvars) {
	return mv_from_int(pvars->pctx->nr);
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
mv_t rval_evaluator_FNR_func(void* pvstate, variables_t* pvars) {
	return mv_from_int(pvars->pctx->fnr);
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
mv_t rval_evaluator_FILENAME_func(void* pvstate, variables_t* pvars) {
	return mv_from_string_no_free(pvars->pctx->filename);
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
mv_t rval_evaluator_FILENUM_func(void* pvstate, variables_t* pvars) {
	return mv_from_int(pvars->pctx->filenum);
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
mv_t rval_evaluator_PI_func(void* pvstate, variables_t* pvars) {
	return mv_from_float(M_PI);
}
static void rval_evaluator_PI_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_M_PI() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_PI_func;
	pevaluator->pfree_func = rval_evaluator_PI_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_E_func(void* pvstate, variables_t* pvars) {
	return mv_from_float(M_E);
}
static void rval_evaluator_E_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_M_E() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_E_func;
	pevaluator->pfree_func = rval_evaluator_E_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_IPS_func(void* pvstate, variables_t* pvars) {
	return mv_from_string_no_free(pvars->pctx->ips);
}
static void rval_evaluator_IPS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_IPS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_IPS_func;
	pevaluator->pfree_func = rval_evaluator_IPS_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_IFS_func(void* pvstate, variables_t* pvars) {
	return mv_from_string_no_free(pvars->pctx->ifs);
}
static void rval_evaluator_IFS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_IFS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_IFS_func;
	pevaluator->pfree_func = rval_evaluator_IFS_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_IRS_func(void* pvstate, variables_t* pvars) {
	context_t* pctx = pvars->pctx;
	return mv_from_string_no_free(
		pctx->auto_line_term_detected
			? pctx->auto_line_term
			: pctx->irs
	);
}
static void rval_evaluator_IRS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_IRS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_IRS_func;
	pevaluator->pfree_func = rval_evaluator_IRS_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_OPS_func(void* pvstate, variables_t* pvars) {
	return mv_from_string_no_free(pvars->pctx->ops);
}
static void rval_evaluator_OPS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_OPS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_OPS_func;
	pevaluator->pfree_func = rval_evaluator_OPS_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_OFS_func(void* pvstate, variables_t* pvars) {
	return mv_from_string_no_free(pvars->pctx->ofs);
}
static void rval_evaluator_OFS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_OFS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_OFS_func;
	pevaluator->pfree_func = rval_evaluator_OFS_free;
	return pevaluator;
}

// ----------------------------------------------------------------
mv_t rval_evaluator_ORS_func(void* pvstate, variables_t* pvars) {
	context_t* pctx = pvars->pctx;
	return mv_from_string_no_free(
		pctx->auto_line_term_detected
			? pctx->auto_line_term
			: pctx->ors
	);
}
static void rval_evaluator_ORS_free(rval_evaluator_t* pevaluator) {
	free(pevaluator);
}
rval_evaluator_t* rval_evaluator_alloc_from_ORS() {
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));
	pevaluator->pvstate = NULL;
	pevaluator->pprocess_func = rval_evaluator_ORS_func;
	pevaluator->pfree_func = rval_evaluator_ORS_free;
	return pevaluator;
}

// ================================================================
rval_evaluator_t* rval_evaluator_alloc_from_context_variable(char* variable_name) {
	if        (streq(variable_name, "NF"))       { return rval_evaluator_alloc_from_NF();
	} else if (streq(variable_name, "NR"))       { return rval_evaluator_alloc_from_NR();
	} else if (streq(variable_name, "FNR"))      { return rval_evaluator_alloc_from_FNR();
	} else if (streq(variable_name, "FILENAME")) { return rval_evaluator_alloc_from_FILENAME();
	} else if (streq(variable_name, "FILENUM"))  { return rval_evaluator_alloc_from_FILENUM();
	} else if (streq(variable_name, "M_PI"))     { return rval_evaluator_alloc_from_M_PI();
	} else if (streq(variable_name, "M_E"))      { return rval_evaluator_alloc_from_M_E();
	} else if (streq(variable_name, "IPS"))      { return rval_evaluator_alloc_from_IPS();
	} else if (streq(variable_name, "IFS"))      { return rval_evaluator_alloc_from_IFS();
	} else if (streq(variable_name, "IRS"))      { return rval_evaluator_alloc_from_IRS();
	} else if (streq(variable_name, "OPS"))      { return rval_evaluator_alloc_from_OPS();
	} else if (streq(variable_name, "OFS"))      { return rval_evaluator_alloc_from_OFS();
	} else if (streq(variable_name, "ORS"))      { return rval_evaluator_alloc_from_ORS();
	} else  { return NULL;
	}
}

// ================================================================
typedef struct _rval_evaluator_from_local_variable_state_t {
	int vardef_frame_relative_index;
} rval_evaluator_from_local_variable_state_t;

mv_t rval_evaluator_from_local_variable_func(void* pvstate, variables_t* pvars) {
	rval_evaluator_from_local_variable_state_t* pstate = pvstate;
	local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
	mv_t val = local_stack_frame_get_terminal_from_nonindexed(pframe, pstate->vardef_frame_relative_index);
	return mv_copy(&val);
}

static void rval_evaluator_from_local_variable_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_from_local_variable_state_t* pstate = pevaluator->pvstate;
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_local_variable(int vardef_frame_relative_index) {
	rval_evaluator_from_local_variable_state_t* pstate = mlr_malloc_or_die(
		sizeof(rval_evaluator_from_local_variable_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	MLR_INTERNAL_CODING_ERROR_IF(vardef_frame_relative_index == MD_UNUSED_INDEX);
	pstate->vardef_frame_relative_index = vardef_frame_relative_index;
	pevaluator->pprocess_func    = rval_evaluator_from_local_variable_func;
	pevaluator->pfree_func       = rval_evaluator_from_local_variable_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ----------------------------------------------------------------
typedef struct _rval_evaluator_mv_state_t {
	mv_t literal;
} rval_evaluator_mv_state_t;

mv_t rval_evaluator_mv_process(void* pvstate, variables_t* pvars) {
	rval_evaluator_mv_state_t* pstate = pvstate;
	return mv_copy(&pstate->literal);

}
static void rval_evaluator_mv_free(rval_evaluator_t* pevaluator) {
	rval_evaluator_mv_state_t* pstate = pevaluator->pvstate;
	mv_free(&pstate->literal);
	free(pstate);
	free(pevaluator);
}

rval_evaluator_t* rval_evaluator_alloc_from_mlrval(mv_t* pval) {
	rval_evaluator_mv_state_t* pstate = mlr_malloc_or_die(sizeof(rval_evaluator_mv_state_t));
	rval_evaluator_t* pevaluator = mlr_malloc_or_die(sizeof(rval_evaluator_t));

	pstate->literal = mv_copy(pval);
	pevaluator->pprocess_func = rval_evaluator_mv_process;
	pevaluator->pfree_func = rval_evaluator_mv_free;

	pevaluator->pvstate = pstate;
	return pevaluator;
}

// ================================================================
// Type-inferenced srec-field getters

// ----------------------------------------------------------------
mv_t get_srec_value_string_only(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, field_name);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		rv = mv_ref_type_infer_string(lrec_get(pinrec, field_name));
		rv = mv_copy(&rv);
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t get_srec_value_string_float(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, field_name);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		rv = mv_ref_type_infer_string_or_float(lrec_get(pinrec, field_name));
		rv = mv_copy(&rv);
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t get_srec_value_string_float_int(char* field_name, lrec_t* pinrec, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, field_name);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		rv = mv_ref_type_infer_string_or_float_or_int(lrec_get(pinrec, field_name));
		rv = mv_copy(&rv);
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t get_copy_srec_value_string_only_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, pentry->key);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		if (pentry->value == NULL) {
			rv = mv_absent();
		} else if (*pentry->value == 0) {
			rv = mv_empty();
		} else {
			rv = mv_from_string_with_free(mlr_strdup_or_die(pentry->value));
		}
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t get_copy_srec_value_string_float_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, pentry->key);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		if (pentry->value == NULL) {
			rv = mv_absent();
		} else if (*pentry->value == 0) {
			rv = mv_empty();
		} else {
			double fltv;
			if (mlr_try_float_from_string(pentry->value, &fltv)) {
				rv = mv_from_float(fltv);
			} else {
				rv = mv_from_string_with_free(mlr_strdup_or_die(pentry->value));
			}
		}
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t get_copy_srec_value_string_float_int_aux(lrece_t* pentry, lhmsmv_t* ptyped_overlay) {
	// See comments in rval_evaluator.h and mapper_put.c regarding the typed-overlay map.
	mv_t* poverlay = lhmsmv_get(ptyped_overlay, pentry->key);
	mv_t rv;
	if (poverlay != NULL) {
		// The lrec-evaluator logic will free its inputs and allocate new outputs, so we must copy
		// a value here to feed into that. Otherwise the typed-overlay map would have its contents
		// freed out from underneath it by the evaluator functions.
		rv = mv_copy(poverlay);
	} else {
		if (pentry->value == NULL) {
			rv = mv_absent();
		} else if (*pentry->value == 0) {
			rv = mv_empty();
		} else {
			long long intv;
			double fltv;
			if (mlr_try_int_from_string(pentry->value, &intv)) {
				rv = mv_from_int(intv);
			} else if (mlr_try_float_from_string(pentry->value, &fltv)) {
				rv = mv_from_float(fltv);
			} else {
				rv = mv_from_string_with_free(mlr_strdup_or_die(pentry->value));
			}
		}
	}
	return rv;
}
