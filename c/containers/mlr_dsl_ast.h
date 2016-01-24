// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef MLR_DSL_AST_H
#define MLR_DSL_AST_H
#include "sllv.h"

#define MLR_DSL_AST_NODE_TYPE_STRNUM_LITERAL    0xaa00 // string     or number
#define MLR_DSL_AST_NODE_TYPE_BOOLEAN_LITERAL   0xaa55 // true/false
#define MLR_DSL_AST_NODE_TYPE_REGEXI            0xaaaa
#define MLR_DSL_AST_NODE_TYPE_FIELD_NAME        0xbbbb
#define MLR_DSL_AST_NODE_TYPE_OOSVAR_NAME       0xbb44
#define MLR_DSL_AST_NODE_TYPE_FUNCTION_NAME     0xcccc
#define MLR_DSL_AST_NODE_TYPE_OPERATOR          0xdd77
#define MLR_DSL_AST_NODE_TYPE_SREC_ASSIGNMENT   0xdd33
#define MLR_DSL_AST_NODE_TYPE_OOSVAR_ASSIGNMENT 0xdddd
#define MLR_DSL_AST_NODE_TYPE_CONTEXT_VARIABLE  0xeeee
#define MLR_DSL_AST_NODE_TYPE_STRIPPED_AWAY     0xffff
#define MLR_DSL_AST_NODE_TYPE_FILTER            0xcc00
#define MLR_DSL_AST_NODE_TYPE_GATE              0xcc33
#define MLR_DSL_AST_NODE_TYPE_EMIT              0xcc66

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
