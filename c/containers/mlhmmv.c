// ================================================================
// Array-only (open addressing) multi-level hash map, with linear probing for collisions.
// All keys, and terminal-level values, are mlrvals.
//
// Notes:
// * null key is not supported.
// * null value is not supported.
//
// See also:
// * http://en.wikipedia.org/wiki/Hash_table
// * http://docs.oracle.com/javase/6/docs/api/java/util/Map.html
// ================================================================

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/mlhmmv.h"
#include "lib/mvfuncs.h"

// ================================================================
// Allow compile-time override, e.g using gcc -D.

#ifndef LOAD_FACTOR
#define LOAD_FACTOR          0.7
#endif

#ifndef ENLARGEMENT_FACTOR
#define ENLARGEMENT_FACTOR   2
#endif

#define OCCUPIED             0xa4
#define DELETED              0xb8
#define EMPTY                0xce

// ================================================================
static int  mlhmmv_hash_func(mv_t* plevel_key);
static void json_decimal_print       (FILE* ostream, char* s, int json_apply_ofmt_to_floats, double parsed);
static void json_print_string_escaped(FILE* ostream, char* s);

// ----------------------------------------------------------------
static void mlhmmv_level_init(mlhmmv_level_t  *plevel, int length);

// ----------------------------------------------------------------
static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pideal_index);

static mlhmmv_level_entry_t* mlhmmv_level_look_up_and_ref_entry(
	mlhmmv_level_t* plevel, sllmve_t* prestkeys, int* perror);

static mlhmmv_level_entry_t* mlhmmv_level_get_next_level_entry(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pindex);
static mlhmmv_xvalue_t*      mlhmmv_level_get_next_level_xvalue(mlhmmv_level_t* plevel, mv_t* plevel_key);

static mlhmmv_level_t* mlhmmv_level_ref_or_create(mlhmmv_level_t* plevel, sllmve_t* prest_keys);
static mlhmmv_level_t* mlhmmv_level_get_or_create_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys);

static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel);
static void mlhmmv_level_move(mlhmmv_level_t* plevel, mv_t* plevel_key, mlhmmv_xvalue_t* plevel_value);

static void mlhmmv_level_put_xvalue_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mlhmmv_xvalue_t* pvalue);
static void mlhmmv_level_put_terminal_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value);

static void mlhmmv_level_to_lrecs_across_records(
	mlhmmv_level_t* plevel,
	char*           prefix,
	sllmve_t*       prestnames,
	lrec_t*         ptemplate,
	sllv_t*         poutrecs,
	int             do_full_prefixing,
	char*           flatten_separator);

static void mlhmmv_level_to_lrec_within_record(
	mlhmmv_level_t* plevel,
	char*           prefix,
	lrec_t*         poutrec,
	int             do_full_prefixing,
	char*           flatten_separator);

static void mlhhmv_levels_to_lrecs_lashed_across_records(
	mlhmmv_level_t** plevels,
	char**           prefixes,
	int              num_levels,
	sllmve_t*        prestnames,
	lrec_t*          ptemplate,
	sllv_t*          poutrecs,
	int              do_full_prefixing,
	char*            flatten_separator);

static void mlhhmv_levels_to_lrecs_lashed_within_records(
	mlhmmv_level_t** pplevels,
	char**           prefixes,
	int              num_levels,
	lrec_t*          poutrec,
	int              do_full_prefixing,
	char*            flatten_separator);

static void mlhmmv_level_print_single_line(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_keys_always, int quote_values_always, int json_apply_ofmt_to_floats,
	FILE* ostream);

// ----------------------------------------------------------------
static void mlhmmv_root_put_xvalue(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, mlhmmv_xvalue_t* pvalue);

// ================================================================
typedef int mlhmmv_typed_hash_func(mv_t* pa);
static int mlhmmv_string_hash_func(mv_t* pa) {
	return mlr_string_hash_func(pa->u.strv);
}
static int mlhmmv_int_hash_func(mv_t* pa) {
	return pa->u.intv;
}
static int mlhmmv_other_hash_func(mv_t* pa) {
	fprintf(stderr, "%s: map keys must be of type %s or %s; got %s.\n",
		MLR_GLOBALS.bargv0,
		mt_describe_type(MT_STRING),
		mt_describe_type(MT_INT),
		mt_describe_type(pa->type));
	exit(1);
}
static mlhmmv_typed_hash_func* mlhmmv_hash_func_dispositions[MT_DIM] = {
	/*ERROR*/  mlhmmv_other_hash_func,
	/*ABSENT*/ mlhmmv_other_hash_func,
	/*EMPTY*/  mlhmmv_other_hash_func,
	/*STRING*/ mlhmmv_string_hash_func,
	/*INT*/    mlhmmv_int_hash_func,
	/*FLOAT*/  mlhmmv_other_hash_func,
	/*BOOL*/   mlhmmv_other_hash_func,
};
static int mlhmmv_hash_func(mv_t* pa) {
	return mlhmmv_hash_func_dispositions[pa->type](pa);
}

// ================================================================
void mlhmmv_print_terminal(mv_t* pmv, int quote_keys_always, int quote_values_always, int json_apply_ofmt_to_floats,
	FILE* ostream)
{
	char* level_value_string = mv_alloc_format_val(pmv);
	if (quote_values_always) {
		json_print_string_escaped(ostream, level_value_string);
	} else if (pmv->type == MT_STRING) {
		double parsed;
		if (mlr_try_float_from_string(level_value_string, &parsed)) {
			json_decimal_print(ostream, level_value_string, json_apply_ofmt_to_floats, parsed);
		} else if (streq(level_value_string, "true") || streq(level_value_string, "false")) {
			fprintf(ostream, "%s", level_value_string);
		} else {
			json_print_string_escaped(ostream, level_value_string);
		}
	} else {
		fprintf(ostream, "%s", level_value_string);
	}
	free(level_value_string);
}

// 0.123 is valid JSON; .123 is not. Meanwhile Miller is a format-converter tool
// so if there is perfectly legitimate CSV/DKVP/etc. data to be JSON-formatted,
// we make it JSON-compliant.
//
// Precondition: the caller has already checked that the string represents a number.
static void json_decimal_print(FILE* ostream, char* s, int json_apply_ofmt_to_floats, double parsed) {
	if (json_apply_ofmt_to_floats) {
		fprintf(ostream, MLR_GLOBALS.ofmt, parsed);
	} else {
		if (s[0] == '.') {
			fprintf(ostream, "0%s", s);
		} else if (s[0] == '-' && s[1] == '.') {
			fprintf(ostream, "-0.%s", &s[2]);
		} else {
			fprintf(ostream, "%s", s);
		}
	}
}

static void json_print_string_escaped(FILE* ostream, char* s) {
	fputc('"', ostream);
	for (char* p = s; *p; p++) {
		char c = *p;
		switch (c) {
		case '"':
			fputc('\\', ostream);
			fputc('"', ostream);
			break;
		case '\\':
			fputc('\\', ostream);
			fputc('\\', ostream);
			break;
		case '\n':
			fputc('\\', ostream);
			fputc('n', ostream);
			break;
		case '\r':
			fputc('\\', ostream);
			fputc('r', ostream);
			break;
		case '\t':
			fputc('\\', ostream);
			fputc('t', ostream);
			break;
		case '\b':
			fputc('\\', ostream);
			fputc('b', ostream);
			break;
		case '\f':
			fputc('\\', ostream);
			fputc('f', ostream);
			break;
		default:
			fputc(c, ostream);
			break;
		}
	}
	fputc('"', ostream);
}

// ================================================================
mlhmmv_xvalue_t mlhmmv_xvalue_alloc_empty_map() {
	mlhmmv_xvalue_t xval = (mlhmmv_xvalue_t) {
		.is_terminal = FALSE,
		.terminal_mlrval= mv_absent(),
		.pnext_level = mlhmmv_level_alloc()
	};
	return xval;
}

