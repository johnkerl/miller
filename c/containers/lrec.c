#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"
#include "containers/lrec.h"

#define SB_ALLOC_LENGTH 256

static lrece_t* lrec_find_entry(lrec_t* prec, char* key);
static void lrec_link_at_head(lrec_t* prec, lrece_t* pe);
static void lrec_link_at_tail(lrec_t* prec, lrece_t* pe);

static void lrec_unbacked_free(lrec_t* prec);
static void lrec_free_single_line_backing(lrec_t* prec);
static void lrec_free_csv_backing(lrec_t* prec);
static void lrec_free_multiline_backing(lrec_t* prec);

// ----------------------------------------------------------------
lrec_t* lrec_unbacked_alloc() {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->pfree_backing_func = lrec_unbacked_free;
	return prec;
}

lrec_t* lrec_dkvp_alloc(char* line) {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->psingle_line = line;
	prec->pfree_backing_func = lrec_free_single_line_backing;
	return prec;
}

lrec_t* lrec_nidx_alloc(char* line) {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->psingle_line  = line;
	prec->pfree_backing_func = lrec_free_single_line_backing;
	return prec;
}

lrec_t* lrec_csvlite_alloc(char* data_line) {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->psingle_line = data_line;
	prec->pfree_backing_func = lrec_free_csv_backing;
	return prec;
}

lrec_t* lrec_csv_alloc(char* data_line) {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->psingle_line = data_line;
	prec->pfree_backing_func = lrec_free_csv_backing;
	return prec;
}

lrec_t* lrec_xtab_alloc(slls_t* pxtab_lines) {
	lrec_t* prec = mlr_malloc_or_die(sizeof(lrec_t));
	memset(prec, 0, sizeof(lrec_t));
	prec->pxtab_lines = pxtab_lines;
	prec->pfree_backing_func = lrec_free_multiline_backing;
	return prec;
}

// ----------------------------------------------------------------
static void lrec_free_contents(lrec_t* prec) {
	for (lrece_t* pe = prec->phead; pe != NULL; /*pe = pe->pnext*/) {
		if (pe->free_flags & FREE_ENTRY_KEY)
			free(pe->key);
		if (pe->free_flags & FREE_ENTRY_VALUE)
			free(pe->value);
		lrece_t* ope = pe;
		pe = pe->pnext;
		free(ope);
	}
	prec->pfree_backing_func(prec);
}

// ----------------------------------------------------------------
void lrec_clear(lrec_t* prec) {
	if (prec == NULL)
		return;
	lrec_free_contents(prec);
	memset(prec, 0, sizeof(lrec_t));
	prec->pfree_backing_func = lrec_unbacked_free;
}

// ----------------------------------------------------------------
void lrec_free(lrec_t* prec) {
	if (prec == NULL)
		return;
	lrec_free_contents(prec);
	free(prec);
}

// ----------------------------------------------------------------
lrec_t* lrec_copy(lrec_t* pinrec) {
	lrec_t* poutrec = lrec_unbacked_alloc();
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		lrec_put(poutrec, mlr_strdup_or_die(pe->key), mlr_strdup_or_die(pe->value),
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}
	return poutrec;
}

// ----------------------------------------------------------------
void lrec_put(lrec_t* prec, char* key, char* value, char free_flags) {
	lrece_t* pe = lrec_find_entry(prec, key);

	if (pe != NULL) {
		if (pe->free_flags & FREE_ENTRY_VALUE) {
			free(pe->value);
		}
		if (free_flags & FREE_ENTRY_KEY)
			free(key);
		pe->value = value;
		if (free_flags & FREE_ENTRY_VALUE)
			pe->free_flags |= FREE_ENTRY_VALUE;
		else
			pe->free_flags &= ~FREE_ENTRY_VALUE;
	} else {
		pe = mlr_malloc_or_die(sizeof(lrece_t));
		pe->key         = key;
		pe->value       = value;
		pe->free_flags  = free_flags;
		pe->quote_flags = 0;

		if (prec->phead == NULL) {
			pe->pprev   = NULL;
			pe->pnext   = NULL;
			prec->phead = pe;
			prec->ptail = pe;
		} else {
			pe->pprev   = prec->ptail;
			pe->pnext   = NULL;
			prec->ptail->pnext = pe;
			prec->ptail = pe;
		}
		prec->field_count++;
	}
}

