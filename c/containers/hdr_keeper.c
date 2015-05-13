#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/hdr_keeper.h"

hdr_keeper_t* hdr_keeper_alloc(char* line, slls_t* pkeys) {
	hdr_keeper_t* phdr_keeper = mlr_malloc_or_die(sizeof(hdr_keeper_t));
	phdr_keeper->line  = line;
	phdr_keeper->pkeys = pkeys;

	return phdr_keeper;
}

void hdr_keeper_free(hdr_keeper_t* phdr_keeper) {
	free(phdr_keeper->line);
	slls_free(phdr_keeper->pkeys);
	free(phdr_keeper);
}
