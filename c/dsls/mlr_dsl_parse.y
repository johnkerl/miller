// vim: set filetype=none:
// (Lemon files have .y extensions like Yacc files but are not Yacc.)

%include {
#include <stdio.h>
#include <string.h>
#include <math.h>
#include "../lib/mlrutil.h"
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
md_body ::= md_statement_list(B). {
	past->proot = B;
}

// ================================================================
// ================================================================
// NEW GRAMMAR
// ================================================================
// ================================================================

// xxx need top-level w/ include/exclude begin/end. or, again, reject @ cst.

// ----------------------------------------------------------------
// Given "$a=1;$b=2;$c=3": since this is a bottom-up parser, we get first the "$a=1", then
// "$a=1;$b=2", then "$a=1;$b=2;$c=3", then finally realize that's the top level, or it's embedded
// within a cond-block, or for-loop body, etc.

// * On the "$a=1" we make a sub-AST called "list" with child $a=1.
// * On the "$b=2" we append the next argument to get "list" having children $a=1 and $b=2.
// * On the "$c=3" we append the next argument to get "list" having children $a=1, $b=2, and $c=3.
//
// We handle statements of the form ' ; ; ' by parsing the empty spaces around the semicolons as NOP nodes.
// But, the NOP nodes are immediately stripped here and are not included in the AST we return.

md_statement_list(A) ::= md_statement(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		A = mlr_dsl_ast_node_alloc_zary("list", MD_AST_NODE_TYPE_STATEMENT_LIST);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("list", MD_AST_NODE_TYPE_STATEMENT_LIST, B);
	}
}
md_statement_list(A) ::= md_statement_list(B) MD_TOKEN_SEMICOLON md_statement(C). {
	if (C->type == MD_AST_NODE_TYPE_NOP) {
		A = B;
	} else {
		A = mlr_dsl_ast_node_append_arg(B, C);
	}
	//mlr_dsl_ast_node_print(C);
}

// This allows for trailing semicolon, as well as empty string (or whitespace) between semicolons:
md_statement(A) ::= . {
    A = mlr_dsl_ast_node_alloc_zary("nop", MD_AST_NODE_TYPE_NOP);
}

// Begin/end (non-nestable)
md_statement ::= md_begin_block.
md_statement ::= md_end_block.

// Nested control structures:
md_statement ::= md_cond_block.
md_statement ::= md_while_block.
md_statement ::= md_do_while_block.
md_statement ::= md_for_loop_full_srec.
md_statement ::= md_for_loop_full_oosvar.
md_statement ::= md_for_loop_oosvar.
md_statement ::= md_if_chain.

// Not valid in begin/end since they refer to srecs:
md_statement ::= md_srec_assignment.
md_statement ::= md_srec_indirect_assignment.
md_statement ::= md_oosvar_from_full_srec_assignment.
md_statement ::= md_full_srec_from_oosvar_assignment.

// Valid in begin/end since they don't refer to srecs (although the RHSs might):
md_statement ::= md_bare_boolean.
md_statement ::= md_oosvar_assignment.
md_statement ::= md_filter.
md_statement ::= md_unset.
md_statement ::= md_emitf.
md_statement ::= md_emitp.
md_statement ::= md_emit.
md_statement ::= md_dump.

// Valid only within for/while, but we accept them here syntactically and reject them in the AST-to-CST
// conversion, where we can produce much more informative error messages:
md_statement(A) ::= MD_TOKEN_BREAK(O). {
	A = mlr_dsl_ast_node_alloc(O->text, MD_AST_NODE_TYPE_BREAK);
}
md_statement(A) ::= MD_TOKEN_CONTINUE(O). {
	A = mlr_dsl_ast_node_alloc(O->text, MD_AST_NODE_TYPE_BREAK);
}

// ================================================================
md_begin_block(A) ::= MD_TOKEN_BEGIN(O) MD_TOKEN_LBRACE md_statement_list(B) MD_TOKEN_RBRACE. {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_BEGIN, B);
}
md_end_block(A)   ::= MD_TOKEN_END(O)   MD_TOKEN_LBRACE md_statement_list(B) MD_TOKEN_RBRACE. {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_END, B);
}

// ----------------------------------------------------------------
md_cond_block(A) ::= md_rhs(B) MD_TOKEN_LBRACE md_statement_list(C) MD_TOKEN_RBRACE. {
	//A = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_alloc_binary("cond", MD_AST_NODE_TYPE_CONDITIONAL_BLOCK, B, C);
}

