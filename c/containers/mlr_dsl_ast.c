#include <string.h>
#include "lib/mlrutil.h"
#include "containers/mlr_dsl_ast.h"

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
	pnode->text = mlr_strdup_or_die(text);
	pnode->type = type;
	pnode->pchildren = NULL;
	return pnode;
}

// ----------------------------------------------------------------
mlr_dsl_ast_node_t* mlr_dsl_ast_node_copy(mlr_dsl_ast_node_t* pother) {
	mlr_dsl_ast_node_t* pnode = mlr_dsl_ast_node_alloc(pother->text, pother->type);
	return pnode;
}

// ----------------------------------------------------------------
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
mlr_dsl_ast_node_t* mlr_dsl_ast_node_prepend_arg(
	mlr_dsl_ast_node_t* pa, mlr_dsl_ast_node_t* pb)
{
	if (pa->pchildren == NULL)
		pa->pchildren = sllv_alloc();
	sllv_prepend(pa->pchildren, pb);
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

mlr_dsl_ast_node_t* mlr_dsl_ast_node_set_function_name(
	mlr_dsl_ast_node_t* pa, char* name)
{
	free(pa->text);
	pa->text = mlr_strdup_or_die(name);
	return pa;
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
	fprintf(o, "text=\"%s\", type=%s%s\n",
		pnode->text,
		mlr_dsl_ast_node_describe_type(pnode->type),
		(pnode->pchildren != NULL) ? ":" : ".");
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
char* mlr_dsl_ast_node_describe_type(mlr_dsl_ast_node_type_t type) {
	switch(type) {
	case MD_AST_NODE_TYPE_STATEMENT_LIST:                   return "statement_list";                   break;
	case MD_AST_NODE_TYPE_BEGIN:                            return "begin";                            break;
	case MD_AST_NODE_TYPE_END:                              return "end";                              break;
	case MD_AST_NODE_TYPE_STRING_LITERAL:                   return "string_literal";                   break;
	case MD_AST_NODE_TYPE_STRNUM_LITERAL:                   return "strnum_literal";                   break;
	case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:                  return "boolean_literal";                  break;
	case MD_AST_NODE_TYPE_REGEXI:                           return "regexi";                           break;
	case MD_AST_NODE_TYPE_FIELD_NAME:                       return "field_name";                       break;
	case MD_AST_NODE_TYPE_INDIRECT_FIELD_NAME:              return "indirect_field_name";              break;
	case MD_AST_NODE_TYPE_FULL_SREC:                        return "full_srec";                        break;
	case MD_AST_NODE_TYPE_OOSVAR_KEYLIST:                   return "oosvar_keylist";                   break;
	case MD_AST_NODE_TYPE_FULL_OOSVAR:                      return "full_oosvar";                      break;
	case MD_AST_NODE_TYPE_NON_SIGIL_NAME:                   return "non_sigil_name";                   break;
	case MD_AST_NODE_TYPE_OPERATOR:                         return "operator";                         break;
	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:                  return "srec_assignment";                  break;
	case MD_AST_NODE_TYPE_INDIRECT_SREC_ASSIGNMENT:         return "indirect_srec_assignment";         break;
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:                return "oosvar_assignment";                break;
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT: return "oosvar_from_full_srec_assignment"; break;
	case MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT: return "full_srec_from_oosvar_assignment"; break;
	case MD_AST_NODE_TYPE_CONTEXT_VARIABLE:                 return "context_variable";                 break;
	case MD_AST_NODE_TYPE_STRIPPED_AWAY:                    return "stripped_away";                    break;
	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:                return "conditional_block";                break;
	case MD_AST_NODE_TYPE_FILTER:                           return "filter";                           break;
	case MD_AST_NODE_TYPE_UNSET:                            return "unset";                            break;
	case MD_AST_NODE_TYPE_EMITF:                            return "emitf";                            break;
	case MD_AST_NODE_TYPE_EMITX:                            return "emitx";                            break;
	case MD_AST_NODE_TYPE_EMITP:                            return "emitp";                            break;
	case MD_AST_NODE_TYPE_EMIT:                             return "emit";                             break;
	case MD_AST_NODE_TYPE_EMITX_LASHED:                     return "emitx_lashed";                     break;
	case MD_AST_NODE_TYPE_EMITP_LASHED:                     return "emitp_lashed";                     break;
	case MD_AST_NODE_TYPE_EMIT_LASHED:                      return "emit_lashed";                      break;
	case MD_AST_NODE_TYPE_DUMP:                             return "dump";                             break;
	case MD_AST_NODE_TYPE_EDUMP:                            return "edump";                            break;
	case MD_AST_NODE_TYPE_PRINT:                            return "print";                            break;
	case MD_AST_NODE_TYPE_EPRINT:                           return "eprint";                           break;
	case MD_AST_NODE_TYPE_ALL:                              return "all";                              break;
	case MD_AST_NODE_TYPE_ENV:                              return "env";                              break;
	case MD_AST_NODE_TYPE_WHILE:                            return "while";                            break;
	case MD_AST_NODE_TYPE_DO_WHILE:                         return "do_while";                         break;
	case MD_AST_NODE_TYPE_FOR_SREC:                         return "for_srec";                         break;
	case MD_AST_NODE_TYPE_FOR_OOSVAR:                       return "for_oosvar";                       break;
	case MD_AST_NODE_TYPE_FOR_VARIABLES:                    return "for_variables";                    break;
	case MD_AST_NODE_TYPE_BOUND_VARIABLE:                   return "bound_variable";                   break;
	case MD_AST_NODE_TYPE_IN:                               return "in";                               break;
	case MD_AST_NODE_TYPE_BREAK:                            return "break";                            break;
	case MD_AST_NODE_TYPE_CONTINUE:                         return "continue";                         break;
	case MD_AST_NODE_TYPE_IF_HEAD:                          return "if_head";                          break;
	case MD_AST_NODE_TYPE_IF_ITEM:                          return "if_item";                          break;
	default: return "???";
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
