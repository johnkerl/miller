#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

// xxx cmt/docxref re no irs for xtab, & two or more ifses separates records.
typedef struct _lrec_reader_mmap_xtab_state_t {
	char* ifs;
	char  ips;
	int   ifslen;
	int   allow_repeat_ips;
} lrec_reader_mmap_xtab_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_xtab_process_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_single_ifs(phandle, pstate->ifs[0], pstate->ips, pstate->allow_repeat_ips);
}

static lrec_t* lrec_reader_mmap_xtab_process_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab_multi_ifs(phandle, pstate->ifs, pstate->ips, pstate->ifslen, pstate->allow_repeat_ips);
}

static void lrec_reader_mmap_xtab_sof(void* pvstate) {
}

lrec_reader_t* lrec_reader_mmap_xtab_alloc(char* ifs, char ips, int allow_repeat_ips) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_xtab_state_t));
	pstate->ifs                 = ifs;
	pstate->ips                 = ips;
	pstate->ifslen              = strlen(pstate->ifs);
	pstate->allow_repeat_ips    = allow_repeat_ips;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_mmap_vopen;
	plrec_reader->pclose_func   = &file_reader_mmap_vclose;

	plrec_reader->pprocess_func = (pstate->ifslen == 1)
		? lrec_reader_mmap_xtab_process_single_ifs
		: lrec_reader_mmap_xtab_process_multi_ifs;

	plrec_reader->psof_func     = &lrec_reader_mmap_xtab_sof;
	plrec_reader->pfree_func    = NULL;

	return plrec_reader;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_xtab_single_ifs(file_reader_mmap_state_t* phandle, char ifs, char ips, int allow_repeat_ips) {

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
		for (p = line; p < phandle->eof && *p; ) {
			if (*p == ifs) {
				*p = 0;
				phandle->sol = p+1;
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

		lrec_put_no_free(prec, key, value);

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

lrec_t* lrec_parse_mmap_xtab_multi_ifs(file_reader_mmap_state_t* phandle, char* ifs, char ips, int ifslen, int allow_repeat_ips)
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
		for (p = line; p < phandle->eof && *p; ) {
			if (streqn(p, ifs, ifslen)) {
				*p = 0;
				phandle->sol = p + ifslen;
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

		lrec_put_no_free(prec, key, value);

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
