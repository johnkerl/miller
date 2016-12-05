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
boxed_xval_t i_x_length_xfunc(boxed_xval_t* pxval1) {
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

// ----------------------------------------------------------------
boxed_xval_t variadic_mapsum_xfunc(boxed_xval_t* pbxvals, int nxvals) {
	// xxx to-do optmization: transfer arg 1 if it's ephemeral
	mlhmmv_xvalue_t sum = mlhmmv_xvalue_alloc_empty_map();
	for (int i = 0; i < nxvals; i++) {
		if (pbxvals[i].xval.is_terminal)
			continue;
		mlhmmv_level_t* plevel = pbxvals[i].xval.pnext_level;
		for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
			// xxx do refs/copies correctly
			mlhmmv_xvalue_t xval_copy = mlhmmv_xvalue_copy(&pe->level_xvalue);
			sllmve_t e = (sllmve_t) { .value = pe->level_key, .free_flags = 0, .pnext = NULL };
			mlhmmv_level_put_xvalue(sum.pnext_level, &e, &xval_copy);
		}
	}
	return box_ephemeral_xval(sum);
}

// ----------------------------------------------------------------
boxed_xval_t variadic_mapdiff_xfunc(boxed_xval_t* pbxvals, int nxvals) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}
