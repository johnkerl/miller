#ifndef MULTI_LREC_WRITER_H
#define MULTI_LREC_WRITER_H

#include <stdio.h>
#include "cli/mlrcli.h"
#include "containers/lhmsv.h"
#include "containers/sllv.h"
#include "output/lrec_writers.h"
#include "output/file_output_mode.h"
#include "lib/context.h"

// ----------------------------------------------------------------
// This is the value struct for the hashmap:
typedef struct _lrec_writer_and_fp_t {
	lrec_writer_t* plrec_writer;
	char* filename_or_command;
	FILE* output_stream;
	int is_popen;
} lrec_writer_and_fp_t;

typedef struct _multi_lrec_writer_t {
	lhmsv_t* pnames_to_lrec_writers_and_fps;
	cli_writer_opts_t* pwriter_opts;
} multi_lrec_writer_t;

// ----------------------------------------------------------------
multi_lrec_writer_t* multi_lrec_writer_alloc(cli_writer_opts_t* pwriter_opts);

void multi_lrec_writer_free(multi_lrec_writer_t* pmlw, context_t* pctx);

void multi_lrec_writer_output_srec(multi_lrec_writer_t* pmlw, lrec_t* poutrec, char* filename_or_command,
	file_output_mode_t file_output_mode, int flush_every_record, context_t* pctx);

void multi_lrec_writer_output_list(multi_lrec_writer_t* pmlw, sllv_t* poutrecs, char* filename_or_command,
	file_output_mode_t file_output_mode, int flush_every_record, context_t* pctx);

void multi_lrec_writer_drain(multi_lrec_writer_t* pmlw, context_t* pctx);

#endif // MULTI_LREC_WRITER_H
