#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "output/lrec_writers.h"

// xxx bag up in popts->writer_opts
lrec_writer_t*  lrec_writer_alloc(char* fmtdesc, char* ors, char* ofs, char* ops,
	int headerless_csv_output, int oquoting,
	int left_align_pprint,
	int right_justify_xtab_value,
	char* json_flatten_separator,
	int quote_json_values_always,
	int stack_json_output_vertically,
	int wrap_json_output_in_outer_list)
{
	if (streq(fmtdesc, "dkvp")) {
		return lrec_writer_dkvp_alloc(ors, ofs, ops);

	} else if (streq(fmtdesc, "json")) {
		return lrec_writer_json_alloc(stack_json_output_vertically,
			wrap_json_output_in_outer_list, quote_json_values_always, json_flatten_separator);

	} else if (streq(fmtdesc, "csv")) {
		return lrec_writer_csv_alloc(ors, ofs, oquoting,
			headerless_csv_output);

	} else if (streq(fmtdesc, "csvlite")) {
		return lrec_writer_csvlite_alloc(ors, ofs, headerless_csv_output);

	} else if (streq(fmtdesc, "markdown")) {
		return lrec_writer_markdown_alloc(ors);

	} else if (streq(fmtdesc, "nidx")) {
		return lrec_writer_nidx_alloc(ors, ofs);

	} else if (streq(fmtdesc, "xtab")) {
		return lrec_writer_xtab_alloc(ofs, ops, right_justify_xtab_value);

	} else if (streq(fmtdesc, "pprint")) {
		if (strlen(ofs) != 1) {
			fprintf(stderr, "%s: OFS for PPRINT format must be single-character; got \"%s\".\n",
				MLR_GLOBALS.bargv0, ofs);
			return NULL;
		} else {
			return lrec_writer_pprint_alloc(ors, ofs[0], left_align_pprint);
		}

	} else {
		return NULL;
	}
}

void lrec_writer_print_all(lrec_writer_t* pwriter, FILE* fp, sllv_t* poutrecs) {
	while (poutrecs->phead != NULL) {
		lrec_t* poutrec = sllv_pop(poutrecs);
		pwriter->pprocess_func(pwriter->pvstate, fp, poutrec);
	}
}
