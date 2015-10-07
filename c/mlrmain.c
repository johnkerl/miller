#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "cli/mlrcli.h"
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "input/lrec_readers.h"
#include "mapping/mappers.h"
#include "output/lrec_writers.h"
#include "stream/stream.h"

int main(int argc, char** argv) {

//	OSX lldb has issues with single-quoted args:
//	if (argc == 2 && streq(argv[1], "--lldb-workaround")) {
//		char* nargv[] = { "mlr-dbg", "put", "$x=$x*2", "../data/small", NULL };
//		argv = nargv;
//		argc = 0;
//		while (argv[argc] != NULL)
//			argc++;
//	}

	mlr_global_init(argv[0], NULL, NULL);
	cli_opts_t* popts = parse_command_line(argc, argv);
	mlr_global_init(argv[0], popts->ofmt, popts);

	lrec_reader_t* plrec_reader = popts->plrec_reader;
	sllv_t*        pmapper_list = popts->pmapper_list;
	lrec_writer_t* plrec_writer = popts->plrec_writer;
	char**         filenames    = popts->filenames;

	int ok = do_stream_chained(filenames, plrec_reader, pmapper_list, plrec_writer, popts->ofmt);

	cli_opts_free(popts);

	return ok ? 0 : 1;
}
