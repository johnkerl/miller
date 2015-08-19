#include <stdlib.h>
#include "containers/mixutil.h"
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

// See https://github.com/johnkerl/miller/issues/4
// Temporary status:
// * --csv     from the command line maps into the (existing) csvlite I/O
// * --csvlite from the command line maps into the (existing) csvlite I/O
// * --csvex   from the command line maps into the (new & experimental & unadvertised) rfc-csv I/O
// Ultimate status:
// * --csvlite from the command line will maps into the csvlite I/O
// * --csv     from the command line will maps into the rfc-csv I/O

typedef struct _lrec_writer_csvlite_state_t {
	int  onr;
	char ors;
	char ofs;
	long long num_header_lines_output;
	slls_t* plast_header_output;
} lrec_writer_csvlite_state_t;

// ----------------------------------------------------------------
// xxx cmt mem-mgmt

static void lrec_writer_csvlite_process(FILE* output_stream, lrec_t* prec, void* pvstate) {
	if (prec == NULL)
		return;
	lrec_writer_csvlite_state_t* pstate = pvstate;
	char ors = pstate->ors;
	char ofs = pstate->ofs;

	if (pstate->plast_header_output != NULL) {
		// xxx make a fcn to compare these w/o copy: put it in mixutil.
		if (!lrec_keys_equal_list(prec, pstate->plast_header_output)) {
			slls_free(pstate->plast_header_output);
			pstate->plast_header_output = NULL;
			if (pstate->num_header_lines_output > 0LL)
				fputc(ors, output_stream);
		}
	}

	if (pstate->plast_header_output == NULL) {
		int nf = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			if (nf > 0)
				fputc(ofs, output_stream);
			fputs(pe->key, output_stream);
			nf++;
		}
		fputc(ors, output_stream);
		pstate->plast_header_output = mlr_copy_keys_from_record(prec);
		pstate->num_header_lines_output++;
	}

	int nf = 0;
	for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
		if (nf > 0)
			fputc(ofs, output_stream);
		fputs(pe->value, output_stream);
		nf++;
	}
	fputc(ors, output_stream);
	pstate->onr++;

	lrec_free(prec); // xxx cmt mem-mgmt
}

static void lrec_writer_csvlite_free(void* pvstate) {
	lrec_writer_csvlite_state_t* pstate = pvstate;
	if (pstate->plast_header_output != NULL) {
		slls_free(pstate->plast_header_output);
		pstate->plast_header_output = NULL;
	}
}

lrec_writer_t* lrec_writer_csvlite_alloc(char ors, char ofs) {
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_csvlite_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_csvlite_state_t));
	pstate->onr                     = 0;
	pstate->ors                     = ors;
	pstate->ofs                     = ofs;
	pstate->num_header_lines_output = 0LL;
	pstate->plast_header_output     = NULL;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = lrec_writer_csvlite_process;
	plrec_writer->pfree_func    = lrec_writer_csvlite_free;

	return plrec_writer;
}
