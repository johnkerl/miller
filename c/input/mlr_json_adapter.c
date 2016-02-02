#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "input/mlr_json_adapter.h"

// xxx validate func: return lrec or null.
// xxx mem-mgmt semantics
static lrec_t* validate_millerable_object(json_value_t* pjson_object, char* flatten_sep);

// ----------------------------------------------------------------
// xxx transfer func:
// input: top-level json value
// input: current sllv of object
// output: appended sllv
// json value will be freed, or transferred to the sllv
// xxx define work done here: *not* recursing into JSON objects. just ascertaining that they *are* JSON objects.
// xxx define pointer-ownership ... the sllv should not free the strings.
// xxx why not make lrecs right here -- want to be able to produce data up to the bad point (or not ...)

int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json, char* flatten_sep) {
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
			sllv_append(precords, validate_millerable_object(pnext_level_json, flatten_sep));
		}
		// xxx free the pointer-array?!? put this logic as a method inside json.c/h.
		ptop_level_json->u.array.length = 0;
	} else if (ptop_level_json->type == JSON_OBJECT) {
		sllv_append(precords, validate_millerable_object(ptop_level_json, flatten_sep));
		return TRUE;
	} else {
		fprintf(stderr,
			"%s: found non-terminal (type %s) at top level. This is valid but unmillerable JSON.\n",
			MLR_GLOBALS.argv0, json_describe_type(ptop_level_json->type));
		return FALSE;
	}
	return TRUE;
}

// JSON_NONE
// JSON_OBJECT
// JSON_ARRAY
// JSON_INTEGER
// JSON_DOUBLE
// JSON_STRING
// JSON_BOOLEAN
// JSON_NULL

// ----------------------------------------------------------------
// xxx validate func: return object or null.
// xxx rename

lrec_t* validate_millerable_object(json_value_t* pjson, char* flatten_sep) {
	// xxx redundantly assert this is of type JSON_OBJECT? or just note as precondition?
	lrec_t* prec = lrec_unbacked_alloc();
	int n = pjson->u.array.length;
	for (int i = 0; i < n; i++) {
		json_object_entry_t* pobject_entry = &pjson->u.object.values[i];
		char* key = (char*)pobject_entry->name;
		char* lrec_value = NULL;
		char free_flags = NO_FREE;

		json_value_t* pjson_value = pobject_entry->pvalue;
		switch (pjson_value->type) {

		case JSON_NONE:
			lrec_put(prec, key, "", NO_FREE);
			break;
		case JSON_NULL:
			lrec_put(prec, key, "", NO_FREE);
			break;
		case JSON_STRING:
			lrec_value = pjson_value->u.string.ptr;
			lrec_put(prec, key, lrec_value, NO_FREE);
			break;
		case JSON_BOOLEAN:
			lrec_value = pjson_value->u.boolean ? "true" : "false";
			lrec_put(prec, key, lrec_value, NO_FREE);
			break;
		case JSON_OBJECT:
			// xxx not yet
			fprintf(stderr,
				"%s: found nested JSON object. This is valid but unmillerable JSON.\n",
				MLR_GLOBALS.argv0);
			exit(1);
			break;
		case JSON_ARRAY:
			fprintf(stderr,
				"%s: found array item within JSON object. This is valid but unmillerable JSON.\n",
				MLR_GLOBALS.argv0);
			break;
		case JSON_INTEGER:
			lrec_value = mlr_alloc_string_from_ll(pjson_value->u.integer);
			lrec_put(prec, key, lrec_value, free_flags | FREE_ENTRY_VALUE);
			break;
		case JSON_DOUBLE:
			lrec_value = mlr_alloc_string_from_double(pjson_value->u.dbl, MLR_GLOBALS.ofmt);
			lrec_put(prec, key, lrec_value, free_flags | FREE_ENTRY_VALUE);
			break;
		default:
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
			break;
		}

	}
	return prec; // xxx temp
}
