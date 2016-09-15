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

	pstack->ppframes = mlr_malloc_or_die(pstack->num_allocated * sizeof(bind_stack_frame_t*));
	memset(pstack->ppframes, 0, pstack->num_allocated * sizeof(bind_stack_frame_t*));

	pstack->pbase_frame = bind_stack_frame_alloc_unfenced();
	bind_stack_push(pstack, pstack->pbase_frame);

	return pstack;
}

// ----------------------------------------------------------------
void bind_stack_free(bind_stack_t* pstack) {
	if (pstack == NULL)
		return;

	bind_stack_frame_free(pstack->pbase_frame);
	free(pstack->ppframes);
	free(pstack);
}

// ----------------------------------------------------------------
bind_stack_frame_t* bind_stack_frame_alloc_unfenced() {
	bind_stack_frame_t* pframe = mlr_malloc_or_die(sizeof(bind_stack_frame_t));
	pframe->pbindings = lhmsmv_alloc();
	pframe->fenced = FALSE;
	return pframe;
}

// ----------------------------------------------------------------
bind_stack_frame_t* bind_stack_frame_alloc_fenced() {
	bind_stack_frame_t* pframe = mlr_malloc_or_die(sizeof(bind_stack_frame_t));
	pframe->pbindings = lhmsmv_alloc();
	pframe->fenced = TRUE;
	return pframe;
}

// ----------------------------------------------------------------
void bind_stack_frame_free(bind_stack_frame_t* pframe) {
	lhmsmv_free(pframe->pbindings);
	free(pframe);
}

// ----------------------------------------------------------------
void bind_stack_push(bind_stack_t* pstack, bind_stack_frame_t* pframe) {
	if (pstack->num_used >= pstack->num_allocated) {
		pstack->num_allocated += INITIAL_SIZE;
		pstack->ppframes = mlr_realloc_or_die(pstack->ppframes,
			pstack->num_allocated * sizeof(bind_stack_frame_t*));
	}
	pstack->ppframes[pstack->num_used] = pframe;
	pstack->num_used++;
}

// ----------------------------------------------------------------
void bind_stack_pop(bind_stack_t* pstack) {
	if (pstack->num_used <= 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	lhmsmv_clear(pstack->ppframes[pstack->num_used-1]->pbindings);
	pstack->num_used--;
}

// ----------------------------------------------------------------
mv_t* bind_stack_resolve(bind_stack_t* pstack, char* key) {
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		mv_t* pval = lhmsmv_get(pstack->ppframes[i]->pbindings, key);
		if (pval != NULL)
			return pval;
		if (pstack->ppframes[i]->fenced)
			break;
	}
	return NULL;
}

// ----------------------------------------------------------------
void bind_stack_set(bind_stack_t* pstack, char* name, mv_t* pmv) {
	bind_stack_frame_t* ptop_frame = pstack->ppframes[pstack->num_used - 1];
	lhmsmv_put(ptop_frame->pbindings, name, pmv, FREE_ENTRY_VALUE);
}

// ----------------------------------------------------------------
void bind_stack_clear(bind_stack_t* pstack) {
	if (pstack->num_used <= 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	lhmsmv_clear(pstack->ppframes[pstack->num_used-1]->pbindings);
}

// ----------------------------------------------------------------
void bind_stack_print(bind_stack_t* pstack) {
	printf("BIND STACK BEGIN (#frames %d):\n", pstack->num_used);
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		printf("-- FRAME %d (fenced:%d):\n", i, pstack->ppframes[i]->fenced);
		lhmsmv_dump(pstack->ppframes[i]->pbindings);
	}
	printf("BIND STACK END\n");
}
