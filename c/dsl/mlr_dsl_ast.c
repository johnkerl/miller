#include <string.h>
#include "lib/mlrutil.h"
#include "dsl/mlr_dsl_ast.h"

// ----------------------------------------------------------------
mlr_dsl_ast_t* mlr_dsl_ast_alloc() {
	mlr_dsl_ast_t* past = mlr_malloc_or_die(sizeof(mlr_dsl_ast_t));
	past->proot = NULL;
	return past;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc(char* text, mlr_dsl_ast_node_type_t type) {
	mlr_dsl_ast_node_t* pnode = (mlr_dsl_ast_node_t*)mlr_malloc_or_die(
		sizeof(mlr_dsl_ast_node_t));

	pnode->text      = mlr_strdup_or_die(text);
	pnode->type      = type;
	pnode->pchildren = NULL;

	pnode->vardef_subframe_relative_index = MD_UNUSED_INDEX;
	pnode->vardef_subframe_index          = MD_UNUSED_INDEX;
	pnode->vardef_frame_relative_index    = MD_UNUSED_INDEX;
	pnode->subframe_var_count             = MD_UNUSED_INDEX;
	pnode->max_subframe_depth             = MD_UNUSED_INDEX;
	pnode->max_var_depth                  = MD_UNUSED_INDEX;

	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother) {
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(pother->text, pother->type);
	return pnode;
}

// ----------------------------------------------------------------
// This is used within the Lemon parser before bind-stack allocation is done.
// It does not copy the indices at each node: text and type.
mlr_dsl_ast_node_t* mlr_dsl_ast_tree_copy(mlr_dsl_ast_node_t* pold) {
	mlr_dsl_ast_node_t* pnew = mlr_dsl_ast_node_copy(pold);
	if (pold->pchildren != NULL) {
		pnew->pchildren = sllv_alloc();
		for (sllve_t* pe = pold->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			sllv_append(pnew->pchildren, mlr_dsl_ast_tree_copy(pchild));
		}
	}
	return pnew;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_zary(char* text, mlr_dsl_ast_node_type_t type)
{
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_unary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa)
{
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_binary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb)
{
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	sllv_append(pnode->pchildren, pb);
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_ternary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc)
{
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	sllv_append(pnode->pchildren, pb);
	sllv_append(pnode->pchildren, pc);
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_alloc_quaternary(char* text, mlr_dsl_ast_node_type_t type,
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb, mlr_dsl_ast_node_t* pc, mlr_dsl_ast_node_t* pd)
{
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	sllv_append(pnode->pchildren, pb);
	sllv_append(pnode->pchildren, pc);
	sllv_append(pnode->pchildren, pd);
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_prepend_arg(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb)
{
	if (pa->pchildren == NULL)
		pa->pchildren = sllv_alloc();
	sllv_push(pa->pchildren, pb);
	return pa;
}

mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb)
{
	if (pa->pchildren == NULL)
		pa->pchildren = sllv_alloc();
	sllv_append(pa->pchildren, pb);
	return pa;
}

mlr_dsl_ast_node_t* mlr_dsl_ast_node_append_arg_to_second_child(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb)
{
	MLR_INTERNAL_CODING_ERROR_IF(pa->pchildren == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pa->pchildren->phead == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pa->pchildren->phead->pnext == NULL);
	MLR_INTERNAL_CODING_ERROR_IF(pa->pchildren->phead->pnext->pvvalue == NULL);

	mlr_dsl_ast_node_append_arg(pa->pchildren->phead->pnext->pvvalue, pb);

	return pa;
}

mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(
	mlr_dsl_ast_node_t* pa, char* name)
{
	free(pa->text);
	pa->text = mlr_strdup_or_die(name);
	return pa;
}

// ----------------------------------------------------------------
void mlr_dsl_ast_node_replace_text(mlr_dsl_ast_node_t* pa, char* text) {
	if (pa->text != NULL) {
		free(pa->text);
	}
	pa->text = mlr_strdup_or_die(text);
}

// ----------------------------------------------------------------
int mlr_dsl_ast_node_type_to_type_mask(mlr_dsl_ast_node_type_t type) {
	switch(type) {

	case MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION:     return TYPE_MASK_ANY;
	case MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION:     return TYPE_MASK_NUMERIC;
	case MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION:         return TYPE_MASK_INT;
	case MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION:       return TYPE_MASK_FLOAT;
	case MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION:     return TYPE_MASK_BOOLEAN;
	case MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION:      return TYPE_MASK_STRING;

	case MD_AST_NODE_TYPE_UNTYPED_PARAMETER_DEFINITION: return TYPE_MASK_ANY;
	case MD_AST_NODE_TYPE_NUMERIC_PARAMETER_DEFINITION: return TYPE_MASK_NUMERIC;
	case MD_AST_NODE_TYPE_INT_PARAMETER_DEFINITION:     return TYPE_MASK_INT;
	case MD_AST_NODE_TYPE_FLOAT_PARAMETER_DEFINITION:   return TYPE_MASK_FLOAT;
	case MD_AST_NODE_TYPE_BOOLEAN_PARAMETER_DEFINITION: return TYPE_MASK_BOOLEAN;
	case MD_AST_NODE_TYPE_STRING_PARAMETER_DEFINITION:  return TYPE_MASK_STRING;
	case MD_AST_NODE_TYPE_MAP_PARAMETER_DEFINITION:     return TYPE_MASK_MAP;

	default: MLR_INTERNAL_CODING_ERROR();               return 0; // not reached
	}
}

// ----------------------------------------------------------------
int mlr_dsl_ast_node_cannot_be_bare_boolean(mlr_dsl_ast_node_t* pnode) {
	switch (pnode->type) {
	case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:
	case MD_AST_NODE_TYPE_FIELD_NAME:
	case MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME:
	case MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME:
	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:
	case MD_AST_NODE_TYPE_NON_SIGIL_NAME:
	case MD_AST_NODE_TYPE_OPERATOR:
	case MD_AST_NODE_TYPE_ENV:
	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:
		return FALSE;
		break;
	default:
		return TRUE;
		break;
	}
}

// ----------------------------------------------------------------
void mlr_dsl_ast_print(mlr_dsl_ast_t* past) {
	printf("AST ROOT:\n");
	if (past->proot == NULL) {
		printf("(null)\n");
	} else {
		mlr_dsl_ast_node_print(past->proot);
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_ast_node_print_aux(mlr_dsl_ast_node_t* pnode, int level, FILE* o) {
	if (pnode == NULL)
		return;
	for (int i = 0; i < level; i++)
		fprintf(o, "    ");
	fprintf(o, "text=\"%s\", type=%s%s",
		pnode->text,
		mlr_dsl_ast_node_describe_type(pnode->type),
		(pnode->pchildren != NULL) ? ":" : ".");

	if (pnode->vardef_subframe_relative_index != MD_UNUSED_INDEX)
		fprintf(o, " vardef_subframe_relative_index=%d", pnode->vardef_subframe_relative_index);
	if (pnode->vardef_subframe_index != MD_UNUSED_INDEX)
		fprintf(o, " vardef_subframe_index=%d", pnode->vardef_subframe_index);
	if (pnode->vardef_frame_relative_index != MD_UNUSED_INDEX)
		fprintf(o, " vardef_frame_relative_index=%d", pnode->vardef_frame_relative_index);
	if (pnode->subframe_var_count != MD_UNUSED_INDEX)
		fprintf(o, " subframe_var_count=%d", pnode->subframe_var_count);
	if (pnode->max_subframe_depth != MD_UNUSED_INDEX)
		fprintf(o, " max_subframe_depth=%d", pnode->max_subframe_depth);
	if (pnode->max_var_depth != MD_UNUSED_INDEX)
		fprintf(o, " max_var_depth=%d", pnode->max_var_depth);

	fprintf(o, "\n");

	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_print_aux(pe->pvvalue, level + 1, o);
		}
	}
}

void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode) {
	mlr_dsl_ast_node_print_aux(pnode, 0, stdout);
}

void mlr_dsl_ast_node_fprint(mlr_dsl_ast_node_t* pnode, FILE* o) {
	mlr_dsl_ast_node_print_aux(pnode, 0, o);
}

// ----------------------------------------------------------------
static void mlr_dsl_ast_node_pretty_fprint_aux(mlr_dsl_ast_node_t* pnode, FILE* o) {
	if (pnode == NULL)
		return;

	if (pnode->pchildren != NULL) {
		fprintf(o, "(");
	}
	if (pnode->type == MD_AST_NODE_TYPE_STRING_LITERAL || pnode->type == MD_AST_NODE_TYPE_REGEXI) {
		fprintf(o, "\"%s\"", pnode->text);
	} else {
		fprintf(o, "%s", pnode->text);
	}

	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			fprintf(o, " ");
			mlr_dsl_ast_node_pretty_fprint_aux(pe->pvvalue, o);
		}
		fprintf(o, ")");
	}
}

