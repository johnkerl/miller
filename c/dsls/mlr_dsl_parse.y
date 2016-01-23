// vim: set filetype=none:
// (Lemon files have .y extensions like Yacc files but are not Yacc.)

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

// ----------------------------------------------------------------
mlr_dsl_body       ::= mlr_dsl_statements.

mlr_dsl_statements ::= mlr_dsl_statement.
mlr_dsl_statements ::= mlr_dsl_statement MLR_DSL_SEMICOLON mlr_dsl_statements.

mlr_dsl_statement ::= mlr_dsl_srec_assignment.
mlr_dsl_statement ::= mlr_dsl_top_level_oosvar_assignment.
mlr_dsl_statement ::= mlr_dsl_bare_boolean.
mlr_dsl_statement ::= mlr_dsl_record_filter.
mlr_dsl_statement ::= mlr_dsl_expression_gate.
mlr_dsl_statement ::= mlr_dsl_top_level_emit.

mlr_dsl_statement ::= mlr_dsl_begin_oosvar_assignment.
mlr_dsl_statement ::= mlr_dsl_begin_bare_boolean.
mlr_dsl_statement ::= mlr_dsl_begin_filter.
mlr_dsl_statement ::= mlr_dsl_begin_gate.
mlr_dsl_statement ::= mlr_dsl_begin_emit.

mlr_dsl_statement ::= mlr_dsl_end_oosvar_assignment.
mlr_dsl_statement ::= mlr_dsl_end_bare_boolean.
mlr_dsl_statement ::= mlr_dsl_end_filter.
mlr_dsl_statement ::= mlr_dsl_end_gate.
mlr_dsl_statement ::= mlr_dsl_end_emit.

// ================================================================
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.

mlr_dsl_srec_assignment(A)  ::= MLR_DSL_FIELD_NAME(B) MLR_DSL_ASSIGN(O) mlr_dsl_ternary(C). {
	// Replace "$field.name" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	B = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT, B, C);
	sllv_add(pasts, A);
}
mlr_dsl_srec_assignment(A)  ::= MLR_DSL_BRACKETED_FIELD_NAME(B) MLR_DSL_ASSIGN(O) mlr_dsl_ternary(C). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	B = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT, B, C);
	sllv_add(pasts, A);
}

mlr_dsl_top_level_oosvar_assignment(A) ::= mlr_dsl_oosvar_assignment(B). {
	A = B;
	sllv_add(pasts, A);
}

mlr_dsl_oosvar_assignment(A)  ::= mlr_dsl_oosvar_name(B) MLR_DSL_ASSIGN(O) mlr_dsl_ternary(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}

mlr_dsl_bare_boolean(A) ::= mlr_dsl_ternary(B). {
	A = B;
	sllv_add(pasts, A);
}

mlr_dsl_record_filter(A) ::= MLR_DSL_FILTER(O) mlr_dsl_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_FILTER, B);
	sllv_add(pasts, A);
}
mlr_dsl_expression_gate(A) ::= MLR_DSL_GATE(O) mlr_dsl_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_GATE, B);
	sllv_add(pasts, A);
}

mlr_dsl_top_level_emit(A) ::= mlr_dsl_emit(B). {
	A = B;
	sllv_add(pasts, A);
}

// Given "emit @a,@b,@c": since this is a bottom-up parser, we get first the "@a",
// then "@a,@b", then "@a,@b,@c", then finally "emit @a,@b,@c". So:
// * On the "@a" we make a function sub-AST called "emit @a".
// * On the "@b" we append the next argument to get "emit @a,@b".
// * On the "@c" we append the next argument to get "emit @a,@b,@c".
mlr_dsl_emit(A) ::= MLR_DSL_EMIT(O) mlr_dsl_emit_args(B). {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "emit $a," -- use some non-empty-args expr.
mlr_dsl_emit_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("emit", MLR_DSL_AST_NODE_TYPE_EMIT);
}
mlr_dsl_emit_args(A) ::= mlr_dsl_oosvar_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("emit", MLR_DSL_AST_NODE_TYPE_EMIT, B);
}
mlr_dsl_emit_args(A) ::= mlr_dsl_emit_args(B) MLR_DSL_COMMA mlr_dsl_oosvar_name(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ================================================================
mlr_dsl_begin_oosvar_assignment(A)  ::= MLR_DSL_BEGIN(X) mlr_dsl_oosvar_assignment(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_BEGIN, B);
	sllv_add(pasts, A);
}

mlr_dsl_begin_bare_boolean(A) ::= MLR_DSL_BEGIN(X) mlr_dsl_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_BEGIN, B);
	sllv_add(pasts, A);
}

