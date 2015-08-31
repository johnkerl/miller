#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/peek_file_reader.h"

#define PEEK_BUF_LEN             32
#define STRING_BUILDER_INIT_SIZE 1024
#define FIXED_LINE_LEN           1024

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
static char* read_line_fgetc(FILE* fp, char* irs, int irs_len) {
	char* line = mlr_malloc_or_die(FIXED_LINE_LEN);
	char* p = line;
	while (TRUE) {
		int c = fgetc(fp);
		if (c == EOF) {
			if (p == line) {
				return NULL;
			} else {
				*(p++) = 0;
				return line;
			}
		} else if (c == irs[0]) {
			*(p++) = 0;
			return line;
		} else {
			*(p++) = c;
		}
	}
}

static int read_file_fgetc(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_fgetc(fp, irs, irs_len);
		if (line == NULL)
			break;
		bc += strlen(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_fgetc_psb(FILE* fp, string_builder_t* psb, char* irs, int irs_len) {
	while (TRUE) {
		int c = fgetc(fp);
		if (c == EOF) {
			if (sb_is_empty(psb))
				return NULL;
			else
				return sb_finish(psb);
		} else if (c == irs[0]) {
			return sb_finish(psb);
		} else {
			sb_append_char(psb, c);
		}
	}
}

static int read_file_fgetc_psb(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_fgetc_psb(fp, psb, irs, irs_len);
		if (line == NULL)
			break;
		bc += strlen(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_mmap_psb(file_reader_mmap_state_t* ph, string_builder_t* psb, char* irs, int irs_len) {
	char *p = ph->sol;
	while (TRUE) {
		if (p == ph->eof) {
			ph->sol = p;
			if (sb_is_empty(psb))
				return NULL;
			else
				return sb_finish(psb);
		} else if (*p == irs[0]) {
			ph->sol = p+1;
			return sb_finish(psb);
		} else {
			sb_append_char(psb, *p);
			p++;
		}
	}
}

static int read_file_mmap_psb(char* filename) {
	file_reader_mmap_state_t* ph = file_reader_mmap_open(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_mmap_psb(ph, psb, irs, irs_len);
		if (line == NULL)
			break;
		bc += strlen(line);
	}
	file_reader_mmap_close(ph);
	return bc;
}

// ================================================================
static char* read_line_pfr_psb(peek_file_reader_t* pfr, string_builder_t* psb, char* irs, int irs_len) {
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

static int read_file_pfr_psb(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	peek_file_reader_t* pfr = pfr_alloc(fp, PEEK_BUF_LEN);
	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_pfr_psb(pfr, psb, irs, irs_len);
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
	int nreps = 1;
	if (argc != 2 && argc != 3)
		usage(argv[0]);
	char* filename = argv[1];
	if (argc == 3)
		(void)sscanf(argv[2], "%d", &nreps);

	double s, e, t;
	int bc;

	for (int i = 0; i < nreps; i++) {
		s = get_systime();
		bc = read_file_mlr_get_line(filename);
		e = get_systime();
		t = e - s;
		printf("type=getdelim,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_fgetc(filename);
		e = get_systime();
		t = e - s;
		printf("type=fgetc,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_fgetc_psb(filename);
		e = get_systime();
		t = e - s;
		printf("type=fgetc_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_mmap_psb(filename);
		e = get_systime();
		t = e - s;
		printf("type=mmap_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_pfr_psb(filename);
		e = get_systime();
		t = e - s;
		printf("type=pfr_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);
	}

	return 0;
}
