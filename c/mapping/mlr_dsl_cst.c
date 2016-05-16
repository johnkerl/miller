#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "mlr_dsl_cst.h"

static mlr_dsl_ast_node_t* get_list_for_block(mlr_dsl_ast_node_t* pnode);

static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* past, int type_inferencing, int begin_end_only);
static mlr_dsl_cst_statement_t* cst_statement_alloc_blank();
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static mlr_dsl_cst_statement_t* cst_statement_alloc_srec_assignment(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_oosvar_assignment(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_oosvar_from_full_srec_assignment(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_full_srec_from_oosvar_assignment(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_unset(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_emitf(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_emit_or_emitp(mlr_dsl_ast_node_t* past, int type_inferencing,
	int do_full_prefixing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_conditional_block(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_while(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_do_while(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_for_srec(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_if_head(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_if_item(mlr_dsl_ast_node_t* pexprnode,
	mlr_dsl_ast_node_t* plistnode, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_filter(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_dump(mlr_dsl_ast_node_t* past, int type_inferencing);
static mlr_dsl_cst_statement_t* cst_statement_alloc_bare_boolean(mlr_dsl_ast_node_t* past, int type_inferencing);

static mlr_dsl_cst_statement_vararg_t* mlr_dsl_cst_statement_vararg_alloc(
	char*             emitf_or_unset_srec_field_name,
	rval_evaluator_t* pemitf_arg_evaluator,
	sllv_t*           punset_oosvar_keylist_evaluators);

static void cst_statement_vararg_free(mlr_dsl_cst_statement_vararg_t* pvararg);

static void mlr_dsl_cst_node_handle_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_oosvar_to_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_full_srec_from_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_unset(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_unset_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_emitf(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_emitp(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_emitp_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_emit(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_emit_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_dump(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_filter(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_conditional_block(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_while(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_do_while(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_do_while(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_for_srec(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_if_head(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static void mlr_dsl_cst_node_handle_bare_boolean(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack);

static sllv_t* allocate_keylist_evaluators_from_oosvar_node(mlr_dsl_ast_node_t* pnode, int type_inferencing);

// ----------------------------------------------------------------
// For mlr filter, which takes a reduced subset of mlr-put syntax:
// * The root node of the AST must be a statement list (as for put).
// * The list must have one child node.
// * That child node must not be a braced statement (begin, end, for, cond, etc.)
// * The child node must evaluate to boolean, although this is fully enforced only
//   during stream processing.

mlr_dsl_ast_node_t* extract_filterable_statement(mlr_dsl_ast_t* past, int type_inferencing) {
	mlr_dsl_ast_node_t* proot = past->proot;

	if (proot == NULL) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d: null root node.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	if (proot->pchildren->phead == NULL) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d: null left child node.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	if (proot->pchildren->phead->pnext != NULL) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d: extraneous right child node.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	if (proot->type != MD_AST_NODE_TYPE_STATEMENT_LIST) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_LIST),
			mlr_dsl_ast_node_describe_type(proot->type));
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft = proot->pchildren->phead->pvvalue;

	return pleft;
}

// ----------------------------------------------------------------
// Example:

// $ mlr -n put -v '#begin{@a=1;@b=2};$m=2;$n=4;end{@y=5;@z=6}'
// AST ROOT:
// list (statement_list):
//     begin (begin):
//         list (statement_list):
//             = (oosvar_assignment):
//                 a (oosvar_name).
//                 1 (strnum_literal).
//             = (oosvar_assignment):
//                 b (oosvar_name).
//                 2 (strnum_literal).
//     = (srec_assignment):
//         m (field_name).
//         2 (strnum_literal).
//     = (srec_assignment):
//         n (field_name).
//         4 (strnum_literal).
//     end (end):
//         list (statement_list):
//             = (oosvar_assignment):
//                 y (oosvar_name).
//                 5 (strnum_literal).
//             = (oosvar_assignment):
//                 z (oosvar_name).
//                 6 (strnum_literal).

mlr_dsl_cst_t* mlr_dsl_cst_alloc(mlr_dsl_ast_t* past, int type_inferencing) {
	// Root node is not populated on empty-string input to the parser.
	if (past->proot == NULL) {
		past->proot = mlr_dsl_ast_node_alloc_zary("list", MD_AST_NODE_TYPE_STATEMENT_LIST);
	}

	mlr_dsl_cst_t* pcst = mlr_malloc_or_die(sizeof(mlr_dsl_cst_t));

	if (past->proot->type != MD_AST_NODE_TYPE_STATEMENT_LIST) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected root node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_LIST),
			mlr_dsl_ast_node_describe_type(past->proot->type));
		exit(1);
	}

	pcst->pbegin_statements = sllv_alloc();
	pcst->pmain_statements  = sllv_alloc();
	pcst->pend_statements   = sllv_alloc();
	mlr_dsl_ast_node_t* plistnode = NULL;
	////xxx delete all these printfs
	////printf("AST->CST STUB\n");
	for (sllve_t* pe = past->proot->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		switch (pnode->type) {
		case MD_AST_NODE_TYPE_BEGIN:
			////printf("GOT BEGIN:%s\n", mlr_dsl_ast_node_describe_type(pnode->type));
			plistnode = get_list_for_block(pnode);
			for (sllve_t* pe = plistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pchild = pe->pvvalue;
				////printf("-- WITHIN BEGIN:%s\n", mlr_dsl_ast_node_describe_type(pchild->type));
				sllv_append(pcst->pbegin_statements, cst_statement_alloc(pchild, type_inferencing, TRUE));
			}
			break;
		case MD_AST_NODE_TYPE_END:
			////printf("GOT END:%s\n", mlr_dsl_ast_node_describe_type(pnode->type));
			plistnode = get_list_for_block(pnode);
			for (sllve_t* pe = plistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pchild = pe->pvvalue;
				////printf("-- WITHIN END:%s\n", mlr_dsl_ast_node_describe_type(pchild->type));
				sllv_append(pcst->pend_statements, cst_statement_alloc(pchild, type_inferencing, TRUE));
			}
			break;
		default:
			////printf("GOT MAIN:%s\n", mlr_dsl_ast_node_describe_type(pnode->type));
			sllv_append(pcst->pmain_statements, cst_statement_alloc(pnode, type_inferencing, FALSE));
			break;
		}
	}
	return pcst;
}

// ----------------------------------------------------------------
// For begin, end, cond: there must be one child node, of type list.
static mlr_dsl_ast_node_t* get_list_for_block(mlr_dsl_ast_node_t* pnode) {
	if (pnode->pchildren->phead == NULL) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d: null left child node.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	if (pnode->pchildren->phead->pnext != NULL) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d: extraneous right child node.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	mlr_dsl_ast_node_t* pleft = pnode->pchildren->phead->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_STATEMENT_LIST) {
		fprintf(stderr,
			"%s: internal coding error detected in file %s at line %d:\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		fprintf(stderr,
			"expected node type %s but found %s.\n",
			mlr_dsl_ast_node_describe_type(MD_AST_NODE_TYPE_STATEMENT_LIST),
			mlr_dsl_ast_node_describe_type(pleft->type));
		exit(1);
	}
	return pleft;
}

// ----------------------------------------------------------------
void mlr_dsl_cst_free(mlr_dsl_cst_t* pcst) {
	if (pcst == NULL)
		return;
	for (sllve_t* pe = pcst->pbegin_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pmain_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	for (sllve_t* pe = pcst->pend_statements->phead; pe != NULL; pe = pe->pnext)
		cst_statement_free(pe->pvvalue);
	sllv_free(pcst->pbegin_statements);
	sllv_free(pcst->pmain_statements);
	sllv_free(pcst->pend_statements);

	free(pcst);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc(mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int begin_end_only)
{
    switch(pnode->type) {

	case MD_AST_NODE_TYPE_BEGIN:
		fprintf(stderr, "%s: begin statements are only valid at top level.\n", MLR_GLOBALS.argv0);
		exit(1);
		break;
	case MD_AST_NODE_TYPE_END:
		fprintf(stderr, "%s: begin statements are only valid at top level.\n", MLR_GLOBALS.argv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK: // xxx rename to ..._COND
		return cst_statement_alloc_conditional_block(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_WHILE:
		return cst_statement_alloc_while(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_DO_WHILE:
		return cst_statement_alloc_do_while(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_FOR_SREC:
		return cst_statement_alloc_for_srec(pnode, type_inferencing);
		break;

	// xxx for-oosvar

	case MD_AST_NODE_TYPE_IF_HEAD:
		return cst_statement_alloc_if_head(pnode, type_inferencing);
		break;

	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:
		return cst_statement_alloc_srec_assignment(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:
		return cst_statement_alloc_oosvar_assignment(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
		return cst_statement_alloc_oosvar_from_full_srec_assignment(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT:
		return cst_statement_alloc_full_srec_from_oosvar_assignment(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_UNSET:
		return cst_statement_alloc_unset(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_EMITF:
		return cst_statement_alloc_emitf(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_EMITP:
		return cst_statement_alloc_emit_or_emitp(pnode, type_inferencing, TRUE);
		break;
	case MD_AST_NODE_TYPE_EMIT:
		return cst_statement_alloc_emit_or_emitp(pnode, type_inferencing, FALSE);
		break;
	case MD_AST_NODE_TYPE_FILTER:
		return cst_statement_alloc_filter(pnode, type_inferencing);
		break;
	case MD_AST_NODE_TYPE_DUMP:
		return cst_statement_alloc_dump(pnode, type_inferencing);
		break;
	default:
		return cst_statement_alloc_bare_boolean(pnode, type_inferencing);
		break;
	}
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc_blank() {
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));

	pstatement->phandler                         = NULL;
	pstatement->poosvar_lhs_keylist_evaluators   = NULL;
	pstatement->srec_lhs_field_name              = NULL;
	pstatement->prhs_evaluator                   = NULL;
	pstatement->poosvar_rhs_keylist_evaluators   = NULL;
	pstatement->pemit_oosvar_namelist_evaluators = NULL;
	pstatement->pvarargs                         = NULL;
	pstatement->pblock_statements                = NULL;
	pstatement->pif_chain_statements             = NULL;
	pstatement->pbound_variables                      = NULL;

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc_srec_assignment(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	if ((past->pchildren == NULL) || (past->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_FIELD_NAME) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	} else if (pleft->pchildren != NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_srec_assignment;
	pstatement->srec_lhs_field_name = pleft->text;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, type_inferencing);
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_oosvar_assignment(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pleft->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	sllv_t* poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pleft, type_inferencing);

	int is_oosvar_to_oosvar = FALSE;

	if (pleft->type == MD_AST_NODE_TYPE_OOSVAR_NAME || pleft->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		if (pright->type == MD_AST_NODE_TYPE_OOSVAR_NAME || pright->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			is_oosvar_to_oosvar = TRUE;
		}
	}

	if (is_oosvar_to_oosvar) {
		pstatement->phandler = mlr_dsl_cst_node_handle_oosvar_to_oosvar_assignment;
		pstatement->poosvar_rhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pright,
			type_inferencing);
	} else {
		pstatement->phandler = mlr_dsl_cst_node_handle_oosvar_assignment;
		pstatement->poosvar_rhs_keylist_evaluators = NULL;
	}

	pstatement->poosvar_lhs_keylist_evaluators = poosvar_lhs_keylist_evaluators;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, type_inferencing);

	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_oosvar_from_full_srec_assignment(
	mlr_dsl_ast_node_t* past, int type_inferencing)
{
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pleft->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	if (pright->type != MD_AST_NODE_TYPE_FULL_SREC) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_oosvar_from_full_srec_assignment;
	pstatement->poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pleft,
		type_inferencing);
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_full_srec_from_oosvar_assignment(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_FULL_SREC) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	if (pright->type != MD_AST_NODE_TYPE_OOSVAR_NAME && pright->type != MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_full_srec_from_oosvar_assignment;
	pstatement->poosvar_rhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pright,
		type_inferencing);
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_unset(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	pstatement->phandler = mlr_dsl_cst_node_handle_unset;
	pstatement->pvarargs = sllv_alloc();
	for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;

		if (pnode->type == MD_AST_NODE_TYPE_ALL) {
			// The grammar allows only 'unset all', not 'unset @x, all, $y'.
			// So if 'all' appears at all, it's the only name.
			pstatement->phandler = mlr_dsl_cst_node_handle_unset_all;

		} else if (pnode->type == MD_AST_NODE_TYPE_FIELD_NAME) {
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				pnode->text,
				NULL,
				NULL));

		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME || pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				NULL,
				NULL,
				allocate_keylist_evaluators_from_oosvar_node(pnode, type_inferencing)));

		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.argv0, __FILE__, __LINE__);
			exit(1);
		}
	}
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_emitf(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	// Loop over oosvar names to emit in e.g. 'emitf @a, @b, @c'.
	pstatement->pvarargs = sllv_alloc();
	for (sllve_t* pe = past->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;
		sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
			pnode->text,
			rval_evaluator_alloc_from_ast(pnode, type_inferencing),
			NULL));
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_emitf;
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_emit_or_emitp(mlr_dsl_ast_node_t* past, int type_inferencing,
	int do_full_prefixing)
{
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;

	// The grammar allows only 'emit all', not 'emit @x, all, $y'.
	// So if 'all' appears at all, it's the only name.
	if (pnode->type == MD_AST_NODE_TYPE_ALL) {

		sllv_t* pemit_oosvar_namelist_evaluators = sllv_alloc();
		for (sllve_t* pe = past->pchildren->phead->pnext; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			sllv_append(pemit_oosvar_namelist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
		}

		pstatement->phandler = do_full_prefixing
			? mlr_dsl_cst_node_handle_emit_all
			: mlr_dsl_cst_node_handle_emitp_all;
		pstatement->pemit_oosvar_namelist_evaluators = pemit_oosvar_namelist_evaluators;

	} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME || pnode->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
		// First argument is oosvar name (e.g. @sums) or keyed ooosvar name (e.g. @sums[$group]). Remainings
		// evaluate to string, e.g. 'emit @sums, "color", "shape"'.
		mlr_dsl_ast_node_t* pnamenode = past->pchildren->phead->pvvalue;

		sllv_t* pemit_oosvar_namelist_evaluators = sllv_alloc();
		for (sllve_t* pe = past->pchildren->phead->pnext; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			sllv_append(pemit_oosvar_namelist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
		}

		pstatement->pvarargs = sllv_alloc();
		sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
			pnamenode->text,
			NULL,
			NULL));

		pstatement->phandler = do_full_prefixing
			? mlr_dsl_cst_node_handle_emit
			: mlr_dsl_cst_node_handle_emitp;

		pstatement->poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pnamenode,
			type_inferencing);
		pstatement->pemit_oosvar_namelist_evaluators = pemit_oosvar_namelist_evaluators;

	} else {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.argv0, __FILE__, __LINE__);
		exit(1);
	}
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_conditional_block(mlr_dsl_ast_node_t* pnode,
	int type_inferencing)
{
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	// Left node is the AST for the boolean expression.
	// Right node is a list of statements to be executed if the left evaluates to true.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		// xxx stub last arg
		mlr_dsl_cst_statement_t *pstatement = cst_statement_alloc(pbody_ast_node, type_inferencing, FALSE);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_conditional_block;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pleft, type_inferencing);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_while(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	// Left child node is the AST for the boolean expression.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = cst_statement_alloc(pbody_ast_node, type_inferencing, FALSE); // xxx stub
		////printf("BODY %s\n", mlr_dsl_ast_node_describe_type(pbody_ast_node->type));
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_while;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pleft, type_inferencing);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_do_while(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	// Left child node is the list of statements in the body.
	// Right child node is the AST for the boolean expression.
	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = pleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = cst_statement_alloc(pbody_ast_node, type_inferencing, FALSE); // xxx stub
		////printf("BODY %s\n", mlr_dsl_ast_node_describe_type(pbody_ast_node->type));
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_do_while;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, type_inferencing);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'for (k,v in $*) { $x=1; $y=2 }'
// list (statement_list):
//     for (for-srec):
//         variables (for-variables):
//             k (non_sigil_name).
//             v (non_sigil_name).
//         list (statement_list):
//             = (srec_assignment):
//                 x (field_name).
//                 1 (strnum_literal).
//             = (srec_assignment):
//                 y (field_name).
//                 2 (strnum_literal).

static mlr_dsl_cst_statement_t* cst_statement_alloc_for_srec(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	// Left child node is list of bound variables.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = past->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = past->pchildren->phead->pnext->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		// xxx also elsewhere, invalidate. cmt there this is done at the CST
		// rather than AST-parse level since we can give better error messages
		// (and, a simpler Lemon grammar).
		if (pbody_ast_node->type == MD_AST_NODE_TYPE_CONTINUE) {
			printf("continue alloc stub!\n");
		} else if (pbody_ast_node->type == MD_AST_NODE_TYPE_BREAK) {
			printf("break alloc stub!\n");
		} else {
			// xxx stub 3rd arg
			sllv_append(pblock_statements, cst_statement_alloc(pbody_ast_node, type_inferencing, FALSE));
		}
	}

	mlr_dsl_ast_node_t* pknode = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvnode = pleft->pchildren->phead->pnext->pvvalue;

	pstatement->phandler = mlr_dsl_cst_node_handle_for_srec;
	pstatement->pblock_statements = pblock_statements;
	pstatement->for_srec_k_name   = pknode->text;
	pstatement->for_srec_v_name   = pvnode->text;
	pstatement->pbound_variables  = lhmsmv_alloc();

	return pstatement;
}

// Example parser-input:
//
//   if (NR == 9) {
//       $x = 10;
//       $x = 11
//   } elif (NR == 12) {
//       $x = 13;
//       $x = 14
//   } else {
//       $x = 15;
//       $x = 16
//   };
//
// Corresponding parser-output AST:
//   if_head (if_head):
//       if (if_item):
//           == (operator):
//               NR (context_variable).
//               9 (strnum_literal).
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   10 (strnum_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   11 (strnum_literal).
//       elif (if_item):
//           == (operator):
//               NR (context_variable).
//               12 (strnum_literal).
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   13 (strnum_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   14 (strnum_literal).
//       else (if_item):
//           list (statement_list):
//               = (srec_assignment):
//                   x (field_name).
//                   15 (strnum_literal).
//               = (srec_assignment):
//                   x (field_name).
//                   16 (strnum_literal).

// xxx rename pasts to pnodes thruout

static mlr_dsl_cst_statement_t* cst_statement_alloc_if_head(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	sllv_t* pif_chain_statements = sllv_alloc();
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		// For if and elif:
		// * Left subnode is the AST for the boolean expression.
		// * Right subnode is a list of statements to be executed if the left evaluates to true.
		// For else:
		// * Sole subnode is a list of statements to be executed.
		mlr_dsl_ast_node_t* pitemnode = pe->pvvalue;
		mlr_dsl_ast_node_t* pexprnode = NULL;
		mlr_dsl_ast_node_t* plistnode = NULL;
		if (pitemnode->pchildren->length == 2) {
			pexprnode = pitemnode->pchildren->phead->pvvalue;
			plistnode = pitemnode->pchildren->phead->pnext->pvvalue;
		} else {
			pexprnode = NULL;
			plistnode = pitemnode->pchildren->phead->pvvalue;
		}

		sllv_append(pif_chain_statements, cst_statement_alloc_if_item(pexprnode, plistnode,
			type_inferencing/*, xxx FALSE*/));
	}

	pstatement->phandler = mlr_dsl_cst_node_handle_if_head;
	pstatement->pif_chain_statements = pif_chain_statements;
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_if_item(mlr_dsl_ast_node_t* pexprnode,
	mlr_dsl_ast_node_t* plistnode, int type_inferencing)
{
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = plistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		// xxx stub last arg
		mlr_dsl_cst_statement_t *pstatement = cst_statement_alloc(pbody_ast_node, type_inferencing, FALSE);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->phandler = NULL; // handled by the containing if-head evaluator
	pstatement->prhs_evaluator = pexprnode != NULL
		? rval_evaluator_alloc_from_ast(pexprnode, type_inferencing) // if-statement or elif-statement
		: rval_evaluator_alloc_from_boolean(TRUE); // else-statement
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* cst_statement_alloc_filter(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	mlr_dsl_ast_node_t* pnode = past->pchildren->phead->pvvalue;

	pstatement->phandler = mlr_dsl_cst_node_handle_filter;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pnode, type_inferencing);
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_dump(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	pstatement->phandler = mlr_dsl_cst_node_handle_dump;
	return pstatement;
}

static mlr_dsl_cst_statement_t* cst_statement_alloc_bare_boolean(mlr_dsl_ast_node_t* past, int type_inferencing) {
	mlr_dsl_cst_statement_t* pstatement = cst_statement_alloc_blank();

	pstatement->phandler = mlr_dsl_cst_node_handle_bare_boolean;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(past, type_inferencing);
	return pstatement;
}

// ----------------------------------------------------------------
static void cst_statement_free(mlr_dsl_cst_statement_t* pstatement) {

	if (pstatement->poosvar_lhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstatement->poosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstatement->poosvar_lhs_keylist_evaluators);
	}

	if (pstatement->prhs_evaluator != NULL) {
		pstatement->prhs_evaluator->pfree_func(pstatement->prhs_evaluator);
	}

	if (pstatement->poosvar_rhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstatement->poosvar_rhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstatement->poosvar_rhs_keylist_evaluators);
	}

	if (pstatement->pemit_oosvar_namelist_evaluators != NULL) {
		for (sllve_t* pe = pstatement->pemit_oosvar_namelist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstatement->pemit_oosvar_namelist_evaluators);
	}

	if (pstatement->pvarargs != NULL) {
		for (sllve_t* pe = pstatement->pvarargs->phead; pe != NULL; pe = pe->pnext)
			cst_statement_vararg_free(pe->pvvalue);
		sllv_free(pstatement->pvarargs);
	}

	if (pstatement->pblock_statements != NULL) {
		for (sllve_t* pe = pstatement->pblock_statements->phead; pe != NULL; pe = pe->pnext)
			cst_statement_free(pe->pvvalue);
		sllv_free(pstatement->pblock_statements);
	}

	free(pstatement);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_vararg_t* mlr_dsl_cst_statement_vararg_alloc(
	char*             emitf_or_unset_srec_field_name,
	rval_evaluator_t* pemitf_arg_evaluator,
	sllv_t*           punset_oosvar_keylist_evaluators)
{
	mlr_dsl_cst_statement_vararg_t* pvararg = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_vararg_t));
	pvararg->emitf_or_unset_srec_field_name = emitf_or_unset_srec_field_name == NULL
		? NULL : mlr_strdup_or_die(emitf_or_unset_srec_field_name);
	pvararg->punset_oosvar_keylist_evaluators  = punset_oosvar_keylist_evaluators;
	pvararg->pemitf_arg_evaluator            = pemitf_arg_evaluator;
	return pvararg;
}

static void cst_statement_vararg_free(mlr_dsl_cst_statement_vararg_t* pvararg) {
	if (pvararg == NULL)
		return;
	free(pvararg->emitf_or_unset_srec_field_name);

	if (pvararg->punset_oosvar_keylist_evaluators != NULL) {
		for (sllve_t* pe = pvararg->punset_oosvar_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pvararg->punset_oosvar_keylist_evaluators);
	}

	if (pvararg->pemitf_arg_evaluator != NULL)
		pvararg->pemitf_arg_evaluator->pfree_func(pvararg->pemitf_arg_evaluator);

	free(pvararg);
}

// ----------------------------------------------------------------
void mlr_dsl_cst_handle(
	sllv_t*          pcst_statements,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		pstatement->phandler(pstatement, poosvars, pinrec, ptyped_overlay, ppregex_captures,
			pctx, pshould_emit_rec, poutrecs, oosvar_flatten_separator, pbind_stack);
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	char* srec_lhs_field_name = pnode->srec_lhs_field_name;
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	mv_t* pval = mlr_malloc_or_die(sizeof(mv_t));
	*pval = val;

	// Write typed mlrval output to the typed overlay rather than into the lrec (which holds only
	// string values).
	//
	// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the input had
	// "x"=>"abc","y"=>"def" but the previous pass through this loop set "y"=>7.4 and "z"=>"ghi" then an
	// expression right-hand side referring to $y would get the floating-point value 7.4. So we don't need
	// to do lrec_put here, and moreover should not for two reasons: (1) there is a performance hit of doing
	// throwaway number-to-string formatting -- it's better to do it once at the end; (2) having the string
	// values doubly owned by the typed overlay and the lrec would result in double frees, or awkward
	// bookkeeping. However, the NR variable evaluator reads prec->field_count, so we need to put something
	// here. And putting something statically allocated minimizes copying/freeing.
	if (mv_is_present(pval)) {
		// xxx to do: replace the typed overlay with an mlhmmv entirely.
		mv_t* pold = lhmsv_get(ptyped_overlay, srec_lhs_field_name);
		if (pold != NULL) {
			mv_free(pold);
			free(pold);
		}
		lhmsv_put(ptyped_overlay, srec_lhs_field_name, pval, FREE_ENTRY_VALUE);
		lrec_put(pinrec, srec_lhs_field_name, "bug", NO_FREE);
	} else {
		mv_free(pval);
		free(pval);
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay,
		poosvars, ppregex_captures, pctx, prhs_evaluator->pvstate);

	if (mv_is_present(&rhs_value)) {
		int all_non_null_or_error = TRUE;
		sllmv_t* pmvkeys = evaluate_list(pnode->poosvar_lhs_keylist_evaluators,
			pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
		if (all_non_null_or_error)
			mlhmmv_put_terminal(poosvars, pmvkeys, &rhs_value);
		sllmv_free(pmvkeys);
	}
	mv_free(&rhs_value);
}

// ----------------------------------------------------------------
// All assignments produce a mlrval on the RHS and store it on the left -- except if both LHS and RHS
// are oosvars in which case there are recursive copies, or in case of $* on the LHS or RHS.

static void mlr_dsl_cst_node_handle_oosvar_to_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int lhs_all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pnode->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &lhs_all_non_null_or_error);

	if (lhs_all_non_null_or_error) {
		int rhs_all_non_null_or_error = TRUE;
		sllmv_t* prhskeys = evaluate_list(pnode->poosvar_rhs_keylist_evaluators,
			pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &rhs_all_non_null_or_error);
		if (rhs_all_non_null_or_error) {
			mlhmmv_copy(poosvars, plhskeys, prhskeys);
		}
		sllmv_free(prhskeys);
	}

	sllmv_free(plhskeys);
}

static void mlr_dsl_cst_node_handle_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pnode->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {

		mlhmmv_level_t* plevel = mlhmmv_get_or_create_level(poosvars, plhskeys);
		if (plevel != NULL) {

			mlhmmv_clear_level(plevel);

			for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
				mv_t k = mv_from_string(pe->key, NO_FREE); // mlhmmv_put_terminal_from_level will copy
				sllmve_t e = { .value = k, .free_flags = 0, .pnext = NULL };
				mv_t* pomv = lhmsv_get(ptyped_overlay, pe->key);
				if (pomv != NULL) {
					mlhmmv_put_terminal_from_level(plevel, &e, pomv);
				} else {
					mv_t v = mv_from_string(pe->value, NO_FREE); // mlhmmv_put_terminal_from_level will copy
					mlhmmv_put_terminal_from_level(plevel, &e, &v);
				}
			}

		}
	}
	sllmv_free(plhskeys);
}

static void mlr_dsl_cst_node_handle_full_srec_from_oosvar_assignment(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	lrec_clear(pinrec);
	for (lhmsve_t* pe = ptyped_overlay->phead; pe != NULL; pe = pe->pnext) {
		mv_t* pmv = pe->pvvalue;
		mv_free(pmv);
	}
	lhmsv_clear(ptyped_overlay);

	int all_non_null_or_error = TRUE;
	sllmv_t* prhskeys = evaluate_list(pnode->poosvar_rhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {
		int error = 0;
		mlhmmv_level_t* plevel = mlhmmv_get_level(poosvars, prhskeys, &error);
		if (plevel != NULL) {
			for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
				if (pentry->level_value.is_terminal) {
					char* skey = mv_alloc_format_val(&pentry->level_key);
					mv_t* pval = mv_alloc_copy(&pentry->level_value.u.mlrval);

					// Write typed mlrval output to the typed overlay rather than into the lrec
					// (which holds only string values).
					//
					// The rval_evaluator reads the overlay in preference to the lrec. E.g. if the
					// input had "x"=>"abc","y"=>"def" but a previous statement had set "y"=>7.4 and
					// "z"=>"ghi", then an expression right-hand side referring to $y would get the
					// floating-point value 7.4. So we don't need to lrec_put the value here, and
					// moreover should not for two reasons: (1) there is a performance hit of doing
					// throwaway number-to-string formatting -- it's better to do it once at the
					// end; (2) having the string values doubly owned by the typed overlay and the
					// lrec would result in double frees, or awkward bookkeeping. However, the NR
					// variable evaluator reads prec->field_count, so we need to put something here.
					// And putting something statically allocated minimizes copying/freeing.
					mv_t* pold = lhmsv_get(ptyped_overlay, skey);
					if (pold != NULL) {
						// xxx to do: replace the typed overlay with an mlhmmv entirely.
						mv_free(pold);
						free(pold);
					}
					lhmsv_put(ptyped_overlay, mlr_strdup_or_die(skey), pval, FREE_ENTRY_KEY);
					lrec_put(pinrec, skey, "bug", FREE_ENTRY_KEY);
				}
			}
		}
	}
	sllmv_free(prhskeys);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_unset(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	for (sllve_t* pf = pnode->pvarargs->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_vararg_t* pvararg = pf->pvvalue;
		if (pvararg->punset_oosvar_keylist_evaluators != NULL) {

			int all_non_null_or_error = TRUE;
			sllmv_t* pmvkeys = evaluate_list(pvararg->punset_oosvar_keylist_evaluators,
				pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);

			if (all_non_null_or_error)
				mlhmmv_remove(poosvars, pmvkeys);
			sllmv_free(pmvkeys);
		} else {
			lrec_remove(pinrec, pvararg->emitf_or_unset_srec_field_name);
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_unset_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	sllmv_t* pempty = sllmv_alloc();
	mlhmmv_remove(poosvars, pempty);
	sllmv_free(pempty);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_emitf(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	lrec_t* prec_to_emit = lrec_unbacked_alloc();
	for (sllve_t* pf = pnode->pvarargs->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_vararg_t* pvararg = pf->pvvalue;
		char* emitf_or_unset_srec_field_name = pvararg->emitf_or_unset_srec_field_name;
		rval_evaluator_t* pemitf_arg_evaluator = pvararg->pemitf_arg_evaluator;

		// This is overkill ... the grammar allows only for oosvar names as args to emit.  So we could bypass
		// that and just hashmap-get keyed by emitf_or_unset_srec_field_name here.
		mv_t val = pemitf_arg_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, pemitf_arg_evaluator->pvstate);

		if (val.type == MT_STRING) {
			// Ownership transfer from (newly created) mlrval to (newly created) lrec.
			lrec_put(prec_to_emit, emitf_or_unset_srec_field_name, val.u.strv, val.free_flags);
		} else {
			char free_flags = NO_FREE;
			char* string = mv_format_val(&val, &free_flags);
			lrec_put(prec_to_emit, emitf_or_unset_srec_field_name, string, free_flags);
		}

	}
	sllv_append(poutrecs, prec_to_emit);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_emitp(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pnode->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pnode->pemit_oosvar_namelist_evaluators,
			pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mlhmmv_to_lrecs(poosvars, pmvkeys, pmvnames, poutrecs, FALSE, oosvar_flatten_separator);
		}
		sllmv_free(pmvnames);
	}
	sllmv_free(pmvkeys);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_emitp_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pnode->pemit_oosvar_namelist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_all_to_lrecs(poosvars, pmvnames, poutrecs, FALSE, oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_emit(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pnode->poosvar_lhs_keylist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pnode->pemit_oosvar_namelist_evaluators,
			pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mlhmmv_to_lrecs(poosvars, pmvkeys, pmvnames, poutrecs, TRUE, oosvar_flatten_separator);
		}
		sllmv_free(pmvnames);
	}
	sllmv_free(pmvkeys);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_emit_all(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pnode->pemit_oosvar_namelist_evaluators,
		pinrec, ptyped_overlay, poosvars, ppregex_captures, pctx, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_all_to_lrecs(poosvars, pmvnames, poutrecs, TRUE, oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_dump(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	mlhmmv_print_json_stacked(poosvars, FALSE);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_filter(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		*pshould_emit_rec = val.u.boolv;
	} else {
		*pshould_emit_rec = FALSE;
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_conditional_block(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		if (val.u.boolv) {
			mlr_dsl_cst_handle(pnode->pblock_statements,
				poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs,
				oosvar_flatten_separator, pbind_stack);
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_while(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	while (TRUE) {
		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				mlr_dsl_cst_handle(pnode->pblock_statements,
					poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs,
					oosvar_flatten_separator, pbind_stack);
			} else {
				break;
			}
		} else {
			break;
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_do_while(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	while (TRUE) {
		mlr_dsl_cst_handle(pnode->pblock_statements,
			poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs,
			oosvar_flatten_separator, pbind_stack);

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (!val.u.boolv) {
				break;
			}
		} else {
			break;
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_for_srec(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	bind_stack_push(pbind_stack, pnode->pbound_variables);
	for (lrece_t* pe = pinrec->phead; pe != NULL; pe = pe->pnext) {
		// xxx beware string/integer for lrec keys ...
		// Copy, not pointer-reference, in case of srec-unset in loop body.
		mv_t mvkey = mv_from_string_with_free(mlr_strdup_or_die(pe->key));

		mv_t* poverlay = lhmsv_get(ptyped_overlay, pe->key);
		mv_t mvval = (poverlay != NULL)
			? mv_copy(poverlay)
			: mv_from_string_with_free(mlr_strdup_or_die(pe->value));

		lhmsmv_put(pnode->pbound_variables, pnode->for_srec_k_name, &mvkey, FREE_ENTRY_VALUE);
		lhmsmv_put(pnode->pbound_variables, pnode->for_srec_v_name, &mvval, FREE_ENTRY_VALUE);

		mlr_dsl_cst_handle(pnode->pblock_statements,
			poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs,
			oosvar_flatten_separator, pbind_stack);
	}
	// xxx break/continue-handling (needs to be in rval evluators w/ stack of brk/ctu flags @ context
	bind_stack_pop(pbind_stack);
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_if_head(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	for (sllve_t* pe = pnode->pif_chain_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pitemnode = pe->pvvalue;
		rval_evaluator_t* prhs_evaluator = pitemnode->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
			ppregex_captures, pctx, prhs_evaluator->pvstate);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				mlr_dsl_cst_handle(pitemnode->pblock_statements,
					poosvars, pinrec, ptyped_overlay, ppregex_captures, pctx, pshould_emit_rec, poutrecs,
					oosvar_flatten_separator, pbind_stack);

				break;
			}
		}
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_cst_node_handle_bare_boolean(
	mlr_dsl_cst_statement_t* pnode,
	mlhmmv_t*        poosvars,
	lrec_t*          pinrec,
	lhmsv_t*         ptyped_overlay,
	string_array_t** ppregex_captures,
	context_t*       pctx,
	int*             pshould_emit_rec,
	sllv_t*          poutrecs,
	char*            oosvar_flatten_separator,
	bind_stack_t*    pbind_stack)
{
	rval_evaluator_t* prhs_evaluator = pnode->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(pinrec, ptyped_overlay, poosvars,
		ppregex_captures, pctx, prhs_evaluator->pvstate);
	if (mv_is_non_null(&val))
		mv_set_boolean_strict(&val);
}

// ----------------------------------------------------------------
// Example ASTs, with and without indexing on the left-hand-side oosvar name:
//
// % mlr put -v '@x[1]["2"][$3][@4]=5' /dev/null
// = (oosvar_assignment):
//     [] (oosvar_level_key):
//         [] (oosvar_level_key):
//             [] (oosvar_level_key):
//                 [] (oosvar_level_key):
//                     x (oosvar_name).
//                     1 (strnum_literal).
//                 2 (strnum_literal).
//             3 (field_name).
//         4 (oosvar_name).
//     5 (strnum_literal).
//
// % mlr put -v '@x = $y'
// = (oosvar_assignment):
//     x (oosvar_name).
//     y (field_name).

// $ mlr put -q -v 'end{emit @v, "a", "b","c"}' ...
// emit (emit):
//     v (oosvar_name).
//     a (strnum_literal).
//     b (strnum_literal).
//     c (strnum_literal).

// mlr put -q -v 'end{emit @v[1][2], "a", "b","c"}' ...
// emit (emit):
//     [] (oosvar_level_key):
//         [] (oosvar_level_key):
//             v (oosvar_name).
//             1 (strnum_literal).
//         2 (strnum_literal).
//     a (strnum_literal).
//     b (strnum_literal).
//     c (strnum_literal).

// pnode is input; pkeylist_evaluators is appended to.
static sllv_t* allocate_keylist_evaluators_from_oosvar_node(mlr_dsl_ast_node_t* pnode, int type_inferencing) {
	sllv_t* pkeylist_evaluators = sllv_alloc();

	if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_NAME) {
		sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pnode->text));
	} else {
		mlr_dsl_ast_node_t* pwalker = pnode;
		while (TRUE) {
			// Bracket operators come in from the right. So the highest AST node is the rightmost index,
			// and the lowest is the oosvar name. Hence sllv_prepend rather than sllv_append.
			if (pwalker->type == MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY) {
				mlr_dsl_ast_node_t* pkeynode = pwalker->pchildren->phead->pnext->pvvalue;
				sllv_prepend(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, type_inferencing));
			} else {
				// Oosvar expressions are of the form '@name[$index1][@index2+3][4]["five"].  The first one
				// (name) is special: syntactically, it's outside the brackets, although that issue is for the
				// parser to handle. Here, it's special since it's always a string, never an expression that
				// evaluates to string.
				sllv_prepend(pkeylist_evaluators, rval_evaluator_alloc_from_string(pwalker->text));
			}
			if (pwalker->pchildren == NULL)
				break;
			pwalker = pwalker->pchildren->phead->pvvalue;
		}
	}
	return pkeylist_evaluators;
}
