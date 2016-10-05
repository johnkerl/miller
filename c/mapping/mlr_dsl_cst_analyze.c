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
	paast->pmain_block   = mlr_dsl_ast_node_alloc_zary("main_block", MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	paast->pend_blocks   = sllv_alloc();

	if (past->proot->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected root node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_BLOCK),
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
// xxx under construction

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
//         g = n;        #
//         h = b         #
//     }                 #
//                       #
//     b = 7;            # LHS A.2
//     i = z;            # LHS A.4 RHS unresolved
// }                     #
// ================================================================

// ----------------------------------------------------------------
// * local-var defines are certainly for the current frame
// * local-var writes need backtracing (if not found at the current frame)
// * local-var reads  need backtracing (if not found at the current frame)
// * unresolved-read needs special handling -- maybe a root-level mv_absent at index 0?
//
// ----------------------------------------------------------------
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
//
// ----------------------------------------------------------------
// Population:
// * in-order AST traversal
// * note statement-list nodes are only every so often in the full AST
// * at each node:
//     if is local-var LHS:
//       if explicit:
//         lhmsi_put(name, ++fridx)
//       else:
//         resolve up ...
//     else if is statement-list:
//       allocate a frame struct
//       attach it to the node
//       recurse & have the recursion populate it
//       pop the frame struct but leave it attached to the node
//     else:
//       nothing to do here.
// ----------------------------------------------------------------

// func f(a,b) {
//     return a+b;
// }
// subr s(n) {
//     print n;
// }
// begin {
//     local x = 1;
// }
// end {
//     local z = 3;
// }
// local y = 2;

// ANALYZED AST:
//
// FUNCTION DEFINITION:
// text="f", type=FUNC_DEF:
//     text="f", type=NON_SIGIL_NAME:
//         text="a", type=NON_SIGIL_NAME.
//         text="b", type=NON_SIGIL_NAME.
//     text="list", type=STATEMENT_BLOCK:
//         text="return_value", type=RETURN_VALUE:
//             text="+", type=OPERATOR:
//                 text="a", type=BOUND_VARIABLE.
//                 text="b", type=BOUND_VARIABLE.
//
// SUBROUTINE DEFINITION:
// text="s", type=SUBR_DEF:
//     text="s", type=NON_SIGIL_NAME:
//         text="n", type=NON_SIGIL_NAME.
//     text="list", type=STATEMENT_BLOCK:
//         text="print", type=PRINT:
//             text="n", type=BOUND_VARIABLE.
//             text=">", type=FILE_WRITE:
//                 text="stdout", type=STDOUT:
//
// BEGIN-BLOCK:
// text="begin", type=BEGIN:
//     text="list", type=STATEMENT_BLOCK:
//         text="local", type=LOCAL:
//             text="x", type=BOUND_VARIABLE.
//             text="1", type=STRNUM_LITERAL.
//
// END-BLOCK:
// text="end", type=END:
//     text="list", type=STATEMENT_BLOCK:
//         text="local", type=LOCAL:
//             text="z", type=BOUND_VARIABLE.
//             text="3", type=STRNUM_LITERAL.
//
// MAIN BLOCK:
// text="main_block", type=STATEMENT_BLOCK:
//     text="local", type=LOCAL:
//         text="y", type=BOUND_VARIABLE.
//         text="2", type=STRNUM_LITERAL.


// ================================================================
static void analyzed_ast_allocate_locals_for_func_subr_block(mlr_dsl_ast_node_t* pnode);
static void analyzed_ast_allocate_locals_for_begin_end_block(mlr_dsl_ast_node_t* pnode);
static void analyzed_ast_allocate_locals_for_main_block(mlr_dsl_ast_node_t* pnode);
static void analyzed_ast_allocate_locals_for_statement_list(mlr_dsl_ast_node_t* pnode, sllv_t* pframe_group);
static void analyzed_ast_allocate_locals_for_node(mlr_dsl_ast_node_t* pnode, sllv_t* pframe_group);

// ----------------------------------------------------------------
void analyzed_ast_allocate_locals(analyzed_ast_t* paast) {
	printf("\n"); // xxx temp
	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		analyzed_ast_allocate_locals_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		analyzed_ast_allocate_locals_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		analyzed_ast_allocate_locals_for_begin_end_block(pe->pvvalue);
	}
	analyzed_ast_allocate_locals_for_main_block(paast->pmain_block);
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		analyzed_ast_allocate_locals_for_begin_end_block(pe->pvvalue);
	}
}

