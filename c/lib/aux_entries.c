#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"

static int lecat_main(int argc, char** argv);
static int lecat_stream(FILE* input_stream, int do_color);

static int termcvt_main(int argc, char** argv);

// ================================================================
void do_aux_entries(int argc, char** argv) {
	if (argc < 2) {
		return;
	}
	if (streq(argv[1], "lecat")) {
		exit(lecat_main(argc, argv));
	}
	if (streq(argv[1], "termcvt")) {
		exit(termcvt_main(argc, argv));
	}
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
		ssize_t linelen = getdelim(&line, &linecap, inend, input_stream);
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
