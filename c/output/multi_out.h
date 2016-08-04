#ifndef MULTI_OUT_H
#define MULTI_OUT_H

#include <stdio.h>
#include "containers/lhmsv.h"
#include "output/file_output_mode.h"

// ----------------------------------------------------------------
typedef struct _multi_out_t {
	// xxx to do: bound the number of open files and LRU them.
	lhmsv_t* pnames_to_fps;
} multi_out_t;

// ----------------------------------------------------------------
multi_out_t* multi_out_alloc();
void  multi_out_close(multi_out_t* pmo);
void  multi_out_free(multi_out_t* pmo);
FILE* multi_out_get(multi_out_t* pmo, char* filename, file_output_mode_t file_output_mode);

#endif // MULTI_OUT_H
