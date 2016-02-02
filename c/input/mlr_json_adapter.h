// Transfers data from the JSON parser to Miller records

#ifndef MLR_JSON_ADAPTER_H
#define MLR_JSON_ADAPTER_H

#include "input/json_parser.h"
#include "containers/lrec.h"


// ----------------------------------------------------------------
// input: current sllv of lrec
// input: top-level json value
// output: appended sllv
// xxx mem-mgmt semantics

int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json);

#endif // MLR_JSON_ADAPTER_H
