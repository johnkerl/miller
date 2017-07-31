%{
#include "stdio.h"  // for fileno() prototype
#include "string.h" // for strdup() prototype
#include "mlr_dsl_parse.h"
#include "../dsl/mlr_dsl_ast.h"
// http://flex.sourceforge.net/manual/Extra-Data.html

// NOTE: any new keywords defined here should also be documented
// in dsl/mlr_dsl_cst.c's mlr_dsl_keyword_usage() et al.
%}

%option reentrant
%option noyywrap
%option extra-type="struct _mlr_dsl_ast_node_t **"

%%

";" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_SEMICOLON;
}
"{" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LBRACE;
}
"}" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RBRACE;
}

"=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_SREC_ASSIGNMENT);
	return MD_TOKEN_ASSIGN;
}
"filter" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FILTER);
	return MD_TOKEN_FILTER;
}
"unset" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_UNSET);
	return MD_TOKEN_UNSET;
}
"tee" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_TEE);
	return MD_TOKEN_TEE;
}
"emitf" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMITF);
	return MD_TOKEN_EMITF;
}
"emitp" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMITP);
	return MD_TOKEN_EMITP;
}
"emit" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EMIT);
	return MD_TOKEN_EMIT;
}
"dump" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_DUMP);
	return MD_TOKEN_DUMP;
}
"edump" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EDUMP);
	return MD_TOKEN_EDUMP;
}
"print" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_PRINT);
	return MD_TOKEN_PRINT;
}
"printn" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_PRINTN);
	return MD_TOKEN_PRINTN;
}
"eprint" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EPRINT);
	return MD_TOKEN_EPRINT;
}
"eprintn" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_EPRINTN);
	return MD_TOKEN_EPRINTN;
}
"stdout" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STDOUT);
	return MD_TOKEN_STDOUT;
}
"stderr" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STDERR);
	return MD_TOKEN_STDERR;
}

"func" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_FUNC_DEF;
}
"subr" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_SUBR_DEF;
}
"call" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_SUBR_CALL;
}
"return" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RETURN;
}
"begin" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_BEGIN;
}
"end" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_END;
}
"all" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_ALL);
	return MD_TOKEN_ALL;
}

"do" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_DO_WHILE);
	return MD_TOKEN_DO;
}
"while" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_WHILE);
	return MD_TOKEN_WHILE;
}
"for" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_FOR;
}
"in" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IN);
	return MD_TOKEN_IN;
}
"break" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BREAK);
	return MD_TOKEN_BREAK;
}
"continue" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTINUE);
	return MD_TOKEN_CONTINUE;
}
"if" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_IF;
}
"elif" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_ELIF;
}
"else" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_IF_ITEM);
	return MD_TOKEN_ELSE;
}

"?" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_QUESTION_MARK;
}

":" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_COLON;
}

"||" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_OR;
}
"||=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_OR_EQUALS;
}
"^^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_XOR;
}
"^^=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_XOR_EQUALS;
}
"&&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_AND;
}
"&&=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_AND_EQUALS;
}

"=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MATCHES;
}
"!=~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOES_NOT_MATCH;
}
"==" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_EQ;
}
"!=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_NE;
}

">" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_GT;
}
">=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_GE;
}
"<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LT;
}
"<=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LE;
}

"|" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_OR;
}
"|=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_OR_EQUALS;
}

"^" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_XOR;
}
"^=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_XOR_EQUALS;
}

"&" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_AND;
}
"&=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_AND_EQUALS;
}

"<<" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_LSH;
}
"<<=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_LSH_EQUALS;
}
">>" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_RSH;
}
">>=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_RSH_EQUALS;
}

"+" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_PLUS;
}
"+=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_PLUS_EQUALS;
}
"-" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MINUS;
}
"-=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MINUS_EQUALS;
}
"\." {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOT;
}
"\.=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DOT_EQUALS;
}
"*" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_TIMES;
}
"*=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_TIMES_EQUALS;
}
"/" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DIVIDE;
}
"/=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_DIVIDE_EQUALS;
}
"//" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_DIVIDE;
}
"//=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_DIVIDE_EQUALS;
}

