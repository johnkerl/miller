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
	pframe->pvars = mlr_malloc_or_die(size * sizeof(mlrval_and_type_mask_t));
	for (int i = 0; i < size; i++) {
		mlrval_and_type_mask_t* pentry = &pframe->pvars[i];
		pentry->mlrval = mv_absent();
		pentry->name = NULL;
		// Any type can be written here, unless otherwise specified by a typed definition
		pentry->type_mask = TYPE_MASK_ANY;
	}

	return pframe;
}

// ----------------------------------------------------------------
local_stack_frame_t* local_stack_frame_alloc(int size) {
	return _local_stack_alloc(size, FALSE);
}

// ----------------------------------------------------------------
void local_stack_frame_free(local_stack_frame_t* pframe) {
	if (pframe == NULL)
		return;
	for (int i = 0; i < pframe->size; i++) {
		mv_free(&pframe->pvars[i].mlrval);
	}
	free(pframe->pvars);
	free(pframe);
}

// ----------------------------------------------------------------
local_stack_frame_t* local_stack_frame_enter(local_stack_frame_t* pframe) {
	if (!pframe->in_use) {
		pframe->in_use = TRUE;
		LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME NON-EPH ENTER %p %d\n", pframe, pframe->size));
		return pframe;
	} else {
		local_stack_frame_t* prv = _local_stack_alloc(pframe->size, TRUE);
		LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME EPH ENTER %p/%p %d\n", pframe, prv, pframe->size));
		prv->in_use = TRUE;
		return prv;
	}
}

// ----------------------------------------------------------------
void local_stack_frame_exit (local_stack_frame_t* pframe) {
	MLR_INTERNAL_CODING_ERROR_UNLESS(mv_is_absent(&pframe->pvars[0].mlrval));
	for (int i = 0; i < pframe->size; i++)
		mv_free(&pframe->pvars[i].mlrval);
	if (!pframe->ephemeral) {
		pframe->in_use = FALSE;
		LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME NON-EPH EXIT %p %d\n", pframe, pframe->size));
	} else {
		local_stack_frame_free(pframe);
		LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME EPH EXIT %p %d\n", pframe, pframe->size));
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
	sllv_push(pstack->pframes, pframe);
}

local_stack_frame_t* local_stack_pop(local_stack_t* pstack) {
	return sllv_pop(pstack->pframes);
}

// ----------------------------------------------------------------
static int local_stack_bounds_check_announce_first_call = TRUE;

void local_stack_bounds_check(local_stack_frame_t* pframe, char* op, int set, int vardef_frame_relative_index) {
	if (local_stack_bounds_check_announce_first_call) {
		fprintf(stderr, "%s: local-stack bounds-checking is enabled\n", MLR_GLOBALS.bargv0);
		local_stack_bounds_check_announce_first_call = FALSE;
	}
	if (vardef_frame_relative_index < 0) {
		fprintf(stderr, "OP=%s FRAME=%p IDX=%d/%d STACK UNDERFLOW\n",
			op, pframe, vardef_frame_relative_index, pframe->size);
		exit(1);
	}
	if (set && vardef_frame_relative_index == 0) {
		fprintf(stderr, "OP=%s FRAME=%p IDX=%d/%d ABSENT WRITE\n",
			op, pframe, vardef_frame_relative_index, pframe->size);
		exit(1);
	}
	if (vardef_frame_relative_index >= pframe->size) {
		fprintf(stderr, "OP=%s FRAME=%p IDX=%d/%d STACK OVERFLOW\n",
			op, pframe, vardef_frame_relative_index, pframe->size);
		exit(1);
	}
}

// ----------------------------------------------------------------
void local_stack_frame_throw_type_mismatch(mlrval_and_type_mask_t* pentry, mv_t* pval) {
	MLR_INTERNAL_CODING_ERROR_IF(pentry->name == NULL);
	char* sval = mv_alloc_format_val(pval);
	fprintf(stderr, "%s: %s type assertion for variable %s unmet by value [%s] with type %s.\n",
		MLR_GLOBALS.bargv0, type_mask_to_desc(pentry->type_mask), pentry->name,
		sval, mt_describe_type_simple(pval->type));
	free(sval);
	exit(1);
}
