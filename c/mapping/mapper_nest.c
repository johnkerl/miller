#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "lib/mlrregex.h"
#include "lib/string_builder.h"
#include "containers/lhmss.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ================================================================
#define SB_ALLOC_LENGTH 128

typedef struct _mapper_nest_state_t {
	ap_state_t* pargp;

	char* field_name;
	char* nested_fs;
	char* nested_ps;
	int   nested_ps_length;

	lhmslv_t* other_keys_to_other_values_to_buckets;
	string_builder_t* psb;
	regex_t regex;
} mapper_nest_state_t;

typedef struct _nest_bucket_t {
	lrec_t*  prepresentative;
	sllv_t* pairs;
} nest_bucket_t;

static void      mapper_nest_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_nest_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_nest_alloc(ap_state_t* pargp,
	char* field_name, char* nested_fs, char* nested_ps,
	int do_explode, int do_pairs, int do_across_fields);
static void    mapper_nest_free(mapper_t* pmapper);

static sllv_t* mapper_nest_explode_pairs_across_fields   (lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t* mapper_nest_explode_pairs_across_records  (lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t* mapper_nest_explode_values_across_fields  (lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t* mapper_nest_explode_values_across_records (lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t* mapper_nest_implode_values_across_fields  (lrec_t* pinrec, context_t* pctx, void* pvstate);
static sllv_t* mapper_nest_implode_values_across_records (lrec_t* pinrec, context_t* pctx, void* pvstate);

static nest_bucket_t* nest_bucket_alloc(lrec_t* prepresentative);
static void nest_bucket_free(nest_bucket_t* pbucket);

// ----------------------------------------------------------------
mapper_setup_t mapper_nest_setup = {
	.verb = "nest",
	.pusage_func = mapper_nest_usage,
	.pparse_func = mapper_nest_parse_cli
};

// ----------------------------------------------------------------
static void mapper_nest_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(o, "Explodes specified field values into separate fields/records, or reverses this.\n");
	fprintf(o, "Options:\n");
	fprintf(o, "  --explode,--implode   One is required.\n");
	fprintf(o, "  --values,--pairs      One is required.\n");
	fprintf(o, "  --across-records,--across-fields One is required.\n");
	fprintf(o, "  -f {field name}       Required.\n");
	fprintf(o, "  --nested-fs {string}  Defaults to \";\". Field separator for nested values.\n");
	fprintf(o, "  --nested-ps {string}  Defaults to \":\". Pair separator for nested key-value pairs.\n");

	fprintf(o, "\n");
	fprintf(o, "Examples:\n");

	fprintf(o, "\n");
	fprintf(o, "  %s %s --explode --values --across-records -f x\n", argv0, verb);
	fprintf(o, "  with input record \"x=a;b;c,y=d\" produces output records\n");
	fprintf(o, "    \"x=a,y=d\"\n");
	fprintf(o, "    \"x=b,y=d\"\n");
	fprintf(o, "    \"x=c,y=d\"\n");
	fprintf(o, "  Use --implode to do the reverse.\n");

	fprintf(o, "\n");
	fprintf(o, "  %s %s --explode --values --across-fields -f x\n", argv0, verb);
	fprintf(o, "  with input record \"x=a;b;c,y=d\" produces output records\n");
	fprintf(o, "    \"y=d,x_1=a,x_2=b,x_3=c\"\n");
	fprintf(o, "  Use --implode to do the reverse.\n");

	fprintf(o, "\n");
	fprintf(o, "  %s %s --explode --pairs --across-records -f x\n", argv0, verb);
	fprintf(o, "  with input record \"x=a:1;b:2;c:3,y=d\" produces output records\n");
	fprintf(o, "    \"y=d,a=1\"\n");
	fprintf(o, "    \"y=d,b=2\"\n");
	fprintf(o, "    \"y=d,c=3\"\n");

	fprintf(o, "\n");
	fprintf(o, "  %s %s --explode --pairs --across-fields -f x\n", argv0, verb);
	fprintf(o, "  with input record \"x=a:1;b:2;c:3,y=d\" produces output records\n");
	fprintf(o, "    \"y=d,a=a,b=2,c=3\"\n");
	fprintf(o, "\n");
	fprintf(o, "Notes:\n");
	fprintf(o, "* With --pairs, --implode doesn't make sense since the original field name has\n");
	fprintf(o, "  been lost.\n");
	fprintf(o, "* The combination \"--implode --values --across-records\" is non-streaming:\n");
	fprintf(o, "  no output records are produced until all input records have been read. In\n");
	fprintf(o, "  particular, this means it won't work in tail -f contexts. But all other flag\n");
	fprintf(o, "  combinations result in streaming (tail -f friendly) data processing.\n");
	fprintf(o, "* It's up to you to ensure that the nested-fs is distinct from your data's IFS:\n");
	fprintf(o, "  e.g. by default the former is semicolon and the latter is comma.\n");
}

static mapper_t* mapper_nest_parse_cli(int* pargi, int argc, char** argv) {
	char* field_name = NULL;
	char* nested_fs = ";";
	char* nested_ps = ":";
	int   do_explode       = NEITHER_TRUE_NOR_FALSE;
	int   do_pairs         = NEITHER_TRUE_NOR_FALSE;
	int   do_across_fields = NEITHER_TRUE_NOR_FALSE;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_flag(pstate, "-f",               &field_name);
	ap_define_string_flag(pstate, "--nested-fs",      &nested_fs);
	ap_define_string_flag(pstate, "--nested-ps",      &nested_ps);
	ap_define_true_flag(pstate,   "--explode",        &do_explode);
	ap_define_false_flag(pstate,  "--implode",        &do_explode);
	ap_define_true_flag(pstate,   "--pairs",          &do_pairs);
	ap_define_false_flag(pstate,  "--values",         &do_pairs);
	ap_define_true_flag(pstate,   "--across-fields",  &do_across_fields);
	ap_define_false_flag(pstate,  "--across-records", &do_across_fields);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}

	if (field_name == NULL) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (do_explode == NEITHER_TRUE_NOR_FALSE) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (do_pairs == NEITHER_TRUE_NOR_FALSE) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (do_across_fields == NEITHER_TRUE_NOR_FALSE) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}
	if (do_pairs == TRUE && do_explode == FALSE) {
		mapper_nest_usage(stderr, argv[0], verb);
		return NULL;
	}

	return mapper_nest_alloc(pstate, field_name, nested_fs, nested_ps, do_explode, do_pairs, do_across_fields);
}

// ----------------------------------------------------------------
static mapper_t* mapper_nest_alloc(ap_state_t* pargp,
	char* field_name, char* nested_fs, char* nested_ps,
	int do_explode, int do_pairs, int do_across_fields)
{
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_nest_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_nest_state_t));

	pstate->pargp      = pargp;
	pstate->field_name = field_name;
	pstate->nested_fs  = mlr_unbackslash(nested_fs);
	pstate->nested_ps  = mlr_unbackslash(nested_ps);
	pstate->nested_ps_length = strlen(pstate->nested_ps);

	if (do_explode) {
		if (do_pairs) {
			pmapper->pprocess_func = do_across_fields
				? mapper_nest_explode_pairs_across_fields
				: mapper_nest_explode_pairs_across_records;
		} else {
			pmapper->pprocess_func = do_across_fields
				? mapper_nest_explode_values_across_fields
				: mapper_nest_explode_values_across_records;
		}
	} else {
		if (do_pairs) {
			// Should have been caught in CLI-parser.
			fprintf(stderr, "%s: internal coding error detected at file %s line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		} else {
			pmapper->pprocess_func = do_across_fields
				? mapper_nest_implode_values_across_fields
				: mapper_nest_implode_values_across_records;
		}
	}
	pstate->other_keys_to_other_values_to_buckets = lhmslv_alloc();
	pstate->psb = sb_alloc(SB_ALLOC_LENGTH);
	char* pattern = mlr_malloc_or_die(strlen(field_name) + 12);
	sprintf(pattern, "^%s_[0-9]+$", field_name);
	regcomp_or_die(&pstate->regex, pattern, REG_NOSUB);
	free(pattern);

	pmapper->pfree_func = mapper_nest_free;

	pmapper->pvstate = (void*)pstate;
	return pmapper;
}

