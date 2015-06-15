#ifndef LREC_READER_MMAP_H
#define LREC_READER_MMAP_H

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"
#include "input/file_reader_mmap.h"

typedef lrec_t* lrec_reader_mmap_process_func_t(void* pvhandle, void* pvstate, context_t* pctx);
typedef void    lrec_reader_mmap_sof_func_t(void* pvstate);

typedef struct _lrec_reader_mmap_t {
	void*                            pvstate;
	lrec_reader_mmap_process_func_t* pprocess_func;
	lrec_reader_mmap_sof_func_t*     psof_func;
} lrec_reader_mmap_t;

#endif // LREC_READER_MMAP_H
