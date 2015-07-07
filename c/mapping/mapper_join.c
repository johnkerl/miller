#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "input/lrec_readers.h"
#include "cli/argparse.h"

// xxx comment
#define OPTION_UNSPECIFIED ((char)0xff)

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

// ----------------------------------------------------------------

typedef struct _join_bucket_t {
	slls_t* pjoin_values;
	sllv_t* precords;
	int was_paired;
} join_bucket_t;

typedef struct _join_bucket_keeper_t {
	lrec_reader_t* plrec_reader;
	void*          pvhandle;
	context_t*     pctx;

	int            state;
	slls_t*        pleft_field_values;
	join_bucket_t* pbucket;
	lrec_t*        prec_peek;
	int            leof;

} join_bucket_keeper_t;

typedef struct _mapper_join_opts_t {
	// xxx prefix for left  non-join field names
	// xxx prefix for right non-join field names
	slls_t*  pleft_field_names;
	slls_t*  pright_field_names;
	slls_t*  poutput_field_names;
	int      allow_unsorted_input;
	int      emit_pairables;
	int      emit_left_unpairables;
	int      emit_right_unpairables;

	char*    left_file_name;

	// These allow the joiner to have its own different format/delimiter for
	// the left-file:
	char*    input_file_format;
	char     irs;
	char     ifs;
	char     ips;
	int      allow_repeat_ifs;
	int      allow_repeat_ips;
	char*    ifmt;
	int      use_mmap_for_read;
} mapper_join_opts_t;

typedef struct _mapper_join_state_t {
	mapper_join_opts_t* popts;

	hss_t*   pleft_field_name_set;
	hss_t*   pright_field_name_set;

	// xxx cmt for sorted
	join_bucket_keeper_t* pjoin_bucket_keeper;

	// xxx key_field -> join_field (or left_field?) thruout
	lhmslv_t* pbuckets_by_key_field_names; // For unsorted input

} mapper_join_state_t;

// xxx reorder declarations & bodies ... use more prototypes if necessary.

static void merge_options(mapper_join_opts_t* popts);
static void ingest_left_file(mapper_join_state_t* pstate);

// ----------------------------------------------------------------
static join_bucket_keeper_t* join_bucket_keeper_alloc(mapper_join_opts_t* popts) {
	join_bucket_keeper_t* pkeeper = mlr_malloc_or_die(sizeof(join_bucket_keeper_t));

	merge_options(popts);
	pkeeper->plrec_reader = lrec_reader_alloc(popts->input_file_format, popts->use_mmap_for_read,
		popts->irs, popts->ifs, popts->allow_repeat_ifs, popts->ips, popts->allow_repeat_ips);

	pkeeper->pvhandle = pkeeper->plrec_reader->popen_func(popts->left_file_name);
	pkeeper->plrec_reader->psof_func(pkeeper->plrec_reader->pvstate);

	pkeeper->pctx = mlr_malloc_or_die(sizeof(context_t));
	pkeeper->pctx->nr       = 0; // xxx make an init func & use it here & in stream.c?
	pkeeper->pctx->fnr      = 0; // xxx incr this in the readers ...
	pkeeper->pctx->filenum  = 1;
	pkeeper->pctx->filename = popts->left_file_name;

	pkeeper->pbucket               = mlr_malloc_or_die(sizeof(join_bucket_t));
	pkeeper->pbucket->pjoin_values = NULL;
	pkeeper->pbucket->precords     = sllv_alloc();
	pkeeper->pbucket->was_paired   = FALSE;
	pkeeper->prec_peek             = NULL;
	pkeeper->leof                  = FALSE;
	pkeeper->state                 = LEFT_STATE_0_PREFILL;

	return pkeeper;
}

static void join_bucket_keeper_free(join_bucket_keeper_t* pkeeper) {
	if (pkeeper->pbucket->pjoin_values != NULL)
		slls_free(pkeeper->pbucket->pjoin_values);
	if (pkeeper->pbucket != NULL)
		if (pkeeper->pbucket->precords != NULL)
			sllv_free(pkeeper->pbucket->precords);
	free(pkeeper);
	pkeeper->plrec_reader->pclose_func(pkeeper->pvhandle);
}

