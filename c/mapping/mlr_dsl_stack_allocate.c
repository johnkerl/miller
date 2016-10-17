#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/free_flags.h"
#include "containers/lhmsi.h"
#include "mapping/mlr_dsl_blocked_ast.h"
#include "mapping/context_flags.h"

// xxx check 'frame' vs. 'subframe' throughout in comments as well as code

// ================================================================
// This is a two-pass stack allocator for the Miller DSL.
//
// ----------------------------------------------------------------
// CONTEXT:
//
// In the initial Miller implementation of local variables (e.g.  for-loop
// indices or frame-local variables), I used a hashmap from string name to
// mlrval at each level. This was very easy to code and it worked, but it had
// the property that every local-variable read or write involved a hashmap
// lookup to locate each local variable on the stack.  For compute-intensive
// work this resulted in 80% or more of the compute time being used for the
// hashmap accesses. It doesn't make sense to always be asking "where is
// variable 'a'"? at runtime (maybe ten million times) since this can be figured
// out ahead of time.
//
// The Miller DSL allows for recursive functions and subroutines, but within
// those, stack layout is knowable at parse time.
//
// ----------------------------------------------------------------
// EXAMPLE:
//
//                       # ---- FUNC FRAME: defcount 7 {absent-RHS,a,b,c,i,j,y}
//                       # To be noted below: absent-RHS is at slot 0 of top level.
// func f(a, b, c) {     # Args define locals 1,2,3 at current level.
//     local i = 24;     # Explicitly define local 4 at current level.
//     j = 25;           # Implicitly define local 5 at current level.
//                       #
//                       # ---- IF FRAME: defcount 1 {k}
//     if (a == 26) {    # Read local 1, up 1 level.
//         local k = 27; # Explicitly define local 0 at this level.
//         j = 28;       # LHS is local 5 up one level.
//                       #
//                       #
//     } else {          # ---- ELSE FRAME: defcount 1 {n}
//         n = b;        # Implicitly define local 0 at this level.
//     }                 #
//                       #
//     y = 7;            # LHS is local 6 at current level.
//     i = z;            # LHS is local 4 at current level;
//                       #   RHS is unresolved -> slot 0 at current level.
// }                     #
//
// Notes:
//
// * Pass 1 computes frame-relative indices and upstack-level counts,
//   as in the example, for each local variable.
//
// * Pass 2 computes absolute indices for each local variable. These
//   aren't computable in pass 1 due to the example 'y = 7' assignment
//   above: the number of local variables in an upper level can change
//   after the invocation of a child level, so total frame size is not
//   known until all AST nodes in the top-level block have been visited.
//
// * Pass 2 also computes the max depth, counting number of variables, so
//   that for each top-level block we can allocate an array of mlrvals which
//   will be reused on every invocation. (For recursive function calls this will
//   be dynamically allocated.)
//
// * Slot 0 of the top level is reserved for an absent-null for unresolved
//   names on reads.
//
// * The tree-traversal order is done correctly so that if a variable is read
//   before it is defined, then read again after it is defined, then the first
//   read gets absent-null and the second gets the defined value. This also
//   requires the concrete-syntax-tree implementation to initialize the
//   stack to mv_absent on each invocation.
// ================================================================

// ----------------------------------------------------------------
// xxx to do:

// * maybe move ast from containers to mapping?

// * 'semantic analysis': use this to describe CST-build-time checks
// * 'object binding': use this to describe linking func/subr defs and callsites
// * separate verbosity for allocator? and invoke it in UT cases specific to this?
//   -> (note allocation marks in the AST will be printed regardless)

// ================================================================
// Pass-1 stack-frame container: simply a hashmap from name to position on the
// frame relative to the curly-braced statement block (top-level, for-loop,
// if-statement, else-statement, etc.).

typedef struct _stkalc_subframe_t {
	int var_count;
	lhmsi_t* pnames_to_indices;
} stkalc_subframe_t;

// ----------------------------------------------------------------
// Pass-1 stack-frame methods

static      stkalc_subframe_t* stkalc_subframe_alloc();
static void stkalc_subframe_free(stkalc_subframe_t* pframe);
static int  stkalc_subframe_test_and_get(stkalc_subframe_t* pframe, char* name, int* pvalue);
static int  stkalc_subframe_get(stkalc_subframe_t* pframe, char* name);
static int  stkalc_subframe_add(stkalc_subframe_t* pframe, char* name);

