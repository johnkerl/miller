#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "file_reader_stdio.h"

// ----------------------------------------------------------------
void* file_reader_stdio_vopen(void* pvstate, char* prepipe, char* filename) {
	FILE* input_stream = stdin;

	if (!streq(filename, "-")) {
		input_stream = fopen(filename, "r");
		if (input_stream == NULL) {
			fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n", MLR_GLOBALS.argv0, filename);
			perror(filename);
			exit(1);
		}
	}
	return input_stream;
}

// ----------------------------------------------------------------
void file_reader_stdio_vclose(void* pvstate, void* pvhandle) {
	FILE* input_stream = pvhandle;
	if (input_stream != stdin)
		fclose(input_stream);
}
