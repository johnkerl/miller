#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/writers.h"

typedef struct _writer_dkvp_state_t {
	char rs;
	char fs;
	char ps;
} writer_dkvp_state_t;

// ----------------------------------------------------------------
static void writer_dkvp_func(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	writer_dkvp_state_t* pstate = pvstate;
	char rs = pstate->rs;
	char fs = pstate->fs;
	char ps = pstate->ps;

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(fs, output_stream);
		fputs(pe->key, output_stream);
		fputc(ps, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputc(rs, output_stream);
	lrec_free(prec); // xxx cmt mem-mgmt
}

static void writer_dkvp_free(void* pvstate) {
}

writer_t* writer_dkvp_alloc(char rs, char fs, char ps) {
	writer_t* pwriter = mlr_malloc_or_die(sizeof(writer_t));

	writer_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(writer_dkvp_state_t));
	pstate->rs = rs;
	pstate->fs = fs;
	pstate->ps = ps;
	pwriter->pvstate = (void*)pstate;

	pwriter->pwriter_func = &writer_dkvp_func;
	pwriter->pfree_func   = &writer_dkvp_free;

	return pwriter;
}
