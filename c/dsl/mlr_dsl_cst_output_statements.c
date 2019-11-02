#include <stdlib.h>
#include "lib/mlr_globals.h"
#include "lib/mlrutil.h"
#include "keylist_evaluators.h"
#include "mlr_dsl_cst.h"
#include "context_flags.h"

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
		MLR_INTERNAL_CODING_ERROR();
		return MODE_WRITE; // not reached
	}
}

// ================================================================
typedef struct _print_state_t {
	rxval_evaluator_t*   prhs_xevaluator;
	FILE*                stdfp;
	file_output_mode_t   file_output_mode;
	rval_evaluator_t*    poutput_filename_evaluator;
	int                  flush_every_record;
	multi_out_t*         pmulti_out;
	char*                print_terminator;
} print_state_t;

static mlr_dsl_cst_statement_handler_t handle_print;
static mlr_dsl_cst_statement_freer_t free_print;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_print(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	char*               print_terminator)
{
	print_state_t* pstate = mlr_malloc_or_die(sizeof(print_state_t));

	pstate->prhs_xevaluator            = NULL;
	pstate->stdfp                      = NULL;
	pstate->poutput_filename_evaluator = NULL;
	pstate->pmulti_out                 = NULL;

	MLR_INTERNAL_CODING_ERROR_IF((pnode->pchildren == NULL) || (pnode->pchildren->length != 2));
	mlr_dsl_ast_node_t* pvalue_node = pnode->pchildren->phead->pvvalue;
	pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(pvalue_node, pcst->pfmgr,
		type_inferencing, context_flags);
	pstate->print_terminator = print_terminator;

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;
	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) {
		pstate->stdfp = stdout;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstate->stdfp = stderr;
	} else {
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstate->pmulti_out = multi_out_alloc();
	}
	pstate->flush_every_record = pcst->flush_every_record;

	return mlr_dsl_cst_statement_valloc(
		pnode,
		handle_print,
		free_print,
		pstate);
}

// ----------------------------------------------------------------
static void free_print(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	print_state_t* pstate = pstatement->pvstate;

	if (pstate->prhs_xevaluator != NULL) {
		pstate->prhs_xevaluator->pfree_func(pstate->prhs_xevaluator);
	}

	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}

	if (pstate->pmulti_out != NULL) {
		multi_out_close(pstate->pmulti_out);
		multi_out_free(pstate->pmulti_out);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_print(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	print_state_t* pstate = pstatement->pvstate;

	rxval_evaluator_t* prhs_xevaluator = pstate->prhs_xevaluator;
	boxed_xval_t bxval = prhs_xevaluator->pprocess_func(prhs_xevaluator->pvstate, pvars);

	char sfree_flags = NO_FREE;
	char* sval = "{is-a-map}";

	if (bxval.xval.is_terminal) {
		sval = mv_alloc_format_val(&bxval.xval.terminal_mlrval);
		sfree_flags = FREE_ENTRY_VALUE;
	}

	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	if (poutput_filename_evaluator == NULL) {
		fprintf(pstate->stdfp, "%s%s", sval, pstate->print_terminator);
		if (pstate->flush_every_record)
			fflush(pstate->stdfp);
	} else {
		mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

		char fn_free_flags;
		char* filename = mv_format_val(&filename_mv, &fn_free_flags);

		FILE* outfp = multi_out_get(pstate->pmulti_out, filename, pstate->file_output_mode);
		fprintf(outfp, "%s%s", sval, pstate->print_terminator);
		if (pstate->flush_every_record)
			fflush(outfp);

		if (fn_free_flags)
			free(filename);
		mv_free(&filename_mv);
	}

	if (sfree_flags) {
		free(sval);
	}
	if (bxval.is_ephemeral) {
		mlhmmv_xvalue_free(&bxval.xval);
	}
}

// ================================================================
typedef struct _tee_state_t {
	FILE*                stdfp;
	file_output_mode_t   file_output_mode;
	rval_evaluator_t*    poutput_filename_evaluator;
	int                  flush_every_record;
	lrec_writer_t*       psingle_lrec_writer;
	multi_lrec_writer_t* pmulti_lrec_writer;
} tee_state_t;

static mlr_dsl_cst_statement_handler_t handle_tee_to_stdfp;
static mlr_dsl_cst_statement_handler_t handle_tee_to_file;
static mlr_dsl_cst_statement_freer_t free_tee;

static lrec_t* handle_tee_common(
	tee_state_t*   pstate,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs);

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_tee(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	tee_state_t* pstate = mlr_malloc_or_die(sizeof(tee_state_t));

	pstate->stdfp                      = NULL;
	pstate->poutput_filename_evaluator = NULL;
	pstate->psingle_lrec_writer        = NULL;
	pstate->pmulti_lrec_writer         = NULL;

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;

	pstate->flush_every_record = pcst->flush_every_record;
	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstate->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;

		return mlr_dsl_cst_statement_valloc(
			pnode,
			handle_tee_to_stdfp,
			free_tee,
			pstate);

	} else {
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);

		return mlr_dsl_cst_statement_valloc(
			pnode,
			handle_tee_to_file,
			free_tee,
			pstate);
	}
}