// ----------------------------------------------------------------
static void analyzed_ast_allocate_locals_for_func_subr_block(mlr_dsl_ast_node_t* pnode) {
	// xxx make a keystroke-saver, use it here, & use it from the cst builder as well
	if (pnode->type != MD_AST_NODE_TYPE_SUBR_DEF && pnode->type != MD_AST_NODE_TYPE_FUNC_DEF) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	//	xxx assert two children of desired type

	long long index_count = 0;

	lhmsi_t* pnames_to_indices = lhmsi_alloc();
	sllv_t* pframe_group = sllv_alloc();
	sllv_prepend(pframe_group, pnames_to_indices);

	printf("ALLOCATING LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	mlr_dsl_ast_node_t* pdef_name_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* plist_node = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pdef_name_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
		if (!lhmsi_has_key(pnames_to_indices, pparameter_node->text)) {
			// xxx wrap in a class
			printf("ALLOCATING PARAMETER %s = %lld\n", pparameter_node->text, index_count);
			lhmsi_put(pnames_to_indices, pparameter_node->text, index_count, NO_FREE);
			index_count++;
		}
	}
	analyzed_ast_allocate_locals_for_statement_list(plist_node, pframe_group);

	sllv_pop(pframe_group);
	sllv_free(pframe_group);
	lhmsi_free(pnames_to_indices);
}

// ----------------------------------------------------------------
static void analyzed_ast_allocate_locals_for_begin_end_block(mlr_dsl_ast_node_t* pnode) {
	if (pnode->type != MD_AST_NODE_TYPE_BEGIN && pnode->type != MD_AST_NODE_TYPE_END) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	printf("ALLOCATING LOCALS FOR %s BLOCK\n", pnode->text);
	lhmsi_t* pnames_to_indices = lhmsi_alloc();
	sllv_t* pframe_group = sllv_alloc();
	sllv_prepend(pframe_group, pnames_to_indices);

	analyzed_ast_allocate_locals_for_statement_list(pnode->pchildren->phead->pvvalue, pframe_group);

	sllv_pop(pframe_group);
	sllv_free(pframe_group);
	lhmsi_free(pnames_to_indices);

}

// ----------------------------------------------------------------
static void analyzed_ast_allocate_locals_for_main_block(mlr_dsl_ast_node_t* pnode) {
//	xxx assert node type

	printf("ALLOCATING LOCALS FOR MAIN BLOCK\n");
	lhmsi_t* pnames_to_indices = lhmsi_alloc();
	sllv_t* pframe_group = sllv_alloc();
	sllv_prepend(pframe_group, pnames_to_indices);

	analyzed_ast_allocate_locals_for_statement_list(pnode, pframe_group);

	sllv_pop(pframe_group);
	sllv_free(pframe_group);
	lhmsi_free(pnames_to_indices);
}

// ----------------------------------------------------------------
// xxx this becomes easier (and less contextful) if there are separate BOUNDVAR ast node types
// for boundvar @ rhs, boundvar @ for-loop bind, @ arg bind, @ local-def, and @ lhs/assign.
static void analyzed_ast_allocate_locals_for_statement_list(mlr_dsl_ast_node_t* pnode, sllv_t* pframe_group) {
	if (pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		analyzed_ast_allocate_locals_for_node(pchild, pframe_group);
	}
}

static void analyzed_ast_allocate_locals_for_node(mlr_dsl_ast_node_t* pnode, sllv_t* pframe_group) {
	if (pnode->type == MD_AST_NODE_TYPE_BOUND_VARIABLE) {
		// make method
		if (!lhmsi_has_key(pframe_group->phead->pvvalue, pnode->text)) {
			// xxx track count
			lhmsi_put(pframe_group->phead->pvvalue, pnode->text, 999, NO_FREE);
			for (int i = 0; i < pframe_group->length; i++)
				printf("::  ");
			printf("BOOP [%s]!\n", pnode->text);
		}
	}
	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;

			// xxx special case for triple-for: only the body statement list is curly braced.
			// the triple elements are not. maybe make some ast-node-type help here in the parser?
			if (pchild->type == MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
				lhmsi_t* pnames_to_indices = lhmsi_alloc();
				for (int i = 0; i < pframe_group->length; i++)
					printf("::  ");
				printf("PUSH FRAME %s\n", pchild->text);
				sllv_prepend(pframe_group, pnames_to_indices);

				analyzed_ast_allocate_locals_for_statement_list(pchild, pframe_group);

				sllv_pop(pframe_group);
				for (int i = 0; i < pframe_group->length; i++)
					printf("::  ");
				printf("POP FRAME %s\n", pchild->text);
				lhmsi_free(pnames_to_indices);
			} else {
				analyzed_ast_allocate_locals_for_node(pchild, pframe_group);
			}
		}
	}
}
