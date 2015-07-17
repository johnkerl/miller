#include <stdlib.h>
#include "lib/mlrutil.h"
#include "lib/mlr_globals.h"
#include "mapping/context.h"
#include "containers/mixutil.h"
#include "containers/join_bucket_keeper.h"
#include "input/lrec_readers.h"

#define LEFT_STATE_0_PREFILL     0
#define LEFT_STATE_1_FULL        1
#define LEFT_STATE_2_LAST_BUCKET 2
#define LEFT_STATE_3_EOF         3

// ----------------------------------------------------------------
// xxx left state:
// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true
// ----------------------------------------------------------------

// Private methods
static int join_bucket_keeper_get_state(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_initial_fill(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired);
static void join_bucket_keeper_fill(join_bucket_keeper_t* pkeeper);
static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired);

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc(
	char* left_file_name,
	char* input_file_format,
	int   use_mmap_for_read,
	char  irs,
	char  ifs,
	int   allow_repeat_ifs,
	char  ips,
	int   allow_repeat_ips,
	slls_t* pleft_field_names
) {

	lrec_reader_t* plrec_reader = lrec_reader_alloc(input_file_format, use_mmap_for_read,
		irs, ifs, allow_repeat_ifs, ips, allow_repeat_ips);

	void* pvhandle = plrec_reader->popen_func(left_file_name); // xxx move this ...
	plrec_reader->psof_func(plrec_reader->pvstate); // xxx move this ...

	context_t* pctx = mlr_malloc_or_die(sizeof(context_t));
	context_init(pctx, left_file_name);

	return join_bucket_keeper_alloc_from_reader(plrec_reader, pvhandle, pctx, pleft_field_names);
}

// ----------------------------------------------------------------
join_bucket_keeper_t* join_bucket_keeper_alloc_from_reader(
	lrec_reader_t* plrec_reader,
	void* pvhandle,
	context_t* pctx,
	slls_t* pleft_field_names
) {

	join_bucket_keeper_t* pkeeper = mlr_malloc_or_die(sizeof(join_bucket_keeper_t));

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
	if (pkeeper->pbucket->pleft_field_values != NULL)
		slls_free(pkeeper->pbucket->pleft_field_values);
	if (pkeeper->pbucket->precords != NULL)
		sllv_free(pkeeper->pbucket->precords);
	free(pkeeper);
	pkeeper->plrec_reader->pclose_func(pkeeper->pvhandle);
}

// ----------------------------------------------------------------
// xxx cmt re who frees, and which
void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	*ppbucket_paired        = NULL;
	*ppbucket_left_unpaired = NULL;
	int cmp = 0;

	if (pkeeper->state == LEFT_STATE_0_PREFILL) {
		join_bucket_keeper_initial_fill(pkeeper);
		pkeeper->state = join_bucket_keeper_get_state(pkeeper);
	}

	if (pright_field_values != NULL) {
		if (pkeeper->state == LEFT_STATE_1_FULL || pkeeper->state == LEFT_STATE_2_LAST_BUCKET) {
			cmp = slls_compare_lexically(pkeeper->pbucket->pleft_field_values, pright_field_values);
			if (cmp < 0) {
				// Advance left until match or LEOF.
				join_bucket_keeper_advance_to(pkeeper, pright_field_values, ppbucket_paired, ppbucket_left_unpaired);
			} else if (cmp == 0) {
				pkeeper->pbucket->was_paired = TRUE;
				*ppbucket_paired = pkeeper->pbucket->precords;
			} else {
				// No match and no need to advance left; return null lists.
			}
		} else if (pkeeper->state != LEFT_STATE_3_EOF) {
			fprintf(stderr, "%s: internal coding error: failed transition from prefill state.\n",
				MLR_GLOBALS.argv0);
			exit(1);
		}

	} else { // Return the final left-unpaireds after right EOF.
		join_bucket_keeper_drain(pkeeper, pright_field_values, ppbucket_paired, ppbucket_left_unpaired);
	}

	pkeeper->state = join_bucket_keeper_get_state(pkeeper);
}

// ----------------------------------------------------------------
// xxx left state:
// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true

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
	pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
		pkeeper->plrec_reader->pvstate, pkeeper->pctx);
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
	pkeeper->prec_peek = NULL;
	while (TRUE) {
		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
			pkeeper->plrec_reader->pvstate, pkeeper->pctx);
		if (pkeeper->prec_peek == NULL) {
			pkeeper->leof = TRUE;
			break;
		}
		// xxx make a function to compare w/o copy
		slls_t* pnext_field_values = mlr_selected_values_from_record(pkeeper->prec_peek,
			pkeeper->pleft_field_names);
		int cmp = slls_compare_lexically(pkeeper->pbucket->pleft_field_values, pnext_field_values);
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
// * xxx
// * xxx

