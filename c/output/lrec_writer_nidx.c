#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_nidx_state_t {
	char* ors;
	char* ofs;
} lrec_writer_nidx_state_t;

static void lrec_writer_nidx_free(lrec_writer_t* pwriter, context_t* pctx);
static void lrec_writer_nidx_process(void* pvstate, FILE* output_stream, lrec_t* prec, char* ors);
static void lrec_writer_nidx_process_auto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx);
static void lrec_writer_nidx_process_nonauto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_nidx_alloc(char* ors, char* ofs) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_nidx_state_t));
	pstate->ors = ors;
	pstate->ofs = ofs;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = streq(ors, "auto")
		? lrec_writer_nidx_process_auto_ors
		: lrec_writer_nidx_process_nonauto_ors;
	plrec_writer->pfree_func    = lrec_writer_nidx_free;

	return plrec_writer;
}

static void lrec_writer_nidx_free(lrec_writer_t* pwriter, context_t* pctx) {
	free(pwriter->pvstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_nidx_process_auto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx) {
	lrec_writer_nidx_process(pvstate, output_stream, prec, pctx->auto_line_term);
}

static void lrec_writer_nidx_process_nonauto_ors(void* pvstate, FILE* output_stream, lrec_t* prec, context_t* pctx) {
	lrec_writer_nidx_state_t* pstate = pvstate;
	lrec_writer_nidx_process(pvstate, output_stream, prec, pstate->ors);
}

static void lrec_writer_nidx_process(void* pvstate, FILE* output_stream, lrec_t* prec, char* ors) {
	if (prec == NULL)
		return;
	lrec_writer_nidx_state_t* pstate = pvstate;
	char* ofs = pstate->ofs;

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputs(ofs, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputs(ors, output_stream);
	lrec_free(prec); // end of baton-pass
}
