#include <stdio.h>
#include <stdlib.h>
#include "lemon_assert.h"

void lemon_assert(char *file, int line) {
	fprintf(stderr,"Assertion failed on line %d of file \"%s\"\n", line, file);
	exit(1);
}
