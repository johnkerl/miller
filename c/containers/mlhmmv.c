// ================================================================
// Array-only (open addressing) multi-level hash map, with linear probing for collisions.
// All keys, and terminal-level values, are mlrvals.
//
// Notes:
// * null key is not supported.
// * null value is not supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"

static mlhmmv_level_t* mlhmmv_level_alloc();
static void            mlhmmv_level_init(mlhmmv_level_t *plevel, int length);
static void            mlhmmv_level_free(mlhmmv_level_t* plevel);

static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key);

static void mlhmmv_level_put(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);
static void mlhmmv_level_put_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);
static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel);
static void mlhmmv_level_move(mlhmmv_level_t* plevel, mv_t* plevel_key, mlhmmv_level_value_t* plevel_value);

static void mlhmmv_level_print(mlhmmv_level_t* plevel, int depth);

static int mlhmmv_hash_func(mv_t* plevel_key);
static int mlhmmv_key_equals(mv_t* pa, mv_t* pb);

// ----------------------------------------------------------------
// Allow compile-time override, e.g using gcc -D.
#ifndef INITIAL_ARRAY_LENGTH
#define INITIAL_ARRAY_LENGTH 16
#endif

#ifndef LOAD_FACTOR
#define LOAD_FACTOR          0.7
#endif

#ifndef ENLARGEMENT_FACTOR
#define ENLARGEMENT_FACTOR   2
#endif

#define OCCUPIED             0xa4
#define DELETED              0xb8
#define EMPTY                0xce

// ----------------------------------------------------------------
mlhmmv_t* mlhmmv_alloc() {
	mlhmmv_t* pmap = mlr_malloc_or_die(sizeof(mlhmmv_t));
	pmap->proot_level = mlhmmv_level_alloc();
	return pmap;
}

static mlhmmv_level_t* mlhmmv_level_alloc() {
	mlhmmv_level_t* plevel = mlr_malloc_or_die(sizeof(mlhmmv_level_t));
	mlhmmv_level_init(plevel, INITIAL_ARRAY_LENGTH);
	return plevel;
}

static void mlhmmv_level_init(mlhmmv_level_t *plevel, int length) {
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->array_length = length;

	plevel->entries      = (mlhmmv_level_entry_t*)mlr_malloc_or_die(sizeof(mlhmmv_level_entry_t) * length);
	// Don't do mlhmmv_level_entry_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.

	plevel->states       = (mlhmmv_level_entry_state_t*)mlr_malloc_or_die(sizeof(mlhmmv_level_entry_state_t) * length);
	memset(plevel->states, EMPTY, length);

	plevel->phead        = NULL;
	plevel->ptail        = NULL;
}

// ----------------------------------------------------------------
void mlhmmv_free(mlhmmv_t* pmap) {
	if (pmap == NULL)
		return;
	mlhmmv_level_free(pmap->proot_level);
	free(pmap);
}

static void mlhmmv_level_free(mlhmmv_level_t* plevel) {
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		if (!pentry->level_value.is_terminal) {
			mlhmmv_level_free(pentry->level_value.u.pnext_level);
		}
		mv_free(&pentry->level_key);
	}
	free(plevel->entries);
	free(plevel->states);
	plevel->entries      = NULL;
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->array_length = 0;

	free(plevel);
}

