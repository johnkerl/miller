#include <stdlib.h>
#include "lib/mlrutil.h"
#include "containers/header_keeper.h"

header_keeper_t* header_keeper_alloc(char* line, slls_t* pkeys) {
	header_keeper_t* pheader_keeper = mlr_malloc_or_die(sizeof(header_keeper_t));
	pheader_keeper->line  = line;
	pheader_keeper->pkeys = pkeys;

	return pheader_keeper;
}

void header_keeper_free(header_keeper_t* pheader_keeper) {
	free(pheader_keeper->line);
	slls_free(pheader_keeper->pkeys);
	free(pheader_keeper);
}
