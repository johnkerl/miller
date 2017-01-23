// ================================================================
// These evaluate right-hand-side values (rvals) and return mlrvals (mv_t).
// This is for scalar-valued contexts: almost all expressions except for rxval
// contexts.
// ================================================================

#ifndef RVAL_EVALUATOR_H
#define RVAL_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsmv.h"
#include "containers/mlhmmv.h"
#include "containers/mvfuncs.h"
#include "containers/boxed_xval.h"
#include "containers/local_stack.h"
#include "containers/loop_stack.h"
#include "lib/string_array.h"
#include "dsl/variables.h"


struct _rval_evaluator_t;  // forward reference for method declarations

typedef mv_t rval_evaluator_process_func_t(void* pvstate, variables_t* pvars);

typedef void rval_evaluator_free_func_t(struct _rval_evaluator_t*);

typedef struct _rval_evaluator_t {
	void* pvstate;
	rval_evaluator_process_func_t* pprocess_func;
	rval_evaluator_free_func_t*    pfree_func;
} rval_evaluator_t;

#endif // RVAL_EVALUATOR_H
