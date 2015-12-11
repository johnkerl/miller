#ifndef MLRUTIL_H
#define MLRUTIL_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <regex.h>
#include "string_builder.h"

#define TRUE  1
#define FALSE 0
#define NEITHER_TRUE_NOR_FALSE -1

//#define MLR_MALLOC_TRACE

// ----------------------------------------------------------------
//int mlr_canonical_mod(int a, int n);
static inline int mlr_canonical_mod(int a, int n) {
	int r = a % n;
	if (r >= 0)
		return r;
	else
		return r+n;
}

// ----------------------------------------------------------------
// strcmp computes signs; we don't need that -- only equality or inequality.
static inline int streq(char* a, char* b) {
#if 0 // performance comparison
	return !strcmp(a, b);
#else
	while (*a && *b) {
		if (*a != *b)
			return FALSE;
		a++;
		b++;
	}
	if (*a || *b)
		return FALSE;
	return TRUE;
#endif
}

// strncmp computes signs; we don't need that -- only equality or inequality.
static inline int streqn(char* a, char* b, int n) {
#if 0 // performance comparison
	return !strncmp(a, b, n);
#else
	while (*a && *b) {
		if (n-- <= 0) {
			return TRUE;
		}
		if (*a != *b) {
			return FALSE;
		}
		a++;
		b++;
	}
	if (*a || *b) {
		return (n <= 0) ? TRUE : FALSE;
	}
	return TRUE;
#endif
}

// ----------------------------------------------------------------
int mlr_bsearch_double_for_insert(double* array, int size, double value);

// seconds since the epoch
double get_systime();

void*  mlr_malloc_or_die(size_t size);
void*  mlr_realloc_or_die(void *ptr, size_t size);
static inline char * mlr_strdup_or_die(const char *s1) {
	char* s2 = strdup(s1);
	if (s2 == NULL) {
		fprintf(stderr, "malloc/strdup failed\n");
		exit(1);
	}
#ifdef MLR_MALLOC_TRACE
	fprintf(stderr, "STRDUP size=%d,p=%p\n", (int)strlen(s2), s2);
#endif
	return s2;
}

// The caller should free the return values from each of these.
char* mlr_alloc_string_from_double(double value, char* fmt);
char* mlr_alloc_string_from_ull(unsigned long long value);
char* mlr_alloc_string_from_ll(long long value);
char* mlr_alloc_string_from_ll_and_format(long long value, char* fmt);
char* mlr_alloc_string_from_int(int value);

char* mlr_alloc_hexfmt_from_ll(long long value);

double mlr_double_from_string_or_die(char* string);
long long mlr_int_from_string_or_die(char* string);
int    mlr_try_float_from_string(char* string, double* pval);
int    mlr_try_int_from_string(char* string, long long* pval);

// Inefficient and intended for call-rarely use. The caller should free the return values.
char* mlr_paste_2_strings(char* s1, char* s2);
char* mlr_paste_3_strings(char* s1, char* s2, char* s3);
char* mlr_paste_4_strings(char* s1, char* s2, char* s3, char* s4);
char* mlr_paste_5_strings(char* s1, char* s2, char* s3, char* s4, char* s5);

int mlr_string_hash_func(char *str);
int mlr_string_pair_hash_func(char* str1, char* str2);

// portable timegm replacement
time_t mlr_timegm (struct tm *tm);

int strlen_for_utf8_display(char* str);
int string_starts_with(char* string, char* prefix);
// If pstrlen is non-null, after return it will contain strlen(string) for
// convenience of the caller.
int string_ends_with(char* string, char* suffix, int* pstringlen);

int mlr_imax2(int a, int b);
int power_of_two_ceil(int n);

// The caller should free the return value. Maps two-character sequences such as
// "\t", "\n", "\\" to single characters such as tab, newline, backslash, etc.
char* mlr_unbackslash(char* input);

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

// If there is a match, input is freed and return value is dynamically
// allocated.  If not, input is returned.  So in either case, the caller should
// free the return value, and it is assumed that the input has been dynamically
// allocated. The by-reference all-captured flag is true on return if all \1,
// etc.  were satisfiable by parenthesized capture groups.
char*  regex_sub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement, int* pmatched, int* pall_captured);
char* regex_gsub(char* input, regex_t* pregex, string_builder_t* psb, char* replacement, int* pmatched, int* pall_captured);

#endif // MLRUTIL_H
