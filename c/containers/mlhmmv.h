// ================================================================
// Array-only (open addressing) multi-level hash map, with linear probing for collisions.
// All keys, and terminal-level values, are mlrvals.
//
// xxx note about all data being copied, none pointer-reffed with free-flags?
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
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllmv_t* pmvkeys);

//// Unit-test hook
//int mlhmmv_check_counts(mlhmmv_t* pmap);

void mlhmmv_print(mlhmmv_t* pmap);


#endif // MLHMMV_H
