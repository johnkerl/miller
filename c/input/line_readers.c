#include <stdio.h>
#include "lib/mlrutil.h"
#include "input/line_readers.h"

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128

// ----------------------------------------------------------------
char* mlr_get_cline(FILE* fp, char irs) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = getdelim(&line, &linecap, irs, fp);
	if (linelen <= 0) {
		if (line != NULL)
			free(line);
		return NULL;
	}
	if (line[linelen-1] == irs) { // chomp
		line[linelen-1] = 0;
		linelen--;
	}
	return line;
}

char* mlr_get_cline_with_length(FILE* fp, char irs, int* plength) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = getdelim(&line, &linecap, irs, fp);
	if (linelen <= 0) {
		if (line != NULL)
			free(line);
		*plength = 0;
		return NULL;
	}
	if (line[linelen-1] == irs) { // chomp
		line[linelen-1] = 0;
		linelen--;
	}
	*plength = linelen;
	return line;
}

// ----------------------------------------------------------------
// Only for direct performance comparison against getdelim()
char* mlr_get_cline2(FILE* fp, char irs) {
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
		c = getc_unlocked(fp);
		if (c == irs) {
			*p = 0;
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
		return NULL;
	} else {
		return line;
	}
}

// ----------------------------------------------------------------
// 0 1 2 3
// a b c 0
char* mlr_get_sline(FILE* fp, char* irs, int irslen) {
	size_t linecap = INITIAL_SIZE;
	char* restrict line = mlr_malloc_or_die(INITIAL_SIZE);
	char* restrict p = line;
	int eof = FALSE;
	int c;
	int irslenm1 = irslen - 1;
	int irslast = irs[irslenm1];

	while (TRUE) {
		size_t offset = p - line;
		if (offset >= linecap) {
			linecap = linecap << 1;
			line = mlr_realloc_or_die(line, linecap);
			p = line + offset;
		}
		c = getc_unlocked(fp);
		if (c == EOF) {
			if (p == line)
				eof = TRUE;
			*p = 0;
			break;
		} else if (c == irslast) {
			// Example: delim="abc". last='c'. Already have read "ab" into line. p-line=2.
			// Now reading 'c'.
			if ((offset >= irslenm1) && streqn(p-irslenm1, irs, irslenm1)) {
				p -= irslenm1;
				*p = 0;
				break;
			} else {
				*(p++) = c;
			}
		} else {
			*(p++) = c;
		}
	}

	if (eof) {
		free(line);
		return NULL;
	} else {
		return line;
	}
}
