#ifndef PEEK_FILE_READER_H
#define PEEK_FILE_READER_H

#include <stdio.h>

typedef struct _peek_file_reader_t {
	//FILE* fp;
	//int   peekbuflen;
	//char* peekbuf;
	//int   npeeked;
} peek_file_reader_t;

peek_file_reader_t* pfr_alloc(int maxnpeek);
// xxx needing contextual comments here.
//int  old_pfr_at_eof(peek_file_reader_t* pfr);
//int  old_pfr_next_is(peek_file_reader_t* pfr, char* string, int len);
//char old_pfr_read_char(peek_file_reader_t* pfr);
//void old_pfr_advance_by(peek_file_reader_t* pfr, int len);
void old_pfr_free(peek_file_reader_t* pfr);

void pfr_dump(peek_file_reader_t* pfr);

#endif // PEEK_FILE_READER_H
