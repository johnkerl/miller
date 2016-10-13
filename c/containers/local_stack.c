#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/local_stack.h"

// ----------------------------------------------------------------
local_stack_t* local_stack_alloc(int size) {
	local_stack_t* pstack = mlr_malloc_or_die(sizeof(local_stack_t));

	pstack->in_use = FALSE;
	pstack->size = size;
	pstack->pvars = mlr_malloc_or_die(size * sizeof(mv_t));
	for (int i = 0; i < size; i++) {
		pstack->pvars[i] = mv_absent();
	}

	return pstack;
}

// ----------------------------------------------------------------
void local_stack_free(local_stack_t* pstack) {
}

//// ----------------------------------------------------------------
//void local_stack_free(local_stack_t* pstack) {
//	if (pstack == NULL)
//		return;
//
//	local_stack_frame_free(pstack->pbase_frame);
//	free(pstack->ppframes);
//	free(pstack);
//}

//// ----------------------------------------------------------------
//static inline local_stack_frame_t* local_stack_frame_alloc(int fenced, int ephemeral) {
//	local_stack_frame_t* pframe = mlr_malloc_or_die(sizeof(local_stack_frame_t));
//	pframe->pbindings = lhmsmv_alloc();
//	pframe->fenced    = fenced;
//	pframe->ephemeral = ephemeral;
//	pframe->in_use    = FALSE;
//	return pframe;
//}
//
//local_stack_frame_t* local_stack_frame_alloc_unfenced() {
//	return local_stack_frame_alloc(FALSE, FALSE);
//}
//
//local_stack_frame_t* local_stack_frame_alloc_fenced() {
//	return local_stack_frame_alloc(TRUE, FALSE);
//}
//
//// ----------------------------------------------------------------
//void local_stack_frame_free(local_stack_frame_t* pframe) {
//	lhmsmv_free(pframe->pbindings);
//	free(pframe);
//}
//
//// ----------------------------------------------------------------
//// xxx cmt
//local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe) {
//	if (pframe->in_use) {
//		local_stack_frame_t* pephemeral = local_stack_frame_alloc(pframe->fenced, TRUE);
//		pephemeral->in_use = TRUE;
//		return pephemeral;
//	} else {
//		pframe->in_use = TRUE;
//		return pframe;
//	}
//}
//
//void local_stack_frame_exit(local_stack_frame_t* pframe) {
//	if (pframe->ephemeral) {
//		local_stack_frame_free(pframe);
//	} else {
//		lhmsmv_clear(pframe->pbindings);
//		pframe->in_use = FALSE;
//	}
//}
//
//// ----------------------------------------------------------------
//void local_stack_push(local_stack_t* pstack, local_stack_frame_t* pframe) {
//	if (pstack->num_used >= pstack->num_allocated) {
//		pstack->num_allocated += INITIAL_SIZE;
//		pstack->ppframes = mlr_realloc_or_die(pstack->ppframes,
//			pstack->num_allocated * sizeof(local_stack_frame_t*));
//	}
//	pstack->ppframes[pstack->num_used] = pframe;
//	pstack->num_used++;
//}
//
//// ----------------------------------------------------------------
//local_stack_frame_t* local_stack_pop(local_stack_t* pstack) {
//	if (pstack->num_used <= 0) {
//		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
//		exit(1);
//	}
//
//	local_stack_frame_t* pframe = pstack->ppframes[pstack->num_used-1];
//	pstack->num_used--;
//	return pframe;
//}
//
//// ----------------------------------------------------------------
//// xxx cmt
//mv_t* local_stack_resolve(local_stack_t* pstack, char* name) {
//	for (int i = pstack->num_used - 1; i >= 0; i--) {
//		mv_t* pval = lhmsmv_get(pstack->ppframes[i]->pbindings, name);
//		if (pval != NULL) {
//			return pval;
//		}
//		if (pstack->ppframes[i]->fenced) {
//			break;
//		}
//	}
//	return NULL;
//}
//
//// ----------------------------------------------------------------
//// xxx
////
//// run_mlr --opprint --from $indir/abixy put '
////     func f(x) {
////         local a = 1;
////         if (NR > 5) {
////             a = 2;
////         }
////         return a;
////     }
////
////     func g(x) {
////         local b = 1;
////         if (NR > 5) {
////             local b = 2;
////         }
////         return b;
////     }
////     $of = f(NR);
////     $og = g(NR);
//// '
//
//void local_stack_define(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags) {
//	local_stack_frame_t* ptop_frame = pstack->ppframes[pstack->num_used - 1];
//	lhmsmv_put(ptop_frame->pbindings, name, pmv, free_flags);
//}
//
//void local_stack_set(local_stack_t* pstack, char* name, mv_t* pmv, char free_flags) {
//
//	// xxx comment
//	// xxx test thoroughly
//	for (int i = pstack->num_used - 1; i >= 0; i--) {
//		local_stack_frame_t* pframe = pstack->ppframes[i];
//		if (lhmsmv_get(pframe->pbindings, name)) {
//			lhmsmv_put(pframe->pbindings, name, pmv, free_flags);
//			return;
//		}
//		if (pstack->ppframes[i]->fenced) {
//			break;
//		}
//	}
//
//	local_stack_frame_t* ptop_frame = pstack->ppframes[pstack->num_used - 1];
//	lhmsmv_put(ptop_frame->pbindings, name, pmv, free_flags);
//}
//
//// ----------------------------------------------------------------
//void local_stack_clear(local_stack_t* pstack) {
//	if (pstack->num_used <= 0) {
//		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
//		exit(1);
//	}
//	lhmsmv_clear(pstack->ppframes[pstack->num_used-1]->pbindings);
//}
//
//// ----------------------------------------------------------------
//void local_stack_print(local_stack_t* pstack) {
//	printf("BIND STACK BEGIN (#frames %d):\n", pstack->num_used);
//	for (int i = pstack->num_used - 1; i >= 0; i--) {
//		printf("-- FRAME %d (fenced:%d):\n", i, pstack->ppframes[i]->fenced);
//		lhmsmv_dump(pstack->ppframes[i]->pbindings);
//	}
//	printf("BIND STACK END\n");
//}
//