static void mapper_nest_free(mapper_t* pmapper) {
	mapper_nest_state_t* pstate = pmapper->pvstate;

	if (pstate->other_keys_to_other_values_to_buckets != NULL) {
		for (lhmslve_t* pe = pstate->other_keys_to_other_values_to_buckets->phead; pe != NULL; pe = pe->pnext) {
			lhmslv_t* other_values_to_buckets = pe->pvvalue;
			for (lhmslve_t* pf = other_values_to_buckets->phead; pf != NULL; pf = pf->pnext) {
				nest_bucket_t* pbucket = pf->pvvalue;
				nest_bucket_free(pbucket);
			}
			lhmslv_free(other_values_to_buckets);
		}
		lhmslv_free(pstate->other_keys_to_other_values_to_buckets);
	}

	free(pstate->nested_fs);
	free(pstate->nested_ps);
	regfree(&pstate->regex);
	ap_free(pstate->pargp);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_explode_values_across_records(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	char* field_value = lrec_get(pinrec, pstate->field_name);
	if (field_value == NULL) {
		return sllv_single(pinrec);
	}

	sllv_t* poutrecs = sllv_alloc();
	char* sep = pstate->nested_fs;
	int i = 1;
	for (char* piece = strtok(field_value, sep); piece != NULL; piece = strtok(NULL, sep), i++) {
		lrec_t* poutrec = lrec_copy(pinrec);
		lrec_put(poutrec, pstate->field_name, mlr_strdup_or_die(piece), FREE_ENTRY_VALUE);
		sllv_append(poutrecs, poutrec);
	}
	lrec_free(pinrec);
	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_implode_values_across_records(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	if (pinrec != NULL) { // Not end of input stream
		char* pfree_flags = NULL;
		char free_flags = 0;
		char* field_value = lrec_get_ext(pinrec, pstate->field_name, &pfree_flags);
		if (field_value == NULL)
			return sllv_single(pinrec);

		// Retain the field_value, and responsibility for freeing it; then, remove it from the input record.
		free_flags = *pfree_flags;
		*pfree_flags &= ~FREE_ENTRY_VALUE;
		lrec_remove(pinrec, pstate->field_name);

		slls_t* other_keys = mlr_reference_keys_from_record(pinrec);
		lhmslv_t* other_values_to_buckets = lhmslv_get(pstate->other_keys_to_other_values_to_buckets, other_keys);
		if (other_values_to_buckets == NULL) {
			other_values_to_buckets = lhmslv_alloc();
			lhmslv_put(pstate->other_keys_to_other_values_to_buckets,
				slls_copy(other_keys), other_values_to_buckets, FREE_ENTRY_KEY);
		}

		slls_t* other_values = mlr_reference_values_from_record(pinrec);
		nest_bucket_t* pbucket = lhmslv_get(other_values_to_buckets, other_values);
		if (pbucket == NULL) {
			pbucket = nest_bucket_alloc(pinrec);
			lhmslv_put(other_values_to_buckets, slls_copy(other_values), pbucket, FREE_ENTRY_KEY);
		} else {
			lrec_free(pinrec);
		}
		lrec_t* pair = lrec_unbacked_alloc();
		lrec_put(pair, pstate->field_name, mlr_strdup_or_die(field_value), FREE_ENTRY_VALUE);
		sllv_append(pbucket->pairs, pair);

		if (free_flags & FREE_ENTRY_VALUE)
			free(field_value);
		slls_free(other_values);
		slls_free(other_keys);

		return NULL;

	} else { // end of input stream
		sllv_t* poutrecs = sllv_alloc();

		for (lhmslve_t* pe = pstate->other_keys_to_other_values_to_buckets->phead; pe != NULL; pe = pe->pnext) {
			lhmslv_t* other_values_to_buckets = pe->pvvalue;
			for (lhmslve_t* pf = other_values_to_buckets->phead; pf != NULL; pf = pf->pnext) {
				nest_bucket_t* pbucket = pf->pvvalue;
				lrec_t* poutrec = pbucket->prepresentative;
				pbucket->prepresentative = NULL; // ownership transfer
				for (sllve_t* pg = pbucket->pairs->phead; pg != NULL; pg = pg->pnext) {
					lrec_t* pr = pg->pvvalue;
					sb_append_string(pstate->psb, pr->phead->value);
					if (pg->pnext != NULL)
						sb_append_string(pstate->psb, pstate->nested_fs);
				}
				lrec_put(poutrec, pstate->field_name, sb_finish(pstate->psb), FREE_ENTRY_VALUE);
				sllv_append(poutrecs, poutrec);
			}
		}

		sllv_append(poutrecs, NULL);
		return poutrecs;
	}
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_explode_pairs_across_records(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	char* pfree_flags = NULL;
	char free_flags = 0;
	char* field_value = lrec_get_ext(pinrec, pstate->field_name, &pfree_flags);
	if (field_value == NULL) {
		return sllv_single(pinrec);
	}

	// Retain the field_value, and responsibility for freeing it; then, remove it from the input record.
	free_flags = *pfree_flags;
	*pfree_flags &= ~FREE_ENTRY_VALUE;
	lrec_remove(pinrec, pstate->field_name);

	sllv_t* poutrecs = sllv_alloc();
	char* sep = pstate->nested_fs;
	for (char* piece = strtok(field_value, sep); piece != NULL; piece = strtok(NULL, sep)) {
		lrec_t* poutrec = lrec_copy(pinrec);
		char* found_sep = strstr(piece, pstate->nested_ps);
		if (found_sep != NULL) { // there is a pair
			*found_sep = 0;
			lrec_put(poutrec, mlr_strdup_or_die(piece), mlr_strdup_or_die(found_sep + pstate->nested_ps_length),
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		} else { // there is not a pair
			lrec_put(poutrec, pstate->field_name, mlr_strdup_or_die(piece), FREE_ENTRY_VALUE);
		}
		sllv_append(poutrecs, poutrec);
	}

	if (free_flags & FREE_ENTRY_VALUE)
		free(field_value);
	lrec_free(pinrec);
	return poutrecs;
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_explode_values_across_fields(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	char* field_value = lrec_get(pinrec, pstate->field_name);
	if (field_value == NULL) {
		return sllv_single(pinrec);
	}

	char* sep = pstate->nested_fs;
	int i = 1;
	for (char* piece = strtok(field_value, sep); piece != NULL; piece = strtok(NULL, sep), i++) {
		char  istring_free_flags;
		char* istring = make_nidx_key(i, &istring_free_flags);
		char* new_key = mlr_paste_3_strings(pstate->field_name, "_", istring);
		if (istring_free_flags & FREE_ENTRY_KEY)
			free(istring);
		lrec_put(pinrec, new_key, mlr_strdup_or_die(piece), FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
	}
	lrec_remove(pinrec, pstate->field_name);
	return sllv_single(pinrec);;
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_implode_values_across_fields(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	if (pinrec == NULL) // end of input stream
		return sllv_single(NULL);

	int field_count = 0;
	for (lrece_t* pe = pinrec->phead; pe != NULL; /* increment in loop */) {
		if (regmatch_or_die(&pstate->regex, pe->key, 0, NULL)) {
			if (field_count > 0)
				sb_append_string(pstate->psb, pstate->nested_fs);
			sb_append_string(pstate->psb, pe->value);
			field_count++;

			lrece_t* pnext = pe->pnext;
			lrec_unlink_and_free(pinrec, pe);
			pe = pnext;

		} else {
			pe = pe->pnext;
		}
	}

	if (field_count > 0)
		lrec_put(pinrec, pstate->field_name, sb_finish(pstate->psb), FREE_ENTRY_VALUE);

	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static sllv_t* mapper_nest_explode_pairs_across_fields(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	if (pinrec == NULL) // End of input stream
		return sllv_single(NULL);
	mapper_nest_state_t* pstate = (mapper_nest_state_t*)pvstate;

	char* pfree_flags = NULL;
	char free_flags = 0;
	char* field_value = lrec_get_ext(pinrec, pstate->field_name, &pfree_flags);
	if (field_value == NULL) {
		return sllv_single(pinrec);
	}

	// Retain the field_value, and responsibility for freeing it; then, remove it from the input record.
	free_flags = *pfree_flags;
	*pfree_flags &= ~FREE_ENTRY_VALUE;
	lrec_remove(pinrec, pstate->field_name);

	char* sep = pstate->nested_fs;
	for (char* piece = strtok(field_value, sep); piece != NULL; piece = strtok(NULL, sep)) {
		char* found_sep = strstr(piece, pstate->nested_ps);
		if (found_sep != NULL) { // there is a pair
			*found_sep = 0;
			lrec_put(pinrec, mlr_strdup_or_die(piece), mlr_strdup_or_die(found_sep + pstate->nested_ps_length),
				FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
		} else { // there is not a pair
			lrec_put(pinrec, pstate->field_name, mlr_strdup_or_die(piece), FREE_ENTRY_VALUE);
		}
	}

	if (free_flags & FREE_ENTRY_VALUE)
		free(field_value);
	return sllv_single(pinrec);
}

// ----------------------------------------------------------------
static nest_bucket_t* nest_bucket_alloc(lrec_t* prepresentative) {
	nest_bucket_t* pbucket = mlr_malloc_or_die(sizeof(nest_bucket_t));
	pbucket->prepresentative = prepresentative;
	pbucket->pairs = sllv_alloc();
	return pbucket;
}
static void nest_bucket_free(nest_bucket_t* pbucket) {
	lrec_free(pbucket->prepresentative);
	for (sllve_t* pe = pbucket->pairs->phead; pe != NULL; pe = pe->pnext) {
		lrec_t* pair = pe->pvvalue;
		lrec_free(pair);
	}
	sllv_free(pbucket->pairs);
	free(pbucket);
}
