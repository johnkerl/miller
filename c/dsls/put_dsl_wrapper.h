#ifndef PUT_DSL_WRAPPER_H
#define PUT_DSL_WRAPPER_H
#include  "../containers/mlr_dsl_ast.h"
#include  "../containers/sllv.h"

// Returns linked list of mlr_dsl_ast_node_t*.
sllv_t* put_dsl_parse(char* string);

#endif // PUT_DSL_WRAPPER_H
