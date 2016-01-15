#ifndef MLRREGEX_H
#define MLRREGEX_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <regex.h>
#include "mlrutil.h"
#include "string_builder.h"
#include "string_array.h"

// Succeeds or aborts the process. cflag REG_EXTENDED is already included.
// Returns its first argument (after compilation).
regex_t* regcomp_or_die(regex_t* pregex, char* regex_string, int cflags);
// Always uses cflags with REG_EXTENDED.
// If the regex_string is of the form a.*b, compiles it using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b", compiles a.*b using cflags without REG_ICASE.
// If the regex_string is of the form "a.*b"i, compiles a.*b using cflags with REG_ICASE.
regex_t* regcomp_or_die_quoted(regex_t* pregex, char* regex_string, int cflags);

// Returns TRUE for match, FALSE for no match, and aborts the process if
// regexec returns anything else.
int regmatch_or_die(const regex_t* pregex, const char* restrict match_string,
	size_t nmatch, regmatch_t pmatch[restrict]);

// If there is a match, the return value is dynamically allocated and returned.
// If not, the input is returned.  So the caller should free the return value
// if matched == TRUE.  The by-reference all-captured flag is true on return if
// all \1, etc.  were satisfiable by parenthesized capture groups.
char* regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int* pmatched, int* pall_captured);

char* regex_gsub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement, int* pmatched, int* pall_captured, unsigned char *pfree_flags);

char* interpolate_regex_captures(char* input, string_array_t* pregex_captures, int* pwas_allocated);

#endif // MLRREGEX_H
