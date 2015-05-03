#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/readers.h"

typedef struct _reader_nidx_state_t {
	char irs;
	char ifs;
} reader_nidx_state_t;

// ----------------------------------------------------------------
// xxx repeated ifs ...
//

static lrec_t* reader_nidx_func(FILE* input_stream, void* pvstate, context_t* pctx) {
	reader_nidx_state_t* pstate = pvstate;
	char* line = mlr_get_line(input_stream, pstate->irs);
	if (line == NULL)
		return NULL;
	else
		return lrec_parse_nidx(line, pstate->ifs, FALSE);
}


// No-op for stateless readers such as this one.
static void reset_nidx_func(void* pvstate) {
}

static void reader_nidx_free_func(void* pvstate) {
}

reader_t* reader_nidx_alloc(char irs, char ifs) {
	reader_t* preader = mlr_malloc_or_die(sizeof(reader_t));

	reader_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(reader_nidx_state_t));
	pstate->irs      = irs;
	pstate->ifs      = ifs;
	preader->pvstate = (void*)pstate;

	preader->preader_func = &reader_nidx_func;
	preader->preset_func  = &reset_nidx_func;
	preader->pfree_func   = &reader_nidx_free_func;

	return preader;
}
