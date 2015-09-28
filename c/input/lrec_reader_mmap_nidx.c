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
} lrec_reader_mmap_nidx_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_single_irs_single_ifs(phandle, pstate->irs[0], pstate->ifs[0],
			pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_single_irs_multi_ifs(phandle, pstate->irs[0], pstate->ifs,
			pstate->ifslen, pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_multi_irs_single_ifs(phandle, pstate->irs, pstate->ifs[0],
			pstate->irslen, pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_nidx_multi_irs_multi_ifs(phandle, pstate->irs, pstate->ifs,
			pstate->irslen, pstate->ifslen, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_nidx_sof(void* pvstate) {
}

lrec_reader_t* lrec_reader_mmap_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_nidx_state_t));
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->irslen                   = strlen(pstate->irs);
	pstate->ifslen                   = strlen(pstate->ifs);
	pstate->allow_repeat_ifs         = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;

	if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_nidx_process_single_irs_single_ifs
			: lrec_reader_mmap_nidx_process_single_irs_multi_ifs;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? lrec_reader_mmap_nidx_process_multi_irs_single_ifs
			: lrec_reader_mmap_nidx_process_multi_irs_multi_ifs;
	}

	plrec_reader->psof_func     = lrec_reader_mmap_nidx_sof;
	plrec_reader->pfree_func    = NULL;

	return plrec_reader;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_nidx_single_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char irs, char ifs, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;
			phandle->sol = p+1;
			break;
		} else if (*p == ifs) {
			*p = 0;

			idx++;
			key = make_nidx_key(idx, &free_flags);
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

lrec_t* lrec_parse_mmap_nidx_single_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char irs, char* ifs, int ifslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;
			phandle->sol = p+1;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;

			idx++;
			key = make_nidx_key(idx, &free_flags);
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

lrec_t* lrec_parse_mmap_nidx_multi_irs_single_ifs(file_reader_mmap_state_t *phandle,
	char* irs, char ifs, int irslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			break;
		} else if (*p == ifs) {
			*p = 0;

			idx++;
			key = make_nidx_key(idx, &free_flags);
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

lrec_t* lrec_parse_mmap_nidx_multi_irs_multi_ifs(file_reader_mmap_state_t *phandle,
	char* irs, char* ifs, int irslen, int ifslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			*p = 0;

			idx++;
			key = make_nidx_key(idx, &free_flags);
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}
