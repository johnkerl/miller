#ifndef READERS_H
#define READERS_H
#include "input/reader.h"

reader_t* reader_csv_alloc(char rs, char fs, int allow_repeat_ifs);
reader_t* reader_dkvp_alloc(char rs, char fs, char ps, int allow_repeat_ifs);
reader_t* reader_nidx_alloc(char rs, char fs, int allow_repeat_ifs);
reader_t* reader_xtab_alloc(char ps, int allow_repeat_ips);

#endif // READERS_H
