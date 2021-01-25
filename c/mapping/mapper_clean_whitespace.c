#include "lib/mlrutil.h"
#include "lib/mvfuncs.h"
#include "containers/sllv.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

#define RENAME_SB_ALLOC_LENGTH 16

typedef struct _mapper_clean_whitespace_state_t {
	ap_state_t* pargp;
	int      do_keys;
	int      do_values;
} mapper_clean_whitespace_state_t;

static void      mapper_clean_whitespace_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_clean_whitespace_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __);
static mapper_t* mapper_clean_whitespace_alloc(ap_state_t* pargp, int do_keys, int do_values);
static void      mapper_clean_whitespace_free(mapper_t* pmapper, context_t* _);
static sllv_t*   mapper_clean_whitespace_kvprocess(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_clean_whitespace_kprocess(lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t*   mapper_clean_whitespace_vprocess(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_clean_whitespace_setup = {
	.verb = "clean-whitespace",
	.pusage_func = mapper_clean_whitespace_usage,
	.pparse_func = mapper_clean_whitespace_parse_cli,
	.ignores_input = FALSE,
};

// ----------------------------------------------------------------
static void mapper_clean_whitespace_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "For each record, for each field in the record, whitespace-cleans the keys and\n");
	fprintf(o, "values. Whitespace-cleaning entails stripping leading and trailing whitespace,\n");
	fprintf(o, "and replacing multiple whitespace with singles. For finer-grained control,\n");
	fprintf(o, "please see the DSL functions lstrip, rstrip, strip, collapse_whitespace,\n");
	fprintf(o, "and clean_whitespace.\n");
	fprintf(o, "\n");
	fprintf(o, "Options:\n");
	fprintf(o, "-k|--keys-only    Do not touch values.\n");
	fprintf(o, "-v|--values-only  Do not touch keys.\n");
	fprintf(o, "It is an error to specify -k as well as -v -- to clean keys and values,\n");
	fprintf(o, "leave off -k as well as -v.\n");
}

static mapper_t* mapper_clean_whitespace_parse_cli(int* pargi, int argc, char** argv,
	cli_reader_opts_t* _, cli_writer_opts_t* __)
{
	int kflag = FALSE;
	int vflag = FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_true_flag(pstate, "-k", &kflag);
	ap_define_true_flag(pstate, "--keys-only", &kflag);
	ap_define_true_flag(pstate, "-v", &vflag);
	ap_define_true_flag(pstate, "--values-only", &vflag);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_clean_whitespace_usage(stderr, argv[0], verb);
		return NULL;
	}

	int do_keys = TRUE;
	int do_values = TRUE;
	if (kflag && vflag) {
		mapper_clean_whitespace_usage(stderr, argv[0], verb);
		return NULL;
	} else if (kflag) {
		do_values = FALSE;
	} else if (vflag) {
		do_keys = FALSE;
	}

	return mapper_clean_whitespace_alloc(pstate, do_keys, do_values);
}

// ----------------------------------------------------------------
static mapper_t* mapper_clean_whitespace_alloc(ap_state_t* pargp, int do_keys, int do_values) {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_clean_whitespace_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_clean_whitespace_state_t));

	pstate->pargp = pargp;
	if (do_keys && do_values) {
		pmapper->pprocess_func = mapper_clean_whitespace_kvprocess;
	} else if (do_keys) {
		pmapper->pprocess_func = mapper_clean_whitespace_kprocess;
	} else if (do_values) {
		pmapper->pprocess_func = mapper_clean_whitespace_vprocess;
	}
	pmapper->pfree_func = mapper_clean_whitespace_free;

	pmapper->pvstate = (void*)pstate;
	return pmapper;
}

static void mapper_clean_whitespace_free(mapper_t* pmapper, context_t* _) {
	mapper_clean_whitespace_state_t* pstate = pmapper->pvstate;
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_clean_whitespace_kvprocess(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		lrec_t* poutrec = lrec_unbacked_alloc();
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			mv_t old_key   = mv_from_string_with_free(mlr_strdup_or_die(pe->key));
			mv_t old_value = mv_from_string_with_free(mlr_strdup_or_die(pe->value));
			mv_t new_key   = s_s_clean_whitespace_func(&old_key);
			mv_t new_value = s_s_clean_whitespace_func(&old_value);
			char free_flags = 0;
			if (new_key.free_flags & FREE_ENTRY_VALUE)
				free_flags |= FREE_ENTRY_KEY;
			if (new_value.free_flags & FREE_ENTRY_VALUE)
				free_flags |= FREE_ENTRY_VALUE;
			lrec_put(poutrec, new_key.u.strv, new_value.u.strv, free_flags);
		}
		lrec_free(pinrec);
		return sllv_single(poutrec);
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_clean_whitespace_kprocess(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		lrec_t* poutrec = lrec_unbacked_alloc();
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			mv_t old_key = mv_from_string_with_free(mlr_strdup_or_die(pe->key));
			mv_t value   = mv_from_string_with_free(mlr_strdup_or_die(pe->value));
			mv_t new_key = s_s_clean_whitespace_func(&old_key);
			char free_flags = FREE_ENTRY_VALUE;
			if (new_key.free_flags & FREE_ENTRY_VALUE)
				free_flags |= FREE_ENTRY_KEY;
			lrec_put(poutrec, new_key.u.strv, value.u.strv, free_flags);
		}
		lrec_free(pinrec);
		return sllv_single(poutrec);
	}
	else {
		return sllv_single(NULL);
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_clean_whitespace_vprocess(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec != NULL) {
		lrec_t* poutrec = lrec_unbacked_alloc();
		for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
			mv_t key       = mv_from_string_with_free(mlr_strdup_or_die(pe->key));
			mv_t old_value = mv_from_string_with_free(mlr_strdup_or_die(pe->value));
			mv_t new_value = s_s_clean_whitespace_func(&old_value);
			char free_flags = FREE_ENTRY_KEY;
			if (new_value.free_flags & FREE_ENTRY_VALUE)
				free_flags |= FREE_ENTRY_VALUE;
			lrec_put(poutrec, key.u.strv, new_value.u.strv, free_flags);
		}
		lrec_free(pinrec);
		return sllv_single(poutrec);
	}
	else {
		return sllv_single(NULL);
	}
}
