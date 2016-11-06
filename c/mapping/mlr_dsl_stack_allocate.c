#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/free_flags.h"
#include "containers/lhmsi.h"
#include "mapping/mlr_dsl_blocked_ast.h"
#include "mapping/context_flags.h"

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
// variable 'a'"? at runtime (maybe ten million times doing hash-map lookups)
// since this can be figured out ahead of time.
//
// The Miller DSL allows for recursive functions and subroutines, but within
// those, stack layout is knowable at parse time.
//
// ----------------------------------------------------------------
// TERMINOLOGY
//
// * A top-level statement block starts with 'begin', 'end', 'func', 'subr',
//   or is the remaining collection of statements outside of any of those which
//   is called the 'main' statement block.
//
// * A stack frame is all the locals for a top-level statement block including
//   anything local to a scoped block within the top-level block. For example,
//   locals may be defined within if/while/do-while blocks. Also, variables can
//   be local to a triple-for, e.g. 'for (var i = 0; i < 10; i += 1) { ... }'.
//   As well, locals are bound to variables in for-loops over stream records
//   or out-of-stream variables: 'for (k, v in $*) { ... }' or
//   'for ((k1, k2), v in @*) { ... }'. Lastly, function arguments are local
//   to the function/subroutine definition site.
//
// * A subframe is the locals defined within any of the above: locals
//   defined within an if/while/for/etc. block. Or, the locals defined
//   in the top-level block have their own subframe.
//
// ----------------------------------------------------------------
// CONVENTION: Local variables are bound to indices within the stack frame, but,
// local variable may be read before they are defined.  Slot 0 of each frame is
// an absent-null which is reserved for this purpose.
//
// ----------------------------------------------------------------
// EXAMPLE:
//
//                       # ---- FUNC TOP-LEVEL SUBFRAME 0: defcount 7 {absent-RHS,a,b,c,i,j,y}
//                       # Absent-null-RHS is at slot 0 of top level.
// func f(a, b, c) {     # Args define locals 1,2,3 at current level.
//     var i = 24;       # Explicitly define local 4 at current level.
//     j = 25;           # Implicitly define local 5 at current level.
//                       #
//                       # ---- IF-STATEMENT SUBFRAME 1: defcount 1 {k}
//     if (a == 26) {    # Read of local 1, found at top level (subframe 0).
//         var k = 27;   # Explicitly define local 0 at this level.
//         j = 28;       # LHS is local 5, found at top level.
//                       #
//                       # Note that the 'if' and the 'else' are both at subframe
//                       # depth 1, while each defines a different number of locals.
//                       # For this function the max variable depth is 9:
//                       # 7 at top level plus 1 = 8 if the if is taken,
//                       # 7 at top level plus 2 = 9 if the else is taken,
//                       #
//     } else {          # ---- ELSE-STATEMENT SUBFRAME 1: defcount 2 {n, u}
//         n = b;        # Implicitly define local 0 at this level.
//         var u = 4;    # Implicitly define local 1 at this level.
//     }                 #
//                       #
//     y = 7;            # LHS is local 6 at current level.
//     i = z;            # LHS is local 4 at current level;
//                       # RHS is unresolved -> slot 0 at current level.
// }                     #
//
// Notes:
//
// * Pass 1 computes frame-relative indices and subframe-level counts,
//   as in the example, for each local variable.
//
// * Pass 2 computes absolute indices for each local variable. These
//   aren't computable in pass 1 due to the example 'y = 7' assignment
//   above: the number of local variables in an upper level can change
//   after the invocation of a child level, so the total frame size is not
//   known until all AST nodes in the top-level block have been visited.
//
// * Pass 2 also computes the max depth, counting number of variables, so
//   that for each top-level block we can allocate an array of mlrvals which
//   will be reused on every invocation. (For recursive function calls an entire
//   frame be dynamically allocated.)
//
// * Slot 0 of the top level is reserved for an absent-null for unresolved
//   names on reads.
//
// * The tree-traversal order is done correctly so that if a variable is read
//   before it is defined, then read again after it is defined, then the first
//   read gets absent-null and the second gets the defined value.
// ================================================================

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

