#include <stdlib.h>
#include "containers/mixutil.h"
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_markdown_state_t {
	int   onr;
	char* ors;
	long long num_header_lines_output;
	slls_t* plast_header_output;
} lrec_writer_markdown_state_t;

static void lrec_writer_markdown_free(lrec_writer_t* pwriter);
static void lrec_writer_markdown_process(FILE* output_stream, lrec_t* prec, void* pvstate);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_markdown_alloc(char* ors) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_markdown_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_markdown_state_t));
	pstate->onr                     = 0;
	pstate->ors                     = ors;
	pstate->num_header_lines_output = 0LL;
	pstate->plast_header_output     = NULL;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = lrec_writer_markdown_process;
	plrec_writer->pfree_func    = lrec_writer_markdown_free;

	return plrec_writer;
}

static void lrec_writer_markdown_free(lrec_writer_t* pwriter) {
	lrec_writer_markdown_state_t* pstate = pwriter->pvstate;
	slls_free(pstate->plast_header_output);
	free(pstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_markdown_process(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_markdown_state_t* pstate = pvstate;
	char* ors = pstate->ors;

	if (pstate->plast_header_output != NULL) {
		if (!lrec_keys_equal_list(prec, pstate->plast_header_output)) {
			slls_free(pstate->plast_header_output);
			pstate->plast_header_output = NULL;
			if (pstate->num_header_lines_output > 0LL)
				fputs(ors, output_stream);
		}
	}

	if (pstate->plast_header_output == NULL) {
		fputc('|', output_stream);
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			fputc(' ', output_stream);
			fputs(pe->key, output_stream);
			fputs(" |", output_stream);
		}
		fputs(ors, output_stream);

		fputc('|', output_stream);
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			fputs(" --- |", output_stream);
		}
		fputs(ors, output_stream);

		pstate->plast_header_output = mlr_copy_keys_from_record(prec);
		pstate->num_header_lines_output++;
	}

	fputc('|', output_stream);
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		fputc(' ', output_stream);
		fputs(pe->value, output_stream);
		fputs(" |", output_stream);
	}
	fputs(ors, output_stream);
	pstate->onr++;

	lrec_free(prec); // end of baton-pass
}