void lrec_put_ext(lrec_t* prec, char* key, char* value, char free_flags, char quote_flags) {
	lrece_t* pe = lrec_find_entry(prec, key);

	if (pe != NULL) {
		if (pe->free_flags & FREE_ENTRY_VALUE) {
			free(pe->value);
		}
		if (free_flags & FREE_ENTRY_KEY)
			free(key);
		pe->value = value;
		if (free_flags & FREE_ENTRY_VALUE)
			pe->free_flags |= FREE_ENTRY_VALUE;
		else
			pe->free_flags &= ~FREE_ENTRY_VALUE;
	} else {
		pe = mlr_malloc_or_die(sizeof(lrece_t));
		pe->key         = key;
		pe->value       = value;
		pe->free_flags  = free_flags;
		pe->quote_flags = quote_flags;

		if (prec->phead == NULL) {
			pe->pprev   = NULL;
			pe->pnext   = NULL;
			prec->phead = pe;
			prec->ptail = pe;
		} else {
			pe->pprev   = prec->ptail;
			pe->pnext   = NULL;
			prec->ptail->pnext = pe;
			prec->ptail = pe;
		}
		prec->field_count++;
	}
}

void lrec_prepend(lrec_t* prec, char* key, char* value, char free_flags) {
	lrece_t* pe = lrec_find_entry(prec, key);

	if (pe != NULL) {
		if (pe->free_flags & FREE_ENTRY_VALUE) {
			free(pe->value);
		}
		pe->value = value;
		pe->free_flags &= ~FREE_ENTRY_VALUE;
		if (free_flags & FREE_ENTRY_VALUE)
			pe->free_flags |= FREE_ENTRY_VALUE;
	} else {
		pe = mlr_malloc_or_die(sizeof(lrece_t));
		pe->key         = key;
		pe->value       = value;
		pe->free_flags  = free_flags;
		pe->quote_flags = 0;

		if (prec->phead == NULL) {
			pe->pprev   = NULL;
			pe->pnext   = NULL;
			prec->phead = pe;
			prec->ptail = pe;
		} else {
			pe->pnext   = prec->phead;
			pe->pprev   = NULL;
			prec->phead->pprev = pe;
			prec->phead = pe;
		}
		prec->field_count++;
	}
}

lrece_t* lrec_put_after(lrec_t* prec, lrece_t* pd, char* key, char* value, char free_flags) {
	lrece_t* pe = lrec_find_entry(prec, key);

	if (pe != NULL) { // Overwrite
		if (pe->free_flags & FREE_ENTRY_VALUE) {
			free(pe->value);
		}
		pe->value = value;
		pe->free_flags &= ~FREE_ENTRY_VALUE;
		if (free_flags & FREE_ENTRY_VALUE)
			pe->free_flags |= FREE_ENTRY_VALUE;
	} else { // Insert after specified entry
		pe = mlr_malloc_or_die(sizeof(lrece_t));
		pe->key         = key;
		pe->value       = value;
		pe->free_flags  = free_flags;
		pe->quote_flags = 0;

		if (pd->pnext == NULL) { // Append at end of list
			pd->pnext = pe;
			pe->pprev = pd;
			pe->pnext = NULL;
			prec->ptail = pe;

		} else {
			lrece_t* pf = pd->pnext;
			pd->pnext = pe;
			pf->pprev = pe;
			pe->pprev = pd;
			pe->pnext = pf;
		}

		prec->field_count++;
	}
	return pe;
}

// ----------------------------------------------------------------
char* lrec_get(lrec_t* prec, char* key) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe != NULL) {
		return pe->value;
	} else {
		return NULL;
	}
}

