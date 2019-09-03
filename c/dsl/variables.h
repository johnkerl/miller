// ================================================================
// Record state is in multiple parts:
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
// See also the comments above mapper_put.c for more information about left-hand sides (lvals).
// ================================================================

#ifndef VARIABLES_H
#define VARIABLES_H

#include "containers/lrec.h"
#include "containers/lhmsmv.h"
#include "lib/string_array.h"
#include "containers/mlhmmv.h"
#include "lib/context.h"
#include "containers/local_stack.h"
#include "dsl/return_state.h"

// Context for DSL evaluation
typedef struct _variables_t {
	lrec_t*          pinrec;
	lhmsmv_t*        ptyped_overlay;
	string_array_t** ppregex_captures;
	mlhmmv_root_t*   poosvars;
	context_t*       pctx;
	local_stack_t*   plocal_stack;
	loop_stack_t*    ploop_stack;
	return_state_t   return_state;
	int              trace_execution;
	int              json_quote_int_keys;
	int              json_quote_non_string_values;
	int              json_apply_ofmt_to_floats;
} variables_t;

#endif // VARIABLES_H
