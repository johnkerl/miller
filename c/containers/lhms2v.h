// ================================================================
// Array-only (open addressing) string-pair-to-void-star linked hash map with
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

#ifndef LHMS2V_H
#define LHMS2V_H

// ----------------------------------------------------------------
typedef struct _lhms2ve_t {
	int   ideal_index;
	char* key1;
	char* key2;
	void* pvvalue;
	struct _lhms2ve_t *pprev;
	struct _lhms2ve_t *pnext;
} lhms2ve_t;

typedef unsigned char lhms2ve_state_t;

// ----------------------------------------------------------------
typedef struct _lhms2v_t {
	int              num_occupied;
	int              num_freed;
	int              array_length;
	lhms2ve_t*       entries;
	lhms2ve_state_t* states;
	lhms2ve_t*       phead;
	lhms2ve_t*       ptail;
} lhms2v_t;

lhms2v_t* lhms2v_alloc();
void   lhms2v_free(lhms2v_t* pmap);
void*  lhms2v_put(lhms2v_t* pmap, char* key1, char* key2, void* pvvalue);
void*  lhms2v_get(lhms2v_t* pmap, char* key1, char* key2);
int    lhms2v_has_key(lhms2v_t* pmap, char* key1, char* key2);
void*  lhms2v_remove(lhms2v_t* pmap, char* key1, char* key2);
void   lhms2v_clear(lhms2v_t* pmap);
int    lhms2v_size(lhms2v_t* pmap);

// Unit-test hook
int    lhms2v_check_counts(lhms2v_t* pmap);

#endif // LHMS2V_H
