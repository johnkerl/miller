// vim: set filetype=none:
// (Lemon files have .y extensions like Yacc files but are not Yacc.)

%include {
#include <stdio.h>
#include <string.h>
#include <math.h>
#include "../lib/mlrutil.h"
#include "../dsl/mlr_dsl_ast.h"
#include "../containers/sllv.h"

// ================================================================
// AST:
// * parens, commas, semis, line endings, whitespace are all stripped away
// * variable names and literal values remain as leaf nodes of the AST
// * = + - * / ** {function names} remain as non-leaf nodes of the AST
// CST: See mlr_dsl_cst.c
//
// Note: This parser accepts many things that are invalid, e.g.
//
// * begin{end{}} -- begin/end not at top level
// * begin{$x=1} -- references to stream records at begin/end
// * break/continue outside of for/while/do-while
// * return outside of a function definition
// * $x=x -- boundvars outside of for-loop variable bindings
//
// All of the above are enforced by the CST builder's semantic-analysis logic,
// which takes this parser's output AST as input.  This is done (a) to keep this
// grammar from being overly complex, and (b) so we can get much more
// informative error messages in C than in Lemon ('syntax error').
//
// The parser hooks all build up an abstract syntax tree specifically for the CST builder.
// For clearer visuals on what the ASTs look like:
// * See mlr_dsl_cst.c
// * See reg_test/run's filter -v and put -v outputs, e.g. in reg_test/expected/out
// * Do "mlr -n put -v 'your expression goes here'"
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
	fprintf(stderr, "mlr DSL: syntax error at \"%s\"\n", TOKEN->text);

//	This is confusing and (as is) worse than nothing.
//	Ideally we want to show the position within the input of the syntax error.
//
//	int n = sizeof(yyTokenName) / sizeof(yyTokenName[0]);
//	for (int i = 0; i < n; ++i) {
//			int a = yy_find_shift_action(pparser, (YYCODETYPE)i);
//			if (a < YYNSTATE + YYNRULE) {
//				fprintf(stderr, "Possible token \"%s\"\n", yyTokenName[i]);
//			}
//	}

}

// ================================================================
md_body ::= md_statement_block(B). {
	past->proot = B;
}

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

md_statement_block(A) ::= md_statement_braceless(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = mlr_dsl_ast_node_alloc_zary("block", MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("block", MD_AST_NODE_TYPE_STATEMENT_BLOCK, B);
	}
}

md_statement_block(A) ::= md_statement_braceful(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = mlr_dsl_ast_node_alloc_zary("block", MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("block", MD_AST_NODE_TYPE_STATEMENT_BLOCK, B);
	}
}

// This could also be done with the list on the left and statement on the right. However,
// curly-brace-terminated statements *end* with a semicolon; they don't start with one. So this seems
// to be the right way to differentiate.
md_statement_block(A) ::= md_statement_braceless(B) MD_TOKEN_SEMICOLON md_statement_block(C). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = C;
	} else {
		A = mlr_dsl_ast_node_prepend_arg(C, B);
	}
}

md_statement_block(A) ::= md_statement_braceful(B) md_statement_block(C). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = C;
	} else {
		A = mlr_dsl_ast_node_prepend_arg(C, B);
	}
}

// This allows for trailing semicolon, as well as empty string (or whitespace) between semicolons:
md_statement_braceless(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("nop", MD_AST_NODE_TYPE_NOP);
}

// Local-variable definitions at the current scope
md_statement_braceless(A) ::= md_untyped_local_definition(B).    { A = B; }
md_statement_braceless(A) ::= md_numeric_local_definition(B).    { A = B; }
md_statement_braceless(A) ::= md_int_local_definition(B).        { A = B; }
md_statement_braceless(A) ::= md_float_local_definition(B).      { A = B; }
md_statement_braceless(A) ::= md_boolean_local_definition(B).    { A = B; }
md_statement_braceless(A) ::= md_string_local_definition(B).     { A = B; }
md_statement_braceless(A) ::= md_map_local_definition(B).        { A = B; }
md_statement_braceless(A) ::= md_nonindexed_local_assignment(B). { A = B; }
md_statement_braceless(A) ::= md_indexed_local_assignment(B).    { A = B; }

// For user-defined functions
md_statement_braceless(A) ::= MD_TOKEN_RETURN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("return_value", MD_AST_NODE_TYPE_RETURN_VALUE, B);
}
md_statement_braceless(A) ::= MD_TOKEN_RETURN md_map_literal(B). {
	A = mlr_dsl_ast_node_alloc_unary("return_value", MD_AST_NODE_TYPE_RETURN_VALUE, B);
}
md_statement_braceless(A) ::= MD_TOKEN_RETURN MD_TOKEN_FULL_SREC(B). {
	A = mlr_dsl_ast_node_alloc_unary("return_value", MD_AST_NODE_TYPE_RETURN_VALUE, B);
}
md_statement_braceless(A) ::= MD_TOKEN_RETURN MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_unary("return_value", MD_AST_NODE_TYPE_RETURN_VALUE, B);
}
// For user-defined subroutines
md_statement_braceless(A) ::= MD_TOKEN_RETURN. {
	A = mlr_dsl_ast_node_alloc_zary("return_void", MD_AST_NODE_TYPE_RETURN_VOID);
}

// Begin/end
md_statement_braceful(A) ::= md_func_block(B).  { A = B; }
md_statement_braceful(A) ::= md_subr_block(B).  { A = B; }
md_statement_braceful(A) ::= md_begin_block(B). { A = B; }
md_statement_braceful(A) ::= md_end_block(B).   { A = B; }

// Nested control structures:
md_statement_braceful(A) ::= md_cond_block(B).                    { A = B; }
md_statement_braceful(A) ::= md_while_block(B).                   { A = B; }
md_statement_braceful(A) ::= md_for_loop_full_srec(B).            { A = B; }
md_statement_braceful(A) ::= md_for_loop_full_srec_key_only(B).   { A = B; }
md_statement_braceful(A) ::= md_for_loop_full_oosvar(B).          { A = B; }
md_statement_braceful(A) ::= md_for_loop_full_oosvar_key_only(B). { A = B; }
md_statement_braceful(A) ::= md_for_loop_oosvar(B).               { A = B; }
md_statement_braceful(A) ::= md_for_loop_oosvar_key_only(B).      { A = B; }
md_statement_braceful(A) ::= md_for_loop_local_map(B).            { A = B; }
md_statement_braceful(A) ::= md_for_loop_local_map_key_only(B).   { A = B; }
md_statement_braceful(A) ::= md_for_loop_map_literal(B).          { A = B; }
md_statement_braceful(A) ::= md_for_loop_map_literal_key_only(B). { A = B; }
md_statement_braceful(A) ::= md_for_loop_func_retval(B).            { A = B; }
md_statement_braceful(A) ::= md_for_loop_func_retval_key_only(B).   { A = B; }
md_statement_braceful(A) ::= md_triple_for(B).                    { A = B; }
md_statement_braceful(A) ::= md_if_chain(B).                      { A = B; }

