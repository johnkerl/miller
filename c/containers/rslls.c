#include <string.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/rslls.h"

// []
// [a]
// [a,b]
// [a,b,c]
// [x,x,x]

// ----------------------------------------------------------------
rslls_t* rslls_alloc() {
	rslls_t* plist = mlr_malloc_or_die(sizeof(rslls_t));
	plist->phead  = NULL;
	plist->ptail  = NULL;
	plist->length = 0;
	return plist;
}

// ----------------------------------------------------------------
void rslls_reset(rslls_t* plist) {
	if (plist == NULL)
		return;
	if (plist->phead == NULL)
		return;
	rsllse_t* pnode = plist->phead;
	while (pnode != NULL) {
		if (pnode->free_flag & RSLLS_FREE_ENTRY_VALUE)
			free(pnode->value);
		pnode->value = NULL;
		pnode = pnode->pnext;
	}
	plist->ptail = plist->phead;
	plist->length = 0;
}

// ----------------------------------------------------------------
void rslls_free(rslls_t* plist) {
	if (plist == NULL)
		return;
	rsllse_t* pnode = plist->phead;
	while (pnode != NULL) {
		rsllse_t* pdel = pnode;
		pnode = pnode->pnext;
		if (pdel->free_flag & RSLLS_FREE_ENTRY_VALUE)
			free(pdel->value);
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;
	free(plist);
}

// ----------------------------------------------------------------
static inline void rslls_add(rslls_t* plist, char* value, char free_flag) {
	if (plist->ptail == NULL) {
		// First add on new list
		rsllse_t* pnode = mlr_malloc_or_die(sizeof(rsllse_t));
		pnode->value = value;
		pnode->free_flag = free_flag;
		pnode->pnext = NULL;
		plist->phead = pnode;
		plist->ptail = pnode;
	} else if (plist->ptail->value == NULL) {
		// Subsequent add on reused list
		plist->ptail->value = value;
		plist->ptail->free_flag = free_flag;
		if (plist->ptail->pnext != NULL)
			plist->ptail = plist->ptail->pnext;
	} else {
		// Append at end of list
		rsllse_t* pnode = mlr_malloc_or_die(sizeof(rsllse_t));
		pnode->value = value;
		pnode->free_flag = free_flag;
		pnode->pnext = NULL;
		plist->ptail->pnext = pnode;
		plist->ptail = pnode;
	}
	plist->length++;
}

void rslls_add_with_free(rslls_t* plist, char* value) {
	rslls_add(plist, value, RSLLS_FREE_ENTRY_VALUE);
}
void rslls_add_no_free(rslls_t* plist, char* value) {
	rslls_add(plist, value, 0);
}

void rslls_print(rslls_t* plist) {
	if (plist == NULL) {
		printf("NULL");
	} else {
		int i = 0;
		for (rsllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext, i++) {
			if (i > 0)
				printf(",");
			if (pe->value == NULL)
				printf("NULL");
			else
				printf("%s", pe->value);
		}
	}
}
