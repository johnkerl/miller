#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "loop_stack.h"

#define INITIAL_SIZE 32

// Example states:

// num_allocated = 4         num_allocated = 4         num_allocated = 4         num_allocated = 4
// num_used = 1              num_used = 2              num_used = 3              num_used = 4
// num_used_minus_one = 0    num_used_minus_one = 1    num_used_minus_one = 2    num_used_minus_one = 3
//
// +---+                     +---+                     +---+                     +---+
// | 2 | 0  <--- top         | 2 | 0                   | 2 | 0                   | 2 | 0
// +---+                     +---+                     +---+                     +---+
// |///| 1                   | 0 | 1  <--- top         | 0 | 1                   | 0 |
// +---+                     +---+                     +---+                     +---+
// |///| 2                   |///| 2                   | 4 | 2  <--- top         | 4 |
// +---+                     +---+                     +---+                     +---+
// |///| 3                   |///| 3                   |///| 3                   | 6 | 3  <--- top
// +---+                     +---+                     +---+                     +---+

// ----------------------------------------------------------------
loop_stack_t* loop_stack_alloc() {
	loop_stack_t* pstack = mlr_malloc_or_die(sizeof(loop_stack_t));

	// Guard zone of one. As noted in the header file, set/get are intentionally not bounds-checked.
	// If set is called without push, or after final pop, we can at least not corrupt other code.
	pstack->num_used_minus_one = 0;
	pstack->num_allocated = INITIAL_SIZE;

	pstack->pframes = mlr_malloc_or_die(pstack->num_allocated * sizeof(int));
	memset(pstack->pframes, 0, pstack->num_allocated * sizeof(int));

	return pstack;
}

// ----------------------------------------------------------------
void loop_stack_free(loop_stack_t* pstack) {
	if (pstack == NULL)
		return;
	free(pstack->pframes);
	free(pstack);
}

// ----------------------------------------------------------------
void loop_stack_push(loop_stack_t* pstack) {
	if (pstack->num_used_minus_one >= pstack->num_allocated - 1) {
		pstack->num_allocated += INITIAL_SIZE;
		pstack->pframes = mlr_realloc_or_die(pstack->pframes, pstack->num_allocated);
	}
	pstack->num_used_minus_one++;
	pstack->pframes[pstack->num_used_minus_one] = 0;
}

// ----------------------------------------------------------------
int loop_stack_pop(loop_stack_t* pstack) {
	if (pstack->num_used_minus_one <= 0) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	int rv = pstack->pframes[pstack->num_used_minus_one];
	pstack->num_used_minus_one--;
	return rv;
}

// ----------------------------------------------------------------
// Not bounds-checked, as noted in the header file.
void loop_stack_set(loop_stack_t* pstack, int bits) {
	pstack->pframes[pstack->num_used_minus_one] |= bits;
}
void loop_stack_clear(loop_stack_t* pstack, int bits) {
	pstack->pframes[pstack->num_used_minus_one] &= ~bits;
}

// ----------------------------------------------------------------
// Not bounds-checked, as noted in the header file.
int loop_stack_get(loop_stack_t* pstack) {
	return pstack->pframes[pstack->num_used_minus_one];
}
