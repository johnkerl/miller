#ifndef LREC_EVALUATORS_H
#define LREC_EVALUATORS_H
#include <stdio.h>
#include "containers/mlr_dsl_ast.h"
#include "mapping/lrec_evaluator.h"

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* proot);
void lrec_evaluator_describe_functions(FILE* output_stream);

#endif // LREC_FEVALUATORS_H
