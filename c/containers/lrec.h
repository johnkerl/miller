// ================================================================
// This is a hashless implementation of insertion-ordered key-value pairs for
// Miller's fundamental record data structure.  It implements the same
// interface as the hashed version (see lhmss.h).
//
// Design:
//
// * It keeps a doubly-linked list of key-value pairs.
// * No hash functions are computed when the map is written to or read from.
// * Gets are implemented by sequential scan through the list: given a key,
//   the key-value pairs are scanned through until a match is (or is not) found.
// * Performance improvement of 10-15% percent over lhmss is found (for test data).
//
// Motivation:
//
// * The use case for records in Miller is that *all* fields are read from
//   strings & written to strings (split/join), while only *some* fields are
//   operated on.
//
// * Meanwhile there are few repeated accesses to a given record: the
//   access-to-construct ratio is quite low for Miller data records.  Miller
//   instantiates thousands, millions, billions of records (depending on the
//   input data) but accesses each record only once per mapping operation.
//   (This is in contrast to accumulator hashmaps which are repeatedly accessed
//   during a stats run.)
//
// * The hashed impl computes hashsums for *all* fields whether operated on or not,
//   for the benefit of the *few* fields looked up during the mapping operation.
//
// * The hashless impl only keeps string pointers.  Lookups are done at runtime
//   doing prefix search on the key names. Assuming field names are distinct,
//   this is just a few char-ptr accesses which (in experiments) turn out to
//   offer about a 10-15% performance improvement.
//
// * Added benefit: the field-rename operation (preserving field order) becomes
//   trivial.
//
// Notes:
// * null key is not supported.
// * null value is supported.
// ================================================================

#ifndef LREC_H
#define LREC_H

#include "containers/sllv.h"
#include "containers/header_keeper.h"

#define LREC_FREE_ENTRY_KEY        0x08
#define LREC_FREE_ENTRY_VALUE      0x80

struct _lrec_t; // forward reference
typedef struct _lrec_t lrec_t;

typedef void lrec_free_func_t(lrec_t* prec);

// ----------------------------------------------------------------
typedef struct _lrece_t {
	char* key;
	char* value;
	// These indicate whether the key/value should be freed on lrec_free().
	// Affirmative example: key/value is strdup of something.
	// Negative example: key/value are pointers into a line the memory
	// management of which is separately managed.
	// Another negative example: key/value is a string literal, e.g. "".
	char free_flags;

	struct _lrece_t *pprev;
	struct _lrece_t *pnext;
} lrece_t;

struct _lrec_t {
	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	int      field_count;
	lrece_t* phead;
	lrece_t* ptail;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// See comments above free_flags. Used to track a mallocked pointer to be
	// freed at lrec_free().

	// E.g. for NIDX, DKVP, and CSV formats (header handled separately in the
	// latter case).
	char* psingle_line;

	// For XTAB format.
	slls_t* pxtab_lines;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Format-dependent virtual-function pointer:
	lrec_free_func_t* pfree_backing_func;
};

// ----------------------------------------------------------------
lrec_t* lrec_unbacked_alloc();
lrec_t* lrec_dkvp_alloc(char* line);
lrec_t* lrec_nidx_alloc(char* line);
lrec_t* lrec_csvlite_alloc(char* data_line);
lrec_t* lrec_csv_alloc(char* data_line);
lrec_t* lrec_xtab_alloc(slls_t* pxtab_lines);

void  lrec_put_no_free(lrec_t* prec, char* key, char* value);
void  lrec_put(lrec_t* prec, char* key, char* value, char free_flags);

char* lrec_get(lrec_t* prec, char* key);
void  lrec_remove(lrec_t* prec, char* key);
void  lrec_rename(lrec_t* prec, char* old_key, char* new_key, int new_needs_freeing);
void  lrec_move_to_head(lrec_t* prec, char* key);
void  lrec_move_to_tail(lrec_t* prec, char* key);

void  lrec_free(lrec_t* prec);

void lrec_print(lrec_t* prec);
void lrec_dump(lrec_t* prec);
void lrec_dump_titled(char* msg, lrec_t* prec);

// NIDX data are keyed by one-up field index which is not explicitly contained
// in the file, e.g. line "a b c" splits to an lrec with "{"1" => "a", "2" =>
// "b", "3" => "c"}. This function creates the keys, avoiding redundant memory
// allocation for most-used keys such as "1", "2", ... up to 100 or so. In case
// of large idx, free_flags & LREC_FREE_ENTRY_KEY will indicate that the key
// was dynamically allocated.
char* make_nidx_key(int idx, char* pfree_flags);

// For unit-test.
lrec_t* lrec_literal_1(char* k1, char* v1);
lrec_t* lrec_literal_2(char* k1, char* v1, char* k2, char* v2);
lrec_t* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3);
lrec_t* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4);

#endif // LREC_H
