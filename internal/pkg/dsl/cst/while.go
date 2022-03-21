// ================================================================
// This is for while/do-while loops
// ================================================================

package cst

import (
	"fmt"

	"github.com/johnkerl/miller/internal/pkg/dsl"
	"github.com/johnkerl/miller/internal/pkg/lib"
	"github.com/johnkerl/miller/internal/pkg/parsing/token"
	"github.com/johnkerl/miller/internal/pkg/runtime"
)

// ================================================================
type WhileLoopNode struct {
	conditionNode      IEvaluable
	conditionToken     *token.Token
	statementBlockNode *StatementBlockNode
}

func NewWhileLoopNode(
	conditionNode IEvaluable,
	conditionToken *token.Token,
	statementBlockNode *StatementBlockNode,
) *WhileLoopNode {
	return &WhileLoopNode{
		conditionNode:      conditionNode,
		conditionToken:     conditionToken,
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
	conditionToken := astNode.Children[0].Token
	statementBlockNode, err := root.BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewWhileLoopNode(
		conditionNode,
		conditionToken,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (node *WhileLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	for {
		condition := node.conditionNode.Evaluate(state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			return nil, fmt.Errorf(
				"mlr: conditional expression did not evaluate to boolean%s.",
				dsl.TokenToLocationInfo(node.conditionToken),
			)
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
	conditionToken     *token.Token
}

func NewDoWhileLoopNode(
	statementBlockNode *StatementBlockNode,
	conditionNode IEvaluable,
	conditionToken *token.Token,
) *DoWhileLoopNode {
	return &DoWhileLoopNode{
		statementBlockNode: statementBlockNode,
		conditionNode:      conditionNode,
		conditionToken:     conditionToken,
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
	conditionToken := astNode.Children[1].Token

	return NewDoWhileLoopNode(
		statementBlockNode,
		conditionNode,
		conditionToken,
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
			return nil, fmt.Errorf(
				"mlr: conditional expression did not evaluate to boolean%s.",
				dsl.TokenToLocationInfo(node.conditionToken),
			)
		}
		if boolValue == false {
			break
		}
	}
	return nil, nil
}
