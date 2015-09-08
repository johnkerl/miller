#include <stdio.h>
#include <ctype.h>
#include "input/peek_file_reader.h"

// ----------------------------------------------------------------
void pfr_dump(peek_file_reader_t* pfr) {
	printf("======================== pfr at %p\n", pfr);
	printf("  peekbuflen = %d\n", pfr->peekbuflen);
	printf("  npeeked    = %d\n", pfr->npeeked);
	int sob = pfr->sob;
	int eob = (pfr->sob + pfr->npeeked) & pfr->peekbuflenmask;
	for (int i = 0; i < pfr->peekbuflen; i++) {
		char c = pfr->peekbuf[i];
		char* sdesc = (i == sob) ? "START" : "";
		char* edesc = (i == eob) ? "END" : "";
		char* occdesc = "";
		if (sob <= eob) {
			if (sob <= i && i < eob)
				occdesc = "OCC";
		} else {
			if (i < eob || sob <= i)
				occdesc = "OCC";
		}
		printf("  %-5s %-5s %-3s i=%2d c=%c [%02x]\n",
			sdesc, edesc, occdesc, i, isprint((unsigned char)c) ? c : ' ', c);
	}
	printf("------------------------\n");
}