static stkalc_subframe_t* stkalc_subframe_alloc();

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
// local-variable usage they assign a subframe-relative index and a subframe-depth
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

static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_statement_list(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);


static void pass_1_for_local_definition(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_local_assignment(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_local_read(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_srec_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_map_key_only_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_map_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_triple_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

static void pass_1_for_non_terminal_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace);

// Pass-2 helper methods for the main entry point to this file.
static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode, int trace);

static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int subframe_depth, int var_count_below_subframe, int var_count_at_subframe, int* pmax_var_depth,
	int* subframe_var_count_belows, int max_subframe_depth,
	int trace);

// ================================================================
// Utility methods
static void leader_print(int depth);

// ================================================================
// Main entry point for the bind-stack allocator

void blocked_ast_allocate_locals(blocked_ast_t* paast, int trace) {

	// Pass 1
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

	// Pass 2
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
	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_SUBR_DEF && pnode->type != MD_AST_NODE_TYPE_FUNC_DEF);

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);
	int max_subframe_depth = 1;

	mlr_dsl_ast_node_t* pdef_name_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* plist_node = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pdef_name_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
		stkalc_subframe_group_mutate_node_for_define(pframe_group, pparameter_node, "PARAMETER", trace);
	}
	pass_1_for_statement_block(plist_node, pframe_group, &max_subframe_depth, trace);
	pnode->subframe_var_count = pframe->var_count;
	pnode->max_subframe_depth = max_subframe_depth;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d max_subframe_depth=%d\n",
			pnode->text, pnode->subframe_var_count, pnode->max_subframe_depth);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_begin_end_block(mlr_dsl_ast_node_t* pnode, int trace) {
	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR %s BLOCK\n", pnode->text);
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_BEGIN && pnode->type != MD_AST_NODE_TYPE_END);

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);
	int max_subframe_depth = 1;

	pass_1_for_statement_block(pnode->pchildren->phead->pvvalue, pframe_group, &max_subframe_depth, trace);
	pnode->subframe_var_count = pframe->var_count;
	pnode->max_subframe_depth = max_subframe_depth;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d max_subframe_depth=%d\n",
			pnode->text, pnode->subframe_var_count, pnode->max_subframe_depth);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_main_block(mlr_dsl_ast_node_t* pnode, int trace) {
	if (trace) {
		printf("\n");
		printf("ALLOCATING RELATIVE (PASS-1) LOCALS FOR MAIN BLOCK\n");
	}

	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK);

	stkalc_subframe_t* pframe = stkalc_subframe_alloc();
	stkalc_subframe_group_t* pframe_group = stkalc_subframe_group_alloc(pframe, trace);
	int max_subframe_depth = 1;

	pass_1_for_statement_block(pnode, pframe_group, &max_subframe_depth, trace);
	pnode->subframe_var_count = pframe->var_count;
	pnode->max_subframe_depth = max_subframe_depth;
	if (trace) {
		printf("BLOCK %s subframe_var_count=%d max_subframe_depth=%d\n",
			pnode->text, pnode->subframe_var_count, pnode->max_subframe_depth);
	}

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	stkalc_subframe_group_free(pframe_group);
}

// ----------------------------------------------------------------
// Curly-braced bodies of if/while/for/etc.

static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK);
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, pmax_subframe_depth, trace);
	}
}

