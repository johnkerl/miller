#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_json_state_t {
	unsigned long long counter;

	int quote_json_values_always;
	char* before_records_at_start_of_stream;
	char* between_records_after_start_of_stream;
	char* after_records_at_end_of_stream;
	int stack_vertically;

} lrec_writer_json_state_t;

static void lrec_writer_json_free(lrec_writer_t* pwriter);
static void lrec_writer_json_process(FILE* output_stream, lrec_t* prec, void* pvstate);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_json_alloc(int stack_vertically, int wrap_json_output_in_outer_list,
	int quote_json_values_always)
{
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_json_state_t));
	pstate->quote_json_values_always = quote_json_values_always;
	pstate->counter = 0;

	pstate->before_records_at_start_of_stream     = wrap_json_output_in_outer_list ? "[\n" : "";
	pstate->between_records_after_start_of_stream = wrap_json_output_in_outer_list ? ","   : "";
	pstate->after_records_at_end_of_stream        = wrap_json_output_in_outer_list ? "]\n" : "";
	pstate->stack_vertically                      = stack_vertically;

	plrec_writer->pvstate       = (void*)pstate;
	plrec_writer->pprocess_func = lrec_writer_json_process;
	plrec_writer->pfree_func    = lrec_writer_json_free;

	return plrec_writer;
}

static void lrec_writer_json_free(lrec_writer_t* pwriter) {
	free(pwriter->pvstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_json_process(FILE* output_stream, lrec_t* prec, void* pvstate) {
	lrec_writer_json_state_t* pstate = pvstate;
	if (prec != NULL) { // not end of record stream
		if (pstate->counter++ == 0)
			printf("%s", pstate->before_records_at_start_of_stream);
		else
			printf("%s", pstate->between_records_after_start_of_stream);
		mlhmmv_t* pmap = mlhmmv_alloc();

		char* flatten_sep = ":"; // xxx temp; needs to be parameterized

		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			char* lkey = pe->key;
			char* lvalue = pe->value;

			sllmv_t* pmvkeys = sllmv_alloc();
			for (char* piece = strtok(lkey, flatten_sep); piece != NULL; piece = strtok(NULL, flatten_sep)) {
				mv_t mvkey = mv_from_string(piece, NO_FREE);
				sllmv_add(pmvkeys, &mvkey);
			}
			mv_t mvval = mv_from_string(lvalue, NO_FREE);
			mlhmmv_put(pmap, pmvkeys, &mvval);
			sllmv_free(pmvkeys);
		}

		if (pstate->stack_vertically)
			mlhmmv_print_stacked(pmap, pstate->quote_json_values_always);
		else
			mlhmmv_print_single_line(pmap, pstate->quote_json_values_always);

		mlhmmv_free(pmap);

	} else { // end of record stream
		fputs(pstate->after_records_at_end_of_stream, output_stream);
	}
}
