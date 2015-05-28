#include <stdio.h>
#include <stdlib.h>
#include "filter_dsl_wrapper.h"
#include "filter_dsl_lexer.h"
#include "filter_dsl_parse.h"
#include "../lib/mlrutil.h"
#include "../containers/mlr_dsl_ast.h"

void *filter_dsl_lemon_parser_alloc(void *(*mallocProc)(size_t));

// ----------------------------------------------------------------
// http://flex.sourceforge.net/manual/Init-and-Destroy-Functions.html
// http://flex.sourceforge.net/manual/Extra-Data.html

static mlr_dsl_ast_node_holder_t* filter_dsl_parse_inner(yyscan_t scanner, void* pvparser,
	mlr_dsl_ast_node_t** ppnode)
{
	int lex_code;
	int parse_code;
	mlr_dsl_ast_node_holder_t* past = mlr_malloc_or_die(sizeof(mlr_dsl_ast_node_holder_t));
	past->proot = NULL;
	do {
		lex_code = filter_dsl_lexer_lex(scanner);
		mlr_dsl_ast_node_t* plexed_node = *ppnode;
		parse_code = filter_dsl_lemon_parser_parse_token(pvparser, lex_code, plexed_node, past);
		if (parse_code == 0)
			return NULL;
	} while (lex_code > 0);
	if (-1 == lex_code) {
		fprintf(stderr, "The scanner encountered an error.\n");
		return NULL;
	}
	parse_code = filter_dsl_lemon_parser_parse_token(pvparser, 0, NULL, past);

	if (parse_code == 0)
		return NULL;
	return past;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_holder_t* filter_dsl_parse(char* string) {
	mlr_dsl_ast_node_t* pnode = NULL;
	yyscan_t scanner;
	filter_dsl_lexer_lex_init_extra(&pnode, &scanner);
	void* pvparser = filter_dsl_lemon_parser_alloc(malloc);

	YY_BUFFER_STATE buf = NULL;
	if (string == NULL) {
		filter_dsl_lexer_set_in(stdin, scanner);
	} else {
		YY_BUFFER_STATE buf = filter_dsl_lexer__scan_string(string, scanner);
		filter_dsl_lexer__switch_to_buffer (buf, scanner);
	}

	mlr_dsl_ast_node_holder_t* past = filter_dsl_parse_inner(scanner, pvparser, &pnode);

	if (buf != NULL)
		filter_dsl_lexer__delete_buffer(buf, scanner);

	filter_dsl_lexer_lex_destroy(scanner);
	filter_dsl_lemon_parser_free(pvparser, free);

	return past;
}

// ----------------------------------------------------------------
#ifdef __FILTER_DSL_MAIN__
static int main_single(char* string) {
	mlr_dsl_ast_node_holder_t* past = filter_dsl_parse(string);
	if (past == NULL) {
		printf("filter_dsl main syntax error!\n");
		return 1;
	} else {
		mlr_dsl_ast_node_print(past->proot);
		printf("filter_dsl main parse OK\n");
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
#endif // __FILTER_DSL_MAIN__