char* lrec_get_pff(lrec_t* prec, char* key, char** ppfree_flags) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe != NULL) {
		*ppfree_flags = &pe->free_flags;
		return pe->value;
	} else {
		*ppfree_flags = NULL;
		return NULL;
	}
}

char* lrec_get_ext(lrec_t* prec, char* key, lrece_t** ppentry) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe != NULL) {
		*ppentry = pe;
		return pe->value;
	} else {
		*ppentry = NULL;;
		return NULL;
	}
}

// ----------------------------------------------------------------
lrece_t* lrec_get_pair_by_position(lrec_t* prec, int position) { // 1-up not 0-up
	if (position <= 0 || position > prec->field_count) {
		return NULL;
	}
	int sought_index = position - 1;
	int found_index = 0;
	lrece_t* pe = NULL;
	for (
		found_index = 0, pe = prec->phead;
		pe != NULL;
		found_index++, pe = pe->pnext
	) {
		if (found_index == sought_index) {
			return pe;
		}
	}
	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
		MLR_GLOBALS.bargv0, __FILE__, __LINE__);
	exit(1);
}

char* lrec_get_key_by_position(lrec_t* prec, int position) { // 1-up not 0-up
	lrece_t* pe = lrec_get_pair_by_position(prec, position);
	if (pe == NULL) {
		return NULL;
	} else {
		return pe->key;
	}
}

char* lrec_get_value_by_position(lrec_t* prec, int position) { // 1-up not 0-up
	lrece_t* pe = lrec_get_pair_by_position(prec, position);
	if (pe == NULL) {
		return NULL;
	} else {
		return pe->value;
	}
}

// ----------------------------------------------------------------
void lrec_remove(lrec_t* prec, char* key) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe == NULL)
		return;

	lrec_unlink(prec, pe);

	if (pe->free_flags & FREE_ENTRY_KEY) {
		free(pe->key);
	}
	if (pe->free_flags & FREE_ENTRY_VALUE) {
		free(pe->value);
	}

	free(pe);
}

// ----------------------------------------------------------------
void lrec_remove_by_position(lrec_t* prec, int position) { // 1-up not 0-up
	lrece_t* pe = lrec_get_pair_by_position(prec, position);
	if (pe == NULL)
		return;

	lrec_unlink(prec, pe);

	if (pe->free_flags & FREE_ENTRY_KEY) {
		free(pe->key);
	}
	if (pe->free_flags & FREE_ENTRY_VALUE) {
		free(pe->value);
	}

	free(pe);
}

// Before:
//   "x" => "3"
//   "y" => "4"  <-- pold
//   "z" => "5"  <-- pnew
//
// Rename y to z
//
// After:
//   "x" => "3"
//   "z" => "4"
//
void lrec_rename(lrec_t* prec, char* old_key, char* new_key, int new_needs_freeing) {

	lrece_t* pold = lrec_find_entry(prec, old_key);
	if (pold != NULL) {
		lrece_t* pnew = lrec_find_entry(prec, new_key);

		if (pnew == NULL) { // E.g. rename "x" to "y" when "y" is not present
			if (pold->free_flags & FREE_ENTRY_KEY) {
				free(pold->key);
				pold->key = new_key;
				if (!new_needs_freeing)
					pold->free_flags &= ~FREE_ENTRY_KEY;
			} else {
				pold->key = new_key;
				if (new_needs_freeing)
					pold->free_flags |=  FREE_ENTRY_KEY;
			}

		} else { // E.g. rename "x" to "y" when "y" is already present
			if (pnew->free_flags & FREE_ENTRY_VALUE) {
				free(pnew->value);
			}
			if (pold->free_flags & FREE_ENTRY_KEY) {
				free(pold->key);
				pold->free_flags &= ~FREE_ENTRY_KEY;
			}
			pold->key = new_key;
			if (new_needs_freeing)
				pold->free_flags |=  FREE_ENTRY_KEY;
			else
				pold->free_flags &= ~FREE_ENTRY_KEY;
			lrec_unlink(prec, pnew);
			free(pnew);
		}
	}
}

