#ifndef MLRVAL_H
#define MLRVAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
#include "../lib/mlrutil.h"
#include "../lib/mlrregex.h"
#include "../lib/free_flags.h"

// ================================================================
// MT for Miller type -- highly abbreviated here since these are
// spelled out a lot in rval_evaluators.c.
//
// ================================================================
// NOTE: mlrval functions invalidate their arguments. In particular, dynamically
// allocated strings input to these functions will either be freed, or will
// have their ownership transferred to the output mlrval.
//
// This is because the primary purpose of mlrvals is for evaluation of abstract
// syntax trees defined by the DSL for put and filter. Example AST:
//
//   $ mlr put -v '$z = $x . $y . "sum"' /dev/null
//   = (operator):
//       z (field_name).
//       . (operator):
//           . (operator):
//               x (field_name).
//               y (field_name).
//           sum (literal).
//
// * Given an lrec with fields named "x" and "y", there will be pointers to x
//   and y's field values from the input-data stream -- either to mmapped data
//   from a file, or pointers into dynamically allocated lines from stdio.
//
// * The from-field-name mlrvals for x and y values will point into lrec memory
//   but will have their own free-flags unset (since freeing of lrec memory is
//   the job of the lrec instance).
//
// * The dot operator will do any necessary freeing of the x and y mlrval
//   strings -- none in this case since they are direct references to field
//   values. The output of $x . $y, by contrast, will be dynamically
//   allocated.
//
// * The "sum" literal string is a pointer ultimately into argv[].
//   The from-literal mlrval will not have its free-flag set.
//
// * The concatenation of $x . $y and "sum" will dynamically allocated.
//   The $x . $y input string will be freed; the "sum" string won't be
//   since it wasn't owned by the from-literal mlrval.
//
// * The result of this outer concatenation will be stored in the $z field of
//   the current record, with ownership for the dynamically allocated string
//   transferred to the lrec instance.
//
// There is also some use of mlrvals in mixed float/int handling inside various
// mappers (e.g. stats1). There the use is much simpler: accumulation of
// numeric quantities, ultimately formatted as a string for output.
//
// ================================================================
//
// Many functions here use the naming convention x_yz_name or x_:
//
// * The first letter indicates return type.
//
// * The letters between the underscores indicate argument types, and their count indicates arity.
//
// * The following abbreviations apply:
//   o a: MT_ABSENT
//   o v: MT_EMPTY (v for void; e is for error)
//   o e: MT_ERROR
//   o b: MT_BOOLEAN
//   o f: MT_FLOAT
//   o i: MT_INT
//   o s: MT_STRING
//   o r: regular expression
//   o n: Numeric, i.e. MT_INT or MT_FLOAT
//   o x: any of the above.
//   o z: used for zero-argument functions, e.g. f_z_urand takes no arguments and returns MT_FLOAT.
//
// * If a function takes arguments of type x then that indicates it has a disposition vector/matrix
//   (or switch statements, or if-else statements) allowing it to handle various types.
//
// * If it takes arguments of type n then that indicates it is up to the caller to pass only numeric types.
//
// * If it takes arguments of type s then that indicates it is up to the caller to pass only strings.
//
// ================================================================


// Among other things, these defines are used in mlrval.c to index disposition matrices.
// So, if the numeric values are changed, all the matrices must be as well.

// Two kinds of null: absent (key not present in a record) and void (key present with empty value).
// Note void is an acceptable string (empty string) but not an acceptable number.
// Void-valued mlrvals have u.strv = "".
#define MT_ERROR    0 // E.g. error encountered in one eval & it propagates up the AST.
#define MT_ABSENT   1 // No such key, e.g. $z in 'x=,y=2'
#define MT_EMPTY    2 // Empty value, e.g. $x in 'x=,y=2'
#define MT_STRING   3
#define MT_INT      4
#define MT_FLOAT    5
#define MT_BOOLEAN  6
#define MT_DIM      7

