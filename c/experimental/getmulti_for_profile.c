#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/peek_file_reader.h"

#define PEEK_BUF_LEN             32
#define STRING_BUILDER_INIT_SIZE 1024

static char* mlr_get_line2(FILE* input_stream, char* irs, int irs_len,
	peek_file_reader_t* pfr, string_builder_t* psb)
{
	while (TRUE) {
		if (pfr_at_eof(pfr)) {
			if (sb_is_empty(psb))
				return NULL;
			else
				return sb_finish(psb);
		} else if (pfr_next_is(pfr, irs, irs_len)) {
			if (!pfr_advance_past(pfr, irs)) {
				fprintf(stderr, "%s: Internal coding error: IRS found and lost.\n", MLR_GLOBALS.argv0);
				exit(1);
			}
			return sb_finish(psb);
		} else {
			sb_append_char(psb, pfr_read_char(pfr));
		}
	}
}

int main(void) {
	FILE* input_stream = stdin;
	peek_file_reader_t* pfr = pfr_alloc(input_stream, PEEK_BUF_LEN);
	string_builder_t  sb;
	string_builder_t* psb = &sb;
	sb_init(&sb, STRING_BUILDER_INIT_SIZE);

	while (1) {
		char* line = mlr_get_line2(stdin, "\r\n", 2, pfr, psb);
		if (line == NULL)
			break;
		fputs(line, stdout);
		free(line);
	}
	return 0;
}
