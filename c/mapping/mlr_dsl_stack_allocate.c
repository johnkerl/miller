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

// xxx put 'pass_1' and 'pass_2' in the function names

// * maybe move ast from containers to mapping?

// * make a separate file for tree-reorg part into top-level blocks

// * 'semantic analysis': use this to describe CST-build-time checks
// * 'object binding': use this to describe linking func/subr defs and callsites
// * separate verbosity for allocator? and invoke it in UT cases specific to this?
//   -> (note allocation marks in the AST will be printed regardless)

// * nodestash:
// @ localvar: fridx & upcount; then frgridx
// @ statement block: frct; then maxdepth (default #def NONESUCH @ ctor; respect @ ast-node printer)

// pass 1:
// @localvar put fridx & upcount
// @exit from statement block put frct

// pass 2:
// @localvar map fridx & upcount to relidx (>=0 for in-frame, <0 for upframe)
// @base put maxdepth

// 1 frame_relative_index
// 1 upstack_frame_count
// 2 frame_var_count
// 2 recursive_max_var_count
// 2 absolute_index

// ================================================================
typedef struct _stkalc_frame_t {
	long long index_count;
	lhmsi_t* pnames_to_indices;
} stkalc_frame_t;

static      stkalc_frame_t* stkalc_frame_alloc();
static void stkalc_frame_free(stkalc_frame_t* pframe);
static int  stkalc_frame_has(stkalc_frame_t* pframe, char* name);
static int  stkalc_frame_get(stkalc_frame_t* pframe, char* name);
static void stkalc_frame_add(stkalc_frame_t* pframe, char* desc, char* name, int depth, int verbose);

// ----------------------------------------------------------------
typedef struct _stkalc_frame_group_t {
	sllv_t* plist;
} stkalc_frame_group_t;

static                 stkalc_frame_group_t* stkalc_frame_group_alloc(stkalc_frame_t* pframe);
static void            stkalc_frame_group_free(stkalc_frame_group_t* pframe_group);
static void            stkalc_frame_group_push(stkalc_frame_group_t* pframe_group, stkalc_frame_t* pframe);
static stkalc_frame_t* stkalc_frame_group_pop(stkalc_frame_group_t* pframe_group);

static void stkalc_frame_group_mark_for_define(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose);

static void stkalc_frame_group_mark_for_write(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose);

static void stkalc_frame_group_mark_for_read(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose);

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_func_subr_block(mlr_dsl_ast_node_t* pnode);
static void blocked_ast_allocate_locals_for_begin_end_block(mlr_dsl_ast_node_t* pnode);
static void blocked_ast_allocate_locals_for_main_block(mlr_dsl_ast_node_t* pnode);
static void blocked_ast_allocate_locals_for_statement_block(mlr_dsl_ast_node_t* pnode,
	stkalc_frame_group_t* pframe_group);
static void blocked_ast_allocate_locals_for_node(mlr_dsl_ast_node_t* pnode,
	stkalc_frame_group_t* pframe_group);

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
// * analysis tree w/ pointers to statement-block nodes?
//   this way analysis points to AST.
//
// ----------------------------------------------------------------
// Population:
// * in-order AST traversal
// * note statement-block nodes are only every so often in the full AST
// * at each node:
//     if is local-var LHS:
//       if explicit:
//         lhmsi_put(name, ++fridx)
//       else:
//         resolve up ...
//     else if is statement-block:
//       allocate a frame struct
//       attach it to the node
//       recurse & have the recursion populate it
//       pop the frame struct but leave it attached to the node
//     else:
//       nothing to do here.
// ----------------------------------------------------------------

