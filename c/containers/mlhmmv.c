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

static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pideal_index);

static void mlhmmv_level_put(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);
static void mlhmmv_level_put_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);
static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel);
static void mlhmmv_level_move(mlhmmv_level_t* plevel, mv_t* plevel_key, mlhmmv_level_value_t* plevel_value);

static mlhmmv_level_value_t* mlhmmv_level_get(mlhmmv_level_t* pmap, sllmve_t* prest_keys);

static void mlhmmv_level_print_stacked(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_values_always);
static void mlhmmv_level_print_single_line(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_values_always);

static int mlhmmv_hash_func(mv_t* plevel_key);

// ----------------------------------------------------------------
// Allow compile-time override, e.g using gcc -D.

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
	mlhmmv_level_init(plevel, MLHMMV_INITIAL_ARRAY_LENGTH);
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
static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pideal_index) {
	int hash = mlhmmv_hash_func(plevel_key);
	int index = mlr_canonical_mod(hash, plevel->array_length);
	*pideal_index = index;
	int num_tries = 0;

	while (TRUE) {
		mlhmmv_level_entry_t* pentry = &plevel->entries[index];
		if (plevel->states[index] == OCCUPIED) {
			mv_t* ekey = &pentry->level_key;
			// Existing key found in chain.
			if (mv_equals_si(plevel_key, ekey))
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
	mv_t* plevel_key = &prest_keys->value;
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == EMPTY) { // End of chain.
		pentry->ideal_index = ideal_index;
		pentry->level_key = mv_copy(plevel_key);

		if (prest_keys->pnext == NULL) {
			pentry->level_value.is_terminal = TRUE;
			pentry->level_value.u.mlrval = mv_copy(pterminal_value);
		} else {
			pentry->level_value.is_terminal = FALSE;
			pentry->level_value.u.pnext_level = mlhmmv_level_alloc();
		}
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) { // First entry at this level
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {                     // Subsequent entry at this level
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}

		plevel->num_occupied++;
		if (prest_keys->pnext != NULL) {
			// RECURSE
			mlhmmv_level_put(pentry->level_value.u.pnext_level, prest_keys->pnext, pterminal_value);
		}

	} else if (plevel->states[index] == OCCUPIED) { // Existing key found in chain
		if (prest_keys->pnext == NULL) { // Place the terminal at this level
			if (pentry->level_value.is_terminal) {
				mv_free(&pentry->level_value.u.mlrval);
			} else {
				mlhmmv_level_free(pentry->level_value.u.pnext_level);
			}
			pentry->level_value.is_terminal = TRUE;
			pentry->level_value.u.mlrval = mv_copy(pterminal_value);

		} else { // The terminal will be placed at a deeper level
			if (pentry->level_value.is_terminal) {
				mv_free(&pentry->level_value.u.mlrval);
				pentry->level_value.is_terminal = FALSE;
				pentry->level_value.u.pnext_level = mlhmmv_level_alloc();
			}
			// RECURSE
			mlhmmv_level_put(pentry->level_value.u.pnext_level, prest_keys->pnext, pterminal_value);
		}

	} else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
		exit(1);
	}
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
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		pentry->level_value = *plevel_value;

	} else if (plevel->states[index] == EMPTY) {
		// End of chain.
		pentry->ideal_index = ideal_index;
		pentry->level_key = *plevel_key;
		// For the put API, we copy data passed in. But for internal enlarges, we just need to move pointers around.
		pentry->level_value = *plevel_value;
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) {
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}
		plevel->num_occupied++;
	}
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllmv_t* pmvkeys, int* perror) {
	*perror = MLHMMV_ERROR_NONE;
	sllmve_t* prest_keys = pmvkeys->phead;
	if (prest_keys == NULL) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	mlhmmv_level_t* plevel = pmap->proot_level;
	mlhmmv_level_value_t* plevel_value = mlhmmv_level_get(plevel, prest_keys);
	while (prest_keys->pnext != NULL) {
		if (plevel_value == NULL) {
			return NULL;
		}
		if (plevel_value->is_terminal) {
			*perror = MLHMMV_ERROR_KEYLIST_TOO_DEEP;
			return NULL;
		}
		plevel = plevel_value->u.pnext_level;
		prest_keys = prest_keys->pnext;
		plevel_value = mlhmmv_level_get(plevel, prest_keys);
	}
	if (plevel_value == NULL) {
		return NULL;
	}
	if (!plevel_value->is_terminal) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	return &plevel_value->u.mlrval;
}

