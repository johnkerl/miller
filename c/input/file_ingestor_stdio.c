#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <sys/stat.h>
#include <sys/mman.h>
#include "lib/mlrutil.h"
#include "lib/mlrescape.h"
#include "lib/mlr_globals.h"
#include "file_ingestor_stdio.h"

// ----------------------------------------------------------------
void* file_ingestor_stdio_vopen(void* pvstate, char* prepipe, char* filename) {
	char* file_contents_buffer = NULL;
	size_t file_size = 0;

	if (prepipe == NULL) {
		if (streq(filename, "-")) {
			file_contents_buffer = read_fp_into_memory(stdin, &file_size);
			if (file_contents_buffer == NULL) {
				fprintf(stderr, "%s: Couldn't open standard input for read.\n", MLR_GLOBALS.bargv0);
				exit(1);
			}
		} else {
			file_contents_buffer = read_file_into_memory(filename, &file_size);
			if (file_contents_buffer == NULL) {
				fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n", MLR_GLOBALS.bargv0, filename);
				exit(1);
			}
		}

	} else {
		char* escaped_filename = alloc_file_name_escaped_for_popen(filename);
		char* command = mlr_malloc_or_die(strlen(prepipe) + 3 + strlen(escaped_filename) + 1);
		if (streq(filename, "-"))
			sprintf(command, "%s", prepipe);
		else
			sprintf(command, "%s < %s", prepipe, escaped_filename);
		FILE* input_stream = popen(command, "r");
		if (input_stream == NULL) {
			fprintf(stderr, "%s: Couldn't popen \"%s\" for read.\n", MLR_GLOBALS.bargv0, command);
			perror(command);
			exit(1);
		}
		file_contents_buffer = read_fp_into_memory(input_stream, &file_size);
		if (file_contents_buffer == NULL) {
			fprintf(stderr, "%s: Couldn't open popen output for read: \"%s\".\n", MLR_GLOBALS.bargv0, command);
			exit(1);
		}
		pclose(input_stream);
		free(escaped_filename);
		free(command);
	}
	file_ingestor_stdio_state_t* pstate = mlr_malloc_or_die(sizeof(file_ingestor_stdio_state_t));
	pstate->sof = file_contents_buffer;
	pstate->eof = &file_contents_buffer[file_size];
	return pstate;
}

// ----------------------------------------------------------------
void file_ingestor_stdio_nop_vclose(void* pvstate, void* pvhandle, char* prepipe) {
}

// ----------------------------------------------------------------
void file_ingestor_stdio_vclose(void* pvstate, void* pvhandle, char* prepipe) {
	char* file_contents_buffer = pvhandle;
	free(file_contents_buffer);
	file_ingestor_stdio_state_t* pstate = pvstate;
	free(pstate);
}
