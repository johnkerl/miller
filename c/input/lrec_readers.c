#include "lib/mlrutil.h"
#include "input/lrec_readers.h"
#include "input/byte_readers.h"

lrec_reader_t*  lrec_reader_alloc(char* fmtdesc, int use_mmap, char irs, char ifs, int allow_repeat_ifs,
	char ips, int allow_repeat_ips)
{
	// xxx refactor for https://github.com/johnkerl/miller/issues/51 et al.
	byte_reader_t* pbr = use_mmap ? mmap_byte_reader_alloc() : stdio_byte_reader_alloc();

	if (streq(fmtdesc, "dkvp")) {
		if (use_mmap)
			return lrec_reader_mmap_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
		else
			return lrec_reader_stdio_dkvp_alloc(irs, ifs, ips, allow_repeat_ifs);
	} else if (streq(fmtdesc, "csv")) {
		return lrec_reader_csv_alloc(pbr, irs, ifs);
	} else if (streq(fmtdesc, "csvlite")) {
		if (use_mmap)
			return lrec_reader_mmap_csvlite_alloc(irs, ifs, allow_repeat_ifs);
		else
			return lrec_reader_stdio_csvlite_alloc(irs, ifs, allow_repeat_ifs);
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