// ----------------------------------------------------------------
static void free_tee(mlr_dsl_cst_statement_t* pstatement, context_t* pctx) {
	tee_state_t* pstate = pstatement->pvstate;

	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}

	if (pstate->psingle_lrec_writer != NULL) {
		pstate->psingle_lrec_writer->pfree_func(pstate->psingle_lrec_writer, pctx);
	}

	if (pstate->pmulti_lrec_writer != NULL) {
		multi_lrec_writer_drain(pstate->pmulti_lrec_writer, pctx);
		multi_lrec_writer_free(pstate->pmulti_lrec_writer, pctx);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_tee_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	tee_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->psingle_lrec_writer == NULL)
		pstate->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	lrec_t* pcopy = handle_tee_common(pstate, pvars, pcst_outputs);

	// The writer frees the lrec
	pstate->psingle_lrec_writer->pprocess_func(pstate->psingle_lrec_writer->pvstate,
		pstate->stdfp, pcopy, pvars->pctx);
	if (pstate->flush_every_record)
		fflush(pstate->stdfp);
}

// ----------------------------------------------------------------
static void handle_tee_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	tee_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->pmulti_lrec_writer == NULL)
		pstate->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

	lrec_t* pcopy = handle_tee_common(pstate, pvars, pcst_outputs);

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	// The writer frees the lrec
	multi_lrec_writer_output_srec(pstate->pmulti_lrec_writer, pcopy, filename,
		pstate->file_output_mode, pstate->flush_every_record, pvars->pctx);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

// ----------------------------------------------------------------
static lrec_t* handle_tee_common(
	tee_state_t*   pstate,
	variables_t*   pvars,
	cst_outputs_t* pcst_outputs)
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

// ================================================================
// Most statements have one item, except emit and emitf.
struct _emitf_item_t;
typedef void emitf_item_handler_t(
	struct _emitf_item_t* pemitf_item,
	variables_t*          pvars,
	cst_outputs_t*        pcst_outputs);

typedef struct _emitf_item_t {
	char*                 srec_field_name;
	rval_evaluator_t*     parg_evaluator;
} emitf_item_t;

static emitf_item_t* alloc_emitf_item(char* srec_field_name, rval_evaluator_t* parg_evaluator);
static void free_emitf_item(emitf_item_t* pemitf_item);

// ----------------------------------------------------------------
typedef struct _emitf_state_t {
	FILE*                stdfp;
	file_output_mode_t   file_output_mode;
	rval_evaluator_t*    poutput_filename_evaluator;
	int                  flush_every_record;
	lrec_writer_t*       psingle_lrec_writer;
	multi_lrec_writer_t* pmulti_lrec_writer;
	sllv_t*              pemitf_items;
} emitf_state_t;

static mlr_dsl_cst_statement_handler_t handle_emitf;
static mlr_dsl_cst_statement_freer_t free_emitf;

static void handle_emitf(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs);

static void handle_emitf_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs);

static void handle_emitf_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs);

static void handle_emitf_common(
	emitf_state_t* pstate,
	variables_t*   pvars,
	sllv_t*        poutrecs);

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_emitf(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	emitf_state_t* pstate = mlr_malloc_or_die(sizeof(emitf_state_t));

	pstate->stdfp                      = NULL;
	pstate->poutput_filename_evaluator = NULL;
	pstate->psingle_lrec_writer        = NULL;
	pstate->pmulti_lrec_writer         = NULL;
	pstate->pemitf_items               = NULL;

	mlr_dsl_ast_node_t* pnamesnode = pnode->pchildren->phead->pvvalue;

	// Loop over oosvar names to emit in e.g. 'emitf @a, @b, @c'.
	pstate->pemitf_items = sllv_alloc();
	for (sllve_t* pe = pnamesnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
		mlr_dsl_ast_node_t* pwalker = pe->pvvalue;

		char* name = NULL;
		switch(pwalker->type) {
		case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
			name = ((mlr_dsl_ast_node_t*)(pwalker->pchildren->phead->pvvalue))->text;
			break;
		case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
			name = pwalker->text;
			break;
		case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
			name = pwalker->text;
			break;
		default:
			MLR_INTERNAL_CODING_ERROR();
			break;
		}
		sllv_append(pstate->pemitf_items,
			alloc_emitf_item(
				name,
				rval_evaluator_alloc_from_ast(pwalker, pcst->pfmgr, type_inferencing, context_flags)));
	}

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	mlr_dsl_cst_statement_handler_t* phandler = NULL;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		phandler = handle_emitf;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		pstate->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
		phandler = handle_emitf_to_stdfp;
	} else {
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		phandler = handle_emitf_to_file;
	}
	pstate->flush_every_record = pcst->flush_every_record;

	return mlr_dsl_cst_statement_valloc(
		pnode,
		phandler,
		free_emitf,
		pstate);
}

static emitf_item_t* alloc_emitf_item(char* srec_field_name, rval_evaluator_t* parg_evaluator) {
	emitf_item_t* pemitf_item = mlr_malloc_or_die(sizeof(emitf_item_t));
	pemitf_item->srec_field_name = srec_field_name;
	pemitf_item->parg_evaluator  = parg_evaluator;
	return pemitf_item;
}

static void free_emitf_item(emitf_item_t* pemitf_item) {
	pemitf_item->parg_evaluator->pfree_func(pemitf_item->parg_evaluator);
	free(pemitf_item);
}

