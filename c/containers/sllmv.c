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
		if (pdel->free_flags & FREE_ENTRY_VALUE)
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

static void sllmv_add(sllmv_t* plist, mv_t* pvalue, char free_flags) {
	sllmve_t* pnode = mlr_malloc_or_die(sizeof(sllmve_t));
	pnode->value = *pvalue; // struct assignment
	pnode->free_flags = free_flags;
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

void sllmv_add_with_free(sllmv_t* plist, mv_t* pvalue) {
	sllmv_add(plist, pvalue, FREE_ENTRY_VALUE);
}

void sllmv_add_no_free(sllmv_t* plist, mv_t* pvalue) {
	sllmv_add(plist, pvalue, NO_FREE);
}

// ----------------------------------------------------------------
sllmv_t* sllmv_single_no_free(mv_t* pvalue) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add_no_free(psllmv, pvalue);
	return psllmv;
}

sllmv_t* sllmv_single_with_free(mv_t* pvalue) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add_with_free(psllmv, pvalue);
	return psllmv;
}

sllmv_t* sllmv_double_with_free(mv_t* pvalue1, mv_t* pvalue2) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add_with_free(psllmv, pvalue1);
	sllmv_add_with_free(psllmv, pvalue2);
	return psllmv;
}

sllmv_t* sllmv_triple_with_free(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add_with_free(psllmv, pvalue1);
	sllmv_add_with_free(psllmv, pvalue2);
	sllmv_add_with_free(psllmv, pvalue3);
	return psllmv;
}

sllmv_t* sllmv_quadruple_with_free(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3, mv_t* pvalue4) {
	sllmv_t* psllmv = sllmv_alloc();
	sllmv_add_with_free(psllmv, pvalue1);
	sllmv_add_with_free(psllmv, pvalue2);
	sllmv_add_with_free(psllmv, pvalue3);
	sllmv_add_with_free(psllmv, pvalue4);
	return psllmv;
}

// ----------------------------------------------------------------
void sllmv_print(sllmv_t* plist) {
	sllmve_tail_print(plist->phead);
}

void sllmve_tail_print(sllmve_t* pnode) {
	printf("[");
	int i = 0;
	for (sllmve_t* pe = pnode; pe != NULL; pe = pe->pnext, i++) {
		char* string = mv_alloc_format_val(&pe->value);
		if (i > 0)
			printf(", ");
		printf(" ");
		printf("%s", string);
		free(string);
	}
	printf("]\n");
}
