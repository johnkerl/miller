// ================================================================
// This is a hashless implementation of insertion-ordered key-value pairs.
// It implements the same interface as the hashed version (see lhmss.h).
//
// Design:
//
// * It keeps a doubly-linked list of key-value pairs.
// * No hash functions are computed when the map is written to or read from.
// * Gets are implemented by sequential scan through the list: given a key,
//   scan through the key-value pairs until a match is (or is not) found.
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
//   (This is in contrast to hashmaps which are repeatedly accessed during a
//   stats run.)
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
#include "containers/hdr_keeper.h"

#define LREC_FREE_ENTRY_KEY        0x08
#define LREC_FREE_ENTRY_VALUE      0x80

// xxx cmt for virtual-function ptrs; not for API use
struct _lrec_t; // forward reference
typedef struct _lrec_t lrec_t;

typedef void lrec_free_func_t(lrec_t* prec);

// ----------------------------------------------------------------
typedef struct _lrece_t {
	char* key;
	char* value;
	// xxx cmt
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
	// xxx cmt why (baton-handoff memory management)
	// xxx cmt csv hdr-alloc handled at reader

	char* psingle_line;

	char* pcsv_data_line; // xxx merge w/ single-line

	slls_t* pxtab_lines;

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// xxx cmt virtual-function pointers:

	lrec_free_func_t* pfree_backing_func;
};

// ----------------------------------------------------------------
lrec_t* lrec_unbacked_alloc();
lrec_t* lrec_dkvp_alloc(char* line);
lrec_t* lrec_nidx_alloc(char* line);
lrec_t* lrec_csv_alloc(char* data_line);
lrec_t* lrec_xtab_alloc(slls_t* pxtab_lines);

void  lrec_put_no_free(lrec_t* prec, char* key, char* value);
void  lrec_put(lrec_t* prec, char* key, char* value, char free_flags);

char* lrec_get(lrec_t* prec, char* key);
void  lrec_remove(lrec_t* prec, char* key);
void  lrec_rename(lrec_t* prec, char* old_key, char* new_key);
void  lrec_set_name(lrec_t* prec, lrece_t* pfield, char* new_key);
void  lrec_move_to_head(lrec_t* prec, char* key);
void  lrec_move_to_tail(lrec_t* prec, char* key);

void  lrec_free(lrec_t* prec);

void lrec_dump(lrec_t* prec);
void lrec_dump_titled(char* msg, lrec_t* prec);

#endif // LREC_H
