// ================================================================
// Zero-indexed max-heap of doubles.
// John Kerl
// 2012-06-02
// ================================================================

#ifndef DHEAP_H
#define DHEAP_H

#define DHEAP_INIT_ALLOC_SIZE 1024 // Power of two
typedef struct _dheap_t {
	int  n;
	int  alloc_size;
	int  is_malloced;
	double *elements;
} dheap_t;

dheap_t *dheap_alloc();
dheap_t *dheap_from_array(double *array, int n);
void dheap_free(dheap_t *pheap);

void dheap_print(dheap_t *pdheap);
void dheap_check(dheap_t *pdheap, char *file, int line);

void dheap_add(dheap_t *pdheap, double v);
double dheap_remove(dheap_t *pdheap);

#endif // DHEAP_H
