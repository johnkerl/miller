#ifndef TYPE_INFERENCE_H
#define TYPE_INFERENCE_H

// Intentionally not an enum
#define TYPE_INFER_STRING_FLOAT_INT 0xce08
#define TYPE_INFER_STRING_FLOAT     0xce09
#define TYPE_INFER_STRING_ONLY      0xce0a

static inline char* type_inferencing_to_string(int ti) {
	switch(ti) {
	case TYPE_INFER_STRING_FLOAT_INT:
		return "TYPE_INFER_STRING_FLOAT_INT";
		break;
	case TYPE_INFER_STRING_FLOAT:
		return "TYPE_INFER_STRING_FLOAT";
		break;
	case TYPE_INFER_STRING_ONLY:
		return "TYPE_INFER_STRING_ONLY";
		break;
	default:
		return "???";
		break;
	}
}

#endif // TYPE_INFERENCE_H