static void free_emitf(mlr_dsl_cst_statement_t* pstatement, context_t* pctx) {
	emitf_state_t* pstate = pstatement->pvstate;

	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}

	if (pstate->psingle_lrec_writer != NULL) {
		pstate->psingle_lrec_writer->pfree_func(pstate->psingle_lrec_writer, pctx);
	}

	if (pstate->pmulti_lrec_writer != NULL) {
		multi_lrec_writer_drain(pstate->pmulti_lrec_writer, pctx);
		multi_lrec_writer_free(pstate->pmulti_lrec_writer, pctx);
	}

	if (pstate->pemitf_items != NULL) {
		for (sllve_t* pe = pstate->pemitf_items->phead; pe != NULL; pe = pe->pnext)
			free_emitf_item(pe->pvvalue);
		sllv_free(pstate->pemitf_items);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_emitf(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emitf_state_t* pstate = pstatement->pvstate;
	handle_emitf_common(pstate, pvars, pcst_outputs->poutrecs);
}

static void handle_emitf_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emitf_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->psingle_lrec_writer == NULL)
		pstate->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emitf_common(pstate, pvars, poutrecs);

	lrec_writer_print_all(pstate->psingle_lrec_writer, pstate->stdfp, poutrecs, pvars->pctx);
	if (pstate->flush_every_record)
		fflush(pstate->stdfp);
	sllv_free(poutrecs);
}

static void handle_emitf_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emitf_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->pmulti_lrec_writer == NULL)
		pstate->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);

	sllv_t* poutrecs = sllv_alloc();

	handle_emitf_common(pstate, pvars, poutrecs);

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	multi_lrec_writer_output_list(pstate->pmulti_lrec_writer, poutrecs, filename,
		pstate->file_output_mode, pstate->flush_every_record, pvars->pctx);

	sllv_free(poutrecs);
	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

static void handle_emitf_common(
	emitf_state_t* pstate,
	variables_t*   pvars,
	sllv_t*        poutrecs)
{
	lrec_t* prec_to_emit = lrec_unbacked_alloc();
	for (sllve_t* pf = pstate->pemitf_items->phead; pf != NULL; pf = pf->pnext) {
		emitf_item_t* pemitf_item = pf->pvvalue;
		char* srec_field_name = pemitf_item->srec_field_name;
		rval_evaluator_t* parg_evaluator = pemitf_item->parg_evaluator;

		// This is overkill ... the grammar allows only for oosvar names as args to emit.  So we could bypass
		// that and just hashmap-get keyed by srec_field_name here.
		mv_t val = parg_evaluator->pprocess_func(parg_evaluator->pvstate, pvars);

		if (val.type == MT_STRING) {
			// Ownership transfer from (newly created) mlrval to (newly created) lrec.
			lrec_put(prec_to_emit, srec_field_name, val.u.strv, val.free_flags);
		} else {
			char free_flags = NO_FREE;
			char* string = mv_format_val(&val, &free_flags);
			lrec_put(prec_to_emit, srec_field_name, string, free_flags);
		}

	}
	sllv_append(poutrecs, prec_to_emit);
}

// ================================================================
struct _emit_state_t; // Forward reference
typedef void record_emitter_t(
	struct _emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator);

typedef struct _emit_state_t {
	rval_evaluator_t*  poutput_filename_evaluator;
	FILE*              stdfp;
	file_output_mode_t file_output_mode;
	sllv_t*            pemit_namelist_evaluators;
	int                do_full_prefixing;

	record_emitter_t*  precord_emitter;

	// For map literals
	rxval_evaluator_t* prhs_xevaluator;

	// For local variables
	char* localvar_name;
	int   localvar_frame_relative_index;

	// For oosvars and localvars: indices ["a", 1, $2] in 'for (k,v in @a[1][$2]) {...}'.
	sllv_t* pemit_keylist_evaluators;

	lrec_writer_t* psingle_lrec_writer; // emit/tee to stdout/stderr
	multi_lrec_writer_t* pmulti_lrec_writer; // emit-to-file

	int flush_every_record;
} emit_state_t;

static mlr_dsl_cst_statement_handler_t handle_emit;
static mlr_dsl_cst_statement_handler_t handle_emit_to_stdfp;
static mlr_dsl_cst_statement_handler_t handle_emit_to_file;

static void record_emitter_from_oosvar(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator);

static void record_emitter_from_local_variable(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator);

static void record_emitter_from_full_srec(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator);

static void record_emitter_from_ephemeral_map(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator);

static mlr_dsl_cst_statement_handler_t handle_emit_all;
static mlr_dsl_cst_statement_handler_t handle_emit_all_to_stdfp;
static mlr_dsl_cst_statement_handler_t handle_emit_all_to_file;

static mlr_dsl_cst_statement_freer_t free_emit;

// ----------------------------------------------------------------
// $ mlr -n put -v 'emit @a[2][3], "x", "y", "z"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="a", type=string_literal.
//                 text="2", type=numeric_literal.
//                 text="3", type=numeric_literal.
//             text="emit_namelist", type=emit:
//                 text="x", type=numeric_literal.
//                 text="y", type=numeric_literal.
//                 text="z", type=numeric_literal.
//         text="stream", type=stream:
//
// $ mlr -n put -v 'emit all, "x", "y", "z"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="all", type=all.
//             text="emit_namelist", type=emit:
//                 text="x", type=numeric_literal.
//                 text="y", type=numeric_literal.
//                 text="z", type=numeric_literal.
//         text="stream", type=stream:

