#ifndef STREAM_H
#define STREAM_H

#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

// filenames must be null-terminated
int do_stream_chained(char** filenames, int use_mmap_reader, lrec_reader_t* plrec_reader, lrec_reader_mmap_t* plrec_reader_mmap,
	sllv_t* pmapper_list, lrec_writer_t* plrec_writer, char* ofmt);

#endif // STREAM_H
