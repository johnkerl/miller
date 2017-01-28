#include <stdio.h>
#include "lib/context.h"

void context_init(context_t* pctx, char* first_file_name) {
	pctx->nr       = 0;
	pctx->fnr      = 0;
	pctx->filenum  = 1;
	pctx->filename = first_file_name;
	pctx->auto_irs = "\n"; // xxx default to "\r\n" on Windows
}

void context_print(context_t* pctx, char* indent) {
	printf("%spctx at %p:\n", indent, pctx);
	printf("%s  nr       = %lld\n", indent, pctx->nr);
	printf("%s  fnr      = %lld\n", indent, pctx->fnr);
	printf("%s  filenum  = %d\n", indent, pctx->filenum);
	if (pctx->filename == NULL) {
		printf("%s  filename = null\n", indent);
	} else {
		printf("%s  filename = \"%s\"\n", indent, pctx->filename);
	}
}
