#include "mlr_json.h"

// ----------------------------------------------------------------
// xxx transfer func:
// input: top-level json value
// input: current sllv of object
// output: appended sllv
// json value will be freed, or transferred to the sllv

void transfer_objects(json_value_t* ptop_level_json, sllv_t* pobjects) {
	if (ptop_level_json->type == JSON_ARRAY) {
	} else if (ptop_level_json->type == JSON_OBJECT) {
	} else {
	}
}

		// xxx
		//switch(parsed_top_level_json->type) {
		//case JSON_ARRAY:
		//	for each {
		//		validate & add it
		//	}
		//	break;
		//case JSON_OBJECT:
		//	validate & add it
		//	break;
		//default:
		//	break;
		//}

// JSON_NONE
// JSON_OBJECT
// JSON_ARRAY
// JSON_INTEGER
// JSON_DOUBLE
// JSON_STRING
// JSON_BOOLEAN
// JSON_NULL

// ----------------------------------------------------------------
// xxx validate func: return object or die

json_value_t* validate_millerable_object(json_value_t* pjson) {
	return NULL; // xxx temp
}
