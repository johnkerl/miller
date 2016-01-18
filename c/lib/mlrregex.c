#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include <sys/time.h>
#include "lib/mlrregex.h"
#include "lib/mlr_globals.h"

// ----------------------------------------------------------------
// Succeeds or aborts the process. cflag REG_EXTENDED is already included.
regex_t* regcomp_or_die(regex_t* pregex, char* regex_string, int cflags) {
	cflags |= REG_EXTENDED;
	int rc = regcomp(pregex, regex_string, cflags);
	if (rc != 0) {
		size_t nbytes = regerror(rc, pregex, NULL, 0);
		char* errbuf = malloc(nbytes);
		(void)regerror(rc, pregex, errbuf, nbytes);
		fprintf(stderr, "%s: could not compile regex \"%s\" : %s\n",
			MLR_GLOBALS.argv0, regex_string, errbuf);
		exit(1);
	}
	return pregex;
}

// Always uses cflags with REG_EXTENDED.
// If the regex_string is of the form a.*b, compiles it using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b", compiles a.*b using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b"i, compiles a.*b using cflags with REG_ICASE.
regex_t* regcomp_or_die_quoted(regex_t* pregex, char* orig_regex_string, int cflags) {
	cflags |= REG_EXTENDED;
	if (string_starts_with(orig_regex_string, "\"")) {
		char* regex_string = mlr_strdup_or_die(orig_regex_string);
		int len = 0;
		if (string_ends_with(regex_string, "\"", &len)) {
			regex_string[len-1] = 0;
		} else if (string_ends_with(regex_string, "\"i", &len)) {
			regex_string[len-2] = 0;
			cflags |= REG_ICASE;
		} else {
			fprintf(stderr, "%s: imbalanced double-quote in regex [%s].\n",
				MLR_GLOBALS.argv0, regex_string);
			exit(1);
		}
		regcomp_or_die(pregex, regex_string+1, cflags);
		free(regex_string);
	} else {
		regcomp_or_die(pregex, orig_regex_string, cflags);
	}
	return pregex;
}

// Returns TRUE for match, FALSE for no match, and aborts the process if
// regexec returns anything else.
int regmatch_or_die(const regex_t* pregex, const char* restrict match_string,
	size_t nmatchmax, regmatch_t pmatch[restrict])
{
	int rc = regexec(pregex, match_string, nmatchmax, pmatch, 0);
	if (rc == 0) {
		return TRUE;
	} else if (rc == REG_NOMATCH) {
		return FALSE;
	} else {
		size_t nbytes = regerror(rc, pregex, NULL, 0);
		char* errbuf = malloc(nbytes);
		(void)regerror(rc, pregex, errbuf, nbytes);
		printf("regexec failure: %s\n", errbuf);
		exit(1);
	}
}

// If there is a match, the return value is dynamically allocated.  If not, the
// input is returned.
//
// Capture-group example:
// sed: $ echo '<<abcdefg>>'|sed 's/ab\(.\)d\(..\)g/AYEBEE\1DEE\2GEE/' gives <<AYEBEEcDEEefGEE>>
// mlr: echo 'x=<<abcdefg>>' | mlr put '$x = sub($x, "ab(.)d(..)g", "AYEBEE\1DEE\2GEE")' x=<<AYEBEEcDEEefGEE>>

char* regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int* pmatched, int *pall_captured)
{
	const size_t nmatchmax = 10; // Capture-groups \1 through \9 supported, along with entire-string match
	regmatch_t matches[nmatchmax];
	if (pall_captured)
		*pall_captured = TRUE;

	*pmatched = regmatch_or_die(pregex, input, nmatchmax, matches);
	if (!*pmatched) {
		return input;
	} else {
		sb_append_chars(psb, input, 0, matches[0].rm_so-1);
		char* p = replacement;
		while (*p) {
			if (p[0] == '\\' && isdigit(p[1])) {
				int idx = p[1] - '0';
				regmatch_t* pmatch = &matches[idx];
				if (pmatch->rm_so == -1) {
					if (pall_captured)
						*pall_captured = FALSE;
					sb_append_chars(psb, p, 0, 1);
				} else {
					sb_append_chars(psb, input, matches[idx].rm_so, matches[idx].rm_eo-1);
				}
				p += 2;
			} else {
				sb_append_char(psb, *p);
				p++;
			}
		}
		sb_append_chars(psb, input, matches[0].rm_eo, strlen(input));

		return sb_finish(psb);
	}
}