// ----------------------------------------------------------------
void mlhmmv_xvalue_reset(mlhmmv_xvalue_t* pxvalue) {
	pxvalue->is_terminal     = TRUE;
	pxvalue->terminal_mlrval = mv_absent();
	pxvalue->pnext_level     = NULL;
}

// ----------------------------------------------------------------
mlhmmv_xvalue_t mlhmmv_xvalue_copy(mlhmmv_xvalue_t* pvalue) {
	if (pvalue->is_terminal) {
		return (mlhmmv_xvalue_t) {
			.is_terminal = TRUE,
			.terminal_mlrval = mv_copy(&pvalue->terminal_mlrval),
			.pnext_level = NULL,
		};

	} else {
		mlhmmv_level_t* psrc_level = pvalue->pnext_level;
		mlhmmv_level_t* pdst_level = mlr_malloc_or_die(sizeof(mlhmmv_level_t));

		mlhmmv_level_init(pdst_level, MLHMMV_INITIAL_ARRAY_LENGTH);

		for (
			mlhmmv_level_entry_t* psubentry = psrc_level->phead;
			psubentry != NULL;
			psubentry = psubentry->pnext)
		{
			mlhmmv_xvalue_t next_value = mlhmmv_xvalue_copy(&psubentry->level_xvalue);
			mlhmmv_level_put_xvalue_singly_keyed(pdst_level, &psubentry->level_key, &next_value);
		}

		return (mlhmmv_xvalue_t) {
			.is_terminal = FALSE,
			.terminal_mlrval = mv_absent(),
			.pnext_level = pdst_level,
		};
	}
}

// ----------------------------------------------------------------
void mlhmmv_xvalue_free(mlhmmv_xvalue_t* pxvalue) {
	if (pxvalue->is_terminal) {
		mv_free(&pxvalue->terminal_mlrval);
	} else if (pxvalue->pnext_level != NULL) {
		mlhmmv_level_free(pxvalue->pnext_level);
		pxvalue->pnext_level = NULL;
	}
}

// ----------------------------------------------------------------
char* mlhmmv_xvalue_describe_type_simple(mlhmmv_xvalue_t* pxvalue) {
	if (pxvalue->is_terminal) {
		return mt_describe_type_simple(pxvalue->terminal_mlrval.type);
	} else {
		return "map";
	}
}

sllv_t* mlhmmv_xvalue_copy_keys_indexed(mlhmmv_xvalue_t* pmvalue, sllmv_t* pmvkeys) { // xxx code dedupe
	int error;
	if (pmvkeys == NULL || pmvkeys->length == 0) {
		return mlhmmv_xvalue_copy_keys_nonindexed(pmvalue);
	} else if (pmvalue->is_terminal) { // xxx copy this check up to oosvar case too
		return sllv_alloc();
	} else {
		mlhmmv_level_entry_t* pfromentry = mlhmmv_level_look_up_and_ref_entry(pmvalue->pnext_level,
			pmvkeys->phead, &error);
		if (pfromentry != NULL) {
			return mlhmmv_xvalue_copy_keys_nonindexed(&pfromentry->level_xvalue);
		} else {
			return sllv_alloc();
		}
	}
}

sllv_t* mlhmmv_xvalue_copy_keys_nonindexed(mlhmmv_xvalue_t* pvalue) {
	sllv_t* pkeys = sllv_alloc();

	if (!pvalue->is_terminal) {
		mlhmmv_level_t* pnext_level = pvalue->pnext_level;
		for (mlhmmv_level_entry_t* pe = pnext_level->phead; pe != NULL; pe = pe->pnext) {
			mv_t* p = mv_alloc_copy(&pe->level_key);
			sllv_append(pkeys, p);
		}
	}

	return pkeys;
}

// ----------------------------------------------------------------
void mlhmmv_xvalues_to_lrecs_lashed(mlhmmv_xvalue_t** ptop_values, int num_submaps, mv_t* pbasenames, sllmv_t* pnames,
	sllv_t* poutrecs, int do_full_prefixing, char* flatten_separator)
{

	// First is primary and rest are lashed to it (lookups with same keys as primary).
	if (ptop_values[0] == NULL) {
		// No such entry in the mlhmmv results in no output records
	} else if (ptop_values[0]->is_terminal && mv_is_present(&ptop_values[0]->terminal_mlrval)) {
		lrec_t* poutrec = lrec_unbacked_alloc();
		for (int i = 0; i < num_submaps; i++) {
			// E.g. '@v = 3' at the top level of the mlhmmv.
			if (ptop_values[i]->is_terminal && mv_is_present(&ptop_values[i]->terminal_mlrval)) {
				lrec_put(poutrec,
					mv_alloc_format_val(&pbasenames[i]),
					mv_alloc_format_val(&ptop_values[i]->terminal_mlrval), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			}
		}
		sllv_append(poutrecs, poutrec);
	} else {
		// E.g. '@v = {...}' at the top level of the map, where the map is an oosvar
		// submap keyed by oosvar-name 'v', or a localvar value, or a map literal.  This
		// needs to be flattened down to an lrec which is a list of key-value pairs.  We
		// recursively invoke mlhhmv_levels_to_lrecs_lashed_across_records for each of the
		// name-list entries, one map level deeper each call, then from there invoke
		// mlhhmv_levels_to_lrecs_lashed_within_records on any remaining map levels.
		lrec_t* ptemplate = lrec_unbacked_alloc();

		mlhmmv_level_t** ppnext_levels = mlr_malloc_or_die(num_submaps * sizeof(mlhmmv_level_t*));
		char** oosvar_names = mlr_malloc_or_die(num_submaps * sizeof(char*));
		for (int i = 0; i < num_submaps; i++) {
			if (ptop_values[i] == NULL || ptop_values[i]->is_terminal) {
				ppnext_levels[i] = NULL;
				oosvar_names[i] = NULL;
			} else {
				ppnext_levels[i] = ptop_values[i]->pnext_level;
				oosvar_names[i] = mv_alloc_format_val(&pbasenames[i]);
			}
		}

		mlhhmv_levels_to_lrecs_lashed_across_records(ppnext_levels, oosvar_names, num_submaps,
			pnames->phead, ptemplate, poutrecs, do_full_prefixing, flatten_separator);

		for (int i = 0; i < num_submaps; i++) {
			free(oosvar_names[i]);
		}
		free(oosvar_names);
		free(ppnext_levels);

		lrec_free(ptemplate);
	}
}

// ================================================================
mlhmmv_level_t* mlhmmv_level_alloc() {
	mlhmmv_level_t* plevel = mlr_malloc_or_die(sizeof(mlhmmv_level_t));
	mlhmmv_level_init(plevel, MLHMMV_INITIAL_ARRAY_LENGTH);
	return plevel;
}

// ----------------------------------------------------------------
static void mlhmmv_level_init(mlhmmv_level_t *plevel, int length) {
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->array_length = length;

	plevel->entries      = (mlhmmv_level_entry_t*)mlr_malloc_or_die(sizeof(mlhmmv_level_entry_t) * length);
	// Don't do mlhmmv_level_entry_clear() of all entries at init time, since this has a
	// drastic effect on the time needed to construct an empty map (and miller
	// constructs an awful lot of those). The attributes there are don't-cares
	// if the corresponding entry state is EMPTY. They are set on put, and
	// mutated on remove.

	plevel->states       = (mlhmmv_level_entry_state_t*)mlr_malloc_or_die(sizeof(mlhmmv_level_entry_state_t) * length);
	memset(plevel->states, EMPTY, length);

	plevel->phead        = NULL;
	plevel->ptail        = NULL;
}

// ----------------------------------------------------------------
void mlhmmv_level_free(mlhmmv_level_t* plevel) {
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		mv_free(&pentry->level_key);
		if (pentry->level_xvalue.is_terminal) {
			mv_free(&pentry->level_xvalue.terminal_mlrval);
		} else {
			mlhmmv_level_free(pentry->level_xvalue.pnext_level);
		}
	}
	free(plevel->entries);
	free(plevel->states);
	plevel->entries      = NULL;
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->array_length = 0;
	free(plevel);
}

