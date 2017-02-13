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

static int do_stream_chained_in_place(char* prepipe, slls_t* filenames, sllv_t* pmapper_list,
	context_t* pctx, cli_opts_t* popts);

static int do_stream_chained_to_stdout(char* prepipe, slls_t* filenames, sllv_t* pmapper_list,
	context_t* pctx, cli_opts_t* popts);

static int do_file_chained(char* prepipe, char* filename, context_t* pctx,
	lrec_reader_t* plrec_reader, sllv_t* pmapper_list, lrec_writer_t* plrec_writer, FILE* output_stream,
	cli_opts_t* popts);

static sllv_t* chain_map(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head);

static void drive_lrec(lrec_t* pinrec, context_t* pctx, sllve_t* pmapper_list_head, lrec_writer_t* plrec_writer,
	FILE* output_stream);

typedef void progress_indicator_t(context_t* pctx, long long nr_progress_mod);
static void null_progress_indicator(context_t* pctx, long long nr_progress_mod);
static void stderr_progress_indicator(context_t* pctx, long long nr_progress_mod);

// ----------------------------------------------------------------
int do_stream_chained(char* prepipe, slls_t* filenames, sllv_t* pmapper_list, context_t* pctx, cli_opts_t* popts) {
	if (popts->do_in_place) {
		return do_stream_chained_in_place(prepipe, filenames, pmapper_list, pctx, popts);
	} else {
		return do_stream_chained_to_stdout(prepipe, filenames, pmapper_list, pctx, popts);
	}
}

