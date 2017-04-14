#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/lrec.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "containers/join_bucket_keeper.h"
#include "mapping/mappers.h"
#include "input/lrec_readers.h"

// ----------------------------------------------------------------
typedef struct _mapper_join_opts_t {
	char*    left_prefix;
	char*    right_prefix;
	slls_t*  poutput_join_field_names;
	slls_t*  pleft_join_field_names;
	slls_t*  pright_join_field_names;
	char*    output_join_field_names_cli_backing;
	char*    left_join_field_names_cli_backing;
	char*    right_join_field_names_cli_backing;
	int      allow_unsorted_input;
	int      emit_pairables;
	int      emit_left_unpairables;
	int      emit_right_unpairables;

	char*    prepipe;
	char*    left_file_name;

	// These allow the joiner to have its own different format/delimiter for
	// the left-file:
	cli_reader_opts_t reader_opts;
} mapper_join_opts_t;

typedef struct _mapper_join_state_t {

	mapper_join_opts_t* popts;

	hss_t*   pleft_field_name_set;
	hss_t*   pright_field_name_set;

	// For sorted input
	join_bucket_keeper_t* pjoin_bucket_keeper;

	// For unsorted input
	lhmslv_t* pleft_buckets_by_join_field_values;
	sllv_t*   pleft_unpaired_records;

} mapper_join_state_t;

// ----------------------------------------------------------------
static void mapper_join_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_join_alloc(mapper_join_opts_t* popts);
static void mapper_join_free(mapper_t* pmapper, context_t* _);
static void ingest_left_file(mapper_join_state_t* pstate);
static void mapper_join_form_pairs(sllv_t* pleft_records, lrec_t* pright_rec, mapper_join_state_t* pstate,
	sllv_t* pout_recs);
static sllv_t* mapper_join_process_sorted(lrec_t* pright_rec, context_t* pctx, void* pvstate);
static sllv_t* mapper_join_process_unsorted(lrec_t* pright_rec, context_t* pctx, void* pvstate);

