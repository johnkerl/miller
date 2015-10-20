#include <string.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/dvector.h"

// ----------------------------------------------------------------
dvector_t* dvector_alloc(int initial_capacity) {
	int capacity = initial_capacity;
	dvector_t* pdvector = mlr_malloc_or_die(sizeof(dvector_t));
	pdvector->data     = mlr_malloc_or_die(capacity*sizeof(double));
	pdvector->size     = 0;
	pdvector->capacity = capacity;
	return pdvector;
}

// ----------------------------------------------------------------
void dvector_free(dvector_t* pdvector) {
	if (pdvector == NULL)
		return;
	free(pdvector->data);
	pdvector->data = NULL;
	pdvector->size = 0;
	pdvector->capacity = 0;
	free(pdvector);
}

void dvector_append(dvector_t* pdvector, double value) {
	if (pdvector->size >= pdvector->capacity) {
		pdvector->capacity = (int)(pdvector->capacity * 2);
		pdvector->data = (double*)mlr_realloc_or_die(pdvector->data,
			pdvector->capacity*sizeof(double));
	}
	pdvector->data[pdvector->size++] = value;
}
