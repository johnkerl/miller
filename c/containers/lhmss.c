// ================================================================
// Array-only (open addressing) string-to-string linked hash map with linear
// probing for collisions.
//
// Keys and values are not strduped.
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
#include "containers/lhmss.h"
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
static void lhmss_put_no_enlarge(lhmss_t* pmap, char* key, char* value, char free_flags);
static void lhmss_enlarge(lhmss_t* pmap);

static void lhmss_init(lhmss_t *pmap, int length) {
	pmap->num_occupied = 0;
	pmap->num_freed    = 0;
	pmap->array_length = length;

	pmap->entries      = (lhmsse_t*)mlr_malloc_or_die(sizeof(lhmsse_t) * length);
	// Don't do lhmsse_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.

	pmap->states = (lhmsse_state_t*)mlr_malloc_or_die(sizeof(lhmsse_state_t) * length);
	memset(pmap->states, EMPTY, length);

	pmap->phead = NULL;
	pmap->ptail = NULL;
}

lhmss_t* lhmss_alloc() {
	lhmss_t* pmap = mlr_malloc_or_die(sizeof(lhmss_t));
	lhmss_init(pmap, INITIAL_ARRAY_LENGTH);
	return pmap;
}

lhmss_t* lhmss_copy(lhmss_t* pmap) {
	lhmss_t* pnew = lhmss_alloc();
	for (lhmsse_t* pe = pmap->phead; pe != NULL; pe = pe->pnext)
		lhmss_put(pnew, mlr_strdup_or_die(pe->key), mlr_strdup_or_die(pe->value), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	return pnew;
}

void lhmss_free(lhmss_t* pmap) {
	if (pmap == NULL)
		return;
	for (lhmsse_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		if (pe->free_flags & FREE_ENTRY_KEY)
			free(pe->key);
		if (pe->free_flags & FREE_ENTRY_VALUE)
			free(pe->value);
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
static int lhmss_find_index_for_key(lhmss_t* pmap, char* key, int* pideal_index) {
	int hash = mlr_string_hash_func(key);
	int index = mlr_canonical_mod(hash, pmap->array_length);
	*pideal_index = index;
	int num_tries = 0;

	while (TRUE) {
		lhmsse_t* pe = &pmap->entries[index];
		if (pmap->states[index] == OCCUPIED) {
			char* ekey = pe->key;
			// Existing key found in chain.
			if (streq(key, ekey))
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
void lhmss_put(lhmss_t* pmap, char* key, char* value, char free_flags) {
	if ((pmap->num_occupied + pmap->num_freed) >= (pmap->array_length*LOAD_FACTOR))
		lhmss_enlarge(pmap);
	lhmss_put_no_enlarge(pmap, key, value, free_flags);
}

static void lhmss_put_no_enlarge(lhmss_t* pmap, char* key, char* value, char free_flags) {
	int ideal_index = 0;
	int index = lhmss_find_index_for_key(pmap, key, &ideal_index);
	lhmsse_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		if (pe->free_flags & FREE_ENTRY_KEY)
			free(key);
		if (pe->free_flags & FREE_ENTRY_VALUE)
			free(pe->value);
		pe->value = value;
		if (free_flags & FREE_ENTRY_VALUE)
			pe->free_flags |= FREE_ENTRY_VALUE;
		else
			pe->free_flags &= ~FREE_ENTRY_VALUE;

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
		fprintf(stderr, "%s: lhmss_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
char* lhmss_get(lhmss_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmss_find_index_for_key(pmap, key, &ideal_index);
	lhmsse_t* pe = &pmap->entries[index];

	if (pmap->states[index] == OCCUPIED) {
		return pe->value;
	} else if (pmap->states[index] == EMPTY) {
		return NULL;
	} else {
		fprintf(stderr, "%s: lhmss_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
int lhmss_has_key(lhmss_t* pmap, char* key) {
	int ideal_index = 0;
	int index = lhmss_find_index_for_key(pmap, key, &ideal_index);

	if (pmap->states[index] == OCCUPIED)
		return TRUE;
	else if (pmap->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "%s: lhmss_find_index_for_key did not find end of chain.\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
void  lhmss_rename(lhmss_t* pmap, char* old_key, char* new_key) {
	fprintf(stderr, "rename is not supported in the hashed-record impl.\n");
	exit(1);
}

// ----------------------------------------------------------------
static void lhmss_enlarge(lhmss_t* pmap) {
	lhmsse_t*       old_entries = pmap->entries;
	lhmsse_state_t* old_states  = pmap->states;
	lhmsse_t*       old_head    = pmap->phead;

	lhmss_init(pmap, pmap->array_length*ENLARGEMENT_FACTOR);

	for (lhmsse_t* pe = old_head; pe != NULL; pe = pe->pnext) {
		lhmss_put_no_enlarge(pmap, pe->key, pe->value, pe->free_flags);
	}
	free(old_entries);
	free(old_states);
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

void lhmss_dump(lhmss_t* pmap) {
	for (int index = 0; index < pmap->array_length; index++) {
		lhmsse_t* pe = &pmap->entries[index];

		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		const char* value_string = (pe == NULL) ? "none" :
			pe->value == NULL ? "null" :
			pe->value;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s | value: %12s |\n",
			get_state_name(pmap->states[index]), index, pe->ideal_index, key_string, value_string);
	}
	printf("+\n");
	printf("| phead: %p | ptail %p\n", pmap->phead, pmap->ptail);
	printf("+\n");
	for (lhmsse_t* pe = pmap->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		const char* value_string = (pe == NULL) ? "none" :
			pe->value == NULL ? "null" :
			pe->value;
		printf(
		"| prev: %p curr: %p next: %p | nidx: %6d | key: %12s | value: %12s |\n",
			pe->pprev, pe, pe->pnext,
			pe->ideal_index, key_string, value_string);
	}
}

// ----------------------------------------------------------------
int lhmss_check_counts(lhmss_t* pmap) {
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