// ================================================================
// Pass-1 frame-group container: a linked list with current frame at the head
// and top-level frame at the tail. Within a given top-level block there is a
// tree of curly-braced statement blocks, e.g. a function-definition might have
// an if-statement, else-if, else-if, else, each with its own frame. But during
// pass 1 we only maintain a list from the current frame being analyzed up to
// its parents; sibling branches are not simultaneously stored in this data
// structure.

typedef struct _stkalc_subframe_group_t {
	sllv_t* plist;
} stkalc_subframe_group_t;

// ----------------------------------------------------------------
// Pass-1 stack-frame-group methods

static      stkalc_subframe_group_t* stkalc_subframe_group_alloc(stkalc_subframe_t* pframe, int trace);
static void stkalc_subframe_group_free(stkalc_subframe_group_t* pframe_group);
static void stkalc_subframe_group_push(stkalc_subframe_group_t* pframe_group, stkalc_subframe_t* pframe);
static stkalc_subframe_t* stkalc_subframe_group_pop(stkalc_subframe_group_t* pframe_group);

// Pass-1 stack-frame-group node-mutator methods: given an AST node containing a
// local-variable usage they assign a frame-relative index and a frame-depth
// counter (how many frames deep into the top-level statement block the node
// is).

static void stkalc_subframe_group_mutate_node_for_define(stkalc_subframe_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int trace);

static void stkalc_subframe_group_mutate_node_for_write(stkalc_subframe_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int trace);

static void stkalc_subframe_group_mutate_node_for_read(stkalc_subframe_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int trace);

// ================================================================
// Pass-1 helper methods for the main entry point to this file.

static void pass_1_for_func_subr_block(mlr_dsl_ast_node_t* pnode, int trace);
static void pass_1_for_begin_end_block(mlr_dsl_ast_node_t* pnode, int trace);
static void pass_1_for_main_block(mlr_dsl_ast_node_t* pnode, int trace);
static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_statement_list(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);

static void pass_1_for_local_definition(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_local_assignment(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_local_read(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_srec_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_oosvar_key_only_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_oosvar_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_triple_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);
static void pass_1_for_non_terminal_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace);

// Pass-2 helper methods for the main entry point to this file.
static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode, int trace);
static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int frame_depth, int var_count_below_subframe, int var_count_at_subframe, int* pmax_var_depth, int trace);

// ================================================================
// Utility methods
static void leader_print(int depth);

// ================================================================
// Main entry point for the bind-stack allocator

void blocked_ast_allocate_locals(blocked_ast_t* paast, int trace) {

	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_func_subr_block(pe->pvvalue, trace);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_func_subr_block(pe->pvvalue, trace);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_begin_end_block(pe->pvvalue, trace);
	}
	{
		pass_1_for_main_block(paast->pmain_block, trace);
	}
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_begin_end_block(pe->pvvalue, trace);
	}

	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue, trace);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue, trace);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue, trace);
	}
	{
		pass_2_for_top_level_block(paast->pmain_block, trace);
	}
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue, trace);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_func_subr_block(mlr_dsl_ast_node_t* pnode, int trace) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_SUBR_DEF && pnode->type != MD_AST_NODE_TYPE_FUNC_DEF);
	// xxx assert two children of desired type

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);

	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	}
	mlr_dsl_ast_node_t* pdef_name_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* plist_node = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pdef_name_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
		stkalc_subframe_group_mutate_node_for_define(pframe_group, pparameter_node, "PARAMETER", trace);
	}
	pass_1_for_statement_block(plist_node, pframe_group, trace);
	pnode->subframe_var_count = pframe->var_count;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_begin_end_block(mlr_dsl_ast_node_t* pnode, int trace) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_BEGIN && pnode->type != MD_AST_NODE_TYPE_END);

	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR %s BLOCK\n", pnode->text);
	}

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);

	pass_1_for_statement_block(pnode->pchildren->phead->pvvalue, pframe_group, trace);
	pnode->subframe_var_count = pframe->var_count;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_main_block(mlr_dsl_ast_node_t* pnode, int trace) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK);

	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR MAIN BLOCK\n");
	}

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);

	pass_1_for_statement_block(pnode, pframe_group, trace);
	pnode->subframe_var_count = pframe->var_count;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
