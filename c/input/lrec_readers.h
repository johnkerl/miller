#ifndef LREC_READERS_H
#define LREC_READERS_H
#include "input/lrec_reader_stdio.h"
#include "input/lrec_reader_mmap.h"

lrec_reader_stdio_t*  lrec_reader_csv_stdio_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_dkvp_stdio_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_nidx_stdio_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_stdio_t* lrec_reader_xtab_stdio_alloc(char ips, int allow_repeat_ips);

lrec_reader_mmap_t*  lrec_reader_csv_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_dkvp_mmap_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_nidx_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_mmap_t* lrec_reader_xtab_mmap_alloc(char irs, char ips, int allow_repeat_ips);

#endif // LREC_READERS_H
