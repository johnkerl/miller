#include "context.h"

void context_init(context_t* pctx, char* first_file_name) {
	pctx->nr       = 0;
	pctx->fnr      = 0;
	pctx->filenum  = 1;
	pctx->filename = first_file_name;
}
