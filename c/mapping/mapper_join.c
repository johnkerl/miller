#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "input/lrec_reader_stdio.h" // xxx temp
#include "input/lrec_readers.h" // xxx temp
#include "cli/argparse.h"

typedef struct _mapper_join_state_t {
	slls_t*  pleft_field_names;
	slls_t*  pright_field_names;
	slls_t*  poutput_field_names;
	// xxx prefix for left  non-join field names
	// xxx prefix for right non-join field names
	hss_t*   pleft_field_name_set;
	hss_t*   pright_field_name_set;
	int      emit_pairables;
	int      emit_left_unpairables;
	int      emit_right_unpairables;
	char*    left_file_name;
	lhmslv_t* pbuckets_by_key_field_names; // For unsorted input
} mapper_join_state_t;

static void ingest_left_file(mapper_join_state_t* pstate);

typedef struct _join_bucket_t {
	sllv_t* precords;
	int was_paired;
} join_bucket_t;

// xxx need several pieces of info to construct the reader.
// xxx may as well put all of cli_opts into MLR_GLOBALS.cli_opts_for_join.
// xxx 2nd step, allow the joiner to have its own different format/delimiter.
//  --rs      --irs
//  --fs      --ifs
//  --ps      --ips
//  --dkvp    --idkvp
//  --nidx    --inidx
//  --csv     --icsv
//  --pprint  --ipprint
//  --xtab    --ixtab

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_unsorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	if (pright_rec == NULL) { // End of input record stream
		// xxx stub: if emit_left_unpairables then dump them out here.
		// requires was-used hashset/hashmap on buckets.
		return sllv_single(NULL);
	}
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	if (pstate->pbuckets_by_key_field_names == NULL) // First call
		ingest_left_file(pstate);

	slls_t* pright_field_values = mlr_selected_values_from_record(pright_rec, pstate->pright_field_names);
	join_bucket_t* pleft_bucket = lhmslv_get(pstate->pbuckets_by_key_field_names, pright_field_values);
	if (pleft_bucket == NULL) {
		// right unpairable
		// xxx stub:
		return NULL;
	} else if (pstate->emit_pairables) {
		sllv_t* pout_records = sllv_alloc();
		pleft_bucket->was_paired = TRUE;
		for (sllve_t* pe = pleft_bucket->precords->phead; pe != NULL; pe = pe->pnext) {
			lrec_t* pleft_rec = pe->pvdata;
			lrec_t* pout_rec = lrec_unbacked_alloc();

			// add the joined-on fields
			// xxx do this conditionally on cli flags
			sllse_t* pg = pstate->pleft_field_names->phead;
			sllse_t* ph = pstate->pright_field_names->phead;
			sllse_t* pi = pstate->poutput_field_names->phead;
			for ( ; pg != NULL && ph != NULL && pi != NULL; pg = pg->pnext, ph = ph->pnext, pi = pi->pnext) {
				char* v = lrec_get(pleft_rec, pg->value);
				lrec_put(pout_rec, pi->value, strdup(v), 0);
			}

			// add the left-record fields not already added
			// xxx do this conditionally on cli flags
			for (lrece_t* pl = pleft_rec->phead; pl != NULL; pl = pl->pnext) {
				if (!hss_has(pstate->pleft_field_name_set, pl->key))
					lrec_put(pout_rec, strdup(pl->key), strdup(pl->value), 0);
			}

			// add the right-record fields not already added
			// xxx do this conditionally on cli flags
			for (lrece_t* pr = pright_rec->phead; pr != NULL; pr = pr->pnext) {
				if (!hss_has(pstate->pright_field_name_set, pr->key))
					lrec_put(pout_rec, strdup(pr->key), strdup(pr->value), LREC_FREE_ENTRY_KEY|LREC_FREE_ENTRY_VALUE);
			}

			sllv_add(pout_records, pout_rec);
		}
		return pout_records;
	} else {
		// xxx stub
		return NULL;
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_sorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	if (pright_rec == NULL) // End of input record stream
		return sllv_single(NULL);
	//mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	return sllv_single(pright_rec);
}

// ----------------------------------------------------------------
static void mapper_join_free(void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;
	if (pstate->pleft_field_names != NULL)
		slls_free(pstate->pleft_field_names);
	if (pstate->pright_field_names != NULL)
		slls_free(pstate->pright_field_names);
	if (pstate->poutput_field_names != NULL)
		slls_free(pstate->poutput_field_names);
}

// ----------------------------------------------------------------
// xxx void-abstract the stdio/mmap readers. this is a pita. also neaten up stream.c.
// xxx have a reader factory which can be called here.
// xxx for the moment, just dev with hard-coded separator & format parameters.

