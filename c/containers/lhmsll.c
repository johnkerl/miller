// ================================================================
// Array-only (open addressing) string-to-string linked hash map with linear
// probing for collisions.
//
// Keys are not strduped.
//
// John Kerl 2012-08-13
//
// Notes:
// * null key is not supported.
// * null value is supported.
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
#include "containers/lhmsll.h"
#include "lib/free_flags.h"

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
static void lhmsll_put_no_enlarge(lhmsll_t* pmap, char* key, int value, char free_flags);
static void lhmsll_enlarge(lhmsll_t* pmap);

// ================================================================
static void lhmsll_init(lhmsll_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhmslle_t*)mlr_malloc_or_die(sizeof(lhmslle_t) * length);
	// Don't do lhmslle_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.

	pmap->states       = (lhmslle_state_t*)mlr_malloc_or_die(sizeof(lhmslle_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead        = NULL;
	pmap->ptail        = NULL;
}

lhmsll_t* lhmsll_alloc() {
	lhmsll_t* pmap = mlr_malloc_or_die(sizeof(lhmsll_t));
	lhmsll_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

lhmsll_t* lhmsll_copy(lhmsll_t* pmap) {
	lhmsll_t* pnew = lhmsll_alloc();
	for (lhmslle_t* pe = pmap->phead; pe != NULL; pe = pe->pnext)
		lhmsll_put(pnew, mlr_strdup_or_die(pe->key), pe->value, FREE_ENTRY_KEY);
	return pnew;
}

void lhmsll_free(lhmsll_t* pmap) {
	if (pmap == NULL)
		return;
	for (lhmslle_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		if (pe->free_flags & FREE_ENTRY_KEY)
			free(pe->key);
	}
	free(pmap->entries);
	free(pmap->states);
	pmap->entries      = NULL;
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = 0;
	free(pmap);
}

// ----------------------------------------------------------------
// Used by get() and remove().
// Returns >=0 for where the key is *or* should go (end of chain).
static int lhmsll_find_index_for_key(lhmsll_t* pmap, char* key, int* pideal_index) {
	int hash = mlr_string_hash_func(key);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	*pideal_index = index;
	int num_tries = 0;

	while (TRUE) {
		lhmslle_t* pe = &pmap->entries[index];
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
				"%s: internal coding error: table full even after enlargement.\n", MLR_GLOBALS.bargv0);
			exit(1);
		}

		// Linear probing.
		if (++index >= pmap->array_length)
			index = 0;
	}
	MLR_INTERNAL_CODING_ERROR();
	return -1; // not reached
}

// ----------------------------------------------------------------
void lhmsll_put(lhmsll_t* pmap, char* key, int value, char free_flags) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhmsll_enlarge(pmap);
	lhmsll_put_no_enlarge(pmap, key, value, free_flags);
}

static void lhmsll_put_no_enlarge(lhmsll_t* pmap, char* key, int value, char free_flags) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);
	lhmslle_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		pe->value = value;

	} else if (pmap->states[index] == EMPTY) {
		// End of chain.
		pe->ideal_index = ideal_index;
		pe->key = key;
		pe->value = value;
		pe->free_flags = free_flags;
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

	} else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
long long lhmsll_get(lhmsll_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);
	lhmslle_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe->value;
	else if (pmap->states[index] == EMPTY)
		return -999; // caller must do lhmsll_has_key to check validity
	else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhmsll_test_and_get(lhmsll_t* pmap, char* key, long long* pval) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);
	lhmslle_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		*pval = pe->value;
		return TRUE;
	} else if (pmap->states[index] == EMPTY) {
		return FALSE;
	} else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

int lhmsll_test_and_increment(lhmsll_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);
	lhmslle_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		pe->value++;
		return TRUE;
	} else if (pmap->states[index] == EMPTY) {
		return FALSE;
	} else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

lhmslle_t* lhmsll_get_entry(lhmsll_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);
	lhmslle_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe;
	else if (pmap->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhmsll_has_key(lhmsll_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmsll_find_index_for_key(pmap, key, &ideal_index);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "%s: lhmsll_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
void  lhmsll_rename(lhmsll_t* pmap, char* old_key, char* new_key) {
	fprintf(stderr, "rename is not supported in the hashed-record impl.\n");
	exit(1);
}

// ----------------------------------------------------------------
static void lhmsll_enlarge(lhmsll_t* pmap) {
	lhmslle_t*       old_entries = pmap->entries;
	lhmslle_state_t* old_states  = pmap->states;
	lhmslle_t*       old_head    = pmap->phead;

	lhmsll_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhmslle_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		lhmsll_put_no_enlarge(pmap, pe->key, pe->value, pe->free_flags);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
int lhmsll_check_counts(lhmsll_t* pmap) {
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
		return FALSE;
	}
	if (ndel != pmap->num_freed) {
		fprintf(stderr,
			"deleted-count mismatch:  actual %d != cached  %d.\n",
				ndel, pmap->num_freed);
		return FALSE;
	}
	return TRUE;
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

void lhmsll_print(lhmsll_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhmslle_t* pe = &pmap->entries[index];

		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | value: %8lld |\n",
			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, pe->value);
	}
	printf("+\n");
	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
	printf("+\n");
	for (lhmslle_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | value: %8lld |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key_string, pe->value);
	}
}
