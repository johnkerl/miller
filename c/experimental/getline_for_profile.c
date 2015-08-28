#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"

int main(void) {
	while (1) {
		char* line = mlr_get_line(stdin, '\n');
		if (line == NULL)
			break;
		fputs(line, stdout);
		fputc('\n', stdout);
		free(line);
	}
	return 0;
}
