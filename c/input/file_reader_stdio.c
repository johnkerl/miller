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

	if (prepipe == NULL) {
		if (!streq(filename, "-")) {
			input_stream = fopen(filename, "r");
			if (input_stream == NULL) {
				fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n", MLR_GLOBALS.bargv0, filename);
				perror(filename);
				exit(1);
			}
		}
	} else {
		char* command = mlr_malloc_or_die(strlen(prepipe) + 5 + strlen(filename) + 1);
		if (streq(filename, "-"))
			sprintf(command, "%s", prepipe);
		else
			sprintf(command, "%s < '%s'", prepipe, filename);
		input_stream = popen(command, "r");
		if (input_stream == NULL) {
			fprintf(stderr, "%s: Couldn't popen \"%s\" for read.\n", MLR_GLOBALS.bargv0, command);
			perror(command);
			exit(1);
		}
		free(command);
	}
	return input_stream;
}

// ----------------------------------------------------------------
void file_reader_stdio_vclose(void* pvstate, void* pvhandle, char* prepipe) {
	FILE* input_stream = pvhandle;
	if (prepipe == NULL) {
		if (input_stream != stdin)
			fclose(input_stream);
	} else {
		pclose(input_stream);
	}
}