// Non-curly-braced triple-for starts/continuations/updates statement lists.
static void pass_1_for_statement_list(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	MLR_INTERNAL_CODING_ERROR_IF(pnode->type != MD_AST_NODE_TYPE_STATEMENT_LIST);
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group, pmax_subframe_depth, trace);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	if (pnode->type == MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION) { // LHS
		pass_1_for_local_definition(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT) { // LHS
		pass_1_for_local_assignment(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT) { // LHS
		pass_1_for_local_assignment(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE) { // RHS
		pass_1_for_local_read(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE) { // RHS
		pass_1_for_local_read(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_SREC) {
		pass_1_for_srec_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY) {
		pass_1_for_map_key_only_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_OOSVAR) {
		pass_1_for_map_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_LOCAL_MAP_KEY_ONLY) {
		pass_1_for_map_key_only_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_LOCAL_MAP) {
		pass_1_for_map_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->type == MD_AST_NODE_TYPE_TRIPLE_FOR) {
		pass_1_for_triple_for_loop(pnode, pframe_group, pmax_subframe_depth, trace);

	} else if (pnode->pchildren != NULL) {
		pass_1_for_non_terminal_node(pnode, pframe_group, pmax_subframe_depth, trace);

	}
}

// ----------------------------------------------------------------
static void pass_1_for_local_definition(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;

	// RHS must exist for non-map types ('int x = 3', not 'int x') but must not for map types
	// ('map x', not 'map x = {}').
	if (pnode->pchildren->phead->pnext) {
		mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
		pass_1_for_node(pvaluenode, pframe_group, pmax_subframe_depth, trace);
	}
	// Do the LHS after the RHS, in case 'var nonesuch = nonesuch'
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pnamenode, "DEFINE", trace);
}

// ----------------------------------------------------------------
static void pass_1_for_local_assignment(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
	pass_1_for_node(pvaluenode, pframe_group, pmax_subframe_depth, trace);
	// Do the LHS after the RHS, in case 'var nonesuch = nonesuch'
	stkalc_subframe_group_mutate_node_for_write(pframe_group, pnamenode, "WRITE", trace);
}

// ----------------------------------------------------------------
static void pass_1_for_local_read(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	stkalc_subframe_group_mutate_node_for_read(pframe_group, pnode, "READ", trace);
}

// ----------------------------------------------------------------
// for (k,v in $*) { ... }: the k and v are scoped to the curly-brace block.
static void pass_1_for_srec_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH SUBFRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);
	if (*pmax_subframe_depth < pframe_group->plist->length)
		*pmax_subframe_depth = pframe_group->plist->length;

	mlr_dsl_ast_node_t* pvarsnode  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pblocknode = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pknode = pvarsnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvnode = pvarsnode->pchildren->phead->pnext->pvvalue;
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pknode, "FOR-BIND", trace);
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pvnode, "FOR-BIND", trace);

	pass_1_for_statement_block(pblocknode, pframe_group, pmax_subframe_depth, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP SUBFRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_map_key_only_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	mlr_dsl_ast_node_t* pkeynode     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pkeylistnode = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblocknode   = pnode->pchildren->phead->pnext->pnext->pvvalue;

	// The keylistnode is outside the block binding. In particular if there
	// are any localvar reads in there, they shouldn't read from forloop
	// boundvars.
	//
	// Example: 'for(a in @b[c][d]) { var e = a}': the c and d
	// should be obtained from the enclosing scope.

	if (pkeylistnode->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE) {
		// For-local-map: e.g. 'for(a in b[c][d])'.
		pass_1_for_node(pkeylistnode, pframe_group, pmax_subframe_depth, trace);
	}
	if (pkeylistnode->pchildren != NULL) {
		for (sllve_t* pe = pkeylistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_1_for_node(pchild, pframe_group, pmax_subframe_depth, trace);
		}
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH SUBFRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);
	if (*pmax_subframe_depth < pframe_group->plist->length)
		*pmax_subframe_depth = pframe_group->plist->length;

	stkalc_subframe_group_mutate_node_for_define(pframe_group, pkeynode, "FOR-BIND", trace);
	pass_1_for_statement_block(pblocknode, pframe_group, pmax_subframe_depth, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP SUBFRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_map_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	mlr_dsl_ast_node_t* pvarsnode    = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pkeylistnode = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblocknode   = pnode->pchildren->phead->pnext->pnext->pvvalue;

	mlr_dsl_ast_node_t* pkeysnode    = pvarsnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvalnode     = pvarsnode->pchildren->phead->pnext->pvvalue;

	// The keylistnode is outside the block binding. In particular if there
	// are any localvar reads in there, they shouldn't read from forloop
	// boundvars.
	//
	// Example: 'for(k, v in @a[b][c]) { var d = k; var e = v }': the b and c
	// should be obtained from the enclosing scope.
	if (pkeylistnode->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE) {
		// For-local-map: e.g. 'for(a,b in c[d][e])'.
		pass_1_for_node(pkeylistnode, pframe_group, pmax_subframe_depth, trace);
	}
	if (pkeylistnode->pchildren == NULL) {
		pass_1_for_node(pkeylistnode, pframe_group, pmax_subframe_depth, trace);
	} else {
		for (sllve_t* pe = pkeylistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_1_for_node(pchild, pframe_group, pmax_subframe_depth, trace);
		}
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH SUBFRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);
	if (*pmax_subframe_depth < pframe_group->plist->length)
		*pmax_subframe_depth = pframe_group->plist->length;

	for (sllve_t* pe = pkeysnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		stkalc_subframe_group_mutate_node_for_define(pframe_group, pkeynode, "FOR-BIND", trace);
	}
	stkalc_subframe_group_mutate_node_for_define(pframe_group, pvalnode, "FOR-BIND", trace);
	pass_1_for_statement_block(pblocknode, pframe_group, pmax_subframe_depth, trace);
	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP SUBFRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_triple_for_loop(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	mlr_dsl_ast_node_t* pstarts_node        = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pcontinuations_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pupdates_node       = pnode->pchildren->phead->pnext->pnext->pvvalue;
	mlr_dsl_ast_node_t* pblock_node         = pnode->pchildren->phead->pnext->pnext->pnext->pvvalue;

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("PUSH SUBFRAME %s\n", pnode->text);
	}
	stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
	stkalc_subframe_group_push(pframe_group, pnext_subframe);
	if (*pmax_subframe_depth < pframe_group->plist->length)
		*pmax_subframe_depth = pframe_group->plist->length;

	pass_1_for_statement_list(pstarts_node, pframe_group, pmax_subframe_depth, trace);
	pass_1_for_statement_list(pcontinuations_node, pframe_group, pmax_subframe_depth, trace);
	pass_1_for_statement_list(pupdates_node, pframe_group, pmax_subframe_depth, trace);
	pass_1_for_statement_block(pblock_node, pframe_group, pmax_subframe_depth, trace);

	pnode->subframe_var_count = pnext_subframe->var_count;

	stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("POP SUBFRAME %s subframe_var_count=%d\n", pnode->text, pnode->subframe_var_count);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_non_terminal_node(mlr_dsl_ast_node_t* pnode, stkalc_subframe_group_t* pframe_group,
	int* pmax_subframe_depth, int trace)
{
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;

		if (pchild->type == MD_AST_NODE_TYPE_STATEMENT_BLOCK) {

			if (trace) {
				leader_print(pframe_group->plist->length);
				printf("PUSH SUBFRAME %s\n", pchild->text);
			}

			stkalc_subframe_t* pnext_subframe = stkalc_subframe_alloc();
			stkalc_subframe_group_push(pframe_group, pnext_subframe);
			if (*pmax_subframe_depth < pframe_group->plist->length)
				*pmax_subframe_depth = pframe_group->plist->length;

			pass_1_for_statement_block(pchild, pframe_group, pmax_subframe_depth, trace);
			pchild->subframe_var_count = pnext_subframe->var_count;

			stkalc_subframe_free(stkalc_subframe_group_pop(pframe_group));

			if (trace) {
				leader_print(pframe_group->plist->length);
				printf("POP SUBFRAME %s subframe_var_count=%d\n", pnode->text, pchild->subframe_var_count);
			}

		} else {
			pass_1_for_node(pchild, pframe_group, pmax_subframe_depth, trace);
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
	sllv_push(pframe_group->plist, pframe);
	stkalc_subframe_add(pframe, "");
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
	sllv_push(pframe_group->plist, pframe);
}

static stkalc_subframe_t* stkalc_subframe_group_pop(stkalc_subframe_group_t* pframe_group) {
	return sllv_pop(pframe_group->plist);
}

// 'var x = 1' always applies to the current subframe.
static void stkalc_subframe_group_mutate_node_for_define(stkalc_subframe_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int trace)
{
	stkalc_subframe_t* pframe = pframe_group->plist->phead->pvvalue;
	pnode->vardef_subframe_index = pframe_group->plist->length - 1;
	if (stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->vardef_subframe_relative_index)) {
		fprintf(stderr, "%s: redefinition of variable %s in the same scope.\n",
			MLR_GLOBALS.bargv0, pnode->text);
		exit(1);
	} else {
		pnode->vardef_subframe_relative_index = stkalc_subframe_add(pframe, pnode->text);
	}
	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("ADD %s %s @ %ds%d\n", desc, pnode->text,
			pnode->vardef_subframe_relative_index, pnode->vardef_subframe_index);
	}
}

// 'x = 1' is one of two things: (1) already defined in a higher subframe and
// referenced in the current subframe; (2) not defined in a higher subframe, in
// which case it is hereby defined in the current subframe.
static void stkalc_subframe_group_mutate_node_for_write(stkalc_subframe_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int trace)
{
	char* op = "REUSE";
	int found = FALSE;

	// Search for definitions in current & higher subframes.
	pnode->vardef_subframe_index = pframe_group->plist->length - 1;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, pnode->vardef_subframe_index--) {
		stkalc_subframe_t* pframe = pe->pvvalue;
		if (stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->vardef_subframe_relative_index)) {
			found = TRUE;
			break;
		}
	}

	// If not found, define locally.
	if (!found) {
		pnode->vardef_subframe_index = pframe_group->plist->length - 1;
		stkalc_subframe_t* pframe = pframe_group->plist->phead->pvvalue;
		pnode->vardef_subframe_relative_index = stkalc_subframe_add(pframe, pnode->text);
		op = "ADD";
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("%s %s %s @ %ds%d\n", op, desc, pnode->text,
			pnode->vardef_subframe_relative_index, pnode->vardef_subframe_index);
	}
}

