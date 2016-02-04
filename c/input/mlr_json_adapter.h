#ifndef MLR_JSON_ADAPTER_H
#define MLR_JSON_ADAPTER_H

#include "input/json_parser.h"
#include "containers/lrec.h"

// Given parsed JSON, constructs a list of lrecs with string values pointing into the parsed JSON.
// This is done for efficiency, to avoid data copying. It also means the parsed JSON should not be
// freed until the lrecs are freed.
int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json, char* flatten_sep);

#endif // MLR_JSON_ADAPTER_H
