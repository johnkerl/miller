#include <stdio.h>
#include <string.h>
#include "input/old_peek_file_reader.h"

int main(int argc, char** argv) {
	FILE* fp = stdin;
	old_peek_file_reader_t* pfr = old_pfr_alloc(fp, 32);

	printf("@eof = %d\n", old_pfr_at_eof(pfr));
	old_pfr_dump(pfr);
	printf("read 0x%02x\n", (unsigned)old_pfr_read_char(pfr));
	old_pfr_dump(pfr);
	char* s = "//";
	old_pfr_dump(pfr);
	printf("next is %s = %d\n", s, old_pfr_next_is(pfr, s, strlen(s)));
	old_pfr_dump(pfr);
	char c = old_pfr_read_char(pfr);
	printf("read %c [0x%02x]\n", c, (unsigned)c);
	old_pfr_dump(pfr);

	printf("@eof = %d\n", old_pfr_at_eof(pfr));

	old_pfr_free(pfr);
	return 0;
}