// The right-hand side of '$a = b' is one of two things: (1) already defined in a higher
// subframe and referenced in the current subframe; (2) not defined in a higher subframe,
// in which case the RHS evaluates to absent-null.  An absent-null is always kept at index
// 0 in the frame. This is an important assumption to be tracked across modules, including
// here as well as the CST-node handlers. It's tested in the test-dsl-stack-allocation.mlr
// regtest case.
static void stkalc_subframe_group_mutate_node_for_read(stkalc_subframe_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int trace)
{
	char* op = "PRESENT";
	int found = FALSE;

	// Search for definitions in current & higher subframes.
	pnode->vardef_subframe_index = pframe_group->plist->length - 1;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, pnode->vardef_subframe_index--) {
		stkalc_subframe_t* pframe = pe->pvvalue;
		if (stkalc_subframe_test_and_get(pframe, pnode->text, &pnode->vardef_subframe_relative_index)) {
			found = TRUE;
			break;
		}
	}

	// Absent-null is indexed in this stack allocator by the "" variable name.
	if (!found) {
		stkalc_subframe_t* plast = pframe_group->plist->ptail->pvvalue;
		pnode->vardef_subframe_relative_index = stkalc_subframe_get(plast, "");
		pnode->vardef_subframe_index = 0;
		op = "ABSENT";
	}

	if (trace) {
		leader_print(pframe_group->plist->length);
		printf("%s %s %s @ %ds%d\n", desc, pnode->text, op,
			pnode->vardef_subframe_relative_index, pnode->vardef_subframe_index);
	}
}

