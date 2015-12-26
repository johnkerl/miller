// ================================================================
// Array-only (open addressing) string-to-string linked hash map with linear
// probing for collisions.
//
// Keys are strduped.
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
#include "containers/lhmsi.h"

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
static void lhmsi_put_no_enlarge(lhmsi_t* pmap, char* key, int value);
static void lhmsi_enlarge(lhmsi_t* pmap);

// ================================================================
static void lhmsi_init(lhmsi_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhmsie_t*)mlr_malloc_or_die(sizeof(lhmsie_t) * length);
	// Don't do lhmsie_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.
	//for (int i = 0; i < length; i++)
	//	lhmsie_clear(&entries[i]);

	pmap->states       = (lhmsie_state_t*)mlr_malloc_or_die(sizeof(lhmsie_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead        = NULL;
	pmap->ptail        = NULL;
}

lhmsi_t* lhmsi_alloc() {
	lhmsi_t* pmap = mlr_malloc_or_die(sizeof(lhmsi_t));
	lhmsi_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

lhmsi_t* lhmsi_copy(lhmsi_t* pmap) {
	lhmsi_t* pnew = lhmsi_alloc();
	for (lhmsie_t* pe = pmap->phead; pe != NULL; pe = pe->pnext)
		lhmsi_put(pnew, pe->key, pe->value);
	return pnew;
}

void lhmsi_free(lhmsi_t* pmap) {
	if (pmap == NULL)
		return;
	for (lhmsie_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
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
// Returns >0 for where the key is *or* should go (end of chain).
static int lhmsi_find_index_for_key(lhmsi_t* pmap, char* key) {
	int hash = mlr_string_hash_func(key);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	int num_tries = 0;

	while (TRUE) {
		lhmsie_t* pe = &pmap->entries[index];
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
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.argv0, __FILE__, __LINE__);
	exit(1);
}

// ----------------------------------------------------------------
void lhmsi_put(lhmsi_t* pmap, char* key, int value) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhmsi_enlarge(pmap);
	lhmsi_put_no_enlarge(pmap, key, value);
}

static void lhmsi_put_no_enlarge(lhmsi_t* pmap, char* key, int value) {
	int index = lhmsi_find_index_for_key(pmap, key);
	lhmsie_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (streq(pe->key, key)) {
			pe->value = value;
		}
	}
	else if (pmap->states[index] == EMPTY) {
		// End of chain.
		pe->ideal_index = mlr_canonical_mod(mlr_string_hash_func(key), pmap->array_length);
		pe->key = mlr_strdup_or_die(key);
		pe->value = value;
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
		fprintf(stderr, "lhmsi_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
long long lhmsi_get(lhmsi_t* pmap, char* key) {
	int index = lhmsi_find_index_for_key(pmap, key);
	lhmsie_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe->value;
	else if (pmap->states[index] == EMPTY)
		return -999; // caller must do lhmsi_has_key to check validity
	else {
		fprintf(stderr, "lhmsi_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

lhmsie_t* lhmsi_get_entry(lhmsi_t* pmap, char* key) {
	int index = lhmsi_find_index_for_key(pmap, key);
	lhmsie_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED)
		return pe;
	else if (pmap->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "lhmsi_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhmsi_has_key(lhmsi_t* pmap, char* key) {
	int index = lhmsi_find_index_for_key(pmap, key);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "lhmsi_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void lhmsi_remove(lhmsi_t* pmap, char* key) {
	int index = lhmsi_find_index_for_key(pmap, key);
	lhmsie_t* pe = &pmap->entries[index];
	if (pmap->states[index] == OCCUPIED) {
		pe->ideal_index = -1;
		pe->key         = NULL;
		pe->value       = -999;
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
		fprintf(stderr, "lhmsi_find_index_for_key did not find end of chain.\n");
		exit(1);
	}
}

// ----------------------------------------------------------------
void  lhmsi_rename(lhmsi_t* pmap, char* old_key, char* new_key) {
	fprintf(stderr, "rename is not supported in the hashed-record impl.\n");
	exit(1);
}

// ----------------------------------------------------------------
static void lhmsie_clear(lhmsie_t *pentry) {
	pentry->ideal_index = -1;
	pentry->key         = NULL;
	pentry->value       = -1;
	pentry->pprev       = NULL;
	pentry->pnext       = NULL;
}

void lhmsi_clear(lhmsi_t* pmap) {
	for (int i = 0; i < pmap->array_length; i++) {
		lhmsie_clear(&pmap->entries[i]);
		pmap->states[i] = EMPTY;
	}
	pmap->num_occupied = 0;
	pmap->num_freed = 0;
}

// ----------------------------------------------------------------
static void lhmsi_enlarge(lhmsi_t* pmap) {
	lhmsie_t*       old_entries = pmap->entries;
	lhmsie_state_t* old_states  = pmap->states;
	lhmsie_t*       old_head    = pmap->phead;

	lhmsi_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhmsie_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		// xxx implement free-flags here so we can do this without the redundant strdups
		lhmsi_put_no_enlarge(pmap, pe->key, pe->value);
		free(pe->key);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
int lhmsi_check_counts(lhmsi_t* pmap) {
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

void lhmsi_print(lhmsi_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhmsie_t* pe = &pmap->entries[index];

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
	for (lhmsie_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | value: %8lld |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key_string, pe->value);
	}
}