md_statement_braceless(A) ::= MD_TOKEN_SUBR_CALL md_fcn_or_subr_call(B). {
	A = mlr_dsl_ast_node_alloc_unary("subr_call", MD_AST_NODE_TYPE_SUBR_CALLSITE, B);
}

// Not valid in begin/end since they refer to srecs:
md_statement_braceless(A) ::= md_srec_assignment(B).                  { A = B; }
md_statement_braceless(A) ::= md_srec_indirect_assignment(B).         { A = B; }
md_statement_braceless(A) ::= md_srec_positional_name_assignment(B).  { A = B; }
md_statement_braceless(A) ::= md_srec_positional_value_assignment(B). { A = B; }
md_statement_braceless(A) ::= md_oosvar_from_full_srec_assignment(B). { A = B; }
md_statement_braceless(A) ::= md_full_srec_assignment(B).             { A = B; }
md_statement_braceless(A) ::= md_env_assignment(B).                   { A = B; }

// Valid in begin/end since they don't refer to srecs (although the RHSs might):
md_statement_braceless(A) ::= md_do_while_block(B).         { A = B; }
md_statement_braceless(A) ::= md_bare_boolean(B).           { A = B; }
md_statement_braceless(A) ::= md_oosvar_assignment(B).      { A = B; }
md_statement_braceless(A) ::= md_full_oosvar_assignment(B). { A = B; }
md_statement_braceless(A) ::= md_filter(B).                 { A = B; }
md_statement_braceless(A) ::= md_unset(B).                  { A = B; }

md_statement_braceless(A) ::= md_tee_write(B).              { A = B; }
md_statement_braceless(A) ::= md_tee_append(B).             { A = B; }
md_statement_braceless(A) ::= md_tee_pipe(B).               { A = B; }
md_statement_braceless(A) ::= md_emitf(B).                  { A = B; }
md_statement_braceless(A) ::= md_emitf_write(B).            { A = B; }
md_statement_braceless(A) ::= md_emitf_append(B).           { A = B; }
md_statement_braceless(A) ::= md_emitf_pipe(B).             { A = B; }
md_statement_braceless(A) ::= md_emitp(B).                  { A = B; }
md_statement_braceless(A) ::= md_emitp_write(B).            { A = B; }
md_statement_braceless(A) ::= md_emitp_append(B).           { A = B; }
md_statement_braceless(A) ::= md_emitp_pipe(B).             { A = B; }
md_statement_braceless(A) ::= md_emit(B).                   { A = B; }
md_statement_braceless(A) ::= md_emit_write(B).             { A = B; }
md_statement_braceless(A) ::= md_emit_append(B).            { A = B; }
md_statement_braceless(A) ::= md_emit_pipe(B).              { A = B; }
md_statement_braceless(A) ::= md_emitp_lashed(B).           { A = B; }
md_statement_braceless(A) ::= md_emitp_lashed_write(B).     { A = B; }
md_statement_braceless(A) ::= md_emitp_lashed_append(B).    { A = B; }
md_statement_braceless(A) ::= md_emitp_lashed_pipe(B).      { A = B; }
md_statement_braceless(A) ::= md_emit_lashed(B).            { A = B; }
md_statement_braceless(A) ::= md_emit_lashed_write(B).      { A = B; }
md_statement_braceless(A) ::= md_emit_lashed_append(B).     { A = B; }
md_statement_braceless(A) ::= md_emit_lashed_pipe(B).       { A = B; }

md_statement_braceless(A) ::= md_dump(B).                   { A = B; }
md_statement_braceless(A) ::= md_dump_write(B).             { A = B; }
md_statement_braceless(A) ::= md_dump_append(B).            { A = B; }
md_statement_braceless(A) ::= md_dump_pipe(B).              { A = B; }
md_statement_braceless(A) ::= md_edump(B).                  { A = B; }
md_statement_braceless(A) ::= md_print(B).                  { A = B; }
md_statement_braceless(A) ::= md_eprint(B).                 { A = B; }
md_statement_braceless(A) ::= md_print_write(B).            { A = B; }
md_statement_braceless(A) ::= md_print_append(B).           { A = B; }
md_statement_braceless(A) ::= md_print_pipe(B).             { A = B; }
md_statement_braceless(A) ::= md_printn(B).                 { A = B; }
md_statement_braceless(A) ::= md_eprintn(B).                { A = B; }
md_statement_braceless(A) ::= md_printn_write(B).           { A = B; }
md_statement_braceless(A) ::= md_printn_append(B).          { A = B; }
md_statement_braceless(A) ::= md_printn_pipe(B).            { A = B; }

// Valid only within for/while, but we accept them here syntactically and reject them in the AST-to-CST
// conversion, where we can produce much more informative error messages:
md_statement_braceless(A) ::= MD_TOKEN_BREAK(O). {
	A = mlr_dsl_ast_node_alloc(O->text, MD_AST_NODE_TYPE_BREAK);
}
md_statement_braceless(A) ::= MD_TOKEN_CONTINUE(O). {
	A = mlr_dsl_ast_node_alloc(O->text, MD_AST_NODE_TYPE_CONTINUE);
}

// ================================================================
// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

md_func_block(C) ::= MD_TOKEN_FUNC_DEF
	MD_TOKEN_NON_SIGIL_NAME(F) MD_TOKEN_LPAREN md_func_or_subr_parameter_list(A) MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_set_function_name(A, F->text);
	mlr_dsl_ast_node_replace_text(B, "func_block");
	C = mlr_dsl_ast_node_alloc_binary(F->text, MD_AST_NODE_TYPE_FUNC_DEF, A, B);
}

md_func_block(C) ::= MD_TOKEN_FUNC_DEF
	MD_TOKEN_NON_SIGIL_NAME(F) MD_TOKEN_LPAREN md_func_or_subr_parameter_list(A) MD_TOKEN_RPAREN
	MD_TOKEN_COLON md_typedecl(M)
	MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_set_function_name(A, F->text);
	mlr_dsl_ast_node_replace_text(B, "func_block");
	C = mlr_dsl_ast_node_alloc_ternary(F->text, MD_AST_NODE_TYPE_FUNC_DEF, A, B, M);
}

md_subr_block(C) ::= MD_TOKEN_SUBR_DEF
	MD_TOKEN_NON_SIGIL_NAME(F) MD_TOKEN_LPAREN md_func_or_subr_parameter_list(A) MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE.
{
	A = mlr_dsl_ast_node_set_function_name(A, F->text);
	mlr_dsl_ast_node_replace_text(B, "subr_block");
	C = mlr_dsl_ast_node_alloc_binary(F->text, MD_AST_NODE_TYPE_SUBR_DEF, A, B);
}

