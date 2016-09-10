#ifndef BIND_STACK_H
#define BIND_STACK_H

#include "containers/mlrval.h"
#include "containers/lhmsmv.h"

// Bound & scoped variables for use in for-loops, function bodies, and
// subroutine bodies.  A frame is 'fenced' if it's at entry to a
// function/subroutine body. Example:
//
// mlr ... put -f '
//   func f(x) {
//     return k;
//   }
//   for (k,v in $*) {
//     $[k] = f(v)
//   }
// '
//
// There is a bind-stack created on entering the put. A frame is pushed in the
// for-loop.  A second frame is pushed in the call to f(x). The former should
// have access to k; the latter should not.

typedef struct _bind_stack_frame_t {
	lhmsmv_t*  bindings;
	int        fenced;
} bind_stack_frame_t;

typedef struct _bind_stack_t {
	int num_used;
	int num_allocated;
	bind_stack_frame_t* pframes;
} bind_stack_t;

bind_stack_t* bind_stack_alloc();
void bind_stack_free(bind_stack_t* pstack);

void bind_stack_push(bind_stack_t* pstack, lhmsmv_t* bindings);
void bind_stack_push_fenced(bind_stack_t* pstack, lhmsmv_t* bindings);
void bind_stack_pop(bind_stack_t* pstack);
mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key);

void bind_stack_print(bind_stack_t* pstack);

#endif // BIND_STACK_H
