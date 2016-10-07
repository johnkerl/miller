#ifndef MLR_DSL_BLOCKED_AST_H
#define MLR_DSL_BLOCKED_AST_H

#include "containers/mlr_dsl_ast.h"
#include "containers/sllv.h"

// ================================================================
// The Lemon parser produces a single raw abstract syntax tree.  This container
// simply has the subtrees organized by top-level statement blocks. This is just
// a minor reorganization but it makes stack allocation and CST-builds simpler.
// ================================================================

typedef struct _blocked_ast_t {
	sllv_t* pfunc_defs;
	sllv_t* psubr_defs;
	sllv_t* pbegin_blocks;
	mlr_dsl_ast_node_t* pmain_block;
	sllv_t* pend_blocks;
} blocked_ast_t;

// This strips nodes off the raw AST and transfers them to the analyzed AST.
blocked_ast_t* blocked_ast_alloc(mlr_dsl_ast_t* past);
void blocked_ast_free(blocked_ast_t* paast);

#endif // MLR_DSL_BLOCKED_AST_H