md_func_or_subr_parameter_list(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME);
}
md_func_or_subr_parameter_list(A) ::= md_func_or_subr_non_empty_parameter_list(B). {
	A = B;
}
md_func_or_subr_non_empty_parameter_list(A) ::= md_func_or_subr_parameter(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME, B);
}
md_func_or_subr_non_empty_parameter_list(A) ::= md_func_or_subr_parameter(B) MD_TOKEN_COMMA. {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME, B);
}
md_func_or_subr_non_empty_parameter_list(A) ::= md_func_or_subr_parameter(B) MD_TOKEN_COMMA
	md_func_or_subr_non_empty_parameter_list(C).
{
	A = mlr_dsl_ast_node_prepend_arg(C, B);
}

md_func_or_subr_parameter(A) ::= MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_UNTYPED_PARAMETER_DEFINITION);
}
md_func_or_subr_parameter(A) ::= md_typedecl(T) MD_TOKEN_NON_SIGIL_NAME(N). {
	A = mlr_dsl_ast_node_alloc(N->text, T->type);
}

md_typedecl(A) ::= MD_TOKEN_VAR(B).     { A = B; A->type = MD_AST_NODE_TYPE_UNTYPED_PARAMETER_DEFINITION; }
md_typedecl(A) ::= MD_TOKEN_NUMERIC(B). { A = B; A->type = MD_AST_NODE_TYPE_NUMERIC_PARAMETER_DEFINITION; }
md_typedecl(A) ::= MD_TOKEN_INT(B).     { A = B; A->type = MD_AST_NODE_TYPE_INT_PARAMETER_DEFINITION;     }
md_typedecl(A) ::= MD_TOKEN_FLOAT(B).   { A = B; A->type = MD_AST_NODE_TYPE_FLOAT_PARAMETER_DEFINITION;   }
md_typedecl(A) ::= MD_TOKEN_STRING(B).  { A = B; A->type = MD_AST_NODE_TYPE_STRING_PARAMETER_DEFINITION;  }
md_typedecl(A) ::= MD_TOKEN_BOOLEAN(B). { A = B; A->type = MD_AST_NODE_TYPE_BOOLEAN_PARAMETER_DEFINITION; }
md_typedecl(A) ::= MD_TOKEN_MAP(B).     { A = B; A->type = MD_AST_NODE_TYPE_MAP_PARAMETER_DEFINITION;     }

// ================================================================
md_begin_block(A) ::= MD_TOKEN_BEGIN(O) MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE. {
	mlr_dsl_ast_node_replace_text(B, "begin_block");
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_BEGIN, B);
}
md_end_block(A)   ::= MD_TOKEN_END(O)   MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE. {
	mlr_dsl_ast_node_replace_text(B, "end_block");
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_END, B);
}

// ----------------------------------------------------------------
md_cond_block(A) ::= md_rhs(B) MD_TOKEN_LBRACE md_statement_block(C) MD_TOKEN_RBRACE. {
	mlr_dsl_ast_node_replace_text(C, "cond_block");
	A = mlr_dsl_ast_node_alloc_binary("cond", MD_AST_NODE_TYPE_CONDITIONAL_BLOCK, B, C);
}

// ----------------------------------------------------------------
md_while_block(A) ::=
	MD_TOKEN_WHILE(O)
		MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN
		MD_TOKEN_LBRACE md_statement_block(C) MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(C, "while_block");
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_WHILE, B, C);
}

// ----------------------------------------------------------------
md_do_while_block(A) ::=
	MD_TOKEN_DO(O)
		MD_TOKEN_LBRACE md_statement_block(B) MD_TOKEN_RBRACE
	MD_TOKEN_WHILE
		MD_TOKEN_LPAREN md_rhs(C) MD_TOKEN_RPAREN.
{
	mlr_dsl_ast_node_replace_text(B, "do_while_block");
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
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_full_srec_block");
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

// for(k in $*) { ... }
md_for_loop_full_srec_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_IN MD_TOKEN_FULL_SREC
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_full_srec_block");
	A = mlr_dsl_ast_node_alloc_binary(
		F->text,
		MD_AST_NODE_TYPE_FOR_SREC_KEY_ONLY,
		mlr_dsl_ast_node_alloc_unary(
			"variables",
			MD_AST_NODE_TYPE_FOR_VARIABLES,
			K
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
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_full_oosvar_block");
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
		MD_TOKEN_LPAREN md_for_map_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN MD_TOKEN_FULL_OOSVAR
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_full_oosvar_block");
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

// for(k in @*) { ... }
md_for_loop_full_oosvar_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K)
		MD_TOKEN_IN
		MD_TOKEN_FULL_OOSVAR
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_full_oosvar_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY,
		K,
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
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_oosvar_block");
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
		MD_TOKEN_LPAREN md_for_map_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_oosvar_keylist(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_oosvar_block");
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
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION);
}
md_for_loop_index(A) ::= MD_TOKEN_NUMERIC MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION);
}
md_for_loop_index(A) ::= MD_TOKEN_INT MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION);
}
md_for_loop_index(A) ::= MD_TOKEN_FLOAT MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION);
}
md_for_loop_index(A) ::= MD_TOKEN_STRING MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION);
}
md_for_loop_index(A) ::= MD_TOKEN_BOOLEAN MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION);
}

md_for_map_keylist(A) ::= md_for_loop_index(K). {
	A = mlr_dsl_ast_node_alloc_unary("key_variables", MD_AST_NODE_TYPE_FOR_VARIABLES, K);
}
md_for_map_keylist(A) ::= md_for_map_keylist(L) MD_TOKEN_COMMA md_for_loop_index(K). {
	A = mlr_dsl_ast_node_append_arg(L, K);
}

// for(k in @o[1][2]) { ... }
md_for_loop_oosvar_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K)
		MD_TOKEN_IN md_oosvar_keylist(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_oosvar_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY,
		K,
		O,
		S
	);
}

