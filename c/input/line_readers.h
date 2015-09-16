#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>

// xxx cmt mem mgt
char* mlr_get_line(FILE* input_stream, char irs);
char* mlr_get_line_multi_delim(FILE* input_stream, char* irs);

// xxx temp
size_t mlr_getcdelim(char ** restrict ppline, size_t * restrict plinecap, int delimiter, FILE * restrict fp);
size_t mlr_getsdelim(char ** restrict ppline, size_t * restrict plinecap, char* delimiter, int delimlen, FILE * restrict fp);

#endif // LINE_READERS_H