// ----------------------------------------------------------------
// xxx rename
void blocked_ast_allocate_locals(blocked_ast_t* paast) {
	printf("\n"); // xxx temp
	for (sllve_t* pe = paast->pfunc_defs->phead; pe != NULL; pe = pe->pnext) {
		blocked_ast_allocate_locals_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->psubr_defs->phead; pe != NULL; pe = pe->pnext) {
		blocked_ast_allocate_locals_for_func_subr_block(pe->pvvalue);
	}
	for (sllve_t* pe = paast->pbegin_blocks->phead; pe != NULL; pe = pe->pnext) {
		blocked_ast_allocate_locals_for_begin_end_block(pe->pvvalue);
	}
	blocked_ast_allocate_locals_for_main_block(paast->pmain_block);
	for (sllve_t* pe = paast->pend_blocks->phead; pe != NULL; pe = pe->pnext) {
		blocked_ast_allocate_locals_for_begin_end_block(pe->pvvalue);
	}
}

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_func_subr_block(mlr_dsl_ast_node_t* pnode) {
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
		stkalc_frame_group_mark_for_define(pframe_group, pparameter_node, "PARAMETER", TRUE/*xxx temp*/);
	}
	blocked_ast_allocate_locals_for_statement_block(plist_node, pframe_group);
	pnode->frame_var_count = pframe->index_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count);

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_begin_end_block(mlr_dsl_ast_node_t* pnode) {
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

	blocked_ast_allocate_locals_for_statement_block(pnode->pchildren->phead->pvvalue, pframe_group);
	pnode->frame_var_count = pframe->index_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count); // xxx fix name

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_main_block(mlr_dsl_ast_node_t* pnode) {
//	xxx assert node type

	printf("\n");
	printf("ALLOCATING LOCALS FOR MAIN BLOCK\n");

	// xxx make this a one-liner
	stkalc_frame_t* pframe = stkalc_frame_alloc();
	stkalc_frame_group_t* pframe_group = stkalc_frame_group_alloc(pframe);

	blocked_ast_allocate_locals_for_statement_block(pnode, pframe_group);
	pnode->frame_var_count = pframe->index_count;
	printf("BLK %s frct=%d\n", pnode->text, pnode->frame_var_count); // xxx fix name

	stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
	stkalc_frame_group_free(pframe_group);
}

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_statement_block(mlr_dsl_ast_node_t* pnode,
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
		blocked_ast_allocate_locals_for_node(pchild, pframe_group);
	}
}

// ----------------------------------------------------------------
static void blocked_ast_allocate_locals_for_node(mlr_dsl_ast_node_t* pnode,
	stkalc_frame_group_t* pframe_group)
{
	// xxx make separate functions

	if (pnode->type == MD_AST_NODE_TYPE_FOR_SREC) { // xxx comment

		for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
			printf("::  ");
		printf("PUSH FRAME %s\n", pnode->text);
		stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
		stkalc_frame_group_push(pframe_group, pnext_frame);

		mlr_dsl_ast_node_t* pvarsnode  = pnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pblocknode = pnode->pchildren->phead->pnext->pvvalue;

		mlr_dsl_ast_node_t* pknode = pvarsnode->pchildren->phead->pvvalue;
		mlr_dsl_ast_node_t* pvnode = pvarsnode->pchildren->phead->pnext->pvvalue;
		stkalc_frame_group_mark_for_define(pframe_group, pknode, "FOR-BIND", TRUE/*xxx temp*/);
		stkalc_frame_group_mark_for_define(pframe_group, pvnode, "FOR-BIND", TRUE/*xxx temp*/);

		blocked_ast_allocate_locals_for_statement_block(pblocknode, pframe_group);
		pblocknode->frame_var_count = pnext_frame->index_count;

		stkalc_frame_free(stkalc_frame_group_pop(pframe_group));

		for (int i = 0; i < pframe_group->plist->length; i++)
			printf("::  ");
		printf("POP FRAME %s frct=%d\n", pnode->text, pblocknode->frame_var_count);

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
			blocked_ast_allocate_locals_for_node(pchild, pframe_group);
		}

		for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
			printf("::  ");
		printf("PUSH FRAME %s\n", pnode->text);
		stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
		stkalc_frame_group_push(pframe_group, pnext_frame);

		for (sllve_t* pe = pkeysnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			stkalc_frame_group_mark_for_define(pframe_group, pkeynode, "FOR-BIND", TRUE/*xxx temp*/);
		}
		stkalc_frame_group_mark_for_define(pframe_group, pvalnode, "FOR-BIND", TRUE/*xxx temp*/);
		blocked_ast_allocate_locals_for_statement_block(pblocknode, pframe_group);
		// xxx make accessor ...
		pblocknode->frame_var_count = pnext_frame->index_count;

		stkalc_frame_free(stkalc_frame_group_pop(pframe_group));
		for (int i = 0; i < pframe_group->plist->length; i++)
			printf("::  ");
		printf("POP FRAME %s frct=%d\n", pnode->text, pblocknode->frame_var_count);

	} else if (pnode->type == MD_AST_NODE_TYPE_LOCAL_DEFINITION) {
		// xxx decide on preorder vs. postorder
		mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;

		stkalc_frame_group_mark_for_define(pframe_group, pnamenode, "DEFINE", TRUE/*xxx temp*/);
		mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
		blocked_ast_allocate_locals_for_node(pvaluenode, pframe_group);

	} else if (pnode->type == MD_AST_NODE_TYPE_LOCAL_ASSIGNMENT) { // xxx rename
		mlr_dsl_ast_node_t* pnamenode = pnode->pchildren->phead->pvvalue;
		stkalc_frame_group_mark_for_write(pframe_group, pnamenode, "WRITE", TRUE/*xxx temp*/);
		mlr_dsl_ast_node_t* pvaluenode = pnode->pchildren->phead->pnext->pvvalue;
		blocked_ast_allocate_locals_for_node(pvaluenode, pframe_group);

	} else if (pnode->type == MD_AST_NODE_TYPE_BOUND_VARIABLE) {
		stkalc_frame_group_mark_for_read(pframe_group, pnode, "READ", TRUE/*xxx temp*/);

	} else if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;

			if (pchild->type == MD_AST_NODE_TYPE_STATEMENT_BLOCK) {

				for (int i = 0; i < pframe_group->plist->length; i++) // xxx temp
					printf("::  ");
				printf("PUSH FRAME %s\n", pchild->text);

				stkalc_frame_t* pnext_frame = stkalc_frame_alloc();
				stkalc_frame_group_push(pframe_group, pnext_frame);

				blocked_ast_allocate_locals_for_statement_block(pchild, pframe_group);
				pchild->frame_var_count = pnext_frame->index_count;

				stkalc_frame_free(stkalc_frame_group_pop(pframe_group));

				for (int i = 0; i < pframe_group->plist->length; i++)
					printf("::  ");
				printf("POP FRAME %s frct=%d\n", pnode->text, pchild->frame_var_count);

			} else {
				blocked_ast_allocate_locals_for_node(pchild, pframe_group);
			}
		}
	}
}