// ----------------------------------------------------------------
// for(k, v in o[1][2]) { ... }
md_for_loop_local_map(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_local_map_keylist(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_local_map_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_LOCAL_MAP,
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

// for((k1, k2), v in o[1][2]) { ... }
// for((k1, k2, k3), v in o[1][2]) { ... }
md_for_loop_local_map(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		MD_TOKEN_LPAREN md_for_map_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_local_map_keylist(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_local_map_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_LOCAL_MAP,
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

// for(k in o[1][2]) { ... }
md_for_loop_local_map_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K)
		MD_TOKEN_IN md_local_map_keylist(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_local_map_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_LOCAL_MAP_KEY_ONLY,
		K,
		O,
		S
	);
}

// ----------------------------------------------------------------
// for(k, v in o[1][2]) { ... }
md_for_loop_map_literal(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_map_literal(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_map_literal_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_MAP_LITERAL,
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

// for((k1, k2), v in o[1][2]) { ... }
// for((k1, k2, k3), v in o[1][2]) { ... }
md_for_loop_map_literal(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		MD_TOKEN_LPAREN md_for_map_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_map_literal(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_map_literal_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_MAP_LITERAL,
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

// for(k in o[1][2]) { ... }
md_for_loop_map_literal_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K)
		MD_TOKEN_IN md_map_literal(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_map_literal_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_MAP_LITERAL_KEY_ONLY,
		K,
		O,
		S
	);
}

// ----------------------------------------------------------------
// for(k, v in o[1][2]) { ... }
md_for_loop_func_retval(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K) MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_fcn_or_subr_call(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_func_retval_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_FUNC_RETVAL,
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

// for((k1, k2), v in o[1][2]) { ... }
// for((k1, k2, k3), v in o[1][2]) { ... }
md_for_loop_func_retval(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		MD_TOKEN_LPAREN md_for_map_keylist(L) MD_TOKEN_RPAREN MD_TOKEN_COMMA md_for_loop_index(V)
		MD_TOKEN_IN md_fcn_or_subr_call(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_func_retval_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_FUNC_RETVAL,
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

// for(k in o[1][2]) { ... }
md_for_loop_func_retval_key_only(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_for_loop_index(K)
		MD_TOKEN_IN md_fcn_or_subr_call(O)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(S)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(S, "for_loop_func_retval_block");
	A = mlr_dsl_ast_node_alloc_ternary(
		F->text,
		MD_AST_NODE_TYPE_FOR_FUNC_RETVAL_KEY_ONLY,
		K,
		O,
		S
	);
}

// ----------------------------------------------------------------
md_triple_for(A) ::=
	MD_TOKEN_FOR(F) MD_TOKEN_LPAREN
		md_triple_for_start(S)
			MD_TOKEN_SEMICOLON
		md_triple_for_continuation(C)
			MD_TOKEN_SEMICOLON
		md_triple_for_update(U)
	MD_TOKEN_RPAREN
	MD_TOKEN_LBRACE
		md_statement_block(L)
	MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(L, "triple_for_block");
	A = mlr_dsl_ast_node_alloc_quaternary(F->text, MD_AST_NODE_TYPE_TRIPLE_FOR, S, C, U, L);
}

md_triple_for_start(A) ::= md_statement_braceless(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = mlr_dsl_ast_node_alloc_zary("triple_for_start_statements", MD_AST_NODE_TYPE_STATEMENT_LIST);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("triple_for_start_statements", MD_AST_NODE_TYPE_STATEMENT_LIST, B);
	}
}
md_triple_for_start(A) ::= md_triple_for_start(B) MD_TOKEN_COMMA md_statement_braceless(C). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = C;
	} else {
		A = mlr_dsl_ast_node_append_arg(B, C);
	}
}

md_triple_for_continuation(A) ::= md_statement_braceless(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = mlr_dsl_ast_node_alloc_zary("triple_for_continuation_statements", MD_AST_NODE_TYPE_STATEMENT_LIST);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("triple_for_continuation_statements", MD_AST_NODE_TYPE_STATEMENT_LIST, B);
	}
}
md_triple_for_continuation(A) ::= md_triple_for_continuation(B) MD_TOKEN_COMMA md_statement_braceless(C). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = C;
	} else {
		A = mlr_dsl_ast_node_append_arg(B, C);
	}
}

md_triple_for_update(A) ::= md_statement_braceless(B). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = mlr_dsl_ast_node_alloc_zary("triple_for_update_statements", MD_AST_NODE_TYPE_STATEMENT_LIST);
	} else {
		A = mlr_dsl_ast_node_alloc_unary("triple_for_update_statements", MD_AST_NODE_TYPE_STATEMENT_LIST, B);
	}
}
md_triple_for_update(A) ::= md_triple_for_update(B) MD_TOKEN_COMMA md_statement_braceless(C). {
	if (B->type == MD_AST_NODE_TYPE_NOP) {
		mlr_dsl_ast_node_free(B);
		A = C;
	} else {
		A = mlr_dsl_ast_node_append_arg(B, C);
	}
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
		MD_TOKEN_LBRACE md_statement_block(C) MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(C, "if_block");
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_IF_ITEM, B, C);
}

md_elif_block(A) ::=
	MD_TOKEN_ELIF(O)
		MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN
		MD_TOKEN_LBRACE md_statement_block(C) MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(C, "elif_block");
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_IF_ITEM, B, C);
}

md_else_block(A) ::=
	MD_TOKEN_ELSE(O)
		MD_TOKEN_LBRACE md_statement_block(C) MD_TOKEN_RBRACE.
{
	mlr_dsl_ast_node_replace_text(C, "else_block");
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
md_untyped_local_definition(A) ::= MD_TOKEN_VAR(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION, N, B);
}
md_untyped_local_definition(A) ::= MD_TOKEN_VAR(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN MD_TOKEN_FULL_SREC(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION, N, B);
}
md_untyped_local_definition(A) ::= MD_TOKEN_VAR(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION, N, B);
}
md_untyped_local_definition(A) ::= MD_TOKEN_VAR(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION, N, C);
}

md_numeric_local_definition(A) ::= MD_TOKEN_NUMERIC(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION, N, B);
}
md_int_local_definition(A) ::= MD_TOKEN_INT(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION, N, B);
}
md_float_local_definition(A) ::= MD_TOKEN_FLOAT(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION, N, B);
}
md_boolean_local_definition(A) ::= MD_TOKEN_BOOLEAN(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION, N, B);
}
md_string_local_definition(A) ::= MD_TOKEN_STRING(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION, N, B);
}

md_map_local_definition(A) ::= MD_TOKEN_MAP(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION, N, C);
}
md_map_local_definition(A) ::= MD_TOKEN_MAP(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION, N, C);
}
md_map_local_definition(A) ::= MD_TOKEN_MAP(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION, N, C);
}
md_map_local_definition(A) ::= MD_TOKEN_MAP(T) md_nonindexed_local_variable(N) MD_TOKEN_ASSIGN md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(T->text, MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION, N, C);
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B, C);
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B, C);
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_ASSIGN(O) md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B, C);
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_ASSIGN(O) md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B, C);
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B, C);
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B, C);
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B, C);
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B, C);
}

