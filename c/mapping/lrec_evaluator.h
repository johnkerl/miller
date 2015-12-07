#ifndef LREC_EVALUATOR_H
#define LREC_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/mlrval.h"

typedef mv_t lrec_evaluator_func_t(lrec_t* prec, context_t* pctx, void* pvstate);

typedef struct _lrec_evaluator_t {
	void* pvstate;
	lrec_evaluator_func_t* pevaluator_func;
	// xxx needs a pfree_func too
} lrec_evaluator_t;

#endif // LREC_EVALUATOR_H
