#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>
#include "lib/mlrutil.h"
#include "lib/mlrmath.h"
#include "lib/mlr_globals.h"
#include "input/byte_reader.h"

// This is a ring-buffered peekahead file/string reader.

// Note: Throughout Miller as a general rule I treat struct attributes as if
// they were private attributes. However, for performance, parse_trie_ring_match
// accesses this ring buffer directly.
typedef struct _peek_file_reader_t {
	byte_reader_t* pbr;
	int   peekbuflen;
	int   peekbuflenmask;
	char* peekbuf;
	int   sob; // start of ring-buffer
	int   npeeked;
} peek_file_reader_t;

// ----------------------------------------------------------------
static inline peek_file_reader_t* pfr_alloc(byte_reader_t* pbr, int maxnpeek) {
	peek_file_reader_t* pfr = mlr_malloc_or_die(sizeof(peek_file_reader_t));
	pfr->pbr            = pbr;
	pfr->peekbuflen     = power_of_two_ceil(maxnpeek);
	pfr->peekbuflenmask = pfr->peekbuflen - 1;
	// The peek-buffer doesn't contain null-terminated C strings, but we
	// nonetheless null-terminate the buffer with an extra never-touched byte
	// so that print statements in the debugger, etc. will be nice.
	pfr->peekbuf        =  mlr_malloc_or_die(pfr->peekbuflen + 1);
	memset(pfr->peekbuf, 0, pfr->peekbuflen + 1);
	pfr->sob            =  0;
	pfr->npeeked        =  0;

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
	pfr->sob     = 0;
	pfr->npeeked = 0;
}

// ----------------------------------------------------------------
static inline char pfr_peek_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked < 1) {
		pfr->peekbuf[pfr->sob] = pfr->pbr->pread_func(pfr->pbr);
		pfr->npeeked++;
	}
	return pfr->peekbuf[pfr->sob];
}

// ----------------------------------------------------------------
static inline char pfr_read_char(peek_file_reader_t* pfr) {
	if (pfr->npeeked < 1) {
		return pfr->pbr->pread_func(pfr->pbr);
	} else {
		char c = pfr->peekbuf[pfr->sob];
		pfr->sob = (pfr->sob + 1) & pfr->peekbuflenmask;
		pfr->npeeked--;
		return c;
	}
}

// ----------------------------------------------------------------
static inline void pfr_buffer_by(peek_file_reader_t* pfr, int len) {
	while (pfr->npeeked < len) {
		pfr->peekbuf[(pfr->sob + pfr->npeeked++) & pfr->peekbuflenmask] = pfr->pbr->pread_func(pfr->pbr);
	}
}

// ----------------------------------------------------------------
static inline void pfr_advance_by(peek_file_reader_t* pfr, int len) {
	if (len > pfr->npeeked) {
		fprintf(stderr, "%s: internal coding error: advance-by %d exceeds buffer depth %d.\n",
			MLR_GLOBALS.bargv0, len, pfr->npeeked);
		exit(1);
	}
	pfr->sob = (pfr->sob + len) & pfr->peekbuflenmask;
	pfr->npeeked -= len;
}

// ----------------------------------------------------------------
void pfr_print(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
