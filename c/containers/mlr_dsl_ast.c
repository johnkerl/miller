#include <string.h>
#include "lib/mlrutil.h"
#include "containers/mlr_dsl_ast.h"

// ----------------------------------------------------------------
mlr_dsl_ast_t* mlr_dsl_ast_alloc() {
	mlr_dsl_ast_t* past = mlr_malloc_or_die(sizeof(mlr_dsl_ast_t));
	past->pbegin_statements = sllv_alloc();
	past->pmain_statements  = sllv_alloc();
	past->pend_statements   = sllv_alloc();
	past->proot             = NULL;
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
	// xxx old grammar
	printf("AST BEGIN STATEMENTS (%llu):\n", past->pbegin_statements->length);
	for (sllve_t* pe = past->pbegin_statements->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_print(pe->pvvalue);

	printf("AST MAIN STATEMENTS (%llu):\n", past->pmain_statements->length);
	for (sllve_t* pe = past->pmain_statements->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_print(pe->pvvalue);

	printf("AST END STATEMENTS (%llu):\n", past->pend_statements->length);
	for (sllve_t* pe = past->pend_statements->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_print(pe->pvvalue);

	// xxx new grammar
	if (past->proot != NULL) {
		printf("AST ROOT:\n");
		mlr_dsl_ast_node_print(past->proot);
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_ast_node_print_aux(mlr_dsl_ast_node_t* pnode, int level) {
	if (pnode == NULL)
		return;
	for (int i = 0; i < level; i++)
		printf("    ");
	printf("%s (%s)%s\n",
		pnode->text,
		mlr_dsl_ast_node_describe_type(pnode->type),
		(pnode->pchildren != NULL) ? ":" : ".");
	if (pnode->pchildren != NULL) {
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_print_aux(pe->pvvalue, level + 1);
		}
	}
}

void mlr_dsl_ast_node_print(mlr_dsl_ast_node_t* pnode) {
	mlr_dsl_ast_node_print_aux(pnode, 0);
}

// ----------------------------------------------------------------
char* mlr_dsl_ast_node_describe_type(mlr_dsl_ast_node_type_t type) {
	switch(type) {
	case MD_AST_NODE_TYPE_STATEMENT_LIST:     return "statement_list";     break;
	case MD_AST_NODE_TYPE_BEGIN:              return "begin";              break;
	case MD_AST_NODE_TYPE_END:                return "end";                break;
	case MD_AST_NODE_TYPE_STRNUM_LITERAL:     return "strnum_literal";     break;
	case MD_AST_NODE_TYPE_BOOLEAN_LITERAL:    return "boolean_literal";    break;
	case MD_AST_NODE_TYPE_REGEXI:             return "regexi";             break;
	case MD_AST_NODE_TYPE_FIELD_NAME:         return "field_name";         break;
	case MD_AST_NODE_TYPE_FULL_SREC:          return "full_srec";          break;
	case MD_AST_NODE_TYPE_OOSVAR_NAME:        return "oosvar_name";        break;
	case MD_AST_NODE_TYPE_OOSVAR_LEVEL_KEY:   return "oosvar_level_key";   break;
	case MD_AST_NODE_TYPE_NON_SIGIL_NAME:     return "non_sigil_name";     break;
	case MD_AST_NODE_TYPE_OPERATOR:           return "operator";           break;
	case MD_AST_NODE_TYPE_SREC_ASSIGNMENT:    return "srec_assignment";    break;
	case MD_AST_NODE_TYPE_OOSVAR_ASSIGNMENT:  return "oosvar_assignment";  break;
	case MD_AST_NODE_TYPE_OOSVAR_FROM_FULL_SREC_ASSIGNMENT:  return "oosvar_from_full_srec_assignment";  break;
	case MD_AST_NODE_TYPE_FULL_SREC_FROM_OOSVAR_ASSIGNMENT:  return "full_srec_from_oosvar_assignment";  break;
	case MD_AST_NODE_TYPE_CONTEXT_VARIABLE:   return "context_variable";   break;
	case MD_AST_NODE_TYPE_STRIPPED_AWAY:      return "stripped_away";      break;
	case MD_AST_NODE_TYPE_CONDITIONAL_BLOCK:  return "conditional_block";  break;
	case MD_AST_NODE_TYPE_FILTER:             return "filter";             break;
	case MD_AST_NODE_TYPE_UNSET:              return "unset";              break;
	case MD_AST_NODE_TYPE_EMITF:              return "emitf";              break;
	case MD_AST_NODE_TYPE_EMITP:              return "emitp";              break;
	case MD_AST_NODE_TYPE_EMIT:               return "emit";               break;
	case MD_AST_NODE_TYPE_DUMP:               return "dump";               break;
	case MD_AST_NODE_TYPE_ALL:                return "all";                break;
	case MD_AST_NODE_TYPE_ENV:                return "env";                break;
	case MD_AST_NODE_TYPE_WHILE:              return "while";              break;
	case MD_AST_NODE_TYPE_FOR_SREC:           return "for-srec";           break;
	case MD_AST_NODE_TYPE_FOR_VARIABLES:      return "for-variables";      break;
	case MD_AST_NODE_TYPE_IN:                 return "in";                 break;
	case MD_AST_NODE_TYPE_BREAK:              return "break";              break;
	case MD_AST_NODE_TYPE_CONTINUE:           return "continue";           break;
	case MD_AST_NODE_TYPE_IFCHAIN:            return "ifchain";            break;
	default: return "???";
	}
}

// ----------------------------------------------------------------
static void mlr_dsl_ast_free_statement_list(sllv_t* plist) {
	for (sllve_t* pe = plist->phead; pe != NULL; pe = pe->pnext)
		mlr_dsl_ast_node_free(pe->pvvalue);
	sllv_free(plist);
}

void mlr_dsl_ast_free(mlr_dsl_ast_t* past) {
	mlr_dsl_ast_free_statement_list(past->pbegin_statements);
	mlr_dsl_ast_free_statement_list(past->pmain_statements);
	mlr_dsl_ast_free_statement_list(past->pend_statements);
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
