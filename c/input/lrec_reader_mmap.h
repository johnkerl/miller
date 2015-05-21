#ifndef LREC_READER_MMAP_H
#define LREC_READER_MMAP_H

// xxx rename x all to indicate line/file reader vs. lrec reader

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"
#include "input/file_reader_mmap.h"

// xxx rename mmap_state to mmap_handle?
typedef lrec_t* lrec_reader_mmap_func_t(mmap_reader_state_t* pfh, void* pvstate, context_t* pctx);
// xxx rename to sof_resetter or some such
typedef void    reset_mmap_func_t(void* pvstate);

typedef struct _lrec_reader_mmap_t {
	void*               pvstate;
	lrec_reader_mmap_func_t* plrec_reader_func;
	reset_mmap_func_t*  preset_func;
} lrec_reader_mmap_t;

#endif // LREC_READER_MMAP_H
