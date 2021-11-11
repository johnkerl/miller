// ================================================================
// Shows warnings for things like uninitialized variables. These are things
// that are statically computable from the AST by itself -- it confines itself
// to local-variable analysis. There are other uninitialization issues
// detectable only at runtime, which would benefit from a 'strict mode'.
// ================================================================

package cst

import (
	"fmt"
	"os"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
// Returns true if there are no warnings.
func WarnOnAST(
	ast *dsl.AST,
) bool {
	variableNamesWrittenTo := make(map[string]bool)
	inAssignment := false
	ok := true

	if ast.RootNode.Children != nil {
		for _, astChild := range ast.RootNode.Children {
			ok1 := warnOnASTAux(
				astChild,
				variableNamesWrittenTo,
				inAssignment,
			)
			// Don't end early on first warning; tree-walk to list them all.
			ok = ok1 && ok
		}
	}

	return ok
}

// ----------------------------------------------------------------
// Example ASTs:
//
// $ mlr -n put -v 'z = x + y'
// DSL EXPRESSION:
// z = x + y
//
// AST:
// * statement block
//     * assignment "="
//         * local variable "z"
//         * operator "+"
//             * local variable "x"
//             * local variable "y"
//
// $ mlr -n put -v 'z[i] = x + y'
// DSL EXPRESSION:
// z[i] = x + y
//
// AST:
// * statement block
//     * assignment "="
//         * array or map index access "[]"
//             * local variable "z"
//             * local variable "i"
//         * operator "+"
//             * local variable "x"
//             * local variable "y"
//
// $ mlr -n put -v 'func f(n) { return n+1}'
// DSL EXPRESSION:
// func f(n) { return n+1}
//
// AST:
// * statement block
//     * function definition "f"
//         * parameter list
//             * parameter
//                 * parameter name "n"
//         * statement block
//             * return "return"
//                 * operator "+"
//                     * local variable "n"
//                     * int literal "1"

// Returns true if there are no warnings
func warnOnASTAux(
	astNode *dsl.ASTNode,
	variableNamesWrittenTo map[string]bool,
	inAssignment bool,
) bool {

	ok := true

	// Check local-variable references, and see if they're reads or writes
	// based on the AST parenting of this node.
	if astNode.Type == dsl.NodeTypeLocalVariable {
		variableName := string(astNode.Token.Lit)
		if inAssignment {
			variableNamesWrittenTo[variableName] = true
		} else {
			if !variableNamesWrittenTo[variableName] {
				// TODO: this would be much more useful with line numbers. :(
				// That would be a big of work with the parser.  Fortunately,
				// Miller is designed around low-keystroke little expressions
				// -- not thousands of lines of Miller-DSL source code -- so
				// people can look at their few lines of Miller-DSL code and
				// spot their error.
				fmt.Fprintf(
					os.Stderr,
					"Variable name %s might not have been assigned yet.\n",
					variableName,
				)
				ok = false
			}
		}

	} else if astNode.Type == dsl.NodeTypeBeginBlock {
		// Locals are confined to begin/end blocks and func/subr blocks.
		// Reset for this part of the treewalk.
		variableNamesWrittenTo = make(map[string]bool)
	} else if astNode.Type == dsl.NodeTypeEndBlock {
		// Locals are confined to begin/end blocks and func/subr blocks.
		// Reset for this part of the treewalk.
		variableNamesWrittenTo = make(map[string]bool)

	} else if astNode.Type == dsl.NodeTypeNamedFunctionDefinition {
		// Locals are confined to begin/end blocks and func/subr blocks.  Reset
		// for this part of the treewalk, except mark the parameters as
		// defined.
		variableNamesWrittenTo = noteParametersForWarnings(astNode)
	} else if astNode.Type == dsl.NodeTypeSubroutineDefinition {
		// Locals are confined to begin/end blocks and func/subr blocks.  Reset
		// for this part of the treewalk, except mark the parameters as
		// defined.
		variableNamesWrittenTo = noteParametersForWarnings(astNode)
	}

	// Treewalk to check the rest of the AST below this node.

	if astNode.Children != nil {
		for i, astChild := range astNode.Children {
			childInAssignment := inAssignment

			if astNode.Type == dsl.NodeTypeAssignment && i == 0 {
				// LHS of assignment statements
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopOneVariable && i == 0 {
				// The 'k' in 'for (k in $*)'
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopTwoVariable && (i == 0 || i == 1) {
				// The 'k' and 'v' in 'for (k,v in $*)'
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopMultivariable && (i == 0 || i == 1) {
				// The 'k1', 'k2', and 'v' in 'for ((k1,k2),v in $*)'
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeParameterList {
				childInAssignment = true
			} else if inAssignment && astNode.Type == dsl.NodeTypeArrayOrMapIndexAccess {
				// In 'z[i] = 1', the 'i' is a read and the 'z' is a write.
				//
				// mlr --from r put -v -W 'z[i] = 1'
				// DSL EXPRESSION:
				// z[i]=1
				//
				// AST:
				// * statement block
				//     * assignment "="
				//         * array or map index access "[]"
				//             * local variable "z"
				//             * local variable "i"
				//         * int literal "1"
				if i == 0 {
					childInAssignment = true
				} else {
					childInAssignment = false
				}
			}
			ok1 := warnOnASTAux(
				astChild,
				variableNamesWrittenTo,
				childInAssignment,
			)
			// Don't end early on first error; tree-walk to list them all.
			ok = ok1 && ok
		}
	}

	return ok
}

// ----------------------------------------------------------------
// Given a func/subr block, find the names of its parameters.  All the
// lib.InternalCodingErrorIf parts are shape-assertions to make sure this code
// is in sync with the BNF grammar which builds the AST from a Miller-DSL
// source string.
func noteParametersForWarnings(
	astNode *dsl.ASTNode,
) map[string]bool {

	variableNamesWrittenTo := make(map[string]bool)

	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeNamedFunctionDefinition &&
			astNode.Type != dsl.NodeTypeSubroutineDefinition)
	lib.InternalCodingErrorIf(len(astNode.Children) < 1)
	parameterListNode := astNode.Children[0]

	lib.InternalCodingErrorIf(parameterListNode.Type != dsl.NodeTypeParameterList)

	for _, parameterNode := range parameterListNode.Children {
		lib.InternalCodingErrorIf(parameterNode.Type != dsl.NodeTypeParameter)
		lib.InternalCodingErrorIf(len(parameterNode.Children) != 1)
		parameterNameNode := parameterNode.Children[0]
		lib.InternalCodingErrorIf(parameterNameNode.Type != dsl.NodeTypeParameterName)
		parameterName := string(parameterNameNode.Token.Lit)
		variableNamesWrittenTo[parameterName] = true
	}

	return variableNamesWrittenTo
}
