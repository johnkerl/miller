#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include <time.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/mlrval.h"

// ================================================================
// See important notes at the top of mlrval.h.
// ================================================================

// ----------------------------------------------------------------
char* mt_describe_type(int type) {
	switch (type) {
	case     MT_ERROR:   return "MT_ERROR";   break;
	case     MT_ABSENT:  return "MT_ABSENT";  break;
	case     MT_EMPTY:   return "MT_EMPTY";   break;
	case     MT_STRING:  return "MT_STRING";  break;
	case     MT_INT:     return "MT_INT";     break;
	case     MT_FLOAT:   return "MT_FLOAT";   break;
	case     MT_BOOLEAN: return "MT_BOOLEAN"; break;
	default:             return "???";        break;
	}
}

char* mt_describe_type_simple(int type) {
	switch (type) {
	case MT_ERROR:   return "error";  break;
	case MT_ABSENT:  return "absent"; break;
	case MT_EMPTY:   return "empty";  break;
	case MT_STRING:  return "string"; break;
	case MT_INT:     return "int";    break;
	case MT_FLOAT:   return "float";  break;
	case MT_BOOLEAN: return "bool";   break;
	default:         return "???";    break;
	}
}

// ----------------------------------------------------------------
// See comments in header file
char* mv_alloc_format_val(mv_t* pval) {
	switch(pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		return mlr_strdup_or_die("");
		break;
	case MT_ERROR:
		return mlr_strdup_or_die("(error)");
		break;
	case MT_STRING:
		return mlr_strdup_or_die(pval->u.strv);
		break;
	case MT_BOOLEAN:
		return mlr_strdup_or_die(pval->u.boolv ? "true" : "false");
		break;
	case MT_FLOAT:
		return mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		return mlr_alloc_string_from_ll(pval->u.intv);
		break;
	default:
		return mlr_strdup_or_die("???");
		break;
	}
}

char* mv_alloc_format_val_quoting_strings(mv_t* pval) {
	switch(pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		return mlr_strdup_quoted_or_die("");
		break;
	case MT_ERROR:
		return mlr_strdup_or_die("(error)");
		break;
	case MT_STRING:
		return mlr_strdup_quoted_or_die(pval->u.strv);
		break;
	case MT_BOOLEAN:
		return mlr_strdup_or_die(pval->u.boolv ? "true" : "false");
		break;
	case MT_FLOAT:
		return mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		return mlr_alloc_string_from_ll(pval->u.intv);
		break;
	default:
		return mlr_strdup_or_die("???");
		break;
	}
}

// ----------------------------------------------------------------
// See comments in header file
char* mv_maybe_alloc_format_val(mv_t* pval, char* pfree_flags) {
	switch(pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		*pfree_flags = NO_FREE;
		return "";
		break;
	case MT_ERROR:
		*pfree_flags = NO_FREE;
		return "(error)";
		break;
	case MT_STRING:
		*pfree_flags = NO_FREE;
		return pval->u.strv;
		break;
	case MT_BOOLEAN:
		*pfree_flags = NO_FREE;
		return pval->u.boolv ? "true" : "false";
		break;
	case MT_FLOAT:
		*pfree_flags = FREE_ENTRY_VALUE;
		return mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		*pfree_flags = FREE_ENTRY_VALUE;
		return mlr_alloc_string_from_ll(pval->u.intv);
		break;
	default:
		*pfree_flags = NO_FREE;
		return "???";
		break;
	}
}

// ----------------------------------------------------------------
// See comments in header file
char* mv_format_val(mv_t* pval, char* pfree_flags) {
	char* rv = NULL;
	switch(pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		*pfree_flags = NO_FREE;
		rv = "";
		break;
	case MT_ERROR:
		*pfree_flags = NO_FREE;
		rv = "(error)";
		break;
	case MT_STRING:
		// Ownership transfer to the caller
		*pfree_flags = pval->free_flags;
		rv = pval->u.strv;
		*pval = mv_empty();
		break;
	case MT_BOOLEAN:
		*pfree_flags = NO_FREE;
		rv = pval->u.boolv ? "true" : "false";
		break;
	case MT_FLOAT:
		*pfree_flags = FREE_ENTRY_VALUE;
		rv = mlr_alloc_string_from_double(pval->u.fltv, MLR_GLOBALS.ofmt);
		break;
	case MT_INT:
		*pfree_flags = FREE_ENTRY_VALUE;
		rv = mlr_alloc_string_from_ll(pval->u.intv);
		break;
	default:
		*pfree_flags = NO_FREE;
		rv = "???";
		break;
	}
	return rv;
}

// ----------------------------------------------------------------
// See comments in header file
char* mv_describe_val(mv_t val) {
	char* stype = mt_describe_type(val.type);
	char* strv  = mv_alloc_format_val(&val);
	char* desc  = mlr_malloc_or_die(strlen(stype) + strlen(strv) + 4);
	sprintf(desc, "[%s] %s", stype, strv);
	free(strv);
	return desc;
}

// ----------------------------------------------------------------
void mv_set_boolean_strict(mv_t* pval) {
	if (pval->type != MT_BOOLEAN) {
		char* desc = mt_describe_type(pval->type);
		fprintf(stderr, "Expression does not evaluate to boolean: got %s.\n", desc);
		exit(1);
	}
}