// ----------------------------------------------------------------
// Used by get() and remove().
// Returns >0 for where the key is *or* should go (end of chain).
static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key) {
	int hash = mlhmmv_hash_func(plevel_key);
	int index = mlr_canonical_mod(hash, plevel->array_length);
	int num_tries = 0;

	while (TRUE) {
		mlhmmv_level_entry_t* pentry = &plevel->entries[index];
		if (plevel->states[index] == OCCUPIED) {
			mv_t* ekey = &pentry->level_key;
			// Existing key found in chain.
			if (mlhmmv_key_equals(plevel_key, ekey))
				return index;
		} else if (plevel->states[index] == EMPTY) {
			return index;
		}

		// If the current entry has been freed, i.e. previously occupied,
		// the sought index may be further down the chain.  So we must
		// continue looking.
		if (++num_tries >= plevel->array_length) {
			fprintf(stderr,
				"%s: Coding error:  table full even after enlargement.\n", MLR_GLOBALS.argv0);
			exit(1);
		}

		// Linear probing.
		if (++index >= plevel->array_length)
			index = 0;
	}
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.argv0, __FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
// Example: keys = ["a", 2, "c"] and value = 4.
void mlhmmv_put(mlhmmv_t* pmap, sllmv_t* pmvkeys, mv_t* pterminal_value) {
	mlhmmv_level_put(pmap->proot_level, pmvkeys->phead, pterminal_value);
}
// Example on recursive calls:
// * level = map, rest_keys = ["a", 2, "c"] , terminal value = 4.
// * level = map["a"], rest_keys = [2, "c"] , terminal value = 4.
// * level = map["a"][2], rest_keys = ["c"] , terminal value = 4.
static void mlhmmv_level_put(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value) {
	if ((plevel->num_occupied + plevel->num_freed) >= (plevel->array_length * LOAD_FACTOR))
		mlhmmv_level_enlarge(plevel);
	mlhmmv_level_put_no_enlarge(plevel, prest_keys, pterminal_value);
}

static void mlhmmv_level_put_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value) {
	mv_t* plevel_key = prest_keys->pvalue;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (mlhmmv_key_equals(&pentry->level_key, plevel_key)) {
			if (pentry->level_value.is_terminal)
				mv_free(&pentry->level_value.u.mlrval);
			else
				// xxx no, not always
				mlhmmv_level_free(pentry->level_value.u.pnext_level);
			if (prest_keys->pnext == NULL) {
				pentry->level_value.is_terminal = TRUE;
				pentry->level_value.u.mlrval = *pterminal_value;
			} else {
				pentry->level_value.is_terminal = FALSE;
				pentry->level_value.u.pnext_level = mlhmmv_level_alloc();
				mlhmmv_level_put(pentry->level_value.u.pnext_level, prest_keys->pnext, pterminal_value);
			}
			return;
			// xxx make sllmve_t pmv -> mv !
		}
	} else if (plevel->states[index] == EMPTY) {
		// End of chain.
		pentry->ideal_index = mlr_canonical_mod(mlhmmv_hash_func(plevel_key), plevel->array_length);
		pentry->level_key = *plevel_key;

		if (prest_keys->pnext == NULL) {
			pentry->level_value.is_terminal = TRUE;
			pentry->level_value.u.mlrval = *pterminal_value;
		} else {
			pentry->level_value.is_terminal = FALSE;
			pentry->level_value.u.pnext_level = mlhmmv_level_alloc();
			mlhmmv_level_put(pentry->level_value.u.pnext_level, prest_keys->pnext, pterminal_value);
		}
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) {
			pentry->pprev   = NULL;
			pentry->pnext   = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {
			pentry->pprev   = plevel->ptail;
			pentry->pnext   = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}
		plevel->num_occupied++;
		return;
	}
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
		exit(1);
	}
	// This one is to appease a compiler warning about control reaching the end
	// of a non-void function
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.argv0, __FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
// This is done only on map-level enlargement.
// Example:
// * level = map["a"], rest_keys = [2, "c"] ,   terminal_value = 4.
//                     rest_keys = ["e", "f"] , terminal_value = 7.
//                     rest_keys = [6] ,        terminal_value = "g".
//
// which is to say for the purposes of this routine
//
// * level = map["a"], level_key = 2,   level_value = non-terminal ["c"] => terminal_value = 4.
//                     level_key = "e", level_value = non-terminal ["f"] => terminal_value = 7.
//                     level_key = 6,   level_value = terminal_value = "g".

static void mlhmmv_level_move(mlhmmv_level_t* plevel, mv_t* plevel_key, mlhmmv_level_value_t* plevel_value) {
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (mlhmmv_key_equals(&pentry->level_key, plevel_key)) {
			pentry->level_value = *plevel_value;
			return;
		}
	} else if (plevel->states[index] == EMPTY) {
		// End of chain.
		pentry->ideal_index = mlr_canonical_mod(mlhmmv_hash_func(plevel_key), plevel->array_length);
		pentry->level_key = *plevel_key;
		// For the put API, we copy data passed in. But for internal enlarges, we just need to move pointers around.
		pentry->level_value = *plevel_value;
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) {
			pentry->pprev   = NULL;
			pentry->pnext   = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {
			pentry->pprev   = plevel->ptail;
			pentry->pnext   = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}
		plevel->num_occupied++;
		return;
	}
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
		exit(1);
	}
	// This one is to appease a compiler warning about control reaching the end
	// of a non-void function
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.argv0, __FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllmv_t* pmvkeys) {
	return NULL; // xxx stub
}
//	int index = mlhmmv_level_find_index_for_key(plevel, level_key);
//	mlhmmv_level_entry_t* pentry = &pmap->entries[index];
//
//	if (pmap->states[index] == OCCUPIED)
//		return pentry->pvalue;
//	else if (pmap->states[index] == EMPTY)
//		return NULL;
//	else {
//		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
//		exit(1);
//	}

