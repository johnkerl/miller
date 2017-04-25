#ifndef XVFUNCS_H
#define XVFUNCS_H

// ================================================================
// Functions on extended values, namely, scalars and maps.
//
// NOTE ON EPHEMERALITY OF MAPS:
//
// Most functions here free their inputs. E.g. for string concatenation, the
// output which is returned is the concatenation of the two inputs which are
// freed. This is true for functions on scalars as well as functions on maps
// (the latter in this header file).
//
// However, maps are treated differently from scalars in that some maps are
// referenced, rather than copied, within the concrete syntax tree.  E.g. in
// 'mapsum({3:4},@v)', the map literal {3:4} is ephemeral to the expression and
// must be freed during evaluation, but the @v part is referenced to the @v data
// structure and is only copied on write.  The boxed_xval_t's is_ephemeral flag
// tracks the difference between the two cases.
// ================================================================

#include "../lib/mlrutil.h"
#include "../containers/mlhmmv.h"
#include "../containers/boxed_xval.h"

// ----------------------------------------------------------------
typedef boxed_xval_t xv_variadic_func_t(
	boxed_xval_t* pbxvals,
	int           nxvals);

typedef boxed_xval_t xv_zary_func_t();

typedef boxed_xval_t xv_unary_func_t(
	boxed_xval_t* pbxval1);

typedef boxed_xval_t xv_binary_func_t(
	boxed_xval_t* pbxval1,
	boxed_xval_t* pbxval2);

typedef boxed_xval_t xv_ternary_func_t(
	boxed_xval_t* pbxval1,
	boxed_xval_t* pbxval2,
	boxed_xval_t* pbxval3);

// ----------------------------------------------------------------
// Most functions here free their inputs. E.g. for string concatenation, the
// output which is returned is the concatenation of the two inputs which are
// freed. For another example, is_string frees its input and returns the boolean
// value of the result. These functions, by contrast, only return a boolean for
// the outcome of the test but do not free the inputs. The intended usage is for
// type-assertion checks.  E.g. in '$b = asserting_string($a)', if $a is a
// string it is assigned to $b, else an error is thrown.

static inline boxed_xval_t b_x_is_present_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal || mv_is_present(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_absent_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_absent(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_map_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal
		)
	);
}

static inline boxed_xval_t b_x_is_not_map_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal
		)
	);
}

static inline boxed_xval_t b_x_is_numeric_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_numeric(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_int_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_int(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_float_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_float(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_boolean_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_boolean(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_string_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_string(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_null_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_null(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_not_null_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!(pbxval1->xval.is_terminal && mv_is_null(&pbxval1->xval.terminal_mlrval))
		)
	);
}

static inline boxed_xval_t b_x_is_empty_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_empty(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_is_not_empty_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!(pbxval1->xval.is_terminal && mv_is_empty(&pbxval1->xval.terminal_mlrval))
		)
	);
}

static inline boxed_xval_t b_x_is_empty_map_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal && pbxval1->xval.pnext_level->num_occupied == 0
		)
	);
}

static inline boxed_xval_t b_x_is_nonempty_map_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal && pbxval1->xval.pnext_level->num_occupied != 0
		)
	);
}

static inline boxed_xval_t s_x_typeof_no_free_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
	    mv_from_string(
			mlhmmv_xvalue_describe_type_simple(&pbxval1->xval), NO_FREE
		)
	);
}

// ----------------------------------------------------------------
static inline boxed_xval_t b_x_is_present_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_present_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_absent_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_absent_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_map_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_map_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_not_map_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_not_map_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_numeric_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_numeric_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_int_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_int_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_float_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_float_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_boolean_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_boolean_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_string_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_string_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_null_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_null_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_not_null_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_not_null_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_empty_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_empty_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_not_empty_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_not_empty_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_empty_map_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_empty_map_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t b_x_is_nonempty_map_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = b_x_is_nonempty_map_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

static inline boxed_xval_t s_x_typeof_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = s_x_typeof_no_free_xfunc(pbxval1);
	if (pbxval1->is_ephemeral)
	  mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
}

// ----------------------------------------------------------------
boxed_xval_t b_xx_haskey_xfunc(
	boxed_xval_t* pmapval,
	boxed_xval_t* pkeyval);

boxed_xval_t i_x_length_xfunc(
	boxed_xval_t* pbxval1);

boxed_xval_t i_x_depth_xfunc(
	boxed_xval_t* pbxval1);

boxed_xval_t i_x_leafcount_xfunc(
	boxed_xval_t* pbxval1);

boxed_xval_t variadic_mapsum_xfunc(
	boxed_xval_t* pbxvals, int nxvals);

boxed_xval_t variadic_mapdiff_xfunc(
	boxed_xval_t* pbxvals, int nxvals);

boxed_xval_t variadic_mapexcept_xfunc(
	boxed_xval_t* pbxvals, int nxvals);

boxed_xval_t variadic_maponly_xfunc(
	boxed_xval_t* pbxvals, int nxvals);

boxed_xval_t m_ss_splitnv_xfunc(
	boxed_xval_t* pstringval,
	boxed_xval_t* psepval);

boxed_xval_t m_ss_splitnvx_xfunc(
	boxed_xval_t* pstringval,
	boxed_xval_t* psepval);

boxed_xval_t m_sss_splitkv_xfunc(
	boxed_xval_t* pstringval,
	boxed_xval_t* ppairsepval,
	boxed_xval_t* plistsepval);

boxed_xval_t m_sss_splitkvx_xfunc(
	boxed_xval_t* pstringval,
	boxed_xval_t* ppairsepval,
	boxed_xval_t* plistsepval);

boxed_xval_t s_ms_joink_xfunc(
	boxed_xval_t* pmapval,
	boxed_xval_t* psepval);

boxed_xval_t s_ms_joinv_xfunc(
	boxed_xval_t* pmapval,
	boxed_xval_t* psepval);

boxed_xval_t s_mss_joinkv_xfunc(
	boxed_xval_t* pmapval,
	boxed_xval_t* ppairsepval,
	boxed_xval_t* plistsepval);

#endif // XVFUNCS_H
