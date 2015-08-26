#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// ----------------------------------------------------------------
static int do_stream(char* file_name, int do_add) {
	FILE* input_stream  = stdin;
	FILE* output_stream = stdout;

	if (strcmp(file_name, "-")) {
		input_stream = fopen(file_name, "r");
		if (input_stream == NULL) {
			perror(file_name);
			return 0;
		}
	}

	while (1) {
		char* line = NULL;
		size_t linecap = 0;
		ssize_t linelen = getdelim(&line, &linecap, '\n', input_stream);
		if (linelen <= 0) {
			break;
		}
		if (do_add) {
			// replace "\n" with "\r\n" unless already there
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
		} else {
			// replace "\r\n" with "\n" if there
			if (linelen >= 2 && line[linelen-2] == '\r' && line[linelen-1] == '\n') {
				line[linelen-2] = '\n';
				line[linelen-1] = '\0';
			}
			fputs(line, output_stream);
		}
		free(line);
	}
	if (input_stream != stdin)
		fclose(input_stream);

	return 1;
}

// ================================================================
int main(int argc, char** argv) {
	int ok = 1;
	int do_add = 1;
	if (argc >= 2 && (!strcmp(argv[1], "-h") || !strcmp(argv[1], "--help"))) {
		printf("Usage: %s [-u] {zero or more file names}\n", argv[0]);
		printf("Without -u, replaces LF with CRLF in input stream.\n");
		printf("With    -u, replaces CRLF with LF in input stream.\n");
		printf("Zero file names means read from standard input.\n");
		printf("Output is always to standard output; files are not written in-place.\n");
		exit(1);
	}
	if (argc >= 2 && !strcmp(argv[1], "-u")) {
		do_add = 0;
		argc--;
		argv++;
	}
	if (argc == 1) {
		ok = ok && do_stream("-", do_add);
	} else {
		for (int argi = 1; argi < argc; argi++) {
		    ok = do_stream(argv[argi], do_add);
		}
	}
	return ok ? 0 : 1;
}
