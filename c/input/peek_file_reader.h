#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>
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

peek_file_reader_t* pfr_alloc(byte_reader_t* pbr, int maxnpeek);
void pfr_free(peek_file_reader_t* pfr);

char pfr_peek_char(peek_file_reader_t* pfr);
void pfr_buffer_by(peek_file_reader_t* pfr, int len);
void pfr_advance_by(peek_file_reader_t* pfr, int len);

void pfr_dump(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
