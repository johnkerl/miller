#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>
#include "lib/context.h"

// Notes:
// * The caller should free the return value.
// * The line-terminator is not returned as part of the string.
// * Null is returned at EOF.

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

#endif // LINE_READERS_H
