#include <string.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/percentile_keeper.h"

#define INITIAL_CAPACITY 10000
#define GROWTH_FACTOR    2.0

// ----------------------------------------------------------------
percentile_keeper_t* percentile_keeper_alloc() {
	int capacity = INITIAL_CAPACITY;
	percentile_keeper_t* ppercentile_keeper = mlr_malloc_or_die(sizeof(percentile_keeper_t));
	ppercentile_keeper->data     = mlr_malloc_or_die(capacity*sizeof(double));
	ppercentile_keeper->size     = 0;
	ppercentile_keeper->capacity = capacity;
	ppercentile_keeper->sorted   = FALSE;
	return ppercentile_keeper;
}

// ----------------------------------------------------------------
void percentile_keeper_free(percentile_keeper_t* ppercentile_keeper) {
	if (ppercentile_keeper == NULL)
		return;
	free(ppercentile_keeper->data);
	ppercentile_keeper->data = NULL;
	ppercentile_keeper->size = 0;
	ppercentile_keeper->capacity = 0;
	free(ppercentile_keeper);
}

// ----------------------------------------------------------------
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, double value) {
	if (ppercentile_keeper->size >= ppercentile_keeper->capacity) {
		ppercentile_keeper->capacity = (int)(ppercentile_keeper->capacity * GROWTH_FACTOR);
		ppercentile_keeper->data = (double*)mlr_realloc_or_die(ppercentile_keeper->data,
			ppercentile_keeper->capacity*sizeof(double));
	}
	ppercentile_keeper->data[ppercentile_keeper->size++] = value;
	ppercentile_keeper->sorted = FALSE;
}

// ----------------------------------------------------------------
static int double_comparator(const void* pva, const void* pvb) {
	const double* pa = pva;
	const double* pb = pvb;
	double d = *pa - *pb;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}

static int compute_index(int n, double p) {
	int index = p*n/100.0;
	if (index < 0)
		index = 0;
	else if (index >= n)
		index = n-1;
	// xxx need to try harder on round-up/round-down cases?
	return index;
}

double percentile_keeper_emit(percentile_keeper_t* ppercentile_keeper, double percentile) {
	if (!ppercentile_keeper->sorted) {
		qsort(ppercentile_keeper->data, ppercentile_keeper->size, sizeof(double), double_comparator);
		ppercentile_keeper->sorted = TRUE;
	}
	return ppercentile_keeper->data[compute_index(ppercentile_keeper->size, percentile)];
}

// ----------------------------------------------------------------
void percentile_keeper_print(percentile_keeper_t* ppercentile_keeper) {
	printf("percentile_keeper dump:\n");
	for (int i = 0; i < ppercentile_keeper->size; i++)
		printf("[%02d] %.8lf\n", i, ppercentile_keeper->data[i]);
}

