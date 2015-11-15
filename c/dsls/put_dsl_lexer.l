%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "put_dsl_parse.h"
#include "../containers/mlr_dsl_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html
%}

%option reentrant
%option noyywrap
%option extra-type="struct _mlr_dsl_ast_node_t **"

%%
";" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return PUT_DSL_SEMICOLON;
}

"=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_ASSIGN;
}

"||" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LOGICAL_OR;
}
"^^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LOGICAL_XOR;
}
"&&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LOGICAL_AND;
}

"=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_MATCHES;
}
"!=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_DOES_NOT_MATCH;
}
"==" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_EQ;
}
"!=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_NE;
}

">" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_GT;
}
">=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_GE;
}
"<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LT;
}
"<=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LE;
}

"|" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_OR;
}

"^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_XOR;
}

"&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_AND;
}

"<<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_LSH;
}
">>" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_RSH;
}

"+" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_PLUS;
}
"-" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_MINUS;
}
"\." {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_DOT;
}

"*" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_TIMES;
}
"/" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_DIVIDE;
}
"//" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_INT_DIVIDE;
}
"%" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_MOD;
}

"!" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_LOGICAL_NOT;
}

"~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_BITWISE_NOT;
}

"**" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return PUT_DSL_POW;
}

\$[a-zA-Z_0-9]+ {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	return PUT_DSL_FIELD_NAME;
}

[0-9]+|[0-9]+\.[0-9]*|[0-9]*\.[0-9]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_LITERAL);
	return PUT_DSL_NUMBER;
}
0x[0-9a-fA-F]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_LITERAL);
	return PUT_DSL_NUMBER;
}

\"([^\\\"]|\\.)*\" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_LITERAL);
	return PUT_DSL_STRING;
}

\"([^\\\"]|\\.)*\"\i {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_REGEXI);
	return PUT_DSL_REGEXI;
}

"NF" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"NR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"FNR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"FILENAME" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"FILENUM" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"PI" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}
"E" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return PUT_DSL_CONTEXT_VARIABLE;
}

[a-zA-Z_][a-zA-Z_0-9]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	return PUT_DSL_FCN_NAME;
}

"(" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return PUT_DSL_LPAREN;
}
")" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return PUT_DSL_RPAREN;
}

"," {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return PUT_DSL_COMMA;
}

[ \t\r\n] { }

. {
	return -1;
}
%%
