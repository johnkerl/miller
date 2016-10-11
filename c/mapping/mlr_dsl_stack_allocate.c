#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "containers/free_flags.h"
#include "containers/lhmsi.h"
#include "mapping/mlr_dsl_blocked_ast.h"
#include "mapping/context_flags.h"

// ================================================================
// xxx make a summary comment here
// ================================================================

// ----------------------------------------------------------------
// xxx to do:

// * maybe move ast from containers to mapping?

// * make a separate file for tree-reorg part into top-level blocks

// * 'semantic analysis': use this to describe CST-build-time checks
// * 'object binding': use this to describe linking func/subr defs and callsites
// * separate verbosity for allocator? and invoke it in UT cases specific to this?
//   -> (note allocation marks in the AST will be printed regardless)

// ================================================================
typedef struct _stkalc_frame_t {
	long long var_count;
	lhmsi_t* pnames_to_indices;
} stkalc_frame_t;

static      stkalc_frame_t* stkalc_frame_alloc();
static void stkalc_frame_free(stkalc_frame_t* pframe);
static int  stkalc_frame_has(stkalc_frame_t* pframe, char* name);
static int  stkalc_frame_get(stkalc_frame_t* pframe, char* name);
static int  stkalc_frame_add(stkalc_frame_t* pframe, char* desc, char* name, int verbose);

// ----------------------------------------------------------------
typedef struct _stkalc_frame_group_t {
	sllv_t* plist;
} stkalc_frame_group_t;

static      stkalc_frame_group_t* stkalc_frame_group_alloc(stkalc_frame_t* pframe);
static void stkalc_frame_group_free(stkalc_frame_group_t* pframe_group);
static void stkalc_frame_group_push(stkalc_frame_group_t* pframe_group, stkalc_frame_t* pframe);
static stkalc_frame_t* stkalc_frame_group_pop(stkalc_frame_group_t* pframe_group);

static void stkalc_frame_group_mark_node_for_define(stkalc_frame_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int verbose);

static void stkalc_frame_group_mark_node_for_write(stkalc_frame_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int verbose);

static void stkalc_frame_group_mark_node_for_read(stkalc_frame_group_t* pframe_group,
	mlr_dsl_ast_node_t* pnode, char* desc, int verbose);

// ----------------------------------------------------------------
static void pass_1_for_func_subr_block(mlr_dsl_ast_node_t* pnode);
static void pass_1_for_begin_end_block(mlr_dsl_ast_node_t* pnode);
static void pass_1_for_main_block(mlr_dsl_ast_node_t* pnode);
static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode, stkalc_frame_group_t* pframe_group);
static void pass_1_for_node(mlr_dsl_ast_node_t* pnode, stkalc_frame_group_t* pframe_group);

static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode);
static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int frame_depth, int var_count_below_frame, int var_count_at_frame, int* pmax_var_depth);

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
void blocked_ast_allocate_locals(blocked_ast_t* paast) {

	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_begin_end_block(pe->pvvalue);
	}
	{
		pass_1_for_main_block(paast->pmain_block);
	}
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_1_for_begin_end_block(pe->pvvalue);
	}

	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue);
	}
	{
		pass_2_for_top_level_block(paast->pmain_block);
	}
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		pass_2_for_top_level_block(pe->pvvalue);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_func_subr_block(mlr_dsl_ast_node_t* pnode) {
	// xxx make a keystroke-saver, use it here, & use it from the cst builder as well
	if (pnode->type != MD_AST_NODE_TYPE_SUBR_DEF && pnode->type != MD_AST_NODE_TYPE_FUNC_DEF) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	//	xxx assert two children of desired type

	stkalc_frame_t* pframe = stkalc_frame_alloc();
	stkalc_frame_group_t* pframe_group = stkalc_frame_group_alloc(pframe);

	printf("\n");
	printf("ALLOCATING LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	mlr_dsl_ast_node_t* pdef_name_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* plist_node = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pdef_name_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pparameter_node = pe->pvvalue;
		stkalc_frame_group_mark_node_for_define(pframe_group, pparameter_node, "PARAMETER", TRUE/*xxx temp*/);
	}
	pass_1_for_statement_block(plist_node, pframe_group);
	pnode->frame_var_count = pframe->var_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count);

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_begin_end_block(mlr_dsl_ast_node_t* pnode) {
	if (pnode->type != MD_AST_NODE_TYPE_BEGIN && pnode->type != MD_AST_NODE_TYPE_END) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	printf("\n");
	printf("ALLOCATING LOCALS FOR %s BLOCK\n", pnode->text);

	stkalc_frame_t* pframe = stkalc_frame_alloc();
	stkalc_frame_group_t* pframe_group = stkalc_frame_group_alloc(pframe);

	pass_1_for_statement_block(pnode->pchildren->phead->pvvalue, pframe_group);
	pnode->frame_var_count = pframe->var_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count); // xxx fix name

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_main_block(mlr_dsl_ast_node_t* pnode) {
	if (pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	printf("\n");
	printf("ALLOCATING LOCALS FOR MAIN BLOCK\n");

	stkalc_frame_t* pframe = stkalc_frame_alloc();
	stkalc_frame_group_t* pframe_group = stkalc_frame_group_alloc(pframe);

	pass_1_for_statement_block(pnode, pframe_group);
	pnode->frame_var_count = pframe->var_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count); // xxx fix name

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void pass_1_for_statement_block(mlr_dsl_ast_node_t* pnode,
	stkalc_frame_group_t* pframe_group)
{
	if (pnode->type != MD_AST_NODE_TYPE_STATEMENT_BLOCK) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pchild = pe->pvvalue;
		pass_1_for_node(pchild, pframe_group);
	}
}

