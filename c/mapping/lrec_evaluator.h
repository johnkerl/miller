#ifndef LREC_EVALUATOR_H
#define LREC_EVALUATOR_H

#include "lib/context.h"
#include "containers/lrec.h"
#include "containers/lhmsv.h"
#include "containers/mlrval.h"
#include "lib/string_array.h"

struct _lrec_evaluator_t; // forward reference for method declarations

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
// See also the comments above mapper_put.c for more information.

typedef mv_t lrec_evaluator_process_func_t(lrec_t* prec, lhmsv_t* ptyped_overlay, string_array_t* pregex_captures,
	context_t* pctx, void* pvstate);

typedef void lrec_evaluator_free_func_t(struct _lrec_evaluator_t*);

typedef struct _lrec_evaluator_t {
	void* pvstate;
	lrec_evaluator_process_func_t* pprocess_func;
	lrec_evaluator_free_func_t*    pfree_func;
} lrec_evaluator_t;

#endif // LREC_EVALUATOR_H
