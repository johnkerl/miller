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

#include "lib/mlrval.h"
#include "containers/sllmv.h"
#include "containers/sllv.h"
#include "containers/lrec.h"

#define MLHMMV_ERROR_NONE                0x0000
#define MLHMMV_ERROR_KEYLIST_TOO_DEEP    0xdeef
#define MLHMMV_ERROR_KEYLIST_TOO_SHALLOW 0x58a1

// This is made visible here in the API so the unit-tester can be sure to exercise the resize logic.
#define MLHMMV_INITIAL_ARRAY_LENGTH 16

// ----------------------------------------------------------------
void mlhmmv_print_terminal(mv_t* pmv, int quote_keys_always, int quote_values_always, FILE* ostream);

// ----------------------------------------------------------------
struct _mlhmmv_level_t; // forward reference

// The 'x' is for extended: this can hold a scalar or a map.
typedef struct _mlhmmv_xvalue_t {
	struct _mlhmmv_level_t* pnext_level;
	mv_t terminal_mlrval;
	char is_terminal;
} mlhmmv_xvalue_t;

void            mlhmmv_xvalue_reset(mlhmmv_xvalue_t* pxvalue);
mlhmmv_xvalue_t mlhmmv_xvalue_alloc_empty_map();
mlhmmv_xvalue_t mlhmmv_xvalue_copy(mlhmmv_xvalue_t* pxvalue);
void            mlhmmv_xvalue_free(mlhmmv_xvalue_t* pxvalue);

char* mlhmmv_xvalue_describe_type_simple(mlhmmv_xvalue_t* pxvalue);

static inline int mlhmmv_xvalue_is_absent_and_nonterminal(mlhmmv_xvalue_t* pxvalue) {
	return (pxvalue->is_terminal && mv_is_absent(&pxvalue->terminal_mlrval));
}

static inline int mlhmmv_xvalue_is_present_and_nonterminal(mlhmmv_xvalue_t* pxvalue) {
	return (pxvalue->is_terminal && mv_is_present(&pxvalue->terminal_mlrval));
}

// Used by for-loops over map-valued local variables
sllv_t* mlhmmv_xvalue_copy_keys_indexed   (mlhmmv_xvalue_t* pxvalue, sllmv_t* pmvkeys);
sllv_t* mlhmmv_xvalue_copy_keys_nonindexed(mlhmmv_xvalue_t* pxvalue);

void mlhmmv_xvalues_to_lrecs_lashed(
	mlhmmv_xvalue_t** ptop_values,
	int               num_submaps,
	mv_t*             pbasenames,
	sllmv_t*          pnames,
	sllv_t*           poutrecs,
	int               do_full_prefixing,
	char*             flatten_separator);

// ----------------------------------------------------------------
typedef struct _mlhmmv_level_entry_t {
	int     ideal_index;
	mv_t    level_key;
	mlhmmv_xvalue_t level_xvalue; // terminal mlrval, or another hashmap
	struct _mlhmmv_level_entry_t *pprev;
	struct _mlhmmv_level_entry_t *pnext;
} mlhmmv_level_entry_t;

typedef unsigned char mlhmmv_level_entry_state_t;

