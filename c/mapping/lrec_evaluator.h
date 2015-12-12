#ifndef LREC_EVALUATOR_H
#define LREC_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/mlrval.h"

typedef mv_t lrec_evaluator_process_func_t(lrec_t* prec, context_t* pctx, void* pvstate);
typedef void lrec_evaluator_free_func_t(void* pvstate);

typedef struct _lrec_evaluator_t {
	void* pvstate;
	lrec_evaluator_process_func_t* pprocess_func;
	lrec_evaluator_free_func_t*    pfree_func;
} lrec_evaluator_t;

#endif // LREC_EVALUATOR_H
