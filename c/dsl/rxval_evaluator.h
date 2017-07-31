// ================================================================
// These evaluate right-hand-side extended values (rxvals) and return the same.
//
// For scalar (non-extended) right-hand side values, everything is ephemeral as
// it propagates up through the concrete syntax tree: e.g. in '$c = $a . $b' the
// $a and $b are copied out as ephemerals; in the concat function their
// concatenation is computed and the ephemeral input arguments are freed; then
// the result is stored in field $c.
//
// But for extended values (here) everything is copy-on-write:
// expression-intermediate values are not always ephemeral.  This is due to the
// size of the data involved. We can do dump or emit of a nested hashmap stored
// in an oosvar or local without copying it; we can do mapdiff of two map-valued
// variables while not modifying or copying either argument.
//
// The boxed_xval_t decorates mlhmmv_value_t (extended value) with an
// is_ephemeral flag.  The mlhmmv_value_t in turn has a map or a scalar.
// ================================================================

#ifndef RXVAL_EVALUATOR_H
#define RXVAL_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsmv.h"
#include "containers/mlhmmv.h"
#include "lib/mvfuncs.h"
#include "containers/boxed_xval.h"
#include "containers/local_stack.h"
#include "containers/loop_stack.h"
#include "lib/string_array.h"

// ----------------------------------------------------------------
struct _rxval_evaluator_t;  // forward reference for method declarations

typedef boxed_xval_t rxval_evaluator_process_func_t(void* pvstate, variables_t* pvars);

typedef void rxval_evaluator_free_func_t(struct _rxval_evaluator_t*);

typedef struct _rxval_evaluator_t {
	void* pvstate;
	rxval_evaluator_process_func_t* pprocess_func;
	rxval_evaluator_free_func_t*    pfree_func;
} rxval_evaluator_t;

#endif // RXVAL_EVALUATOR_H
