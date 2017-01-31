#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Just a catter but with line-ending chars called out in red

// ----------------------------------------------------------------
static int do_stream(FILE* input_stream) {
	while (1) {
		int c = fgetc(input_stream);
		if (c == EOF)
			break;
		if (c == '\r') {
			printf("\033[31;01m"); // xterm red
			printf(" [CR]");
			printf("\033[0m");
		} else if (c == '\n') {
			printf("\033[32;01m"); // xterm green
			printf(" [LF]\n");
			printf("\033[0m");
		} else {
			putchar(c);
		}
	}
	return 1;
}

// ----------------------------------------------------------------
int main(int argc, char** argv) {
	int ok = 1;

	if (argc == 1) {
		ok = ok && do_stream(stdin);
	} else {
		for (int argi = 1; argi < argc; argi++) {
			char* file_name = argv[argi];
			FILE* input_stream = fopen(file_name, "r");
			if (input_stream == NULL) {
				perror(file_name);
				exit(1);
			}
			ok = do_stream(input_stream);
			fclose(input_stream);
		}
	}
	return ok ? 0 : 1;
}
