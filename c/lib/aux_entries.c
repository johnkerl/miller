#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"

static int lecat(int argc, char** argv);
static int lecat_stream(FILE* input_stream, int do_color);

// ----------------------------------------------------------------
void do_aux_entries(int argc, char** argv) {
	if (argc < 2) {
		return;
	}
	if (streq(argv[1], "lecat")) {
		exit(lecat(argc, argv));
	}
}


// ----------------------------------------------------------------
static int lecat(int argc, char** argv) {
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
