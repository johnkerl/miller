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
char * mlr_arch_strsep(char **pstring, const char *delim) {
#ifdef MLR_ON_MSYS2
	return strtok_r(*pstring, delim, pstring);
#else
	return strsep(pstring, delim);
#endif
}

// ----------------------------------------------------------------
#ifdef MLR_ON_MSYS2

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128
static int local_getdelim(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream) {
	size_t linecap = INITIAL_SIZE;
	char* restrict line = mlr_malloc_or_die(INITIAL_SIZE);
	char* restrict p = line;
	int eof = FALSE;
	int c;

	while (TRUE) {
		size_t offset = p - line;
		if (offset >= linecap) {
			linecap = linecap << 1;
			line = mlr_realloc_or_die(line, linecap);
			p = line + offset;
		}
		c = mlr_arch_getc(stream);
		if (c == delimiter) {
			*p = 0;
			p++;
			break;
		} else if (c == EOF) {
			if (p == line)
				eof = TRUE;
			*p = 0;
			break;
		} else {
			*(p++) = c;
		}
	}

	if (eof) {
		free(line);
		*pline = NULL;
		*plinecap = 0;
		return -1;
	} else {
		*pline = line;
		*plinecap = linecap;
		return p - line;
	}
}
#endif

// ----------------------------------------------------------------
ssize_t mlr_arch_getdelim(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream) {
#ifndef MLR_ON_MSYS2
	return getdelim(pline, plinecap, delimiter, stream);
#else
	return local_getdelim(pline, plinecap, delimiter, stream);
#endif
}
