// ================================================================
// Array-only (open addressing) string-to-string hash set with linear probing
// for collisions.
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

#include "lib/mlrutil.h"
#include "containers/hss.h"

// ----------------------------------------------------------------
#define INITIAL_ARRAY_LENGTH 128
#define LOAD_FACTOR          0.7

#define OCCUPIED 444
#define DELETED  555
#define EMPTY    666

// ================================================================
static void hsse_clear(hsse_t *pentry) {
	pentry->key         = NULL;
	pentry->state       = EMPTY;
	pentry->ideal_index = -1;
}

// ----------------------------------------------------------------
static hsse_t* hss_make_alloc_array(int length) {
	hsse_t* array = (hsse_t*)mlr_malloc_or_die(sizeof(hsse_t) * length);
	for (int i = 0; i < length; i++)
		hsse_clear(&array[i]);
	return array;
}

static void hss_init(hss_t *pset, int length) {
	pset->num_occupied = 0;
	pset->num_freed    = 0;
	pset->array_length = length;
	pset->array        = hss_make_alloc_array(length);
}

hss_t* hss_alloc() {
	hss_t* pset = mlr_malloc_or_die(sizeof(hss_t));
	hss_init(pset, INITIAL_ARRAY_LENGTH);
	return pset;
}

void hss_free(hss_t* pset) {
	if (pset == NULL)
		return;
	free(pset->array);
	pset->array = NULL;
	pset->num_occupied = 0;
	pset->num_freed    = 0;
	pset->array_length = 0;
	free(pset);
}

// ----------------------------------------------------------------
// Used by get() and remove().
// Returns >0 for where the key is *or* should go (end of chain).
static int hss_find_index_for_key(hss_t* pset, char* key) {
	int hash = mlr_string_hash_func(key);
	int index = mlr_canonical_mod(hash, pset->array_length);
	int num_tries = 0;
	int done = 0;

	while (!done) {
		hsse_t* pe = &pset->array[index];
		if (pe->state == OCCUPIED) {
			char* ekey = pe->key;
			// Existing key found in chain.
			if (streq(key, ekey))
				return index;
		}
		else if (pe->state == EMPTY) {
			return index;
		}

		// If the current entry has been freed, i.e. previously occupied,
		// the sought index may be further down the chain.  So we must
		// continue looking.
		if (++num_tries >= pset->array_length) {
			fprintf(stderr,
				"Coding error:  table full even after enlargement.");
			exit(1);
		}

		// Linear probing.
		if (++index >= pset->array_length)
			index = 0;
	}
	return -1; // xxx not reached
}

// ----------------------------------------------------------------
static void hss_enlarge(hss_t* pset);

void hss_add(hss_t* pset, char* key) {
	if ((pset->num_occupied + pset->num_freed) >= (pset->array_length*LOAD_FACTOR))
		hss_enlarge(pset);

	int index = hss_find_index_for_key(pset, key);
	hsse_t* pe = &pset->array[index];

	if (pe->state == OCCUPIED) {
		// Existing key found in chain
		if (streq(pe->key, key)) {
			return;
		}
		// xxx this looks wrong
	}
	else if (pe->state == EMPTY) {
		// End of chain.
		pe->key = key;
		pe->state = OCCUPIED;
		pe->ideal_index = mlr_canonical_mod(mlr_string_hash_func(key), pset->array_length);
		pset->num_occupied++;
	}
	else {
		fprintf(stderr, "hss_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
static void hss_enlarge(hss_t* pset) {
	int old_array_length = pset->array_length;
	hsse_t* old_array = pset->array;

	hss_init(pset, pset->array_length*2);

	for (int index = 0; index < old_array_length; index++) {
		hsse_t e = old_array[index];
		if (e.state == OCCUPIED)
			hss_add(pset, e.key);
	}
}

// ----------------------------------------------------------------
int hss_has(hss_t* pset, char* key) {
	int index = hss_find_index_for_key(pset, key);
	hsse_t* pe = &pset->array[index];

	if (pe->state == OCCUPIED)
		return TRUE;
	else if (pe->state == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "hss_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
void hss_remove(hss_t* pset, char* key) {
	int index = hss_find_index_for_key(pset, key);
	hsse_t* pe = &pset->array[index];
	if (pe->state == OCCUPIED) {
		pe->key          = NULL;
		pe->state        = DELETED;
		pe->ideal_index  = -1;
		pset->num_freed++;
		pset->num_occupied--;
	}
	else if (pe->state == EMPTY) {
	}
	else {
		fprintf(stderr, "hss_find_index_for_key did not find end of chain");
		exit(1);
	}
}

// ----------------------------------------------------------------
void hss_clear(hss_t* pset) {
	for (int i = 0; i < pset->array_length; i++) {
		hsse_clear(&pset->array[i]);
		pset->num_occupied = 0;
		pset->num_freed = 0;
	}
}

int hss_size(hss_t* pset) {
	return pset->num_occupied;
}

// ----------------------------------------------------------------
void hss_check_counts(hss_t* pset) {
	int nocc = 0;
	int ndel = 0;
	for (int index = 0; index < pset->array_length; index++) {
		hsse_t* pe = &pset->array[index];
		if (pe->state == OCCUPIED)
			nocc++;
		else if (pe->state == DELETED)
			ndel++;
	}
	if (nocc != pset->num_occupied) {
		fprintf(stderr,
			"occupancy-count mismatch:  actual %d != cached  %d",
				nocc, pset->num_occupied);
		exit(1);
	}
	if (ndel != pset->num_freed) {
		fprintf(stderr,
			"freed-count mismatch:  actual %d != cached  %d",
				ndel, pset->num_freed);
		exit(1);
	}
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

void hss_dump(hss_t* pset) {
	for (int index = 0; index < pset->array_length; index++) {
		hsse_t* pe = &pset->array[index];

		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;

		printf(
		"| stt: %-8s  | idx: %6d | nidx: %6d | key: %12s |\n",
			get_state_name(pe->state), index, pe->ideal_index, key_string);
	}
}
