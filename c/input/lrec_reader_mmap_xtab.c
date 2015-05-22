#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/lrec_parsers.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_xtab_state_t {
	char irs;
	char ips; // xxx make me real
	int allow_repeat_ips;
	int at_eof;
	// xxx need to remember EOF for subsequent read
} lrec_reader_mmap_xtab_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_xtab_func(file_reader_mmap_state_t* phandle, void* pvstate, context_t* pctx) {
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;

	if (pstate->at_eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab(phandle, pstate->irs, pstate->ips, pstate->allow_repeat_ips);
}

// xxx rename resets to sof_reset or some such
static void reset_xtab_func(void* pvstate) {
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;
	pstate->at_eof = FALSE;
}

lrec_reader_mmap_t* lrec_reader_mmap_xtab_alloc(char irs, char ips, int allow_repeat_ips) {
	lrec_reader_mmap_t* plrec_reader_stdio = mlr_malloc_or_die(sizeof(lrec_reader_mmap_t));

	lrec_reader_mmap_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_xtab_state_t));
	//pstate->ips              = ips;
	//pstate->allow_repeat_ips = allow_repeat_ips;
	pstate->irs              = irs;
	pstate->ips              = ' ';
	pstate->allow_repeat_ips = TRUE;
	pstate->at_eof           = FALSE;
	plrec_reader_stdio->pvstate         = (void*)pstate;

	plrec_reader_stdio->pprocess_func = &lrec_reader_mmap_xtab_func;
	plrec_reader_stdio->psof_func  = &reset_xtab_func;

	return plrec_reader_stdio;
}
