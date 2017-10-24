// ================================================================
// Note: there are multiple process methods with a lot of code duplication.
// This is intentional. Much of Miller's measured processing time is in the
// lrec-reader process methods. This is code which needs to execute on every
// byte of input and even moving a single runtime if-statement into a
// function-pointer assignment at alloc time can have noticeable effects on
// performance (5-10% in some cases).
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"

typedef struct _lrec_reader_mmap_dkvp_state_t {
	char* irs;
	char* ifs;
	char* ips;
	int   irslen;
	int   ifslen;
	int   ipslen;
	int   allow_repeat_ifs;
	int   do_auto_line_term;
} lrec_reader_mmap_dkvp_state_t;

static void    lrec_reader_mmap_dkvp_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_dkvp_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx);
static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_dkvp_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs,
	char* comment_string)
{
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_dkvp_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_dkvp_state_t));
	pstate->irs              = irs;
	pstate->ifs              = ifs;
	pstate->ips              = ips;
	pstate->irslen           = strlen(irs);
	pstate->ifslen           = strlen(ifs);
	pstate->ipslen           = strlen(ips);
	pstate->allow_repeat_ifs = allow_repeat_ifs;
	pstate->do_auto_line_term      = FALSE;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;
	if (streq(irs, "auto")) {
		// Auto means either lines end in "\n" or "\r\n" (LF or CRLF).  In
		// either case the final character is "\n". Then for autodetect we
		// simply check if there's a character in the line before the '\n', and
		// if that is '\r'.
		pstate->do_auto_line_term = TRUE;
		pstate->irs = "\n";
		pstate->irslen = 1;
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_mmap_dkvp_process_single_irs_single_others
			: lrec_reader_mmap_dkvp_process_single_irs_multi_others;
	} else if (pstate->irslen == 1) {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_mmap_dkvp_process_single_irs_single_others
			: lrec_reader_mmap_dkvp_process_single_irs_multi_others;
	} else {
		plrec_reader->pprocess_func = (pstate->ifslen == 1 && pstate->ipslen == 1)
			? lrec_reader_mmap_dkvp_process_multi_irs_single_others
			: lrec_reader_mmap_dkvp_process_multi_irs_multi_others;
	}
	plrec_reader->psof_func   = lrec_reader_mmap_dkvp_sof;
	plrec_reader->pfree_func  = lrec_reader_mmap_dkvp_free;

	return plrec_reader;
}

static void lrec_reader_mmap_dkvp_free(lrec_reader_t* preader) {
	free(preader->pvstate);
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_mmap_dkvp_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_single_irs_single_others(phandle, pstate->irs[0], pstate->ifs[0], pstate->ips[0],
			pstate->allow_repeat_ifs, pstate->do_auto_line_term, pctx);
}

static lrec_t* lrec_reader_mmap_dkvp_process_single_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_single_irs_multi_others(phandle, pstate->irs[0], pstate->ifs, pstate->ips,
			pstate->ifslen, pstate->ipslen, pstate->allow_repeat_ifs, pstate->do_auto_line_term, pctx);
}

static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_single_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_multi_irs_single_others(phandle, pstate->irs, pstate->ifs[0], pstate->ips[0],
			pstate->irslen, pstate->allow_repeat_ifs, pctx);
}

