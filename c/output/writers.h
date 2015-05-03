#ifndef WRITERS_H
#define WRITERS_H
#include "output/writer.h"

writer_t* writer_csv_alloc(char rs, char fs);
writer_t* writer_dkvp_alloc(char rs, char fs, char ps);
writer_t* writer_nidx_alloc(char rs, char fs);
writer_t* writer_pprint_alloc(int left_align);
writer_t* writer_xtab_alloc();

#endif // WRITERS_H
