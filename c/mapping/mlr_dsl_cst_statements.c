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

// ----------------------------------------------------------------
// Normally I use function typedefs only to simplify function-pointer usage.
// Here, though, they're used to reduce keystrokes in the (very large) number
// of static-function prototypes here at the top of this file.

typedef mlr_dsl_cst_statement_t* cst_statement_allocator_t(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

typedef void cst_statement_handler_t(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs);

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_blank();
void mlr_dsl_cst_statement_free(mlr_dsl_cst_statement_t* pstatement);

static cst_statement_allocator_t alloc_return_value; // UDF
static cst_statement_allocator_t alloc_return_void;  // Subroutine
static cst_statement_allocator_t alloc_subr_callsite;

static cst_statement_allocator_t alloc_local_variable_definition;
static cst_statement_allocator_t alloc_local_variable_assignment;
static cst_statement_allocator_t alloc_srec_assignment;
static cst_statement_allocator_t alloc_indirect_srec_assignment;
static cst_statement_allocator_t alloc_oosvar_assignment;
static cst_statement_allocator_t alloc_oosvar_from_full_srec_assignment;
static cst_statement_allocator_t alloc_full_srec_from_oosvar_assignment;
static cst_statement_allocator_t alloc_env_assignment;
static cst_statement_allocator_t alloc_unset;

static cst_statement_allocator_t alloc_conditional_block;
static cst_statement_allocator_t alloc_if_head;
static cst_statement_allocator_t alloc_while;
static cst_statement_allocator_t alloc_do_while;
static cst_statement_allocator_t alloc_for_srec;
static cst_statement_allocator_t alloc_for_oosvar;
static cst_statement_allocator_t alloc_for_oosvar_key_only;
static cst_statement_allocator_t alloc_triple_for;
static cst_statement_allocator_t alloc_break;
static cst_statement_allocator_t alloc_continue;
static cst_statement_allocator_t alloc_filter;

static cst_statement_allocator_t alloc_bare_boolean;
static mlr_dsl_cst_statement_t* alloc_final_filter(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 negate_final_filter,
	int                 type_inferencing,
	int                 context_flags);

static cst_statement_allocator_t alloc_tee;
static cst_statement_allocator_t alloc_emitf;

static mlr_dsl_cst_statement_t* alloc_emit(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 do_full_prefixing);

static mlr_dsl_cst_statement_t* alloc_emit_lashed(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 do_full_prefixing);

static cst_statement_allocator_t alloc_dump;

static mlr_dsl_cst_statement_t* alloc_print(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	char*               print_terminator);

static file_output_mode_t file_output_mode_from_ast_node_type(mlr_dsl_ast_node_type_t mlr_dsl_ast_node_type);

static mlr_dsl_cst_statement_t* alloc_if_item(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pexprnode,
	mlr_dsl_ast_node_t* plistnode,
	int                 type_inferencing,
	int                 context_flags);

static mlr_dsl_cst_statement_vararg_t* mlr_dsl_cst_statement_vararg_alloc(
	char*             emitf_or_unset_srec_field_name,
	rval_evaluator_t* punset_srec_field_name_evaluator,
	rval_evaluator_t* pemitf_arg_evaluator,
	sllv_t*           punset_oosvar_keylist_evaluators);

static sllv_t* allocate_keylist_evaluators_from_oosvar_node(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags);

static void cst_statement_vararg_free(mlr_dsl_cst_statement_vararg_t* pvararg);

static void handle_statement_list_with_break_continue(
	sllv_t*        pcst_statements,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

static cst_statement_handler_t handle_subr_callsite;
static cst_statement_handler_t handle_return_void;
static cst_statement_handler_t handle_return_value;

static cst_statement_handler_t handle_srec_assignment;
static cst_statement_handler_t handle_indirect_srec_assignment;
static cst_statement_handler_t handle_oosvar_assignment;
static cst_statement_handler_t handle_oosvar_to_oosvar_assignment;
static cst_statement_handler_t handle_oosvar_from_full_srec_assignment;
static cst_statement_handler_t handle_full_srec_from_oosvar_assignment;
static cst_statement_handler_t handle_oosvar_assignment;
static cst_statement_handler_t handle_env_assignment;
static cst_statement_handler_t handle_local_variable_definition;
static cst_statement_handler_t handle_local_variable_assignment;
static cst_statement_handler_t handle_unset;
static cst_statement_handler_t handle_unset_all;

static cst_statement_handler_t handle_filter;
static cst_statement_handler_t handle_final_filter;
static cst_statement_handler_t handle_conditional_block;
static cst_statement_handler_t handle_while;
static cst_statement_handler_t handle_do_while;
static cst_statement_handler_t handle_for_srec;
static cst_statement_handler_t handle_for_oosvar;
static cst_statement_handler_t handle_for_oosvar_key_only;
static cst_statement_handler_t handle_triple_for;
static cst_statement_handler_t handle_break;
static cst_statement_handler_t handle_continue;
static cst_statement_handler_t handle_if_head;
static cst_statement_handler_t handle_bare_boolean;

static void handle_for_oosvar_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_value_t           submap,
	sllse_t*                 prest_for_k_names);

static void handle_unset_vararg_oosvar(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs);

static void handle_unset_vararg_full_srec(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs);

static void handle_unset_vararg_srec_field_name(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs);

static void handle_unset_vararg_indirect_srec_field_name(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs);


static cst_statement_handler_t handle_tee_to_stdfp;
static cst_statement_handler_t handle_tee_to_file;
static lrec_t*                 handle_tee_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs);

static cst_statement_handler_t handle_emitf;
static cst_statement_handler_t handle_emitf_to_stdfp;
static cst_statement_handler_t handle_emitf_to_file;
static void handle_emitf_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  poutrecs);

static cst_statement_handler_t handle_emit;
static cst_statement_handler_t handle_emit_to_stdfp;
static cst_statement_handler_t handle_emit_to_file;
static void handle_emit_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  pcst_outputs,
	char*                    f);

static cst_statement_handler_t handle_emit_lashed;
static cst_statement_handler_t handle_emit_lashed_to_stdfp;
static cst_statement_handler_t handle_emit_lashed_to_file;
static void handle_emit_lashed_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  pcst_outputs,
	char*                    f);

static cst_statement_handler_t handle_emit_all;
static cst_statement_handler_t handle_emit_all_to_stdfp;
static cst_statement_handler_t handle_emit_all_to_file;
static cst_statement_handler_t handle_dump;
static cst_statement_handler_t handle_dump_to_file;
static cst_statement_handler_t handle_print;

// ================================================================
cst_top_level_statement_block_t* cst_top_level_statement_block_alloc(int max_var_depth) {
	cst_top_level_statement_block_t* pblock = mlr_malloc_or_die(sizeof(cst_top_level_statement_block_t));

	pblock->max_var_depth = max_var_depth;
	pblock->pstack        = local_stack_alloc(max_var_depth);
	pblock->pstatements   = sllv_alloc();

	return pblock;
}

// ----------------------------------------------------------------
void cst_top_level_statement_block_free(cst_top_level_statement_block_t* pblock) {
	if (pblock == NULL)
		return;
	local_stack_free(pblock->pstack);
	for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_free(pe->pvvalue);
	}
	sllv_free(pblock->pstatements);
	free(pblock);
}

// ----------------------------------------------------------------
cst_statement_block_t* cst_statement_block_alloc(int frame_var_count) {
	cst_statement_block_t* pblock = mlr_malloc_or_die(sizeof(cst_statement_block_t));

	pblock->frame_var_count = frame_var_count;
	pblock->pstatements     = sllv_alloc();

	return pblock;
}

// ----------------------------------------------------------------
void cst_statement_block_free(cst_statement_block_t* pblock) {
	if (pblock == NULL)
		return;
	for (sllve_t* pe = pblock->pstatements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_free(pe->pvvalue);
	}
	sllv_free(pblock->pstatements);
	free(pblock);
}

