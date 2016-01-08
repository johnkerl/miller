// ================================================================
// Array-only (open addressing) string-list-to-void-star linked hash map with
// linear probing for collisions.
//
// John Kerl 2014-12-22
//
// Notes:
// * null key is not supported.
// * null value is supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#ifndef LHMSLV_H
#define LHMSLV_H

#include "containers/slls.h"

// ----------------------------------------------------------------
typedef struct _lhmslve_t {
	int     ideal_index;
	slls_t* key;
	void*   pvvalue;
	char    free_flags;
	struct _lhmslve_t *pprev;
	struct _lhmslve_t *pnext;
} lhmslve_t;

typedef unsigned char lhmslve_state_t;

// ----------------------------------------------------------------
typedef struct _lhmslv_t {
	int              num_occupied;
	int              num_freed;
	int              array_length;
	lhmslve_t*       entries;
	lhmslve_state_t* states;
	lhmslve_t*       phead;
	lhmslve_t*       ptail;
} lhmslv_t;

lhmslv_t* lhmslv_alloc();
void   lhmslv_free(lhmslv_t* pmap);
void*  lhmslv_put(lhmslv_t* pmap, slls_t* key, void* pvvalue, char free_flags);
void*  lhmslv_get(lhmslv_t* pmap, slls_t* key);
int    lhmslv_has_key(lhmslv_t* pmap, slls_t* key);
int    lhmslv_size(lhmslv_t* pmap);

// Unit-test hook
int lhmslv_check_counts(lhmslv_t* pmap);

#endif // LHMSLV_H