// Cases:
// 1. Rename field at position 3 from "x" to "y when "y" does not exist elsewhere in the srec
// 2. Rename field at position 3 from "x" to "y when "y" does     exist elsewhere in the srec
// Note: position is 1-up not 0-up
void  lrec_rename_at_position(lrec_t* prec, int position, char* new_key, int new_needs_freeing){
	lrece_t* pe = lrec_get_pair_by_position(prec, position);
	if (pe == NULL) {
		if (new_needs_freeing) {
			free(new_key);
		}
		return;
	}

	lrece_t* pother = lrec_find_entry(prec, new_key);

	if (pe->free_flags & FREE_ENTRY_KEY) {
		free(pe->key);
	}
	pe->key = new_key;
	if (new_needs_freeing) {
		pe->free_flags |= FREE_ENTRY_KEY;
	} else {
		pe->free_flags &= ~FREE_ENTRY_KEY;
	}
	if (pother != NULL) {
		lrec_unlink(prec, pother);
		free(pother);
	}
}

// ----------------------------------------------------------------
void lrec_move_to_head(lrec_t* prec, char* key) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe == NULL)
		return;

	lrec_unlink(prec, pe);
	lrec_link_at_head(prec, pe);
}

void lrec_move_to_tail(lrec_t* prec, char* key) {
	lrece_t* pe = lrec_find_entry(prec, key);
	if (pe == NULL)
		return;

	lrec_unlink(prec, pe);
	lrec_link_at_tail(prec, pe);
}

// ----------------------------------------------------------------
void lrec_unlink(lrec_t* prec, lrece_t* pe) {
	if (pe == prec->phead) {
		if (pe == prec->ptail) {
			prec->phead = NULL;
			prec->ptail = NULL;
		} else {
			prec->phead = pe->pnext;
			pe->pnext->pprev = NULL;
		}
	} else {
		pe->pprev->pnext = pe->pnext;
		if (pe == prec->ptail) {
			prec->ptail = pe->pprev;
		} else {
			pe->pnext->pprev = pe->pprev;
		}
	}
	prec->field_count--;
}

void lrec_unlink_and_free(lrec_t* prec, lrece_t* pe) {
	if (pe->free_flags & FREE_ENTRY_KEY)
		free(pe->key);
	if (pe->free_flags & FREE_ENTRY_VALUE)
		free(pe->value);
	lrec_unlink(prec, pe);
	free(pe);
}

// ----------------------------------------------------------------
static void lrec_link_at_head(lrec_t* prec, lrece_t* pe) {

	if (prec->phead == NULL) {
		pe->pprev   = NULL;
		pe->pnext   = NULL;
		prec->phead = pe;
		prec->ptail = pe;
	} else {
		// [b,c,d] + a
		pe->pprev   = NULL;
		pe->pnext   = prec->phead;
		prec->phead->pprev = pe;
		prec->phead = pe;
	}
	prec->field_count++;
}

static void lrec_link_at_tail(lrec_t* prec, lrece_t* pe) {

	if (prec->phead == NULL) {
		pe->pprev   = NULL;
		pe->pnext   = NULL;
		prec->phead = pe;
		prec->ptail = pe;
	} else {
		pe->pprev   = prec->ptail;
		pe->pnext   = NULL;
		prec->ptail->pnext = pe;
		prec->ptail = pe;
	}
	prec->field_count++;
}

// ----------------------------------------------------------------
void lrec_dump(lrec_t* prec) {
	lrec_dump_fp(prec, stdout);
}
void lrec_dump_fp(lrec_t* prec, FILE* fp) {
	if (prec == NULL) {
		fprintf(fp, "NULL\n");
		return;
	}
	fprintf(fp, "field_count = %d\n", prec->field_count);
	fprintf(fp, "| phead: %16p | ptail %16p\n", prec->phead, prec->ptail);
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		const char* key_string = (pe == NULL) ? "none" :
			pe->key == NULL ? "null" :
			pe->key;
		const char* value_string = (pe == NULL) ? "none" :
			pe->value == NULL ? "null" :
			pe->value;
		fprintf(fp,
		"| prev: %16p curr: %16p next: %16p | key: %12s | value: %12s |\n",
			pe->pprev, pe, pe->pnext,
			key_string, value_string);
	}
}

