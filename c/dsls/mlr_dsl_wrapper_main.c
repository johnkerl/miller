#include <stdio.h>
#include <stdlib.h>
#include "mlr_dsl_wrapper.h"
#include "mlr_dsl_lexer.h"
#include "mlr_dsl_parse.h"
#include "../lib/mlrutil.h"
#include "../containers/mlr_dsl_ast.h"
#include "../containers/sllv.h"

// ----------------------------------------------------------------
int main(int argc, char** argv) {
	int trace_parse = FALSE;
	int argi = 1;
	if (argc >= 2 && streq(argv[1], "-t")) {
		argi++;
		trace_parse = TRUE;
	}
	if ((argc - argi) != 1) {
		fprintf(stderr, "Usage: %s [-t] {expression}\n", argv[0]);
		exit(1);
	}

	mlr_dsl_ast_t* past = mlr_dsl_parse(argv[argi], trace_parse);
	if (past == NULL) {
		printf("syntax error\n");
	} else {
		mlr_dsl_ast_print(past);
	}

	return 0;
}
