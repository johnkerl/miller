#ifndef LREC_EVALUATORS_H
#define LREC_EVALUATORS_H
#include "containers/mlr_dsl_ast.h"
#include "mapping/lrec_evaluator.h"

lrec_evaluator_t* lrec_evaluator_alloc_from_ast(mlr_dsl_ast_node_t* proot);

#endif // LREC_FEVALUATORS_H
