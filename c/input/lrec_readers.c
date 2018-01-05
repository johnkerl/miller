#include "cli/comment_handling.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/lrec_readers.h"
#include "input/byte_readers.h"

lrec_reader_t*  lrec_reader_alloc(cli_reader_opts_t* popts) {
	if (streq(popts->ifile_fmt, "gen")) {
		generator_opts_t* pgopts = &popts->generator_opts;
		return lrec_reader_gen_alloc(pgopts->field_name, pgopts->start, pgopts->stop, pgopts->step);
	} else if (streq(popts->ifile_fmt, "dkvp")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_dkvp_alloc(popts->irs, popts->ifs, popts->ips, popts->allow_repeat_ifs,
				popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_dkvp_alloc(popts->irs, popts->ifs, popts->ips, popts->allow_repeat_ifs,
				popts->comment_handling, popts->comment_string);
	} else if (streq(popts->ifile_fmt, "csv")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_csv_alloc(popts->irs, popts->ifs, popts->use_implicit_csv_header,
				popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_csv_alloc(popts->irs, popts->ifs, popts->use_implicit_csv_header,
				popts->comment_handling, popts->comment_string);
	} else if (streq(popts->ifile_fmt, "csvlite")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_csvlite_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->use_implicit_csv_header, popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_csvlite_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->use_implicit_csv_header, popts->comment_handling, popts->comment_string);
	} else if (streq(popts->ifile_fmt, "nidx")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_nidx_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_nidx_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->comment_handling, popts->comment_string);
	} else if (streq(popts->ifile_fmt, "xtab")) {
		// Use stdio-xtab for comment handling; not supported in the mmap-xtab reader.
		if (popts->use_mmap_for_read && popts->comment_string == NULL)
			return lrec_reader_mmap_xtab_alloc(popts->ifs, popts->ips, popts->allow_repeat_ips,
				popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_xtab_alloc(popts->ifs, popts->ips, popts->allow_repeat_ips,
				popts->comment_handling, popts->comment_string);
	} else if (streq(popts->ifile_fmt, "json")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_json_alloc(popts->input_json_flatten_separator,
				popts->json_array_ingest, popts->irs, popts->comment_handling, popts->comment_string);
		else
			return lrec_reader_stdio_json_alloc(popts->input_json_flatten_separator,
				popts->json_array_ingest, popts->irs, popts->comment_handling, popts->comment_string);
	} else {
		return NULL;
	}
}

lrec_reader_t* lrec_reader_alloc_or_die(cli_reader_opts_t* popts) {
	lrec_reader_t* plrec_reader = lrec_reader_alloc(popts);
	if (plrec_reader == NULL) {
		fprintf(stderr, "%s: unrecognized input-file format \"%s\".\n",
			MLR_GLOBALS.bargv0, popts->ifile_fmt);
		exit(1);
	}
	return plrec_reader;
}