void mlr_dsl_ast_node_pretty_fprint(mlr_dsl_ast_node_t* pnode, FILE* o) {
	mlr_dsl_ast_node_pretty_fprint_aux(pnode, o);
	fprintf(o, "\n");
}

// ----------------------------------------------------------------
char* mlr_dsl_ast_node_describe_type(mlr_dsl_ast_node_type_t type) {
	switch(type) {
	case MD_AST_NODE_TYPE_STATEMENT_BLOCK:                       return "STATEMENT_BLOCK";                       break;
	case MD_AST_NODE_TYPE_STATEMENT_LIST:                        return "STATEMENT_LIST";                        break;
	case MD_AST_NODE_TYPE_FUNC_DEF:                              return "FUNC_DEF";                              break;
	case MD_AST_NODE_TYPE_FUNCTION_CALLSITE:                     return "FUNCTION_CALLSITE";                     break;
	case MD_AST_NODE_TYPE_INDEXED_FUNCTION_CALLSITE:             return "INDEXED_FUNCTION_CALLSITE";             break;
	case MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST:           return "MD_AST_NODE_TYPE_INDEXED_FUNCTION_INDEX_LIST"; break;
	case MD_AST_NODE_TYPE_SUBR_DEF:                              return "SUBR_DEF";                              break;
	case MD_AST_NODE_TYPE_SUBR_CALLSITE:                         return "SUBR_CALLSITE";                         break;
	case MD_AST_NODE_TYPE_UNTYPED_LOCAL_DEFINITION:              return "UNTYPED_LOCAL_DEFINITION";              break;
	case MD_AST_NODE_TYPE_NUMERIC_LOCAL_DEFINITION:              return "NUMERIC_LOCAL_DEFINITION";              break;
	case MD_AST_NODE_TYPE_INT_LOCAL_DEFINITION:                  return "INT_LOCAL_DEFINITION";                  break;
	case MD_AST_NODE_TYPE_FLOAT_LOCAL_DEFINITION:                return "FLOAT_LOCAL_DEFINITION";                break;
	case MD_AST_NODE_TYPE_BOOLEAN_LOCAL_DEFINITION:              return "BOOLEAN_LOCAL_DEFINITION";              break;
	case MD_AST_NODE_TYPE_STRING_LOCAL_DEFINITION:               return "STRING_LOCAL_DEFINITION";               break;
	case MD_AST_NODE_TYPE_MAP_LOCAL_DEFINITION:                  return "MAP_LOCAL_DEFINITION";                  break;
	case MD_AST_NODE_TYPE_UNTYPED_PARAMETER_DEFINITION:          return "UNTYPED_PARAMETER_DEFINITION";          break;
	case MD_AST_NODE_TYPE_NUMERIC_PARAMETER_DEFINITION:          return "NUMERIC_PARAMETER_DEFINITION";          break;
	case MD_AST_NODE_TYPE_INT_PARAMETER_DEFINITION:              return "INT_PARAMETER_DEFINITION";              break;
	case MD_AST_NODE_TYPE_FLOAT_PARAMETER_DEFINITION:            return "FLOAT_PARAMETER_DEFINITION";            break;
	case MD_AST_NODE_TYPE_BOOLEAN_PARAMETER_DEFINITION:          return "BOOLEAN_PARAMETER_DEFINITION";          break;
	case MD_AST_NODE_TYPE_STRING_PARAMETER_DEFINITION:           return "STRING_PARAMETER_DEFINITION";           break;
	case MD_AST_NODE_TYPE_MAP_PARAMETER_DEFINITION:              return "MAP_PARAMETER_DEFINITION";              break;
	case MD_AST_NODE_TYPE_RETURN_VALUE:                          return "RETURN_VALUE";                          break;
	case MD_AST_NODE_TYPE_RETURN_VOID:                           return "RETURN_VOID";                           break;
	case MD_AST_NODE_TYPE_BEGIN:                                 return "BEGIN";                                 break;
	case MD_AST_NODE_TYPE_END:                                   return "END";                                   break;
	case MD_AST_NODE_TYPE_STRING_LITERAL:                        return "STRING_LITERAL";                        break;
	case MD_AST_NODE_TYPE_NUMERIC_LITERAL:                       return "NUMERIC_LITERAL";                       break;
	case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:                       return "BOOLEAN_LITERAL";                       break;
	case MD_AST_NODE_TYPE_MAP_LITERAL:                           return "MAP_LITERAL";                           break;
	case MD_AST_NODE_TYPE_MAP_LITERAL_PAIR:                      return "MAP_LITERAL_PAIR";                      break;
	case MD_AST_NODE_TYPE_MAP_LITERAL_KEY:                       return "MAP_LITERAL_KEY";                       break;
	case MD_AST_NODE_TYPE_MAP_LITERAL_VALUE:                     return "MAP_LITERAL_VALUE";                     break;
	case MD_AST_NODE_TYPE_REGEXI:                                return "REGEXI";                                break;
	case MD_AST_NODE_TYPE_FIELD_NAME:                            return "FIELD_NAME";                            break;
	case MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME:                   return "INDIRECT_FIELD_NAME";                   break;
	case MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME:                  return "POSITIONAL_SREC_NAME";                  break;
	case MD_AST_NODE_TYPE_FULL_SREC:                             return "FULL_SREC";                             break;
	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:                        return "OOSVAR_KEYLIST";                        break;
	case MD_AST_NODE_TYPE_FULL_OOSVAR:                           return "FULL_OOSVAR";                           break;
	case MD_AST_NODE_TYPE_NON_SIGIL_NAME:                        return "NON_SIGIL_NAME";                        break;
	case MD_AST_NODE_TYPE_OPERATOR:                              return "OPERATOR";                              break;
	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_ASSIGNMENT:           return "NONINDEXED_LOCAL_ASSIGNMENT";           break;
	case MD_AST_NODE_TYPE_INDEXED_LOCAL_ASSIGNMENT:              return "INDEXED_LOCAL_ASSIGNMENT";              break;
	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:                       return "SREC_ASSIGNMENT";                       break;
	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:              return "INDIRECT_SREC_ASSIGNMENT";              break;
	case MD_AST_NODE_TYPE_POSITIONAL_SREC_NAME_ASSIGNMENT:       return "POSITIONAL_SREC_NAME_ASSIGNMENT";       break;
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:                     return "OOSVAR_ASSIGNMENT";                     break;
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:      return "OOSVAR_FROM_FULL_SREC_ASSIGNMENT";      break;
	case MD_AST_NODE_TYPE_FULL_OOSVAR_ASSIGNMENT:                return "FULL_OOSVAR_ASSIGNMENT";                break;
	case MD_AST_NODE_TYPE_FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT: return "FULL_OOSVAR_FROM_FULL_SREC_ASSIGNMENT"; break;
	case MD_AST_NODE_TYPE_FULL_SREC_ASSIGNMENT:                  return "FULL_SREC_ASSIGNMENT";                  break;
	case MD_AST_NODE_TYPE_ENV_ASSIGNMENT:                        return "ENV_ASSIGNMENT";                        break;
	case MD_AST_NODE_TYPE_CONTEXT_VARIABLE:                      return "CONTEXT_VARIABLE";                      break;
	case MD_AST_NODE_TYPE_STRIPPED_AWAY:                         return "STRIPPED_AWAY";                         break;
	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:                     return "CONDITIONAL_BLOCK";                     break;
	case MD_AST_NODE_TYPE_FILTER:                                return "FILTER";                                break;
	case MD_AST_NODE_TYPE_UNSET:                                 return "UNSET";                                 break;
	case MD_AST_NODE_TYPE_PIPE:                                  return "PIPE";                                  break;
	case MD_AST_NODE_TYPE_FILE_WRITE:                            return "FILE_WRITE";                            break;
	case MD_AST_NODE_TYPE_FILE_APPEND:                           return "FILE_APPEND";                           break;
	case MD_AST_NODE_TYPE_TEE:                                   return "TEE";                                   break;
	case MD_AST_NODE_TYPE_EMITF:                                 return "EMITF";                                 break;
	case MD_AST_NODE_TYPE_EMITP:                                 return "EMITP";                                 break;
	case MD_AST_NODE_TYPE_EMIT:                                  return "EMIT";                                  break;
	case MD_AST_NODE_TYPE_EMITP_LASHED:                          return "EMITP_LASHED";                          break;
	case MD_AST_NODE_TYPE_EMIT_LASHED:                           return "EMIT_LASHED";                           break;
	case MD_AST_NODE_TYPE_DUMP:                                  return "DUMP";                                  break;
	case MD_AST_NODE_TYPE_EDUMP:                                 return "EDUMP";                                 break;
	case MD_AST_NODE_TYPE_PRINT:                                 return "PRINT";                                 break;
	case MD_AST_NODE_TYPE_PRINTN:                                return "PRINTN";                                break;
	case MD_AST_NODE_TYPE_EPRINT:                                return "EPRINT";                                break;
	case MD_AST_NODE_TYPE_EPRINTN:                               return "EPRINTN";                               break;
	case MD_AST_NODE_TYPE_STDOUT:                                return "STDOUT";                                break;
	case MD_AST_NODE_TYPE_STDERR:                                return "STDERR";                                break;
	case MD_AST_NODE_TYPE_STREAM:                                return "STREAM";                                break;
	case MD_AST_NODE_TYPE_ALL:                                   return "ALL";                                   break;
	case MD_AST_NODE_TYPE_ENV:                                   return "ENV";                                   break;
	case MD_AST_NODE_TYPE_WHILE:                                 return "WHILE";                                 break;
	case MD_AST_NODE_TYPE_DO_WHILE:                              return "DO_WHILE";                              break;
	case MD_AST_NODE_TYPE_FOR_SREC:                              return "FOR_SREC";                              break;
	case MD_AST_NODE_TYPE_FOR_SREC_KEY_ONLY:                     return "FOR_SREC_KEY_ONLY";                     break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR:                            return "FOR_OOSVAR";                            break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR_KEY_ONLY:                   return "FOR_OOSVAR_KEY_ONLY";                   break;
	case MD_AST_NODE_TYPE_FOR_LOCAL_MAP:                         return "FOR_LOCAL_MAP";                         break;
	case MD_AST_NODE_TYPE_FOR_LOCAL_MAP_KEY_ONLY:                return "FOR_LOCAL_MAP_KEY_ONLY";                break;
	case MD_AST_NODE_TYPE_FOR_MAP_LITERAL:                       return "FOR_MAP_LITERAL";                       break;
	case MD_AST_NODE_TYPE_FOR_MAP_LITERAL_KEY_ONLY:              return "FOR_MAP_LITERAL_KEY_ONLY";              break;
	case MD_AST_NODE_TYPE_FOR_FUNC_RETVAL:                       return "FOR_FUNC_RETVAL";                       break;
	case MD_AST_NODE_TYPE_FOR_FUNC_RETVAL_KEY_ONLY:              return "FOR_FUNC_RETVAL_KEY_ONLY";              break;
	case MD_AST_NODE_TYPE_FOR_VARIABLES:                         return "FOR_VARIABLES";                         break;
	case MD_AST_NODE_TYPE_TRIPLE_FOR:                            return "TRIPLE_FOR";                            break;
	case MD_AST_NODE_TYPE_NONINDEXED_LOCAL_VARIABLE:             return "NONINDEXED_LOCAL_VARIABLE";             break;
	case MD_AST_NODE_TYPE_INDEXED_LOCAL_VARIABLE:                return "INDEXED_LOCAL_VARIABLE";                break;
	case MD_AST_NODE_TYPE_IN:                                    return "IN";                                    break;
	case MD_AST_NODE_TYPE_BREAK:                                 return "BREAK";                                 break;
	case MD_AST_NODE_TYPE_CONTINUE:                              return "CONTINUE";                              break;
	case MD_AST_NODE_TYPE_IF_HEAD:                               return "IF_HEAD";                               break;
	case MD_AST_NODE_TYPE_IF_ITEM:                               return "IF_ITEM";                               break;

	default: return "UNRECOGNIZED_AST_NODE_TYPE";
	}
}

// ----------------------------------------------------------------
void mlr_dsl_ast_free(mlr_dsl_ast_t* past) {
	mlr_dsl_ast_node_free(past->proot);
	free(past);
}

// ----------------------------------------------------------------
void mlr_dsl_ast_node_free(mlr_dsl_ast_node_t* pnode) {
	if (pnode->pchildren) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pchild = pe->pvvalue;
			mlr_dsl_ast_node_free(pchild);
		}
		sllv_free(pnode->pchildren);
	}
	free(pnode->text);
	free(pnode);
}
