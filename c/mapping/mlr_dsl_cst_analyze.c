#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// xxx make a summary comment here

// ----------------------------------------------------------------
analyzed_ast_t* analyzed_ast_alloc(mlr_dsl_ast_t* past) {
	analyzed_ast_t* paast = mlr_malloc_or_die(sizeof(analyzed_ast_t));

	paast->pfunc_defs    = sllv_alloc();
	paast->psubr_defs    = sllv_alloc();
	paast->pbegin_blocks = sllv_alloc();
	paast->pmain_block   = mlr_dsl_ast_node_alloc_zary("main_block", MD_AST_NODE_TYPE_STATEMENT_LIST);
	paast->pend_blocks   = sllv_alloc();

	if (past->proot->type != MD_AST_NODE_TYPE_STATEMENT_LIST) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected root node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_LIST),
			mlr_dsl_ast_node_describe_type(past->proot->type));
		exit(1);
	}

	// xxx walk through the AST and transfer nodes
	// xxx replace with while/pop loop
	for (sllve_t* pe = past->proot->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		switch (pnode->type) {
		case MD_AST_NODE_TYPE_FUNC_DEF:
			break;
		case MD_AST_NODE_TYPE_SUBR_DEF:
			break;
		case MD_AST_NODE_TYPE_BEGIN:
			break;
		case MD_AST_NODE_TYPE_END:
			break;
		default:
			break;
		}
		//printf("\n");
		//printf("================================================================\n");
		//mlr_dsl_ast_node_print(pnode);
		//printf("\n");
	}

	return paast;
}

// ----------------------------------------------------------------
void analyzed_ast_free(analyzed_ast_t* paast) {
	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_free(pe->pvvalue);
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_free(pe->pvvalue);
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_free(pe->pvvalue);
	mlr_dsl_ast_node_free(paast->pmain_block);
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_free(pe->pvvalue);

	sllv_free(paast->pfunc_defs);
	sllv_free(paast->psubr_defs);
	sllv_free(paast->pbegin_blocks);
	sllv_free(paast->pend_blocks);
	free(paast);

}
