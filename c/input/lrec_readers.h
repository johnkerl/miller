#ifndef LREC_READERS_H
#define LREC_READERS_H
#include "input/lrec_reader.h"
#include "input/lrec_reader_mmap.h"

lrec_reader_t*  lrec_reader_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_t* lrec_reader_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
lrec_reader_t* lrec_reader_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
lrec_reader_t* lrec_reader_xtab_alloc(char ips, int allow_repeat_ips);

reader_mmap_t*  reader_csv_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_mmap_t* reader_dkvp_mmap_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
reader_mmap_t* reader_nidx_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_mmap_t* reader_xtab_mmap_alloc(char irs, char ips, int allow_repeat_ips);

#endif // LREC_READERS_H
