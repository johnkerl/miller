#ifndef LREC_WRITERS_H
#define LREC_WRITERS_H
#include "output/lrec_writer.h"

lrec_writer_t* lrec_writer_csv_alloc(char* ors, char* ofs, int oquoting);
lrec_writer_t* lrec_writer_csvlite_alloc(char* ors, char* ofs);
lrec_writer_t* lrec_writer_dkvp_alloc(char* ors, char* ofs, char* ops);
lrec_writer_t* lrec_writer_nidx_alloc(char* ors, char* ofs);
lrec_writer_t* lrec_writer_pprint_alloc(char* ors, char*ofs, int left_align);
lrec_writer_t* lrec_writer_xtab_alloc(char* ors, char* ofs);

#endif // LREC_WRITERS_H