// ----------------------------------------------------------------
md_while_block(A) ::=
	MD_TOKEN_WHILE(O)
		MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN
		MD_TOKEN_LBRACE md_statement_list(C) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_WHILE, B, C);
}

// ----------------------------------------------------------------
md_do_while_block(A) ::=
	MD_TOKEN_DO(O)
		MD_TOKEN_LBRACE md_statement_list(B) MD_TOKEN_RBRACE
	MD_TOKEN_WHILE
		MD_TOKEN_LPAREN md_rhs(C) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DO_WHILE, B, C);
}

// ----------------------------------------------------------------
// for(k, v in $*) { ... }
md_for_loop_full_srec(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN MD_TOKEN_FULL_SREC
	MD_TOKEN_RPAREN
    MD_TOKEN_LBRACE
    	md_statement_list(S)
    MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_binary(
		F->text,
		MD_AST_NODE_TYPE_FOR_SREC,
		mlr_dsl_ast_node_alloc_binary(
			"variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			K,
			V
		),
		S
	);
}

// for(k, v in @*) { ... }
md_for_loop_full_oosvar(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN MD_TOKEN_FULL_OOSVAR
	MD_TOKEN_RPAREN
    MD_TOKEN_LBRACE
    	md_statement_list(S)
    MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR,
		mlr_dsl_ast_node_alloc_binary(
			"key_and_value_variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			mlr_dsl_ast_node_alloc_unary(
				"key_variables",
				MD_AST_NODE_TYPE_FOR_VARIABLES,
				K
			),
			V
		),
		mlr_dsl_ast_node_alloc_zary("empty_keylist", MD_AST_NODE_TYPE_OOSVAR_KEYLIST),
		S
	);
}

// for((k1, k2), v in @*) { ... }
// for((k1, k2, k3), v in @*) { ... }
md_for_loop_full_oosvar(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		MD_TOKEN_LPAREN md_for_oosvar_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN MD_TOKEN_FULL_OOSVAR
	MD_TOKEN_RPAREN
    MD_TOKEN_LBRACE
    	md_statement_list(S)
    MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR,
		mlr_dsl_ast_node_alloc_binary(
			"key_and_value_variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			L,
			V
		),
		mlr_dsl_ast_node_alloc_zary("empty_keylist", MD_AST_NODE_TYPE_OOSVAR_KEYLIST),
		S
	);
}

// for(k, v in @o[1][2]) { ... }
md_for_loop_oosvar(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_oosvar_keylist(O)
	MD_TOKEN_RPAREN
    MD_TOKEN_LBRACE
    	md_statement_list(S)
    MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR,
		mlr_dsl_ast_node_alloc_binary(
			"key_and_value_variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			mlr_dsl_ast_node_alloc_unary(
				"key_variables",
				MD_AST_NODE_TYPE_FOR_VARIABLES,
				K
			),
			V
		),
		O,
		S
	);
}

// for((k1, k2), v in @o[1][2]) { ... }
// for((k1, k2, k3), v in @o[1][2]) { ... }
md_for_loop_oosvar(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		MD_TOKEN_LPAREN md_for_oosvar_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_oosvar_keylist(O)
	MD_TOKEN_RPAREN
    MD_TOKEN_LBRACE
    	md_statement_list(S)
    MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR,
		mlr_dsl_ast_node_alloc_binary(
			"key_and_value_variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			L,
			V
		),
		O,
		S
	);
}

md_for_loop_index(A) ::= MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_NON_SIGIL_NAME);
}
md_for_oosvar_keylist(A) ::= MD_TOKEN_NON_SIGIL_NAME(K). {
	A = mlr_dsl_ast_node_alloc_unary("key_variables", MD_AST_NODE_TYPE_FOR_VARIABLES, K);
}
md_for_oosvar_keylist(A) ::= md_for_oosvar_keylist(L) MD_TOKEN_COMMA MD_TOKEN_NON_SIGIL_NAME(K). {
	A = mlr_dsl_ast_node_append_arg(L, K);
}

// ----------------------------------------------------------------
// Cases:
//   if elif*
//   if elif* else

md_if_chain(A) ::= md_if_elif_star(B) . {
	A = B;
}
md_if_chain(A) ::= md_if_elif_star(B) md_else_block(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
md_if_elif_star(A) ::= md_if_block(B). {
	A = mlr_dsl_ast_node_alloc_unary("if_head", MD_AST_NODE_TYPE_IF_HEAD, B);
}
md_if_elif_star(A) ::= md_if_elif_star(B) md_elif_block(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

md_if_block(A) ::=
	MD_TOKEN_IF(O)
		MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN
		MD_TOKEN_LBRACE md_statement_list(C) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_IF_ITEM, B, C);
}

