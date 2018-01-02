#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>
#include "cli/comment_handling.h"
#include "lib/context.h"

// Notes:
// * The caller should free the return value.
// * The line-terminator is not returned as part of the string.
// * Null is returned at EOF.
// * Simiar to getdelim but customized for Miller: in particular, support for autodetected line endings (LF/CRLF).
//   Also, exists on Windows MSYS2 where there isn't a getdelim.
// * Line-length reuses previous length for initial buffer-size allocation. Pass MLR_ALLOC_READ_LINE_INITIAL_SIZE
//   on first call. On subsequent calls, buffer-size allocations will adapt to the file's line-lengths.

#define MLR_ALLOC_READ_LINE_INITIAL_SIZE 128

char* mlr_alloc_read_line_single_delimiter(
	FILE*      fp,
	int        delimiter,
	size_t*    pold_then_new_strlen,
	int        do_auto_line_term,
	context_t* pctx);

char* mlr_alloc_read_line_multiple_delimiter(
	FILE*      fp,
	char*      delimiter,
	int        delimiter_length,
	size_t*    pold_then_new_strlen);

char* mlr_alloc_read_line_single_delimiter_stripping_comments(
	FILE*              fp,
	int                delimiter,
	size_t*            pold_then_new_strlen,
	int                do_auto_line_term,
	comment_handling_t comment_handling,
	char*              comment_string,
	context_t*         pctx);

char* mlr_alloc_read_line_multiple_delimiter_stripping_comments(
	FILE*              fp,
	char*              delimiter,
	int                delimiter_length,
	size_t*            pold_then_new_strlen,
	comment_handling_t comment_handling,
	char*              comment_string);

#endif // LINE_READERS_H
