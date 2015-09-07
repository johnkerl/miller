#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/peek_file_reader.h"

// ----------------------------------------------------------------
peek_file_reader_t* pfr_alloc(byte_reader_t* pbr, int maxnpeek) {
	peek_file_reader_t* pfr = mlr_malloc_or_die(sizeof(peek_file_reader_t));
	pfr->pbr        =  pbr;
	pfr->peekbuflen =  maxnpeek + 1;
	pfr->peekbuf    =  mlr_malloc_or_die(pfr->peekbuflen);
	memset(pfr->peekbuf, 0, pfr->peekbuflen);
	pfr->npeeked    =  0;

	return pfr;
}

// ----------------------------------------------------------------
void pfr_free(peek_file_reader_t* pfr) {
	if (pfr == NULL)
		return;
	free(pfr->peekbuf);
	free(pfr);
}

// ----------------------------------------------------------------
void pfr_reset(peek_file_reader_t* pfr) {
	memset(pfr->peekbuf, 0, pfr->peekbuflen);
	pfr->npeeked    =  0;
}

// ----------------------------------------------------------------
char pfr_peek_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked < 1) {
		pfr->peekbuf[pfr->npeeked++] = pfr->pbr->pread_func(pfr->pbr);
	}
	return pfr->peekbuf[0];
}

// ----------------------------------------------------------------
char pfr_read_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked < 1) {
		return pfr->pbr->pread_func(pfr->pbr);
	} else {
		// xxx to do: make this a ring buffer to avoid the shifts.
		char c = pfr->peekbuf[0];
		for (int i = 1; i < pfr->npeeked; i++)
			pfr->peekbuf[i-1] = pfr->peekbuf[i];
		pfr->npeeked--;
		return c;
	}
}

// ----------------------------------------------------------------
void pfr_buffer_by(peek_file_reader_t* pfr, int len) {
	while (pfr->npeeked < len) {
		pfr->peekbuf[pfr->npeeked++] = pfr->pbr->pread_func(pfr->pbr);
	}
}

// ----------------------------------------------------------------
void pfr_advance_by(peek_file_reader_t* pfr, int len) {
	if (len > pfr->npeeked) {
		fprintf(stderr, "%s: internal coding error: advance-by %d exceeds buffer depth %d.\n",
			MLR_GLOBALS.argv0, len, pfr->npeeked);
		exit(1);
	}
	for (int i = len; i < pfr->npeeked; i++)
		pfr->peekbuf[i-len] = pfr->peekbuf[i];
	pfr->npeeked -= len;
}

// ----------------------------------------------------------------
void pfr_dump(peek_file_reader_t* pfr) {
	// xxx stub
	printf("PFR DUMP STUB. MAYBE REMOVE THIS ROUTINE.\n");
}
