#include <stdio.h>
#include "lib/mlrutil.h"
#include "input/line_readers.h"

// ----------------------------------------------------------------
char* mlr_get_line(FILE* input_stream, char irs) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = getdelim(&line, &linecap, irs, input_stream);
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
char* mlr_get_line_multi_delim(FILE* input_stream, char* irs) {
	char* line = NULL;
	size_t linecap = 0;
	// xxx move irslen into api to cache the strlen
	ssize_t linelen = mlr_getsdelim(&line, &linecap, irs, strlen(irs), input_stream);
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
// xxx under construction

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128

// xxx look up what restrict do ... should i be using these more often?
// xxx limited semantics: initial linep & linecapp are disregarded; doesn't return the delimiter in the string.
// xxx getcdelim is just for comparison to getdelim. getsdelim is the deliverable.
size_t mlr_getcdelim(char ** restrict ppline, size_t * restrict plinecap, int delimiter, FILE * restrict fp) {
	size_t linecap = INITIAL_SIZE;
	char* restrict pline = mlr_malloc_or_die(INITIAL_SIZE);
	char* restrict p = pline;
	int eof = FALSE;
	int c;

	while (TRUE) {
		if ((p-pline) >= linecap) {
			linecap = linecap << 1;
			pline = realloc(pline, linecap); // xxx mlr_realloc_or_die
			p = pline;
		}
		c = getc_unlocked(fp);
		if (c == delimiter) {
			*p = 0;
			break;
		} else if (c == EOF) {
			if (p == pline)
				eof = TRUE;
			*p = 0;
			break;
		} else {
			*(p++) = c;
		}
	}

	if (eof) {
		free(pline);
		*ppline = NULL;
		return -1;
	} else {
		*ppline = pline;
		*plinecap = linecap;
		return p - pline;
	}
}

size_t mlr_getsdelim(char ** restrict ppline, size_t * restrict plinecap, char* delimiter, int delimlen,
	FILE * restrict fp)
{
	size_t linecap = INITIAL_SIZE;
	char* restrict pline = mlr_malloc_or_die(INITIAL_SIZE);
	char* restrict p = pline;
	int eof = FALSE;
	int c;
	int delimlen1 = delimlen - 1;
	int delimlast = delimiter[delimlen1];

	while (TRUE) {
		if ((p-pline) >= linecap) {
			linecap = linecap << 1;
			pline = realloc(pline, linecap); // xxx mlr_realloc_or_die
			p = pline;
		}
		c = getc_unlocked(fp);
		if (c == delimlast) {
			// Example: delim="abc". last='c'. Already have read "ab" into pline. p-pline=2.
			// Now reading 'c'.
			// xxx make a memeq
			if (((p-pline) >= delimlen1) && !strncmp(p-delimlen1, delimiter, delimlen1)) {
				*p = 0;
				break;
			} else {
				*(p++) = c;
			}
		} else if (c == EOF) {
			if (p == pline)
				eof = TRUE;
			*p = 0;
			break;
		} else {
			*(p++) = c;
		}
	}

	if (eof) {
		free(pline);
		*ppline = NULL;
		return -1;
	} else {
		*ppline = pline;
		*plinecap = linecap;
		return p - pline;
	}
}