// ================================================================
// The parser accepts many things that are invalid, e.g.
// * begin{end{}} -- begin/end not at top level
// * begin{$x=1} -- references to stream records at begin/end
// * break/continue outside of for/while/do-while
// * $x=x -- boundvars outside of for-loop variable bindings
//
// All of the above are enforced here by the CST builder, which takes the parser's output AST as
// input.  This is done (a) to keep the parser from being overly complex, and (b) so we can get much
// more informative error messages in C than in Lemon ('syntax error').
//
// In this file we set up left-hand sides for assignments, as well as right-hand sides for emit and
// unset.  Most right-hand sides are set up in rval_expr_evaluators.c so the context_flags are
// passed through to there as well.

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_FUNC_DEF:
		fprintf(stderr, "%s: func statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_SUBR_DEF:
		fprintf(stderr, "%s: subr statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_BEGIN:
		fprintf(stderr, "%s: begin statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_END:
		fprintf(stderr, "%s: end statements are only valid at top level.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_RETURN_VALUE:
		if (!(context_flags & IN_FUNC_DEF)) {
			fprintf(stderr, "%s: return-value statements are only valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_return_value(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_RETURN_VOID:
		if (!(context_flags & IN_SUBR_DEF)) {
			fprintf(stderr, "%s: return-void statements are only valid within subr blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_return_void(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_SUBR_CALLSITE:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: subroutine calls are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_subr_callsite(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:
		return alloc_conditional_block(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_IF_HEAD:
		return alloc_if_head(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_WHILE:
		return alloc_while(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_DO_WHILE:
		return alloc_do_while(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_SREC:
		return alloc_for_srec(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR:
		return alloc_for_oosvar(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY:
		return alloc_for_oosvar_key_only(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;
	case MD_AST_NODE_TYPE_TRIPLE_FOR:
		return alloc_triple_for(pcst, pnode, type_inferencing, context_flags | IN_BREAKABLE);
		break;

	case MD_AST_NODE_TYPE_BREAK:
		if (!(context_flags & IN_BREAKABLE)) {
			fprintf(stderr, "%s: break statements are only valid within for, while, or do-while.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_break(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_CONTINUE:
		if (!(context_flags & IN_BREAKABLE)) {
			fprintf(stderr, "%s: break statements are only valid within for, while, or do-while.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_continue(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_LOCAL_DEFINITION:
		return alloc_local_variable_definition(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_LOCAL_ASSIGNMENT:
		return alloc_local_variable_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_indirect_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:
		return alloc_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments from $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_oosvar_from_full_srec_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT:
		if (context_flags & IN_BEGIN_OR_END) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within begin or end blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: assignments to $-variables are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_full_srec_from_oosvar_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_ENV_ASSIGNMENT:
		return alloc_env_assignment(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_UNSET:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: unset statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_unset(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_TEE:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: tee statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_tee(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_EMITF:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitf statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emitf(pcst, pnode, type_inferencing, context_flags);
		break;
	case MD_AST_NODE_TYPE_EMITP:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitp statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit(pcst, pnode, type_inferencing, context_flags, TRUE);
		break;
	case MD_AST_NODE_TYPE_EMIT:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emit statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit(pcst, pnode, type_inferencing, context_flags, FALSE);
		break;

	case MD_AST_NODE_TYPE_EMITP_LASHED:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emitp statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit_lashed(pcst, pnode, type_inferencing, context_flags, TRUE);
		break;
	case MD_AST_NODE_TYPE_EMIT_LASHED:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: emit statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_emit_lashed(pcst, pnode, type_inferencing, context_flags, FALSE);
		break;

	case MD_AST_NODE_TYPE_FILTER:
		if (context_flags & IN_FUNC_DEF) {
			fprintf(stderr, "%s: filter statements are not valid within func blocks.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		if (context_flags & IN_MLR_FILTER) {
			fprintf(stderr, "%s filter: expressions must not also contain the \"filter\" keyword.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		return alloc_filter(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_DUMP:
		return alloc_dump(pcst, pnode, type_inferencing, context_flags);
		break;

	case MD_AST_NODE_TYPE_PRINT:
		return alloc_print(pcst, pnode, type_inferencing, context_flags, "\n");
		break;

	case MD_AST_NODE_TYPE_PRINTN:
		return alloc_print(pcst, pnode, type_inferencing, context_flags, "");
		break;

	default:
		return alloc_bare_boolean(pcst, pnode, type_inferencing, context_flags);
		break;
	}
}

// ----------------------------------------------------------------
// mlr put and mlr filter are almost entirely the same code. The key difference is that the final
// statement for the latter must be a bare boolean expression.

mlr_dsl_cst_statement_t* mlr_dsl_cst_alloc_final_filter_statement(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int negate_final_filter, int type_inferencing, int context_flags)
{
	switch(pnode->type) {

	case MD_AST_NODE_TYPE_FILTER:
		fprintf(stderr, "%s filter: expressions must not also contain the \"filter\" keyword.\n",
			MLR_GLOBALS.bargv0);
		exit(1);
		break;

	case MD_AST_NODE_TYPE_FUNC_DEF:
	case MD_AST_NODE_TYPE_SUBR_DEF:
	case MD_AST_NODE_TYPE_BEGIN:
	case MD_AST_NODE_TYPE_END:
	case MD_AST_NODE_TYPE_RETURN_VALUE:
	case MD_AST_NODE_TYPE_RETURN_VOID:
	case MD_AST_NODE_TYPE_SUBR_CALLSITE:
	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:
	case MD_AST_NODE_TYPE_IF_HEAD:
	case MD_AST_NODE_TYPE_WHILE:
	case MD_AST_NODE_TYPE_DO_WHILE:
	case MD_AST_NODE_TYPE_FOR_SREC:
	case MD_AST_NODE_TYPE_FOR_OOSVAR:
	case MD_AST_NODE_TYPE_TRIPLE_FOR:
	case MD_AST_NODE_TYPE_BREAK:
	case MD_AST_NODE_TYPE_CONTINUE:
	case MD_AST_NODE_TYPE_LOCAL_DEFINITION:
	case MD_AST_NODE_TYPE_LOCAL_ASSIGNMENT:
	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:
	case MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT:
	case MD_AST_NODE_TYPE_UNSET:
	case MD_AST_NODE_TYPE_TEE:
	case MD_AST_NODE_TYPE_EMITF:
	case MD_AST_NODE_TYPE_EMITP:
	case MD_AST_NODE_TYPE_EMIT:
	case MD_AST_NODE_TYPE_EMITP_LASHED:
	case MD_AST_NODE_TYPE_EMIT_LASHED:
	case MD_AST_NODE_TYPE_DUMP:
	case MD_AST_NODE_TYPE_PRINT:
	case MD_AST_NODE_TYPE_PRINTN:
		fprintf(stderr, "%s: filter expressions must end in a final boolean statement.\n", MLR_GLOBALS.bargv0);
		exit(1);
		break;

	default:
		return alloc_final_filter(pcst, pnode, negate_final_filter, type_inferencing, context_flags);
		break;
	}
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_blank() {
	mlr_dsl_cst_statement_t* pstatement = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_t));

	pstatement->pnode_handler                           = NULL;
	pstatement->subr_callsite_argument_evaluators       = NULL;
	pstatement->subr_callsite_arguments                 = NULL;
	pstatement->psubr_callsite                          = NULL;
	pstatement->psubr_defsite                           = NULL;
	pstatement->local_variable_name                     = NULL;
	pstatement->preturn_evaluator                       = NULL;
	pstatement->pblock_handler                          = NULL;
	pstatement->poosvar_lhs_keylist_evaluators          = NULL;
	pstatement->pemit_keylist_evaluators                = NULL;
	pstatement->num_emit_keylist_evaluators             = 0;
	pstatement->ppemit_keylist_evaluators               = NULL;
	pstatement->local_lhs_variable_name                 = NULL;
	pstatement->srec_lhs_field_name                     = NULL;
	pstatement->env_lhs_name                            = NULL;
	pstatement->psrec_lhs_evaluator                     = NULL;
	pstatement->prhs_evaluator                          = NULL;
	pstatement->stdfp                                   = NULL;
	pstatement->print_terminator                        = "\n";
	pstatement->poutput_filename_evaluator              = NULL;
	pstatement->file_output_mode                        = MODE_WRITE;
	pstatement->pmulti_out                              = NULL;
	pstatement->psingle_lrec_writer                     = NULL;
	pstatement->pmulti_lrec_writer                      = NULL;
	pstatement->poosvar_rhs_keylist_evaluators          = NULL;
	pstatement->pemit_oosvar_namelist_evaluators        = NULL;
	pstatement->pvarargs                                = NULL;
	pstatement->do_full_prefixing                       = FALSE;
	pstatement->flush_every_record                      = TRUE;
	pstatement->pblock_statements                       = NULL;
	pstatement->pif_chain_statements                    = NULL;
	pstatement->pfor_oosvar_k_names                     = NULL;
	pstatement->for_v_name                              = NULL;
	pstatement->ptype_infererenced_srec_field_getter    = NULL;
	pstatement->ptriple_for_start_statements            = NULL;
	pstatement->ptriple_for_pre_continuation_statements = NULL;
	pstatement->ptriple_for_continuation_evaluator      = NULL;
	pstatement->ptriple_for_update_statements           = NULL;
	pstatement->pframe                                  = NULL;
	pstatement->negate_final_filter                     = FALSE;

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_local_variable_definition(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	mlr_dsl_ast_node_t* pname_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvalue_node = pnode->pchildren->phead->pnext->pvvalue;
	pstatement->local_variable_name = pname_node->text;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pvalue_node, pcst->pfmgr,
		type_inferencing, context_flags);
	pstatement->pnode_handler = handle_local_variable_definition;
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_local_variable_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	if ((pnode->pchildren == NULL) || (pnode->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_LOCAL_VARIABLE) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	} else if (pleft->pchildren != NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->pnode_handler = handle_local_variable_assignment;
	pstatement->local_lhs_variable_name = pleft->text;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_return_value(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	mlr_dsl_ast_node_t* prhs_node = pnode->pchildren->phead->pvvalue;
	pstatement->preturn_evaluator = rval_evaluator_alloc_from_ast(prhs_node, pcst->pfmgr,
		type_inferencing, context_flags);
	pstatement->pnode_handler = handle_return_value;
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_return_void(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	pstatement->pnode_handler = handle_return_void;
	return pstatement;
}

// ----------------------------------------------------------------
// xxx rename
static mlr_dsl_cst_statement_t* alloc_subr_callsite(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pname_node = pnode->pchildren->phead->pvvalue;
	int callsite_arity = pname_node->pchildren->length;

	pstatement->psubr_callsite = subr_callsite_alloc(pname_node->text, callsite_arity,
		type_inferencing, context_flags);
	// Remember this callsite to be resolved later, after all subroutine definitions have been done.
	sllv_append(pcst->psubr_callsite_statements_to_resolve, pstatement);

	pstatement->subr_callsite_argument_evaluators = mlr_malloc_or_die(callsite_arity * sizeof(rval_evaluator_t*));
	pstatement->subr_callsite_arguments = mlr_malloc_or_die(callsite_arity * sizeof(mv_t));

	int i = 0;
	for (sllve_t* pe = pname_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pargument_node = pe->pvvalue;
		pstatement->subr_callsite_argument_evaluators[i] = rval_evaluator_alloc_from_ast(pargument_node,
			pcst->pfmgr, type_inferencing, context_flags);
	}

	pstatement->pnode_handler = handle_subr_callsite;
	return pstatement;
}

subr_callsite_t* subr_callsite_alloc(char* name, int arity, int type_inferencing, int context_flags) {
	subr_callsite_t* psubr_callsite  = mlr_malloc_or_die(sizeof(subr_callsite_t));
	psubr_callsite->name             = mlr_strdup_or_die(name);
	psubr_callsite->arity            = arity;
	psubr_callsite->type_inferencing = type_inferencing;
	psubr_callsite->context_flags    = context_flags;
	return psubr_callsite;
}

void subr_callsite_free(subr_callsite_t* psubr_callsite) {
	if (psubr_callsite == NULL)
		return;
	free(psubr_callsite->name);
	free(psubr_callsite);
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	if ((pnode->pchildren == NULL) || (pnode->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_FIELD_NAME) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	} else if (pleft->pchildren != NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->pnode_handler = handle_srec_assignment;
	pstatement->srec_lhs_field_name = pleft->text;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
// $ mlr --from ../data/small put -v '$[@x] = 1'
// AST ROOT:
// text="list", type=statement_list:
//     text="=", type=indirect_srec_assignment:
//         text="oosvar_keylist", type=oosvar_keylist:
//             text="x", type=string_literal.
//         text="1", type=strnum_literal.

static mlr_dsl_cst_statement_t* alloc_indirect_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	if ((pnode->pchildren == NULL) || (pnode->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	pstatement->pnode_handler = handle_indirect_srec_assignment;
	pstatement->psrec_lhs_evaluator = rval_evaluator_alloc_from_ast(pleft,  pcst->pfmgr,
		type_inferencing, context_flags);
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr,
		type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_oosvar_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	sllv_t* poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pcst, pleft,
		type_inferencing, context_flags);

	if (pleft->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST && pright->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
		pstatement->pnode_handler = handle_oosvar_to_oosvar_assignment;
		pstatement->poosvar_rhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pcst, pright,
			type_inferencing, context_flags);
	} else {
		pstatement->pnode_handler = handle_oosvar_assignment;
		pstatement->poosvar_rhs_keylist_evaluators = NULL;
	}

	pstatement->poosvar_lhs_keylist_evaluators = poosvar_lhs_keylist_evaluators;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr,
		type_inferencing, context_flags);

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_oosvar_from_full_srec_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	if (pright->type != MD_AST_NODE_TYPE_FULL_SREC) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->pnode_handler = handle_oosvar_from_full_srec_assignment;
	pstatement->poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pcst, pleft,
		type_inferencing, context_flags);
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_full_srec_from_oosvar_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_FULL_SREC) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	if (pright->type != MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->pnode_handler = handle_full_srec_from_oosvar_assignment;
	pstatement->poosvar_rhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pcst, pright,
		type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_env_assignment(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	if ((pnode->pchildren == NULL) || (pnode->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	if (pleft->type != MD_AST_NODE_TYPE_ENV) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	} else if (pleft->pchildren == NULL) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	} else if (pleft->pchildren->length != 2) {
		fprintf(stderr, "%s: coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	mlr_dsl_ast_node_t* pnamenode = pleft->pchildren->phead->pnext->pvvalue;

	pstatement->pnode_handler = handle_env_assignment;
	pstatement->env_lhs_name = pnamenode->text;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_unset(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	pstatement->pnode_handler = handle_unset;
	pstatement->pvarargs = sllv_alloc();
	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnode = pe->pvvalue;

		if (pnode->type == MD_AST_NODE_TYPE_ALL || pnode->type == MD_AST_NODE_TYPE_FULL_OOSVAR) {
			// The grammar allows only 'unset all', not 'unset @x, all, $y'.
			// So if 'all' appears at all, it's the only name. Likewise with '@*'.
			pstatement->pnode_handler = handle_unset_all;

		} else if (pnode->type == MD_AST_NODE_TYPE_FULL_SREC) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables is not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				NULL,
				NULL,
				NULL,
				NULL));

		} else if (pnode->type == MD_AST_NODE_TYPE_FIELD_NAME) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables is not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				pnode->text,
				NULL,
				NULL,
				NULL));

		} else if (pnode->type == MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME) {
			if (context_flags & IN_BEGIN_OR_END) {
				fprintf(stderr, "%s: unset of $-variables are not valid within begin or end blocks.\n",
					MLR_GLOBALS.bargv0);
				exit(1);
			}
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				NULL,
				rval_evaluator_alloc_from_ast(pnode->pchildren->phead->pvvalue, pcst->pfmgr,
					type_inferencing, context_flags),
				NULL,
				NULL));

		} else if (pnode->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {
			sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
				NULL,
				NULL,
				NULL,
				allocate_keylist_evaluators_from_oosvar_node(pcst, pnode,
					type_inferencing, context_flags)));

		} else {
			fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
				MLR_GLOBALS.bargv0, __FILE__, __LINE__);
			exit(1);
		}
	}
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_while(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left child node is the AST for the boolean expression.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->pframe = bind_stack_frame_alloc_unfenced();
	pstatement->pnode_handler = handle_while;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pleft, pcst->pfmgr, type_inferencing, context_flags);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_do_while(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left child node is the list of statements in the body.
	// Right child node is the AST for the boolean expression.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = pleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->pframe = bind_stack_frame_alloc_unfenced();
	pstatement->pnode_handler = handle_do_while;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pright, pcst->pfmgr, type_inferencing, context_flags);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'for (k,v in $*) { $x=1; $y=2 }'
// AST ROOT:
// text="list", type=statement_list:
//     text="for", type=for_srec:
//         text="variables", type=for_variables:
//             text="k", type=non_sigil_name.
//             text="v", type=non_sigil_name.
//         text="list", type=statement_list:
//             text="=", type=srec_assignment:
//                 text="x", type=field_name.
//                 text="1", type=strnum_literal.
//             text="=", type=srec_assignment:
//                 text="y", type=field_name.
//                 text="2", type=strnum_literal.

static mlr_dsl_cst_statement_t* alloc_for_srec(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left child node is list of bound variables.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pknode = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pvnode = pleft->pchildren->phead->pnext->pvvalue;

	if (streq(pknode->text, pvnode->text)) {
		fprintf(stderr, "%s: duplicate for-loop boundvars \"%s\" and \"%s\".\n",
			MLR_GLOBALS.bargv0, pknode->text, pvnode->text);
		exit(1);
	}
	pstatement->for_srec_k_name = pknode->text;
	pstatement->for_v_name = pvnode->text;

	sllv_t* pblock_statements = sllv_alloc();
	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	pstatement->pnode_handler = handle_for_srec;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;
	pstatement->pblock_statements = pblock_statements;
	pstatement->pframe = bind_stack_frame_alloc_unfenced();
	pstatement->ptype_infererenced_srec_field_getter =
		(type_inferencing == TYPE_INFER_STRING_ONLY)      ? get_srec_value_string_only_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT)     ? get_srec_value_string_float_aux :
		(type_inferencing == TYPE_INFER_STRING_FLOAT_INT) ? get_srec_value_string_float_int_aux :
		NULL;
	if (pstatement->ptype_infererenced_srec_field_getter == NULL) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	return pstatement;
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'for((k1,k2,k3),v in @a["4"][$5]) { $6 = 7; $8 = 9}'
// AST ROOT:
// text="list", type=statement_list:
//     text="for", type=for_oosvar:
//         text="key_and_value_variables", type=for_variables:
//             text="key_variables", type=for_variables:
//                 text="k1", type=non_sigil_name.
//                 text="k2", type=non_sigil_name.
//                 text="k3", type=non_sigil_name.
//             text="v", type=non_sigil_name.
//         text="oosvar_keylist", type=oosvar_keylist:
//             text="a", type=string_literal.
//             text="4", type=strnum_literal.
//             text="5", type=field_name.
//         text="list", type=statement_list:
//             text="=", type=srec_assignment:
//                 text="6", type=field_name.
//                 text="7", type=strnum_literal.
//             text="=", type=srec_assignment:
//                 text="8", type=field_name.
//                 text="9", type=strnum_literal.

static mlr_dsl_cst_statement_t* alloc_for_oosvar(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left child node is list of bound variables.
	//   Left subnode is namelist for key boundvars.
	//   Right subnode is name for value boundvar.
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubleft  = pleft->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* psubright = pleft->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	pstatement->pfor_oosvar_k_names = slls_alloc();
	int ok = TRUE;
	hss_t* pnameset = hss_alloc();
	for (sllve_t* pe = psubleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pnamenode = pe->pvvalue;
		slls_append_with_free(pstatement->pfor_oosvar_k_names, mlr_strdup_or_die(pnamenode->text));
		if (hss_has(pnameset, pnamenode->text)) {
			fprintf(stderr, "%s: duplicate for-loop boundvar \"%s\".\n",
				MLR_GLOBALS.bargv0, pnamenode->text);
			ok = FALSE;
		}
		hss_add(pnameset, pnamenode->text);
	}
	pstatement->for_v_name = psubright->text;
	if (hss_has(pnameset, psubright->text)) {
		fprintf(stderr, "%s: duplicate for-loop boundvar \"%s\".\n",
			MLR_GLOBALS.bargv0, psubright->text);
		ok = FALSE;
	}
	hss_add(pnameset, psubright->text);
	if (!ok) {
		fprintf(stderr, "Boundvars: ");
		for (sllve_t* pe = psubleft->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pnamenode = pe->pvvalue;
			fprintf(stderr, "\"%s\", ", pnamenode->text);
		}
		fprintf(stderr, "\"%s\".\n", psubright->text);
		exit(1);
	}

	hss_free(pnameset);

	pstatement->poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(
		pcst, pmiddle, type_inferencing, context_flags);

	sllv_t* pblock_statements = sllv_alloc();
	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}
	pstatement->pblock_statements = pblock_statements;
	pstatement->pframe = bind_stack_frame_alloc_unfenced();

	pstatement->pnode_handler = handle_for_oosvar;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;

	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_for_oosvar_key_only(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left child node is single bound variable
	// Middle child node is keylist for basepoint in the oosvar mlhmmv.
	// Right child node is the list of statements in the body.
	mlr_dsl_ast_node_t* pleft     = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pmiddle   = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pright    = pnode->pchildren->phead->pnext->pnext->pvvalue;

	pstatement->pfor_oosvar_k_names = slls_alloc();
	slls_append_with_free(pstatement->pfor_oosvar_k_names, mlr_strdup_or_die(pleft->text));

	pstatement->poosvar_lhs_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(
		pcst, pmiddle, type_inferencing, context_flags);

	sllv_t* pblock_statements = sllv_alloc();
	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pblock_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}
	pstatement->pblock_statements = pblock_statements;
	pstatement->pframe = bind_stack_frame_alloc_unfenced();

	pstatement->pnode_handler = handle_for_oosvar_key_only;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;

	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_triple_for(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pstart_statements_node        = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pcontinuation_statements_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pupdate_statements_node       = pnode->pchildren->phead->pnext->pnext->pvvalue;
	mlr_dsl_ast_node_t* pbody_statements_node         = pnode->pchildren->phead->pnext->pnext->pnext->pvvalue;

	pstatement->ptriple_for_start_statements = sllv_alloc();
	for (sllve_t* pe = pstart_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pstatement->ptriple_for_start_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags & ~IN_BREAKABLE));
	}

	// Continuation statements are split into the final boolean, and the statements before (if any).
	pstatement->ptriple_for_pre_continuation_statements = sllv_alloc();
	// Empty continuation for triple-for is an implicit TRUE.
	if (pcontinuation_statements_node->pchildren->length == 0) {
		pstatement->ptriple_for_continuation_evaluator = rval_evaluator_alloc_from_boolean(TRUE);
	} else {
		for (sllve_t* pe = pcontinuation_statements_node->pchildren->phead; pe != NULL && pe->pnext != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
			sllv_append(pstatement->ptriple_for_pre_continuation_statements,
				mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
				type_inferencing, context_flags & ~IN_BREAKABLE));
		}
		mlr_dsl_ast_node_t* pfinal_continuation_statement_node = pcontinuation_statements_node->pchildren->ptail->pvvalue;
		if (mlr_dsl_ast_node_cannot_be_bare_boolean(pfinal_continuation_statement_node)) {
			fprintf(stderr,
				"%s: the final triple-for continutation statement must be a bare boolean.\n",
				MLR_GLOBALS.bargv0);
			exit(1);
		}
		pstatement->ptriple_for_continuation_evaluator = rval_evaluator_alloc_from_ast(pfinal_continuation_statement_node,
			pcst->pfmgr, type_inferencing, (context_flags & ~IN_BREAKABLE) | IN_TRIPLE_FOR_CONTINUE);
	}

	pstatement->ptriple_for_update_statements = sllv_alloc();
	for (sllve_t* pe = pupdate_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pstatement->ptriple_for_update_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags & ~IN_BREAKABLE));
	}

	pstatement->pblock_statements = sllv_alloc();
	for (sllve_t* pe = pbody_statements_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		sllv_append(pstatement->pblock_statements, mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags));
	}

	pstatement->pframe = bind_stack_frame_alloc_unfenced();

	pstatement->pnode_handler = handle_triple_for;
	pstatement->pblock_handler = handle_statement_list_with_break_continue;

	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_break(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	pstatement->pnode_handler = handle_break;
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_continue(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	pstatement->pnode_handler = handle_continue;
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_conditional_block(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	// Left node is the AST for the boolean expression.
	// Right node is a list of statements to be executed if the left evaluates to true.
	mlr_dsl_ast_node_t* pleft  = pnode->pchildren->phead->pvvalue;
	sllv_t* pblock_statements = sllv_alloc();

	mlr_dsl_ast_node_t* pright = pnode->pchildren->phead->pnext->pvvalue;
	for (sllve_t* pe = pright->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->pframe = bind_stack_frame_alloc_unfenced();
	pstatement->pnode_handler = handle_conditional_block;
	pstatement->pblock_handler = (context_flags & IN_BREAKABLE)
		? handle_statement_list_with_break_continue
		: mlr_dsl_cst_handle_statement_list;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pleft, pcst->pfmgr, type_inferencing, context_flags);
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

// ----------------------------------------------------------------
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

static mlr_dsl_cst_statement_t* alloc_if_head(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

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

		sllv_append(pif_chain_statements, alloc_if_item(pcst, pexprnode, plistnode,
			type_inferencing, context_flags));
	}

	pstatement->pnode_handler = handle_if_head;
	pstatement->pblock_handler = (context_flags & IN_BREAKABLE)
		?  handle_statement_list_with_break_continue
		: mlr_dsl_cst_handle_statement_list;
	pstatement->pif_chain_statements = pif_chain_statements;
	return pstatement;
}

static mlr_dsl_cst_statement_t* alloc_if_item(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pexprnode,
	mlr_dsl_ast_node_t* plistnode, int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	sllv_t* pblock_statements = sllv_alloc();

	for (sllve_t* pe = plistnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pbody_ast_node = pe->pvvalue;
		mlr_dsl_cst_statement_t *pstatement = mlr_dsl_cst_alloc_statement(pcst, pbody_ast_node,
			type_inferencing, context_flags);
		sllv_append(pblock_statements, pstatement);
	}

	pstatement->pframe = bind_stack_frame_alloc_unfenced();
	pstatement->pnode_handler = NULL; // handled by the containing if-head evaluator
	pstatement->prhs_evaluator = pexprnode != NULL
		? rval_evaluator_alloc_from_ast(pexprnode, pcst->pfmgr,
			type_inferencing, context_flags) // if-statement or elif-statement
		: rval_evaluator_alloc_from_boolean(TRUE); // else-statement
	pstatement->pblock_statements = pblock_statements;
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_filter(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pchild = pnode->pchildren->phead->pvvalue;

	pstatement->pnode_handler = handle_filter;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pchild, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_bare_boolean(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	pstatement->pnode_handler = handle_bare_boolean;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pnode, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_final_filter(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int negate_final_filter, int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	pstatement->pnode_handler = handle_final_filter;
	pstatement->negate_final_filter = negate_final_filter;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pnode, pcst->pfmgr, type_inferencing, context_flags);
	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_tee(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;

	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->pnode_handler = handle_tee_to_stdfp;
		pstatement->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pnode_handler = handle_tee_to_file;
	}
	pstatement->flush_every_record = pcst->flush_every_record;

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_emitf(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pnamesnode = pnode->pchildren->phead->pvvalue;

	// Loop over oosvar names to emit in e.g. 'emitf @a, @b, @c'.
	pstatement->pvarargs = sllv_alloc();
	for (sllve_t* pe = pnamesnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pwalker = pe->pvvalue;
		mlr_dsl_ast_node_t* pchild = pwalker->pchildren->phead->pvvalue;
		// This could be enforced in the lemon parser but it's easier to do it here.
		sllv_append(pstatement->pvarargs, mlr_dsl_cst_statement_vararg_alloc(
			pchild->text,
			NULL,
			rval_evaluator_alloc_from_ast(pwalker, pcst->pfmgr, type_inferencing, context_flags),
			NULL));
	}

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		pstatement->pnode_handler = handle_emitf;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->pnode_handler = handle_emitf_to_stdfp;
		pstatement->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pnode_handler = handle_emitf_to_file;
	}
	pstatement->flush_every_record = pcst->flush_every_record;

	return pstatement;
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'emit @a[2][3], "x", "y", "z"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="a", type=string_literal.
//                 text="2", type=strnum_literal.
//                 text="3", type=strnum_literal.
//             text="emit_namelist", type=emit:
//                 text="x", type=strnum_literal.
//                 text="y", type=strnum_literal.
//                 text="z", type=strnum_literal.
//         text="stream", type=stream:
//
// $ mlr -n put -v 'emit all, "x", "y", "z"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="all", type=all.
//             text="emit_namelist", type=emit:
//                 text="x", type=strnum_literal.
//                 text="y", type=strnum_literal.
//                 text="z", type=strnum_literal.
//         text="stream", type=stream:

static mlr_dsl_cst_statement_t* alloc_emit(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags, int do_full_prefixing)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pemit_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pkeylist_node = pemit_node->pchildren->phead->pvvalue;

	int output_all = FALSE;
	// The grammar allows only 'emit all', not 'emit @x, all, $y'.
	// So if 'all' appears at all, it's the only name.
	if (pkeylist_node->type == MD_AST_NODE_TYPE_ALL || pkeylist_node->type == MD_AST_NODE_TYPE_FULL_OOSVAR) {
		output_all = TRUE;

		sllv_t* pemit_oosvar_namelist_evaluators = sllv_alloc();
		if (pemit_node->pchildren->length == 2) {
			mlr_dsl_ast_node_t* pnamelist_node = pemit_node->pchildren->phead->pnext->pvvalue;
			for (sllve_t* pe = pnamelist_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
				sllv_append(pemit_oosvar_namelist_evaluators,
					rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr, type_inferencing, context_flags));
			}
		}
		pstatement->pemit_oosvar_namelist_evaluators = pemit_oosvar_namelist_evaluators;

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {

		pstatement->pemit_keylist_evaluators = allocate_keylist_evaluators_from_oosvar_node(pcst, pkeylist_node,
			type_inferencing, context_flags);

		sllv_t* pemit_oosvar_namelist_evaluators = sllv_alloc();
		if (pemit_node->pchildren->length == 2) {
			mlr_dsl_ast_node_t* pnamelist_node = pemit_node->pchildren->phead->pnext->pvvalue;
			for (sllve_t* pe = pnamelist_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
				mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
				sllv_append(pemit_oosvar_namelist_evaluators,
					rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr, type_inferencing, context_flags));
			}
		}
		pstatement->pemit_oosvar_namelist_evaluators = pemit_oosvar_namelist_evaluators;

	} else {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}

	pstatement->do_full_prefixing = do_full_prefixing;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		pstatement->pnode_handler = output_all ? handle_emit_all : handle_emit;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->pnode_handler = output_all ? handle_emit_all_to_stdfp : handle_emit_to_stdfp;
		pstatement->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pnode_handler = output_all ? handle_emit_all_to_file : handle_emit_to_file;
	}
	pstatement->flush_every_record = pcst->flush_every_record;

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_emit_lashed(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags, int do_full_prefixing)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* pemit_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;

	mlr_dsl_ast_node_t* pkeylists_node = pemit_node->pchildren->phead->pvvalue;

	pstatement->num_emit_keylist_evaluators = pkeylists_node->pchildren->length;
	pstatement->ppemit_keylist_evaluators = mlr_malloc_or_die(pstatement->num_emit_keylist_evaluators
		* sizeof(sllv_t*));
	int i = 0;
	for (sllve_t* pe = pkeylists_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pkeylist_node = pe->pvvalue;
		pstatement->ppemit_keylist_evaluators[i] = allocate_keylist_evaluators_from_oosvar_node(pcst,
			pkeylist_node, type_inferencing, context_flags);
	}

	sllv_t* pemit_oosvar_namelist_evaluators = sllv_alloc();
	if (pemit_node->pchildren->length == 2) {
		mlr_dsl_ast_node_t* pnamelist_node = pemit_node->pchildren->phead->pnext->pvvalue;
		for (sllve_t* pe = pnamelist_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			sllv_append(pemit_oosvar_namelist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr, type_inferencing, context_flags));
		}
	}
	pstatement->pemit_oosvar_namelist_evaluators = pemit_oosvar_namelist_evaluators;

	pstatement->do_full_prefixing = do_full_prefixing;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		pstatement->pnode_handler = handle_emit_lashed;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->pnode_handler = handle_emit_lashed_to_stdfp;
		pstatement->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pnode_handler = handle_emit_lashed_to_file;
	}

	return pstatement;
}

// ----------------------------------------------------------------
// Example ASTs, with and without indexing on the left-hand-side oosvar name:

// $ mlr -n put -v '@x[1]["2"][$3][@4]=5'
// AST ROOT:
// text="list", type=statement_list:
//     text="=", type=oosvar_assignment:
//         text="oosvar_keylist", type=oosvar_keylist:
//             text="x", type=string_literal.
//             text="1", type=strnum_literal.
//             text="2", type=strnum_literal.
//             text="3", type=field_name.
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="4", type=string_literal.
//         text="5", type=strnum_literal.
//
// $ mlr -n put -v '@x = $y'
// AST ROOT:
// text="list", type=statement_list:
//     text="=", type=oosvar_assignment:
//         text="oosvar_keylist", type=oosvar_keylist:
//             text="x", type=string_literal.
//         text="y", type=field_name.
//
// $ mlr -n put -q -v 'emit @v, "a", "b", "c"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="v", type=string_literal.
//             text="emit_namelist", type=emit:
//                 text="a", type=strnum_literal.
//                 text="b", type=strnum_literal.
//                 text="c", type=strnum_literal.
//         text="stream", type=stream:
//
// $ mlr -n put -q -v 'emit @v[1][2], "a", "b","c"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="v", type=string_literal.
//                 text="1", type=strnum_literal.
//                 text="2", type=strnum_literal.
//             text="emit_namelist", type=emit:
//                 text="a", type=strnum_literal.
//                 text="b", type=strnum_literal.
//                 text="c", type=strnum_literal.
//         text="stream", type=stream:

// pnode is input; pkeylist_evaluators is appended to.
static sllv_t* allocate_keylist_evaluators_from_oosvar_node(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	sllv_t* pkeylist_evaluators = sllv_alloc();

	for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
		if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
		} else {
			sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr,
				type_inferencing, context_flags));
		}
	}
	return pkeylist_evaluators;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_dump(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;
	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) {
		pstatement->pnode_handler = handle_dump;
		pstatement->stdfp = stdout;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->pnode_handler = handle_dump;
		pstatement->stdfp = stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pmulti_out = multi_out_alloc();
		pstatement->pnode_handler = handle_dump_to_file;
	}
	pstatement->flush_every_record = pcst->flush_every_record;

	return pstatement;
}

// ----------------------------------------------------------------
static mlr_dsl_cst_statement_t* alloc_print(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags, char* print_terminator)
{
	if ((pnode->pchildren == NULL) || (pnode->pchildren->length != 2)) {
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
	mlr_dsl_cst_statement_t* pstatement = alloc_blank();
	mlr_dsl_ast_node_t* pvalue_node = pnode->pchildren->phead->pvvalue;
	pstatement->prhs_evaluator = rval_evaluator_alloc_from_ast(pvalue_node, pcst->pfmgr,
		type_inferencing, context_flags);
	pstatement->print_terminator = print_terminator;

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;
	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) {
		pstatement->stdfp = stdout;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstatement->stdfp = stderr;
	} else {
		pstatement->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstatement->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstatement->pmulti_out = multi_out_alloc();
	}
	pstatement->flush_every_record = pcst->flush_every_record;
	pstatement->pnode_handler = handle_print;

	return pstatement;
}

// ----------------------------------------------------------------
static file_output_mode_t file_output_mode_from_ast_node_type(mlr_dsl_ast_node_type_t mlr_dsl_ast_node_type) {
	switch(mlr_dsl_ast_node_type) {
	case MD_AST_NODE_TYPE_FILE_APPEND:
		return MODE_APPEND;
	case MD_AST_NODE_TYPE_PIPE:
		return MODE_PIPE;
	case MD_AST_NODE_TYPE_FILE_WRITE:
		return MODE_WRITE;
	default:
		fprintf(stderr, "%s: internal coding error detected in file %s at line %d.\n",
			MLR_GLOBALS.bargv0, __FILE__, __LINE__);
		exit(1);
	}
}

// ----------------------------------------------------------------
void mlr_dsl_cst_statement_free(mlr_dsl_cst_statement_t* pstatement) {

	if (pstatement->subr_callsite_argument_evaluators != NULL) {
		for (int i = 0; i < pstatement->psubr_callsite->arity; i++) {
			rval_evaluator_t* phandler = pstatement->subr_callsite_argument_evaluators[i];
			phandler->pfree_func(phandler);
		}
		free(pstatement->subr_callsite_argument_evaluators);
	}

	if (pstatement->subr_callsite_arguments != NULL) {
		// mv_frees already done in bind_stack_free at call exit
		free(pstatement->subr_callsite_arguments);
	}
	subr_callsite_free(pstatement->psubr_callsite);

	if (pstatement->preturn_evaluator != NULL) {
		pstatement->preturn_evaluator->pfree_func(pstatement->preturn_evaluator);
	}

	if (pstatement->poosvar_lhs_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstatement->poosvar_lhs_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstatement->poosvar_lhs_keylist_evaluators);
	}

	if (pstatement->pemit_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstatement->pemit_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstatement->pemit_keylist_evaluators);
	}

	if (pstatement->ppemit_keylist_evaluators != NULL) {
		for (int i = 0; i < pstatement->num_emit_keylist_evaluators; i++) {
			sllv_t* pemit_keylist_evaluators = pstatement->ppemit_keylist_evaluators[i];
			for (sllve_t* pe = pemit_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
				rval_evaluator_t* phandler = pe->pvvalue;
				phandler->pfree_func(phandler);
			}
			sllv_free(pemit_keylist_evaluators);
		}
		free(pstatement->ppemit_keylist_evaluators);
	}

	if (pstatement->psrec_lhs_evaluator != NULL) {
		pstatement->psrec_lhs_evaluator->pfree_func(pstatement->psrec_lhs_evaluator);
	}

	if (pstatement->prhs_evaluator != NULL) {
		pstatement->prhs_evaluator->pfree_func(pstatement->prhs_evaluator);
	}

	if (pstatement->poutput_filename_evaluator != NULL) {
		pstatement->poutput_filename_evaluator->pfree_func(pstatement->poutput_filename_evaluator);
	}

	if (pstatement->pmulti_out != NULL) {
		multi_out_close(pstatement->pmulti_out);
		multi_out_free(pstatement->pmulti_out);
	}

	if (pstatement->pmulti_lrec_writer != NULL) {
		multi_lrec_writer_drain(pstatement->pmulti_lrec_writer);
		multi_lrec_writer_free(pstatement->pmulti_lrec_writer);
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
			mlr_dsl_cst_statement_free(pe->pvvalue);
		sllv_free(pstatement->pblock_statements);
	}

	if (pstatement->pif_chain_statements != NULL) {
		for (sllve_t* pe = pstatement->pif_chain_statements->phead; pe != NULL; pe = pe->pnext)
			mlr_dsl_cst_statement_free(pe->pvvalue);
		sllv_free(pstatement->pif_chain_statements);
	}

	if (pstatement->pfor_oosvar_k_names != NULL) {
		slls_free(pstatement->pfor_oosvar_k_names);
	}

    if (pstatement->ptriple_for_start_statements != NULL) {
		for (sllve_t* pe = pstatement->ptriple_for_start_statements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			mlr_dsl_cst_statement_free(pstatement);
		}
		sllv_free(pstatement->ptriple_for_start_statements);
	}

	if (pstatement->ptriple_for_continuation_evaluator != NULL) {
		pstatement->ptriple_for_continuation_evaluator->pfree_func(pstatement->ptriple_for_continuation_evaluator);
	}

    if (pstatement->ptriple_for_update_statements != NULL) {
		for (sllve_t* pe = pstatement->ptriple_for_update_statements->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
			mlr_dsl_cst_statement_free(pstatement);
		}
		sllv_free(pstatement->ptriple_for_update_statements);
	}


	if (pstatement->pframe != NULL) {
		bind_stack_frame_free(pstatement->pframe);
	}

	free(pstatement);
}

// ================================================================
static mlr_dsl_cst_statement_vararg_t* mlr_dsl_cst_statement_vararg_alloc(
	char*             emitf_or_unset_srec_field_name,
	rval_evaluator_t* punset_srec_field_name_evaluator,
	rval_evaluator_t* pemitf_arg_evaluator,
	sllv_t*           punset_oosvar_keylist_evaluators)
{
	mlr_dsl_cst_statement_vararg_t* pvararg = mlr_malloc_or_die(sizeof(mlr_dsl_cst_statement_vararg_t));
	pvararg->punset_handler = NULL;
	pvararg->emitf_or_unset_srec_field_name = emitf_or_unset_srec_field_name == NULL
		? NULL : mlr_strdup_or_die(emitf_or_unset_srec_field_name);
	pvararg->punset_oosvar_keylist_evaluators = punset_oosvar_keylist_evaluators;
	pvararg->punset_srec_field_name_evaluator = punset_srec_field_name_evaluator;
	pvararg->pemitf_arg_evaluator             = pemitf_arg_evaluator;

	if (pvararg->punset_oosvar_keylist_evaluators != NULL) {
		pvararg->punset_handler = handle_unset_vararg_oosvar;
	} else if (pvararg->punset_srec_field_name_evaluator != NULL) {
		pvararg->punset_handler = handle_unset_vararg_indirect_srec_field_name;
	} else if (pvararg->emitf_or_unset_srec_field_name != NULL) {
		pvararg->punset_handler = handle_unset_vararg_srec_field_name;
	} else {
		pvararg->punset_handler = handle_unset_vararg_full_srec;
	}

	return pvararg;
}

static void cst_statement_vararg_free(mlr_dsl_cst_statement_vararg_t* pvararg) {
	if (pvararg == NULL)
		return;
	free(pvararg->emitf_or_unset_srec_field_name);

	if (pvararg->punset_srec_field_name_evaluator != NULL) {
		pvararg->punset_srec_field_name_evaluator->pfree_func(pvararg->punset_srec_field_name_evaluator);
	}

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

// ================================================================
// Top-level entry point, e.g. from mapper_put.
void mlr_dsl_cst_handle_top_level_statement_blocks(
	sllv_t*      ptop_level_blocks, // block bodies for begins, main, ends
	variables_t* pvars,
	cst_outputs_t* pcst_outputs)
{
	for (sllve_t* pe = ptop_level_blocks->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_handle_top_level_statement_block(pe->pvvalue, pvars, pcst_outputs);
	}
}

void mlr_dsl_cst_handle_top_level_statement_block(
	cst_top_level_statement_block_t* ptop_level_block,
	variables_t* pvars,
	cst_outputs_t* pcst_outputs)
{
	// xxx check in-use
	local_stack_enter(ptop_level_block->pstack);

	// xxx adapt callee to also handle local stack
	mlr_dsl_cst_handle_statement_list(ptop_level_block->pstatements, pvars, pcst_outputs);

	bind_stack_clear(pvars->pbind_stack); // clear the baseframe // xxx rm
	local_stack_exit(ptop_level_block->pstack);
}

// ================================================================
// xxx copy to ..._for_filter
// xxx rename to ..._for_put

void mlr_dsl_cst_handle_statement_block(
	sllv_t*        pcst_block,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	mlr_dsl_cst_handle_statement_list(pcst_block, pvars, pcst_outputs);
	bind_stack_clear(pvars->pbind_stack); // clear the baseframe
}

// This is for statement lists not recursively contained within a loop body -- including the
// main/begin/end statements.  Since there is no containing loop body, there is no need to check
// for break or continue flags after each statement.
void mlr_dsl_cst_handle_statement_list(
	sllv_t*        pcst_statements,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		pstatement->pnode_handler(pstatement, pvars, pcst_outputs);
		// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
		if (pvars->return_state.returned) {
			break;
		}
	}
}

// This is for statement lists recursively contained within a loop body.
// It checks for break or continue flags after each statement.
static void handle_statement_list_with_break_continue(
	sllv_t*        pcst_statements,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
{
	for (sllve_t* pe = pcst_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pstatement = pe->pvvalue;
		pstatement->pnode_handler(pstatement, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) != 0) {
			break;
		}
		// The UDF/subroutine executor will clear the flag, and consume the retval if there is one.
		if (pvars->return_state.returned) {
			break;
		}
	}
}

// ----------------------------------------------------------------
static void handle_subr_callsite(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for (int i = 0; i < pstatement->psubr_callsite->arity; i++) {
		rval_evaluator_t* pev = pstatement->subr_callsite_argument_evaluators[i];
		pstatement->subr_callsite_arguments[i] = pev->pprocess_func(pev->pvstate, pvars);
	}

	mlr_dsl_cst_execute_subroutine(pstatement->psubr_defsite, pvars, pcst_outputs,
		pstatement->psubr_callsite->arity, pstatement->subr_callsite_arguments);
}

// ----------------------------------------------------------------
static void handle_return_void(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	pvars->return_state.returned = TRUE;
}

// ----------------------------------------------------------------
static void handle_return_value(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	pvars->return_state.retval = pstatement->preturn_evaluator->pprocess_func(
		pstatement->preturn_evaluator->pvstate, pvars);
	pvars->return_state.returned = TRUE;
}

// ----------------------------------------------------------------
static void handle_local_variable_definition(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_present(&val)) {
		bind_stack_define(pvars->pbind_stack, pstatement->local_variable_name, &val, FREE_ENTRY_VALUE);
	} else {
		mv_free(&val);
	}
}

// ----------------------------------------------------------------
static void handle_local_variable_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&val)) {
		bind_stack_set(pvars->pbind_stack, pstatement->local_lhs_variable_name, &val, FREE_ENTRY_VALUE);
	} else {
		mv_free(&val);
	}
}

// ----------------------------------------------------------------
static void handle_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	char* srec_lhs_field_name = pstatement->srec_lhs_field_name;
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

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
	if (mv_is_present(&val)) {
		lhmsmv_put(pvars->ptyped_overlay, srec_lhs_field_name, &val, FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, srec_lhs_field_name, "bug", NO_FREE);
	} else {
		mv_free(&val);
	}
}

// ----------------------------------------------------------------
static void handle_indirect_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* plhs_evaluator = pstatement->psrec_lhs_evaluator;
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	mv_t lval = plhs_evaluator->pprocess_func(plhs_evaluator->pvstate, pvars);
	char free_flags;
	char* srec_lhs_field_name = mv_format_val(&lval, &free_flags);

	mv_t rval = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

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
	if (mv_is_present(&rval)) {
		lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(srec_lhs_field_name), &rval,
			FREE_ENTRY_KEY|FREE_ENTRY_VALUE);
		lrec_put(pvars->pinrec, mlr_strdup_or_die(srec_lhs_field_name), "bug", FREE_ENTRY_KEY | FREE_ENTRY_KEY);
	} else {
		mv_free(&rval);
	}

	if (free_flags)
		free(srec_lhs_field_name);
}

// ----------------------------------------------------------------
static void handle_oosvar_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t rhs_value = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&rhs_value)) {
		int all_non_null_or_error = TRUE;
		sllmv_t* pmvkeys = evaluate_list(pstatement->poosvar_lhs_keylist_evaluators, pvars,
			&all_non_null_or_error);
		if (all_non_null_or_error)
			mlhmmv_put_terminal(pvars->poosvars, pmvkeys, &rhs_value);
		sllmv_free(pmvkeys);
	}
	mv_free(&rhs_value);
}

