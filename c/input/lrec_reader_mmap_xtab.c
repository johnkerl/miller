#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_xtab_state_t {
	char* ifs;
	char* ips;
	int   ifslen;
	int   ipslen;
	int   allow_repeat_ips;
} lrec_reader_mmap_xtab_state_t;

static void    lrec_reader_mmap_xtab_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_xtab_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_xtab_process_single_ifs_single_ips(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_xtab_process_single_ifs_multi_ips(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_xtab_process_multi_ifs_single_ips(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_xtab_process_multi_ifs_multi_ips(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_xtab_alloc(char* ifs, char* ips, int allow_repeat_ips) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_xtab_state_t));
	pstate->ifs                 = ifs;
	pstate->ips                 = ips;
	pstate->ifslen              = strlen(pstate->ifs);
	pstate->ipslen              = strlen(pstate->ips);
	pstate->allow_repeat_ips    = allow_repeat_ips;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;

	if (pstate->ifslen == 1) {
		plrec_reader->pprocess_func = (pstate->ipslen == 1)
			? lrec_reader_mmap_xtab_process_single_ifs_single_ips
			: lrec_reader_mmap_xtab_process_single_ifs_multi_ips;
	} else {
		plrec_reader->pprocess_func = (pstate->ipslen == 1)
			? lrec_reader_mmap_xtab_process_multi_ifs_single_ips
			: lrec_reader_mmap_xtab_process_multi_ifs_multi_ips;
	}

	plrec_reader->psof_func     = lrec_reader_mmap_xtab_sof;
	plrec_reader->pfree_func    = lrec_reader_mmap_xtab_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
static void lrec_reader_mmap_xtab_free(lrec_reader_t* preader) {
	free(preader->pvstate);
	free(preader);
}

static void lrec_reader_mmap_xtab_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_xtab_process_single_ifs_single_ips(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_single_ifs_single_ips(phandle, pstate->ifs[0], pstate->ips[0],
			pstate->allow_repeat_ips);
}

static lrec_t* lrec_reader_mmap_xtab_process_single_ifs_multi_ips(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_single_ifs_multi_ips(phandle, pstate->ifs[0], pstate->ips, pstate->ipslen,
			pstate->allow_repeat_ips);
}

static lrec_t* lrec_reader_mmap_xtab_process_multi_ifs_single_ips(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_multi_ifs_single_ips(phandle, pstate->ifs, pstate->ips[0], pstate->ifslen,
			pstate->allow_repeat_ips);
}

static lrec_t* lrec_reader_mmap_xtab_process_multi_ifs_multi_ips(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_multi_ifs_multi_ips(phandle, pstate->ifs, pstate->ips, pstate->ifslen,
			pstate->ipslen, pstate->allow_repeat_ips);
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_xtab_single_ifs_single_ips(file_reader_mmap_state_t* phandle, char ifs, char ips,
	int allow_repeat_ips)
{
	while (phandle->sol < phandle->eof && *phandle->sol == ifs)
		phandle->sol++;

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* p;

		// Construct one field
		int saw_eol = FALSE;
		for (p = line; p < phandle->eof && *p; ) {
			if (*p == ifs) {
				*p = 0;
				phandle->sol = p+1;
				saw_eol = TRUE;
				break;
			} else if (*p == ips) {
				key = line;
				*p = 0;

				p++;
				if (allow_repeat_ips) {
					while (*p == ips)
						p++;
				}
				value = p;
			} else {
				p++;
			}
		}
		if (p >= phandle->eof)
			phandle->sol = p+1;

		if (saw_eol) {
			// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate
			// the C string so it's OK to retain a pointer to that.
			lrec_put(prec, key, value, NO_FREE);
		} else {
			// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null
			// character to terminate the C string: if the file size is not a multiple of the OS page size it'll work
			// (it's our copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at
			// EOF is one byte past the page and that will segv us.
			char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
			lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
		}

		if (phandle->sol >= phandle->eof || *phandle->sol == ifs)
			break;
	}
	if (prec->field_count == 0) {
		lrec_free(prec);
		return NULL;
	} else {
		return prec;
	}
}

lrec_t* lrec_parse_mmap_xtab_single_ifs_multi_ips(file_reader_mmap_state_t* phandle, char ifs, char* ips, int ipslen,
	int allow_repeat_ips)
{
	while (phandle->sol < phandle->eof && *phandle->sol == ifs)
		phandle->sol++;

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* p;

		// Construct one field
		int saw_eol = FALSE;
		for (p = line; p < phandle->eof && *p; ) {
			if (*p == ifs) {
				*p = 0;
				phandle->sol = p+1;
				saw_eol = TRUE;
				break;
			} else if (streqn(p, ips, ipslen)) {
				key = line;
				*p = 0;

				p += ipslen;
				if (allow_repeat_ips) {
					while (streqn(p, ips, ipslen))
						p += ipslen;
				}
				value = p;
			} else {
				p++;
			}
		}
		if (p >= phandle->eof)
			phandle->sol = p+1;

		if (saw_eol) {
			// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate
			// the C string so it's OK to retain a pointer to that.
			lrec_put(prec, key, value, NO_FREE);
		} else {
			// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null
			// character to terminate the C string: if the file size is not a multiple of the OS page size it'll work
			// (it's our copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at
			// EOF is one byte past the page and that will segv us.
			char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
			lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
		}

		if (phandle->sol >= phandle->eof || *phandle->sol == ifs)
			break;
	}
	if (prec->field_count == 0) {
		lrec_free(prec);
		return NULL;
	} else {
		return prec;
	}
}

lrec_t* lrec_parse_mmap_xtab_multi_ifs_single_ips(file_reader_mmap_state_t* phandle, char* ifs, char ips, int ifslen,
	int allow_repeat_ips)
{
	while (phandle->sol < phandle->eof && streqn(phandle->sol, ifs, ifslen))
		phandle->sol += ifslen;

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* p;

		// Construct one field
		int saw_eol = FALSE;
		for (p = line; p < phandle->eof && *p; ) {
			if (streqn(p, ifs, ifslen)) {
				*p = 0;
				phandle->sol = p + ifslen;
				saw_eol = TRUE;
				break;
			} else if (*p == ips) {
				key = line;
				*p = 0;

				p++;
				if (allow_repeat_ips) {
					while (*p == ips)
						p++;
				}
				value = p;
			} else {
				p++;
			}
		}
		if (p >= phandle->eof)
			phandle->sol = p+1;

		if (saw_eol) {
			// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate
			// the C string so it's OK to retain a pointer to that.
			lrec_put(prec, key, value, NO_FREE);
		} else {
			// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null
			// character to terminate the C string: if the file size is not a multiple of the OS page size it'll work
			// (it's our copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at
			// EOF is one byte past the page and that will segv us.
			char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
			lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
		}

		if (phandle->sol >= phandle->eof || streqn(phandle->sol, ifs, ifslen))
			break;
	}
	if (prec->field_count == 0) {
		lrec_free(prec);
		return NULL;
	} else {
		return prec;
	}
}

