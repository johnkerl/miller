
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
filter_dsl_body(A) ::= fdsl_logical_or_term(B). {
	A = B;
	past->proot = A;
}

// ----------------------------------------------------------------
fdsl_logical_or_term(A) ::= fdsl_logical_xor_term(B). {
	A = B;
}
fdsl_logical_or_term(A) ::= fdsl_logical_or_term(B) FILTER_DSL_LOGICAL_OR(O) fdsl_logical_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_logical_xor_term(A) ::= fdsl_logical_and_term(B). {
	A = B;
}
fdsl_logical_xor_term(A) ::= fdsl_logical_xor_term(B) FILTER_DSL_LOGICAL_XOR(O) fdsl_logical_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_logical_and_term(A) ::= fdsl_eqne_term(B). {
	A = B;
}
fdsl_logical_and_term(A) ::= fdsl_logical_and_term(B) FILTER_DSL_LOGICAL_AND(O) fdsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_eqne_term(A) ::= fdsl_cmp_term(B). {
	A = B;
}
fdsl_eqne_term(A) ::= fdsl_eqne_term(B) FILTER_DSL_MATCHES(O) fdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_eqne_term(A) ::= fdsl_eqne_term(B) FILTER_DSL_DOES_NOT_MATCH(O) fdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_eqne_term(A) ::= fdsl_eqne_term(B) FILTER_DSL_EQ(O) fdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_eqne_term(A) ::= fdsl_eqne_term(B) FILTER_DSL_NE(O) fdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_cmp_term(A) ::= fdsl_bitwise_or_term(B). {
	A = B;
}
fdsl_cmp_term(A) ::= fdsl_cmp_term(B) FILTER_DSL_GT(O) fdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_cmp_term(A) ::= fdsl_cmp_term(B) FILTER_DSL_GE(O) fdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_cmp_term(A) ::= fdsl_cmp_term(B) FILTER_DSL_LT(O) fdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_cmp_term(A) ::= fdsl_cmp_term(B) FILTER_DSL_LE(O) fdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_bitwise_or_term(A) ::= fdsl_bitwise_xor_term(B). {
	A = B;
}
fdsl_bitwise_or_term(A) ::= fdsl_bitwise_or_term(B) FILTER_DSL_BITWISE_OR(O) fdsl_bitwise_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_bitwise_xor_term(A) ::= fdsl_bitwise_and_term(B). {
	A = B;
}
fdsl_bitwise_xor_term(A) ::= fdsl_bitwise_xor_term(B) FILTER_DSL_BITWISE_XOR(O) fdsl_bitwise_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_bitwise_and_term(A) ::= fdsl_bitwise_shift_term(B). {
	A = B;
}
fdsl_bitwise_and_term(A) ::= fdsl_bitwise_and_term(B) FILTER_DSL_BITWISE_AND(O) fdsl_bitwise_shift_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_bitwise_shift_term(A) ::= fdsl_addsubdot_term(B). {
	A = B;
}
fdsl_bitwise_shift_term(A) ::= fdsl_bitwise_shift_term(B) FILTER_DSL_BITWISE_LSH(O) fdsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_bitwise_shift_term(A) ::= fdsl_bitwise_shift_term(B) FILTER_DSL_BITWISE_RSH(O) fdsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_addsubdot_term(A) ::= fdsl_muldiv_term(B). {
	A = B;
}
fdsl_addsubdot_term(A) ::= fdsl_addsubdot_term(B) FILTER_DSL_PLUS(O) fdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_addsubdot_term(A) ::= fdsl_addsubdot_term(B) FILTER_DSL_MINUS(O) fdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_addsubdot_term(A) ::= fdsl_addsubdot_term(B) FILTER_DSL_DOT(O) fdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_muldiv_term(A) ::= fdsl_unary_bitwise_op_term(B). {
	A = B;
}
fdsl_muldiv_term(A) ::= fdsl_muldiv_term(B) FILTER_DSL_TIMES(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_muldiv_term(A) ::= fdsl_muldiv_term(B) FILTER_DSL_DIVIDE(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_muldiv_term(A) ::= fdsl_muldiv_term(B) FILTER_DSL_INT_DIVIDE(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
fdsl_muldiv_term(A) ::= fdsl_muldiv_term(B) FILTER_DSL_MOD(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
fdsl_unary_bitwise_op_term(A) ::= fdsl_pow_term(B). {
	A = B;
}
fdsl_unary_bitwise_op_term(A) ::= FILTER_DSL_PLUS(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
fdsl_unary_bitwise_op_term(A) ::= FILTER_DSL_MINUS(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
fdsl_unary_bitwise_op_term(A) ::= FILTER_DSL_LOGICAL_NOT(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
fdsl_unary_bitwise_op_term(A) ::= FILTER_DSL_BITWISE_NOT(O) fdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}

// ----------------------------------------------------------------
fdsl_pow_term(A) ::= fdsl_atom_or_fcn(B). {
	A = B;
}
fdsl_pow_term(A) ::= fdsl_atom_or_fcn(B) FILTER_DSL_POW(O) fdsl_pow_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.
fdsl_atom_or_fcn(A) ::= FILTER_DSL_FIELD_NAME(B). {
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
fdsl_atom_or_fcn(A) ::= FILTER_DSL_BRACKETED_FIELD_NAME(B). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
fdsl_atom_or_fcn(A) ::= FILTER_DSL_NUMBER(B). {
	A = B;
}

// Strip off the leading/trailing double quotes which are included by the lexer
fdsl_atom_or_fcn(A) ::= FILTER_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

// Strip off the leading '"' and trailing '"i' which are included by the lexer
fdsl_atom_or_fcn(A) ::= FILTER_DSL_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

fdsl_atom_or_fcn(A) ::= FILTER_DSL_CONTEXT_VARIABLE(B). {
	A = B;
}
fdsl_atom_or_fcn(A) ::= FILTER_DSL_LPAREN fdsl_logical_or_term(B) FILTER_DSL_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

fdsl_atom_or_fcn(A) ::= FILTER_DSL_FCN_NAME(O) FILTER_DSL_LPAREN filter_dsl_fcn_args(B) FILTER_DSL_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// xxx need to invalidate "f(10,)" -- use some non-empty-args expr.
filter_dsl_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
}

filter_dsl_fcn_args(A) ::= fdsl_logical_or_term(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
}
filter_dsl_fcn_args(A) ::= filter_dsl_fcn_args(B) FILTER_DSL_COMMA fdsl_logical_or_term(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
