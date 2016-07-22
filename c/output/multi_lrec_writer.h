#ifndef MULTI_LREC_WRITER_H
#define MULTI_LREC_WRITER_H

#include <stdio.h>
#include "containers/lhmsv.h"
#include "containers/sllv.h"
#include "output/lrec_writers.h"

// ----------------------------------------------------------------
typedef struct _lrec_writer_and_fp_t {
	lrec_writer_t* plrec_writer;
	FILE* output_stream;
} lrec_writer_and_fp_t;

typedef struct _multi_lrec_writer_t {
	// xxx to do: bound the number of open files and LRU them.
	lhmsv_t* pnames_to_lrec_writers_and_fps;
} multi_lrec_writer_t;

// ----------------------------------------------------------------
multi_lrec_writer_t* multi_lrec_writer_alloc();
void multi_lrec_writer_free(multi_lrec_writer_t* pmlw);

void multi_lrec_writer_write(multi_lrec_writer_t* pmlw, sllv_t* poutrecs, char* filename, int flush_every_record);
void multi_lrec_writer_append(multi_lrec_writer_t* pmlw, sllv_t* poutrecs, char* filename, int flush_every_record);
void multi_lrec_writer_drain(multi_lrec_writer_t* pmlw);

#endif // MULTI_LREC_WRITER_H
