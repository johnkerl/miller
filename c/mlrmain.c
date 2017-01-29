#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "cli/mlrcli.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"
#include "stream/stream.h"

int main(int argc, char** argv) {

	mlr_global_init(argv[0], NULL);
	cli_opts_t* popts = parse_command_line(argc, argv);
	mlr_global_init(argv[0], popts->ofmt);

	char*          prepipe      = popts->reader_opts.prepipe;
	lrec_reader_t* plrec_reader = popts->plrec_reader;
	sllv_t*        pmapper_list = popts->pmapper_list;
	lrec_writer_t* plrec_writer = popts->plrec_writer;
	slls_t*        filenames    = popts->filenames;

	context_t ctx = {
		.nr        = 0,
		.fnr       = 0,
		.filenum   = 0,
		.filename  = NULL,
		.force_eof = FALSE,

		.ips       = popts->reader_opts.ips,
		.ifs       = popts->reader_opts.ifs,
		.irs       = popts->reader_opts.irs,
		.ops       = popts->writer_opts.ops,
		.ofs       = popts->writer_opts.ofs,
		.ors       = popts->writer_opts.ors,
		.auto_line_term  = popts->reader_opts.irs,
	};

	int ok = do_stream_chained(prepipe, filenames, plrec_reader, pmapper_list, plrec_writer, &ctx, popts);

	cli_opts_free(popts, &ctx);

	return ok ? 0 : 1;
}
