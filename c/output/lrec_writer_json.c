#include <stdlib.h>
#include "lib/mlrutil.h"
#include "output/lrec_writers.h"

// xxx under construction

typedef struct _lrec_writer_json_state_t {
	unsigned long long counter;

	int quote_json_values_always;
	char* before_records_at_start_of_stream;
	char* between_records_after_start_of_stream;
	char* left_curly_at_start_of_record;
	char* between_fields;
	char* right_curly_at_end_of_record;
	char* after_records_at_end_of_stream;

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

	pstate->before_records_at_start_of_stream     = wrap_json_output_in_outer_list ? "[\n"   : "";
	pstate->between_records_after_start_of_stream = wrap_json_output_in_outer_list ? ","     : "";
	pstate->left_curly_at_start_of_record         = stack_vertically               ? "{\n  " : "{ ";
	pstate->between_fields                        = stack_vertically               ? ",\n  " : ", ";
	pstate->right_curly_at_end_of_record          = stack_vertically               ? "\n}\n" : " }\n";
	pstate->after_records_at_end_of_stream        = wrap_json_output_in_outer_list ? "]\n"   : "";

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

		if (pstate->counter == 0)
			fputs(pstate->before_records_at_start_of_stream, output_stream);
		else
			fputs(pstate->between_records_after_start_of_stream, output_stream);

		fputs(pstate->left_curly_at_start_of_record, output_stream);

		int nf = 0;
		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			if (nf > 0)
				fputs(pstate->between_fields, output_stream);

			fputs("\"", output_stream);
			fputs(pe->key, output_stream);
			fputs("\"", output_stream);

			if (pstate->quote_json_values_always) {
				fputs(": \"", output_stream);
				fputs(pe->value, output_stream);
				fputs("\"", output_stream);
			} else {
				double unused;
				if (mlr_try_float_from_string(pe->value, &unused)) {
					fputs(": ", output_stream);
					fputs(pe->value, output_stream);
				} else {
					fputs(": \"", output_stream);
					fputs(pe->value, output_stream);
					fputs("\"", output_stream);
				}
			}
			nf++;
		}
		fputs(pstate->right_curly_at_end_of_record, output_stream);

		pstate->counter++;
		lrec_free(prec); // end of baton-pass

	} else { // end of record stream
		fputs(pstate->after_records_at_end_of_stream, output_stream);
	}
}
