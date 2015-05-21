#ifndef LREC_WRITERS_H
#define LREC_WRITERS_H
#include "output/lrec_writer.h"

lrec_writer_t* lrec_writer_csv_alloc(char rs, char fs);
lrec_writer_t* lrec_writer_dkvp_alloc(char rs, char fs, char ps);
lrec_writer_t* lrec_writer_nidx_alloc(char rs, char fs);
lrec_writer_t* lrec_writer_pprint_alloc(int left_align);
lrec_writer_t* lrec_writer_xtab_alloc();

#endif // LREC_WRITERS_H
