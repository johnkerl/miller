#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "lib/context.h"
#include "containers/mixutil.h"
#include "containers/join_bucket_keeper.h"
#include "input/lrec_readers.h"

// ================================================================
// JOIN_BUCKET_KEEPER
//
// This data structure supports Miller's sorted (double-streaming) join.  It is
// perhaps best explained by first comparing with the unsorted (half-streaming)
// case.
//
// In both cases, we have left and right join keys. Suppose the left file has
// data with field name "L" to be joined with right-file(s) data with field
// name "R". For the unsorted case (see mapper_join.c) the entire left file is
// first loaded into buckets of record-lists, one for each distinct value of L.
// E.g. given the following:
//
//   +-----+-----+
//   |  L  |  R  |
//   + --- + --- +
//   |  a  |  a  |
//   |  c  |  b  |
//   |  a  |  f  |
//   |  b  |     |
//   |  c  |     |
//   |  d  |     |
//   |  a  |     |
//   +-----+-----+
//
// the left file is bucketed as
//
//   +-----+     +-----+     +-----+     +-----+
//   |  L  |     |  L  |     |  L  |     |  L  |
//   + --- +     + --- +     + --- +     + --- +
//   |  a  |     |  c  |     |  b  |     |  d  |
//   |  a  |     |  c  |     +-----+     +-----+
//   |  a  |     + --- +
//   + --- +
//
// Then the right file is processed one record at a time (hence
// "half-streaming"). The pairings are easy:
// * the right record with R=a is paired with the L=a bucket,
// * the right record with R=b is paired with the L=b bucket,
// * the right record with R=f is unpaired, and
// * the left records with L=c and L=d are unpaired.
//
// ----------------------------------------------------------------
// Now for the sorted (doubly-streaming) case. Here we require that the left
// and right files be already sorted (lexically ascending) by the join fields.
// Then the example inputs look like this:
//
//   +-----+-----+
//   |  L  |  R  |
//   + --- + --- +
//   |  a  |  a  |
//   |  a  |  b  |
//   |  a  |  f  |
//   |  b  |     |
//   |  c  |     |
//   |  c  |     |
//   |  d  |     |
//   +-----+-----+
//
// The right file is still read one record at a time. It's the job of this
// join_bucket_keeper class to keep track of the left-file buckets, one bucket
// at a time.  This includes all records with same values for the join
// field(s), e.g. the three L=a records, as well as a "peek" record which is
// either the next record with a different join value (e.g. the L=b record), or
// an end-of-file indicator.
//
// If a right-file record has join field matching the current left-file bucket,
// then it's paired with all records in that bucket. Otherwise the
// join_bucket_keeper needs to either stay with the current bucket or advance
// to the next one, depending whether the current right-file record's
// join-field values compare lexically with the the left-file bucket's
// join-field values.
//
// Examples:
//
// +-----------+-----------+-----------+-----------+-----------+-----------+
// |  L    R   |   L   R   |   L   R   |   L   R   |   L   R   |   L   R   |
// + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  +
// |       a   |       a   |   e       |       a   |   e   e   |   e   e   |
// |       b   |   e       |   e       |   e   e   |   e       |   e   e   |
// |   e       |   e       |   e       |   e       |   e       |   e       |
// |   e       |   e       |       f   |   e       |       f   |   g   g   |
// |   e       |       f   |   g       |   g       |   g       |   g       |
// |   g       |   g       |   g       |   g       |   g       |           |
// |   g       |   g       |       h   |           |           |           |
// +-----------+-----------+-----------+-----------+-----------+-----------+
//
// In all these examples, the join_bucket_keeper goes through these steps:
// * bucket is empty, peek rec has L=e
// * bucket is L=e records, peek rec has L=g
// * bucket is L=g records, peek rec is null (due to EOF)
// * bucket is empty, peek rec is null (due to EOF)
//
// Example 1:
// * left-bucket is empty and left-peek has L=e
// * right record has R=a; join_bucket_keeper does not advance
// * right record has R=b; join_bucket_keeper does not advance
// * right end of file; all left records are unpaired.
//
// Example 2:
// * left-bucket is empty and left-peek has L=e
// * right record has R=a; join_bucket_keeper does not advance
// * right record has R=f; left records with L=e are unpaired.
// * etc.
//
// ================================================================

// ----------------------------------------------------------------
#define LEFT_STATE_0_PREFILL     0
#define LEFT_STATE_1_FULL        1
#define LEFT_STATE_2_LAST_BUCKET 2
#define LEFT_STATE_3_EOF         3

// ----------------------------------------------------------------
// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true
// ----------------------------------------------------------------

