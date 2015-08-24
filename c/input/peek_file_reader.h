#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>

typedef struct _peek_file_reader_t {
	FILE* fp;
	int   peekbuflen;
	char* peekbuf;
	int   npeeked;
} peek_file_reader_t;

// The caller should fclose the fp, since presumably it will have opened it. We
// could have our constructor fopen (taking not fp but filename as argument)
// and the destructor fclose but that would break reading from stdin.

peek_file_reader_t* pfr_alloc(FILE* fp, int maxnpeek);
// xxx needing contextual comments here.
int  pfr_at_eof(peek_file_reader_t* pfr);
int  pfr_next_is(peek_file_reader_t* pfr, char* string, int len);
char pfr_read_char(peek_file_reader_t* pfr);
int  pfr_advance_past(peek_file_reader_t* pfr, char* string);
void pfr_free(peek_file_reader_t* pfr);

void pfr_dump(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