mlr_dsl_begin_filter(A) ::= MLR_DSL_BEGIN(X) MLR_DSL_FILTER(O) mlr_dsl_ternary(B). {
	B = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_FILTER, B);
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_BEGIN, B);
	sllv_add(pasts, A);
}

mlr_dsl_begin_gate(A) ::= MLR_DSL_BEGIN(X) MLR_DSL_GATE(O) mlr_dsl_ternary(B). {
	B = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_GATE, B);
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_BEGIN, B);
	sllv_add(pasts, A);
}

mlr_dsl_begin_emit(A) ::= MLR_DSL_BEGIN(X) mlr_dsl_emit(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_BEGIN, B);
	sllv_add(pasts, A);
}


mlr_dsl_end_oosvar_assignment(A)  ::= MLR_DSL_END(X) mlr_dsl_oosvar_assignment(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_END, B);
	sllv_add(pasts, A);
}

mlr_dsl_end_bare_boolean(A) ::= MLR_DSL_END(X) mlr_dsl_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_END, B);
	sllv_add(pasts, A);
}

mlr_dsl_end_filter(A) ::= MLR_DSL_END(X) MLR_DSL_FILTER(O) mlr_dsl_ternary(B). {
	B = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_FILTER, B);
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_END, B);
	sllv_add(pasts, A);
}

mlr_dsl_end_gate(A) ::= MLR_DSL_END(X) MLR_DSL_GATE(O) mlr_dsl_ternary(B). {
	B = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_GATE, B);
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_END, B);
	sllv_add(pasts, A);
}

mlr_dsl_end_emit(A) ::= MLR_DSL_END(X) mlr_dsl_emit(B). {
	A = mlr_dsl_ast_node_alloc_unary(X->text, MLR_DSL_AST_NODE_TYPE_END, B);
	sllv_add(pasts, A);
}

// ================================================================
mlr_dsl_ternary(A) ::= mlr_dsl_logical_or_term(B) MLR_DSL_QUESTION_MARK mlr_dsl_ternary(C) MLR_DSL_COLON mlr_dsl_ternary(D). {
	A = mlr_dsl_ast_node_alloc_ternary("? :", MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C, D);
}

mlr_dsl_ternary(A) ::= mlr_dsl_logical_or_term(B). {
	A = B;
}

