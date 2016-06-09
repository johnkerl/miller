// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef MLR_DSL_AST_H
#define MLR_DSL_AST_H
#include "sllv.h"

// ----------------------------------------------------------------
typedef enum _mlr_dsl_ast_node_type_t {
	MD_AST_NODE_TYPE_STATEMENT_LIST,
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
	MD_AST_NODE_TYPE_EMITP,
	MD_AST_NODE_TYPE_EMIT,
	MD_AST_NODE_TYPE_DUMP,
	MD_AST_NODE_TYPE_PRINT,
	MD_AST_NODE_TYPE_ALL,
	MD_AST_NODE_TYPE_NOP, // only for parser internals; should not be in the AST returned by the parser
	MD_AST_NODE_TYPE_WHILE,
	MD_AST_NODE_TYPE_DO_WHILE,
	MD_AST_NODE_TYPE_FOR_SREC,
	MD_AST_NODE_TYPE_FOR_OOSVAR,
	MD_AST_NODE_TYPE_FOR_VARIABLES,
	MD_AST_NODE_TYPE_BOUND_VARIABLE,
	MD_AST_NODE_TYPE_IN,
	MD_AST_NODE_TYPE_BREAK,
	MD_AST_NODE_TYPE_CONTINUE,
	MD_AST_NODE_TYPE_IF_HEAD,
	MD_AST_NODE_TYPE_IF_ITEM,
} mlr_dsl_ast_node_type_t;

typedef struct _mlr_dsl_ast_node_t {
	char*                   text;
	mlr_dsl_ast_node_type_t type;
	sllv_t*                 pchildren;
} mlr_dsl_ast_node_t;

typedef struct _mlr_dsl_ast_t {
	mlr_dsl_ast_node_t* proot;
} mlr_dsl_ast_t;

// ----------------------------------------------------------------
mlr_dsl_ast_t* mlr_dsl_ast_alloc();

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc(char* text, mlr_dsl_ast_node_type_t type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_zary(char* text, mlr_dsl_ast_node_type_t type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_unary(char* text, mlr_dsl_ast_node_type_t type, mlr_dsl_ast_node_t* pa);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_binary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_ternary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother);
// These are so the parser can expand '$x += 1' to '$x = $x + 1', etc.
mlr_dsl_ast_node_t* mlr_dsl_ast_tree_copy(mlr_dsl_ast_node_t* pother);

// See comments in mlr_dsl_parse.y for this seemingly awkward syntax wherein
// we change the function name after having set it up. This is a consequence of
// bottom-up DSL parsing.
mlr_dsl_ast_node_t* mlr_dsl_ast_node_prepend_arg(mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg(mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(mlr_dsl_ast_node_t* pa, char* name);

void mlr_dsl_ast_print(mlr_dsl_ast_t* past);
void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode);
void mlr_dsl_ast_node_fprint(mlr_dsl_ast_node_t* pnode, FILE* o);
char* mlr_dsl_ast_node_describe_type(mlr_dsl_ast_node_type_t type);

void mlr_dsl_ast_node_free(mlr_dsl_ast_node_t* pnode);

void mlr_dsl_ast_free(mlr_dsl_ast_t* past);

#endif // MLR_DSL_AST_H
