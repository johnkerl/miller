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

void transfer_objects(json_value_t* ptop_level_json, sllv_t* pobjects);

		// xxx
		//switch(parsed_top_level_json->type) {
		//case json_array:
		//	for each {
		//		validate & add it
		//	}
		//	break;
		//case json_object:
		//	validate & add it
		//	break;
		//default:
		//	break;
		//}

// ----------------------------------------------------------------
// xxx validate func: return object or die

json_value_t* validate_millerable_object(json_value_t* pjson);


#endif // MLR_JSON_H
