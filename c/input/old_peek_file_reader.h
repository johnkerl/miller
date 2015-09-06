#ifndef OLD_PEEK_FILE_READER_H
#define OLD_PEEK_FILE_READER_H

#include <stdio.h>

typedef struct _old_peek_file_reader_t {
	FILE* fp;
	int   peekbuflen;
	char* peekbuf;
	int   npeeked;
} old_peek_file_reader_t;

// The caller should fclose the fp, since presumably it will have opened it. We
// could have our constructor do the fopen (taking not fp but filename as
// argument) and the destructor do the fclose but that would break reading from
// stdin.

old_peek_file_reader_t* old_pfr_alloc(FILE* fp, int maxnpeek);
// xxx needing contextual comments here.
int  old_pfr_at_eof(old_peek_file_reader_t* pfr);
int  old_pfr_next_is(old_peek_file_reader_t* pfr, char* string, int len);
char old_pfr_read_char(old_peek_file_reader_t* pfr);
void old_pfr_advance_by(old_peek_file_reader_t* pfr, int len);
void old_pfr_free(old_peek_file_reader_t* pfr);

void old_pfr_dump(old_peek_file_reader_t* pfr);

#endif // OLD_PEEK_FILE_READER_H