// ----------------------------------------------------------------
void mv_set_float_strict(mv_t* pval) {
	double fltv = 0.0;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		break;
	case MT_ERROR:
		break;
	case MT_FLOAT:
		break;
	case MT_STRING:
		if (!mlr_try_float_from_string(pval->u.strv, &fltv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_float(fltv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	case MT_INT:
		pval ->type = MT_FLOAT;
		pval->u.fltv = (double)pval->u.intv;
		break;
	case MT_BOOLEAN:
		pval->type = MT_ERROR;
		pval->u.intv = 0LL;
		break;
	default:
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_float_nullable(mv_t* pval) {
	double fltv = 0.0;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		break;
	case MT_ERROR:
		break;
	case MT_FLOAT:
		break;
	case MT_INT:
		pval ->type = MT_FLOAT;
		pval->u.fltv = (double)pval->u.intv;
		break;
	case MT_BOOLEAN:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		if (*pval->u.strv == '\0') {
			nval = mv_empty();
		} else if (!mlr_try_float_from_string(pval->u.strv, &fltv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_float(fltv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	default:
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_int_nullable(mv_t* pval) {
	long long intv = 0LL;
	mv_t nval = mv_error();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		break;
	case MT_ERROR:
		break;
	case MT_INT:
		break;
	case MT_FLOAT:
		pval ->type = MT_INT;
		pval->u.intv = (long long)pval->u.fltv;
		break;
	case MT_BOOLEAN:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		if (*pval->u.strv == '\0') {
			nval = mv_empty();
		} else if (!mlr_try_int_from_string(pval->u.strv, &intv)) {
			// keep nval = mv_error()
		} else {
			nval = mv_from_int(intv);
		}
		mv_free(pval);
		*pval = nval;
		break;
	default:
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
}

// ----------------------------------------------------------------
void mv_set_number_nullable(mv_t* pval) {
	mv_t nval = mv_empty();
	switch (pval->type) {
	case MT_ABSENT:
	case MT_EMPTY:
		break;
	case MT_ERROR:
		break;
	case MT_INT:
		break;
	case MT_FLOAT:
		break;
	case MT_BOOLEAN:
		pval->type = MT_ERROR;
		pval->u.intv = 0;
		break;
	case MT_STRING:
		nval = mv_scan_number_nullable(pval->u.strv);
		mv_free(pval);
		*pval = nval;
		break;
	default:
		MLR_INTERNAL_CODING_ERROR();
		break;
	}
}

mv_t mv_scan_number_nullable(char* string) {
	double fltv = 0.0;
	long long intv = 0LL;
	mv_t rv = mv_empty();
	if (*string == '\0') {
		// keep rv = mv_empty();
	} else if (mlr_try_int_from_string(string, &intv)) {
		rv = mv_from_int(intv);
	} else if (mlr_try_float_from_string(string, &fltv)) {
		rv = mv_from_float(fltv);
	} else {
		rv = mv_error();
	}
	return rv;
}

mv_t mv_scan_number_or_die(char* string) {
	mv_t rv = mv_scan_number_nullable(string);
	if (!mv_is_numeric(&rv)) {
		fprintf(stderr, "%s: couldn't parse \"%s\" as number.\n",
			MLR_GLOBALS.bargv0, string);
		exit(1);
	}
	return rv;
}

// ----------------------------------------------------------------
mv_t mv_ref_type_infer_string(char* string) {
	if (string == NULL) {
		return mv_absent();
	} else if (*string == 0) {
		return mv_empty();
	} else {
		return mv_from_string(string, NO_FREE);
	}
}

mv_t mv_ref_type_infer_string_or_float(char* string) {
	if (string == NULL) {
		return mv_absent();
	} else if (*string == 0) {
		return mv_empty();
	} else {
		double fltv;
		if (mlr_try_float_from_string(string, &fltv)) {
			return mv_from_float(fltv);
		} else {
			return mv_from_string(string, NO_FREE);
		}
	}
}

mv_t mv_ref_type_infer_string_or_float_or_int(char* string) {
	if (string == NULL) {
		return mv_absent();
	} else if (*string == 0) {
		return mv_empty();
	} else {
		long long intv;
		double fltv;
		if (mlr_try_int_from_string(string, &intv)) {
			return mv_from_int(intv);
		} else if (mlr_try_float_from_string(string, &fltv)) {
			return mv_from_float(fltv);
		} else {
			return mv_from_string(string, NO_FREE);
		}
	}
}

mv_t mv_copy_type_infer_string_or_float_or_int(char* string) {
	if (string == NULL) {
		return mv_absent();
	} else if (*string == 0) {
		return mv_empty();
	} else {
		long long intv;
		double fltv;
		if (mlr_try_int_from_string(string, &intv)) {
			return mv_from_int(intv);
		} else if (mlr_try_float_from_string(string, &fltv)) {
			return mv_from_float(fltv);
		} else {
			return mv_from_string(mlr_strdup_or_die(string), FREE_ENTRY_VALUE);
		}
	}
}
