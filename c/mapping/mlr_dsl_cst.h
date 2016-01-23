#ifndef MLR_DSL_CST_H
#define MLR_DSL_CST_H

#include "lrec_evaluators.h"

// ================================================================
// Concrete syntax tree (CST) derived from an abstract syntax tree (AST).
//
// At present (January 2016) the grammar for 'mlr put' is flat: there is no looping or branching (other than the ternary
// operator, and the gate statement).  Statements are of the form:
//
// * Assignment to oosvar (out-of-stream variables, prefixed with @ sigil)
//
// * Assignment to srec (in-stream records, with field names prefixed with $ sigil)
//
// * filter statements: if false, do not pass the current record from the input stream to the output stream
//
// * gate statements: if false, stop executing further statements
//
// * bare-boolean statements: no-ops unless they have side effects: namely, the matches/does-not-match
//   operators =~ and !=~ setting regex captures \1, \2, etc.
//
// * emit statements: these place oosvar key-value pairs into the output stream.  These can be of the form 'emit @a;
//   emit @b' which produce separate records such as a=3 and b=4, or of the form 'emit @a, @b' which produce records
//   such as a=3,b=4.
//
// This means that as of the present, the CSTs are just list of left-hand sides and right-hand sides.  The left-hand
// side is an output field name, with an is-oosvar flag. For filter/gate/bare-boolean these are unused; for assignments
// and emits, the output field names are used. The right-hand sides are lrec-evaluators which take an input record and
// oosvar-map, and produce a mlrval which can be assigned. Additionally, for all but emit, there is a single LHS name
// and single RHS lrec-evaluator per statement.
//
// Further, these statements are organized into three groups:
//
// * begin: executed once, before the first input record is read.
//
// * main : executed for each input record.
//
// * end :  executed once, after the last input record is read.

// ================================================================

// begin list of:
// main  list of:
// end   list of:
//
// * node type
// * list of:
//   o output field name
//   o evaluator
//   o is_oosvar

xypedef struct _mlr_dsl_cst_statement_item_t {
	lrec_evaluator_t* pevaluator;
	char* output_field_name;
	int is_oosvar;
} mlr_dsl_cst_statement_item_t;

//typedef struct _mlr_dsl_cst_statement_t {
//	int ast_statement_type;
//	sllv_t* pstatement_items;
//} mlr_dsl_cst_statement_t;

#endif // MLR_DSL_CST_H
