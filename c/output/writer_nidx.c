#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/writers.h"

typedef struct _writer_nidx_state_t {
	char rs;
	char fs;
} writer_nidx_state_t;

// ----------------------------------------------------------------
static void writer_nidx_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	writer_nidx_state_t* pstate = pvstate;
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

static void writer_nidx_free(void* pvstate) {
}

writer_t* writer_nidx_alloc(char rs, char fs) {
	writer_t* pwriter = mlr_malloc_or_die(sizeof(writer_t));

	writer_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(writer_nidx_state_t));
	pstate->rs = rs;
	pstate->fs = fs;
	pwriter->pvstate = (void*)pstate;

	pwriter->pwriter_func = &writer_nidx_func;
	pwriter->pfree_func   = &writer_nidx_free;

	return pwriter;
}
