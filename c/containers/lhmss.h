// ================================================================
// Array-only (open addressing) string-to-string linked hash map with linear
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

#ifndef LHMSS_H
#define LHMSS_H

#include "containers/sllv.h"

// ----------------------------------------------------------------
typedef struct _lhmsse_t {
	int   ideal_index;
	char  free_flags;
	char* key;
	char* value;
	struct _lhmsse_t *pprev;
	struct _lhmsse_t *pnext;
} lhmsse_t;

typedef unsigned char lhmsse_state_t;

typedef struct _lhmss_t {
	int             num_occupied;
	int             num_freed;
	int             array_length;
	lhmsse_t*       entries;
	lhmsse_state_t* states;
	lhmsse_t*       phead;
	lhmsse_t*       ptail;
} lhmss_t;

// ----------------------------------------------------------------
lhmss_t* lhmss_alloc();
lhmss_t* lhmss_copy(lhmss_t* pmap);
void  lhmss_free(lhmss_t* pmap);
void  lhmss_put(lhmss_t* pmap, char* key, char* value, char free_flags);
char* lhmss_get(lhmss_t* pmap, char* key);
int   lhmss_has_key(lhmss_t* pmap, char* key);
void  lhmss_rename(lhmss_t* pmap, char* old_key, char* new_key);

void lhmss_print(lhmss_t* pmap);

// Unit-test hook
int lhmss_check_counts(lhmss_t* pmap);

#endif // LHMSS_H
