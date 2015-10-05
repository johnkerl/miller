#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

// ----------------------------------------------------------------
// Note: If OPS is single-character then we can do alignment of the form
//
//   ab  123
//   def 4567
//
// On the other hand, if it's multi-character, we won't be able to align
// neatly in all cases. Yet we do allow multi-character OPS, just without
// repetition: if someone wants to use OPS ": " and format data as
//
//   ab: 123
//   def: 4567
//
// then they can do that.
// ----------------------------------------------------------------

typedef struct _lrec_writer_xtab_state_t {
	char* ofs;
	char* ops;
	int   opslen;
	long long record_count;
} lrec_writer_xtab_state_t;

static void lrec_writer_xtab_free(void* pvstate);
static void lrec_writer_xtab_process_aligned(FILE* output_stream, lrec_t* prec, void* pvstate);
static void lrec_writer_xtab_process_unaligned(FILE* output_stream, lrec_t* prec, void* pvstate);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_xtab_alloc(char* ofs, char* ops) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_xtab_state_t));
	pstate->ofs          = ofs;
	pstate->ops          = ops;
	pstate->opslen       = strlen(ops);
	pstate->record_count = 0LL;

	plrec_writer->pvstate       = pstate;
	plrec_writer->pprocess_func = (pstate->opslen == 1)
		? lrec_writer_xtab_process_aligned
		: lrec_writer_xtab_process_unaligned;
	plrec_writer->pfree_func    = lrec_writer_xtab_free;

	return plrec_writer;
}

static void lrec_writer_xtab_free(void* pvstate) {
}

// ----------------------------------------------------------------
static void lrec_writer_xtab_process_aligned(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_xtab_state_t* pstate = pvstate;
	if (pstate->record_count > 0LL)
		fputs(pstate->ofs, output_stream);
	pstate->record_count++;

	int max_key_width = 1;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		int key_width = strlen_for_utf8_display(pe->key);
		if (key_width > max_key_width)
			max_key_width = key_width;
	}

	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		// "%-*s" fprintf format isn't correct for non-ASCII UTF-8
		fprintf(output_stream, "%s", pe->key);
		int d = max_key_width - strlen_for_utf8_display(pe->key);
		for (int i = 0; i < d; i++)
			fputs(pstate->ops, output_stream);
		fprintf(output_stream, "%s%s%s", pstate->ops, pe->value, pstate->ofs);
	}
	lrec_free(prec); // xxx cmt mem-mgmt
}

static void lrec_writer_xtab_process_unaligned(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_xtab_state_t* pstate = pvstate;
	if (pstate->record_count > 0LL)
		fputs(pstate->ofs, output_stream);
	pstate->record_count++;

	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		// "%-*s" fprintf format isn't correct for non-ASCII UTF-8
		fprintf(output_stream, "%s%s%s%s", pe->key, pstate->ops, pe->value, pstate->ofs);
	}
	lrec_free(prec); // xxx cmt mem-mgmt
}
