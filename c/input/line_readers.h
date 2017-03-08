#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>

// Notes:
// * The caller should free the return value.
// * The line-terminator is not returned as part of the string.
// * Null is returned at EOF.

// Get a line terminated by a single character, e.g. '\n' (LF).
char* mlr_get_cline(FILE* input_stream, char irs);
// *plength is an output reference argument which, after return, contains the strlen
// of the return value (i.e. not counting the null-terminator character).
char* mlr_get_cline_with_length(FILE* input_stream, char irs, int* plength);

// Get a line terminated by multiple characters, e.g. '\r\n' (CRLF).
// The irslen is simply to cache the result of what would otherwise be a
// redundant call to strlen() on every invocation.
char*  mlr_get_sline(FILE* input_stream, char* irs, int irslen);

// getdelim is built-in on OSX and modern unix-like OSs. For MSYS2, we need to
// roll our own. The function is exposed publicly here, rather than privately
// inside mlr_arch.c, for unit-testing visibility.
ssize_t local_getdelim(char** restrict pline, size_t* restrict plinecap, int delimiter, FILE* restrict stream);

// xxx type up comments:
// * in delimiter (single/multiple)
// * in fp
// * -in do_auto_line_term- separate variant
// * -inout pctx- separate variant
// * out line
// * out reached eof
// * inout strlen (old/new). DEFAULT_SIZE @ first call
// * inout linecap (old/new) DEFAULT_SIZE @ first call
//
// reuse linecap on subsequent calls. power of two above last readlen.
// work autodetect deeper into the callstack

#define MLR_ALLOC_READ_LINE_INITIAL_SIZE 128
char* mlr_alloc_read_line_single_delimiter(
	FILE*   fp,
	int     delimiter,
	size_t* pold_then_new_strlen,
	size_t* pnew_linecap);

#endif // LINE_READERS_H
