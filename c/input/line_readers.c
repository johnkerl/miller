#include <stdio.h>
#include "lib/mlrutil.h"
#include "input/line_readers.h"

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128

// xxx should i be using restrict more often?

// ----------------------------------------------------------------
char* mlr_get_cline(FILE* fp, char irs) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = getdelim(&line, &linecap, irs, fp);
	if (linelen <= 0) {
		return NULL;
	}
	if (line[linelen-1] == '\n') { // chomp
		line[linelen-1] = 0;
		linelen--;
	}

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
		if ((p-line) >= linecap) {
			linecap = linecap << 1;
			line = realloc(line, linecap); // xxx mlr_realloc_or_die
			p = line;
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
		if ((p-line) >= linecap) {
			linecap = linecap << 1;
			line = realloc(line, linecap); // xxx mlr_realloc_or_die
			p = line;
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
			// xxx make a memneq
			if (((p-line) >= irslenm1) && !strncmp(p-irslenm1, irs, irslenm1)) {
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