// ================================================================
static stkalc_frame_t* stkalc_frame_alloc() {
	stkalc_frame_t* pframe = mlr_malloc_or_die(sizeof(stkalc_frame_t));
	pframe->index_count = 0;
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

static void stkalc_frame_add(stkalc_frame_t* pframe, char* desc, char* name, int depth, int verbose) {
	lhmsi_put(pframe->pnames_to_indices, name, pframe->index_count, NO_FREE);
	pframe->index_count++;
}

// ================================================================
static stkalc_frame_group_t* stkalc_frame_group_alloc(stkalc_frame_t* pframe) {
	stkalc_frame_group_t* pframe_group = mlr_malloc_or_die(sizeof(stkalc_frame_group_t));
	pframe_group->plist = sllv_alloc();
	sllv_prepend(pframe_group->plist, pframe);
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

static void stkalc_frame_group_mark_for_define(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose)
{
	char* op = "REUSE";
	stkalc_frame_t* pframe = pframe_group->plist->phead->pvvalue;
	if (!stkalc_frame_has(pframe, pnode->text)) {
		stkalc_frame_add(pframe, desc, pnode->text, pframe_group->plist->length, verbose);
		op = "ADD";
	}
	pnode->upstack_frame_count = 0;
	pnode->frame_relative_index = stkalc_frame_get(pframe, pnode->text);
	if (verbose) {
		for (int i = 1; i < pframe_group->plist->length; i++) {
			printf("::  ");
		}
		printf("::  %s %s %s @ %du%d\n", op, desc, pnode->text,
			pnode->frame_relative_index, pnode->upstack_frame_count);
	}
}

static void stkalc_frame_group_mark_for_write(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
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
		stkalc_frame_add(pframe, desc, pnode->text, pnode->upstack_frame_count, verbose);
		// xxx temp
		pnode->frame_relative_index = pframe->index_count;
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

static void stkalc_frame_group_mark_for_read(stkalc_frame_group_t* pframe_group, mlr_dsl_ast_node_t* pnode,
	char* desc, int verbose)
{
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

	if (found) {
		if (verbose) {
			for (int i = 1; i < pframe_group->plist->length; i++) {
				printf("::  ");
			}
			printf("::  %s %s @ %du%d\n", desc, pnode->text, pnode->frame_relative_index, pnode->upstack_frame_count);
		}
	} else {
		if (verbose) {
			for (int i = 1; i < pframe_group->plist->length; i++) {
				printf("::  ");
			}
			printf("::  %s %s ABSENT\n", desc, pnode->text);
		}
	}
}
