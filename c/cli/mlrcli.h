// ================================================================
// Miller command-line parsing
// ================================================================

#ifndef MLRCLI_H
#define MLRCLI_H

#include "containers/sllv.h"
#include "input/lrec_reader.h"
#include "mapping/mapper.h"
#include "output/lrec_writer.h"

typedef struct _cli_opts_t {
	char* irs;
	char* ifs;
	char* ips;
	int   allow_repeat_ifs;
	int   allow_repeat_ips;
	int   use_implicit_csv_header;
	int   headerless_csv_output;
	int   use_mmap_for_read;
	char* ifile_fmt;
	char* ofile_fmt;

	char* ors;
	char* ofs;
	char* ops;

	int   right_justify_xtab_value;

	char* ofmt;
	int   oquoting;

	lrec_reader_t* plrec_reader;
	sllv_t*        pmapper_list;
	lrec_writer_t* plrec_writer;

	char*  prepipe;   // Command for popen on input, e.g. "zcat -cf <". Can be null.
	char** filenames; // null-terminated

} cli_opts_t;

cli_opts_t* parse_command_line(int argc, char** argv);
void cli_opts_free(cli_opts_t* popts);

#endif // MLRCLI_H
