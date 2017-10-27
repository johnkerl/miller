#ifndef MLR_JSON_ADAPTER_H
#define MLR_JSON_ADAPTER_H

#include "input/json_parser.h"
#include "containers/lrec.h"

// Given parsed JSON, constructs a list of lrecs with string values pointing into the parsed JSON.
// This is done for efficiency, to avoid data copying. It also means the parsed JSON should not be
// freed until the lrecs are freed.
//
// Default behavior on arrays is to fatal. There is a command-line option to skip them.  Miller
// doesn't have an array object in its DSL, only maps, and converting JSON arrays to int-keyed maps
// poses problems of irreversibility. (Namely, 'mlr --json cat foo.json' when foo.json contains
// arrays would result in output differing from input.)
int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json, char* flatten_sep,
	json_array_ingest_t json_array_ingest);

// * The buffer is an entire JSON blob, e.g. contents from stdio read or mmap; peof-psof is the file size so peof is one
//   byte *after* the last valid file byte.
// * The buffer is not assumed to be null-terminated.
// * Any lines beginning with comment_string are modified by poking space characters up to line_term.
void mlr_json_strip_comments(char* psof, char* peof, char* comment_string, char* line_term);

#endif // MLR_JSON_ADAPTER_H
