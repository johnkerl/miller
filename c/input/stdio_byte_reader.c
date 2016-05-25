#include <stdio.h>
#include <string.h>
#include "input/byte_readers.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"

typedef struct _stdio_byte_reader_state_t {
	char* filename;
	FILE* fp;
} stdio_byte_reader_state_t;

static int stdio_byte_reader_open_func(struct _byte_reader_t* pbr, char* prepipe, char* filename);
static int stdio_byte_reader_read_func(struct _byte_reader_t* pbr);
static void stdio_byte_reader_close_func(struct _byte_reader_t* pbr, char* prepipe);

// ----------------------------------------------------------------
byte_reader_t* stdio_byte_reader_alloc() {
	byte_reader_t* pbr = mlr_malloc_or_die(sizeof(byte_reader_t));

	pbr->pvstate     = NULL;
	pbr->popen_func  = stdio_byte_reader_open_func;
	pbr->pread_func  = stdio_byte_reader_read_func;
	pbr->pclose_func = stdio_byte_reader_close_func;

	return pbr;
}

void stdio_byte_reader_free(byte_reader_t* pbr) {
	stdio_byte_reader_state_t* pstate = pbr->pvstate;
	if (pstate != NULL) {
		free(pstate->filename); // null-ok semantics
	}
	free(pbr);
}

// ----------------------------------------------------------------
static int stdio_byte_reader_open_func(struct _byte_reader_t* pbr, char* prepipe, char* filename) {
	stdio_byte_reader_state_t* pstate = mlr_malloc_or_die(sizeof(stdio_byte_reader_state_t));

	pstate->filename = mlr_strdup_or_die(filename);

	if (prepipe == NULL) {
		if (streq(pstate->filename, "-")) {
			pstate->fp = stdin;
		} else {
			pstate->fp = fopen(filename, "r");
			if (pstate->fp == NULL) {
				perror("fopen");
				fprintf(stderr, "%s: Couldn't fopen \"%s\" for read.\n", MLR_GLOBALS.bargv0, filename);
				exit(1);
			}
		}
	} else {
		char* command = mlr_malloc_or_die(strlen(prepipe) + 3 + strlen(filename) + 1);
		if (streq(filename, "-"))
			sprintf(command, "%s", prepipe);
		else
			sprintf(command, "%s < %s", prepipe, filename);
		pstate->fp = popen(command, "r");
		if (pstate->fp == NULL) {
			fprintf(stderr, "%s: Couldn't popen \"%s\" for read.\n", MLR_GLOBALS.bargv0, command);
			perror(command);
			exit(1);
		}
		free(command);
	}

	pbr->pvstate = pstate;
	return TRUE;
}

static int stdio_byte_reader_read_func(struct _byte_reader_t* pbr) {
	stdio_byte_reader_state_t* pstate = pbr->pvstate;
	int c = getc_unlocked(pstate->fp);
	if (c == EOF && ferror(pstate->fp)) {
		perror("fread");
		fprintf(stderr, "%s: Read error on file \"%s\".\n", MLR_GLOBALS.bargv0, pstate->filename);
		exit(1);
	}
	return c;
}

static void stdio_byte_reader_close_func(struct _byte_reader_t* pbr, char* prepipe) {
	stdio_byte_reader_state_t* pstate = pbr->pvstate;
	if (prepipe == NULL) {
		if (pstate->fp != stdin)
			fclose(pstate->fp);
	} else {
		pclose(pstate->fp);
	}
	free(pstate->filename);
	free(pstate);
	pbr->pvstate = NULL;
}
