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
typedef int mv_comparator_func_t(const mv_t* pa, const mv_t* pb);
static int mv_ff_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.fltv - pb->u.fltv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_fi_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.fltv - pb->u.intv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_if_comparator(const mv_t* pa, const mv_t* pb) {
	double d = pa->u.intv - pb->u.fltv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
static int mv_ii_comparator(const mv_t* pa, const mv_t* pb) {
	long long d = pa->u.intv - pb->u.intv;
	return (d < 0) ? -1 : (d > 0) ? 1 : 0;
}
// We assume mv_t's coming into percentile keeper are int or double -- in particular, non-null.
static mv_comparator_func_t* mv_comparator_dispositions[MT_MAX][MT_MAX] = {
	//         NULL   ERROR BOOL  FLOAT             INT               STRING
	/*NULL*/   {NULL, NULL, NULL, NULL,             NULL,             NULL},
	/*ERROR*/  {NULL, NULL, NULL, NULL,             NULL,             NULL},
	/*BOOL*/   {NULL, NULL, NULL, NULL,             NULL,             NULL},
	/*FLOAT*/  {NULL, NULL, NULL, mv_ff_comparator, mv_fi_comparator, NULL},
	/*INT*/    {NULL, NULL, NULL, mv_if_comparator, mv_ii_comparator, NULL},
	/*STRING*/ {NULL, NULL, NULL, NULL,             NULL,             NULL},
};
static int mv_comparator(const void* pva, const void* pvb) {
	const mv_t* pa = pva;
	const mv_t* pb = pvb;
	return mv_comparator_dispositions[pa->type][pb->type](pa, pb);
}

static int compute_index(int n, double p) {
	int index = p*n/100.0;
	if (index < 0)
		index = 0;
	else if (index >= n)
		index = n-1;
	return index;
}

// See also https://github.com/johnkerl/miller/issues/14 which requests an interpolation option.
mv_t percentile_keeper_emit(percentile_keeper_t* ppercentile_keeper, double percentile) {
	if (!ppercentile_keeper->sorted) {
		qsort(ppercentile_keeper->data, ppercentile_keeper->size, sizeof(mv_t), mv_comparator);
		ppercentile_keeper->sorted = TRUE;
	}
	return ppercentile_keeper->data[compute_index(ppercentile_keeper->size, percentile)];
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
