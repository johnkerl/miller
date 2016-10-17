#ifndef LOCAL_STACK_H
#define LOCAL_STACK_H

#include "containers/mlrval.h"
#include "containers/sllv.h"

// ================================================================
// Bound & scoped variables for use in for-loops, function bodies, and
// subroutine bodies. Indices of local variables, and max-depth for top-level
// statement blocks, are compted by the stack-allocator which marks up the AST
// before the CST is built from it.
//
// A convention shared between the stack-allocator and this data structure is
// that slot 0 is an absent-null which is used for reads of undefined (or
// as-yet-undefined) local variables.
// ================================================================

// ================================================================
typedef struct _local_stack_frame_t {
	int in_use;
	int ephemeral;
	int size;
	int subframe_base;
	mv_t* pvars;
	// xxx has-absent-read flag ...
} local_stack_frame_t;

// ----------------------------------------------------------------
// A stack is allocated for a top-level statement block: begin, end, or main, or
// user-defined function/subroutine. (The latter two may be called recursively
// in which case the in_use flag notes the need to allocate a new stack.)

local_stack_frame_t* local_stack_frame_alloc(int size);
void local_stack_frame_free(local_stack_frame_t* pframe);

// ----------------------------------------------------------------
#define LOCAL_STACK_BOUNDS_CHECK_ENABLE // xxx disable
#ifdef LOCAL_STACK_BOUNDS_CHECK_ENABLE

static void local_stack_bounds_check(local_stack_frame_t* pframe, char* op, int frame_relative_index) {
	if (frame_relative_index < 0) {
		fprintf(stderr, "OP=%s FRAME=%p IDX=%d/%d STACK UNDERFLOW\n", op, pframe, frame_relative_index, pframe->size);
		exit(1);
	}
	if (frame_relative_index >= pframe->size) {
		fprintf(stderr, "OP=%s FRAME=%p IDX=%d/%d STACK OVERFLOW\n", op, pframe, frame_relative_index, pframe->size);
		exit(1);
	}
}
#define LOCAL_STACK_BOUNDS_CHECK(pframe, op, frame_relative_index) \
	local_stack_bounds_check((pframe), (op), (frame_relative_index))

#else

#define LOCAL_STACK_BOUNDS_CHECK(pframe, op, frame_relative_index)

#endif

// ----------------------------------------------------------------
// Sets/clear the in-use flag for top-level statement blocks, and verifies
// the contract for absent-null at slot 0.

// For non-recursive functions/subroutines the enter method sets the in-use flag
// and returns its argument; the exit method clears that flag. For recursively
// invoked functions/subroutines the enter method returns another stack of the
// same size, and the exit method frees that.
local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe);
void local_stack_frame_exit(local_stack_frame_t* pframe);

static inline mv_t* local_stack_frame_get(local_stack_frame_t* pframe, int frame_relative_index) {
	LOCAL_STACK_BOUNDS_CHECK(pframe, "GET", frame_relative_index);
	return &pframe->pvars[frame_relative_index];
}
static inline void local_stack_frame_set(local_stack_frame_t* pframe, int frame_relative_index, mv_t val) {
	LOCAL_STACK_BOUNDS_CHECK(pframe, "SET", frame_relative_index);
	pframe->pvars[frame_relative_index] = val;
}

// ----------------------------------------------------------------
// Frames are entered/exited for each curly-braced statement block, including
// the top-level block itself as well as ifs/fors/whiles.

static inline void local_stack_subframe_enter(local_stack_frame_t* pframe, int count) {
	// xxx try to avoid with absent-read flag at stack-allocator ...
	mv_t* psubframe = &pframe->pvars[pframe->subframe_base];
	for (int i = 0; i < count; i++) {
		psubframe[i] = mv_absent();
	}
	pframe->subframe_base += count;
}
static inline void local_stack_subframe_exit(local_stack_frame_t* pframe, int count) {
	pframe->subframe_base -= count;
}

// ================================================================
typedef struct _local_stack_t {
	sllv_t* pframes;
} local_stack_t;

local_stack_t* local_stack_alloc();
void local_stack_free(local_stack_t* pstack);

void local_stack_push(local_stack_t* pstack, local_stack_frame_t* pframe);

local_stack_frame_t* local_stack_pop(local_stack_t* pstack);

static inline local_stack_frame_t* local_stack_get_top_frame(local_stack_t* pstack) {
	return pstack->pframes->phead->pvvalue;
}

#endif // LOCAL_STACK_H
