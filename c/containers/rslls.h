// ================================================================
// Reusable singly linked list of string, with tail for append.
//
// Unlike slls which frees all memory on the destructor, this rslls class is
// designed to be reusable: there is a reset method which doesn't free all node
// structures.  The context is the CSV-reader which extends a singly linked
// list for every data line. Now, all CSV lines have nominally the same number
// of fields (say N), so using slls, the CSV reader would be allocating and
// freeing N nodes on every line. Mingled in with other mallocs and frees, this
// results in needless heap fragmentation. Here, by contrast, as a performance
// optimization, the CSV reader can keep and reuse the nodes of a list, only
// changing the value-pointers on each CSV line.
//
// This means that while an slls iteration looks like
//
//   for (sllse_t* pe = plist->phead; pe != NULL; pe = pe->pnext) {
//     ...
//   }
//
// an rslls iteration must also check length:
//
//   unsigned long long i = 0;
//   for (rsllse_t* pe = plist->phead; i < plist->length && pe != NULL; pe = pe->pnext, i++) {
//     ...
//   }
//
// ================================================================

#ifndef RSLLS_H
#define RSLLS_H

#include <stdio.h>
#include "containers/free_flags.h"

typedef struct _rsllse_t {
	char* value;
	char  free_flag;
	char  quote_flag;
	struct _rsllse_t *pnext;
} rsllse_t;

typedef struct _rslls_t {
	rsllse_t *phead;
	rsllse_t *ptail;
	unsigned long long length;
} rslls_t;

rslls_t* rslls_alloc();
void    rslls_reset(rslls_t* plist);
void    rslls_free(rslls_t* plist);
void    rslls_append(rslls_t* plist, char* value, char free_flag, char quote_flag);

void    rslls_print(rslls_t* plist);

#endif // RSLLS_H