// ----------------------------------------------------------------
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B, C);
}
md_srec_indirect_assignment(A)  ::=
	MD_TOKEN_DOLLAR_SIGN MD_TOKEN_LEFT_BRACKET md_rhs(B) MD_TOKEN_RIGHT_BRACKET
	MD_TOKEN_ASSIGN(O) md_rhs(C).
{
	A = mlr_dsl_ast_node_alloc_binary(
		O->text,
		MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT,
		B,
		C
	);
}
md_srec_positional_name_assignment(A)  ::=
	MD_TOKEN_DOLLAR_SIGN
		MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET
			md_rhs(B)
		MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET
	MD_TOKEN_ASSIGN(O) md_rhs(C).
{
	A = mlr_dsl_ast_node_alloc_binary(
		O->text,
		MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME_ASSIGNMENT,
		B,
		C
	);
}
// '$[[[3]]] = "new"' is shorthand for '$[ $[[3]] ] = "new"'.
// Note that '$[[3]]' is key at srec position 3 and '$[[[3]]]' is value at srec position 3.
md_srec_positional_value_assignment(A)  ::=
	MD_TOKEN_DOLLAR_SIGN
		MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET
			md_rhs(B)
		MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET
	MD_TOKEN_ASSIGN(O) md_rhs(C).
{
	A = mlr_dsl_ast_node_alloc_binary(
		O->text,
		MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT,
		mlr_dsl_ast_node_alloc_unary(
			"positional_srec_field_name",
			MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME,
			B
		),
		C
	);
}

md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}

md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}

md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B, C);
}

md_full_oosvar_assignment(A)  ::= MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT, B, C);
}

md_full_oosvar_assignment(A)  ::= MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_ASSIGN(O) md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT, B, C);
}

md_full_oosvar_assignment(A)  ::= MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT, B, C);
}

md_full_oosvar_assignment(A)  ::= MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT, B, C);
}

md_full_oosvar_assignment(A)  ::= MD_TOKEN_FULL_OOSVAR(B) MD_TOKEN_ASSIGN(O) md_fcn_or_subr_call(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT, B, C);
}

md_oosvar_from_full_srec_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT, B, C);
}

md_full_srec_assignment(A)  ::= MD_TOKEN_FULL_SREC(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text,
		MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT, B, C);
}
md_full_srec_assignment(A)  ::= MD_TOKEN_FULL_SREC(B) MD_TOKEN_ASSIGN(O) md_map_literal(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT, B, C);
}
md_full_srec_assignment(A)  ::= MD_TOKEN_FULL_SREC(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_SREC(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT, B, C);
}
md_full_srec_assignment(A)  ::= MD_TOKEN_FULL_SREC(B) MD_TOKEN_ASSIGN(O) MD_TOKEN_FULL_OOSVAR(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT, B, C);
}

md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_ASSIGN(O) md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B, C);
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

md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_OPLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_OMINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_OTIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("./", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_srec_assignment(A)  ::= md_field_name(B) MD_TOKEN_INT_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_SREC_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
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

md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_OPLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_OMINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_OTIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("./", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_oosvar_assignment(A)  ::= md_oosvar_keylist(B) MD_TOKEN_INT_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_LOGICAL_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("||", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_LOGICAL_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_LOGICAL_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_BITWISE_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("|", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_BITWISE_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_BITWISE_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_BITWISE_LSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("<<", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_BITWISE_RSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(">>", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_PLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_MINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_DOT_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_TIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("/", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_INT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_MOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_POW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_OPLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_OMINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_OTIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("./", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_nonindexed_local_assignment(A)  ::= md_nonindexed_local_variable(B) MD_TOKEN_INT_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_LOGICAL_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("||", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_LOGICAL_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_LOGICAL_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_BITWISE_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("|", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_BITWISE_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_BITWISE_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_BITWISE_LSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("<<", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_BITWISE_RSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(">>", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_PLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_MINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_DOT_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_TIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("/", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_INT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_MOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_POW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_OPLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_OMINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_OTIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("./", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_indexed_local_assignment(A)  ::= md_indexed_local_variable(B) MD_TOKEN_INT_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

// ----------------------------------------------------------------
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_LOGICAL_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("||", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_LOGICAL_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_LOGICAL_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_BITWISE_OR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("|", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_BITWISE_XOR_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("^", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_BITWISE_AND_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("&", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_BITWISE_LSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("<<", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_BITWISE_RSH_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(">>", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_PLUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("+", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_MINUS_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("-", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_DOT_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_TIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("/", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_INT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_MOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_POW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}

md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_OTIMES_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".*", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_ODIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary("./", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_OINT_DIVIDE_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".//", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_OMOD_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".%", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
}
md_env_assignment(A)  ::= md_env_index(B) MD_TOKEN_OPOW_EQUALS md_rhs(C). {
	A = mlr_dsl_ast_node_alloc_binary("=", MD_AST_NODE_TYPE_ENV_ASSIGNMENT, B,
		mlr_dsl_ast_node_alloc_binary(".**", MD_AST_NODE_TYPE_OPERATOR, mlr_dsl_ast_tree_copy(B) , C));
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
md_unset_args(A) ::= md_positional_srec_name(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= MD_TOKEN_FULL_SREC(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= md_nonindexed_local_variable(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}
md_unset_args(A) ::= md_indexed_local_variable(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_UNSET, B);
}

md_unset_args(A) ::= md_unset_args(B) MD_TOKEN_COMMA md_field_name(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
md_unset_args(A) ::= md_unset_args(B) MD_TOKEN_COMMA md_oosvar_keylist(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_tee_write(A) ::= MD_TOKEN_TEE(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA MD_TOKEN_FULL_SREC(M). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_TEE,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F),
		M);
}

md_tee_append(A) ::= MD_TOKEN_TEE(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA MD_TOKEN_FULL_SREC(M). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_TEE,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F),
		M);
}

md_tee_pipe(A) ::= MD_TOKEN_TEE(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA MD_TOKEN_FULL_SREC(M). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_TEE,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P),
		M);
}

// ----------------------------------------------------------------
// Given "emitf @a,@b,@c": since this is a bottom-up parser, we get first the "@a",
// then "@a,@b", then "@a,@b,@c", then finally "emit @a,@b,@c". So:
// * On the "@a" we make a sub-AST called "temp @a" (although we could call it "emit").
// * On the "@b" we append the next argument to get "temp @a,@b".
// * On the "@c" we append the next argument to get "temp @a,@b,@c".
// * On the "emit" we change the name to get "emit @a,@b,@c".

md_emitf(A) ::= MD_TOKEN_EMITF(O) md_emitf_args(B). {
	B = mlr_dsl_ast_node_set_function_name(B, O->text);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITF, B,
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}
// Need to invalidate "emit @a," -- use some non-empty-args expr.
md_emitf_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("temp", MD_AST_NODE_TYPE_EMITF);
}
md_emitf_args(A) ::= md_oosvar_keylist(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMITF, B);
}
md_emitf_args(A) ::= md_nonindexed_local_variable(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMITF, B);
}
md_emitf_args(A) ::= md_indexed_local_variable(B). {
	A = mlr_dsl_ast_node_alloc_unary("temp", MD_AST_NODE_TYPE_EMITF, B);
}
md_emitf_args(A) ::= md_emitf_args(B) MD_TOKEN_COMMA md_oosvar_keylist(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
md_emitf_args(A) ::= md_emitf_args(B) MD_TOKEN_COMMA md_nonindexed_local_variable(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
md_emitf_args(A) ::= md_emitf_args(B) MD_TOKEN_COMMA md_indexed_local_variable(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

md_emitf_write(A) ::= MD_TOKEN_EMITF(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA md_emitf_args(B). {
	B = mlr_dsl_ast_node_set_function_name(B, O->text);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITF, B,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}

md_emitf_append(A) ::= MD_TOKEN_EMITF(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA md_emitf_args(B). {
	B = mlr_dsl_ast_node_set_function_name(B, O->text);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITF, B,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}

md_emitf_pipe(A) ::= MD_TOKEN_EMITF(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA md_emitf_args(B). {
	B = mlr_dsl_ast_node_set_function_name(B, O->text);
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITF, B,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}


// ----------------------------------------------------------------
md_emitp(A) ::= MD_TOKEN_EMITP(O) md_emittable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}

md_emitp(A) ::= MD_TOKEN_EMITP(O) md_emittable(B) MD_TOKEN_COMMA md_emitp_namelist(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP, B, C),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}

md_emitp_write(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
		F));
}
md_emitp_write(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emitp_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP, B, C),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
		F));
}

md_emitp_append(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
		F));
}
md_emitp_append(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emitp_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP, B, C),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
		F));
}

md_emitp_pipe(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP, B),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
		P));
}
md_emitp_pipe(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emitp_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP, B, C),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
		P));
}

// ----------------------------------------------------------------
md_emitp_namelist(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("emitp_namelist", MD_AST_NODE_TYPE_EMITP, B);
}
md_emitp_namelist(A) ::= md_emitp_namelist(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ================================================================
md_emit(A) ::= MD_TOKEN_EMIT(O) md_emittable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}

md_emit(A) ::= MD_TOKEN_EMIT(O) md_emittable(B) MD_TOKEN_COMMA md_emit_namelist(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT, B, C),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}


md_emit_write(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
		F));
}
md_emit_write(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emit_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT, B, C),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
		F));
}

