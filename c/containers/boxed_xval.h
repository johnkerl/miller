// ================================================================
// This is for map-valued contexts: LHS/RHS of assignments,
// UDF/subroutine arguments, and UDF return values.

// The is_ephemeral flag is TRUE for map-literals, function return values, and
// data copied out of srecs.  It is FALSE when the pointer is into an existing
// data structure's memory (e.g. oosvars or locals).
// ================================================================

#ifndef BOXED_XVAL_H
#define BOXED_XVAL_H

#include "../lib/mlrval.h"
#include "../containers/mlhmmv.h"

// ----------------------------------------------------------------
typedef struct _boxed_xval_t {
	mlhmmv_xvalue_t xval;
	char is_ephemeral;
} boxed_xval_t;

// ----------------------------------------------------------------
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

#endif // BOXED_XVAL_H
