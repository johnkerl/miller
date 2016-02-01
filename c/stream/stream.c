#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"

static int do_file_chained(char* prepipe, char* filename, context_t* pctx,
	lrec_reader_t* plrec_reader, sllv_t* pmapper_list, lrec_writer_t* plrec_writer, FILE* output_stream);

static sllv_t* chain_map(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head);

static void drive_lrec(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head, lrec_writer_t* plrec_writer,
	FILE* output_stream);

// ----------------------------------------------------------------
int do_stream_chained(char* prepipe, char** filenames, lrec_reader_t* plrec_reader, sllv_t* pmapper_list,
	lrec_writer_t* plrec_writer, char* ofmt)
{
	FILE* output_stream = stdout;

	if (pmapper_list->length < 1) { // Should not have been allowed by the CLI parser.
		fprintf(stderr, "%s: internal coding error detected at file %s line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	context_t ctx = { .nr = 0, .fnr = 0, .filenum = 0, .filename = NULL };
	int ok = 1;
	if (*filenames == NULL) {
		ctx.filenum++;
		ctx.filename = "(stdin)";
		ctx.fnr = 0;
		ok = do_file_chained(prepipe, "-", &ctx, plrec_reader, pmapper_list, plrec_writer, output_stream) && ok;
	} else {
		for (char** pfilename = filenames; *pfilename != NULL; pfilename++) {
			ctx.filenum++;
			ctx.filename = *pfilename;
			ctx.fnr = 0;
			ok = do_file_chained(prepipe, *pfilename, &ctx, plrec_reader, pmapper_list,
				plrec_writer, output_stream) && ok;
		}
	}

	// Mappers and writers receive end-of-stream notifications via null input record.
	// Do that, now that data from all input file(s) have been exhausted.
	drive_lrec(NULL, &ctx, pmapper_list->phead, plrec_writer, output_stream);

	// Drain the pretty-printer.
	plrec_writer->pprocess_func(output_stream, NULL, plrec_writer->pvstate);

	return ok;
}

// ----------------------------------------------------------------
static int do_file_chained(char* prepipe, char* filename, context_t* pctx,
	lrec_reader_t* plrec_reader, sllv_t* pmapper_list, lrec_writer_t* plrec_writer, FILE* output_stream)
{
	void* pvhandle = plrec_reader->popen_func(plrec_reader->pvstate, prepipe, filename);

	// Start-of-file hook, e.g. expecting CSV headers on input.
	plrec_reader->psof_func(plrec_reader->pvstate);

	while (1) {
		lrec_t* pinrec = plrec_reader->pprocess_func(plrec_reader->pvstate, pvhandle, pctx);
		if (pinrec == NULL)
			break;
		pctx->nr++;
		pctx->fnr++;
		drive_lrec(pinrec, pctx, pmapper_list->phead, plrec_writer, output_stream);
	}

	plrec_reader->pclose_func(plrec_reader->pvstate, pvhandle, prepipe);
	return 1;
}

// ----------------------------------------------------------------
static void drive_lrec(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head, lrec_writer_t* plrec_writer,
	FILE* output_stream)
{
	sllv_t* outrecs = chain_map(pinrec, pctx, pmapper_list_head);
	if (outrecs != NULL) {
		for (sllve_t* pe = outrecs->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* poutrec = pe->pvvalue;
			if (poutrec != NULL) // writer frees records (sllv void-star payload)
				plrec_writer->pprocess_func(output_stream, poutrec, plrec_writer->pvstate);
		}
		sllv_free(outrecs); // we free the list
	}
}

// ----------------------------------------------------------------
// Map a single input record (maybe null at end of input stream) to zero or
// more output records.
//
// Return: list of lrec_t*. Input: lrec_t* and list of mapper_t*.

static sllv_t* chain_map(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head) {
	mapper_t* pmapper = pmapper_list_head->pvvalue;
	sllv_t* outrecs = pmapper->pprocess_func(pinrec, pctx, pmapper->pvstate);
	if (pmapper_list_head->pnext == NULL) {
		return outrecs;
	} else if (outrecs == NULL) { // end of input stream
		return NULL;
	} else {
		sllv_t* nextrecs = sllv_alloc();

		for (sllve_t* pe = outrecs->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* poutrec = pe->pvvalue;
			sllv_t* nextrecsi = chain_map(poutrec, pctx, pmapper_list_head->pnext);
			sllv_transfer(nextrecs, nextrecsi);
			sllv_free(nextrecsi);
		}
		sllv_free(outrecs);

		return nextrecs;
	}
}
