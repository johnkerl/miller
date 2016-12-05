#ifndef XVFUNCS_H
#define XVFUNCS_H

// ================================================================
// Functions on extended values, namely, mlrvals/hashmaps.
// ================================================================

// xxx need memory-transfer semantics
// xxx make the xvfuncs API entirely in terms of boxed_xval_t's?

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

// ----------------------------------------------------------------
// xxx hook all into fmgr

static inline boxed_xval_t b_x_ispresent_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal || mv_is_present(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isabsent_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_absent(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_ismap_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal
		)
	);
}

static inline boxed_xval_t b_x_isscalar_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_present(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isnumeric_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_numeric(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isint_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_int(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isfloat_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_float(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isboolean_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_boolean(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isstring_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_string(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isnull_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_null(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isnotnull_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!(pbxval1->xval.is_terminal && mv_is_null(&pbxval1->xval.terminal_mlrval))
		)
	);
}

static inline boxed_xval_t b_x_isempty_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal && mv_is_empty(&pbxval1->xval.terminal_mlrval)
		)
	);
}

static inline boxed_xval_t b_x_isnotempty_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!(pbxval1->xval.is_terminal && mv_is_empty(&pbxval1->xval.terminal_mlrval))
		)
	);
}

static inline boxed_xval_t b_x_isemptymap_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			!pbxval1->xval.is_terminal && pbxval1->xval.pnext_level->num_occupied == 0
		)
	);
}

static inline boxed_xval_t b_x_isnonemptymap_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_bool(
			pbxval1->xval.is_terminal || pbxval1->xval.pnext_level->num_occupied != 0
		)
	);
}

static inline boxed_xval_t b_x_typeof_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
	    mv_from_string(
			mlhmmv_xvalue_describe_type_simple(&pbxval1->xval), NO_FREE
		)
	);
}

// ----------------------------------------------------------------
boxed_xval_t b_x_haskey_xfunc(boxed_xval_t* pmapval, boxed_xval_t* pkeyval);

boxed_xval_t b_x_length_xfunc(boxed_xval_t* pbxval1);

// xxx to do (non-inline):
//boxed_xval_t i_m_depth_xfunc(boxed_xval_t* pbxval1);
//boxed_xval_t i_m_deepcount_xfunc(boxed_xval_t* pbxval1);
//boxed_xval_t m_mm_mapsum_xfunc(boxed_xval_t* pbxval1, boxed_xval_t* pbxval2);
//boxed_xval_t m_mm_mapdiff_xfunc(boxed_xval_t* pbxval1, boxed_xval_t* pbxval2);
//boxed_xval_t m_ss_split_xfunc(boxed_xval_t* pbxval1, boxed_xval_t* pbxval2);
//boxed_xval_t s_ms_join_xfunc(boxed_xval_t* pbxval1, boxed_xval_t* pbxval2);

#endif // XVFUNCS_H
