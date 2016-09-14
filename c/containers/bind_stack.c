#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/bind_stack.h"

#define INITIAL_SIZE 32

// ----------------------------------------------------------------
bind_stack_t* bind_stack_alloc() {
	bind_stack_t* pstack = mlr_malloc_or_die(sizeof(bind_stack_t));

	pstack->num_used = 0;
	pstack->num_allocated = INITIAL_SIZE;

	pstack->pframes = mlr_malloc_or_die(pstack->num_allocated * sizeof(bind_stack_frame_t));
	memset(pstack->pframes, 0, pstack->num_allocated * sizeof(bind_stack_frame_t));

	pstack->pbase_bindings = lhmsmv_alloc();
	bind_stack_push(pstack, pstack->pbase_bindings);

	return pstack;
}

// ----------------------------------------------------------------
void bind_stack_free(bind_stack_t* pstack) {
	if (pstack == NULL)
		return;

	lhmsmv_free(pstack->pbase_bindings);
	free(pstack->pframes);
	free(pstack);
}

// ----------------------------------------------------------------
static void bind_stack_push_fenced_aux(bind_stack_t* pstack, lhmsmv_t* bindings, int fenced) {
	if (pstack->num_used >= pstack->num_allocated) {
		pstack->num_allocated += INITIAL_SIZE;
		pstack->pframes = mlr_realloc_or_die(pstack->pframes,
			pstack->num_allocated * sizeof(bind_stack_frame_t));
	}
	pstack->pframes[pstack->num_used].bindings = bindings;
	pstack->pframes[pstack->num_used].fenced = fenced;
	pstack->num_used++;
}

void bind_stack_push(bind_stack_t* pstack, lhmsmv_t* bindings) {
	bind_stack_push_fenced_aux(pstack, bindings, FALSE);
}

void bind_stack_push_fenced(bind_stack_t* pstack, lhmsmv_t* bindings) {
	bind_stack_push_fenced_aux(pstack, bindings, TRUE);
}

// ----------------------------------------------------------------
void bind_stack_pop(bind_stack_t* pstack) {
	if (pstack->num_used <= 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	pstack->num_used--;
}

// ----------------------------------------------------------------
mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key) {
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		mv_t* pval = lhmsmv_get(pstack->pframes[i].bindings, key);
		if (pval != NULL)
			return pval;
		if (pstack->pframes[i].fenced)
			break;
	}
	return NULL;
}

// ----------------------------------------------------------------
void bind_stack_set(bind_stack_t* pstack, char* name, mv_t* pmv) {
	bind_stack_frame_t* ptop_frame = &pstack->pframes[pstack->num_used - 1];
	lhmsmv_put(ptop_frame->bindings, name, pmv, FREE_ENTRY_VALUE);
}

// ----------------------------------------------------------------
void bind_stack_print(bind_stack_t* pstack) {
	printf("BIND STACK BEGIN (#frames %d):\n", pstack->num_used);
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		printf("-- FRAME %d (fenced:%d):\n", i, pstack->pframes[i].fenced);
		lhmsmv_dump(pstack->pframes[i].bindings);
	}
	printf("BIND STACK END\n");
}
