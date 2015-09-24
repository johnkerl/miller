#ifndef LREC_EVALUATORS_H
#define LREC_EVALUATORS_H
#include <stdio.h>
#include "containers/mlr_dsl_ast.h"
#include "mapping/lrec_evaluator.h"

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* proot);
void lrec_evaluator_list_functions(FILE* output_stream);
// Pass function_name == NULL to get usage for all functions:
void lrec_evaluator_function_usage(FILE* output_stream, char* function_name);

int test_lrec_evaluators_main(int argc, char **argv);

#endif // LREC_FEVALUATORS_H
