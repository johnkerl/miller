#include "containers/type_decl.h"

// ----------------------------------------------------------------
char* type_mask_to_desc(int type_mask) {
	switch(type_mask) {
	case TYPE_MASK_ERROR:   return "error";   break;
	case TYPE_MASK_ABSENT:  return "absent";  break;
	case TYPE_MASK_EMPTY:   return "empty";   break;
	case TYPE_MASK_STRING:  return "string";  break;
	case TYPE_MASK_INT:     return "int";     break;
	case TYPE_MASK_FLOAT:   return "float";   break;
	case TYPE_MASK_BOOLEAN: return "boolean"; break;
	case TYPE_MASK_MAP:     return "map";     break;
	case TYPE_MASK_NUMERIC: return "numeric"; break;
	case TYPE_MASK_PRESENT: return "present"; break;
	case TYPE_MASK_ANY:     return "any";     break;
	default:                return "???";     break;
	}
}
