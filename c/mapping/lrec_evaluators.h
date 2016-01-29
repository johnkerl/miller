#ifndef LREC_EVALUATORS_H
#define LREC_EVALUATORS_H
#include <stdio.h>
#include "containers/mlr_dsl_ast.h"
#include "mapping/lrec_evaluator.h"

#define TYPE_INFER_STRING_FLOAT_INT 0xce08
#define TYPE_INFER_STRING_FLOAT     0xce09
#define TYPE_INFER_STRING_ONLY      0xce0a

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* past, int type_inferencing);
lrec_evaluator_t* lrec_evaluator_alloc_from_moosvar_name(char* moosvar_name); // xxx still need this external?
lrec_evaluator_t* lrec_evaluator_alloc_from_strnum_literal(char* string, int type_inferencing);

void lrec_evaluator_list_functions(FILE* output_stream, char* leader);
// Pass function_name == NULL to get usage for all functions:
void lrec_evaluator_function_usage(FILE* output_stream, char* function_name);
void lrec_evaluator_list_all_functions_raw(FILE* output_stream);

int test_lrec_evaluators_main(int argc, char **argv);

#endif // LREC_FEVALUATORS_H
