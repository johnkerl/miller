// ================================================================
// Array-only (open addressing) multi-level hash map, with linear probing for collisions.
// All keys, and terminal-level values, are mlrvals. All data passed into the put method
// are copied; no pointers in this data structure reference anything external.
//
// Notes:
// * null key is not supported.
// * null value is not supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#ifndef MLHMMV_H
#define MLHMMV_H

#include "containers/mlrval.h"
#include "containers/sllmv.h"

#define MLHMMV_ERROR_NONE                0x0000
#define MLHMMV_ERROR_KEYLIST_TOO_DEEP    0xdeef
#define MLHMMV_ERROR_KEYLIST_TOO_SHALLOW 0x58a1

// This is made visible here in the API so the unit-tester can be sure to exercise the resize logic.
#define MLHMMV_INITIAL_ARRAY_LENGTH 16

struct _mlhmmv_level_t; // forward reference

// ----------------------------------------------------------------
typedef struct _mlhmmv_level_value_t {
	int is_terminal;
	union {
		mv_t mlrval;
		struct _mlhmmv_level_t* pnext_level;
	} u;
} mlhmmv_level_value_t;

typedef struct _mlhmmv_level_entry_t {
	int     ideal_index;
	mv_t    level_key;
	mlhmmv_level_value_t level_value; // terminal mlrval, or another hashmap
	struct _mlhmmv_level_entry_t *pprev;
	struct _mlhmmv_level_entry_t *pnext;
} mlhmmv_level_entry_t;

typedef unsigned char mlhmmv_level_entry_state_t;

// ----------------------------------------------------------------
typedef struct _mlhmmv_level_t {
	int                         num_occupied;
	int                         num_freed;
	int                         array_length;
	mlhmmv_level_entry_t*       entries;
	mlhmmv_level_entry_state_t* states;
	mlhmmv_level_entry_t*       phead;
	mlhmmv_level_entry_t*       ptail;
} mlhmmv_level_t;

// ----------------------------------------------------------------
typedef struct _mlhmmv_t {
	mlhmmv_level_t* proot_level;
} mlhmmv_t;

mlhmmv_t* mlhmmv_alloc();
void  mlhmmv_free(mlhmmv_t* pmap);
// pmvkeys is a list of mlhmmv_level_value_t
void  mlhmmv_put(mlhmmv_t* pmap, sllmv_t* pmvkeys, mv_t* pterminal_value);

// If the return value is non-null, error will be MLHMMV_ERROR_NONE.  If the
// return value is null, the error will be MLHMMV_ERROR_KEYLIST_TOO_DEEP or
// MLHMMV_ERROR_KEYLIST_TOO_SHALLOW, or MLHMMV_ERROR_NONE if the keylist matches
// map depth but the entry is not found.
//
// Note: this returns a pointer to the map's data, not to a copy.
// The caller shouldn't free it, or modify it.
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllmv_t* pmvkeys, int* perror);

// Made public for the oosvar-emitter.
mlhmmv_level_value_t* mlhmmv_level_get(mlhmmv_level_t* pmap, mv_t* plevel_key);

void mlhmmv_print_json_stacked(mlhmmv_t* pmap, int quote_values_always);
void mlhmmv_print_json_single_line(mlhmmv_t* pmap, int quote_values_always);

#endif // MLHMMV_H
