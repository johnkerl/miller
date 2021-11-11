// ================================================================
// This is for while/do-while loops
// ================================================================

package cst

import (
	"errors"

	"mlr/internal/pkg/dsl"
	"mlr/internal/pkg/lib"
	"mlr/internal/pkg/runtime"
)

// ================================================================
type WhileLoopNode struct {
	conditionNode      IEvaluable
	statementBlockNode *StatementBlockNode
}

func NewWhileLoopNode(
	conditionNode IEvaluable,
	statementBlockNode *StatementBlockNode,
) *WhileLoopNode {
	return &WhileLoopNode{
		conditionNode:      conditionNode,
		statementBlockNode: statementBlockNode,
	}
}

func (root *RootNode) BuildWhileLoopNode(astNode *dsl.ASTNode) (*WhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	conditionNode, err := root.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	statementBlockNode, err := root.BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewWhileLoopNode(
		conditionNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (node *WhileLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for {
		condition := node.conditionNode.Evaluate(state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("mlr: conditional expression did not evaluate to boolean.")
		}
		if boolValue != true {
			break
		}
		blockExitPayload, err := node.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
				break
			}
			// If continue, keep going -- this means the body was exited
			// early but we keep going at this level
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
				return blockExitPayload, nil
			}
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
				return blockExitPayload, nil
			}
		}
		// TODO: handle return statements
		// TODO: runtime errors for any other types
	}
	return nil, nil
}

// ================================================================
type DoWhileLoopNode struct {
	statementBlockNode *StatementBlockNode
	conditionNode      IEvaluable
}

func NewDoWhileLoopNode(
	statementBlockNode *StatementBlockNode,
	conditionNode IEvaluable,
) *DoWhileLoopNode {
	return &DoWhileLoopNode{
		statementBlockNode: statementBlockNode,
		conditionNode:      conditionNode,
	}
}

func (root *RootNode) BuildDoWhileLoopNode(astNode *dsl.ASTNode) (*DoWhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDoWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	statementBlockNode, err := root.BuildStatementBlockNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	conditionNode, err := root.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewDoWhileLoopNode(
		statementBlockNode,
		conditionNode,
	), nil
}

// ----------------------------------------------------------------
func (node *DoWhileLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for {
		blockExitPayload, err := node.statementBlockNode.Execute(state)
		if err != nil {
			return nil, err
		}
		if blockExitPayload != nil {
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_BREAK {
				break
			}
			// If continue, keep going -- this means the body was exited
			// early but we keep going at this level
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VOID {
				return blockExitPayload, nil
			}
			if blockExitPayload.blockExitStatus == BLOCK_EXIT_RETURN_VALUE {
				return blockExitPayload, nil
			}
		}
		// TODO: handle return statements
		// TODO: runtime errors for any other types

		condition := node.conditionNode.Evaluate(state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("mlr: conditional expression did not evaluate to boolean.")
		}
		if boolValue == false {
			break
		}
	}
	return nil, nil
}
