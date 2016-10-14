#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/bind_stack.h"

#define INITIAL_SIZE 32

// ----------------------------------------------------------------
// Private to .c file
struct _bind_stack_frame_t {
	lhmsmv_t*  pbindings;
	char       fenced;
	char       ephemeral;
	char       in_use;
};

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
static inline bind_stack_frame_t* bind_stack_frame_alloc(int fenced, int ephemeral) {
	bind_stack_frame_t* pframe = mlr_malloc_or_die(sizeof(bind_stack_frame_t));
	pframe->pbindings = lhmsmv_alloc();
	pframe->fenced    = fenced;
	pframe->ephemeral = ephemeral;
	pframe->in_use    = FALSE;
	return pframe;
}

bind_stack_frame_t* bind_stack_frame_alloc_unfenced() {
	return bind_stack_frame_alloc(FALSE, FALSE);
}

bind_stack_frame_t* bind_stack_frame_alloc_fenced() {
	return bind_stack_frame_alloc(TRUE, FALSE);
}

// ----------------------------------------------------------------
void bind_stack_frame_free(bind_stack_frame_t* pframe) {
	lhmsmv_free(pframe->pbindings);
	free(pframe);
}

// ----------------------------------------------------------------
// xxx cmt
bind_stack_frame_t* bind_stack_frame_enter(bind_stack_frame_t* pframe) {
	if (pframe->in_use) {
		bind_stack_frame_t* pephemeral = bind_stack_frame_alloc(pframe->fenced, TRUE);
		pephemeral->in_use = TRUE;
		return pephemeral;
	} else {
		pframe->in_use = TRUE;
		return pframe;
	}
}

void bind_stack_frame_exit(bind_stack_frame_t* pframe) {
	if (pframe->ephemeral) {
		bind_stack_frame_free(pframe);
	} else {
		lhmsmv_clear(pframe->pbindings);
		pframe->in_use = FALSE;
	}
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
bind_stack_frame_t* bind_stack_pop(bind_stack_t* pstack) {
	MLR_INTERNAL_CODING_ERROR_IF(pstack->num_used <= 0);

	bind_stack_frame_t* pframe = pstack->ppframes[pstack->num_used-1];
	pstack->num_used--;
	return pframe;
}

// ----------------------------------------------------------------
// xxx cmt
mv_t* bind_stack_resolve(bind_stack_t* pstack, char* name) {
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		mv_t* pval = lhmsmv_get(pstack->ppframes[i]->pbindings, name);
		if (pval != NULL) {
			return pval;
		}
		if (pstack->ppframes[i]->fenced) {
			break;
		}
	}
	return NULL;
}

// ----------------------------------------------------------------
// xxx
//
// run_mlr --opprint --from $indir/abixy put '
//     func f(x) {
//         local a = 1;
//         if (NR > 5) {
//             a = 2;
//         }
//         return a;
//     }
//
//     func g(x) {
//         local b = 1;
//         if (NR > 5) {
//             local b = 2;
//         }
//         return b;
//     }
//     $of = f(NR);
//     $og = g(NR);
// '

void bind_stack_define(bind_stack_t* pstack, char* name, mv_t* pmv, char free_flags) {
	bind_stack_frame_t* ptop_frame = pstack->ppframes[pstack->num_used - 1];
	lhmsmv_put(ptop_frame->pbindings, name, pmv, free_flags);
}

void bind_stack_set(bind_stack_t* pstack, char* name, mv_t* pmv, char free_flags) {

	// xxx comment
	// xxx test thoroughly
	for (int i = pstack->num_used - 1; i >= 0; i--) {
		bind_stack_frame_t* pframe = pstack->ppframes[i];
		if (lhmsmv_get(pframe->pbindings, name)) {
			lhmsmv_put(pframe->pbindings, name, pmv, free_flags);
			return;
		}
		if (pstack->ppframes[i]->fenced) {
			break;
		}
	}

	bind_stack_frame_t* ptop_frame = pstack->ppframes[pstack->num_used - 1];
	lhmsmv_put(ptop_frame->pbindings, name, pmv, free_flags);
}

// ----------------------------------------------------------------
void bind_stack_clear(bind_stack_t* pstack) {
	MLR_INTERNAL_CODING_ERROR_IF(pstack->num_used <= 0);
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

