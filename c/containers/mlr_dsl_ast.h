// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef MLR_DSL_AST_H
#define MLR_DSL_AST_H
#include "sllv.h"

#define MLR_DSL_AST_NODE_TYPE_STRNUM_LITERAL     0xda00
#define MLR_DSL_AST_NODE_TYPE_BOOLEAN_LITERAL    0xda11
#define MLR_DSL_AST_NODE_TYPE_REGEXI             0xda22
#define MLR_DSL_AST_NODE_TYPE_FIELD_NAME         0xda33
#define MLR_DSL_AST_NODE_TYPE_OOSVAR_NAME        0xda44
#define MLR_DSL_AST_NODE_TYPE_MOOSVAR_NAME       0xda55
#define MLR_DSL_AST_NODE_TYPE_MOOSVAR_INDEX      0xda5c
#define MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME      0xda66
#define MLR_DSL_AST_NODE_TYPE_OPERATOR           0xda77
#define MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT    0xda88
#define MLR_DSL_AST_NODE_TYPE_OOSVAR_ASSIGNMENT  0xda99
#define MLR_DSL_AST_NODE_TYPE_MOOSVAR_ASSIGNMENT 0xdaaa
#define MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE   0xdabb
#define MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY      0xdacc
#define MLR_DSL_AST_NODE_TYPE_FILTER             0xdadd
#define MLR_DSL_AST_NODE_TYPE_GATE               0xdaee
#define MLR_DSL_AST_NODE_TYPE_EMIT               0xdaff

typedef struct _mlr_dsl_ast_t {
	sllv_t* pbegin_statements;
	sllv_t* pmain_statements;
	sllv_t* pend_statements;
} mlr_dsl_ast_t;

typedef struct _mlr_dsl_ast_node_t {
	char*   text;
	int     type;
	sllv_t* pchildren;
} mlr_dsl_ast_node_t;

mlr_dsl_ast_t* mlr_dsl_ast_alloc();

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc(char* text, int type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_zary(char* text, int type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_unary(char* text, int type, mlr_dsl_ast_node_t* pa);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_binary(char* text, int type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_ternary(char* text, int type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother);

// See comments in mlr_dsl_parse.y for this seemingly awkward syntax wherein
// we change the function name after having set it up. This is a consequence of
// bottom-up DSL parsing.
mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(
	mlr_dsl_ast_node_t* pa, char* name);

void mlr_dsl_ast_print(mlr_dsl_ast_t* past);
void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode);
char* mlr_dsl_ast_node_describe_type(int type);

void mlr_dsl_ast_node_free(mlr_dsl_ast_node_t* pnode);

void mlr_dsl_ast_free(mlr_dsl_ast_t* past);

#endif // MLR_DSL_AST_H
