#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_json_state_t {
	unsigned long long counter;
	char* json_flatten_separator;

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
	int quote_json_values_always, char* json_flatten_separator)
{
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_json_state_t));
	pstate->quote_json_values_always = quote_json_values_always;
	pstate->counter = 0;
	pstate->json_flatten_separator = json_flatten_separator;

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

		// Use the mlhmmv printer since it naturally handles Miller-to-JSON key deconcatenation:
		// e.g. 'a:x=1,a:y=2' maps to '{"a":{"x":1,"y":2}}'.
		mlhmmv_t* pmap = mlhmmv_alloc();

		char* sep = pstate->json_flatten_separator;

		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			// strdup since strtok is destructive and CSV/PPRINT header fields
			// are shared across multiple records
			char* lkey = mlr_strdup_or_die(pe->key);
			char* lvalue = pe->value;

			sllmv_t* pmvkeys = sllmv_alloc();
			for (char* piece = strtok(lkey, sep); piece != NULL; piece = strtok(NULL, sep)) {
				mv_t mvkey = mv_from_string(piece, NO_FREE);
				sllmv_add(pmvkeys, &mvkey);
			}
			mv_t mvval = mv_from_string(lvalue, NO_FREE);
			mlhmmv_put(pmap, pmvkeys, &mvval);
			sllmv_free(pmvkeys);
			free(lkey);
		}

		if (pstate->stack_vertically)
			mlhmmv_print_json_stacked(pmap, pstate->quote_json_values_always);
		else
			mlhmmv_print_json_single_line(pmap, pstate->quote_json_values_always);

		mlhmmv_free(pmap);

	} else { // end of record stream
		fputs(pstate->after_records_at_end_of_stream, output_stream);
	}
}