void lrec_dump_titled(char* msg, lrec_t* prec) {
	printf("%s:\n", msg);
	lrec_dump(prec);
	printf("\n");
}

void lrec_pointer_dump(lrec_t* prec) {
	printf("prec %p\n", prec);
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		printf("  pe %p k %p v %p\n", pe, pe->key, pe->value);
	}
}

// ----------------------------------------------------------------
static void lrec_unbacked_free(lrec_t* prec) {
}

static void lrec_free_single_line_backing(lrec_t* prec) {
	free(prec->psingle_line);
}

static void lrec_free_csv_backing(lrec_t* prec) {
	free(prec->psingle_line);
}

static void lrec_free_multiline_backing(lrec_t* prec) {
	slls_free(prec->pxtab_lines);
}

// ================================================================

// ----------------------------------------------------------------
// Note on efficiency:
//
// I was imagining/hoping that strcmp has additional optimizations (e.g.
// hand-coded in assembly), so I don't *want* to re-implement it (i.e. I
// probably can't outperform it).
//
// But actual experiments show I get about a 1-2% performance gain doing it
// myself (on my particular system).

static lrece_t* lrec_find_entry(lrec_t* prec, char* key) {
#if 1
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		char* pa = pe->key;
		char* pb = key;
		while (*pa && *pb && (*pa == *pb)) {
			pa++;
			pb++;
		}
		if (*pa == 0 && *pb == 0)
			return pe;
	}
	return NULL;
#else
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext)
		if (streq(pe->key, key))
			return pe;
	return NULL;
#endif
}

// ----------------------------------------------------------------
lrec_t* lrec_literal_1(char* k1, char* v1) {
	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put(prec, k1, v1, NO_FREE);
	return prec;
}

lrec_t* lrec_literal_2(char* k1, char* v1, char* k2, char* v2) {
	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put(prec, k1, v1, NO_FREE);
	lrec_put(prec, k2, v2, NO_FREE);
	return prec;
}

lrec_t* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3) {
	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put(prec, k1, v1, NO_FREE);
	lrec_put(prec, k2, v2, NO_FREE);
	lrec_put(prec, k3, v3, NO_FREE);
	return prec;
}

lrec_t* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4) {
	lrec_t* prec = lrec_unbacked_alloc();
	lrec_put(prec, k1, v1, NO_FREE);
	lrec_put(prec, k2, v2, NO_FREE);
	lrec_put(prec, k3, v3, NO_FREE);
	lrec_put(prec, k4, v4, NO_FREE);
	return prec;
}

void lrec_print(lrec_t* prec) {
	FILE* output_stream = stdout;
	char ors = '\n';
	char ofs = ',';
	char ops = '=';
	if (prec == NULL) {
		fputs("NULL", output_stream);
		fputc(ors, output_stream);
		return;
	}
	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(ofs, output_stream);
		fputs(pe->key, output_stream);
		fputc(ops, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputc(ors, output_stream);
}

char* lrec_sprint(lrec_t* prec, char* ors, char* ofs, char* ops) {
	string_builder_t* psb = sb_alloc(SB_ALLOC_LENGTH);
	if (prec == NULL) {
		sb_append_string(psb, "NULL");
	} else {
		int nf = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			if (nf > 0)
				sb_append_string(psb, ofs);
			sb_append_string(psb, pe->key);
			sb_append_string(psb, ops);
			sb_append_string(psb, pe->value);
			nf++;
		}
		sb_append_string(psb, ors);
	}
	char* rv = sb_finish(psb);
	sb_free(psb);
	return rv;
}
