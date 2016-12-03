#include <stdio.h>
#include <stdlib.h>
#include "mlr_dsl_wrapper.h"
#include "mlr_dsl_lexer.h"
#include "mlr_dsl_parse.h"
#include "../lib/mlrutil.h"
#include "../dsl/mlr_dsl_ast.h"
#include "../containers/sllv.h"

// These prototypes are copied out manually from mlr_dsl_parse.c. With some
// more work I could have Lemon autogenerate these prototypes into
// mlr_dsl_parse.h.

void *mlr_dsl_lemon_parser_alloc(void *(*mallocProc)(size_t));

int mlr_dsl_lemon_parser_parse_token(
	void *pvparser,              /* The parser */
	int yymajor,                 /* The major token code number */
	mlr_dsl_ast_node_t* yyminor, /* The value for the token */
	mlr_dsl_ast_t* past);        /* Optional %extra_argument parameter */
void mlr_dsl_lemon_parser_free(
	void *pvparser,             /* The parser to be deleted */
	void (*freeProc)(void*));   /* Function used to reclaim memory */

void mlr_dsl_ParseTrace(FILE *TraceFILE, char *zTracePrompt);

// ----------------------------------------------------------------
// http://flex.sourceforge.net/manual/Init-and-Destroy-Functions.html
// http://flex.sourceforge.net/manual/Extra-Data.html

// Returns linked list of mlr_dsl_ast_node_t*.
static mlr_dsl_ast_t* mlr_dsl_parse_inner(yyscan_t scanner, void* pvparser, mlr_dsl_ast_node_t** ppnode,
	int trace_parse)
{
	int lex_code;
	int parse_code;
	mlr_dsl_ast_t* past = mlr_dsl_ast_alloc();
	if (trace_parse)
		mlr_dsl_ParseTrace(stderr, "[DSLTRACE] ");
	do {
		lex_code = mlr_dsl_lexer_lex(scanner);
		mlr_dsl_ast_node_t* plexed_node = *ppnode;
		parse_code = mlr_dsl_lemon_parser_parse_token(pvparser, lex_code, plexed_node, past);
		if (parse_code == 0) {
			//mlr_dsl_ast_node_print(plexed_node);
			return NULL;
		}
	} while (lex_code > 0);
	if (-1 == lex_code) {
		fprintf(stderr, "The scanner encountered an error.\n");
		return NULL;
	}
	parse_code = mlr_dsl_lemon_parser_parse_token(pvparser, 0, NULL, past);

	if (parse_code == 0)
		return NULL;
	return past;
}

// ----------------------------------------------------------------
// Returns linked list of mlr_dsl_ast_node_t*.
mlr_dsl_ast_t* mlr_dsl_parse(char* string, int trace_parse) {
	mlr_dsl_ast_node_t* pnode = NULL;
	yyscan_t scanner;
	mlr_dsl_lexer_lex_init_extra(&pnode, &scanner);
	void* pvparser = mlr_dsl_lemon_parser_alloc(malloc);

	YY_BUFFER_STATE buf = NULL;
	if (string == NULL) {
		mlr_dsl_lexer_set_in(stdin, scanner);
	} else {
		YY_BUFFER_STATE buf = mlr_dsl_lexer__scan_string(string, scanner);
		mlr_dsl_lexer__switch_to_buffer (buf, scanner);
	}

	mlr_dsl_ast_t* past = mlr_dsl_parse_inner(scanner, pvparser, &pnode, trace_parse);

	if (buf != NULL)
		mlr_dsl_lexer__delete_buffer(buf, scanner);

	mlr_dsl_lexer_lex_destroy(scanner);
	mlr_dsl_lemon_parser_free(pvparser, free);

	return past;
}

// ----------------------------------------------------------------
void yytestcase(int ignored) {
}
