%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "ex0_parse.h"
#include "./ex_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html

%}

%option reentrant
%option noyywrap
%option extra-type="struct _ex_ast_node_t **"

%%

[0-9]+|[0-9]+\.[0-9]*|[0-9]*\.[0-9]+|[0-9]+[eE][0-9]+|[0-9]+[eE]-[0-9]+|[0-9]+\.[0-9]*[eE][0-9]+|[0-9]+\.[0-9]*[eE]-[0-9]+|[0-9]*\.[0-9]+[eE][0-9]+|[0-9]*\.[0-9]+[eE]-[0-9]+ {
	// Leading minus sign is handled via the unary-minus operator, not here.
	// 123
	// 123. 123.4
	// .234
	// 1e2
	// 1e-2
	// 1.2e3 1.e3
	// 1.2e-3 1.e-3
	// .2e3
	// .2e-3 1.e-3
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRNUM_LITERAL);
	return MD_TOKEN_NUMBER;
}

"+" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_PLUS;
}
"-" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MINUS;
}
"*" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_TIMES;
}
"/" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DIVIDE;
}

[ \t\r\n] { }

. {
	return -1;
}
%%
