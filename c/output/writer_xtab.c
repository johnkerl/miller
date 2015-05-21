#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "output/writers.h"

typedef struct _lrec_writer_xtab_state_t {
	long long record_count;
} lrec_writer_xtab_state_t;

// ----------------------------------------------------------------
static void lrec_writer_xtab_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_xtab_state_t* pstate = pvstate;
	if (pstate->record_count > 0LL)
		fprintf(output_stream, "\n");
	pstate->record_count++;

	int max_key_width = 1;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		int key_width = strlen(pe->key);
		if (key_width > max_key_width)
			max_key_width = key_width;
	}

	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		fprintf(output_stream, "%-*s %s\n", max_key_width, pe->key, pe->value);
	}
	lrec_free(prec); // xxx cmt mem-mgmt
}

static void lrec_writer_xtab_free(void* pvstate) {
}

lrec_writer_t* lrec_writer_xtab_alloc() {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_xtab_state_t));
	pstate->record_count = 0LL;
	plrec_writer->pvstate = pstate;

	plrec_writer->plrec_writer_func = &lrec_writer_xtab_func;
	plrec_writer->pfree_func   = &lrec_writer_xtab_free;

	return plrec_writer;
}
