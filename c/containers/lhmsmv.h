// ================================================================
// Array-only (open addressing) string-to-mlrval linked hash map with linear
// probing for collisions.
//
// John Kerl 2012-08-13
//
// Notes:
// * null key is not supported.
// * null value is supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#ifndef LHMSMV_H
#define LHMSMV_H

#include "containers/sllv.h"
#include "lib/mlrval.h"

// ----------------------------------------------------------------
typedef struct _lhmsmve_t {
	int   ideal_index;
	char  free_flags;
	char* key;
	mv_t  value;
	struct _lhmsmve_t *pprev;
	struct _lhmsmve_t *pnext;
} lhmsmve_t;

typedef unsigned char lhmsmve_state_t;

typedef struct _lhmsmv_t {
	int              num_occupied;
	int              num_freed;
	int              array_length;
	lhmsmve_t*       entries;
	lhmsmve_state_t* states;
	lhmsmve_t*       phead;
	lhmsmve_t*       ptail;
} lhmsmv_t;

// ----------------------------------------------------------------
lhmsmv_t* lhmsmv_alloc();
lhmsmv_t* lhmsmv_copy(lhmsmv_t* pmap);
void  lhmsmv_clear(lhmsmv_t* pmap);
void  lhmsmv_free(lhmsmv_t* pmap);

void  lhmsmv_put(lhmsmv_t* pmap, char* key, mv_t* pvalue, char free_flags);
mv_t* lhmsmv_get(lhmsmv_t* pmap, char* key);
int   lhmsmv_has_key(lhmsmv_t* pmap, char* key);

void  lhmsmv_dump(lhmsmv_t* pmap);
int lhmsmv_check_counts(lhmsmv_t* pmap);

#endif // LHMSMV_H