// ----------------------------------------------------------------
static void pass_1_for_node(mlr_dsl_ast_node_t* pnode,
	stkalc_frame_group_t* pframe_group)
{
	// xxx make separate functions

	if (pnode->type == MD_AST_NODE_TYPE_LOCAL_DEFINITION) {
		// xxx decide on preorder vs. postorder
		mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;

		stkalc_frame_group_mark_node_for_define(pframe_group, pnamenode, "DEFINE", TRUE/*xxx temp*/);
		mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
		pass_1_for_node(pvaluenode, pframe_group);

	} else if (pnode->type == MD_AST_NODE_TYPE_LOCAL_ASSIGNMENT) { // xxx rename
		mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;
		stkalc_frame_group_mark_node_for_write(pframe_group, pnamenode, "WRITE", TRUE/*xxx temp*/);
		mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
		pass_1_for_node(pvaluenode, pframe_group);

	} else if (pnode->type == MD_AST_NODE_TYPE_BOUND_VARIABLE) {
		stkalc_frame_group_mark_node_for_read(pframe_group, pnode, "READ", TRUE/*xxx temp*/);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_SREC) { // xxx comment

		for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
			printf("::  ");
		printf("PUSH FRAME %s\n", pnode->text);
		stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
		stkalc_frame_group_push(pframe_group, pnext_frame);

		mlr_dsl_ast_node_t* pvarsnode  = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pblocknode = pnode->pchildren->phead->pnext->pvvalue;

		mlr_dsl_ast_node_t* pknode = pvarsnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pvnode = pvarsnode->pchildren->phead->pnext->pvvalue;
		stkalc_frame_group_mark_node_for_define(pframe_group, pknode, "FOR-BIND", TRUE/*xxx temp*/);
		stkalc_frame_group_mark_node_for_define(pframe_group, pvnode, "FOR-BIND", TRUE/*xxx temp*/);

		pass_1_for_statement_block(pblocknode, pframe_group);
		pnode->frame_var_count = pnext_frame->var_count;

		stkalc_frame_free(stkalc_frame_group_pop(pframe_group));

		for (int i = 0; i < pframe_group->plist->length; i++)
			printf("::  ");
		printf("POP FRAME %s frct=%d\n", pnode->text, pnode->frame_var_count);

	} else if (pnode->type == MD_AST_NODE_TYPE_FOR_OOSVAR) { // xxx comment

		mlr_dsl_ast_node_t* pvarsnode    = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pkeylistnode = pnode->pchildren->phead->pnext->pvvalue;
		mlr_dsl_ast_node_t* pblocknode   = pnode->pchildren->phead->pnext->pnext->pvvalue;

		mlr_dsl_ast_node_t* pkeysnode    = pvarsnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pvalnode     = pvarsnode->pchildren->phead->pnext->pvvalue;

		// xxx note keylistnode is outside the block binding. in particular if there are any localvar reads
		// in there they shouldn't read from forloop boundvars.
		for (sllve_t* pe = pkeylistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_1_for_node(pchild, pframe_group);
		}

		for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
			printf("::  ");
		printf("PUSH FRAME %s\n", pnode->text);
		stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
		stkalc_frame_group_push(pframe_group, pnext_frame);

		for (sllve_t* pe = pkeysnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			stkalc_frame_group_mark_node_for_define(pframe_group, pkeynode, "FOR-BIND", TRUE/*xxx temp*/);
		}
		stkalc_frame_group_mark_node_for_define(pframe_group, pvalnode, "FOR-BIND", TRUE/*xxx temp*/);
		pass_1_for_statement_block(pblocknode, pframe_group);
		// xxx make accessor ...
		pnode->frame_var_count = pnext_frame->var_count;

		stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
		for (int i = 0; i < pframe_group->plist->length; i++)
			printf("::  ");
		printf("POP FRAME %s frct=%d\n", pnode->text, pnode->frame_var_count);

	} else if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;

			if (pchild->type == MD_AST_NODE_TYPE_STATEMENT_BLOCK) {

				for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
					printf("::  ");
				printf("PUSH FRAME %s\n", pchild->text);

				stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
				stkalc_frame_group_push(pframe_group, pnext_frame);

				pass_1_for_statement_block(pchild, pframe_group);
				pchild->frame_var_count = pnext_frame->var_count;

				stkalc_frame_free(stkalc_frame_group_pop(pframe_group));

				for (int i = 0; i < pframe_group->plist->length; i++)
					printf("::  ");
				printf("POP FRAME %s frct=%d\n", pnode->text, pchild->frame_var_count);

			} else {
				pass_1_for_node(pchild, pframe_group);
			}
		}
	}
}

