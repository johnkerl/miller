#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/writers.h"

typedef struct _lrec_writer_nidx_state_t {
	char rs;
	char fs;
} lrec_writer_nidx_state_t;

// ----------------------------------------------------------------
static void lrec_writer_nidx_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_nidx_state_t* pstate = pvstate;
	char rs = pstate->rs;
	char fs = pstate->fs;

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(fs, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputc(rs, output_stream);
	lrec_free(prec); // xxx cmt mem-mgmt
}

static void lrec_writer_nidx_free(void* pvstate) {
}

lrec_writer_t* lrec_writer_nidx_alloc(char rs, char fs) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_nidx_state_t));
	pstate->rs = rs;
	pstate->fs = fs;
	plrec_writer->pvstate = (void*)pstate;

	plrec_writer->plrec_writer_func = &lrec_writer_nidx_func;
	plrec_writer->pfree_func   = &lrec_writer_nidx_free;

	return plrec_writer;
}
