#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "cli/json_array_ingest.h"
#include "input/mlr_json_adapter.h"

static lrec_t* validate_millerable_object(json_value_t* pjson_object, char* flatten_sep,
	json_array_ingest_t json_array_ingest);
static int populate_from_nested_object(lrec_t* prec, json_value_t* pjson_object, char* prefix, char* flatten_sep,
	json_array_ingest_t json_array_ingest);
static int populate_from_nested_array(lrec_t* prec, json_value_t* pjson_array, char* prefix, char* flatten_sep,
	json_array_ingest_t json_array_ingest);

// ----------------------------------------------------------------
int reference_json_objects_as_lrecs(sllv_t* precords, json_value_t* ptop_level_json, char* flatten_sep,
	json_array_ingest_t json_array_ingest)
{
	if (ptop_level_json->type == JSON_ARRAY) {
		int n = ptop_level_json->u.array.length;
		for (int i = 0; i < n; i++) {
			json_value_t* pnext_level_json = ptop_level_json->u.array.values[i];
			if (pnext_level_json->type != JSON_OBJECT) {
				fprintf(stderr,
					"%s: found non-object (type %s) within top-level array. This is valid but unmillerable JSON.\n",
					MLR_GLOBALS.bargv0, json_describe_type(ptop_level_json->type));
				return FALSE;
			}
			lrec_t* prec = validate_millerable_object(pnext_level_json, flatten_sep, json_array_ingest);
			if (prec == NULL)
				return FALSE;
			sllv_append(precords, prec);
		}
	} else if (ptop_level_json->type == JSON_OBJECT) {
		lrec_t* prec = validate_millerable_object(ptop_level_json, flatten_sep, json_array_ingest);
		if (prec == NULL)
			return FALSE;
		sllv_append(precords, prec);
	} else {
		fprintf(stderr,
			"%s: found non-terminal (type %s) at top level. This is valid but unmillerable JSON.\n",
			MLR_GLOBALS.bargv0, json_describe_type(ptop_level_json->type));
		return FALSE;
	}
	return TRUE;
}

// ----------------------------------------------------------------
// Returns NULL if the JSON object is not millerable, else returns a new lrec with string pointers
// backed by the JSON object.
//
// Precondition: the JSON value is assumed to have already been checked to be of type JSON_OBJECT.

lrec_t* validate_millerable_object(json_value_t* pjson, char* flatten_sep, json_array_ingest_t json_array_ingest) {
	lrec_t* prec = lrec_unbacked_alloc();
	int n = pjson->u.array.length;
	for (int i = 0; i < n; i++) {
		json_object_entry_t* pobject_entry = &pjson->u.object.p.values[i];
		char* key = (char*)pobject_entry->name;
		char* prefix = NULL;

		json_value_t* pjson_value = pobject_entry->pvalue;
		switch (pjson_value->type) {

		case JSON_NONE:
			lrec_put(prec, key, "", NO_FREE);
			break;
		case JSON_NULL:
			lrec_put(prec, key, "", NO_FREE);
			break;

		case JSON_OBJECT:
			// This could be made more efficient ... the string length is in the json_value_t.
			prefix = mlr_paste_2_strings(key, flatten_sep);
			if (!populate_from_nested_object(prec, pjson_value, prefix, flatten_sep, json_array_ingest))
				return NULL;
			free(prefix);
			break;
		case JSON_ARRAY:
			switch (json_array_ingest) {
			case JSON_ARRAY_INGEST_FATAL:
				fprintf(stderr,
					"%s: found array item within JSON object. This is valid but unmillerable JSON.\n"
					"Use --json-skip-arrays-on-input to exclude these from input without fataling.\n"
					"Or, --json-map-arrays-on-input to convert them to integer-indexed maps.\n",
					MLR_GLOBALS.bargv0);
				return NULL;
				break;
			case JSON_ARRAY_INGEST_AS_MAP:
				prefix = mlr_paste_2_strings(key, flatten_sep);
				if (!populate_from_nested_array(prec, pjson_value, prefix, flatten_sep, json_array_ingest)) {
					free(prefix);
					return NULL;
				}
				free(prefix);
				break;
			// xxx other cases!
			default:
				break;
			}
			break;

		case JSON_STRING:
			lrec_put(prec, key, pjson_value->u.string.ptr, NO_FREE);
			break;

		case JSON_BOOLEAN:
			lrec_put(prec, key, pjson_value->u.boolean.sval, NO_FREE);
			break;
		case JSON_INTEGER:
			lrec_put(prec, key, pjson_value->u.integer.sval, NO_FREE);
			break;
		case JSON_DOUBLE:
			lrec_put(prec, key, pjson_value->u.dbl.sval, NO_FREE);
			break;
		default:
			MLR_INTERNAL_CODING_ERROR();
			break;
		}

	}
	return prec;
}

