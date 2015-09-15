#ifndef LINE_READERS_H
#define LINE_READERS_H

#include <stdio.h>

size_t mlr_getcdelim(char ** restrict ppline, size_t * restrict plinecap, int delimiter, FILE * restrict fp);
size_t mlr_getsdelim(char ** restrict ppline, size_t * restrict plinecap, char* delimiter, int delimlen, FILE * restrict fp);

#endif // LINE_READERS_H