// Store a mlrval into the mlhmmv_xvalue without copying, implicitly transferring
// ownership of the mlrval's free_flags. This means the mlrval will be freed
// when the mlhmmv_xvalue is freed, so the caller should make a copy first if
// necessary.
//
// This is a hot path for non-map local-variable assignments.
static inline mlhmmv_xvalue_t mlhmmv_xvalue_wrap_terminal(mv_t val) {
	return (mlhmmv_xvalue_t) {.is_terminal = TRUE, .terminal_mlrval = val, .pnext_level = NULL};
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

mlhmmv_level_t* mlhmmv_level_alloc();
void mlhmmv_level_free(mlhmmv_level_t* plevel);

void mlhmmv_level_clear(mlhmmv_level_t* plevel);
void mlhmmv_level_remove(mlhmmv_level_t* plevel, sllmve_t* prestkeys);
int mlhmmv_level_has_key(mlhmmv_level_t* plevel, mv_t* plevel_key);

mv_t* mlhmmv_level_look_up_and_ref_terminal(
	mlhmmv_level_t* plevel,
	sllmv_t*        pmvkeys,
	int*            perror);

mlhmmv_xvalue_t* mlhmmv_level_look_up_and_ref_xvalue(
	mlhmmv_level_t* plevel,
	sllmv_t*        pmvkeys,
	int*            perror);

mlhmmv_level_t* mlhmmv_level_put_empty_map(
	mlhmmv_level_t* plevel,
	mv_t*           pkey);

void mlhmmv_level_put_xvalue(
	mlhmmv_level_t*  plevel,
	sllmve_t*        prest_keys,
	mlhmmv_xvalue_t* pvalue);

void mlhmmv_level_put_xvalue_singly_keyed(
	mlhmmv_level_t*  plevel,
	mv_t*            pkey,
	mlhmmv_xvalue_t* pvalue);

void mlhmmv_level_put_terminal(
	mlhmmv_level_t* plevel,
	sllmve_t*       prest_keys,
	mv_t*           pterminal_value);

void mlhmmv_level_put_terminal_singly_keyed(
	mlhmmv_level_t* plevel,
	mv_t*           pkey,
	mv_t*           pterminal_value);

void mlhmmv_level_to_lrecs(
	mlhmmv_level_t* plevel,
	sllmv_t*        pkeys,
	sllmv_t*        pnames,
	sllv_t*         poutrecs,
	int             do_full_prefixing,
	char*           flatten_separator);

void mlhmmv_level_print_stacked(
	mlhmmv_level_t* plevel,
	int             depth,
	int             do_final_comma,
	int             quote_keys_always,
	int             quote_values_always,
	char*           line_indent,
	char*           line_term,
	FILE*           ostream);

// ----------------------------------------------------------------
typedef struct _mlhmmv_root_t {
	mlhmmv_xvalue_t root_xvalue;
} mlhmmv_root_t;

mlhmmv_root_t* mlhmmv_root_alloc();

void mlhmmv_root_free(mlhmmv_root_t* pmap);

void mlhmmv_root_clear(mlhmmv_root_t* pmap);

// If the return value is non-null, error will be MLHMMV_ERROR_NONE.  If the
// return value is null, the error will be MLHMMV_ERROR_KEYLIST_TOO_DEEP or
// MLHMMV_ERROR_KEYLIST_TOO_SHALLOW, or MLHMMV_ERROR_NONE if the keylist matches
// map depth but the entry is not found.
//
// Note: this returns a pointer to the map's data, not to a copy.
// The caller shouldn't free it, or modify it.
mv_t* mlhmmv_root_look_up_and_ref_terminal(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, int* perror);

// These are an optimization for assignment from full srec, e.g. '@records[$key1][$key2] = $*'.
// Using mlhmmv_root_look_up_or_create_then_ref_level, the CST logic can get or create the @records[$key1][$key2]
// level of the mlhmmv, then copy values there.
mlhmmv_level_t* mlhmmv_root_look_up_or_create_then_ref_level(mlhmmv_root_t* pmap, sllmv_t* pmvkeys);

void mlhmmv_root_put_terminal(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, mv_t* pterminal_value);

// For for-loop-over-oosvar, wherein we need to copy the submap before iterating over it
// (since the iteration may modify it). If the keys don't index a submap, then the return
// value has is_terminal = TRUE and pnext_level = NULL.
mlhmmv_xvalue_t mlhmmv_root_copy_xvalue(mlhmmv_root_t* pmap, sllmv_t* pmvkeys);

// Used by for-loops over oosvars. Return value is an array of ephemeral mlrvals.
sllv_t* mlhmmv_root_copy_keys_from_submap(mlhmmv_root_t* pmap, sllmv_t* pmvkeys);

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
void mlhmmv_root_remove(mlhmmv_root_t* pmap, sllmv_t* pmvkeys);

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

// For 'emit all' and 'emitp all' in the DSL
void mlhmmv_root_all_to_lrecs(mlhmmv_root_t* pmap, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator);

// For 'emit' and 'emitp' in the DSL
void mlhmmv_root_partial_to_lrecs(mlhmmv_root_t* pmap, sllmv_t* pkeys, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator);

// For 'dump' in the DSL; also used by the lrec-to-JSON writer.
void mlhmmv_root_print_json_stacked(mlhmmv_root_t* pmap,
	int quote_keys_always, int quote_values_always, char* line_indent, char* line_term,
	FILE* ostream);
void mlhmmv_root_print_json_single_lines(mlhmmv_root_t* pmap, int quote_keys_always, int quote_values_always,
	char* line_term, FILE* ostream);

// Used for emit of localvars. Puts the xvalue in a single-key-value-pair map
// keyed by the specified name. The xvalue is referenced, not copied.
mlhmmv_root_t* mlhmmv_wrap_name_and_xvalue(mv_t* pname, mlhmmv_xvalue_t* pxval);

// Used for takedown of the temporary map returned by mlhmmv_wrap_name_and_xvalue. Since the xvalue there
// is referenced, not copied, mlhmmv_xvalue_free would prematurely free the xvalue. This method releases
// the xvalue so that the remaining, map-internal structures can be freed correctly.
void mlhmmv_unwrap_name_and_xvalue(mlhmmv_root_t* pmap);

#endif // MLHMMV_H
