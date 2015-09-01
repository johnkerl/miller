#include <stdio.h>
#include <ctype.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "peek_file_reader.h"

// xxx comment about efficiency here: enough to deliver rfc-csv feature with performance tuning still tbd.

// tripartite ascii art w/ chars-to-caller; the peekbuf; pending data in the fp.
// label in particular eof handling.

// ----------------------------------------------------------------
peek_file_reader_t* pfr_alloc(FILE* fp, int maxnpeek) {
	peek_file_reader_t* pfr = mlr_malloc_or_die(sizeof(peek_file_reader_t));
	pfr->fp         = fp;
	pfr->peekbuflen = maxnpeek + 1;
	pfr->peekbuf    = mlr_malloc_or_die(pfr->peekbuflen);
	memset(pfr->peekbuf, 0, pfr->peekbuflen);
	pfr->npeeked    = 0;

	// Pre-read one char into the peekbuf so that we can say pfr_at_eof
	// right away on the first call on an empty file.
	// getc_unlocked() is appropriate since Miller is single-threaded.
	pfr->peekbuf[pfr->npeeked++] = getc_unlocked(pfr->fp); // maybe EOF

	return pfr;
}

// ----------------------------------------------------------------
int pfr_at_eof(peek_file_reader_t* pfr) {
	return pfr->npeeked >= 1 && pfr->peekbuf[0] == EOF;
}

// ----------------------------------------------------------------
// xxx inline this for perf.
int pfr_next_is(peek_file_reader_t* pfr, char* string, int len) {
	// xxx abend on len > peekbuflen
	while (pfr->npeeked < len) {
		char c = getc_unlocked(pfr->fp); // maybe EOF
		pfr->peekbuf[pfr->npeeked++] = c;
	}
	// xxx make a memeq, inlined.
	return memcmp(string, pfr->peekbuf, len) == 0;
}

// ----------------------------------------------------------------
char pfr_read_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked == 1 && pfr->peekbuf[0] == EOF) {
		return EOF;
	} else if (pfr->npeeked == 0) {
		pfr->peekbuf[0] = getc_unlocked(pfr->fp); // maybe EOF
		pfr->npeeked = 1;
		return pfr->peekbuf[0];
	} else {
		char c = pfr->peekbuf[0];
		pfr->npeeked--;
		for (int i = 0; i < pfr->npeeked; i++)
			pfr->peekbuf[i] = pfr->peekbuf[i+1];
		return c;
	}
}

// ----------------------------------------------------------------
int pfr_advance_past(peek_file_reader_t* pfr, char* string) {
	for (char* p = string; *p; p++)
		if (pfr_read_char(pfr) != *p)
			return FALSE;
	return TRUE;
}

// ----------------------------------------------------------------
void pfr_free(peek_file_reader_t* pfr) {
	if (pfr == NULL)
		return;
	free(pfr->peekbuf);
	pfr->fp      = NULL;
	pfr->peekbuf = NULL;
	free(pfr);
}

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
