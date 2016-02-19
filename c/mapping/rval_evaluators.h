#ifndef RVAL_EVALUATORS_H
#define RVAL_EVALUATORS_H
#include <stdio.h>
#include "containers/mlr_dsl_ast.h"
#include "mapping/rval_evaluator.h"

#define TYPE_INFER_STRING_FLOAT_INT 0xce08
#define TYPE_INFER_STRING_FLOAT     0xce09
#define TYPE_INFER_STRING_ONLY      0xce0a

rval_evaluator_t* rval_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* past, int type_inferencing);
rval_evaluator_t* rval_evaluator_alloc_from_string(char* string);

void rval_evaluator_list_functions(FILE* output_stream, char* leader);
// Pass function_name == NULL to get usage for all functions:
void rval_evaluator_function_usage(FILE* output_stream, char* function_name);
void rval_evaluator_list_all_functions_raw(FILE* output_stream);

int test_rval_evaluators_main(int argc, char **argv);

#endif // LREC_FEVALUATORS_H
