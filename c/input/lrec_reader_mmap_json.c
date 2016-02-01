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
	sllv_t* parsed_json_objects;
	int num_records;
	int record_index;
} lrec_reader_mmap_json_state_t;

static void    lrec_reader_mmap_json_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_json_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_json_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_json_state_t));
	pstate->parsed_json_objects = NULL;
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
	for (sllve_t* pe = pstate->parsed_json_objects->phead; pe != NULL; pe = pe->pnext) {
		json_value_t* parsed_json_object = pe->pvvalue;
		json_value_free(parsed_json_object);
	}
	sllv_free(pstate->parsed_json_objects);
	free(pstate);
	free(preader);
}

// xxx cmt non-streaming; ingest-all here.
// xxx need an eof hook too!! or .... free on successive sofs, then on free ...
static void lrec_reader_mmap_json_sof(void* pvstate, void* pvhandle) {
	lrec_reader_mmap_json_state_t* pstate = pvstate;
	file_reader_mmap_state_t* phandle = pvhandle;
	json_char* json_input = (json_char*)phandle->sol;
	json_value_t* parsed_top_level_json;
	json_char error_buf[JSON_ERROR_MAX];
	json_settings_t settings = {
		.setting_flags = JSON_ENABLE_SEQUENTIAL_OBJECTS,
		.max_memory = 0
	};

	if (pstate->parsed_json_objects != NULL) {
		for (sllve_t* pe = pstate->parsed_json_objects->phead; pe != NULL; pe = pe->pnext) {
			json_value_t* parsed_json_object = pe->pvvalue;
			json_value_free(parsed_json_object);
		}
		// xxx make an sllv_free_with_callback & use it throughout
		sllv_free(pstate->parsed_json_objects);
	}
	pstate->parsed_json_objects = sllv_alloc();

	// xxx comment support missing outer [], as jq does.

	json_char* item_start = json_input;
	int length = phandle->eof - phandle->sol;;

	while (TRUE) {
		parsed_top_level_json = json_parse_ex(item_start, length, error_buf, &item_start, &settings);

		if (parsed_top_level_json == NULL) {
			fprintf(stderr, "Unable to parse JSON data: %s\n", error_buf);
			exit(1);
		}

		// xxx stub
		sllv_append(pstate->parsed_json_objects, parsed_top_level_json);

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
	//file_reader_mmap_state_t* phandle = pvhandle;
	//lrec_reader_mmap_json_state_t* pstate = pvstate;
	return NULL; // xxx eof temp stub
}
