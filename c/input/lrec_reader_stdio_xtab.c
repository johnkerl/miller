#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_xtab_state_t {
	char ips; // xxx make me real
	int allow_repeat_ips;
	int at_eof;
} lrec_reader_stdio_xtab_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_xtab_process(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_xtab_state_t* pstate = pvstate;

	if (pstate->at_eof)
		return NULL;

	slls_t* pxtab_lines = slls_alloc();

	while (TRUE) {
		char* line = mlr_get_cline(input_stream, '\n'); // xxx parameterize
		if (line == NULL) { // EOF
			// EOF or blank line terminates the stanza.
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

static void lrec_reader_stdio_xtab_sof(void* pvstate) {
	lrec_reader_stdio_xtab_state_t* pstate = pvstate;
	pstate->at_eof = FALSE;
}

static void lrec_reader_stdio_xtab_free(void* pvstate) {
}

lrec_reader_t* lrec_reader_stdio_xtab_alloc(char ips, int allow_repeat_ips) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_xtab_state_t));
	//pstate->ips              = ips;
	//pstate->allow_repeat_ips = allow_repeat_ips;
	pstate->ips              = ' ';
	pstate->allow_repeat_ips = TRUE;
	pstate->at_eof           = FALSE;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_stdio_vopen;
	plrec_reader->pclose_func   = &file_reader_stdio_vclose;
	plrec_reader->pprocess_func = &lrec_reader_stdio_xtab_process;
	plrec_reader->psof_func     = &lrec_reader_stdio_xtab_sof;
	plrec_reader->pfree_func    = &lrec_reader_stdio_xtab_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_xtab(slls_t* pxtab_lines, char ips, int allow_repeat_ips) {
	lrec_t* prec = lrec_xtab_alloc(pxtab_lines);

	for (sllse_t* pe = pxtab_lines->phead; pe != NULL; pe = pe->pnext) {
		char* line = pe->value;
		char* p = line;
		char* key = p;

		while (*p != 0 && *p != ips)
			p++;
		if (*p == 0) {
			lrec_put_no_free(prec, key, "");
		} else {
			while (*p != 0 && *p == ips) {
				*p = 0;
				p++;
			}
			lrec_put_no_free(prec, key, p);
		}
	}

	return prec;
}
