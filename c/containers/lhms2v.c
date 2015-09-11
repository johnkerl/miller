// ================================================================
// Array-only (open addressing) string-list-to-void-star linked hash map with
// linear probing for collisions.
//
// John Kerl 2014-12-22
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
#include "containers/lhms2v.h"

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
static void* lhms2v_put_no_enlarge(lhms2v_t* pmap, char* key1, char* key2, void* pvvalue);
static void lhms2v_enlarge(lhms2v_t* pmap);

// ================================================================
static void lhms2ve_clear(lhms2ve_t *pentry) {
	pentry->ideal_index = -1;
	pentry->key1        = NULL;
	pentry->key2        = NULL;
	pentry->pvvalue     = NULL;
	pentry->pprev       = NULL;
	pentry->pnext       = NULL;
}

static void lhms2v_init(lhms2v_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhms2ve_t*)mlr_malloc_or_die(sizeof(lhms2ve_t) * length);
	// Don't do lhms2ve_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.
	//for (int i = 0; i < length; i++)
	//	lhms2ve_clear(&entries[i]);

	pmap->states       = (lhms2ve_state_t*)mlr_malloc_or_die(sizeof(lhms2ve_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead        = NULL;
	pmap->ptail        = NULL;
}

lhms2v_t* lhms2v_alloc() {
	lhms2v_t* pmap = mlr_malloc_or_die(sizeof(lhms2v_t));
	lhms2v_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

// xxx cmt re memmgt
void lhms2v_free(lhms2v_t* pmap) {
	if (pmap == NULL)
		return;
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
static int lhms2v_find_index_for_key(lhms2v_t* pmap, char* key1, char* key2) {
	int hash = mlr_string_pair_hash_func(key1, key2);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	int num_tries = 0;
	int done = 0;

	while (!done) {
		lhms2ve_t* pe = &pmap->entries[index];
		if (pmap->states[index] == OCCUPIED) {
			char* ekey1 = pe->key1;
			char* ekey2 = pe->key2;
			// Existing key found in chain.
			if (streq(key1, ekey1) && streq(key2, ekey2))
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
				"Coding error:  table full even after enlargement.");
			exit(1);
		}

		// Linear probing.
		if (++index >= pmap->array_length)
			index = 0;
	}
	return -1; // xxx not reached
}

// ----------------------------------------------------------------
void* lhms2v_put(lhms2v_t* pmap, char* key1, char* key2, void* pvvalue) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhms2v_enlarge(pmap);
	return lhms2v_put_no_enlarge(pmap, key1, key2, pvvalue);
}

static void* lhms2v_put_no_enlarge(lhms2v_t* pmap, char* key1, char* key2, void* pvvalue) {
	int index = lhms2v_find_index_for_key(pmap, key1, key2);
	lhms2ve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (streq(pe->key1, key1) && streq(pe->key2, key2)) {
			pe->pvvalue = pvvalue;
			return pvvalue;
		}
	}
	else if (pmap->states[index] == EMPTY) {
		// End of chain.
		pe->ideal_index = mlr_canonical_mod(mlr_string_pair_hash_func(key1, key2), pmap->array_length);
		// xxx comment memmgt.
		pe->key1 = strdup(key1);
		pe->key2 = strdup(key2);
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
		return pvvalue;
	}
	else {
		fprintf(stderr, "lhms2v_find_index_for_key did not find end of chain");
		exit(1);
	}
	return NULL; // xxx not reached
}

// ----------------------------------------------------------------
void* lhms2v_get(lhms2v_t* pmap, char* key1, char* key2) {
	int index = lhms2v_find_index_for_key(pmap, key1, key2);
	lhms2ve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe->pvvalue;
	else if (pmap->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "lhms2v_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhms2v_has_key(lhms2v_t* pmap, char* key1, char* key2) {
	int index = lhms2v_find_index_for_key(pmap, key1, key2);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "lhms2v_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
void* lhms2v_remove(lhms2v_t* pmap, char* key1, char* key2) {
	int index = lhms2v_find_index_for_key(pmap, key1, key2);
	lhms2ve_t* pe = &pmap->entries[index];
	if (pmap->states[index] == OCCUPIED) {
		void* pvvalue   = pe->pvvalue;
		pe->ideal_index = -1;
		pe->key1        = NULL;
		pe->key2        = NULL;
		pe->pvvalue     = NULL;
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
		return pvvalue;
	}
	else if (pmap->states[index] == EMPTY) {
		return NULL;
	}
	else {
		fprintf(stderr, "lhms2v_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
void lhms2v_clear(lhms2v_t* pmap) {
	for (int i = 0; i < pmap->array_length; i++) {
		lhms2ve_clear(&pmap->entries[i]);
		pmap->num_occupied = 0;
		pmap->num_freed = 0;
	}
}

int lhms2v_size(lhms2v_t* pmap) {
	return pmap->num_occupied;
}

// ----------------------------------------------------------------
static void lhms2v_enlarge(lhms2v_t* pmap) {
	lhms2ve_t*       old_entries = pmap->entries;
	lhms2ve_state_t* old_states  = pmap->states;
	lhms2ve_t*       old_head    = pmap->phead;

	lhms2v_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhms2ve_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		lhms2v_put_no_enlarge(pmap, pe->key1, pe->key2, pe->pvvalue);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
void lhms2v_check_counts(lhms2v_t* pmap) {
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
			"occupancy-count mismatch:  actual %d != cached  %d",
				nocc, pmap->num_occupied);
		exit(1);
	}
	if (ndel != pmap->num_freed) {
		fprintf(stderr,
			"deleted-count mismatch:  actual %d != cached  %d",
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

void lhms2v_dump(lhms2v_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhms2ve_t* pe = &pmap->entries[index];

		const char* key1_string = (pe == NULL) ? "none" :
			pe->key1 == NULL ? "null" :
			pe->key1;
		const char* key2_string = (pe == NULL) ? "none" :
			pe->key2 == NULL ? "null" :
			pe->key2;
		const char* value_string = (pe == NULL) ? "none" :
			pe->pvvalue == NULL ? "null" :
			pe->pvvalue;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key1: %12s | key2: %12s | pvvalue: %12s |\n",
			get_state_name(pmap->states[index]), index, pe->ideal_index,
			key1_string, key2_string, value_string);
	}
	printf("+\n");
	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
	printf("+\n");
	for (lhms2ve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key1_string = (pe == NULL) ? "none" :
			pe->key1 == NULL ? "null" :
			pe->key1;
		const char* key2_string = (pe == NULL) ? "none" :
			pe->key2 == NULL ? "null" :
			pe->key2;
		const char* value_string = (pe == NULL) ? "none" :
			pe->pvvalue == NULL ? "null" :
			pe->pvvalue;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key1: %12s | key2: %12s | pvvalue: %12s |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key1_string, key2_string, value_string);
	}
}