// ----------------------------------------------------------------
// xxx left state:
// (0) pre-fill:    Lv == null, peek == null, leof = false
// (1) midstream:   Lv != null, peek != null, leof = false
// (2) last bucket: Lv != null, peek == null, leof = true
// (3) leof:        Lv == null, peek == null, leof = true

static int join_bucket_keeper_state(join_bucket_keeper_t* pkeeper) {
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

// xxx cmt re who frees
static void join_bucket_keeper_emit(join_bucket_keeper_t* pkeeper, slls_t* pright_field_values,
	sllv_t** ppbucket_paired, sllv_t** ppbucket_left_unpaired)
{
	*ppbucket_paired        = NULL;
	*ppbucket_left_unpaired = NULL;
	int cmp = 0;

	if (pkeeper->state == LEFT_STATE_0_PREFILL) {
		// try fill Lv & peek; next state is 1,2,3 & continue from there.
		pkeeper->state = join_bucket_keeper_state(pkeeper);
	}

	switch (pkeeper->state) {
	case LEFT_STATE_1_FULL:
	case LEFT_STATE_2_LAST_BUCKET: // Intentional fall-through
		cmp = slls_compare_lexically(pkeeper->pleft_field_values, pright_field_values);
		if (cmp < 0) {
			//     Lunp <- bucket
			//     paired = null
			//     next state is 2 / 3 respectively
		} else if (cmp == 0) {
			//     Lunp = null
			//     paired = bucket
			//     next state is 1 / 2 respectively
		} else {
			//     Lunp = null
			//     paired = null
			//     next state is 1 / 2 respectively
		}

		break;

	case LEFT_STATE_3_EOF:
		break;

	default:
		fprintf(stderr, "%s: internal coding error: failed transition from prefill state.\n",
			MLR_GLOBALS.argv0);
		exit(1);
		break;
	}

	pkeeper->state = join_bucket_keeper_state(pkeeper);

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
//			*ppbucket_paired = pkeeper->pbucket->precords;
//		} else {
//			*ppbucket_left_unpaired = pkeeper->pbucket->precords;
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

}

// xxx need a drain-hook for returning the final left-unpaireds after right EOF.

// ----------------------------------------------------------------

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

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_sorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	// This can't be done in the CLI-parser since it requires information which
	// isn't known until after the CLI-parser is called.
	if (pstate->pjoin_bucket_keeper == NULL)
		pstate->pjoin_bucket_keeper = join_bucket_keeper_alloc(pstate->popts);
	join_bucket_keeper_t* pkeeper = pstate->pjoin_bucket_keeper; // keystroke-saver

	if (pright_rec == NULL) { // End of input record stream
		if (!pstate->popts->emit_left_unpairables) {
			return sllv_single(NULL);
		}

		// xxx stub
		return sllv_single(NULL);
	}

	slls_t* pright_field_values = mlr_selected_values_from_record(pright_rec, pstate->popts->pright_field_names);
	sllv_t* pleft_bucket = NULL;
	sllv_t* pbucket_left_unpaired = NULL;

	join_bucket_keeper_emit(pkeeper, pright_field_values, &pleft_bucket, &pbucket_left_unpaired);

	// xxx can we have left & right unpaired on the same call? work out some cases & ascii them in.
	sllv_t* pout_recs = sllv_alloc();

	if (pstate->popts->emit_left_unpairables) {
		if (pbucket_left_unpaired != NULL && pbucket_left_unpaired->length >= 0) {
			sllv_add_all(pout_recs, pbucket_left_unpaired);
			sllv_free(pbucket_left_unpaired);
		}
	}

	if (pstate->popts->emit_right_unpairables) {
		if (pleft_bucket->length == 0) {
			sllv_add(pout_recs, pright_rec);
		}
	}

	if (pstate->popts->emit_pairables) {

		// xxx make a method here shared between sorted & unsorted
		for (sllve_t* pe = pleft_bucket->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* pleft_rec = pe->pvdata;
			lrec_t* pout_rec = lrec_unbacked_alloc();

			// add the joined-on fields
			sllse_t* pg = pstate->popts->pleft_field_names->phead;
			sllse_t* ph = pstate->popts->pright_field_names->phead;
			sllse_t* pi = pstate->popts->poutput_field_names->phead;
			for ( ; pg != NULL && ph != NULL && pi != NULL; pg = pg->pnext, ph = ph->pnext, pi = pi->pnext) {
				char* v = lrec_get(pleft_rec, pg->value);
				lrec_put(pout_rec, pi->value, strdup(v), LREC_FREE_ENTRY_VALUE);
			}

			// add the left-record fields not already added
			for (lrece_t* pl = pleft_rec->phead; pl != NULL; pl = pl->pnext) {
				if (!hss_has(pstate->pleft_field_name_set, pl->key))
					lrec_put(pout_rec, strdup(pl->key), strdup(pl->value), LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			}

			// add the right-record fields not already added
			for (lrece_t* pr = pright_rec->phead; pr != NULL; pr = pr->pnext) {
				if (!hss_has(pstate->pright_field_name_set, pr->key))
					lrec_put(pout_rec, strdup(pr->key), strdup(pr->value), LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			}

			sllv_add(pout_recs, pout_rec);
		}

	}

	return pout_recs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_unsorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	if (pright_rec == NULL) { // End of input record stream
		if (pstate->popts->emit_left_unpairables) {
			sllv_t* poutrecs = sllv_alloc();
			for (lhmslve_t* pe = pstate->pbuckets_by_key_field_names->phead; pe != NULL; pe = pe->pnext) {
				join_bucket_t* pbucket = pe->pvvalue;
				if (!pbucket->was_paired) {
					for (sllve_t* pf = pbucket->precords->phead; pf != NULL; pf = pf->pnext) {
						lrec_t* pleft_rec = pf->pvdata;
						sllv_add(poutrecs, pleft_rec);
					}
				}
			}
			sllv_add(poutrecs, NULL);
			return poutrecs;
		} else {
			return sllv_single(NULL);
		}
	}

	// This can't be done in the CLI-parser since it requires information which
	// isn't known until after the CLI-parser is called.
	if (pstate->pbuckets_by_key_field_names == NULL) // First call
		ingest_left_file(pstate);

	slls_t* pright_field_values = mlr_selected_values_from_record(pright_rec, pstate->popts->pright_field_names);
	join_bucket_t* pleft_bucket = lhmslv_get(pstate->pbuckets_by_key_field_names, pright_field_values);
	if (pleft_bucket == NULL) {
		if (pstate->popts->emit_right_unpairables) {
			return sllv_single(pright_rec);
		} else {
			return NULL;
		}
	} else if (pstate->popts->emit_pairables) {
		sllv_t* pout_records = sllv_alloc();
		pleft_bucket->was_paired = TRUE;
		for (sllve_t* pe = pleft_bucket->precords->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* pleft_rec = pe->pvdata;
			lrec_t* pout_rec = lrec_unbacked_alloc();

			// add the joined-on fields
			sllse_t* pg = pstate->popts->pleft_field_names->phead;
			sllse_t* ph = pstate->popts->pright_field_names->phead;
			sllse_t* pi = pstate->popts->poutput_field_names->phead;
			for ( ; pg != NULL && ph != NULL && pi != NULL; pg = pg->pnext, ph = ph->pnext, pi = pi->pnext) {
				char* v = lrec_get(pleft_rec, pg->value);
				lrec_put(pout_rec, pi->value, strdup(v), LREC_FREE_ENTRY_VALUE);
			}

			// add the left-record fields not already added
			for (lrece_t* pl = pleft_rec->phead; pl != NULL; pl = pl->pnext) {
				if (!hss_has(pstate->pleft_field_name_set, pl->key))
					lrec_put(pout_rec, strdup(pl->key), strdup(pl->value), LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			}

			// add the right-record fields not already added
			for (lrece_t* pr = pright_rec->phead; pr != NULL; pr = pr->pnext) {
				if (!hss_has(pstate->pright_field_name_set, pr->key))
					lrec_put(pout_rec, strdup(pr->key), strdup(pr->value), LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			}

			sllv_add(pout_records, pout_rec);
		}
		return pout_records;
	} else {
		return NULL;
	}
}

// ----------------------------------------------------------------
static void mapper_join_free(void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;
	if (pstate->popts->pleft_field_names != NULL)
		slls_free(pstate->popts->pleft_field_names);
	if (pstate->popts->pright_field_names != NULL)
		slls_free(pstate->popts->pright_field_names);
	if (pstate->popts->poutput_field_names != NULL)
		slls_free(pstate->popts->poutput_field_names);

	if (pstate->pjoin_bucket_keeper != NULL)
		join_bucket_keeper_free(pstate->pjoin_bucket_keeper);
}

// ----------------------------------------------------------------
// Format and separator flags are passed to mapper_join in MLR_GLOBALS rather
// than on the stack, since the latter would require complicating the interface
// for all the other mappers which don't do their own file I/O.  (Also, while
// some of the information needed to construct an lrec_reader is available on
// the command line before the mapper-allocators are called, some is not
// available until after.  Hence our obtaining these flags after mapper-alloc.)

static void merge_options(mapper_join_opts_t* popts) {
	if (popts->input_file_format == NULL)
		popts->input_file_format = MLR_GLOBALS.popts->ifmt;
	if (popts->irs               == OPTION_UNSPECIFIED)
		popts->irs = MLR_GLOBALS.popts->irs;
	if (popts->ifs               == OPTION_UNSPECIFIED)
		popts->ifs = MLR_GLOBALS.popts->ifs;
	if (popts->ips               == OPTION_UNSPECIFIED)
		popts->ips = MLR_GLOBALS.popts->ips;
	if (popts->allow_repeat_ifs  == OPTION_UNSPECIFIED)
		popts->allow_repeat_ifs = MLR_GLOBALS.popts->allow_repeat_ifs;
	if (popts->allow_repeat_ips  == OPTION_UNSPECIFIED)
		popts->allow_repeat_ips = MLR_GLOBALS.popts->allow_repeat_ips;
	if (popts->use_mmap_for_read == OPTION_UNSPECIFIED)
		popts->use_mmap_for_read = MLR_GLOBALS.popts->use_mmap_for_read;
}

static void ingest_left_file(mapper_join_state_t* pstate) {
	mapper_join_opts_t* popts = pstate->popts;
	merge_options(popts);

	lrec_reader_t* plrec_reader = lrec_reader_alloc(popts->input_file_format, popts->use_mmap_for_read,
		popts->irs, popts->ifs, popts->allow_repeat_ifs, popts->ips, popts->allow_repeat_ips);

	void* pvhandle = plrec_reader->popen_func(pstate->popts->left_file_name);
	plrec_reader->psof_func(plrec_reader->pvstate);

	context_t ctx = { .nr = 0, .fnr = 0, .filenum = 1, .filename = pstate->popts->left_file_name };
	context_t* pctx = &ctx;

	pstate->pbuckets_by_key_field_names = lhmslv_alloc();

	while (TRUE) {
		lrec_t* pleft_rec = plrec_reader->pprocess_func(pvhandle, plrec_reader->pvstate, pctx);
		if (pleft_rec == NULL)
			break;

		slls_t* pleft_field_values = mlr_selected_values_from_record(pleft_rec, pstate->popts->pleft_field_names);
		join_bucket_t* pbucket = lhmslv_get(pstate->pbuckets_by_key_field_names, pleft_field_values);
		if (pbucket == NULL) { // New key-field-value: new bucket and hash-map entry
			slls_t* pkey_field_values_copy = slls_copy(pleft_field_values);
			join_bucket_t* pbucket = mlr_malloc_or_die(sizeof(join_bucket_t));
			pbucket->precords = sllv_alloc();
			pbucket->was_paired = FALSE;
			sllv_add(pbucket->precords, pleft_rec);
			lhmslv_put(pstate->pbuckets_by_key_field_names, pkey_field_values_copy, pbucket);
		} else { // Previously seen key-field-value: append record to bucket
			sllv_add(pbucket->precords, pleft_rec);
		}
	}

	plrec_reader->pclose_func(pvhandle);
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_alloc(mapper_join_opts_t* popts)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_join_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_join_state_t));
	pstate->popts = popts;
	pstate->pleft_field_name_set        = hss_from_slls(popts->pleft_field_names);
	pstate->pright_field_name_set       = hss_from_slls(popts->pright_field_names);
	pstate->pbuckets_by_key_field_names = NULL;
	pstate->pjoin_bucket_keeper = NULL;

	pmapper->pvstate = (void*)pstate;
	if (popts->allow_unsorted_input) {
		pmapper->pprocess_func = mapper_join_process_unsorted;
	} else {
		pmapper->pprocess_func = mapper_join_process_sorted;
	}
	pmapper->pfree_func = mapper_join_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_join_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "xxx write me up.\n");
	fprintf(stdout, "-f {left file name}\n");
	fprintf(stdout, "-l {a,b,c}\n");
	fprintf(stdout, "-r {a,b,c}\n");
	fprintf(stdout, "-j {a,b,c}\n");
	fprintf(stdout, "--np\n");
	fprintf(stdout, "--ul\n");
	fprintf(stdout, "--ur\n");
	fprintf(stdout, "-u\n");
	fprintf(stdout, "-e EMPTY\n"); // xxx implement this
	fprintf(stdout, "-i {xxx ffmt}\n");
	fprintf(stdout, "-ifs {etc.}\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv) {
	mapper_join_opts_t* popts = mlr_malloc_or_die(sizeof(mapper_join_opts_t));
	popts->left_file_name          = NULL;
	popts->pleft_field_names       = NULL;
	popts->pright_field_names      = NULL;
	popts->poutput_field_names     = NULL;
	popts->allow_unsorted_input    = FALSE;
	popts->emit_pairables          = TRUE;
	popts->emit_left_unpairables   = FALSE;
	popts->emit_right_unpairables  = FALSE;

	popts->input_file_format = NULL;
	popts->irs               = OPTION_UNSPECIFIED;
	popts->ifs               = OPTION_UNSPECIFIED;
	popts->ips               = OPTION_UNSPECIFIED;
	popts->allow_repeat_ifs  = OPTION_UNSPECIFIED;
	popts->allow_repeat_ips  = OPTION_UNSPECIFIED;
	popts->use_mmap_for_read = OPTION_UNSPECIFIED;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate,      "-f",   &popts->left_file_name);
	ap_define_string_list_flag(pstate, "-l",   &popts->pleft_field_names);
	ap_define_string_list_flag(pstate, "-r",   &popts->pright_field_names);
	ap_define_string_list_flag(pstate, "-j",   &popts->poutput_field_names);
	ap_define_false_flag(pstate,       "--np", &popts->emit_pairables);
	ap_define_true_flag(pstate,        "--ul", &popts->emit_left_unpairables);
	ap_define_true_flag(pstate,        "--ur", &popts->emit_right_unpairables);
	ap_define_true_flag(pstate,        "-u",   &popts->allow_unsorted_input);

	ap_define_string_flag(pstate, "-i",         &popts->input_file_format);
	ap_define_char_flag(pstate,   "--irs",      &popts->irs);
	ap_define_char_flag(pstate,   "--ifs",      &popts->ifs);
	ap_define_char_flag(pstate,   "--ips",      &popts->ips);
	ap_define_true_flag(pstate,   "--repifs",   &popts->allow_repeat_ifs);
	ap_define_true_flag(pstate,   "--repips",   &popts->allow_repeat_ips);
	ap_define_true_flag(pstate,   "--use-mmap", &popts->use_mmap_for_read);
	ap_define_false_flag(pstate,  "--no-mmap",  &popts->use_mmap_for_read);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (popts->left_file_name == NULL) {
		fprintf(stderr, "%s %s: need left file name\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (!popts->emit_pairables && !popts->emit_left_unpairables && !popts->emit_right_unpairables) {
		fprintf(stderr, "%s %s: all emit flags are unset; no output is possible.\n",
			MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (popts->poutput_field_names == NULL) {
		fprintf(stderr, "%s %s: need output field names\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}
	if (popts->pleft_field_names == NULL)
		popts->pleft_field_names = slls_copy(popts->poutput_field_names);
	if (popts->pright_field_names == NULL)
		popts->pright_field_names = slls_copy(popts->pleft_field_names);

	int llen = popts->pleft_field_names->length;
	int rlen = popts->pright_field_names->length;
	int olen = popts->poutput_field_names->length;
	if (llen != rlen || llen != olen) {
		fprintf(stderr,
			"%s %s: must have equal left,right,output field-name lists; got lengths %d,%d,%d.\n",
			MLR_GLOBALS.argv0, verb, llen, rlen, olen);
		exit(1);
	}

	return mapper_join_alloc(popts);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_join_setup = {
	.verb = "join",
	.pusage_func = mapper_join_usage,
	.pparse_func = mapper_join_parse_cli,
};