// Curly-braced bodies of if/while/for/etc.

static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, trace);
	}
}

// Non-curly-braced triple-for starts/continuations/updates statement lists.
static void pass_1_for_statement_list(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_LIST);
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, trace);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	if (pnode->type == MD_AST_NODE_TYPE_LOCAL_DEFINITION) {
		pass_1_for_local_definition(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_LOCAL_ASSIGNMENT) {
		pass_1_for_local_assignment(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_LOCAL_VARIABLE) { // RHS
		pass_1_for_local_read(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_SREC) {
		pass_1_for_srec_for_loop(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY) { // xxx comment
		pass_1_for_oosvar_key_only_for_loop(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_OOSVAR) { // xxx comment
		pass_1_for_oosvar_for_loop(pnode, pframe_group, trace);
	} else if (pnode->type == MD_AST_NODE_TYPE_TRIPLE_FOR) { // xxx comment
		pass_1_for_triple_for_loop(pnode, pframe_group, trace);
	} else if (pnode->pchildren != NULL) {
		pass_1_for_non_terminal_node(pnode, pframe_group, trace);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_local_definition(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;

	mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
	pass_1_for_node(pvaluenode, pframe_group, trace);
	// Do the LHS after the RHS, in case 'local nonesuch = nonesuch'
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pnamenode, "DEFINE", trace);
}

// ----------------------------------------------------------------
static void pass_1_for_local_assignment(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
	pass_1_for_node(pvaluenode, pframe_group, trace);
	// Do the LHS after the RHS, in case 'local nonesuch = nonesuch'
	stkalc_subframe_group_mutate_node_for_write(pframe_group, pnamenode, "WRITE", trace);
}

// ----------------------------------------------------------------
static void pass_1_for_local_read(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	stkalc_subframe_group_mutate_node_for_read(pframe_group, pnode, "READ", trace);
}

// ----------------------------------------------------------------
// for (k,v in $*) { ... }: the k and v are scoped to the curly-brace block.
static void pass_1_for_srec_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH FRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);

	mlr_dsl_ast_node_t* pvarsnode  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pblocknode = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pknode = pvarsnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvnode = pvarsnode->pchildren->phead->pnext->pvvalue;
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pknode, "FOR-BIND", trace);
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pvnode, "FOR-BIND", trace);

	pass_1_for_statement_block(pblocknode, pframe_group, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP FRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_oosvar_key_only_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {

	mlr_dsl_ast_node_t* pkeynode     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pkeylistnode = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblocknode   = pnode->pchildren->phead->pnext->pnext->pvvalue;

	// The keylistnode is outside the block binding. In particular if there
	// are any localvar reads in there, they shouldn't read from forloop
	// boundvars.
	//
	// Example: 'for(a in @b[c][d]) { local e = a}': the c and d
	// should be obtained from the enclosing scope.
	for (sllve_t* pe = pkeylistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, trace);
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH FRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);

	stkalc_subframe_group_mutate_node_for_define(pframe_group, pkeynode, "FOR-BIND", trace);
	pass_1_for_statement_block(pblocknode, pframe_group, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP FRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_oosvar_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {

	mlr_dsl_ast_node_t* pvarsnode    = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pkeylistnode = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblocknode   = pnode->pchildren->phead->pnext->pnext->pvvalue;

	mlr_dsl_ast_node_t* pkeysnode    = pvarsnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvalnode     = pvarsnode->pchildren->phead->pnext->pvvalue;

	// The keylistnode is outside the block binding. In particular if there
	// are any localvar reads in there, they shouldn't read from forloop
	// boundvars.
	//
	// Example: 'for(k, v in @a[b][c]) { local d = k; local e = v }': the b and c
	// should be obtained from the enclosing scope.
	for (sllve_t* pe = pkeylistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, trace);
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH FRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);

	for (sllve_t* pe = pkeysnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		stkalc_subframe_group_mutate_node_for_define(pframe_group, pkeynode, "FOR-BIND", trace);
	}
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pvalnode, "FOR-BIND", trace);
	pass_1_for_statement_block(pblocknode, pframe_group, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP FRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_triple_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	mlr_dsl_ast_node_t* pstarts_node        = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pcontinuations_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pupdates_node       = pnode->pchildren->phead->pnext->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblock_node         = pnode->pchildren->phead->pnext->pnext->pnext->pvvalue;

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH FRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);

	pass_1_for_statement_list(pstarts_node, pframe_group, trace);
	pass_1_for_statement_list(pcontinuations_node, pframe_group, trace);
	pass_1_for_statement_list(pupdates_node, pframe_group, trace);
	pass_1_for_statement_block(pblock_node, pframe_group, trace);

	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP FRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_non_terminal_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group, int trace) {
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;

		if (pchild->type == MD_AST_NODE_TYPE_STATEMENT_BLOCK) {

			if (trace) {
				leader_print(pframe_group->plist->length);
				printf("PUSH FRAME %s\n", pchild->text);
			}

			stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
			stkalc_subframe_group_push(pframe_group, pnext_subframe);

			pass_1_for_statement_block(pchild, pframe_group, trace);
			pchild->subframe_var_count = pnext_subframe->var_count;

			stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));

			if (trace) {
				leader_print(pframe_group->plist->length);
				printf("POP FRAME %s subframe_var_count=%d\n", pnode->text, pchild->subframe_var_count);
			}

		} else {
			pass_1_for_node(pchild, pframe_group, trace);
		}
	}
}

