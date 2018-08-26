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
	size_t nmatchmax, regmatch_t pmatch[restrict]);

// The return value is dynamically allocated even if there is no match, i.e. when output
// equals input.  The by-reference all-captured flag is true on return if all \1, etc.
// were satisfiable by parenthesized capture groups.
char* regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int* pmatched, int* pall_captured);

char* regex_gsub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement,
	int* pmatched, int* pall_captured, char *pfree_flags);

// The return value is dynamically allocated if there is a match, else it returns null.
char* regex_extract(char* input, regex_t* pregex);

// The regex library gives us an array of match pointers into the input string. This function strdups them
// out into separate storage, to implement "\0", "\1", "\2", etc. regex-captures for the =~ and !=~ operators.
// If the regex-captures array is null, it is allocated; otherwise it is resized. If the input regex does not
// match the regex, then the regex-captures array will be non-null but will have length 0.
void save_regex_captures(string_array_t** ppregex_captures, char* input, regmatch_t matches[], int nmatchmax);

// Given an array of regex-captures and an input string, interpolates the matches. E.g. if capture 1 is "abc"
// and capture 2 is "def" and the input is "hello \1 goodbye \2", then the output is a newly allocated string
// with value "hello abc goodbye def".  The was-allocated flag is an output flag: if true upon return, there was
// something modified and the returnv value should be freed; if false, nothing was modified and the input string was
// returned as the function value.
//
// To avoid performance regressions in non-match cases, this function quickly returns if the regex-captures array is
// NULL. See comments in mapper_put.c for more information.
char* interpolate_regex_captures(char* input, string_array_t* pregex_captures, int* pwas_allocated);

#endif // MLRREGEX_H
