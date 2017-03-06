#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlr_arch.h"
#include "lib/mlrutil.h"

static int aux_help_main(int argc, char** argv);
static int lecat_main(int argc, char** argv);
static int lecat_stream(FILE* input_stream, int do_color);
static int termcvt_main(int argc, char** argv);
static int show_temp_suffixes_main(int argc, char** argv);
static int hex_main(int argc, char** argv);

typedef int aux_main_t(int argc, char**argv);

typedef struct _aux_lookup_entry_t {
	char* name;
	aux_main_t* phandler;
} aux_lookup_entry_t;

static aux_lookup_entry_t aux_lookup_table[] = {

	{ "aux_help",           aux_help_main           },
	{ "lecat",              lecat_main              },
	{ "termcvt",            termcvt_main            },
	{ "show_temp_suffixes", show_temp_suffixes_main },
	{ "hex",                hex_main                },

};

static int aux_lookup_table_size = sizeof(aux_lookup_table) / sizeof(aux_lookup_table[0]);

// ================================================================
void do_aux_entries(int argc, char** argv) {
	if (argc < 2) {
		return;
	}

	// xxx to do: list-all option.
	for (int i = 0; i < aux_lookup_table_size; i++) {
		if (streq(argv[1], aux_lookup_table[i].name)) {
			exit(aux_lookup_table[i].phandler(argc, argv));
		}
	}
}

// ----------------------------------------------------------------
static int aux_help_main(int argc, char** argv) {
	for (int i = 0; i < aux_lookup_table_size; i++) {
		printf("%s\n", aux_lookup_table[i].name);
	}
	return 0;
}

// ----------------------------------------------------------------
static int lecat_main(int argc, char** argv) {
	int ok = 1;
	int do_color = TRUE;

	// 'mlr' and 'lecat' are already argv[0] and argv[1].
	int argb = 2;
	if (argc >= 3 && argv[argb][0] == '-') {
		if (streq(argv[argb], "--mono")) {
			do_color = FALSE;
			argb++;
		} else {
			fprintf(stderr, "%s %s: unrecognized option \"%s\".\n",
				argv[0], argv[1], argv[argb]);
			return 1;
		}
	}

	if (argb == argc) {
		ok = ok && lecat_stream(stdin, do_color);
	} else {
		for (int argi = argb; argi < argc; argi++) {
			char* file_name = argv[argi];
			FILE* input_stream = fopen(file_name, "r");
			if (input_stream == NULL) {
				perror(file_name);
				exit(1);
			}
			ok = lecat_stream(input_stream, do_color);
			fclose(input_stream);
		}
	}
	return ok ? 0 : 1;
}

// ----------------------------------------------------------------
static int lecat_stream(FILE* input_stream, int do_color) {
	while (1) {
		int c = fgetc(input_stream);
		if (c == EOF)
			break;
		if (c == '\r') {
			if (do_color)
				printf("\033[31;01m"); // xterm red
			printf("[CR]");
			if (do_color)
				printf("\033[0m");
		} else if (c == '\n') {
			if (do_color)
				printf("\033[32;01m"); // xterm green
			printf("[LF]\n");
			if (do_color)
				printf("\033[0m");
		} else {
			putchar(c);
		}
	}
	return 1;
}

// ================================================================
typedef void line_cvt_func_t(char* line, ssize_t linelen, FILE* output_stream);

static void cr_to_crlf(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen == 1) {
		if (line[0] == '\r') {
			fputc('\r', output_stream);
			fputc('\n', output_stream);
		} else {
			fputc(line[0], output_stream);
		}
	} else {
		if (line[linelen-2] == '\r' && line[linelen-1] == '\n') {
			fputs(line, output_stream);
		} else if (line[linelen-1] == '\r') {
			fputs(line, output_stream);
			fputc('\n', output_stream);
		} else {
			fputs(line, output_stream);
		}
	}
}

static void lf_to_crlf(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen == 1) {
		if (line[0] == '\n') {
			fputc('\r', output_stream);
			fputc('\n', output_stream);
		} else {
			fputc(line[0], output_stream);
		}
	} else {
		if (line[linelen-2] == '\r' && line[linelen-1] == '\n') {
			fputs(line, output_stream);
		} else if (line[linelen-1] == '\n') {
			line[linelen-1] = '\r';
			fputs(line, output_stream);
			fputc('\n', output_stream);
		} else {
			fputs(line, output_stream);
		}
	}
}

