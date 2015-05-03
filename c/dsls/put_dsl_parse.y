
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

put_dsl_body ::= put_dsl_assignments.                // For scan-from-string
put_dsl_body ::= put_dsl_assignments PUT_DSL_EOL. // For scan-from-stdin

// ----------------------------------------------------------------
put_dsl_assignments ::= put_dsl_assignment.
put_dsl_assignments ::= put_dsl_assignment PUT_DSL_SEMICOLON put_dsl_assignments.

// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.
put_dsl_assignment(A)  ::= PUT_DSL_FIELD_NAME(B) PUT_DSL_ASSIGN(O) put_dsl_expr(C). {
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	B = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);

	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
	sllv_add(pasts, A);
}

// ----------------------------------------------------------------
put_dsl_expr(A) ::= put_dsl_term(B). {
	A = B;
}
put_dsl_expr(A) ::= put_dsl_term(B) PUT_DSL_PLUS(O) put_dsl_expr(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
put_dsl_expr(A) ::= put_dsl_term(B) PUT_DSL_MINUS(O) put_dsl_expr(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
put_dsl_expr(A) ::= PUT_DSL_MINUS(O) put_dsl_expr(B). {
	A = mlr_dsl_ast_node_alloc_unary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B);
}
put_dsl_expr(A) ::= put_dsl_term(B) PUT_DSL_DOT(O) put_dsl_expr(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
put_dsl_term(A) ::= put_dsl_factor(B). {
	A = B;
}
put_dsl_term(A) ::= put_dsl_factor(B) PUT_DSL_TIMES(O) put_dsl_term(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}
put_dsl_term(A) ::= put_dsl_factor(B) PUT_DSL_DIVIDE(O) put_dsl_term(C).{
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
put_dsl_factor(A) ::= put_dsl_expitem(B). {
	A = B;
}
put_dsl_factor(A) ::= put_dsl_expitem(B) PUT_DSL_POW(O) put_dsl_factor(C). {
	A = mlr_dsl_ast_node_alloc_binary(O->text, MLR_DSL_AST_NODE_TYPE_OPERATOR, B, C);
}

// ----------------------------------------------------------------
// In the grammar provided to the user, field names are of the form "$x".  But
// within Miller internally, field names are of the form "x".  We coded the
// lexer to give us field names with leading "$" so we can confidently strip it
// off here.
put_dsl_expitem(A) ::= PUT_DSL_FIELD_NAME(B). {
	//A = B;
	char* dollar_name = B->text;
	char* no_dollar_name = &dollar_name[1];
	A = mlr_dsl_ast_node_alloc(no_dollar_name, B->type);
}
put_dsl_expitem(A) ::= PUT_DSL_NUMBER(B). {
	A = B;
}
put_dsl_expitem(A) ::= PUT_DSL_STRING(B). {
	char* input = B->text;
	char* stripped = &input[1];
	// xxx make/call method
	int len = strlen(input);
	stripped[len-2] = 0;
	A = mlr_dsl_ast_node_alloc(stripped, B->type);
}
put_dsl_expitem(A) ::= PUT_DSL_CONTEXT_VARIABLE(B). {
	A = B;
}

put_dsl_expitem(A) ::= PUT_DSL_LPAREN put_dsl_expr(B) PUT_DSL_RPAREN. {
	A = B;
}

// Given "f(a,b,c)": since this is a bottom-up parser, we get first the "a",
// then "a,b", then "a,b,c", then finally "f(a,b,c)". So:
// * On the "a" we make a function sub-AST called "anon(a)".
// * On the "b" we append the next argument to get "anon(a,b)".
// * On the "c" we append the next argument to get "anon(a,b,c)".
// * On the "f" we change the function name to get "f(a,b,c)".

put_dsl_expitem(A) ::= PUT_DSL_FCN_NAME(O) PUT_DSL_LPAREN put_dsl_fcn_args(B) PUT_DSL_RPAREN. {
	A = mlr_dsl_ast_node_set_function_name(B, O->text);
}
// xxx need to invalidate "f(10,)" -- use some non-empty-args expr.
put_dsl_fcn_args(A) ::= . {
	A = mlr_dsl_ast_node_alloc_zary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME);
}

put_dsl_fcn_args(A) ::= put_dsl_expr(B). {
	A = mlr_dsl_ast_node_alloc_unary("anon", MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME, B);
}
put_dsl_fcn_args(A) ::= put_dsl_fcn_args(B) PUT_DSL_COMMA put_dsl_expr(C). {
	A = mlr_dsl_ast_node_append_arg(B, C);
}
