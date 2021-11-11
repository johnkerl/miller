// ================================================================
// This is for if/elif/elif/else chains.
// ================================================================

package cst

import (
	"errors"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
	"mlr/internal/pkg/types"
)

// ----------------------------------------------------------------
type IfChainNode struct {
	ifItems []*IfItem
}

func NewIfChainNode(ifItems []*IfItem) *IfChainNode {
	return &IfChainNode{
		ifItems: ifItems,
	}
}

// ----------------------------------------------------------------
// For each if/elif/elif/else portion: the conditional part (...) and the
// statement-block part {...}. For "else", the conditional is nil.
type IfItem struct {
	conditionNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

// ----------------------------------------------------------------
// Sample AST:

// DSL EXPRESSION:
// if (NR == 1) { $z = 100 } elif (NR == 2) { $z = 200 } elif (NR == 3) { $z = 300 } else { $z = 900 }
// AST:
// * StatementBlock
//     * IfChain
//         * IfItem "if"
//             * Operator "=="
//                 * ContextVariable "NR"
//                 * IntLiteral "1"
//             * StatementBlock
//                 * Assignment "="
//                     * DirectFieldValue "z"
//                     * IntLiteral "100"
//         * IfItem "elif"
//             * Operator "=="
//                 * ContextVariable "NR"
//                 * IntLiteral "2"
//             * StatementBlock
//                 * Assignment "="
//                     * DirectFieldValue "z"
//                     * IntLiteral "200"
//         * IfItem "elif"
//             * Operator "=="
//                 * ContextVariable "NR"
//                 * IntLiteral "3"
//             * StatementBlock
//                 * Assignment "="
//                     * DirectFieldValue "z"
//                     * IntLiteral "300"
//         * IfItem "else"
//             * StatementBlock
//                 * Assignment "="
//                     * DirectFieldValue "z"
//                     * IntLiteral "900"

func (root *RootNode) BuildIfChainNode(astNode *dsl.ASTNode) (*IfChainNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeIfChain)

	ifItems := make([]*IfItem, 0)

	astChildren := astNode.Children

	for _, astChild := range astChildren {
		lib.InternalCodingErrorIf(astChild.Type != dsl.NodeTypeIfItem)
		token := string(astChild.Token.Lit) // "if", "elif", "else"
		if token == "if" || token == "elif" {
			lib.InternalCodingErrorIf(len(astChild.Children) != 2)
			conditionNode, err := root.BuildEvaluableNode(astChild.Children[0])
			if err != nil {
				return nil, err
			}
			statementBlockNode, err := root.BuildStatementBlockNode(astChild.Children[1])
			if err != nil {
				return nil, err
			}
			ifItem := &IfItem{
				conditionNode:      conditionNode,
				statementBlockNode: statementBlockNode,
			}
			ifItems = append(ifItems, ifItem)

		} else if token == "else" {
			lib.InternalCodingErrorIf(len(astChild.Children) != 1)
			var conditionNode IEvaluable = nil
			statementBlockNode, err := root.BuildStatementBlockNode(astChild.Children[0])
			if err != nil {
				return nil, err
			}
			ifItem := &IfItem{
				conditionNode:      conditionNode,
				statementBlockNode: statementBlockNode,
			}
			ifItems = append(ifItems, ifItem)

		} else {
			lib.InternalCodingErrorIf(true)
		}
	}
	ifChainNode := NewIfChainNode(ifItems)
	return ifChainNode, nil
}

// ----------------------------------------------------------------
func (node *IfChainNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for _, ifItem := range node.ifItems {
		condition := types.MLRVAL_TRUE
		if ifItem.conditionNode != nil {
			condition = ifItem.conditionNode.Evaluate(state)
		}
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("mlr: conditional expression did not evaluate to boolean.")
		}
		if boolValue == true {
			blockExitPayload, err := ifItem.statementBlockNode.Execute(state)
			if err != nil {
				return nil, err
			}
			// Pass break/continue out of the if-block since they apply to the
			// containing for/while/etc.
			if blockExitPayload != nil {
				return blockExitPayload, nil
			}
			break
		}
	}
	return nil, nil
}
