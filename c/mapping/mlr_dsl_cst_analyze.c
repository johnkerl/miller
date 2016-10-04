#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// xxx make a summary comment here

// ----------------------------------------------------------------
// xxx have ast freed back where it was (for callsite-balance) but w/ has-been-exfoliated comment
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

	sllv_t* pnodelist = past->proot->pchildren;
	while (pnodelist->phead) {
		mlr_dsl_ast_node_t* pnode = sllv_pop(pnodelist);
		switch (pnode->type) {
		case MD_AST_NODE_TYPE_FUNC_DEF:
			sllv_append(paast->pfunc_defs, pnode);
			break;
		case MD_AST_NODE_TYPE_SUBR_DEF:
			sllv_append(paast->psubr_defs, pnode);
			break;
		case MD_AST_NODE_TYPE_BEGIN:
			sllv_append(paast->pbegin_blocks, pnode);
			break;
		case MD_AST_NODE_TYPE_END:
			sllv_append(paast->pend_blocks, pnode);
			break;
		default:
			sllv_append(paast->pmain_block->pchildren, pnode);
			break;
		}
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

// ================================================================
//                       # ---- FUNC FRAME: defcount 5 {a,b,c,i,j}
// func f(a, b, c) {     # arg define A.1,A.2,A.3
//     local i = 1;      # explicit define A.4
//     j = 2;            # implicit define A.5
//                       #
//                       # ---- IF FRAME: defcount 2 {k,m}
//     if (a == 3) {     # RHS A.1
//         local k = 4;  # explicit define B.1
//         j = 5;        # LHS A.5
//         m = 6;        # implicit define B.2
//         k = a;        # LHS B.1 RHS A.1
//         k = i;        # LHS B.1 RHS A.4
//         m = k;        # LHS B.2 RHS B.1
//                       #
//                       # ---- ELSE FRAME: defcount 3 {n,g,h}
//     } else {          #
//         n = b;        #
//         g = n         #
//         h = b         #
//     }                 #
//                       #
//     b = 7;            # LHS A.2
//     i = z;            # LHS A.4 RHS unresolved
// }                     #
// ================================================================

// * local-var defines are certainly for the current frame
// * local-var writes need backtracing (if not found at the current frame)
// * local-var reads  need backtracing (if not found at the current frame)
// * unresolved-read needs special handling -- maybe a root-level mv_absent at index 0?
//
// One frame per curly-braced block
// One framegroup per block (funcdef, subrdef, begin, end, main)
// -> has maxdepth attrs
//
// frame_t at analysis phase:
// * hss I guess. or better: lhmsi. has size attr.
//
// frame_t at run phase:
// * numvars attr
// * indices refer to frame_group's array:
//   o non-negative indices are local
//   o negative indices are locals within ancestor node(s)
//
// frame_group_t at analysis phase:
// * this is a tree
// * each node has nameset/defcount for its locals
// * each node also has max defcount for its transitive children?
//
// frame_group_t at run phase:
// * maxdepth attr
// * array of mlrvals
// * optionally (or always?) a single slot is the undef.
//
// storage options:
// * decorate ast_node_t w/ default-null pointer to allocated analysis info?
//   this way AST points to analysis.
// * analysis tree w/ pointers to statement-list nodes?
//   this way analysis points to AST.
