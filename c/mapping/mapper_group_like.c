#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/slls.h"
#include "containers/lhmslv.h"
#include "containers/mixutil.h"
#include "mapping/mappers.h"

typedef struct _mapper_group_like_state_t {
	// map from list of string to list of record
	lhmslv_t* precords_by_key_field_names;
} mapper_group_like_state_t;

static void      mapper_group_like_usage(FILE* o, char* argv0, char* verb);
static mapper_t* mapper_group_like_parse_cli(int* pargi, int argc, char** argv);
static mapper_t* mapper_group_like_alloc();
static void      mapper_group_like_free(mapper_t* pmapper);
static sllv_t*   mapper_group_like_process(lrec_t* pinrec, context_t* pctx, void* pvstate);

// ----------------------------------------------------------------
mapper_setup_t mapper_group_like_setup = {
	.verb = "group-like",
	.pusage_func = mapper_group_like_usage,
	.pparse_func = mapper_group_like_parse_cli,
};

// ----------------------------------------------------------------
static void mapper_group_like_usage(FILE* o, char* argv0, char* verb) {
	fprintf(o, "Usage: %s %s\n", argv0, verb);
	fprintf(o, "Outputs records in batches having identical field names.\n");
}

static mapper_t* mapper_group_like_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 1) {
		mapper_group_like_usage(stderr, argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_group_like_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
static mapper_t* mapper_group_like_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_group_like_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_group_like_state_t));
	pstate->precords_by_key_field_names = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_group_like_process;
	pmapper->pfree_func    = mapper_group_like_free;

	return pmapper;
}

static void mapper_group_like_free(mapper_t* pmapper) {
	mapper_group_like_state_t* pstate = pmapper->pvstate;
	// lhmslv_free will free the hashmap keys; we need to free the void-star hashmap values.
	for (lhmslve_t* pa = pstate->precords_by_key_field_names->phead; pa != NULL; pa = pa->pnext) {
		sllv_t* plist = pa->pvvalue;
		sllv_free(plist);
	}
	lhmslv_free(pstate->precords_by_key_field_names);
	free(pstate);
	free(pmapper);
}

// ----------------------------------------------------------------
static sllv_t* mapper_group_like_process(lrec_t* pinrec, context_t* pctx, void* pvstate) {
	mapper_group_like_state_t* pstate = pvstate;
	if (pinrec != NULL) {
		slls_t* pkey_field_names = mlr_reference_keys_from_record(pinrec);
		sllv_t* plist = lhmslv_get(pstate->precords_by_key_field_names, pkey_field_names);
		if (plist == NULL) {
			plist = sllv_alloc();
			sllv_add(plist, pinrec);
			lhmslv_put(pstate->precords_by_key_field_names, slls_copy(pkey_field_names), plist);
		} else {
			sllv_add(plist, pinrec);
		}
		slls_free(pkey_field_names);
		return NULL;
	} else {
		sllv_t* poutput = sllv_alloc();
		for (lhmslve_t* pe = pstate->precords_by_key_field_names->phead; pe != NULL; pe = pe->pnext) {
			sllv_t* plist = pe->pvvalue;
			sllv_transfer(poutput, plist);
		}
		sllv_add(poutput, NULL);
		return poutput;
	}
}
