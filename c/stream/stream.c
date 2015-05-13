#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/readers.h"
#include "mapping/mappers.h"
#include "output/writers.h"

static int do_file_chained(char* filename, context_t* pctx,
	reader_t* preader, sllv_t* pmapper_list, writer_t* pwriter, FILE* output_stream);

static sllv_t* chain_map(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head);

static void drive_lrec(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head, writer_t* pwriter, FILE* output_stream);

// ----------------------------------------------------------------
int do_stream_chained(char** filenames, reader_t* preader, sllv_t* pmapper_list, writer_t* pwriter, char* ofmt) {
	FILE* output_stream = stdout;

	context_t ctx = { .nr = 0, .fnr = 0, .filenum = 0, .filename = NULL };
	int ok = 1;
	if (*filenames == NULL) {
		ctx.filenum++;
		ctx.filename = "(stdin)";
		ctx.fnr = 0;
		ok = do_file_chained("-", &ctx, preader, pmapper_list, pwriter, output_stream) && ok;
	} else {
		for (char** pfilename = filenames; *pfilename != NULL; pfilename++) {
			ctx.filenum++;
			ctx.filename = *pfilename;
			ctx.fnr = 0;
			// Start-of-file hook, e.g. expecting CSV headers on input.
			preader->preset_func(preader->pvstate);
		    ok = do_file_chained(*pfilename, &ctx, preader, pmapper_list, pwriter, output_stream) && ok;
		}
	}

	// Mappers and writers receive end-of-stream notifications via null input record.
	// Do that, now that data from all input file(s) have been exhausted.
	drive_lrec(NULL, &ctx, pmapper_list->phead, pwriter, output_stream);

	// Drain the pretty-printer.
	pwriter->pwriter_func(output_stream, NULL, pwriter->pvstate);

	return ok;
}

// ----------------------------------------------------------------
static int do_file_chained(char* filename, context_t* pctx,
	reader_t* preader, sllv_t* pmapper_list, writer_t* pwriter, FILE* output_stream)
{
	FILE* input_stream = stdin;

	// xxx assert pmapper_list non-empty ...

	if (!streq(filename, "-")) {
		input_stream = fopen(filename, "r");
		if (input_stream == NULL) {
			fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n", MLR_GLOBALS.argv0, filename);
			perror(filename);
			return 0;
		}
	}

	while (1) {
		lrec_t* pinrec = preader->preader_func(input_stream, preader->pvstate, pctx);
		if (pinrec == NULL)
			break;
		pctx->nr++;
		pctx->fnr++;
		drive_lrec(pinrec, pctx, pmapper_list->phead, pwriter, output_stream);
	}
	if (input_stream != stdin)
		fclose(input_stream);

	return 1;
}

// ----------------------------------------------------------------
static void drive_lrec(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head, writer_t* pwriter, FILE* output_stream) {
	sllv_t* outrecs = chain_map(pinrec, pctx, pmapper_list_head);
	if (outrecs != NULL) {
		for (sllve_t* pe = outrecs->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* poutrec = pe->pvdata;
			if (poutrec != NULL)
				pwriter->pwriter_func(output_stream, poutrec, pwriter->pvstate);
			// doc & encode convention that writer frees.
		}
		sllv_free(outrecs); // xxx cmt mem-mgmt
	}
}

// ----------------------------------------------------------------
// Map a single input record (maybe null at end of input stream) to zero or more output record.
// Return: list of lrec_t*. Input: lrec_t* and list of mapper_t*.
// xxx need to figure out mem-mgmt here
static sllv_t* chain_map(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head) {
	mapper_t* pmapper = pmapper_list_head->pvdata;
	sllv_t* outrecs = pmapper->pmapper_process_func(pinrec, pctx, pmapper->pvstate);
	if (pmapper_list_head->pnext == NULL) {
		return outrecs;
	} else if (outrecs == NULL) { // xxx cmt
		return NULL;
	} else {
		sllv_t* nextrecs = sllv_alloc();

		for (sllve_t* pe = outrecs->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* poutrec = pe->pvdata;
			sllv_t* nextrecsi = chain_map(poutrec, pctx, pmapper_list_head->pnext);
			nextrecs = sllv_append(nextrecs, nextrecsi);
		}
		sllv_free(outrecs);

		return nextrecs;
	}
}
