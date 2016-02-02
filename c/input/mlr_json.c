#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/mlr_json.h"

// ----------------------------------------------------------------
// xxx transfer func:
// input: top-level json value
// input: current sllv of object
// output: appended sllv
// json value will be freed, or transferred to the sllv
// xxx define work done here: *not* recursing into JSON objects. just ascertaining that they *are* JSON objects.
// xxx define pointer-ownership ... do *not* call it transfer_objects but rather reference_objects. the sllv
//   should not free the strings.
// xxx why not make lrecs right here -- want to be able to produce data up to the bad point (or not ...)

int transfer_objects(json_value_t* ptop_level_json, sllv_t* pobjects) {
	if (ptop_level_json->type == JSON_ARRAY) {
		int n = ptop_level_json->u.array.length;
		for (int i = 0; i < n; i++) {
			json_value_t* pnext_level_json = ptop_level_json->u.array.values[i];
			if (pnext_level_json->type != JSON_OBJECT) {
				fprintf(stderr,
					"%s: found non-object (type %s) within top-level array. This is valid but unmillerable JSON.\n",
					MLR_GLOBALS.argv0, json_describe_type(ptop_level_json->type));
				return FALSE;
			}
			sllv_append(pobjects, validate_millerable_object(pnext_level_json));
		}
		// xxx free the pointer-array?!? put this logic as a method inside json.c/h.
		ptop_level_json->u.array.length = 0;
	} else if (ptop_level_json->type == JSON_OBJECT) {
		sllv_append(pobjects, validate_millerable_object(ptop_level_json));
		return TRUE;
	} else {
		fprintf(stderr,
			"%s: found non-terminal (type %s) at top level. This is valid but unmillerable JSON.\n",
			MLR_GLOBALS.argv0, json_describe_type(ptop_level_json->type));
		return FALSE;
	}
	return TRUE;
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
	// xxx redundantly assert this is of type JSON_OBJECT? or just note as precondition?
	int n = pjson->u.array.length;
	for (int i = 0; i < n; i++) {
		json_object_entry_t* pobject_entry = &pjson->u.object.values[i];
		char* key = (char*)pobject_entry->name;
		json_value_t* pvalue = pobject_entry->value;
		if (pvalue->type == JSON_ARRAY || pvalue->type == JSON_OBJECT) {
			fprintf(stderr,
				"%s: found nested non-object (type %s). This is valid but unmillerable JSON.\n",
				MLR_GLOBALS.argv0, json_describe_type(pvalue->type));
		}

		printf("xxx temp key=%s\n", key);

//typedef struct _json_object_entry_t {
	 //json_char * name;
	 //unsigned int name_length;
	 //struct _json_value_t * value;
//} json_object_entry_t;
	}
	return NULL; // xxx temp
}
