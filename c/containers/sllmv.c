#include "lib/mlrutil.h"
#include "containers/sllmv.h"

// ----------------------------------------------------------------
sllmv_t* sllmv_alloc() {
	sllmv_t* plist = mlr_malloc_or_die(sizeof(sllmv_t));
	plist->phead  = NULL;
	plist->ptail  = NULL;
	plist->length = 0;
	return plist;
}

// ----------------------------------------------------------------
void sllmv_free(sllmv_t* plist) {
	if (plist == NULL)
		return;
	sllmve_t* pnode = plist->phead;
	while (pnode != NULL) {
		sllmve_t* pdel = pnode;
		pnode = pnode->pnext;
		mv_free(&pdel->value);
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;

	free(plist);
}

// ----------------------------------------------------------------
// Mlrvals are small structs and we do struct assignment from argument
// to list storage. For all but string mlrvals, this is a copy.
// For string mlrvals, it is pointer assignment without string duplication.
// This is intentional (for performance); callees are advised.
void sllmv_add(sllmv_t* plist, mv_t* pvalue) {
	sllmve_t* pnode = mlr_malloc_or_die(sizeof(sllmve_t));
	pnode->value = *pvalue; // struct assignment
	if (plist->ptail == NULL) {
		pnode->pnext = NULL;
		plist->phead = pnode;
		plist->ptail = pnode;
	} else {
		pnode->pnext = NULL;
		plist->ptail->pnext = pnode;
		plist->ptail = pnode;
	}
	plist->length++;
}

// ----------------------------------------------------------------
sllmv_t* sllmv_single(mv_t* pvalue) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add(psllmv, pvalue);
	return psllmv;
}

sllmv_t* sllmv_double(mv_t* pvalue1, mv_t* pvalue2) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add(psllmv, pvalue1);
	sllmv_add(psllmv, pvalue2);
	return psllmv;
}

sllmv_t* sllmv_triple(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add(psllmv, pvalue1);
	sllmv_add(psllmv, pvalue2);
	sllmv_add(psllmv, pvalue3);
	return psllmv;
}

sllmv_t* sllmv_quadruple(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3, mv_t* pvalue4) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add(psllmv, pvalue1);
	sllmv_add(psllmv, pvalue2);
	sllmv_add(psllmv, pvalue3);
	sllmv_add(psllmv, pvalue4);
	return psllmv;
}

// ----------------------------------------------------------------
void sllmv_print(sllmv_t* plist) {
	printf("[");
	for (sllmve_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
		char* string = mv_alloc_format_val(&pe->value);
		if (pe != plist->phead)
			printf(", ");
		printf("%s", string);
		free(string);
	}
	printf("]\n");
}
// xxx merge
void sllmve_tail_print(sllmve_t* pnode) {
	printf("[");
	for (sllmve_t* pe = pnode; pe != NULL; pe = pe->pnext) {
		char* string = mv_alloc_format_val(&pe->value);
		// xxx if (pe != plist->phead)
			// printf(", ");
		printf(" ");
		printf("%s", string);
		free(string);
	}
	printf("]\n");
}