md_emit_append(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
		F));
}
md_emit_append(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emit_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT, B, C),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
		F));
}

md_emit_pipe(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	md_emittable(B).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT, B),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
		P));
}
md_emit_pipe(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	md_emittable(B) MD_TOKEN_COMMA md_emit_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT, B, C),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
		P));
}

// ----------------------------------------------------------------
md_emittable(A) ::= MD_TOKEN_ALL(B). {
	A = B;
}
md_emittable(A) ::= MD_TOKEN_FULL_OOSVAR(B). {
	A = B;
}
md_emittable(A) ::= md_oosvar_keylist(B). {
	A = B;
}
md_emittable(A) ::= md_nonindexed_local_variable(B). {
	A = B;
}
md_emittable(A) ::= md_indexed_local_variable(B). {
	A = B;
}
md_emittable(A) ::= MD_TOKEN_FULL_SREC(B). {
	A = B;
}
md_emittable(A) ::= md_map_literal(B). {
	A = B;
}
md_emittable(A) ::= md_fcn_or_subr_call(B). {
	A = B;
}

// ----------------------------------------------------------------
md_emit_namelist(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("emit_namelist", MD_AST_NODE_TYPE_EMIT, B);
}
md_emit_namelist(A) ::= md_emit_namelist(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_emitp_lashed(A) ::= MD_TOKEN_EMITP(O)
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}
md_emitp_lashed(A) ::= MD_TOKEN_EMITP(O)
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emitp_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B, C),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}

md_emitp_lashed_write(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}
md_emitp_lashed_write(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emitp_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}

md_emitp_lashed_append(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}
md_emitp_lashed_append(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emitp_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}

md_emitp_lashed_pipe(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}
md_emitp_lashed_pipe(A) ::= MD_TOKEN_EMITP(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emitp_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emitp_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMITP_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
md_emitp_lashed_keylists(A) ::= md_emittable(B). {
	A = mlr_dsl_ast_node_alloc_unary("lashed_keylists", MD_AST_NODE_TYPE_EMITP_LASHED, B);
}
md_emitp_lashed_keylists(A) ::= md_emitp_lashed_keylists(B) MD_TOKEN_COMMA md_emittable(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
md_emitp_lashed_namelist(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("lashed_namelist", MD_AST_NODE_TYPE_EMITP_LASHED, B);
}
md_emitp_lashed_namelist(A) ::= md_emitp_lashed_namelist(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_emit_lashed(A) ::= MD_TOKEN_EMIT(O)
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}
md_emit_lashed(A) ::= MD_TOKEN_EMIT(O)
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emit_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B, C),
		mlr_dsl_ast_node_alloc_zary("stream", MD_AST_NODE_TYPE_STREAM));
}

md_emit_lashed_write(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}
md_emit_lashed_write(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emit_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}

md_emit_lashed_append(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}
md_emit_lashed_append(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emit_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}

md_emit_lashed_pipe(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN.
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}
md_emit_lashed_pipe(A) ::= MD_TOKEN_EMIT(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA
	MD_TOKEN_LPAREN md_emit_lashed_keylists(B) MD_TOKEN_RPAREN
	MD_TOKEN_COMMA md_emit_lashed_namelist(C).
{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED,
		mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_EMIT_LASHED, B, C),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
md_emit_lashed_keylists(A) ::= md_emittable(B). {
	A = mlr_dsl_ast_node_alloc_unary("lashed_keylists", MD_AST_NODE_TYPE_EMIT_LASHED, B);
}
md_emit_lashed_keylists(A) ::= md_emit_lashed_keylists(B) MD_TOKEN_COMMA md_emittable(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
md_emit_lashed_namelist(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("lashed_namelist", MD_AST_NODE_TYPE_EMIT_LASHED, B);
}
md_emit_lashed_namelist(A) ::= md_emit_lashed_namelist(B) MD_TOKEN_COMMA md_rhs(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}

// ----------------------------------------------------------------
md_dump(A) ::= MD_TOKEN_DUMP(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)),
		mlr_dsl_ast_node_alloc("all", MD_AST_NODE_TYPE_FULL_OOSVAR));
}
md_edump(A) ::= MD_TOKEN_EDUMP(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)),
		mlr_dsl_ast_node_alloc("all", MD_AST_NODE_TYPE_FULL_OOSVAR));
}
md_dump_write(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_GT md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			F),
		mlr_dsl_ast_node_alloc("all", MD_AST_NODE_TYPE_FULL_OOSVAR));
}
md_dump_append(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_BITWISE_RSH md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
			F),
		mlr_dsl_ast_node_alloc("all", MD_AST_NODE_TYPE_FULL_OOSVAR));
}
md_dump_pipe(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_BITWISE_OR md_rhs(P). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
			P),
		mlr_dsl_ast_node_alloc("all", MD_AST_NODE_TYPE_FULL_OOSVAR));
}