// ----------------------------------------------------------------
void mlhmmv_level_clear(mlhmmv_level_t* plevel) {
	if (plevel->phead == NULL)
		return;

	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		if (pentry->level_xvalue.is_terminal) {
			mv_free(&pentry->level_xvalue.terminal_mlrval);
		} else {
			mlhmmv_level_free(pentry->level_xvalue.pnext_level);
		}
		mv_free(&pentry->level_key);
	}
	plevel->num_occupied = 0;
	plevel->num_freed    = 0;
	plevel->phead        = NULL;
	plevel->ptail        = NULL;

	memset(plevel->states, EMPTY, plevel->array_length);
}

// ----------------------------------------------------------------
int mlhmmv_level_has_key(mlhmmv_level_t* plevel, mv_t* plevel_key) {
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	if (plevel->states[index] == OCCUPIED)
		return TRUE;
	else if (plevel->states[index] == EMPTY)
		return FALSE;
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
// Returns >=0 for where the key is *or* should go (end of chain).

static int mlhmmv_level_find_index_for_key(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pideal_index) {
	int hash = mlhmmv_hash_func(plevel_key);
	int index = mlr_canonical_mod(hash, plevel->array_length);
	*pideal_index = index;
	int num_tries = 0;

	while (TRUE) {
		mlhmmv_level_entry_t* pentry = &plevel->entries[index];
		if (plevel->states[index] == OCCUPIED) {
			mv_t* ekey = &pentry->level_key;
			// Existing key found in chain.
			if (mv_equals_si(plevel_key, ekey))
				return index;
		} else if (plevel->states[index] == EMPTY) {
			return index;
		}

		// If the current entry has been freed, i.e. previously occupied,
		// the sought index may be further down the chain.  So we must
		// continue looking.
		if (++num_tries >= plevel->array_length) {
			fprintf(stderr,
				"%s: Coding error:  table full even after enlargement.\n", MLR_GLOBALS.bargv0);
			exit(1);
		}

		// Linear probing.
		if (++index >= plevel->array_length)
			index = 0;
	}
	MLR_INTERNAL_CODING_ERROR();
	return -1; // not reached
}

// ----------------------------------------------------------------
static mlhmmv_level_entry_t* mlhmmv_level_look_up_and_ref_entry(mlhmmv_level_t* plevel,
	sllmve_t* prestkeys, int* perror)
{
	if (perror)
		*perror = MLHMMV_ERROR_NONE;
	if (prestkeys == NULL) {
		if (perror)
			*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	mlhmmv_level_entry_t* plevel_entry = mlhmmv_level_get_next_level_entry(plevel, &prestkeys->value, NULL);
	while (prestkeys->pnext != NULL) {
		if (plevel_entry == NULL) {
			return NULL;
		}
		if (plevel_entry->level_xvalue.is_terminal) {
			if (perror)
				*perror = MLHMMV_ERROR_KEYLIST_TOO_DEEP;
			return NULL;
		}
		plevel = plevel_entry->level_xvalue.pnext_level;
		prestkeys = prestkeys->pnext;
		plevel_entry = mlhmmv_level_get_next_level_entry(plevel_entry->level_xvalue.pnext_level,
			&prestkeys->value, NULL);
	}
	return plevel_entry;
}

// ----------------------------------------------------------------
static mlhmmv_level_entry_t* mlhmmv_level_get_next_level_entry(mlhmmv_level_t* plevel, mv_t* plevel_key, int* pindex) {
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (pindex != NULL)
		*pindex = index;

	if (plevel->states[index] == OCCUPIED)
		return pentry;
	else if (plevel->states[index] == EMPTY)
		return NULL;
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
static mlhmmv_xvalue_t* mlhmmv_level_get_next_level_xvalue(mlhmmv_level_t* plevel, mv_t* plevel_key) {
	if (plevel == NULL)
		return NULL;
	mlhmmv_level_entry_t* pentry = mlhmmv_level_get_next_level_entry(plevel, plevel_key, NULL);
	if (pentry == NULL)
		return NULL;
	else
		return &pentry->level_xvalue;
}

// ----------------------------------------------------------------
mv_t* mlhmmv_level_look_up_and_ref_terminal(mlhmmv_level_t* plevel, sllmv_t* pmvkeys, int* perror) {
	mlhmmv_level_entry_t* plevel_entry = mlhmmv_level_look_up_and_ref_entry(plevel, pmvkeys->phead, perror);
	if (plevel_entry == NULL) {
		return NULL;
	}
	if (!plevel_entry->level_xvalue.is_terminal) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	return &plevel_entry->level_xvalue.terminal_mlrval;
}

// ----------------------------------------------------------------
mlhmmv_xvalue_t* mlhmmv_level_look_up_and_ref_xvalue(mlhmmv_level_t* pstart_level, sllmv_t* pmvkeys, int* perror) {
	*perror = MLHMMV_ERROR_NONE;
	sllmve_t* prest_keys = pmvkeys->phead;
	if (prest_keys == NULL) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	mlhmmv_level_t* plevel = pstart_level;
	if (plevel == NULL) {
		return NULL;
	}
	mlhmmv_level_entry_t* plevel_entry = mlhmmv_level_get_next_level_entry(plevel, &prest_keys->value, NULL);
	while (prest_keys->pnext != NULL) {
		if (plevel_entry == NULL) {
			return NULL;
		} else {
			plevel = plevel_entry->level_xvalue.pnext_level;
			prest_keys = prest_keys->pnext;
			if (plevel == NULL)
				return NULL;
			plevel_entry = mlhmmv_level_get_next_level_entry(plevel, &prest_keys->value, NULL);
		}
	}
	if (plevel_entry == NULL) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_DEEP;
		return NULL;
	}
	return &plevel_entry->level_xvalue;
}

// ----------------------------------------------------------------
static mlhmmv_level_t* mlhmmv_level_ref_or_create(mlhmmv_level_t* plevel, sllmve_t* prest_keys) {
	if ((plevel->num_occupied + plevel->num_freed) >= (plevel->array_length * LOAD_FACTOR))
		mlhmmv_level_enlarge(plevel);
	return mlhmmv_level_get_or_create_no_enlarge(plevel, prest_keys);
}

// ----------------------------------------------------------------
static mlhmmv_level_t* mlhmmv_level_get_or_create_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys) {
	mv_t* plevel_key = &prest_keys->value;
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == EMPTY) { // End of chain.

		plevel->states[index] = OCCUPIED;
		plevel->num_occupied++;
		pentry->ideal_index = ideal_index;
		pentry->level_key = mv_copy(plevel_key);
		pentry->level_xvalue.is_terminal = FALSE;
		pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();

		if (plevel->phead == NULL) { // First entry at this level
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {                     // Subsequent entry at this level
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}

		if (prest_keys->pnext != NULL) {
			// RECURSE
			return mlhmmv_level_ref_or_create(pentry->level_xvalue.pnext_level, prest_keys->pnext);
		} else {
			return pentry->level_xvalue.pnext_level;
		}

	} else if (plevel->states[index] == OCCUPIED) { // Existing key found in chain

		if (pentry->level_xvalue.is_terminal) {
			mv_free(&pentry->level_xvalue.terminal_mlrval);
			pentry->level_xvalue.is_terminal = FALSE;
			pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();
		}
		if (prest_keys->pnext == NULL) {
			return pentry->level_xvalue.pnext_level;
		} else { // RECURSE
			return mlhmmv_level_ref_or_create(pentry->level_xvalue.pnext_level, prest_keys->pnext);
		}

	} else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
mlhmmv_level_t* mlhmmv_level_put_empty_map(mlhmmv_level_t* plevel, mv_t* pkey) {
	int error;
	mv_t absent = mv_absent();

	sllmve_t e = {
		.value      = *pkey,
		.free_flags = 0,
		.pnext      = NULL
	};
	mlhmmv_level_put_terminal(plevel, &e, &absent); // xxx optimize to avoid 2nd lookup
	sllmv_t s = {
		.phead = &e,
		.ptail = &e,
		.length = 1
	};
	mlhmmv_xvalue_t* pxval = mlhmmv_level_look_up_and_ref_xvalue(plevel, &s, &error);
	*pxval = mlhmmv_xvalue_alloc_empty_map();
	return pxval->pnext_level;
}

// ----------------------------------------------------------------
static void mlhmmv_level_enlarge(mlhmmv_level_t* plevel) {
	mlhmmv_level_entry_t*       old_entries = plevel->entries;
	mlhmmv_level_entry_state_t* old_states  = plevel->states;
	mlhmmv_level_entry_t*       old_head    = plevel->phead;

	mlhmmv_level_init(plevel, plevel->array_length*ENLARGEMENT_FACTOR);

	for (mlhmmv_level_entry_t* pentry = old_head; pentry != NULL; pentry = pentry->pnext) {
		mlhmmv_level_move(plevel, &pentry->level_key, &pentry->level_xvalue);
	}
	free(old_entries);
	free(old_states);
}

// ----------------------------------------------------------------
// This is done only on map-level enlargement.
// Example:
// * level = map["a"], rest_keys = [2, "c"] ,   terminal_value = 4.
//                     rest_keys = ["e", "f"] , terminal_value = 7.
//                     rest_keys = [6] ,        terminal_value = "g".
//
// which is to say for the purposes of this routine
//
// * level = map["a"], level_key = 2,   level_xvalue = non-terminal ["c"] => terminal_value = 4.
//                     level_key = "e", level_xvalue = non-terminal ["f"] => terminal_value = 7.
//                     level_key = 6,   level_xvalue = terminal_value = "g".

static void mlhmmv_level_move(mlhmmv_level_t* plevel, mv_t* plevel_key, mlhmmv_xvalue_t* plevel_value) {
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == OCCUPIED) {
		// Existing key found in chain; put value.
		pentry->level_xvalue = *plevel_value;

	} else if (plevel->states[index] == EMPTY) {
		// End of chain.
		pentry->ideal_index = ideal_index;
		pentry->level_key = *plevel_key;
		// For the put API, we copy data passed in. But for internal enlarges, we just need to move pointers around.
		pentry->level_xvalue = *plevel_value;
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) {
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}
		plevel->num_occupied++;
	}
	else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
void mlhmmv_level_put_xvalue(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mlhmmv_xvalue_t* pvalue) { // xxx 'copy' into name
	if ((plevel->num_occupied + plevel->num_freed) >= (plevel->array_length * LOAD_FACTOR))
		mlhmmv_level_enlarge(plevel);
	mlhmmv_level_put_xvalue_no_enlarge(plevel, prest_keys, pvalue);
}

void mlhmmv_level_put_xvalue_singly_keyed(mlhmmv_level_t* plevel, mv_t* pkey, mlhmmv_xvalue_t* pvalue) {
	sllmve_t e = { .value = *pkey, .free_flags = 0, .pnext = NULL };
	mlhmmv_level_put_xvalue(plevel, &e, pvalue);
}

// ----------------------------------------------------------------
static void mlhmmv_level_put_xvalue_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys,
	mlhmmv_xvalue_t* pvalue)
{
	mv_t* plevel_key = &prest_keys->value;
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == EMPTY) { // End of chain.
		pentry->ideal_index = ideal_index;
		pentry->level_key = mv_copy(plevel_key); // (xxx weird & needs explaining) key is copied ...

		if (prest_keys->pnext == NULL) {
			pentry->level_xvalue = *pvalue; // (xxx weird & needs explaining) ... but the value is not copied
		} else {
			pentry->level_xvalue.is_terminal = FALSE;
			pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();
		}
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) { // First entry at this level
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {                     // Subsequent entry at this level
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}

		plevel->num_occupied++;
		if (prest_keys->pnext != NULL) {
			// RECURSE
			mlhmmv_level_put_xvalue(pentry->level_xvalue.pnext_level, prest_keys->pnext, pvalue);
		}

	} else if (plevel->states[index] == OCCUPIED) { // Existing key found in chain
		if (prest_keys->pnext == NULL) { // Place the terminal at this level
			if (pentry->level_xvalue.is_terminal) {
				mv_free(&pentry->level_xvalue.terminal_mlrval);
			} else {
				mlhmmv_level_free(pentry->level_xvalue.pnext_level);
			}
			pentry->level_xvalue = *pvalue;

		} else { // The terminal will be placed at a deeper level
			if (pentry->level_xvalue.is_terminal) {
				mv_free(&pentry->level_xvalue.terminal_mlrval);
				pentry->level_xvalue.is_terminal = FALSE;
				pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();
			}
			// RECURSE
			mlhmmv_level_put_xvalue(pentry->level_xvalue.pnext_level, prest_keys->pnext, pvalue);
		}

	} else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
