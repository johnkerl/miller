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
#include "containers/sllv.h"

//// ----------------------------------------------------------------
//typedef struct _mlhmmve_t {
//	int     ideal_index;
//	mv_t*   pvalue;
//	struct _mlhmmve_t *pprev;
//	struct _mlhmmve_t *pnext;
//} mlhmmve_t;
//
//typedef unsigned char mlhmmve_state_t;

// ----------------------------------------------------------------
typedef struct _mlhmmv_level_t {
//	int              num_occupied;
//	int              num_freed;
//	int              array_length;
//	mlhmmv_level_t*  entries;
//	mlhmmv_state_t*  states;
//	mlhmmv_level_t*  phead;
//	mlhmmv_level_t*  ptail;
} mlhmmv_level_t;

// ----------------------------------------------------------------
typedef struct _mlhmmv_t {
//	mlhmmv_level_t* proot_level;
} mlhmmv_t;

mlhmmv_t* mlhmmv_alloc();
void  mlhmmv_free(mlhmmv_t* pmap);
void  mlhmmv_put(mlhmmv_t* pmap, sllv_t* pmvkeys, mv_t* pvalue);
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllv_t* pmvkeys);
int   mlhmmv_has_keys(mlhmmv_t* pmap, sllv_t* pmvkeys);

//// Unit-test hook
//int mlhmmv_check_counts(mlhmmv_t* pmap);

#endif // MLHMMV_H
