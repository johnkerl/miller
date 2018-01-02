// ================================================================
// Note: there are multiple process methods with a lot of code duplication.
// This is intentional. Much of Miller's measured processing time is in the
// lrec-reader process methods. This is code which needs to execute on every
// byte of input and even moving a single runtime if-statement into a
// function-pointer assignment at alloc time can have noticeable effects on
// performance (5-10% in some cases).
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include "cli/comment_handling.h"
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_dkvp_state_t {
	char*  irs;
	char*  ifs;
	char*  ips;
	int    irslen;
	int    ifslen;
	int    ipslen;
	int    allow_repeat_ifs;
	comment_handling_t comment_handling;
	char*  comment_string;
	size_t line_length;
} lrec_reader_stdio_dkvp_state_t;

static void    lrec_reader_stdio_dkvp_free(lrec_reader_t* preader);
static void    lrec_reader_stdio_dkvp_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others_auto_line_term(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others_auto_line_term(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle,
	context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle,
	context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs,
	comment_handling_t comment_handling, char* comment_string)
{
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->ipslen           = strlen(ips);
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	pstate->comment_handling = comment_handling;
	pstate->comment_string   = comment_string;
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
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_stdio_dkvp_process_single_irs_single_others_auto_line_term
			: lrec_reader_stdio_dkvp_process_single_irs_multi_others_auto_line_term;
	} else if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_dkvp_process_single_irs_single_others
			: &lrec_reader_stdio_dkvp_process_single_irs_multi_others;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_dkvp_process_multi_irs_single_others
			: &lrec_reader_stdio_dkvp_process_multi_irs_multi_others;
	}
	plrec_reader->psof_func     = lrec_reader_stdio_dkvp_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_dkvp_free;

	return plrec_reader;
}

static void lrec_reader_stdio_dkvp_free(lrec_reader_t* preader) {
	free(preader->pvstate);
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others_auto_line_term(
	void* pvstate, void* pvhandle, context_t* pctx)
{
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;

	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
			&pstate->line_length, TRUE, pctx)
		: mlr_alloc_read_line_single_delimiter_stripping_comments(input_stream, pstate->irs[0],
			&pstate->line_length, TRUE, pstate->comment_handling, pstate->comment_string, pctx);

	if (line == NULL) {
		return NULL;
	} else {
		return lrec_parse_stdio_dkvp_single_sep(line, pstate->ifs[0], pstate->ips[0], pstate->allow_repeat_ifs);
	}
}

static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others_auto_line_term(
	void* pvstate, void* pvhandle, context_t* pctx)
{
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;

	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0], &pstate->line_length, TRUE, pctx)
		: mlr_alloc_read_line_single_delimiter_stripping_comments(input_stream, pstate->irs[0], &pstate->line_length,
			TRUE, pstate->comment_handling, pstate->comment_string, pctx);
	if (line == NULL) {
		return NULL;
	} else {
		return lrec_parse_stdio_dkvp_multi_sep(line, pstate->ifs, pstate->ips, pstate->ifslen, pstate->ipslen,
			pstate->allow_repeat_ifs);
	}
}

static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0], &pstate->line_length, FALSE, pctx)
		: mlr_alloc_read_line_single_delimiter_stripping_comments(input_stream, pstate->irs[0], &pstate->line_length,
			FALSE, pstate->comment_handling, pstate->comment_string, pctx);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_single_sep(line, pstate->ifs[0], pstate->ips[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
			&pstate->line_length, FALSE, pctx)
		: mlr_alloc_read_line_single_delimiter_stripping_comments(input_stream, pstate->irs[0],
			&pstate->line_length, FALSE, pstate->comment_handling, pstate->comment_string, pctx);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_multi_sep(line, pstate->ifs, pstate->ips, pstate->ifslen, pstate->ipslen,
			pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
			&pstate->line_length)
		: mlr_alloc_read_line_multiple_delimiter_stripping_comments(input_stream, pstate->irs, pstate->irslen,
			&pstate->line_length, pstate->comment_handling, pstate->comment_string);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_single_sep(line, pstate->ifs[0], pstate->ips[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = pstate->comment_handling == COMMENTS_ARE_DATA
		? mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
			&pstate->line_length)
		: mlr_alloc_read_line_multiple_delimiter_stripping_comments(input_stream, pstate->irs, pstate->irslen,
			&pstate->line_length, pstate->comment_handling, pstate->comment_string);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_multi_sep(line, pstate->ifs, pstate->ips, pstate->ifslen, pstate->ipslen,
			pstate->allow_repeat_ifs);
}

// ----------------------------------------------------------------
// "abc=def,ghi=jkl"
//      P     F     P
//      S     S     S
// "abc" "def" "ghi" "jkl"

// I couldn't find a performance gain using stdlib index(3) ... *maybe* even a
// fraction of a percent *slower*.

lrec_t* lrec_parse_stdio_dkvp_single_sep(char* line, char ifs, char ips, int allow_repeat_ifs) {
	lrec_t* prec = lrec_dkvp_alloc(line);

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	char* p = line;

	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; *p; ) {
		if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
			}

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			key = p;
			value = p;
		} else if (*p == ips && !saw_ps) {
			*p = 0;
			p++;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0 || value <= key) {
			char  free_flags = 0;
			lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put(prec, key, value, NO_FREE);
		}
	}

	return prec;
}

lrec_t* lrec_parse_stdio_dkvp_multi_sep(char* line, char* ifs, char* ips, int ifslen, int ipslen,
	int allow_repeat_ifs)
{
	lrec_t* prec = lrec_dkvp_alloc(line);

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	char* p = line;

	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
			}

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			key = p;
			value = p;
		} else if (streqn(p, ips, ipslen) && !saw_ps) {
			*p = 0;
			p += ipslen;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0 || value <= key) {
			char  free_flags = 0;
			lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put(prec, key, value, NO_FREE);
		}
	}

	return prec;
}