// ================================================================
static stkalc_frame_t* stkalc_frame_alloc() {
	stkalc_frame_t* pframe = mlr_malloc_or_die(sizeof(stkalc_frame_t));
	pframe->var_count = 0;
	pframe->pnames_to_indices = lhmsi_alloc();
	return pframe;
}

static void stkalc_frame_free(stkalc_frame_t* pframe) {
	if (pframe == NULL)
		return;
	lhmsi_free(pframe->pnames_to_indices);
	free(pframe);
}

static int stkalc_frame_has(stkalc_frame_t* pframe, char* name) {
	return lhmsi_has_key(pframe->pnames_to_indices, name);
}

static int stkalc_frame_get(stkalc_frame_t* pframe, char* name) {
	return lhmsi_get(pframe->pnames_to_indices, name);
}

static int stkalc_frame_add(stkalc_frame_t* pframe, char* desc, char* name, int verbose) {
	int rv = pframe->var_count;
	lhmsi_put(pframe->pnames_to_indices, name, pframe->var_count, NO_FREE);
	pframe->var_count++;
	return rv;
}

// ================================================================
static stkalc_frame_group_t* stkalc_frame_group_alloc(stkalc_frame_t* pframe) {
	stkalc_frame_group_t* pframe_group = mlr_malloc_or_die(sizeof(stkalc_frame_group_t));
	pframe_group->plist = sllv_alloc();
	sllv_prepend(pframe_group->plist, pframe);
	stkalc_frame_add(pframe, "FOR_ABSENT", "", /*xxx temp*/TRUE);
	return pframe_group;
}

static void stkalc_frame_group_free(stkalc_frame_group_t* pframe_group) {
	if (pframe_group == NULL)
		return;
	while (pframe_group->plist->phead != NULL) {
		stkalc_frame_free(sllv_pop(pframe_group->plist));
	}
	sllv_free(pframe_group->plist);
	free(pframe_group);
}

static void stkalc_frame_group_push(stkalc_frame_group_t* pframe_group, stkalc_frame_t* pframe) {
	sllv_prepend(pframe_group->plist, pframe);
}

static stkalc_frame_t* stkalc_frame_group_pop(stkalc_frame_group_t* pframe_group) {
	return sllv_pop(pframe_group->plist);
}

static void stkalc_frame_group_mark_node_for_define(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose)
{
	char* op = "REUSE";
	stkalc_frame_t* pframe = pframe_group->plist->phead->pvvalue;
	pnode->upstack_frame_count = 0;
	if (stkalc_frame_has(pframe, pnode->text)) {
		pnode->frame_relative_index = stkalc_frame_get(pframe, pnode->text);
	} else {
		// xxx this factorization is gross.
		pnode->frame_relative_index = stkalc_frame_add(pframe, desc, pnode->text, verbose);
		op = "ADD";
	}
	if (verbose) {
		for (int i = 1; i < pframe_group->plist->length; i++) {
			printf("::  ");
		}
		printf("::  %s %s %s @ %du%d\n", op, desc, pnode->text,
			pnode->frame_relative_index, pnode->upstack_frame_count);
	}
}

