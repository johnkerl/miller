#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MYBUFSIZ 8192
static char iobuf[MYBUFSIZ];

// ----------------------------------------------------------------
static int do_stream(char* file_name) {
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
		char* line = fgets(iobuf, BUFSIZ, input_stream);
		if (line == NULL)
			break;
		fputs(line, output_stream);
	}
	if (input_stream != stdin)
		fclose(input_stream);

	return 1;
}

// ================================================================
int main(int argc, char** argv) {
	int ok = 1;
	if (argc == 1) {
		ok = ok && do_stream("-");
	} else {
		for (int argi = 1; argi < argc; argi++) {
		    ok = do_stream(argv[argi]);
		}
	}
	return ok ? 0 : 1;
}
