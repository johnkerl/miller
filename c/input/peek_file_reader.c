#include "lib/mlrutil.h"
#include "peek_file_reader.h"

// typedef struct _peek_file_reader_t {
//     FILE* fp;
//     int   peekbuflen;
//     char* peekbuf;
//     int   npeeked;
// } peek_file_reader_t;

int pfr_at_eof(peek_file_reader_t* pfr) {
	return TRUE; // xxx stub
}

int pfr_next_is(peek_file_reader_t* pfr, char* string, int len) {
	return FALSE; // xxx stub
}

char pfr_read_char(peek_file_reader_t* pfr) {
	// xxx stub
	return 'x';
}

void pfr_advance_past(peek_file_reader_t* pfr, char* string) {
	// xxx stub
}

void pfr_close(peek_file_reader_t* pfr) {
	fclose(pfr->fp);
	pfr->fp = NULL;
}
