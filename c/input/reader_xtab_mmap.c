#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/mmap.h"
#include "input/readers.h"

typedef struct _reader_xtab_mmap_state_t {
	char irs;
	char ips; // xxx make me real
	int allow_repeat_ips;
	int at_eof;
	// xxx need to remember EOF for subsequent read
} reader_xtab_mmap_state_t;

// ----------------------------------------------------------------
static lrec_t* reader_xtab_mmap_func(mmap_reader_state_t* phandle, void* pvstate, context_t* pctx) {
//	reader_xtab_mmap_state_t* pstate = pvstate;

//	if (pstate->at_eof)
		return NULL;

//	slls_t* pxtab_lines = slls_alloc();

//	while (TRUE) {
//		char* line = mlr_get_line(input_stream, '\n'); // xxx parameterize
//		if (line == NULL) { // EOF
//			// xxx cmt EOF terminates the stanza etc.
//			pstate->at_eof = TRUE;
//			if (pxtab_lines->length == 0) {
//				return NULL;
//			} else {
//				return lrec_parse_xtab(pxtab_lines, pstate->ips, pstate->allow_repeat_ips);
//			}
//		} else if (*line == '\0') {
//			if (pxtab_lines->length > 0) { // xxx make an is_empty_modulo_whitespace()
//				return lrec_parse_xtab(pxtab_lines, pstate->ips, pstate->allow_repeat_ips);
//			}
//		} else {
//			slls_add_with_free(pxtab_lines, line);
//		}
//	}
}

// xxx rename resets to sof_reset or some such
static void reset_xtab_func(void* pvstate) {
	reader_xtab_mmap_state_t* pstate = pvstate;
	pstate->at_eof = FALSE;
}

reader_mmap_t* reader_xtab_mmap_alloc(char irs, char ips, int allow_repeat_ips) {
	reader_mmap_t* preader = mlr_malloc_or_die(sizeof(reader_mmap_t));

	reader_xtab_mmap_state_t* pstate = mlr_malloc_or_die(sizeof(reader_xtab_mmap_state_t));
	//pstate->ips              = ips;
	//pstate->allow_repeat_ips = allow_repeat_ips;
	pstate->irs              = irs;
	pstate->ips              = ' ';
	pstate->allow_repeat_ips = TRUE;
	pstate->at_eof           = FALSE;
	preader->pvstate         = (void*)pstate;

	preader->preader_func = &reader_xtab_mmap_func;
	preader->preset_func  = &reset_xtab_func;

	return preader;
}
