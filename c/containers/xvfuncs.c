#include "../containers/xvfuncs.h"
#include "../lib/string_builder.h"
#include "../lib/free_flags.h"

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
// Precondition: the caller has ensured that both arguments are string-valued terminals.
boxed_xval_t m_ss_splitnv_xfunc(boxed_xval_t* pstringval, boxed_xval_t* psepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* sep = psepval->xval.terminal_mlrval.u.strv;

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = strsep(&walker, sep)) != NULL) {
		mv_t key = mv_from_int(i);
		mv_t val = mv_type_infer_string_or_float_or_int(piece, NO_FREE);
		mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		i++;
	}
	free(input);
	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that both arguments are string-valued terminals.
boxed_xval_t m_ss_splitnvx_xfunc(boxed_xval_t* pstringval, boxed_xval_t* psepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* sep = psepval->xval.terminal_mlrval.u.strv;

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = strsep(&walker, sep)) != NULL) {
		mv_t key = mv_from_int(i);
		mv_t val = mv_type_infer_string(piece, NO_FREE);
		mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		i++;
	}

	free(input);
	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that all arguments are string-valued terminals.
boxed_xval_t m_sss_splitkv_xfunc(boxed_xval_t* pstringval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* listsep = plistsepval->xval.terminal_mlrval.u.strv;
	char* pairsep = ppairsepval->xval.terminal_mlrval.u.strv;

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = strsep(&walker, listsep)) != NULL) {
		char* xxx_rename = piece;
		char* left = strsep(&xxx_rename, pairsep);
		if (xxx_rename == NULL) {
			mv_t key = mv_from_int(i);
			mv_t val = mv_type_infer_string_or_float_or_int(left, NO_FREE);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		} else {
			char* right = strsep(&xxx_rename, pairsep);
			mv_t key = mv_from_string(left, NO_FREE);
			mv_t val = mv_type_infer_string_or_float_or_int(right, NO_FREE);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		}
		i++;
	}

	free(input);
	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that all arguments are string-valued terminals.
boxed_xval_t m_sss_splitkvx_xfunc(boxed_xval_t* pstringval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* listsep = plistsepval->xval.terminal_mlrval.u.strv;
	char* pairsep = ppairsepval->xval.terminal_mlrval.u.strv;

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = strsep(&walker, listsep)) != NULL) {
		char* xxx_rename = piece;
		char* left = strsep(&xxx_rename, pairsep);
		if (xxx_rename == NULL) {
			mv_t key = mv_from_int(i);
			mv_t val = mv_type_infer_string(left, NO_FREE);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		} else {
			char* right = strsep(&xxx_rename, pairsep);
			mv_t key = mv_from_string(left, NO_FREE);
			mv_t val = mv_type_infer_string(right, NO_FREE);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		}
		i++;
	}

	free(input);
	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
#define SB_JOIN_ALLOC_SIZE 128

// ----------------------------------------------------------------
// Precondition: the caller has ensured that the separator is a string-valued terminal.
boxed_xval_t s_ms_joink_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	if (pmapval->xval.is_terminal) {
		return box_ephemeral_val(mv_absent());
	}

	string_builder_t* psb = sb_alloc(SB_JOIN_ALLOC_SIZE);
	for (mlhmmv_level_entry_t* pentry = pmapval->xval.pnext_level->phead; pentry != NULL; pentry = pentry->pnext) {
		// The string_builder object will copy the string so we can point into source string space,
		// without a copy here, when possible.
		char free_flags = 0;
		char* sval = mv_maybe_alloc_format_val(&pentry->level_key, &free_flags);
		sb_append_string(psb, sval);
		if (free_flags)
			free(sval);
		if (pentry->pnext != NULL) {
			sb_append_string(psb, psepval->xval.terminal_mlrval.u.strv);
		}
	}

	char* sretval = sb_finish(psb);
	sb_free(psb);
	return box_ephemeral_val(
		mv_from_string(sretval, FREE_ENTRY_VALUE)
	);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that the separator is a string-valued terminal.
boxed_xval_t s_ms_joinv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* psepval) {
	if (pmapval->xval.is_terminal) {
		return box_ephemeral_val(mv_absent());
	}

	string_builder_t* psb = sb_alloc(SB_JOIN_ALLOC_SIZE);
	for (mlhmmv_level_entry_t* pentry = pmapval->xval.pnext_level->phead; pentry != NULL; pentry = pentry->pnext) {
		if (pentry->level_xvalue.is_terminal) {
			// The string_builder object will copy the string so we can point into source string space,
			// without a copy here, when possible.
			char free_flags = 0;
			char* sval = mv_maybe_alloc_format_val(&pentry->level_xvalue.terminal_mlrval, &free_flags);
			sb_append_string(psb, sval);
			if (free_flags)
				free(sval);
			if (pentry->pnext != NULL) {
				sb_append_string(psb, psepval->xval.terminal_mlrval.u.strv);
			}
		}
	}

	char* sretval = sb_finish(psb);
	sb_free(psb);
	return box_ephemeral_val(
		mv_from_string(sretval, FREE_ENTRY_VALUE)
	);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that the separators are string-valued terminals.
boxed_xval_t s_mss_joinkv_xfunc(boxed_xval_t* pmapval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	if (pmapval->xval.is_terminal) {
		return box_ephemeral_val(mv_absent());
	}

	string_builder_t* psb = sb_alloc(SB_JOIN_ALLOC_SIZE);
	for (mlhmmv_level_entry_t* pentry = pmapval->xval.pnext_level->phead; pentry != NULL; pentry = pentry->pnext) {
		if (pentry->level_xvalue.is_terminal) {
			// The string_builder object will copy the string so we can point into source string space,
			// without a copy here, when possible.

			char kfree_flags = 0;
			char* skval = mv_maybe_alloc_format_val(&pentry->level_key, &kfree_flags);
			sb_append_string(psb, skval);
			if (kfree_flags)
				free(skval);

			sb_append_string(psb, ppairsepval->xval.terminal_mlrval.u.strv);

			char vfree_flags = 0;
			char* svval = mv_maybe_alloc_format_val(&pentry->level_xvalue.terminal_mlrval, &vfree_flags);
			sb_append_string(psb, svval);
			if (vfree_flags)
				free(svval);

			if (pentry->pnext != NULL) {
				sb_append_string(psb, plistsepval->xval.terminal_mlrval.u.strv);
			}
		}
	}

	char* sretval = sb_finish(psb);
	sb_free(psb);
	return box_ephemeral_val(
		mv_from_string(sretval, FREE_ENTRY_VALUE)
	);
}
