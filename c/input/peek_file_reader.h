#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>

typedef struct _peek_file_reader_t {
	FILE* fp;
	int   peekbuflen;
	char* peekbuf;
	int   npeeked;
} peek_file_reader_t;

// xxx needing contextual comments here.
int  pfr_at_eof(peek_file_reader_t* pfr);
int  pfr_next_is(peek_file_reader_t* pfr, char* string, int len);
char pfr_read_char(peek_file_reader_t* pfr);
void pfr_advance_past(peek_file_reader_t* pfr, char* string);
void pfr_close(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
