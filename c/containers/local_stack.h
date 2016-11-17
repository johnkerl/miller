#ifndef LOCAL_STACK_H
#define LOCAL_STACK_H

#include "containers/mlrval.h"
#include "containers/type_decl.h"
#include "containers/sllv.h"
#include "containers/mlhmmv.h"

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
typedef struct _local_stack_frame_entry_t {
	char* name; // For type-check error messages. Not strduped; the caller must ensure extent.
	mlhmmv_value_t value;
	int   type_mask;
} local_stack_frame_entry_t;

typedef struct _local_stack_frame_t {
	int in_use;
	int ephemeral;
	int size;
	int subframe_base;
	local_stack_frame_entry_t* pvars;
} local_stack_frame_t;

// ----------------------------------------------------------------
// A stack is allocated for a top-level statement block: begin, end, or main, or
// user-defined function/subroutine. (The latter two may be called recursively
// in which case the in_use flag notes the need to allocate a new stack.)

local_stack_frame_t* local_stack_frame_alloc(int size);
void local_stack_frame_free(local_stack_frame_t* pframe);

// ================================================================
//#define LOCAL_STACK_TRACE_ENABLE
//#define LOCAL_STACK_BOUNDS_CHECK_ENABLE

void local_stack_bounds_check(local_stack_frame_t* pframe, char* op, int set, int vardef_frame_relative_index);
#ifdef LOCAL_STACK_BOUNDS_CHECK_ENABLE
#define LOCAL_STACK_BOUNDS_CHECK(pframe, op, set, vardef_frame_relative_index) \
	local_stack_bounds_check((pframe), (op), (set), (vardef_frame_relative_index))
#else
#define LOCAL_STACK_BOUNDS_CHECK(pframe, op, set, vardef_frame_relative_index)
#endif

#ifdef LOCAL_STACK_TRACE_ENABLE
#define LOCAL_STACK_TRACE(p) p
#else
#define LOCAL_STACK_TRACE(p)
#endif

// These are unconditional. With the single added character 'X' they can be
// used to focus verbosity at specific callsites for dev/debug.
#define LOCAL_STACK_BOUNDS_CHECKX(pframe, op, set, vardef_frame_relative_index) \
	local_stack_bounds_check((pframe), (op), (set), (vardef_frame_relative_index))

#define LOCAL_STACK_TRACEX(p) p

// ----------------------------------------------------------------
// Sets/clears the in-use flag for top-level statement blocks, and verifies the
// contract for absent-null at slot 0.

// For non-recursive functions/subroutines the enter method sets the in-use flag
// and returns its argument; the exit method clears that flag. For recursively
// invoked functions/subroutines the enter method returns another stack of the
// same size, and the exit method frees that.
//
// The reason we don't simply always allocate is that begin/main/end statements
// are never recursive, and most functions and subroutines are not recursive, so
// most of the time there will be a single frame for each. We allocate that once
// at startup, reuse it on every record, and free it at exit -- rather than
// allocating and freeing frames on every record.

local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe);
void local_stack_frame_exit(local_stack_frame_t* pframe);
void local_stack_frame_throw_type_mismatch(local_stack_frame_entry_t* pentry, mv_t* pval);
void local_stack_frame_throw_type_xmismatch(local_stack_frame_entry_t* pentry, mlhmmv_value_t* pxval); // xxx temp

// ----------------------------------------------------------------
static inline mv_t local_stack_frame_get_scalar_from_nonindexed(local_stack_frame_t* pframe,
	int vardef_frame_relative_index)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p GET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "GET", FALSE, vardef_frame_relative_index);
	// xxx encapsulate
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];
	mlhmmv_value_t* pvalue = &pentry->value;
	if (pvalue != NULL && pvalue->is_terminal) {
		return pvalue->mlrval;
	} else {
		return mv_absent();
	}
}

mv_t local_stack_frame_get_scalar_from_indexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, sllmv_t* pmvkeys);

mlhmmv_value_t* local_stack_frame_get_extended_from_indexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, sllmv_t* pmvkeys);

// ----------------------------------------------------------------
static inline void local_stack_frame_define_scalar(local_stack_frame_t* pframe, char* variable_name,
	int vardef_frame_relative_index, int type_mask, mv_t val)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "DEFINE", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	pentry->name = variable_name; // no strdup, for performance -- caller must ensure extent
	pentry->type_mask = type_mask;

	if (!(type_mask_from_mv(&val) & pentry->type_mask)) { // xxx temp
		local_stack_frame_throw_type_mismatch(pentry, &val);
	}

	// xxx temp -- make a single method
	mv_free(&pentry->value.mlrval); // xxx temp -- make value-free
	if (mv_is_absent(&val)) {
		mv_free(&val); // xxx confusing ownership semantics
	} else {
		pentry->value = mlhmmv_value_transfer_terminal(val); // xxx deep-copy?
	}
}

