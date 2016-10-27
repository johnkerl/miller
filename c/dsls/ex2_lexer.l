%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "ex2_parse.h"
#include "./ex_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html

%}

%option reentrant
%option noyywrap
%option extra-type="struct _ex_ast_node_t **"

%%

"=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_SREC_ASSIGNMENT);
	return MD_TOKEN_ASSIGN;
}

\$[a-zA-Z_0-9]+ {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FIELD_NAME);
	return MD_TOKEN_FIELD_NAME;
}

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
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	return MD_TOKEN_NUMBER;
}
0x[0-9a-fA-F]+ {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	return MD_TOKEN_NUMBER;
}

[ \t\r\n] { }

. {
	return -1;
}
%%
