// ================================================================
// Array-only (open addressing) string-to-int linked hash map with linear
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

#ifndef LHMSLL_H
#define LHMSLL_H

// ----------------------------------------------------------------
typedef struct _lhmslle_t {
	int   ideal_index;
	char* key;
	long long value;
	char  free_flags;
	struct _lhmslle_t *pprev;
	struct _lhmslle_t *pnext;
} lhmslle_t;

typedef unsigned char lhmslle_state_t;

typedef struct _lhmsll_t {
	int              num_occupied;
	int              num_freed;
	int              array_length;
	lhmslle_t*       entries;
	lhmslle_state_t* states;
	lhmslle_t*       phead;
	lhmslle_t*       ptail;
} lhmsll_t;

// ----------------------------------------------------------------
lhmsll_t* lhmsll_alloc();

lhmsll_t* lhmsll_copy(lhmsll_t* pmap);

void  lhmsll_free(lhmsll_t* pmap);
void  lhmsll_put(lhmsll_t* pmap, char* key, int value, char free_flags);
long long lhmsll_get(lhmsll_t* pmap, char* key); // caller must do lhmsll_has_key to check validity
int lhmsll_test_and_get(lhmsll_t* pmap, char* key, long long* pval); // *pval undefined if return is FALSE
int lhmsll_test_and_increment(lhmsll_t* pmap, char* key); // increments value only if mapping exists
lhmslle_t* lhmsll_get_entry(lhmsll_t* pmap, char* key);
int   lhmsll_has_key(lhmsll_t* pmap, char* key);

// Unit-test hook
int lhmsll_check_counts(lhmsll_t* pmap);

#endif // LHMSLL_H