// ----------------------------------------------------------------
// All assignments produce a mlrval on the RHS and store it on the left -- except if both LHS and RHS
// are oosvars in which case there are recursive copies, or in case of $* on the LHS or RHS.

static void handle_oosvar_to_oosvar_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	int lhs_all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pstatement->poosvar_lhs_keylist_evaluators, pvars,
		&lhs_all_non_null_or_error);

	if (lhs_all_non_null_or_error) {
		int rhs_all_non_null_or_error = TRUE;
		sllmv_t* prhskeys = evaluate_list(pstatement->poosvar_rhs_keylist_evaluators, pvars,
			&rhs_all_non_null_or_error);
		if (rhs_all_non_null_or_error) {
			mlhmmv_copy(pvars->poosvars, plhskeys, prhskeys);
		}
		sllmv_free(prhskeys);
	}

	sllmv_free(plhskeys);
}

static void handle_oosvar_from_full_srec_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* plhskeys = evaluate_list(pstatement->poosvar_lhs_keylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {

		mlhmmv_level_t* plevel = mlhmmv_get_or_create_level(pvars->poosvars, plhskeys);
		if (plevel != NULL) {

			mlhmmv_clear_level(plevel);

			for (lrece_t* pe = pvars->pinrec->phead; pe != NULL; pe = pe->pnext) {
				mv_t k = mv_from_string(pe->key, NO_FREE); // mlhmmv_put_terminal_from_level will copy
				sllmve_t e = { .value = k, .free_flags = 0, .pnext = NULL };
				mv_t* pomv = lhmsmv_get(pvars->ptyped_overlay, pe->key);
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

static void handle_full_srec_from_oosvar_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	lrec_clear(pvars->pinrec);
	lhmsmv_clear(pvars->ptyped_overlay);

	int all_non_null_or_error = TRUE;
	sllmv_t* prhskeys = evaluate_list(pstatement->poosvar_rhs_keylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		int error = 0;
		mlhmmv_level_t* plevel = mlhmmv_get_level(pvars->poosvars, prhskeys, &error);
		if (plevel != NULL) {
			for (mlhmmv_level_entry_t* pentry = plevel->phead; pentry != NULL; pentry = pentry->pnext) {
				if (pentry->level_value.is_terminal) {
					char* skey = mv_alloc_format_val(&pentry->level_key);
					mv_t val = mv_copy(&pentry->level_value.u.mlrval);

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
					lhmsmv_put(pvars->ptyped_overlay, mlr_strdup_or_die(skey), &val,
						FREE_ENTRY_KEY | FREE_ENTRY_VALUE);
					lrec_put(pvars->pinrec, skey, "bug", FREE_ENTRY_KEY);
				}
			}
		}
	}
	sllmv_free(prhskeys);
}

// ----------------------------------------------------------------
static void handle_env_assignment(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	char* env_lhs_name = pstatement->env_lhs_name;
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);

	if (mv_is_present(&val)) {
		setenv(env_lhs_name, mlr_strdup_or_die(mv_alloc_format_val(&val)), 1);
	}
	mv_free(&val);
}

// ----------------------------------------------------------------
static void handle_unset(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for (sllve_t* pf = pstatement->pvarargs->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_vararg_t* pvararg = pf->pvvalue;
		pvararg->punset_handler(pvararg, pvars, pcst_outputs);
	}
}

static void handle_unset_vararg_oosvar(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pvararg->punset_oosvar_keylist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error)
		mlhmmv_remove(pvars->poosvars, pmvkeys);
	sllmv_free(pmvkeys);
}

static void handle_unset_vararg_full_srec(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs)
{
	lrec_clear(pvars->pinrec);
}

static void handle_unset_vararg_srec_field_name(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs)
{
	lrec_remove(pvars->pinrec, pvararg->emitf_or_unset_srec_field_name);
}

static void handle_unset_vararg_indirect_srec_field_name(
	mlr_dsl_cst_statement_vararg_t* pvararg,
	variables_t*                    pvars,
	cst_outputs_t*                  pcst_outputs)
{
	rval_evaluator_t* pevaluator = pvararg->punset_srec_field_name_evaluator;
	mv_t nameval = pevaluator->pprocess_func(pevaluator->pvstate, pvars);
	char free_flags = NO_FREE;
	char* field_name = mv_maybe_alloc_format_val(&nameval, &free_flags);
	lrec_remove(pvars->pinrec, field_name);
	if (free_flags & FREE_ENTRY_VALUE)
		free(field_name);
	mv_free(&nameval);
}

// ----------------------------------------------------------------
static void handle_unset_all(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	sllmv_t* pempty = sllmv_alloc();
	mlhmmv_remove(pvars->poosvars, pempty);
	sllmv_free(pempty);
}

// ----------------------------------------------------------------
static void handle_filter(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		*pcst_outputs->pshould_emit_rec = val.u.boolv;
	} else {
		*pcst_outputs->pshould_emit_rec = FALSE;
	}
}

