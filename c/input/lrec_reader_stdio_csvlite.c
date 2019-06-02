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
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

// Idea of pheader_keepers: each header_keeper object retains the input-line backing
// and the slls_t for a CSV header line which is used by one or more CSV data
// lines.  Meanwhile some mappers retain input records from the entire data
// stream, including header-schema changes in the input stream. This means we
// need to keep headers intact as long as any lrecs are pointing to them.  One
// option is reference-counting which I experimented with; it was messy and
// error-prone. The approach used here is to keep a hash map from header-schema
// to header_keeper object. The current pheader_keeper is a pointer into one of
// those.  Then when the reader is freed, all the header-keepers are freed.

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

typedef struct _lrec_reader_stdio_csvlite_state_t {
	long long  ifnr;
	long long  ilno; // Line-level, not record-level as in context_t
	char*  irs;
	char*  ifs;
	int    irslen;
	int    ifslen;
	int    allow_repeat_ifs;
	int    do_auto_line_term;
	int    use_implicit_csv_header;
	int    allow_ragged_csv_input;
	size_t line_length;
	comment_handling_t comment_handling;
	char*  comment_string;

	int  expect_header_line_next;
	header_keeper_t* pheader_keeper;
	lhmslv_t*     pheader_keepers;
} lrec_reader_stdio_csvlite_state_t;

static void    lrec_reader_stdio_csvlite_free(lrec_reader_t* preader);
static void    lrec_reader_stdio_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_stdio_csvlite_process(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_csvlite_alloc(char* irs, char* ifs, int allow_repeat_ifs, int use_implicit_csv_header,
	int allow_ragged_csv_input, comment_handling_t comment_handling, char* comment_string)
{
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_csvlite_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_csvlite_state_t));
	pstate->ifnr                    = 0LL;
	pstate->irs                     = irs;
	pstate->ifs                     = ifs;
	pstate->irslen                  = strlen(irs);
	pstate->ifslen                  = strlen(ifs);
	pstate->allow_repeat_ifs        = allow_repeat_ifs;
	pstate->do_auto_line_term       = FALSE;
	pstate->use_implicit_csv_header = use_implicit_csv_header;
	pstate->allow_ragged_csv_input  = allow_ragged_csv_input;
	// This is used to track nominal line length over the file read. Bootstrap with a default length.
	pstate->line_length             = MLR_ALLOC_READ_LINE_INITIAL_SIZE;
	pstate->comment_handling        = comment_handling;
	pstate->comment_string          = comment_string;

	pstate->expect_header_line_next = use_implicit_csv_header  ? FALSE : TRUE;
	pstate->pheader_keeper          = NULL;
	pstate->pheader_keepers         = lhmslv_alloc();

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
		pstate->do_auto_line_term = TRUE;
	}
	plrec_reader->pprocess_func = lrec_reader_stdio_csvlite_process;
	plrec_reader->psof_func     = lrec_reader_stdio_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_csvlite_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_csvlite_free(lrec_reader_t* preader) {
	lrec_reader_stdio_csvlite_state_t* pstate = preader->pvstate;
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		header_keeper_t* pheader_keeper = pe->pvvalue;
		header_keeper_free(pheader_keeper);
	}
	lhmslv_free(pstate->pheader_keepers);
	free(pstate);
	free(preader);
}

