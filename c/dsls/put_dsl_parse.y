
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
%extra_argument {sllv_t* pasts}

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

pdsl_body ::= pdsl_assignments.                // For scan-from-string

// ----------------------------------------------------------------
pdsl_assignments ::= pdsl_assignment.
pdsl_assignments ::= pdsl_assignment PUT_DSL_SEMICOLON pdsl_assignments.

// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.

pdsl_assignment(A)  ::= PUT_DSL_FIELD_NAME(B) PUT_DSL_ASSIGN(O) pdsl_logical_or_term(C). {
	// Replace "$field.name" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	B = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	sllv_add(pasts, A);
}

pdsl_assignment(A)  ::= PUT_DSL_BRACKETED_FIELD_NAME(B) PUT_DSL_ASSIGN(O) pdsl_logical_or_term(C). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	B = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	sllv_add(pasts, A);
}

// ----------------------------------------------------------------
pdsl_logical_or_term(A) ::= pdsl_logical_xor_term(B). {
	A = B;
}
pdsl_logical_or_term(A) ::= pdsl_logical_or_term(B) PUT_DSL_LOGICAL_OR(O) pdsl_logical_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_logical_xor_term(A) ::= pdsl_logical_and_term(B). {
	A = B;
}
pdsl_logical_xor_term(A) ::= pdsl_logical_xor_term(B) PUT_DSL_LOGICAL_XOR(O) pdsl_logical_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_logical_and_term(A) ::= pdsl_eqne_term(B). {
	A = B;
}
pdsl_logical_and_term(A) ::= pdsl_logical_and_term(B) PUT_DSL_LOGICAL_AND(O) pdsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_eqne_term(A) ::= pdsl_cmp_term(B). {
	A = B;
}
pdsl_eqne_term(A) ::= pdsl_eqne_term(B) PUT_DSL_MATCHES(O) pdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_eqne_term(A) ::= pdsl_eqne_term(B) PUT_DSL_DOES_NOT_MATCH(O) pdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_eqne_term(A) ::= pdsl_eqne_term(B) PUT_DSL_EQ(O) pdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_eqne_term(A) ::= pdsl_eqne_term(B) PUT_DSL_NE(O) pdsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_cmp_term(A) ::= pdsl_bitwise_or_term(B). {
	A = B;
}
pdsl_cmp_term(A) ::= pdsl_cmp_term(B) PUT_DSL_GT(O) pdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_cmp_term(A) ::= pdsl_cmp_term(B) PUT_DSL_GE(O) pdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_cmp_term(A) ::= pdsl_cmp_term(B) PUT_DSL_LT(O) pdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_cmp_term(A) ::= pdsl_cmp_term(B) PUT_DSL_LE(O) pdsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_bitwise_or_term(A) ::= pdsl_bitwise_xor_term(B). {
	A = B;
}
pdsl_bitwise_or_term(A) ::= pdsl_bitwise_or_term(B) PUT_DSL_BITWISE_OR(O) pdsl_bitwise_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_bitwise_xor_term(A) ::= pdsl_bitwise_and_term(B). {
	A = B;
}
pdsl_bitwise_xor_term(A) ::= pdsl_bitwise_xor_term(B) PUT_DSL_BITWISE_XOR(O) pdsl_bitwise_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_bitwise_and_term(A) ::= pdsl_bitwise_shift_term(B). {
	A = B;
}
pdsl_bitwise_and_term(A) ::= pdsl_bitwise_and_term(B) PUT_DSL_BITWISE_AND(O) pdsl_bitwise_shift_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_bitwise_shift_term(A) ::= pdsl_addsubdot_term(B). {
	A = B;
}
pdsl_bitwise_shift_term(A) ::= pdsl_bitwise_shift_term(B) PUT_DSL_BITWISE_LSH(O) pdsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_bitwise_shift_term(A) ::= pdsl_bitwise_shift_term(B) PUT_DSL_BITWISE_RSH(O) pdsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_addsubdot_term(A) ::= pdsl_muldiv_term(B). {
	A = B;
}
pdsl_addsubdot_term(A) ::= pdsl_addsubdot_term(B) PUT_DSL_PLUS(O) pdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_addsubdot_term(A) ::= pdsl_addsubdot_term(B) PUT_DSL_MINUS(O) pdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_addsubdot_term(A) ::= pdsl_addsubdot_term(B) PUT_DSL_DOT(O) pdsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_muldiv_term(A) ::= pdsl_unary_bitwise_op_term(B). {
	A = B;
}
pdsl_muldiv_term(A) ::= pdsl_muldiv_term(B) PUT_DSL_TIMES(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_muldiv_term(A) ::= pdsl_muldiv_term(B) PUT_DSL_DIVIDE(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_muldiv_term(A) ::= pdsl_muldiv_term(B) PUT_DSL_INT_DIVIDE(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
pdsl_muldiv_term(A) ::= pdsl_muldiv_term(B) PUT_DSL_MOD(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
pdsl_unary_bitwise_op_term(A) ::= pdsl_pow_term(B). {
	A = B;
}
pdsl_unary_bitwise_op_term(A) ::= PUT_DSL_PLUS(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
pdsl_unary_bitwise_op_term(A) ::= PUT_DSL_MINUS(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
pdsl_unary_bitwise_op_term(A) ::= PUT_DSL_LOGICAL_NOT(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
pdsl_unary_bitwise_op_term(A) ::= PUT_DSL_BITWISE_NOT(O) pdsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}

// ----------------------------------------------------------------
pdsl_pow_term(A) ::= pdsl_atom_or_fcn(B). {
	A = B;
}
pdsl_pow_term(A) ::= pdsl_atom_or_fcn(B) PUT_DSL_POW(O) pdsl_pow_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}


// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.
pdsl_atom_or_fcn(A) ::= PUT_DSL_FIELD_NAME(B). {
	// not:
	// A = B;
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
pdsl_atom_or_fcn(A) ::= PUT_DSL_BRACKETED_FIELD_NAME(B). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
pdsl_atom_or_fcn(A) ::= PUT_DSL_NUMBER(B). {
	A = B;
}

pdsl_atom_or_fcn(A) ::= PUT_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}
pdsl_atom_or_fcn(A) ::= PUT_DSL_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

pdsl_atom_or_fcn(A) ::= PUT_DSL_CONTEXT_VARIABLE(B). {
	A = B;
}

pdsl_atom_or_fcn(A) ::= PUT_DSL_LPAREN pdsl_logical_or_term(B) PUT_DSL_RPAREN. {
	A = B;
}
///pdsl_atom_or_fcn(A) ::= PUT_DSL_MINUS(O) pdsl_atom_or_fcn(B). {
	///A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B);
///}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

pdsl_atom_or_fcn(A) ::= PUT_DSL_FCN_NAME(O) PUT_DSL_LPAREN pdsl_fcn_args(B) PUT_DSL_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// xxx need to invalidate "f(10,)" -- use some non-empty-args expr.
pdsl_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
}

pdsl_fcn_args(A) ::= pdsl_logical_or_term(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
}
pdsl_fcn_args(A) ::= pdsl_fcn_args(B) PUT_DSL_COMMA pdsl_logical_or_term(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
