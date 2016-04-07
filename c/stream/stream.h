#ifndef STREAM_H
#define STREAM_H

#include "containers/slls.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

int do_stream_chained(char* prepipe, slls_t* filenames, lrec_reader_t* plrec_reader, sllv_t* pmapper_list,
	lrec_writer_t* plrec_writer, char* ofmt, long long nr_progress_mod);

#endif // STREAM_H
