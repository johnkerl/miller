#ifndef LREC_PARSERS_H
#define LREC_PARSERS_H

#include "input/file_reader_mmap.h"
#include "containers/slls.h"
#include "containers/sllv.h"
#include "containers/lrec.h"
#include "containers/hdr_keeper.h"

// xxx cmt/arg re freeing .......
lrec_t* lrec_parse_stdio_nidx(char* line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_dkvp(char* line, char ifs, char ips, int allow_repeat_ifs);
slls_t* split_csv_header_line(char* line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_csv(hdr_keeper_t* phdr_keeper, char* data_line, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_stdio_xtab(slls_t* pxtab_lines, char ips, int allow_repeat_ips);

lrec_t* lrec_parse_mmap_nidx(file_reader_mmap_state_t* phandle, char irs, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_dkvp(file_reader_mmap_state_t *phandle, char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_csv(hdr_keeper_t* phdr_keeper, file_reader_mmap_state_t* phandle, char irs, char ifs, int allow_repeat_ifs);
lrec_t* lrec_parse_mmap_xtab(file_reader_mmap_state_t* phandle, char irs, char ips, int allow_repeat_ips);

#endif // LREC_PARSERS_H
