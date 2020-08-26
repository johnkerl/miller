// ================================================================
// This is a hashless implementation of insertion-ordered key-value pairs for
// Miller's fundamental record data structure.
//
// Design:
//
// * It keeps a doubly-linked list of key-value pairs.
//
// * No hash functions are computed when the map is written to or read from.
//
// * Gets are implemented by sequential scan through the list: given a key,
//   the key-value pairs are scanned through until a match is (or is not) found.
//
// * Performance improvement of 10-15% percent over lhmss was found in the C
//   impleentation (for test data).
//
// Motivation:
//
// * The use case for records in Miller is that *all* fields are read from
//   strings & written to strings (split/join), while only *some* fields are
//   operated on.
//
// * Meanwhile there are few repeated accesses to a given record: the
//   access-to-construct ratio is quite low for Miller data records.  Miller
//   instantiates thousands, millions, billions of records (depending on the
//   input data) but accesses each record only once per mapping operation.
//   (This is in contrast to accumulator hashmaps which are repeatedly accessed
//   during a stats run.)
//
// * The hashed impl computes hashsums for *all* fields whether operated on or not,
//   for the benefit of the *few* fields looked up during the mapping operation.
//
// * The hashless impl only keeps string pointers.  Lookups are done at runtime
//   doing prefix search on the key names. Assuming field names are distinct,
//   this is just a few char-ptr accesses which (in experiments) turn out to
//   offer about a 10-15% performance improvement.
//
// * Added benefit: the field-rename operation (preserving field order) becomes
//   trivial.
//
// Notes:
// * nil key is not supported.
// * nil value is not supported.
// ================================================================

package containers

import (
	"bytes"
	"fmt"
)

// ----------------------------------------------------------------
type Lrec struct {
	fieldCount int
	phead *lrecEntry
	ptail *lrecEntry
}

type lrecEntry struct {
	key *string
	value *string
	pprev *lrecEntry
	pnext *lrecEntry
}

// ----------------------------------------------------------------
func LrecAlloc() *Lrec {
	return &Lrec {
		0,
		nil,
		nil,
	}
}

// ----------------------------------------------------------------
// xxx to do: take an ostream arg