// ================================================================
mlr_dsl_logical_or_term(A) ::= mlr_dsl_logical_xor_term(B). {
	A = B;
}
mlr_dsl_logical_or_term(A) ::= mlr_dsl_logical_or_term(B) MLR_DSL_LOGICAL_OR(O) mlr_dsl_logical_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_logical_xor_term(A) ::= mlr_dsl_logical_and_term(B). {
	A = B;
}
mlr_dsl_logical_xor_term(A) ::= mlr_dsl_logical_xor_term(B) MLR_DSL_LOGICAL_XOR(O) mlr_dsl_logical_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_logical_and_term(A) ::= mlr_dsl_eqne_term(B). {
	A = B;
}
mlr_dsl_logical_and_term(A) ::= mlr_dsl_logical_and_term(B) MLR_DSL_LOGICAL_AND(O) mlr_dsl_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_eqne_term(A) ::= mlr_dsl_cmp_term(B). {
	A = B;
}
mlr_dsl_eqne_term(A) ::= mlr_dsl_eqne_term(B) MLR_DSL_MATCHES(O) mlr_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_eqne_term(A) ::= mlr_dsl_eqne_term(B) MLR_DSL_DOES_NOT_MATCH(O) mlr_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_eqne_term(A) ::= mlr_dsl_eqne_term(B) MLR_DSL_EQ(O) mlr_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_eqne_term(A) ::= mlr_dsl_eqne_term(B) MLR_DSL_NE(O) mlr_dsl_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_cmp_term(A) ::= mlr_dsl_bitwise_or_term(B). {
	A = B;
}
mlr_dsl_cmp_term(A) ::= mlr_dsl_cmp_term(B) MLR_DSL_GT(O) mlr_dsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_cmp_term(A) ::= mlr_dsl_cmp_term(B) MLR_DSL_GE(O) mlr_dsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_cmp_term(A) ::= mlr_dsl_cmp_term(B) MLR_DSL_LT(O) mlr_dsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_cmp_term(A) ::= mlr_dsl_cmp_term(B) MLR_DSL_LE(O) mlr_dsl_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_bitwise_or_term(A) ::= mlr_dsl_bitwise_xor_term(B). {
	A = B;
}
mlr_dsl_bitwise_or_term(A) ::= mlr_dsl_bitwise_or_term(B) MLR_DSL_BITWISE_OR(O) mlr_dsl_bitwise_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_bitwise_xor_term(A) ::= mlr_dsl_bitwise_and_term(B). {
	A = B;
}
mlr_dsl_bitwise_xor_term(A) ::= mlr_dsl_bitwise_xor_term(B) MLR_DSL_BITWISE_XOR(O) mlr_dsl_bitwise_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_bitwise_and_term(A) ::= mlr_dsl_bitwise_shift_term(B). {
	A = B;
}
mlr_dsl_bitwise_and_term(A) ::= mlr_dsl_bitwise_and_term(B) MLR_DSL_BITWISE_AND(O) mlr_dsl_bitwise_shift_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_bitwise_shift_term(A) ::= mlr_dsl_addsubdot_term(B). {
	A = B;
}
mlr_dsl_bitwise_shift_term(A) ::= mlr_dsl_bitwise_shift_term(B) MLR_DSL_BITWISE_LSH(O) mlr_dsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_bitwise_shift_term(A) ::= mlr_dsl_bitwise_shift_term(B) MLR_DSL_BITWISE_RSH(O) mlr_dsl_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_addsubdot_term(A) ::= mlr_dsl_muldiv_term(B). {
	A = B;
}
mlr_dsl_addsubdot_term(A) ::= mlr_dsl_addsubdot_term(B) MLR_DSL_PLUS(O) mlr_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_addsubdot_term(A) ::= mlr_dsl_addsubdot_term(B) MLR_DSL_MINUS(O) mlr_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_addsubdot_term(A) ::= mlr_dsl_addsubdot_term(B) MLR_DSL_DOT(O) mlr_dsl_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_muldiv_term(A) ::= mlr_dsl_unary_bitwise_op_term(B). {
	A = B;
}
mlr_dsl_muldiv_term(A) ::= mlr_dsl_muldiv_term(B) MLR_DSL_TIMES(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_muldiv_term(A) ::= mlr_dsl_muldiv_term(B) MLR_DSL_DIVIDE(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_muldiv_term(A) ::= mlr_dsl_muldiv_term(B) MLR_DSL_INT_DIVIDE(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
mlr_dsl_muldiv_term(A) ::= mlr_dsl_muldiv_term(B) MLR_DSL_MOD(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
mlr_dsl_unary_bitwise_op_term(A) ::= mlr_dsl_pow_term(B). {
	A = B;
}
mlr_dsl_unary_bitwise_op_term(A) ::= MLR_DSL_PLUS(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
mlr_dsl_unary_bitwise_op_term(A) ::= MLR_DSL_MINUS(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
mlr_dsl_unary_bitwise_op_term(A) ::= MLR_DSL_LOGICAL_NOT(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}
mlr_dsl_unary_bitwise_op_term(A) ::= MLR_DSL_BITWISE_NOT(O) mlr_dsl_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, C);
}

// ----------------------------------------------------------------
mlr_dsl_pow_term(A) ::= mlr_dsl_atom_or_fcn(B). {
	A = B;
}
mlr_dsl_pow_term(A) ::= mlr_dsl_atom_or_fcn(B) MLR_DSL_POW(O) mlr_dsl_pow_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}


// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.

mlr_dsl_atom_or_fcn(A) ::= mlr_dsl_field_name(B). {
	A = B;
}
mlr_dsl_field_name(A) ::= MLR_DSL_FIELD_NAME(B). {
	// not:
	// A = B;
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
mlr_dsl_field_name(A) ::= MLR_DSL_BRACKETED_FIELD_NAME(B). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}

mlr_dsl_atom_or_fcn(A) ::= mlr_dsl_oosvar_name(B). {
	A = B;
}
mlr_dsl_oosvar_name(A) ::= MLR_DSL_OOSVAR_NAME(B). {
	// not:
	// A = B;
	char* at_name = B->text;
	char* no_at_name = &at_name[1];
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}
mlr_dsl_oosvar_name(A) ::= MLR_DSL_BRACKETED_OOSVAR_NAME(B). {
	// Replace "${field.name}" with just "field.name"
	char* at_name = B->text;
	char* no_at_name = &at_name[2];
	int len = strlen(no_at_name);
	if (len > 0)
		no_at_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}

mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_NUMBER(B). {
	A = B;
}
mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_TRUE(B). {
	A = B;
}
mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_FALSE(B). {
	A = B;
}

mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}
mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_CONTEXT_VARIABLE(B). {
	A = B;
}

mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_LPAREN mlr_dsl_logical_or_term(B) MLR_DSL_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

mlr_dsl_atom_or_fcn(A) ::= MLR_DSL_FCN_NAME(O) MLR_DSL_LPAREN mlr_dsl_fcn_args(B) MLR_DSL_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "f(10,)" -- use some non-empty-args expr.
mlr_dsl_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
}

mlr_dsl_fcn_args(A) ::= mlr_dsl_logical_or_term(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
}
mlr_dsl_fcn_args(A) ::= mlr_dsl_fcn_args(B) MLR_DSL_COMMA mlr_dsl_logical_or_term(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
