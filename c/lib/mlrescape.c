#include "lib/mlrescape.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"

// Avoids shell-injection cases by replacing single-quote with backslash single-quote,
// then wrapping the entire result in initial and final single-quote.
char* alloc_file_name_escaped_for_popen(char* filename) {
	string_builder_t* psb = sb_alloc(strlen(filename));

	sb_append_char(psb, '\'');
	for (char* p = filename; *p; p++) {
		char c = *p;
		//if (c == '\'') {
			//sb_append_char(psb, '\\');
		//}
		//sb_append_char(psb, c);
		if (c == '\'') {
			sb_append_char(psb, '\'');
			sb_append_char(psb, '\\');
			sb_append_char(psb, '\'');
			sb_append_char(psb, '\'');
		} else {
			sb_append_char(psb, c);
		}
	}

	sb_append_char(psb, '\'');
	char* rv = sb_finish(psb);
	sb_free(psb);
	return rv;
}
