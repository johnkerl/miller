// Shows warnings for things like uninitialized variables. These are things
// that are statically computable from the AST by itself -- it confines itself
// to local-variable analysis. There are other uninitialization issues
// detectable only at runtime, which would benefit from a 'strict mode'.

package cst

import (
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/lib"

	"github.com/johnkerl/pgpg/go/lib/pkg/asts"
)

// Returns true if there are no warnings.
func WarnOnAST(
	ast *asts.AST,
) bool {
	variableNamesWrittenTo := make(map[string]bool)
	inAssignment := false
	ok := true

	for _, astChild := range ast.RootNode.Children {
		ok1 := warnOnASTAux(
			astChild,
			variableNamesWrittenTo,
			inAssignment,
		)
		// Don't end early on first warning; tree-walk to list them all.
		ok = ok1 && ok
	}

	return ok
}

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
	astNode *asts.ASTNode,
	variableNamesWrittenTo map[string]bool,
	inAssignment bool,
) bool {

	ok := true

	// Check local-variable references, and see if they're reads or writes
	// based on the AST parenting of this node.
	if astNode.Type == asts.NodeType(NodeTypeLocalVariable) {
		variableName := tokenLit(astNode)
		if inAssignment {
			variableNamesWrittenTo[variableName] = true
		} else {
			if !variableNamesWrittenTo[variableName] {
				fmt.Fprintf(
					os.Stderr,
					"Variable name %s might not have been assigned yet%s.\n",
					variableName,
					pgpgTokenToLocationInfo(astNode.Token),
				)
				ok = false
			}
		}

	} else if astNode.Type == asts.NodeType(NodeTypeBeginBlock) {
		// Locals are confined to begin/end blocks and func/subr blocks.
		// Reset for this part of the treewalk.
		variableNamesWrittenTo = make(map[string]bool)
	} else if astNode.Type == asts.NodeType(NodeTypeEndBlock) {
		// Locals are confined to begin/end blocks and func/subr blocks.
		// Reset for this part of the treewalk.
		variableNamesWrittenTo = make(map[string]bool)

	} else if astNode.Type == asts.NodeType(NodeTypeNamedFunctionDefinition) {
		// Locals are confined to begin/end blocks and func/subr blocks.  Reset
		// for this part of the treewalk, except mark the parameters as
		// defined.
		variableNamesWrittenTo = noteParametersForWarnings(astNode)
	} else if astNode.Type == asts.NodeType(NodeTypeSubroutineDefinition) {
		// Locals are confined to begin/end blocks and func/subr blocks.  Reset
		// for this part of the treewalk, except mark the parameters as
		// defined.
		variableNamesWrittenTo = noteParametersForWarnings(astNode)
	}

	// Treewalk to check the rest of the AST below this node.

	for i, astChild := range astNode.Children {
		childInAssignment := inAssignment

		if (astNode.Type == asts.NodeType(NodeTypeAssignment) || astNode.Type == asts.NodeType(NodeTypeCompoundAssignment)) && i == 0 {
			// LHS of assignment statements
			childInAssignment = true
		} else if astNode.Type == asts.NodeType(NodeTypeForLoopOneVariable) && i == 0 {
			// The 'k' in 'for (k in $*)'
			childInAssignment = true
		} else if astNode.Type == asts.NodeType(NodeTypeForLoopTwoVariable) && (i == 0 || i == 1) {
			// The 'k' and 'v' in 'for (k,v in $*)'
			childInAssignment = true
		} else if astNode.Type == asts.NodeType(NodeTypeForLoopMultivariable) && (i == 0 || i == 1) {
			// The 'k1', 'k2', and 'v' in 'for ((k1,k2),v in $*)'
			childInAssignment = true
		} else if astNode.Type == asts.NodeType(NodeTypeParameterList) {
			childInAssignment = true
		} else if inAssignment && astNode.Type == asts.NodeType(NodeTypeArrayOrMapIndexAccess) {
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

	return ok
}

// collectParameterNodes flattens ParameterList (which may nest in PGPG grammar) into Parameter nodes.
func collectParameterNodes(parameterListNode *asts.ASTNode) []*asts.ASTNode {
	var out []*asts.ASTNode
	for _, ch := range parameterListNode.Children {
		if ch.Type == asts.NodeType(NodeTypeParameterList) {
			out = append(out, collectParameterNodes(ch)...)
		} else if ch.Type == asts.NodeType(NodeTypeParameter) {
			out = append(out, ch)
		}
	}
	return out
}

// Given a func/subr block, find the names of its parameters.  All the
// lib.InternalCodingErrorIf parts are shape-assertions to make sure this code
// is in sync with the BNF grammar which builds the AST from a Miller-DSL
// source string.
//
// PGPG: Parameter has one child which is LocalVariable (not ParameterName).
func noteParametersForWarnings(
	astNode *asts.ASTNode,
) map[string]bool {

	variableNamesWrittenTo := make(map[string]bool)

	lib.InternalCodingErrorIf(
		astNode.Type != asts.NodeType(NodeTypeNamedFunctionDefinition) &&
			astNode.Type != asts.NodeType(NodeTypeSubroutineDefinition))
	lib.InternalCodingErrorIf(len(astNode.Children) < 1)
	parameterListNode := astNode.Children[0]

	lib.InternalCodingErrorIf(parameterListNode.Type != asts.NodeType(NodeTypeParameterList))

	// PGPG: ParameterList can nest (FuncParams wraps FuncOrSubrParameterList, which builds list recursively).
	for _, ch := range collectParameterNodes(parameterListNode) {
		lib.InternalCodingErrorIf(len(ch.Children) != 1)
		paramNameNode := ch.Children[0]
		lib.InternalCodingErrorIf(paramNameNode.Type != asts.NodeType(NodeTypeLocalVariable))
		variableNamesWrittenTo[tokenLit(paramNameNode)] = true
	}

	return variableNamesWrittenTo
}
