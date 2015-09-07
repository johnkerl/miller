#include <stdio.h>
#include <ctype.h>
#include "input/peek_file_reader.h"

// ----------------------------------------------------------------
void pfr_dump(peek_file_reader_t* pfr) {
	printf("======================== pfr at %p\n", pfr);
	printf("  peekbuflen = %d\n", pfr->peekbuflen);
	printf("  npeeked    = %d\n", pfr->npeeked);
	for (int i = 0; i < pfr->npeeked; i++) {
		char c = pfr->peekbuf[i];
		printf("  i=%d c=%c [%02x]\n", i, isprint((unsigned char)c) ? c : ' ', c);
	}
	printf("------------------------\n");
}
