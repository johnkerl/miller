#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "lib/string_builder.h"
#include "input/peek_file_reader.h"

// See https://github.com/johnkerl/miller/issues/4
// Temporary status:
// * --csv     from the command line maps into the (existing) csvlite I/O
// * --csvlite from the command line maps into the (existing) csvlite I/O
// * --csvex   from the command line maps into the (new & experimental & unadvertised) rfc-csv I/O
// Ultimate status:
// * --csvlite from the command line will maps into the csvlite I/O
// * --csv     from the command line will maps into the rfc-csv I/O

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
	long long  ifnr; // xxx cmt w/r/t pctx
	long long  ilno; // xxx cmt w/r/t pctx
	char* irs;
	char* ifs;
	// xxx parameterize dquote_irs
	// xxx parameterize dquote_ifs
	// xxx parameterize irs_len
	// xxx parameterize ifs_len
	// xxx parameterize dquote_irs_len
	// xxx parameterize dquote_ifs_len
	// xxx parameterize maxlen of all of those; for the pfr buf
	//int  allow_repeat_ifs;

	string_builder_t    sb;
	string_builder_t*   psb;
	peek_file_reader_t* pfr;

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
		pstate->pfr = pfr_alloc((FILE*)pvhandle, 32); // xxx set up via max of all terminds
	}

	record_wrapper_t rwrapper;
	while (TRUE) {
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

		pstate->ifnr++;
		return paste_header_and_data(pstate, rwrapper.contents);
	}
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
	if (pfr_at_eof(pstate->pfr)) {
		wrapper.contents = NULL;
		wrapper.termind = TERMIND_EOF;
		return wrapper;
	} else if (pfr_next_is(pstate->pfr, "\"", 1)) {
		return get_csv_field_dquoted(pstate);
	} else {
		return get_csv_field_not_dquoted(pstate);
	}
}

static field_wrapper_t get_csv_field_not_dquoted(lrec_reader_stdio_csv_state_t* pstate) {
	// xxx need pfr_advance_past_or_die ...
	// xxx "\"," etc. will be encoded in the rfc_csv_reader_t ctor -- this is just sketch
	while (TRUE) {
		if (pfr_at_eof(pstate->pfr)) {
			return (field_wrapper_t) {
				.contents = sb_is_empty(pstate->psb) ? NULL: sb_finish(pstate->psb),
				.termind = TERMIND_EOF
			};
		} else if (pfr_next_is(pstate->pfr, ",\xff", 2)) {
			if (!pfr_advance_past(pstate->pfr, ",\xff")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
		} else if (pfr_next_is(pstate->pfr, ",", 1)) {
			if (!pfr_advance_past(pstate->pfr, ",")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
		} else if (pfr_next_is(pstate->pfr, "\r\n", 2)) {
			if (!pfr_advance_past(pstate->pfr, "\r\n")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
		} else {
			sb_append_char(pstate->psb, pfr_read_char(pstate->pfr));
		}
	}
}

static field_wrapper_t get_csv_field_dquoted(lrec_reader_stdio_csv_state_t* pstate) {
	// xxx need pfr_advance_past_or_die ...
	if (!pfr_advance_past(pstate->pfr, "\"")) {
		fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
		exit(1);
	}
	while (TRUE) {
		if (pfr_at_eof(pstate->pfr)) {
			// xxx imbalanced-dquote error
			fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
			exit(1);
		} else if (pfr_next_is(pstate->pfr, "\"\xff", 2)) {
			if (!pfr_advance_past(pstate->pfr, "\"\xff")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_EOF };
		} else if (pfr_next_is(pstate->pfr, "\",", 2)) {
			if (!pfr_advance_past(pstate->pfr, "\",")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_FS };
		} else if (pfr_next_is(pstate->pfr, "\"\r\n", 3)) {
			if (!pfr_advance_past(pstate->pfr, "\"\r\n")) {
				fprintf(stderr, "xxx k0d3 me up b04k3n b04k3n b04ken %d\n", __LINE__);
				exit(1);
			}
			return (field_wrapper_t) { .contents = sb_finish(pstate->psb), .termind = TERMIND_RS };
		} else {
			sb_append_char(pstate->psb, pfr_read_char(pstate->pfr));
		}
	}
}

static lrec_t* paste_header_and_data(lrec_reader_stdio_csv_state_t* pstate, slls_t* pdata_fields) {
	if (pstate->pheader_keeper->pkeys->length != pdata_fields->length) {
		// xxx incorporate ctx/ilno/etc.
		fprintf(stderr, "Header/data length mismatch: %d != %d.\n",
			pstate->pheader_keeper->pkeys->length, pdata_fields->length);
		exit(1);
	}
	lrec_t* prec = lrec_unbacked_alloc();
	sllse_t* ph = pstate->pheader_keeper->pkeys->phead;
	sllse_t* pd = pdata_fields->phead;
	for ( ; ph != NULL && pd != NULL; ph = ph->pnext, pd = pd->pnext) {
		// xxx reduce the copies here
		lrec_put(prec, ph->value, strdup(pd->value), LREC_FREE_ENTRY_VALUE);
	}
	return prec;
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_sof(void* pvstate) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	pstate->ifnr = 0LL;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = TRUE;
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_csv_free(void* pvstate) {
	lrec_reader_stdio_csv_state_t* pstate = pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	pfr_free(pstate->pfr);
}

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_csv_alloc(char irs, char ifs, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_csv_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_csv_state_t));
	pstate->ifnr                      = 0LL;
	pstate->irs                       = "\r\n"; // xxx multi-byte the cli irs/ifs/etc, and integrate here
	pstate->ifs                       = ",";   // xxx multi-byte the cli irs/ifs/etc, and integrate here
	//pstate->allow_repeat_ifs          = allow_repeat_ifs;

	sb_init(&pstate->sb, 1024); // xxx #define at top
	pstate->psb                       = &pstate->sb;
	pstate->pfr                       = NULL;

	pstate->expect_header_line_next   = TRUE;
	pstate->pheader_keeper            = NULL;
	pstate->pheader_keepers           = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_stdio_vopen;
	plrec_reader->pclose_func   = &file_reader_stdio_vclose;
	plrec_reader->pprocess_func = &lrec_reader_stdio_csv_process;
	plrec_reader->psof_func     = &lrec_reader_stdio_sof;
	plrec_reader->pfree_func    = &lrec_reader_stdio_csv_free;

	return plrec_reader;
}
