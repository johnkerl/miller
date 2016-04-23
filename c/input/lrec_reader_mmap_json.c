// ================================================================
// Unlike other Miller record-readers, there is no streaming for JSON input: no
// records are processed until EOF is seen. See also
// https://github.com/johnkerl/miller/issues/99.
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/file_reader_mmap.h"
#include "input/lrec_readers.h"
#include "input/json_parser.h"
#include "input/mlr_json_adapter.h"

typedef struct _lrec_reader_mmap_json_state_t {
	// The list of top-level JSON objects is backed by the file contents. The records are in turn
	// backed by the top-level JSON objects. This means the latter should not be freed while
	// the records are in used. (This is done to reduce data copies, for performance: we can
	// manipulate pointers to strings rather than copying strings.)
	//
	// In particular, in the multifile-input case, we need to keep *all* parsed JSON (and
	// not free one file's data when we proceed to the next) since records with pointers
	// into the parsed JSON may still be in use -- e.g. mlr sort.
	sllv_t* ptop_level_json_objects;
	sllv_t* precords;
	char* json_flatten_separator;
} lrec_reader_mmap_json_state_t;

static void    lrec_reader_mmap_json_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_json_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_json_alloc(char* json_flatten_separator) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_json_state_t));
	pstate->ptop_level_json_objects = sllv_alloc();
	pstate->precords                = sllv_alloc();
	pstate->json_flatten_separator  = json_flatten_separator;

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

	if (pstate->ptop_level_json_objects != NULL) {
		for (sllve_t* pe = pstate->ptop_level_json_objects->phead; pe != NULL; pe = pe->pnext) {
			json_value_t* top_level_json_object = pe->pvvalue;
			json_free_value(top_level_json_object);
		}
		sllv_free(pstate->ptop_level_json_objects);
		pstate->ptop_level_json_objects = NULL;
	}
	if (pstate->precords != NULL) {
		for (sllve_t* pf = pstate->precords->phead; pf != NULL; pf = pf->pnext) {
			lrec_t* prec = pf->pvvalue;
			lrec_free(prec);
		}
		sllv_free(pstate->precords);
	}

	free(pstate);
	free(preader);
}

// The mmap-JSON lrec-reader is non-streaming: we ingest all records here in the start-of-file hook.
// Then in the process method we pop one lrec off the list at a time, until they are all exhausted.
// This is in contrast to other Miller lrec-readers.
//
// It would be possible to extend the streaming framework to also have an end-of-file hook
// which we could use here to free parsed-JSON data. However, we simply leverage the start-of-file
// hook for the *next* file (if any) or the free method (if not): these free parsed-JSON structures
// from the previous file (if any).
static void lrec_reader_mmap_json_sof(void* pvstate, void* pvhandle) {
	lrec_reader_mmap_json_state_t* pstate = pvstate;
	file_reader_mmap_state_t* phandle = pvhandle;
	json_char* json_input = (json_char*)phandle->sol;
	json_value_t* parsed_top_level_json;
	json_char error_buf[JSON_ERROR_MAX];

	// This enables us to handle input of the form
	//
	//   { "a" : 1 }
	//   { "b" : 2 }
	//   { "c" : 3 }
	//
	// in addition to
	//
	// [
	//   { "a" : 1 }
	//   { "b" : 2 }
	//   { "c" : 3 }
	// ]
	//
	// This is in line with what jq can handle. In this case, json_parse will return
	// once for each top-level item and will give us back a pointer to the start of
	// the rest of the input stream, so we can call json_parse on the rest until it is
	// all exhausted.

	json_char* item_start = json_input;
	int length = phandle->eof - phandle->sol;

	while (TRUE) {
		parsed_top_level_json = json_parse(item_start, length, error_buf, &item_start);
		if (parsed_top_level_json == NULL) {
			fprintf(stderr, "%s: Unable to parse JSON data: %s\n", MLR_GLOBALS.argv0, error_buf);
			exit(1);
		}

		sllv_append(pstate->ptop_level_json_objects, parsed_top_level_json);

		// The lrecs have their string pointers pointing into the parsed-JSON objects (for
		// efficiency) so it's important we not free the latter until our free method.
		if (!reference_json_objects_as_lrecs(pstate->precords, parsed_top_level_json, pstate->json_flatten_separator)) {
			fprintf(stderr, "%s: Unable to parse JSON data.\n", MLR_GLOBALS.argv0);
			exit(1);
		}

		if (item_start == NULL)
			break;
		if (*item_start == 0)
			break;
		length -= (item_start - json_input);
		json_input = item_start;

	}
}

// ----------------------------------------------------------------
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx) {
	lrec_reader_mmap_json_state_t* pstate = pvstate;
	return sllv_pop(pstate->precords);
}