"\.+" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OPLUS;
}
"\.+=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OPLUS_EQUALS;
}
"\.-" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OMINUS;
}
"\.-=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OMINUS_EQUALS;
}
"\.*" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OTIMES;
}
"\.*=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_OTIMES_EQUALS;
}
"\./" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_ODIVIDE;
}
"\./=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_ODIVIDE_EQUALS;
}
"\.//" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_ODIVIDE;
}
"\.//=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_INT_ODIVIDE_EQUALS;
}

"%" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MOD;
}
"%=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_MOD_EQUALS;
}

"!" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_LOGICAL_NOT;
}

"~" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_BITWISE_NOT;
}

"**" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_POW;
}
"**=" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_OPERATOR);
	return MD_TOKEN_POW_EQUALS;
}

\$\* {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FULL_SREC);
	return MD_TOKEN_FULL_SREC;
}
\$ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_DOLLAR_SIGN;
}

\$[a-zA-Z_0-9]+ {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FIELD_NAME);
	return MD_TOKEN_FIELD_NAME;
}
\$\{([^\}]|\\.)*\} {
	// Note: the parser depends on the dollar sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FIELD_NAME);
	return MD_TOKEN_BRACED_FIELD_NAME;
}

@[a-zA-Z_0-9]+ {
	// Note: the parser depends on the at-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LITERAL);
	return MD_TOKEN_UNBRACED_OOSVAR_NAME;
}
@\{([^\}]|\\.)*\} {
	// Note: the parser depends on the percent-sign being here. If this is changed,
	// that needs to be changed as well.
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LITERAL);
	return MD_TOKEN_BRACED_OOSVAR_NAME;
}
@ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_AT_SIGN;
}
\@\* {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FULL_OOSVAR);
	return MD_TOKEN_FULL_OOSVAR;
}

"\[" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LEFT_BRACKET;
}
"\]" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
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
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	return MD_TOKEN_NUMBER;
}
0x[0-9a-fA-F]+ {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LITERAL);
	return MD_TOKEN_NUMBER;
}

\"([^\\\"]|\\.)*\" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LITERAL);
	return MD_TOKEN_STRING;
}

\"([^\\\"]|\\.)*\"\i {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_REGEXI);
	return MD_TOKEN_REGEXI;
}

"map" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_MAP;
}
"var" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_VAR;
}
"num" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION);
	return MD_TOKEN_NUMERIC;
}
"int" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION);
	return MD_TOKEN_INT;
}
"float" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION);
	return MD_TOKEN_FLOAT;
}
"bool" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION);
	return MD_TOKEN_BOOLEAN;
}
"str" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION);
	return MD_TOKEN_STRING;
}

"true" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MD_TOKEN_TRUE;
}
"false" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_BOOLEAN_LITERAL);
	return MD_TOKEN_FALSE;
}

"IPS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"IFS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"IRS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"OPS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"OFS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"ORS" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}

"NF" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"NR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FNR" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FILENAME" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"FILENUM" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"PI" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"E" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_CONTEXT_VARIABLE);
	return MD_TOKEN_CONTEXT_VARIABLE;
}
"ENV" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_ENV);
	return MD_TOKEN_ENV;
}

[a-zA-Z_][a-zA-Z_0-9]* {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_NON_SIGIL_NAME);
	return MD_TOKEN_NON_SIGIL_NAME;
}

"(" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_LPAREN;
}
")" {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_RPAREN;
}

"," {
	*yyextra = mlr_dsl_ast_node_alloc(yytext, MD_AST_NODE_TYPE_STRIPPED_AWAY);
	return MD_TOKEN_COMMA;
}

#.*\n { }

[ \t\r\n] { }

. {
	return -1;
}
%%
