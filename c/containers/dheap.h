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
	char  is_malloced; // xxx bool_t
	double *elements;
} dheap_t;

dheap_t *dheap_alloc();
dheap_t *dheap_from_array(double *array, int n);
void dheap_free(dheap_t *pheap);

void dheap_add(dheap_t *pdheap, double v);
double dheap_remove(dheap_t *pdheap);

// For debug
void dheap_print(dheap_t *pdheap);
// For unit test
int dheap_check(dheap_t *pdheap, char *file, int line);

#endif // DHEAP_H
