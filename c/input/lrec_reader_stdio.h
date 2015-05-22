#ifndef LREC_READER_STDIO_H
#define LREC_READER_STDIO_H

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"

// xxx rename to lrec_reader
typedef lrec_t* lrec_reader_stdio_func_t(FILE* fp, void* pvstate, context_t* pctx);
// xxx rename to sof_resetter or some such
typedef void    reset_func_t(void* pvstate);
typedef void    lrec_reader_stdio_free_func_t(void* pvstate);

typedef struct _lrec_reader_stdio_t {
	void*               pvstate;
	lrec_reader_stdio_func_t*      plrec_reader_stdio_func;
	reset_func_t*       preset_func;
	lrec_reader_stdio_free_func_t* pfree_func;
} lrec_reader_stdio_t;

#endif // LREC_READER_STDIO_H
