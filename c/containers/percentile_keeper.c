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
	ppercentile_keeper->data     = mlr_malloc_or_die(capacity*sizeof(mv_t));
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
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, mv_t value) {
	if (ppercentile_keeper->size >= ppercentile_keeper->capacity) {
		ppercentile_keeper->capacity = (int)(ppercentile_keeper->capacity * GROWTH_FACTOR);
		ppercentile_keeper->data = (mv_t*)mlr_realloc_or_die(ppercentile_keeper->data,
			ppercentile_keeper->capacity*sizeof(mv_t));
	}
	ppercentile_keeper->data[ppercentile_keeper->size++] = value;
	ppercentile_keeper->sorted = FALSE;
}

// ----------------------------------------------------------------
// https://en.wikipedia.org/wiki/Percentile

static int compute_index(int n, double p) {
	int index = p*n/100.0;
	if (index < 0)
		index = 0;
	else if (index >= n)
		index = n-1;
	return index;
}

static int compute_index_nearest_rank(int n, double p) {
	int index = (int)(ceil((p/100.0)*n)) - 1;
	if (index < 0)
		index = 0;
	else if (index >= n)
		index = n-1;
	return index;
}

static mv_t get_percentile_linearly_interpolated(mv_t* array, int n, double p) {
	double findex = (p/100.0)*n - 1.0;
	int iindex = (int)floor(findex);
	// xxx make corner-case UTs
	if (iindex < 0)
		iindex = 0;
	else if (iindex >= n)
		iindex = n-1;
	if (iindex >= n-1) {
		return array[iindex];
	} else {
		// array[iindex] + frac * (array[iindex+1] - array[iindex]);
		mv_t frac = mv_from_float(findex - iindex);
		printf("n=%d findex=%lf iiindex=%d frac=%lf\n", n, findex, iindex, findex-iindex);
		mv_t* pa = &array[iindex];
		mv_t* pb = &array[iindex+1];
		mv_t diff = x_xx_minus_func(pa, pb);
		mv_t prod = x_xx_times_func(&frac, &diff);
		mv_t rv = x_xx_plus_func(pa, &prod);
		return rv;
	}
}

// ----------------------------------------------------------------
mv_t percentile_keeper_emit(percentile_keeper_t* ppercentile_keeper, double percentile) {
	if (!ppercentile_keeper->sorted) {
		qsort(ppercentile_keeper->data, ppercentile_keeper->size, sizeof(mv_t), mv_nn_comparator);
		ppercentile_keeper->sorted = TRUE;
	}
	return ppercentile_keeper->data[compute_index(ppercentile_keeper->size, percentile)];
}

mv_t percentile_keeper_emit_nearest_rank(percentile_keeper_t* ppercentile_keeper, double percentile) {
	if (!ppercentile_keeper->sorted) {
		qsort(ppercentile_keeper->data, ppercentile_keeper->size, sizeof(mv_t), mv_nn_comparator);
		ppercentile_keeper->sorted = TRUE;
	}
	return ppercentile_keeper->data[compute_index_nearest_rank(ppercentile_keeper->size, percentile)];
}

mv_t percentile_keeper_emit_linearly_interpolated(percentile_keeper_t* ppercentile_keeper, double percentile) {
	if (!ppercentile_keeper->sorted) {
		qsort(ppercentile_keeper->data, ppercentile_keeper->size, sizeof(mv_t), mv_nn_comparator);
		ppercentile_keeper->sorted = TRUE;
	}
	return get_percentile_linearly_interpolated(ppercentile_keeper->data, ppercentile_keeper->size, percentile);
}

// ----------------------------------------------------------------
void percentile_keeper_print(percentile_keeper_t* ppercentile_keeper) {
	printf("percentile_keeper dump:\n");
	for (int i = 0; i < ppercentile_keeper->size; i++) {
		mv_t* pa = &ppercentile_keeper->data[i];
		if (pa->type == MT_FLOAT)
			printf("[%02d] %.8lf\n", i, ppercentile_keeper->data[i].u.fltv);
		else
			printf("[%02d] %8lld\n", i, ppercentile_keeper->data[i].u.intv);
	}
}
