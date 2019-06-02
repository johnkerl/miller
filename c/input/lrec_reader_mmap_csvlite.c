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
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

// ----------------------------------------------------------------
// Multi-file cases:
//
// a,a        a,b        c          d
// -- FILE1:  -- FILE1:  -- FILE1:  -- FILE1:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// -- FILE2:  -- FILE2:
// a,b,c      d,e,f,g    a,b,c      d,e,f
// 7,8,9      3,4,5,6    7,8,9      3,4,5
// --OUTPUT:  --OUTPUT:  --OUTPUT:  --OUTPUT:
// a,b,c      a,b,c      a,b,c      a,b,c
// 1,2,3      1,2,3      1,2,3      1,2,3
// 4,5,6      4,5,6      4,5,6      4,5,6
// 7,8,9                 7,8,9
//            d,e,f,g               d,e,f
//            3,4,5,6               3,4,5
// ----------------------------------------------------------------

typedef struct _lrec_reader_mmap_csvlite_state_t {
	long long  ifnr;
	long long  ilno; // Line-level, not record-level as in context_t
	char* irs;
	char* ifs;
	int   irslen;
	int   ifslen;
	int   allow_repeat_ifs;
	int   do_auto_line_term;
	int   use_implicit_header;
	comment_handling_t comment_handling;
	char* comment_string;
	int   comment_string_length;

	int  expect_header_line_next;
	header_keeper_t* pheader_keeper;
	lhmslv_t*     pheader_keepers;
} lrec_reader_mmap_csvlite_state_t;

static void    lrec_reader_mmap_csvlite_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_csvlite_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_csvlite_process_single_seps(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_csvlite_process_multi_seps(void* pvstate, void* pvhandle, context_t* pctx);

static slls_t* lrec_reader_mmap_csvlite_get_header_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx);

static slls_t* lrec_reader_mmap_csvlite_get_header_multi_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate);

static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza);

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza);

static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza);

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza);

static int handle_comment_line_single_irs(
	file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate,
	char irs);

static int handle_comment_line_multi_irs(
	file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_csvlite_alloc(char* irs, char* ifs, int allow_repeat_ifs, int use_implicit_header,
	comment_handling_t comment_handling, char* comment_string)
{
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_csvlite_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_csvlite_state_t));
	pstate->ifnr                     = 0LL;
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->irslen                   = strlen(irs);
	pstate->ifslen                   = strlen(ifs);
	pstate->allow_repeat_ifs         = allow_repeat_ifs;
	pstate->do_auto_line_term        = FALSE;
	pstate->use_implicit_header      = use_implicit_header;
	pstate->comment_handling         = comment_handling;
	pstate->comment_string           = comment_string;
	pstate->comment_string_length    = comment_string == NULL ? 0 : strlen(comment_string);

	pstate->expect_header_line_next  = use_implicit_header ? FALSE : TRUE;
	pstate->pheader_keeper           = NULL;
	pstate->pheader_keepers          = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;

	if (streq(irs, "auto")) {
		// Auto means either lines end in "\n" or "\r\n" (LF or CRLF).  In
		// either case the final character is "\n". Then for autodetect we
		// simply check if there's a character in the line before the '\n', and
		// if that is '\r'.
		pstate->do_auto_line_term = TRUE;
		pstate->irs = "\n";
		pstate->irslen = 1;
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_csvlite_process_single_seps
			: lrec_reader_mmap_csvlite_process_multi_seps;
	} else {
		plrec_reader->pprocess_func = (pstate->irslen == 1 && pstate->ifslen == 1)
			? lrec_reader_mmap_csvlite_process_single_seps
			: lrec_reader_mmap_csvlite_process_multi_seps;
	}

	plrec_reader->psof_func     = lrec_reader_mmap_csvlite_sof;
	plrec_reader->pfree_func    = lrec_reader_mmap_csvlite_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
static void lrec_reader_mmap_csvlite_free(lrec_reader_t* preader) {
	lrec_reader_mmap_csvlite_state_t* pstate = preader->pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	lhmslv_free(pstate->pheader_keepers);
	free(pstate);
	free(preader);
}

