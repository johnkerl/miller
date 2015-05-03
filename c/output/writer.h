#ifndef WRITER_H
#define WRITER_H

#include <stdio.h>
#include "containers/lrec.h"

typedef void writer_func_t(FILE* fp, lrec_t* prec, void* pvstate);
typedef void writer_free_func_t(void* pvstate);

typedef struct _writer_t {
	void*               pvstate;
	writer_func_t*      pwriter_func;
	writer_free_func_t* pfree_func;
} writer_t;

#endif // WRITER_H
