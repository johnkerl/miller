#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/lrec_readers.h"
#include "input/byte_readers.h"

lrec_reader_t*  lrec_reader_alloc(char* fmtdesc, int use_mmap, char* irs, char* ifs, int allow_repeat_ifs,
	char* ips, int allow_repeat_ips, int use_implicit_csv_header, char* json_flatten_separator)
{
	if (streq(fmtdesc, "dkvp")) {
		if (use_mmap)
			return lrec_reader_mmap_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
		else
			return lrec_reader_stdio_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
	} else if (streq(fmtdesc, "csv")) {
		if (use_mmap)
			return lrec_reader_mmap_csv_alloc(irs, ifs, use_implicit_csv_header);
		else
			return lrec_reader_stdio_csv_alloc(irs, ifs, use_implicit_csv_header);
	} else if (streq(fmtdesc, "csvlite")) {
		if (use_mmap)
			return lrec_reader_mmap_csvlite_alloc(irs, ifs, allow_repeat_ifs, use_implicit_csv_header);
		else
			return lrec_reader_stdio_csvlite_alloc(irs, ifs, allow_repeat_ifs, use_implicit_csv_header);
	} else if (streq(fmtdesc, "nidx")) {
		if (use_mmap)
			return lrec_reader_mmap_nidx_alloc(irs, ifs, allow_repeat_ifs);
		else
			return lrec_reader_stdio_nidx_alloc(irs, ifs, allow_repeat_ifs);
	} else if (streq(fmtdesc, "xtab")) {
		if (use_mmap)
			return lrec_reader_mmap_xtab_alloc(ifs, ips, allow_repeat_ips);
		else
			return lrec_reader_stdio_xtab_alloc(ifs, ips, allow_repeat_ips);
	} else if (streq(fmtdesc, "json")) {
		if (use_mmap)
			return lrec_reader_mmap_json_alloc(json_flatten_separator);
		else
			fprintf(stderr, "The JSON stdio reader is not yet implemented.\n");
			return NULL; // xxx stub
	} else {
		return NULL;
	}
}
