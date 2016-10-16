#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/local_stack.h"

// ================================================================
static local_stack_frame_t* _local_stack_alloc(int size, int ephemeral) {
	local_stack_frame_t* pframe = mlr_malloc_or_die(sizeof(local_stack_frame_t));

	pframe->in_use = FALSE;
	pframe->ephemeral = ephemeral;
	pframe->size = size;
	pframe->subframe_base = 0;
	pframe->pvars = mlr_malloc_or_die(size * sizeof(mv_t));
	for (int i = 0; i < size; i++) {
		pframe->pvars[i] = mv_absent();
	}

// xxx #ifdef DEBUG
	// xxx printf(">>LOCAL_STACK_ALLOC %d %p\n", size, pframe);
// xxx #endif
	return pframe;
}

// ----------------------------------------------------------------
local_stack_frame_t* local_stack_frame_alloc(int size) {
	return _local_stack_alloc(size, FALSE);
}

// ----------------------------------------------------------------
void local_stack_frame_free(local_stack_frame_t* pframe) {
// xxx #ifdef DEBUG
	// xxx printf(">>LOCAL_STACK_FREE  %d %p\n", pframe->size, pframe);
// xxx #endif
	if (pframe == NULL)
		return;
	for (int i = 0; i < pframe->size; i++) {
		mv_free(&pframe->pvars[i]);
	}
	free(pframe->pvars);
	free(pframe);
}

// ----------------------------------------------------------------
// xxx cmt
local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe) {
	if (!pframe->in_use) {
// xxx #ifdef DEBUG
	// xxx printf(">>LOCAL_STACK_FRAME_ENTER NONREC %d %p\n", pframe->size, pframe);
// xxx #endif
		pframe->in_use = TRUE;
		return pframe;
	} else {
		local_stack_frame_t* prv = _local_stack_alloc(pframe->size, TRUE);
// xxx #ifdef DEBUG
	// xxx printf(">>LOCAL_STACK_FRAME_ENTER REC %d %p/%p\n", pframe->size, pframe, prv);
// xxx #endif
		prv->in_use = TRUE;
		return prv;
	}
}

// ----------------------------------------------------------------
void local_stack_frame_exit (local_stack_frame_t* pframe) {
	MLR_INTERNAL_CODING_ERROR_UNLESS(mv_is_absent(&pframe->pvars[0]));
	if (!pframe->ephemeral) {
		pframe->in_use = FALSE;
	} else {
		local_stack_frame_free(pframe);
	}
}

// ================================================================
local_stack_t* local_stack_alloc() {
	local_stack_t* pstack = mlr_malloc_or_die(sizeof(local_stack_t));

	pstack->pframes = sllv_alloc();

	return pstack;
}

// ----------------------------------------------------------------
void local_stack_free(local_stack_t* pstack) {
	if (pstack == NULL)
		return;
	for (sllve_t* pe = pstack->pframes->phead; pe != NULL; pe = pe->pnext) {
		local_stack_frame_free(pe->pvvalue);
	}
	sllv_free(pstack->pframes);
	free(pstack);
}

// ----------------------------------------------------------------
void local_stack_push(local_stack_t* pstack, local_stack_frame_t* pframe) {
	// xxx rename to sllv_push throughout
	sllv_prepend(pstack->pframes, pframe);
}

local_stack_frame_t* local_stack_pop(local_stack_t* pstack) {
	return sllv_pop(pstack->pframes);
}
