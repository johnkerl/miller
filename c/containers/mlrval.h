#ifndef MLR_VAL_H
#define MLR_VAL_H

#include <math.h>
#include <string.h>
#include <ctype.h>
#include <regex.h>
#include "lib/mlrmath.h"
#include "lib/mlrutil.h"
#include "lib/mtrand.h"
#include "lib/string_builder.h"
#include "lib/string_array.h"
#include "containers/free_flags.h"

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
//   o b: MT_BOOL
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
#define MT_BOOL     6
#define MT_DIM      7

#define MV_SB_ALLOC_LENGTH 32

#define ISO8601_TIME_FORMAT "%Y-%m-%dT%H:%M:%SZ"
#define ISO8601_DATE_FORMAT "%Y-%m-%d"

typedef struct _mv_t {
	union {
		char*      strv;  // MT_STRING and MT_EMPTY
		long long  intv;  // MT_INT, and == 0 for MT_ABSENT and MT_ERROR
		double     fltv;  // MT_FLOAT
		int        boolv; // MT_BOOL
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
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = b};
}
static inline mv_t mv_from_true() {
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = TRUE};
}
static inline mv_t mv_from_false() {
	return (mv_t) {.type = MT_BOOL, .free_flags = NO_FREE, .u.boolv = FALSE};
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
}

// ----------------------------------------------------------------
// TYPE-TESTERS

static inline int mv_is_string_or_empty(mv_t* pval) {
	return pval->type == MT_STRING || pval->type == MT_EMPTY;
}
static inline int mv_is_numeric(mv_t* pval) {
	return pval->type == MT_INT || pval->type == MT_FLOAT;
}
static inline int mv_is_null(mv_t* pval) {
	return MT_ERROR < pval->type && pval->type <= MT_EMPTY;
}
static inline int mv_is_null_or_error(mv_t* pval) {
	return pval->type <= MT_EMPTY;
}
static inline int mv_is_non_null(mv_t* pval) {
	return MT_ERROR < pval->type && pval->type > MT_EMPTY;
}
static inline int mv_is_absent(mv_t* pval) {
	return pval->type == MT_ABSENT;
}
static inline int mv_is_present(mv_t* pval) {
	return pval->type != MT_ABSENT;
}
static inline int mv_is_empty(mv_t* pval) {
	return pval->type == MT_EMPTY;
}
static inline int mv_is_not_empty(mv_t* pval) {
	return pval->type != MT_EMPTY;
}

// ----------------------------------------------------------------
// AUXILIARY METHODS

char* mt_describe_type(int type);
char* mt_describe_type_simple(int type);

// Allocates memory which the caller must free; does not modify the mlrval.
// Returns no reference to the mlrval's data.  Suitable for getting data out of
// a mlrval which might be about to be freed.
char* mv_alloc_format_val(mv_t* pval);

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

// ----------------------------------------------------------------
// FUNCTION-OF-MLRVAL TYPES

typedef mv_t mv_zary_func_t();
typedef mv_t mv_unary_func_t(mv_t* pval1);
typedef mv_t mv_binary_func_t(mv_t* pval1, mv_t* pval2);
typedef mv_t mv_binary_arg3_capture_func_t(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
typedef mv_t mv_binary_arg2_regex_func_t(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);
typedef mv_t mv_ternary_func_t(mv_t* pval1, mv_t* pval2, mv_t* pval3);
typedef mv_t mv_ternary_arg2_regex_func_t(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);

// ----------------------------------------------------------------
// FUNCTIONS OF MLRVALS

static inline mv_t b_b_not_func(mv_t* pval1) {
	return mv_from_bool(!pval1->u.boolv);
}

static inline mv_t b_bb_or_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv || pval2->u.boolv);
}
static inline mv_t b_bb_and_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv && pval2->u.boolv);
}
static inline mv_t b_bb_xor_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_bool(pval1->u.boolv ^ pval2->u.boolv);
}

// ----------------------------------------------------------------
static inline mv_t f_z_urand_func() {
	return mv_from_float(get_mtrand_double()); // mtrand.h
}
static inline mv_t i_z_urand32_func() {
	return mv_from_float(get_mtrand_int32()); // mtrand.h
}
static inline mv_t f_z_systime_func() {
	return mv_from_float(get_systime()); // mlrutil.h
}

