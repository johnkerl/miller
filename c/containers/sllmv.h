// ================================================================
// Singly-linked list of mlrval, with tail for append.
// Everything is copied; nothing is referenced.
// ================================================================

#ifndef SLLMV_H
#define SLLMV_H

#include "mlrval.h"

typedef struct _sllmve_t {
	mv_t* pvalue;
	struct _sllmve_t *pnext;
} sllmve_t;

typedef struct _sllmv_t {
	sllmve_t *phead;
	sllmve_t *ptail;
	int length;
} sllmv_t;

sllmv_t* sllmv_alloc();
void     sllmv_free(sllmv_t* plist);
void     sllmv_add(sllmv_t* plist, mv_t* pvalue);

sllmv_t* sllmv_single(mv_t* pvalue);
sllmv_t* sllmv_double(mv_t* pvalue1, mv_t* pvalue2);
sllmv_t* sllmv_triple(mv_t* pvalue1, mv_t* pvalue2, mv_t* pvalue3);

void sllmv_print(sllmv_t* plist);

#endif // SLLMV_H