// Private methods
static int join_bucket_keeper_get_state(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_initial_fill(join_bucket_keeper_t* pkeeper,
	sllv_t** pprecords_left_unpaired);
static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired);
static void join_bucket_keeper_fill(join_bucket_keeper_t* pkeeper, sllv_t** pprecords_left_unpaired);
static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired);
static char* describe_state(int state);

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc(
	char* prepipe,
	char* left_file_name,
	cli_reader_opts_t* popts,
	slls_t* pleft_field_names
) {
	lrec_reader_t* plrec_reader = lrec_reader_alloc(popts);
	return join_bucket_keeper_alloc_from_reader(plrec_reader, prepipe, left_file_name, pleft_field_names);
}

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc_from_reader(
	lrec_reader_t* plrec_reader,
	char*          prepipe,
	char*          left_file_name,
	slls_t*        pleft_field_names)
{
	join_bucket_keeper_t* pkeeper = mlr_malloc_or_die(sizeof(join_bucket_keeper_t));

	void* pvhandle = plrec_reader->popen_func(plrec_reader->pvstate, prepipe, left_file_name);
	plrec_reader->psof_func(plrec_reader->pvstate, pvhandle);

	context_t* pctx = mlr_malloc_or_die(sizeof(context_t));
	context_init_from_first_file_name(pctx, left_file_name);

	pkeeper->plrec_reader = plrec_reader;
	pkeeper->pvhandle = pvhandle;
	pkeeper->pctx = pctx;

	pkeeper->pleft_field_names           = pleft_field_names;

	pkeeper->pbucket                     = mlr_malloc_or_die(sizeof(join_bucket_t));
	pkeeper->pbucket->precords           = sllv_alloc();
	pkeeper->pbucket->pleft_field_values = NULL;
	pkeeper->pbucket->was_paired         = FALSE;

	pkeeper->prec_peek                   = NULL;
	pkeeper->leof                        = FALSE;
	pkeeper->state                       = LEFT_STATE_0_PREFILL;

	return pkeeper;
}

// ----------------------------------------------------------------
void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper, char* prepipe) {
	if (pkeeper == NULL)
		return;
	slls_free(pkeeper->pbucket->pleft_field_values);
	sllv_free(pkeeper->pbucket->precords);
	free(pkeeper->pbucket);
	pkeeper->plrec_reader->pclose_func(pkeeper->plrec_reader->pvstate, pkeeper->pvhandle, prepipe);
	pkeeper->plrec_reader->pfree_func(pkeeper->plrec_reader);
	lrec_free(pkeeper->prec_peek);
	free(pkeeper->pctx);
	free(pkeeper);
}

// ----------------------------------------------------------------
void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	*pprecords_paired        = NULL;
	*pprecords_left_unpaired = NULL;

	if (pkeeper->state == LEFT_STATE_0_PREFILL) {
		join_bucket_keeper_initial_fill(pkeeper, pprecords_left_unpaired);
		pkeeper->state = join_bucket_keeper_get_state(pkeeper);
	}

	if (pright_field_values != NULL) { // Not right EOF
		if (pkeeper->state == LEFT_STATE_1_FULL || pkeeper->state == LEFT_STATE_2_LAST_BUCKET) {
			int cmp = slls_compare_lexically(pkeeper->pbucket->pleft_field_values, pright_field_values);
			if (cmp < 0) {
				// Advance left until match or left EOF.
				join_bucket_keeper_advance_to(pkeeper, pright_field_values, pprecords_paired, pprecords_left_unpaired);
			} else if (cmp == 0) {
				pkeeper->pbucket->was_paired = TRUE;
				*pprecords_paired = pkeeper->pbucket->precords;
			} else {
				// No match and no need to advance left; return null lists.
			}
		} else if (pkeeper->state != LEFT_STATE_3_EOF) {
			fprintf(stderr, "%s: internal coding error: failed transition from prefill state.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}

	} else { // Right EOF: return the final left-unpaireds.
		join_bucket_keeper_drain(pkeeper, pright_field_values, pprecords_paired, pprecords_left_unpaired);
	}

	pkeeper->state = join_bucket_keeper_get_state(pkeeper);
}

// ----------------------------------------------------------------
static int join_bucket_keeper_get_state(join_bucket_keeper_t* pkeeper) {
	if (pkeeper->pbucket->pleft_field_values == NULL) {
		if (pkeeper->leof)
			return LEFT_STATE_3_EOF;
		else
			return LEFT_STATE_0_PREFILL;
	} else {
		if (pkeeper->prec_peek == NULL)
			return LEFT_STATE_2_LAST_BUCKET;
		else
			return LEFT_STATE_1_FULL;
	}
}

static void join_bucket_keeper_initial_fill(join_bucket_keeper_t* pkeeper,
	sllv_t** pprecords_left_unpaired)
{
	while (TRUE) {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
			pkeeper->pvhandle, pkeeper->pctx);
		if (pkeeper->prec_peek == NULL) {
			break;
		}
		if (record_has_all_keys(pkeeper->prec_peek, pkeeper->pleft_field_names)) {
			break;
		} else {
			if (*pprecords_left_unpaired == NULL)
				*pprecords_left_unpaired = sllv_alloc();
			sllv_append(*pprecords_left_unpaired, pkeeper->prec_peek);
		}
	}

	if (pkeeper->prec_peek == NULL) {
		pkeeper->leof = TRUE;
		return;
	}
	join_bucket_keeper_fill(pkeeper, pprecords_left_unpaired);
}