char* regex_gsub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int *pmatched, int* pall_captured, unsigned char* pfree_flags)
{
	const size_t nmatchmax = 10;
	regmatch_t matches[nmatchmax];
	*pmatched = FALSE;
	*pall_captured = TRUE;
	*pfree_flags = NO_FREE;

	int   match_start = 0;
	char* current_input = input;

	while (TRUE) {
		int matched = regmatch_or_die(pregex, &current_input[match_start], nmatchmax, matches);
		if (!matched) {
			return current_input;
		}
		*pmatched = TRUE;

		sb_append_chars(psb, current_input, 0, match_start + matches[0].rm_so-1);

		char* p = replacement;
		int len1 = psb->used_length;
		while (*p) {
			if (p[0] == '\\' && isdigit(p[1]) && p[1] != '0') {
				int idx = p[1] - '0';
				regmatch_t* pmatch = &matches[idx];
				if (pmatch->rm_so == -1) {
					*pall_captured = FALSE;
					sb_append_chars(psb, p, 0, 1);
				} else {
					sb_append_chars(psb, &current_input[match_start], matches[idx].rm_so, matches[idx].rm_eo-1);
				}
				p += 2;
			} else {
				sb_append_char(psb, *p);
				p++;
			}
		}

		int replen = psb->used_length - len1;
		sb_append_chars(psb, current_input, match_start + matches[0].rm_eo, strlen(current_input));

		char* next_input = sb_finish(psb);
		if (*pfree_flags & FREE_ENTRY_VALUE)
			free(current_input);
		current_input = next_input;
		*pfree_flags = FREE_ENTRY_VALUE;

		match_start += matches[0].rm_so + replen;
	}
}

// ----------------------------------------------------------------
void copy_regex_captures(string_array_t* pregex_captures_1_up, char* input, regmatch_t matches[], int nmatchmax) {
	// Slot 0 is the entire input string.
	// Slots 1 and up are substring matches for parenthesized capture expressions (if any).
	// Example regex "a(.*)e" with input string "abcde": slot 1 points to "bcd" and n = 1.
	// Slot 2 has rm_so == -1.
	// If all allocated slots have matches then there is no slot with -1's.
	int n = nmatchmax - 1;
	for (int i = 1; i < nmatchmax; i++) {
		if (matches[i].rm_so == -1) {
			n = i - 1;
			break;
		}
	}
	string_array_realloc(pregex_captures_1_up, n+1);
	for (int i = 1; i <= n; i++) {
		int len = matches[i].rm_eo - matches[i].rm_so;
		char* dst = mlr_malloc_or_die(len + 1);
		memcpy(dst, &input[matches[i].rm_so], len);
		dst[len] = 0;
		pregex_captures_1_up->strings[i] = dst;
	}
	pregex_captures_1_up->strings_need_freeing = TRUE;
}

// ----------------------------------------------------------------
char* interpolate_regex_captures(char* input, string_array_t* pregex_captures_1_up, int* pwas_allocated) {
	*pwas_allocated = FALSE;
	if (pregex_captures_1_up == NULL || pregex_captures_1_up->length == 0)
		return input;

	string_builder_t* psb = sb_alloc(32);

	char* p = input;
	while (*p) {
		if (p[0] == '\\' && isdigit(p[1]) && p[1] != '0') {
			int idx = p[1] - '0';
			if (idx < pregex_captures_1_up->length) {
				*pwas_allocated = TRUE;
				sb_append_string(psb, pregex_captures_1_up->strings[idx]);
			} else {
				sb_append_char(psb, p[0]);
				sb_append_char(psb, p[1]);
			}
			p += 2;
		} else {
			sb_append_char(psb, *p);
			p++;
		}
	}

	if (*pwas_allocated) {
		char* output = sb_finish(psb);
		sb_free(psb);
		return output;
	} else {
		sb_free(psb);
		return input;
	}
}

