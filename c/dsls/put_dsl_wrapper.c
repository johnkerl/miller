#include <stdio.h>
#include <stdlib.h>
#include "put_dsl_wrapper.h"
#include "put_dsl_lexer.h"
#include "put_dsl_parse.h"
#include "../lib/mlrutil.h"
#include "../containers/mlr_dsl_ast.h"
#include "../containers/sllv.h"

void *put_dsl_lemon_parser_alloc(void *(*mallocProc)(size_t));

// ----------------------------------------------------------------
// http://flex.sourceforge.net/manual/Init-and-Destroy-Functions.html
// http://flex.sourceforge.net/manual/Extra-Data.html

// Returns linked list of mlr_dsl_ast_node_t*.
static sllv_t* put_dsl_parse_inner(yyscan_t scanner, void* pvparser, mlr_dsl_ast_node_t** ppnode) {
	int lex_code;
	int parse_code;
	sllv_t* pasts = sllv_alloc();
	do {
		lex_code = put_dsl_lexer_lex(scanner);
		mlr_dsl_ast_node_t* plexed_node = *ppnode;
		parse_code = put_dsl_lemon_parser_parse_token(pvparser, lex_code, plexed_node, pasts);
		if (parse_code == 0)
			return NULL;
	} while (lex_code > 0);
	if (-1 == lex_code) {
		fprintf(stderr, "The scanner encountered an error.\n");
		return NULL;
	}
	parse_code = put_dsl_lemon_parser_parse_token(pvparser, 0, NULL, pasts);

	if (parse_code == 0)
		return NULL;
	return pasts;
}

// ----------------------------------------------------------------
// Returns linked list of mlr_dsl_ast_node_t*.
sllv_t* put_dsl_parse(char* string) {
	mlr_dsl_ast_node_t* pnode = NULL;
	yyscan_t scanner;
	put_dsl_lexer_lex_init_extra(&pnode, &scanner);
	void* pvparser = put_dsl_lemon_parser_alloc(malloc);

	YY_BUFFER_STATE buf = NULL;
	if (string == NULL) {
		put_dsl_lexer_set_in(stdin, scanner);
	} else {
		YY_BUFFER_STATE buf = put_dsl_lexer__scan_string(string, scanner);
		put_dsl_lexer__switch_to_buffer (buf, scanner);
	}

	sllv_t* pasts = put_dsl_parse_inner(scanner, pvparser, &pnode);

	if (buf != NULL)
		put_dsl_lexer__delete_buffer(buf, scanner);

	put_dsl_lexer_lex_destroy(scanner);
	put_dsl_lemon_parser_free(pvparser, free);

	return pasts;
}

// ----------------------------------------------------------------
#ifdef __PUT_DSL_MAIN__
static int main_single(char* string) {
	sllv_t* pasts = put_dsl_parse(string);
	if (pasts == NULL || pasts->length == 0) {
		printf("put_dsl main syntax error!\n");
		return 1;
	} else {
		printf("#AST = %d\n", pasts->length);
		for (sllve_t* pe = pasts->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_print(pe->pvdata);
		}
		printf("put_dsl main parse OK\n");
		return 0;
	}
}

int main(int argc, char** argv) {
	int shellrc = 0;
	if (argc == 1) {
		printf("> ");
		shellrc |= main_single(NULL);
	} else {
		for (int argi = 1; argi < argc; argi++) {
			shellrc |= main_single(argv[argi]);
		}
	}

	return shellrc;
}
#endif // __PUT_DSL_MAIN__

// ----------------------------------------------------------------
void yytestcase(int ignored) {
}
