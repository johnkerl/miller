// vim: set filetype=none:
// (Lemon files have .y extensions like Yacc files but are not Yacc.)

%include {
#include <stdio.h>
#include <string.h>
#include <math.h>
#include "../containers/mlr_dsl_ast.h"
#include "../containers/sllv.h"

// ================================================================
// AST:
// * parens, commas, semis, line endings, whitespace are all stripped away
// * variable names and literal values remain as leaf nodes of the AST
// * = + - * / ** {function names} remain as non-leaf nodes of the AST
// CST: See the md_cst.c
// ================================================================

}

%token_type     {mlr_dsl_ast_node_t*}
%default_type   {mlr_dsl_ast_node_t*}
%extra_argument {mlr_dsl_ast_t* past}

//void token_destructor(mlr_dsl_ast_node_t t) {
//	printf("In token_destructor t->text(%s)=t->type(%lf)\n", t->text, t->type);
//}

//%token_destructor {token_destructor($$);}

%parse_accept {
}

// The caller is expected to provide more context.
%syntax_error {
	fprintf(stderr, "mlr DSL: syntax error.\n");
}

// ================================================================
md_body       ::= md_statements.

md_statements ::= md_statement.
md_statements ::= md_statement MD_TOKEN_SEMICOLON md_statements.

// This allows for trailing semicolon, as well as empty string (or whitespace) between semicolons:
md_statement ::= .

md_statement ::= md_main_srec_assignment.
md_statement ::= md_main_oosvar_assignment.
md_statement ::= md_main_moosvar_assignment.
md_statement ::= md_main_bare_boolean.
md_statement ::= md_main_filter.
md_statement ::= md_main_gate.
md_statement ::= md_main_emit.
md_statement ::= md_main_dump.

// E.g. 'begin { emit @count }'
md_statement ::= md_begin_block.
// E.g. 'begin emit @count'
md_statement ::= md_begin_solo_oosvar_assignment.
md_statement ::= md_begin_solo_moosvar_assignment.
md_statement ::= md_begin_solo_bare_boolean.
md_statement ::= md_begin_solo_filter.
md_statement ::= md_begin_solo_gate.
md_statement ::= md_begin_solo_emit.
md_statement ::= md_begin_solo_dump.

// E.g. 'end { emit @count }'
md_statement ::= md_end_block.
// E.g. 'end emit @count'
md_statement ::= md_end_solo_oosvar_assignment.
md_statement ::= md_end_solo_moosvar_assignment.
md_statement ::= md_end_solo_bare_boolean.
md_statement ::= md_end_solo_filter.
md_statement ::= md_end_solo_gate.
md_statement ::= md_end_solo_emit.
md_statement ::= md_end_solo_dump.

// ----------------------------------------------------------------
// This looks redundant to the above, but it avoids having pathologies such as nested 'begin { begin { ... } }'.

md_begin_block ::= MD_TOKEN_BEGIN MD_TOKEN_LEFT_BRACE md_begin_block_statements MD_TOKEN_RIGHT_BRACE.

md_begin_block_statements ::= md_begin_block_statement.
md_begin_block_statements ::= md_begin_block_statement MD_TOKEN_SEMICOLON md_begin_block_statements.

// This allows for trailing semicolon, as well as empty string (or whitespace) between semicolons:
md_begin_block_statement ::= .
md_begin_block_statement ::= md_begin_block_oosvar_assignment.
md_begin_block_statement ::= md_begin_block_moosvar_assignment.
md_begin_block_statement ::= md_begin_block_bare_boolean.
md_begin_block_statement ::= md_begin_block_filter.
md_begin_block_statement ::= md_begin_block_gate.
md_begin_block_statement ::= md_begin_block_emit.
md_begin_block_statement ::= md_begin_block_dump.

md_end_block ::= MD_TOKEN_END MD_TOKEN_LEFT_BRACE md_end_block_statements MD_TOKEN_RIGHT_BRACE.

md_end_block_statements ::= md_end_block_statement.
md_end_block_statements ::= md_end_block_statement MD_TOKEN_SEMICOLON md_end_block_statements.

// This allows for trailing semicolon, as well as empty string (or whitespace) between semicolons:
md_end_block_statement ::= .
md_end_block_statement ::= md_end_block_oosvar_assignment.
md_end_block_statement ::= md_end_block_moosvar_assignment.
md_end_block_statement ::= md_end_block_bare_boolean.
md_end_block_statement ::= md_end_block_filter.
md_end_block_statement ::= md_end_block_gate.
md_end_block_statement ::= md_end_block_emit.
md_end_block_statement ::= md_end_block_dump.


// ================================================================
// These are top-level; they update the AST top-level statement-lists.

