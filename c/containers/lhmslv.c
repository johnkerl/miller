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
#include "containers/lhmslv.h"

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
static void* lhmslv_put_no_enlarge(lhmslv_t* pmap, slls_t* key, void* pvvalue);
static void lhmslv_enlarge(lhmslv_t* pmap);

// ================================================================
static void lhmslve_clear(lhmslve_t *pentry) {
	pentry->ideal_index = -1;
	pentry->key         = NULL;
	pentry->pvvalue     = NULL;
	pentry->pprev       = NULL;
	pentry->pnext       = NULL;
}

static void lhmslv_init(lhmslv_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhmslve_t*)mlr_malloc_or_die(sizeof(lhmslve_t) * length);
	// Don't do lhmslve_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.
	//for (int i = 0; i < length; i++)
	//	lhmslve_clear(&entries[i]);

	pmap->states       = (lhmslve_state_t*)mlr_malloc_or_die(sizeof(lhmslve_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead        = NULL;
	pmap->ptail        = NULL;
}

lhmslv_t* lhmslv_alloc() {
	lhmslv_t* pmap = mlr_malloc_or_die(sizeof(lhmslv_t));
	lhmslv_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

// xxx cmt re memmgt
void lhmslv_free(lhmslv_t* pmap) {
	if (pmap == NULL)
		return;
	// xxx free each key
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
static int lhmslv_find_index_for_key(lhmslv_t* pmap, slls_t* key) {
	int hash = slls_hash_func(key);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	int num_tries = 0;

	while (TRUE) {
		lhmslve_t* pe = &pmap->entries[index];
		if (pmap->states[index] == OCCUPIED) {
			slls_t* ekey = pe->key;
			// Existing key found in chain.
			if (slls_equals(key, ekey))
				return index;
		}
		else if (pmap->states[index] == EMPTY) {
			return index;
		}

		// If the current entry has been freed, i.e. previously occupied,
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
	fprintf(stderr, "Miller: coding error detected in file %s at line %d.\n",
		__FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
void* lhmslv_put(lhmslv_t* pmap, slls_t* key, void* pvvalue) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhmslv_enlarge(pmap);
	return lhmslv_put_no_enlarge(pmap, key, pvvalue);
}

static void* lhmslv_put_no_enlarge(lhmslv_t* pmap, slls_t* key, void* pvvalue) {
	int index = lhmslv_find_index_for_key(pmap, key);
	lhmslve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (slls_equals(pe->key, key)) {
			pe->pvvalue = pvvalue;
			return pvvalue;
		}
	}
	else if (pmap->states[index] == EMPTY) {
		// End of chain.
		pe->ideal_index = mlr_canonical_mod(slls_hash_func(key), pmap->array_length);
		// xxx comment memmgt.
		pe->key = key;
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
		fprintf(stderr, "lhmslv_find_index_for_key did not find end of chain\n");
		exit(1);
	}
	// This one is to appease a compiler warning about control reaching the end
	// of a non-void function
	fprintf(stderr, "Miller: coding error detected in file %s at line %d.\n",
		__FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
void* lhmslv_get(lhmslv_t* pmap, slls_t* key) {
	int index = lhmslv_find_index_for_key(pmap, key);
	lhmslve_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe->pvvalue;
	else if (pmap->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "lhmslv_find_index_for_key did not find end of chain\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhmslv_has_key(lhmslv_t* pmap, slls_t* key) {
	int index = lhmslv_find_index_for_key(pmap, key);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "lhmslv_find_index_for_key did not find end of chain\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void* lhmslv_remove(lhmslv_t* pmap, slls_t* key) {
	int index = lhmslv_find_index_for_key(pmap, key);
	lhmslve_t* pe = &pmap->entries[index];
	if (pmap->states[index] == OCCUPIED) {
		void* pvvalue   = pe->pvvalue;
		pe->ideal_index = -1;
		pe->key         = NULL;
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
		fprintf(stderr, "lhmslv_find_index_for_key did not find end of chain\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void lhmslv_clear(lhmslv_t* pmap) {
	for (int i = 0; i < pmap->array_length; i++) {
		lhmslve_clear(&pmap->entries[i]);
		pmap->states[i] = EMPTY;
	}
	pmap->num_occupied = 0;
	pmap->num_freed = 0;
}

int lhmslv_size(lhmslv_t* pmap) {
	return pmap->num_occupied;
}

// ----------------------------------------------------------------
static void lhmslv_enlarge(lhmslv_t* pmap) {
	lhmslve_t*       old_entries = pmap->entries;
	lhmslve_state_t* old_states  = pmap->states;
	lhmslve_t*       old_head    = pmap->phead;

	lhmslv_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhmslve_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		lhmslv_put_no_enlarge(pmap, pe->key, pe->pvvalue);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
int lhmslv_check_counts(lhmslv_t* pmap) {
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
			"occupancy-count mismatch:  actual %d != cached  %d\n",
				nocc, pmap->num_occupied);
		return FALSE;
	}
	if (ndel != pmap->num_freed) {
		fprintf(stderr,
			"freed-count mismatch:  actual %d != cached  %d\n",
				ndel, pmap->num_freed);
		return FALSE;
	}
	return TRUE;
}

// ----------------------------------------------------------------
static char* get_state_name(int state) {
	switch(state) {
	case OCCUPIED: return "occupied"; break;
	case DELETED:  return "freed";  break;
	case EMPTY:    return "empty";    break;
	default:       return "?????";    break;
	}
}

void lhmslv_print(lhmslv_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhmslve_t* pe = &pmap->entries[index];

		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			slls_join(pe->key, ",");
		const char* value_string = (pe == NULL) ? "none" :
			pe->pvvalue == NULL ? "null" :
			pe->pvvalue;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | pvvalue: %12s |\n",
			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, value_string);
	}
	printf("+\n");
	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
	printf("+\n");
	for (lhmslve_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			slls_join(pe->key, ",");
		const char* value_string = (pe == NULL) ? "none" :
			pe->pvvalue == NULL ? "null" :
			pe->pvvalue;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | pvvalue: %12s |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key_string, value_string);
	}
}
