#include <stdio.h>
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "input/line_readers.h"

// Use powers of two exclusively, to help avoid heap fragmentation
#define INITIAL_SIZE 128

// ----------------------------------------------------------------
char* mlr_get_cline(FILE* fp, char irs) {
	char* line = NULL;
	size_t linecap = 0;
	ssize_t linelen = mlr_arch_getdelim(&line, &linecap, irs, fp);
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
	ssize_t linelen = mlr_arch_getdelim(&line, &linecap, irs, fp);
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
		c = mlr_arch_getc(fp);
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

// ----------------------------------------------------------------
ssize_t local_getdelim(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream) {
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
			*(p++) = delimiter;
			break;
		} else if (c == EOF) {
			if (p == line)
				eof = TRUE;
			break;
		} else {
			*(p++) = c;
		}
	}

	// xxx check length
	size_t offset = p - line;
	if (offset >= linecap) {
		linecap = linecap + 1;
		line = mlr_realloc_or_die(line, linecap);
		p = line + offset;
	}
	p[1] = 0;

	*pline = line;
	*plinecap = linecap;
	if (eof) {
		**pline = 0;
		return -1;
	} else {
		return p - line;
	}
}

// xxx rm:
// getline:
// in delimiter (single/multiple)
// in fp
// -in do_auto_line_term- separate variant
// -inout pctx- separate variant
// out line
// inout strlen (old/new). DEFAULT_SIZE @ first call
// inout linecap (old/new) DEFAULT_SIZE @ first call
//
// reuse linecap on subsequent calls. power of two above last readlen.
// work autodetect deeper into the callstack

// ----------------------------------------------------------------
char* mlr_alloc_read_line_single_delimiter(
	FILE*   fp,
	int     delimiter,
	size_t* pold_then_new_strlen,
	size_t* pnew_linecap)
{
	size_t linecap = power_of_two_above(*pold_then_new_strlen);
	char* restrict line = mlr_malloc_or_die(linecap);
	char* restrict p = line;
	int reached_eof = FALSE;
	int c;
	int nread = 0;

	// [#|#|#|#]
	// [a|b|c|0]

	while (TRUE) {
		size_t offset = p - line;
		if (offset >= linecap) {
			linecap = linecap << 1;
			line = mlr_realloc_or_die(line, linecap);
			p = line + offset;
		}
		c = mlr_arch_getc(fp);
		if (c == EOF) {
			*p = 0;
			reached_eof = TRUE;
			break;
		} else if (c == delimiter) {
			nread++;
			*p = 0;
			break;
		} else {
			nread++;
			*(p++) = c;
		}
	}

//	// The line-ending '\n' won't be included in the field buffer.
//	if (do_auto_line_term) {
//		if (p > line && *p == '\r') {
//			*(p--) = 0;
//			context_set_autodetected_crlf(pctx);
//		} else {
//			context_set_autodetected_lf(pctx);
//		}
//	}

	int linelen = p - line;

	// linelen excludes line-ending characters.
	// nread   includes line-ending characters.
	if (nread == 0 && reached_eof) {
		line = NULL;
		linelen = 0;
	}
	*pold_then_new_strlen = linelen;
	*pnew_linecap = linecap;

	return line;
}
