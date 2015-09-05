#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/old_peek_file_reader.h"

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
#define TERMIND_RS  0x1111
#define TERMIND_FS  0x2222
#define TERMIND_EOF 0x3333

typedef struct _field_wrapper_t {
	char* contents;
	int   termind;
} field_wrapper_t;

typedef struct _record_wrapper_t {
	slls_t* contents;
	int   at_eof;
} record_wrapper_t;

// ----------------------------------------------------------------
typedef struct _lrec_reader_stdio_csv_state_t {
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
	old_peek_file_reader_t* pfr;

	int                 expect_header_line_next;
	header_keeper_t*    pheader_keeper;
	lhmslv_t*           pheader_keepers;

} lrec_reader_stdio_csv_state_t;

// ----------------------------------------------------------------
static record_wrapper_t lrec_reader_stdio_csv_get_record(lrec_reader_stdio_csv_state_t* pstate);

static field_wrapper_t get_csv_field(lrec_reader_stdio_csv_state_t* pstate);
static field_wrapper_t get_csv_field_not_dquoted(lrec_reader_stdio_csv_state_t* pstate);
static field_wrapper_t get_csv_field_dquoted(lrec_reader_stdio_csv_state_t* pstate);
static lrec_t*         paste_header_and_data(lrec_reader_stdio_csv_state_t* pstate, slls_t* pdata_fields);

// ----------------------------------------------------------------
// xxx needs abend on null lhs. etc.

static lrec_t* lrec_reader_stdio_csv_process(void* pvhandle, void* pvstate, context_t* pctx) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	if (pstate->pfr == NULL) {
		pstate->pfr = pfr_alloc((FILE*)pvhandle, pstate->peek_buf_len);
	}

	record_wrapper_t rwrapper;

	if (pstate->expect_header_line_next) {
		rwrapper = lrec_reader_stdio_csv_get_record(pstate);

		if (rwrapper.contents == NULL && rwrapper.at_eof)
			return NULL;
		pstate->ilno++;

		pstate->expect_header_line_next = FALSE;

		pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, rwrapper.contents);
		if (pstate->pheader_keeper == NULL) {
			pstate->pheader_keeper = header_keeper_alloc(NULL, rwrapper.contents);
			lhmslv_put(pstate->pheader_keepers, rwrapper.contents, pstate->pheader_keeper);
		} else { // Re-use the header-keeper in the header cache
			slls_free(rwrapper.contents);
		}
	}

	rwrapper = lrec_reader_stdio_csv_get_record(pstate);
	if (rwrapper.contents == NULL && rwrapper.at_eof)
		return NULL;

	pstate->ilno++;
	return paste_header_and_data(pstate, rwrapper.contents);
}

static record_wrapper_t lrec_reader_stdio_csv_get_record(lrec_reader_stdio_csv_state_t* pstate) {
	slls_t* pfields = slls_alloc();
	record_wrapper_t rwrapper;
	rwrapper.contents = pfields;
	rwrapper.at_eof = FALSE;
	while (TRUE) {
		field_wrapper_t fwrapper = get_csv_field(pstate);
		if (fwrapper.termind == TERMIND_EOF)
			rwrapper.at_eof = TRUE;
		if (fwrapper.contents != NULL)
			slls_add_with_free(pfields, fwrapper.contents);
		if (fwrapper.termind != TERMIND_FS)
			break;
	}
	if (pfields->length == 0 && rwrapper.at_eof) {
		slls_free(pfields);
		rwrapper.contents = NULL;
	}
	return rwrapper;
}

static field_wrapper_t get_csv_field(lrec_reader_stdio_csv_state_t* pstate) {
	field_wrapper_t wrapper;
	if (old_pfr_at_eof(pstate->pfr)) {
		wrapper.contents = NULL;
		wrapper.termind = TERMIND_EOF;
		return wrapper;
	} else if (old_pfr_next_is(pstate->pfr, pstate->dquote, pstate->dquote_len)) {
		old_pfr_advance_by(pstate->pfr, pstate->dquote_len);
		return get_csv_field_dquoted(pstate);
	} else {
		return get_csv_field_not_dquoted(pstate);
	}
}

static field_wrapper_t get_csv_field_not_dquoted(lrec_reader_stdio_csv_state_t* pstate) {
	while (TRUE) {
		if (old_pfr_at_eof(pstate->pfr)) {
			return (field_wrapper_t) {
				.contents = sb_is_empty(pstate->psb) ? NULL: sb_finish(pstate->psb),
				.termind = TERMIND_EOF
			};
		} else if (old_pfr_next_is(pstate->pfr, pstate->ifs_eof, pstate->ifs_eof_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->ifs_eof_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
		} else if (old_pfr_next_is(pstate->pfr, pstate->ifs, pstate->ifs_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->ifs_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
		} else if (old_pfr_next_is(pstate->pfr, pstate->irs, pstate->irs_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->irs_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
		} else if (old_pfr_next_is(pstate->pfr, pstate->dquote, pstate->dquote_len)) {
			fprintf(stderr, "%s: non-compliant field-internal double-quote at line %lld.\n",
				MLR_GLOBALS.argv0, pstate->ilno);
			exit(1);
		} else {
			sb_append_char(pstate->psb, old_pfr_read_char(pstate->pfr));
		}
	}
}

static field_wrapper_t get_csv_field_dquoted(lrec_reader_stdio_csv_state_t* pstate) {
	while (TRUE) {
		if (old_pfr_at_eof(pstate->pfr)) {
			fprintf(stderr, "%s: imbalanced double-quote at line %lld.\n", MLR_GLOBALS.argv0, pstate->ilno);
			exit(1);
		} else if (old_pfr_next_is(pstate->pfr, pstate->dquote_eof, pstate->dquote_eof_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->dquote_eof_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
		} else if (old_pfr_next_is(pstate->pfr, pstate->dquote_ifs, pstate->dquote_ifs_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->dquote_ifs_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
		} else if (old_pfr_next_is(pstate->pfr, pstate->dquote_irs, pstate->dquote_irs_len)) {
			old_pfr_advance_by(pstate->pfr, pstate->dquote_irs_len);
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
		} else if (old_pfr_next_is(pstate->pfr, pstate->dquote_dquote, pstate->dquote_dquote_len)) {
			// "" inside a dquoted field is an escape for "
			old_pfr_advance_by(pstate->pfr, pstate->dquote_dquote_len);
			sb_append_string(pstate->psb, pstate->dquote);
		} else {
			sb_append_char(pstate->psb, old_pfr_read_char(pstate->pfr));
		}
	}
}

static lrec_t* paste_header_and_data(lrec_reader_stdio_csv_state_t* pstate, slls_t* pdata_fields) {
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
static void lrec_reader_stdio_csv_sof(void* pvstate) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = TRUE;
	pstate->pfr = NULL;
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_csv_free(void* pvstate) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	old_pfr_free(pstate->pfr);
}

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_csv_alloc(char irs, char ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_csv_state_t));
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
	pstate->pfr                       = NULL;

	pstate->expect_header_line_next   = TRUE;
	pstate->pheader_keeper            = NULL;
	pstate->pheader_keepers           = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_stdio_vopen;
	plrec_reader->pclose_func   = &file_reader_stdio_vclose;
	plrec_reader->pprocess_func = &lrec_reader_stdio_csv_process;
	plrec_reader->psof_func     = &lrec_reader_stdio_csv_sof;
	plrec_reader->pfree_func    = &lrec_reader_stdio_csv_free;

	return plrec_reader;
}