static void crlf_to_cr(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen >= 2 && line[linelen-2] == '\r' && line[linelen-1] == '\n') {
		line[linelen-2] = '\r';
		line[linelen-1] = '\0';
	}
	fputs(line, output_stream);
}

static void crlf_to_lf(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen >= 2 && line[linelen-2] == '\r' && line[linelen-1] == '\n') {
		line[linelen-2] = '\n';
		line[linelen-1] = '\0';
	}
	fputs(line, output_stream);
}

static void cr_to_lf(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen >= 1 && line[linelen-1] == '\r') {
		line[linelen-1] = '\n';
	}
	fputs(line, output_stream);
}

static void lf_to_cr(char* line,  ssize_t linelen, FILE* output_stream) {
	if (linelen >= 1 && line[linelen-1] == '\n') {
		line[linelen-1] = '\r';
	}
	fputs(line, output_stream);
}

// ----------------------------------------------------------------
static int do_stream(FILE* input_stream, FILE* output_stream, char inend, line_cvt_func_t* pcvt_func) {
	while (1) {
		char* line = NULL;
		size_t linecap = 0;
		ssize_t linelen = mlr_arch_getdelim(&line, &linecap, inend, input_stream);
		if (linelen <= 0) {
			break;
		}

		pcvt_func(line, linelen, output_stream);

		free(line);
	}
	return 1;
}

// ----------------------------------------------------------------
static void termcvt_usage(char* argv0, char* argv1) {
	printf("Usage: %s %s [option] {zero or more file names}\n", argv0, argv1);
	printf("Option (exactly one is required):\n");
	printf("--cr2crlf\n");
	printf("--lf2crlf\n");
	printf("--crlf2cr\n");
	printf("--crlf2lf\n");
	printf("--cr2lf\n");
	printf("--lf2cr\n");
	printf("-h or --help: print this message\n");
	printf("Zero file names means read from standard input.\n");
	printf("Output is always to standard output; files are not written in-place.\n");

	exit(1);
}

// ----------------------------------------------------------------
static int termcvt_main(int argc, char** argv) {
	int ok = 1;
	char inend = '\n';
	line_cvt_func_t* pcvt_func = lf_to_crlf;

	// argv[0] is 'mlr'
	// argv[1] is 'termcvt'
	// argv[2] is '--some-option'
	// argv[3] and above are filenames
	if (argc < 3)
		termcvt_usage(argv[0], argv[1]);
	char* opt = argv[2];
	if (!strcmp(opt, "-h")) {
		termcvt_usage(argv[0], argv[1]);
	} else if (!strcmp(opt, "--help")) {
		termcvt_usage(argv[0], argv[1]);
	} else if (!strcmp(opt, "--cr2crlf")) {
		pcvt_func = cr_to_crlf;
		inend = '\r';
	} else if (!strcmp(opt, "--lf2crlf")) {
		pcvt_func = lf_to_crlf;
		inend = '\n';
	} else if (!strcmp(opt, "--crlf2cr")) {
		pcvt_func = crlf_to_cr;
		inend = '\n';
	} else if (!strcmp(opt, "--crlf2lf")) {
		pcvt_func = crlf_to_lf;
		inend = '\n';
	} else if (!strcmp(opt, "--cr2lf")) {
		pcvt_func = cr_to_lf;
		inend = '\r';
	} else if (!strcmp(opt, "--lf2cr")) {
		pcvt_func = lf_to_cr;
		inend = '\n';
	} else {
		termcvt_usage(argv[0], argv[1]);
	}

	if (argc == 3) {
		ok = ok && do_stream(stdin, stdout, inend, pcvt_func);
	} else {
		for (int argi = 3; argi < argc; argi++) {
			char* file_name = argv[argi];
			FILE* input_stream = fopen(file_name, "r");
			if (input_stream == NULL) {
				perror(file_name);
				exit(1);
			}
			ok = do_stream(input_stream, stdout, inend, pcvt_func);
			fclose(input_stream);
		}
	}
	return ok ? 0 : 1;
}

// ----------------------------------------------------------------
static int show_temp_suffixes_main(int argc, char** argv) {
	for (int i = 0; i < 1000; i++) {
		char* in = "a";
		char* out = alloc_suffixed_temp_file_name(in);
		printf("%s\n", out);
		free(out);

	}
	return 0;
}

