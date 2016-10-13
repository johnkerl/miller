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
	if (!mv_is_absent(&pstack->pvars[0])) {
		fprintf(stderr, "%s: Internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	for (int i = 0; i < pstack->size; i++) {
		mv_free(&pstack->pvars[i]);
	}
	free(pstack->pvars);
	free(pstack);
}

// ----------------------------------------------------------------
void local_stack_frame_enter(local_stack_t* pstack, int count) {
	// xxx try to avoid with absent-read flag at stack-allocator ...
	mv_t* pframe = &pstack->pvars[pstack->frame_base];
	for (int i = 0; i < count; i++) {
		pframe[i] = mv_absent();
	}
	pstack->frame_base += count;
}

// ----------------------------------------------------------------
void local_stack_frame_exit(local_stack_t* pstack, int count) {
	pstack->frame_base -= count;
}
