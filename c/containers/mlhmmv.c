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

//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
//
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"

static mlhmmv_level_t* mlhmmv_level_alloc();
static void            mlhmmv_level_init(mlhmmv_level_t *plevel, int length);
static void            mlhmmv_level_free(mlhmmv_level_t* plevel);

static void mlhmmv_put_aux(mlhmmv_level_t* plevel, sllmve_t* pmvkeys, mv_t* pvalue);
static void mlhmmv_level_put_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* pmvkeys, mv_t* pvalue, int do_copy);
static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel);

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

// ----------------------------------------------------------------
#define OCCUPIED 0xa4
#define DELETED  0xb8
#define EMPTY    0xce

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

	plevel->entries      = (mlhmmv_entry_t*)mlr_malloc_or_die(sizeof(mlhmmv_entry_t) * length);
	// Don't do mlhmmv_entry_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.

	plevel->states       = (mlhmmv_entry_state_t*)mlr_malloc_or_die(sizeof(mlhmmv_entry_state_t) * length);
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

	// xxx recurse

	for (mlhmmv_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext)
		mv_free(&pe->key);
	free(plevel->entries);
	free(plevel->states);
	plevel->entries      = NULL;
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->array_length = 0;

	free(plevel);
}

//// ----------------------------------------------------------------
//// Used by get() and remove().
//// Returns >0 for where the key is *or* should go (end of chain).
//static int mlhmmv_find_index_for_key(mlhmmv_t* pmap, slls_t* key) {
//	int hash = mlhmmv_hash_func(key);
//	int index = mlr_canonical_mod(hash, pmap->array_length);
//	int num_tries = 0;
//
//	while (TRUE) {
//		mlhmmv_entry_t* pe = &pmap->entries[index];
//		if (pmap->states[index] == OCCUPIED) {
//			slls_t* ekey = pe->key;
//			// Existing key found in chain.
//			if (slls_equals(key, ekey))
//				return index;
//		}
//		else if (pmap->states[index] == EMPTY) {
//			return index;
//		}
//
//		// If the current entry has been freed, i.e. previously occupied,
//		// the sought index may be further down the chain.  So we must
//		// continue looking.
//		if (++num_tries >= pmap->array_length) {
//			fprintf(stderr,
//				"Coding error:  table full even after enlargement.\n");
//			exit(1);
//		}
//
//		// Linear probing.
//		if (++index >= pmap->array_length)
//			index = 0;
//	}
//	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//		MLR_GLOBALS.argv0, __FILE__, __LINE__);
//	exit(1);
//}

// ----------------------------------------------------------------
void mlhmmv_put(mlhmmv_t* pmap, sllmv_t* pmvkeys, mv_t* pvalue) {
	mlhmmv_put_aux(pmap->proot_level, pmvkeys->phead, pvalue);
}
static void mlhmmv_put_aux(mlhmmv_level_t* plevel, sllmve_t* pmvkeys, mv_t* pvalue) {
	if ((plevel->num_occupied + plevel->num_freed) >= (plevel->array_length * LOAD_FACTOR))
		mlhmmv_level_enlarge(plevel);
	mlhmmv_level_put_no_enlarge(plevel, pmvkeys, pvalue, TRUE);
}


static void mlhmmv_level_put_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* pmvkeys, mv_t* pvalue, int do_copy) {
//	int index = mlhmmv_find_index_for_key(plevel, key);
//	mlhmmv_entry_t* pe = &plevel->entries[index];
//
//	if (plevel->states[index] == OCCUPIED) {
//		// Existing key found in chain; put value.
//		if (slls_equals(pe->key, key)) {
//			pe->pvalue = pvalue;
//			return pvalue;
//		}
//	}
//	else if (plevel->states[index] == EMPTY) {
//		// End of chain.
//		pe->ideal_index = mlr_canonical_mod(mlhmmv_hash_func(key), plevel->array_length);
//		pe->key = key;
//		// For the put API, we copy data passed in. But for internal enlarges, we just need to move pointers around.
//		pe->pvalue = do_copy ? mv_alloc_copy(pvalue) : pvalue;
//		plevel->states[index] = OCCUPIED;
//
//		if (plevel->phead == NULL) {
//			pe->pprev   = NULL;
//			pe->pnext   = NULL;
//			plevel->phead = pe;
//			plevel->ptail = pe;
//		} else {
//			pe->pprev   = plevel->ptail;
//			pe->pnext   = NULL;
//			plevel->ptail->pnext = pe;
//			plevel->ptail = pe;
//		}
//		plevel->num_occupied++;
//		return pvalue;
//	}
//	else {
//		fprintf(stderr, "mlhmmv_find_index_for_key did not find end of chain\n");
//		exit(1);
//	}
//	// This one is to appease a compiler warning about control reaching the end
//	// of a non-void function
//	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//		MLR_GLOBALS.argv0, __FILE__, __LINE__);
//	exit(1);
}

