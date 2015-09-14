#include <stdio.h>
#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "mapping/context.h"
#include "containers/mixutil.h"
#include "containers/join_bucket_keeper.h"
#include "input/lrec_readers.h"

// xxx overview here ...

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
static void join_bucket_keeper_initial_fill(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired);
static void join_bucket_keeper_fill(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired);

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc(
	char* left_file_name,
	char* input_file_format,
	int   use_mmap_for_read,
	char* irs,
	char* ifs,
	int   allow_repeat_ifs,
	char* ips,
	int   allow_repeat_ips,
	slls_t* pleft_field_names
) {

	lrec_reader_t* plrec_reader = lrec_reader_alloc(input_file_format, use_mmap_for_read,
		irs, ifs, allow_repeat_ifs, ips, allow_repeat_ips);

	return join_bucket_keeper_alloc_from_reader(plrec_reader, left_file_name, pleft_field_names);
}

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc_from_reader(
	lrec_reader_t* plrec_reader,
	char*          left_file_name,
	slls_t*        pleft_field_names)
{

	join_bucket_keeper_t* pkeeper = mlr_malloc_or_die(sizeof(join_bucket_keeper_t));

	void* pvhandle = plrec_reader->popen_func(plrec_reader->pvstate, left_file_name);
	plrec_reader->psof_func(plrec_reader->pvstate);

	context_t* pctx = mlr_malloc_or_die(sizeof(context_t));
	context_init(pctx, left_file_name);

	pkeeper->plrec_reader = plrec_reader;
	pkeeper->pvhandle = pvhandle;
	pkeeper->pctx = pctx;

	pkeeper->pleft_field_names           = slls_copy(pleft_field_names); // xxx be sure the caller frees its own

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
void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper) {
	if (pkeeper == NULL)
		return;
	if (pkeeper->pbucket->pleft_field_values != NULL)
		slls_free(pkeeper->pbucket->pleft_field_values);
	if (pkeeper->pbucket->precords != NULL)
		sllv_free(pkeeper->pbucket->precords);
	pkeeper->plrec_reader->pclose_func(pkeeper->plrec_reader->pvstate, pkeeper->pvhandle);
	free(pkeeper);
}

// ----------------------------------------------------------------
void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	*pprecords_paired        = NULL;
	*pprecords_left_unpaired = NULL;
	int cmp = 0;

	if (pkeeper->state == LEFT_STATE_0_PREFILL) {
		join_bucket_keeper_initial_fill(pkeeper);
		pkeeper->state = join_bucket_keeper_get_state(pkeeper);
	}

	if (pright_field_values != NULL) { // Not right EOF
		if (pkeeper->state == LEFT_STATE_1_FULL || pkeeper->state == LEFT_STATE_2_LAST_BUCKET) {
			cmp = slls_compare_lexically(pkeeper->pbucket->pleft_field_values, pright_field_values);
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
				MLR_GLOBALS.argv0);
			exit(1);
		}

	} else { // Right EOF: return the final left-unpaireds.
		join_bucket_keeper_drain(pkeeper, pright_field_values, pprecords_paired, pprecords_left_unpaired);
	}

	pkeeper->state = join_bucket_keeper_get_state(pkeeper);
}

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

static void join_bucket_keeper_initial_fill(join_bucket_keeper_t* pkeeper) {
	pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
		pkeeper->pvhandle, pkeeper->pctx);
	if (pkeeper->prec_peek == NULL) {
		pkeeper->leof = TRUE;
		return;
	}
	join_bucket_keeper_fill(pkeeper);
}

// xxx preconditions:
// * prec_peek != NULL
static void join_bucket_keeper_fill(join_bucket_keeper_t* pkeeper) {
	pkeeper->pbucket->pleft_field_values = mlr_selected_values_from_record(pkeeper->prec_peek,
		pkeeper->pleft_field_names);
	sllv_add(pkeeper->pbucket->precords, pkeeper->prec_peek);
	pkeeper->pbucket->was_paired = FALSE;
	pkeeper->prec_peek = NULL;
	while (TRUE) {
		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
			pkeeper->pvhandle, pkeeper->pctx);
		if (pkeeper->prec_peek == NULL) {
			pkeeper->leof = TRUE;
			break;
		}

		int cmp = slls_lrec_compare_lexically(
			pkeeper->pbucket->pleft_field_values,
			pkeeper->prec_peek,
			pkeeper->pleft_field_names);

		if (cmp != 0) {
			break;
		}
		sllv_add(pkeeper->pbucket->precords, pkeeper->prec_peek);
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
//   else ... parallel initial_fill. :)
// Post-conditions:
// * xxx

static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	if (pkeeper->pbucket->was_paired) {
		sllv_free(pkeeper->pbucket->precords);
	} else {
		*pprecords_left_unpaired = pkeeper->pbucket->precords;
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

	int cmp = lrec_slls_compare_lexically(pkeeper->prec_peek, pkeeper->pleft_field_names, pright_field_values);
	if (cmp < 0) {
		// keep seeking & filling the bucket until = or >; this may or may not end up being a match.

		if (*pprecords_left_unpaired == NULL)
			*pprecords_left_unpaired = sllv_alloc();

		while (TRUE) {
			sllv_add(*pprecords_left_unpaired, pkeeper->prec_peek);
			pkeeper->prec_peek = NULL;

			pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
				pkeeper->pvhandle, pkeeper->pctx);
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
		join_bucket_keeper_fill(pkeeper);
		pkeeper->pbucket->was_paired = TRUE;
		*pprecords_paired = pkeeper->pbucket->precords;
	} else if (cmp > 0) {
		join_bucket_keeper_fill(pkeeper);
	}
}

static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** pprecords_paired, sllv_t** pprecords_left_unpaired)
{
	// 1. Any records already in pkeeper->pbucket->precords (current bucket)
	if (pkeeper->pbucket->was_paired) {
		*pprecords_left_unpaired = sllv_alloc();
	} else {
		*pprecords_left_unpaired = pkeeper->pbucket->precords;
	}
	// 2. Peek-record, if any
	if (pkeeper->prec_peek != NULL) {
		sllv_add(*pprecords_left_unpaired, pkeeper->prec_peek);
		pkeeper->prec_peek = NULL;
	}
	// 3. Remainder of left input stream
	while (TRUE) {
		lrec_t* prec = pkeeper->plrec_reader->pprocess_func(pkeeper->plrec_reader->pvstate,
			pkeeper->pvhandle, pkeeper->pctx);
		if (prec == NULL)
			break;
		sllv_add(*pprecords_left_unpaired, prec);
	}

	pkeeper->pbucket->precords = NULL;
}

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
	printf("  state = %d\n", pkeeper->state);
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
		printf("%s  precords (length=%d):\n", indent, pbucket->precords->length);
		lrec_print_list_with_prefix(pbucket->precords, "      ");
	}
	printf("%s  was_paired = %d\n", indent, pbucket->was_paired);
}
