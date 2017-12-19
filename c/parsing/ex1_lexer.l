%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "ex1_parse.h"
#include "./ex_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html

%}

%option reentrant
%option noyywrap
%option extra-type="struct _ex_ast_node_t **"

%%

";" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_SEMICOLON;
}
"{" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LBRACE;
}
"}" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RBRACE;
}

"=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_SREC_ASSIGNMENT);
	return MD_TOKEN_ASSIGN;
}

"filter" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FILTER);
	return MD_TOKEN_FILTER;
}
"unset" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_UNSET);
	return MD_TOKEN_UNSET;
}
"emitf" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMITF);
	return MD_TOKEN_EMITF;
}
"emitp" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMITP);
	return MD_TOKEN_EMITP;
}
"emit" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMIT);
	return MD_TOKEN_EMIT;
}
"dump" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_DUMP);
	return MD_TOKEN_DUMP;
}
"edump" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EDUMP);
	return MD_TOKEN_EDUMP;
}
"print" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_PRINT);
	return MD_TOKEN_PRINT;
}
"printn" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_PRINTN);
	return MD_TOKEN_PRINTN;
}
"eprint" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EPRINT);
	return MD_TOKEN_EPRINT;
}
"eprintn" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EPRINTN);
	return MD_TOKEN_EPRINTN;
}
"begin" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_BEGIN;
}
"end" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_END;
}
"all" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_ALL);
	return MD_TOKEN_ALL;
}

"do" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_DO_WHILE);
	return MD_TOKEN_DO;
}
"while" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_WHILE);
	return MD_TOKEN_WHILE;
}
"for" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_FOR;
}
"in" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IN);
	return MD_TOKEN_IN;
}
"break" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BREAK);
	return MD_TOKEN_BREAK;
}
"continue" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTINUE);
	return MD_TOKEN_CONTINUE;
}
"if" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_IF;
}
"elif" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_ELIF;
}
"else" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_ELSE;
}

"?" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_QUESTION_MARK;
}

":" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_COLON;
}

"||" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_OR;
}
"||=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_OR_EQUALS;
}
"^^" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_XOR;
}
"^^=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_XOR_EQUALS;
}
"&&" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_AND;
}
"&&=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_AND_EQUALS;
}

"=~" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MATCHES;
}
"!=~" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOES_NOT_MATCH;
}
"==" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_EQ;
}
"!=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_NE;
}

">" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_GT;
}
">=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_GE;
}
"<" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LT;
}
"<=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LE;
}

"|" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_OR;
}
"|=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_OR_EQUALS;
}

"^" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_XOR;
}
"^=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_XOR_EQUALS;
}

"&" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_AND;
}
"&=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_AND_EQUALS;
}

"<<" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_LSH;
}
"<<=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_LSH_EQUALS;
}
">>" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_RSH;
}
">>=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_RSH_EQUALS;
}

"+" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_PLUS;
}
"+=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_PLUS_EQUALS;
}
"-" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MINUS;
}
"-=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MINUS_EQUALS;
}
"\." {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOT;
}
"\.=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOT_EQUALS;
}

"*" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_TIMES;
}
"*=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_TIMES_EQUALS;
}
"/" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DIVIDE;
}
"/=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DIVIDE_EQUALS;
}
"//" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_DIVIDE;
}
"//=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_DIVIDE_EQUALS;
}
"%" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MOD;
}
"%=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MOD_EQUALS;
}

"!" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_NOT;
}

"~" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_NOT;
}

"**" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_POW;
}
"**=" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_POW_EQUALS;
}

\$\* {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FULL_SREC);
	return MD_TOKEN_FULL_SREC;
}
\$ {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_DOLLAR_SIGN;
}

\$[a-zA-Z_0-9]+ {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FIELD_NAME);
	return MD_TOKEN_FIELD_NAME;
}
\$\{([^\}]|\\.)*\} {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FIELD_NAME);
	return MD_TOKEN_BRACED_FIELD_NAME;
}

@[a-zA-Z_0-9]+ {
	// Note: the parser depends on the at-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LITERAL);
	return MD_TOKEN_UNBRACED_OOSVAR_NAME;
}
@\{([^\}]|\\.)*\} {
	// Note: the parser depends on the percent-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LITERAL);
	return MD_TOKEN_BRACED_OOSVAR_NAME;
}
@ {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_AT_SIGN;
}
\@\* {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FULL_OOSVAR);
	return MD_TOKEN_FULL_OOSVAR;
}

"\[" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LEFT_BRACKET;
}
"\]" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RIGHT_BRACKET;
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

\"([^\\\"]|\\.)*\" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	return MD_TOKEN_STRING;
}

\"([^\\\"]|\\.)*\"\i {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_REGEXI);
	return MD_TOKEN_REGEXI;
}

true {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MD_TOKEN_TRUE;
}
false {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MD_TOKEN_FALSE;
}

"NF" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"NR" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FNR" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FILENAME" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FILENUM" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"M_PI" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"M_E" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"ENV" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_ENV);
	return MD_TOKEN_ENV;
}

[a-zA-Z_][a-zA-Z_0-9]* {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NON_SIGIL_NAME);
	return MD_TOKEN_NON_SIGIL_NAME;
}

"(" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LPAREN;
}
")" {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RPAREN;
}

"," {
	*yyextra = ex_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_COMMA;
}

[ \t\r\n] { }

. {
	return -1;
}
%%
