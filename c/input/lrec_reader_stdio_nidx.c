// ================================================================
// Note: there are multiple process methods with a lot of code duplication.
// This is intentional. Much of Miller's measured processing time is in the
// lrec-reader process methods. This is code which needs to execute on every
// byte of input and even moving a single runtime if-statement into a
// function-pointer assignment at alloc time can have noticeable effects on
// performance (5-10% in some cases).
// ================================================================

#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_nidx_state_t {
	char*  irs;
	char*  ifs;
	int    irslen;
	int    ifslen;
	int    allow_repeat_ifs;
	size_t line_length;
} lrec_reader_stdio_nidx_state_t;

static void    lrec_reader_stdio_nidx_free(lrec_reader_t* preader);
static void    lrec_reader_stdio_nidx_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs_auto_line_term(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs_auto_line_term(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs, char* comment_string) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_nidx_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	// This is used to track nominal line length over the file read. Bootstrap with a default length.
	pstate->line_length      = MLR_ALLOC_READ_LINE_INITIAL_SIZE;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_stdio_vopen;
	plrec_reader->pclose_func   = file_reader_stdio_vclose;
	if (streq(irs, "auto")) {
		// Auto means either lines end in "\n" or "\r\n" (LF or CRLF).  In
		// either case the final character is "\n". Then for autodetect we
		// simply check if there's a character in the line before the '\n', and
		// if that is '\r'.
		pstate->irs = "\n";
		pstate->irslen = 1;
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_stdio_nidx_process_single_irs_single_ifs_auto_line_term
			: lrec_reader_stdio_nidx_process_single_irs_multi_ifs_auto_line_term;
	} else if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_nidx_process_single_irs_single_ifs
			: &lrec_reader_stdio_nidx_process_single_irs_multi_ifs;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_nidx_process_multi_irs_single_ifs
			: &lrec_reader_stdio_nidx_process_multi_irs_multi_ifs;
	}
	plrec_reader->psof_func     = lrec_reader_stdio_nidx_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_nidx_free;

	return plrec_reader;
}

static void lrec_reader_stdio_nidx_free(lrec_reader_t* preader) {
	free(preader->pvstate);
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_nidx_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs_auto_line_term(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;

	char* line = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
		&pstate->line_length, TRUE, pctx);

	if (line == NULL) {
		return NULL;
	} else {
		return lrec_parse_stdio_nidx_single_sep(line, pstate->ifs[0], pstate->allow_repeat_ifs);
	}
}

static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs_auto_line_term(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;

	char* line = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
		&pstate->line_length, TRUE, pctx);

	if (line == NULL) {
		return NULL;
	} else {
		return lrec_parse_stdio_nidx_multi_sep(line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
	}
}

static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;

	char* line = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
		&pstate->line_length, FALSE, pctx);

	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_single_sep(line, pstate->ifs[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;

	char* line = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
		&pstate->line_length, FALSE, pctx);

	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_multi_sep(line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
		&pstate->line_length);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_single_sep(line, pstate->ifs[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
		&pstate->line_length);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_multi_sep(line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_nidx_single_sep(char* line, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int idx = 0;
	char  free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; *p; ) {
		if (*p == ifs) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			value = p;
		} else {
			p++;
		}
	}
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = low_int_to_string(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_nidx_multi_sep(char* line, char* ifs, int ifslen, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int  idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			value = p;
		} else {
			p++;
		}
	}
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = low_int_to_string(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}