typedef struct _mv_t {
	union {
		char*      strv;  // MT_STRING and MT_EMPTY
		long long  intv;  // MT_INT, and == 0 for MT_ABSENT and MT_ERROR
		double     fltv;  // MT_FLOAT
		int        boolv; // MT_BOOLEAN
	} u;
	unsigned char type;
	char free_flags;
} mv_t;

// ----------------------------------------------------------------
#define NULL_OR_ERROR_OUT_FOR_STRINGS(val) { \
	if ((val).type < MT_EMPTY) \
		return val; \
}

#define NULL_OR_ERROR_OUT_FOR_NUMBERS(val) { \
	if ((val).type <= MT_EMPTY) \
		return val; \
}

#define EMPTY_OR_ERROR_OUT_FOR_NUMBERS(val) { \
	if ((val).type == MT_ERROR || (val).type == MT_EMPTY) \
		return val; \
}

// ----------------------------------------------------------------
// CONSTRUCTORS

static inline mv_t mv_from_float(double d) {
	return (mv_t) {.type = MT_FLOAT, .free_flags = NO_FREE, .u.fltv = d};
}

static inline mv_t mv_from_int(long long i) {
	return (mv_t) {.type = MT_INT, .free_flags = NO_FREE, .u.intv = i};
}

static inline mv_t mv_from_bool(int b) {
	return (mv_t) {.type = MT_BOOLEAN, .free_flags = NO_FREE, .u.boolv = b};
}
static inline mv_t mv_from_true() {
	return (mv_t) {.type = MT_BOOLEAN, .free_flags = NO_FREE, .u.boolv = TRUE};
}
static inline mv_t mv_from_false() {
	return (mv_t) {.type = MT_BOOLEAN, .free_flags = NO_FREE, .u.boolv = FALSE};
}

static inline mv_t mv_from_string_with_free(char* s) {
	return (mv_t) {.type = MT_STRING, .free_flags = FREE_ENTRY_VALUE, .u.strv = s};
}
static inline mv_t mv_from_string_no_free(char* s) {
	return (mv_t) {.type = MT_STRING, .free_flags = NO_FREE, .u.strv = s};
}
static inline mv_t mv_from_string(char* s, char free_flags) {
	return (mv_t) {.type = MT_STRING, .free_flags = free_flags, .u.strv = s};
}

static inline mv_t mv_absent() { return (mv_t) {.type = MT_ABSENT, .free_flags = NO_FREE, .u.intv = 0};  }
static inline mv_t mv_empty()  { return (mv_t) {.type = MT_EMPTY,  .free_flags = NO_FREE, .u.strv = ""}; }
static inline mv_t mv_error()  { return (mv_t) {.type = MT_ERROR,  .free_flags = NO_FREE, .u.intv = 0};  }

static inline mv_t mv_copy(mv_t* pval) {
	if (pval->type == MT_STRING) {
		return mv_from_string_with_free(mlr_strdup_or_die(pval->u.strv));
	} else {
		return *pval;
	}
}

static inline mv_t* mv_alloc_copy(mv_t* pold) {
	mv_t* pnew = mlr_malloc_or_die(sizeof(mv_t));
	*pnew = mv_copy(pold);
	return pnew;
}

// ----------------------------------------------------------------
// DESTRUCTOR

static inline void mv_free(mv_t* pval) {
	if ((pval->type) == MT_STRING && (pval->free_flags & FREE_ENTRY_VALUE)) {
		free(pval->u.strv);
		pval->u.strv = NULL;
	}
	pval->type = MT_ABSENT;
}

// ----------------------------------------------------------------
// For stack-clear: set to ABSENT, freeing previous value if necessary
static inline void mv_reset(mv_t* pval) {
	if ((pval->type) == MT_STRING && (pval->free_flags & FREE_ENTRY_VALUE)) {
		free(pval->u.strv);
	}
	*pval = mv_absent();
}

