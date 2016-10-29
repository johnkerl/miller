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
#include "containers/sllv.h"
#include "containers/lrec.h"

#define MLHMMV_ERROR_NONE                0x0000
#define MLHMMV_ERROR_KEYLIST_TOO_DEEP    0xdeef
#define MLHMMV_ERROR_KEYLIST_TOO_SHALLOW 0x58a1

// This is made visible here in the API so the unit-tester can be sure to exercise the resize logic.
#define MLHMMV_INITIAL_ARRAY_LENGTH 16

struct _mlhmmv_level_t; // forward reference

// ----------------------------------------------------------------
typedef struct _mlhmmv_value_t {
	int is_terminal;
	union {
		mv_t mlrval;
		struct _mlhmmv_level_t* pnext_level;
	} u;
} mlhmmv_value_t;

typedef struct _mlhmmv_level_entry_t {
	int     ideal_index;
	mv_t    level_key;
	mlhmmv_value_t level_value; // terminal mlrval, or another hashmap
	struct _mlhmmv_level_entry_t *pprev;
	struct _mlhmmv_level_entry_t *pnext;
} mlhmmv_level_entry_t;

typedef unsigned char mlhmmv_level_entry_state_t;

// Store a mlrval into the mlhmmv_value without copying, implicitly transferring
// ownership of the mlrval's free_flags. This means the mlrval will be freed
// when the mlhmmv_value is freed, so the caller should make a copy first if
// necessary.
//
// This is a hot path for non-map local-variable assignments.
static inline mlhmmv_value_t mlhmmv_value_transfer_terminal(mv_t val) { // xxx temp?
	return (mlhmmv_value_t) {.is_terminal = TRUE, .u.mlrval = val};
}

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

void mlhmmv_free(mlhmmv_t* pmap);

void mlhmmv_put_terminal(mlhmmv_t* pmap, sllmv_t* pmvkeys, mv_t* pterminal_value);

// If the return value is non-null, error will be MLHMMV_ERROR_NONE.  If the
// return value is null, the error will be MLHMMV_ERROR_KEYLIST_TOO_DEEP or
// MLHMMV_ERROR_KEYLIST_TOO_SHALLOW, or MLHMMV_ERROR_NONE if the keylist matches
// map depth but the entry is not found.
//
// Note: this returns a pointer to the map's data, not to a copy.
// The caller shouldn't free it, or modify it.
mv_t* mlhmmv_get_terminal(mlhmmv_t* pmap, sllmv_t* pmvkeys, int* perror);
mv_t* mlhmmv_get_terminal_from_level(mlhmmv_level_t* plevel, sllmv_t* pmvkeys, int* perror);

// These are an optimization for assignment from full srec, e.g. '@records[$key1][$key2] = $*'.
// Using mlhmmv_get_or_create_level, the CST logic can get or create the @records[$key1][$key2]
// level of the mlhmmv, then copy values there.
mlhmmv_level_t* mlhmmv_get_or_create_level(mlhmmv_t* pmap, sllmv_t* pmvkeys);
void mlhmmv_put_terminal_from_level(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);

// This is an assignment for assignment to full srec, e.g. '$* = @records[$key1][$key2]'.
// The CST logic can use this function to get the @records[$key1][$key2] level of the mlhmmv,
// then copy values from there.
mlhmmv_level_t* mlhmmv_get_level(mlhmmv_t* pmap, sllmv_t* pmvkeys, int* perror);

// For oosvar-to-oosvar assignment.
void mlhmmv_copy(mlhmmv_t* pmap, sllmv_t* ptokeys, sllmv_t* pfromkeys);

// For for-loop-over-oosvar, wherein we need to copy the submap before iterating over it
// (since the iteration may modify it). If the keys don't index a submap, then the return
// value has is_terminal = TRUE and pnext_level = NULL.
mlhmmv_value_t mlhmmv_copy_submap(mlhmmv_t* pmap, sllmv_t* pmvkeys);
void mlhmmv_free_submap(mlhmmv_value_t submap);

sllv_t* mlhmmv_copy_keys_from_submap(mlhmmv_t* pmap, sllmv_t* pmvkeys);

