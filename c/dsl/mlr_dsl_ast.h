// ================================================================
// Miller abstract syntax tree for put and filter.
// ================================================================

#ifndef MLR_DSL_AST_H
#define MLR_DSL_AST_H
#include "../containers/sllv.h"
#include "../containers/type_decl.h"

// ----------------------------------------------------------------
typedef enum _mlr_dsl_ast_node_type_t {
	MD_AST_NODE_TYPE_STATEMENT_BLOCK,
	MD_AST_NODE_TYPE_STATEMENT_LIST,
	MD_AST_NODE_TYPE_FUNC_DEF,
	MD_AST_NODE_TYPE_SUBR_DEF,
	MD_AST_NODE_TYPE_FUNCTION_CALLSITE,
	MD_AST_NODE_TYPE_INDEXED_FUNCTION_CALLSITE,
	MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST,
	MD_AST_NODE_TYPE_SUBR_CALLSITE,
	MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION,
	MD_AST_NODE_TYPE_UNTYPED_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_NUMERIC_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_INT_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_FLOAT_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_BOOLEAN_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_STRING_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_MAP_PARAMETER_DEFINITION,
	MD_AST_NODE_TYPE_RETURN_VALUE,
	MD_AST_NODE_TYPE_RETURN_VOID,
	MD_AST_NODE_TYPE_BEGIN,
	MD_AST_NODE_TYPE_END,
	MD_AST_NODE_TYPE_STRING_LITERAL,
	MD_AST_NODE_TYPE_NUMERIC_LITERAL,
	MD_AST_NODE_TYPE_BOOLEAN_LITERAL,
	MD_AST_NODE_TYPE_MAP_LITERAL,
	MD_AST_NODE_TYPE_MAP_LITERAL_PAIR,
	MD_AST_NODE_TYPE_MAP_LITERAL_KEY,
	MD_AST_NODE_TYPE_MAP_LITERAL_VALUE,
	MD_AST_NODE_TYPE_REGEXI,
	MD_AST_NODE_TYPE_FIELD_NAME, // E.g. value = $x
	MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME, // E.g. value = $[@x]
	MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME, // E.g. key = $[[3]]
	MD_AST_NODE_TYPE_FULL_SREC,
	MD_AST_NODE_TYPE_OOSVAR_KEYLIST,
	MD_AST_NODE_TYPE_FULL_OOSVAR,
	MD_AST_NODE_TYPE_NON_SIGIL_NAME,
	MD_AST_NODE_TYPE_OPERATOR,
	MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT,
	MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT,
	MD_AST_NODE_TYPE_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT, // E.g. $[@y] = newvalue
	MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME_ASSIGNMENT, // E.g. $[[3]] = newkey
	MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT,
	MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT,
	MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT,
	MD_AST_NODE_TYPE_ENV_ASSIGNMENT,
	MD_AST_NODE_TYPE_CONTEXT_VARIABLE,
	MD_AST_NODE_TYPE_ENV,
	MD_AST_NODE_TYPE_STRIPPED_AWAY,
	MD_AST_NODE_TYPE_CONDITIONAL_BLOCK,
	MD_AST_NODE_TYPE_FILTER,
	MD_AST_NODE_TYPE_UNSET,
	MD_AST_NODE_TYPE_PIPE,
	MD_AST_NODE_TYPE_FILE_WRITE,
	MD_AST_NODE_TYPE_FILE_APPEND,
	MD_AST_NODE_TYPE_TEE,
	MD_AST_NODE_TYPE_EMITF,
	MD_AST_NODE_TYPE_EMITP,
	MD_AST_NODE_TYPE_EMIT,
	MD_AST_NODE_TYPE_EMITP_LASHED,
	MD_AST_NODE_TYPE_EMIT_LASHED,
	MD_AST_NODE_TYPE_DUMP,
	MD_AST_NODE_TYPE_EDUMP,
	MD_AST_NODE_TYPE_PRINT,
	MD_AST_NODE_TYPE_PRINTN,
	MD_AST_NODE_TYPE_EPRINT,
	MD_AST_NODE_TYPE_EPRINTN,
	MD_AST_NODE_TYPE_STDOUT,
	MD_AST_NODE_TYPE_STDERR,
	MD_AST_NODE_TYPE_STREAM,
	MD_AST_NODE_TYPE_ALL,
	MD_AST_NODE_TYPE_NOP, // only for parser internals; should not be in the AST returned by the parser
	MD_AST_NODE_TYPE_WHILE,
	MD_AST_NODE_TYPE_DO_WHILE,
	MD_AST_NODE_TYPE_FOR_SREC,
	MD_AST_NODE_TYPE_FOR_SREC_KEY_ONLY,
	MD_AST_NODE_TYPE_FOR_OOSVAR,
	MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY,
	MD_AST_NODE_TYPE_FOR_LOCAL_MAP,
	MD_AST_NODE_TYPE_FOR_LOCAL_MAP_KEY_ONLY,
	MD_AST_NODE_TYPE_FOR_MAP_LITERAL,
	MD_AST_NODE_TYPE_FOR_MAP_LITERAL_KEY_ONLY,
	MD_AST_NODE_TYPE_FOR_FUNC_RETVAL,
	MD_AST_NODE_TYPE_FOR_FUNC_RETVAL_KEY_ONLY,
	MD_AST_NODE_TYPE_FOR_VARIABLES,
	MD_AST_NODE_TYPE_TRIPLE_FOR,
	MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE,
	MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE,
	MD_AST_NODE_TYPE_IN,
	MD_AST_NODE_TYPE_BREAK,
	MD_AST_NODE_TYPE_CONTINUE,
	MD_AST_NODE_TYPE_IF_HEAD,
	MD_AST_NODE_TYPE_IF_ITEM,
} mlr_dsl_ast_node_type_t;