// ----------------------------------------------------------------
// TYPE-TESTERS

static inline int mv_is_string_or_empty(mv_t* pval) {
	return pval->type == MT_STRING || pval->type == MT_EMPTY;
}
static inline int mv_is_numeric(mv_t* pval) {
	return pval->type == MT_INT || pval->type == MT_FLOAT;
}
static inline int mv_is_int(mv_t* pval) {
	return pval->type == MT_INT;
}
static inline int mv_is_float(mv_t* pval) {
	return pval->type == MT_FLOAT;
}
static inline int mv_is_boolean(mv_t* pval) {
	return pval->type == MT_BOOLEAN;
}
static inline int mv_is_string(mv_t* pval) {
	return pval->type == MT_STRING || pval->type == MT_EMPTY;
}
static inline int mv_is_error(mv_t* pval) {
	return pval->type == MT_ERROR;
}
static inline int mv_is_absent(mv_t* pval) {
	return pval->type == MT_ABSENT;
}
static inline int mv_is_present(mv_t* pval) {
	return pval->type != MT_ABSENT;
}
static inline int mv_is_empty(mv_t* pval) {
	return pval->type == MT_EMPTY || (pval->type == MT_STRING && *pval->u.strv == 0);
}
static inline int mv_is_not_empty(mv_t* pval) {
	return pval->type != MT_EMPTY;
}
static inline int mv_is_null(mv_t* pval) {
	return mv_is_absent(pval) || mv_is_empty(pval);
}
static inline int mv_is_null_or_error(mv_t* pval) {
	return mv_is_null(pval) || pval->type == MT_EMPTY;
}
static inline int mv_is_non_null(mv_t* pval) {
	return !mv_is_null(pval);
}

// ----------------------------------------------------------------
// AUXILIARY METHODS

char* mt_describe_type(int type);
char* mt_describe_type_simple(int type);

// Allocates memory which the caller must free; does not modify the mlrval.
// Returns no reference to the mlrval's data.  Suitable for getting data out of
// a mlrval which might be about to be freed.
char* mv_alloc_format_val(mv_t* pval);
char* mv_alloc_format_val_quoting_strings(mv_t* pval);

// Returns a reference to the mlrval's data if the mlrval is MT_STRING.
// Does not modify the mlrval. Suitable only for read-only string-formatting
// of the mlrval while it still exists and hasn't been freed yet.
char* mv_maybe_alloc_format_val(mv_t* pval, char* pfree_flags);

// If the mlrval is MT_STRING, returns that and invalidates the argument.
// This is suitable for baton-pass-out (end of evaluation chain).
char* mv_format_val(mv_t* pval, char* pfree_flags);

// Output string includes type and value information (e.g. for debug).
// The caller must free the return value.
char* mv_describe_val(mv_t val);

void mv_set_boolean_strict(mv_t* pval);
void mv_set_float_strict(mv_t* pval);
void mv_set_float_nullable(mv_t* pval);
void mv_set_int_nullable(mv_t* pval);

// int or float:
void mv_set_number_nullable(mv_t* pval);
mv_t mv_scan_number_nullable(char* string);
mv_t mv_scan_number_or_die(char* string);

// Each of the following three
// Type-inferencing for the following three functions, respectively:
//   "x" => "x", "3" => "3"
//   "x" => "x", "3" => 3.0
//   "x" => "x", "3" => 3
// In common to all three:
// * Null string -> mv_absent
// * Empty string -> mv_empty
// * Non-numeric -> string-valued mlrval with storage pointing
//   to the char* (no copy is done).
mv_t mv_ref_type_infer_string(char* string);
mv_t mv_ref_type_infer_string_or_float(char* string);
mv_t mv_ref_type_infer_string_or_float_or_int(char* string);
mv_t mv_copy_type_infer_string_or_float_or_int(char* string); // strdups if retval is MT_STRING

#endif // MLRVAL_H