// ================================================================
static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode, int trace) {
	int subframe_depth           = 0;
	int var_count_below_subframe = 0;
	int var_count_at_subframe    = 0;
	int max_var_depth            = 0;

	int max_subframe_depth = pnode->max_subframe_depth;
	MLR_INTERNAL_CODING_ERROR_IF(max_subframe_depth == MD_UNUSED_INDEX);
	int* subframe_var_count_belows = mlr_malloc_or_die(max_subframe_depth * sizeof(int));
	for (int i = 0; i < pnode->max_subframe_depth; i++)
		subframe_var_count_belows[i] = MD_UNUSED_INDEX;


	if (trace) {
		printf("\n");
		printf("ALLOCATING ABSOLUTE (PASS-2) LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	}
	pass_2_for_node(pnode, subframe_depth, var_count_below_subframe, var_count_at_subframe, &max_var_depth,
		subframe_var_count_belows, max_subframe_depth, trace);
	pnode->max_var_depth = max_var_depth;
	free(subframe_var_count_belows);
}

// ----------------------------------------------------------------
static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int subframe_depth, int var_count_below_subframe, int var_count_at_subframe, int* pmax_var_depth,
	int* subframe_var_count_belows, int max_subframe_depth,
	int trace)
{
	if (pnode->subframe_var_count != MD_UNUSED_INDEX) {
		var_count_below_subframe += var_count_at_subframe;
		var_count_at_subframe = pnode->subframe_var_count;
		int depth = var_count_below_subframe + var_count_at_subframe;
		if (depth > *pmax_var_depth) {
			*pmax_var_depth = depth;
		}
		subframe_var_count_belows[subframe_depth] = var_count_below_subframe;
		if (trace) {
			leader_print(subframe_depth);
			printf("SUBFRAME [%s] var_count_below=%d var_count_at=%d max_var_depth_so_far=%d subframe_depth=%d\n",
				pnode->text, var_count_below_subframe, var_count_at_subframe, *pmax_var_depth, subframe_depth);
		}
		subframe_depth++;
	}

	if (pnode->vardef_subframe_relative_index != MD_UNUSED_INDEX) {

		MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_subframe_index < 0);
		MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_subframe_index >= max_subframe_depth);
		pnode->vardef_frame_relative_index = subframe_var_count_belows[pnode->vardef_subframe_index]
			+ pnode->vardef_subframe_relative_index;

		if (trace) {
			leader_print(subframe_depth);

			printf("NODE %s %ds%d (",
				pnode->text,
				pnode->vardef_subframe_relative_index,
				pnode->vardef_subframe_index);

			for (int i = 0; i < subframe_depth; i++) {
				if (i > 0)
					printf(",");
				printf("%d:%d", i, subframe_var_count_belows[i]);
			}
			printf(") -> %d\n",
				pnode->vardef_frame_relative_index);

		}
		MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index < 0);
		MLR_INTERNAL_CODING_ERROR_IF(pnode->vardef_frame_relative_index > *pmax_var_depth);
	}

	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_2_for_node(pchild, subframe_depth,
				var_count_below_subframe, var_count_at_subframe, pmax_var_depth,
				subframe_var_count_belows, max_subframe_depth,
				trace);
		}
	}
	if (pnode->subframe_var_count != MD_UNUSED_INDEX) {
		subframe_depth--;
		subframe_var_count_belows[subframe_depth] = MD_UNUSED_INDEX;
	}
}

// ================================================================
const char* STKALC_TRACE_LEADER = "    ";
static void leader_print(int depth) {
	for (int i = 0; i < depth; i++)
		printf("%s", STKALC_TRACE_LEADER);
}
