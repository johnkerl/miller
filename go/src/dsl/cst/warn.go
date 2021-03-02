// ================================================================
// TODO: comment
// ================================================================

package cst

import (
	"fmt"

	"miller/src/dsl"
	"miller/src/lib"
)

// ----------------------------------------------------------------
func WarnOnAST(
	ast *dsl.AST,
) error {
	variableNamesWrittenTo := make(map[string]bool)
	inAssignment := false

	if ast.RootNode.Children != nil {
		for _, astChild := range ast.RootNode.Children {
			err := warnOnASTAux(
				astChild,
				variableNamesWrittenTo,
				inAssignment,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'z = x + y'
// DSL EXPRESSION:
// z = x + y
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
// AST:
// * statement block
//     * assignment "="
//         * array or map index access "[]"
//             * local variable "z"
//             * local variable "i"
//         * operator "+"
//             * local variable "x"
//             * local variable "y"

// $ mlr -n put -v 'func f(n) { return n+1}'
// DSL EXPRESSION:
// func f(n) { return n+1}
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
//
// Variable name n might not have been assigned yet.

func warnOnASTAux(
	astNode *dsl.ASTNode,
	variableNamesWrittenTo map[string]bool,
	inAssignment bool,
) error {

	if astNode.Type == dsl.NodeTypeLocalVariable {
		// xxx todo: 'z[i] = x + y' where z is set but i is not.
		// z is a write but i is a read :^/

		variableName := string(astNode.Token.Lit)
		if inAssignment {
			variableNamesWrittenTo[variableName] = true
		} else {
			if !variableNamesWrittenTo[variableName] {
				// xxx todo: this would be much more useful with line numbers. :(
				fmt.Printf(
					"Variable name %s might not have been assigned yet.\n",
					variableName,
				)
			}
		}
	} else if astNode.Type == dsl.NodeTypeBeginBlock {
		// TODO: comment why reset
		variableNamesWrittenTo = make(map[string]bool)
	} else if astNode.Type == dsl.NodeTypeEndBlock {
		// TODO: comment why reset
		variableNamesWrittenTo = make(map[string]bool)
	} else if astNode.Type == dsl.NodeTypeFunctionDefinition {
		// TODO: comment why reset
		// TODO: propagate parameter-list back from first node
		var err error = nil
		variableNamesWrittenTo, err = noteParametersForWarnings(astNode)
		if err != nil {
			return err
		}
	} else if astNode.Type == dsl.NodeTypeSubroutineDefinition {
		// TODO: comment why reset
		// TODO: propagate parameter-list back from first node
		var err error = nil
		variableNamesWrittenTo, err = noteParametersForWarnings(astNode)
		if err != nil {
			return err
		}
	}

	//  - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
	// Treewalk

	if astNode.Children != nil {
		for i, astChild := range astNode.Children {
			childInAssignment := inAssignment

			if astNode.Type == dsl.NodeTypeAssignment && i == 0 {
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopOneVariable && i == 0 {
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopTwoVariable && (i == 0 || i == 1) {
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeForLoopMultivariable && (i == 0 || i == 1) {
				childInAssignment = true
			} else if astNode.Type == dsl.NodeTypeParameterList {
				childInAssignment = true
			}
			err := warnOnASTAux(
				astChild,
				variableNamesWrittenTo,
				childInAssignment,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ----------------------------------------------------------------
// $ mlr -n put -v 'func f(n) { return n+1}'
// DSL EXPRESSION:
// func f(n) { return n+1}
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
//
// Variable name n might not have been assigned yet.

func noteParametersForWarnings(
	astNode *dsl.ASTNode,
) (map[string]bool, error) {

	variableNamesWrittenTo := make(map[string]bool)

	lib.InternalCodingErrorIf(
		astNode.Type != dsl.NodeTypeFunctionDefinition &&
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

	return variableNamesWrittenTo, nil
}
