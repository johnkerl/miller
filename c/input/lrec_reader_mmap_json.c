// ================================================================
// mmap: easy pointer math
// stdio from file: stat, alloc, read. libify this.
// stdio from stdin: realloc w/ page-size fread. libify this.

// note @ mlr -h: no streaing for JSON input. No records are processed until EOF is seen.

// paginated:
//   json parse || error msg
// produce sllv of items

// sllv processing:
//   insist sllv.length == 1 & is array & each array item is an object,
//   or each sllv item is an object
// for each item:
//   loop over k/v pairs in the object and insist on level-1 only.
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"
#include "input/json.h"

typedef struct _lrec_reader_mmap_json_state_t {
	json_value_t* parsed_json;
	int num_records;
	int record_index;
} lrec_reader_mmap_json_state_t;

static void    lrec_reader_mmap_json_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_json_sof(void* pvstate);
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_json_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_json_state_t));
	pstate->num_records         = 0;
	pstate->record_index        = 0;

	plrec_reader->pvstate       = (void*)pstate;
	plrec_reader->popen_func    = file_reader_mmap_vopen;
	plrec_reader->pclose_func   = file_reader_mmap_vclose;
	plrec_reader->pprocess_func = lrec_reader_mmap_json_process;
	plrec_reader->psof_func     = lrec_reader_mmap_json_sof;
	plrec_reader->pfree_func    = lrec_reader_mmap_json_free;

	return plrec_reader;
}

static void lrec_reader_mmap_json_free(lrec_reader_t* preader) {
	lrec_reader_mmap_json_state_t* pstate = preader->pvstate;
	json_value_free(pstate->parsed_json);
	pstate->parsed_json = NULL;
	free(pstate);
	free(preader);
}

// xxx cmt non-streaming; ingest-all here.
static void lrec_reader_mmap_json_sof(void* pvstate) {
// xxx parse
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx) {
	//file_reader_mmap_state_t* phandle = pvhandle;
	//lrec_reader_mmap_json_state_t* pstate = pvstate;
	return NULL; // xxx eof temp stub
}
