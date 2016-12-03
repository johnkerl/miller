#include <stdio.h>
#include <stdlib.h>
#include "ex2_wrapper.h"
#include "ex2_lexer.h"
#include "ex2_parse.h"
#include "../lib/mlrutil.h"
#include "./ex_ast.h"
#include "../containers/sllv.h"

// These prototypes are copied out manually from ex2_parse.c. With some
// more work I could have Lemon autogenerate these prototypes into
// ex2_parse.h.

void *ex2_lemon_parser_alloc(void *(*mallocProc)(size_t));

int ex2_lemon_parser_parse_token(
	void *pvparser,              /* The parser */
	int yymajor,                 /* The major token code number */
	ex_ast_node_t* yyminor, /* The value for the token */
	ex_ast_t* past);        /* Optional %extra_argument parameter */
void ex2_lemon_parser_free(
	void *pvparser,             /* The parser to be deleted */
	void (*freeProc)(void*));   /* Function used to reclaim memory */

void ex2_ParseTrace(FILE *TraceFILE, char *zTracePrompt);

// ----------------------------------------------------------------
// http://flex.sourceforge.net/manual/Init-and-Destroy-Functions.html
// http://flex.sourceforge.net/manual/Extra-Data.html

// Returns linked list of ex_ast_node_t*.
static ex_ast_t* ex2_parse_inner(yyscan_t scanner, void* pvparser, ex_ast_node_t** ppnode,
	int trace_parse)
{
	int lex_code;
	int parse_code;
	ex_ast_t* past = ex_ast_alloc();
	if (trace_parse)
		ex2_ParseTrace(stderr, "[DSLTRACE] ");
	do {
		lex_code = ex2_lexer_lex(scanner);
		ex_ast_node_t* plexed_node = *ppnode;
		parse_code = ex2_lemon_parser_parse_token(pvparser, lex_code, plexed_node, past);
		if (parse_code == 0)
			return NULL;
	} while (lex_code > 0);
	if (-1 == lex_code) {
		fprintf(stderr, "The scanner encountered an error.\n");
		return NULL;
	}
	parse_code = ex2_lemon_parser_parse_token(pvparser, 0, NULL, past);

	if (parse_code == 0)
		return NULL;
	return past;
}

// ----------------------------------------------------------------
// Returns linked list of ex_ast_node_t*.
ex_ast_t* ex2_parse(char* string, int trace_parse) {
	ex_ast_node_t* pnode = NULL;
	yyscan_t scanner;
	ex2_lexer_lex_init_extra(&pnode, &scanner);
	void* pvparser = ex2_lemon_parser_alloc(malloc);

	YY_BUFFER_STATE buf = NULL;
	if (string == NULL) {
		ex2_lexer_set_in(stdin, scanner);
	} else {
		YY_BUFFER_STATE buf = ex2_lexer__scan_string(string, scanner);
		ex2_lexer__switch_to_buffer (buf, scanner);
	}

	ex_ast_t* past = ex2_parse_inner(scanner, pvparser, &pnode, trace_parse);

	if (buf != NULL)
		ex2_lexer__delete_buffer(buf, scanner);

	ex2_lexer_lex_destroy(scanner);
	ex2_lemon_parser_free(pvparser, free);

	return past;
}

// ----------------------------------------------------------------
void yytestcase(int ignored) {
}

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

	ex_ast_t* past = ex2_parse(argv[argi], trace_parse);
	if (past == NULL) {
		printf("syntax error\n");
	} else {
		ex_ast_print(past);
	}

	return 0;
}
