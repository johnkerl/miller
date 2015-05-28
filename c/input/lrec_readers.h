#ifndef LREC_READERS_H
#define LREC_READERS_H
#include "input/lrec_reader_stdio.h"
#include "input/lrec_reader_mmap.h"

// ----------------------------------------------------------------
// Primary entry points
lrec_reader_stdio_t*  lrec_reader_stdio_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_xtab_alloc(char ips, int allow_repeat_ips);

lrec_reader_mmap_t*  lrec_reader_mmap_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_xtab_alloc(char irs, char ips, int allow_repeat_ips);

// ----------------------------------------------------------------
// These entry points are made public for unit test
lrec_t* lrec_parse_stdio_nidx(char* line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_dkvp(char* line, char ifs, char ips, int allow_repeat_ifs);
slls_t* split_csv_header_line(char* line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_csv_data_line(header_keeper_t* pheader_keeper, char* data_line, char ifs,
	int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_xtab(slls_t* pxtab_lines, char ips, int allow_repeat_ips);

lrec_t* lrec_parse_mmap_nidx(file_reader_mmap_state_t* phandle, char irs, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_dkvp(file_reader_mmap_state_t *phandle, char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_xtab(file_reader_mmap_state_t* phandle, char irs, char ips, int allow_repeat_ips);

#endif // LREC_READERS_H