md_main_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_ASSIGN(O) md_ternary(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B, C);
	sllv_add(past->pmain_statements, A);
}
md_main_oosvar_assignment(A) ::= md_oosvar_assignment(B). {
	A = B;
	sllv_add(past->pmain_statements, A);
}
md_main_moosvar_assignment(A) ::= md_moosvar_assignment(B). {
	A = B;
	sllv_add(past->pmain_statements, A);
}
md_main_bare_boolean(A) ::= md_ternary(B). {
	A = B;
	sllv_add(past->pmain_statements, A);
}
md_main_filter(A) ::= MD_TOKEN_FILTER(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
	sllv_add(past->pmain_statements, A);
}
md_main_gate(A) ::= MD_TOKEN_GATE(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_GATE, B);
	sllv_add(past->pmain_statements, A);
}
md_main_emit(A) ::= md_emit(B). {
	A = B;
	sllv_add(past->pmain_statements, A);
}
md_main_dump(A) ::= md_dump(B). {
	A = B;
	sllv_add(past->pmain_statements, A);
}

// ----------------------------------------------------------------
// These are top-level; they update the AST top-level statement-lists.

md_begin_solo_oosvar_assignment(A)  ::= MD_TOKEN_BEGIN md_oosvar_assignment(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_moosvar_assignment(A)  ::= MD_TOKEN_BEGIN md_moosvar_assignment(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_bare_boolean(A) ::= MD_TOKEN_BEGIN md_ternary(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_filter(A) ::= MD_TOKEN_BEGIN MD_TOKEN_FILTER(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_gate(A) ::= MD_TOKEN_BEGIN MD_TOKEN_GATE(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_GATE, B);
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_emit(A) ::= MD_TOKEN_BEGIN md_emit(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_solo_dump(A) ::= MD_TOKEN_BEGIN md_dump(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}

md_begin_block_oosvar_assignment(A)  ::= md_oosvar_assignment(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_block_moosvar_assignment(A)  ::= md_moosvar_assignment(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}

md_begin_block_bare_boolean(A) ::= md_ternary(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_block_filter(A) ::= MD_TOKEN_FILTER(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
	sllv_add(past->pbegin_statements, A);
}
md_begin_block_gate(A) ::= MD_TOKEN_GATE(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_GATE, B);
	sllv_add(past->pbegin_statements, A);
}
md_begin_block_emit(A) ::= md_emit(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}
md_begin_block_dump(A) ::= md_dump(B). {
	A = B;
	sllv_add(past->pbegin_statements, A);
}

// ----------------------------------------------------------------
// These are top-level; they update the AST top-level statement-lists.

md_end_solo_oosvar_assignment(A)  ::= MD_TOKEN_END md_oosvar_assignment(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_solo_moosvar_assignment(A)  ::= MD_TOKEN_END md_moosvar_assignment(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_solo_bare_boolean(A) ::= MD_TOKEN_END md_ternary(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_solo_filter(A) ::= MD_TOKEN_END MD_TOKEN_FILTER(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
	sllv_add(past->pend_statements, A);
}
md_end_solo_gate(A) ::= MD_TOKEN_END MD_TOKEN_GATE(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_GATE, B);
	sllv_add(past->pend_statements, A);
}
md_end_solo_emit(A) ::= MD_TOKEN_END md_emit(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_solo_dump(A) ::= MD_TOKEN_END md_dump(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}

md_end_block_oosvar_assignment(A)  ::= md_oosvar_assignment(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_block_moosvar_assignment(A)  ::= md_moosvar_assignment(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_block_bare_boolean(A) ::= md_ternary(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_block_filter(A) ::= MD_TOKEN_FILTER(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
	sllv_add(past->pend_statements, A);
}
md_end_block_gate(A) ::= MD_TOKEN_GATE(O) md_ternary(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_GATE, B);
	sllv_add(past->pend_statements, A);
}
md_end_block_emit(A) ::= md_emit(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}
md_end_block_dump(A) ::= md_dump(B). {
	A = B;
	sllv_add(past->pend_statements, A);
}

// ----------------------------------------------------------------
md_oosvar_assignment(A)  ::= md_oosvar_name(B) MD_TOKEN_ASSIGN(O) md_ternary(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}
md_moosvar_assignment(A)  ::= md_moosvar_name(B) MD_TOKEN_ASSIGN(O) md_ternary(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_MOOSVAR_ASSIGNMENT, B, C);
}
md_moosvar_assignment(A)  ::= md_keyed_moosvar_name(B) MD_TOKEN_ASSIGN(O) md_ternary(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_MOOSVAR_ASSIGNMENT, B, C);
}

// ----------------------------------------------------------------
// Given "emit @a,@b,@c": since this is a bottom-up parser, we get first the "@a",
// then "@a,@b", then "@a,@b,@c", then finally "emit @a,@b,@c". So:
// * On the "@a" we make a sub-AST called "temp @a" (although we could call it "emit").
// * On the "@b" we append the next argument to get "temp @a,@b".
// * On the "@c" we append the next argument to get "temp @a,@b,@c".
// * On the "emit" we change the name to get "emit @a,@b,@c".

md_emit(A) ::= MD_TOKEN_EMIT(O) md_emit_args(B). {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "emit $a," -- use some non-empty-args expr.
md_emit_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("temp", MD_AST_NODE_TYPE_EMIT);
}

md_emit_args(A) ::= md_oosvar_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMIT, B);
}
md_emit_args(A) ::= md_emit_args(B) MD_TOKEN_COMMA md_oosvar_name(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

md_emit_args(A) ::= md_moosvar_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMIT, B);
}
md_emit_args(A) ::= md_emit_args(B) MD_TOKEN_COMMA md_moosvar_name(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
// Temporary dev/debug hook for moosvars
md_dump(A) ::= MD_TOKEN_DUMP(O). {
	A = mlr_dsl_ast_node_alloc_zary(O->text, MD_AST_NODE_TYPE_DUMP);
}

// ================================================================
md_ternary(A) ::= md_logical_or_term(B) MD_TOKEN_QUESTION_MARK md_ternary(C) MD_TOKEN_COLON md_ternary(D). {
	A = mlr_dsl_ast_node_alloc_ternary("? :", MD_AST_NODE_TYPE_OPERATOR, B, C, D);
}

md_ternary(A) ::= md_logical_or_term(B). {
	A = B;
}

// ================================================================
md_logical_or_term(A) ::= md_logical_xor_term(B). {
	A = B;
}
md_logical_or_term(A) ::= md_logical_or_term(B) MD_TOKEN_LOGICAL_OR(O) md_logical_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_logical_xor_term(A) ::= md_logical_and_term(B). {
	A = B;
}
md_logical_xor_term(A) ::= md_logical_xor_term(B) MD_TOKEN_LOGICAL_XOR(O) md_logical_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_logical_and_term(A) ::= md_eqne_term(B). {
	A = B;
}
md_logical_and_term(A) ::= md_logical_and_term(B) MD_TOKEN_LOGICAL_AND(O) md_eqne_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_eqne_term(A) ::= md_cmp_term(B). {
	A = B;
}
md_eqne_term(A) ::= md_eqne_term(B) MD_TOKEN_MATCHES(O) md_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_eqne_term(A) ::= md_eqne_term(B) MD_TOKEN_DOES_NOT_MATCH(O) md_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_eqne_term(A) ::= md_eqne_term(B) MD_TOKEN_EQ(O) md_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_eqne_term(A) ::= md_eqne_term(B) MD_TOKEN_NE(O) md_cmp_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_cmp_term(A) ::= md_bitwise_or_term(B). {
	A = B;
}
md_cmp_term(A) ::= md_cmp_term(B) MD_TOKEN_GT(O) md_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_cmp_term(A) ::= md_cmp_term(B) MD_TOKEN_GE(O) md_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_cmp_term(A) ::= md_cmp_term(B) MD_TOKEN_LT(O) md_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_cmp_term(A) ::= md_cmp_term(B) MD_TOKEN_LE(O) md_bitwise_or_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_bitwise_or_term(A) ::= md_bitwise_xor_term(B). {
	A = B;
}
md_bitwise_or_term(A) ::= md_bitwise_or_term(B) MD_TOKEN_BITWISE_OR(O) md_bitwise_xor_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_bitwise_xor_term(A) ::= md_bitwise_and_term(B). {
	A = B;
}
md_bitwise_xor_term(A) ::= md_bitwise_xor_term(B) MD_TOKEN_BITWISE_XOR(O) md_bitwise_and_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_bitwise_and_term(A) ::= md_bitwise_shift_term(B). {
	A = B;
}
md_bitwise_and_term(A) ::= md_bitwise_and_term(B) MD_TOKEN_BITWISE_AND(O) md_bitwise_shift_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_bitwise_shift_term(A) ::= md_addsubdot_term(B). {
	A = B;
}
md_bitwise_shift_term(A) ::= md_bitwise_shift_term(B) MD_TOKEN_BITWISE_LSH(O) md_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_bitwise_shift_term(A) ::= md_bitwise_shift_term(B) MD_TOKEN_BITWISE_RSH(O) md_addsubdot_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_addsubdot_term(A) ::= md_muldiv_term(B). {
	A = B;
}
md_addsubdot_term(A) ::= md_addsubdot_term(B) MD_TOKEN_PLUS(O) md_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_addsubdot_term(A) ::= md_addsubdot_term(B) MD_TOKEN_MINUS(O) md_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_addsubdot_term(A) ::= md_addsubdot_term(B) MD_TOKEN_DOT(O) md_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_muldiv_term(A) ::= md_unary_bitwise_op_term(B). {
	A = B;
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_TIMES(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_DIVIDE(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_INT_DIVIDE(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_MOD(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
md_unary_bitwise_op_term(A) ::= md_pow_term(B). {
	A = B;
}
md_unary_bitwise_op_term(A) ::= MD_TOKEN_PLUS(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_OPERATOR, C);
}
md_unary_bitwise_op_term(A) ::= MD_TOKEN_MINUS(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_OPERATOR, C);
}
md_unary_bitwise_op_term(A) ::= MD_TOKEN_LOGICAL_NOT(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_OPERATOR, C);
}
md_unary_bitwise_op_term(A) ::= MD_TOKEN_BITWISE_NOT(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_OPERATOR, C);
}

// ----------------------------------------------------------------
md_pow_term(A) ::= md_atom_or_fcn(B). {
	A = B;
}
md_pow_term(A) ::= md_atom_or_fcn(B) MD_TOKEN_POW(O) md_pow_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}


// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.

md_atom_or_fcn(A) ::= md_field_name(B). {
	A = B;
}
md_field_name(A) ::= MD_TOKEN_FIELD_NAME(B). {
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
md_field_name(A) ::= MD_TOKEN_BRACED_FIELD_NAME(B). {
	// Replace "${field.name}" with just "field.name"
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[2];
	int len = strlen(no_dollar_name);
	if (len > 0)
		no_dollar_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}

md_atom_or_fcn(A) ::= md_oosvar_name(B). {
	A = B;
}
md_oosvar_name(A) ::= MD_TOKEN_OOSVAR_NAME(B). {
	char* at_name = B->text;
	char* no_at_name = &at_name[1];
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}
md_oosvar_name(A) ::= MD_TOKEN_BRACED_OOSVAR_NAME(B). {
	// Replace "@{field.name}" with just "field.name"
	char* at_name = B->text;
	char* no_at_name = &at_name[2];
	int len = strlen(no_at_name);
	if (len > 0)
		no_at_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}

md_atom_or_fcn(A) ::= md_keyed_moosvar_name(B). {
	A = B;
}
md_atom_or_fcn(A) ::= md_moosvar_name(B). {
	A = B;
}

md_keyed_moosvar_name(A) ::= md_moosvar_name(B) MD_TOKEN_LEFT_BRACKET md_ternary(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_binary("[]", MD_AST_NODE_TYPE_MOOSVAR_LEVEL_KEY, B, C);
}
md_keyed_moosvar_name(A) ::= md_keyed_moosvar_name(B) MD_TOKEN_LEFT_BRACKET md_ternary(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_binary("[]", MD_AST_NODE_TYPE_MOOSVAR_LEVEL_KEY, B, C);
}

md_moosvar_name(A) ::= MD_TOKEN_MOOSVAR_NAME(B). {
	char* at_name = B->text;
	char* no_at_name = &at_name[2];
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}
md_moosvar_name(A) ::= MD_TOKEN_BRACED_MOOSVAR_NAME(B). {
	// Replace "@%{field.name}" with just "field.name"
	char* at_name = B->text;
	char* no_at_name = &at_name[3];
	int len = strlen(no_at_name);
	if (len > 0)
		no_at_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc(no_at_name, B->type);
}

md_atom_or_fcn(A) ::= MD_TOKEN_NUMBER(B). {
	A = B;
}
md_atom_or_fcn(A) ::= MD_TOKEN_TRUE(B). {
	A = B;
}
md_atom_or_fcn(A) ::= MD_TOKEN_FALSE(B). {
	A = B;
}

md_atom_or_fcn(A) ::= MD_TOKEN_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}
md_atom_or_fcn(A) ::= MD_TOKEN_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}

md_atom_or_fcn(A) ::= MD_TOKEN_CONTEXT_VARIABLE(B). {
	A = B;
}

md_atom_or_fcn(A) ::= MD_TOKEN_LPAREN md_logical_or_term(B) MD_TOKEN_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

md_atom_or_fcn(A) ::= MD_TOKEN_FCN_NAME(O) MD_TOKEN_LPAREN md_fcn_args(B) MD_TOKEN_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "f(10,)" -- use some non-empty-args expr.
md_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MD_AST_NODE_TYPE_FUNCTION_NAME);
}

md_fcn_args(A) ::= md_logical_or_term(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_FUNCTION_NAME, B);
}
md_fcn_args(A) ::= md_fcn_args(B) MD_TOKEN_COMMA md_logical_or_term(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