//	// xxx collect into a join_bucket_t
//  join_bucket_t* pbucket {
//	  slls_t*        pleft_field_values;
//	  sllv_t*        precords;
//	  int            was_paired;
//	}
//
//	lrec_t*        prec_peek;
//	int            leof;
//	int            state;
//} join_bucket_keeper_t;

static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	if (!pkeeper->pbucket->was_paired) {
		*ppbucket_left_unpaired = pkeeper->pbucket->precords;
	} else {
		sllv_free(pkeeper->pbucket->precords);
		*ppbucket_left_unpaired = sllv_alloc();
	}
	pkeeper->pbucket->precords = sllv_alloc();

	if (pkeeper->prec_peek == NULL) { // left EOF
		return;
	}

	// +-----------+-----------+-----------+-----------+-----------+-----------+
	// |  L    R   |   L   R   |   L   R   |   L   R   |   L   R   |   L   R   |
	// + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  + ---  ---  +
	// |       a   |       a   |   e       |       a   |   e   e   |   e   e   |
	// |       b   |   e       |   e       |   e   e   |   e       |   e   e   |
	// |   e       |   e       |   e       |   e       |   e       |   e       |
	// |   e       |   e       |       f * |   e       |       f * |   g   g * |
	// |   e       |       f * |   g p     |   g       |   g p     |   g       |
	// |   g       |   g p     |   g       |   g       |   g       |           |
	// |   g       |   g       |       h * |           |           | (no p)    |
	// +-----------+-----------+-----------+-----------+-----------+-----------+
	// |           |   #####   |   #####   |           |   #####   |   #####   |
	// +-----------+-----------+-----------+-----------+-----------+-----------+

	// xxx method to compare without extracting slls
	slls_t* pleft_field_values = mlr_selected_values_from_record(pkeeper->prec_peek, pkeeper->pleft_field_names);
	int cmp = slls_compare_lexically(pleft_field_values, pright_field_values);
	if (cmp < 0) {
		// keep seeking & filling the bucket until = or >; this may or may not end up being a match.

		while (TRUE) {
			sllv_add(*ppbucket_left_unpaired, pkeeper->prec_peek);
			pkeeper->prec_peek = NULL;

			pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
				pkeeper->plrec_reader->pvstate, pkeeper->pctx);
			if (pkeeper->prec_peek == NULL) {
				pkeeper->leof = TRUE;
				break;
			}
			// xxx make a function to compare w/o copy
			slls_t* pnext_field_values = mlr_selected_values_from_record(pkeeper->prec_peek,
				pkeeper->pleft_field_names);
			cmp = slls_compare_lexically(pkeeper->pbucket->pleft_field_values, pnext_field_values);
			if (cmp >= 0)
				break;
		}

	}

	if (cmp == 0) {
		join_bucket_keeper_fill(pkeeper);
		pkeeper->pbucket->was_paired = TRUE;
		*ppbucket_paired = pkeeper->pbucket->precords;
	} else if (cmp > 0) {
		// keep seeking & filling the bucket until change of lvals; try to get a rec peek.
		// this will not be a match.
		join_bucket_keeper_fill(pkeeper);
		pkeeper->pbucket->was_paired = FALSE;
		*ppbucket_paired = pkeeper->pbucket->precords;

		// xxx undup some code from initial_fill?
	}
}

static void join_bucket_keeper_drain(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	// 1. Any records already in pkeeper->pbucket->precords (current bucket)
	if (pkeeper->pbucket->was_paired) {
		*ppbucket_paired = pkeeper->pbucket->precords;
		*ppbucket_left_unpaired = sllv_alloc();
	} else {
		*ppbucket_left_unpaired = pkeeper->pbucket->precords;
	}
	// 2. Peek-record, if any
	if (pkeeper->prec_peek != NULL) {
		sllv_add(*ppbucket_left_unpaired, pkeeper->prec_peek);
		pkeeper->prec_peek = NULL;
	}
	// 3. Remainder of left input stream
	while (TRUE) {
		lrec_t* prec = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
			pkeeper->plrec_reader->pvstate, pkeeper->pctx);
		if (prec == NULL)
			break;
		sllv_add(*ppbucket_left_unpaired, prec);
	}

	pkeeper->pbucket->precords = NULL;
}

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

// Cases:
// * 1st emit, right row <  1st left row
// * 1st emit, right row == 1st left row
// * 1st emit, right row >  1st left row
// * subsequent emit, right row <  1st left row
// * subsequent emit, right row == 1st left row
// * subsequent emit, right row >  1st left row
// * new left EOF, right row <  1st left row
// * new left EOF, right row == 1st left row
// * new left EOF, right row >  1st left row
// * old left EOF, right row <  1st left row
// * old left EOF, right row == 1st left row
// * old left EOF, right row >  1st left row

