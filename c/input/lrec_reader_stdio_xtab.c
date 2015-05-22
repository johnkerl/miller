#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_xtab_state_t {
	char ips; // xxx make me real
	int allow_repeat_ips;
	int at_eof;
	// xxx need to remember EOF for subsequent read
} lrec_reader_stdio_xtab_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_xtab_func(FILE* input_stream, void* pvstate, context_t* pctx) {
	lrec_reader_stdio_xtab_state_t* pstate = pvstate;

	if (pstate->at_eof)
		return NULL;

	slls_t* pxtab_lines = slls_alloc();

	while (TRUE) {
		char* line = mlr_get_line(input_stream, '\n'); // xxx parameterize
		if (line == NULL) { // EOF
			// xxx cmt EOF terminates the stanza etc.
			pstate->at_eof = TRUE;
			if (pxtab_lines->length == 0) {
				return NULL;
			} else {
				return lrec_parse_stdio_xtab(pxtab_lines, pstate->ips, pstate->allow_repeat_ips);
			}
		} else if (*line == '\0') {
			free(line);
			if (pxtab_lines->length > 0) { // xxx make an is_empty_modulo_whitespace()
				return lrec_parse_stdio_xtab(pxtab_lines, pstate->ips, pstate->allow_repeat_ips);
			}
		} else {
			slls_add_with_free(pxtab_lines, line);
		}
	}
}

// xxx rename resets to sof_reset or some such
static void reset_xtab_func(void* pvstate) {
	lrec_reader_stdio_xtab_state_t* pstate = pvstate;
	pstate->at_eof = FALSE;
}

static void lrec_reader_stdio_xtab_free(void* pvstate) {
}

lrec_reader_stdio_t* lrec_reader_stdio_xtab_alloc(char ips, int allow_repeat_ips) {
	lrec_reader_stdio_t* plrec_reader_stdio = mlr_malloc_or_die(sizeof(lrec_reader_stdio_t));

	lrec_reader_stdio_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_xtab_state_t));
	//pstate->ips              = ips;
	//pstate->allow_repeat_ips = allow_repeat_ips;
	pstate->ips              = ' ';
	pstate->allow_repeat_ips = TRUE;
	pstate->at_eof           = FALSE;
	plrec_reader_stdio->pvstate         = (void*)pstate;

	plrec_reader_stdio->plrec_reader_stdio_func = &lrec_reader_stdio_xtab_func;
	plrec_reader_stdio->preset_func  = &reset_xtab_func;
	plrec_reader_stdio->pfree_func   = &lrec_reader_stdio_xtab_free;;

	return plrec_reader_stdio;
}