mlr_dsl_cst_statement_t* alloc_emit(
	mlr_dsl_cst_t*      pcst,
	mlr_dsl_ast_node_t* pnode,
	int                 type_inferencing,
	int                 context_flags,
	int                 do_full_prefixing)
{
	emit_state_t* pstate = mlr_malloc_or_die(sizeof(emit_state_t));

	pstate->poutput_filename_evaluator    = NULL;
	pstate->stdfp                         = NULL;
	pstate->precord_emitter               = NULL;
	pstate->prhs_xevaluator               = NULL;
	pstate->localvar_name                 = NULL;
	pstate->localvar_frame_relative_index = MD_UNUSED_INDEX;
	pstate->pemit_namelist_evaluators     = NULL;
	pstate->pemit_keylist_evaluators      = NULL;
	pstate->psingle_lrec_writer           = NULL;
	pstate->pmulti_lrec_writer            = NULL;

	mlr_dsl_ast_node_t* pemit_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pkeylist_node = pemit_node->pchildren->phead->pvvalue;

	// Name note: difference between keylist and namelist: in emit @a[$b]["c"], "d", @e,
	// the keylist is ["a", $b, "c"] and the namelist is ["d", @e].

	// xxx why not rxval_evaluator here??

	int output_all = FALSE;
	// The grammar allows only 'emit all', not 'emit @x, all, $y'.
	// So if 'all' appears at all, it's the only name.
	if (pkeylist_node->type == MD_AST_NODE_TYPE_ALL || pkeylist_node->type == MD_AST_NODE_TYPE_FULL_OOSVAR) {
		output_all = TRUE;

		pstate->precord_emitter = record_emitter_from_oosvar;

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_OOSVAR_KEYLIST) {

		pstate->pemit_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
			pkeylist_node, pcst->pfmgr, type_inferencing, context_flags);

		pstate->precord_emitter = record_emitter_from_oosvar;

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE) {
		pstate->precord_emitter = record_emitter_from_local_variable;

		MLR_INTERNAL_CODING_ERROR_IF(pkeylist_node->vardef_frame_relative_index == MD_UNUSED_INDEX);
		pstate->localvar_name = pkeylist_node->text;
		pstate->localvar_frame_relative_index = pkeylist_node->vardef_frame_relative_index;
		pstate->pemit_keylist_evaluators = sllv_alloc();

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE) {
		pstate->precord_emitter = record_emitter_from_local_variable;

		MLR_INTERNAL_CODING_ERROR_IF(pkeylist_node->vardef_frame_relative_index == MD_UNUSED_INDEX);
		pstate->localvar_name = pkeylist_node->text;
		pstate->localvar_frame_relative_index = pkeylist_node->vardef_frame_relative_index;
		pstate->pemit_keylist_evaluators = allocate_keylist_evaluators_from_ast_node(
			pkeylist_node, pcst->pfmgr, type_inferencing, context_flags);

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_FULL_SREC) {
		pstate->precord_emitter = record_emitter_from_full_srec;

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_FUNCTION_CALLSITE) {
		pstate->precord_emitter = record_emitter_from_ephemeral_map;
		pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(
			pkeylist_node, pcst->pfmgr, type_inferencing, context_flags);

	// xxx indexed function callsite -- ?

	} else if (pkeylist_node->type == MD_AST_NODE_TYPE_MAP_LITERAL) {
		pstate->precord_emitter = record_emitter_from_ephemeral_map;
		pstate->prhs_xevaluator = rxval_evaluator_alloc_from_ast(
			pkeylist_node, pcst->pfmgr, type_inferencing, context_flags);

	} else {
		MLR_INTERNAL_CODING_ERROR();
	}

	pstate->pemit_namelist_evaluators = sllv_alloc();
	if (pemit_node->pchildren->length == 2) {
		mlr_dsl_ast_node_t* pnamelist_node = pemit_node->pchildren->phead->pnext->pvvalue;
		for (sllve_t* pe = pnamelist_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			sllv_append(pstate->pemit_namelist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr, type_inferencing, context_flags));
		}
	}

	mlr_dsl_cst_statement_handler_t* phandler = NULL;

	pstate->do_full_prefixing = do_full_prefixing;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		phandler = output_all ? handle_emit_all : handle_emit;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		phandler = output_all ? handle_emit_all_to_stdfp : handle_emit_to_stdfp;
		pstate->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		phandler = output_all ? handle_emit_all_to_file : handle_emit_to_file;
	}
	pstate->flush_every_record = pcst->flush_every_record;

	return mlr_dsl_cst_statement_valloc(
		pnode,
		phandler,
		free_emit,
		pstate);
}

// ----------------------------------------------------------------
static void free_emit(mlr_dsl_cst_statement_t* pstatement, context_t* pctx) {
	emit_state_t* pstate = pstatement->pvstate;

	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}

	if (pstate->prhs_xevaluator != NULL) {
		pstate->prhs_xevaluator->pfree_func(pstate->prhs_xevaluator);
	}

	if (pstate->pemit_namelist_evaluators != NULL) {
		for (sllve_t* pe = pstate->pemit_namelist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstate->pemit_namelist_evaluators);
	}

	if (pstate->pemit_keylist_evaluators != NULL) {
		for (sllve_t* pe = pstate->pemit_keylist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstate->pemit_keylist_evaluators);
	}

	if (pstate->psingle_lrec_writer != NULL) {
		pstate->psingle_lrec_writer->pfree_func(pstate->psingle_lrec_writer, pctx);
	}

	if (pstate->pmulti_lrec_writer != NULL) {
		multi_lrec_writer_drain(pstate->pmulti_lrec_writer, pctx);
		multi_lrec_writer_free(pstate->pmulti_lrec_writer, pctx);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_emit(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;
	pstate->precord_emitter(pstate, pvars, pcst_outputs->poutrecs, pcst_outputs->oosvar_flatten_separator);
}

// ----------------------------------------------------------------
static void handle_emit_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;
	sllv_t* poutrecs = sllv_alloc();

	pstate->precord_emitter(pstate, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->psingle_lrec_writer == NULL)
		pstate->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	lrec_writer_print_all(pstate->psingle_lrec_writer, pstate->stdfp, poutrecs, pvars->pctx);
	if (pstate->flush_every_record)
		fflush(pstate->stdfp);

	sllv_free(poutrecs);
}

