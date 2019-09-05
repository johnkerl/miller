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

#include "lib/free_flags.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/hss.h"
#include "containers/header_keeper.h"

#define FIELD_QUOTED_ON_INPUT 0x02

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
	char quote_flags;

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

void lrec_clear(lrec_t* prec);
void  lrec_free(lrec_t* prec);
lrec_t* lrec_copy(lrec_t* pinrec);

// The only difference between lrec_put and lrec_prepend is that the latter
// adds to the end of the record, while the former adds to the beginning.
//
// For both, the key/value respectively will be freed by lrec_free if the
// corresponding bits are set in the free_flags.
//
// * If a string literal or other non-allocated pointer (e.g. mmapped memory
//   from a file reader) is passed in, the free flag should not be set.
//
// * If dynamically allocated pointers are passed in, then either:
//
//   o The respective free_flag(s) should be set and the caller should be sure
//     not to also free (else, there will be heap corruption due to
//     double-free), or
//
//   o The respective free_flag(s) should not be set and the caller should
//     free the memory (else, there will be a memory leak).
void  lrec_put(lrec_t* prec, char* key, char* value, char free_flags);
void  lrec_put_ext(lrec_t* prec, char* key, char* value, char free_flags, char quote_flags);
// Like lrec_put: if key is present, modify value. But if not, add new field at start of record, not at end.
void  lrec_prepend(lrec_t* prec, char* key, char* value, char free_flags);
// Like lrec_put: if key is present, modify value. But if not, add new field after specified entry, not at end.
// Returns a pointer to the added/modified node.
lrece_t*  lrec_put_after(lrec_t* prec, lrece_t* pd, char* key, char* value, char free_flags);

char* lrec_get(lrec_t* prec, char* key);
lrece_t* lrec_get_pair_by_position(lrec_t* prec, int position); // 1-up not 0-up
char* lrec_get_key_by_position(lrec_t* prec, int position); // 1-up not 0-up
char* lrec_get_value_by_position(lrec_t* prec, int position); // 1-up not 0-up

// This returns a pointer to the lrec's free-flags so that the caller can do ownership-transfer
// of about-to-be-removed key-value pairs.
char* lrec_get_pff(lrec_t* prec, char* key, char** ppfree_flags);

// This returns a pointer to the entry so the caller can update it directly without needing
// to do another field-scan on subsequent lrec_put etc. This is a performance optimization;
// it also allows mlr nest --explode to do explode-in-place rather than explode-at-end.
char* lrec_get_ext(lrec_t* prec, char* key, lrece_t** ppentry);

void  lrec_remove(lrec_t* prec, char* key);
void  lrec_remove_by_position(lrec_t* prec, int position); // 1-up not 0-up
void  lrec_rename(lrec_t* prec, char* old_key, char* new_key, int new_needs_freeing);
void  lrec_rename_at_position(lrec_t* prec, int position, char* new_key, int new_needs_freeing); // 1-up not 0-up
void  lrec_move_to_head(lrec_t* prec, char* key);
void  lrec_move_to_tail(lrec_t* prec, char* key);
// Renames the first n fields where n is the length of pnames.
// The hash-set argument is for efficient dedupe.
// Assumes as a precondition that pnames_as_list has no duplicates.
// If the new labels include any field names existing later on in the record, those are unset.
// For example, input record "a=1,b=2,c=3,d=4,e=5" with labels "d,x,f" results in output record "d=1,x=2,f=3,e=5".
void  lrec_label(lrec_t* prec, slls_t* pnames_as_list, hss_t* pnames_as_set);

void lrece_update_value(lrece_t* pe, char* new_value, int new_needs_freeing);

// For lrec-internal use:
void lrec_unlink(lrec_t* prec, lrece_t* pe);
// May be used for removing fields from a record while iterating over it:
void lrec_unlink_and_free(lrec_t* prec, lrece_t* pe);

void lrec_print(lrec_t* prec);
void lrec_dump(lrec_t* prec);
void lrec_dump_fp(lrec_t* prec, FILE* fp);
void lrec_dump_titled(char* msg, lrec_t* prec);
void lrec_pointer_dump(lrec_t* prec);
// The caller should free the return value
char* lrec_sprint(lrec_t* prec, char* ors, char* ofs, char* ops);

// NIDX data are keyed by one-up field index which is not explicitly contained
// in the file, e.g. line "a b c" splits to an lrec with "{"1" => "a", "2" =>
// "b", "3" => "c"}. This function creates the keys, avoiding redundant memory
// allocation for most-used keys such as "1", "2", ... up to 100 or so. In case
// of large idx, free_flags & FREE_ENTRY_KEY will indicate that the key
// was dynamically allocated.
char* low_int_to_string(int idx, char* pfree_flags);

// For unit-test.
lrec_t* lrec_literal_1(char* k1, char* v1);
lrec_t* lrec_literal_2(char* k1, char* v1, char* k2, char* v2);
lrec_t* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3);
lrec_t* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4);

#endif // LREC_H
