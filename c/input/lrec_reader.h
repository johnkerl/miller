#ifndef LREC_READER_H
#define LREC_READER_H

#include <stdio.h>
#include "lib/context.h"
#include "containers/lrec.h"

struct _lrec_reader_t; // forward reference for method declarations

// The void* pvhandle is either FILE* for stdio readers or file_reader_mmap_state_t* for mmap readers.
typedef void*   lrec_reader_open_func_t(void* pvstate, char* prepipe, char* filename);
typedef void    lrec_reader_close_func_t(void* pvstate, void* pvhandle, char* prepipe);
typedef lrec_t* lrec_reader_process_func_t(void* pvstate, void* pvhandle, context_t* pctx);
typedef void    lrec_reader_sof_func_t(void* pvstate, void* pvhandle);
typedef void    lrec_reader_free_func_t(struct _lrec_reader_t* preader);

typedef struct _lrec_reader_t {
	void*                       pvstate;
	lrec_reader_open_func_t*    popen_func;
	lrec_reader_close_func_t*   pclose_func;
	lrec_reader_process_func_t* pprocess_func;
	lrec_reader_sof_func_t*     psof_func;
	lrec_reader_free_func_t*    pfree_func; // virtual destructor
} lrec_reader_t;

#endif // LREC_READER_H
