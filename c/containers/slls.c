#include <string.h>
#include "lib/mlrutil.h"
#include "containers/slls.h"

// ----------------------------------------------------------------
slls_t* slls_alloc() {
	slls_t* plist = mlr_malloc_or_die(sizeof(slls_t));
	plist->phead  = NULL;
	plist->ptail  = NULL;
	plist->length = 0;
	return plist;
}

// ----------------------------------------------------------------
slls_t* slls_copy(slls_t* pold) {
	slls_t* pnew = slls_alloc();
	for (sllse_t* pe = pold->phead; pe != NULL; pe = pe->pnext)
		slls_add_with_free(pnew, strdup(pe->value));
	return pnew;
}

// ----------------------------------------------------------------
void slls_free(slls_t* plist) {
	sllse_t* pnode = plist->phead;
	while (pnode != NULL) {
		sllse_t* pdel = pnode;
		pnode = pnode->pnext;
		if (pdel->free_flag & SLLS_FREE_ENTRY_VALUE)
			free(pdel->value);
		free(pdel);
	}
	plist->phead  = NULL;
	plist->ptail  = 0;
	plist->length = 0;
	free(plist);
}

// ----------------------------------------------------------------
slls_t* slls_single_with_free(char* value) {
	slls_t* pslls = slls_alloc();
	slls_add_with_free(pslls, value);
	return pslls;
}

// ----------------------------------------------------------------
slls_t* slls_single_no_free(char* value) {
	slls_t* pslls = slls_alloc();
	slls_add_no_free(pslls, value);
	return pslls;
}

// ----------------------------------------------------------------
// xxx move to API
static inline void slls_add(slls_t* plist, char* value, char free_flag) {
	sllse_t* pnode = mlr_malloc_or_die(sizeof(sllse_t));
	pnode->value = value;
	pnode->free_flag = free_flag;

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

void slls_add_with_free(slls_t* plist, char* value) {
	slls_add(plist, value, SLLS_FREE_ENTRY_VALUE);
}
void slls_add_no_free(slls_t* plist, char* value) {
	slls_add(plist, value, 0);
}

// ----------------------------------------------------------------
int slls_equals(slls_t* pa, slls_t* pb) {
	if (pa->length != pb->length)
		return FALSE;
	sllse_t* pea = pa->phead;
	sllse_t* peb = pb->phead;
	for ( ; pea != NULL && peb != NULL; pea = pea->pnext, peb = peb->pnext) {
		if (!streq(pea->value, peb->value))
			return FALSE;
	}
	return TRUE;
}

// ----------------------------------------------------------------
slls_t* slls_from_line(char* line, char ifs, int allow_repeat_ifs) {
	slls_t* plist = slls_alloc();
	if (*line == 0) // empty string splits to empty list
		return plist;

	char* start = line;
	for (char* p = line; *p; p++) {
		if (*p == ifs) {
			*p = 0;
			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			slls_add_no_free(plist, start);
			start = p;
		}
	}
	slls_add_no_free(plist, start);

	return plist;
}

// ----------------------------------------------------------------
// xxx cmt for debug. inefficient. or fix that.
// xxx rename to slls_alloc_join
char* slls_join(slls_t* plist, char fs) {
	int len = 0;
	for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext)
		len += strlen(pe->value) + 1; // include space for fs and null-terminator
	char* output = mlr_malloc_or_die(len);
	char sep[2] = {fs, 0};
	*output = 0;
	for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
		strcat(output, pe->value);
		if (pe->pnext != NULL) {
			strcat(output, sep);
		}
	}

	return output;
}

void slls_debug_print(slls_t* plist, FILE* output_stream) {
	char* string = slls_join(plist, ',');
	fprintf(output_stream, "%s\n", string);
	free(string);
}

// ----------------------------------------------------------------
void slls_reverse(slls_t* plist) {
	if (plist->phead == NULL)
		return;

	sllse_t* pnewhead = NULL;
	sllse_t* pnewtail = plist->phead;
	sllse_t* p = plist->phead;
	sllse_t* q = p->pnext;
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
// xxx fix this so ["ab","c"] doesn't hash to the same as ["a","bc"].
int slls_hash_func(slls_t *plist) {
	unsigned long hash = 5381;
	int c;

	for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
		char* str = pe->value;
		while ((c = *str++) != 0)
			hash = ((hash << 5) + hash) + c; /* hash * 33 + c */
	}

	return (int)hash;
}

// ----------------------------------------------------------------
#ifdef __SLLS_MAIN__
int main(int argc, char** argv)  {
	slls_t* plist = slls_from_line(argv[1], ',', FALSE);
	printf("len=%d\n", plist->length);
	for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext)
		printf("val=%s\n", (char*)pe->value);
	return 0;
}
#endif
