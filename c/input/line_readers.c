#include <stdio.h>
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "input/line_readers.h"

// ----------------------------------------------------------------
char* mlr_alloc_read_line_single_delimiter(
	FILE*      fp,
	int        delimiter,
	size_t*    pold_then_new_strlen,
	int        do_auto_line_term,
	context_t* pctx)
{
	size_t linecap = power_of_two_above(*pold_then_new_strlen + 1); // +1 for null-terminator
	char* restrict line = mlr_malloc_or_die(linecap);
	char* restrict p = line;
	int reached_eof = FALSE;
	int c;
	int nread = 0;

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

	if (do_auto_line_term) {
		char* q = p - 1;
		if (q >= line && *q == '\r') {
			*q = 0;
			context_set_autodetected_crlf(pctx);
			p = q;
		} else {
			context_set_autodetected_lf(pctx);
		}
	}

	// linelen excludes line-ending characters.
	// nread   includes line-ending characters.
	int linelen = p - line;
	if (nread == 0 && reached_eof) {
		free(line);
		line = NULL;
		linelen = 0;
	}
	*pold_then_new_strlen = linelen;

	return line;
}

// ----------------------------------------------------------------
char* mlr_alloc_read_line_multiple_delimiter(
	FILE*      fp,
	char*      delimiter,
	int        delimiter_length,
	size_t*    pold_then_new_strlen)
{
	size_t linecap = power_of_two_above(*pold_then_new_strlen + 1); // +1 for null-terminator
	char* line = mlr_malloc_or_die(linecap);
	char* p = line; // points to null-terminator in (chomped) output string
	char* q = line; // points to end of line in (non-chomped) data read from file
	int reached_eof = FALSE;
	int c;
	int nread = 0;
	int dlm1 = delimiter_length - 1;
	char delimend = delimiter[dlm1];

	while (TRUE) {
		size_t offset = q - line;
		if (offset >= linecap-1) {
			linecap = linecap << 1;
			line = mlr_realloc_or_die(line, linecap);
			q = line + offset;
		}
		c = mlr_arch_getc(fp);
		if (c == EOF) {
			*q = 0;
			reached_eof = TRUE;
			p = q;
			break;
		} else if (c == delimend) {
			// For efficiency, do a single-character test to see if we've seen
			// the last character in the line-ending sequence. If we have, then
			// strcmp back to see if we've seen the entire line-ending sequence.
			//
			// This function exists separately from in order to avoid the performance
			// penalty of this strcmp.
			nread++;
			*(q++) = c;
			p = q - delimiter_length;
			if (q - line >= delimiter_length && memcmp(p, delimiter, delimiter_length) == 0) {
				*p = 0;
				break;
			}
		} else {
			nread++;
			*(q++) = c;
		}
	}

	// linelen excludes line-ending characters.
	// nread   includes line-ending characters.
	int linelen = p - line;
	if (nread == 0 && reached_eof) {
		free(line);
		line = NULL;
		linelen = 0;
	}
	*pold_then_new_strlen = linelen;

	return line;
}