// ----------------------------------------------------------------
static void handle_final_filter(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		*pcst_outputs->pshould_emit_rec = val.u.boolv ^ pstatement->negate_final_filter;
	} else {
		*pcst_outputs->pshould_emit_rec = FALSE;
	}
}

// ----------------------------------------------------------------
static void handle_conditional_block(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_non_null(&val)) {
		mv_set_boolean_strict(&val);
		if (val.u.boolv) {
			pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);
		}
	}
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_if_head(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	for (sllve_t* pe = pstatement->pif_chain_statements->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_cst_statement_t* pitemnode = pe->pvvalue;
		rval_evaluator_t* prhs_evaluator = pitemnode->prhs_evaluator;

		mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pitemnode->pframe));

				pstatement->pblock_handler(pitemnode->pblock_statements, pvars, pcst_outputs);

				bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
				break;
			}
		}
	}
}

// ----------------------------------------------------------------
static void handle_while(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	loop_stack_push(pvars->ploop_stack);
	while (TRUE) {
		mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (val.u.boolv) {
				pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);
				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
					break;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}
			} else {
				break;
			}
		} else {
			break;
		}
	}
	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_do_while(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	loop_stack_push(pvars->ploop_stack);

	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	while (TRUE) {
		pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			// don't skip the boolean test
		}

		mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
		if (mv_is_non_null(&val)) {
			mv_set_boolean_strict(&val);
			if (!val.u.boolv) {
				break;
			}
		} else {
			break;
		}
	}
	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_for_srec(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	loop_stack_push(pvars->ploop_stack);
	// Copy the lrec for the very likely case that it is being updated inside the for-loop.
	lrec_t* pcopyrec = lrec_copy(pvars->pinrec);
	lhmsmv_t* pcopyoverlay = lhmsmv_copy(pvars->ptyped_overlay);

	for (lrece_t* pe = pcopyrec->phead; pe != NULL; pe = pe->pnext) {

		mv_t mvval = pstatement->ptype_infererenced_srec_field_getter(pe, pcopyoverlay);
		mv_t mvkey = mv_from_string_no_free(pe->key);

		bind_stack_set(pvars->pbind_stack, pstatement->for_srec_k_name, &mvkey, FREE_ENTRY_VALUE);
		bind_stack_set(pvars->pbind_stack, pstatement->for_v_name, &mvval, FREE_ENTRY_VALUE);

		pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);
		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}
	}
	lhmsmv_free(pcopyoverlay);
	lrec_free(pcopyrec);
	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_for_oosvar(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	loop_stack_push(pvars->ploop_stack);

	// Evaluate the keylist: e.g. in 'for ((k1, k2), v in @a[3][$4]) { ... }', find the value of $4 for
	// the current record.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* plhskeylist = evaluate_list(pstatement->poosvar_lhs_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {

		// Locate and copy the submap indexed by the keylist. E.g. in 'for ((k1, k2), v in @a[3][$4]) { ... }', the
		// submap is indexed by ["a", 3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.
		mlhmmv_value_t submap = mlhmmv_copy_submap(pvars->poosvars, plhskeylist);

		if (!submap.is_terminal && submap.u.pnext_level != NULL) {
			// Recurse over the for-k-names, e.g. ["k1", "k2"], on each call descending one level
			// deeper into the submap.  Note there must be at least one k-name so we are assuming
			// the for-loop within handle_for_oosvar_aux was gone through once & thus
			// handle_statement_list_with_break_continue was called through there.

			handle_for_oosvar_aux(pstatement, pvars, pcst_outputs, submap,
				pstatement->pfor_oosvar_k_names->phead);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}
		}

		mlhmmv_free_submap(submap);
	}
	sllmv_free(plhskeylist);

	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

