#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_dkvp_state_t {
	char irs;
	char ifs;
	char ips;
	int  allow_repeat_ifs;
} lrec_reader_mmap_dkvp_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_dkvp_process(file_reader_mmap_state_t* phandle, void* pvstate, context_t* pctx) {
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof) // xxx encapsulate a method for this ...
		return NULL;
	else
		return lrec_parse_mmap_dkvp(phandle, pstate->irs, pstate->ifs, pstate->ips, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_dkvp_sof(void* pvstate) {
}

lrec_reader_mmap_t* lrec_reader_mmap_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs) {
	lrec_reader_mmap_t* plrec_reader_mmap = mlr_malloc_or_die(sizeof(lrec_reader_mmap_t));

	lrec_reader_mmap_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->allow_repeat_ifs = allow_repeat_ifs;

	plrec_reader_mmap->pvstate       = (void*)pstate;
	plrec_reader_mmap->pprocess_func = &lrec_reader_mmap_dkvp_process;
	plrec_reader_mmap->psof_func     = &lrec_reader_mmap_dkvp_sof;

	return plrec_reader_mmap;
}
