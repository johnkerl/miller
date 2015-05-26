#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_nidx_state_t {
	char irs;
	char ifs;
	int  allow_repeat_ifs;
} lrec_reader_stdio_nidx_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_nidx_process(FILE* input_stream, void* pvstate, context_t* pctx) {
	lrec_reader_stdio_nidx_state_t* pstate = pvstate;
	char* line = mlr_get_line(input_stream, pstate->irs);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_nidx(line, pstate->ifs, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_nidx_sof(void* pvstate) {
}

static void lrec_reader_stdio_nidx_free(void* pvstate) {
}

lrec_reader_stdio_t* lrec_reader_stdio_nidx_alloc(char irs, char ifs, int allow_repeat_ifs) {
	lrec_reader_stdio_t* plrec_reader_stdio = mlr_malloc_or_die(sizeof(lrec_reader_stdio_t));

	lrec_reader_stdio_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_nidx_state_t));
	pstate->irs                 = irs;
	pstate->ifs                 = ifs;
	pstate->allow_repeat_ifs    = allow_repeat_ifs;

	plrec_reader_stdio->pvstate       = (void*)pstate;
	plrec_reader_stdio->pprocess_func = &lrec_reader_stdio_nidx_process;
	plrec_reader_stdio->psof_func     = &lrec_reader_stdio_nidx_sof;
	plrec_reader_stdio->pfree_func    = &lrec_reader_stdio_nidx_free;

	return plrec_reader_stdio;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_stdio_nidx(char* line, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int idx = 0;
	char* key        = NULL;
	char* value      = line;
	char  free_flags = 0;

	for (char* p = line; *p; ) {
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
	key = make_nidx_key(idx, &free_flags);
	lrec_put(prec, key, value, free_flags);

	return prec;
}
