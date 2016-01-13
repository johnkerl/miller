// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef MLR_DSL_AST_H
#define MLR_DSL_AST_H
#include "sllv.h"

#define MLR_DSL_AST_NODE_TYPE_LITERAL          0xaa00
#define MLR_DSL_AST_NODE_TYPE_REGEXI           0xaaaa
#define MLR_DSL_AST_NODE_TYPE_FIELD_NAME       0xbbbb
#define MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME    0xcccc
#define MLR_DSL_AST_NODE_TYPE_OPERATOR         0xdddd
#define MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE 0xeeee
#define MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY    0xffff

typedef struct _mlr_dsl_ast_node_t {
	char*   text;
	int     type;
	sllv_t* pchildren;
} mlr_dsl_ast_node_t;

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc(char* text, int type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_zary(char* text, int type);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_unary(char* text, int type,
	mlr_dsl_ast_node_t* pa);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_binary(char* text, int type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_ternary(char* text, int type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc);

// See comments in put_dsl_parse.y for this seemingly awkward syntax wherein
// we change the function name after having set it up. This is a consequence of
// bottom-up DSL parsing.
mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(
	mlr_dsl_ast_node_t* pa, char* name);

void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode);
char* mlr_dsl_ast_node_describe_type(int type);

void mlr_dsl_ast_node_free(mlr_dsl_ast_node_t* pnode);

#endif // MLR_DSL_AST_H
