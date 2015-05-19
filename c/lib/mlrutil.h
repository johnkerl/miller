#ifndef MLRUTIL_H
#define MLRUTIL_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define TRUE  1
#define FALSE 0

// ----------------------------------------------------------------
//int mlr_canonical_mod(int a, int n);
static inline int mlr_canonical_mod(int a, int n) {
	int r = a % n;
	if (r >= 0)
		return r;
	else
		return r+n;
}

int mlr_bsearch_double_for_insert(double* array, int size, double value);

// seconds since the epoch
double get_systime();

// ----------------------------------------------------------------
void* mlr_malloc_or_die(size_t size);

// ----------------------------------------------------------------
static inline int streq(char* a, char* b) {
	return !strcmp(a, b);
}

// xxx cmt mem mgt
char* mlr_alloc_string_from_double(double value, char* fmt);
char* mlr_alloc_string_from_ull(unsigned long long value);
char* mlr_alloc_string_from_int(int value);

double mlr_double_from_string_or_die(char* string);
int    mlr_try_double_from_string(char* string, double* pval);

// xxx cmt infrequently used; also cmt mem mgt
char* mlr_paste_2_strings(char* s1, char* s2);
char* mlr_paste_3_strings(char* s1, char* s2, char* s3);
char* mlr_paste_4_strings(char* s1, char* s2, char* s3, char* s4);
char* mlr_paste_5_strings(char* s1, char* s2, char* s3, char* s4, char* s5);

int mlr_string_hash_func(char *str);
int mlr_string_pair_hash_func(char* str1, char* str2);

// xxx cmt mem mgt
char* mlr_get_line(FILE* input_stream, char rs);

#endif // MLRUTIL_H