// ----------------------------------------------------------------
static void lrec_reader_stdio_sof(void* pvstate, void* pvhandle) {
	lrec_reader_stdio_csvlite_state_t* pstate = pvstate;
	pstate->ifnr = 0LL;
	pstate->ilno = 0LL;
	pstate->expect_header_line_next = pstate->use_implicit_csv_header ? FALSE : TRUE;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_csvlite_process(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_csvlite_state_t* pstate = pvstate;

	while (TRUE) {
		if (pstate->expect_header_line_next) {
			while (TRUE) {
				char* hline = NULL;
				if (pstate->comment_handling == COMMENTS_ARE_DATA) {
					if (pstate->irslen == 1)
						hline = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
							&pstate->line_length, pstate->do_auto_line_term, pctx);
					else
						hline = mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
							&pstate->line_length);
					if (hline != NULL)
						pstate->ilno++;
				} else {
					int num_lines_comment_skipped = 0;
					if (pstate->irslen == 1)
						hline = mlr_alloc_read_line_single_delimiter_stripping_comments_aux(input_stream, pstate->irs[0],
							&pstate->line_length, pstate->do_auto_line_term,
							pstate->comment_handling, pstate->comment_string, &num_lines_comment_skipped, pctx);
					else
						hline = mlr_alloc_read_line_multiple_delimiter_stripping_comments_aux(input_stream,
							pstate->irs, pstate->irslen, &pstate->line_length,
							pstate->comment_handling, pstate->comment_string, &num_lines_comment_skipped);
					pstate->ilno += num_lines_comment_skipped;
					if (hline != NULL)
						pstate->ilno++;
				}

				if (hline == NULL) // EOF
					return NULL;

				slls_t* pheader_fields = (pstate->ifslen == 1)
					? split_csvlite_header_line_single_ifs(hline, pstate->ifs[0], pstate->allow_repeat_ifs)
					: split_csvlite_header_line_multi_ifs(hline, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
				if (pheader_fields->length == 0) {
					pstate->expect_header_line_next = TRUE;
					if (pstate->pheader_keeper != NULL) {
						pstate->pheader_keeper = NULL;
					}
				} else {
					for (sllse_t* pe = pheader_fields->phead; pe != NULL; pe = pe->pnext) {
						if (*pe->value == 0) {
							fprintf(stderr, "%s: unacceptable empty CSV key at file \"%s\" line %lld.\n",
								MLR_GLOBALS.bargv0, pctx->filename, pstate->ilno);
							exit(1);
						}
					}

					pstate->expect_header_line_next = FALSE;

					pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
					if (pstate->pheader_keeper == NULL) {
						pstate->pheader_keeper = header_keeper_alloc(hline, pheader_fields);
						lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper,
							NO_FREE); // freed by header-keeper
					} else { // Re-use the header-keeper in the header cache
						slls_free(pheader_fields);
						free(hline);
					}
					break;
				}
			}
		}

		char* line = NULL;

		if (pstate->comment_handling == COMMENTS_ARE_DATA) {
			if (pstate->irslen == 1)
				line = mlr_alloc_read_line_single_delimiter(input_stream, pstate->irs[0],
					&pstate->line_length, pstate->do_auto_line_term, pctx);
			else
				line = mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->irs, pstate->irslen,
					&pstate->line_length);
		} else {
			if (pstate->irslen == 1)
				line = mlr_alloc_read_line_single_delimiter_stripping_comments(input_stream, pstate->irs[0],
					&pstate->line_length, pstate->do_auto_line_term,
					pstate->comment_handling, pstate->comment_string, pctx);
			else
				line = mlr_alloc_read_line_multiple_delimiter_stripping_comments(input_stream,
					pstate->irs, pstate->irslen, &pstate->line_length,
					pstate->comment_handling, pstate->comment_string);
		}

		if (line == NULL) // EOF
			return NULL;
		pstate->ilno++;

		if (!*line) {
			if (pstate->pheader_keeper != NULL) {
				pstate->pheader_keeper = NULL;
				pstate->expect_header_line_next = TRUE;
				free(line);
				continue;
			}
		} else {
			pstate->ifnr++;
			if (pstate->ifslen == 1) {
				return pstate->use_implicit_csv_header
					? lrec_parse_stdio_csvlite_data_line_single_ifs_implicit_header(
						pstate->pheader_keeper, pctx->filename, pstate->ilno, line,
						pstate->ifs[0], pstate->allow_repeat_ifs)
					: lrec_parse_stdio_csvlite_data_line_single_ifs(pstate->pheader_keeper, pctx->filename,
						pstate->ilno, line, pstate->ifs[0], pstate->allow_repeat_ifs, pstate->allow_ragged_csv_input);
			} else {
				return pstate->use_implicit_csv_header
					? lrec_parse_stdio_csvlite_data_line_multi_ifs_implicit_header(
						pstate->pheader_keeper, pctx->filename, pstate->ilno, line,
						pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs)
					: lrec_parse_stdio_csvlite_data_line_multi_ifs(pstate->pheader_keeper, pctx->filename,
						pstate->ilno, line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs,
						pstate->allow_ragged_csv_input);
			}
		}
	}
}

// ----------------------------------------------------------------
slls_t* split_csvlite_header_line_single_ifs(char* line, char ifs, int allow_repeat_ifs) {
	slls_t* plist = slls_alloc();
	if (*line == 0) // empty string splits to empty list
		return plist;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* start = p;
	for ( ; *p; p++) {
		if (*p == ifs) {
			*p = 0;
			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			slls_append_no_free(plist, start);
			start = p;
		}
	}
	if (allow_repeat_ifs && *start == 0) {
		; // OK
	} else {
		slls_append_no_free(plist, start);
	}

	return plist;
}

