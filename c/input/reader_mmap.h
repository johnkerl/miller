#ifndef READER_MMAP_H
#define READER_MMAP_H

// xxx rename x all to indicate line/file reader vs. lrec reader

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"
#include "input/mmap.h"

// xxx rename mmap_state to mmap_handle?
typedef lrec_t* reader_mmap_func_t(mmap_reader_state_t* pfh, void* pvstate, context_t* pctx);
// xxx rename to sof_resetter or some such
typedef void    reset_mmap_func_t(void* pvstate);

typedef struct _reader_mmap_t {
	void*               pvstate;
	reader_mmap_func_t* preader_func;
	reset_mmap_func_t*  preset_func;
} reader_mmap_t;

#endif // READER_MMAP_H
