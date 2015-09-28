#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "containers/sllv.h"

// lrec_reader_t impl for unit-test.
typedef struct _lrec_reader_in_memory_state_t {
	sllv_t* precords;
} lrec_reader_in_memory_state_t;

// ----------------------------------------------------------------
static lrec_t* lrec_reader_in_memory_process(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_in_memory_state_t* pstate = pvstate;

	if (pstate->precords->phead == NULL)
		return NULL;
	else
		return sllv_pop(pstate->precords);
}

// No-op for stateless readers such as this one.
static void lrec_reader_in_memory_sof(void* pvstate) {
}

// No-op for stateless readers such as this one.
static void lrec_reader_in_memory_free(void* pvstate) {
}

static void* lrec_reader_in_memory_vopen(void* pvstate, char* filename) {
	return NULL;
}

static void lrec_reader_in_memory_vclose(void* pvstate, void* pvhandle) {
}

lrec_reader_t* lrec_reader_in_memory_alloc(sllv_t* precords) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_in_memory_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_in_memory_state_t));
	pstate->precords = precords;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = &lrec_reader_in_memory_vopen;
	plrec_reader->pclose_func   = &lrec_reader_in_memory_vclose;
	plrec_reader->pprocess_func = &lrec_reader_in_memory_process;
	plrec_reader->psof_func     = &lrec_reader_in_memory_sof;
	plrec_reader->pfree_func    = &lrec_reader_in_memory_free;

	return plrec_reader;
}
