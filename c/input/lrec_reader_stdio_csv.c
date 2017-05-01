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
#include <ctype.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"
#include "input/file_reader_stdio.h"
#include "input/byte_readers.h"
#include "input/lrec_readers.h"
#include "input/peek_file_reader.h"
#include "containers/rslls.h"
#include "containers/lhmslv.h"
#include "containers/parse_trie.h"

// Idea of pheader_keepers: each header_keeper object retains the input-line backing
// and the slls_t for a CSV header line which is used by one or more CSV data
// lines.  Meanwhile some mappers retain input records from the entire data
// stream, including header-schema changes in the input stream. This means we
// need to keep headers intact as long as any lrecs are pointing to them.  One
// option is reference-counting which I experimented with; it was messy and
// error-prone. The approach used here is to keep a hash map from header-schema
// to header_keeper object. The current pheader_keeper is a pointer into one of
// those.  Then when the reader is freed, all the header-keepers are freed.

// ----------------------------------------------------------------
#define STRING_BUILDER_INIT_SIZE 1024

// AKA "token"
#define EOF_STRIDX           0x2000
#define IRS_STRIDX           0x2001
#define IFS_EOF_STRIDX       0x2002
#define IFS_STRIDX           0x2003
#define DQUOTE_STRIDX        0x2004
#define DQUOTE_IRS_STRIDX    0x2005
#define DQUOTE_IRS2_STRIDX   0x2006 // alternate line-ending for autodetect LF/CRLF
#define DQUOTE_IFS_STRIDX    0x2007
#define DQUOTE_EOF_STRIDX    0x2008
#define DQUOTE_DQUOTE_STRIDX 0x2009

// ----------------------------------------------------------------
typedef struct _lrec_reader_stdio_csv_state_t {
	// Input line number is not the same as the record-counter in context_t,
	// which counts records.
	long long  ilno;

	char* eof;
	char* irs;
	char* ifs_eof;
	char* ifs;
	int   do_auto_line_term;
	char* dquote;
	char* dquote_irs;
	char* dquote_irs2;
	char* dquote_ifs;
	char* dquote_eof;
	char* dquote_dquote;
	int   dquotelen;

	rslls_t*            pfields;
	string_builder_t*   psb;
	byte_reader_t*      pbr;
	peek_file_reader_t* pfr;

	parse_trie_t*       pno_dquote_parse_trie;
	parse_trie_t*       pdquote_parse_trie;

	int                 expect_header_line_next;
	int                 use_implicit_header;
	header_keeper_t*    pheader_keeper;
	lhmslv_t*           pheader_keepers;

} lrec_reader_stdio_csv_state_t;

static void    lrec_reader_stdio_csv_free(lrec_reader_t* preader);
static void    lrec_reader_stdio_csv_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_stdio_csv_process(void* pvstate, void* pvhandle, context_t* pctx);
static int     lrec_reader_stdio_csv_get_fields(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pfields,
	context_t* pctx);
static lrec_t* paste_indices_and_data(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pdata_fields,
	context_t* pctx);
static lrec_t* paste_header_and_data(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pdata_fields,
	context_t* pctx);
