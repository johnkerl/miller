#include <stdio.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/top_keeper.h"

// ----------------------------------------------------------------
top_keeper_t* top_keeper_alloc(int capacity) {
	top_keeper_t* ptop_keeper = mlr_malloc_or_die(sizeof(top_keeper_t));
	// xxx mk func for neg-cap check; use here & elsewhere
	ptop_keeper->top_values   = mlr_malloc_or_die(capacity*sizeof(double));
	ptop_keeper->top_precords = mlr_malloc_or_die(capacity*sizeof(lrec_t*));
	ptop_keeper->size         = 0;
	ptop_keeper->capacity     = capacity;
	return ptop_keeper;
}

// ----------------------------------------------------------------
void top_keeper_free(top_keeper_t* ptop_keeper) {
	if (ptop_keeper == NULL)
		return;
	free(ptop_keeper->top_values);
	free(ptop_keeper->top_precords);
	ptop_keeper->top_values = NULL;
	ptop_keeper->top_precords = NULL;
	ptop_keeper->size = 0;
	ptop_keeper->capacity = 0;
	free(ptop_keeper);
}

// ----------------------------------------------------------------
// Cases:
// 1. array size <  capacity
//    * find destidx
//    * if destidx == size
//        put it there
//      else
//        shift down & insert
//      increment size
//
// 2. array size == capacity
//    * find destidx
//    * if destidx == size
//        discard
//      else
//        shift down & insert

// capacity = 10, size = 6, destidx = 3     capacity = 10, size = 10, destidx = 3
// [0 #]   [0 #]                            [0 #]   [0 #]
// [1 #]   [1 #]                            [1 #]   [1 #]
// [2 #]   [2 #]                            [2 #]   [2 #]
// [3 #]*  [3 X]                            [3 #]*  [3 X]
// [4 #]   [4 #]                            [4 #]   [4 #]
// [5 #]   [5 #]                            [5 #]   [5 #]
// [6  ]   [6 #]                            [6 #]   [6 #]
// [7  ]   [7  ]                            [7 #]   [7 #]
// [8  ]   [8  ]                            [8 #]   [8 #]
// [9  ]   [9  ]                            [9 #]   [9 #]

void top_keeper_add(top_keeper_t* ptop_keeper, double value, lrec_t* prec) {
	int destidx = mlr_bsearch_double_for_insert(ptop_keeper->top_values, ptop_keeper->size, value);
	if (ptop_keeper->size < ptop_keeper->capacity) {
		for (int i = ptop_keeper->size-1; i >= destidx; i--) {
			ptop_keeper->top_values[i+1]   = ptop_keeper->top_values[i];
			ptop_keeper->top_precords[i+1] = ptop_keeper->top_precords[i];
		}
		ptop_keeper->top_values[destidx]   = value;
		ptop_keeper->top_precords[destidx] = prec;
		ptop_keeper->size++;
	} else {
		if (destidx >= ptop_keeper->capacity)
			return;
		for (int i = ptop_keeper->size-2; i >= destidx; i--) {
			ptop_keeper->top_values[i+1]   = ptop_keeper->top_values[i];
			ptop_keeper->top_precords[i+1] = ptop_keeper->top_precords[i];
		}
		ptop_keeper->top_values[destidx]   = value;
		ptop_keeper->top_precords[destidx] = prec; // xxx copy?? xxx free on shift-off?!?
	}
}

// ----------------------------------------------------------------
void top_keeper_print(top_keeper_t* ptop_keeper) {
	printf("top_keeper dump:\n");
	for (int i = 0; i < ptop_keeper->size; i++)
		printf("[%02d] %.8lf\n", i, ptop_keeper->top_values[i]);
	for (int i = ptop_keeper->size; i < ptop_keeper->capacity; i++)
		printf("[%02d] ---\n", i);
}
