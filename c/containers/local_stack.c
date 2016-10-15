#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/local_stack.h"

// ----------------------------------------------------------------
static local_stack_t* _local_stack_alloc(int size, int ephemeral) {
	local_stack_t* pstack = mlr_malloc_or_die(sizeof(local_stack_t));

	pstack->in_use = FALSE;
	pstack->ephemeral = ephemeral;
	pstack->size = size;
	pstack->frame_base = 0;
	pstack->pvars = mlr_malloc_or_die(size * sizeof(mv_t));
	for (int i = 0; i < size; i++) {
		pstack->pvars[i] = mv_absent();
	}

	return pstack;
}

// ----------------------------------------------------------------
local_stack_t* local_stack_alloc(int size) {
	return _local_stack_alloc(size, FALSE);
}

// ----------------------------------------------------------------
void local_stack_free(local_stack_t* pstack) {
	if (pstack == NULL)
		return;
	for (int i = 0; i < pstack->size; i++) {
		mv_free(&pstack->pvars[i]);
	}
	free(pstack->pvars);
	free(pstack);
}

// ----------------------------------------------------------------
// xxx cmt
local_stack_t* local_stack_enter(local_stack_t* pstack) {
	if (!pstack->in_use) {
		pstack->in_use = TRUE;
		return pstack;
	} else {
		local_stack_t* prv = _local_stack_alloc(pstack->size, TRUE);
		prv->in_use = TRUE;
		return prv;
	}
}

// ----------------------------------------------------------------
void local_stack_exit (local_stack_t* pstack) {
	MLR_INTERNAL_CODING_ERROR_UNLESS(mv_is_absent(&pstack->pvars[0]));
	if (!pstack->ephemeral) {
		pstack->in_use = FALSE;
	} else {
		local_stack_free(pstack);
	}
}
