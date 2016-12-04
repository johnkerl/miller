// ================================================================
// These evaluate right-hand-side values (rvals) and return mlrvals (mv_t).
//
// Record state is in three parts here:
//
// * The lrec is read for input fields; output fields are written to the typed-overlay map. It is up to the
//   caller to format the typed-overlay field values as strings and write them into the lrec.
//
// * Typed-overlay values are read in favor to the lrec: e.g. if the lrec has "x"=>"abc" and the typed overlay
//   has "x"=>3.7 then the evaluators will be presented with 3.7 for the value of the field named "x".
//
// * The =~ and !=~ operators populate the regex-captures array from \1, \2, etc. in the regex; the from-literal
//   evaluator interpolates those into output. Example:
//
//   o echo x=abc_def | mlr put '$name =~ "(.*)_(.*)"; $left = "\1"; $right = "\2"'
//
//   o The =~ resizes the regex-captures array to length 3 (1-up length 2), then copies "abc" to index 1
//     and "def" to index 2.
//
//   o The second expression writes "left"=>"abc" to the typed-overlay map; the third expression writes "right"=>"def"
//     to the typed-overlay map. The \1 and \2 get "abc" and "def" interpolated from the regex-captures array.
//
//   o It is up to mapper_put to write "left"=>"abc" and "right"=>"def" into the lrec.
//
// See also the comments above mapper_put.c for more information about left-hand sides (lvals).
// ================================================================

#ifndef RVAL_EVALUATOR_H
#define RVAL_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsmv.h"
#include "containers/mlhmmv.h"
#include "containers/mvfuncs.h"
#include "containers/local_stack.h"
#include "containers/loop_stack.h"
#include "lib/string_array.h"

// ----------------------------------------------------------------
// xxx needs to be in a different header file
typedef struct _return_state_t {
	mlhmmv_xvalue_t retval;
	int returned;
} return_state_t;

// ----------------------------------------------------------------
// xxx needs to be in a different header file
typedef struct _variables_t {
	lrec_t*          pinrec;
	lhmsmv_t*        ptyped_overlay;
	mlhmmv_root_t*        poosvars;
	string_array_t** ppregex_captures;
	context_t*       pctx;
	local_stack_t*   plocal_stack;
	loop_stack_t*    ploop_stack;
	return_state_t   return_state;
	int              trace_execution;
} variables_t;

// ----------------------------------------------------------------
// This is for scalar-valued contexts: almost all expressions except
// for rxval contexts (see below).

struct _rval_evaluator_t;  // forward reference for method declarations

typedef mv_t rval_evaluator_process_func_t(void* pvstate, variables_t* pvars);

typedef void rval_evaluator_free_func_t(struct _rval_evaluator_t*);

typedef struct _rval_evaluator_t {
	void* pvstate;
	rval_evaluator_process_func_t* pprocess_func;
	rval_evaluator_free_func_t*    pfree_func;
} rval_evaluator_t;

// ----------------------------------------------------------------
// This is for map-valued contexts: LHS/RHS of assignments,
// UDF/subroutine arguments, and UDF return values.

// The is_ephemeral flag is TRUE for map-literals, function return values, and
// data copied out of srecs.  It is FALSE when the pointer is into an existing
// data structure's memory (e.g. oosvars or locals).
typedef struct _boxed_xval_t {
	mlhmmv_xvalue_t xval;
	int is_ephemeral;
} boxed_xval_t;

static inline boxed_xval_t box_ephemeral_val(mv_t val) {
	return (boxed_xval_t) {
		.xval = mlhmmv_xvalue_wrap_terminal(val),
		.is_ephemeral = TRUE,
	};
}

static inline boxed_xval_t box_non_ephemeral_val(mv_t val) {
	return (boxed_xval_t) {
		.xval = mlhmmv_xvalue_wrap_terminal(val),
		.is_ephemeral = FALSE,
	};
}

static inline boxed_xval_t box_ephemeral_xval(mlhmmv_xvalue_t xval) {
	return (boxed_xval_t) {
		.xval = xval,
		.is_ephemeral = TRUE,
	};
}

static inline boxed_xval_t box_non_ephemeral_xval(mlhmmv_xvalue_t xval) {
	return (boxed_xval_t) {
		.xval = xval,
		.is_ephemeral = FALSE,
	};
}

// ----------------------------------------------------------------
struct _rxval_evaluator_t;  // forward reference for method declarations

typedef boxed_xval_t rxval_evaluator_process_func_t(void* pvstate, variables_t* pvars);

typedef void rxval_evaluator_free_func_t(struct _rxval_evaluator_t*);

typedef struct _rxval_evaluator_t {
	void* pvstate;
	rxval_evaluator_process_func_t* pprocess_func;
	rxval_evaluator_free_func_t*    pfree_func;
} rxval_evaluator_t;

#endif // RVAL_EVALUATOR_H