// ----------------------------------------------------------------
static int do_stream_chained_in_place(char* prepipe, slls_t* filenames, sllv_t* pmapper_list,
	context_t* pctx, cli_opts_t* popts)
{
	MLR_INTERNAL_CODING_ERROR_IF(pmapper_list->length < 1); // Should not have been allowed by the CLI parser.
	// xxx make these more clear to the user.
	MLR_INTERNAL_CODING_ERROR_IF(filenames == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(filenames->length == 0);

	int ok = 1;

	// Read from each file name in turn
	for (sllse_t* pe = filenames->phead; pe != NULL; pe = pe->pnext) {

		lrec_writer_t* plrec_writer = lrec_writer_alloc_or_die(&popts->writer_opts);
		lrec_reader_t* plrec_reader = lrec_reader_alloc(&popts->reader_opts);
		if (plrec_reader == NULL) { // xxx funcify
			fprintf(stderr, "%s: unrecognized input-file format \"%s\".\n",
				MLR_GLOBALS.bargv0, popts->reader_opts.ifile_fmt);
			// xxx main_usage_short(stderr, MLR_GLOBALS.bargv0);
			exit(1);
		}

		char* filename = pe->value;

		char* foo = mlr_malloc_or_die(strlen(filename) + 32);
		sprintf(foo, "%s.tmp", filename); // xxx needs uuid

		FILE* output_stream = fopen(foo, "wb"); // xxx stub

		pctx->filenum++;
		pctx->filename = filename;
		pctx->fnr = 0;
		ok = do_file_chained(prepipe, filename, pctx, plrec_reader, pmapper_list, plrec_writer,
			output_stream, popts) && ok;
		if (pctx->force_eof == TRUE) // e.g. mlr head
			continue;

		// Mappers and writers receive end-of-stream notifications via null input record.
		// Do that, now that data from all input file(s) have been exhausted.
		drive_lrec(NULL, pctx, pmapper_list->phead, plrec_writer, output_stream);

		// Drain the pretty-printer.
		plrec_writer->pprocess_func(plrec_writer->pvstate, output_stream, NULL, pctx);

		// xxx needs mapper-reset logic
		// xxx needs writer-reset logic. or move reader/writer alloc/frees inside this file.

		fclose(output_stream);

		int rc = rename(foo, filename);
		if (rc != 0) {
			perror("rename");
			fprintf(stderr, "%s: Could not rename \"%s\" to \"%s\".\n",
				MLR_GLOBALS.bargv0, foo, filename);
			exit(1);
		}

		free(foo);

		plrec_reader->pfree_func(plrec_reader);
		plrec_writer->pfree_func(plrec_writer, pctx);
	}

	return ok;
}

// ----------------------------------------------------------------
static int do_stream_chained_to_stdout(char* prepipe, slls_t* filenames, sllv_t* pmapper_list,
	context_t* pctx, cli_opts_t* popts)
{
	FILE* output_stream = stdout;

	lrec_writer_t* plrec_writer = lrec_writer_alloc_or_die(&popts->writer_opts);
	lrec_reader_t* plrec_reader = lrec_reader_alloc(&popts->reader_opts);
	if (plrec_reader == NULL) { // xxx funcify
		fprintf(stderr, "%s: unrecognized input-file format \"%s\".\n",
			MLR_GLOBALS.bargv0, popts->reader_opts.ifile_fmt);
		// xxx main_usage_short(stderr, MLR_GLOBALS.bargv0);
		exit(1);
	}

	MLR_INTERNAL_CODING_ERROR_IF(pmapper_list->length < 1); // Should not have been allowed by the CLI parser.

	int ok = 1;
	if (filenames == NULL) {
		// No input at all
	} else if (filenames->length == 0) {
		// Zero file names means read from standard input
		pctx->filenum++;
		pctx->filename = "(stdin)";
		pctx->fnr = 0;
		ok = do_file_chained(prepipe, "-", pctx, plrec_reader, pmapper_list, plrec_writer,
			output_stream, popts) && ok;
	} else {
		// Read from each file name in turn
		for (sllse_t* pe = filenames->phead; pe != NULL; pe = pe->pnext) {
			char* filename = pe->value;
			pctx->filenum++;
			pctx->filename = filename;
			pctx->fnr = 0;
			ok = do_file_chained(prepipe, filename, pctx, plrec_reader, pmapper_list, plrec_writer,
				output_stream, popts) && ok;
			if (pctx->force_eof == TRUE) // e.g. mlr head
				break;
		}
	}

	// Mappers and writers receive end-of-stream notifications via null input record.
	// Do that, now that data from all input file(s) have been exhausted.
	drive_lrec(NULL, pctx, pmapper_list->phead, plrec_writer, output_stream);

	// Drain the pretty-printer.
	plrec_writer->pprocess_func(plrec_writer->pvstate, output_stream, NULL, pctx);


	plrec_reader->pfree_func(plrec_reader);
	plrec_writer->pfree_func(plrec_writer, pctx);

	return ok;
}

// ----------------------------------------------------------------
static int do_file_chained(char* prepipe, char* filename, context_t* pctx,
	lrec_reader_t* plrec_reader, sllv_t* pmapper_list, lrec_writer_t* plrec_writer, FILE* output_stream,
	cli_opts_t* popts)
{
	void* pvhandle = plrec_reader->popen_func(plrec_reader->pvstate, prepipe, filename);
	progress_indicator_t* pindicator = popts->nr_progress_mod == 0LL
		? null_progress_indicator
		: stderr_progress_indicator;

	// Start-of-file hook, e.g. expecting CSV headers on input.
	plrec_reader->psof_func(plrec_reader->pvstate, pvhandle);

	while (1) {
		lrec_t* pinrec = plrec_reader->pprocess_func(plrec_reader->pvstate, pvhandle, pctx);
		if (pinrec == NULL)
			break;
		if (pctx->force_eof == TRUE) { // e.g. mlr head
			lrec_free(pinrec);
			break;
		}
		pctx->nr++;
		pctx->fnr++;

		pindicator(pctx, popts->nr_progress_mod);

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
				plrec_writer->pprocess_func(plrec_writer->pvstate, output_stream, poutrec, pctx);
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

// ----------------------------------------------------------------
static void stderr_progress_indicator(context_t* pctx, long long nr_progress_mod) {
	long long remainder = pctx->nr % nr_progress_mod;
	if (remainder == 0) {
		fprintf(stderr, "NR=%lld FNR=%lld FILENAME=%s\n", pctx->nr, pctx->fnr, pctx->filename);
	}
}

static void null_progress_indicator(context_t* pctx, long long nr_progress_mod) {
}
