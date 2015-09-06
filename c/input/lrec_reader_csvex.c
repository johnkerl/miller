#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/string_builder.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_stdio.h"
#include "input/byte_reader.h"
#include "input/lrec_readers.h"
#include "input/peek_file_reader.h"

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

// ----------------------------------------------------------------
typedef struct _lrec_reader_csvex_state_t {
	// Input line number is not the same as the record-counter in context_t,
	// which counts records.
	long long  ilno;

	char* irs;
	char* ifs;
	char* dquote_irs;
	char* dquote_ifs;
	char* dquote_eof;
	char* dquote;
	char* dquote_dquote;
	char* ifs_eof;

	int   irs_len;
	int   ifs_len;
	int   dquote_irs_len;
	int   dquote_ifs_len;
	int   dquote_eof_len;
	int   dquote_len;
	int   dquote_dquote_len;
	int   ifs_eof_len;

	int   peek_buf_len;

	string_builder_t    sb;
	string_builder_t*   psb;
	byte_reader_t*      pbr;
	peek_file_reader_t* pfr;

	int                 expect_header_line_next;
	header_keeper_t*    pheader_keeper;
	lhmslv_t*           pheader_keepers;

} lrec_reader_csvex_state_t;

static slls_t* lrec_reader_csvex_get_fields(lrec_reader_csvex_state_t* pstate);
static lrec_t* paste_header_and_data(lrec_reader_csvex_state_t* pstate, slls_t* pdata_fields);

// ----------------------------------------------------------------
// xxx needs abend on null lhs. etc.

static lrec_t* lrec_reader_csvex_process(void* pvhandle, void* pvstate, context_t* pctx) {
	lrec_reader_csvex_state_t* pstate = pvstate;

//	xxx byte-reader open ...
//	if (pstate->pfr == NULL) {
//		pstate->pfr = pfr_alloc((FILE*)pvhandle, pstate->peek_buf_len);
//	}

	if (pstate->expect_header_line_next) {
		slls_t* pheader_fields = lrec_reader_csvex_get_fields(pstate);
		if (pheader_fields == NULL)
			return NULL;
		pstate->ilno++;
		pstate->expect_header_line_next = FALSE;

		pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
		if (pstate->pheader_keeper == NULL) {
			pstate->pheader_keeper = header_keeper_alloc(NULL, pheader_fields);
			lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper);
		} else { // Re-use the header-keeper in the header cache
			slls_free(pheader_fields);
		}
	}
	pstate->ilno++;

	slls_t* pdata_fields = lrec_reader_csvex_get_fields(pstate);
	return paste_header_and_data(pstate, pdata_fields);
}

static slls_t* lrec_reader_csvex_get_fields(lrec_reader_csvex_state_t* pstate) {

//	if @ eof: return null

//	while (TRUE) { // loop over fields in record

//		if (peek char is dquote) {
//			advance past the dquote

//			while (TRUE) { // loop over characters in field
//				pfr peek up to maxlen
//				rc = parse_trie_match(pstate->pdquote_parse_trie, xxx buf, ...);
//				if (rc) {
//					switch(stridx) {
//					case DQUOTE_EOF:
//						...; end of record
//						break;
//					case DQUOTE_IFS:
//						...; end of field
//						break;
//					case DQUOTE_IRS:
//						...; end of record
//						break;
//					case DQUOTE_DQUOTE:
//						...; sb append char '"'
//						break;
//					default:
//						...; sb append char of pfr->read_char
//						break;
//					}
//				}
//			}

//		} else {

//			while (TRUE) { // loop over characters in field
//				pfr peek up to maxlen
//				rc = parse_trie_match(pstate->pno_dquote_parse_trie, xxx buf, ...);
//				if (rc) {
//					switch(stridx) {
//					case EOF:
//						...;
//						break;
//					case IFS:
//						...;
//						break;
//					case IRS:
//						...;
//						break;
//					case DQUOTE:
//						...;
//						break;
//					default:
//						...;
//						break;
//					}
//				}
//			}

//		}
//	}

	return NULL; // xxx stub
}

//static field_wrapper_t get_csvex_field(lrec_reader_csvex_state_t* pstate) {
//	field_wrapper_t wrapper;
//	if (pfr_at_eof(pstate->pfr)) {
//		wrapper.contents = NULL;
//		wrapper.termind = TERMIND_EOF;
//		return wrapper;
//	} else if (pfr_next_is(pstate->pfr, pstate->dquote, pstate->dquote_len)) {
//		pfr_advance_by(pstate->pfr, pstate->dquote_len);
//		return get_csvex_field_dquoted(pstate);
//	} else {
//		return get_csvex_field_not_dquoted(pstate);
//	}
//}

//static field_wrapper_t get_csvex_field_not_dquoted(lrec_reader_csvex_state_t* pstate) {
//	while (TRUE) {
//		if (pfr_at_eof(pstate->pfr)) {
//			return (field_wrapper_t) {
//				.contents = sb_is_empty(pstate->psb) ? NULL: sb_finish(pstate->psb),
//				.termind = TERMIND_EOF
//			};
//		} else if (pfr_next_is(pstate->pfr, pstate->ifs_eof, pstate->ifs_eof_len)) {
//			pfr_advance_by(pstate->pfr, pstate->ifs_eof_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
//		} else if (pfr_next_is(pstate->pfr, pstate->ifs, pstate->ifs_len)) {
//			pfr_advance_by(pstate->pfr, pstate->ifs_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
//		} else if (pfr_next_is(pstate->pfr, pstate->irs, pstate->irs_len)) {
//			pfr_advance_by(pstate->pfr, pstate->irs_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
//		} else if (pfr_next_is(pstate->pfr, pstate->dquote, pstate->dquote_len)) {
//			fprintf(stderr, "%s: non-compliant field-internal double-quote at line %lld.\n",
//				MLR_GLOBALS.argv0, pstate->ilno);
//			exit(1);
//		} else {
//			sb_append_char(pstate->psb, pfr_read_char(pstate->pfr));
//		}
//	}
//}