// ----------------------------------------------------------------
int mlhmmv_has_keys(mlhmmv_t* pmap, sllmv_t* pmvkeys) {
//	int index = mlhmmv_level_find_index_for_key(plevel, level_key);
//
//	if (pmap->states[index] == OCCUPIED)
//		return TRUE;
//	else if (pmap->states[index] == EMPTY)
		return FALSE;
//	else {
//		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
//		exit(1);
//	}
}

// ----------------------------------------------------------------
// Example:
// * level = map["a"], rest_keys = [2, "c"] ,   terminal_value = 4.
//                     rest_keys = ["e", "f"] , terminal_value = 7.
//                     rest_keys = [6] ,        terminal_value = "g".
//
// which is to say for the purposes of this routine
//
// * level = map["a"], level_key = 2,   level_value = non-terminal ["c"] => terminal_value = 4.
//                     level_key = "e", level_value = non-terminal ["f"] => terminal_value = 7.
//                     level_key = 6,   level_value = terminal_value = "g".

static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel) {
	mlhmmv_level_entry_t*       old_entries = plevel->entries;
	mlhmmv_level_entry_state_t* old_states  = plevel->states;
	mlhmmv_level_entry_t*       old_head    = plevel->phead;

	mlhmmv_level_init(plevel, plevel->array_length*ENLARGEMENT_FACTOR);

	for (mlhmmv_level_entry_t* pentry = old_head; pentry != NULL; pentry = pentry->pnext) {
		mlhmmv_level_move(plevel, &pentry->level_key, &pentry->level_value);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
// Example output:
//
// {
//   0 => {
//     fghij => {
//       0 => 17
//     }
//   }
//   3 => 4
//   abcde => {
//     -6 => 7
//   }
// }

void mlhmmv_print(mlhmmv_t* pmap) {
	mlhmmv_level_print(pmap->proot_level, 0);
}

static void mlhmmv_level_print(mlhmmv_level_t* plevel, int depth) {
	static char* leader = "  ";
	// Top-level opening brace on a line by itself; subsequents on the same line after the level key.
	if (depth == 0)
		printf("{\n");
	for (int index = 0; index < plevel->array_length; index++) {
		if (plevel->states[index] != OCCUPIED)
			continue;

		mlhmmv_level_entry_t* pentry = &plevel->entries[index];

		for (int i = 0; i <= depth; i++)
			printf("%s", leader);
		char* level_key_string = mv_alloc_format_val(&pentry->level_key);
		printf("%s =>", level_key_string);
		free(level_key_string);

		if (pentry->level_value.is_terminal) {
			char* level_value_string = mv_alloc_format_val(&pentry->level_value.u.mlrval);
			printf(" %s\n", level_value_string);
			free(level_value_string);
		} else {
			printf(" {\n");
			mlhmmv_level_print(pentry->level_value.u.pnext_level, depth + 1);
		}
	}
	for (int i = 0; i < depth; i++)
		printf("%s", leader);
	printf("}\n");
}

// ----------------------------------------------------------------
typedef int mv_typed_hash_func_t(mv_t* pa);

static int mv_int_hash_func(mv_t* pa) {
	return pa->u.intv;
}
static int mv_string_hash_func(mv_t* pa) {
	return mlr_string_hash_func(pa->u.strv);
}
static int mv_other_hash_func(mv_t* pa) {
	fprintf(stderr, "%s: mlhmmv: cannot hash %s, only %s or %s.\n",
		MLR_GLOBALS.argv0,
		mt_describe_type(pa->type),
		mt_describe_type(MT_STRING),
		mt_describe_type(MT_INT));
	exit(1);
}

static mv_typed_hash_func_t* hash_func_dispositions[MT_MAX] = {
	/*NULL*/   mv_other_hash_func,
	/*ERROR*/  mv_other_hash_func,
	/*BOOL*/   mv_other_hash_func,
	/*FLOAT*/  mv_other_hash_func,
	/*INT*/    mv_int_hash_func,
	/*STRING*/ mv_string_hash_func,
};

static int mlhmmv_hash_func(mv_t* pa) {
	return (hash_func_dispositions[pa->type])(pa);
}

static int mlhmmv_key_equals(mv_t* pa, mv_t* pb) {
	if (pa->type == MT_INT) {
		return (pb->type == MT_INT) ? pa->u.intv == pb->u.intv : FALSE;
	} else {
		return (pb->type == MT_STRING) ? streq(pa->u.strv, pb->u.strv) : FALSE;
	}
}