// ================================================================
static stkalc_subframe_t* stkalc_subframe_alloc() {
	stkalc_subframe_t* pframe = mlr_malloc_or_die(sizeof(stkalc_subframe_t));
	pframe->var_count = 0;
	pframe->pnames_to_indices = lhmsi_alloc();
	return pframe;
}

static void stkalc_subframe_free(stkalc_subframe_t* pframe) {
	if (pframe == NULL)
		return;
	lhmsi_free(pframe->pnames_to_indices);
	free(pframe);
}

static int stkalc_subframe_test_and_get(stkalc_subframe_t* pframe, char* name, int* pvalue) {
	return lhmsi_test_and_get(pframe->pnames_to_indices, name, pvalue);
}

static int stkalc_subframe_get(stkalc_subframe_t* pframe, char* name) {
	return lhmsi_get(pframe->pnames_to_indices, name);
}

static int stkalc_subframe_add(stkalc_subframe_t* pframe, char* name) {
	int rv = pframe->var_count;
	lhmsi_put(pframe->pnames_to_indices, name, pframe->var_count, NO_FREE);
	pframe->var_count++;
	return rv;
}

// ================================================================
static stkalc_subframe_group_t* stkalc_subframe_group_alloc(stkalc_subframe_t* pframe, int trace) {
	stkalc_subframe_group_t* pframe_group = mlr_malloc_or_die(sizeof(stkalc_subframe_group_t));
	pframe_group->plist = sllv_alloc();
	sllv_prepend(pframe_group->plist, pframe);
	stkalc_subframe_add(pframe, ""); // xxx add to trace
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("ADD FOR ABSENT s @ %du%d\n", 0, 0);
	}
	return pframe_group;
}

static void stkalc_subframe_group_free(stkalc_subframe_group_t* pframe_group) {
	if (pframe_group == NULL)
		return;
	while (pframe_group->plist->phead != NULL) {
		stkalc_subframe_free(sllv_pop(pframe_group->plist));
	}
	sllv_free(pframe_group->plist);
	free(pframe_group);
}

static void stkalc_subframe_group_push(stkalc_subframe_group_t* pframe_group, stkalc_subframe_t* pframe) {
	sllv_prepend(pframe_group->plist, pframe);
}

static stkalc_subframe_t* stkalc_subframe_group_pop(stkalc_subframe_group_t* pframe_group) {
	return sllv_pop(pframe_group->plist);
}

static void stkalc_subframe_group_mutate_node_for_define(stkalc_subframe_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int trace)
{
	char* op = "REUSE";
	stkalc_subframe_t* pframe = pframe_group->plist->phead->pvvalue;
	pnode->upstack_subframe_count = 0;
	if (!stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->subframe_relative_index)) {
		pnode->subframe_relative_index = stkalc_subframe_add(pframe, pnode->text);
		op = "ADD";
	}
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("%s %s %s @ %du%d\n", op, desc, pnode->text,
			pnode->subframe_relative_index, pnode->upstack_subframe_count);
	}
}