slls_t* split_csvlite_header_line_multi_ifs(char* line, char* ifs, int ifslen, int allow_repeat_ifs) {
	slls_t* plist = slls_alloc();
	if (*line == 0) // empty string splits to empty list
		return plist;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* start = p;
	for ( ; *p; p++) {
		if (streqn(p, ifs, ifslen)) {
			*p = 0;
			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			slls_append_no_free(plist, start);
			start = p;
		}
	}
	if (allow_repeat_ifs && *start == 0) {
		; // OK
	} else {
		slls_append_no_free(plist, start);
	}

	return plist;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_csvlite_data_line_single_ifs(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char ifs, int allow_repeat_ifs, int allow_ragged_csv_input)
{
	lrec_t* prec = lrec_csvlite_alloc(data_line);
	char* p = data_line;

	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	int idx = 0;

	sllse_t* pe = pheader_keeper->pkeys->phead;
	for ( ; *p; ) {
		if (*p == ifs) {
			*p = 0;
			idx++;
			if (pe == NULL) {
				// Data line has more fields than the header line did
				if (allow_ragged_csv_input) {
					char free_flags = NO_FREE;
					key = low_int_to_string(idx, &free_flags);
					lrec_put(prec, key, value, free_flags);
				} else{
					fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
						MLR_GLOBALS.bargv0, filename, ilno);
					exit(1);
				}
			} else {
				key = pe->value;
				pe = pe->pnext;
				lrec_put(prec, key, value, NO_FREE);
			}
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
	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else if (pe == NULL) {
		// Data line has more fields than the header line did
		if (allow_ragged_csv_input) {
			idx++;
			char free_flags = NO_FREE;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);
		} else{
			fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
				MLR_GLOBALS.bargv0, filename, ilno);
			exit(1);
		}
	} else {
		key = pe->value;
		lrec_put(prec, key, value, NO_FREE);
		if (pe->pnext != NULL) {
			// Header line has more fields than the data line did
			if (allow_ragged_csv_input) {
				for (pe = pe->pnext ; pe != NULL; pe = pe->pnext) {
					key = pe->value;
					lrec_put(prec, key, "", NO_FREE);
				}
			} else {
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.bargv0, filename, ilno);
				exit(1);
			}
		}
	}

	return prec;
}

lrec_t* lrec_parse_stdio_csvlite_data_line_multi_ifs(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char* ifs, int ifslen, int allow_repeat_ifs, int allow_ragged_csv_input)
{
	lrec_t* prec = lrec_csvlite_alloc(data_line);
	char* p = data_line;

	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	int idx = 0;

	sllse_t* pe = pheader_keeper->pkeys->phead;
	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
			*p = 0;
			idx++;
			if (pe == NULL) {
				// Data line has more fields than the header line did
				if (allow_ragged_csv_input) {
					char free_flags = NO_FREE;
					key = low_int_to_string(idx, &free_flags);
					lrec_put(prec, key, value, free_flags);
				} else{
					fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
						MLR_GLOBALS.bargv0, filename, ilno);
					exit(1);
				}
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
	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else if (pe == NULL) {
		// Data line has more fields than the header line did
		if (allow_ragged_csv_input) {
			char free_flags = NO_FREE;
			key = low_int_to_string(idx, &free_flags);
			lrec_put(prec, key, value, free_flags);
		} else{
			fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
				MLR_GLOBALS.bargv0, filename, ilno);
			exit(1);
		}
	} else {
		key = pe->value;
		lrec_put(prec, key, value, NO_FREE);
		if (pe->pnext != NULL) {
			// Header line has more fields than the data line did
			if (allow_ragged_csv_input) {
				for (pe = pe->pnext ; pe != NULL; pe = pe->pnext) {
					key = pe->value;
					lrec_put(prec, key, "", NO_FREE);
				}
			} else {
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.bargv0, filename, ilno);
				exit(1);
			}
		}
	}

	return prec;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_csvlite_data_line_single_ifs_implicit_header(header_keeper_t* pheader_keeper, char* filename, long long ilno,
	char* data_line, char ifs, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_csvlite_alloc(data_line);
	char* p = data_line;

	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char free_flags = NO_FREE;
	char* value = p;

	int idx = 0;
	for ( ; *p; ) {
		if (*p == ifs) {
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
	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = low_int_to_string(++idx, &free_flags);
		lrec_put(prec, key, value, NO_FREE);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

lrec_t* lrec_parse_stdio_csvlite_data_line_multi_ifs_implicit_header(header_keeper_t* pheader_keeper, char* filename,
	long long ilno, char* data_line, char* ifs, int ifslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_csvlite_alloc(data_line);
	char* p = data_line;

	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	char free_flags;

	int idx = 0;
	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
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
	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = low_int_to_string(++idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}
