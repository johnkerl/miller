#ifndef LREC_READER_STDIO_H
#define LREC_READER_STDIO_H

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"

typedef lrec_t* lrec_reader_stdio_process_func_t(FILE* fp, void* pvstate, context_t* pctx);
// xxx rename to sof_func or some such
typedef void    lrec_reader_stdio_sof_func_t(void* pvstate);
typedef void    lrec_reader_stdio_free_func_t(void* pvstate);

typedef struct _lrec_reader_stdio_t {
	void*               pvstate;
	lrec_reader_stdio_process_func_t* pprocess_func;
	lrec_reader_stdio_sof_func_t*   psof_func;
	lrec_reader_stdio_free_func_t*    pfree_func;
} lrec_reader_stdio_t;

#endif // LREC_READER_STDIO_H
