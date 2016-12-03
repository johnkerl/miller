#include <stdio.h>
#include <stdlib.h>
#include "lemon_memory.h"

// Report an out-of-memory condition and abort.
void memory_error() {
	fprintf(stderr, "Out of memory.  Aborting.\n");
	exit(1);
}