// ----------------------------------------------------------------
static void handle_emit_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->pmulti_lrec_writer == NULL)
		pstate->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	pstate->precord_emitter(pstate, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	multi_lrec_writer_output_list(pstate->pmulti_lrec_writer, poutrecs, filename,
		pstate->file_output_mode, pstate->flush_every_record, pvars->pctx);
	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

// ----------------------------------------------------------------
static void record_emitter_from_oosvar(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pemit_keylist_evaluators, pvars, &keys_all_non_null_or_error);
	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars,
			&names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mlhmmv_root_partial_to_lrecs(pvars->poosvars, pmvkeys, pmvnames, poutrecs,
				pstate->do_full_prefixing, oosvar_flatten_separator);
		}
		sllmv_free(pmvnames);
	}
	sllmv_free(pmvkeys);
}

static void record_emitter_from_local_variable(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator)
{
	int keys_all_non_null_or_error = TRUE;
	sllmv_t* pmvkeys = evaluate_list(pstate->pemit_keylist_evaluators, pvars, &keys_all_non_null_or_error);

	mv_t name = mv_from_string(pstate->localvar_name, NO_FREE);
	sllmv_prepend_no_free(pmvkeys, &name);

	if (keys_all_non_null_or_error) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars,
			&names_all_non_null_or_error);
		if (names_all_non_null_or_error) {

			local_stack_frame_t* pframe = local_stack_get_top_frame(pvars->plocal_stack);
			mlhmmv_xvalue_t* pmval = local_stack_frame_ref_extended_from_indexed(pframe,
				pstate->localvar_frame_relative_index, NULL);
			if (pmval != NULL) {
				// Temporarily wrap the localvar in a parent map whose single key is the variable name.
				mlhmmv_root_t* pmap = mlhmmv_wrap_name_and_xvalue(&name, pmval);

				mlhmmv_root_partial_to_lrecs(pmap, pmvkeys, pmvnames, poutrecs,
					pstate->do_full_prefixing, oosvar_flatten_separator);

				mlhmmv_unwrap_name_and_xvalue(pmap);
			}
		}
		sllmv_free(pmvnames);
	}
	sllmv_free(pmvkeys);
}

static void record_emitter_from_full_srec(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator)
{
	sllv_append(poutrecs, lrec_copy(pvars->pinrec));
}

static void record_emitter_from_ephemeral_map(
	emit_state_t* pstate,
	variables_t*  pvars,
	sllv_t*       poutrecs,
	char*         oosvar_flatten_separator)
{
	rxval_evaluator_t* prhs_xevaluator = pstate->prhs_xevaluator;
	boxed_xval_t boxed_xval = prhs_xevaluator->pprocess_func(prhs_xevaluator->pvstate, pvars);
	sllmv_t* pmvkeys = sllmv_alloc();

	if (!boxed_xval.xval.is_terminal) {
		int names_all_non_null_or_error = TRUE;
		sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars,
			&names_all_non_null_or_error);
		if (names_all_non_null_or_error) {
			mv_t name = mv_from_string("_", NO_FREE);
			sllmv_prepend_no_free(pmvkeys, &name);

			// Temporarily wrap the localvar in a parent map whose single key is the variable name.
			mlhmmv_root_t* pmap = mlhmmv_wrap_name_and_xvalue(&name, &boxed_xval.xval);

			mlhmmv_root_partial_to_lrecs(pmap, pmvkeys, pmvnames, poutrecs,
				pstate->do_full_prefixing, oosvar_flatten_separator);

			mlhmmv_unwrap_name_and_xvalue(pmap);
		}
		sllmv_free(pmvnames);
	}

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(&boxed_xval.xval);
	}

	sllmv_free(pmvkeys);
}

// ----------------------------------------------------------------
static void handle_emit_all(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_root_all_to_lrecs(pvars->poosvars, pmvnames, pcst_outputs->poutrecs,
			pstate->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);
}

// ----------------------------------------------------------------
static void handle_emit_all_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->psingle_lrec_writer == NULL)
		pstate->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_root_all_to_lrecs(pvars->poosvars, pmvnames, poutrecs,
			pstate->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}
	sllmv_free(pmvnames);

	lrec_writer_print_all(pstate->psingle_lrec_writer, pstate->stdfp, poutrecs, pvars->pctx);
	if (pstate->flush_every_record)
		fflush(pstate->stdfp);
	sllv_free(poutrecs);
}

// ----------------------------------------------------------------
static void handle_emit_all_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->pmulti_lrec_writer == NULL)
		pstate->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();
	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	int all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars, &all_non_null_or_error);
	if (all_non_null_or_error) {
		mlhmmv_root_all_to_lrecs(pvars->poosvars, pmvnames, poutrecs,
			pstate->do_full_prefixing, pcst_outputs->oosvar_flatten_separator);
	}

	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);
	multi_lrec_writer_output_list(pstate->pmulti_lrec_writer, poutrecs, filename,
		pstate->file_output_mode, pstate->flush_every_record, pvars->pctx);
	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
	sllmv_free(pmvnames);
}

// ================================================================
struct _emit_lashed_item_t; // Forward reference

typedef struct _emit_lashed_item_t {
	rval_evaluator_t*  pbasename_evaluator;
	rxval_evaluator_t* pitem_xevaluator;
} emit_lashed_item_t;

