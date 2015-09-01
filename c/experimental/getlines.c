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

static int read_file_fgetc_fixed_len(char* filename) {
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
static char* read_line_getc_unlocked(FILE* fp, char* irs, int irs_len) {
	char* line = mlr_malloc_or_die(FIXED_LINE_LEN);
	char* p = line;
	while (TRUE) {
		int i = getc_unlocked(fp);
		char c = i;
		if (i == EOF) {
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

static int read_file_getc_unlocked_fixed_len(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_getc_unlocked(fp, irs, irs_len);
		if (line == NULL)
			break;
		bc += strlen(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_getc_unlocked_psb(FILE* fp, string_builder_t* psb, char* irs, int irs_len) {
	while (TRUE) {
		int c = getc_unlocked(fp);
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

static int read_file_getc_unlocked_psb(char* filename) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";
	int irs_len = strlen(irs);

	int bc = 0;

	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	while (TRUE) {
		char* line = read_line_getc_unlocked_psb(fp, psb, irs, irs_len);
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
		bc = read_file_fgetc_fixed_len(filename);
		e = get_systime();
		t = e - s;
		printf("type=fgetc_fixed_len,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_getc_unlocked_fixed_len(filename);
		e = get_systime();
		t = e - s;
		printf("type=getc_unlocked_fixed_len,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_getc_unlocked_psb(filename);
		e = get_systime();
		t = e - s;
		printf("type=getc_unlocked_psb,t=%.6lf,n=%d\n", t, bc);
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

// $ ./getl ../data/big.csv 5|tee x
//
// $  mlr --opprint sort -nr t then step -a delta -f t x
// type                    t        n        t_delta
// fgetc_fixed_len         3.166140 55888899 3.166140
// fgetc_fixed_len         3.029210 55888899 -0.136930
// fgetc_fixed_len         3.001850 55888899 -0.027360
// fgetc_psb               2.984247 55888899 -0.017603
// fgetc_psb               2.952416 55888899 -0.031831
// fgetc_fixed_len         2.951750 55888899 -0.000666
// fgetc_fixed_len         2.931093 55888899 -0.020657
// fgetc_psb               2.839564 55888899 -0.091529
// fgetc_psb               2.819264 55888899 -0.020300
// fgetc_psb               2.806522 55888899 -0.012742
// getc_unlocked_fixed_len 0.829920 55888899 -1.976602
// pfr_psb                 0.790989 55888900 -0.038931
// pfr_psb                 0.736122 55888900 -0.054867
// pfr_psb                 0.707881 55888900 -0.028241
// pfr_psb                 0.692827 55888900 -0.015054
// pfr_psb                 0.689040 55888900 -0.003787
// getc_unlocked_fixed_len 0.612850 55888899 -0.076190
// getc_unlocked_fixed_len 0.586335 55888899 -0.026515
// getc_unlocked_fixed_len 0.518139 55888899 -0.068196
// getc_unlocked_fixed_len 0.500122 55888899 -0.018017
// getdelim                0.379211 55888899 -0.120911
// getc_unlocked_psb       0.312675 55888899 -0.066536
// getc_unlocked_psb       0.303223 55888899 -0.009452
// getc_unlocked_psb       0.302722 55888899 -0.000501
// getc_unlocked_psb       0.295561 55888899 -0.007161
// mmap_psb                0.291486 55888899 -0.004075
// mmap_psb                0.280870 55888899 -0.010616
// getc_unlocked_psb       0.270824 55888899 -0.010046
// mmap_psb                0.264112 55888899 -0.006712
// mmap_psb                0.256946 55888899 -0.007166
// mmap_psb                0.247112 55888899 -0.009834
// getdelim                0.134491 55888899 -0.112621
// getdelim                0.124980 55888899 -0.009511
// getdelim                0.124031 55888899 -0.000949
// getdelim                0.123838 55888899 -0.000193

// $ mlr --opprint stats1 -a min,mean,max,stddev -f t -g type then sort -n t_mean x
// type                    t_min    t_mean   t_max    t_stddev
// getdelim                0.123838 0.177310 0.379211 0.112953
// mmap_psb                0.247112 0.268105 0.291486 0.017964
// getc_unlocked_psb       0.270824 0.297001 0.312675 0.015846
// getc_unlocked_fixed_len 0.500122 0.609473 0.829920 0.131760
// pfr_psb                 0.689040 0.723372 0.790989 0.042090
// fgetc_psb               2.806522 2.880403 2.984247 0.081905
// fgetc_fixed_len         2.931093 3.016009 3.166140 0.092539


// type                    t_min    t_mean   t_max    t_stddev
// getdelim                0.123838 0.177310 0.379211 0.112953

// mmap_psb                0.247112 0.268105 0.291486 0.017964
// getc_unlocked_psb       0.270824 0.297001 0.312675 0.015846
// getc_unlocked_fixed_len 0.500122 0.609473 0.829920 0.131760
// pfr_psb                 0.689040 0.723372 0.790989 0.042090

// fgetc_psb               2.806522 2.880403 2.984247 0.081905
// fgetc_fixed_len         2.931093 3.016009 3.166140 0.092539
