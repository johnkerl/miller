#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "containers/lrec.h"
#include "containers/lrec_parsers.h"

// ----------------------------------------------------------------
static char* static_nidx_keys[] = {
	"0",   "1",  "2",  "3",  "4",  "5",  "6",  "7",  "8",  "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "100"
};

static char* make_nidx_key(int idx, char* pfree_flags) {
	if ((0 <= idx) && (idx <= 100)) {
		*pfree_flags = 0;
		return static_nidx_keys[idx];
	} else {
		char buf[32];
		sprintf(buf, "%d", idx);
		*pfree_flags = LREC_FREE_ENTRY_KEY;
		return strdup(buf);
	}
}

lrec_t* lrec_parse_nidx(char* line, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_nidx_alloc(line);

	int idx = 0;
	char* key        = NULL;
	char* value      = line;
	char  free_flags = 0;

	for (char* p = line; *p; ) {
		if (*p == ifs) {
			*p = 0;

			idx++; key = make_nidx_key(idx, &free_flags);
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

lrec_t* lrec_parse_nidx_mmap(mmap_reader_state_t *phandle, char irs, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_unbacked_alloc();

	char* line  = phandle->sol;
	int idx = 0;
	char* key   = NULL;
	char* value = line;
	char* eol   = NULL;
	char free_flags = 0;

	for (char* p = line; *p; ) {
		if (*p == irs) {
			*p = 0;
			eol = p;
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

lrec_t* lrec_parse_dkvp(char* line, char ifs, char ips, int allow_repeat_ifs) {
	lrec_t* prec = lrec_dkvp_alloc(line);

	char* key   = line;
	char* value = line;

	// It would be easier to split the line on field separator (e.g. ","), then
	// split each key-value pair on pair separator (e.g. "="). But, that
	// requires two passes through the data. Here we do it in one pass.

	int idx = 0;
	for (char* p = line; *p; ) {
		if (*p == ifs) {
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

lrec_t* lrec_parse_dkvp_mmap(mmap_reader_state_t *phandle, char irs, char ifs, char ips, int allow_repeat_ifs) {
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

// ----------------------------------------------------------------
// xxx cmt mem-mgt
slls_t* split_csv_header_line(char* line, char ifs, int allow_repeat_ifs) {
	slls_t* plist = slls_alloc();
	if (*line == 0) // empty string splits to empty list
		return plist;

	char* start = line;
	for (char* p = line; *p; p++) {
		if (*p == ifs) {
			*p = 0;
			p++;
			// xxx hoist loop invariant at the cost of some code duplication
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			slls_add_no_free(plist, start);
			start = p;
		}
	}
	slls_add_no_free(plist, start);

	return plist;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_csv(hdr_keeper_t* phdr_keeper, char* data_line, char ifs, int allow_repeat_ifs) {
	lrec_t* prec = lrec_csv_alloc(data_line);
	char* key = NULL;
	char* value = data_line;

	// xxx needs hdr/data length check!!!!!!
	// xxx needs pe-non-null (hdr-empty) check:
	sllse_t* pe = phdr_keeper->pkeys->phead;
	for (char* p = data_line; *p; ) {
		if (*p == ifs) {
			*p = 0;

			if (pe == NULL) { // xxx to do: get file-name/line-number context in here
				fprintf(stderr, "Header-data length mismatch!\n");
				exit(1);
			}
			key = pe->value;
			lrec_put_no_free(prec, key, value);

			p++;
			if (allow_repeat_ifs) {
				while (*p == ifs)
					p++;
			}
			value = p;
			pe = pe->pnext;
		} else {
			p++;
		}
	}
	if (pe == NULL) {
		fprintf(stderr, "Header-data length mismatch!\n");
		exit(1);
	}
	key = pe->value;
	lrec_put_no_free(prec, key, value);
	if (pe->pnext != NULL) {
		fprintf(stderr, "Header-data length mismatch!\n");
		exit(1);
	}

	return prec;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_csv_mmap(hdr_keeper_t* phdr_keeper, mmap_reader_state_t* phandle, char irs, char ifs, int allow_repeat_ifs) {
	return NULL; // xxx stub
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_xtab(slls_t* pxtab_lines, char ips, int allow_repeat_ips) {
	lrec_t* prec = lrec_xtab_alloc(pxtab_lines);

	for (sllse_t* pe = pxtab_lines->phead; pe != NULL; pe = pe->pnext) {
		char* line = pe->value;
		char* p = line;
		char* key = p;

		while (*p != 0 && *p != ips)
			p++;
		if (*p == 0) {
			lrec_put_no_free(prec, key, "");
		} else {
			while (*p != 0 && *p == ips) {
				*p = 0;
				p++;
			}
			lrec_put_no_free(prec, key, p);
		}
	}

	return prec;
}

// ----------------------------------------------------------------
lrec_t* lrec_parse_xtab_mmap(mmap_reader_state_t* phandle, char irs, char ips, int allow_repeat_ips) {

	if (phandle->sol >= phandle->eof)
		return NULL;

	lrec_t* prec = lrec_unbacked_alloc();
	while (*phandle->sol == irs)
		phandle->sol++;

	if (phandle->sol >= phandle->eof)
		return NULL;

	// Loop over fields, one per line
	while (TRUE) {
		char* line  = phandle->sol;
		char* key   = line;
		char* value = "";
		char* eol   = NULL;
		char* p;

		// Construct one field
		for (p = line; *p; ) {
			if (*p == irs) {
				*p = 0;
				eol = p;
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
