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
		return NULL;
	}
	else {
		sllv_t* poutput = sllv_alloc();
		for (lhmslve_t* pe = pstate->precords_by_key_field_names->phead; pe != NULL; pe = pe->pnext) {
			sllv_t* plist = pe->pvvalue;
			for (sllve_t* pf = plist->phead; pf != NULL; pf = pf->pnext) {
				sllv_add(poutput, pf->pvdata);
			}
		}
		sllv_add(poutput, NULL);
		return poutput;
	}
}

// ----------------------------------------------------------------
static void mapper_group_like_free(void* pvstate) {
	mapper_group_like_state_t* pstate = (mapper_group_like_state_t*)pvstate;
	if (pstate->precords_by_key_field_names != NULL)
		// xxx check for full recursive free
		// xxx in lhmslv & more general outermost readme, articulate the philosophy that containers
		// will free contents except for void-stars.
		lhmslv_free(pstate->precords_by_key_field_names);
}

static mapper_t* mapper_group_like_alloc() {
	mapper_t* pmapper = mlr_malloc_or_die(sizeof(mapper_t));

	mapper_group_like_state_t* pstate = mlr_malloc_or_die(sizeof(mapper_group_like_state_t));
	pstate->precords_by_key_field_names = lhmslv_alloc();

	pmapper->pvstate       = pstate;
	pmapper->pprocess_func = mapper_group_like_process;
	pmapper->pfree_func    = mapper_group_like_free;

	return pmapper;
}

// ----------------------------------------------------------------
static void mapper_group_like_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s\n", argv0, verb);
	fprintf(stdout, "Outputs records in batches having identical field names.\n");
}
static mapper_t* mapper_group_like_parse_cli(int* pargi, int argc, char** argv) {
	if ((argc - *pargi) < 1) {
		mapper_group_like_usage(argv[0], argv[*pargi]);
		return NULL;
	}
	mapper_t* pmapper = mapper_group_like_alloc();
	*pargi += 1;
	return pmapper;
}

// ----------------------------------------------------------------
mapper_setup_t mapper_group_like_setup = {
	.verb = "group-like",
	.pusage_func = mapper_group_like_usage,
	.pparse_func = mapper_group_like_parse_cli,
};
