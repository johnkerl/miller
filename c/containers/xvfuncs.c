#include "../containers/xvfuncs.h"

// ----------------------------------------------------------------
mlhmmv_xvalue_t b_x_haskey_xfunc(mlhmmv_xvalue_t* pmapval, mlhmmv_xvalue_t* pkeyval) {
	if (pmapval->is_terminal) {
		return mlhmmv_xvalue_wrap_terminal(mv_from_bool(FALSE));
	} else if (!pkeyval->is_terminal) {
		return mlhmmv_xvalue_wrap_terminal(mv_from_bool(FALSE));
	} else {
		return mlhmmv_xvalue_wrap_terminal(
			mv_from_bool(
				mlhmmv_level_has_key(pmapval->pnext_level, &pkeyval->terminal_mlrval)
			)
		);
	}
}

// ----------------------------------------------------------------
mlhmmv_xvalue_t b_x_length_xfunc(mlhmmv_xvalue_t* pxval1) {
	if (pxval1->is_terminal) {
		return mlhmmv_xvalue_wrap_terminal(
			mv_from_int(1)
		);
	} else {
		return mlhmmv_xvalue_wrap_terminal(
			mv_from_int(
				pxval1->pnext_level->num_occupied
			)
		);
	}
}
