// ================================================================
// This is for while/do-while loops
// ================================================================

package cst

import (
	"errors"

	"miller/src/dsl"
	"miller/src/lib"
	"miller/src/runtime"
	"miller/src/types"
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

func (this *RootNode) BuildWhileLoopNode(astNode *dsl.ASTNode) (*WhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	conditionNode, err := this.BuildEvaluableNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	statementBlockNode, err := this.BuildStatementBlockNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewWhileLoopNode(
		conditionNode,
		statementBlockNode,
	), nil
}

// ----------------------------------------------------------------
func (this *WhileLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	var condition types.Mlrval
	for {
		this.conditionNode.Evaluate(&condition, state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("Miller: conditional expression did not evaluate to boolean.")
		}
		if boolValue != true {
			break
		}
		blockExitPayload, err := this.statementBlockNode.Execute(state)
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

func (this *RootNode) BuildDoWhileLoopNode(astNode *dsl.ASTNode) (*DoWhileLoopNode, error) {
	lib.InternalCodingErrorIf(astNode.Type != dsl.NodeTypeDoWhileLoop)
	lib.InternalCodingErrorIf(len(astNode.Children) != 2)

	statementBlockNode, err := this.BuildStatementBlockNode(astNode.Children[0])
	if err != nil {
		return nil, err
	}
	conditionNode, err := this.BuildEvaluableNode(astNode.Children[1])
	if err != nil {
		return nil, err
	}

	return NewDoWhileLoopNode(
		statementBlockNode,
		conditionNode,
	), nil
}

// ----------------------------------------------------------------
func (this *DoWhileLoopNode) Execute(state *runtime.State) (*BlockExitPayload, error) {
	var condition types.Mlrval
	for {
		blockExitPayload, err := this.statementBlockNode.Execute(state)
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

		this.conditionNode.Evaluate(&condition, state)
		boolValue, isBool := condition.GetBoolValue()
		if !isBool {
			// TODO: line-number/token info for the DSL expression.
			return nil, errors.New("Miller: conditional expression did not evaluate to boolean.")
		}
		if boolValue == false {
			break
		}
	}
	return nil, nil
}