// Example on recursive calls:
// * level = map, rest_keys = ["a", 2, "c"] , terminal value = 4.
// * level = map["a"], rest_keys = [2, "c"] , terminal value = 4.
// * level = map["a"][2], rest_keys = ["c"] , terminal value = 4.

void mlhmmv_level_put_terminal(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value) {
	if ((plevel->num_occupied + plevel->num_freed) >= (plevel->array_length * LOAD_FACTOR))
		mlhmmv_level_enlarge(plevel);
	mlhmmv_level_put_terminal_no_enlarge(plevel, prest_keys, pterminal_value);
}

void mlhmmv_level_put_terminal_singly_keyed(mlhmmv_level_t* plevel, mv_t* pkey, mv_t* pterminal_value) {
	sllmve_t e = { .value = *pkey, .free_flags = 0, .pnext = NULL };
	mlhmmv_level_put_terminal(plevel, &e, pterminal_value);
}

// ----------------------------------------------------------------
static void mlhmmv_level_put_terminal_no_enlarge(mlhmmv_level_t* plevel, sllmve_t* prest_keys, mv_t* pterminal_value) {
	mv_t* plevel_key = &prest_keys->value;
	int ideal_index = 0;
	int index = mlhmmv_level_find_index_for_key(plevel, plevel_key, &ideal_index);
	mlhmmv_level_entry_t* pentry = &plevel->entries[index];

	if (plevel->states[index] == EMPTY) { // End of chain.
		pentry->ideal_index = ideal_index;
		pentry->level_key = mv_copy(plevel_key); // <------------------------- copy key

		if (prest_keys->pnext == NULL) {
			pentry->level_xvalue.is_terminal = TRUE;
			pentry->level_xvalue.terminal_mlrval = mv_copy(pterminal_value); // <--------------- copy value
		} else {
			pentry->level_xvalue.is_terminal = FALSE;
			pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();
		}
		plevel->states[index] = OCCUPIED;

		if (plevel->phead == NULL) { // First entry at this level
			pentry->pprev = NULL;
			pentry->pnext = NULL;
			plevel->phead = pentry;
			plevel->ptail = pentry;
		} else {                     // Subsequent entry at this level
			pentry->pprev = plevel->ptail;
			pentry->pnext = NULL;
			plevel->ptail->pnext = pentry;
			plevel->ptail = pentry;
		}

		plevel->num_occupied++;
		if (prest_keys->pnext != NULL) {
			// RECURSE
			mlhmmv_level_put_terminal(pentry->level_xvalue.pnext_level, prest_keys->pnext, pterminal_value);
		}

	} else if (plevel->states[index] == OCCUPIED) { // Existing key found in chain
		if (prest_keys->pnext == NULL) { // Place the terminal at this level
			if (pentry->level_xvalue.is_terminal) {
				mv_free(&pentry->level_xvalue.terminal_mlrval);
			} else {
				mlhmmv_level_free(pentry->level_xvalue.pnext_level);
			}
			pentry->level_xvalue.is_terminal = TRUE;
			pentry->level_xvalue.terminal_mlrval = mv_copy(pterminal_value); // <--------------- copy value

		} else { // The terminal will be placed at a deeper level
			if (pentry->level_xvalue.is_terminal) {
				mv_free(&pentry->level_xvalue.terminal_mlrval);
				pentry->level_xvalue.is_terminal = FALSE;
				pentry->level_xvalue.pnext_level = mlhmmv_level_alloc();
			}
			// RECURSE
			mlhmmv_level_put_terminal(pentry->level_xvalue.pnext_level, prest_keys->pnext, pterminal_value);
		}

	} else {
		fprintf(stderr, "%s: mlhmmv_level_find_index_for_key did not find end of chain\n", MLR_GLOBALS.bargv0);
		exit(1);
	}
}

