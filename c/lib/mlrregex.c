#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <ctype.h>
#include <sys/time.h>
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mlr_globals.h"
#include "lib/free_flags.h"

// ----------------------------------------------------------------
// Succeeds or aborts the process. cflag REG_EXTENDED is already included.
//
// Reason for the double-backslashing routine: Miller DSL literals are unbackslashed, e.g. the
// two-character sequence "\t" is converted to a tab character, and users need to type "\\t" to get
// a backslash followed by a t. Well and good, but the system regex library handles backslashes not
// quite as I want. Namely, without double-backslashing,
//
//   echo 'x=a\tb' | mlr put '$x=sub($x,"\\t","TAB")'
//
// (note: not echo -e, but just plain echo) outputs
//
//   a\TABb
//
// while
//
//   echo 'x=a\tb' | mlr put '$x=sub($x,"\\\\t","TAB")'
//
// outputs
//
//   aTABb
//
// Using double-backslashing, backslashes can be escaped as the regex library requires, before I call regcomp:
//
//   echo 'x=a\tb' | mlr put '$x=sub($x,"\\t","TAB")'
//
// outputs
//
//   aTABb
//
// as desired.

regex_t* regcomp_or_die(regex_t* pregex, char* regex_string, int cflags) {
	cflags |= REG_EXTENDED;
	char* doubly_backslashed = mlr_alloc_double_backslash(regex_string);
	int rc = regcomp(pregex, doubly_backslashed, cflags);
	free(doubly_backslashed);
	if (rc != 0) {
		size_t nbytes = regerror(rc, pregex, NULL, 0);
		char* errbuf = malloc(nbytes);
		(void)regerror(rc, pregex, errbuf, nbytes);
		fprintf(stderr, "%s: could not compile regex \"%s\" : %s\n",
			MLR_GLOBALS.bargv0, regex_string, errbuf);
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
				MLR_GLOBALS.bargv0, regex_string);
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

// Capture-group example:
// sed: $ echo '<<abcdefg>>'|sed 's/ab\(.\)d\(..\)g/AYEBEE\1DEE\2GEE/' gives <<AYEBEEcDEEefGEE>>
// mlr: echo 'x=<<abcdefg>>' | mlr put '$x = sub($x, "ab(.)d(..)g", "AYEBEE\1DEE\2GEE")' x=<<AYEBEEcDEEefGEE>>

char* regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int* pmatched, int *pall_captured)
{
	const size_t nmatchmax = 10; // Capture-groups \1 through \9 supported, along with entire-string match \0
	regmatch_t matches[nmatchmax];
	if (pall_captured)
		*pall_captured = TRUE;

	*pmatched = regmatch_or_die(pregex, input, nmatchmax, matches);
	if (!*pmatched) {
		return mlr_strdup_or_die(input);
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
					// implicitly append empty string by doing nothing at all --
					// we don't need to write:
					// sb_append_string(psb, "");
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
	int *pmatched, int* pall_captured, char* pfree_flags)
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
			if (input == current_input) {
				*pfree_flags = FREE_ENTRY_VALUE;
				return mlr_strdup_or_die(current_input);
			} else {
				return current_input;
			}
		}
		*pmatched = TRUE;

		sb_append_chars(psb, current_input, 0, match_start + matches[0].rm_so-1);

		char* p = replacement;
		int len1 = psb->used_length;
		while (*p) {
			if (p[0] == '\\' && isdigit(p[1])) {
				int idx = p[1] - '0';
				regmatch_t* pmatch = &matches[idx];
				if (pmatch->rm_so == -1) {
					*pall_captured = FALSE;
					// implicitly append empty string by doing nothing at all --
					// we don't need to write:
					// sb_append_string(psb, "");
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
char* regextract(char* input, regex_t* pregex) {
	const size_t nmatchmax = 1;
	regmatch_t matches[nmatchmax];

	int matched = regmatch_or_die(pregex, input, nmatchmax, matches);
	if (!matched) {
		return NULL;
	}
	regmatch_t* pmatch = &matches[0];
	int len = pmatch->rm_eo - pmatch->rm_so;
	return mlr_alloc_string_from_char_range(&input[pmatch->rm_so], len);
}

// ----------------------------------------------------------------
char* regextract_or_else(char* input, regex_t* pregex, char* default_value) {
	const size_t nmatchmax = 1;
	regmatch_t matches[nmatchmax];

	int matched = regmatch_or_die(pregex, input, nmatchmax, matches);
	if (!matched) {
		return mlr_strdup_or_die(default_value);
	}
	regmatch_t* pmatch = &matches[0];
	int len = pmatch->rm_eo - pmatch->rm_so;
	return mlr_alloc_string_from_char_range(&input[pmatch->rm_so], len);
}

// ----------------------------------------------------------------
// Slot 0 is the entire matched input string.
// Slots 1 and up are substring matches for parenthesized capture expressions (if any).
// Example regex "a(.*)e" with input string "abcde": slot 1 points to "bcd" and match_count = 2.
// Slot 2 has rm_so == -1.
// (If all allocated slots have matches then there is no slot with -1's.)

// Input "abcde"
// Regex "a(.*)e"
// matches[0].rm_so =  0, matches[0].rm_eo =  5
// matches[1].rm_so =  1, matches[1].rm_eo =  4
// matches[2].rm_so = -1, matches[2].rm_eo = -1
//
// pregex_captures->length = 2
// pregex_captures->strings[0] = "abcde"
// pregex_captures->strings[1] = "bcd"
//
// Note that even if there is no match, a non-null zero-length regex-captures array is returned (by reference).
// This is important: see the comments in mapper_put for details.

void save_regex_captures(string_array_t** ppregex_captures, char* input, regmatch_t matches[], int nmatchmax) {
	int match_count = 0;
	match_count = 0;
	// In fully occupied case, there will be no slots with -1's.
	// Using optional regex captures, one slot may have rm_so == rm_eo == -1 (i.e. trivial) while a subsequent slot
	// may be non-trivial. So we need to check all slots.
	for (int i = 0; i < nmatchmax; i++) {
		if (matches[i].rm_so != -1) {
			match_count = i + 1;
		}
	}
	if (*ppregex_captures != NULL)
		string_array_realloc(*ppregex_captures, match_count);
	else
		*ppregex_captures = string_array_alloc(match_count);
	string_array_t* pregex_captures = *ppregex_captures;
	if (match_count >= 1) {
		for (int i = 0; i < match_count; i++) {
			int len = matches[i].rm_eo - matches[i].rm_so;
			pregex_captures->strings[i] = mlr_alloc_string_from_char_range(&input[matches[i].rm_so], len);
		}
		pregex_captures->strings_need_freeing = TRUE;
	}
}

// ----------------------------------------------------------------
// Using the above example:
// Input "abcde"
// Regex "a(.*)e"
//
// pregex_captures->length = 2
// pregex_captures->strings[0] = "abcde"
// pregex_captures->strings[1] = "bcd"
//
// "\0" should be replaced with "abcde".
// "\1" should be replaced with "bcd".
// "\2" through "\9" should be replaced with "".

char* interpolate_regex_captures(char* input, string_array_t* pregex_captures, int* pwas_allocated) {
	*pwas_allocated = FALSE;

	string_builder_t* psb = sb_alloc(32);

	char* p = input;
	while (*p) {
		if (p[0] == '\\' && isdigit(p[1])) {
			*pwas_allocated = TRUE;
			int idx = p[1] - '0';
			if (idx < pregex_captures->length)
				sb_append_string(psb, pregex_captures->strings[idx]);
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
