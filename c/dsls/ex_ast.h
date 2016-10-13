// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef EX_AST_H
#define EX_AST_H
#include "../containers/sllv.h"

// ----------------------------------------------------------------
typedef enum _ex_ast_node_type_t {
	MD_AST_NODE_TYPE_STATEMENT_BLOCK,
	MD_AST_NODE_TYPE_BEGIN,
	MD_AST_NODE_TYPE_END,
	MD_AST_NODE_TYPE_STRING_LITERAL,
	MD_AST_NODE_TYPE_STRNUM_LITERAL,
	MD_AST_NODE_TYPE_BOOLEAN_LITERAL,
	MD_AST_NODE_TYPE_REGEXI,
	MD_AST_NODE_TYPE_FIELD_NAME, // E.g. $x
	MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME, // E.g. $[@x]
	MD_AST_NODE_TYPE_FULL_SREC,
	MD_AST_NODE_TYPE_OOSVAR_KEYLIST,
	MD_AST_NODE_TYPE_FULL_OOSVAR,
	MD_AST_NODE_TYPE_NON_SIGIL_NAME,
	MD_AST_NODE_TYPE_OPERATOR,
	MD_AST_NODE_TYPE_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT,
	MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT,
	MD_AST_NODE_TYPE_CONTEXT_VARIABLE,
	MD_AST_NODE_TYPE_ENV,
	MD_AST_NODE_TYPE_STRIPPED_AWAY,
	MD_AST_NODE_TYPE_CONDITIONAL_BLOCK,
	MD_AST_NODE_TYPE_FILTER,
	MD_AST_NODE_TYPE_UNSET,
	MD_AST_NODE_TYPE_EMITF,
	MD_AST_NODE_TYPE_EMITF_WRITE,
	MD_AST_NODE_TYPE_EMITF_APPEND,
	MD_AST_NODE_TYPE_EMITP,
	MD_AST_NODE_TYPE_EMITP_WRITE,
	MD_AST_NODE_TYPE_EMITP_APPEND,
	MD_AST_NODE_TYPE_EMIT,
	MD_AST_NODE_TYPE_EMIT_WRITE,
	MD_AST_NODE_TYPE_EMIT_APPEND,
	MD_AST_NODE_TYPE_EMITP_LASHED,
	MD_AST_NODE_TYPE_EMITP_LASHED_WRITE,
	MD_AST_NODE_TYPE_EMITP_LASHED_APPEND,
	MD_AST_NODE_TYPE_EMIT_LASHED,
	MD_AST_NODE_TYPE_EMIT_LASHED_WRITE,
	MD_AST_NODE_TYPE_EMIT_LASHED_APPEND,
	MD_AST_NODE_TYPE_DUMP,
	MD_AST_NODE_TYPE_DUMP_WRITE,
	MD_AST_NODE_TYPE_DUMP_APPEND,
	MD_AST_NODE_TYPE_EDUMP,
	MD_AST_NODE_TYPE_PRINT,
	MD_AST_NODE_TYPE_PRINT_WRITE,
	MD_AST_NODE_TYPE_PRINT_APPEND,
	MD_AST_NODE_TYPE_PRINTN,
	MD_AST_NODE_TYPE_PRINTN_WRITE,
	MD_AST_NODE_TYPE_PRINTN_APPEND,
	MD_AST_NODE_TYPE_EPRINT,
	MD_AST_NODE_TYPE_EPRINTN,
	MD_AST_NODE_TYPE_ALL,
	MD_AST_NODE_TYPE_NOP, // only for parser internals; should not be in the AST returned by the parser
	MD_AST_NODE_TYPE_WHILE,
	MD_AST_NODE_TYPE_DO_WHILE,
	MD_AST_NODE_TYPE_FOR_SREC,
	MD_AST_NODE_TYPE_FOR_OOSVAR,
	MD_AST_NODE_TYPE_FOR_VARIABLES,
	MD_AST_NODE_TYPE_LOCAL_VARIABLE,
	MD_AST_NODE_TYPE_IN,
	MD_AST_NODE_TYPE_BREAK,
	MD_AST_NODE_TYPE_CONTINUE,
	MD_AST_NODE_TYPE_IF_HEAD,
	MD_AST_NODE_TYPE_IF_ITEM,
} ex_ast_node_type_t;

typedef struct _ex_ast_node_t {
	char*                   text;
	ex_ast_node_type_t type;
	sllv_t*                 pchildren;
} ex_ast_node_t;

typedef struct _ex_ast_t {
	ex_ast_node_t* proot;
} ex_ast_t;

// ----------------------------------------------------------------
ex_ast_t* ex_ast_alloc();

ex_ast_node_t* ex_ast_node_alloc(char* text, ex_ast_node_type_t type);

ex_ast_node_t* ex_ast_node_alloc_zary(char* text, ex_ast_node_type_t type);

ex_ast_node_t* ex_ast_node_alloc_unary(char* text, ex_ast_node_type_t type, ex_ast_node_t* pa);

ex_ast_node_t* ex_ast_node_alloc_binary(char* text, ex_ast_node_type_t type,
	ex_ast_node_t* pa, ex_ast_node_t* pb);

ex_ast_node_t* ex_ast_node_alloc_ternary(char* text, ex_ast_node_type_t type,
	ex_ast_node_t* pa, ex_ast_node_t* pb, ex_ast_node_t* pc);

ex_ast_node_t* ex_ast_node_copy(ex_ast_node_t* pother);
// These are so the parser can expand '$x += 1' to '$x = $x + 1', etc.
ex_ast_node_t* ex_ast_tree_copy(ex_ast_node_t* pother);

// See comments in ex_parse.y for this seemingly awkward syntax wherein
// we change the function name after having set it up. This is a consequence of
// bottom-up DSL parsing.
ex_ast_node_t* ex_ast_node_prepend_arg(ex_ast_node_t* pa, ex_ast_node_t* pb);
ex_ast_node_t* ex_ast_node_append_arg(ex_ast_node_t* pa, ex_ast_node_t* pb);
ex_ast_node_t* ex_ast_node_set_function_name(ex_ast_node_t* pa, char* name);

void ex_ast_print(ex_ast_t* past);
void ex_ast_node_print(ex_ast_node_t* pnode);
void ex_ast_node_fprint(ex_ast_node_t* pnode, FILE* o);
char* ex_ast_node_describe_type(ex_ast_node_type_t type);

void ex_ast_node_free(ex_ast_node_t* pnode);

void ex_ast_free(ex_ast_t* past);

#endif // EX_AST_H