static void stkalc_frame_group_mark_node_for_write(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose)
{
	char* op = "REUSE";
	int found = FALSE;
	// xxx loop. if not found, fall back to top frame.
	pnode->upstack_frame_count = 0;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, pnode->upstack_frame_count++) {
		stkalc_frame_t* pframe = pe->pvvalue;
		if (stkalc_frame_has(pframe, pnode->text)) {
			found = TRUE; // xxx dup
			pnode->frame_relative_index = stkalc_frame_get(pframe, pnode->text);
			break;
		}
	}

	if (!found) {
		pnode->upstack_frame_count = 0;
		stkalc_frame_t* pframe = pframe_group->plist->phead->pvvalue;
		pnode->frame_relative_index = stkalc_frame_add(pframe, desc, pnode->text, verbose);
		// xxx temp
		op = "ADD";
	}

	if (verbose) {
		for (int i = 1; i < pframe_group->plist->length; i++) {
			printf("::  ");
		}
		printf("::  %s %s %s @ %du%d\n", op, desc, pnode->text,
			pnode->frame_relative_index, pnode->upstack_frame_count);
	}
}

// xxx make this very clear in the header somehow ... this is an assumption to be tracked across modules.
static void stkalc_frame_group_mark_node_for_read(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose)
{
	char* op = "PRESENT";
	int found = FALSE;
	// xxx loop. if not found, fall back to top frame.
	int upstack_frame_count = 0;
	for (sllve_t* pe = pframe_group->plist->phead; pe != NULL; pe = pe->pnext, upstack_frame_count++) {
		stkalc_frame_t* pframe = pe->pvvalue;
		if (stkalc_frame_has(pframe, pnode->text)) {
			found = TRUE; // xxx dup
			pnode->frame_relative_index = stkalc_frame_get(pframe, pnode->text);
			break;
		}
	}

	// xxx if not found: go to the tail & use the "" entry
	if (!found) {
		stkalc_frame_t* plast = pframe_group->plist->ptail->pvvalue;
		pnode->frame_relative_index = stkalc_frame_get(plast, "");
		upstack_frame_count = pframe_group->plist->length - 1;
		op = "ABSENT";
	}

	if (verbose) {
		for (int i = 1; i < pframe_group->plist->length; i++) {
			printf("::  ");
		}
		printf("::  %s %s %s @ %du%d\n", desc, pnode->text, op, pnode->frame_relative_index, upstack_frame_count);
	}
	pnode->upstack_frame_count = upstack_frame_count;

}

// ================================================================
static void pass_2_for_top_level_block(mlr_dsl_ast_node_t* pnode) {
	int frame_depth = 0;
	int var_count_below_frame = 0;
	int var_count_at_frame = 0;
	int max_var_depth   = 0;
	printf("\n");
	printf("ABSOLUTIZING LOCALS FOR DEFINITION BLOCK [%s]\n", pnode->text);
	pass_2_for_node(pnode, frame_depth, var_count_below_frame, var_count_at_frame, &max_var_depth);
	pnode->max_var_depth = max_var_depth;
}

static void pass_2_for_node(mlr_dsl_ast_node_t* pnode,
	int frame_depth, int var_count_below_frame, int var_count_at_frame, int* pmax_var_depth)
{
	if (pnode->frame_var_count != MD_UNUSED_INDEX) {
		var_count_below_frame += var_count_at_frame;
		var_count_at_frame = pnode->frame_var_count;
		int depth = var_count_below_frame + var_count_at_frame;
		frame_depth++;
		if (depth > *pmax_var_depth) {
			*pmax_var_depth = depth;
		}
		// xxx funcify
		for (int i = 1; i < frame_depth; i++)
			printf("::  ");
		printf("FRAME [%s] var_count_below=%d var_count_at=%d max_var_depth=%d\n",
			pnode->text, var_count_below_frame, var_count_at_frame, *pmax_var_depth);
	}

	if (pnode->frame_relative_index != MD_UNUSED_INDEX) {
		pnode->absolute_index = var_count_below_frame + pnode->frame_relative_index;
		for (int i = 0; i < frame_depth; i++)
			printf("::  ");
		printf("NODE %s %du%d -> %d\n",
			pnode->text,
			pnode->frame_relative_index,
			pnode->upstack_frame_count,
			pnode->absolute_index);
	}

	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			pass_2_for_node(pchild, frame_depth,
				var_count_below_frame, var_count_at_frame, pmax_var_depth);
		}
	}
}
