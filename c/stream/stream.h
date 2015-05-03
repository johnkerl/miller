#ifndef STREAM_H
#define STREAM_H

#include "containers/sllv.h"
#include "input/readers.h"
#include "mapping/mappers.h"
#include "output/writers.h"

// filenames must be null-terminated
int do_stream(char** filenames, reader_t* preader, mapper_t* pmapper, writer_t* pwriter, char* ofmt);
int do_stream_chained(char* argv0, char** filenames, reader_t* preader, sllv_t* pmapper_list, writer_t* pwriter, char* ofmt);

#endif // STREAM_H
