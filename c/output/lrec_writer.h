#ifndef LREC_WRITER_H
#define LREC_WRITER_H

#include <stdio.h>
#include "containers/lrec.h"
#include "lib/context.h"

struct _lrec_writer_t; // forward reference for method declarations

typedef void lrec_writer_process_func_t(void* pvstate, FILE* fp, lrec_t* prec, context_t* pctx);
typedef void lrec_writer_free_func_t(struct _lrec_writer_t* pwriter, context_t* pctx);

typedef struct _lrec_writer_t {
	void*                       pvstate;
	lrec_writer_process_func_t* pprocess_func;
	lrec_writer_free_func_t*    pfree_func; // virtual destructor
} lrec_writer_t;

#endif // LREC_WRITER_H
