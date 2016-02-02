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
#include "input/json_parser.h"
#include "input/mlr_json_adapter.h"

typedef struct _lrec_reader_mmap_json_state_t {
	// xxx just have a list of top-level objects and a list of lrecs?
	// xxx cmt re 3 layers of backing
	sllv_t* ptop_level_json_objects;
	sllv_t* precords;
	// xxx parameterize
	char* flatten_sep;
} lrec_reader_mmap_json_state_t;

static void    lrec_reader_mmap_json_free(lrec_reader_t* preader);
static void    lrec_reader_mmap_json_sof(void* pvstate, void* pvhandle);
static lrec_t* lrec_reader_mmap_json_process(void* pvstate, void* pvhandle, context_t* pctx);

// ----------------------------------------------------------------
lrec_reader_t* lrec_reader_mmap_json_alloc(char* irs, char* ifs, char* ips, int allow_repeat_ifs) {
	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));

	lrec_reader_mmap_json_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_mmap_json_state_t));
	pstate->ptop_level_json_objects = NULL;
	pstate->precords                = NULL;
	pstate->flatten_sep             = ":";

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
			json_value_free(top_level_json_object);
		}
		sllv_free(pstate->ptop_level_json_objects);
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

	// xxx make an sllv_free_with_callback & use it throughout

	if (pstate->ptop_level_json_objects != NULL) {
		for (sllve_t* pe = pstate->ptop_level_json_objects->phead; pe != NULL; pe = pe->pnext) {
			json_value_t* top_level_json_object = pe->pvvalue;
			json_value_free(top_level_json_object);
		}
		sllv_free(pstate->ptop_level_json_objects);
	}
	if (pstate->precords != NULL) {
		for (sllve_t* pf = pstate->precords->phead; pf != NULL; pf = pf->pnext) {
			lrec_t* prec = pf->pvvalue;
			lrec_free(prec);
		}
		sllv_free(pstate->precords);
	}

	pstate->ptop_level_json_objects = sllv_alloc();
	pstate->precords = sllv_alloc();

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
		//sllv_append(pstate->parsed_json_objects, parsed_top_level_json);
		// xxx swap arg order
		reference_json_objects_as_lrecs(pstate->precords, parsed_top_level_json, pstate->flatten_sep);

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
