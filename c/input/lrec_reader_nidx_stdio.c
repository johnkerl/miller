#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_nidx_stdio_state_t {
	char irs;
	char ifs;
	int  allow_repeat_ifs;
} lrec_reader_nidx_stdio_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_nidx_stdio_func(FILE* input_stream, void* pvstate, context_t* pctx) {
	lrec_reader_nidx_stdio_state_t* pstate = pvstate;
	char* line = mlr_get_line(input_stream, pstate->irs);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_nidx_stdio(line, pstate->ifs, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void reset_nidx_func(void* pvstate) {
}

static void lrec_reader_nidx_stdio_free_func(void* pvstate) {
}

lrec_reader_stdio_t* lrec_reader_nidx_stdio_alloc(char irs, char ifs, int allow_repeat_ifs) {
	lrec_reader_stdio_t* plrec_reader_stdio = mlr_malloc_or_die(sizeof(lrec_reader_stdio_t));

	lrec_reader_nidx_stdio_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_nidx_stdio_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	plrec_reader_stdio->pvstate         = (void*)pstate;

	plrec_reader_stdio->plrec_reader_stdio_func = &lrec_reader_nidx_stdio_func;
	plrec_reader_stdio->preset_func  = &reset_nidx_func;
	plrec_reader_stdio->pfree_func   = &lrec_reader_nidx_stdio_free_func;

	return plrec_reader_stdio;
}
