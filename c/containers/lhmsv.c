// ================================================================
// Array-only (open addressing) string-to-void linked hash map with linear
// probing for collisions.
//
// John Kerl 2012-08-13
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

#include "lib/mlrutil.h"
#include "containers/lhmsv.h"

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
static void lhmsv_put_no_enlarge(lhmsv_t* pmap, char* key, void* pvvalue);
static void lhmsv_enlarge(lhmsv_t* pmap);

static void lhmsv_init(lhmsv_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhmsve_t*)mlr_malloc_or_die(sizeof(lhmsve_t) * length);
	// Don't do lhmsve_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.
	//for (int i = 0; i < length; i++)
	//	lhmsve_clear(&entries[i]);

	pmap->states       = (lhmsve_state_t*)mlr_malloc_or_die(sizeof(lhmsve_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead        = NULL;
	pmap->ptail        = NULL;
}

lhmsv_t* lhmsv_alloc() {
	lhmsv_t* pmap = mlr_malloc_or_die(sizeof(lhmsv_t));
	lhmsv_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

void lhmsv_free(lhmsv_t* pmap) {
	if (pmap == NULL)
		return;
	for (lhmsve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		free(pe->key);
	}
	free(pmap->entries);
	pmap->entries      = NULL;
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = 0;
	free(pmap);
}

// ----------------------------------------------------------------
// Used by get() and remove().
// Returns >0 for where the key is *or* should go (end of chain).
static int lhmsv_find_index_for_key(lhmsv_t* pmap, char* key) {
	int hash = mlr_string_hash_func(key);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	int num_tries = 0;
	int done = 0;

	while (!done) {
		lhmsve_t* pe = &pmap->entries[index];
		if (pmap->states[index] == OCCUPIED) {
			char* ekey = pe->key;
			// Existing key found in chain.
			if (streq(key, ekey))
				return index;
		}
		else if (pmap->states[index] == EMPTY) {
			return index;
		}

		// If the current entry has been deleted, i.e. previously occupied,
		// the sought index may be further down the chain.  So we must
		// continue looking.
		if (++num_tries >= pmap->array_length) {
			fprintf(stderr,
				"Coding error:  table full even after enlargement.\n");
			exit(1);
		}

		// Linear probing.
		if (++index >= pmap->array_length)
			index = 0;
	}
	return -1; // xxx not reached
}

// ----------------------------------------------------------------
void lhmsv_put(lhmsv_t* pmap, char* key, void* pvvalue) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhmsv_enlarge(pmap);
	lhmsv_put_no_enlarge(pmap, key, pvvalue);
}

static void lhmsv_put_no_enlarge(lhmsv_t* pmap, char* key, void* pvvalue) {
	int index = lhmsv_find_index_for_key(pmap, key);
	lhmsve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (streq(pe->key, key)) {
			pe->pvvalue = pvvalue;
		}
	}
	else if (pmap->states[index] == EMPTY) {
		// End of chain.
		pe->ideal_index = mlr_canonical_mod(mlr_string_hash_func(key), pmap->array_length);
		// xxx comment all malloced.
		pe->key = strdup(key);
		pe->pvvalue = pvvalue;
		pmap->states[index] = OCCUPIED;

		if (pmap->phead == NULL) {
			pe->pprev   = NULL;
			pe->pnext   = NULL;
			pmap->phead = pe;
			pmap->ptail = pe;
		} else {
			pe->pprev   = pmap->ptail;
			pe->pnext   = NULL;
			pmap->ptail->pnext = pe;
			pmap->ptail = pe;
		}
		pmap->num_occupied++;
	}
	else {
		fprintf(stderr, "lhmsv_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void* lhmsv_get(lhmsv_t* pmap, char* key) {
	int index = lhmsv_find_index_for_key(pmap, key);
	lhmsve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe->pvvalue;
	else if (pmap->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "lhmsv_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
int  lhmsv_has_key(lhmsv_t* pmap, char* key) {
	int index = lhmsv_find_index_for_key(pmap, key);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "lhmsv_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void lhmsv_remove(lhmsv_t* pmap, char* key) {
	int index = lhmsv_find_index_for_key(pmap, key);
	lhmsve_t* pe = &pmap->entries[index];
	if (pmap->states[index] == OCCUPIED) {
		pe->ideal_index     = -1;
		pe->key             = NULL;
		pe->pvvalue         = NULL;
		pmap->states[index] = DELETED;

		if (pe == pmap->phead) {
			if (pe == pmap->ptail) {
				pmap->phead = NULL;
				pmap->ptail = NULL;
			} else {
				pmap->phead = pe->pnext;
			}
		} else {
			pe->pprev->pnext = pe->pnext;
			pe->pnext->pprev = pe->pprev;
		}

		pmap->num_freed++;
		pmap->num_occupied--;
		return;
	}
	else if (pmap->states[index] == EMPTY) {
		return;
	}
	else {
		fprintf(stderr, "lhmsv_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
static void lhmsv_enlarge(lhmsv_t* pmap) {
	lhmsve_t*       old_entries = pmap->entries;
	lhmsve_state_t* old_states  = pmap->states;
	lhmsve_t*       old_head    = pmap->phead;

	lhmsv_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhmsve_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		lhmsv_put_no_enlarge(pmap, pe->key, pe->pvvalue);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
void lhmsv_check_counts(lhmsv_t* pmap) {
	int nocc = 0;
	int ndel = 0;
	for (int index = 0; index < pmap->array_length; index++) {
		if (pmap->states[index] == OCCUPIED)
			nocc++;
		else if (pmap->states[index] == DELETED)
			ndel++;
	}
	if (nocc != pmap->num_occupied) {
		fprintf(stderr,
			"occupancy-count mismatch:  actual %d != cached  %d.\n",
				nocc, pmap->num_occupied);
		exit(1);
	}
	if (ndel != pmap->num_freed) {
		fprintf(stderr,
			"deleted-count mismatch:  actual %d != cached  %d.\n",
				ndel, pmap->num_freed);
		exit(1);
	}
}

// ----------------------------------------------------------------
static char* get_state_name(int state) {
	switch(state) {
	case OCCUPIED: return "occupied"; break;
	case DELETED:  return "deleted";  break;
	case EMPTY:    return "empty";    break;
	default:       return "?????";    break;
	}
}

void lhmsv_dump(lhmsv_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhmsve_t* pe = &pmap->entries[index];

		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | pvvalue: %p |\n",
			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, pe->pvvalue);
	}
	printf("+\n");
	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
	printf("+\n");
	for (lhmsve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | pvvalue: %p |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key_string, pe->pvvalue);
	}
}
