#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>

// xxx cmt mem mgt
// xxx cmt semantics w/ not returning the terminator, and null at eof
// xxx maybe return the line-length by reference? it's available in the function bodies.

char*  mlr_get_cline(FILE* input_stream, char irs);

// Only for performance comparison
char* mlr_get_cline2(FILE* input_stream, char irs);

// The irslen is simply to cache the result of what would otherwise be a
// redundant call to strlen() on every invocation.
char*  mlr_get_sline(FILE* input_stream, char* irs, int irslen);

#endif // LINE_READERS_H
