// ================================================================
// Miller command-line parsing
// ================================================================

#ifndef MLRCLI_H
#define MLRCLI_H

#include "containers/sllv.h"
#include "input/lrec_reader.h"
#include "mapping/mapper.h"
#include "output/lrec_writer.h"

// xxx move to another header file ...
#define QUOTE_ALL     0xb1
#define QUOTE_NONE    0xb2
#define QUOTE_MINIMAL 0xb3
#define QUOTE_NUMERIC 0xb4

typedef struct _cli_opts_t {
	char  irs;
	char  ifs;
	char  ips;
	int   allow_repeat_ifs;
	int   allow_repeat_ips;
	int   use_mmap_for_read;
	char* ifile_fmt;
	char* ofile_fmt;

	char* ors;
	char* ofs;
	char* ops;

	char* ofmt;
	int   oquoting;

	lrec_reader_t* plrec_reader;
	sllv_t*        pmapper_list;
	lrec_writer_t* plrec_writer;

	char** filenames; // null-terminated

} cli_opts_t;

cli_opts_t* parse_command_line(int argc, char** argv);
void cli_opts_free(cli_opts_t* popts);

#endif // MLRCLI_H
