// ================================================================
// Singly linked list of string, with tail for append.
// ================================================================

#ifndef SLLS_H
#define SLLS_H

#include <stdio.h>
#include "lib/free_flags.h"

typedef struct _sllse_t {
	char* value;
	char  free_flag;
	struct _sllse_t *pnext;
} sllse_t;

typedef struct _slls_t {
	sllse_t *phead;
	sllse_t *ptail;
	unsigned long long length;
} slls_t;

slls_t* slls_alloc();
int     slls_size(slls_t* plist);
slls_t* slls_copy(slls_t* pold);
void    slls_free(slls_t* plist);
slls_t* slls_single_with_free(char* value);
slls_t* slls_single_no_free(char* value);
void    slls_append_with_free(slls_t* plist, char* value);
void    slls_append_no_free(slls_t* plist, char* value);
void    slls_append(slls_t* plist, char* value, char free_flag);
int     slls_equals(slls_t* pa, slls_t* pb);
slls_t* slls_from_line(char* line, char ifs, int allow_repeat_ifs);

void    slls_reverse(slls_t* plist);
int     slls_hash_func(slls_t *plist);
int     slls_compare_lexically(slls_t* pa, slls_t* pb);

void    slls_sort(slls_t* plist);

// Debug routines:
char*   slls_join(slls_t* plist, char* ofs);
void    slls_print(slls_t* plist);
void    slls_print_quoted(slls_t* plist);

#endif // SLLS_H