static void handle_for_oosvar_aux(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs,
	mlhmmv_value_t           submap,
	sllse_t*                 prest_for_k_names)
{
	if (prest_for_k_names != NULL) { // Keep recursing over remaining k-names

		if (submap.is_terminal) {
			// The submap was too shallow for the user-specified k-names; there are no terminals here.
		} else {
			// Loop over keys at this submap level:
			for (mlhmmv_level_entry_t* pe = submap.u.pnext_level->phead; pe != NULL; pe = pe->pnext) {
				// Bind the k-name to the entry-key mlrval:
				bind_stack_set(pvars->pbind_stack, prest_for_k_names->value, &pe->level_key, NO_FREE);
				// Recurse into the next-level submap:
				handle_for_oosvar_aux(pstatement, pvars, pcst_outputs, pe->level_value, prest_for_k_names->pnext);

				if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
					// Bit cleared in recursive caller
					return;
				} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
					loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
				}

			}
		}

	} else { // End of recursion: k-names have all been used up

		if (!submap.is_terminal) {
			// The submap was too deep for the user-specified k-names; there are no terminals here.
		} else {
			// Bind the v-name to the terminal mlrval:
			bind_stack_set(pvars->pbind_stack, pstatement->for_v_name, &submap.u.mlrval, NO_FREE);
			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);
		}

	}
}

