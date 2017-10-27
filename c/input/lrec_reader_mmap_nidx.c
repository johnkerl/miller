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
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_nidx_state_t {
	char* irs;
	char* ifs;
	int   irslen;
	int   ifslen;
	int   allow_repeat_ifs;
	int   do_auto_line_term;
	char* comment_string;
	int   comment_string_length;
} lrec_reader_mmap_nidx_state_t;

static void    lrec_reader_mmap_nidx_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_nidx_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);

static lrec_t* lrec_parse_mmap_nidx_single_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char irs, char ifs, lrec_reader_mmap_nidx_state_t* pstate, context_t* pctx);

static lrec_t* lrec_parse_mmap_nidx_single_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char irs, lrec_reader_mmap_nidx_state_t* pstate, context_t* pctx);

static lrec_t* lrec_parse_mmap_nidx_multi_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char ifs, lrec_reader_mmap_nidx_state_t* pstate);

static lrec_t* lrec_parse_mmap_nidx_multi_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	lrec_reader_mmap_nidx_state_t* pstate);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs, char* comment_string) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_nidx_state_t));
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->irslen                   = strlen(pstate->irs);
	pstate->ifslen                   = strlen(pstate->ifs);
	pstate->allow_repeat_ifs         = allow_repeat_ifs;
	pstate->do_auto_line_term        = FALSE;
	pstate->comment_string           = comment_string;
	pstate->comment_string_length    = comment_string == NULL ? 0 : strlen(comment_string);

	plrec_reader->pvstate     = (void*)pstate;
	plrec_reader->popen_func  = file_reader_mmap_vopen;
	plrec_reader->pclose_func = file_reader_mmap_vclose;

	if (streq(irs, "auto")) {
		// Auto means either lines end in "\n" or "\r\n" (LF or CRLF).  In
		// either case the final character is "\n". Then for autodetect we
		// simply check if there's a character in the line before the '\n', and
		// if that is '\r'.
		pstate->do_auto_line_term = TRUE;
		pstate->irs = "\n";
		pstate->irslen = 1;
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_nidx_process_single_irs_single_ifs
			: lrec_reader_mmap_nidx_process_single_irs_multi_ifs;
	} else if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_nidx_process_single_irs_single_ifs
			: lrec_reader_mmap_nidx_process_single_irs_multi_ifs;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_nidx_process_multi_irs_single_ifs
			: lrec_reader_mmap_nidx_process_multi_irs_multi_ifs;
	}

	plrec_reader->psof_func     = lrec_reader_mmap_nidx_sof;
	plrec_reader->pfree_func    = lrec_reader_mmap_nidx_free;

	return plrec_reader;
}

static void lrec_reader_mmap_nidx_free(lrec_reader_t* preader) {
	free(preader->pvstate);
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_nidx_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_single_irs_single_ifs(phandle, pstate->irs[0], pstate->ifs[0], pstate, pctx);
}

static lrec_t* lrec_reader_mmap_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_single_irs_multi_ifs(phandle, pstate->irs[0], pstate, pctx);
}

static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_multi_irs_single_ifs(phandle, pstate->ifs[0], pstate);
}

static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_multi_irs_multi_ifs(phandle, pstate);
}

// ----------------------------------------------------------------
static lrec_t* lrec_parse_mmap_nidx_single_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char irs, char ifs, lrec_reader_mmap_nidx_state_t* pstate, context_t* pctx)
{
	// Skip comment lines
	if (pstate->comment_string != NULL) {
		while ((phandle->eof - phandle->sol) >= pstate->comment_string_length
		&& streqn(phandle->sol, pstate->comment_string, pstate->comment_string_length))
		{
			phandle->sol += pstate->comment_string_length;
			while (phandle->sol < phandle->eof && *phandle->sol != irs) {
				phandle->sol++;
			}
			if (phandle->sol < phandle->eof && *phandle->sol == irs) {
				phandle->sol++;
			}
		}
	}
	if (phandle->sol >= phandle->eof)
		return NULL;

	char* line  = phandle->sol;
	lrec_t* prec = lrec_unbacked_alloc();

	int idx = 0;
	char free_flags = NO_FREE;

	char* p = line;
	if (pstate->allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;

			if (pstate->do_auto_line_term) {
				if (p > line && p[-1] == '\r') {
					p[-1] = 0;
					context_set_autodetected_crlf(pctx);
				} else {
					context_set_autodetected_lf(pctx);
				}
			}

			phandle->sol = p+1;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p++;
			if (pstate->allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			value = p;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (pstate->allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(idx, &free_flags);

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, free_flags);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, free_flags|FREE_ENTRY_VALUE);
	}

	return prec;
}

static lrec_t* lrec_parse_mmap_nidx_single_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char irs, lrec_reader_mmap_nidx_state_t* pstate, context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* ifs = pstate->ifs;
	int ifslen = pstate->ifslen;

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = NO_FREE;

	char* p = line;
	if (pstate->allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;

			if (pstate->do_auto_line_term) {
				if (p > line && p[-1] == '\r') {
					p[-1] = 0;
					context_set_autodetected_crlf(pctx);
				} else {
					context_set_autodetected_lf(pctx);
				}
			}

			phandle->sol = p+1;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p += ifslen;
			if (pstate->allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			value = p;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (pstate->allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(idx, &free_flags);

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, free_flags);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, free_flags|FREE_ENTRY_VALUE);
	}

	return prec;
}

static lrec_t* lrec_parse_mmap_nidx_multi_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char ifs, lrec_reader_mmap_nidx_state_t* pstate)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = NO_FREE;

	char* p = line;
	if (pstate->allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;

	char* irs = pstate->irs;
	int irslen = pstate->irslen;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p++;
			if (pstate->allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			value = p;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (pstate->allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(idx, &free_flags);

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, free_flags);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, free_flags|FREE_ENTRY_VALUE);
	}

	return prec;
}

static lrec_t* lrec_parse_mmap_nidx_multi_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	lrec_reader_mmap_nidx_state_t* pstate)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = NO_FREE;

	char* ifs = pstate->ifs;
	int ifslen = pstate->ifslen;
	char* irs = pstate->irs;
	int irslen = pstate->irslen;

	char* p = line;
	if (pstate->allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;

			idx++;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);

			p += ifslen;
			if (pstate->allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			value = p;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (pstate->allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(idx, &free_flags);

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, free_flags);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, free_flags|FREE_ENTRY_VALUE);
	}

	return prec;
}