mapper_setup_t mapper_join_setup = {
	.verb = "join",
	.pusage_func = mapper_join_usage,
	.pparse_func = mapper_join_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_join_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Joins records from specified left file name with records from all file names\n");
	fprintf(o, "at the end of the Miller argument list.\n");
	fprintf(o, "Functionality is essentially the same as the system \"join\" command, but for\n");
	fprintf(o, "record streams.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "  -f {left file name}\n");
	fprintf(o, "  -j {a,b,c}   Comma-separated join-field names for output\n");
	fprintf(o, "  -l {a,b,c}   Comma-separated join-field names for left input file;\n");
	fprintf(o, "               defaults to -j values if omitted.\n");
	fprintf(o, "  -r {a,b,c}   Comma-separated join-field names for right input file(s);\n");
	fprintf(o, "               defaults to -j values if omitted.\n");
	fprintf(o, "  --lp {text}  Additional prefix for non-join output field names from\n");
	fprintf(o, "               the left file\n");
	fprintf(o, "  --rp {text}  Additional prefix for non-join output field names from\n");
	fprintf(o, "               the right file(s)\n");
	fprintf(o, "  --np         Do not emit paired records\n");
	fprintf(o, "  --ul         Emit unpaired records from the left file\n");
	fprintf(o, "  --ur         Emit unpaired records from the right file(s)\n");
	fprintf(o, "  -s|--sorted-input  Require sorted input: records must be sorted\n");
	fprintf(o, "               lexically by their join-field names, else not all records will\n");
	fprintf(o, "               be paired. The only likely use case for this is with a left\n");
	fprintf(o, "               file which is too big to fit into system memory otherwise.\n");
	fprintf(o, "  -u           Enable unsorted input. (This is the default even without -u.)\n");
	fprintf(o, "               In this case, the entire left file will be loaded into memory.\n");

	fprintf(o, "  --prepipe {command} As in main input options; see %s --help for details.\n",
		MLR_GLOBALS.bargv0);
	fprintf(o, "               If you wish to use a prepipe command for the main input as well\n");
	fprintf(o, "               as here, it must be specified there as well as here.\n");
	fprintf(o, "File-format options default to those for the right file names on the Miller\n");
	fprintf(o, "argument list, but may be overridden for the left file as follows. Please see\n");
	fprintf(o, "the main \"%s --help\" for more information on syntax for these arguments.\n", argv0);
	fprintf(o, "  -i {one of csv,dkvp,nidx,pprint,xtab}\n");
	fprintf(o, "  --irs {record-separator character}\n");
	fprintf(o, "  --ifs {field-separator character}\n");
	fprintf(o, "  --ips {pair-separator character}\n");
	fprintf(o, "  --repifs\n");
	fprintf(o, "  --repips\n");
	fprintf(o, "  --use-mmap\n");
	fprintf(o, "  --no-mmap\n");
	fprintf(o, "Please use \"%s --usage-separator-options\" for information on specifying separators.\n",
		argv0);
	fprintf(o, "Please see http://johnkerl.org/miller/doc/reference.html for more information\n");
	fprintf(o, "including examples.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* pmain_reader_opts, cli_writer_opts_t* __)
{
	mapper_join_opts_t* popts = mlr_malloc_or_die(sizeof(mapper_join_opts_t));

	cli_reader_opts_init(&popts->reader_opts);

	popts->left_prefix                         = NULL;
	popts->right_prefix                        = NULL;
	popts->prepipe                             = NULL;
	popts->left_file_name                      = NULL;
	popts->poutput_join_field_names            = NULL;
	popts->pleft_join_field_names              = NULL;
	popts->pright_join_field_names             = NULL;
	popts->output_join_field_names_cli_backing = NULL;
	popts->left_join_field_names_cli_backing   = NULL;
	popts->right_join_field_names_cli_backing  = NULL;
	popts->emit_pairables                      = TRUE;
	popts->emit_left_unpairables               = FALSE;
	popts->emit_right_unpairables              = FALSE;
	popts->allow_unsorted_input                = TRUE;

	int argi = *pargi;
	char* verb = argv[argi++];

	for (; argi < argc; /* variable increment: 1 or 2 depending on flag */) {

		if (argv[argi][0] != '-') {
			break; // No more flag options to process

		} else if (cli_handle_reader_options(argv, argc, &argi, &popts->reader_opts)) {
			// handled

		} else if (streq(argv[argi], "--prepipe")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->prepipe = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "-f")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->left_file_name = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "-j")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->output_join_field_names_cli_backing = mlr_strdup_or_die(argv[argi+1]);
			popts->poutput_join_field_names = slls_from_line(popts->output_join_field_names_cli_backing, ',', FALSE);
			argi += 2;

		} else if (streq(argv[argi], "-l")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->left_join_field_names_cli_backing   = mlr_strdup_or_die(argv[argi+1]);
			popts->pleft_join_field_names = slls_from_line(popts->left_join_field_names_cli_backing, ',', FALSE);
			argi += 2;

		} else if (streq(argv[argi], "-r")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->right_join_field_names_cli_backing  = mlr_strdup_or_die(argv[argi+1]);
			popts->pright_join_field_names = slls_from_line(popts->right_join_field_names_cli_backing, ',', FALSE);
			argi += 2;

		} else if (streq(argv[argi], "--lp")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->left_prefix = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "--rp")) {
			if ((argc - argi) < 2) {
				mapper_join_usage(stderr, argv[0], verb);
				return NULL;
			}
			popts->right_prefix = argv[argi+1];
			argi += 2;

		} else if (streq(argv[argi], "--np")) {
			popts->emit_pairables = FALSE;
			argi += 1;

		} else if (streq(argv[argi], "--ul")) {
			popts->emit_left_unpairables = TRUE;
			argi += 1;

		} else if (streq(argv[argi], "--ur")) {
			popts->emit_right_unpairables = TRUE;
			argi += 1;

		} else if (streq(argv[argi], "-u")) {
			popts->allow_unsorted_input = TRUE;
			argi += 1;

		} else if (streq(argv[argi], "--sorted-input") || streq(argv[argi], "-s")) {
			popts->allow_unsorted_input = FALSE;
			argi += 1;

		} else {
			mapper_join_usage(stderr, argv[0], verb);
			return NULL;
		}

	}

	cli_merge_reader_opts(&popts->reader_opts, pmain_reader_opts);

	// popen is a stdio construct, not an mmap construct, and it can't be supported here.
	if (popts->prepipe != NULL)
		popts->reader_opts.use_mmap_for_read = FALSE;

	if (popts->left_file_name == NULL) {
		fprintf(stderr, "%s %s: need left file name\n", MLR_GLOBALS.bargv0, verb);
		mapper_join_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (!popts->emit_pairables && !popts->emit_left_unpairables && !popts->emit_right_unpairables) {
		fprintf(stderr, "%s %s: all emit flags are unset; no output is possible.\n",
			MLR_GLOBALS.bargv0, verb);
		mapper_join_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (popts->poutput_join_field_names == NULL) {
		fprintf(stderr, "%s %s: need output field names\n", MLR_GLOBALS.bargv0, verb);
		mapper_join_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (popts->pleft_join_field_names == NULL) {
		popts->pleft_join_field_names = slls_copy(popts->poutput_join_field_names);
	}
	if (popts->pright_join_field_names == NULL) {
		popts->pright_join_field_names = slls_copy(popts->poutput_join_field_names);
	}

	int llen = popts->pleft_join_field_names->length;
	int rlen = popts->pright_join_field_names->length;
	int olen = popts->poutput_join_field_names->length;
	if (llen != rlen || llen != olen) {
		fprintf(stderr,
			"%s %s: must have equal left,right,output field-name lists; got lengths %d,%d,%d.\n",
			MLR_GLOBALS.bargv0, verb, llen, rlen, olen);
		exit(1);
	}

	*pargi = argi;

	return mapper_join_alloc(popts);
}

// ----------------------------------------------------------------
static mapper_t* mapper_join_alloc(mapper_join_opts_t* popts) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_join_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_join_state_t));
	pstate->popts                              = popts;
	pstate->pleft_field_name_set               = hss_from_slls(popts->pleft_join_field_names);
	pstate->pright_field_name_set              = hss_from_slls(popts->pright_join_field_names);

	pstate->pjoin_bucket_keeper = join_bucket_keeper_alloc(
		popts->prepipe,
		popts->left_file_name,
		&popts->reader_opts,
		popts->pleft_join_field_names);

	pstate->pleft_buckets_by_join_field_values = NULL;
	pstate->pleft_unpaired_records             = NULL;

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
static void mapper_join_free(mapper_t* pmapper, context_t* _) {
	mapper_join_state_t* pstate = pmapper->pvstate;

	if (pstate->pleft_buckets_by_join_field_values != NULL) {
		for (lhmslve_t* pe = pstate->pleft_buckets_by_join_field_values->phead; pe != NULL; pe = pe->pnext) {
			join_bucket_t* pbucket = pe->pvvalue;
			slls_free(pbucket->pleft_field_values);
			if (pbucket->precords)
				while (pbucket->precords->phead)
					lrec_free(sllv_pop(pbucket->precords));
			sllv_free(pbucket->precords);
			free(pbucket);
		}
		lhmslv_free(pstate->pleft_buckets_by_join_field_values);
	}

	// The void-star payload, which is lrec_t*'s, should have been sllv_transferred out.
	// Misses should be detected by valgrind --leak-check=full, e.g. reg_test/run --valgrind.
	sllv_free(pstate->pleft_unpaired_records);

	join_bucket_keeper_free(pstate->pjoin_bucket_keeper, pstate->popts->prepipe);

	slls_free(pstate->popts->poutput_join_field_names);
	slls_free(pstate->popts->pleft_join_field_names);
	slls_free(pstate->popts->pright_join_field_names);

	free(pstate->popts->output_join_field_names_cli_backing);
	free(pstate->popts->left_join_field_names_cli_backing);
	free(pstate->popts->right_join_field_names_cli_backing);

	hss_free(pstate->pleft_field_name_set);
	hss_free(pstate->pright_field_name_set);

	free(pstate->popts);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_sorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	join_bucket_keeper_t* pkeeper = pstate->pjoin_bucket_keeper; // keystroke-saver

	sllv_t* pleft_records = NULL;
	sllv_t* pbucket_left_unpaired = NULL;
	sllv_t* pout_recs = sllv_alloc();

	if (pright_rec == NULL) { // End of input record stream
		join_bucket_keeper_emit(pkeeper, NULL, &pleft_records, &pbucket_left_unpaired);
		if (pstate->popts->emit_left_unpairables) {
			sllv_transfer(pout_recs, pbucket_left_unpaired);
		} else {
			if (pbucket_left_unpaired)
				while (pbucket_left_unpaired->phead)
					lrec_free(sllv_pop(pbucket_left_unpaired));
		}
		// pleft_records are not caller-owned; we don't free them.
		sllv_free(pbucket_left_unpaired);
		sllv_append(pout_recs, NULL);
		return pout_recs;
	}

	int produced_right = FALSE;
	slls_t* pright_field_values = mlr_reference_selected_values_from_record(pright_rec,
		pstate->popts->pright_join_field_names);

	if (pright_field_values != NULL) {
		join_bucket_keeper_emit(pkeeper, pright_field_values, &pleft_records, &pbucket_left_unpaired);
		slls_free(pright_field_values);
	}

	if (pstate->popts->emit_left_unpairables) {
		if (pbucket_left_unpaired != NULL && pbucket_left_unpaired->length > 0) {
			sllv_transfer(pout_recs, pbucket_left_unpaired);
		}
	}

	if (pstate->popts->emit_right_unpairables) {
		if (pleft_records == NULL || pleft_records->length == 0) {
			sllv_append(pout_recs, pright_rec);
			produced_right = TRUE;
		}
	}

	if (pstate->popts->emit_pairables && pleft_records != NULL) {
		mapper_join_form_pairs(pleft_records, pright_rec, pstate, pout_recs);
	}

	if (!produced_right)
		lrec_free(pright_rec);
	// pleft_records are not caller-owned; we don't free them.
	if (pbucket_left_unpaired)
		while (pbucket_left_unpaired->phead)
			lrec_free(sllv_pop(pbucket_left_unpaired));
	sllv_free(pbucket_left_unpaired);
	return pout_recs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_join_process_unsorted(lrec_t* pright_rec, context_t* pctx, void* pvstate) {
	mapper_join_state_t* pstate = (mapper_join_state_t*)pvstate;

	if (pstate->pleft_unpaired_records == NULL) // First call
		pstate->pleft_unpaired_records = sllv_alloc();

	// This can't be done in the CLI-parser since it requires information which
	// isn't known until after the CLI-parser is called.
	if (pstate->pleft_buckets_by_join_field_values == NULL) // First call
		ingest_left_file(pstate);

	if (pright_rec == NULL) { // End of input record stream
		if (pstate->popts->emit_left_unpairables) {
			sllv_t* poutrecs = sllv_alloc();
			if (pstate->pleft_buckets_by_join_field_values != NULL) { // E.g. empty right input
				for (lhmslve_t* pe = pstate->pleft_buckets_by_join_field_values->phead; pe != NULL; pe = pe->pnext) {
					join_bucket_t* pbucket = pe->pvvalue;
					if (!pbucket->was_paired) {
						sllv_transfer(poutrecs, pbucket->precords);
					}
				}
			}
			sllv_transfer(poutrecs, pstate->pleft_unpaired_records);
			sllv_append(poutrecs, NULL);
			return poutrecs;
		} else {
			while (pstate->pleft_unpaired_records->phead) {
				lrec_t* prec = sllv_pop(pstate->pleft_unpaired_records);
				lrec_free(prec);
			}
			return sllv_single(NULL);
		}
	}

	slls_t* pright_field_values = mlr_reference_selected_values_from_record(pright_rec,
		pstate->popts->pright_join_field_names);
	if (pright_field_values != NULL) {
		join_bucket_t* pleft_bucket = lhmslv_get(pstate->pleft_buckets_by_join_field_values, pright_field_values);
		slls_free(pright_field_values);
		if (pleft_bucket == NULL) {
			if (pstate->popts->emit_right_unpairables) {
				return sllv_single(pright_rec);
			} else {
				lrec_free(pright_rec);
				return NULL;
			}
		} else if (pstate->popts->emit_pairables) {
			sllv_t* pout_recs = sllv_alloc();
			pleft_bucket->was_paired = TRUE;
			mapper_join_form_pairs(pleft_bucket->precords, pright_rec, pstate, pout_recs);
			lrec_free(pright_rec);
			return pout_recs;
		} else {
			pleft_bucket->was_paired = TRUE;
			lrec_free(pright_rec);
			return NULL;
		}
	} else {
		if (pstate->popts->emit_right_unpairables) {
			return sllv_single(pright_rec);
		} else {
			lrec_free(pright_rec);
			return NULL;
		}
	}
}

// ----------------------------------------------------------------
// This could be optimized in several ways:
// * Store the prefix length instead of computing its strlen inside
//   mlr_paste_2_strings
// * Don't do the if-statement on each call: prefix is null or non-null
//   at the time the mapper is constructed. Use a function pointer.
static inline char* compose_keys(char* prefix, char* key) {
	if (prefix == NULL)
		return mlr_strdup_or_die(key);
	else
		return mlr_paste_2_strings(prefix, key);
}

static void mapper_join_form_pairs(sllv_t* pleft_records, lrec_t* pright_rec, mapper_join_state_t* pstate,
	sllv_t* pout_recs)
{
	for (sllve_t* pe = pleft_records->phead; pe != NULL; pe = pe->pnext) {
		lrec_t* pleft_rec = pe->pvvalue;
		lrec_t* pout_rec = lrec_unbacked_alloc();

		// add the joined-on fields
		sllse_t* pg = pstate->popts->pleft_join_field_names->phead;
		sllse_t* ph = pstate->popts->pright_join_field_names->phead;
		sllse_t* pi = pstate->popts->poutput_join_field_names->phead;
		for ( ; pg != NULL && ph != NULL && pi != NULL; pg = pg->pnext, ph = ph->pnext, pi = pi->pnext) {
			char* v = lrec_get(pleft_rec, pg->value);
			if (v != NULL) {
				lrec_put(pout_rec, pi->value, mlr_strdup_or_die(v), FREE_ENTRY_VALUE);
			}
		}

		// add the left-record fields not already added
		for (lrece_t* pl = pleft_rec->phead; pl != NULL; pl = pl->pnext) {
			if (!hss_has(pstate->pleft_field_name_set, pl->key)) {
				lrec_put(pout_rec, compose_keys(pstate->popts->left_prefix, pl->key),
					mlr_strdup_or_die(pl->value), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			}
		}

		// add the right-record fields not already added
		for (lrece_t* pr = pright_rec->phead; pr != NULL; pr = pr->pnext) {
			if (!hss_has(pstate->pright_field_name_set, pr->key)) {
				lrec_put(pout_rec, compose_keys(pstate->popts->right_prefix, pr->key),
					mlr_strdup_or_die(pr->value), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
			}
		}

		sllv_append(pout_recs, pout_rec);
	}
}

// ----------------------------------------------------------------
static void ingest_left_file(mapper_join_state_t* pstate) {
	mapper_join_opts_t* popts = pstate->popts;

	lrec_reader_t* plrec_reader = lrec_reader_alloc(&popts->reader_opts);

	void* pvhandle = plrec_reader->popen_func(plrec_reader->pvstate, pstate->popts->prepipe,
		pstate->popts->left_file_name);
	plrec_reader->psof_func(plrec_reader->pvstate, pvhandle);

	context_t ctx = {
		.nr             = 0,
		.fnr            = 0,
		.filenum        = 1,
		.filename       = pstate->popts->left_file_name,
		.force_eof      = FALSE,
		.ips            = NULL,
		.ifs            = NULL,
		.irs            = NULL,
		.ops            = NULL,
		.ofs            = NULL,
		.ors            = NULL,
		.auto_line_term = NULL,
	};
	context_t* pctx = &ctx;

	pstate->pleft_buckets_by_join_field_values = lhmslv_alloc();

	while (TRUE) {
		lrec_t* pleft_rec = plrec_reader->pprocess_func(plrec_reader->pvstate, pvhandle, pctx);
		if (pleft_rec == NULL)
			break;
		// Mmapped input records are backed by their storage, i.e. they contain pointers into
		// mmaped file data. After the lrec reader is freed they will be invalid. So in this
		// ingestor we need to copy.
		lrec_t* pleft_copy = lrec_copy(pleft_rec);

		slls_t* pleft_field_values = mlr_reference_selected_values_from_record(pleft_copy,
			pstate->popts->pleft_join_field_names);
		if (pleft_field_values != NULL) {
			join_bucket_t* pbucket = lhmslv_get(pstate->pleft_buckets_by_join_field_values, pleft_field_values);
			if (pbucket == NULL) { // New key-field-value: new bucket and hash-map entry
				slls_t* pkey_field_values_copy = slls_copy(pleft_field_values);
				join_bucket_t* pbucket = mlr_malloc_or_die(sizeof(join_bucket_t));
				pbucket->precords = sllv_alloc();
				pbucket->was_paired = FALSE;
				pbucket->pleft_field_values = slls_copy(pleft_field_values);
				lhmslv_put(pstate->pleft_buckets_by_join_field_values, pkey_field_values_copy, pbucket,
					FREE_ENTRY_KEY);
				sllv_append(pbucket->precords, pleft_copy);
			} else { // Previously seen key-field-value: append record to bucket
				sllv_append(pbucket->precords, pleft_copy);
			}
			slls_free(pleft_field_values);
		} else {
			sllv_append(pstate->pleft_unpaired_records, pleft_copy);
		}
		lrec_free(pleft_rec);
	}

	plrec_reader->pclose_func(plrec_reader->pvstate, pvhandle, pstate->popts->prepipe);

	plrec_reader->pfree_func(plrec_reader);
}
