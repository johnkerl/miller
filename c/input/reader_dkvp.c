#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/readers.h"

typedef struct _reader_dkvp_state_t {
	char irs;
	char ifs;
	char ips;
	int  allow_repeat_ifs;
} reader_dkvp_state_t;

// ----------------------------------------------------------------
static lrec_t* reader_dkvp_func(FILE* input_stream, void* pvstate, context_t* pctx) {
	reader_dkvp_state_t* pstate = pvstate;

	char* line = mlr_get_line(input_stream, pstate->irs);

	if (line == NULL)
		return NULL;
	else
		return lrec_parse_dkvp(line, pstate->ifs, pstate->ips, FALSE);
}

// No-op for stateless readers such as this one.
static void reset_dkvp_func(void* pvstate) {
}

// No-op for stateless readers such as this one.
static void reader_dkvp_free(void* pvstate) {
}

reader_t* reader_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs) {
	reader_t* preader = mlr_malloc_or_die(sizeof(reader_t));

	reader_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(reader_dkvp_state_t));
	pstate->irs = irs;
	pstate->ifs = ifs;
	pstate->ips = ips;
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	preader->pvstate = (void*)pstate;

	preader->preader_func = &reader_dkvp_func;
	preader->preset_func  = &reset_dkvp_func;
	preader->pfree_func   = &reader_dkvp_free;;

	return preader;
}
