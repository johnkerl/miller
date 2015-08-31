#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/peek_file_reader.h"

#define PEEK_BUF_LEN             32
#define STRING_BUILDER_INIT_SIZE 1024

// ================================================================
static FILE* fopen_or_die(char* filename) {
	FILE* fp = fopen(filename, "r");
	if (fp == NULL) {
		perror("fopen");
		fprintf(stderr, "Couldn't open \"%s\" for read; exiting.\n", filename);
		exit(1);
	}
	return fp;
}

// ================================================================
static int read_file_mlr_get_line(char* filename) {
	FILE* fp = fopen_or_die(filename);
	int bc = 0;
	while (1) {
		char* line = mlr_get_line(fp, '\n');
		if (line == NULL)
			break;
		bc += strlen(line);
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_pfr_and_psb(peek_file_reader_t* pfr, string_builder_t* psb, char* irs, int irs_len) {
	while (TRUE) {
		if (pfr_at_eof(pfr)) {
			if (sb_is_empty(psb))
				return NULL;
			else
				return sb_finish(psb);
		} else if (pfr_next_is(pfr, irs, irs_len)) {
			if (!pfr_advance_past(pfr, irs)) {
				fprintf(stderr, "%s: Internal coding error: IRS found and lost.\n", MLR_GLOBALS.argv0);
				exit(1);
			}
			return sb_finish(psb);
		} else {
			sb_append_char(psb, pfr_read_char(pfr));
		}
	}
}

static int read_file_pfr_and_psb(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	peek_file_reader_t* pfr = pfr_alloc(fp, PEEK_BUF_LEN);
	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_pfr_and_psb(pfr, psb, irs, irs_len);
		if (line == NULL)
			break;
		bc += strlen(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static void usage(char* argv0) {
	fprintf(stderr, "Usage: %s {filename}\n", argv0);
	exit(1);
}

int main(int argc, char** argv) {
	if (argc != 2)
		usage(argv[0]);
	char* filename = argv[1];

	double s1 = get_systime();
	int bc1 = read_file_mlr_get_line(filename);
	double e1 = get_systime();
	double t1 = e1 - s1;
	printf("type=getdelim,t=%.6lf,n=%d\n", t1, bc1);

	double s2 = get_systime();
	int bc2 = read_file_pfr_and_psb(filename);
	double e2 = get_systime();
	double t2 = e2 - s2;
	printf("type=getdelim,t=%.6lf,n=%d\n", t2, bc2);

	return 0;
}