// ----------------------------------------------------------------
// Example: the JSON object has { "a": { "b" : 1, "c" : 2 } }. Then we add "a:b" => "1" and "a:c" => "2"
// to the lrec.

static int populate_from_nested_object(lrec_t* prec, json_value_t* pjson_object, char* prefix, char* flatten_sep,
	json_array_ingest_t json_array_ingest)
{
	int n = pjson_object->u.object.length;
	for (int i = 0; i < n; i++) {
		json_object_entry_t* pobject_entry = &pjson_object->u.object.p.values[i];
		char* json_key = (char*)pobject_entry->name;
		json_value_t* pjson_value = pobject_entry->pvalue;
		char* lrec_key = mlr_paste_2_strings(prefix, json_key);
		char* next_prefix = NULL;

		switch (pjson_value->type) {
		case JSON_NONE:
			lrec_put(prec, lrec_key, "", FREE_ENTRY_KEY);
			break;
		case JSON_NULL:
			lrec_put(prec, lrec_key, "", FREE_ENTRY_KEY);
			break;
		case JSON_STRING:
			lrec_put(prec, lrec_key, pjson_value->u.string.ptr, FREE_ENTRY_KEY);
			break;
		case JSON_BOOLEAN:
			lrec_put(prec, lrec_key, pjson_value->u.boolean.sval, FREE_ENTRY_KEY);
			break;
		case JSON_OBJECT:
			next_prefix = mlr_paste_2_strings(lrec_key, flatten_sep);
			if (!populate_from_nested_object(prec, pjson_value, next_prefix, flatten_sep, json_array_ingest))
				return FALSE;
			free(next_prefix);
			free(lrec_key);
			break;
		case JSON_ARRAY:
			switch (json_array_ingest) {
			case JSON_ARRAY_INGEST_FATAL:
				fprintf(stderr,
					"%s: found array item within JSON object. This is valid but unmillerable JSON.\n"
					"Use --json-skip-arrays-on-input to exclude these from input without fataling.\n"
					"Or, --json-map-arrays-on-input to convert them to integer-indexed maps.\n",
					MLR_GLOBALS.bargv0);
				free(lrec_key);
				return FALSE;
				break;
			case JSON_ARRAY_INGEST_AS_MAP:
				next_prefix = mlr_paste_2_strings(lrec_key, flatten_sep);
				if (!populate_from_nested_array(prec, pjson_value, next_prefix, flatten_sep, json_array_ingest)) {
					free(next_prefix);
					free(lrec_key);
					return FALSE;
				}
				free(next_prefix);
				free(lrec_key);
				break;
			// xxx other cases!
			default:
				free(lrec_key);
				break;
			}
			break;
		case JSON_INTEGER:
			lrec_put(prec, lrec_key, pjson_value->u.integer.sval, FREE_ENTRY_KEY);
			break;
		case JSON_DOUBLE:
			lrec_put(prec, lrec_key, pjson_value->u.dbl.sval, FREE_ENTRY_KEY);
			break;
		default:
			MLR_INTERNAL_CODING_ERROR();
			break;
		}
	}
	return TRUE;
}

