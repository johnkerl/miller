#ifndef WRITER_H
#define WRITER_H

#include <stdio.h>
#include "containers/lrec.h"

typedef void lrec_writer_func_t(FILE* fp, lrec_t* prec, void* pvstate);
typedef void lrec_writer_free_func_t(void* pvstate);

typedef struct _lrec_writer_t {
	void*               pvstate;
	lrec_writer_func_t*      plrec_writer_func;
	lrec_writer_free_func_t* pfree_func;
} lrec_writer_t;

#endif // WRITER_H
