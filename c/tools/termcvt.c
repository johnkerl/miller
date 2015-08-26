#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// .    CR LF CRLF
// CR   -  1  1
// LF   2  -  2
// CRLF x  x  -

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

// ================================================================
static void usage(char* argv0) {
	printf("Usage: %s [option] {zero or more file names}\n", argv0);
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

int main(int argc, char** argv) {
	int ok = 1;
	char inend = '\n';
	line_cvt_func_t* pcvt_func = lf_to_crlf;

	if (argc < 2)
		usage(argv[0]);
	char* opt = argv[1];
	if (!strcmp(opt, "-h")) {
		usage(opt);
	} else if (!strcmp(opt, "--help")) {
		usage(opt);
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
		usage(argv[0]);
	}

	if (argc == 2) {
		ok = ok && do_stream(stdin, stdout, inend, pcvt_func);
	} else {
		for (int argi = 2; argi < argc; argi++) {
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
