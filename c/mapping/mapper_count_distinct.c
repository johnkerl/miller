#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "containers/sllv.h"
#include "containers/lhmslv.h"
#include "containers/lhmsv.h"
#include "mapping/mappers.h"
#include "cli/argparse.h"

// ----------------------------------------------------------------
static void mapper_count_distinct_usage(char* argv0, char* verb) {
	fprintf(stdout, "Usage: %s %s [options]\n", argv0, verb);
	fprintf(stdout, "-f {a,b,c}   Field names for distinct count.\n");
}

// ----------------------------------------------------------------
static mapper_t* mapper_count_distinct_parse_cli(int* pargi, int argc, char** argv) {
	slls_t* pfield_names  = NULL;

	char* verb = argv[(*pargi)++];

	ap_state_t* pstate = ap_alloc();
	ap_define_string_list_flag(pstate, "-f", &pfield_names);

	if (!ap_parse(pstate, verb, pargi, argc, argv)) {
		mapper_count_distinct_usage(argv[0], verb);
		return NULL;
	}

	if (pfield_names == NULL) {
		mapper_count_distinct_usage(argv[0], verb);
		return NULL;
	}

	return mapper_uniq_alloc(pfield_names, TRUE);
}

// ----------------------------------------------------------------
mapper_setup_t mapper_count_distinct_setup = {
	.verb = "count-distinct",
	.pusage_func = mapper_count_distinct_usage,
	.pparse_func = mapper_count_distinct_parse_cli,
};
