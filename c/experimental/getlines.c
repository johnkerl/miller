#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"
#include "lib/mlrdatetime.h"
#include "input/file_reader_stdio.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/byte_readers.h"
#include "input/line_readers.h"
#include "input/peek_file_reader.h"
#include "containers/parse_trie.h"

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
static int read_file_mlr_alloc_read_line_single_delimiter_no_autodetect(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	int bc = 0;
	size_t line_length = MLR_ALLOC_READ_LINE_INITIAL_SIZE;
	while (1) {
		char* line = mlr_alloc_read_line_single_delimiter(
			fp, '\n', &line_length, FALSE, NULL);
		if (line == NULL)
			break;
		bc += strlen(line);
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static int read_file_mlr_alloc_read_line_single_delimiter_with_autodetect(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	int bc = 0;
	size_t line_length = MLR_ALLOC_READ_LINE_INITIAL_SIZE;
	context_t ctx;
	context_init_from_first_file_name(&ctx, "fake-file-name");

	while (1) {
		char* line = mlr_alloc_read_line_single_delimiter(
			fp, '\n', &line_length, TRUE, &ctx);
		if (line == NULL)
			break;
		bc += strlen(line);
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static int read_file_mlr_alloc_read_line_multiple_delimiter(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	int bc = 0;
	size_t line_length = MLR_ALLOC_READ_LINE_INITIAL_SIZE;

	while (1) {
		char* line = mlr_alloc_read_line_multiple_delimiter(fp, "\n", 1, &line_length);
		if (line == NULL)
			break;
		bc += strlen(line);
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_fgetc(FILE* fp, char* irs) {
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

static int read_file_fgetc_fixed_len(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";

	int bc = 0;

	while (TRUE) {
		char* line = read_line_fgetc(fp, irs);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_getc_unlocked(FILE* fp, char* irs) {
	char* line = mlr_malloc_or_die(FIXED_LINE_LEN);
	char* p = line;
	while (TRUE) {
		int i = mlr_arch_getc(fp);
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

static int read_file_getc_unlocked_fixed_len(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	char* irs = "\n";

	int bc = 0;

	while (TRUE) {
		char* line = read_line_getc_unlocked(fp, irs);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
		free(line);
	}
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_getc_unlocked_psb(FILE* fp, string_builder_t* psb, char* irs) {
	while (TRUE) {
		int c = mlr_arch_getc(fp);
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

static int read_file_getc_unlocked_psb(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	string_builder_t* psb = sb_alloc(STRING_BUILDER_INIT_SIZE);
	char* irs = "\n";
	int bc = 0;

	while (TRUE) {
		char* line = read_line_getc_unlocked_psb(fp, psb, irs);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
		free(line);
	}
	sb_free(psb);
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_fgetc_psb(FILE* fp, string_builder_t* psb, char* irs) {
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

static int read_file_fgetc_psb(char* filename, int do_write) {
	FILE* fp = fopen_or_die(filename);
	string_builder_t* psb = sb_alloc(STRING_BUILDER_INIT_SIZE);
	char* irs = "\n";
	int bc = 0;

	while (TRUE) {
		char* line = read_line_fgetc_psb(fp, psb, irs);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
		free(line);
	}
	sb_free(psb);
	fclose(fp);
	return bc;
}

// ================================================================
static char* read_line_mmap_psb(file_reader_mmap_state_t* ph, string_builder_t* psb, char* irs) {
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

static int read_file_mmap_psb(char* filename, int do_write) {
	file_reader_mmap_state_t* ph = file_reader_mmap_open(NULL, filename);
	string_builder_t* psb = sb_alloc(STRING_BUILDER_INIT_SIZE);
	char* irs = "\n";
	int bc = 0;

	while (TRUE) {
		char* line = read_line_mmap_psb(ph, psb, irs);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
	}
	sb_free(psb);
	file_reader_mmap_close(ph, NULL);
	return bc;
}

// ================================================================
#define IRS_TOKEN    11
#define EOF_TOKEN    22
#define IRSEOF_TOKEN 33

static char* read_line_pfr_psb(peek_file_reader_t* pfr, string_builder_t* psb, parse_trie_t* ptrie) {
	int rc, stridx, matchlen;
	while (TRUE) {
		pfr_buffer_by(pfr, ptrie->maxlen);
		rc = parse_trie_ring_match(ptrie, pfr->peekbuf, pfr->sob, pfr->npeeked, pfr->peekbuflenmask,
			&stridx, &matchlen);
		if (rc) {
			pfr_advance_by(pfr, matchlen);
			switch(stridx) {
			case IRS_TOKEN:
				return sb_finish(psb);
				break;
			case IRSEOF_TOKEN:
				return sb_finish(psb);
				break;
			case EOF_TOKEN:
				return NULL;
				break;
			}
		} else {
			sb_append_char(psb, pfr_read_char(pfr));
		}
	}
}

static int read_file_pfr_psb(char* filename, int do_write) {
	byte_reader_t* pbr = stdio_byte_reader_alloc();
	string_builder_t* psb = sb_alloc(STRING_BUILDER_INIT_SIZE);
	pbr->popen_func(pbr, NULL, filename);

	peek_file_reader_t* pfr = pfr_alloc(pbr, PEEK_BUF_LEN);

	parse_trie_t* ptrie = parse_trie_alloc();
	parse_trie_add_string(ptrie, "\n", IRS_TOKEN);
	parse_trie_add_string(ptrie, "\xff", EOF_TOKEN);
	parse_trie_add_string(ptrie, "\n\xff", IRSEOF_TOKEN);

	int bc = 0;

	while (TRUE) {
		char* line = read_line_pfr_psb(pfr, psb, ptrie);
		if (line == NULL)
			break;
		if (do_write) {
			fputs(line, stdout);
			fputc('\n', stdout);
		}
		bc += strlen(line);
		free(line);
	}
	sb_free(psb);
	pbr->pclose_func(pbr, NULL);
	return bc;
}

// ================================================================
static void usage(char* argv0) {
	fprintf(stderr, "Usage: %s {filename}\n", argv0);
	exit(1);
}

int main(int argc, char** argv) {
	int nreps = 1;
	int do_write = 0;
	if (argc != 2 && argc != 3 && argc != 4)
		usage(argv[0]);
	char* filename = argv[1];
	if (argc >= 3)
		(void)sscanf(argv[2], "%d", &nreps);
	if (argc >= 4)
		(void)sscanf(argv[3], "%d", &do_write);

	double s, e, t;
	int bc;

	for (int i = 0; i < nreps; i++) {

		s = get_systime();
		bc = read_file_mlr_alloc_read_line_single_delimiter_no_autodetect(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=single_delim_no_auto,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_mlr_alloc_read_line_single_delimiter_with_autodetect(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=single_delim_with_auto,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_mlr_alloc_read_line_multiple_delimiter(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=multiple_delim,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_fgetc_fixed_len(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=fgetc_fixed_len,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_getc_unlocked_fixed_len(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=getc_unlocked_fixed_len,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_getc_unlocked_psb(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=getc_unlocked_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_fgetc_psb(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=fgetc_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_mmap_psb(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=mmap_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);

		s = get_systime();
		bc = read_file_pfr_psb(filename, do_write);
		e = get_systime();
		t = e - s;
		printf("type=pfr_psb,t=%.6lf,n=%d\n", t, bc);
		fflush(stdout);
	}

	return 0;
}

// ================================================================
// $ ./getl ../data/big.csv 5|tee x

// $ mlr --opprint cat then sort -n t x
// type                    t        n         type                    t        n
// getdelim                0.118618 55888899  getdelim                0.118057 55888899
// getdelim                0.121467 55888899  getdelim                0.118727 55888899
// getdelim                0.121943 55888899  getdelim                0.119609 55888899
// getdelim                0.124756 55888899  getdelim                0.122506 55888899
// getdelim                0.127039 55888899  getdelim                0.123099 55888899
// getc_unlocked_fixed_len 0.167563 55888899  getc_unlocked_fixed_len 0.168109 55888899
// getc_unlocked_fixed_len 0.167803 55888899  getc_unlocked_fixed_len 0.168392 55888899
// getc_unlocked_fixed_len 0.168808 55888899  getc_unlocked_fixed_len 0.169387 55888899
// getc_unlocked_fixed_len 0.168980 55888899  getc_unlocked_fixed_len 0.178484 55888899
// getc_unlocked_fixed_len 0.176187 55888899  getc_unlocked_fixed_len 0.182793 55888899
// getc_unlocked_psb       0.238986 55888899  getc_unlocked_psb       0.293240 55888899
// getc_unlocked_psb       0.241325 55888899  getc_unlocked_psb       0.298449 55888899
// getc_unlocked_psb       0.246466 55888899  getc_unlocked_psb       0.298508 55888899
// getc_unlocked_psb       0.247592 55888899  getc_unlocked_psb       0.301125 55888899
// getc_unlocked_psb       0.248112 55888899  mmap_psb                0.313239 55888899
// mmap_psb                0.250021 55888899  mmap_psb                0.315061 55888899
// mmap_psb                0.254118 55888899  mmap_psb                0.315517 55888899
// mmap_psb                0.257428 55888899  mmap_psb                0.316790 55888899
// mmap_psb                0.261807 55888899  mmap_psb                0.320654 55888899
// mmap_psb                0.264367 55888899  getc_unlocked_psb       0.326494 55888899
// pfr_psb                 0.760035 55888900  pfr_psb                 0.417141 55888899
// pfr_psb                 0.765121 55888900  pfr_psb                 0.439269 55888899
// pfr_psb                 0.768731 55888900  pfr_psb                 0.439342 55888899
// pfr_psb                 0.771937 55888900  pfr_psb                 0.447218 55888899
// pfr_psb                 0.780460 55888900  pfr_psb                 0.453839 55888899
// fgetc_fixed_len         2.516459 55888899  fgetc_psb               2.476543 55888899
// fgetc_fixed_len         2.522877 55888899  fgetc_psb               2.477130 55888899
// fgetc_fixed_len         2.587373 55888899  fgetc_psb               2.484007 55888899
// fgetc_psb               2.590090 55888899  fgetc_psb               2.484495 55888899
// fgetc_psb               2.590536 55888899  fgetc_fixed_len         2.493730 55888899
// fgetc_fixed_len         2.608356 55888899  fgetc_fixed_len         2.528333 55888899
// fgetc_psb               2.623930 55888899  fgetc_fixed_len         2.533535 55888899
// fgetc_fixed_len         2.624310 55888899  fgetc_fixed_len         2.555377 55888899
// fgetc_psb               2.637269 55888899  fgetc_fixed_len         2.736391 55888899
//                                            fgetc_psb               2.743828 55888899

// $ mlr --opprint cat then stats1 -a min,max,stddev,mean -f t -g type then sort -n t_mean x
// type                    t_min    t_max    t_stddev t_mean
// getdelim                0.118618 0.127039 0.003232 0.122765
// getc_unlocked_fixed_len 0.167563 0.176187 0.003585 0.169868
// getc_unlocked_psb       0.238986 0.248112 0.004091 0.244496
// mmap_psb                0.250021 0.264367 0.005768 0.257548
// pfr_psb                 0.760035 0.780460 0.007667 0.769257
// fgetc_fixed_len         2.516459 2.624310 0.049478 2.571875
// fgetc_psb               2.590090 2.680364 0.037489 2.624438

// type                    t_min    t_max    t_stddev t_mean
// getdelim                0.118057 0.123099 0.002271 0.120400
// getc_unlocked_fixed_len 0.168109 0.182793 0.006768 0.173433
// getc_unlocked_psb       0.293240 0.326494 0.013134 0.303563
// mmap_psb                0.313239 0.320654 0.002771 0.316252
// pfr_psb                 0.417141 0.453839 0.013830 0.439362
// fgetc_psb               2.476543 2.743828 0.117803 2.533201
// fgetc_fixed_len         2.493730 2.736391 0.095892 2.569473

// ----------------------------------------------------------------
// Analysis:
// * getdelim is good; fatal flaw is single-char line-terminator
//   o maybe i could cobble up a line-stacked iterator which
//     consumes usually 1, sometimes (double-quote case) multiple
//     delim lines to make up the data for a given record
// * as before, maybe a 5-10% improvement mmap over stdio.
//   worth doing as a second-level refinement.
// * getc_unlocked vs. fgetc, no-brainer for this single-threaded code.
// * string-builder is a little than fixed-length malloc, as expected
//   -- it's adding value.
// ! peek_file_reader is where the optimization opportunities are