//static field_wrapper_t get_csvex_field_dquoted(lrec_reader_csvex_state_t* pstate) {
//	while (TRUE) {
//		if (pfr_at_eof(pstate->pfr)) {
//			fprintf(stderr, "%s: imbalanced double-quote at line %lld.\n", MLR_GLOBALS.argv0, pstate->ilno);
//			exit(1);
//		} else if (pfr_next_is(pstate->pfr, pstate->dquote_eof, pstate->dquote_eof_len)) {
//			pfr_advance_by(pstate->pfr, pstate->dquote_eof_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
//		} else if (pfr_next_is(pstate->pfr, pstate->dquote_ifs, pstate->dquote_ifs_len)) {
//			pfr_advance_by(pstate->pfr, pstate->dquote_ifs_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
//		} else if (pfr_next_is(pstate->pfr, pstate->dquote_irs, pstate->dquote_irs_len)) {
//			pfr_advance_by(pstate->pfr, pstate->dquote_irs_len);
//			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
//		} else if (pfr_next_is(pstate->pfr, pstate->dquote_dquote, pstate->dquote_dquote_len)) {
//			// "" inside a dquoted field is an escape for "
//			pfr_advance_by(pstate->pfr, pstate->dquote_dquote_len);
//			sb_append_string(pstate->psb, pstate->dquote);
//		} else {
//			sb_append_char(pstate->psb, pfr_read_char(pstate->pfr));
//		}
//	}
//}

static lrec_t* paste_header_and_data(lrec_reader_csvex_state_t* pstate, slls_t* pdata_fields) {
	if (pstate->pheader_keeper->pkeys->length != pdata_fields->length) {
		fprintf(stderr, "%s: Header/data length mismatch: %d != %d at line %lld.\n",
			MLR_GLOBALS.argv0, pstate->pheader_keeper->pkeys->length, pdata_fields->length, pstate->ilno);
		exit(1);
	}
	lrec_t* prec = lrec_unbacked_alloc();
	sllse_t* ph = pstate->pheader_keeper->pkeys->phead;
	sllse_t* pd = pdata_fields->phead;
	for ( ; ph != NULL && pd != NULL; ph = ph->pnext, pd = pd->pnext) {
		lrec_put_no_free(prec, ph->value, pd->value);
	}
	return prec;
}

// ----------------------------------------------------------------
static void lrec_reader_csvex_sof(void* pvstate) {
	lrec_reader_csvex_state_t* pstate = pvstate;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = TRUE;
	pstate->pfr = NULL;
}

// ----------------------------------------------------------------
static void lrec_reader_csvex_free(void* pvstate) {
	lrec_reader_csvex_state_t* pstate = pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	pfr_free(pstate->pfr);
}

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_csvex_alloc(byte_reader_t* pbr, char irs, char ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_csvex_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_csvex_state_t));
	pstate->ilno                      = 0LL;
	pstate->irs                       = "\r\n"; // xxx multi-byte the cli irs/ifs/etc, and integrate here
	pstate->ifs                       = ",";    // xxx multi-byte the cli irs/ifs/etc, and integrate here

	pstate->dquote_irs                = mlr_paste_2_strings("\"", pstate->irs);
	pstate->dquote_ifs                = mlr_paste_2_strings("\"", pstate->ifs);
	pstate->dquote_eof                = "\"\xff";
	pstate->dquote                    = "\"";
	pstate->dquote_dquote             = "\"\"";
	pstate->ifs_eof                   = mlr_paste_2_strings(pstate->ifs, "\xff");

	pstate->irs_len                   = strlen(pstate->irs);
	pstate->ifs_len                   = strlen(pstate->ifs);
	pstate->dquote_irs_len            = strlen(pstate->dquote_irs);
	pstate->dquote_ifs_len            = strlen(pstate->dquote_ifs);
	pstate->dquote_eof_len            = strlen(pstate->dquote_eof);
	pstate->dquote_len                = strlen(pstate->dquote);
	pstate->dquote_dquote_len         = strlen(pstate->dquote_dquote);
	pstate->ifs_eof_len               = strlen(pstate->ifs_eof);

	pstate->peek_buf_len              = pstate->irs_len;
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->ifs_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->dquote_irs_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->dquote_ifs_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->dquote_eof_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->dquote_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->dquote_dquote_len);
	pstate->peek_buf_len              = mlr_imax2(pstate->peek_buf_len, pstate->ifs_eof_len);
	pstate->peek_buf_len             += 2;

	sb_init(&pstate->sb, STRING_BUILDER_INIT_SIZE);
	pstate->psb                       = &pstate->sb;
	pstate->pbr                       = pbr;
	pstate->pfr                       = NULL;

	// xxx allocate the parse-tries here -- one for dquote only,
	// the second for non-dquote-after-that, the third for dquoted-after-that.

	pstate->expect_header_line_next   = TRUE;
	pstate->pheader_keeper            = NULL;
	pstate->pheader_keepers           = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_stdio_vopen;
	plrec_reader->pclose_func   = &file_reader_stdio_vclose;
	plrec_reader->pprocess_func = &lrec_reader_csvex_process;
	plrec_reader->psof_func     = &lrec_reader_csvex_sof;
	plrec_reader->pfree_func    = &lrec_reader_csvex_free;

	return plrec_reader;
}
