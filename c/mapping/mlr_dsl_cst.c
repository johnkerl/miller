#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/hss.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

// ================================================================
// The Lemon parser in dsls/mlr_dsl_parse.y builds up an abstract syntax tree
// specifically for the CST builder here.
//
// For clearer visuals on what the ASTs look like:
// * See dsls/mlr_dsl_parse.y
// * See reg_test/run's filter -v and put -v outputs, e.g. in reg_test/expected/out
// * Do "mlr -n put -v 'your expression goes here'"
// ================================================================

static mlr_dsl_ast_node_t* get_list_for_block(mlr_dsl_ast_node_t* pnode);
mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_final_filter_statement(mlr_dsl_cst_t* pcst,
	mlr_dsl_ast_node_t* pnode, int negate_final_filter, int type_inferencing, int context_flags);
static void mlr_dsl_cst_resolve_subr_callsites(mlr_dsl_cst_t* pcst);

// ----------------------------------------------------------------
// Main entry point for AST-to-CST for mlr put and mlr filter.
//
// Example AST (using put -v):
//
// $ mlr -n put -v '#begin{@a=1;@b=2};$m=2;$n=4;end{@y=5;@z=6}'
// AST ROOT:
// text="list", type=statement_list:
//     text="begin", type=begin:
//         text="list", type=statement_list:
//             text="=", type=oosvar_assignment:
//                 text="oosvar_keylist", type=oosvar_keylist:
//                     text="a", type=string_literal.
//                 text="1", type=numeric_literal.
//             text="=", type=oosvar_assignment:
//                 text="oosvar_keylist", type=oosvar_keylist:
//                     text="b", type=string_literal.
//                 text="2", type=numeric_literal.
//     text="=", type=srec_assignment:
//         text="m", type=field_name.
//         text="2", type=numeric_literal.
//     text="=", type=srec_assignment:
//         text="n", type=field_name.
//         text="4", type=numeric_literal.
//     text="end", type=end:
//         text="list", type=statement_list:
//             text="=", type=oosvar_assignment:
//                 text="oosvar_keylist", type=oosvar_keylist:
//                     text="y", type=string_literal.
//                 text="5", type=numeric_literal.
//             text="=", type=oosvar_assignment:
//                 text="oosvar_keylist", type=oosvar_keylist:
//                     text="z", type=string_literal.
//                 text="6", type=numeric_literal.

mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int print_ast, int trace_stack_allocation,
	int type_inferencing, int flush_every_record,
	int do_final_filter, int negate_final_filter) // for mlr filter
{
	int context_flags = do_final_filter ? IN_MLR_FILTER : 0;
	// The root node is not populated on empty-string input to the parser.
	if (past->proot == NULL) {
		if (do_final_filter) {
			fprintf(stderr, "%s: filter statement must not be empty.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		past->proot = mlr_dsl_ast_node_alloc_zary("list", MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	}

	mlr_dsl_cst_t* pcst = mlr_malloc_or_die(sizeof(mlr_dsl_cst_t));

	pcst->paast = blocked_ast_alloc(past);

	// Assign local-variable names to indices within frame-stack.
	blocked_ast_allocate_locals(pcst->paast, trace_stack_allocation);

	pcst->pfmgr          = fmgr_alloc();
	pcst->psubr_defsites = lhmsv_alloc();
	pcst->psubr_callsite_statements_to_resolve = sllv_alloc();
	pcst->flush_every_record = flush_every_record;

	if (print_ast) {
		printf("\n");
		printf("BLOCKED AST:\n");
	}

	for (sllve_t* pe = pcst->paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		if (print_ast) {
			printf("\n");
			printf("FUNCTION DEFINITION:\n");
			mlr_dsl_ast_node_print(pnode);
		}
		udf_defsite_state_t* pudf_defsite_state = mlr_dsl_cst_alloc_udf(pcst, pnode,
			type_inferencing, context_flags);
		fmgr_install_udf(pcst->pfmgr, pudf_defsite_state);
	}

	for (sllve_t* pe = pcst->paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		if (print_ast) {
			printf("\n");
			printf("SUBROUTINE DEFINITION:\n");
			mlr_dsl_ast_node_print(pnode);
		}
		subr_defsite_t* psubr_defsite = mlr_dsl_cst_alloc_subroutine(pcst, pnode, type_inferencing, context_flags);
		if (lhmsv_get(pcst->psubr_defsites, psubr_defsite->name)) {
			fprintf(stderr, "%s: subroutine named \"%s\" has already been defined.\n",
				MLR_GLOBALS.bargv0, psubr_defsite->name);
			exit(1);
		}
		lhmsv_put(pcst->psubr_defsites, psubr_defsite->name, psubr_defsite, NO_FREE);
	}

	pcst->pbegin_blocks = sllv_alloc();
	for (sllve_t* pe = pcst->paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		if (print_ast) {
			printf("\n");
			printf("BEGIN-BLOCK:\n");
			mlr_dsl_ast_node_print(pnode);
		}
		MLR_INTERNAL_CODING_ERROR_IF(pnode->max_var_depth == MD_UNUSED_INDEX);
		MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
		cst_top_level_statement_block_t* pblock = cst_top_level_statement_block_alloc(pnode->max_var_depth,
			pnode->subframe_var_count);
		for (sllve_t* pf = pnode->pchildren->phead; pf != NULL; pf = pf->pnext) {
			mlr_dsl_ast_node_t* plistnode = get_list_for_block(pnode);
			for (sllve_t* pg = plistnode->pchildren->phead; pg != NULL; pg = pg->pnext) {
				mlr_dsl_ast_node_t* pchild = pg->pvvalue;
				sllv_append(pblock->pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pchild,
					type_inferencing, context_flags | IN_BEGIN_OR_END));
			}
		}
		sllv_append(pcst->pbegin_blocks, pblock);
	}

	pcst->pend_blocks = sllv_alloc();
	for (sllve_t* pe = pcst->paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		if (print_ast) {
			printf("\n");
			printf("END-BLOCK:\n");
			mlr_dsl_ast_node_print(pnode);
		}
		MLR_INTERNAL_CODING_ERROR_IF(pnode->max_var_depth == MD_UNUSED_INDEX);
		MLR_INTERNAL_CODING_ERROR_IF(pnode->subframe_var_count == MD_UNUSED_INDEX);
		cst_top_level_statement_block_t* pblock = cst_top_level_statement_block_alloc(pnode->max_var_depth,
			pnode->subframe_var_count);
		for (sllve_t* pf = pnode->pchildren->phead; pf != NULL; pf = pf->pnext) {
			mlr_dsl_ast_node_t* plistnode = get_list_for_block(pnode);
			for (sllve_t* pg = plistnode->pchildren->phead; pg != NULL; pg = pg->pnext) {
				mlr_dsl_ast_node_t* pchild = pg->pvvalue;
				sllv_append(pblock->pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pchild,
					type_inferencing, context_flags | IN_BEGIN_OR_END));
			}
		}
		sllv_append(pcst->pend_blocks, pblock);
	}

	if (print_ast) {
		printf("\n");
		printf("MAIN BLOCK:\n");
		mlr_dsl_ast_node_print(pcst->paast->pmain_block);
	}
	MLR_INTERNAL_CODING_ERROR_IF(pcst->paast->pmain_block->max_var_depth == MD_UNUSED_INDEX);
	MLR_INTERNAL_CODING_ERROR_IF(pcst->paast->pmain_block->subframe_var_count == MD_UNUSED_INDEX);
	pcst->pmain_block = cst_top_level_statement_block_alloc(pcst->paast->pmain_block->max_var_depth,
		pcst->paast->pmain_block->subframe_var_count);
	for (sllve_t* pe = pcst->paast->pmain_block->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;

		// The last statement of mlr filter must be a bare boolean.
		if (do_final_filter && pe->pnext == NULL) {
			sllv_append(pcst->pmain_block->pblock->pstatements, mlr_dsl_cst_alloc_final_filter_statement(
				pcst, pnode, negate_final_filter, type_inferencing, context_flags | IN_MLR_FINAL_FILTER));
		} else {
			sllv_append(pcst->pmain_block->pblock->pstatements, mlr_dsl_cst_alloc_statement(pcst, pnode,
				type_inferencing, context_flags));
		}
	}

	if (print_ast) {
		printf("\n");
	}

	// Now that all subroutine/function definitions have been done, resolve
	// their callsites whose locations we stashed during the CST build. (Without
	// this delayed resolution, there could be no recursion, and subroutines
	// could call one another only in the reverse order of their definition.
	// E.g. if 's' is defined and then 't', then t could call s but s could not
	// call t [subroutine not defined yet], and neither could call itself.)
	//
	// This object-binding step is not a full pass through the CST since we've
	// maintained pointers into callsites which we can now directly poke.
	fmgr_resolve_func_callsites(pcst->pfmgr);
	mlr_dsl_cst_resolve_subr_callsites(pcst);

	return pcst;
}

// ----------------------------------------------------------------
void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst) {
	if (pcst == NULL)
		return;

	if (pcst->pbegin_blocks != NULL) {
		for (sllve_t* pe = pcst->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
			cst_top_level_statement_block_free(pe->pvvalue);
		}
		sllv_free(pcst->pbegin_blocks);
	}

	cst_top_level_statement_block_free(pcst->pmain_block);

	if (pcst->pend_blocks != NULL) {
		for (sllve_t* pe = pcst->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
			cst_top_level_statement_block_free(pe->pvvalue);
		}
		sllv_free(pcst->pend_blocks);
	}

	fmgr_free(pcst->pfmgr);

	// Void-star payloads already popped and freed during symbol-resolution phase of CST alloc
	sllv_free(pcst->psubr_callsite_statements_to_resolve);

	if (pcst->psubr_defsites != NULL) {
		for (lhmsve_t* pe = pcst->psubr_defsites->phead; pe != NULL; pe = pe->pnext) {
			subr_defsite_t* psubr_defsite = pe->pvvalue;
			mlr_dsl_cst_free_subroutine(psubr_defsite);
		}
		lhmsv_free(pcst->psubr_defsites);
	}

	blocked_ast_free(pcst->paast);

	free(pcst);
}

// ----------------------------------------------------------------
// For begin, end, cond: there must be one child node, of type list.
static mlr_dsl_ast_node_t* get_list_for_block(mlr_dsl_ast_node_t* pnode) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->pchildren->phead == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pnode->pchildren->phead->pnext != NULL);
	mlr_dsl_ast_node_t* pleft = pnode->pchildren->phead->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_BLOCK),
			mlr_dsl_ast_node_describe_type(pleft->type));
		exit(1);
	}
	return pleft;
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_resolve_subr_callsites(mlr_dsl_cst_t* pcst) {
	while (pcst->psubr_callsite_statements_to_resolve->phead != NULL) {
		mlr_dsl_cst_statement_t* pstatement = sllv_pop(pcst->psubr_callsite_statements_to_resolve);
		mlr_dsl_cst_resolve_subr_callsite(pcst, pstatement);
	}
}