md_dump(A) ::= MD_TOKEN_DUMP(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)),
		B);
}
md_edump(A) ::= MD_TOKEN_EDUMP(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)),
		B);
}
md_dump_write(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_GT md_output_file(F) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			F),
		B);
}
md_dump_append(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_BITWISE_RSH md_output_file(F) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
			F),
		B);
}
md_dump_pipe(A) ::= MD_TOKEN_DUMP(O) MD_TOKEN_BITWISE_OR md_dumpable(P) md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_DUMP,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
			P),
		B);
}

// ----------------------------------------------------------------
md_dumpable(A) ::= MD_TOKEN_ALL(B). {
	A = B;
}
md_dumpable(A) ::= MD_TOKEN_FULL_OOSVAR(B). {
	A = B;
}
md_dumpable(A) ::= MD_TOKEN_FULL_SREC(B). {
	A = B;
}
md_dumpable(A) ::= md_map_literal(B). {
	A = B;
}
md_dumpable(A) ::= md_rhs(B). {
	A = B;
}

// ----------------------------------------------------------------
// Print string
md_print(A) ::= MD_TOKEN_PRINT(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT, B,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)));
}
md_eprint(A) ::= MD_TOKEN_EPRINT(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT, B,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)));
}
md_print_write(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT, C,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			F));
}
md_print_append(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT, C,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
			F));
}
md_print_pipe(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT, C,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
			P));
}

// Print with no string (newline only)
md_print(A) ::= MD_TOKEN_PRINT(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)));
}
md_eprint(A) ::= MD_TOKEN_EPRINT(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)));
}
md_print_write(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_GT md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}
md_print_append(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_BITWISE_RSH md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}
md_print_pipe(A) ::= MD_TOKEN_PRINT(O) MD_TOKEN_BITWISE_OR md_rhs(P). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINT,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}

// Printn string
md_printn(A) ::= MD_TOKEN_PRINTN(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN, B,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)));
}
md_eprintn(A) ::= MD_TOKEN_EPRINTN(O) md_dumpable(B). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN, B,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)));
}
md_printn_write(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_GT md_output_file(F) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN, C,
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			F));
}
md_printn_append(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_BITWISE_RSH md_output_file(F) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN, C,
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND,
			F));
}
md_printn_pipe(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_BITWISE_OR md_rhs(P) MD_TOKEN_COMMA md_dumpable(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN, C,
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE,
			P));
}

// Printn with no string: produces no output but will create a zero-length
// output file, so not quite a no-op.
md_printn(A) ::= MD_TOKEN_PRINTN(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDOUT)));
}
md_eprintn(A) ::= MD_TOKEN_EPRINTN(O). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE,
			mlr_dsl_ast_node_alloc_zary("stdout", MD_AST_NODE_TYPE_STDERR)));
}
md_printn_write(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_GT md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">", MD_AST_NODE_TYPE_FILE_WRITE, F));
}
md_printn_append(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_BITWISE_RSH md_output_file(F). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary(">>", MD_AST_NODE_TYPE_FILE_APPEND, F));
}
md_printn_pipe(A) ::= MD_TOKEN_PRINTN(O) MD_TOKEN_BITWISE_OR md_dumpable(P). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_PRINTN,
		mlr_dsl_ast_node_alloc("", MD_AST_NODE_TYPE_NUMERIC_LITERAL),
		mlr_dsl_ast_node_alloc_unary("|", MD_AST_NODE_TYPE_PIPE, P));
}