// ----------------------------------------------------------------
static inline mv_t f_f_acos_func(mv_t*     pval1) {return mv_from_float( acos     (pval1->u.fltv));}
static inline mv_t f_f_acosh_func(mv_t*    pval1) {return mv_from_float( acosh    (pval1->u.fltv));}
static inline mv_t f_f_asin_func(mv_t*     pval1) {return mv_from_float( asin     (pval1->u.fltv));}
static inline mv_t f_f_asinh_func(mv_t*    pval1) {return mv_from_float( asinh    (pval1->u.fltv));}
static inline mv_t f_f_atan_func(mv_t*     pval1) {return mv_from_float( atan     (pval1->u.fltv));}
static inline mv_t f_f_atanh_func(mv_t*    pval1) {return mv_from_float( atanh    (pval1->u.fltv));}
static inline mv_t f_f_cbrt_func(mv_t*     pval1) {return mv_from_float( cbrt     (pval1->u.fltv));}
static inline mv_t f_f_cos_func(mv_t*      pval1) {return mv_from_float( cos      (pval1->u.fltv));}
static inline mv_t f_f_cosh_func(mv_t*     pval1) {return mv_from_float( cosh     (pval1->u.fltv));}
static inline mv_t f_f_erf_func(mv_t*      pval1) {return mv_from_float( erf      (pval1->u.fltv));}
static inline mv_t f_f_erfc_func(mv_t*     pval1) {return mv_from_float( erfc     (pval1->u.fltv));}
static inline mv_t f_f_exp_func(mv_t*      pval1) {return mv_from_float( exp      (pval1->u.fltv));}
static inline mv_t f_f_expm1_func(mv_t*    pval1) {return mv_from_float( expm1    (pval1->u.fltv));}
static inline mv_t f_f_invqnorm_func(mv_t* pval1) {return mv_from_float( invqnorm (pval1->u.fltv));}
static inline mv_t f_f_log10_func(mv_t*    pval1) {return mv_from_float( log10    (pval1->u.fltv));}
static inline mv_t f_f_log1p_func(mv_t*    pval1) {return mv_from_float( log1p    (pval1->u.fltv));}
static inline mv_t f_f_log_func(mv_t*      pval1) {return mv_from_float( log      (pval1->u.fltv));}
static inline mv_t f_f_qnorm_func(mv_t*    pval1) {return mv_from_float( qnorm    (pval1->u.fltv));}
static inline mv_t f_f_sin_func(mv_t*      pval1) {return mv_from_float( sin      (pval1->u.fltv));}
static inline mv_t f_f_sinh_func(mv_t*     pval1) {return mv_from_float( sinh     (pval1->u.fltv));}
static inline mv_t f_f_sqrt_func(mv_t*     pval1) {return mv_from_float( sqrt     (pval1->u.fltv));}
static inline mv_t f_f_tan_func(mv_t*      pval1) {return mv_from_float( tan      (pval1->u.fltv));}
static inline mv_t f_f_tanh_func(mv_t*     pval1) {return mv_from_float( tanh     (pval1->u.fltv));}

static inline mv_t f_ff_pow_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_float(pow(pval1->u.fltv, pval2->u.fltv));
}

mv_t x_xx_plus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_minus_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_times_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_divide_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_int_divide_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_mod_func(mv_t* pval1, mv_t* pval2);
mv_t x_x_upos_func(mv_t* pval1);
mv_t x_x_uneg_func(mv_t* pval1);

// Bitwise
mv_t x_xx_bxor_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_band_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_bor_func(mv_t* pval1, mv_t* pval2);

mv_t x_x_abs_func(mv_t* pval1);
mv_t x_x_ceil_func(mv_t* pval1);
mv_t x_x_floor_func(mv_t* pval1);
mv_t x_x_round_func(mv_t* pval1);
mv_t x_x_sgn_func(mv_t* pval1);

mv_t x_xx_min_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_max_func(mv_t* pval1, mv_t* pval2);
mv_t x_xx_roundm_func(mv_t* pval1, mv_t* pval2);

mv_t b_x_isnull_func(mv_t* pval1);
mv_t b_x_isnotnull_func(mv_t* pval1);
mv_t b_x_isabsent_func(mv_t* pval1);
mv_t b_x_ispresent_func(mv_t* pval1);
mv_t b_x_isempty_func(mv_t* pval1);
mv_t b_x_isnotempty_func(mv_t* pval1);

mv_t b_x_isnumeric_func(mv_t* pval1);
mv_t b_x_isint_func(mv_t* pval1);
mv_t b_x_isfloat_func(mv_t* pval1);
mv_t b_x_isbool_func(mv_t* pval1);
mv_t b_x_isstring_func(mv_t* pval1);

mv_t i_x_int_func(mv_t* pval1);
mv_t f_x_float_func(mv_t* pval1);
mv_t b_x_boolean_func(mv_t* pval1);
mv_t s_x_string_func(mv_t* pval1);
mv_t s_x_hexfmt_func(mv_t* pval1);
mv_t s_xs_fmtnum_func(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
static inline mv_t f_ff_atan2_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_float(atan2(pval1->u.fltv, pval2->u.fltv));
}

