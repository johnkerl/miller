#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/lrec_readers.h"
#include "input/byte_readers.h"

lrec_reader_t*  lrec_reader_alloc(cli_reader_opts_t* popts) {
	if (streq(popts->ifile_fmt, "dkvp")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_dkvp_alloc(popts->irs, popts->ifs, popts->ips, popts->allow_repeat_ifs);
		else
			return lrec_reader_stdio_dkvp_alloc(popts->irs, popts->ifs, popts->ips, popts->allow_repeat_ifs);
	} else if (streq(popts->ifile_fmt, "csv")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_csv_alloc(popts->irs, popts->ifs, popts->use_implicit_csv_header);
		else
			return lrec_reader_stdio_csv_alloc(popts->irs, popts->ifs, popts->use_implicit_csv_header);
	} else if (streq(popts->ifile_fmt, "csvlite")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_csvlite_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->use_implicit_csv_header);
		else
			return lrec_reader_stdio_csvlite_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs,
				popts->use_implicit_csv_header);
	} else if (streq(popts->ifile_fmt, "nidx")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_nidx_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs);
		else
			return lrec_reader_stdio_nidx_alloc(popts->irs, popts->ifs, popts->allow_repeat_ifs);
	} else if (streq(popts->ifile_fmt, "xtab")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_xtab_alloc(popts->ifs, popts->ips, popts->allow_repeat_ips);
		else
			return lrec_reader_stdio_xtab_alloc(popts->ifs, popts->ips, popts->allow_repeat_ips);
	} else if (streq(popts->ifile_fmt, "json")) {
		if (popts->use_mmap_for_read)
			return lrec_reader_mmap_json_alloc(popts->input_json_flatten_separator);
		else
			return lrec_reader_stdio_json_alloc(popts->input_json_flatten_separator);
	} else {
		return NULL;
	}
}
