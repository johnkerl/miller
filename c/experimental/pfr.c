#include <stdio.h>
#include <string.h>
#include "input/peek_file_reader.h"

int main(int argc, char** argv) {
	FILE* fp = stdin;
	peek_file_reader_t* pfr = pfr_alloc(fp, 32);

	printf("@eof = %d\n", pfr_at_eof(pfr));
	pfr_dump(pfr);
	printf("read 0x%02x\n", (unsigned)pfr_read_char(pfr));
	pfr_dump(pfr);
	char* s = "//";
	pfr_dump(pfr);
	printf("next is %s = %d\n", s, pfr_next_is(pfr, s, strlen(s)));
	pfr_dump(pfr);
	char c = pfr_read_char(pfr);
	printf("read %c [0x%02x]\n", c, (unsigned)c);
	pfr_dump(pfr);

	printf("@eof = %d\n", pfr_at_eof(pfr));

	pfr_free(pfr);
	return 0;
}
