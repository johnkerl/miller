#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_dkvp_state_t {
	char* irs;
	char* ifs;
	char* ips;
	int   irslen;
	int   ifslen;
	int   ipslen;
	int   allow_repeat_ifs;
} lrec_reader_stdio_dkvp_state_t;

static void    lrec_reader_stdio_dkvp_free(void* pvstate);
static void    lrec_reader_stdio_dkvp_sof(void* pvstate);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->ipslen           = strlen(ips);
	pstate->allow_repeat_ifs = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_stdio_vopen;
	plrec_reader->pclose_func   = file_reader_stdio_vclose;
	if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? &lrec_reader_stdio_dkvp_process_single_irs_single_others
			: &lrec_reader_stdio_dkvp_process_single_irs_multi_others;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? &lrec_reader_stdio_dkvp_process_multi_irs_single_others
			: &lrec_reader_stdio_dkvp_process_multi_irs_multi_others;
	}
	plrec_reader->psof_func     = lrec_reader_stdio_dkvp_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_dkvp_free;

	return plrec_reader;
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_free(void* pvstate) {
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_sof(void* pvstate) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	char* line = mlr_get_cline(input_stream, pstate->irs[0]);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_single_sep(line, pstate->ifs[0], pstate->ips[0], pstate->allow_repeat_ifs, pctx);
}

static lrec_t* lrec_reader_stdio_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	char* line = mlr_get_cline(input_stream, pstate->irs[0]);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_multi_sep(line, pstate->ifs, pstate->ips, pstate->ifslen, pstate->ipslen, pstate->allow_repeat_ifs, pctx);
}

static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_get_sline(input_stream, pstate->irs, pstate->irslen);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_single_sep(line, pstate->ifs[0], pstate->ips[0], pstate->allow_repeat_ifs, pctx);
}

static lrec_t* lrec_reader_stdio_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_get_sline(input_stream, pstate->irs, pstate->irslen);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp_multi_sep(line, pstate->ifs, pstate->ips, pstate->ifslen, pstate->ipslen, pstate->allow_repeat_ifs, pctx);
}

// ----------------------------------------------------------------
// "abc=def,ghi=jkl"
//      P     F     P
//      S     S     S
// "abc" "def" "ghi" "jkl"

// I couldn't find a performance gain using stdlib index(3) ... *maybe* even a
// fraction of a percent *slower*.

lrec_t* lrec_parse_stdio_dkvp_single_sep(char* line, char ifs, char ips, int allow_repeat_ifs, context_t* pctx) {
	lrec_t* prec = lrec_dkvp_alloc(line);

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	char* p = line;

	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; *p; ) {
		if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
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
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0 || value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put(prec, key, value, NO_FREE);
		}
	}

	return prec;
}

lrec_t* lrec_parse_stdio_dkvp_multi_sep(char* line, char* ifs, char* ips, int ifslen, int ipslen, int allow_repeat_ifs, context_t* pctx) {
	lrec_t* prec = lrec_dkvp_alloc(line);

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	char* p = line;

	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
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
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
		if (*key == 0 || value <= key) {
			char  free_flags = 0;
			lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
		}
		else {
			lrec_put(prec, key, value, NO_FREE);
		}
	}

	return prec;
}
