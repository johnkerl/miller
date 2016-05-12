#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "containers/mlr_dsl_ast.h"
#include "rval_evaluators.h"

// ================================================================
// Concrete syntax tree (CST) derived from an abstract syntax tree (AST).
//
// Statements are of the form:
//
// * Assignment of mlrval (i.e. result of expression evaluation, e.g. $name or f($x,$y)) to oosvar (out-of-stream
// variables, prefixed with @ sigil)
//
// * Assignment to srec (in-stream records, with field names prefixed with $ sigil)
//
// * Copying full srec ($* syntax) to/from an oosvar
//
// * Oosvar-to-oosvar assignments (recursively if RHS is non-terminal)
//
// * pattern-action statements: boolean expression with curly-braced statements which are executed only
//   when the boolean evaluates to true.
//
// * bare-boolean statements: no-ops unless they have side effects: namely, the matches/does-not-match
//   operators =~ and !=~ setting regex captures \1, \2, etc.
//
// * emit statements: these place oosvar key-value pairs into the output stream.  These can be of the following forms:
//
//   o 'emit @a; emit @b' which produce separate records such as a=3 and b=4
//
//   o 'emitf @a, @b' which produce records such as a=3,b=4
//
//   o For nested maps, 'emit @c, "x", "y"' in which case the first two map levels are pulled out and named "x" and "y"
//   in separate fields. See containers/mlhmmv.h for more information.
//
// Further, these statements are organized into three groups:
//
// * begin: executed once, before the first input record is read.
// * main:  executed for each input record.
// * end:   executed once, after the last input record is read.
//
// The exceptions being, of course, assignment to/from srec is disallowed for begin/end statements since those occur
// before/after stream processing, respectively.
// ================================================================

struct _mlr_dsl_cst_statement_t;

typedef void mlr_dsl_cst_node_evaluator_func_t(
	struct _mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator);

// Most statements have one item, except multi-oosvar emit and multi-oosvar unset.
typedef struct _mlr_dsl_cst_statement_item_t {
	// E.g. emit @a[$b]["c"], "d", @e: keylist is [$b, "c"] and namelist is ["d", @e].
	sllv_t* punset_oosvar_keylist_evaluators; // xxx temp
	char* emitf_or_unset_srec_field_name;     // xxx temp
	rval_evaluator_t* pemitf_arg_evaluator;   // xxx temp
} mlr_dsl_cst_statement_item_t;

typedef struct _mlr_dsl_cst_statement_t {
	mlr_dsl_cst_node_evaluator_func_t* pevaluator;

	// For assignment to oosvar, emit, and emitp
	sllv_t* poosvar_lhs_keylist_evaluators;

	// For assignment to srec
	char* srec_lhs_field_name;

	// For assignments to srec or oosvar, filter, cond, and bare-boolean
	rval_evaluator_t* prhs_evaluator;

	// For assigning full srec from oosvar
	sllv_t* poosvar_rhs_keylist_evaluators;

	// For emit/emitp
	sllv_t* pemit_oosvar_namelist_evaluators;

	// xxx temp
	sllv_t* pitems;

	// For pattern-action blocks, while, for, etc.
	sllv_t* pblock_statements;

} mlr_dsl_cst_statement_t;

// ---------------------------------------------------------------- xxx
// cond-expr {}
// while (expr) {}
// for (k, v in $*) {}
// for (k1, k2, v in @v["a"]) {}
// if (expr) {} elif (expr) {} elif (expr) else {}
// $srec = RHS
// @v["a"] = $*
// $* = @v["a"]
// bare-boolean
// @v["a"] = RHS
// filter expr
// unset
// emitf
// emitp
// emit
// dump
// break
// continue
// ---------------------------------------------------------------- xxx

typedef struct _mlr_dsl_cst_t {
	sllv_t* pbegin_statements;
	sllv_t* pmain_statements;
	sllv_t* pend_statements;
} mlr_dsl_cst_t;

// ----------------------------------------------------------------
mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing);
void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst);

void mlr_dsl_cst_evaluate(
	sllv_t*          pcst_statements, // begin/main/end
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator);

#endif // MLR_DSL_CST_H
