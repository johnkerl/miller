#ifndef LREC_EVALUATOR_H
#define LREC_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "mapping/mlr_val.h"

typedef mv_t lrec_evaluator_func_t(lrec_t* prec, context_t* pctx, void* pvstate);

typedef struct _lrec_evaluator_t {
	void* pvstate;
	// xxx needs a pfree_func too
	lrec_evaluator_func_t* pevaluator_func;
} lrec_evaluator_t;

#endif // LREC_EVALUATOR_H