static emit_lashed_item_t* emit_lashed_item_alloc(
	rval_evaluator_t*  pbasename_evaluator,
	rxval_evaluator_t* pitem_xevaluator)
{
	emit_lashed_item_t* pitem = mlr_malloc_or_die(sizeof(emit_lashed_item_t));
	pitem->pbasename_evaluator = pbasename_evaluator;
	pitem->pitem_xevaluator = pitem_xevaluator;
	return pitem;
}

static void emit_lashed_item_free(emit_lashed_item_t* pitem) {
	pitem->pbasename_evaluator->pfree_func(pitem->pbasename_evaluator);
	pitem->pitem_xevaluator->pfree_func(pitem->pitem_xevaluator);
	free(pitem);
}

// ----------------------------------------------------------------
typedef struct _emit_lashed_state_t {
	int                  num_emit_lashed_items;
	emit_lashed_item_t** ppitems;
	sllv_t*              pemit_namelist_evaluators;

	// Used per-call but allocated once in the constructor:
	mv_t*                pbasenames;
	boxed_xval_t*        pboxed_xvals;
	mlhmmv_xvalue_t**    ptop_values;

	rval_evaluator_t*    poutput_filename_evaluator;
	FILE*                stdfp;
	file_output_mode_t   file_output_mode;
	lrec_writer_t*       psingle_lrec_writer; // emit/tee to stdout/stderr
	multi_lrec_writer_t* pmulti_lrec_writer;  // emit-to-file

	int                  do_full_prefixing;
	int                  flush_every_record;
} emit_lashed_state_t;

static mlr_dsl_cst_statement_handler_t handle_emit_lashed;
static mlr_dsl_cst_statement_handler_t handle_emit_lashed_to_stdfp;
static mlr_dsl_cst_statement_handler_t handle_emit_lashed_to_file;
static void handle_emit_lashed_common(
	emit_lashed_state_t* pstate,
	variables_t*         pvars,
	sllv_t*              poutrecs,
	char*                oosvar_flatten_separator);
static void free_emit_lashed(mlr_dsl_cst_statement_t* pstatement, context_t* pctx);

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_emit_lashed(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags, int do_full_prefixing)
{
	emit_lashed_state_t* pstate = mlr_malloc_or_die(sizeof(emit_lashed_state_t));

	pstate->num_emit_lashed_items      = 0;
	pstate->ppitems                    = NULL;
	pstate->pemit_namelist_evaluators  = NULL;

	pstate->poutput_filename_evaluator = NULL;
	pstate->stdfp                      = NULL;
	pstate->psingle_lrec_writer        = NULL;
	pstate->pmulti_lrec_writer         = NULL;

	mlr_dsl_ast_node_t* pemit_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pnext->pvvalue;
	mlr_dsl_ast_node_t* pkeylists_node = pemit_node->pchildren->phead->pvvalue;

	pstate->num_emit_lashed_items = pkeylists_node->pchildren->length;
	pstate->ppitems = mlr_malloc_or_die(pstate->num_emit_lashed_items * sizeof(emit_lashed_item_t*));
	int i = 0;
	for (sllve_t* pe = pkeylists_node->pchildren->phead; pe != NULL; pe = pe->pnext, i++) {
		mlr_dsl_ast_node_t* pkeylist_node = pe->pvvalue;

		rval_evaluator_t* pbasename_evaluator = NULL;
		switch (pkeylist_node->type) {
		case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:
			pbasename_evaluator = rval_evaluator_alloc_from_string(pkeylist_node->text);
			break;
		case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
			pbasename_evaluator = rval_evaluator_alloc_from_string(
				((mlr_dsl_ast_node_t*)pkeylist_node->pchildren->phead->pvvalue)->text);
			break;
		default:
			pbasename_evaluator = rval_evaluator_alloc_from_string("_");
			break;
		}

		rxval_evaluator_t* pitem_xevaluator = rxval_evaluator_alloc_from_ast(
			pkeylist_node, pcst->pfmgr, type_inferencing, context_flags);

		pstate->ppitems[i] = emit_lashed_item_alloc(pbasename_evaluator, pitem_xevaluator);
	}

	sllv_t* pemit_namelist_evaluators = sllv_alloc();
	if (pemit_node->pchildren->length == 2) {
		mlr_dsl_ast_node_t* pnamelist_node = pemit_node->pchildren->phead->pnext->pvvalue;
		for (sllve_t* pe = pnamelist_node->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			sllv_append(pemit_namelist_evaluators,
				rval_evaluator_alloc_from_ast(pkeynode, pcst->pfmgr, type_inferencing, context_flags));
		}
	}
	pstate->pemit_namelist_evaluators = pemit_namelist_evaluators;

	pstate->pbasenames   = mlr_malloc_or_die(pstate->num_emit_lashed_items * sizeof(mv_t));
	pstate->pboxed_xvals = mlr_malloc_or_die(pstate->num_emit_lashed_items * sizeof(boxed_xval_t));
	pstate->ptop_values  = mlr_malloc_or_die(pstate->num_emit_lashed_items * sizeof(mlhmmv_level_entry_t*));

	mlr_dsl_cst_statement_handler_t* phandler = NULL;
	pstate->do_full_prefixing = do_full_prefixing;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren == NULL
		? NULL
		: poutput_node->pchildren->phead == NULL
		? NULL
		: poutput_node->pchildren->phead->pvvalue;
	if (poutput_node->type == MD_AST_NODE_TYPE_STREAM) {
		phandler = handle_emit_lashed;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT || pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		phandler = handle_emit_lashed_to_stdfp;
		pstate->stdfp = (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) ? stdout : stderr;
	} else {
		pstate->poutput_filename_evaluator = NULL;
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		phandler = handle_emit_lashed_to_file;
	}
	pstate->flush_every_record = pcst->flush_every_record;

	return mlr_dsl_cst_statement_valloc(
		pnode,
		phandler,
		free_emit_lashed,
		pstate);
}