static void stkalc_subframe_group_mutate_node_for_write(stkalc_subframe_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int trace)
{
	char* op = "REUSE";
	int found = FALSE;
	// xxx comment: re loop. if not found, fall back to top frame.
	pnode->upstack_subframe_count = 0;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, pnode->upstack_subframe_count++) {
		stkalc_subframe_t* pframe = pe->pvvalue;
		if (stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->subframe_relative_index)) {
			found = TRUE; // xxx dup
			break;
		}
	}

	if (!found) {
		pnode->upstack_subframe_count = 0;
		stkalc_subframe_t* pframe = pframe_group->plist->phead->pvvalue;
		pnode->subframe_relative_index = stkalc_subframe_add(pframe, pnode->text);
		// xxx temp
		op = "ADD";
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("%s %s %s @ %du%d\n", op, desc, pnode->text,
			pnode->subframe_relative_index, pnode->upstack_subframe_count);
	}
}

// xxx make this very clear in the header somehow ... this is an assumption to be tracked across modules.
static void stkalc_subframe_group_mutate_node_for_read(stkalc_subframe_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int trace)
{
	char* op = "PRESENT";
	int found = FALSE;
	// xxx comment: re loop. if not found, fall back to top frame.
	int upstack_subframe_count = 0;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, upstack_subframe_count++) {
		stkalc_subframe_t* pframe = pe->pvvalue;
		if (stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->subframe_relative_index)) {
			found = TRUE; // xxx dup
			break;
		}
	}

	// xxx if not found: go to the tail & use the "" entry
	if (!found) {
		stkalc_subframe_t* plast = pframe_group->plist->ptail->pvvalue;
		pnode->subframe_relative_index = stkalc_subframe_get(plast, "");
		upstack_subframe_count = pframe_group->plist->length - 1;
		op = "ABSENT";
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("%s %s %s @ %du%d\n", desc, pnode->text, op, pnode->subframe_relative_index, upstack_subframe_count);
	}
	pnode->upstack_subframe_count = upstack_subframe_count;

}

// ================================================================
static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode, int trace) {
	int frame_depth = 0;
	int var_count_below_subframe = 0;
	int var_count_at_subframe = 0;
	int max_var_depth   = 0;
	if (trace) {
		printf("\n");
		printf("ALLOCATING ABSOLUTE (PASS-2) LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
		printf("\n");
		mlr_dsl_ast_node_print(pnode);
		printf("\n");
	}
	pass_2_for_node(pnode, frame_depth, var_count_below_subframe, var_count_at_subframe, &max_var_depth, trace);
	pnode->max_var_depth = max_var_depth;
}

static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int frame_depth, int var_count_below_subframe, int var_count_at_subframe, int* pmax_var_depth, int trace)
{
	if (pnode->subframe_var_count != MD_UNUSED_INDEX) {
		var_count_below_subframe += var_count_at_subframe;
		var_count_at_subframe = pnode->subframe_var_count;
		int depth = var_count_below_subframe + var_count_at_subframe;
		frame_depth++;
		if (depth > *pmax_var_depth) {
			*pmax_var_depth = depth;
		}
		if (trace) {
			leader_print(frame_depth-1);
			printf("FRAME [%s] var_count_below=%d var_count_at=%d max_var_depth_so_far=%d\n",
				pnode->text, var_count_below_subframe, var_count_at_subframe, *pmax_var_depth);
		}
	}

	if (pnode->subframe_relative_index != MD_UNUSED_INDEX) {

		pnode->frame_relative_index = var_count_below_subframe + pnode->subframe_relative_index; // R-O-N-G

		if (trace) {
			leader_print(frame_depth);
			printf("NODE %s %du%d -> %d\n",
				pnode->text,
				pnode->subframe_relative_index,
				pnode->upstack_subframe_count,
				pnode->frame_relative_index);
		}
		MLR_INTERNAL_CODING_ERROR_IF(pnode->frame_relative_index <= 0);
		MLR_INTERNAL_CODING_ERROR_IF(pnode->frame_relative_index >= *pmax_var_depth);
	}

	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_2_for_node(pchild, frame_depth,
				var_count_below_subframe, var_count_at_subframe, pmax_var_depth, trace);
		}
	}
}

// ================================================================
const char* STKALC_TRACE_LEADER = "    ";
static void leader_print(int depth) {
	for (int i = 0; i < depth; i++)
		printf("%s", STKALC_TRACE_LEADER);
}
