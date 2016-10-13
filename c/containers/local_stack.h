#ifndef LOCAL_STACK_H
#define LOCAL_STACK_H

#include "containers/mlrval.h"

// Bound & scoped variables for use in for-loops, function bodies, and
// subroutine bodies. Indices of local variables, and max-depth for top-level
// statement blocks, are compted by the stack-allocator which marks up the AST
// before the CST is built from it.
//
// A convention shared between the stack-allocator and this data structure is
// that slot 0 is an absent-null which is used for reads of undefined (or
// as-yet-undefined) local variables.

// ----------------------------------------------------------------
typedef struct _local_stack_t {
	int in_use;
	int size;
	int frame_base;
	mv_t* pvars;
	// xxx has-absent-read flag ...
} local_stack_t;

// ----------------------------------------------------------------
// A stack is allocated for a top-level statement block: begin, end, or main, or
// user-defined function/subroutine. (The latter two may be called recursively
// in which case the in_use flag notes the need to allocate a new stack.)

local_stack_t* local_stack_alloc(int size);
void local_stack_free(local_stack_t* pstack);

// Frames are entered/exited for each curly-braced statement block, including
// the top-level block itself as well as ifs/fors/whiles.

void local_stack_frame_enter(local_stack_t* pstack, int count);
void local_stack_frame_exit (local_stack_t* pstack, int count);

#endif // LOCAL_STACK_H
