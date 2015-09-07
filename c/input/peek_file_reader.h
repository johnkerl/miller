#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/byte_reader.h"

typedef struct _peek_file_reader_t {
	byte_reader_t* pbr;
	int   peekbuflen;
	char* peekbuf;
	int   npeeked;
} peek_file_reader_t;

// xxx needing contextual comments here.

// xxx to do: try using a ring buffer (power-of-two length >= buflen) instead
// of the current slipback buffer, for performance

// ----------------------------------------------------------------
static inline peek_file_reader_t* pfr_alloc(byte_reader_t* pbr, int maxnpeek) {
	peek_file_reader_t* pfr = mlr_malloc_or_die(sizeof(peek_file_reader_t));
	pfr->pbr        =  pbr;
	pfr->peekbuflen =  maxnpeek + 1;
	pfr->peekbuf    =  mlr_malloc_or_die(pfr->peekbuflen);
	memset(pfr->peekbuf, 0, pfr->peekbuflen);
	pfr->npeeked    =  0;

	return pfr;
}

// ----------------------------------------------------------------
static inline void pfr_free(peek_file_reader_t* pfr) {
	if (pfr == NULL)
		return;
	free(pfr->peekbuf);
	free(pfr);
}

// ----------------------------------------------------------------
static inline void pfr_reset(peek_file_reader_t* pfr) {
	memset(pfr->peekbuf, 0, pfr->peekbuflen);
	pfr->npeeked = 0;
}

// ----------------------------------------------------------------
static inline char pfr_peek_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked < 1) {
		pfr->peekbuf[pfr->npeeked++] = pfr->pbr->pread_func(pfr->pbr);
	}
	return pfr->peekbuf[0];
}

// ----------------------------------------------------------------
static inline char pfr_read_char(peek_file_reader_t* pfr) {
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
static inline void pfr_buffer_by(peek_file_reader_t* pfr, int len) {
	while (pfr->npeeked < len) {
		pfr->peekbuf[pfr->npeeked++] = pfr->pbr->pread_func(pfr->pbr);
	}
}

// ----------------------------------------------------------------
static inline void pfr_advance_by(peek_file_reader_t* pfr, int len) {
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
void pfr_dump(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