// ----------------------------------------------------------------
static void free_emit_lashed(mlr_dsl_cst_statement_t* pstatement, context_t* pctx) {
	emit_lashed_state_t* pstate = pstatement->pvstate;

	if (pstate->ppitems != NULL) {
		for (int i = 0; i < pstate->num_emit_lashed_items; i++) {
			emit_lashed_item_free(pstate->ppitems[i]);
		}
		free(pstate->ppitems);
	}

	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}

	if (pstate->pemit_namelist_evaluators != NULL) {
		for (sllve_t* pe = pstate->pemit_namelist_evaluators->phead; pe != NULL; pe = pe->pnext) {
			rval_evaluator_t* phandler = pe->pvvalue;
			phandler->pfree_func(phandler);
		}
		sllv_free(pstate->pemit_namelist_evaluators);
	}

	free(pstate->pbasenames);
	free(pstate->pboxed_xvals);
	free(pstate->ptop_values);

	if (pstate->psingle_lrec_writer != NULL) {
		pstate->psingle_lrec_writer->pfree_func(pstate->psingle_lrec_writer, pctx);
	}

	if (pstate->pmulti_lrec_writer != NULL) {
		multi_lrec_writer_drain(pstate->pmulti_lrec_writer, pctx);
		multi_lrec_writer_free(pstate->pmulti_lrec_writer, pctx);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_emit_lashed(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_lashed_state_t* pstate = pstatement->pvstate;
	handle_emit_lashed_common(pstate, pvars, pcst_outputs->poutrecs, pcst_outputs->oosvar_flatten_separator);
}

static void handle_emit_lashed_to_stdfp(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_lashed_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->psingle_lrec_writer == NULL)
		pstate->psingle_lrec_writer = lrec_writer_alloc_or_die(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emit_lashed_common(pstate, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	lrec_writer_print_all(pstate->psingle_lrec_writer, pstate->stdfp, poutrecs, pvars->pctx);
	if (pstate->flush_every_record)
		fflush(pstate->stdfp);

	sllv_free(poutrecs);
}

// ----------------------------------------------------------------
static void handle_emit_lashed_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	emit_lashed_state_t* pstate = pstatement->pvstate;

	// The opts aren't complete at alloc time so we need to handle them on first use.
	if (pstate->pmulti_lrec_writer == NULL)
		pstate->pmulti_lrec_writer = multi_lrec_writer_alloc(pcst_outputs->pwriter_opts);

	sllv_t* poutrecs = sllv_alloc();

	handle_emit_lashed_common(pstate, pvars, poutrecs, pcst_outputs->oosvar_flatten_separator);

	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags = 0;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	multi_lrec_writer_output_list(pstate->pmulti_lrec_writer, poutrecs, filename,
		pstate->file_output_mode, pstate->flush_every_record, pvars->pctx);

	sllv_free(poutrecs);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);
}

// ----------------------------------------------------------------
static void handle_emit_lashed_common(
	emit_lashed_state_t* pstate,
	variables_t*         pvars,
	sllv_t*              poutrecs,
	char*                oosvar_flatten_separator)
{
	for (int i = 0; i < pstate->num_emit_lashed_items; i++) {
		rval_evaluator_t* pev = pstate->ppitems[i]->pbasename_evaluator;
		pstate->pbasenames[i] = pev->pprocess_func(pev->pvstate, pvars);
	}

	int names_all_non_null_or_error = TRUE;
	sllmv_t* pmvnames = evaluate_list(pstate->pemit_namelist_evaluators, pvars,
		&names_all_non_null_or_error);
	if (names_all_non_null_or_error) {

		for (int i = 0; i < pstate->num_emit_lashed_items; i++) {
			emit_lashed_item_t* pitem = pstate->ppitems[i];
			pstate->pboxed_xvals[i] = pitem->pitem_xevaluator->pprocess_func(
				pitem->pitem_xevaluator->pvstate, pvars);
			pstate->ptop_values[i] = &pstate->pboxed_xvals[i].xval;
		}

		mlhmmv_xvalues_to_lrecs_lashed(pstate->ptop_values, pstate->num_emit_lashed_items,
			pstate->pbasenames, pmvnames, poutrecs, pstate->do_full_prefixing, oosvar_flatten_separator);

		for (int i = 0; i < pstate->num_emit_lashed_items; i++) {
			if (pstate->pboxed_xvals[i].is_ephemeral) {
				mlhmmv_xvalue_free(&pstate->pboxed_xvals[i].xval);
			}
		}
	}
	sllmv_free(pmvnames);

	for (int i = 0; i < pstate->num_emit_lashed_items; i++) {
		mv_free(&pstate->pbasenames[i]);
	}
}

// ================================================================
struct _dump_state_t; // forward reference

typedef struct _dump_state_t {
	rxval_evaluator_t*    ptarget_xevaluator;

	FILE*                 stdfp;
	file_output_mode_t    file_output_mode;
	rval_evaluator_t*     poutput_filename_evaluator;
	int                   flush_every_record;
	multi_out_t*          pmulti_out;
} dump_state_t;

static mlr_dsl_cst_statement_handler_t handle_dump;
static mlr_dsl_cst_statement_handler_t handle_dump_to_file;
static mlr_dsl_cst_statement_freer_t   free_dump;

