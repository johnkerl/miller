#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"
#include "output/lrec_writers.h"

typedef struct _lrec_writer_json_state_t {
	unsigned long long counter;
	char* output_json_flatten_separator;

	int json_quote_int_keys;
	int json_quote_non_string_values;
	char* line_indent;
	char* before_records_at_start_of_stream1;
	char* between_records_after_start_of_stream;
	char* after_records_at_end_of_stream1;
	char* line_term;
	int stack_vertically;

} lrec_writer_json_state_t;

static void lrec_writer_json_free(lrec_writer_t* pwriter, context_t* pctx);
static void lrec_writer_json_process(void* pvstate, FILE* output_stream, lrec_t* prec,
	char* before_or_after_records, char* line_term);
static void lrec_writer_json_process_auto_line_term_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx);
static void lrec_writer_json_process_auto_line_term_no_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx);
static void lrec_writer_json_process_nonauto_line_term_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx);
static void lrec_writer_json_process_nonauto_line_term_no_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx);

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_json_alloc(int stack_vertically, int wrap_json_output_in_outer_list,
	int json_quote_int_keys, int json_quote_non_string_values,
	char* output_json_flatten_separator, char* line_term)
{
	lrec_writer_t* plrec_writer = mlr_malloc_or_die(sizeof(lrec_writer_t));

	lrec_writer_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_writer_json_state_t));
	pstate->json_quote_int_keys = json_quote_int_keys;
	pstate->json_quote_non_string_values = json_quote_non_string_values;
	pstate->counter = 0;
	pstate->output_json_flatten_separator = output_json_flatten_separator;

	// xxx pending reworked JSON-output logic (not always ending in LF; needing fflush),
	// to be implemented someday if ever. Workaround: pipe output to "jq .".
	//pstate->line_indent                           = wrap_json_output_in_outer_list ? "  "  : "";
	pstate->line_indent                           = wrap_json_output_in_outer_list ? ""  : "";
	pstate->before_records_at_start_of_stream1    = wrap_json_output_in_outer_list ? "[" : "";
	pstate->between_records_after_start_of_stream = wrap_json_output_in_outer_list ? ","   : "";
	pstate->after_records_at_end_of_stream1       = wrap_json_output_in_outer_list ? "]" : "";
	pstate->line_term                             = line_term;
	pstate->stack_vertically                      = stack_vertically;

	plrec_writer->pvstate = (void*)pstate;
	if (streq(line_term, "auto")) {
		plrec_writer->pprocess_func = wrap_json_output_in_outer_list
			? lrec_writer_json_process_auto_line_term_wrap
			: lrec_writer_json_process_auto_line_term_no_wrap;
	} else {
		plrec_writer->pprocess_func = wrap_json_output_in_outer_list
			? lrec_writer_json_process_nonauto_line_term_wrap
			: lrec_writer_json_process_nonauto_line_term_no_wrap;
	}
	plrec_writer->pfree_func    = lrec_writer_json_free;

	return plrec_writer;
}

static void lrec_writer_json_free(lrec_writer_t* pwriter, context_t* pctx) {
	free(pwriter->pvstate);
	free(pwriter);
}

// ----------------------------------------------------------------
static void lrec_writer_json_process_auto_line_term_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx)
{
	lrec_writer_json_process(pvstate, output_stream, prec, pctx->auto_line_term, pctx->auto_line_term);
}

static void lrec_writer_json_process_auto_line_term_no_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx)
{
	lrec_writer_json_process(pvstate, output_stream, prec, "", pctx->auto_line_term);
}

static void lrec_writer_json_process_nonauto_line_term_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx)
{
	lrec_writer_json_state_t* pstate = pvstate;
	lrec_writer_json_process(pvstate, output_stream, prec, pstate->line_term, pstate->line_term);
}

static void lrec_writer_json_process_nonauto_line_term_no_wrap(void* pvstate, FILE* output_stream, lrec_t* prec,
	context_t* pctx)
{
	lrec_writer_json_state_t* pstate = pvstate;
	lrec_writer_json_process(pvstate, output_stream, prec, "", pstate->line_term);
}

static void lrec_writer_json_process(void* pvstate, FILE* output_stream, lrec_t* prec,
	char* before_or_after_records, char* line_term)
{
	lrec_writer_json_state_t* pstate = pvstate;
	if (prec != NULL) { // not end of record stream
		if (pstate->counter++ == 0) {
			fputs(pstate->before_records_at_start_of_stream1, output_stream);
			fputs(before_or_after_records, output_stream);
		} else {
			fputs(pstate->between_records_after_start_of_stream, output_stream);
		}

		// Use the mlhmmv printer since it naturally handles Miller-to-JSON key deconcatenation:
		// e.g. 'a:x=1,a:y=2' maps to '{"a":{"x":1,"y":2}}'.
		mlhmmv_root_t* pmap = mlhmmv_root_alloc();

		char* sep = pstate->output_json_flatten_separator;
		int seplen = strlen(sep);

		for (lrece_t* pe = prec->phead; pe != NULL; pe = pe->pnext) {
			// strdup since strmsep is destructive and CSV/PPRINT header fields
			// are shared across multiple records
			char* lkey = mlr_strdup_or_die(pe->key);
			char* lvalue = pe->value;

			sllmv_t* pmvkeys = sllmv_alloc();
			char* walker = lkey;
			char* piece = NULL;
			while ((piece = mlr_strmsep(&walker, sep, seplen)) != NULL) {
				mv_t mvkey = mv_from_string(piece, NO_FREE);
				sllmv_append_no_free(pmvkeys, &mvkey);
			}
			mv_t mvval = mv_from_string(lvalue, NO_FREE);
			mlhmmv_root_put_terminal(pmap, pmvkeys, &mvval);
			sllmv_free(pmvkeys);
			free(lkey);
		}

		if (pstate->stack_vertically)
			mlhmmv_root_print_json_stacked(pmap, pstate->json_quote_int_keys, pstate->json_quote_non_string_values,
				pstate->line_indent, line_term, output_stream);
		else
			mlhmmv_root_print_json_single_lines(pmap, pstate->json_quote_int_keys,
				pstate->json_quote_non_string_values, line_term, output_stream);

		mlhmmv_root_free(pmap);

		lrec_free(prec); // end of baton-pass

	} else { // end of record stream
		if (pstate->counter == 0) {
			fputs(pstate->before_records_at_start_of_stream1, output_stream);
		}
		fputs(pstate->after_records_at_end_of_stream1, output_stream);
		fputs(before_or_after_records, output_stream);
	}
}