#define MD_UNUSED_INDEX -1000000000

typedef struct _mlr_dsl_ast_node_t {
	char*                   text;
	mlr_dsl_ast_node_type_t type;
	sllv_t*                 pchildren;

	// For bind-stack allocation only in local-var nodes: unused for any other node types.
	int vardef_subframe_relative_index; // pass 1 output: which index in subframe
	int vardef_subframe_index;          // pass 1 output: which subframe the variable is defined in
	int vardef_frame_relative_index;    // pass 2 output: index relative to full stack frame

	// For bind-stack allocation only in statement-block nodes: unused for any other node types.
	int subframe_var_count;
	int max_subframe_depth;
	int max_var_depth;

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

mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_quaternary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc, mlr_dsl_ast_node_t* pd);

mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother);
// These are so the parser can expand '$x += 1' to '$x = $x + 1', etc.
mlr_dsl_ast_node_t* mlr_dsl_ast_tree_copy(mlr_dsl_ast_node_t* pother);

// See comments in mlr_dsl_parse.y for this seemingly awkward syntax wherein
// we change the function name after having set it up. This is a consequence of
// bottom-up DSL parsing.
mlr_dsl_ast_node_t* mlr_dsl_ast_node_prepend_arg(mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg(mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg_to_second_child(mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb);
mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(mlr_dsl_ast_node_t* pa, char* name);

void mlr_dsl_ast_node_replace_text(mlr_dsl_ast_node_t* pa, char* text);

int mlr_dsl_ast_node_type_to_type_mask(mlr_dsl_ast_node_type_t type);

int mlr_dsl_ast_node_cannot_be_bare_boolean(mlr_dsl_ast_node_t* pnode);

void mlr_dsl_ast_print(mlr_dsl_ast_t* past);
void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode);
void mlr_dsl_ast_node_fprint(mlr_dsl_ast_node_t* pnode, FILE* o);
void mlr_dsl_ast_node_pretty_fprint(mlr_dsl_ast_node_t* pnode, FILE* o);
char* mlr_dsl_ast_node_describe_type(mlr_dsl_ast_node_type_t type);

void mlr_dsl_ast_node_free(mlr_dsl_ast_node_t* pnode);

void mlr_dsl_ast_free(mlr_dsl_ast_t* past);

#endif // MLR_DSL_AST_H