md_elif_block(A) ::=
	MD_TOKEN_ELIF(O)
		MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN
		MD_TOKEN_LBRACE md_statement_list(C) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_IF_ITEM, B, C);
}

md_else_block(A) ::=
	MD_TOKEN_ELSE(O)
		MD_TOKEN_LBRACE md_statement_list(C) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_IF_ITEM, C);
}

// ----------------------------------------------------------------
md_bare_boolean(A) ::= md_rhs(B). {
	A = B;
}

// ----------------------------------------------------------------
md_filter(A) ::= MD_TOKEN_FILTER(O) md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_FILTER, B);
}

// ----------------------------------------------------------------
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B, C);
}
md_srec_indirect_assignment(A)  ::=
	MD_TOKEN_DOLLAR_SIGN MD_TOKEN_LEFT_BRACKET md_rhs(B) MD_TOKEN_RIGHT_BRACKET
	MD_TOKEN_ASSIGN(O) md_rhs(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT, B, C);
}

md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}

md_oosvar_from_full_srec_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT, B, C);
}

md_full_srec_from_oosvar_assignment(A)  ::= MD_TOKEN_FULL_SREC(B) MD_TOKEN_ASSIGN(O) md_oosvar_keylist(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT, B, C);
}

// ----------------------------------------------------------------
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_LOGICAL_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("||", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_LOGICAL_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_LOGICAL_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_BITWISE_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("|", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_BITWISE_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_BITWISE_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_BITWISE_LSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("<<", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_BITWISE_RSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(">>", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_PLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_MINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_DOT_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_TIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("/", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_INT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_MOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_POW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}


md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_LOGICAL_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("||", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_LOGICAL_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_LOGICAL_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_BITWISE_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("|", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_BITWISE_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_BITWISE_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_BITWISE_LSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("<<", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_BITWISE_RSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(">>", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_PLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_MINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_DOT_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_TIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("/", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_INT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_MOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_POW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

// ----------------------------------------------------------------
md_unset(A) ::= MD_TOKEN_UNSET(O) MD_TOKEN_ALL(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_UNSET, B);
}
md_unset(A) ::= MD_TOKEN_UNSET(O) MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_UNSET, B);
}
md_unset(A) ::= MD_TOKEN_UNSET(O) md_unset_args(B). {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "emit @a," -- use some non-empty-args expr.
md_unset_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("temp", MD_AST_NODE_TYPE_UNSET);
}

md_unset_args(A) ::= md_field_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= md_indirect_field_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}

md_unset_args(A) ::= md_unset_args(B) MD_TOKEN_COMMA md_field_name(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
md_unset_args(A) ::= md_unset_args(B) MD_TOKEN_COMMA md_oosvar_keylist(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
// Given "emitf @a,@b,@c": since this is a bottom-up parser, we get first the "@a",
// then "@a,@b", then "@a,@b,@c", then finally "emit @a,@b,@c". So:
// * On the "@a" we make a sub-AST called "temp @a" (although we could call it "emit").
// * On the "@b" we append the next argument to get "temp @a,@b".
// * On the "@c" we append the next argument to get "temp @a,@b,@c".
// * On the "emit" we change the name to get "emit @a,@b,@c".

md_emitf(A) ::= MD_TOKEN_EMITF(O) md_emitf_args(B). {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "emit @a," -- use some non-empty-args expr.
md_emitf_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("temp", MD_AST_NODE_TYPE_EMITF);
}
md_emitf_args(A) ::= md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMITF, B);
}
md_emitf_args(A) ::= md_emitf_args(B) MD_TOKEN_COMMA md_oosvar_keylist(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_emitp(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_ALL(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B);
}
md_emitp(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_ALL(B) MD_TOKEN_COMMA md_emitp_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emitp(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B);
}
md_emitp(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_COMMA md_emitp_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emitp(A) ::= MD_TOKEN_EMITP(O) md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B);
}
md_emitp(A) ::= MD_TOKEN_EMITP(O) md_oosvar_keylist(B) MD_TOKEN_COMMA md_emitp_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emitp_args(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMITP, B);
}
md_emitp_args(A) ::= md_emitp_args(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_emit(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_ALL(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B);
}
md_emit(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_ALL(B) MD_TOKEN_COMMA md_emit_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emit(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B);
}
md_emit(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_COMMA md_emit_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emit(A) ::= MD_TOKEN_EMIT(O) md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B);
}
md_emit(A) ::= MD_TOKEN_EMIT(O) md_oosvar_keylist(B) MD_TOKEN_COMMA md_emit_args(C). {
	B = mlr_dsl_ast_node_prepend_arg(C, B);
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}

md_emit_args(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMIT, B);
}
md_emit_args(A) ::= md_emit_args(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_dump(A) ::= MD_TOKEN_DUMP(O). {
	A = mlr_dsl_ast_node_alloc_zary(O->text, MD_AST_NODE_TYPE_DUMP);
}

// ================================================================
// Begin RHS precedence chain

md_rhs(A) ::= md_ternary(B). {
	A = B;
}

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

md_atom_or_fcn(A) ::= md_indirect_field_name(B). {
	A = B;
}
md_indirect_field_name(A) ::= MD_TOKEN_DOLLAR_SIGN MD_TOKEN_LEFT_BRACKET md_rhs(B) MD_TOKEN_RIGHT_BRACKET.  {
	A = mlr_dsl_ast_node_alloc_unary("indirect_field_name", MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME, B);
}

// ----------------------------------------------------------------
md_atom_or_fcn(A) ::= md_oosvar_keylist(B). {
	A = B;
}
md_oosvar_keylist(A) ::= md_oosvar_basename(B). {
	A = B;
}
md_oosvar_keylist(A) ::= md_oosvar_keylist(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// E.g. @name
md_oosvar_basename(A) ::= MD_TOKEN_UNBRACED_OOSVAR_NAME(B). {
	char* at_name = B->text;
	char* no_at_name = &at_name[1];
	A = mlr_dsl_ast_node_alloc_unary("oosvar_keylist", MD_AST_NODE_TYPE_OOSVAR_KEYLIST,
		mlr_dsl_ast_node_alloc(no_at_name, B->type));
}

// E.g. @{name}
md_oosvar_basename(A) ::= MD_TOKEN_BRACED_OOSVAR_NAME(B). {
	// Replace "@%{field.name}" with just "field.name"
	char* at_name = B->text;
	char* no_at_name = &at_name[2];
	int len = strlen(no_at_name);
	if (len > 0)
		no_at_name[len-1] = 0;
	A = mlr_dsl_ast_node_alloc_unary("oosvar_keylist", MD_AST_NODE_TYPE_OOSVAR_KEYLIST,
		mlr_dsl_ast_node_alloc(no_at_name, B->type));
}

// E.g. @["name"]
md_oosvar_basename(A) ::= MD_TOKEN_AT_SIGN MD_TOKEN_LEFT_BRACKET md_rhs(B) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_unary("oosvar_keylist", MD_AST_NODE_TYPE_OOSVAR_KEYLIST, B);
}

// ----------------------------------------------------------------
md_atom_or_fcn(A) ::= MD_TOKEN_NUMBER(B). {
	A = B;
}
md_atom_or_fcn(A) ::= MD_TOKEN_TRUE(B). {
	A = B;
}
md_atom_or_fcn(A) ::= MD_TOKEN_FALSE(B). {
	A = B;
}

md_atom_or_fcn(A) ::= md_string(B). {
	A = B;
}
md_atom_or_fcn(A) ::= md_regexi(B). {
	A = B;
}

md_atom_or_fcn(A) ::= md_bound_variable(B). {
	A = B;
}
md_bound_variable(A) ::= MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_BOUND_VARIABLE);
}

md_string(A) ::= MD_TOKEN_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(mlr_alloc_unbackslash(stripped), B->type);
}
md_regexi(A) ::= MD_TOKEN_REGEXI(B). {
	char* input = B->text;
	char* stripped = &input[1];
	int len = strlen(input);
	stripped[len-3] = 0;
	A = mlr_dsl_ast_node_alloc(mlr_alloc_unbackslash(stripped), B->type);
}

md_atom_or_fcn(A) ::= MD_TOKEN_CONTEXT_VARIABLE(B). {
	A = B;
}
md_atom_or_fcn(A) ::= MD_TOKEN_ENV(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_binary("env", MD_AST_NODE_TYPE_ENV, B, C);
}

md_atom_or_fcn(A) ::= MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

md_atom_or_fcn(A) ::= MD_TOKEN_NON_SIGIL_NAME(O) MD_TOKEN_LPAREN md_fcn_args(B) MD_TOKEN_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// Need to invalidate "f(10,)" -- use some non-empty-args expr.
md_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME);
}
md_fcn_args(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME, B);
}
md_fcn_args(A) ::= md_fcn_args(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
