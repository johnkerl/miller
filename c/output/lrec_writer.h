#ifndef LREC_WRITER_H
#define LREC_WRITER_H

#include <stdio.h>
#include "containers/lrec.h"

typedef void lrec_writer_process_func_t(FILE* fp, lrec_t* prec, void* pvstate);
typedef void lrec_writer_free_func_t(void* pvstate);

typedef struct _lrec_writer_t {
	void*                       pvstate;
	lrec_writer_process_func_t* pprocess_func;
	lrec_writer_free_func_t*    pfree_func;
} lrec_writer_t;

#endif // LREC_WRITER_H
