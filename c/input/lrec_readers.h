#ifndef LREC_READERS_H
#define LREC_READERS_H
#include "input/lrec_reader_stdio.h"
#include "input/lrec_reader_mmap.h"

lrec_reader_stdio_t*  lrec_reader_stdio_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_stdio_xtab_alloc(char ips, int allow_repeat_ips);

lrec_reader_mmap_t*  lrec_reader_mmap_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_mmap_xtab_alloc(char irs, char ips, int allow_repeat_ips);

#endif // LREC_READERS_H