// ----------------------------------------------------------------
md_output_file(A) ::= md_rhs(F).          { A = F; }
md_output_file(A) ::= MD_TOKEN_STDOUT(F). { A = F; }
md_output_file(A) ::= MD_TOKEN_STDERR(F). { A = F; }

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
md_addsubdot_term(A) ::= md_addsubdot_term(B) MD_TOKEN_OPLUS(O) md_muldiv_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_addsubdot_term(A) ::= md_addsubdot_term(B) MD_TOKEN_OMINUS(O) md_muldiv_term(C). {
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

md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_OTIMES(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_ODIVIDE(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, B, C);
}
md_muldiv_term(A) ::= md_muldiv_term(B) MD_TOKEN_INT_ODIVIDE(O) md_unary_bitwise_op_term(C). {
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
md_unary_bitwise_op_term(A) ::= MD_TOKEN_OPLUS(O) md_unary_bitwise_op_term(C). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MD_AST_NODE_TYPE_OPERATOR, C);
}
md_unary_bitwise_op_term(A) ::= MD_TOKEN_OMINUS(O) md_unary_bitwise_op_term(C). {
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

md_atom_or_fcn(A) ::= md_positional_srec_name(B). {
	A = B;
}
md_positional_srec_name(A) ::= MD_TOKEN_DOLLAR_SIGN
	MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET
		md_rhs(B)
	MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET.  {
	A = mlr_dsl_ast_node_alloc_unary(
		"positional_srec_field_name",
		MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME,
		B
	);
}

// '$value = $[[[3]]]' is shorthand for '$value = $[ $[[3]] ]'.
// Note that '$[[3]]' is key at srec position 3 and '$[[[3]]]' is value at srec position 3.
md_atom_or_fcn(A) ::= md_positional_srec_value(B). {
	A = B;
}
md_positional_srec_value(A) ::= MD_TOKEN_DOLLAR_SIGN
	MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET MD_TOKEN_LEFT_BRACKET
		md_rhs(B)
	MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET MD_TOKEN_RIGHT_BRACKET.  {
	A = mlr_dsl_ast_node_alloc_unary(
		"indirect_field_name",
		MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME,
		mlr_dsl_ast_node_alloc_unary(
			"positional_srec_field_name",
			MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME,
			B
		)
	);
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

md_local_map_keylist(A) ::= md_nonindexed_local_variable(B). {
	A = B;
}
md_local_map_keylist(A) ::= md_local_map_keylist(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_append_arg(B, C);
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

md_atom_or_fcn(A) ::= md_nonindexed_local_variable(B). {
	A = B;
}

md_atom_or_fcn(A) ::= md_indexed_local_variable(B). {
	A = B;
}

md_nonindexed_local_variable(A) ::= MD_TOKEN_NON_SIGIL_NAME(B). {
	A = mlr_dsl_ast_node_alloc(B->text, MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE);
}

md_indexed_local_variable(A) ::= md_nonindexed_local_variable(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_unary(B->text, MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE, C);
	mlr_dsl_ast_node_free(B);
}
md_indexed_local_variable(A) ::= md_indexed_local_variable(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_append_arg(B, C);
}


md_atom_or_fcn(A) ::= md_indexed_function_call(B). {
	A = B;
}

// Indexed function calls:
// * First child node is function call with argument list.
// * Second child node is list of indexing expressions.
//
// Example: 'foo(1,2,3)[4][5]' parses to AST
//
// text="foo", type=INDEXED_FUNCTION_CALLSITE:
//     text="foo", type=FUNCTION_CALLSITE:
//         text="1", type=NUMERIC_LITERAL.
//         text="2", type=NUMERIC_LITERAL.
//         text="3", type=NUMERIC_LITERAL.
//     text="indexing", type=MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST:
//         text="4", type=NUMERIC_LITERAL.
//         text="5", type=NUMERIC_LITERAL.
md_indexed_function_call(A) ::= md_fcn_or_subr_call(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_binary(
		B->text,
		MD_AST_NODE_TYPE_INDEXED_FUNCTION_CALLSITE,
		B,
		mlr_dsl_ast_node_alloc_unary(
			"indexing",
			MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST,
			C
		)
	);
}
md_indexed_function_call(A) ::= md_indexed_function_call(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	// Append to second child node which is list of indexing expressions.
	A = mlr_dsl_ast_node_append_arg_to_second_child(B, C);
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
md_atom_or_fcn(A) ::= md_env_index(B). {
	A = B;
}

md_env_index(A) ::= MD_TOKEN_ENV(B) MD_TOKEN_LEFT_BRACKET md_rhs(C) MD_TOKEN_RIGHT_BRACKET. {
	A = mlr_dsl_ast_node_alloc_binary("env", MD_AST_NODE_TYPE_ENV, B, C);
}

md_atom_or_fcn(A) ::= MD_TOKEN_LPAREN md_rhs(B) MD_TOKEN_RPAREN. {
	A = B;
}

md_atom_or_fcn(A) ::= md_fcn_or_subr_call(B). {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

md_fcn_or_subr_call(A) ::= MD_TOKEN_NON_SIGIL_NAME(O) MD_TOKEN_LPAREN md_fcn_arg_list(B) MD_TOKEN_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
	A->type = MD_AST_NODE_TYPE_FUNCTION_CALLSITE;
}
// For most functions it suffices to use the MD_TOKEN_NON_SIGIL_NAME pattern. But
// int and float are keywords in the lexer so we need to spell those out explicitly.
// (They're type-decl keywords but they're also the names of type-conversion functions.)
md_fcn_or_subr_call(A) ::= MD_TOKEN_INT(O) MD_TOKEN_LPAREN md_fcn_arg_list(B) MD_TOKEN_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
	A->type = MD_AST_NODE_TYPE_FUNCTION_CALLSITE;
}
md_fcn_or_subr_call(A) ::= MD_TOKEN_FLOAT(O) MD_TOKEN_LPAREN md_fcn_arg_list(B) MD_TOKEN_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
	A->type = MD_AST_NODE_TYPE_FUNCTION_CALLSITE;
}

md_fcn_arg_list(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME);
}
md_fcn_arg_list(A) ::= md_fcn_non_empty_arg_list(B). {
	A = B;
}

md_fcn_non_empty_arg_list(A) ::= md_fcn_arg(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME, B);
}
md_fcn_non_empty_arg_list(A) ::= md_fcn_arg(B) MD_TOKEN_COMMA. {
	A = mlr_dsl_ast_node_alloc_unary("anon", MD_AST_NODE_TYPE_NON_SIGIL_NAME, B);
}
md_fcn_non_empty_arg_list(A) ::= md_fcn_arg(B) MD_TOKEN_COMMA md_fcn_non_empty_arg_list(C). {
	A = mlr_dsl_ast_node_prepend_arg(C, B);
}

md_fcn_arg(A) ::= md_rhs(B). {
	A = B;
}
md_fcn_arg(A) ::= MD_TOKEN_FULL_SREC(B). {
	A = B;
}
md_fcn_arg(A) ::= MD_TOKEN_FULL_OOSVAR(B). {
	A = B;
}
md_fcn_arg(A) ::= md_map_literal(B). {
	A = B;
}

// ----------------------------------------------------------------
// Map-literals in Miller are JSON-ish.

md_map_literal(A) ::= MD_TOKEN_LBRACE MD_TOKEN_RBRACE. {
	A = mlr_dsl_ast_node_alloc_zary("map_literal", MD_AST_NODE_TYPE_MAP_LITERAL);
}
md_map_literal(A) ::= MD_TOKEN_LBRACE md_map_literal_kv_pairs(B) MD_TOKEN_RBRACE. {
	A = B;
}
md_map_literal_kv_pairs(A) ::= md_map_literal_kv_pair(B). {
	A = mlr_dsl_ast_node_alloc_unary("map_literal", MD_AST_NODE_TYPE_MAP_LITERAL, B);
}
// Allow trailing final comma, especially for multiline map literals
md_map_literal_kv_pairs(A) ::= md_map_literal_kv_pair(B) MD_TOKEN_COMMA. {
	A = mlr_dsl_ast_node_alloc_unary("map_literal", MD_AST_NODE_TYPE_MAP_LITERAL, B);
}
md_map_literal_kv_pairs(A) ::= md_map_literal_kv_pair(B) MD_TOKEN_COMMA md_map_literal_kv_pairs(C). {
	A = mlr_dsl_ast_node_prepend_arg(C, B);
}

md_map_literal_kv_pair(A) ::= md_map_literal_key(B) MD_TOKEN_COLON md_map_literal_value(C). {
	A = mlr_dsl_ast_node_alloc_binary("mappair", MD_AST_NODE_TYPE_MAP_LITERAL_PAIR, B, C);
}
md_map_literal_key(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("mapkey", MD_AST_NODE_TYPE_MAP_LITERAL_KEY, B);
}
md_map_literal_value(A) ::= md_rhs(B). {
	A = mlr_dsl_ast_node_alloc_unary("mapval", MD_AST_NODE_TYPE_MAP_LITERAL_VALUE, B);
}
md_map_literal_value(A) ::= md_map_literal(B). {
	A = mlr_dsl_ast_node_alloc_unary("mapval", MD_AST_NODE_TYPE_MAP_LITERAL_VALUE, B);
}
md_map_literal_value(A) ::= MD_TOKEN_FULL_SREC(B). {
	A = mlr_dsl_ast_node_alloc_unary("mapval", MD_AST_NODE_TYPE_MAP_LITERAL_VALUE, B);
}
md_map_literal_value(A) ::= MD_TOKEN_FULL_OOSVAR(B). {
	A = mlr_dsl_ast_node_alloc_unary("mapval", MD_AST_NODE_TYPE_MAP_LITERAL_VALUE, B);
}