// ----------------------------------------------------------------
static void handle_for_oosvar_key_only(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	loop_stack_push(pvars->ploop_stack);

	// Evaluate the keylist: e.g. in 'for ((k1, k2), v in @a[3][$4]) { ... }', find the value of $4 for
	// the current record.

	int keys_all_non_null_or_error = FALSE;
	sllmv_t* plhskeylist = evaluate_list(pstatement->poosvar_lhs_keylist_evaluators, pvars,
		&keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		// Locate the submap indexed by the keylist and copy its keys. E.g. in 'for (k1 in @a[3][$4]) { ... }', the
		// submap is indexed by ["a", 3, $4].  Copy it for the very likely case that it is being updated inside the
		// for-loop.
		sllv_t* pkeys = mlhmmv_copy_keys_from_submap(pvars->poosvars, plhskeylist);

		for (sllve_t* pe = pkeys->phead; pe != NULL; pe = pe->pnext) {
			// Bind the v-name to the terminal mlrval:
			bind_stack_set(pvars->pbind_stack, pstatement->pfor_oosvar_k_names->phead->value, pe->pvvalue, NO_FREE);

			// Execute the loop-body statements:
			pstatement->pblock_handler(pstatement->pblock_statements, pvars, pcst_outputs);

			if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
				loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			}
			if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
				loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
			}

			mv_free(pe->pvvalue);
		}
		sllv_free(pkeys);
	}
	sllmv_free(plhskeylist);

	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_triple_for(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	bind_stack_push(pvars->pbind_stack, bind_stack_frame_enter(pstatement->pframe));
	loop_stack_push(pvars->ploop_stack);

	// Start statements
	mlr_dsl_cst_handle_statement_list(pstatement->ptriple_for_start_statements, pvars, pcst_outputs);

	while (TRUE) {

		// Continuation statement
		rval_evaluator_t* pev = pstatement->ptriple_for_continuation_evaluator;
		mv_t val = pev->pprocess_func(pev->pvstate, pvars);
		if (mv_is_non_null(&val))
			mv_set_boolean_strict(&val);
		if (!val.u.boolv)
			break;

		// Body statements
		handle_statement_list_with_break_continue(pstatement->pblock_statements, pvars, pcst_outputs);

		if (loop_stack_get(pvars->ploop_stack) & LOOP_BROKEN) {
			loop_stack_clear(pvars->ploop_stack, LOOP_BROKEN);
			break;
		} else if (loop_stack_get(pvars->ploop_stack) & LOOP_CONTINUED) {
			loop_stack_clear(pvars->ploop_stack, LOOP_CONTINUED);
		}

		// Update statements
		mlr_dsl_cst_handle_statement_list(pstatement->ptriple_for_update_statements, pvars, pcst_outputs);
	}

	loop_stack_pop(pvars->ploop_stack);
	bind_stack_frame_exit(bind_stack_pop(pvars->pbind_stack));
}

