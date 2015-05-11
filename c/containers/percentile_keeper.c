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
	if (ppercentile_keeper->data != NULL) {
		free(ppercentile_keeper->data);
		ppercentile_keeper->data = NULL;
	}
	ppercentile_keeper->size = 0;
	ppercentile_keeper->capacity = 0;
	free(ppercentile_keeper);
}

// ----------------------------------------------------------------
void percentile_keeper_ingest(percentile_keeper_t* ppercentile_keeper, double value) {
	if (ppercentile_keeper->size >= ppercentile_keeper->capacity) {
		ppercentile_keeper->capacity = (int)(ppercentile_keeper->capacity * GROWTH_FACTOR);
		// xxx make realloc_or_die
		ppercentile_keeper->data = (double*)realloc(ppercentile_keeper->data,
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

// ================================================================
#ifdef __PERCENTILE_KEEPER_MAIN__
void percentile_keeper_dump(percentile_keeper_t* ppercentile_keeper) {
	for (int i = 0; i < ppercentile_keeper->size; i++)
		printf("[%02d] %.8lf\n", i, ppercentile_keeper->data[i]);
}
int main(int argc, char** argv) {
	char buffer[1024];
	percentile_keeper_t* ppercentile_keeper = percentile_keeper_alloc();
	char* line;
	while ((line = fgets(buffer, sizeof(buffer), stdin)) != NULL) {
		int len = strlen(line);
		if (len >= 1) // xxx write and use a chomp()
			if (line[len-1] == '\n')
				line[len-1] = 0;
		double v;
		if (sscanf(line, "%lf", &v) == 1) {
			percentile_keeper_ingest(ppercentile_keeper, v);
		} else {
			printf("meh? >>%s<<\n", line);
		}
	}
	percentile_keeper_dump(ppercentile_keeper);
	printf("\n");
	double p;
	p = 0.10; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
	p = 0.50; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
	p = 0.90; printf("%.2lf: %.6lf\n", p, percentile_keeper_emit(ppercentile_keeper, p));
	printf("\n");
	percentile_keeper_dump(ppercentile_keeper);
	return 0;
}
#endif // __PERCENTILE_KEEPER_MAIN__
