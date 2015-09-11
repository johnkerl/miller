// ================================================================
// Zero-indexed max-heap of double.
// John Kerl
// 2012-06-02
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/dheap.h"

// ----------------------------------------------------------------
// 1-up: left  child 2*i
//       right child 2*i+1
//       parent      i/2
// 0-up: left  child 2*i+1
//       right child 2*i+2
//       parent      (i-1)/2
// Why: Example of 1-up i=10 l=20 r=21; 0-up i=9 l=19 r=20.
// Or: 0-up i |-> 1-up i+1 |-> 1-up 2*i+2 |-> 0-up 2*i+1.
//   And likewise for right child & parent.

static inline int dheap_left_child_index(int i, int n)
{
	int li = 2*i+1;
	if (li >= n)
		return -1;
	else
		return li;
}

static inline int dheap_right_child_index(int i, int n)
{
	int ri = 2*i+2;
	if (ri >= n)
		return -1;
	else
		return ri;
}

static inline int dheap_parent_index(int i, int n)
{
	if (i == 0)
		return -1;
	else
		return (i-1)/2;
}

static inline void ptr_swap(double *pa, double *pb)
{
	double temp = *pa;
	*pa = *pb;
	*pb = temp;
}

// ================================================================
dheap_t *dheap_alloc()
{
	dheap_t *pdheap = mlr_malloc_or_die(sizeof(dheap_t));
	pdheap->n = 0;
	pdheap->alloc_size = DHEAP_INIT_ALLOC_SIZE;
	pdheap->elements = mlr_malloc_or_die(pdheap->alloc_size*sizeof(double));
	pdheap->is_malloced = 1;
	return pdheap;
}

// ----------------------------------------------------------------
dheap_t *dheap_from_array(double *array, int n)
{
	dheap_t *pdheap = mlr_malloc_or_die(sizeof(dheap_t));
	pdheap->n = 0;
	pdheap->alloc_size = n;
	pdheap->elements = array;
	pdheap->is_malloced = 0;
	return pdheap;
}

// ----------------------------------------------------------------
void dheap_free(dheap_t *pdheap)
{
	if (pdheap == NULL)
		return;
	if (pdheap->elements != NULL)
		if (pdheap->is_malloced)
			free(pdheap->elements);
	pdheap->n = 0;
	pdheap->alloc_size = 0;
	pdheap->elements = NULL;
	free(pdheap);
}

// ================================================================
static void dheap_print_aux(dheap_t *pdheap, int i, int depth)
{
	if (i >= pdheap->n)
		return;
	int w;
	printf("[%04d] ", i);
	for (w = 0; w < depth; w++)
		printf("     ");
	printf("%.8lf\n", pdheap->elements[i]);
	int li = dheap_left_child_index (i, pdheap->n);
	int ri = dheap_right_child_index(i, pdheap->n);
	if (li != -1)
		dheap_print_aux(pdheap, li, depth+1);
	if (ri != -1)
		dheap_print_aux(pdheap, ri, depth+1);
}

void dheap_print(dheap_t *pdheap)
{
	printf("BEGIN DHEAP (n=%d):\n", pdheap->n);
	dheap_print_aux(pdheap, 0, 0);
	printf("END DHEAP\n");
}

// ----------------------------------------------------------------
//          1
//    2           3
//  4    5     6     7
// 8 9 10 11 12 13 14 15

static int dheap_check_aux(dheap_t *pdheap, int i, char *file, int line)
{
	int n = pdheap->n;
	double *pe = pdheap->elements;

	if (i >= n)
		return TRUE;
	int li = dheap_left_child_index (i, pdheap->n);
	int ri = dheap_right_child_index(i, pdheap->n);
	if (li != -1) {
		if (pe[i] < pe[li]) {
			fprintf(stderr, "dheap check fail %s:%d pe[%d]=%lf < pe[%d]=%lf\n",
				file, line, i, pe[i], li, pe[li]);
			return FALSE;
		}
		dheap_check_aux(pdheap, li, file, line);
	}
	if (ri != -1) {
		if (pe[i] < pe[ri]) {
			fprintf(stderr, "dheap check fail %s:%d pe[%d]=%lf < pe[%d]=%lf\n",
				file, line, i, pe[i], ri, pe[ri]);
			return FALSE;
		}
		dheap_check_aux(pdheap, ri, file, line);
	}
	return TRUE;
}

