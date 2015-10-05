#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_dkvp_state_t {
	char* irs;
	char* ifs;
	char* ips;
	int   irslen;
	int   ifslen;
	int   ipslen;
	int   allow_repeat_ifs;
} lrec_reader_mmap_dkvp_state_t;

static void    lrec_reader_mmap_dkvp_sof(void* pvstate);
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->ipslen           = strlen(ips);
	pstate->allow_repeat_ifs = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;
	if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_mmap_dkvp_process_single_irs_single_others
			: lrec_reader_mmap_dkvp_process_single_irs_multi_others;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_mmap_dkvp_process_multi_irs_single_others
			: lrec_reader_mmap_dkvp_process_multi_irs_multi_others;
	}
	plrec_reader->psof_func   = lrec_reader_mmap_dkvp_sof;
	plrec_reader->pfree_func  = NULL;

	return plrec_reader;
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_dkvp_sof(void* pvstate) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_single_irs_single_others(phandle, pstate->irs[0], pstate->ifs[0], pstate->ips[0],
			pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_single_irs_multi_others(phandle, pstate->irs[0], pstate->ifs, pstate->ips,
			pstate->ifslen, pstate->ipslen,
			pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_multi_irs_single_others(phandle, pstate->irs, pstate->ifs[0], pstate->ips[0],
			pstate->irslen,
			pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_multi_irs_multi_others(phandle, pstate->irs, pstate->ifs, pstate->ips,
			pstate->irslen, pstate->ifslen, pstate->ipslen, pstate->allow_repeat_ifs);
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_dkvp_single_irs_single_others(file_reader_mmap_state_t *phandle,
	char irs, char ifs, char ips, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;
			phandle->sol = p+1;
			break;
		} else if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			if (*key == 0) { // xxx to do: get file-name/line-number context in here.
				fprintf(stderr, "Empty key disallowed.\n");
				exit(1);
			}
			idx++;
			if (value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put_no_free(prec, key, value);
			}

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			key = p;
			value = p;
		} else if (*p == ips && !saw_ps) {
			*p = 0;
			p++;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;
	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0) { // xxx to do: get file-name/line-number context in here.
			fprintf(stderr, "Empty key disallowed.\n");
			exit(1);
		}
		if (value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put_no_free(prec, key, value);
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_multi_irs_single_others(file_reader_mmap_state_t *phandle,
	char* irs, char ifs, char ips, int irslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			break;
		} else if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			if (*key == 0) { // xxx to do: get file-name/line-number context in here.
				fprintf(stderr, "Empty key disallowed.\n");
				exit(1);
			}
			idx++;
			if (value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put_no_free(prec, key, value);
			}

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			key = p;
			value = p;
		} else if (*p == ips && !saw_ps) {
			*p = 0;
			p++;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;
	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0) { // xxx to do: get file-name/line-number context in here.
			fprintf(stderr, "Empty key disallowed.\n");
			exit(1);
		}
		if (value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put_no_free(prec, key, value);
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_single_irs_multi_others(file_reader_mmap_state_t *phandle, char irs, char* ifs, char* ips,
	int ifslen, int ipslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;
			phandle->sol = p+1;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			if (*key == 0) { // xxx to do: get file-name/line-number context in here.
				fprintf(stderr, "Empty key disallowed.\n");
				exit(1);
			}
			idx++;
			if (value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put_no_free(prec, key, value);
			}

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			key = p;
			value = p;
		} else if (streqn(p, ips, ipslen) && !saw_ps) {
			*p = 0;
			p += ipslen;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;
	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0) { // xxx to do: get file-name/line-number context in here.
			fprintf(stderr, "Empty key disallowed.\n");
			exit(1);
		}
		if (value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put_no_free(prec, key, value);
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_multi_irs_multi_others(file_reader_mmap_state_t *phandle,
	char* irs, char* ifs, char* ips, int irslen, int ifslen, int ipslen, int allow_repeat_ifs)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			if (*key == 0) { // xxx to do: get file-name/line-number context in here.
				fprintf(stderr, "Empty key disallowed.\n");
				exit(1);
			}
			idx++;
			if (value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put_no_free(prec, key, value);
			}

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			key = p;
			value = p;
		} else if (streqn(p, ips, ipslen) && !saw_ps) {
			*p = 0;
			p += ipslen;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;
	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0) { // xxx to do: get file-name/line-number context in here.
			fprintf(stderr, "Empty key disallowed.\n");
			exit(1);
		}
		if (value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put_no_free(prec, key, value);
		}
	}

	return prec;
}
