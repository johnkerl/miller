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
//
//// ----------------------------------------------------------------
//// Allow compile-time override, e.g using gcc -D.
//#ifndef INITIAL_ARRAY_LENGTH
//#define INITIAL_ARRAY_LENGTH 16
//#endif
//
//#ifndef LOAD_FACTOR
//#define LOAD_FACTOR          0.7
//#endif
//
//#ifndef ENLARGEMENT_FACTOR
//#define ENLARGEMENT_FACTOR   2
//#endif
//
//// ----------------------------------------------------------------
//#define OCCUPIED 0xa4
//#define DELETED  0xb8
//#define EMPTY    0xce
//
//// ----------------------------------------------------------------
//static void* mlhmmv_put_no_enlarge(mlhmmv_t* pmap, slls_t* key, void* pvvalue, char free_flags);
//static void mlhmmv_enlarge(mlhmmv_t* pmap);
//
//// ================================================================
//static void mlhmmv_init(mlhmmv_t *pmap, int length) {
//	pmap->num_occupied = 0;
//	pmap->num_freed    = 0;
//	pmap->array_length = length;
//
//	pmap->entries      = (mlhmmve_t*)mlr_malloc_or_die(sizeof(mlhmmve_t) * length);
//	// Don't do mlhmmve_clear() of all entries at init time, since this has a
//	// drastic effect on the time needed to construct an empty map (and miller
//	// constructs an awful lot of those). The attributes there are don't-cares
//	// if the corresponding entry state is EMPTY. They are set on put, and
//	// mutated on remove.
//
//	pmap->states       = (mlhmmve_state_t*)mlr_malloc_or_die(sizeof(mlhmmve_state_t) * length);
//	memset(pmap->states, EMPTY, length);
//
//	pmap->phead        = NULL;
//	pmap->ptail        = NULL;
//}
//
mlhmmv_t* mlhmmv_alloc() {
	mlhmmv_t* pmap = mlr_malloc_or_die(sizeof(mlhmmv_t));
//	mlhmmv_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

void mlhmmv_free(mlhmmv_t* pmap) {
	if (pmap == NULL)
		return;
//	for (mlhmmve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext)
//		if (pe->free_flags & FREE_ENTRY_KEY)
//			slls_free(pe->key);
//	free(pmap->entries);
//	free(pmap->states);
//	pmap->entries      = NULL;
//	pmap->num_occupied = 0;
//	pmap->num_freed    = 0;
//	pmap->array_length = 0;
	free(pmap);
}

//// ----------------------------------------------------------------
//// Used by get() and remove().
//// Returns >0 for where the key is *or* should go (end of chain).
//static int mlhmmv_find_index_for_key(mlhmmv_t* pmap, slls_t* key) {
//	int hash = slls_hash_func(key);
//	int index = mlr_canonical_mod(hash, pmap->array_length);
//	int num_tries = 0;
//
//	while (TRUE) {
//		mlhmmve_t* pe = &pmap->entries[index];
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
void mlhmmv_put(mlhmmv_t* pmap, sllv_t* pmkeys, mv_t* pvalue) {
//	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
//		mlhmmv_enlarge(pmap);
//	return mlhmmv_put_no_enlarge(pmap, key, pvvalue, free_flags);
}

//static void* mlhmmv_put_no_enlarge(mlhmmv_t* pmap, slls_t* key, void* pvvalue, char free_flags) {
//	int index = mlhmmv_find_index_for_key(pmap, key);
//	mlhmmve_t* pe = &pmap->entries[index];
//
//	if (pmap->states[index] == OCCUPIED) {
//		// Existing key found in chain; put value.
//		if (slls_equals(pe->key, key)) {
//			pe->pvvalue = pvvalue;
//			return pvvalue;
//		}
//	}
//	else if (pmap->states[index] == EMPTY) {
//		// End of chain.
//		pe->ideal_index = mlr_canonical_mod(slls_hash_func(key), pmap->array_length);
//		pe->key = key;
//		pe->free_flags = free_flags;
//		pe->pvvalue = pvvalue;
//		pmap->states[index] = OCCUPIED;
//
//		if (pmap->phead == NULL) {
//			pe->pprev   = NULL;
//			pe->pnext   = NULL;
//			pmap->phead = pe;
//			pmap->ptail = pe;
//		} else {
//			pe->pprev   = pmap->ptail;
//			pe->pnext   = NULL;
//			pmap->ptail->pnext = pe;
//			pmap->ptail = pe;
//		}
//		pmap->num_occupied++;
//		return pvvalue;
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
//}

// ----------------------------------------------------------------
mv_t* mlhmmv_get(mlhmmv_t* pmap, sllv_t* pmvkeys) {
//	int index = mlhmmv_find_index_for_key(pmap, key);
//	mlhmmve_t* pe = &pmap->entries[index];
//
//	if (pmap->states[index] == OCCUPIED)
//		return pe->pvvalue;
//	else if (pmap->states[index] == EMPTY)
		return NULL;
//	else {
//		fprintf(stderr, "%s: mlhmmv_find_index_for_key did not find end of chain\n",
//			MLR_GLOBALS.argv0);
//		exit(1);
//	}
}

// ----------------------------------------------------------------
int mlhmmv_has_keys(mlhmmv_t* pmap, sllv_t* pmvkeys) {
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

//// ----------------------------------------------------------------
//int mlhmmv_size(mlhmmv_t* pmap) {
//	return pmap->num_occupied;
//}
//
//// ----------------------------------------------------------------
//static void mlhmmv_enlarge(mlhmmv_t* pmap) {
//	mlhmmve_t*       old_entries = pmap->entries;
//	mlhmmve_state_t* old_states  = pmap->states;
//	mlhmmve_t*       old_head    = pmap->phead;
//
//	mlhmmv_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);
//
//	for (mlhmmve_t* pe = old_head; pe != NULL; pe = pe->pnext) {
//		mlhmmv_put_no_enlarge(pmap, pe->key, pe->pvvalue, pe->free_flags);
//	}
//	free(old_entries);
//	free(old_states);
//}
//
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
//
//void mlhmmv_print(mlhmmv_t* pmap) {
//	for (int index = 0; index < pmap->array_length; index++) {
//		mlhmmve_t* pe = &pmap->entries[index];
//
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			slls_join(pe->key, ",");
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->pvvalue == NULL ? "null" :
//			pe->pvvalue;
//
//		printf(
//		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | pvvalue: %12s |\n",
//			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, value_string);
//	}
//	printf("+\n");
//	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
//	printf("+\n");
//	for (mlhmmve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			slls_join(pe->key, ",");
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->pvvalue == NULL ? "null" :
//			pe->pvvalue;
//		printf(
//		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | pvvalue: %12s |\n",
//			pe->pprev, pe, pe->pnext,
//			pe->ideal_index, key_string, value_string);
//	}
//}
