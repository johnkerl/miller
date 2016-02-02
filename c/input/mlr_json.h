// Transfers data from the JSON parser to Miller records

#ifndef MLR_JSON_H
#define MLR_JSON_H

#include "input/json.h"
#include "containers/lrec.h"


// ----------------------------------------------------------------
// xxx transfer func:
// input: top-level json value
// input: current sllv of object
// output: appended sllv
// json value will be freed, or transferred to the sllv

int transfer_objects(json_value_t* ptop_level_json, sllv_t* pobjects);

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

// ----------------------------------------------------------------
// xxx validate func: return object or die

json_value_t* validate_millerable_object(json_value_t* pjson);


#endif // MLR_JSON_H