static mlhmmv_level_value_t* mlhmmv_level_get(mlhmmv_level_t* plevel, sllmve_t* prest_keys) {
	mv_t* plevel_key = &prest_keys->value;
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == OCCUPIED)
		return &pentry->level_value;
	else if (plevel->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.argv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
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
// This is simply JSON. Example output:
// {
//   "0":  {
//     "fghij":  {
//       "0":  17
//     }
//   },
//   "3":  4,
//   "abcde":  {
//     "-6":  7
//   }
// }

void mlhmmv_print_stacked(mlhmmv_t* pmap, int quote_values_always) {
	mlhmmv_level_print_stacked(pmap->proot_level, 0, FALSE, quote_values_always);
}

static void mlhmmv_level_print_stacked(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_values_always)
{
	static char* leader = "  ";
	// Top-level opening brace goes on a line by itself; subsequents on the same line after the level key.
	if (depth == 0)
		printf("{\n");
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		for (int i = 0; i <= depth; i++)
			printf("%s", leader);
		char* level_key_string = mv_alloc_format_val(&pentry->level_key);
			printf("\"%s\": ", level_key_string);
		free(level_key_string);

		if (pentry->level_value.is_terminal) {
			char* level_value_string = mv_alloc_format_val(&pentry->level_value.u.mlrval);

			if (quote_values_always) {
				printf("\"%s\"", level_value_string);
			} else if (pentry->level_value.u.mlrval.type == MT_STRING) {
				double unused;
				if (mlr_try_float_from_string(level_value_string, &unused))
					printf("%s", level_value_string);
				else
					printf("\"%s\"", level_value_string);
			} else {
				printf("%s", level_value_string);
			}

			free(level_value_string);
			if (pentry->pnext != NULL)
				printf(",\n");
			else
				printf("\n");
		} else {
			printf("{\n");
			mlhmmv_level_print_stacked(pentry->level_value.u.pnext_level, depth + 1,
				pentry->pnext != NULL, quote_values_always);
		}
	}
	for (int i = 0; i < depth; i++)
		printf("%s", leader);
	if (do_final_comma)
		printf("},\n");
	else
		printf("}\n");
}

// ----------------------------------------------------------------
void mlhmmv_print_single_line(mlhmmv_t* pmap, int quote_values_always) {
	mlhmmv_level_print_single_line(pmap->proot_level, 0, FALSE, quote_values_always);
	printf("\n");
}

static void mlhmmv_level_print_single_line(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_values_always)
{
	// Top-level opening brace goes on a line by itself; subsequents on the same line after the level key.
	if (depth == 0)
		printf("{ ");
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		char* level_key_string = mv_alloc_format_val(&pentry->level_key);
			printf("\"%s\": ", level_key_string);
		free(level_key_string);

		if (pentry->level_value.is_terminal) {
			char* level_value_string = mv_alloc_format_val(&pentry->level_value.u.mlrval);

			if (quote_values_always) {
				printf("\"%s\"", level_value_string);
			} else if (pentry->level_value.u.mlrval.type == MT_STRING) {
				double unused;
				if (mlr_try_float_from_string(level_value_string, &unused))
					printf("%s", level_value_string);
				else
					printf("\"%s\"", level_value_string);
			} else {
				printf("%s", level_value_string);
			}

			free(level_value_string);
			if (pentry->pnext != NULL)
				printf(", ");
		} else {
			printf("{");
			mlhmmv_level_print_single_line(pentry->level_value.u.pnext_level, depth + 1,
				pentry->pnext != NULL, quote_values_always);
		}
	}
	if (do_final_comma)
		printf(" },");
	else
		printf(" }");
}

// ----------------------------------------------------------------
typedef int mlhmmv_typed_hash_func(mv_t* pa);

static int mlhmmv_string_hash_func(mv_t* pa) {
	return mlr_string_hash_func(pa->u.strv);
}
static int mlhmmv_int_hash_func(mv_t* pa) {
	return pa->u.intv;
}
static int mlhmmv_other_hash_func(mv_t* pa) {
	fprintf(stderr, "%s: @-variable keys must be of type %s or %s; got %s.\n",
		MLR_GLOBALS.argv0,
		mt_describe_type(MT_STRING),
		mt_describe_type(MT_INT),
		mt_describe_type(pa->type));
	exit(1);
}
static mlhmmv_typed_hash_func* mlhmmv_hash_func_dispositions[MT_MAX] = {
	/*NULL*/   mlhmmv_other_hash_func,
	/*ERROR*/  mlhmmv_other_hash_func,
	/*BOOL*/   mlhmmv_other_hash_func,
	/*FLOAT*/  mlhmmv_other_hash_func,
	/*INT*/    mlhmmv_int_hash_func,
	/*STRING*/ mlhmmv_string_hash_func,
};

static int mlhmmv_hash_func(mv_t* pa) {
	return mlhmmv_hash_func_dispositions[pa->type](pa);
}
