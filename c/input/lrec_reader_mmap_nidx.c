#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_nidx_state_t {
	char irs;
	char ifs;
	int  allow_repeat_ifs;
} lrec_reader_mmap_nidx_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_nidx_process(void* pvhandle, void* pvstate, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_nidx_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof) // xxx encapsulate a method for this ...
		return NULL;
	else
		return lrec_parse_mmap_nidx(phandle, pstate->irs, pstate->ifs, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_nidx_sof(void* pvstate) {
}

lrec_reader_t* lrec_reader_mmap_nidx_alloc(char irs, char ifs, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_nidx_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_nidx_state_t));
	pstate->irs                      = irs;
	pstate->ifs                      = ifs;
	pstate->allow_repeat_ifs         = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_mmap_vopen;
	plrec_reader->pclose_func   = &file_reader_mmap_vclose;
	plrec_reader->pprocess_func = &lrec_reader_mmap_nidx_process;
	plrec_reader->psof_func     = &lrec_reader_mmap_nidx_sof;
	plrec_reader->pfree_func    = NULL;

	return plrec_reader;
}

lrec_t* lrec_parse_mmap_nidx(file_reader_mmap_state_t *phandle, char irs, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char* key   = NULL;
	char* value = line;
	char free_flags = 0;

	for (char* p = line; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;
			phandle->sol = p+1;
			break;
		} else if (*p == ifs) {
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
