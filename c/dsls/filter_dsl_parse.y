
%include {
#include <stdio.h>
#include <string.h>
#include <math.h>
#include "../containers/mlr_dsl_ast.h"
#include "../containers/sllv.h"

// ================================================================
// CST: See the *.y files
// AST:
// * parens, commas, semis, line endings, whitespace are all stripped away
// * variable names and literal values remain as leaf nodes of the AST
// * = + - * / ** {function names} remain as non-leaf nodes of the AST
// ================================================================

}

%token_type     {mlr_dsl_ast_node_t*}
%default_type   {mlr_dsl_ast_node_t*}
%extra_argument {mlr_dsl_ast_node_holder_t* past}

//void token_destructor(mlr_dsl_ast_node_t t) {
//	printf("In token_destructor t->text(%s)=t->type(%lf)\n", t->text, t->type);
//}

//%token_destructor {token_destructor($$);}

%parse_accept {
	//printf("End of parse.\n");
	//printf("End of parse: proot=%p\n", past->proot);
	//mlr_dsl_ast_node_print(past->proot);
	//printf("End of parse: tree end\n");
}

%syntax_error {
	fprintf(stderr, "Syntax error!\n");
}

// ----------------------------------------------------------------
filter_dsl_body(A) ::= filter_dsl_bool_expr(B). {                 // For scan-from-string
  A = B;
	past->proot = A;
}
filter_dsl_body(A) ::= filter_dsl_bool_expr(B) FILTER_DSL_EOL. {  // For scan-from-stdin
  A = B;
	past->proot = A;
}

filter_dsl_bool_expr(A) ::= filter_dsl_or_term(B). {
	//printf("[FDSL-PARSE] BEXPR-IS-OR-TERM %s\n", B->text);
	A = B;
}
filter_dsl_bool_expr(A) ::= filter_dsl_or_term(B) FILTER_DSL_OR(O) filter_dsl_bool_expr(C). {
	//printf("[FDSL-PARSE] OR %s %s %s\n", B->text, O->text, C->text);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_or_term(A) ::= filter_dsl_and_term(B). {
	//printf("[FDSL-PARSE] OR-TERM-IS-AND-TERM %s\n", B->text);
	A = B;
}
filter_dsl_or_term(A) ::= filter_dsl_and_term(B) FILTER_DSL_AND(O) filter_dsl_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] AND %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_and_term(A) ::= filter_dsl_eqne_term(B). {
	//printf("[FDSL-PARSE] AND-TERM-IS-EQNE-TERM %s\n", B->text);
	A = B;
}
filter_dsl_and_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_EQ(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] EQ %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_and_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_NE(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] NE %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B). {
	//printf("[FDSL-PARSE] EQNE-TERM-IS-CMP-TERM %s\n", B->text);
	A = B;
}
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B) FILTER_DSL_GT(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] GT %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B) FILTER_DSL_GE(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] GE %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B) FILTER_DSL_LT(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] LT %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B) FILTER_DSL_LE(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] LE %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_cmp_term(A) ::= filter_dsl_pmdot_term(B). {
	//printf("[FDSL-PARSE] CMP-TERM-IS-PMDOT-TERM %s\n", B->text);
	A = B;
}
filter_dsl_cmp_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_PLUS(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] PMDOT %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_cmp_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_MINUS(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] PMDOT %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_cmp_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_DOT(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] PMDOT %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_pmdot_term(A) ::= filter_dsl_muldiv_term(B). {
	//printf("[FDSL-PARSE] PMDOT-TERM-IS-MULDIV-TERM %s\n", B->text);
	A = B;
}
filter_dsl_pmdot_term(A) ::= filter_dsl_muldiv_term(B) FILTER_DSL_TIMES(O) filter_dsl_pmdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] MULDIV %s %s %s\n", B->text, O->text, C->text);
}
filter_dsl_pmdot_term(A) ::= filter_dsl_muldiv_term(B) FILTER_DSL_DIVIDE(O) filter_dsl_pmdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] MULDIV %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_muldiv_term(A) ::= filter_dsl_unary_term(B). {
	//printf("[FDSL-PARSE] MULDIV-IS-UNARY-TERM %s\n", B->text);
	A = B;
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_PLUS(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
	//printf("[FDSL-PARSE] UNARY %s %s\n", O->text, C->text);
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_MINUS(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
	//printf("[FDSL-PARSE] UNARY %s %s\n", O->text, C->text);
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_NOT(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
	//printf("[FDSL-PARSE] UNARY %s %s\n", O->text, C->text);
}

// ----------------------------------------------------------------
filter_dsl_unary_term(A) ::= filter_dsl_exp_term(B). {
	//printf("[FDSL-PARSE] UNARY-IS-EXP %s\n", B->text);
	A = B;
}
filter_dsl_unary_term(A) ::= filter_dsl_unary_term(B) FILTER_DSL_POW(O) filter_dsl_exp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	//printf("[FDSL-PARSE] POW %s %s %s\n", B->text, O->text, C->text);
}

// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.
filter_dsl_exp_term(A) ::= FILTER_DSL_FIELD_NAME(B). {
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
	//printf("[FDSL-PARSE] EXP-TERM-IS-FIELD-NAME %s\n", B->text);
}
filter_dsl_exp_term(A) ::= FILTER_DSL_NUMBER(B). {
	A = B;
	//printf("[FDSL-PARSE] EXP-TERM-IS-NUMBER %s\n", B->text);
}
// xxx commment me more
filter_dsl_exp_term(A) ::= FILTER_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	// xxx make/call method
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
	//printf("[FDSL-PARSE] EXP-TERM-IS-STRING %s\n", B->text);
}
filter_dsl_exp_term(A) ::= FILTER_DSL_CONTEXT_VARIABLE(B). {
	A = B;
	//printf("[FDSL-PARSE] EXP-TERM-IS-CTX-VAR %s\n", B->text);
}
filter_dsl_exp_term(A) ::= FILTER_DSL_LPAREN filter_dsl_bool_expr(B) FILTER_DSL_RPAREN. {
	A = B;
	//printf("[FDSL-PARSE] EXP-TERM-PARENS (%s)\n", B->text);
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

filter_dsl_exp_term(A) ::= FILTER_DSL_FCN_NAME(O) FILTER_DSL_LPAREN filter_dsl_fcn_args(B) FILTER_DSL_RPAREN. {
	//printf("[FDSL-PARSE] EXPITEM-IS-FCN %s(%04x)\n", O->text, B->type);
	//printf("[FDSL-PARSE] EXPITEM-IS-FCN %s(%s)\n", O->text, B->text);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
	//mlr_dsl_ast_node_print(A);
}
// xxx need to invalidate "f(10,)" -- use some non-empty-args expr.
filter_dsl_fcn_args(A) ::= . {
	//printf("[FDSL-PARSE] FCNARGS empty\n");
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
	//mlr_dsl_ast_node_print(A);
}

filter_dsl_fcn_args(A) ::= filter_dsl_bool_expr(B). {
	//printf("[FDSL-PARSE] FCNARG %04x\n", B->type);
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
	//mlr_dsl_ast_node_print(A);
}
filter_dsl_fcn_args(A) ::= filter_dsl_fcn_args(B) FILTER_DSL_COMMA filter_dsl_bool_expr(C). {
	//printf("[FDSL-PARSE] FCNARGS %04x, %04x\n", B->type, C->type);
	//printf("[FDSL-PARSE] FCNARGS %s, %s\n", B->text, C->text);
	A = mlr_dsl_ast_node_append_arg(B, C);
	//mlr_dsl_ast_node_print(A);
}
