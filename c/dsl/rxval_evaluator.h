// ================================================================
// These evaluate right-hand-side extended values (rxvals) and return the same.
//
// xxx update comment
//
// See also the comments above mapper_put.c for more information about left-hand sides (lvals).
// ================================================================

#ifndef RXVAL_EVALUATOR_H
#define RXVAL_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsmv.h"
#include "containers/mlhmmv.h"
#include "containers/mvfuncs.h"
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
