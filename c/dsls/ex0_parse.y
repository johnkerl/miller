// vim: set filetype=none:
// (Lemon files have .y extensions like Yacc files but are not Yacc.)

%include {
#include <stdio.h>
#include <string.h>
#include <math.h>
#include "../lib/mlrutil.h"
#include "./ex_ast.h"
#include "../containers/sllv.h"

#define DO_WRITE_APPEND // transitional pending lemon-memory issue

// ================================================================
// AST:
// * parens, commas, semis, line endings, whitespace are all stripped away
// * variable names and literal values remain as leaf nodes of the AST
// * = + - * / ** {function names} remain as non-leaf nodes of the AST
// CST: See ex0_cst.c
//
// Note: This parser accepts many things that are invalid, e.g.
// * begin{end{}} -- begin/end not at top level
// * begin{$x=1} -- references to stream records at begin/end
// * break/continue outside of for/while/do-while
// * $x=x -- boundvars outside of for-loop variable bindings
// All of the above are enforced by the CST builder, which takes this parser's output AST as input.
// This is done (a) to keep this grammar from being overly complex, and (b) so we can get much more
// informative error messages in C than in Lemon ('syntax error').
//
// The parser hooks all build up an abstract syntax tree specifically for the CST builder.
// For clearer visuals on what the ASTs look like:
// * See ex0_cst.c
// * See reg_test/run's filter -v and put -v outputs, e.g. in reg_test/expected/out
// * Do "mlr -n put -v 'your expression goes here'"
// ================================================================

}

%token_type     {ex_ast_node_t*}
%default_type   {ex_ast_node_t*}
%extra_argument {ex_ast_t* past}

//void token_destructor(ex_ast_node_t t) {
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
md_body ::= math_expr(B). {
	past->proot = B;
}

math_expr(C) ::= MD_TOKEN_NUMBER(A) MD_TOKEN_PLUS(O)   MD_TOKEN_NUMBER(B) . {
	C = ex_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, A, B);
}

math_expr(C) ::= MD_TOKEN_NUMBER(A) MD_TOKEN_MINUS(O)  MD_TOKEN_NUMBER(B) . {
	C = ex_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, A, B);
}

math_expr(C) ::= MD_TOKEN_NUMBER(A) MD_TOKEN_TIMES(O)  MD_TOKEN_NUMBER(B) . {
	C = ex_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, A, B);
}

math_expr(C) ::= MD_TOKEN_NUMBER(A) MD_TOKEN_DIVIDE(O) MD_TOKEN_NUMBER(B) . {
	C = ex_ast_node_alloc_binary(O->text, MD_AST_NODE_TYPE_OPERATOR, A, B);
}
