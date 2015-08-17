#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_xtab_state_t {
	char irs;
	char ips; // xxx make me real
	int allow_repeat_ips;
} lrec_reader_mmap_xtab_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_xtab_process(void* pvhandle, void* pvstate, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_xtab_state_t* pstate = pvstate;

	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_xtab(phandle, pstate->irs, pstate->ips, pstate->allow_repeat_ips);
}

static void lrec_reader_mmap_xtab_sof(void* pvstate) {
}

lrec_reader_t* lrec_reader_mmap_xtab_alloc(char irs, char ips, int allow_repeat_ips) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_xtab_state_t));
	//pstate->ips              = ips;
	//pstate->allow_repeat_ips = allow_repeat_ips;
	pstate->irs                 = irs;
	pstate->ips                 = ' ';
	pstate->allow_repeat_ips    = TRUE;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_mmap_vopen;
	plrec_reader->pclose_func   = &file_reader_mmap_vclose;
	plrec_reader->pprocess_func = &lrec_reader_mmap_xtab_process;
	plrec_reader->psof_func     = &lrec_reader_mmap_xtab_sof;
	plrec_reader->pfree_func    = NULL;

	return plrec_reader;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_xtab(file_reader_mmap_state_t* phandle, char irs, char ips, int allow_repeat_ips) {

	while (phandle->sol < phandle->eof && *phandle->sol == irs)
		phandle->sol++;

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* p;

		// Construct one field
		for (p = line; p < phandle->eof && *p; ) {
			if (*p == irs) {
				*p = 0;
				phandle->sol = p+1;
				break;
			} else if (*p == ips) {
				key = line;
				*p = 0;

				p++;
				if (allow_repeat_ips) {
					while (*p == ips)
						p++;
				}
				value = p;
			} else {
				p++;
			}
		}

		lrec_put_no_free(prec, key, value);

		if (phandle->sol >= phandle->eof || *phandle->sol == irs)
			break;
	}
	if (prec->field_count == 0) {
		lrec_free(prec);
		return NULL;
	} else {
		return prec;
	}
}
