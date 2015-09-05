#include "lib/mlrutil.h"
#include "input/lrec_readers.h"

lrec_reader_t*  lrec_reader_alloc(char* fmtdesc, int use_mmap, char irs, char ifs, int allow_repeat_ifs,
	char ips, int allow_repeat_ips)
{
	if (streq(fmtdesc, "dkvp")) {
		if (use_mmap)
			return lrec_reader_mmap_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
		else
			return lrec_reader_stdio_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
	} else if (streq(fmtdesc, "csv")) {
		// xxx not now
		// if (use_mmap)
			//return lrec_reader_mmap_csv_alloc(irs, ifs, allow_repeat_ifs);
		//else
			return lrec_reader_stdio_csv_alloc(irs, ifs, allow_repeat_ifs);
	} else if (streq(fmtdesc, "csvlite")) {
		if (use_mmap)
			return lrec_reader_mmap_csvlite_alloc(irs, ifs, allow_repeat_ifs);
		else
			return lrec_reader_stdio_csvlite_alloc(irs, ifs, allow_repeat_ifs);
	} else if (streq(fmtdesc, "csvex")) {
		return lrec_reader_stdio_csvex_alloc(irs, ifs, allow_repeat_ifs);
	} else if (streq(fmtdesc, "nidx")) {
		if (use_mmap)
			return lrec_reader_mmap_nidx_alloc(irs, ifs, allow_repeat_ifs);
		else
			return lrec_reader_stdio_nidx_alloc(irs, ifs, allow_repeat_ifs);
	} else if (streq(fmtdesc, "xtab")) {
		if (use_mmap)
			return lrec_reader_mmap_xtab_alloc(irs, ips, TRUE/*allow_repeat_ips*/);
		else
			return lrec_reader_stdio_xtab_alloc(ips, TRUE); // xxx parameterize allow_repeat_ips
	} else {
		return NULL;
	}
}
