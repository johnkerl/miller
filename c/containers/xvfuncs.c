#include "../containers/xvfuncs.h"

// ----------------------------------------------------------------
boxed_xval_t b_xx_haskey_xfunc(boxed_xval_t* pmapval, boxed_xval_t* pkeyval) {
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
static int depth_aux(mlhmmv_xvalue_t* pxval) {
	if (pxval->is_terminal) {
		return 0;
	} else {
		int max = 0;
		for (mlhmmv_level_entry_t* pe = pxval->pnext_level->phead; pe != NULL; pe = pe->pnext) {
			int curr = depth_aux(&pe->level_xvalue);
			max = (curr > max) ? curr : max;
		}
		return 1 + max;
	}
}

boxed_xval_t i_x_depth_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_int(
			depth_aux(&pbxval1->xval)
		)
	);
}

// ----------------------------------------------------------------
static int leafcount_aux(mlhmmv_xvalue_t* pxval) {
	if (pxval->is_terminal) {
		return 1;
	} else {
		int sum = 0;
		for (mlhmmv_level_entry_t* pe = pxval->pnext_level->phead; pe != NULL; pe = pe->pnext) {
			sum += leafcount_aux(&pe->level_xvalue);
		}
		return sum;
	}
}

// xxx memmgt
boxed_xval_t i_x_leafcount_xfunc(boxed_xval_t* pbxval1) {
	return box_ephemeral_val(
		mv_from_int(
			leafcount_aux(&pbxval1->xval)
		)
	);
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
	mlhmmv_xvalue_t diff = mlhmmv_xvalue_alloc_empty_map();
	if (nxvals == 0) {
		return box_ephemeral_xval(diff);
	}
	// xxx to-do optmization: transfer arg 1 if it's ephemeral

	// xxx methodize
	int i = 0;
	if (!pbxvals[i].xval.is_terminal) {
		mlhmmv_level_t* plevel = pbxvals[i].xval.pnext_level;
		for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
			// xxx do refs/copies correctly
			mlhmmv_xvalue_t xval_copy = mlhmmv_xvalue_copy(&pe->level_xvalue);
			sllmve_t e = (sllmve_t) { .value = pe->level_key, .free_flags = 0, .pnext = NULL };
			mlhmmv_level_put_xvalue(diff.pnext_level, &e, &xval_copy);
		}
	}

	for (i = 1; i < nxvals; i++) {
		if (!pbxvals[i].xval.is_terminal) {
			mlhmmv_level_t* plevel = pbxvals[i].xval.pnext_level;
			for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
				sllmve_t e = (sllmve_t) { .value = pe->level_key, .free_flags = 0, .pnext = NULL };
				mlhmmv_level_remove(diff.pnext_level, &e);
			}
		}
	}

	return box_ephemeral_xval(diff);
}

// ----------------------------------------------------------------
boxed_xval_t m_ss_splitnv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}

// ----------------------------------------------------------------
boxed_xval_t m_ss_splitkv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}

// ----------------------------------------------------------------
boxed_xval_t s_ms_joink_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}

// ----------------------------------------------------------------
boxed_xval_t s_ms_joinv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}

// ----------------------------------------------------------------
boxed_xval_t s_mss_joinkv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	return box_ephemeral_val(mv_absent()); // xxx stub
}
