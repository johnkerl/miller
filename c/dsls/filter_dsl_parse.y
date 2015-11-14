
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
filter_dsl_body(A) ::= filter_dsl_bool_expr(B). {
  A = B;
	past->proot = A;
}

// ----------------------------------------------------------------
filter_dsl_bool_expr(A) ::= filter_dsl_or_term(B). {
	A = B;
}
filter_dsl_bool_expr(A) ::= filter_dsl_bool_expr(B) FILTER_DSL_OR(O) filter_dsl_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_or_term(A) ::= filter_dsl_and_term(B). {
	A = B;
}
filter_dsl_or_term(A) ::= filter_dsl_or_term(B) FILTER_DSL_AND(O) filter_dsl_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_and_term(A) ::= filter_dsl_eqne_term(B). {
	A = B;
}
filter_dsl_and_term(A) ::= filter_dsl_and_term(B) FILTER_DSL_MATCHES(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_and_term(A) ::= filter_dsl_and_term(B) FILTER_DSL_DOES_NOT_MATCH(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_and_term(A) ::= filter_dsl_and_term(B) FILTER_DSL_EQ(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_and_term(A) ::= filter_dsl_and_term(B) FILTER_DSL_NE(O) filter_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_eqne_term(A) ::= filter_dsl_cmp_term(B). {
	A = B;
}
filter_dsl_eqne_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_GT(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_eqne_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_GE(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_eqne_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_LT(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_eqne_term(A) ::= filter_dsl_eqne_term(B) FILTER_DSL_LE(O) filter_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_cmp_term(A) ::= filter_dsl_bit_or_term(B). {
	A = B;
}
filter_dsl_cmp_term(A) ::= filter_dsl_cmp_term(B) FILTER_DSL_BIT_OR(O) filter_dsl_bit_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_bit_or_term(A) ::= filter_dsl_bit_xor_term(B). {
	A = B;
}
filter_dsl_bit_or_term(A) ::= filter_dsl_bit_or_term(B) FILTER_DSL_BIT_XOR(O) filter_dsl_bit_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_bit_xor_term(A) ::= filter_dsl_bit_and_term(B). {
	A = B;
}
filter_dsl_bit_xor_term(A) ::= filter_dsl_bit_xor_term(B) FILTER_DSL_BIT_AND(O) filter_dsl_bit_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_bit_and_term(A) ::= filter_dsl_pmdot_term(B). {
	A = B;
}
filter_dsl_bit_and_term(A) ::= filter_dsl_bit_and_term(B) FILTER_DSL_PLUS(O) filter_dsl_pmdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_bit_and_term(A) ::= filter_dsl_bit_and_term(B) FILTER_DSL_MINUS(O) filter_dsl_pmdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_bit_and_term(A) ::= filter_dsl_bit_and_term(B) FILTER_DSL_DOT(O) filter_dsl_pmdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_pmdot_term(A) ::= filter_dsl_muldiv_term(B). {
	A = B;
}
filter_dsl_pmdot_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_TIMES(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_pmdot_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_DIVIDE(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_pmdot_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_INT_DIVIDE(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
filter_dsl_pmdot_term(A) ::= filter_dsl_pmdot_term(B) FILTER_DSL_MOD(O) filter_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
filter_dsl_muldiv_term(A) ::= filter_dsl_unary_term(B). {
	A = B;
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_PLUS(O) filter_dsl_unary_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_MINUS(O) filter_dsl_unary_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
filter_dsl_muldiv_term(A) ::= FILTER_DSL_NOT(O) filter_dsl_unary_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}

// ----------------------------------------------------------------
filter_dsl_unary_term(A) ::= filter_dsl_exp_term(B). {
	A = B;
}
filter_dsl_unary_term(A) ::= filter_dsl_unary_term(B) FILTER_DSL_POW(O) filter_dsl_exp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
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
}
filter_dsl_exp_term(A) ::= FILTER_DSL_NUMBER(B). {
	A = B;
}

// Strip off the leading/trailing double quotes which are included by the lexer
filter_dsl_exp_term(A) ::= FILTER_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

// Strip off the leading '"' and trailing '"i' which are included by the lexer
filter_dsl_exp_term(A) ::= FILTER_DSL_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

filter_dsl_exp_term(A) ::= FILTER_DSL_CONTEXT_VARIABLE(B). {
	A = B;
}
filter_dsl_exp_term(A) ::= FILTER_DSL_LPAREN filter_dsl_bool_expr(B) FILTER_DSL_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

filter_dsl_exp_term(A) ::= FILTER_DSL_FCN_NAME(O) FILTER_DSL_LPAREN filter_dsl_fcn_args(B) FILTER_DSL_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// xxx need to invalidate "f(10,)" -- use some non-empty-args expr.
filter_dsl_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
}

filter_dsl_fcn_args(A) ::= filter_dsl_bool_expr(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
}
filter_dsl_fcn_args(A) ::= filter_dsl_fcn_args(B) FILTER_DSL_COMMA filter_dsl_bool_expr(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
