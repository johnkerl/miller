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
		mv_free(pdel->pvalue);
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;

	free(plist);
}

// ----------------------------------------------------------------
void sllmv_add(sllmv_t* plist, mv_t* pvalue) {
	sllmve_t* pnode = mlr_malloc_or_die(sizeof(sllmve_t));
	pnode->pvalue = mv_alloc_copy(pvalue);
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
