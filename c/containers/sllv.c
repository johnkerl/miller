#include "lib/mlrutil.h"
#include "containers/sllv.h"

// ----------------------------------------------------------------
sllv_t* sllv_alloc() {
	sllv_t* plist = mlr_malloc_or_die(sizeof(sllv_t));
	plist->phead  = NULL;
	plist->ptail  = NULL;
	plist->length = 0;
	return plist;
}

// ----------------------------------------------------------------
void sllv_free(sllv_t* plist) {
	sllve_t* pnode = plist->phead;
	while (pnode != NULL) {
		sllve_t* pdel = pnode;
		pnode = pnode->pnext;
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;

	free(plist);
}

// ----------------------------------------------------------------
sllv_t* sllv_single(void* pvdata) {
	sllv_t* psllv = sllv_alloc();
	sllv_add(psllv, pvdata);
	return psllv;
}

// ----------------------------------------------------------------
void sllv_add(sllv_t* plist, void* pvdata) {
	sllve_t* pnode = mlr_malloc_or_die(sizeof(sllve_t));
	pnode->pvdata = pvdata;
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
void sllv_reverse(sllv_t* plist) {
	if (plist->phead == NULL)
		return;

	sllve_t* pnewhead = NULL;
	sllve_t* pnewtail = plist->phead;
	sllve_t* p = plist->phead;
	sllve_t* q = p->pnext;
	while (1) {
		p->pnext = pnewhead;
		pnewhead = p;
		if (q == NULL)
			break;
		p = q;
		q = p->pnext;
	}
	plist->phead = pnewhead;
	plist->ptail = pnewtail;
}

// ----------------------------------------------------------------
void* sllv_pop(sllv_t* plist) {
	// Zero entries in list
	if (plist->phead == NULL)
		return NULL;

	void* pval = plist->phead->pvdata;
	// One entry in list
	if (plist->phead->pnext == NULL) {
		free(plist->phead);
		plist->phead  = NULL;
		plist->ptail  = NULL;
		plist->length = 0;
	}
	// Two or more entries in list
	else {
		sllve_t* pnext = plist->phead->pnext;
		free(plist->phead);
		plist->phead = pnext;
		plist->length--;
	}

	return pval;
}

// ----------------------------------------------------------------
// xxx cmt mem-mgt
// xxx cmt this can create circular lists if called repeatedly
sllv_t* sllv_append(sllv_t* pa, sllv_t* pb) {
	if (pa == NULL || pa->length == 0)
		return pb;
	if (pb == NULL || pb->length == 0)
		return pa;
	pa->length += pb->length;
	pa->ptail->pnext = pb->phead;
	pa->ptail = pb->ptail;
	return pa;
}

// ================================================================
#ifdef __SLLV_MAIN__
static void print_list(sllv_t* pa, char* desc) {
	printf("%s [%d]:\n", desc, pa->length);
	for (sllve_t* pe = pa->phead; pe != NULL; pe = pe->pnext) {
		printf("  %s\n", (char*)pe->pvdata);
	}
}

int main(int argc, char** argv) {
	sllv_t* pa = sllv_alloc();
	sllv_add(pa, "a");
	sllv_add(pa, "b");
	sllv_add(pa, "c");

	sllv_t* pb = sllv_alloc();
	sllv_add(pb, "d");
	sllv_add(pb, "e");

	print_list(pa, "A");
	print_list(pb, "B");

	pa = sllv_append(pa, pb);
	print_list(pa, "A+B");

	return 0;
}
#endif // __SLLV_MAIN__
