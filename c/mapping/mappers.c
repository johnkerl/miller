#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "mapping/mappers.h"

// ----------------------------------------------------------------
void mapper_chain_free(sllv_t* pmapper_chain, context_t* pctx) {
	for (sllve_t* pe = pmapper_chain->phead; pe != NULL; pe = pe->pnext) {
		mapper_t* pmapper = pe->pvvalue;
		pmapper->pfree_func(pmapper, pctx);
	}
	sllv_free(pmapper_chain);
}