// 5x faster than fmt.Print() separately
func (this *Lrec) Print() {
	var buffer bytes.Buffer
	for pe := this.phead; pe != nil; pe = pe.pnext {
		buffer.WriteString(*pe.key)
		buffer.WriteString("=")
		buffer.WriteString(*pe.value)
		if pe.pnext != nil {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("\n")
	fmt.Print(buffer.String())
}

// ----------------------------------------------------------------
func lrecEntryAlloc(key *string, value *string) *lrecEntry {
	return &lrecEntry {
		key,
		value,
		nil,
		nil,
	}
}

// ----------------------------------------------------------------
func (this *Lrec) findEntry(key *string) *lrecEntry {
	for pe := this.phead; pe != nil; pe = pe.pnext {
		if *pe.key == *key {
			return pe
		}
	}
	return nil
}

// ----------------------------------------------------------------
func (this *Lrec) Put(key *string, value *string) {
	pe := this.findEntry(key)
	if pe == nil {
		pe = lrecEntryAlloc(key, value)
		if this.phead == nil {
			this.phead = pe
			this.ptail = pe
		} else {
			pe.pprev = this.ptail
			pe.pnext = nil
			this.ptail.pnext = pe
			this.ptail = pe
		}
		this.fieldCount++
	} else {
		pe.value = value
	}
}

// ----------------------------------------------------------------
func (this *Lrec) PutAtEnd(key *string, value *string) {
	pe := lrecEntryAlloc(key, value)
	if this.phead == nil {
		this.phead = pe
		this.ptail = pe
	} else {
		pe.pprev = this.ptail
		this.ptail.pnext = pe
		this.ptail = pe
	}
	this.fieldCount++
}

// ----------------------------------------------------------------
func (this *Lrec) Get(key *string) *string {
	pe := this.findEntry(key)
	if pe == nil {
		return nil
	} else {
		return pe.value
	}
	return nil
}

func (this *Lrec) Clear() {
	this.fieldCount = 0
	// Assuming everything unreferenced is getting GC'ed by the Go runtime
	this.phead = nil
	this.ptail = nil
}

func (this *Lrec) Copy() *Lrec{
	that := LrecAlloc()
	for pe := this.phead; pe != nil; pe = pe.pnext {
		that.Put(pe.key, pe.value)
	}
	return that
}

//void lrec_prepend(Lrec* prec, char* key, char* value, char free_flags) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//
//	if (pe != NULL) {
//		if (pe->free_flags & FREE_ENTRY_VALUE) {
//			free(pe->value);
//		}
//		pe->value = value;
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//		if (free_flags & FREE_ENTRY_VALUE)
//			pe->free_flags |= FREE_ENTRY_VALUE;
//	} else {
//		pe = mlr_malloc_or_die(sizeof(lrecEntry));
//		pe->key         = key;
//		pe->value       = value;
//		pe->free_flags  = free_flags;
//		pe->quote_flags = 0;
//
//		if (prec->phead == NULL) {
//			pe->pprev   = NULL;
//			pe->pnext   = NULL;
//			prec->phead = pe;
//			prec->ptail = pe;
//		} else {
//			pe->pnext   = prec->phead;
//			pe->pprev   = NULL;
//			prec->phead->pprev = pe;
//			prec->phead = pe;
//		}
//		prec->field_count++;
//	}
//}

//lrecEntry* lrec_put_after(Lrec* prec, lrecEntry* pd, char* key, char* value, char free_flags) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//
//	if (pe != NULL) { // Overwrite
//		if (pe->free_flags & FREE_ENTRY_VALUE) {
//			free(pe->value);
//		}
//		pe->value = value;
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//		if (free_flags & FREE_ENTRY_VALUE)
//			pe->free_flags |= FREE_ENTRY_VALUE;
//	} else { // Insert after specified entry
//		pe = mlr_malloc_or_die(sizeof(lrecEntry));
//		pe->key         = key;
//		pe->value       = value;
//		pe->free_flags  = free_flags;
//		pe->quote_flags = 0;
//
//		if (pd->pnext == NULL) { // Append at end of list
//			pd->pnext = pe;
//			pe->pprev = pd;
//			pe->pnext = NULL;
//			prec->ptail = pe;
//
//		} else {
//			lrecEntry* pf = pd->pnext;
//			pd->pnext = pe;
//			pf->pprev = pe;
//			pe->pprev = pd;
//			pe->pnext = pf;
//		}
//
//		prec->field_count++;
//	}
//	return pe;
//}

//char* lrec_get_ext(Lrec* prec, char* key, lrecEntry** ppentry) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe != NULL) {
//		*ppentry = pe;
//		return pe->value;
//	} else {
//		*ppentry = NULL;;
//		return NULL;
//	}
//}

//// ----------------------------------------------------------------
//lrecEntry* lrec_get_pair_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	if (position <= 0 || position > prec->field_count) {
//		return NULL;
//	}
//	int sought_index = position - 1;
//	int found_index = 0;
//	lrecEntry* pe = NULL;
//	for (
//		found_index = 0, pe = prec->phead;
//		pe != NULL;
//		found_index++, pe = pe->pnext
//	) {
//		if (found_index == sought_index) {
//			return pe;
//		}
//	}
//	fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
//		MLR_GLOBALS.bargv0, __FILE__, __LINE__);
//	exit(1);
//}

//char* lrec_get_key_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->key;
//	}
//}

//char* lrec_get_value_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		return NULL;
//	} else {
//		return pe->value;
//	}
//}

//// ----------------------------------------------------------------
//void lrec_remove(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//
//	free(pe);
//}

//// ----------------------------------------------------------------
//void lrec_remove_by_position(Lrec* prec, int position) { // 1-up not 0-up
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//
//	free(pe);
//}

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
//void lrec_rename(Lrec* prec, char* old_key, char* new_key, int new_needs_freeing) {
//
//	lrecEntry* pold = lrec_find_entry(prec, old_key);
//	if (pold != NULL) {
//		lrecEntry* pnew = lrec_find_entry(prec, new_key);
//
//		if (pnew == NULL) { // E.g. rename "x" to "y" when "y" is not present
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->key = new_key;
//				if (!new_needs_freeing)
//					pold->free_flags &= ~FREE_ENTRY_KEY;
//			} else {
//				pold->key = new_key;
//				if (new_needs_freeing)
//					pold->free_flags |=  FREE_ENTRY_KEY;
//			}
//
//		} else { // E.g. rename "x" to "y" when "y" is already present
//			if (pnew->free_flags & FREE_ENTRY_VALUE) {
//				free(pnew->value);
//			}
//			if (pold->free_flags & FREE_ENTRY_KEY) {
//				free(pold->key);
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			}
//			pold->key = new_key;
//			if (new_needs_freeing)
//				pold->free_flags |=  FREE_ENTRY_KEY;
//			else
//				pold->free_flags &= ~FREE_ENTRY_KEY;
//			lrec_unlink(prec, pnew);
//			free(pnew);
//		}
//	}
//}

// Cases:
// 1. Rename field at position 3 from "x" to "y when "y" does not exist elsewhere in the srec
// 2. Rename field at position 3 from "x" to "y when "y" does     exist elsewhere in the srec
// Note: position is 1-up not 0-up
//void  lrec_rename_at_position(Lrec* prec, int position, char* new_key, int new_needs_freeing){
//	lrecEntry* pe = lrec_get_pair_by_position(prec, position);
//	if (pe == NULL) {
//		if (new_needs_freeing) {
//			free(new_key);
//		}
//		return;
//	}
//
//	lrecEntry* pother = lrec_find_entry(prec, new_key);
//
//	if (pe->free_flags & FREE_ENTRY_KEY) {
//		free(pe->key);
//	}
//	pe->key = new_key;
//	if (new_needs_freeing) {
//		pe->free_flags |= FREE_ENTRY_KEY;
//	} else {
//		pe->free_flags &= ~FREE_ENTRY_KEY;
//	}
//	if (pother != NULL) {
//		lrec_unlink(prec, pother);
//		free(pother);
//	}
//}

//// ----------------------------------------------------------------
//void lrec_move_to_head(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_head(prec, pe);
//}

//void lrec_move_to_tail(Lrec* prec, char* key) {
//	lrecEntry* pe = lrec_find_entry(prec, key);
//	if (pe == NULL)
//		return;
//
//	lrec_unlink(prec, pe);
//	lrec_link_at_tail(prec, pe);
//}

// ----------------------------------------------------------------
// Simply rename the first (at most) n positions where n is the length of pnames.
//
// Possible complications:
//
// * pnames itself contains duplicates -- we require this as invariant-check
//   from the caller since (for performance) we don't want to check this on every
//   record processed.
//
// * pnames has length less than the current record and one of the new names
//   becomes a clash with an existing name.
//
//   Example:
//   - Input record has names "a,b,c,d,e".
//   - pnames is "d,x,f"
//   - We then construct the invalid "d,x,f,d,e" -- we need to detect and unset
//     the second 'd' field.

//void  lrec_label(Lrec* prec, slls_t* pnames_as_list, hss_t* pnames_as_set) {
//	lrecEntry* pe = prec->phead;
//	sllse_t* pn = pnames_as_list->phead;
//
//	// Process the labels list
//	for ( ; pe != NULL && pn != NULL; pe = pe->pnext, pn = pn->pnext) {
//		char* new_name = pn->value;
//
//		if (pe->free_flags & FREE_ENTRY_KEY) {
//			free(pe->key);
//		}
//		pe->key = mlr_strdup_or_die(new_name);;
//		pe->free_flags |= FREE_ENTRY_KEY;
//	}
//
//	// Process the remaining fields in the record beyond those affected by the new-labels list
//	for ( ; pe != NULL; ) {
//		char* name = pe->key;
//		if (hss_has(pnames_as_set, name)) {
//			lrecEntry* pnext = pe->pnext;
//			if (pe->free_flags & FREE_ENTRY_KEY) {
//				free(pe->key);
//			}
//			if (pe->free_flags & FREE_ENTRY_VALUE) {
//				free(pe->value);
//			}
//			lrec_unlink(prec, pe);
//			free(pe);
//			pe = pnext;
//		} else {
//			pe = pe->pnext;
//		}
//	}
//}

//// ----------------------------------------------------------------
//void lrece_update_value(lrecEntry* pe, char* new_value, int new_needs_freeing) {
//	if (pe == NULL) {
//		return;
//	}
//	if (pe->free_flags & FREE_ENTRY_VALUE) {
//		free(pe->value);
//	}
//	pe->value = new_value;
//	if (new_needs_freeing)
//		pe->free_flags |= FREE_ENTRY_VALUE;
//	else
//		pe->free_flags &= ~FREE_ENTRY_VALUE;
//}

//// ----------------------------------------------------------------
//void lrec_unlink(Lrec* prec, lrecEntry* pe) {
//	if (pe == prec->phead) {
//		if (pe == prec->ptail) {
//			prec->phead = NULL;
//			prec->ptail = NULL;
//		} else {
//			prec->phead = pe->pnext;
//			pe->pnext->pprev = NULL;
//		}
//	} else {
//		pe->pprev->pnext = pe->pnext;
//		if (pe == prec->ptail) {
//			prec->ptail = pe->pprev;
//		} else {
//			pe->pnext->pprev = pe->pprev;
//		}
//	}
//	prec->field_count--;
//}

//// ----------------------------------------------------------------
//static void lrec_link_at_head(Lrec* prec, lrecEntry* pe) {
//
//	if (prec->phead == NULL) {
//		pe->pprev   = NULL;
//		pe->pnext   = NULL;
//		prec->phead = pe;
//		prec->ptail = pe;
//	} else {
//		// [b,c,d] + a
//		pe->pprev   = NULL;
//		pe->pnext   = prec->phead;
//		prec->phead->pprev = pe;
//		prec->phead = pe;
//	}
//	prec->field_count++;
//}

//static void lrec_link_at_tail(Lrec* prec, lrecEntry* pe) {
//
//	if (prec->phead == NULL) {
//		pe->pprev   = NULL;
//		pe->pnext   = NULL;
//		prec->phead = pe;
//		prec->ptail = pe;
//	} else {
//		pe->pprev   = prec->ptail;
//		pe->pnext   = NULL;
//		prec->ptail->pnext = pe;
//		prec->ptail = pe;
//	}
//	prec->field_count++;
//}

//// ----------------------------------------------------------------
//void lrec_dump(Lrec* prec) {
//	lrec_dump_fp(prec, stdout);
//}

//void lrec_dump_fp(Lrec* prec, FILE* fp) {
//	if (prec == NULL) {
//		fprintf(fp, "NULL\n");
//		return;
//	}
//	fprintf(fp, "field_count = %d\n", prec->field_count);
//	fprintf(fp, "| phead: %16p | ptail %16p\n", prec->phead, prec->ptail);
//	for (lrecEntry* pe = prec->phead; pe != NULL; pe = pe->pnext) {
//		const char* key_string = (pe == NULL) ? "none" :
//			pe->key == NULL ? "null" :
//			pe->key;
//		const char* value_string = (pe == NULL) ? "none" :
//			pe->value == NULL ? "null" :
//			pe->value;
//		fprintf(fp,
//		"| prev: %16p curr: %16p next: %16p | key: %12s | value: %12s |\n",
//			pe->pprev, pe, pe->pnext,
//			key_string, value_string);
//	}
//}

//void lrec_dump_titled(char* msg, Lrec* prec) {
//	printf("%s:\n", msg);
//	lrec_dump(prec);
//	printf("\n");
//}

//// ----------------------------------------------------------------
//Lrec* lrec_literal_1(char* k1, char* v1) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_2(char* k1, char* v1, char* k2, char* v2) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_3(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	return prec;
//}

//Lrec* lrec_literal_4(char* k1, char* v1, char* k2, char* v2, char* k3, char* v3, char* k4, char* v4) {
//	Lrec* prec = lrec_unbacked_alloc();
//	lrec_put(prec, k1, v1, NO_FREE);
//	lrec_put(prec, k2, v2, NO_FREE);
//	lrec_put(prec, k3, v3, NO_FREE);
//	lrec_put(prec, k4, v4, NO_FREE);
//	return prec;
//}

//void lrec_print(Lrec* prec) {
//	FILE* output_stream = stdout;
//	char ors = '\n';
//	char ofs = ',';
//	char ops = '=';
//	if (prec == NULL) {
//		fputs("NULL", output_stream);
//		fputc(ors, output_stream);
//		return;
//	}
//	int nf = 0;
//	for (lrecEntry* pe = prec->phead; pe != NULL; pe = pe->pnext) {
//		if (nf > 0)
//			fputc(ofs, output_stream);
//		fputs(pe->key, output_stream);
//		fputc(ops, output_stream);
//		fputs(pe->value, output_stream);
//		nf++;
//	}
//	fputc(ors, output_stream);
//}

//char* lrec_sprint(Lrec* prec, char* ors, char* ofs, char* ops) {
//	string_builder_t* psb = sb_alloc(SB_ALLOC_LENGTH);
//	if (prec == NULL) {
//		sb_append_string(psb, "NULL");
//	} else {
//		int nf = 0;
//		for (lrecEntry* pe = prec->phead; pe != NULL; pe = pe->pnext) {
//			if (nf > 0)
//				sb_append_string(psb, ofs);
//			sb_append_string(psb, pe->key);
//			sb_append_string(psb, ops);
//			sb_append_string(psb, pe->value);
//			nf++;
//		}
//		sb_append_string(psb, ors);
//	}
//	char* rv = sb_finish(psb);
//	sb_free(psb);
//	return rv;
//}
