#ifndef READER_H
#define READER_H

#include <stdio.h>
#include "containers/lrec.h"
#include "mapping/context.h"

typedef lrec_t* reader_func_t(FILE* fp, void* pvstate, context_t* pctx);
// xxx rename to sof_resetter or some such
typedef void    reset_func_t(void* pvstate);
typedef void    reader_free_func_t(void* pvstate);

typedef struct _reader_t {
	void*               pvstate;
	reader_func_t*      preader_func;
	reset_func_t*       preset_func;
	reader_free_func_t* pfree_func;
} reader_t;

#endif // READER_H
