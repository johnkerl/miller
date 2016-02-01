#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "input/file_reader_stdio.h"
#include "input/lrec_readers.h"
#include "containers/sllv.h"

// lrec_reader_t impl for unit-test.
typedef struct _lrec_reader_in_memory_state_t {
	sllv_t* precords;
} lrec_reader_in_memory_state_t;

static void    lrec_reader_in_memory_free(lrec_reader_t* preader);
static void    lrec_reader_in_memory_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_in_memory_process(void* pvstate, void* pvhandle, context_t* pctx);
static void*   lrec_reader_in_memory_vopen(void* pvstate, char* prepipe, char* filename);
static void    lrec_reader_in_memory_vclose(void* pvstate, void* pvhandle, char* prepipe);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_in_memory_alloc(sllv_t* precords) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_in_memory_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_in_memory_state_t));
	pstate->precords = precords;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = lrec_reader_in_memory_vopen;
	plrec_reader->pclose_func   = lrec_reader_in_memory_vclose;
	plrec_reader->pprocess_func = lrec_reader_in_memory_process;
	plrec_reader->psof_func     = lrec_reader_in_memory_sof;
	plrec_reader->pfree_func    = lrec_reader_in_memory_free;

	return plrec_reader;
}

static void lrec_reader_in_memory_free(lrec_reader_t* preader) {
	lrec_reader_in_memory_state_t* pstate = preader->pvstate;
	sllv_free(pstate->precords);
	free(pstate);
	free(preader);
}

// No-op for stateless readers such as this one.
static void lrec_reader_in_memory_sof(void* pvstate, void* pvhandle) {
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_in_memory_process(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_in_memory_state_t* pstate = pvstate;

	if (pstate->precords->phead == NULL)
		return NULL;
	else
		return sllv_pop(pstate->precords);
}

static void* lrec_reader_in_memory_vopen(void* pvstate, char* prepipe, char* filename) {
	// popen is a stdio construct, not an mmap construct, and it can't be supported here.
	if (prepipe != NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	return NULL;
}

static void lrec_reader_in_memory_vclose(void* pvstate, void* pvhandle, char* prepipe) {
}