static void lrec_reader_mmap_csvlite_sof(void* pvstate, void* pvhandle) {
	lrec_reader_mmap_csvlite_state_t* pstate = pvstate;
	pstate->ifnr = 0LL;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = pstate->use_implicit_header ? FALSE : TRUE;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_csvlite_process_single_seps(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_csvlite_state_t* pstate = pvstate;

	while (TRUE) {
		if (pstate->expect_header_line_next) {

			slls_t* pheader_fields = lrec_reader_mmap_csvlite_get_header_single_seps(phandle, pstate, pctx);
			if (pheader_fields == NULL) { // EOF
				return NULL;
			}

			for (sllse_t* pe = pheader_fields->phead; pe != NULL; pe = pe->pnext) {
				if (*pe->value == 0) {
					fprintf(stderr, "%s: unacceptable empty CSV key at file \"%s\" line %lld.\n",
						MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
					exit(1);
				}
			}

			pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
			if (pstate->pheader_keeper == NULL) {
				pstate->pheader_keeper = header_keeper_alloc(NULL, pheader_fields);
				lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper,
					NO_FREE); // freed by header-keeper
			} else { // Re-use the header-keeper in the header cache
				slls_free(pheader_fields);
			}
			pstate->expect_header_line_next = FALSE;
		}

		int end_of_stanza = FALSE;
		lrec_t* prec = pstate->use_implicit_header
			?  lrec_reader_mmap_csvlite_get_record_single_seps_implicit_header(phandle, pstate, pctx,
				pstate->pheader_keeper, &end_of_stanza)
			:  lrec_reader_mmap_csvlite_get_record_single_seps(phandle, pstate, pctx,
				pstate->pheader_keeper, &end_of_stanza);
		if (end_of_stanza) {
			pstate->expect_header_line_next = TRUE;
		} else if (prec == NULL) { // EOF
			return NULL;
		} else {
			return prec;
		}
	}
}

static lrec_t* lrec_reader_mmap_csvlite_process_multi_seps(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_csvlite_state_t* pstate = pvstate;

	while (TRUE) {
		if (pstate->expect_header_line_next) {

			slls_t* pheader_fields = lrec_reader_mmap_csvlite_get_header_multi_seps(phandle, pstate);
			if (pheader_fields == NULL) // EOF
				return NULL;

			for (sllse_t* pe = pheader_fields->phead; pe != NULL; pe = pe->pnext) {
				if (*pe->value == 0) {
					fprintf(stderr, "%s: unacceptable empty CSV key at file \"%s\" line %lld.\n",
						MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
					exit(1);
				}
			}

			pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
			if (pstate->pheader_keeper == NULL) {
				pstate->pheader_keeper = header_keeper_alloc(NULL, pheader_fields);
				lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper,
					NO_FREE); // freed by header-keeper
			} else { // Re-use the header-keeper in the header cache
				slls_free(pheader_fields);
			}
			pstate->expect_header_line_next = FALSE;
		}

		int end_of_stanza = FALSE;
		lrec_t* prec = pstate->use_implicit_header
			? lrec_reader_mmap_csvlite_get_record_multi_seps_implicit_header(phandle, pstate, pctx,
				pstate->pheader_keeper, &end_of_stanza)
			: lrec_reader_mmap_csvlite_get_record_multi_seps(phandle, pstate, pctx,
				pstate->pheader_keeper, &end_of_stanza);
		if (end_of_stanza) {
			pstate->expect_header_line_next = TRUE;
		} else if (prec == NULL) { // EOF
			return NULL;
		} else {
			return prec;
		}
	}
}

// ----------------------------------------------------------------
static slls_t* lrec_reader_mmap_csvlite_get_header_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx)
{
	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	slls_t* pheader_names = slls_alloc();

	// Skip blank/comment lines and seek to header line
	while (TRUE) {
		if (phandle->sol < phandle->eof && *phandle->sol == irs) {
			phandle->sol++;
			pstate->ilno++;
			continue;
		}
		if (pstate->comment_string != NULL && handle_comment_line_single_irs(phandle, pstate, irs)) {
			continue;
		}
		break;
	}

	char* p = phandle->sol;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* osol = p;
	char* header_name = p;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;

			if (pstate->do_auto_line_term) {
				if (p > phandle->sol && p[-1] == '\r') {
					p[-1] = 0;
					context_set_autodetected_crlf(pctx);
				} else {
					context_set_autodetected_lf(pctx);
				}
			}

			phandle->sol = p+1;
			pstate->ilno++;
			break;
		} else if (*p == ifs) {
			*p = 0;

			slls_append_no_free(pheader_names, header_name);

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			header_name = p;
		} else {
			p++;
		}
	}
	if (allow_repeat_ifs && *header_name == 0) {
		// OK
	} else if (p == osol) {
		// OK
	} else {
		slls_append_no_free(pheader_names, header_name);
	}

	return pheader_names;
}

