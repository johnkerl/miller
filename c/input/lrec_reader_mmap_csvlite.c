#include <stdio.h>
#include <stdlib.h>
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
	int   use_implicit_header;

	int  expect_header_line_next;
	header_keeper_t* pheader_keeper;
	lhmslv_t*     pheader_keepers;
} lrec_reader_mmap_csvlite_state_t;

static void    lrec_reader_mmap_csvlite_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_csvlite_sof(void* pvstate);
static lrec_t* lrec_reader_mmap_csvlite_process_single_seps(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_csvlite_process_multi_seps(void* pvstate, void* pvhandle, context_t* pctx);

static slls_t* lrec_reader_mmap_csvlite_get_header_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate);

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


// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_csvlite_alloc(char* irs, char* ifs, int allow_repeat_ifs, int use_implicit_header) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_csvlite_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_csvlite_state_t));
	pstate->ifnr                     = 0LL;
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->irslen                   = strlen(irs);
	pstate->ifslen                   = strlen(ifs);
	pstate->allow_repeat_ifs         = allow_repeat_ifs;
	pstate->use_implicit_header      = use_implicit_header;
	pstate->expect_header_line_next  = use_implicit_header ? FALSE : TRUE;
	pstate->pheader_keeper           = NULL;
	pstate->pheader_keepers          = lhmslv_alloc();

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;

	plrec_reader->pprocess_func = (pstate->irslen == 1 && pstate->ifslen == 1)
		? lrec_reader_mmap_csvlite_process_single_seps
		: lrec_reader_mmap_csvlite_process_multi_seps;

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

	// header-fields lists are doubly referenced: as hashmap keys in
	// pstate->pheader_keepers, and within the header_keeper objects.
	// Nullify the keys here to avoid a double free.
	// xxx this could be refactored to be more elegant.
	for (lhmslve_t* pe = pstate->pheader_keepers->phead; pe != NULL; pe = pe->pnext) {
		pe->key = NULL;
	}
	lhmslv_free(pstate->pheader_keepers);

	free(pstate);
	free(preader);
}

static void lrec_reader_mmap_csvlite_sof(void* pvstate) {
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

			slls_t* pheader_fields = lrec_reader_mmap_csvlite_get_header_single_seps(phandle, pstate);
			if (pheader_fields == NULL) // EOF
				return NULL;

			for (sllse_t* pe = pheader_fields->phead; pe != NULL; pe = pe->pnext) {
				if (*pe->value == 0) {
					fprintf(stderr, "%s: unacceptable empty CSV key at file \"%s\" line %lld.\n",
						MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
					exit(1);
				}
			}

			pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
			if (pstate->pheader_keeper == NULL) {
				pstate->pheader_keeper = header_keeper_alloc(NULL, pheader_fields);
				lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper);
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
						MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
					exit(1);
				}
			}

			pstate->pheader_keeper = lhmslv_get(pstate->pheader_keepers, pheader_fields);
			if (pstate->pheader_keeper == NULL) {
				pstate->pheader_keeper = header_keeper_alloc(NULL, pheader_fields);
				lhmslv_put(pstate->pheader_keepers, pheader_fields, pstate->pheader_keeper);
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
	lrec_reader_mmap_csvlite_state_t* pstate)
{
	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	slls_t* pheader_names = slls_alloc();

	while (phandle->sol < phandle->eof && *phandle->sol == irs) {
		phandle->sol++;
		pstate->ilno++;
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
			phandle->sol = p+1;
			pstate->ilno++;
			break;
		} else if (*p == ifs) {
			*p = 0;

			slls_add_no_free(pheader_names, header_name);

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
		slls_add_no_free(pheader_names, header_name);
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

	slls_t* pheader_names = slls_alloc();

	while ((phandle->eof - phandle->sol) >= irslen && streqn(phandle->sol, irs, irslen)) {
		phandle->sol += irslen;
		pstate->ilno++;
	}

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

			slls_add_no_free(pheader_names, header_name);

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
		slls_add_no_free(pheader_names, header_name);
	}

	return pheader_names;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	if (phandle->sol >= phandle->eof)
		return NULL;

	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	sllse_t* pe = pheader_keeper->pkeys->phead;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
			*p = 0;
			phandle->sol = p+1;
			pstate->ilno++;
			break;
		} else if (*p == ifs) {
			*p = 0;
			if (pe == NULL) {
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
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

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else if (pe == NULL) {
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
		exit(1);
	} else {
		key = pe->value;
		lrec_put(prec, key, value, NO_FREE);
		if (pe->pnext != NULL) {
			fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
				MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
			exit(1);
		}
	}

	return prec;
}

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
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
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;
			if (pe == NULL) {
				fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
					MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
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

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else if (pe == NULL) {
		fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
			MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
		exit(1);
	} else {
		key = pe->value;
		lrec_put(prec, key, value, NO_FREE);
		if (pe->pnext != NULL) {
			fprintf(stderr, "%s: Header-data length mismatch in file %s at line %lld.\n",
				MLR_GLOBALS.argv0, pctx->filename, pstate->ilno);
			exit(1);
		}
	}

	return prec;
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_csvlite_get_record_single_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
	if (phandle->sol >= phandle->eof)
		return NULL;

	char irs = pstate->irs[0];
	char ifs = pstate->ifs[0];
	int allow_repeat_ifs = pstate->allow_repeat_ifs;

	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	char  free_flags;
	int idx = 0;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			if (p == line) {
				*pend_of_stanza = TRUE;
				lrec_free(prec);
				return NULL;
			}
			*p = 0;
			phandle->sol = p+1;
			pstate->ilno++;
			break;
		} else if (*p == ifs) {
			*p = 0;
			key = make_nidx_key(++idx, &free_flags);
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

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(++idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

static lrec_t* lrec_reader_mmap_csvlite_get_record_multi_seps_implicit_header(file_reader_mmap_state_t* phandle,
	lrec_reader_mmap_csvlite_state_t* pstate, context_t* pctx, header_keeper_t* pheader_keeper, int* pend_of_stanza)
{
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
	char  free_flags;
	int idx = 0;
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
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;
			key = make_nidx_key(++idx, &free_flags);
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

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(++idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}
