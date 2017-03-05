#include <stdio.h>
#include <stdlib.h>
#include "mlr_globals.h"
#include "mlr_arch.h"

// ----------------------------------------------------------------
int mlr_arch_setenv(const char *name, const char *value) {
#ifdef MLR_ON_MSYS2
	fprintf(stderr, "%s: setenv is not supported on this architecture.\n", MLR_GLOBALS.bargv0);
	exit(1);
#else
	return setenv(name, value, 1 /*overwrite*/);
#endif
}

// ----------------------------------------------------------------
int mlr_arch_unsetenv(const char *name) {
#ifdef MLR_ON_MSYS2
	fprintf(stderr, "%s: unsetenv is not supported on this architecture.\n", MLR_GLOBALS.bargv0);
	exit(1);
#else
	return unsetenv(name);
#endif
}
