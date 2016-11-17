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
	pframe->pvars = mlr_malloc_or_die(size * sizeof(local_stack_frame_entry_t));
	for (int i = 0; i < size; i++) {
		local_stack_frame_entry_t* pentry = &pframe->pvars[i];
		pentry->value.mlrval = mv_absent(); // xxx make method
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
		mv_free(&pframe->pvars[i].value.mlrval); // xxx temp
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
	MLR_INTERNAL_CODING_ERROR_UNLESS(mv_is_absent(&pframe->pvars[0].value.mlrval)); // xxx temp
	for (int i = 0; i < pframe->size; i++)
		mv_free(&pframe->pvars[i].value.mlrval); // xxx temp
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
mv_t local_stack_frame_get_scalar_from_indexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, sllmv_t* pmvkeys)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p GET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "GET", FALSE, vardef_frame_relative_index);

	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];
	mlhmmv_value_t* pmvalue = &pentry->value;
	// xxx encapsulate
	if (pmvalue->is_terminal) {
		return mv_absent();
	} else if (pmvalue->pnext_level == NULL) {
		LOCAL_STACK_TRACE(printf("VALUE IS EMPTY\n"));
		return mv_absent();
	} else {
		int error = 0;
		LOCAL_STACK_TRACE(printf("VALUE IS:\n"));
		LOCAL_STACK_TRACE(printf("PTR IS %p\n", pmvalue->pnext_level));
		LOCAL_STACK_TRACE(mlhmmv_level_print_stacked(pmvalue->pnext_level, 0, TRUE, TRUE, "", stdout));
		mv_t* pval = mlhmmv_get_terminal_from_level(pmvalue->pnext_level, pmvkeys, &error);
		if (pval == NULL) {
			return mv_absent();
		} else {
			return mv_copy(pval); // xxx temp copy?
		}
	}
}

// ----------------------------------------------------------------
mlhmmv_value_t* local_stack_frame_get_extended_from_indexed(local_stack_frame_t* pframe, // xxx rename w/ 'reference' in name
	int vardef_frame_relative_index, sllmv_t* pmvkeys)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p GET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "GET", FALSE, vardef_frame_relative_index);

	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];
	mlhmmv_value_t* pmvalue = &pentry->value;
	if (pmvalue->is_terminal && pmvalue->pnext_level == NULL) {
		LOCAL_STACK_TRACE(printf("VALUE IS EMPTY\n"));
		return NULL;
	} else {
		int error = 0;
		LOCAL_STACK_TRACE(printf("VALUE IS:\n"));
		LOCAL_STACK_TRACE(printf("PTR IS %p\n", pmvalue->pnext_level));
		LOCAL_STACK_TRACE(mlhmmv_level_print_stacked(pmvalue->pnext_level, 0, TRUE, TRUE, "", stdout));

		// Maybe null
		if (pmvkeys == NULL || pmvkeys->length == 0) {
			return pmvalue;
		} else {
			return mlhmmv_get_value_from_level(pmvalue->pnext_level, pmvkeys, &error);
		}
	}
}

// ----------------------------------------------------------------
void local_stack_frame_assign_scalar_indexed(local_stack_frame_t* pframe,
	int vardef_frame_relative_index, sllmv_t* pmvkeys,
	mv_t terminal_value)
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "ASSIGN", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	if (!(TYPE_MASK_MAP & pentry->type_mask)) {
		local_stack_frame_throw_type_mismatch(pentry, &terminal_value);
	}

	mlhmmv_value_t* pmvalue = &pentry->value;

	// xxx encapsulate
	if (pmvalue->is_terminal) {
		mv_free(&pmvalue->mlrval);
		pmvalue->is_terminal = FALSE;
		pmvalue->pnext_level = mlhmmv_level_alloc();
	}
	mlhmmv_put_terminal_from_level(pmvalue->pnext_level, pmvkeys->phead, &terminal_value);

	LOCAL_STACK_TRACE(printf("VALUE IS:\n"));
	LOCAL_STACK_TRACE(mlhmmv_level_print_stacked(pmvalue->pnext_level, 0, TRUE, TRUE, "", stdout));
}

void local_stack_frame_assign_extended_indexed(local_stack_frame_t* pframe, // xxx rename
	int vardef_frame_relative_index, sllmv_t* pmvkeys,
	mlhmmv_value_t new_value) // xxx by ptr
{
	LOCAL_STACK_TRACE(printf("LOCAL STACK FRAME %p SET %d\n", pframe, vardef_frame_relative_index));
	LOCAL_STACK_BOUNDS_CHECK(pframe, "ASSIGN", TRUE, vardef_frame_relative_index);
	local_stack_frame_entry_t* pentry = &pframe->pvars[vardef_frame_relative_index];

	if (!(TYPE_MASK_MAP & pentry->type_mask)) {
		local_stack_frame_throw_type_xmismatch(pentry, &new_value);
	}

	mlhmmv_value_t* pmvalue = &pentry->value;

	// xxx encapsulate
	if (pmvalue->is_terminal) {
		mv_free(&pmvalue->mlrval);
		pmvalue->is_terminal = FALSE;
		pmvalue->pnext_level = mlhmmv_level_alloc();
	}
	mlhmmv_put_value_at_level_aux(pmvalue->pnext_level, pmvkeys->phead, &new_value);

	LOCAL_STACK_TRACE(printf("VALUE IS:\n"));
	LOCAL_STACK_TRACE(mlhmmv_level_print_stacked(pmvalue->pnext_level, 0, TRUE, TRUE, "", stdout));
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
void local_stack_frame_throw_type_mismatch(local_stack_frame_entry_t* pentry, mv_t* pval) {
	MLR_INTERNAL_CODING_ERROR_IF(pentry->name == NULL);
	char* sval = mv_alloc_format_val_quoting_strings(pval); // xxx temp
	fprintf(stderr, "%s: %s type assertion for variable %s unmet by value %s with type %s.\n",
		MLR_GLOBALS.bargv0, type_mask_to_desc(pentry->type_mask), pentry->name,
		sval, mt_describe_type_simple(pval->type)); // xxx temp
	free(sval);
	exit(1);
}

void local_stack_frame_throw_type_xmismatch(local_stack_frame_entry_t* pentry, mlhmmv_value_t* pxval) {
	MLR_INTERNAL_CODING_ERROR_IF(pentry->name == NULL);
	char* sval = mv_alloc_format_val_quoting_strings(&pxval->mlrval); // xxx temp
	fprintf(stderr, "%s: %s type assertion for variable %s unmet by value %s with type %s.\n",
		MLR_GLOBALS.bargv0, type_mask_to_desc(pentry->type_mask), pentry->name,
		sval, mt_describe_type_simple(pxval->mlrval.type)); // xxx temp -- needs xtype
	free(sval);
	exit(1);
}
