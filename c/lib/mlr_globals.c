#include <stdlib.h>
#include "lib/mlr_globals.h"

mlr_globals_t MLR_GLOBALS = { .argv0 = NULL, .ofmt = NULL };
void mlr_global_init(char* argv0, char* ofmt) {
	MLR_GLOBALS.argv0 = argv0;
	MLR_GLOBALS.ofmt  = ofmt;
}