static int populate_from_nested_array(lrec_t* prec, json_value_t* pjson_array, char* prefix, char* flatten_sep,
	json_array_ingest_t json_array_ingest)
{
	int n = pjson_array->u.array.length;
	for (int i = 0; i < n; i++) {
		json_value_t* pjson_value = pjson_array->u.array.values[i];

		char free_flags = NO_FREE;
		char* json_key = low_int_to_string(i, &free_flags);
		char* lrec_key = mlr_paste_2_strings(prefix, json_key);
		if (free_flags)
			free(json_key);
		char* next_prefix = NULL;

		switch (pjson_value->type) {
		case JSON_NONE:
			lrec_put(prec, lrec_key, "", FREE_ENTRY_KEY);
			break;
		case JSON_NULL:
			lrec_put(prec, lrec_key, "", FREE_ENTRY_KEY);
			break;
		case JSON_STRING:
			lrec_put(prec, lrec_key, pjson_value->u.string.ptr, FREE_ENTRY_KEY);
			break;
		case JSON_BOOLEAN:
			lrec_put(prec, lrec_key, pjson_value->u.boolean.sval, FREE_ENTRY_KEY);
			break;
		case JSON_OBJECT:
			next_prefix = mlr_paste_2_strings(lrec_key, flatten_sep);
			if (!populate_from_nested_object(prec, pjson_value, next_prefix, flatten_sep, json_array_ingest))
				return FALSE;
			free(next_prefix);
			free(lrec_key);
			break;
		case JSON_ARRAY:
			switch (json_array_ingest) {
			case JSON_ARRAY_INGEST_FATAL:
				fprintf(stderr,
					"%s: found array item within JSON object. This is valid but unmillerable JSON.\n"
					"Use --json-skip-arrays-on-input to exclude these from input without fataling.\n"
					"Or, --json-map-arrays-on-input to convert them to integer-indexed maps.\n",
					MLR_GLOBALS.bargv0);
				return FALSE;
				break;
			case JSON_ARRAY_INGEST_AS_MAP:
				next_prefix = mlr_paste_2_strings(lrec_key, flatten_sep);
				if (!populate_from_nested_array(prec, pjson_value, next_prefix, flatten_sep, json_array_ingest)) {
					free(lrec_key);
					free(next_prefix);
					return FALSE;
				}
				free(lrec_key);
				free(next_prefix);
				break;
			// xxx other cases!
			default:
				free(lrec_key);
				break;
			}
			break;

		case JSON_INTEGER:
			lrec_put(prec, lrec_key, pjson_value->u.integer.sval, FREE_ENTRY_KEY);
			break;
		case JSON_DOUBLE:
			lrec_put(prec, lrec_key, pjson_value->u.dbl.sval, FREE_ENTRY_KEY);
			break;
		default:
			MLR_INTERNAL_CODING_ERROR();
			break;
		}

	}
	return TRUE;
}

// ----------------------------------------------------------------
// * The buffer is an entire JSON blob, e.g. contents from stdio read or mmap; peof-psof is the file size so peof is one
//   byte *after* the last valid file byte.
// * The buffer is not assumed to be null-terminated.
// * Any lines beginning with comment_string are modified by poking space characters up to line_term.
void mlr_json_strip_comments(char* psof, char* peof, char* comment_string, char* line_term) {
	int comment_string_len = strlen(comment_string);
	int line_term_len = strlen(line_term);
	int at_line_start = TRUE;
	for (char* p = psof; p < peof; /* increment in loop */) {
		if (streqn(p, line_term, line_term_len)) {
			p += line_term_len;
			at_line_start = TRUE;
		} else if (at_line_start && streqn(p, comment_string, comment_string_len)) {
			// Fill with spaces to end of line
			while (p < peof && !streqn(p, line_term, line_term_len)) {
				*p = ' ';
				p++;
			}
			at_line_start = TRUE;
		} else {
			at_line_start = FALSE;
			p++;
		}
	}
}
