#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_stdio_dkvp_state_t {
	char irs;
	char ifs;
	char ips;
	int  allow_repeat_ifs;
} lrec_reader_stdio_dkvp_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_stdio_dkvp_process(void* pvstate, void* pvhandle, context_t* pctx) {
	FILE* input_stream = pvhandle;
	lrec_reader_stdio_dkvp_state_t* pstate = pvstate;

	char* line = mlr_get_line(input_stream, pstate->irs);

	if (line == NULL)
		return NULL;
	else
		return lrec_parse_stdio_dkvp(line, pstate->ifs, pstate->ips, pstate->allow_repeat_ifs);
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_sof(void* pvstate) {
}

// No-op for stateless readers such as this one.
static void lrec_reader_stdio_dkvp_free(void* pvstate) {
}

lrec_reader_t* lrec_reader_stdio_dkvp_alloc(char irs, char ifs, char ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_stdio_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->allow_repeat_ifs = allow_repeat_ifs;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &file_reader_stdio_vopen;
	plrec_reader->pclose_func   = &file_reader_stdio_vclose;
	plrec_reader->pprocess_func = &lrec_reader_stdio_dkvp_process;
	plrec_reader->psof_func     = &lrec_reader_stdio_dkvp_sof;
	plrec_reader->pfree_func    = &lrec_reader_stdio_dkvp_free;

	return plrec_reader;
}

// ----------------------------------------------------------------
// xxx needs checking on repeated occurrences of ps between fs occurrences. don't zero-poke there.
//
// xxx needs abend on null lhs.
//
// etc.

// "abc=def,ghi=jkl"
//      P     F     P
//      S     S     S
// "abc" "def" "ghi" "jkl"

// I couldn't find a performance gain using stdlib index(3) ... *maybe* even a
// fraction of a percent *slower*.

lrec_t* lrec_parse_stdio_dkvp(char* line, char ifs, char ips, int allow_repeat_ifs) {
	lrec_t* prec = lrec_dkvp_alloc(line);

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	char* p = line;

	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;

	for ( ; *p; ) {
		if (*p == ifs) {
			saw_ps = FALSE;
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
		} else if (*p == ips && !saw_ps) {
			*p = 0;
			p++;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0) {
		; // OK
	} else {
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
	}

	return prec;
}
