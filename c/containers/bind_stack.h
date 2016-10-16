// xxx rm entirely
//#ifndef BIND_STACK_H
//#define BIND_STACK_H
//
//#include "containers/mlrval.h"
//#include "containers/lhmsmv.h"
//
//// Bound & scoped variables for use in for-loops, function bodies, and
//// subroutine bodies.  A frame is 'fenced' if it's at entry to a
//// function/subroutine body. Example:
////
//// mlr ... put -f '
////   func f(x) {
////     return k;
////   }
////   for (k,v in $*) {
////     $[k] = f(v)
////   }
//// '
////
//// There is a bind-stack created on entering the put. A frame is pushed in the
//// for-loop.  A second frame is pushed in the call to f(x). The former should
//// have access to k; the latter should not.
//
//// ----------------------------------------------------------------
//// Data private to .c file
//typedef struct _bind_stack_frame_t bind_stack_frame_t;
//
//typedef struct _bind_stack_t {
//	int num_used;
//	int num_allocated;
//	bind_stack_frame_t** ppframes;
//	bind_stack_frame_t* pbase_frame;
//} bind_stack_t;
//
//// ----------------------------------------------------------------
//// Constructors/destructors
//
//bind_stack_t* bind_stack_alloc();
//void bind_stack_free(bind_stack_t* pstack);
//
//bind_stack_frame_t* bind_stack_frame_alloc_unfenced();
//bind_stack_frame_t* bind_stack_frame_alloc_fenced();
//void bind_stack_frame_free(bind_stack_frame_t* pframe);
//
//// xxx comment
//bind_stack_frame_t* bind_stack_frame_enter(bind_stack_frame_t* pframe);
//void bind_stack_frame_exit(bind_stack_frame_t* pframe);
//
//// ----------------------------------------------------------------
//// Scope entry/exit
//
//// To be called on entry to scoped block
//void bind_stack_push(bind_stack_t* pstack, bind_stack_frame_t* pframe);
//
//// To be called on exit from scoped block.
//bind_stack_frame_t* bind_stack_pop(bind_stack_t* pstack);
//
//// ----------------------------------------------------------------
//// Access within scope
//
//// Use of local variables on expression right-hand sides
//// xxx rm mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key);
//
//// Use of local variables on expression left-hand sides
//// The pmv is not copied. You may wish to mv_copy the argument you pass in.
//// The pmv will be freed.
//
//// xxx cmt
//// xxx rm void bind_stack_define(bind_stack_t* pstack, char* name, mv_t* pmv, char free_flags);
//// xxx rm void bind_stack_set(bind_stack_t* pstack, char* name, mv_t* pmv, char free_flags);
//
//// Clears the binding from the top frame without popping it. Useful
//// for clearing the baseframe which is never popped.
//void bind_stack_clear(bind_stack_t* pstack);
//
//// ----------------------------------------------------------------
//// Test/debug
//
//void bind_stack_print(bind_stack_t* pstack);
//
//#endif // BIND_STACK_H
