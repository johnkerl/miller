#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "mlr_globals.h"
#include "mlr_arch.h"
#include "mlrutil.h"

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

// ----------------------------------------------------------------
ssize_t mlr_arch_getdelim(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream) {
#ifndef MLR_ON_MSYS2
	return getdelim(pline, plinecap, delimiter, stream);
#else
	return local_getdelim(pline, plinecap, delimiter, stream);
#endif
}
