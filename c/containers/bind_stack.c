#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/bind_stack.h"

#define INITIAL_SIZE 32

bind_stack_t* bind_stack_alloc() {
	bind_stack_t* pstack = mlr_malloc_or_die(sizeof(bind_stack_t));

	pstack->num_used = 0;
	pstack->num_allocated = INITIAL_SIZE;

	pstack->pframes = mlr_malloc_or_die(pstack->num_allocated * sizeof(lhmsmv_t*));
	memset(pstack->pframes, 0, pstack->num_allocated * sizeof(lhmsmv_t*));

	return pstack;
}

void bind_stack_free(bind_stack_t* pstack) {
	if (pstack == NULL)
		return;

	free(pstack->pframes);
	free(pstack);
}

void bind_stack_push(bind_stack_t* pstack, lhmsmv_t* pframe) {
	if (pstack->num_used > pstack->num_allocated) {
		pstack->num_allocated += INITIAL_SIZE;
		pstack->pframes = mlr_realloc_or_die(pstack->pframes, pstack->num_allocated);
	}
	pstack->pframes[pstack->num_used] = pframe;
	pstack->num_used++;
}

void bind_stack_pop(bind_stack_t* pstack) {
	if (pstack->num_used <= 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
}

mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key) {
	for(int i = pstack->num_used - 1; i >= 0; i--) {
		mv_t* pval = lhmsmv_get(pstack->pframes[i], key);
		if (pval != NULL)
			return pval;
	}
	return NULL;
}

void bind_stack_print(bind_stack_t* pstack) {
	printf("BIND STACK BEGIN (#frames %d):\n", pstack->num_used);
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		printf("-- FRAME %d:\n", i);
		lhmsmv_dump(pstack->pframes[i]);
	}
	printf("BIND STACK END\n");
}
