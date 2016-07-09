#ifndef MULTI_OUT_H
#define MULTI_OUT_H

#include <stdio.h>
#include "containers/lhmsv.h"

// ----------------------------------------------------------------
typedef struct _multi_out_t {
	// xxx to do: bound the number of open files and LRU them.
	lhmsv_t* pnames_to_fps;
} multi_out_t;

// ----------------------------------------------------------------
multi_out_t* multi_out_alloc();
void multi_out_free(multi_out_t* pmo);
FILE* multi_out_get_for_write(multi_out_t* pmo, char* filename);
FILE* multi_out_get_for_append(multi_out_t* pmo, char* filename);

#endif // MULTI_OUT_H
