#include <stdlib.h>
#include <libgen.h>
#include "lib/mlr_globals.h"

mlr_globals_t MLR_GLOBALS = { .bargv0 = "mlr", .ofmt = NULL, .popts = NULL };
void mlr_global_init(char* argv0, char* ofmt, cli_opts_t* popts) {
	MLR_GLOBALS.bargv0 = basename(argv0);
	MLR_GLOBALS.ofmt   = ofmt;
	MLR_GLOBALS.popts  = popts;
}