static inline void local_stack_frame_define_extended(local_stack_frame_t* pframe, char* variable_name, // xxx rename
	int vardef_frame_relative_index, int type_mask, mlhmmv_value_t xval)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "ASSIGN", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	pentry->name = variable_name; // no strdup, for performance -- caller must ensure extent
	pentry->type_mask = type_mask;

	// xxx subroutineize
	if (xval.is_terminal) {
		if (!(type_mask_from_mv(&xval.mlrval) & pentry->type_mask)) { // xxx temp
			local_stack_frame_throw_type_mismatch(pentry, &xval.mlrval);
		}
	} else {
		if (!(TYPE_MASK_MAP & pentry->type_mask)) { // xxx temp
			local_stack_frame_throw_type_mismatch(pentry, &xval.mlrval); // xxx allow xvals
		}
	}

	// xxx temp -- make a single method
	if (xval.is_terminal && mv_is_absent(&xval.mlrval)) {
		mlhmmv_free_submap(xval); // xxx confusing ownership semantics
	} else {
		mlhmmv_free_submap(pentry->value); // xxx rename
		pentry->value = xval;
	}
}


// ----------------------------------------------------------------
static inline void local_stack_frame_assign_scalar_nonindexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, mv_t val)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "ASSIGN", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	if (!(type_mask_from_mv(&val) & pentry->type_mask)) { // xxx temp
		local_stack_frame_throw_type_mismatch(pentry, &val);
	}

	// xxx temp -- make a single method
	mv_free(&pentry->value.mlrval); // xxx temp -- make value-free
	pentry->value = mlhmmv_value_transfer_terminal(val); // xxx deep-copy?
}

static inline void local_stack_frame_assign_extended_nonindexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, mlhmmv_value_t xval)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "ASSIGN", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	// xxx subroutineize
	if (xval.is_terminal) {
		if (!(type_mask_from_mv(&xval.mlrval) & pentry->type_mask)) { // xxx temp
			local_stack_frame_throw_type_mismatch(pentry, &xval.mlrval);
		}
	} else {
		if (!(TYPE_MASK_MAP & pentry->type_mask)) { // xxx temp
			local_stack_frame_throw_type_mismatch(pentry, &xval.mlrval); // xxx allow xvals
		}
	}

	// xxx temp -- make a single method
	mlhmmv_free_submap(pentry->value); // xxx rename
	pentry->value = xval;
}

void local_stack_frame_assign_scalar_indexed(local_stack_frame_t* pframe,
	int vardef_frame_relative_index, sllmv_t* pmvkeys,
	mv_t terminal_value);

void local_stack_frame_assign_extended_indexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, sllmv_t* pmvkeys,
	mlhmmv_value_t terminal_value);

// ----------------------------------------------------------------
// Frames are entered/exited for each curly-braced statement block, including
// the top-level block itself as well as ifs/fors/whiles.

static inline void local_stack_subframe_enter(local_stack_frame_t* pframe, int count) {
	LOCAL_STACK_TRACE(printf("LOCAL STACK SUBFRAME %p ENTER %d->%d\n",
		pframe, pframe->subframe_base, pframe->subframe_base+count));
	local_stack_frame_entry_t* psubframe = &pframe->pvars[pframe->subframe_base];
	for (int i = 0; i < count; i++) {
		LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p CLEAR %d\n", pframe, pframe->subframe_base+i));
		LOCAL_STACK_BOUNDS_CHECK(pframe, "CLEAR", FALSE, pframe->subframe_base+i);
		local_stack_frame_entry_t* pentry = &psubframe[i];
		pentry->value.is_terminal = TRUE; // xxx make inline method in mlhmmv for all of this
		mv_reset(&pentry->value.mlrval);
		pentry->type_mask = TYPE_MASK_ANY;
	}
	pframe->subframe_base += count;
}

static inline void local_stack_subframe_exit(local_stack_frame_t* pframe, int count) {
	LOCAL_STACK_TRACE(printf("LOCAL STACK SUBFRAME %p EXIT  %d->%d\n",
		pframe, pframe->subframe_base, pframe->subframe_base-count));
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
