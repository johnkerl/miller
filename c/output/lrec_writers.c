#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "output/lrec_writers.h"

lrec_writer_t*  lrec_writer_alloc(cli_writer_opts_t* popts) {
	if (streq(popts->ofile_fmt, "dkvp")) {
		return lrec_writer_dkvp_alloc(popts->ors, popts->ofs, popts->ops);

	} else if (streq(popts->ofile_fmt, "json")) {
		return lrec_writer_json_alloc(popts->stack_json_output_vertically,
			popts->wrap_json_output_in_outer_list, popts->json_quote_int_keys,
			popts->json_quote_non_string_values,
			popts->output_json_flatten_separator, popts->ors);

	} else if (streq(popts->ofile_fmt, "csv")) {
		return lrec_writer_csv_alloc(popts->ors, popts->ofs, popts->oquoting,
			popts->headerless_csv_output);

	} else if (streq(popts->ofile_fmt, "csvlite")) {
		return lrec_writer_csvlite_alloc(popts->ors, popts->ofs, popts->headerless_csv_output);

	} else if (streq(popts->ofile_fmt, "markdown")) {
		return lrec_writer_markdown_alloc(popts->ors);

	} else if (streq(popts->ofile_fmt, "nidx")) {
		return lrec_writer_nidx_alloc(popts->ors, popts->ofs);

	} else if (streq(popts->ofile_fmt, "xtab")) {
		return lrec_writer_xtab_alloc(popts->ofs, popts->ops, popts->right_justify_xtab_value);

	} else if (streq(popts->ofile_fmt, "pprint")) {
		if (strlen(popts->ofs) != 1) {
			fprintf(stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
				MLR_GLOBALS.bargv0, popts->ofs);
			return NULL;
		} else {
			return lrec_writer_pprint_alloc(popts->ors, popts->ofs[0], popts->right_align_pprint,
				popts->pprint_barred);
		}

	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
lrec_writer_t* lrec_writer_alloc_or_die(cli_writer_opts_t* popts) {
	lrec_writer_t* plrec_writer = lrec_writer_alloc(popts);
	MLR_INTERNAL_CODING_ERROR_IF(plrec_writer == NULL);
	return plrec_writer;
}

// ----------------------------------------------------------------
void lrec_writer_print_all(lrec_writer_t* pwriter, FILE* fp, sllv_t* poutrecs, context_t* pctx) {
	while (poutrecs->phead != NULL) {
		lrec_t* poutrec = sllv_pop(poutrecs);
		pwriter->pprocess_func(pwriter->pvstate, fp, poutrec, pctx);
	}
}
