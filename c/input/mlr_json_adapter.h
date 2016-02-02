// Transfers data from the JSON parser to Miller records

#ifndef MLR_JSON_ADAPTER_H
#define MLR_JSON_ADAPTER_H

#include "input/json_parser.h"
#include "containers/lrec.h"


// ----------------------------------------------------------------
// xxx fix cmt:
// input: current sllv of lrecs
// input: top-level json value
// output: appended sllv
// xxx define pointer-ownership ... the sllv should not free the strings.

int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json, char* flatten_sep);

#endif // MLR_JSON_ADAPTER_H
