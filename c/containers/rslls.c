#include <string.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/rslls.h"

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
		if (pnode->free_flag & FREE_ENTRY_VALUE)
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
		if (pdel->free_flag & FREE_ENTRY_VALUE)
			free(pdel->value);
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;
	free(plist);
}

// ----------------------------------------------------------------
void rslls_append(rslls_t* plist, char* value, char free_flag, char quote_flag) {
	if (plist->ptail == NULL) {
		// First add on new list
		rsllse_t* pnode = mlr_malloc_or_die(sizeof(rsllse_t));
		pnode->value = value;
		pnode->free_flag = free_flag;
		pnode->quote_flag = quote_flag;
		pnode->pnext = NULL;
		plist->phead = pnode;
		plist->ptail = pnode;
	} else if (plist->ptail->value == NULL) {
		// Subsequent add on reused list
		plist->ptail->value = value;
		plist->ptail->free_flag = free_flag;
		plist->ptail->quote_flag = quote_flag;
		if (plist->ptail->pnext != NULL)
			plist->ptail = plist->ptail->pnext;
	} else {
		// Append at end of list
		rsllse_t* pnode = mlr_malloc_or_die(sizeof(rsllse_t));
		pnode->value = value;
		pnode->free_flag = free_flag;
		pnode->quote_flag = quote_flag;
		pnode->pnext = NULL;
		plist->ptail->pnext = pnode;
		plist->ptail = pnode;
	}
	plist->length++;
}

void rslls_print(rslls_t* plist) {
	if (plist == NULL) {
		printf("NULL");
	} else {
		unsigned long long i = 0;
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