// ----------------------------------------------------------------
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllmv_t* pmvkeys) {
	return NULL; // xxx stub
}
//	int index = mlhmmv_find_index_for_key(pmap, key);
//	mlhmmv_entry_t* pe = &pmap->entries[index];
//
//	if (pmap->states[index] == OCCUPIED)
//		return pe->pvalue;
//	else if (pmap->states[index] == EMPTY)
//		return NULL;
//	else {
//		fprintf(stderr, "%s: mlhmmv_find_index_for_key did not find end of chain\n",
//			MLR_GLOBALS.argv0);
//		exit(1);
//	}

// ----------------------------------------------------------------
int mlhmmv_has_keys(mlhmmv_t* pmap, sllmv_t* pmvkeys) {
//	int index = mlhmmv_find_index_for_key(pmap, key);
//
//	if (pmap->states[index] == OCCUPIED)
//		return TRUE;
//	else if (pmap->states[index] == EMPTY)
		return FALSE;
//	else {
//		fprintf(stderr, "mlhmmv_find_index_for_key did not find end of chain\n");
//		exit(1);
//	}
}

// ----------------------------------------------------------------
static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel) {
//	mlhmmv_entry_t*       old_entries = plevel->entries;
//	mlhmmv_entry_state_t* old_states  = plevel->states;
//	mlhmmv_entry_t*       old_head    = plevel->phead;
//
//	mlhmmv_level_init(plevel, plevel->array_length*ENLARGEMENT_FACTOR);
//
//	for (mlhmmv_entry_t* pe = old_head; pe != NULL; pe = pe->pnext) {
//		mlhmmv_level_put_no_enlarge(plevel, pe->key, pe->pvalue, FALSE);
//	}
//	free(old_entries);
//	free(old_states);
}

//// ----------------------------------------------------------------
//int mlhmmv_check_counts(mlhmmv_t* pmap) {
//	int nocc = 0;
//	int ndel = 0;
//	for (int index = 0; index < pmap->array_length; index++) {
//		if (pmap->states[index] == OCCUPIED)
//			nocc++;
//		else if (pmap->states[index] == DELETED)
//			ndel++;
//	}
//	if (nocc != pmap->num_occupied) {
//		fprintf(stderr,
//			"occupancy-count mismatch:  actual %d != cached  %d\n",
//				nocc, pmap->num_occupied);
//		return FALSE;
//	}
//	if (ndel != pmap->num_freed) {
//		fprintf(stderr,
//			"freed-count mismatch:  actual %d != cached  %d\n",
//				ndel, pmap->num_freed);
//		return FALSE;
//	}
//	return TRUE;
//}
//
//// ----------------------------------------------------------------
//static char* get_state_name(int state) {
//	switch(state) {
//	case OCCUPIED: return "occupied"; break;
//	case DELETED:  return "freed";  break;
//	case EMPTY:    return "empty";    break;
//	default:       return "?????";    break;
//	}
//}

void mlhmmv_print(mlhmmv_t* pmap) {
//	for (int index = 0; index < pmap->array_length; index++) {
//		mlhmmv_entry_t* pe = &pmap->entries[index];
//
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			slls_join(pe->key, ",");
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->pvalue == NULL ? "null" :
//			pe->pvalue;
//
//		printf(
//		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | pvalue: %12s |\n",
//			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, value_string);
//	}
//	printf("+\n");
//	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
//	printf("+\n");
//	for (mlhmmv_entry_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			slls_join(pe->key, ",");
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->pvalue == NULL ? "null" :
//			pe->pvalue;
//		printf(
//		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | pvalue: %12s |\n",
//			pe->pprev, pe, pe->pnext,
//			pe->ideal_index, key_string, value_string);
//	}
}

// ----------------------------------------------------------------
mlhmmv_value_t* mlhmmv_value_from_mv(mv_t* pmv) {
	mlhmmv_value_t* pvalue = mlr_malloc_or_die(sizeof(mlhmmv_value_t));
	pvalue->type = MLHMMV_VALUE_TYPE_TERMINAL;
	pvalue->u.mlrval = mv_copy(pmv);
	return pvalue;
}
