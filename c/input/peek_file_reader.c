#include <stdio.h>
#include <ctype.h>
#include "input/peek_file_reader.h"

// ----------------------------------------------------------------
void pfr_dump(peek_file_reader_t* pfr) {
	printf("======================== pfr at %p\n", pfr);
	printf("  peekbuflen = %d\n", pfr->peekbuflen);
	printf("  npeeked    = %d\n", pfr->npeeked);
	for (int i = 0; i < pfr->peekbuflen; i++) {
		char c = pfr->peekbuf[i];
		char* sdesc = (i == pfr->sob) ? "START" : "";
		char* edesc = (i == pfr->sob + pfr->npeeked - 1) ? "END" : "";
		printf("  %-5s %-5s i=%2d c=%c [%02x]\n",
			sdesc, edesc, i, isprint((unsigned char)c) ? c : ' ', c);
	}
	printf("------------------------\n");
}
