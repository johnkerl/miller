#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/line_readers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_nidx_state_t {
	char* irs;
	char* ifs;
	int   irslen;
	int   ifslen;
	int   allow_repeat_ifs;
} lrec_reader_stdio_nidx_state_t;

static void    lrec_reader_stdio_nidx_free(lrec_reader_t* preader, void* pvstate);
static void    lrec_reader_stdio_nidx_sof(void* pvstate);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_stdio_nidx_alloc(char* irs, char* ifs, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_nidx_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->allow_repeat_ifs = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_stdio_vopen;
	plrec_reader->pclose_func   = file_reader_stdio_vclose;
	if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_nidx_process_single_irs_single_ifs
			: &lrec_reader_stdio_nidx_process_single_irs_multi_ifs;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1)
			? &lrec_reader_stdio_nidx_process_multi_irs_single_ifs
			: &lrec_reader_stdio_nidx_process_multi_irs_multi_ifs;
	}
	plrec_reader->psof_func     = lrec_reader_stdio_nidx_sof;
	plrec_reader->pfree_func    = lrec_reader_stdio_nidx_free;

	return plrec_reader;
}

static void lrec_reader_stdio_nidx_free(lrec_reader_t* preader, void* pvstate) {
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_nidx_sof(void* pvstate) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_nidx_process_single_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	char* line = mlr_get_cline(input_stream, pstate->irs[0]);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_single_sep(line, pstate->ifs[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_single_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	char* line = mlr_get_cline(input_stream, pstate->irs[0]);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_multi_sep(line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_single_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_get_sline(input_stream, pstate->irs, pstate->irslen);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_single_sep(line, pstate->ifs[0], pstate->allow_repeat_ifs);
}

static lrec_t* lrec_reader_stdio_nidx_process_multi_irs_multi_ifs(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	FILE* input_stream = pvhandle;
	char* line = mlr_get_sline(input_stream, pstate->irs, pstate->irslen);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx_multi_sep(line, pstate->ifs, pstate->ifslen, pstate->allow_repeat_ifs);
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_nidx_single_sep(char* line, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int idx = 0;
	char  free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; *p; ) {
		if (*p == ifs) {
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_nidx_multi_sep(char* line, char* ifs, int ifslen, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int  idx = 0;
	char free_flags = 0;

	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = NULL;
	char* value = p;
	for ( ; *p; ) {
		if (streqn(p, ifs, ifslen)) {
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
	idx++;

	if (allow_repeat_ifs && *value == 0) {
		; // OK
	} else {
		key = make_nidx_key(idx, &free_flags);
		lrec_put(prec, key, value, free_flags);
	}

	return prec;
}