// Unset value/submap from a specified level onward, also unsetting any maps which become empty as a result.
// Examples:
//   {
//     "a" : { "x" : 1, "y" : 2 },
//     "b" : { "x" : 3, "y" : 4 },
//   }
// with pmvkeys = ["a"] leaves
//   {
//     "b" : { "x" : 3, "y" : 4 },
//   }
// but with pmvkeys = ["a", "y"] leaves
//   {
//     "a" : { "x" : 1 },
//     "b" : { "x" : 3, "y" : 4 },
//   }
// and with pmvkeys = [] leaves
//   {
//   }
// Now if ["a","x"] is removed from
//   {
//     "a" : { "x" : 1 },
//     "b" : { "x" : 3, "y" : 4 },
//   }
// then
//   {
//     "b" : { "x" : 3, "y" : 4 },
//   }
// is left: unsetting "a":"x" leaves the map at "a" so this is unset as well.
void mlhmmv_remove(mlhmmv_t* pmap, sllmv_t* pmvkeys);

void mlhmmv_clear_level(mlhmmv_level_t* plevel);

// For 'emit' and 'emitp' in the DSL. These allocate lrecs, appended to the poutrecs list.
// * pmap is the base-level oosvar multi-level hashmap.
// * pkeys specify the level in the mlhmmv at which to produce data.
// * pnames is used to pull subsequent-level keys out into separate fields.
// * In case pnames isn't long enough to reach a terminal mlrval level in the mlhmmv,
//   do_full_prefixing specifies whether to concatenate nested mlhmmv keys into single lrec keys.
//
// Examples:

// * pkeys reaches a terminal level:
//
//   $ mlr --opprint put -q '@sum += $x; end { emit @sum }' ../data/small
//   sum
//   4.536294

// * pkeys reaches terminal levels:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emit @sum, "a", "b" }' ../data/small
//   a   b   sum
//   pan pan 0.346790
//   pan wye 0.502626
//   eks pan 0.758680
//   eks wye 0.381399
//   eks zee 0.611784
//   wye wye 0.204603
//   wye pan 0.573289
//   zee pan 0.527126
//   zee wye 0.598554
//   hat wye 0.031442

// * pkeys reaches non-terminal levels: non-prefixed:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emit @sum, "a" }' ../data/small
//   a   pan      wye
//   pan 0.346790 0.502626
//
//   a   pan      wye      zee
//   eks 0.758680 0.381399 0.611784
//
//   a   wye      pan
//   wye 0.204603 0.573289
//
//   a   pan      wye
//   zee 0.527126 0.598554
//
//   a   wye
//   hat 0.031442

// * pkeys reaches non-terminal levels: prefixed:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emitp @sum, "a" }' ../data/small
//   a   sum:pan  sum:wye
//   pan 0.346790 0.502626
//
//   a   sum:pan  sum:wye  sum:zee
//   eks 0.758680 0.381399 0.611784
//
//   a   sum:wye  sum:pan
//   wye 0.204603 0.573289
//
//   a   sum:pan  sum:wye
//   zee 0.527126 0.598554
//
//   a   sum:wye
//   hat 0.031442

void mlhmmv_to_lrecs(mlhmmv_t* pmap, sllmv_t* pkeys, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator);
void mlhmmv_to_lrecs_lashed(mlhmmv_t* pmap, sllmv_t** ppkeys, int num_keylists, sllmv_t* pnames,
	sllv_t* poutrecs, int do_full_prefixing, char* flatten_separator);

// For 'emit all' and 'emitp all' in the DSL
void mlhmmv_all_to_lrecs(mlhmmv_t* pmap, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator);

// For 'dump' in the DSL; also used by the lrec-to-JSON writer.
void mlhmmv_print_json_stacked(mlhmmv_t* pmap, int quote_values_always, char* line_indent, FILE* ostream);
void mlhmmv_print_json_single_line(mlhmmv_t* pmap, int quote_values_always, FILE* ostream);

void mlhmmv_level_print_stacked(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_values_always, char* line_indent, FILE* ostream);

#endif // MLHMMV_H
