#ifndef LREC_EVALUATOR_H
#define LREC_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsv.h"
#include "containers/mlrval.h"
#include "containers/string_array.h"

struct _lrec_evaluator_t; // forward reference for method declarations
// xxx comment here
typedef mv_t lrec_evaluator_process_func_t(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate);
typedef void lrec_evaluator_free_func_t(struct _lrec_evaluator_t*);

typedef struct _lrec_evaluator_t {
	void* pvstate;
	lrec_evaluator_process_func_t* pprocess_func;
	lrec_evaluator_free_func_t*    pfree_func;
} lrec_evaluator_t;

#endif // LREC_EVALUATOR_H