lrec_t* lrec_parse_mmap_xtab_multi_ifs_multi_ips(file_reader_mmap_state_t* phandle, char* ifs, char* ips,
	int ifslen, int ipslen, int allow_repeat_ips)
{
	while (phandle->sol < phandle->eof && streqn(phandle->sol, ifs, ifslen))
		phandle->sol += ifslen;

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* p;

		// Construct one field
		int saw_eol = FALSE;
		for (p = line; p < phandle->eof && *p; ) {
			if (streqn(p, ifs, ifslen)) {
				*p = 0;
				phandle->sol = p + ifslen;
				saw_eol = TRUE;
				break;
			} else if (streqn(p, ips, ipslen)) {
				key = line;
				*p = 0;

				p += ipslen;
				if (allow_repeat_ips) {
					while (streqn(p, ips, ipslen))
						p += ipslen;
				}
				value = p;
			} else {
				p++;
			}
		}
		if (p >= phandle->eof)
			phandle->sol = p+1;

		if (saw_eol) {
			// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate
			// the C string so it's OK to retain a pointer to that.
			lrec_put(prec, key, value, NO_FREE);
		} else {
			// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null
			// character to terminate the C string: if the file size is not a multiple of the OS page size it'll work
			// (it's our copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at
			// EOF is one byte past the page and that will segv us.
			char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
			lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
		}

		if (phandle->sol >= phandle->eof || streqn(phandle->sol, ifs, ifslen))
			break;
	}
	if (prec->field_count == 0) {
		lrec_free(prec);
		return NULL;
	} else {
		return prec;
	}
}
