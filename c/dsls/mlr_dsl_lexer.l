%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "mlr_dsl_parse.h"
#include "../containers/mlr_dsl_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html

// xxx SREC_ASSIGNMENT or OOSREC_ASSIGNMENT; confusing (and irrelevant) labels here

%}

%option reentrant
%option noyywrap
%option extra-type="struct _mlr_dsl_ast_node_t **"

%%
";" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return MLR_DSL_SEMICOLON;
}

"=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT);
	return MLR_DSL_ASSIGN;
}

"filter" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FILTER);
	return MLR_DSL_FILTER;
}
"gate" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_GATE);
	return MLR_DSL_GATE;
}
"emit" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_EMIT);
	return MLR_DSL_EMIT;
}
"begin" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_BEGIN);
	return MLR_DSL_BEGIN;
}
"end" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_END);
	return MLR_DSL_END;
}

"?" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_QUESTION_MARK;
}

":" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_COLON;
}

"||" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LOGICAL_OR;
}
"^^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LOGICAL_XOR;
}
"&&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LOGICAL_AND;
}

"=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_MATCHES;
}
"!=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_DOES_NOT_MATCH;
}
"==" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_EQ;
}
"!=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_NE;
}

">" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_GT;
}
">=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_GE;
}
"<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LT;
}
"<=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LE;
}

"|" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_OR;
}

"^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_XOR;
}

"&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_AND;
}

"<<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_LSH;
}
">>" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_RSH;
}

"+" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_PLUS;
}
"-" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_MINUS;
}
"\." {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_DOT;
}

"*" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_TIMES;
}
"/" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_DIVIDE;
}
"//" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_INT_DIVIDE;
}
"%" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_MOD;
}

"!" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_LOGICAL_NOT;
}

"~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_BITWISE_NOT;
}

"**" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OPERATOR);
	return MLR_DSL_POW;
}

\$[a-zA-Z_0-9]+ {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	return MLR_DSL_FIELD_NAME;
}
\$\{([^\}]|\\.)*\} {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FIELD_NAME);
	return MLR_DSL_BRACKETED_FIELD_NAME;
}

@[a-zA-Z_0-9]+ {
	// Note: the parser depends on the at-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OOSVAR_NAME);
	return MLR_DSL_OOSVAR_NAME;
}
@\{([^\}]|\\.)*\} {
	// Note: the parser depends on the at-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_OOSVAR_NAME);
	return MLR_DSL_BRACKETED_OOSVAR_NAME;
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
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRNUM_LITERAL);
	return MLR_DSL_NUMBER;
}
0x[0-9a-fA-F]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRNUM_LITERAL);
	return MLR_DSL_NUMBER;
}

\"([^\\\"]|\\.)*\" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRNUM_LITERAL);
	return MLR_DSL_STRING;
}

\"([^\\\"]|\\.)*\"\i {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_REGEXI);
	return MLR_DSL_REGEXI;
}

true {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MLR_DSL_TRUE;
}
false {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MLR_DSL_FALSE;
}

"NF" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"NR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"FNR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"FILENAME" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"FILENUM" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"PI" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}
"E" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MLR_DSL_CONTEXT_VARIABLE;
}

[a-zA-Z_][a-zA-Z_0-9]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	return MLR_DSL_FCN_NAME;
}

"(" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return MLR_DSL_LPAREN;
}
")" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return MLR_DSL_RPAREN;
}

"," {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY);
	return MLR_DSL_COMMA;
}

[ \t\r\n] { }

. {
	return -1;
}
%%
