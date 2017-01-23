#include "keylist_evaluators.h"
#include "rval_evaluators.h"

// ----------------------------------------------------------------
// Example ASTs, with and without indexing on the left-hand-side oosvar name:

// $ mlr -n put -v '@x[1]["2"][$3][@4]=5'
// AST ROOT:
// text="list", type=statement_list:
//     text="=", type=oosvar_assignment:
//         text="oosvar_keylist", type=oosvar_keylist:
//             text="x", type=string_literal.
//             text="1", type=numeric_literal.
//             text="2", type=numeric_literal.
//             text="3", type=field_name.
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="4", type=string_literal.
//         text="5", type=numeric_literal.
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
//                 text="a", type=numeric_literal.
//                 text="b", type=numeric_literal.
//                 text="c", type=numeric_literal.
//         text="stream", type=stream:
//
// $ mlr -n put -q -v 'emit @v[1][2], "a", "b","c"'
// AST ROOT:
// text="list", type=statement_list:
//     text="emit", type=emit:
//         text="emit", type=emit:
//             text="oosvar_keylist", type=oosvar_keylist:
//                 text="v", type=string_literal.
//                 text="1", type=numeric_literal.
//                 text="2", type=numeric_literal.
//             text="emit_namelist", type=emit:
//                 text="a", type=numeric_literal.
//                 text="b", type=numeric_literal.
//                 text="c", type=numeric_literal.
//         text="stream", type=stream:

// pnode is input; pkeylist_evaluators is appended to.
sllv_t* allocate_keylist_evaluators_from_ast_node(
	mlr_dsl_ast_node_t* pnode, fmgr_t* pfmgr, int type_inferencing, int context_flags)
{
	sllv_t* pkeylist_evaluators = sllv_alloc();

	if (pnode->pchildren != NULL) { // Non-indexed localvars have no child nodes in the AST.
		for (sllve_t* pe = pnode->pchildren->phead; pe != NULL; pe = pe->pnext) {
			mlr_dsl_ast_node_t* pkeynode = pe->pvvalue;
			if (pkeynode->type == MD_AST_NODE_TYPE_STRING_LITERAL) {
				sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_string(pkeynode->text));
			} else {
				sllv_append(pkeylist_evaluators, rval_evaluator_alloc_from_ast(pkeynode, pfmgr,
					type_inferencing, context_flags));
			}
		}
	}

	return pkeylist_evaluators;
}