int dheap_check(dheap_t *pdheap, char *file, int line)
{
	return dheap_check_aux(pdheap, 1, file, line);
}

// ----------------------------------------------------------------
static void dheap_bubble_up(dheap_t *pdheap, int i)
{
	int pi = dheap_parent_index(i, pdheap->n);
	if (pi == -1)
		return;
	double *pe = pdheap->elements;
	if (pe[pi] < pe[i]) {
		ptr_swap(&pe[pi], &pe[i]);
		dheap_bubble_up(pdheap, pi);
	}
}

// ----------------------------------------------------------------
//          1
//    2           3
//  4    5     6     7
// 8 9 10 11 12         (n=13)

void dheap_add(dheap_t *pdheap, double v)
{
	if (pdheap->n >= pdheap->alloc_size) {
		if (!pdheap->is_malloced) {
			fprintf(stderr, "extension of non-malloced dheap!\n");
			exit(1);
		}
		pdheap->alloc_size *= 2;
		pdheap->elements = (double *)realloc(pdheap->elements,
			pdheap->alloc_size*sizeof(double));
	}

	pdheap->elements[pdheap->n++] = v;
	dheap_bubble_up(pdheap, pdheap->n-1);
}

// ----------------------------------------------------------------
// 1. Replace the root of the dheap with the last element on the last level.
// 2. Compare the new root with its children; if they are in the correct order,
//    stop.
// 3. If not, swap the element with one of its children and return to the
//    previous step. (Swap with its smaller child in a min-dheap and its larger
//    child in a max-dheap.)

static void dheap_bubble_down(dheap_t *pdheap, int i)
{
	int li = dheap_left_child_index(i, pdheap->n);

	if (li == -1) {
		// We add left to right, so this means left and right are both nil.
		return;
	}

	int ri = dheap_right_child_index(i, pdheap->n);
	double *pe = pdheap->elements;

	if (ri == -1) {
		// Right is nil, left is non-nil.
		if (pe[li] > pe[i]) {
			ptr_swap(&pe[li], &pe[i]);
			dheap_bubble_down(pdheap, li);
		}
		return;
	}

	// Now left and right are both non-nil.
	//
	//    P             3
	//  L   R         9   7
	// a b c d       1 2 4 6
	//
	// Cases:
	double *L = &pe[li];
	double *P = &pe[i];
	double *R = &pe[ri];

	if (*L <= *P) {
		if (*R <= *P) {
			// 1. L <= R <= P:  done.
			// 2. R <= L <= P:  done.
			return;
		}
		else if (*P <= *R) {
			// 3. L <= P <= R:  swap P&R; bubble down R.
			ptr_swap(R, P);
			dheap_bubble_down(pdheap, ri);
		}
	}

	else if (*R <= *P && *P <= *L) {
		// 4. R <= P <= L:  swap P&L; bubble down L.
		ptr_swap(L, P);
		dheap_bubble_down(pdheap, li);
	}
	else if (*P <= *R && *R <= *L) {
		// 5. P <= R <= L:  swap P&L; bubble down L.
		ptr_swap(L, P);
		dheap_bubble_down(pdheap, li);
	}
	else if (P <= L && L <= R) {
		// 6. P <= L <= R:  swap P&R; bubble down R.
		ptr_swap(R, P);
		dheap_bubble_down(pdheap, ri);
	}
}

double dheap_remove(dheap_t *pdheap)
{
	if (pdheap->n <= 0) {
		fprintf(stderr, "remove from empty dheap!\n");
		exit(1);
	}

	double rv = pdheap->elements[0];

	pdheap->elements[0] = pdheap->elements[pdheap->n-1];
	pdheap->n--;
	dheap_bubble_down(pdheap, 0);

	return rv;
}

void dheap_sort(double *array, int n)
{
	dheap_t *pdheap = dheap_from_array(array, n);
	int i;

	for (i = 0; i < n; i++)
		dheap_add(pdheap, pdheap->elements[i]);

	for (i = n-1; i >= 0; i--)
		pdheap->elements[i] = dheap_remove(pdheap);

	dheap_free(pdheap);
}
