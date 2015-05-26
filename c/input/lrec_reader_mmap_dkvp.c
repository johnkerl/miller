#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
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

lrec_t* lrec_parse_mmap_dkvp(file_reader_mmap_state_t *phandle, char irs, char ifs, char ips, int allow_repeat_ifs) {
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	char* key   = line;
	char* value = line;
	char* eol   = NULL;

	int idx = 0;
	for (char* p = line; *p; ) {
		if (*p == irs) {
			*p = 0;
			eol = p;
			phandle->sol = p+1;
			break;
		} else if (*p == ifs) {
			*p = 0;

			if (*key == 0) { // xxx to do: get file-name/line-number context in here.
				fprintf(stderr, "Empty key disallowed.\n");
				exit(1);
			}
			idx++;
			if (value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char  free_flags = 0;
				lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put_no_free(prec, key, value);
			}

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			key = p;
			value = p;
		} else if (*p == ips) {
			*p = 0;
			p++;
			value = p;
		} else {
			p++;
		}
	}
	idx++;
	if (*key == 0) { // xxx to do: get file-name/line-number context in here.
		fprintf(stderr, "Empty key disallowed.\n");
		exit(1);
	}
	if (value <= key) {
		char  free_flags = 0;
		lrec_put(prec, make_nidx_key(idx, &free_flags), value, free_flags);
	}
	else {
		lrec_put_no_free(prec, key, value);
	}

	return prec;
}
