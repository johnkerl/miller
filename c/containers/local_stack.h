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
	mv_t* pvars;
} local_stack_t;

// ----------------------------------------------------------------
// Constructors/destructors

local_stack_t* local_stack_alloc(int size);
void local_stack_free(local_stack_t* pstack);

//local_stack_frame_t* local_stack_frame_alloc_unfenced();
//local_stack_frame_t* local_stack_frame_alloc_fenced();
//void local_stack_frame_free(local_stack_frame_t* pframe);
//
//// xxx comment
//local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe);
//void local_stack_frame_exit(local_stack_frame_t* pframe);
//
//// ----------------------------------------------------------------
//// Scope entry/exit
//
//// To be called on entry to scoped block
//void local_stack_push(local_stack_t* pstack, local_stack_frame_t* pframe);
//
//// To be called on exit from scoped block.
//local_stack_frame_t* local_stack_pop(local_stack_t* pstack);
//
//// ----------------------------------------------------------------
//// Access within scope
//
//// Use of local variables on expression right-hand sides
//mv_t* local_stack_resolve(local_stack_t* pstack, char* key);
//
//// Use of local variables on expression left-hand sides
//// The pmv is not copied. You may wish to mv_copy the argument you pass in.
//// The pmv will be freed.
//
//// xxx cmt
//void local_stack_define(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags);
//void local_stack_set(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags);
//
//// Clears the binding from the top frame without popping it. Useful
//// for clearing the baseframe which is never popped.
//void local_stack_clear(local_stack_t* pstack);
//
//// ----------------------------------------------------------------
//// Test/debug
//
//void local_stack_print(local_stack_t* pstack);

#endif // LOCAL_STACK_H
