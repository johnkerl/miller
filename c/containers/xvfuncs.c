#include "../lib/mlr_arch.h"
#include "../lib/string_builder.h"
#include "../lib/free_flags.h"
#include "../containers/xvfuncs.h"

// ----------------------------------------------------------------
boxed_xval_t b_xx_haskey_xfunc(boxed_xval_t* pmapval, boxed_xval_t* pkeyval) {
	boxed_xval_t rv;
	if (pmapval->xval.is_terminal) {
		rv = box_ephemeral_val(mv_from_bool(FALSE));
	} else if (!pkeyval->xval.is_terminal) {
		rv = box_ephemeral_val(mv_from_bool(FALSE));
	} else {
		rv = box_ephemeral_val(
			mv_from_bool(
				mlhmmv_level_has_key(pmapval->xval.pnext_level, &pkeyval->xval.terminal_mlrval)
			)
		);
	}
	if (pmapval->is_ephemeral)
		mlhmmv_xvalue_free(&pmapval->xval);
	if (pkeyval->is_ephemeral)
		mlhmmv_xvalue_free(&pkeyval->xval);
	return rv;
}

// ----------------------------------------------------------------
boxed_xval_t i_x_length_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv;
	if (pbxval1->xval.is_terminal) {
		rv = box_ephemeral_val(
			mv_from_int(1)
		);
	} else {
		rv = box_ephemeral_val(
			mv_from_int(
				pbxval1->xval.pnext_level->num_occupied
			)
		);
	}
	if (pbxval1->is_ephemeral)
		mlhmmv_xvalue_free(&pbxval1->xval);
	return rv;
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
	boxed_xval_t rv = box_ephemeral_val(
		mv_from_int(
			depth_aux(&pbxval1->xval)
		)
	);

	if (pbxval1->is_ephemeral)
		mlhmmv_xvalue_free(&pbxval1->xval);

	return rv;
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

boxed_xval_t i_x_leafcount_xfunc(boxed_xval_t* pbxval1) {
	boxed_xval_t rv = box_ephemeral_val(
		mv_from_int(
			leafcount_aux(&pbxval1->xval)
		)
	);

	if (pbxval1->is_ephemeral)
		mlhmmv_xvalue_free(&pbxval1->xval);

	return rv;
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
	boxed_xval_t rv = box_ephemeral_xval(sum);

	for (int i = 0; i < nxvals; i++) {
		if (pbxvals[i].is_ephemeral)
			mlhmmv_xvalue_free(&pbxvals[i].xval);
	}

	return rv;
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
		if (pbxvals[i].is_ephemeral)
			mlhmmv_xvalue_free(&pbxvals[i].xval);
	}

	for (i = 1; i < nxvals; i++) {
		if (!pbxvals[i].xval.is_terminal) {
			mlhmmv_level_t* plevel = pbxvals[i].xval.pnext_level;
			for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
				sllmve_t e = (sllmve_t) { .value = pe->level_key, .free_flags = 0, .pnext = NULL };
				mlhmmv_level_remove(diff.pnext_level, &e);
			}
		}
		if (pbxvals[i].is_ephemeral)
			mlhmmv_xvalue_free(&pbxvals[i].xval);
	}

	return box_ephemeral_xval(diff);
}

// ----------------------------------------------------------------
// Precondition (validated before we're called): there is at least one argument
// which is the map to be unkeyed.
boxed_xval_t variadic_mapexcept_xfunc(
	boxed_xval_t* pbxvals,
	int nxvals)
{
	MLR_INTERNAL_CODING_ERROR_IF(nxvals < 1);

	boxed_xval_t* pinbxval = &pbxvals[0];
	boxed_xval_t outbxval = box_ephemeral_xval(mlhmmv_xvalue_copy(&pinbxval->xval));

	if (pinbxval->xval.is_terminal) { // non-map
		return outbxval;
	}

	mlhmmv_level_t* poutlevel = outbxval.xval.pnext_level;

	for (int i = 1; i < nxvals; i++) {
		if (pbxvals[i].xval.is_terminal) {
			mv_t* pkey = &pbxvals[i].xval.terminal_mlrval;
			sllmve_t e = (sllmve_t) { .value = *pkey, .free_flags = 0, .pnext = NULL };
			mlhmmv_level_remove(poutlevel, &e);
		}
		if (pbxvals[i].is_ephemeral)
			mlhmmv_xvalue_free(&pbxvals[i].xval);
	}

	return outbxval;
}