// Preconditions:
// * prec_peek != NULL
// * prec_peek has the join keys
static void join_bucket_keeper_fill(join_bucket_keeper_t* pkeeper, sllv_t** pprecords_left_unpaired) {
	slls_t* pleft_field_values = mlr_reference_selected_values_from_record(pkeeper->prec_peek,
		pkeeper->pleft_field_names);
	if (pleft_field_values == NULL) {
		fprintf(stderr, "%s: internal coding error: peek record should have had join keys.\n",
			MLR_GLOBALS.bargv0);
		exit(1);
	}

	pkeeper->pbucket->pleft_field_values = slls_copy(pleft_field_values);
	slls_free(pleft_field_values);
	sllv_append(pkeeper->pbucket->precords, pkeeper->prec_peek);
	pkeeper->pbucket->was_paired = FALSE;
	pkeeper->prec_peek = NULL;
	while (TRUE) {
		// Skip over records not having the join keys. These go straight to the
		// left-unpaired list.
		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
			pkeeper->pvhandle, pkeeper->pctx);
		if (pkeeper->prec_peek == NULL) {
			pkeeper->leof = TRUE;
			break;
		}

		if (record_has_all_keys(pkeeper->prec_peek, pkeeper->pleft_field_names)) {
			int cmp = slls_lrec_compare_lexically(
				pkeeper->pbucket->pleft_field_values,
				pkeeper->prec_peek,
				pkeeper->pleft_field_names);
			if (cmp != 0) {
				break;
			}
			sllv_append(pkeeper->pbucket->precords, pkeeper->prec_peek);

		} else {
			if (*pprecords_left_unpaired == NULL)
				*pprecords_left_unpaired = sllv_alloc();
			sllv_append(*pprecords_left_unpaired, pkeeper->prec_peek);
		}
		pkeeper->prec_peek = NULL;
	}
}

// Pre-conditions:
// * pkeeper->pleft_field_values < pright_field_values.
// * currently in state 1 or 2 so there is a bucket but there may or may not be a peek-record
// * current bucket was/wasn't paired on previous emits but is not paired on this emit.
// Actions:
// * if bucket was never paired, return it to the caller; else discard.
// * consume left input stream, feeding into unpaired, for as long as leftvals < rightvals && !eof.
// * if there is leftrec with vals == rightvals: parallel initial_fill.
//   else, mimic initial_fill.

static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	if (pkeeper->pbucket->was_paired) {
		while (pkeeper->pbucket->precords->phead)
			lrec_free(sllv_pop(pkeeper->pbucket->precords));
		sllv_free(pkeeper->pbucket->precords);
		pkeeper->pbucket->precords = NULL;
	} else {
		if (*pprecords_left_unpaired == NULL) {
			*pprecords_left_unpaired = pkeeper->pbucket->precords;
		} else {
			sllv_transfer(*pprecords_left_unpaired, pkeeper->pbucket->precords);
		}
	}

	pkeeper->pbucket->precords = sllv_alloc();
	if (pkeeper->pbucket->pleft_field_values != NULL) {
		slls_free(pkeeper->pbucket->pleft_field_values);
		pkeeper->pbucket->pleft_field_values = NULL;
	}
	pkeeper->pbucket->was_paired = FALSE;

	if (pkeeper->prec_peek == NULL) { // left EOF
		return;
	}

	// Need a double condition here ... the peek record is either het or hom.
	// (Or, change that: -> ensure elsewhere the peek record is hom.)
	// The former is destined for lunp and shouldn't be lexcmped. The latter
	// should be.

	int cmp = lrec_slls_compare_lexically(pkeeper->prec_peek, pkeeper->pleft_field_names, pright_field_values);
	if (cmp < 0) {
		// keep seeking & filling the bucket until = or >; this may or may not end up being a match.

		if (*pprecords_left_unpaired == NULL)
			*pprecords_left_unpaired = sllv_alloc();

		while (TRUE) {
			sllv_append(*pprecords_left_unpaired, pkeeper->prec_peek);
			pkeeper->prec_peek = NULL;

			while (TRUE) {
				// Skip over records not having the join keys. These go straight to the
				// left-unpaired list.
				pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
					pkeeper->pvhandle, pkeeper->pctx);
				if (pkeeper->prec_peek == NULL)
					break;
				if (record_has_all_keys(pkeeper->prec_peek, pkeeper->pleft_field_names)) {
					break;
				} else {
					if (*pprecords_left_unpaired == NULL)
						*pprecords_left_unpaired = sllv_alloc();
					sllv_append(*pprecords_left_unpaired, pkeeper->prec_peek);
				}
			}


			if (pkeeper->prec_peek == NULL) {
				pkeeper->leof = TRUE;
				break;
			}

			cmp = lrec_slls_compare_lexically(pkeeper->prec_peek, pkeeper->pleft_field_names, pright_field_values);
			if (cmp >= 0)
				break;
		}
	}

	if (cmp == 0) {
		join_bucket_keeper_fill(pkeeper, pprecords_left_unpaired);
		pkeeper->pbucket->was_paired = TRUE;
		*pprecords_paired = pkeeper->pbucket->precords;
	} else if (cmp > 0) {
		join_bucket_keeper_fill(pkeeper, pprecords_left_unpaired);
	}
}

