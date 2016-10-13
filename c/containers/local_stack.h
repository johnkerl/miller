#ifndef LOCAL_STACK_H
#define LOCAL_STACK_H

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

// ----------------------------------------------------------------
// Data private to .c file
typedef struct _local_stack_frame_t local_stack_frame_t;

typedef struct _local_stack_t {
	int num_used;
	int num_allocated;
	local_stack_frame_t** ppframes;
	local_stack_frame_t* pbase_frame;
} local_stack_t;

// ----------------------------------------------------------------
// Constructors/destructors

local_stack_t* local_stack_alloc();
void local_stack_free(local_stack_t* pstack);

local_stack_frame_t* local_stack_frame_alloc_unfenced();
local_stack_frame_t* local_stack_frame_alloc_fenced();
void local_stack_frame_free(local_stack_frame_t* pframe);

// xxx comment
local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe);
void local_stack_frame_exit(local_stack_frame_t* pframe);

// ----------------------------------------------------------------
// Scope entry/exit

// To be called on entry to scoped block
void local_stack_push(local_stack_t* pstack, local_stack_frame_t* pframe);

// To be called on exit from scoped block.
local_stack_frame_t* local_stack_pop(local_stack_t* pstack);

// ----------------------------------------------------------------
// Access within scope

// Use of local variables on expression right-hand sides
mv_t* local_stack_resolve(local_stack_t* pstack, char* key);

// Use of local variables on expression left-hand sides
// The pmv is not copied. You may wish to mv_copy the argument you pass in.
// The pmv will be freed.

// xxx cmt
void local_stack_define(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags);
void local_stack_set(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags);

// Clears the binding from the top frame without popping it. Useful
// for clearing the baseframe which is never popped.
void local_stack_clear(local_stack_t* pstack);

// ----------------------------------------------------------------
// Test/debug

void local_stack_print(local_stack_t* pstack);

#endif // LOCAL_STACK_H