// ----------------------------------------------------------------
// Precondition (validated before we're called): there is at least one argument
// which is the map to be unkeyed.
boxed_xval_t variadic_maponly_xfunc(
	boxed_xval_t* pbxvals,
	int nxvals)
{
	MLR_INTERNAL_CODING_ERROR_IF(nxvals < 1);

	boxed_xval_t* pinbxval = &pbxvals[0];
	boxed_xval_t outbxval = box_ephemeral_xval(mlhmmv_xvalue_alloc_empty_map());

	if (pinbxval->xval.is_terminal) { // non-map
		return outbxval;
	}

	mlhmmv_level_t* pinlevel = pinbxval->xval.pnext_level;
	mlhmmv_level_t* poutlevel = outbxval.xval.pnext_level;

	for (int i = 1; i < nxvals; i++) {

		if (pbxvals[i].xval.is_terminal) {
			mv_t* pkey = &pbxvals[i].xval.terminal_mlrval;
			sllmv_t* pkeylist = sllmv_single_no_free(pkey);
			int unused = 0;
			mlhmmv_xvalue_t* pval = mlhmmv_level_look_up_and_ref_xvalue(pinlevel, pkeylist, &unused);
			if (pval != NULL) {
				mlhmmv_xvalue_t copyval = mlhmmv_xvalue_copy(pval);
				// mlhmmv_level_put_xvalue copies key not value
				mlhmmv_level_put_xvalue(poutlevel, pkeylist->phead, &copyval);
			}
			sllmv_free(pkeylist);
		}
		if (pbxvals[i].is_ephemeral)
			mlhmmv_xvalue_free(&pbxvals[i].xval);

	}

	return outbxval;
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that both arguments are string-valued terminals.
boxed_xval_t m_ss_splitnv_xfunc(boxed_xval_t* pstringval, boxed_xval_t* psepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* sep = psepval->xval.terminal_mlrval.u.strv;
	int seplen = strlen(sep);

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = mlr_strmsep(&walker, sep, seplen)) != NULL) {
		mv_t key = mv_from_int(i);
		mv_t val = mv_ref_type_infer_string_or_float_or_int(piece);
		mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		i++;
	}
	free(input);

	if (pstringval->is_ephemeral)
		mlhmmv_xvalue_free(&pstringval->xval);
	if (psepval->is_ephemeral)
		mlhmmv_xvalue_free(&psepval->xval);

	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that both arguments are string-valued terminals.
boxed_xval_t m_ss_splitnvx_xfunc(boxed_xval_t* pstringval, boxed_xval_t* psepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* sep = psepval->xval.terminal_mlrval.u.strv;
	int seplen = strlen(sep);

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = mlr_strmsep(&walker, sep, seplen)) != NULL) {
		mv_t key = mv_from_int(i);
		// xxx do not ref here
		mv_t val = mv_ref_type_infer_string(piece);
		// xxx make clear the copy/ref semantics for mlhmmv put with scalar value.
		// at the moment this does an mv_copy.
		mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		i++;
	}

	free(input);

	if (pstringval->is_ephemeral)
		mlhmmv_xvalue_free(&pstringval->xval);
	if (psepval->is_ephemeral)
		mlhmmv_xvalue_free(&psepval->xval);

	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that all arguments are string-valued terminals.