static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	// 1. Any records already in pkeeper->pbucket->precords (current bucket)
	if (pkeeper->pbucket->was_paired) {
		if (*pprecords_left_unpaired == NULL)
			*pprecords_left_unpaired = sllv_alloc();
	} else {
		if (*pprecords_left_unpaired == NULL) {
			*pprecords_left_unpaired = pkeeper->pbucket->precords;
		} else {
			sllv_transfer(*pprecords_left_unpaired, pkeeper->pbucket->precords);
			sllv_free(pkeeper->pbucket->precords);
		}
	}
	// 2. Peek-record, if any
	if (pkeeper->prec_peek != NULL) {
		sllv_append(*pprecords_left_unpaired, pkeeper->prec_peek);
		pkeeper->prec_peek = NULL;
	}
	// 3. Remainder of left input stream
	while (TRUE) {
		lrec_t* prec = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
			pkeeper->pvhandle, pkeeper->pctx);
		if (prec == NULL)
			break;
		sllv_append(*pprecords_left_unpaired, prec);
	}

	pkeeper->pbucket->precords = NULL;
}

// ----------------------------------------------------------------
void join_bucket_keeper_print(join_bucket_keeper_t* pkeeper) {
	printf("pbucket at %p:\n", pkeeper);
	printf("  pvhandle = %p\n", pkeeper->pvhandle);
	context_print(pkeeper->pctx, "  ");
	printf("  pleft_field_names = ");
	slls_print(pkeeper->pleft_field_names);
	printf("\n");
	join_bucket_print(pkeeper->pbucket, "  ");
	printf("  prec_peek = ");
	if (pkeeper->prec_peek == NULL) {
		printf("null\n");
	} else {
		lrec_print(pkeeper->prec_peek);
	}
	printf("  leof  = %d\n", pkeeper->leof);
	printf("  state = %s\n", describe_state(pkeeper->state));
}

void join_bucket_keeper_print_aux(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	join_bucket_keeper_print(pkeeper);
	printf("  pright_field_values = ");
	slls_print(pright_field_values);
	printf("\n");
	printf("  precords_paired =\n");
	lrec_print_list_with_prefix(*pprecords_paired, "      ");
	printf("\n");
	printf("  precords_left_unpaired =\n");
	lrec_print_list_with_prefix(*pprecords_left_unpaired, "      ");
	printf("\n");
}

void join_bucket_print(join_bucket_t* pbucket, char* indent) {
	printf("%spbucket at %p:\n", indent, pbucket);
	printf("%s  pleft_field_values = ", indent);
	slls_print(pbucket->pleft_field_values);
	printf("\n");
	if (pbucket->precords == NULL) {
		printf("%s  precords:\n", indent);
		printf("%s    (null)\n", indent);
	} else {
		printf("%s  precords (length=%llu):\n", indent, pbucket->precords->length);
		lrec_print_list_with_prefix(pbucket->precords, "      ");
	}
	printf("%s  was_paired = %d\n", indent, pbucket->was_paired);
}

static char* describe_state(int state) {
	switch (state) {
	case LEFT_STATE_0_PREFILL:     return "LEFT_STATE_0_PREFILL";
	case LEFT_STATE_1_FULL:        return "LEFT_STATE_1_FULL";
	case LEFT_STATE_2_LAST_BUCKET: return "LEFT_STATE_2_LAST_BUCKET";
	case LEFT_STATE_3_EOF:         return "LEFT_STATE_3_EOF";
	default:                       return "???";
	}
}