static void ingest_left_file(mapper_join_state_t* pstate) {

	FILE* input_stream = fopen(pstate->left_file_name, "r");
	if (input_stream == NULL) {
		fprintf(stderr, "%s: Couldn't open \"%s\" for read.\n",
			MLR_GLOBALS.argv0, pstate->left_file_name);
		perror(pstate->left_file_name);
		exit(1);
	}
	context_t ctx = { .nr = 0, .fnr = 0, .filenum = 1, .filename = pstate->left_file_name };
	context_t* pctx = &ctx;

	lrec_reader_stdio_t* plrec_reader_stdio = lrec_reader_stdio_dkvp_alloc('\n', ',', '=', FALSE);
	plrec_reader_stdio->psof_func(plrec_reader_stdio->pvstate);

	pstate->pbuckets_by_key_field_names = lhmslv_alloc();

	while (TRUE) {
		lrec_t* pleft_rec = plrec_reader_stdio->pprocess_func(input_stream, plrec_reader_stdio->pvstate, pctx);
		if (pleft_rec == NULL)
			break;

		slls_t* pleft_field_values = mlr_selected_values_from_record(pleft_rec, pstate->pleft_field_names);
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

	fclose(input_stream);
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_alloc(
	slls_t* pleft_field_names,
	slls_t* pright_field_names,
	slls_t* poutput_field_names,
	int      allow_unsorted_input,
	int      emit_pairables,
	int      emit_left_unpairables,
	int      emit_right_unpairables,
	char*    left_file_name)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_join_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_join_state_t));
	pstate->pleft_field_names           = pleft_field_names;
	pstate->pright_field_names          = pright_field_names;
	pstate->poutput_field_names         = poutput_field_names;
	pstate->pleft_field_name_set        = hss_from_slls(pleft_field_names);
	pstate->pright_field_name_set       = hss_from_slls(pright_field_names);
	pstate->emit_pairables              = emit_pairables;
	pstate->emit_left_unpairables       = emit_left_unpairables;
	pstate->emit_right_unpairables      = emit_right_unpairables;
	pstate->left_file_name              = left_file_name;
	pstate->pbuckets_by_key_field_names = NULL;

	pmapper->pvstate = (void*)pstate;
	if (allow_unsorted_input) {
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
	fprintf(stdout, "-o {a,b,c}\n");
	fprintf(stdout, "--np\n");
	fprintf(stdout, "--ul\n");
	fprintf(stdout, "--ur\n");
	fprintf(stdout, "-u\n");
	fprintf(stdout, "-e EMPTY\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv) {
	char*   left_file_name          = NULL;
	slls_t* pleft_field_names       = NULL;
	slls_t* pright_field_names      = NULL;
	slls_t* poutput_field_names     = NULL;
	int     allow_unsorted_input    = FALSE;
	int     emit_pairables          = TRUE;
	int     emit_left_unpairables   = FALSE;
	int     emit_right_unpairables  = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate,      "-f",   &left_file_name);
	ap_define_string_list_flag(pstate, "-l",   &pleft_field_names);
	ap_define_string_list_flag(pstate, "-r",   &pright_field_names);
	ap_define_string_list_flag(pstate, "-o",   &poutput_field_names);
	ap_define_false_flag(pstate,       "--np", &emit_pairables);
	ap_define_int_flag(pstate,         "--ul", &emit_left_unpairables);
	ap_define_int_flag(pstate,         "--ur", &emit_left_unpairables);
	ap_define_true_flag(pstate,        "-u",   &allow_unsorted_input);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	if (left_file_name == NULL) {
		fprintf(stderr, "%s %s: need left file name\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}
	if (pleft_field_names == NULL) {
		fprintf(stderr, "%s %s: need left field names\n", MLR_GLOBALS.argv0, verb);
		mapper_join_usage(argv[0], verb);
		return NULL;
	}

	// xxx check not all emit-flags are off -- else no output.

	// xxx check lengths equal!!
	if (pright_field_names == NULL)
		pright_field_names = slls_copy(pleft_field_names);
	if (poutput_field_names == NULL)
		poutput_field_names = slls_copy(pleft_field_names);

	return mapper_join_alloc(pleft_field_names, pright_field_names, poutput_field_names,
		allow_unsorted_input, emit_pairables, emit_left_unpairables, emit_right_unpairables,
		left_file_name);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_join_setup = {
	.verb = "join",
	.pusage_func = mapper_join_usage,
	.pparse_func = mapper_join_parse_cli,
};
