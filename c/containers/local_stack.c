#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/local_stack.h"

// ----------------------------------------------------------------
local_stack_t* local_stack_alloc(int size) {
	local_stack_t* pstack = mlr_malloc_or_die(sizeof(local_stack_t));

	pstack->in_use = FALSE;
	pstack->size = size;
	pstack->frame_base = 0;
	pstack->pvars = mlr_malloc_or_die(size * sizeof(mv_t));
	for (int i = 0; i < size; i++) {
		pstack->pvars[i] = mv_absent();
	}

	return pstack;
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