// ----------------------------------------------------------------
static void handle_break(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	loop_stack_set(pvars->ploop_stack, LOOP_BROKEN);
}

// ----------------------------------------------------------------
static void handle_continue(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	loop_stack_set(pvars->ploop_stack, LOOP_CONTINUED);
}

// ----------------------------------------------------------------
static void handle_bare_boolean(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;

	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	if (mv_is_non_null(&val))
		mv_set_boolean_strict(&val);
}

// ----------------------------------------------------------------
static void handle_tee_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->psingle_lrec_writer == NULL)
		pstatement->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	lrec_t* pcopy = handle_tee_common(pstatement, pvars, pcst_outputs);

	// The writer frees the lrec
	pstatement->psingle_lrec_writer->pprocess_func(pstatement->psingle_lrec_writer->pvstate,
		pstatement->stdfp, pcopy);
	if (pstatement->flush_every_record)
		fflush(pstatement->stdfp);
}

static void handle_tee_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->pmulti_lrec_writer == NULL)
		pstatement->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

	lrec_t* pcopy = handle_tee_common(pstatement, pvars, pcst_outputs);

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	// The writer frees the lrec
	multi_lrec_writer_output_srec(pstatement->pmulti_lrec_writer, pcopy, filename,
		pstatement->file_output_mode, pstatement->flush_every_record);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

