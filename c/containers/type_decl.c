#include "containers/type_decl.h"

// ----------------------------------------------------------------
char* type_mask_to_desc(int type_mask) {
	switch(type_mask) {
	case TYPE_MASK_ERROR:   return "error";   break;
	case TYPE_MASK_ANY:     return "any";     break;
	case TYPE_MASK_ABSENT:  return "absent";  break;
	case TYPE_MASK_EMPTY:   return "empty";   break;
	case TYPE_MASK_STRING:  return "string";  break;
	case TYPE_MASK_INT:     return "int";     break;
	case TYPE_MASK_FLOAT:   return "float";   break;
	case TYPE_MASK_BOOLEAN: return "boolean"; break;
	case TYPE_MASK_NUMERIC: return "numeric"; break;
	case TYPE_MASK_MAP:     return "map";     break;
	default:                return "???";     break;
	}
}

// ----------------------------------------------------------------
int type_mask_from_name(char* name) {
	if      (streq(name, "error"))   { return TYPE_MASK_ERROR;   }
	else if (streq(name, "absent"))  { return TYPE_MASK_ABSENT;  }
	else if (streq(name, "empty"))   { return TYPE_MASK_EMPTY;   }
	else if (streq(name, "string"))  { return TYPE_MASK_STRING;  }
	else if (streq(name, "str"))     { return TYPE_MASK_STRING;  }
	else if (streq(name, "int"))     { return TYPE_MASK_INT;     }
	else if (streq(name, "float"))   { return TYPE_MASK_FLOAT;   }
	else if (streq(name, "numeric")) { return TYPE_MASK_NUMERIC; }
	else if (streq(name, "num"))     { return TYPE_MASK_NUMERIC; }
	else if (streq(name, "boolean")) { return TYPE_MASK_BOOLEAN; }
	else if (streq(name, "bool"))    { return TYPE_MASK_BOOLEAN; }
	else if (streq(name, "map"))     { return TYPE_MASK_MAP;     }
	else if (streq(name, "var"))     { return TYPE_MASK_ANY;     }
	else if (streq(name, "any"))     { return TYPE_MASK_ANY;     }

	else {
		MLR_INTERNAL_CODING_ERROR();
		return 0; // not reached
	}
}