static void*   lrec_reader_stdio_csv_open(void* pvstate, char* prepipe, char* filename);
static void    lrec_reader_stdio_csv_close(void* pvstate, void* pvhandle, char* prepipe);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_csv_alloc(char* irs, char* ifs, int use_implicit_header) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_csv_state_t));
	pstate->ilno          = 0LL;

	pstate->do_auto_line_term = FALSE;
	if (streq(irs, "auto")) {
		irs = "\n";
		pstate->do_auto_line_term = TRUE;
	}

	pstate->eof           = "\xff";
	pstate->irs           = irs;
	pstate->ifs           = ifs;
	pstate->ifs_eof       = mlr_paste_2_strings(pstate->ifs, "\xff");
	pstate->dquote        = "\"";

	pstate->dquote_ifs    = mlr_paste_2_strings("\"", pstate->ifs);
	pstate->dquote_eof    = "\"\xff";
	pstate->dquote_dquote = "\"\"";

	pstate->dquotelen     = strlen(pstate->dquote);

	pstate->pno_dquote_parse_trie = parse_trie_alloc();
	parse_trie_add_string(pstate->pno_dquote_parse_trie, pstate->eof,     EOF_STRIDX);
	parse_trie_add_string(pstate->pno_dquote_parse_trie, pstate->irs,     IRS_STRIDX);
	parse_trie_add_string(pstate->pno_dquote_parse_trie, pstate->ifs_eof, IFS_EOF_STRIDX);
	parse_trie_add_string(pstate->pno_dquote_parse_trie, pstate->ifs,     IFS_STRIDX);
	parse_trie_add_string(pstate->pno_dquote_parse_trie, pstate->dquote,  DQUOTE_STRIDX);

	pstate->pdquote_parse_trie = parse_trie_alloc();
	if (pstate->do_auto_line_term) {
		pstate->dquote_irs  = mlr_paste_2_strings("\"", "\n");
		pstate->dquote_irs2 = mlr_paste_2_strings("\"", "\r\n");
		parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_irs,  DQUOTE_IRS_STRIDX);
		parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_irs2, DQUOTE_IRS2_STRIDX);
	} else {
		pstate->dquote_irs  = mlr_paste_2_strings("\"", pstate->irs);
		pstate->dquote_irs2 = NULL;
		parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_irs, DQUOTE_IRS_STRIDX);
	}
	parse_trie_add_string(pstate->pdquote_parse_trie, pstate->eof,           EOF_STRIDX);
	parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_irs,    DQUOTE_IRS_STRIDX);
	parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_ifs,    DQUOTE_IFS_STRIDX);
	parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_eof,    DQUOTE_EOF_STRIDX);
	parse_trie_add_string(pstate->pdquote_parse_trie, pstate->dquote_dquote, DQUOTE_DQUOTE_STRIDX);

	pstate->pfields = rslls_alloc();
	pstate->psb = sb_alloc(STRING_BUILDER_INIT_SIZE);
	pstate->pbr = stdio_byte_reader_alloc();
	pstate->pfr = pfr_alloc(pstate->pbr, mlr_imax2(pstate->pno_dquote_parse_trie->maxlen,
		pstate->pdquote_parse_trie->maxlen));

	pstate->expect_header_line_next   = use_implicit_header ? FALSE : TRUE;
	pstate->use_implicit_header       = use_implicit_header;
	pstate->pheader_keeper            = NULL;
	pstate->pheader_keepers           = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = lrec_reader_stdio_csv_open;
	plrec_reader->pclose_func   = lrec_reader_stdio_csv_close;
	plrec_reader->pprocess_func = lrec_reader_stdio_csv_process;
	plrec_reader->psof_func     = lrec_reader_stdio_csv_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_csv_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_csv_free(lrec_reader_t* preader) {
	lrec_reader_stdio_csv_state_t* pstate = preader->pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	lhmslv_free(pstate->pheader_keepers);
	pfr_free(pstate->pfr);
	parse_trie_free(pstate->pno_dquote_parse_trie);
	parse_trie_free(pstate->pdquote_parse_trie);
	rslls_free(pstate->pfields);
	stdio_byte_reader_free(pstate->pbr);
	sb_free(pstate->psb);
	free(pstate->ifs_eof);
	free(pstate->dquote_irs);
	free(pstate->dquote_irs2);
	free(pstate->dquote_ifs);
	free(pstate);
	free(preader);
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_csv_sof(void* pvstate, void* pvhandle) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = pstate->use_implicit_header ? FALSE : TRUE;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_csv_process(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;

	if (pstate->expect_header_line_next) {
		if (!lrec_reader_stdio_csv_get_fields(pstate, pstate->pfields, pctx))
			return NULL;
		pstate->ilno++;

		slls_t* pheader_fields = slls_alloc();
		int i = 0;
		for (rsllse_t* pe = pstate->pfields->phead; i < pstate->pfields->length && pe != NULL; pe = pe->pnext) {
			if (*pe->value == 0) {
				fprintf(stderr, "%s: unacceptable empty CSV key at file \"%s\" line %lld.\n",
					MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
				exit(1);
			}
			// Transfer pointer-free responsibility from the rslls to the
			// header fields in the header keeper
			if (string_starts_with(pe->value, "\xef\xbb\xbf")) {
				// Strip UTF-8 BOM if any
				slls_append(pheader_fields, mlr_strdup_or_die(&pe->value[3]), FREE_ENTRY_VALUE);
			} else {
				slls_append(pheader_fields, pe->value, pe->free_flag);
			}
			pe->free_flag = 0;
		}
		rslls_reset(pstate->pfields);

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
	int rc = lrec_reader_stdio_csv_get_fields(pstate, pstate->pfields, pctx);
	pstate->ilno++;
	if (rc == FALSE) // EOF
		return NULL;
	else {
		lrec_t* prec =  pstate->use_implicit_header
			? paste_indices_and_data(pstate, pstate->pfields, pctx)
			: paste_header_and_data(pstate, pstate->pfields, pctx);
		rslls_reset(pstate->pfields);
		return prec;
	}
}

static int lrec_reader_stdio_csv_get_fields(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pfields,
	context_t* pctx)
{
	int rc, stridx, matchlen, record_done, field_done;
	peek_file_reader_t* pfr = pstate->pfr;
	string_builder_t*   psb = pstate->psb;
	char* field = NULL;
	int field_length = 0;

	if (pfr_peek_char(pfr) == (char)EOF) // char defaults to unsigned on some platforms
		return FALSE;

	// loop over fields in record
	record_done = FALSE;
	while (!record_done) {
		// Assumption is dquote is "\""
		if (pfr_peek_char(pfr) != pstate->dquote[0]) {

			// Loop over characters in field
			field_done = FALSE;
			while (!field_done) {
				pfr_buffer_by(pfr, pstate->pno_dquote_parse_trie->maxlen);

				rc = parse_trie_ring_match(pstate->pno_dquote_parse_trie,
					pfr->peekbuf, pfr->sob, pfr->npeeked, pfr->peekbuflenmask,
					&stridx, &matchlen);
#ifdef DEBUG_PARSER
				pfr_print(pfr);
#endif
				if (rc) {
#ifdef DEBUG_PARSER
					printf("RC=%d stridx=0x%04x matchlen=%d\n", rc, stridx, matchlen);
#endif
					switch(stridx) {
					case EOF_STRIDX: // end of record
						rslls_append(pfields, sb_finish(psb), FREE_ENTRY_VALUE, 0);
						field_done  = TRUE;
						record_done = TRUE;
						break;
					case IFS_EOF_STRIDX:
						fprintf(stderr, "%s: syntax error: record-ending field separator at line %lld.\n",
							MLR_GLOBALS.bargv0, pstate->ilno);
						exit(1);
						break;
					case IFS_STRIDX: // end of field
						rslls_append(pfields, sb_finish(psb), FREE_ENTRY_VALUE, 0);
						field_done  = TRUE;
						break;
					case IRS_STRIDX: // end of record
						field = sb_finish_with_length(psb, &field_length);

						// The line-ending '\n' won't be included in the field buffer.
						if (pstate->do_auto_line_term) {
							if (field_length > 0 && field[field_length-1] == '\r') {
								field[field_length-1] = 0;
								context_set_autodetected_crlf(pctx);
							} else {
								context_set_autodetected_lf(pctx);
							}
						}

						rslls_append(pfields, field, FREE_ENTRY_VALUE, 0);
						field_done  = TRUE;
						record_done = TRUE;
						break;
					case DQUOTE_STRIDX: // CSV syntax error: fields containing quotes must be fully wrapped in quotes
						fprintf(stderr, "%s: syntax error: unwrapped double quote at line %lld.\n",
							MLR_GLOBALS.bargv0, pstate->ilno);
						exit(1);
						break;
					default:
						fprintf(stderr, "%s: internal coding error: unexpected token %d at line %lld.\n",
							MLR_GLOBALS.bargv0, stridx, pstate->ilno);
						exit(1);
						break;
					}
					pfr_advance_by(pfr, matchlen);
				} else {
#ifdef DEBUG_PARSER
					char c = pfr_read_char(pfr);
					printf("CHAR=%c [%02x]\n", isprint((unsigned char)c) ? c : ' ', (unsigned)c);
					sb_append_char(psb, c);
#else
					sb_append_char(psb, pfr_read_char(pfr));
#endif
				}
			}

		} else {
			pfr_advance_by(pfr, pstate->dquotelen);

			// loop over characters in field
			field_done = FALSE;
			char* field = NULL;
			int field_length = 0;
			while (!field_done) {
				pfr_buffer_by(pfr, pstate->pdquote_parse_trie->maxlen);

				rc = parse_trie_ring_match(pstate->pdquote_parse_trie,
					pfr->peekbuf, pfr->sob, pfr->npeeked, pfr->peekbuflenmask,
					&stridx, &matchlen);

				if (rc) {
					switch(stridx) {
					case EOF_STRIDX: // end of record
						fprintf(stderr, "%s: unmatched double quote at line %lld.\n",
							MLR_GLOBALS.bargv0, pstate->ilno);
						exit(1);
						break;
					case DQUOTE_EOF_STRIDX: // end of record
						rslls_append(pfields, sb_finish(psb), FREE_ENTRY_VALUE, FIELD_QUOTED_ON_INPUT);
						field_done  = TRUE;
						record_done = TRUE;
						break;
					case DQUOTE_IFS_STRIDX: // end of field
						rslls_append(pfields, sb_finish(psb), FREE_ENTRY_VALUE, FIELD_QUOTED_ON_INPUT);
						field_done  = TRUE;
						break;
					case DQUOTE_IRS_STRIDX: // end of record
					case DQUOTE_IRS2_STRIDX: // end of record

						field = sb_finish_with_length(psb, &field_length);

						// The line-ending '\n' won't be included in the field buffer.
						if (pstate->do_auto_line_term) {
							if (field_length > 0 && field[field_length-1] == '\r') {
								field[field_length-1] = 0;
								context_set_autodetected_crlf(pctx);
							} else {
								context_set_autodetected_lf(pctx);
							}
						}

						rslls_append(pfields, field, FREE_ENTRY_VALUE, FIELD_QUOTED_ON_INPUT);
						field_done  = TRUE;
						record_done = TRUE;
						break;
					case DQUOTE_DQUOTE_STRIDX: // RFC-4180 CSV: "" inside a dquoted field is an escape for "
						sb_append_char(psb, pstate->dquote[0]);
						break;
					default:
						fprintf(stderr, "%s: internal coding error: unexpected token %d at line %lld.\n",
							MLR_GLOBALS.bargv0, stridx, pstate->ilno);
						exit(1);
						break;
					}
					pfr_advance_by(pfr, matchlen);
				} else {
					sb_append_char(psb, pfr_read_char(pfr));
				}
			}

		}
	}

	return TRUE;
}

// ----------------------------------------------------------------
static lrec_t* paste_indices_and_data(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pdata_fields,
	context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();
	int idx = 0;
	for (rsllse_t* pd = pdata_fields->phead; pd != NULL; pd = pd->pnext) {
		idx++;
		char free_flags = pd->free_flag;
		char* key = low_int_to_string(idx, &free_flags);
		// Transfer pointer-free responsibility from the rslls to the lrec object
		lrec_put_ext(prec, key, pd->value, free_flags, pd->quote_flag);
		pd->free_flag = 0;
	}
	return prec;
}

// ----------------------------------------------------------------
static lrec_t* paste_header_and_data(lrec_reader_stdio_csv_state_t* pstate, rslls_t* pdata_fields,
	context_t* pctx)
{
	if (pstate->pheader_keeper->pkeys->length != pdata_fields->length) {
		fprintf(stderr, "%s: Header/data length mismatch (%llu != %llu) at file \"%s\" line %lld.\n",
			MLR_GLOBALS.bargv0, pstate->pheader_keeper->pkeys->length, pdata_fields->length,
			pctx->filename, pstate->ilno);
		exit(1);
	}
	lrec_t* prec = lrec_unbacked_alloc();
	sllse_t* ph = pstate->pheader_keeper->pkeys->phead;
	rsllse_t* pd = pdata_fields->phead;
	for ( ; ph != NULL && pd != NULL; ph = ph->pnext, pd = pd->pnext) {
		// Transfer pointer-free responsibility from the rslls to the lrec object
		lrec_put_ext(prec, ph->value, pd->value, pd->free_flag, pd->quote_flag);
		pd->free_flag = 0;
	}
	return prec;
}

// ----------------------------------------------------------------
static void* lrec_reader_stdio_csv_open(void* pvstate, char* prepipe, char* filename) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	pstate->pfr->pbr->popen_func(pstate->pfr->pbr, prepipe, filename);
	pfr_reset(pstate->pfr);
	// Different from the other readers, we keep the file handle within the
	// byte_reader object.
	return NULL;
}

static void lrec_reader_stdio_csv_close(void* pvstate, void* pvhandle, char* prepipe) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	pstate->pfr->pbr->pclose_func(pstate->pfr->pbr, prepipe);
}
