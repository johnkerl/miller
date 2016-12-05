#include "../containers/xvfuncs.h"

// ----------------------------------------------------------------
boxed_xval_t b_x_haskey_xfunc(boxed_xval_t* pmapval, boxed_xval_t* pkeyval) {
	if (pmapval->xval.is_terminal) {
		return box_ephemeral_val(mv_from_bool(FALSE));
	} else if (!pkeyval->xval.is_terminal) {
		return box_ephemeral_val(mv_from_bool(FALSE));
	} else {
		return box_ephemeral_val(
			mv_from_bool(
				mlhmmv_level_has_key(pmapval->xval.pnext_level, &pkeyval->xval.terminal_mlrval)
			)
		);
	}
}

// ----------------------------------------------------------------
boxed_xval_t b_x_length_xfunc(boxed_xval_t* pxval1) {
	if (pxval1->xval.is_terminal) {
		return box_ephemeral_val(
			mv_from_int(1)
		);
	} else {
		return box_ephemeral_val(
			mv_from_int(
				pxval1->xval.pnext_level->num_occupied
			)
		);
	}
}
