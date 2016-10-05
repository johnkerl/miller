#include <string.h>
#include "../lib/mlrutil.h"
#include "ex_ast.h"

// ----------------------------------------------------------------
ex_ast_t* ex_ast_alloc() {
	ex_ast_t* past = mlr_malloc_or_die(sizeof(ex_ast_t));
	past->proot = NULL;
	return past;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_alloc(char* text, ex_ast_node_type_t type) {
	ex_ast_node_t* pnode = (ex_ast_node_t*)mlr_malloc_or_die(
		sizeof(ex_ast_node_t));
	pnode->text = mlr_strdup_or_die(text);
	pnode->type = type;
	pnode->pchildren = NULL;
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_copy(ex_ast_node_t* pother) {
	ex_ast_node_t* pnode = ex_ast_node_alloc(pother->text, pother->type);
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_tree_copy(ex_ast_node_t* pold) {
	ex_ast_node_t* pnew = ex_ast_node_copy(pold);
	if (pold->pchildren != NULL) {
		pnew->pchildren = sllv_alloc();
		for (sllve_t* pe = pold->pchildren->phead; pe != NULL; pe = pe->pnext) {
			ex_ast_node_t* pchild = pe->pvvalue;
			sllv_append(pnew->pchildren, ex_ast_tree_copy(pchild));
		}
	}
	return pnew;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_alloc_zary(char* text, ex_ast_node_type_t type)
{
	ex_ast_node_t* pnode = ex_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_alloc_unary(char* text, ex_ast_node_type_t type,
	ex_ast_node_t* pa)
{
	ex_ast_node_t* pnode = ex_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_alloc_binary(char* text, ex_ast_node_type_t type,
	ex_ast_node_t* pa, ex_ast_node_t* pb)
{
	ex_ast_node_t* pnode = ex_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	sllv_append(pnode->pchildren, pb);
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_alloc_ternary(char* text, ex_ast_node_type_t type,
	ex_ast_node_t* pa, ex_ast_node_t* pb, ex_ast_node_t* pc)
{
	ex_ast_node_t* pnode = ex_ast_node_alloc(text, type);
	pnode->pchildren = sllv_alloc();
	sllv_append(pnode->pchildren, pa);
	sllv_append(pnode->pchildren, pb);
	sllv_append(pnode->pchildren, pc);
	return pnode;
}

// ----------------------------------------------------------------
ex_ast_node_t* ex_ast_node_prepend_arg(
	ex_ast_node_t* pa, ex_ast_node_t* pb)
{
	if (pa->pchildren == NULL)
		pa->pchildren = sllv_alloc();
	sllv_prepend(pa->pchildren, pb);
	return pa;
}

ex_ast_node_t* ex_ast_node_append_arg(
	ex_ast_node_t* pa, ex_ast_node_t* pb)
{
	if (pa->pchildren == NULL)
		pa->pchildren = sllv_alloc();
	sllv_append(pa->pchildren, pb);
	return pa;
}

ex_ast_node_t* ex_ast_node_set_function_name(
	ex_ast_node_t* pa, char* name)
{
	free(pa->text);
	pa->text = mlr_strdup_or_die(name);
	return pa;
}

// ----------------------------------------------------------------
void ex_ast_print(ex_ast_t* past) {
	printf("AST ROOT:\n");
	if (past->proot == NULL) {
		printf("(null)\n");
	} else {
		ex_ast_node_print(past->proot);
	}
}

// ----------------------------------------------------------------
static void ex_ast_node_print_aux(ex_ast_node_t* pnode, int level, FILE* o) {
	if (pnode == NULL)
		return;
	for (int i = 0; i < level; i++)
		fprintf(o, "    ");
	fprintf(o, "text=\"%s\", type=%s%s\n",
		pnode->text,
		ex_ast_node_describe_type(pnode->type),
		(pnode->pchildren != NULL) ? ":" : ".");
	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			ex_ast_node_print_aux(pe->pvvalue, level + 1, o);
		}
	}
}

void ex_ast_node_print(ex_ast_node_t* pnode) {
	ex_ast_node_print_aux(pnode, 0, stdout);
}

void ex_ast_node_fprint(ex_ast_node_t* pnode, FILE* o) {
	ex_ast_node_print_aux(pnode, 0, o);
}

// ----------------------------------------------------------------
char* ex_ast_node_describe_type(ex_ast_node_type_t type) {
	switch(type) {
	case MD_AST_NODE_TYPE_STATEMENT_BLOCK:                   return "statement_list";                   break;
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
	case MD_AST_NODE_TYPE_EMITF_WRITE:                      return "emitf_write";                      break;
	case MD_AST_NODE_TYPE_EMITF_APPEND:                     return "emitf_append";                     break;
	case MD_AST_NODE_TYPE_EMITP:                            return "emitp";                            break;
	case MD_AST_NODE_TYPE_EMITP_WRITE:                      return "emitp_write";                      break;
	case MD_AST_NODE_TYPE_EMITP_APPEND:                     return "emitp_append";                     break;
	case MD_AST_NODE_TYPE_EMIT:                             return "emit";                             break;
	case MD_AST_NODE_TYPE_EMIT_WRITE:                       return "emit_write";                       break;
	case MD_AST_NODE_TYPE_EMIT_APPEND:                      return "emit_append";                      break;
	case MD_AST_NODE_TYPE_EMITP_LASHED:                     return "emitp_lashed";                     break;
	case MD_AST_NODE_TYPE_EMITP_LASHED_WRITE:               return "emitp_lashed_write";               break;
	case MD_AST_NODE_TYPE_EMITP_LASHED_APPEND:              return "emitp_lashed_append";              break;
	case MD_AST_NODE_TYPE_EMIT_LASHED:                      return "emit_lashed";                      break;
	case MD_AST_NODE_TYPE_EMIT_LASHED_WRITE:                return "emit_lashed_write";                break;
	case MD_AST_NODE_TYPE_EMIT_LASHED_APPEND:               return "emit_lashed_append";               break;
	case MD_AST_NODE_TYPE_DUMP:                             return "dump";                             break;
	case MD_AST_NODE_TYPE_DUMP_WRITE:                       return "dump_write";                       break;
	case MD_AST_NODE_TYPE_DUMP_APPEND:                      return "dump_append";                      break;
	case MD_AST_NODE_TYPE_EDUMP:                            return "edump";                            break;
	case MD_AST_NODE_TYPE_PRINT:                            return "print";                            break;
	case MD_AST_NODE_TYPE_PRINT_WRITE:                      return "print_write";                      break;
	case MD_AST_NODE_TYPE_PRINT_APPEND:                     return "print_append";                     break;
	case MD_AST_NODE_TYPE_PRINTN:                           return "printn";                           break;
	case MD_AST_NODE_TYPE_PRINTN_WRITE:                     return "printn_write";                     break;
	case MD_AST_NODE_TYPE_PRINTN_APPEND:                    return "printn_append";                    break;
	case MD_AST_NODE_TYPE_EPRINT:                           return "eprint";                           break;
	case MD_AST_NODE_TYPE_EPRINTN:                          return "eprintn";                          break;
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
void ex_ast_free(ex_ast_t* past) {
	ex_ast_node_free(past->proot);
	free(past);
}

// ----------------------------------------------------------------
void ex_ast_node_free(ex_ast_node_t* pnode) {
	if (pnode->pchildren) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			ex_ast_node_t* pchild = pe->pvvalue;
			ex_ast_node_free(pchild);
		}
		sllv_free(pnode->pchildren);
	}
	free(pnode->text);
	free(pnode);
}
