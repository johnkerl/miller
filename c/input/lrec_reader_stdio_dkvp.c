#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_dkvp_state_t {
	char irs;
	char ifs;
	char ips;
	int  allow_repeat_ifs;
} lrec_reader_stdio_dkvp_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_dkvp_func(FILE* input_stream, void* pvstate, context_t* pctx) {
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;

	char* line = mlr_get_line(input_stream, pstate->irs);

	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp(line, pstate->ifs, pstate->ips, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void reset_dkvp_func(void* pvstate) {
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_free(void* pvstate) {
}

lrec_reader_stdio_t* lrec_reader_stdio_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs) {
	lrec_reader_stdio_t* plrec_reader_stdio = mlr_malloc_or_die(sizeof(lrec_reader_stdio_t));

	lrec_reader_stdio_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_dkvp_state_t));
	pstate->irs = irs;
	pstate->ifs = ifs;
	pstate->ips = ips;
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	plrec_reader_stdio->pvstate = (void*)pstate;

	plrec_reader_stdio->pprocess_func = &lrec_reader_stdio_dkvp_func;
	plrec_reader_stdio->psof_func  = &reset_dkvp_func;
	plrec_reader_stdio->pfree_func   = &lrec_reader_stdio_dkvp_free;;

	return plrec_reader_stdio;
}
