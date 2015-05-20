#ifndef READERS_H
#define READERS_H
#include "input/reader.h"
#include "input/reader_mmap.h"

reader_t*  reader_csv_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_t* reader_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
reader_t* reader_nidx_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_t* reader_xtab_alloc(char ips, int allow_repeat_ips);

reader_mmap_t*  reader_csv_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_mmap_t* reader_dkvp_mmap_alloc(char irs, char ifs, char ips, int allow_repeat_ifs);
reader_mmap_t* reader_nidx_mmap_alloc(char irs, char ifs, int allow_repeat_ifs);
reader_mmap_t* reader_xtab_mmap_alloc(char irs, char ips, int allow_repeat_ips);

#endif // READERS_H