// ----------------------------------------------------------------
mlr_dsl_cst_statement_t* alloc_dump(mlr_dsl_cst_t* pcst, mlr_dsl_ast_node_t* pnode,
	int type_inferencing, int context_flags)
{
	dump_state_t* pstate = mlr_malloc_or_die(sizeof(dump_state_t));

	pstate->ptarget_xevaluator         = NULL;
	pstate->stdfp                      = NULL;
	pstate->poutput_filename_evaluator = NULL;
	pstate->pmulti_out                 = NULL;
	pstate->flush_every_record         = pcst->flush_every_record;

	mlr_dsl_ast_node_t* poutput_node = pnode->pchildren->phead->pvvalue;
	mlr_dsl_ast_node_t* pfilename_node = poutput_node->pchildren->phead->pvvalue;
	mlr_dsl_cst_statement_handler_t* phandler = NULL;

	MLR_INTERNAL_CODING_ERROR_UNLESS(pnode->pchildren->length == 2);

	mlr_dsl_ast_node_t* ptarget_node = pnode->pchildren->phead->pnext->pvvalue;

	pstate->ptarget_xevaluator = rxval_evaluator_alloc_from_ast(
		ptarget_node, pcst->pfmgr, type_inferencing, context_flags);

	if (pfilename_node->type == MD_AST_NODE_TYPE_STDOUT) {
		phandler = handle_dump;
		pstate->stdfp = stdout;
	} else if (pfilename_node->type == MD_AST_NODE_TYPE_STDERR) {
		phandler = handle_dump;
		pstate->stdfp = stderr;
	} else {
		pstate->poutput_filename_evaluator = rval_evaluator_alloc_from_ast(pfilename_node, pcst->pfmgr,
			type_inferencing, context_flags);
		pstate->file_output_mode = file_output_mode_from_ast_node_type(poutput_node->type);
		pstate->pmulti_out = multi_out_alloc();
		phandler = handle_dump_to_file;
	}

	return mlr_dsl_cst_statement_valloc(
		pnode,
		phandler,
		free_dump,
		pstate);
}

// ----------------------------------------------------------------
static void free_dump(mlr_dsl_cst_statement_t* pstatement, context_t* _) {
	dump_state_t* pstate = pstatement->pvstate;

	if (pstate->ptarget_xevaluator != NULL) {
		pstate->ptarget_xevaluator->pfree_func(pstate->ptarget_xevaluator);
	}
	if (pstate->poutput_filename_evaluator != NULL) {
		pstate->poutput_filename_evaluator->pfree_func(pstate->poutput_filename_evaluator);
	}
	if (pstate->pmulti_out != NULL) {
		multi_out_close(pstate->pmulti_out);
		multi_out_free(pstate->pmulti_out);
	}

	free(pstate);
}

// ----------------------------------------------------------------
static void handle_dump(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	dump_state_t* pstate = pstatement->pvstate;

	rxval_evaluator_t* ptarget_xevaluator = pstate->ptarget_xevaluator;
	boxed_xval_t boxed_xval = ptarget_xevaluator->pprocess_func(ptarget_xevaluator->pvstate, pvars);

	if (boxed_xval.xval.is_terminal) {
		mlhmmv_print_terminal(&boxed_xval.xval.terminal_mlrval,
			pvars->json_quote_int_keys, pvars->json_quote_non_string_values,
			pstate->stdfp);
		fprintf(pstate->stdfp, "\n");
	} else {
		mlhmmv_level_print_stacked(boxed_xval.xval.pnext_level, 0, FALSE,
			pvars->json_quote_int_keys, pvars->json_quote_non_string_values,
			"", pvars->pctx->auto_line_term,
			pstate->stdfp);
	}

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(&boxed_xval.xval);
	}
}

// ----------------------------------------------------------------
static void handle_dump_to_file(
	mlr_dsl_cst_statement_t* pstatement,
	variables_t*             pvars,
	cst_outputs_t*           pcst_outputs)
{
	dump_state_t* pstate = pstatement->pvstate;
	rval_evaluator_t* poutput_filename_evaluator = pstate->poutput_filename_evaluator;
	mv_t filename_mv = poutput_filename_evaluator->pprocess_func(poutput_filename_evaluator->pvstate, pvars);
	char fn_free_flags;
	char* filename = mv_format_val(&filename_mv, &fn_free_flags);

	FILE* outfp = multi_out_get(pstate->pmulti_out, filename, pstate->file_output_mode);

	rxval_evaluator_t* ptarget_xevaluator = pstate->ptarget_xevaluator;
	boxed_xval_t boxed_xval = ptarget_xevaluator->pprocess_func(ptarget_xevaluator->pvstate, pvars);

	if (boxed_xval.xval.is_terminal) {
		mlhmmv_print_terminal(&boxed_xval.xval.terminal_mlrval,
			pvars->json_quote_int_keys, pvars->json_quote_non_string_values,
			outfp);
		fprintf(outfp, "\n");
	} else if (boxed_xval.xval.pnext_level != NULL) {
		mlhmmv_level_print_stacked(boxed_xval.xval.pnext_level, 0, FALSE,
			pvars->json_quote_int_keys, pvars->json_quote_non_string_values,
			"", pvars->pctx->auto_line_term,
			outfp);
	}

	if (pstate->flush_every_record)
		fflush(outfp);

	if (fn_free_flags)
		free(filename);
	mv_free(&filename_mv);

	if (boxed_xval.is_ephemeral) {
		mlhmmv_xvalue_free(&boxed_xval.xval);
	}
}
