// ================================================================
// Array-only (open addressing) string-valued hash set with linear probing for
// collisions.
//
// Notes:
// * null key is not supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#ifndef HSS_H
#define HSS_H

// ----------------------------------------------------------------
typedef struct _hsse_t {
	char* key;
	int   state;
	int   ideal_index;
} hsse_t;

// ----------------------------------------------------------------
typedef struct _hss_t {
	int num_occupied;
	int num_freed;
	int array_length;
	hsse_t* array;
} hss_t;

// ----------------------------------------------------------------
hss_t* hss_alloc();
void   hss_free(hss_t* pset);
void   hss_add(hss_t* pset, char* key);
int    hss_has(hss_t* pset, char* key);
void   hss_remove(hss_t* pset, char* key);
void   hss_clear(hss_t* pset);
int    hss_size(hss_t* pset);

#endif // HSS_H
