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
static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values);

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

	pkeeper->pleft_field_names  = slls_copy(pleft_field_names); // xxx be sure the caller frees its own
	pkeeper->pleft_field_values = NULL;
	pkeeper->precords           = sllv_alloc();
	pkeeper->prec_peek          = NULL;
	pkeeper->leof               = FALSE;
	pkeeper->state              = LEFT_STATE_0_PREFILL;

	return pkeeper;
}

// ----------------------------------------------------------------
void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper) {
	if (pkeeper->pleft_field_values != NULL)
		slls_free(pkeeper->pleft_field_values);
	if (pkeeper->precords != NULL)
		sllv_free(pkeeper->precords);
	free(pkeeper);
	pkeeper->plrec_reader->pclose_func(pkeeper->pvhandle);
}

// ----------------------------------------------------------------
// xxx cmt re who frees
void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	*ppbucket_paired        = NULL;
	*ppbucket_left_unpaired = NULL;
	int cmp = 0;

	if (pkeeper->state == LEFT_STATE_0_PREFILL) {
		// try fill Lv & peek; next state is 1,2,3 & continue from there.
		join_bucket_keeper_initial_fill(pkeeper);
		pkeeper->state = join_bucket_keeper_get_state(pkeeper);
	}

	// Return the final left-unpaireds after right EOF.
	if (pright_field_values == NULL) {
		// 1. Any records already in pkeeper->precords (current bucket)
		// 2. Peek-record, if any
		if (pkeeper->prec_peek != NULL) {
			sllv_add(pkeeper->precords, pkeeper->prec_peek);
			pkeeper->prec_peek = NULL;
		}
		// 3. Remainder of left input stream
		while (TRUE) {
			lrec_t* prec = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
				pkeeper->plrec_reader->pvstate, pkeeper->pctx);
			if (prec == NULL)
				break;
			sllv_add(pkeeper->precords, prec);
		}

		*ppbucket_left_unpaired = pkeeper->precords;
		pkeeper->precords = NULL;
		return;
	}

	if (pkeeper->state == LEFT_STATE_1_FULL || pkeeper->state == LEFT_STATE_2_LAST_BUCKET) {
		cmp = slls_compare_lexically(pkeeper->pleft_field_values, pright_field_values);
		if (cmp < 0) {
			*ppbucket_left_unpaired = pkeeper->precords;
			pkeeper->precords = sllv_alloc();
			// Advance left until match or LEOF.
			join_bucket_keeper_advance_to(pkeeper, pright_field_values);

		} else if (cmp == 0) {
			*ppbucket_paired = pkeeper->precords;
		} else {
			// No match and no need to advance left; return null lists.
		}
	} else if (pkeeper->state != LEFT_STATE_3_EOF) {
		fprintf(stderr, "%s: internal coding error: failed transition from prefill state.\n",
			MLR_GLOBALS.argv0);
		exit(1);
	}

	pkeeper->state = join_bucket_keeper_get_state(pkeeper);
}

//	// xxx cmt why rec-peek
//	if (!pkeeper->leof && pkeeper->prec_peek == NULL) {
//		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
//			pkeeper->plrec_reader->pvstate, pkeeper->pctx);
//		if (pkeeper->prec_peek == NULL)
//			pkeeper->leof = TRUE;
//	}
//
//	// xxx not quite: only if *already* leof.
//	if (pkeeper->leof) {
//		int cmp = slls_compare_lexically(pkeeper->pbucket->pjoin_values, pright_field_values);
//		// xxx think through various cases
//		if (/* xxx stub */ cmp == 999) {
//			// rename: "bucket" used at different nesting levels. it's confusing.
//			*ppbucket_paired = pkeeper->precords;
//		} else {
//			*ppbucket_left_unpaired = pkeeper->precords;
//		}
//		pkeeper->pbucket = NULL;
//		return;
//	}
//
//	// xxx now that we've got a peek record: fill the bucket with like records
//	// until there's a non-like on-deck.
//	//
//	// xxx do this only if it's time for a change.
//	//
//	// xxx rename "peek" to "on-deck"?
//
//#if 0
//	sllv_empty(pkeeper->pbucket);
//	while (TRUE) {
//		sllv_add(pkeeper->pbucket, pkeeper->prec_peek);
//		pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
//			pkeeper->plrec_reader->pvstate, pkeeper->pctx);
//		get selected keys
//		if (pkeeper->prec_peek == NULL)
//			pkeeper->leof = TRUE;
//		x
//	}
//#endif
//
//	// xxx stub
//	lrec_t* pleft_rec = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle, pkeeper->plrec_reader->pvstate,
//		pkeeper->pctx);
//	sllv_t* pfoo = sllv_alloc();
//	if (pleft_rec != NULL)
//		sllv_add(pfoo, pleft_rec);
//	*ppbucket_paired = pfoo;

// ----------------------------------------------------------------
// xxx left state:
// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true

static int join_bucket_keeper_get_state(join_bucket_keeper_t* pkeeper) {
	if (pkeeper->pleft_field_values == NULL) {
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
	pkeeper->pleft_field_values = mlr_selected_values_from_record(pkeeper->prec_peek,
		pkeeper->pleft_field_names);

	sllv_add(pkeeper->precords, pkeeper->prec_peek);
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
		int cmp = slls_compare_lexically(pkeeper->pleft_field_values, pnext_field_values);
		if (cmp != 0) {
			break;
		}
		sllv_add(pkeeper->precords, pkeeper->prec_peek);
		pkeeper->prec_peek = NULL;
	}
}

static void join_bucket_keeper_advance_to(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values) {
	// xxx finish coding
	pkeeper->prec_peek = pkeeper->plrec_reader->pprocess_func(pkeeper->pvhandle,
		pkeeper->plrec_reader->pvstate, pkeeper->pctx);
	if (pkeeper->prec_peek == NULL) {
		pkeeper->leof = TRUE;
		return;
	}
	pkeeper->pleft_field_values = mlr_selected_values_from_record(pkeeper->prec_peek,
		pkeeper->pleft_field_names);

	sllv_add(pkeeper->precords, pkeeper->prec_peek);
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
		int cmp = slls_compare_lexically(pkeeper->pleft_field_values, pnext_field_values);
		if (cmp != 0) {
			break;
		}
		sllv_add(pkeeper->precords, pkeeper->prec_peek);
		pkeeper->prec_peek = NULL;
	}
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