static lrec_t* handle_tee_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	lrec_t* pcopy = lrec_copy(pvars->pinrec);

	// Write the output fields from the typed overlay back to the lrec.
	for (lhmsmve_t* pe = pvars->ptyped_overlay->phead; pe != NULL; pe = pe->pnext) {
		char* output_field_name = pe->key;
		mv_t* pval = &pe->value;

		// Ownership transfer from mv_t to lrec.
		if (pval->type == MT_STRING || pval->type == MT_EMPTY) {
			lrec_put(pcopy, output_field_name, mlr_strdup_or_die(pval->u.strv), FREE_ENTRY_VALUE);
		} else {
			char free_flags = NO_FREE;
			char* string = mv_format_val(pval, &free_flags);
			lrec_put(pcopy, output_field_name, string, free_flags);
		}
	}
	return pcopy;
}

// ----------------------------------------------------------------
static void handle_emitf(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	handle_emitf_common(pstatement, pvars, pcst_outputs->poutrecs);
}

static void handle_emitf_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->psingle_lrec_writer == NULL)
		pstatement->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emitf_common(pstatement, pvars, poutrecs);

	lrec_writer_print_all(pstatement->psingle_lrec_writer, pstatement->stdfp, poutrecs);
	if (pstatement->flush_every_record)
		fflush(pstatement->stdfp);
	sllv_free(poutrecs);
}

static void handle_emitf_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->pmulti_lrec_writer == NULL)
		pstatement->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

	sllv_t* poutrecs = sllv_alloc();

	handle_emitf_common(pstatement, pvars, poutrecs);

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	multi_lrec_writer_output_list(pstatement->pmulti_lrec_writer, poutrecs, filename,
		pstatement->file_output_mode, pstatement->flush_every_record);

	sllv_free(poutrecs);
	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

static void handle_emitf_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  poutrecs)
{
	lrec_t* prec_to_emit = lrec_unbacked_alloc();
	for (sllve_t* pf = pstatement->pvarargs->phead; pf != NULL; pf = pf->pnext) {
		mlr_dsl_cst_statement_vararg_t* pvararg = pf->pvvalue;
		char* emitf_or_unset_srec_field_name = pvararg->emitf_or_unset_srec_field_name;
		rval_evaluator_t* pemitf_arg_evaluator = pvararg->pemitf_arg_evaluator;

		// This is overkill ... the grammar allows only for oosvar names as args to emit.  So we could bypass
		// that and just hashmap-get keyed by emitf_or_unset_srec_field_name here.
		mv_t val = pemitf_arg_evaluator->pprocess_func(pemitf_arg_evaluator->pvstate, pvars);

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
static void handle_emit(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	handle_emit_common(pstatement, pvars, pcst_outputs->poutrecs, pcst_outputs->oosvar_flatten_separator);
}

static void handle_emit_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	sllv_t* poutrecs = sllv_alloc();

	handle_emit_common(pstatement, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->psingle_lrec_writer == NULL)
		pstatement->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	lrec_writer_print_all(pstatement->psingle_lrec_writer, pstatement->stdfp, poutrecs);
	if (pstatement->flush_every_record)
		fflush(pstatement->stdfp);

	sllv_free(poutrecs);
}

static void handle_emit_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->pmulti_lrec_writer == NULL)
		pstatement->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emit_common(pstatement, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	multi_lrec_writer_output_list(pstatement->pmulti_lrec_writer, poutrecs, filename,
		pstatement->file_output_mode, pstatement->flush_every_record);
	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

// ----------------------------------------------------------------
static void handle_emit_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  poutrecs,
	char*                    oosvar_flatten_separator)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstatement->pemit_keylist_evaluators, pvars, &keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pstatement->pemit_oosvar_namelist_evaluators, pvars,
			&names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mlhmmv_to_lrecs(pvars->poosvars, pmvkeys, pmvnames, poutrecs,
				pstatement->do_full_prefixing, oosvar_flatten_separator);
		}
		sllmv_free(pmvnames);
	}
	sllmv_free(pmvkeys);
}

// ----------------------------------------------------------------
static void handle_emit_lashed(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	handle_emit_lashed_common(pstatement, pvars, pcst_outputs->poutrecs, pcst_outputs->oosvar_flatten_separator);
}

static void handle_emit_lashed_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->psingle_lrec_writer == NULL)
		pstatement->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emit_lashed_common(pstatement, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	lrec_writer_print_all(pstatement->psingle_lrec_writer, pstatement->stdfp, poutrecs);
	if (pstatement->flush_every_record)
		fflush(pstatement->stdfp);

	sllv_free(poutrecs);
}

static void handle_emit_lashed_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->pmulti_lrec_writer == NULL)
		pstatement->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emit_lashed_common(pstatement, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	multi_lrec_writer_output_list(pstatement->pmulti_lrec_writer, poutrecs, filename,
		pstatement->file_output_mode, pstatement->flush_every_record);

	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

static void handle_emit_lashed_common(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	sllv_t*                  poutrecs,
	char*                    oosvar_flatten_separator)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t** ppmvkeys = evaluate_lists(pstatement->ppemit_keylist_evaluators, pstatement->num_emit_keylist_evaluators,
		pvars, &keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pstatement->pemit_oosvar_namelist_evaluators, pvars,
			&names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mlhmmv_to_lrecs_lashed(pvars->poosvars, ppmvkeys, pstatement->num_emit_keylist_evaluators, pmvnames,
				poutrecs, pstatement->do_full_prefixing, oosvar_flatten_separator);
		}
		sllmv_free(pmvnames);
	}
	for (int i = 0; i < pstatement->num_emit_keylist_evaluators; i++) {
		sllmv_free(ppmvkeys[i]);
	}
	free(ppmvkeys);
}

// ----------------------------------------------------------------
static void handle_emit_all(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstatement->pemit_oosvar_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_all_to_lrecs(pvars->poosvars, pmvnames, pcst_outputs->poutrecs,
			pstatement->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);
}

static void handle_emit_all_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->psingle_lrec_writer == NULL)
		pstatement->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstatement->pemit_oosvar_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_all_to_lrecs(pvars->poosvars, pmvnames, poutrecs,
			pstatement->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);

	lrec_writer_print_all(pstatement->psingle_lrec_writer, pstatement->stdfp, poutrecs);
	if (pstatement->flush_every_record)
		fflush(pstatement->stdfp);
	sllv_free(poutrecs);
}

static void handle_emit_all_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstatement->pmulti_lrec_writer == NULL)
		pstatement->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();
	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstatement->pemit_oosvar_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_all_to_lrecs(pvars->poosvars, pmvnames, poutrecs,
			pstatement->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	multi_lrec_writer_output_list(pstatement->pmulti_lrec_writer, poutrecs, filename,
		pstatement->file_output_mode, pstatement->flush_every_record);
	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
	sllmv_free(pmvnames);
}

// ----------------------------------------------------------------
static void handle_dump(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	mlhmmv_print_json_stacked(pvars->poosvars, FALSE, "", pstatement->stdfp);
}

static void handle_dump_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	FILE* outfp = multi_out_get(pstatement->pmulti_out, filename, pstatement->file_output_mode);
	mlhmmv_print_json_stacked(pvars->poosvars, FALSE, "", outfp);
	if (pstatement->flush_every_record)
		fflush(outfp);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

// ----------------------------------------------------------------
static void handle_print(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	rval_evaluator_t* prhs_evaluator = pstatement->prhs_evaluator;
	mv_t val = prhs_evaluator->pprocess_func(prhs_evaluator->pvstate, pvars);
	char sfree_flags;
	char* sval = mv_format_val(&val, &sfree_flags);

	rval_evaluator_t* poutput_filename_evaluator = pstatement->poutput_filename_evaluator;
	if (poutput_filename_evaluator == NULL) {
		fprintf(pstatement->stdfp, "%s%s", sval, pstatement->print_terminator);
	} else {
		mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

		char fn_free_flags;
		char* filename = mv_format_val(&filename_mv, &fn_free_flags);

		FILE* outfp = multi_out_get(pstatement->pmulti_out, filename, pstatement->file_output_mode);
		fprintf(outfp, "%s%s", sval, pstatement->print_terminator);
		if (pstatement->flush_every_record)
			fflush(outfp);

		if (fn_free_flags)
			free(filename);
		mv_free(&filename_mv);
	}

	if (sfree_flags)
		free(sval);
	mv_free(&val);
}

// xxx split out mlr_dsl_cst_statements.c ?
