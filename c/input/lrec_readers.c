#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/lrec_readers.h"
#include "input/byte_readers.h"

static char xxx_temp_check_single_char_separator(char* name, char* value) {
	if (strlen(value) != 1) {
		fprintf(stderr,
			"%s: multi-character separators are not yet supported for all formats. Got %s=\"%s\".\n",
			MLR_GLOBALS.argv0, name, value);
		exit(1);
	}
	return value[0];
}

lrec_reader_t*  lrec_reader_alloc(char* fmtdesc, int use_mmap, char* irs, char* ifs, int allow_repeat_ifs,
	char* ips, int allow_repeat_ips)
{
	// xxx refactor for https://github.com/johnkerl/miller/issues/51 et al.
	byte_reader_t* pbr = use_mmap ? mmap_byte_reader_alloc() : stdio_byte_reader_alloc();

	if (streq(fmtdesc, "dkvp")) {
		if (use_mmap)
			return lrec_reader_mmap_dkvp_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ifs", ifs),
				xxx_temp_check_single_char_separator("ips", ips),
				allow_repeat_ifs);
		else
			return lrec_reader_stdio_dkvp_alloc(
				irs,
				xxx_temp_check_single_char_separator("ifs", ifs),
				xxx_temp_check_single_char_separator("ips", ips),
				allow_repeat_ifs);
	} else if (streq(fmtdesc, "csv")) {
		return lrec_reader_csv_alloc(pbr, irs, ifs);
	} else if (streq(fmtdesc, "csvlite")) {
		if (use_mmap)
			return lrec_reader_mmap_csvlite_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ifs", ifs),
				allow_repeat_ifs);
		else
			return lrec_reader_stdio_csvlite_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ifs", ifs),
				allow_repeat_ifs);
	} else if (streq(fmtdesc, "nidx")) {
		if (use_mmap)
			return lrec_reader_mmap_nidx_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ifs", ifs),
				allow_repeat_ifs);
		else
			return lrec_reader_stdio_nidx_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ifs", ifs),
				allow_repeat_ifs);
	} else if (streq(fmtdesc, "xtab")) {
		if (use_mmap)
			return lrec_reader_mmap_xtab_alloc(
				xxx_temp_check_single_char_separator("irs", irs),
				xxx_temp_check_single_char_separator("ips", ips),
				TRUE/*allow_repeat_ips*/);
		else
			return lrec_reader_stdio_xtab_alloc(
				xxx_temp_check_single_char_separator("ips", ips),
				TRUE); // xxx parameterize allow_repeat_ips
	} else {
		return NULL;
	}
}
