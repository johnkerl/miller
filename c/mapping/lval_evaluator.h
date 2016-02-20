// ================================================================
// These evaluate left-hand-side values (lvals).
//
// See also rval_evaluators.h for context.
// ================================================================

#ifndef LVAL_EVALUATOR_H
#define LVAL_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsv.h"
#include "containers/mlhmmv.h"
#include "containers/mlrval.h"
#include "lib/string_array.h"

struct _lval_evaluator_t; // forward reference for method declarations

typedef void lval_evaluator_process_func_t(
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	mlhmmv_t*        poosvars,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs);

typedef void lval_evaluator_free_func_t(struct _lval_evaluator_t*);

typedef struct _lval_evaluator_t {
	void* pvstate;
	lval_evaluator_process_func_t* pprocess_func;
	lval_evaluator_free_func_t*    pfree_func;
} lval_evaluator_t;

#endif // LVAL_EVALUATOR_H