static inline mv_t f_fff_logifit_func(mv_t* pval1, mv_t* pval2, mv_t* pval3) {
	double x = pval1->u.fltv;
	double m = pval2->u.fltv;
	double b = pval3->u.fltv;
	return mv_from_float(1.0 / (1.0 + exp(-m*x-b)));
}

static inline mv_t i_ii_urandint_func(mv_t* pval1, mv_t* pval2) {
	long long a = pval1->u.intv;
	long long b = pval2->u.intv;
	long long lo, hi;
	if (a <= b) {
		lo = a;
		hi = b + 1;
	} else {
		lo = b;
		hi = a + 1;
	}
	long long u  = lo + (hi - lo) * get_mtrand_double();
	return mv_from_int(u);
}

static inline mv_t i_ii_bitwise_lsh_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv << pval2->u.intv);
}
static inline mv_t i_ii_bitwise_rsh_func(mv_t* pval1, mv_t* pval2) {
	return mv_from_int(pval1->u.intv >> pval2->u.intv);
}
static inline mv_t i_i_bitwise_not_func(mv_t* pval1) {
	return mv_from_int(~pval1->u.intv);
}

mv_t i_iii_modadd_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modsub_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modmul_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t i_iii_modexp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_s_tolower_func(mv_t* pval1);
mv_t s_s_toupper_func(mv_t* pval1);
mv_t i_s_strlen_func(mv_t* pval1);
mv_t s_x_typeof_func(mv_t* pval1);

mv_t s_ss_dot_func(mv_t* pval1, mv_t* pval2);

mv_t sub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t sub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);
mv_t gsub_no_precomp_func(mv_t* pval1, mv_t* pval2, mv_t* pval3);
mv_t gsub_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, mv_t* pval3);

// ----------------------------------------------------------------
mv_t s_n_sec2gmt_func(mv_t* pval1);
mv_t s_n_sec2gmtdate_func(mv_t* pval1);
mv_t i_s_gmt2sec_func(mv_t* pval1);
mv_t s_ns_strftime_func(mv_t* pval1, mv_t* pval2);
mv_t i_ss_strptime_func(mv_t* pval1, mv_t* pval2);

mv_t s_i_sec2hms_func(mv_t* pval1);
mv_t s_f_fsec2hms_func(mv_t* pval1);
mv_t s_i_sec2dhms_func(mv_t* pval1);
mv_t s_f_fsec2dhms_func(mv_t* pval1);
mv_t i_s_hms2sec_func(mv_t* pval1);
mv_t f_s_hms2fsec_func(mv_t* pval1);
mv_t i_s_dhms2sec_func(mv_t* pval1);
mv_t f_s_dhms2fsec_func(mv_t* pval1);

mv_t time_string_from_seconds(mv_t* psec, char* format);

// ----------------------------------------------------------------
// arg2 evaluates to string via compound expression; regexes compiled on each call
mv_t matches_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
mv_t does_not_match_no_precomp_func(mv_t* pval1, mv_t* pval2, string_array_t** ppregex_captures);
// arg2 is a string, compiled to regex only once at alloc time
mv_t matches_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);
mv_t does_not_match_precomp_func(mv_t* pval1, regex_t* pregex, string_builder_t* psb, string_array_t** ppregex_captures);

// For filter/put DSL:
mv_t eq_op_func(mv_t* pval1, mv_t* pval2);
mv_t ne_op_func(mv_t* pval1, mv_t* pval2);
mv_t gt_op_func(mv_t* pval1, mv_t* pval2);
mv_t ge_op_func(mv_t* pval1, mv_t* pval2);
mv_t lt_op_func(mv_t* pval1, mv_t* pval2);
mv_t le_op_func(mv_t* pval1, mv_t* pval2);

// Assumes inputs are MT_STRING or MT_INT. Nominally intended for mlhmmv which uses only string/int mlrvals.
int mv_equals_si(mv_t* pa, mv_t* pb);

// For non-DSL comparison of mlrvals:
int mv_i_nn_eq(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ne(mv_t* pval1, mv_t* pval2);
int mv_i_nn_gt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_ge(mv_t* pval1, mv_t* pval2);
int mv_i_nn_lt(mv_t* pval1, mv_t* pval2);
int mv_i_nn_le(mv_t* pval1, mv_t* pval2);

// For unit-test keystroke-saving:
int mveq(mv_t* pval1, mv_t* pval2);

// ----------------------------------------------------------------
// For qsort of numeric mlrvals.
int mv_nn_comparator(const void* pva, const void* pvb);

int mlr_bsearch_mv_n_for_insert(mv_t* array, int size, mv_t* pvalue);

#endif // MLR_VAL_H