// ----------------------------------------------------------------
void mlhmmv_level_to_lrecs(mlhmmv_level_t* plevel, sllmv_t* pkeys, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator)
{
	mv_t* pfirstkey = &pkeys->phead->value;

	mlhmmv_level_entry_t* ptop_entry = mlhmmv_level_look_up_and_ref_entry(plevel, pkeys->phead, NULL);

	if (ptop_entry == NULL) {
		// No such entry in the mlhmmv results in no output records
	} else if (ptop_entry->level_xvalue.is_terminal) {
		// E.g. '@v = 3' at the top level of the mlhmmv.
		lrec_t* poutrec = lrec_unbacked_alloc();
		lrec_put(poutrec,
			mv_alloc_format_val(pfirstkey),
			mv_alloc_format_val(&ptop_entry->level_xvalue.terminal_mlrval), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		sllv_append(poutrecs, poutrec);
	} else {
		// E.g. '@v = {...}' at the top level of the mlhmmv: the map value keyed by oosvar-name 'v' is itself a hashmap.
		// This needs to be flattened down to an lrec which is a list of key-value pairs.  We recursively invoke
		// mlhmmv_level_to_lrecs_across_records for each of the name-list entries, one map level deeper each call, then
		// from there invoke mlhmmv_level_to_lrec_within_record on any remaining map levels.
		lrec_t* ptemplate = lrec_unbacked_alloc();
		char* oosvar_name = mv_alloc_format_val(pfirstkey);
		mlhmmv_level_to_lrecs_across_records(ptop_entry->level_xvalue.pnext_level, oosvar_name, pnames->phead,
			ptemplate, poutrecs, do_full_prefixing, flatten_separator);
		free(oosvar_name);
		lrec_free(ptemplate);
	}
}

// ----------------------------------------------------------------
static void mlhmmv_level_to_lrecs_across_records(
	mlhmmv_level_t* plevel,
	char*           prefix,
	sllmve_t*       prestnames,
	lrec_t*         ptemplate,
	sllv_t*         poutrecs,
	int             do_full_prefixing,
	char*           flatten_separator)
{
	if (prestnames != NULL) {
		// If there is a namelist entry, pull it out to its own field on the output lrecs.
		for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
			mlhmmv_xvalue_t* plevel_value = &pe->level_xvalue;
			lrec_t* pnextrec = lrec_copy(ptemplate);
			lrec_put(pnextrec,
				mv_alloc_format_val(&prestnames->value),
				mv_alloc_format_val(&pe->level_key), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			if (plevel_value->is_terminal) {
				lrec_put(pnextrec,
					mlr_strdup_or_die(prefix),
					mv_alloc_format_val(&plevel_value->terminal_mlrval), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
				sllv_append(poutrecs, pnextrec);
			} else {
				mlhmmv_level_to_lrecs_across_records(pe->level_xvalue.pnext_level,
					prefix, prestnames->pnext, pnextrec, poutrecs, do_full_prefixing, flatten_separator);
				lrec_free(pnextrec);
			}
		}

	} else {
		// If there are no more remaining namelist entries, flatten remaining map levels using the join separator
		// (default ":") and use them to create lrec values.
		lrec_t* pnextrec = lrec_copy(ptemplate);
		int emit = TRUE;
		for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
			mlhmmv_xvalue_t* plevel_value = &pe->level_xvalue;
			if (plevel_value->is_terminal) {
				char* temp = mv_alloc_format_val(&pe->level_key);
				char* next_prefix = do_full_prefixing
					? mlr_paste_3_strings(prefix, flatten_separator, temp)
					: mlr_strdup_or_die(temp);
				free(temp);
				lrec_put(pnextrec,
					next_prefix,
					mv_alloc_format_val(&plevel_value->terminal_mlrval),
					FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			} else if (do_full_prefixing) {
				char* temp = mv_alloc_format_val(&pe->level_key);
				char* next_prefix = mlr_paste_3_strings(prefix, flatten_separator, temp);
				free(temp);
				mlhmmv_level_to_lrec_within_record(plevel_value->pnext_level, next_prefix, pnextrec,
					do_full_prefixing, flatten_separator);
				free(next_prefix);
			} else {
				mlhmmv_level_to_lrecs_across_records(pe->level_xvalue.pnext_level,
					prefix, NULL, pnextrec, poutrecs, do_full_prefixing, flatten_separator);
				emit = FALSE;
			}
		}
		if (emit)
			sllv_append(poutrecs, pnextrec);
		else
			lrec_free(pnextrec);
	}
}

// ----------------------------------------------------------------
static void mlhmmv_level_to_lrec_within_record(
	mlhmmv_level_t* plevel,
	char*           prefix,
	lrec_t*         poutrec,
	int             do_full_prefixing,
	char*           flatten_separator)
{
	for (mlhmmv_level_entry_t* pe = plevel->phead; pe != NULL; pe = pe->pnext) {
		mlhmmv_xvalue_t* plevel_value = &pe->level_xvalue;
		char* temp = mv_alloc_format_val(&pe->level_key);
		char* next_prefix = do_full_prefixing
			? mlr_paste_3_strings(prefix, flatten_separator, temp)
			: mlr_strdup_or_die(temp);
		free(temp);
		if (plevel_value->is_terminal) {
			lrec_put(poutrec,
				next_prefix,
				mv_alloc_format_val(&plevel_value->terminal_mlrval),
				FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		} else {
			mlhmmv_level_to_lrec_within_record(plevel_value->pnext_level, next_prefix, poutrec,
				do_full_prefixing, flatten_separator);
			free(next_prefix);
		}
	}
}

// ----------------------------------------------------------------
static void mlhhmv_levels_to_lrecs_lashed_across_records(
	mlhmmv_level_t** pplevels,
	char**           prefixes,
	int              num_levels,
	sllmve_t*        prestnames,
	lrec_t*          ptemplate,
	sllv_t*          poutrecs,
	int              do_full_prefixing,
	char*            flatten_separator)
{
	if (prestnames != NULL) {
		// If there is a namelist entry, pull it out to its own field on the output lrecs.
		// First is iterated over and the rest are lashed (lookups with same keys as primary).
		if (pplevels[0] != NULL) {
			for (mlhmmv_level_entry_t* pe = pplevels[0]->phead; pe != NULL; pe = pe->pnext) {
				mlhmmv_xvalue_t* pfirst_level_value = &pe->level_xvalue;
				lrec_t* pnextrec = lrec_copy(ptemplate);
				lrec_put(pnextrec,
					mv_alloc_format_val(&prestnames->value),
					mv_alloc_format_val(&pe->level_key), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);

				if (pfirst_level_value->is_terminal) {
					for (int i = 0; i < num_levels; i++) {
						mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
						if (plevel_value != NULL && plevel_value->is_terminal) {
							lrec_put(pnextrec,
								mlr_strdup_or_die(prefixes[i]),
								mv_alloc_format_val(&plevel_value->terminal_mlrval), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
						}
					}
					sllv_append(poutrecs, pnextrec);
				} else {
					mlhmmv_level_t** ppnext_levels = mlr_malloc_or_die(num_levels * sizeof(mlhmmv_level_t*));
					for (int i = 0; i < num_levels; i++) {
						mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
						if (plevel_value == NULL || plevel_value->is_terminal) {
							ppnext_levels[i] = NULL;
						} else {
							ppnext_levels[i] = plevel_value->pnext_level;
						}
					}

					mlhhmv_levels_to_lrecs_lashed_across_records(ppnext_levels, prefixes, num_levels,
						prestnames->pnext, pnextrec, poutrecs, do_full_prefixing, flatten_separator);

					free(ppnext_levels);
					lrec_free(pnextrec);
				}
			}
		}

	} else if (pplevels[0] != NULL) {
		// If there are no more remaining namelist entries, flatten remaining map levels using the join separator
		// (default ":") and use them to create lrec values.
		lrec_t* pnextrec = lrec_copy(ptemplate);
		int emit = TRUE;
		for (mlhmmv_level_entry_t* pe = pplevels[0]->phead; pe != NULL; pe = pe->pnext) {
			if (pe->level_xvalue.is_terminal) {
				char* temp = mv_alloc_format_val(&pe->level_key);
				for (int i = 0; i < num_levels; i++) {
					mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
					if (plevel_value != NULL && plevel_value->is_terminal) {
						char* name = do_full_prefixing
							? mlr_paste_3_strings(prefixes[i], flatten_separator, temp)
							: mlr_strdup_or_die(temp);
						lrec_put(pnextrec,
							name,
							mv_alloc_format_val(&plevel_value->terminal_mlrval),
							FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
					}
				}
				free(temp);
			} else if (do_full_prefixing) {
				char* temp = mv_alloc_format_val(&pe->level_key);
				mlhmmv_level_t** ppnext_levels = mlr_malloc_or_die(num_levels * sizeof(mlhmmv_level_t*));
				char** next_prefixes = mlr_malloc_or_die(num_levels * sizeof(char*));
				for (int i = 0; i < num_levels; i++) {
					mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
					if (plevel_value != NULL && !plevel_value->is_terminal) {
						ppnext_levels[i] = plevel_value->pnext_level;
						next_prefixes[i] = mlr_paste_3_strings(prefixes[i], flatten_separator, temp);
					} else {
						ppnext_levels[i] = NULL;
						next_prefixes[i] = NULL;
					}
				}

				mlhhmv_levels_to_lrecs_lashed_within_records(ppnext_levels, next_prefixes, num_levels,
					pnextrec, do_full_prefixing, flatten_separator);

				for (int i = 0; i < num_levels; i++) {
					free(next_prefixes[i]);
				}
				free(next_prefixes);
				free(ppnext_levels);
				free(temp);
			} else {
				mlhmmv_level_t** ppnext_levels = mlr_malloc_or_die(num_levels * sizeof(mlhmmv_level_t*));
				for (int i = 0; i < num_levels; i++) {
					mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
					if (plevel_value->is_terminal) {
						ppnext_levels[i] = NULL;
					} else {
						ppnext_levels[i] = plevel_value->pnext_level;
					}
				}

				mlhhmv_levels_to_lrecs_lashed_across_records(ppnext_levels, prefixes, num_levels,
					NULL, pnextrec, poutrecs, do_full_prefixing, flatten_separator);

				free(ppnext_levels);
				emit = FALSE;
			}
		}
		if (emit)
			sllv_append(poutrecs, pnextrec);
		else
			lrec_free(pnextrec);
	}
}

// ----------------------------------------------------------------
static void mlhhmv_levels_to_lrecs_lashed_within_records(
	mlhmmv_level_t** pplevels,
	char**           prefixes,
	int              num_levels,
	lrec_t*          poutrec,
	int              do_full_prefixing,
	char*            flatten_separator)
{
	for (mlhmmv_level_entry_t* pe = pplevels[0]->phead; pe != NULL; pe = pe->pnext) {
		mlhmmv_xvalue_t* pfirst_level_value = &pe->level_xvalue;
		char* temp = mv_alloc_format_val(&pe->level_key);
		if (pfirst_level_value->is_terminal) {
			for (int i = 0; i < num_levels; i++) {
				mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
				if (plevel_value != NULL && plevel_value->is_terminal) {
					char* name = do_full_prefixing
						? mlr_paste_3_strings(prefixes[i], flatten_separator, temp)
						: mlr_strdup_or_die(temp);
					lrec_put(poutrec,
						name,
						mv_alloc_format_val(&plevel_value->terminal_mlrval),
						FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
				}
			}
		} else {
			mlhmmv_level_t** ppnext_levels = mlr_malloc_or_die(num_levels * sizeof(mlhmmv_level_t*));
			char** next_prefixes = mlr_malloc_or_die(num_levels * sizeof(char*));
			for (int i = 0; i < num_levels; i++) {
				mlhmmv_xvalue_t* plevel_value = mlhmmv_level_get_next_level_xvalue(pplevels[i], &pe->level_key);
				if (plevel_value->is_terminal) {
					ppnext_levels[i] = NULL;
					next_prefixes[i] = NULL;
				} else {
					ppnext_levels[i] = plevel_value->pnext_level;
					next_prefixes[i] = do_full_prefixing
						? mlr_paste_3_strings(prefixes[i], flatten_separator, temp)
						: mlr_strdup_or_die(temp);
				}
			}

			mlhhmv_levels_to_lrecs_lashed_within_records(ppnext_levels, next_prefixes, num_levels,
				poutrec, do_full_prefixing, flatten_separator);

			for (int i = 0; i < num_levels; i++) {
				free(next_prefixes[i]);
			}
			free(next_prefixes);
			free(ppnext_levels);
		}
		free(temp);
	}
}

// ----------------------------------------------------------------
void mlhmmv_level_print_stacked(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_keys_always, int quote_values_always, int json_apply_ofmt_to_floats,
	char* line_indent, char* line_term, FILE* ostream)
{
	if (plevel == NULL) {
		return;
	}
	static char* leader = "  ";
	// Top-level opening brace goes on a line by itself; subsequents on the same line after the level key.
	if (depth == 0)
		fprintf(ostream, "%s{%s", line_indent, line_term);
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		fprintf(ostream, "%s", line_indent);
		for (int i = 0; i <= depth; i++)
			fprintf(ostream, "%s", leader);
		char* level_key_string = mv_alloc_format_val(&pentry->level_key);
		if (quote_keys_always || mv_is_string_or_empty(&pentry->level_key)) {
			json_print_string_escaped(ostream, level_key_string);
		} else {
			fputs(level_key_string, ostream);
		}
		free(level_key_string);
		fprintf(ostream, ": ");

		if (pentry->level_xvalue.is_terminal) {
			mlhmmv_print_terminal(&pentry->level_xvalue.terminal_mlrval, quote_keys_always,
				quote_values_always, json_apply_ofmt_to_floats, ostream);

			if (pentry->pnext != NULL)
				fprintf(ostream, ",%s", line_term);
			else
				fprintf(ostream, "%s", line_term);
		} else {
			fprintf(ostream, "%s{%s", line_indent, line_term);
			mlhmmv_level_print_stacked(pentry->level_xvalue.pnext_level, depth + 1,
				pentry->pnext != NULL, quote_keys_always, quote_values_always, json_apply_ofmt_to_floats,
				line_indent, line_term,
				ostream);
		}
	}
	for (int i = 0; i < depth; i++)
		fprintf(ostream, "%s", leader);
	if (do_final_comma)
		fprintf(ostream, "%s},%s", line_indent, line_term);
	else
		fprintf(ostream, "%s}%s", line_indent, line_term);
}

// ----------------------------------------------------------------
static void mlhmmv_level_print_single_line(mlhmmv_level_t* plevel, int depth,
	int do_final_comma, int quote_keys_always, int quote_values_always, int json_apply_ofmt_to_floats,
	FILE* ostream)
{
	// Top-level opening brace goes on a line by itself; subsequents on the same line after the level key.
	if (depth == 0)
		fprintf(ostream, "{ ");
	for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
		char* level_key_string = mv_alloc_format_val(&pentry->level_key);
		if (quote_keys_always || mv_is_string_or_empty(&pentry->level_key)) {
			json_print_string_escaped(ostream, level_key_string);
		} else {
			fputs(level_key_string, ostream);
		}
		free(level_key_string);
		fprintf(ostream, ": ");

		if (pentry->level_xvalue.is_terminal) {
			char* level_value_string = mv_alloc_format_val(&pentry->level_xvalue.terminal_mlrval);

			if (quote_values_always) {
				json_print_string_escaped(ostream,level_value_string);
			} else if (pentry->level_xvalue.terminal_mlrval.type == MT_STRING) {
				double parsed;
				if (mlr_try_float_from_string(level_value_string, &parsed)) {
					json_decimal_print(ostream, level_value_string, json_apply_ofmt_to_floats, parsed);
				} else if (streq(level_value_string, "true") || streq(level_value_string, "false")) {
					fprintf(ostream, "%s", level_value_string);
				} else {
					json_print_string_escaped(ostream,level_value_string);
				}
			} else {
				fprintf(ostream, "%s", level_value_string);
			}

			free(level_value_string);
			if (pentry->pnext != NULL)
				fprintf(ostream, ", ");
		} else {
			fprintf(ostream, "{");
			mlhmmv_level_print_single_line(pentry->level_xvalue.pnext_level, depth + 1,
				pentry->pnext != NULL, quote_keys_always, quote_values_always, json_apply_ofmt_to_floats, ostream);
		}
	}
	if (do_final_comma)
		fprintf(ostream, " },");
	else
		fprintf(ostream, " }");
}

// ----------------------------------------------------------------
void mlhmmv_level_remove(mlhmmv_level_t* plevel, sllmve_t* prestkeys) {
	if (plevel == NULL) // nonesuch
		return;
	if (prestkeys == NULL) // restkeys too short
		return;

	int index = -1;
	mlhmmv_level_entry_t* pentry = mlhmmv_level_get_next_level_entry(plevel, &prestkeys->value, &index);
	if (pentry == NULL)
		return;

	if (prestkeys->pnext != NULL) {
		// Keep recursing until end of restkeys.
		if (pentry->level_xvalue.is_terminal) // restkeys too long
			return;
		mlhmmv_level_remove(pentry->level_xvalue.pnext_level, prestkeys->pnext);

	} else {
		// 1. Excise the node and its descendants from the storage tree
		if (plevel->states[index] != OCCUPIED) {
			fprintf(stderr, "%s: mlhmmv_root_remove: did not find end of chain.\n", MLR_GLOBALS.bargv0);
			exit(1);
		}

		mv_free(&pentry->level_key);
		pentry->ideal_index = -1;
		plevel->states[index] = DELETED;

		if (pentry == plevel->phead) {
			if (pentry == plevel->ptail) {
				plevel->phead = NULL;
				plevel->ptail = NULL;
			} else {
				plevel->phead = pentry->pnext;
				pentry->pnext->pprev = NULL;
			}
		} else if (pentry == plevel->ptail) {
				plevel->ptail = pentry->pprev;
				pentry->pprev->pnext = NULL;
		} else {
			pentry->pprev->pnext = pentry->pnext;
			pentry->pnext->pprev = pentry->pprev;
		}

		plevel->num_freed++;
		plevel->num_occupied--;

		// 2. Free the memory for the node and its descendants
		if (pentry->level_xvalue.is_terminal) {
			mv_free(&pentry->level_xvalue.terminal_mlrval);
		} else {
			mlhmmv_level_free(pentry->level_xvalue.pnext_level);
		}
	}
}

// ================================================================
mlhmmv_root_t* mlhmmv_root_alloc() {
	mlhmmv_root_t* pmap = mlr_malloc_or_die(sizeof(mlhmmv_root_t));
	pmap->root_xvalue = mlhmmv_xvalue_alloc_empty_map();
	return pmap;
}

// ----------------------------------------------------------------
void mlhmmv_root_free(mlhmmv_root_t* pmap) {
	if (pmap == NULL)
		return;
	mlhmmv_level_free(pmap->root_xvalue.pnext_level);
	free(pmap);
}

// ----------------------------------------------------------------
void mlhmmv_root_clear(mlhmmv_root_t* pmap) {
	mlhmmv_level_clear(pmap->root_xvalue.pnext_level);
}

// ----------------------------------------------------------------
mv_t* mlhmmv_root_look_up_and_ref_terminal(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, int* perror) {
	mlhmmv_level_entry_t* plevel_entry = mlhmmv_level_look_up_and_ref_entry(pmap->root_xvalue.pnext_level,
		pmvkeys->phead, perror);
	if (plevel_entry == NULL) {
		return NULL;
	}
	if (!plevel_entry->level_xvalue.is_terminal) {
		*perror = MLHMMV_ERROR_KEYLIST_TOO_SHALLOW;
		return NULL;
	}
	return &plevel_entry->level_xvalue.terminal_mlrval;
}

// ----------------------------------------------------------------
// Example on recursive calls:
// * level = map, rest_keys = ["a", 2, "c"]
// * level = map["a"], rest_keys = [2, "c"]
// * level = map["a"][2], rest_keys = ["c"]
mlhmmv_level_t* mlhmmv_root_look_up_or_create_then_ref_level(mlhmmv_root_t* pmap, sllmv_t* pmvkeys) {
	return mlhmmv_level_ref_or_create(pmap->root_xvalue.pnext_level, pmvkeys->phead);
}

// ----------------------------------------------------------------
// Example: keys = ["a", 2, "c"] and value = 4.
void mlhmmv_root_put_terminal(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, mv_t* pterminal_value) {
	mlhmmv_level_put_terminal(pmap->root_xvalue.pnext_level, pmvkeys->phead, pterminal_value);
}

// ----------------------------------------------------------------
static void mlhmmv_root_put_xvalue(mlhmmv_root_t* pmap, sllmv_t* pmvkeys, mlhmmv_xvalue_t* pvalue) {
	mlhmmv_level_put_xvalue(pmap->root_xvalue.pnext_level, pmvkeys->phead, pvalue);
}

// ----------------------------------------------------------------
void mlhmmv_root_copy_submap(mlhmmv_root_t* pmap, sllmv_t* ptokeys, sllmv_t* pfromkeys) {
	int error = 0;
	mlhmmv_level_entry_t* pfromentry = mlhmmv_level_look_up_and_ref_entry(pmap->root_xvalue.pnext_level,
		pfromkeys->phead, &error);
	if (pfromentry != NULL) {
		mlhmmv_xvalue_t submap = mlhmmv_xvalue_copy(&pfromentry->level_xvalue);
		mlhmmv_root_put_xvalue(pmap, ptokeys, &submap);
	}
}

// ----------------------------------------------------------------
sllv_t* mlhmmv_root_copy_keys_from_submap(mlhmmv_root_t* pmap, sllmv_t* pmvkeys) {
	int error;
	if (pmvkeys->length == 0) {
		mlhmmv_xvalue_t root_value = (mlhmmv_xvalue_t) {
			.is_terminal = FALSE,
			.terminal_mlrval = mv_absent(),
			.pnext_level = pmap->root_xvalue.pnext_level,
		};
		return mlhmmv_xvalue_copy_keys_nonindexed(&root_value);
	} else {
		mlhmmv_level_entry_t* pfromentry = mlhmmv_level_look_up_and_ref_entry(
			pmap->root_xvalue.pnext_level, pmvkeys->phead, &error);
		if (pfromentry != NULL) {
			return mlhmmv_xvalue_copy_keys_nonindexed(&pfromentry->level_xvalue);
		} else {
			return sllv_alloc();
		}
	}
}

// ----------------------------------------------------------------
// Removes entries from a specified level downward, unsetting any maps which become empty as a result.  For example, if
// e.g. a=>b=>c=>4 and the c level is to be removed, then all up-nodes are emptied out & should be pruned.
// * If restkeys too long (e.g. 'unset $a["b"]["c"]' with data "a":"b":3): do nothing.
// * If restkeys just right: (e.g. 'unset $a["b"]' with data "a":"b":3) remove the terminal mlrval.
// * If restkeys is too short: (e.g. 'unset $a["b"]' with data "a":"b":"c":4): remove the level and all below.

void mlhmmv_root_remove(mlhmmv_root_t* pmap, sllmv_t* prestkeys) {
	if (prestkeys == NULL) {
		return;
	} else if (prestkeys->phead == NULL) {
		mlhmmv_level_free(pmap->root_xvalue.pnext_level);
		pmap->root_xvalue.pnext_level = mlhmmv_level_alloc();
		return;
	} else {
		mlhmmv_level_remove(pmap->root_xvalue.pnext_level, prestkeys->phead);
	}
}

// ----------------------------------------------------------------
// For 'emit all' and 'emitp all'.
void mlhmmv_root_all_to_lrecs(mlhmmv_root_t* pmap, sllmv_t* pnames, sllv_t* poutrecs, int do_full_prefixing,
	char* flatten_separator)
{
	for (mlhmmv_level_entry_t* pentry = pmap->root_xvalue.pnext_level->phead; pentry != NULL; pentry = pentry->pnext) {
		sllmv_t* pkey = sllmv_single_no_free(&pentry->level_key);
		mlhmmv_root_partial_to_lrecs(pmap, pkey, pnames, poutrecs, do_full_prefixing, flatten_separator);
		sllmv_free(pkey);
	}
}

// ----------------------------------------------------------------
// For 'emit' and 'emitp': the latter has do_full_prefixing == TRUE.  These allocate lrecs, appended to the poutrecs
// list.

// * pmap is the base-level oosvar multi-level hashmap.
// * pkeys specify the level in the mlhmmv at which to produce data.
// * pnames is used to pull subsequent-level keys out into separate fields.
// * In case pnames isn't long enough to reach a terminal mlrval level in the mlhmmv,
//   do_full_prefixing specifies whether to concatenate nested mlhmmv keys into single lrec keys.
//
// Examples:

// * pkeys reaches a terminal level:
//
//   $ mlr --opprint put -q '@sum += $x; end { emit @sum }' ../data/small
//   sum
//   4.536294

// * pkeys reaches terminal levels:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emit @sum, "a", "b" }' ../data/small
//   a   b   sum
//   pan pan 0.346790
//   pan wye 0.502626
//   eks pan 0.758680
//   eks wye 0.381399
//   eks zee 0.611784
//   wye wye 0.204603
//   wye pan 0.573289
//   zee pan 0.527126
//   zee wye 0.598554
//   hat wye 0.031442

// * pkeys reaches non-terminal levels: non-prefixed:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emit @sum, "a" }' ../data/small
//   a   pan      wye
//   pan 0.346790 0.502626
//
//   a   pan      wye      zee
//   eks 0.758680 0.381399 0.611784
//
//   a   wye      pan
//   wye 0.204603 0.573289
//
//   a   pan      wye
//   zee 0.527126 0.598554
//
//   a   wye
//   hat 0.031442

// * pkeys reaches non-terminal levels: prefixed:
//
//   $ mlr --opprint put -q '@sum[$a][$b] += $x; end { emitp @sum, "a" }' ../data/small
//   a   sum:pan  sum:wye
//   pan 0.346790 0.502626
//
//   a   sum:pan  sum:wye  sum:zee
//   eks 0.758680 0.381399 0.611784
//
//   a   sum:wye  sum:pan
//   wye 0.204603 0.573289
//
//   a   sum:pan  sum:wye
//   zee 0.527126 0.598554
//
//   a   sum:wye
//   hat 0.031442

void mlhmmv_root_partial_to_lrecs(mlhmmv_root_t* pmap, sllmv_t* pkeys, sllmv_t* pnames, sllv_t* poutrecs,
	int do_full_prefixing, char* flatten_separator)
{
	// There should be at least the oosvar basename, e.g. '@a[b][c]' or '@a[b]' or '@a' but not '@'.
	MLR_INTERNAL_CODING_ERROR_IF(pkeys == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pkeys->phead == NULL);

	mlhmmv_level_to_lrecs(pmap->root_xvalue.pnext_level, pkeys, pnames, poutrecs, do_full_prefixing,
		flatten_separator);
}

// ----------------------------------------------------------------
// This is simply JSON. Example output:
// {
//   "0":  {
//     "fghij":  {
//       "0":  17
//     }
//   },
//   "3":  4,
//   "abcde":  {
//     "-6":  7
//   }
// }

void mlhmmv_root_print_json_stacked(mlhmmv_root_t* pmap, int quote_keys_always, int quote_values_always,
	int json_apply_ofmt_to_floats, char* line_indent, char* line_term, FILE* ostream)
{
	mlhmmv_level_print_stacked(pmap->root_xvalue.pnext_level, 0, FALSE, quote_keys_always,
		quote_values_always, json_apply_ofmt_to_floats, line_indent, line_term, ostream);
}

// ----------------------------------------------------------------
void mlhmmv_root_print_json_single_lines(mlhmmv_root_t* pmap, int quote_keys_always, int quote_values_always,
	int json_apply_ofmt_to_floats, char* line_term, FILE* ostream)
{
	mlhmmv_level_print_single_line(pmap->root_xvalue.pnext_level, 0, FALSE, quote_keys_always,
		quote_values_always, json_apply_ofmt_to_floats, ostream);
	fprintf(ostream, "%s", line_term);
}

// Used for emit of localvars. Puts the xvalue in a single-key-value-pair map
// keyed by the specified name. The xvalue is referenced, not copied.
mlhmmv_root_t* mlhmmv_wrap_name_and_xvalue(mv_t* pname, mlhmmv_xvalue_t* pxval) {
	mlhmmv_root_t* pmap = mlhmmv_root_alloc();
	mlhmmv_level_put_xvalue_singly_keyed(pmap->root_xvalue.pnext_level, pname, pxval);
	return pmap;
}

// Used for takedown of the temporary map returned by mlhmmv_wrap_name_and_xvalue. Since the xvalue there
// is referenced, not copied, mlhmmv_xvalue_free would prematurely free the xvalue. This method releases
// the xvalue so that the remaining, map-internal structures can be freed correctly.
void mlhmmv_unwrap_name_and_xvalue(mlhmmv_root_t* pmap)
{
	mlhmmv_level_t* plevel = pmap->root_xvalue.pnext_level;
	mlhmmv_level_entry_t* pentry = plevel->phead;
	mv_free(&pentry->level_key);
	plevel->phead = NULL;
	plevel->ptail = NULL;
	plevel->num_occupied = 0;
	mlhmmv_root_free(pmap);
}