static slls_t* lrec_reader_mmap_csvlite_get_header_multi_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate)
{
	char* irs    = pstate->irs;
	char* ifs    = pstate->ifs;
	int   irslen = pstate->irslen;
	int   ifslen = pstate->ifslen;
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	// Skip blank/comment lines and seek to header line
	while (TRUE) {
		if ((phandle->eof - phandle->sol) >= irslen && streqn(phandle->sol, irs, irslen)) {
			phandle->sol += irslen;
			pstate->ilno++;
			continue;
		}
		if (pstate->comment_string != NULL && handle_comment_line_multi_irs(phandle, pstate)) {
			continue;
		}
		break;
	}

	slls_t* pheader_names = slls_alloc();

	// Parse the header line
	char* p = phandle->sol;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* osol = p;
	char* header_name = p;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			pstate->ilno++;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;

			slls_append_no_free(pheader_names, header_name);

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			header_name = p;
		} else {
			p++;
		}
	}
	if (allow_repeat_ifs && *header_name == 0) {
		// OK
	} else if (p == osol) {
		// OK
	} else {
		slls_append_no_free(pheader_names, header_name);
	}

	return pheader_names;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	// Skip comment lines
	if (pstate->comment_string != NULL) {
		while (handle_comment_line_single_irs(phandle, pstate, irs))
			;
	}

	if (phandle->sol >= phandle->eof)
		return NULL;

	char* line  = phandle->sol;
	lrec_t* prec = lrec_unbacked_alloc();

	sllse_t* pe = pheader_keeper->pkeys->phead;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
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
			pstate->ilno++;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			*p = 0;
			if (pe == NULL) {
				// Data line has more fields than the header line did
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
				exit(1);
			}
			key = pe->value;
			pe = pe->pnext;
			lrec_put(prec, key, value, NO_FREE);

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
	if (p >= phandle->eof)
		phandle->sol = p+1;

	if (allow_repeat_ifs && *value == 0)
		return prec;

	if (pe == NULL) {
		// Data line has more fields than the header line did
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
		exit(1);
	}
	key = pe->value;

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, NO_FREE);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
	}

	if (pe->pnext != NULL) {
		// Header line has more fields than the data line did
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
		exit(1);
	}

	return prec;
}

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	// Skip comment lines
	if (pstate->comment_string != NULL) {
		while (handle_comment_line_multi_irs(phandle, pstate))
			;
	}
	if (phandle->sol >= phandle->eof)
		return NULL;

	char* irs = pstate->irs;
	char* ifs = pstate->ifs;
	int   irslen = pstate->irslen;
	int   ifslen = pstate->ifslen;
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	lrec_t* prec = lrec_unbacked_alloc();
	char* line  = phandle->sol;

	sllse_t* pe = pheader_keeper->pkeys->phead;
	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
			*p = 0;
			phandle->sol = p + irslen;
			pstate->ilno++;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;
			if (pe == NULL) {
				// Data line has more fields than the header line did
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
				exit(1);
			}
			key = pe->value;
			pe = pe->pnext;
			lrec_put(prec, key, value, NO_FREE);

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
	if (p >= phandle->eof)
		phandle->sol = p+1;

	if (allow_repeat_ifs && *value == 0)
		return prec;

	if (pe == NULL) {
		// Data line has more fields than the header line did
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
		exit(1);
	}
	key = pe->value;

	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		lrec_put(prec, key, value, NO_FREE);
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
		lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
	}

	if (pe->pnext != NULL) {
		// Header line has more fields than the data line did
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
		exit(1);
	}

	return prec;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	// Skip comment lines
	if (pstate->comment_string != NULL) {
		while (handle_comment_line_single_irs(phandle, pstate, irs))
			;
	}
	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();
	char* line  = phandle->sol;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	char  free_flags = NO_FREE;
	int idx = 0;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
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
			pstate->ilno++;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			*p = 0;
			key = low_int_to_string(++idx, &free_flags);
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
	if (p >= phandle->eof)
		phandle->sol = p+1;

	if (allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(++idx, &free_flags);

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

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	// Skip comment lines
	if (pstate->comment_string != NULL) {
		while (handle_comment_line_multi_irs(phandle, pstate))
			;
	}
	if (phandle->sol >= phandle->eof)
		return NULL;

	char* irs = pstate->irs;
	char* ifs = pstate->ifs;
	int   irslen = pstate->irslen;
	int   ifslen = pstate->ifslen;
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	lrec_t* prec = lrec_unbacked_alloc();
	char* line  = phandle->sol;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	char free_flags;
	int idx = 0;
	int saw_rs = FALSE;
	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
			*p = 0;
			phandle->sol = p + irslen;
			pstate->ilno++;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;
			key = low_int_to_string(++idx, &free_flags);
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
	if (p >= phandle->eof)
		phandle->sol = p+1;

	if (allow_repeat_ifs && *value == 0)
		return prec;

	key = low_int_to_string(++idx, &free_flags);

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

// ----------------------------------------------------------------
static int handle_comment_line_single_irs(
	file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate,
	char irs)
{
	if ((phandle->eof - phandle->sol) >= pstate->comment_string_length
	&& streqn(phandle->sol, pstate->comment_string, pstate->comment_string_length))
	{
		if (pstate->comment_handling == PASS_COMMENTS)
			for (int i = 0; i < pstate->comment_string_length; i++)
				fputc(phandle->sol[i], stdout);
		phandle->sol += pstate->comment_string_length;
		while (phandle->sol < phandle->eof && *phandle->sol != irs) {
			if (pstate->comment_handling == PASS_COMMENTS)
				fputc(*phandle->sol, stdout);
			phandle->sol++;
		}
		if (phandle->sol < phandle->eof && *phandle->sol == irs) {
			if (pstate->comment_handling == PASS_COMMENTS)
				fputc(*phandle->sol, stdout);
			phandle->sol++;
		}
		pstate->ilno++;
		return TRUE;
	} else {
		return FALSE;
	}
}

// ----------------------------------------------------------------
static int handle_comment_line_multi_irs(
	file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate)
{
	if ((phandle->eof - phandle->sol) >= pstate->comment_string_length
	&& streqn(phandle->sol, pstate->comment_string, pstate->comment_string_length))
	{
		if (pstate->comment_handling == PASS_COMMENTS)
			for (int i = 0; i < pstate->comment_string_length; i++)
				fputc(phandle->sol[i], stdout);
		phandle->sol += pstate->comment_string_length;
		while ((phandle->eof - phandle->sol >= pstate->irslen) && !streqn(phandle->sol, pstate->irs, pstate->irslen)) {
			if (pstate->comment_handling == PASS_COMMENTS)
				fputc(*phandle->sol, stdout);
			phandle->sol++;
		}
		if ((phandle->eof - phandle->sol >= pstate->irslen) && streqn(phandle->sol, pstate->irs, pstate->irslen)) {
			if (pstate->comment_handling == PASS_COMMENTS)
				for (int i = 0; i < pstate->irslen; i++)
					fputc(phandle->sol[i], stdout);
			phandle->sol += pstate->irslen;
		}
		pstate->ilno++;
		return TRUE;
	} else {
		return FALSE;
	}
}