static lrec_t* lrec_reader_mmap_dkvp_process_multi_irs_multi_others(void* pvstate, void* pvhandle, context_t* pctx) {
	file_reader_mmap_state_t* phandle = pvhandle;
	lrec_reader_mmap_dkvp_state_t* pstate = pvstate;
	if (phandle->sol >= phandle->eof)
		return NULL;
	else
		return lrec_parse_mmap_dkvp_multi_irs_multi_others(phandle, pstate->irs, pstate->ifs, pstate->ips,
			pstate->irslen, pstate->ifslen, pstate->ipslen, pstate->allow_repeat_ifs, pctx);
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_mmap_dkvp_single_irs_single_others(file_reader_mmap_state_t *phandle,
	char irs, char ifs, char ips, int allow_repeat_ifs, int do_auto_line_term, context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;
	int saw_rs = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;

			if (do_auto_line_term) {
				if (p > line && p[-1] == '\r') {
					p[-1] = 0;
					context_set_autodetected_crlf(pctx);
				} else {
					context_set_autodetected_lf(pctx);
				}
			}

			phandle->sol = p+1;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char free_flags = NO_FREE;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
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
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0)
		return prec;

	// There are two ways out of that loop: saw IRS, or saw end of file.
	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof)
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			else
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			if (value >= phandle->eof)
				lrec_put(prec, key, "", NO_FREE);
			else
				lrec_put(prec, key, value, NO_FREE);
		}
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof) {
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, low_int_to_string(idx, &free_flags), copy, free_flags | FREE_ENTRY_VALUE);
			}
		}
		else {
			if (value >= phandle->eof) {
				lrec_put(prec, key, "", NO_FREE);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
			}
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_multi_irs_single_others(file_reader_mmap_state_t *phandle,
	char* irs, char ifs, char ips, int irslen, int allow_repeat_ifs, context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (*p == ifs)
			p++;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;
	int saw_rs = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			saw_rs = TRUE;
			break;
		} else if (*p == ifs) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char free_flags = NO_FREE;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
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
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0)
		return prec;

	// There are two ways out of that loop: saw IRS, or saw end of file.
	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof)
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			else
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			if (value >= phandle->eof)
				lrec_put(prec, key, "", NO_FREE);
			else
				lrec_put(prec, key, value, NO_FREE);
		}
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof) {
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, low_int_to_string(idx, &free_flags), copy, free_flags | FREE_ENTRY_VALUE);
			}
		}
		else {
			if (value >= phandle->eof) {
				lrec_put(prec, key, "", NO_FREE);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
			}
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_single_irs_multi_others(file_reader_mmap_state_t *phandle, char irs, char* ifs, char* ips,
	int ifslen, int ipslen, int allow_repeat_ifs, int do_auto_line_term, context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;
	int saw_rs = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (*p == irs) {
			*p = 0;

			if (do_auto_line_term) {
				if (p > line && p[-1] == '\r') {
					p[-1] = 0;
					context_set_autodetected_crlf(pctx);
				} else {
					context_set_autodetected_lf(pctx);
				}
			}

			phandle->sol = p+1;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char free_flags = NO_FREE;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
			}

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			key = p;
			value = p;
		} else if (streqn(p, ips, ipslen) && !saw_ps) {
			*p = 0;
			p += ipslen;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	*p = 0;
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0)
		return prec;

	// There are two ways out of that loop: saw IRS, or saw end of file.
	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof)
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			else
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			if (value >= phandle->eof)
				lrec_put(prec, key, "", NO_FREE);
			else
				lrec_put(prec, key, value, NO_FREE);
		}
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof) {
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, low_int_to_string(idx, &free_flags), copy, free_flags | FREE_ENTRY_VALUE);
			}
		}
		else {
			if (value >= phandle->eof) {
				lrec_put(prec, key, "", NO_FREE);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
			}
		}
	}

	return prec;
}

lrec_t* lrec_parse_mmap_dkvp_multi_irs_multi_others(file_reader_mmap_state_t *phandle,
	char* irs, char* ifs, char* ips, int irslen, int ifslen, int ipslen, int allow_repeat_ifs, context_t* pctx)
{
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;

	int idx = 0;
	char* p = line;
	if (allow_repeat_ifs) {
		while (streqn(p, ifs, ifslen))
			p += ifslen;
	}
	char* key   = p;
	char* value = p;

	int saw_ps = FALSE;
	int saw_rs = FALSE;

	for ( ; p < phandle->eof && *p; ) {
		if (streqn(p, irs, irslen)) {
			*p = 0;
			phandle->sol = p + irslen;
			saw_rs = TRUE;
			break;
		} else if (streqn(p, ifs, ifslen)) {
			saw_ps = FALSE;
			*p = 0;

			idx++;
			if (*key == 0 || value <= key) {
				// E.g the pair has no equals sign: "a" rather than "a=1" or
				// "a=".  Here we use the positional index as the key. This way
				// DKVP is a generalization of NIDX.
				char free_flags = NO_FREE;
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
			}
			else {
				lrec_put(prec, key, value, NO_FREE);
			}

			p += ifslen;
			if (allow_repeat_ifs) {
				while (streqn(p, ifs, ifslen))
					p += ifslen;
			}
			key = p;
			value = p;
		} else if (streqn(p, ips, ipslen) && !saw_ps) {
			*p = 0;
			p += ipslen;
			value = p;
			saw_ps = TRUE;
		} else {
			p++;
		}
	}
	if (p >= phandle->eof)
		phandle->sol = p+1;
	idx++;

	if (allow_repeat_ifs && *key == 0 && *value == 0)
		return prec;

	// There are two ways out of that loop: saw IRS, or saw end of file.
	if (saw_rs) {
		// Easy and simple case: we read until end of line.  We zero-poked the irs to a null character to terminate the
		// C string so it's OK to retain a pointer to that.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof)
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			else
				lrec_put(prec, low_int_to_string(idx, &free_flags), value, free_flags);
		}
		else {
			if (value >= phandle->eof)
				lrec_put(prec, key, "", NO_FREE);
			else
				lrec_put(prec, key, value, NO_FREE);
		}
	} else {
		// Messier case: we read to end of file without seeing end of line.  We can't always zero-poke a null character
		// to terminate the C string: if the file size is not a multiple of the OS page size it'll work (it's our
		// copy-on-write memory). But if the file size is a multiple of the page size, then zero-poking at EOF is one
		// byte past the page and that will segv us.
		if (*key == 0 || value <= key) {
			char free_flags = NO_FREE;
			if (value >= phandle->eof) {
				lrec_put(prec, low_int_to_string(idx, &free_flags), "", free_flags);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, low_int_to_string(idx, &free_flags), copy, free_flags | FREE_ENTRY_VALUE);
			}
		}
		else {
			if (value >= phandle->eof) {
				lrec_put(prec, key, "", NO_FREE);
			} else {
				char* copy = mlr_alloc_string_from_char_range(value, phandle->eof - value);
				lrec_put(prec, key, copy, FREE_ENTRY_VALUE);
			}
		}
	}

	return prec;
}
