// ================================================================
// Singly-linked list of mlrval, with tail for append.
// Strings inside mlrvals are referenced, not copied.
// ================================================================

#ifndef SLLMV_H
#define SLLMV_H

#include "../lib/mlrval.h"
#include "../lib/free_flags.h"

// ----------------------------------------------------------------
typedef struct _sllmve_t {
	mv_t value;
	char free_flags;
	struct _sllmve_t *pnext;
} sllmve_t;

typedef struct _sllmv_t {
	sllmve_t *phead;
	sllmve_t *ptail;
	int length;
} sllmv_t;

// ----------------------------------------------------------------
sllmv_t* sllmv_alloc();
void     sllmv_free(sllmv_t* plist);
void     sllmv_prepend_with_free(sllmv_t* plist, mv_t* pvalue);
void     sllmv_prepend_no_free(sllmv_t* plist, mv_t* pvalue);
void     sllmv_append_with_free(sllmv_t* plist, mv_t* pvalue);
void     sllmv_append_no_free(sllmv_t* plist, mv_t* pvalue);

sllmv_t* sllmv_single_no_free(mv_t* pvalue);
sllmv_t* sllmv_single_with_free(mv_t* pvalue);
sllmv_t* sllmv_double_with_free(mv_t* pvalue1, mv_t* pvalue2);
sllmv_t* sllmv_triple_with_free(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3);
sllmv_t* sllmv_quadruple_with_free(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3, mv_t* pvalue4);

void sllmv_print(sllmv_t* plist);
void sllmve_tail_print(sllmve_t* pnode);

#endif // SLLMV_H