// ================================================================
// Copyright (c) 1998 John Kerl.
// ================================================================
// This is a simple hex dump with hex offsets to the left, hex data in the
// middle, and ASCII at the right.  This is a subset of the functionality of
// Unix od; I wrote it in my NT days.
//
// Example:
//
// $ d2h $(jot 0 128) | unhex | hex
// 00000000: 00 01 02 03  04 05 06 07  08 09 0a 0b  0c 0d 0e 0f |................|
// 00000010: 10 11 12 13  14 15 16 17  18 19 1a 1b  1c 1d 1e 1f |................|
// 00000020: 20 21 22 23  24 25 26 27  28 29 2a 2b  2c 2d 2e 2f | !"#$%&'()*+,-./|
// 00000030: 30 31 32 33  34 35 36 37  38 39 3a 3b  3c 3d 3e 3f |0123456789:;<=>?|
// 00000040: 40 41 42 43  44 45 46 47  48 49 4a 4b  4c 4d 4e 4f |@ABCDEFGHIJKLMNO|
// 00000050: 50 51 52 53  54 55 56 57  58 59 5a 5b  5c 5d 5e 5f |PQRSTUVWXYZ[\]^_|
// 00000060: 60 61 62 63  64 65 66 67  68 69 6a 6b  6c 6d 6e 6f |`abcdefghijklmno|
// 00000070: 70 71 72 73  74 75 76 77  78 79 7a 7b  7c 7d 7e 7f |pqrstuvwxyz{|}~.|
// ================================================================

#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <fcntl.h>

#define LINE_LENGTH_MAX 8192

static void hex_dump_fp(FILE *in_fp, FILE *out_fp, int do_raw);

//----------------------------------------------------------------------
// 'mlr' and 'hex' are already argv[0] and argv[1].
static int hex_main(int argc, char **argv) {
	char * filename;
	FILE * in_fp;
	FILE * out_fp;
	int do_raw = 0;
	int argi = 2;

	if (argc >= 3 && strcmp(argv[2], "-r") == 0) {
		do_raw = 1;
		argi++;
	}

	int num_file_names = argc - argi;

	if (num_file_names == 0) {
#ifdef WINDOWS
		setmode(fileno(stdin), O_BINARY);
#endif //WINDOWS
		hex_dump_fp(stdin, stdout, do_raw);
	} else {
		for ( ; argi < argc; argi++) {
			if (!do_raw) {
				if (num_file_names > 1)
					printf("%s:\n", argv[argi]);
			}
			filename = argv[argi];
			in_fp    = fopen(filename, "rb");
			out_fp   = stdout;
			if (in_fp == NULL) {
				fprintf(stderr, "Couldn't open \"%s\"; skipping.\n",
					filename);
			}
			else {
				hex_dump_fp(in_fp, out_fp, do_raw);
				fclose(in_fp);
				if (out_fp != stdout)
					fclose(out_fp);

			}
			if (!do_raw) {
				if (num_file_names > 1)
					printf("\n");
			}
		}
	}

	return 0;
}

//----------------------------------------------------------------------
#define bytes_per_clump  4
#define clumps_per_line  4
#define buffer_size     (bytes_per_clump * clumps_per_line)

static void hex_dump_fp(FILE *in_fp, FILE *out_fp, int do_raw) {
	unsigned char buf[buffer_size];
	long num_bytes_read;
	long num_bytes_total = 0;
	int byteno;

	while ((num_bytes_read=fread(buf, sizeof(unsigned char),
		buffer_size, in_fp)) > 0)
	{
		if (!do_raw) {
			printf("%08lx: ", num_bytes_total);
		}

		for (byteno = 0; byteno < num_bytes_read; byteno++) {
			unsigned int temp = buf[byteno];
			printf("%02x ", temp);
			if ((byteno % bytes_per_clump) ==
			(bytes_per_clump - 1))
			{
				if ((byteno > 0) && (byteno < buffer_size-1))
					printf(" ");
			}
		}
		for (byteno = num_bytes_read; byteno < buffer_size; byteno++) {
			printf("   ");
			if ((byteno % bytes_per_clump) ==
			(bytes_per_clump - 1))
			{
				if ((byteno > 0) && (byteno < buffer_size-1))
					printf(" ");
			}
		}

		if (!do_raw) {
			printf("|");
			for (byteno = 0; byteno < num_bytes_read; byteno++) {
				unsigned char temp = buf[byteno];
				if (!isprint(temp))
					temp = '.';
				printf("%c", temp);
			}

			printf("|");
		}

		printf("\n");
		num_bytes_total += num_bytes_read;
	}
}
