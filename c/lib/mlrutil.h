#ifndef MLRUTIL_H
#define MLRUTIL_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include "mtrand.h"

#define TRUE  1
#define FALSE 0
#define NEITHER_TRUE_NOR_FALSE -1

//#define MLR_MALLOC_TRACE

// ----------------------------------------------------------------
#define MLR_INTERNAL_CODING_ERROR() mlr_internal_coding_error(__FILE__, __LINE__)
#define MLR_INTERNAL_CODING_ERROR_IF(v) mlr_internal_coding_error_if(v, __FILE__, __LINE__)
#define MLR_INTERNAL_CODING_ERROR_UNLESS(v) mlr_internal_coding_error_unless(v, __FILE__, __LINE__)
void mlr_internal_coding_error(char* file, int line);
void mlr_internal_coding_error_if(int v, char* file, int line);
void mlr_internal_coding_error_unless(int v, char* file, int line);

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
	while (n > 0 && *a && *b) {
		if (n-- <= 0) {
			return TRUE;
		}
		if (*a != *b) {
			return FALSE;
		}
		a++;
		b++;
	}
	if (n == 0)
		return TRUE;
	if (*a || *b) {
		return FALSE;
	}
	return TRUE;
#endif
}

// ----------------------------------------------------------------
// Like strsep but the sep argument is a multi-character delimiter,
// not a set of single-character delimiters.
char* mlr_strmsep(char **pstring, const char *sep, int seplen);

// ----------------------------------------------------------------
int mlr_bsearch_double_for_insert(double* array, int size, double value);

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
char * mlr_strdup_quoted_or_die(const char *s1);

// The caller should free the return values from each of these.
char* mlr_alloc_string_from_double(double value, char* fmt);
char* mlr_alloc_string_from_ull(unsigned long long value);
char* mlr_alloc_string_from_ll(long long value);
char* mlr_alloc_string_from_ll_and_format(long long value, char* fmt);
char* mlr_alloc_string_from_int(int value);
// The input doesn't include the null-terminator; the output does.
char* mlr_alloc_string_from_char_range(char* start, int num_bytes);

char* mlr_alloc_hexfmt_from_ll(long long value);

double mlr_double_from_string_or_die(char* string);
long long mlr_int_from_string_or_die(char* string);
int    mlr_try_float_from_string(char* string, double* pval);
int    mlr_try_int_from_string(char* string, long long* pval);

// For small integers (as of this writing, 0 .. 100) returns a static string representation.
// For other values, returns a dynamically allocated string representation.
char* low_int_to_string(int idx, char* pfree_flags);

// Inefficient and intended for call-rarely use. The caller should free the return values.
char* mlr_paste_2_strings(char* s1, char* s2);
char* mlr_paste_3_strings(char* s1, char* s2, char* s3);
char* mlr_paste_4_strings(char* s1, char* s2, char* s3, char* s4);
char* mlr_paste_5_strings(char* s1, char* s2, char* s3, char* s4, char* s5);

int mlr_string_hash_func(char *str);
int mlr_string_pair_hash_func(char* str1, char* str2);

int strlen_for_utf8_display(char* str);
int string_starts_with(char* string, char* prefix);
// If pstrlen is non-null, after return it will contain strlen(string) for
// convenience of the caller.
int string_ends_with(char* string, char* suffix, int* pstringlen);

int mlr_imax2(int a, int b);
int mlr_imax3(int a, int b, int c);
int power_of_two_above(int n);

// The caller should free the return value. Maps two-character sequences such as
// "\t", "\n", "\\" to single characters such as tab, newline, backslash, etc.
char* mlr_alloc_unbackslash(char* input);

// Miller DSL literals are unbackslashed: e.g. the two-character sequence "\t" is converted to a tab character, and
// users need to type "\\t" to get a backslash followed by a t. Well and good, but the system regex library handles
// backslashes not quite as I want. Namely, without this function,
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
// Using this function, backslashes can be escaped as the regex library requires, before I call regcomp:
//
//   echo 'x=a\tb' | mlr put '$x=sub($x,"\\t","TAB")'
//
// outputs
//
//   aTABb
//
// as desired.
char* mlr_alloc_double_backslash(char* input);

// The caller should free the return value.
char* read_file_into_memory(char* filename, size_t* psize);
// The caller should free the return value.
char* read_fp_into_memory(FILE* fp, size_t* psize);

// Returns a copy of the filename with random characters attached to the end.
char* alloc_suffixed_temp_file_name(char* filename);

#endif // MLRUTIL_H
