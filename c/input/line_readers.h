#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>

// Notes:
// * The caller should free the return value.
// * The line-terminator is not returned as part of the string.
// * Null is returned at EOF.

// Get a line terminated by a single character, e.g. '\n' (LF).
char*  mlr_get_cline(FILE* input_stream, char irs);

// Only for performance comparison
char* mlr_get_cline2(FILE* input_stream, char irs);

// Get a line terminated by multiple characters, e.g. '\r\n' (CRLF).
// The irslen is simply to cache the result of what would otherwise be a
// redundant call to strlen() on every invocation.
char*  mlr_get_sline(FILE* input_stream, char* irs, int irslen);

#endif // LINE_READERS_H
