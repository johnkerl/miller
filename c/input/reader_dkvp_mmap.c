#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/mmap.h"
#include "input/readers.h"

typedef struct _reader_dkvp_mmap_state_t {
	char irs;
	char ifs;
	char ips;
	int  allow_repeat_ifs;
} reader_dkvp_mmap_state_t;

// ----------------------------------------------------------------
static lrec_t* reader_dkvp_mmap_func(mmap_reader_state_t* phandle, void* pvstate, context_t* pctx) {
	reader_dkvp_mmap_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof) // xxx encapsulate a method for this ...
		return NULL;
	else
		return lrec_parse_dkvp_mmap(phandle, pstate->irs, pstate->ifs, pstate->ips, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void reset_dkvp_mmap_func(void* pvstate) {
}

reader_mmap_t* reader_dkvp_mmap_alloc(char irs, char ifs, char ips, int allow_repeat_ifs) {
	reader_mmap_t* preader = mlr_malloc_or_die(sizeof(reader_mmap_t));

	reader_dkvp_mmap_state_t* pstate = mlr_malloc_or_die(sizeof(reader_dkvp_mmap_state_t));
	pstate->irs = irs;
	pstate->ifs = ifs;
	pstate->ips = ips;
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	preader->pvstate = (void*)pstate;

	preader->preader_func = &reader_dkvp_mmap_func;
	preader->preset_func  = &reset_dkvp_mmap_func;

	return preader;
}