boxed_xval_t m_sss_splitkv_xfunc(boxed_xval_t* pstringval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* listsep = plistsepval->xval.terminal_mlrval.u.strv;
	char* pairsep = ppairsepval->xval.terminal_mlrval.u.strv;
	int listseplen = strlen(listsep);
	int pairseplen = strlen(pairsep);

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = mlr_strmsep(&walker, listsep, listseplen)) != NULL) {
		char* pair = piece;
		char* left = mlr_strmsep(&pair, pairsep, pairseplen);
		if (pair == NULL) {
			mv_t key = mv_from_int(i);
			mv_t val = mv_ref_type_infer_string_or_float_or_int(left);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		} else {
			char* right = mlr_strmsep(&pair, pairsep, pairseplen);
			mv_t key = mv_from_string(left, NO_FREE);
			mv_t val = mv_ref_type_infer_string_or_float_or_int(right);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		}
		i++;
	}

	free(input);

	if (pstringval->is_ephemeral)
		mlhmmv_xvalue_free(&pstringval->xval);
	if (ppairsepval->is_ephemeral)
		mlhmmv_xvalue_free(&ppairsepval->xval);
	if (plistsepval->is_ephemeral)
		mlhmmv_xvalue_free(&plistsepval->xval);

	return box_ephemeral_xval(map);
}

// ----------------------------------------------------------------
// Precondition: the caller has ensured that all arguments are string-valued terminals.
boxed_xval_t m_sss_splitkvx_xfunc(boxed_xval_t* pstringval, boxed_xval_t* ppairsepval, boxed_xval_t* plistsepval) {
	mlhmmv_xvalue_t map = mlhmmv_xvalue_alloc_empty_map();
	char* input = mlr_strdup_or_die(pstringval->xval.terminal_mlrval.u.strv);
	char* listsep = plistsepval->xval.terminal_mlrval.u.strv;
	char* pairsep = ppairsepval->xval.terminal_mlrval.u.strv;
	int listseplen = strlen(listsep);
	int pairseplen = strlen(pairsep);

	int i = 1;
	char* walker = input;
	char* piece;
	while ((piece = mlr_strmsep(&walker, listsep, listseplen)) != NULL) {
		char* pair = piece;
		char* left = mlr_strmsep(&pair, pairsep, pairseplen);
		if (pair == NULL) {
			mv_t key = mv_from_int(i);
			mv_t val = mv_ref_type_infer_string(left);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		} else {
			char* right = mlr_strmsep(&pair, pairsep, pairseplen);
			mv_t key = mv_from_string(left, NO_FREE);
			mv_t val = mv_ref_type_infer_string(right);
			mlhmmv_level_put_terminal_singly_keyed(map.pnext_level, &key, &val);
		}
		i++;
	}

	free(input);

	if (pstringval->is_ephemeral)
		mlhmmv_xvalue_free(&pstringval->xval);
	if (ppairsepval->is_ephemeral)
		mlhmmv_xvalue_free(&ppairsepval->xval);
	if (plistsepval->is_ephemeral)
		mlhmmv_xvalue_free(&plistsepval->xval);

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

	if (pmapval->is_ephemeral)
		mlhmmv_xvalue_free(&pmapval->xval);
	if (psepval->is_ephemeral)
		mlhmmv_xvalue_free(&psepval->xval);

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

	if (pmapval->is_ephemeral)
		mlhmmv_xvalue_free(&pmapval->xval);
	if (psepval->is_ephemeral)
		mlhmmv_xvalue_free(&psepval->xval);

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

	if (pmapval->is_ephemeral)
		mlhmmv_xvalue_free(&pmapval->xval);
	if (ppairsepval->is_ephemeral)
		mlhmmv_xvalue_free(&ppairsepval->xval);
	if (plistsepval->is_ephemeral)
		mlhmmv_xvalue_free(&plistsepval->xval);

	return box_ephemeral_val(
		mv_from_string(sretval, FREE_ENTRY_VALUE)
	);
}
